package models

import (
	"time"

	"github.com/google/uuid"
)

// Job represents a job listing in the database.
type Job struct {
	ID        uuid.UUID  `json:"id"`
	HashKey   string     `json:"hash_key"`
	Title     string     `json:"title"`
	Company   string     `json:"company"`
	Role      *string    `json:"role,omitempty"`
	City      *string    `json:"city,omitempty"`
	SourceURL string     `json:"source_url"`
	CreatedAt time.Time  `json:"created_at"`
}

// JobFilter defines query parameters for listing jobs.
type JobFilter struct {
	Role   string `json:"role,omitempty"`
	City   string `json:"city,omitempty"`
	Search string `json:"search,omitempty"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Page   int    `json:"page"`
}

// Pagination metadata for paginated responses.
type Pagination struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// PaginatedJobsResponse represents the JSON envelope for GET /api/jobs.
type PaginatedJobsResponse struct {
	Data       []Job      `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// ErrorResponse represents a standard JSON error response.
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}
