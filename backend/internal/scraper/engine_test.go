package scraper_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/deadlium/eventbeep/backend/internal/models"
	"github.com/deadlium/eventbeep/backend/internal/scraper"
)

// DeduplicatingJobRepository implements repository.JobRepository and scraper.IngestionStatusChecker.
// It mimics PostgreSQL's `ON CONFLICT (hash_key) DO NOTHING`.
type DeduplicatingJobRepository struct {
	mu   sync.Mutex
	jobs map[string]models.Job // key: HashKey
}

func NewDeduplicatingJobRepository() *DeduplicatingJobRepository {
	return &DeduplicatingJobRepository{
		jobs: make(map[string]models.Job),
	}
}

func (r *DeduplicatingJobRepository) GetJobs(ctx context.Context, filter models.JobFilter) ([]models.Job, int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var result []models.Job
	for _, j := range r.jobs {
		result = append(result, j)
	}
	return result, int64(len(result)), nil
}

func (r *DeduplicatingJobRepository) GetJobByID(ctx context.Context, id uuid.UUID) (*models.Job, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, j := range r.jobs {
		if j.ID == id {
			return &j, nil
		}
	}
	return nil, nil
}

func (r *DeduplicatingJobRepository) CreateJob(ctx context.Context, job *models.Job) error {
	_, err := r.CreateJobWithStatus(ctx, job)
	return err
}

func (r *DeduplicatingJobRepository) CreateJobWithStatus(ctx context.Context, job *models.Job) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.jobs[job.HashKey]; exists {
		// Mimic ON CONFLICT (hash_key) DO NOTHING
		return false, nil
	}

	r.jobs[job.HashKey] = *job
	return true, nil
}

func (r *DeduplicatingJobRepository) Count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.jobs)
}

func TestEngine_IngestFromReader_Deduplication(t *testing.T) {
	repo := NewDeduplicatingJobRepository()
	s := scraper.NewUnstopScraper(nil)
	engine := scraper.NewEngine(repo, s)

	htmlData := loadFixture(t, "unstop_jobs.html")

	// First ingestion
	res, err := engine.IngestFromReader(context.Background(), s, bytes.NewReader(htmlData), "https://unstop.com/jobs")
	if err != nil {
		t.Fatalf("unexpected ingestion error: %v", err)
	}

	// In unstop_jobs.html there are 4 valid parsed job cards, with 1 duplicate (Razorpay / Software Engineer)
	// Expected: 4 scraped, 3 inserted, 1 skipped
	if res.TotalScraped != 4 {
		t.Errorf("expected 4 total scraped, got %d", res.TotalScraped)
	}
	if res.TotalInserted != 3 {
		t.Errorf("expected 3 total inserted, got %d", res.TotalInserted)
	}
	if res.TotalSkipped != 1 {
		t.Errorf("expected 1 skipped due to duplicate, got %d", res.TotalSkipped)
	}
	if repo.Count() != 3 {
		t.Errorf("expected 3 jobs in repository, got %d", repo.Count())
	}

	// Second ingestion with the same fixture: all jobs should be skipped!
	res2, err := engine.IngestFromReader(context.Background(), s, bytes.NewReader(htmlData), "https://unstop.com/jobs")
	if err != nil {
		t.Fatalf("unexpected second ingestion error: %v", err)
	}

	if res2.TotalInserted != 0 {
		t.Errorf("expected 0 inserted on re-ingestion, got %d", res2.TotalInserted)
	}
	if res2.TotalSkipped != 4 {
		t.Errorf("expected all 4 jobs skipped on re-ingestion, got %d", res2.TotalSkipped)
	}
	if repo.Count() != 3 {
		t.Errorf("expected repository count to remain 3, got %d", repo.Count())
	}
}

func TestEngine_MultiScraperIngest(t *testing.T) {
	repo := NewDeduplicatingJobRepository()

	unstopHTML := loadFixture(t, "unstop_jobs.html")
	wellfoundHTML := loadFixture(t, "wellfound_jobs.html")

	mockUnstopClient := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(unstopHTML)),
			}, nil
		},
	}

	mockWellfoundClient := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(wellfoundHTML)),
			}, nil
		},
	}

	unstopScraper := scraper.NewUnstopScraper(mockUnstopClient)
	wellfoundScraper := scraper.NewWellfoundScraper(mockWellfoundClient)

	engine := scraper.NewEngine(repo, unstopScraper, wellfoundScraper)

	if len(engine.Scrapers()) != 2 {
		t.Fatalf("expected 2 registered scrapers, got %d", len(engine.Scrapers()))
	}

	res, err := engine.Ingest(context.Background())
	if err != nil {
		t.Fatalf("unexpected multi-scraper ingestion error: %v", err)
	}

	// Unstop: 4 scraped (3 unique, 1 dupe)
	// Wellfound: 2 scraped (2 unique)
	// Total Scraped: 6
	// Total Inserted: 5
	// Total Skipped: 1
	if res.TotalScraped != 6 {
		t.Errorf("expected 6 total scraped, got %d", res.TotalScraped)
	}
	if res.TotalInserted != 5 {
		t.Errorf("expected 5 total inserted, got %d", res.TotalInserted)
	}
	if res.TotalSkipped != 1 {
		t.Errorf("expected 1 skipped, got %d", res.TotalSkipped)
	}
	if repo.Count() != 5 {
		t.Errorf("expected 5 jobs stored in repo, got %d", repo.Count())
	}
}

func TestEngine_ErrorResilience(t *testing.T) {
	repo := NewDeduplicatingJobRepository()

	failingClient := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("network timeout")
		},
	}

	wellfoundHTML := loadFixture(t, "wellfound_jobs.html")
	workingClient := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(wellfoundHTML)),
			}, nil
		},
	}

	failingScraper := scraper.NewUnstopScraper(failingClient)
	workingScraper := scraper.NewWellfoundScraper(workingClient)

	engine := scraper.NewEngine(repo, failingScraper, workingScraper)
	res, err := engine.Ingest(context.Background())
	if err != nil {
		t.Fatalf("engine should not abort completely on single scraper failure: %v", err)
	}

	if len(res.Errors) == 0 {
		t.Error("expected scraper error to be recorded in IngestionResult.Errors")
	}
	if res.TotalInserted != 2 {
		t.Errorf("expected working scraper's 2 jobs to be inserted, got %d", res.TotalInserted)
	}
}

func TestEngine_RegisterScraper_And_FallbackRepo(t *testing.T) {
	// Simple repo that only implements repository.JobRepository (no IngestionStatusChecker)
	simpleRepo := &simpleMockRepo{jobs: make(map[string]models.Job)}
	engine := scraper.NewEngine(simpleRepo)

	if len(engine.Scrapers()) != 0 {
		t.Errorf("expected 0 scrapers initially, got %d", len(engine.Scrapers()))
	}

	s1 := scraper.NewUnstopScraper(nil, "https://unstop.com/custom")
	s2 := scraper.NewWellfoundScraper(nil, "https://wellfound.com/custom")
	s3 := scraper.NewInternshalaScraper(nil, "https://internshala.com/custom")

	engine.RegisterScraper(s1)
	engine.RegisterScraper(s2)
	engine.RegisterScraper(s3)

	if len(engine.Scrapers()) != 3 {
		t.Errorf("expected 3 scrapers after registration, got %d", len(engine.Scrapers()))
	}
	if s1.SourceURL() != "https://unstop.com/custom" {
		t.Errorf("unexpected s1 source URL: %s", s1.SourceURL())
	}
	if s2.SourceURL() != "https://wellfound.com/custom" {
		t.Errorf("unexpected s2 source URL: %s", s2.SourceURL())
	}
	if s3.SourceURL() != "https://internshala.com/custom" {
		t.Errorf("unexpected s3 source URL: %s", s3.SourceURL())
	}

	// Test fallback repo branch in IngestJobs
	jobs := []scraper.ScrapedJob{
		{
			Title:     "DevOps Intern",
			Company:   "PhonePe",
			SourceURL: "https://example.com/job/devops",
		},
		{
			Title:     "DevOps Intern", // Duplicate in batch
			Company:   "PhonePe",
			SourceURL: "https://example.com/job/devops-2",
		},
	}

	res, err := engine.IngestJobs(context.Background(), jobs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.TotalInserted != 1 || res.TotalSkipped != 1 {
		t.Errorf("expected 1 inserted and 1 skipped, got %d inserted and %d skipped", res.TotalInserted, res.TotalSkipped)
	}

	client := scraper.DefaultHTTPClient()
	if client == nil || client.Timeout == 0 {
		t.Errorf("expected configured DefaultHTTPClient")
	}
}

type simpleMockRepo struct {
	mu   sync.Mutex
	jobs map[string]models.Job
}

func (r *simpleMockRepo) GetJobs(ctx context.Context, filter models.JobFilter) ([]models.Job, int64, error) {
	return nil, 0, nil
}

func (r *simpleMockRepo) GetJobByID(ctx context.Context, id uuid.UUID) (*models.Job, error) {
	return nil, nil
}

func (r *simpleMockRepo) CreateJob(ctx context.Context, job *models.Job) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.jobs[job.HashKey] = *job
	return nil
}

func TestScrapedJob_ToJob(t *testing.T) {
	sj := &scraper.ScrapedJob{
		Title:     "  Backend Intern  ",
		Company:   "  Razorpay  ",
		Role:      "  Engineering  ",
		City:      "  Bangalore  ",
		SourceURL: "https://example.com/job",
	}

	job := sj.ToJob()

	if job.ID == uuid.Nil {
		t.Error("expected non-nil UUID")
	}
	if job.Title != "Backend Intern" {
		t.Errorf("expected trimmed title, got '%s'", job.Title)
	}
	if job.Company != "Razorpay" {
		t.Errorf("expected trimmed company, got '%s'", job.Company)
	}
	if job.Role == nil || *job.Role != "Engineering" {
		t.Errorf("expected trimmed role 'Engineering', got '%v'", job.Role)
	}
	if job.City == nil || *job.City != "Bangalore" {
		t.Errorf("expected trimmed city 'Bangalore', got '%v'", job.City)
	}
	if job.HashKey == "" {
		t.Error("expected populated HashKey")
	}
}
