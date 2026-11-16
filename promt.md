Act as a Staff Software Engineer. We are building an automated job board aggregator to parse fresher jobs and internships in India. 

The architecture is strictly decoupled:
1. Backend: Go (using standard library net/http or Gin) and PostgreSQL.
2. Frontend: Next.js (TypeScript, Tailwind CSS).

Execute the following phases step-by-step. Do not proceed to the next phase until I confirm the current one is working.

PHASE 1: Database & Go Backend Setup
- Create a PostgreSQL schema with a `jobs` table: id (UUID), hash_key (VARCHAR UNIQUE), title (VARCHAR), company (VARCHAR), role (VARCHAR), city (VARCHAR), source_url (VARCHAR), created_at (TIMESTAMP).
- Scaffold a Go project with a Postgres connection pool.
- Create a REST endpoint `GET /api/jobs` that supports pagination and filtering by `role` and `city`.

PHASE 2: The Scraper Engine (Go)
- Write a Go web scraper module using `goquery` or `colly`.
- Target 2 mock or real HTML structures for public job boards (e.g., unstop, wellfound).
- Implement a deduplication function: generate a SHA-256 hash of `lowercase(company_name + job_title)` and use it as `hash_key` to upsert/ignore duplicates in PostgreSQL.

PHASE 3: LLM Integration (Go)
- Write a service that takes the scraped raw title and description.
- Call the Gemini (or OpenAI) API using the system prompt: "Extract 'role' and 'city' into JSON: {\"role\": \"string\", \"city\": \"string\"}".
- Parse the JSON response and update the `role` and `city` columns in the database.
- Create a Go ticker/cron job to run Phase 2 and Phase 3 once every 24 hours.

PHASE 4: Next.js Frontend
- Initialize a Next.js project with Tailwind CSS.
- Build a clean, searchable, and filterable data table that fetches from `GET /api/jobs`.
- Add a search bar (debounced) and dropdowns for City and Role.
- Ensure the UI is responsive and minimal.