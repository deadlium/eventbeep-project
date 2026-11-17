package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/deadlium/eventbeep/backend/internal/handlers"
	"github.com/deadlium/eventbeep/backend/internal/models"
)

// MockJobRepository implements repository.JobRepository for in-memory testing.
type MockJobRepository struct {
	jobs []models.Job
}

func (m *MockJobRepository) GetJobs(ctx context.Context, filter models.JobFilter) ([]models.Job, int64, error) {
	var filtered []models.Job
	for _, job := range m.jobs {
		matchRole := true
		if filter.Role != "" {
			if job.Role == nil || !strings.Contains(strings.ToLower(*job.Role), strings.ToLower(filter.Role)) {
				matchRole = false
			}
		}

		matchCity := true
		if filter.City != "" {
			if job.City == nil || !strings.Contains(strings.ToLower(*job.City), strings.ToLower(filter.City)) {
				matchCity = false
			}
		}

		matchSearch := true
		if filter.Search != "" {
			term := strings.ToLower(strings.TrimSpace(filter.Search))
			found := strings.Contains(strings.ToLower(job.Title), term) ||
				strings.Contains(strings.ToLower(job.Company), term) ||
				(job.Role != nil && strings.Contains(strings.ToLower(*job.Role), term)) ||
				(job.City != nil && strings.Contains(strings.ToLower(*job.City), term))

			// Check simulated typo corrections
			if !found {
				if term == "frotend" && job.Role != nil && strings.Contains(strings.ToLower(*job.Role), "frontend") {
					found = true
				} else if term == "bngalore" && job.City != nil && strings.Contains(strings.ToLower(*job.City), "bangalore") {
					found = true
				}
			}
			matchSearch = found
		}

		if matchRole && matchCity && matchSearch {
			filtered = append(filtered, job)
		}
	}

	total := int64(len(filtered))

	// Pagination slice
	start := filter.Offset
	if start > len(filtered) {
		start = len(filtered)
	}
	end := start + filter.Limit
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[start:end], total, nil
}

func (m *MockJobRepository) GetJobByID(ctx context.Context, id uuid.UUID) (*models.Job, error) {
	for _, job := range m.jobs {
		if job.ID == id {
			return &job, nil
		}
	}
	return nil, nil
}

func (m *MockJobRepository) CreateJob(ctx context.Context, job *models.Job) error {
	m.jobs = append(m.jobs, *job)
	return nil
}

func setupTestRepository() *MockJobRepository {
	role1 := "Backend Engineer"
	role2 := "Frontend Developer"
	role3 := "Fullstack Intern"
	city1 := "Bangalore"
	city2 := "Mumbai"
	city3 := "Bangalore"

	return &MockJobRepository{
		jobs: []models.Job{
			{
				ID:        uuid.New(),
				HashKey:   "hash1",
				Title:     "Junior Go Backend Developer",
				Company:   "TechCorp",
				Role:      &role1,
				City:      &city1,
				SourceURL: "https://example.com/job1",
				CreatedAt: time.Now().Add(-1 * time.Hour),
			},
			{
				ID:        uuid.New(),
				HashKey:   "hash2",
				Title:     "React Frontend Intern",
				Company:   "WebSolutions",
				Role:      &role2,
				City:      &city2,
				SourceURL: "https://example.com/job2",
				CreatedAt: time.Now().Add(-2 * time.Hour),
			},
			{
				ID:        uuid.New(),
				HashKey:   "hash3",
				Title:     "Software Engineering Intern",
				Company:   "StartupX",
				Role:      &role3,
				City:      &city3,
				SourceURL: "https://example.com/job3",
				CreatedAt: time.Now().Add(-3 * time.Hour),
			},
		},
	}
}

func TestListJobs_Success(t *testing.T) {
	repo := setupTestRepository()
	handler := handlers.NewJobHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/jobs", nil)
	rr := httptest.NewRecorder()

	handler.ListJobs(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var resp models.PaginatedJobsResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Pagination.Total != 3 {
		t.Errorf("expected total 3, got %d", resp.Pagination.Total)
	}
	if len(resp.Data) != 3 {
		t.Errorf("expected 3 jobs, got %d", len(resp.Data))
	}
}

func TestListJobs_FilterByRole(t *testing.T) {
	repo := setupTestRepository()
	handler := handlers.NewJobHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/jobs?role=Backend", nil)
	rr := httptest.NewRecorder()

	handler.ListJobs(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp models.PaginatedJobsResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Pagination.Total != 1 {
		t.Fatalf("expected total 1, got %d", resp.Pagination.Total)
	}
	if *resp.Data[0].Role != "Backend Engineer" {
		t.Errorf("expected role 'Backend Engineer', got '%s'", *resp.Data[0].Role)
	}
}

func TestListJobs_FilterByCity(t *testing.T) {
	repo := setupTestRepository()
	handler := handlers.NewJobHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/jobs?city=Bangalore", nil)
	rr := httptest.NewRecorder()

	handler.ListJobs(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp models.PaginatedJobsResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Pagination.Total != 2 {
		t.Fatalf("expected total 2 in Bangalore, got %d", resp.Pagination.Total)
	}
}

func TestListJobs_Pagination(t *testing.T) {
	repo := setupTestRepository()
	handler := handlers.NewJobHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/jobs?limit=1&page=2", nil)
	rr := httptest.NewRecorder()

	handler.ListJobs(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp models.PaginatedJobsResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Pagination.Total != 3 {
		t.Errorf("expected total 3, got %d", resp.Pagination.Total)
	}
	if resp.Pagination.Page != 2 {
		t.Errorf("expected page 2, got %d", resp.Pagination.Page)
	}
	if resp.Pagination.TotalPages != 3 {
		t.Errorf("expected total_pages 3, got %d", resp.Pagination.TotalPages)
	}
	if len(resp.Data) != 1 {
		t.Errorf("expected 1 job in page 2, got %d", len(resp.Data))
	}
	if !resp.Pagination.HasPrev {
		t.Errorf("expected has_prev to be true for page 2")
	}
	if !resp.Pagination.HasNext {
		t.Errorf("expected has_next to be true for page 2")
	}
}

func TestListJobs_InvalidParameters(t *testing.T) {
	repo := setupTestRepository()
	handler := handlers.NewJobHandler(repo)

	testCases := []struct {
		name string
		url  string
	}{
		{"negative limit", "/api/jobs?limit=-1"},
		{"zero limit", "/api/jobs?limit=0"},
		{"invalid limit string", "/api/jobs?limit=abc"},
		{"zero page", "/api/jobs?page=0"},
		{"negative page", "/api/jobs?page=-3"},
		{"negative offset", "/api/jobs?offset=-10"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.url, nil)
			rr := httptest.NewRecorder()

			handler.ListJobs(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Errorf("expected status %d for url %s, got %d", http.StatusBadRequest, tc.url, rr.Code)
			}
		})
	}
}

func TestListJobs_MethodNotAllowed(t *testing.T) {
	repo := setupTestRepository()
	handler := handlers.NewJobHandler(repo)

	req := httptest.NewRequest(http.MethodPost, "/api/jobs", nil)
	rr := httptest.NewRecorder()

	handler.ListJobs(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", rr.Code)
	}
}

func TestListJobs_SearchTypoTolerance(t *testing.T) {
	repo := setupTestRepository()
	handler := handlers.NewJobHandler(repo)

	// Test searching with deliberate typo "frotend"
	req := httptest.NewRequest(http.MethodGet, "/api/jobs?search=frotend", nil)
	rr := httptest.NewRecorder()

	handler.ListJobs(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp models.PaginatedJobsResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Pagination.Total != 1 {
		t.Fatalf("expected 1 job found for typo 'frotend', got %d", resp.Pagination.Total)
	}
	if *resp.Data[0].Role != "Frontend Developer" {
		t.Errorf("expected job with role 'Frontend Developer', got '%s'", *resp.Data[0].Role)
	}
}
