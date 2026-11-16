package scraper_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/deadlium/eventbeep/backend/internal/scraper"
)

func TestWellfoundScraper_ParseHTML(t *testing.T) {
	htmlData := loadFixture(t, "wellfound_jobs.html")
	s := scraper.NewWellfoundScraper(nil)

	jobs, err := s.ParseHTML(bytes.NewReader(htmlData), "https://wellfound.com/jobs")
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if len(jobs) != 2 {
		t.Fatalf("expected 2 parsed jobs, got %d", len(jobs))
	}

	// Job 1: Postman / Junior Backend Developer
	j1 := jobs[0]
	if j1.Title != "Junior Backend Developer" {
		t.Errorf("expected title 'Junior Backend Developer', got '%s'", j1.Title)
	}
	if j1.Company != "Postman" {
		t.Errorf("expected company 'Postman', got '%s'", j1.Company)
	}
	if j1.City != "Bangalore" {
		t.Errorf("expected city 'Bangalore', got '%s'", j1.City)
	}
	if j1.Role != "Software Engineer" {
		t.Errorf("expected role 'Software Engineer', got '%s'", j1.Role)
	}
	if j1.SourceURL != "https://wellfound.com/jobs/101-junior-backend-developer" {
		t.Errorf("unexpected URL: %s", j1.SourceURL)
	}

	// Job 2: Zepto / Growth Associate - Fresher
	j2 := jobs[1]
	if j2.Title != "Growth Associate - Fresher" {
		t.Errorf("expected title 'Growth Associate - Fresher', got '%s'", j2.Title)
	}
	if j2.Company != "Zepto" {
		t.Errorf("expected company 'Zepto', got '%s'", j2.Company)
	}
	if j2.City != "Mumbai" {
		t.Errorf("expected city 'Mumbai', got '%s'", j2.City)
	}
	if j2.Role != "Growth" {
		t.Errorf("expected role 'Growth', got '%s'", j2.Role)
	}
}

func TestWellfoundScraper_ScrapeWithMock(t *testing.T) {
	htmlData := loadFixture(t, "wellfound_jobs.html")

	mockClient := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(htmlData)),
			}, nil
		},
	}

	s := scraper.NewWellfoundScraper(mockClient)
	if s.Name() != "wellfound" {
		t.Errorf("expected scraper name 'wellfound', got '%s'", s.Name())
	}

	jobs, err := s.Scrape(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 2 {
		t.Errorf("expected 2 jobs, got %d", len(jobs))
	}
}
