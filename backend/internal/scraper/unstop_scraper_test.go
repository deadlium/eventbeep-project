package scraper_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/deadlium/eventbeep/backend/internal/scraper"
)

type mockHTTPClient struct {
	doFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.doFunc(req)
}

func loadFixture(t *testing.T, filename string) []byte {
	t.Helper()
	path := filepath.Join("testdata", filename)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read test fixture %s: %v", path, err)
	}
	return data
}

func TestUnstopScraper_ParseHTML(t *testing.T) {
	htmlData := loadFixture(t, "unstop_jobs.html")
	s := scraper.NewUnstopScraper(nil)

	jobs, err := s.ParseHTML(bytes.NewReader(htmlData), "https://unstop.com/jobs")
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	// In unstop_jobs.html:
	// Job 1: Razorpay / Software Engineer - Fresher
	// Job 2: Swiggy / Frontend Developer Intern
	// Job 3: Flipkart / Data Analyst Trainee
	// Job 4: Malformed (ignored)
	// Job 5: Razorpay / Software Engineer - Fresher (duplicate content)
	// ParseHTML extracts all valid job cards (4 valid cards total)
	if len(jobs) != 4 {
		t.Fatalf("expected 4 parsed jobs, got %d", len(jobs))
	}

	// Verify Job 1
	j1 := jobs[0]
	if j1.Title != "Software Engineer - Fresher" {
		t.Errorf("expected title 'Software Engineer - Fresher', got '%s'", j1.Title)
	}
	if j1.Company != "Razorpay" {
		t.Errorf("expected company 'Razorpay', got '%s'", j1.Company)
	}
	if j1.City != "Bangalore" {
		t.Errorf("expected city 'Bangalore', got '%s'", j1.City)
	}
	if j1.Role != "Backend Engineer" {
		t.Errorf("expected role 'Backend Engineer', got '%s'", j1.Role)
	}
	expectedURL := "https://unstop.com/jobs/software-engineer-fresher-1"
	if j1.SourceURL != expectedURL {
		t.Errorf("expected source_url '%s', got '%s'", expectedURL, j1.SourceURL)
	}
	if j1.HashKey != scraper.ComputeHashKey("Razorpay", "Software Engineer - Fresher") {
		t.Errorf("hash mismatch for job 1: %s", j1.HashKey)
	}

	// Verify Job 3 absolute URL preservation
	j3 := jobs[2]
	if j3.SourceURL != "https://unstop.com/jobs/data-analyst-trainee-3" {
		t.Errorf("expected absolute URL to be preserved, got '%s'", j3.SourceURL)
	}

	// Verify Job 5 (duplicate card) produces exact same HashKey as Job 1
	j4 := jobs[3]
	if j4.HashKey != j1.HashKey {
		t.Errorf("expected duplicate job cards to have identical HashKey: %s != %s", j4.HashKey, j1.HashKey)
	}
}

func TestUnstopScraper_ScrapeWithMockClient(t *testing.T) {
	htmlData := loadFixture(t, "unstop_jobs.html")

	mockClient := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("User-Agent") == "" {
				t.Error("expected User-Agent header to be set")
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(htmlData)),
			}, nil
		},
	}

	s := scraper.NewUnstopScraper(mockClient, "https://unstop.com/jobs/mock")
	if s.Name() != "unstop" {
		t.Errorf("expected name 'unstop', got '%s'", s.Name())
	}
	if s.SourceURL() != "https://unstop.com/jobs/mock" {
		t.Errorf("expected custom source url, got '%s'", s.SourceURL())
	}

	jobs, err := s.Scrape(context.Background())
	if err != nil {
		t.Fatalf("unexpected scrape error: %v", err)
	}
	if len(jobs) != 4 {
		t.Errorf("expected 4 jobs, got %d", len(jobs))
	}
}

func TestUnstopScraper_ScrapeHTTPError(t *testing.T) {
	mockClient := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(bytes.NewReader([]byte("Server Error"))),
			}, nil
		},
	}

	s := scraper.NewUnstopScraper(mockClient)
	_, err := s.Scrape(context.Background())
	if err == nil {
		t.Fatal("expected error on HTTP 500, got nil")
	}
}
