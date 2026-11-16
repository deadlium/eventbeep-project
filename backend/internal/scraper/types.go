package scraper

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/deadlium/eventbeep/backend/internal/models"
)

// ScrapedJob represents the normalized data extracted from a job board HTML page.
type ScrapedJob struct {
	Title          string `json:"title"`
	Company        string `json:"company"`
	RawDescription string `json:"raw_description,omitempty"`
	Role           string `json:"role,omitempty"`
	City           string `json:"city,omitempty"`
	SourceURL      string `json:"source_url"`
	HashKey        string `json:"hash_key,omitempty"`
}

// ToJob converts a ScrapedJob into a models.Job with UUID and deterministic HashKey.
func (s *ScrapedJob) ToJob() *models.Job {
	hash := s.HashKey
	if hash == "" {
		hash = ComputeHashKey(s.Company, s.Title)
	}

	var rolePtr *string
	trimmedRole := strings.TrimSpace(s.Role)
	if trimmedRole != "" {
		rolePtr = &trimmedRole
	}

	var cityPtr *string
	trimmedCity := strings.TrimSpace(s.City)
	if trimmedCity != "" {
		cityPtr = &trimmedCity
	}

	return &models.Job{
		ID:        uuid.New(),
		HashKey:   hash,
		Title:     strings.TrimSpace(s.Title),
		Company:   strings.TrimSpace(s.Company),
		Role:      rolePtr,
		City:      cityPtr,
		SourceURL: strings.TrimSpace(s.SourceURL),
		CreatedAt: time.Now(),
	}
}

// Scraper defines the interface that all job board scrapers must satisfy.
type Scraper interface {
	// Name returns a friendly identifier for the scraper (e.g. "unstop", "wellfound").
	Name() string
	// SourceURL returns the target base URL.
	SourceURL() string
	// Scrape executes live HTTP fetch and parsing with context.
	Scrape(ctx context.Context) ([]ScrapedJob, error)
	// ParseHTML parses HTML from an io.Reader for deterministic testing.
	ParseHTML(r io.Reader, baseURL string) ([]ScrapedJob, error)
}
