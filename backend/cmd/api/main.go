package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/deadlium/eventbeep/backend/internal/config"
	"github.com/deadlium/eventbeep/backend/internal/database"
	"github.com/deadlium/eventbeep/backend/internal/llm"
	"github.com/deadlium/eventbeep/backend/internal/repository"
	"github.com/deadlium/eventbeep/backend/internal/routes"
	"github.com/deadlium/eventbeep/backend/internal/scheduler"
	"github.com/deadlium/eventbeep/backend/internal/scraper"
)

func main() {
	log.Println("[server] Starting EventBeep backend service...")

	cfg := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize database connection pool
	var db *database.DB
	var jobRepo repository.JobRepository

	pool, err := database.NewPool(ctx, cfg)
	if err != nil {
		log.Printf("[database] Warning: Could not connect to PostgreSQL: %v", err)
		log.Println("[database] Starting in degraded mode (database unavailable until reconnected).")
	} else {
		db = pool
		defer db.Close()

		// Run automated schema validation/migration
		if err := db.EnsureSchema(ctx); err != nil {
			log.Printf("[database] Schema migration warning: %v", err)
		}

		jobRepo = repository.NewPostgresJobRepository(db.Pool)

		// Initialize Phase 2 scrapers, Phase 3 Gemini LLM entity extractor, and 24h background scheduler
		unstopScraper := scraper.NewUnstopScraper(nil)
		wellfoundScraper := scraper.NewWellfoundScraper(nil)
		internshalaScraper := scraper.NewInternshalaScraper(nil)

		extractor := llm.NewGeminiExtractor(cfg.GeminiAPIKey, llm.WithModel(cfg.GeminiModel))
		scraperEngine := scraper.NewEngine(jobRepo, unstopScraper, wellfoundScraper, internshalaScraper)
		scraperEngine.SetExtractor(extractor)

		if cfg.SchedulerEnabled {
			bgScheduler := scheduler.NewScheduler(jobRepo, scraperEngine, extractor, cfg.SchedulerInterval)
			bgScheduler.Start(ctx)
			defer bgScheduler.Stop()
		}
	}

	router := routes.SetupRouter(db, jobRepo)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Listen for OS signals for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("[server] HTTP listening on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[server] Fatal server error: %v", err)
		}
	}()

	<-stop
	log.Println("[server] Shutting down gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("[server] Server shutdown failed: %v", err)
	}

	log.Println("[server] Service exited cleanly.")
}
