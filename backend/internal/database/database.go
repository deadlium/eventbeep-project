package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/deadlium/eventbeep/backend/internal/config"
)

// DB encapsulates the pgx connection pool.
type DB struct {
	Pool *pgxpool.Pool
}

// NewPool initializes and configures a PostgreSQL connection pool using pgxpool.
func NewPool(ctx context.Context, cfg *config.Config) (*DB, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database configuration: %w", err)
	}

	// Apply production connection pool tuning
	poolConfig.MaxConns = cfg.DBMaxConns
	poolConfig.MinConns = cfg.DBMinConns
	poolConfig.MaxConnIdleTime = cfg.DBMaxIdle
	poolConfig.MaxConnLifetime = cfg.DBMaxLife
	poolConfig.HealthCheckPeriod = 30 * time.Second

	// Connect with timeout
	connectCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(connectCtx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgxpool: %w", err)
	}

	// Verify connection
	if err := pool.Ping(connectCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("[database] PostgreSQL connection pool initialized successfully (max_conns=%d, min_conns=%d)", cfg.DBMaxConns, cfg.DBMinConns)
	return &DB{Pool: pool}, nil
}

// EnsureSchema runs idempotent table and index creation for the jobs table.
func (db *DB) EnsureSchema(ctx context.Context) error {
	schemaSQL := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	CREATE EXTENSION IF NOT EXISTS pg_trgm;

	CREATE TABLE IF NOT EXISTS jobs (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		hash_key VARCHAR(64) UNIQUE NOT NULL,
		title VARCHAR(255) NOT NULL,
		company VARCHAR(255) NOT NULL,
		role VARCHAR(100),
		city VARCHAR(100),
		source_url VARCHAR(1024) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_jobs_hash_key ON jobs(hash_key);
	CREATE INDEX IF NOT EXISTS idx_jobs_role ON jobs(role);
	CREATE INDEX IF NOT EXISTS idx_jobs_city ON jobs(city);
	CREATE INDEX IF NOT EXISTS idx_jobs_created_at ON jobs(created_at DESC);
	`

	_, err := db.Pool.Exec(ctx, schemaSQL)
	if err != nil {
		return fmt.Errorf("failed to execute ensure schema: %w", err)
	}
	log.Println("[database] Database schema and indexes verified.")
	return nil
}

// Ping verifies active connectivity to PostgreSQL.
func (db *DB) Ping(ctx context.Context) error {
	if db.Pool == nil {
		return fmt.Errorf("database pool is not initialized")
	}
	return db.Pool.Ping(ctx)
}

// Close gracefully closes all idle and active pool connections.
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
		log.Println("[database] PostgreSQL connection pool closed.")
	}
}
