package config_test

import (
	"os"
	"testing"

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
}

func TestConfig_CustomEnv(t *testing.T) {
	os.Setenv("PORT", "9090")
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/testdb")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("DATABASE_URL")
	}()

	cfg := config.Load()

	if cfg.Port != "9090" {
		t.Errorf("expected port 9090, got %s", cfg.Port)
	}
	if cfg.DatabaseURL != "postgres://user:pass@localhost:5432/testdb" {
		t.Errorf("expected custom database url, got %s", cfg.DatabaseURL)
	}
}
