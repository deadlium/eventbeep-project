package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/deadlium/eventbeep/backend/internal/config"
)

func TestConfig_Defaults(t *testing.T) {
	os.Unsetenv("PORT")
	os.Unsetenv("DATABASE_URL")

	cfg := config.Load()

	if cfg.Port != "8080" {
		t.Errorf("expected default port 8080, got %s", cfg.Port)
	}
	if cfg.DBMaxConns != 25 {
		t.Errorf("expected default DBMaxConns 25, got %d", cfg.DBMaxConns)
	}
	if cfg.DBMinConns != 5 {
		t.Errorf("expected default DBMinConns 5, got %d", cfg.DBMinConns)
	}
	if cfg.GeminiModel != "gemini-1.5-flash" {
		t.Errorf("expected default GeminiModel gemini-1.5-flash, got %s", cfg.GeminiModel)
	}
	if cfg.SchedulerInterval != 24*time.Hour {
		t.Errorf("expected default SchedulerInterval 24h, got %v", cfg.SchedulerInterval)
	}
	if !cfg.SchedulerEnabled {
		t.Errorf("expected default SchedulerEnabled true")
	}
}

func TestConfig_CustomEnv(t *testing.T) {
	os.Setenv("PORT", "9090")
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/testdb")
	os.Setenv("GEMINI_API_KEY", "test-api-key")
	os.Setenv("GEMINI_MODEL", "gemini-2.0-flash")
	os.Setenv("SCHEDULER_INTERVAL", "12h")
	os.Setenv("SCHEDULER_ENABLED", "false")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("GEMINI_API_KEY")
		os.Unsetenv("GEMINI_MODEL")
		os.Unsetenv("SCHEDULER_INTERVAL")
		os.Unsetenv("SCHEDULER_ENABLED")
	}()

	cfg := config.Load()

	if cfg.Port != "9090" {
		t.Errorf("expected port 9090, got %s", cfg.Port)
	}
	if cfg.DatabaseURL != "postgres://user:pass@localhost:5432/testdb" {
		t.Errorf("expected custom database url, got %s", cfg.DatabaseURL)
	}
	if cfg.GeminiAPIKey != "test-api-key" {
		t.Errorf("expected custom gemini api key, got %s", cfg.GeminiAPIKey)
	}
	if cfg.GeminiModel != "gemini-2.0-flash" {
		t.Errorf("expected custom gemini model, got %s", cfg.GeminiModel)
	}
	if cfg.SchedulerInterval != 12*time.Hour {
		t.Errorf("expected custom scheduler interval 12h, got %v", cfg.SchedulerInterval)
	}
	if cfg.SchedulerEnabled {
		t.Errorf("expected scheduler enabled to be false")
	}
}
