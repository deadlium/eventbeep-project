package scheduler_test

import (
	"context"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/deadlium/eventbeep/backend/internal/llm"
	"github.com/deadlium/eventbeep/backend/internal/models"
	"github.com/deadlium/eventbeep/backend/internal/scheduler"
	"github.com/deadlium/eventbeep/backend/internal/scraper"
)

// MockJobRepository implements repository.JobRepository for scheduler tests.
type MockJobRepository struct {
	mu   sync.Mutex
	jobs map[string]models.Job
}

func NewMockJobRepository() *MockJobRepository {
	return &MockJobRepository{
		jobs: make(map[string]models.Job),
	}
}

func (m *MockJobRepository) GetJobs(ctx context.Context, filter models.JobFilter) ([]models.Job, int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var list []models.Job
	for _, j := range m.jobs {
		list = append(list, j)
	}
	return list, int64(len(list)), nil
}

func (m *MockJobRepository) GetJobByID(ctx context.Context, id uuid.UUID) (*models.Job, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, j := range m.jobs {
		if j.ID == id {
			return &j, nil
		}
	}
	return nil, nil
}

func (m *MockJobRepository) CreateJob(ctx context.Context, job *models.Job) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.jobs[job.HashKey]; exists {
		return nil
	}
	m.jobs[job.HashKey] = *job
	return nil
}

func (m *MockJobRepository) Count() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.jobs)
}

// MockScraper implements scraper.Scraper for testing.
type MockScraper struct {
	name string
	jobs []scraper.ScrapedJob
}

func (m *MockScraper) Name() string {
	return m.name
}

func (m *MockScraper) SourceURL() string {
	return "https://example.com/mock"
}

func (m *MockScraper) Scrape(ctx context.Context) ([]scraper.ScrapedJob, error) {
	return m.jobs, nil
}

func (m *MockScraper) ParseHTML(r io.Reader, baseURL string) ([]scraper.ScrapedJob, error) {
	return m.jobs, nil
}

func TestNewScheduler_DefaultsAndCustom(t *testing.T) {
	repo := NewMockJobRepository()
	engine := scraper.NewEngine(repo)
	extractor := llm.NewFallbackExtractor()

	// Default interval
	sDefault := scheduler.NewScheduler(repo, engine, extractor, 0)
	if sDefault.Interval() != 24*time.Hour {
		t.Errorf("expected default interval 24h, got %v", sDefault.Interval())
	}
	if engine.Extractor() != extractor {
		t.Error("expected extractor to be wired to engine")
	}

	// Custom interval
	sCustom := scheduler.NewScheduler(repo, engine, extractor, 12*time.Hour)
	if sCustom.Interval() != 12*time.Hour {
		t.Errorf("expected custom interval 12h, got %v", sCustom.Interval())
	}
}

func TestScheduler_RunOnce_PipelineOrchestration(t *testing.T) {
	repo := NewMockJobRepository()

	mockScraper := &MockScraper{
		name: "test-board",
		jobs: []scraper.ScrapedJob{
			{
				Title:          "Golang Backend Developer",
				Company:        "Razorpay",
				RawDescription: "Build distributed payment systems in Bangalore.",
				SourceURL:      "https://example.com/job1",
			},
			{
				Title:          "React Frontend Intern",
				Company:        "Swiggy",
				RawDescription: "Work on mobile-first web in Mumbai.",
				SourceURL:      "https://example.com/job2",
			},
			{
				Title:          "Golang Backend Developer", // Duplicate
				Company:        "Razorpay",
				RawDescription: "Duplicate listing",
				SourceURL:      "https://example.com/job1-repost",
			},
		},
	}

	engine := scraper.NewEngine(repo, mockScraper)
	extractor := llm.NewFallbackExtractor()
	s := scheduler.NewScheduler(repo, engine, extractor, 24*time.Hour)

	res, err := s.RunOnce(context.Background())
	if err != nil {
		t.Fatalf("unexpected RunOnce error: %v", err)
	}

	// 3 scraped, 2 unique inserted, 1 duplicate skipped
	if res.TotalScraped != 3 {
		t.Errorf("expected 3 scraped, got %d", res.TotalScraped)
	}
	if res.TotalInserted != 2 {
		t.Errorf("expected 2 inserted, got %d", res.TotalInserted)
	}
	if res.TotalSkipped != 1 {
		t.Errorf("expected 1 skipped, got %d", res.TotalSkipped)
	}
	if repo.Count() != 2 {
		t.Errorf("expected 2 jobs stored in repo, got %d", repo.Count())
	}

	if s.LastRun().IsZero() {
		t.Error("expected LastRun to be updated")
	}
	if s.LastResult() != res {
		t.Error("expected LastResult to match res")
	}

	// Verify LLM enrichment occurred on the stored jobs
	jobs, _, err := repo.GetJobs(context.Background(), models.JobFilter{})
	if err != nil {
		t.Fatalf("failed to get jobs: %v", err)
	}

	foundBackend := false
	foundFrontend := false
	for _, j := range jobs {
		if j.Company == "Razorpay" {
			foundBackend = true
			if j.Role == nil || *j.Role != "Backend Engineer" {
				t.Errorf("expected Razorpay role 'Backend Engineer', got %v", j.Role)
			}
			if j.City == nil || *j.City != "Bangalore" {
				t.Errorf("expected Razorpay city 'Bangalore', got %v", j.City)
			}
		}
		if j.Company == "Swiggy" {
			foundFrontend = true
			if j.Role == nil || *j.Role != "Frontend Engineer" {
				t.Errorf("expected Swiggy role 'Frontend Engineer', got %v", j.Role)
			}
			if j.City == nil || *j.City != "Mumbai" {
				t.Errorf("expected Swiggy city 'Mumbai', got %v", j.City)
			}
		}
	}

	if !foundBackend || !foundFrontend {
		t.Error("expected both enriched jobs to be found in repository")
	}
}

func TestScheduler_StartAndStop(t *testing.T) {
	repo := NewMockJobRepository()
	engine := scraper.NewEngine(repo)
	extractor := llm.NewFallbackExtractor()

	s := scheduler.NewScheduler(repo, engine, extractor, 50*time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.Start(ctx)
	if !s.IsRunning() {
		t.Error("expected scheduler to be running")
	}

	// Double start should be a no-op
	s.Start(ctx)

	// Wait briefly for ticker run
	time.Sleep(120 * time.Millisecond)

	s.Stop()
	if s.IsRunning() {
		t.Error("expected scheduler to stop")
	}

	// Stop when not running should be a no-op
	s.Stop()
}

func TestScheduler_ContextCancellation(t *testing.T) {
	repo := NewMockJobRepository()
	engine := scraper.NewEngine(repo)
	extractor := llm.NewFallbackExtractor()

	s := scheduler.NewScheduler(repo, engine, extractor, 50*time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	s.Start(ctx)

	if !s.IsRunning() {
		t.Error("expected scheduler to be running")
	}

	// Cancel context
	cancel()

	// Wait briefly for goroutine to exit
	time.Sleep(60 * time.Millisecond)

	if s.IsRunning() {
		t.Error("expected scheduler to stop after context cancellation")
	}
}
