package scraper

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/deadlium/eventbeep/backend/internal/llm"
	"github.com/deadlium/eventbeep/backend/internal/models"
	"github.com/deadlium/eventbeep/backend/internal/repository"
)

// IngestionResult summarizes the results of a scrape and ingestion run.
type IngestionResult struct {
	TotalScraped  int          `json:"total_scraped"`
	TotalInserted int          `json:"total_inserted"`
	TotalSkipped  int          `json:"total_skipped"`
	Jobs          []models.Job `json:"jobs,omitempty"`
	Errors        []string     `json:"errors,omitempty"`
}

// IngestionStatusChecker is an optional interface repository implementations can satisfy
// to report whether ON CONFLICT actually skipped insertion.
type IngestionStatusChecker interface {
	CreateJobWithStatus(ctx context.Context, job *models.Job) (bool, error)
}

// Engine coordinates job board scrapers and handles deduplicated persistence.
type Engine struct {
	repo      repository.JobRepository
	scrapers  []Scraper
	extractor llm.Extractor
	mu        sync.RWMutex
}

// NewEngine initializes a new scraper engine with repository and scrapers.
func NewEngine(repo repository.JobRepository, scrapers ...Scraper) *Engine {
	return &Engine{
		repo:     repo,
		scrapers: scrapers,
	}
}

// SetExtractor configures an LLM entity extractor to enrich scraped jobs before insertion.
func (e *Engine) SetExtractor(extractor llm.Extractor) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.extractor = extractor
}

// Extractor returns the configured LLM entity extractor.
func (e *Engine) Extractor() llm.Extractor {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.extractor
}

// RegisterScraper registers an additional scraper with the engine.
func (e *Engine) RegisterScraper(s Scraper) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.scrapers = append(e.scrapers, s)
}

// Scrapers returns a copy of registered scrapers.
func (e *Engine) Scrapers() []Scraper {
	e.mu.RLock()
	defer e.mu.RUnlock()
	list := make([]Scraper, len(e.scrapers))
	copy(list, e.scrapers)
	return list
}

// Ingest executes all registered scrapers and persists jobs idempotently.
func (e *Engine) Ingest(ctx context.Context) (*IngestionResult, error) {
	combinedResult := &IngestionResult{}

	e.mu.RLock()
	scrapers := make([]Scraper, len(e.scrapers))
	copy(scrapers, e.scrapers)
	e.mu.RUnlock()

	for _, s := range scrapers {
		res, err := e.IngestScraper(ctx, s)
		if err != nil {
			combinedResult.Errors = append(combinedResult.Errors, fmt.Sprintf("scraper %s error: %v", s.Name(), err))
			continue
		}
		combinedResult.TotalScraped += res.TotalScraped
		combinedResult.TotalInserted += res.TotalInserted
		combinedResult.TotalSkipped += res.TotalSkipped
		combinedResult.Jobs = append(combinedResult.Jobs, res.Jobs...)
		combinedResult.Errors = append(combinedResult.Errors, res.Errors...)
	}

	return combinedResult, nil
}

// IngestScraper executes a single scraper and stores scraped jobs.
func (e *Engine) IngestScraper(ctx context.Context, s Scraper) (*IngestionResult, error) {
	scrapedJobs, err := s.Scrape(ctx)
	if err != nil {
		return nil, fmt.Errorf("scraping failed for %s: %w", s.Name(), err)
	}

	return e.IngestJobs(ctx, scrapedJobs)
}

// IngestFromReader parses HTML from an io.Reader using a given scraper and ingests the jobs.
// Useful for deterministic fixture ingestion and testing.
func (e *Engine) IngestFromReader(ctx context.Context, s Scraper, r io.Reader, baseURL string) (*IngestionResult, error) {
	scrapedJobs, err := s.ParseHTML(r, baseURL)
	if err != nil {
		return nil, fmt.Errorf("parsing failed for %s: %w", s.Name(), err)
	}

	return e.IngestJobs(ctx, scrapedJobs)
}

// IngestJobs persists scraped jobs to the repository, respecting SHA-256 deduplication.
func (e *Engine) IngestJobs(ctx context.Context, scrapedJobs []ScrapedJob) (*IngestionResult, error) {
	res := &IngestionResult{
		TotalScraped: len(scrapedJobs),
	}

	seenHashes := make(map[string]bool)

	statusChecker, hasStatusChecker := e.repo.(IngestionStatusChecker)

	for _, scraped := range scrapedJobs {
		job := scraped.ToJob()

		// In-batch deduplication check
		if seenHashes[job.HashKey] {
			res.TotalSkipped++
			continue
		}
		seenHashes[job.HashKey] = true

		// Enrich job with LLM extractor if configured
		e.mu.RLock()
		ext := e.extractor
		e.mu.RUnlock()

		if ext != nil {
			entity, err := ext.Extract(ctx, job.Title, scraped.RawDescription)
			if err == nil && entity != nil {
				if entity.Role != "" {
					r := entity.Role
					job.Role = &r
				}
				if entity.City != "" {
					c := entity.City
					job.City = &c
				}
			}
		}

		if hasStatusChecker {
			inserted, err := statusChecker.CreateJobWithStatus(ctx, job)
			if err != nil {
				res.Errors = append(res.Errors, fmt.Sprintf("failed to save job %s: %v", job.Title, err))
				continue
			}
			if inserted {
				res.TotalInserted++
				res.Jobs = append(res.Jobs, *job)
			} else {
				res.TotalSkipped++
			}
		} else {
			if err := e.repo.CreateJob(ctx, job); err != nil {
				res.Errors = append(res.Errors, fmt.Sprintf("failed to save job %s: %v", job.Title, err))
				continue
			}
			res.TotalInserted++
			res.Jobs = append(res.Jobs, *job)
		}
	}

	return res, nil
}
