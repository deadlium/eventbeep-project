package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/deadlium/eventbeep/backend/internal/database"
)

// HealthHandler handles health probe endpoints.
type HealthHandler struct {
	db *database.DB
}

// NewHealthHandler creates a new HealthHandler.
func NewHealthHandler(db *database.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// HealthCheck handles GET /health and GET /api/health.
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	status := "ok"
	dbStatus := "connected"

	if h.db != nil {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := h.db.Ping(ctx); err != nil {
			status = "degraded"
			dbStatus = "unreachable"
		}
	} else {
		dbStatus = "unconfigured"
	}

	statusCode := http.StatusOK
	if status == "degraded" {
		statusCode = http.StatusServiceUnavailable
	}

	writeJSON(w, statusCode, map[string]interface{}{
		"status":    status,
		"database":  dbStatus,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"service":   "eventbeep-backend",
	})
}
