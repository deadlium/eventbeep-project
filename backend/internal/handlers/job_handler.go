package handlers

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/deadlium/eventbeep/backend/internal/models"
	"github.com/deadlium/eventbeep/backend/internal/repository"
)

// JobHandler handles job-related HTTP requests.
type JobHandler struct {
	repo repository.JobRepository
}

// NewJobHandler creates a new JobHandler.
func NewJobHandler(repo repository.JobRepository) *JobHandler {
	return &JobHandler{repo: repo}
}

// ListJobs handles GET /api/jobs.
// Supports query parameters:
//   - role: filter by role (case-insensitive substring)
//   - city: filter by city (case-insensitive substring)
//   - limit: page size (default 20, max 100)
//   - page: 1-indexed page number (default 1)
//   - offset: optional manual byte/record offset
func (h *JobHandler) ListJobs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed", "")
		return
	}

	q := r.URL.Query()

	// Parse limit
	limit := 20
	if limitStr := q.Get("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 {
			writeJSONError(w, http.StatusBadRequest, "invalid limit: must be a positive integer", "")
			return
		}
		if parsedLimit > 100 {
			parsedLimit = 100 // Enforce max page size to protect DB
		}
		limit = parsedLimit
	}

	// Parse page
	page := 1
	if pageStr := q.Get("page"); pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err != nil || parsedPage <= 0 {
			writeJSONError(w, http.StatusBadRequest, "invalid page: must be a positive integer starting at 1", "")
			return
		}
		page = parsedPage
	}

	// Calculate offset
	offset := (page - 1) * limit
	if offsetStr := q.Get("offset"); offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil || parsedOffset < 0 {
			writeJSONError(w, http.StatusBadRequest, "invalid offset: must be a non-negative integer", "")
			return
		}
		offset = parsedOffset
		page = (offset / limit) + 1
	}

	role := strings.TrimSpace(q.Get("role"))
	city := strings.TrimSpace(q.Get("city"))

	filter := models.JobFilter{
		Role:   role,
		City:   city,
		Limit:  limit,
		Offset: offset,
		Page:   page,
	}

	jobs, total, err := h.repo.GetJobs(r.Context(), filter)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to retrieve jobs", err.Error())
		return
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}

	hasNext := page < totalPages
	hasPrev := page > 1 && totalPages > 0

	resp := models.PaginatedJobsResponse{
		Data: jobs,
		Pagination: models.Pagination{
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
			HasNext:    hasNext,
			HasPrev:    hasPrev,
		},
	}

	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
	}
}

func writeJSONError(w http.ResponseWriter, status int, message, details string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(models.ErrorResponse{
		Error:   message,
		Details: details,
	})
}
