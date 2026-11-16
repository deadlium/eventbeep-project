package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"
)

// Config holds runtime configuration for the service.
type Config struct {
	Port              string
	DatabaseURL       string
	DBMaxConns        int32
	DBMinConns        int32
	DBMaxIdle         time.Duration
	DBMaxLife         time.Duration
	Env               string
	GeminiAPIKey      string
	GeminiModel       string
	SchedulerInterval time.Duration
	SchedulerEnabled  bool
}

// Load reads configuration from environment variables with fallback defaults.
func Load() *Config {
	port := getEnv("PORT", "8080")
	env := getEnv("APP_ENV", "development")

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		host := getEnv("DB_HOST", "localhost")
		dbPort := getEnv("DB_PORT", "5432")
		user := getEnv("DB_USER", "postgres")
		pass := getEnv("DB_PASSWORD", "postgres")
		name := getEnv("DB_NAME", "eventbeep")
		sslmode := getEnv("DB_SSLMODE", "disable")

		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			url.QueryEscape(user),
			url.QueryEscape(pass),
			host,
			dbPort,
			name,
			sslmode,
		)
	}

	maxConns := int32(getEnvAsInt("DB_MAX_CONNS", 25))
	minConns := int32(getEnvAsInt("DB_MIN_CONNS", 5))

	geminiKey := os.Getenv("GEMINI_API_KEY")
	geminiModel := getEnv("GEMINI_MODEL", "gemini-1.5-flash")

	schedulerInterval := 24 * time.Hour
	if intervalStr := os.Getenv("SCHEDULER_INTERVAL"); intervalStr != "" {
		if d, err := time.ParseDuration(intervalStr); err == nil && d > 0 {
			schedulerInterval = d
		}
	}

	schedulerEnabled := true
	if enabledStr := os.Getenv("SCHEDULER_ENABLED"); enabledStr != "" {
		if b, err := strconv.ParseBool(enabledStr); err == nil {
			schedulerEnabled = b
		}
	}

	return &Config{
		Port:              port,
		DatabaseURL:       dbURL,
		DBMaxConns:        maxConns,
		DBMinConns:        minConns,
		DBMaxIdle:         15 * time.Minute,
		DBMaxLife:         1 * time.Hour,
		Env:               env,
		GeminiAPIKey:      geminiKey,
		GeminiModel:       geminiModel,
		SchedulerInterval: schedulerInterval,
		SchedulerEnabled:  schedulerEnabled,
	}
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}
