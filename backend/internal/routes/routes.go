package routes

import (
	"log"
	"net/http"
	"time"

	"github.com/deadlium/eventbeep/backend/internal/database"
	"github.com/deadlium/eventbeep/backend/internal/handlers"
	"github.com/deadlium/eventbeep/backend/internal/repository"
)

// SetupRouter configures all HTTP routes, handlers, and middlewares.
func SetupRouter(db *database.DB, jobRepo repository.JobRepository) http.Handler {
	mux := http.NewServeMux()

	healthHandler := handlers.NewHealthHandler(db)
	jobHandler := handlers.NewJobHandler(jobRepo)

	// Health check endpoints
	mux.HandleFunc("GET /health", healthHandler.HealthCheck)
	mux.HandleFunc("GET /api/health", healthHandler.HealthCheck)

	// Job endpoints
	mux.HandleFunc("GET /api/jobs", jobHandler.ListJobs)

	// Wrap mux with middleware chain
	return corsMiddleware(loggingMiddleware(recoveryMiddleware(mux)))
}

// corsMiddleware attaches permissive CORS headers for Next.js frontend consumption.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware records incoming request details and execution duration.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[http] %s %s took %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// recoveryMiddleware catches and handles runtime panics gracefully.
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("[panic] recovered from panic: %v", rec)
				http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
