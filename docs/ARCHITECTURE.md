# EventBeep Architecture Specification

EventBeep is an automated job board aggregator designed to collect, enrich, deduplicate, and display fresher jobs and internships across India.

## 1. High-Level System Architecture

```mermaid
flowchart TD
    subgraph Ingestion ["Ingestion & Scraper Engine (Go)"]
      S1[Scraper: Unstop / Wellfound / etc.]
      DEDUP[SHA-256 Deduplication: lowercase(company+title)]
      LLM[LLM Entity Extractor: Gemini API]
      CRON[24-Hour Cron / Ticker]
    end

    subgraph Storage ["Database (PostgreSQL)"]
      DB[(jobs table with UUID, hash_key, role, city)]
    end

    subgraph API ["REST API Service (Go)"]
      REST[GET /api/jobs?role=&city=&page=&limit=]
    end

    subgraph Client ["Frontend (Next.js & Tailwind CSS)"]
      UI[Responsive Job Board UI]
      FILTER[Debounced Search + Dropdowns]
    end

    CRON --> S1
    S1 --> DEDUP
    DEDUP --> LLM
    LLM --> DB
    REST --> DB
    UI --> FILTER
    FILTER --> REST
```

## 2. Decoupled Components

### 2.1 Backend (Go + PostgreSQL)
- **Runtime**: Go 1.23+ using standard library `net/http` or `gin-gonic/gin`.
- **Database Driver**: `pgx` connection pool for high-concurrency database queries.
- **REST Endpoints**:
  - `GET /api/jobs`: Supports pagination (`page`, `limit`) and query filters (`role`, `city`, `q`).
  - `GET /health`: Health check endpoint.

### 2.2 Scraper & Deduplication Engine (Go)
- **Scraper Engine**: `goquery` / `colly` to parse HTML from target job boards.
- **Deduplication**: Computes SHA-256 of `lowercase(company_name + job_title)` to generate a 64-character hex `hash_key`.
- **Idempotency**: Upsert / `ON CONFLICT (hash_key) DO NOTHING` prevents redundant entries.

### 2.3 LLM Extraction Service (Gemini API)
- **Model**: Google Gemini API (or OpenAI API fallback).
- **Prompt**:
  ```
  Extract 'role' and 'city' into JSON: {"role": "string", "city": "string"}
  ```
- **Execution**: Asynchronously runs on newly discovered jobs and updates database columns.
- **Scheduler**: Go `time.Ticker` background worker triggers every 24 hours.

### 2.4 Frontend (Next.js + TypeScript + Tailwind CSS)
- **Framework**: Next.js App Router with React Server Components and Client Components where interactive.
- **Styling**: Tailwind CSS for a minimal, clean, responsive UI.
- **Features**:
  - Filterable data table fetching from `GET /api/jobs`.
  - Debounced search bar.
  - Dropdown filters for `City` and `Role`.
  - Pagination and empty-state handling.
