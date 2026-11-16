# EventBeep Backend Service

Go REST API and data ingestion service for EventBeep.

## Structure
- `cmd/api/`: Main HTTP server entry point (`GET /api/jobs`, health check)
- `internal/database/`: PostgreSQL connection pool (`pgx`) and query handlers
- `internal/models/`: Job entity and request/response models
- `internal/scraper/`: Web scrapers (`goquery` / `colly`) and SHA-256 deduplication
- `internal/llm/`: Gemini API entity extraction service
- `internal/scheduler/`: 24-hour background ticker / cron
