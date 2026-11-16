package scheduler

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/deadlium/eventbeep/backend/internal/llm"
	"github.com/deadlium/eventbeep/backend/internal/repository"
	"github.com/deadlium/eventbeep/backend/internal/scraper"
)

const DefaultInterval = 24 * time.Hour

// Scheduler manages scheduled background runs of the scraping, deduplication, and LLM enrichment pipeline.
type Scheduler struct {
	repo       repository.JobRepository
	engine     *scraper.Engine
	extractor  llm.Extractor
	interval   time.Duration
	stopCh     chan struct{}
	doneCh     chan struct{}
	running    bool
	mu         sync.RWMutex
	lastRun    time.Time
	lastResult *scraper.IngestionResult
}

// NewScheduler constructs a background scheduler with configured interval and dependencies.
func NewScheduler(repo repository.JobRepository, engine *scraper.Engine, extractor llm.Extractor, interval time.Duration) *Scheduler {
	if interval <= 0 {
		interval = DefaultInterval
	}

	if engine != nil && extractor != nil && engine.Extractor() == nil {
		engine.SetExtractor(extractor)
	}

	return &Scheduler{
		repo:      repo,
		engine:    engine,
		extractor: extractor,
		interval:  interval,
		stopCh:    make(chan struct{}),
		doneCh:    make(chan struct{}),
	}
}

// Interval returns the configured ticker interval.
func (s *Scheduler) Interval() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.interval
}

// IsRunning reports whether the background scheduler ticker is active.
func (s *Scheduler) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// LastRun returns the timestamp of the most recent pipeline execution.
func (s *Scheduler) LastRun() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastRun
}

// LastResult returns the result of the most recent ingestion run.
func (s *Scheduler) LastResult() *scraper.IngestionResult {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastResult
}

// Start launches the background scheduler goroutine.
func (s *Scheduler) Start(ctx context.Context) {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.stopCh = make(chan struct{})
	s.doneCh = make(chan struct{})
	s.mu.Unlock()

	go s.loop(ctx)
}

// loop runs the ticker loop until context cancellation or Stop() call.
func (s *Scheduler) loop(ctx context.Context) {
	defer close(s.doneCh)
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	log.Printf("[scheduler] Started 24h background pipeline worker (interval: %v)", s.interval)

	// Execute initial ingestion run on startup
	if _, err := s.RunOnce(ctx); err != nil {
		log.Printf("[scheduler] Initial ingestion run encountered error: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("[scheduler] Context cancelled; shutting down scheduler...")
			s.mu.Lock()
			s.running = false
			s.mu.Unlock()
			return
		case <-s.stopCh:
			log.Println("[scheduler] Stop signal received; shutting down scheduler...")
			s.mu.Lock()
			s.running = false
			s.mu.Unlock()
			return
		case <-ticker.C:
			log.Println("[scheduler] Ticker triggered scheduled pipeline run...")
			if _, err := s.RunOnce(ctx); err != nil {
				log.Printf("[scheduler] Scheduled ingestion run failed: %v", err)
			}
		}
	}
}

// Stop gracefully stops the background scheduler and waits for it to exit.
func (s *Scheduler) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	close(s.stopCh)
	s.running = false
	s.mu.Unlock()

	<-s.doneCh
	log.Println("[scheduler] Scheduler stopped gracefully.")
}

// RunOnce executes a single iteration of: Scrape -> Deduplicate -> LLM Enrich -> Store in DB.
func (s *Scheduler) RunOnce(ctx context.Context) (*scraper.IngestionResult, error) {
	startTime := time.Now()
	log.Printf("[scheduler] Starting pipeline cycle: Scrape -> Deduplicate -> LLM Enrich -> Store...")

	if s.engine == nil {
		return nil, fmt.Errorf("scraper engine is not configured")
	}

	res, err := s.engine.Ingest(ctx)
	if err != nil {
		return nil, fmt.Errorf("ingestion cycle failed: %w", err)
	}

	s.mu.Lock()
	s.lastRun = startTime
	s.lastResult = res
	s.mu.Unlock()

	log.Printf("[scheduler] Pipeline cycle completed in %v: Scraped=%d, Inserted=%d, Skipped=%d, Errors=%d",
		time.Since(startTime), res.TotalScraped, res.TotalInserted, res.TotalSkipped, len(res.Errors))

	return res, nil
}
