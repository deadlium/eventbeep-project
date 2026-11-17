# EventBeep Backend Service

Go REST API and data ingestion service for EventBeep.

## Structure
- `cmd/api/`: Main HTTP server entry point (`GET /api/jobs`, health check, scheduler bootstrap)
- `internal/database/`: PostgreSQL connection pool (`pgx`) and query handlers
- `internal/models/`: Job entity and request/response models
- `internal/repository/`: Data access layer for jobs with PostgreSQL idempotent upsert (`ON CONFLICT (hash_key) DO NOTHING`)
- `internal/scraper/`: Web scrapers (`goquery`), SHA-256 deduplication engine, and HTML fixture tests
- `internal/llm/`: Gemini API entity extraction service (`{"role": "string", "city": "string"}`) with rule-based fallback
- `internal/scheduler/`: 24-hour background ticker orchestrating Scrape -> Deduplicate -> LLM Enrich -> Store

## Scraper Engine & Deduplication (Phase 2)
- **Scrapers**: Extensible `Scraper` interface with implementations for:
  - `UnstopScraper`: Targets fresher/entry-level listings (`unstop.com/jobs`)
  - `WellfoundScraper`: Targets startup fresher/junior roles (`wellfound.com/jobs`)
  - `InternshalaScraper`: Targets fresher jobs & internships (`internshala.com/jobs`)
- **Deduplication**: Deterministic 64-char hex SHA-256 hash calculated as:
  `SHA256(lowercase(company_name + job_title))`
  Stored in `jobs.hash_key` with unique constraint.
- **Engine**: Ingestion coordinator supporting live scrapes, deterministic HTML fixture parsing (`io.Reader`), and batch deduplication before/during PostgreSQL persistence.

## Gemini LLM Entity Extraction & 24h Scheduler (Phase 3)
- **Gemini Extractor (`internal/llm`)**:
  - Calls Gemini API with structured JSON schema system instructions:
    `Extract 'role' and 'city' into JSON: {"role": "string", "city": "string"}`
  - Parses structured output with markdown code-fence sanitization.
  - Transparent heuristic fallback using compiled regex patterns for tech roles and major Indian tech hubs/remote when `GEMINI_API_KEY` is unset or on transient network/API errors.
- **24h Scheduler (`internal/scheduler`)**:
  - Background ticker orchestrating: `Scrape -> Deduplicate -> LLM Enrich -> Store in DB`.
  - Configurable interval via `SCHEDULER_INTERVAL` (defaults to 24 hours), graceful shutdown on OS signals or context cancellation.

## Running Tests
Run all unit and integration tests with coverage:
```bash
go test -v -cover ./...
```
