package scraper_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/deadlium/eventbeep/backend/internal/scraper"
)

func TestInternshalaScraper_ParseHTML(t *testing.T) {
	htmlData := loadFixture(t, "internshala_jobs.html")
	s := scraper.NewInternshalaScraper(nil)

	jobs, err := s.ParseHTML(bytes.NewReader(htmlData), "https://internshala.com/jobs/fresher-jobs/")
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if len(jobs) != 2 {
		t.Fatalf("expected 2 parsed jobs, got %d", len(jobs))
	}

	// Job 1: Zomato / Full Stack Development Intern
	j1 := jobs[0]
	if j1.Title != "Full Stack Development Intern" {
		t.Errorf("expected title 'Full Stack Development Intern', got '%s'", j1.Title)
	}
	if j1.Company != "Zomato" {
		t.Errorf("expected company 'Zomato', got '%s'", j1.Company)
	}
	if j1.City != "Gurgaon" {
		t.Errorf("expected city 'Gurgaon', got '%s'", j1.City)
	}

	// Job 2: Cred / Python Developer Fresher
	j2 := jobs[1]
	if j2.Title != "Python Developer Fresher" {
		t.Errorf("expected title 'Python Developer Fresher', got '%s'", j2.Title)
	}
	if j2.Company != "Cred" {
		t.Errorf("expected company 'Cred', got '%s'", j2.Company)
	}
	if j2.City != "Bangalore" {
		t.Errorf("expected city 'Bangalore', got '%s'", j2.City)
	}
}

func TestInternshalaScraper_ScrapeWithMock(t *testing.T) {
	htmlData := loadFixture(t, "internshala_jobs.html")

	mockClient := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(htmlData)),
			}, nil
		},
	}

	s := scraper.NewInternshalaScraper(mockClient)
	if s.Name() != "internshala" {
		t.Errorf("expected scraper name 'internshala', got '%s'", s.Name())
	}

	jobs, err := s.Scrape(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 2 {
		t.Errorf("expected 2 jobs, got %d", len(jobs))
	}
}
