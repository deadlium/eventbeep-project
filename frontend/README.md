# EventBeep Frontend (Phase 4)

Responsive, accessible, and high-performance job aggregation web application built with **Next.js 16 (App Router)**, **TypeScript**, and **Tailwind CSS 4**.

Consumes the EventBeep Go REST API (`GET /api/jobs`) to surface fresher opportunities and internships aggregated from across India (Unstop, Wellfound, Internshala), deduplicated with SHA-256 signatures and enriched with Google Gemini LLM.

---

## Features

- **Real-Time Job Table & Mobile Card View**: Responsive layout with role badges, company info, city indicators, discovery timestamps, and external application links.
- **Debounced Search**: Instant filtering across job titles, tech stacks, and company names with a 350ms debounce window.
- **Role & City Dropdowns**: Instant categorization by engineering roles (Backend, Frontend, Full Stack, DevOps, Data/AI, Mobile, QA) and Indian tech hubs (Bangalore, Pune, Hyderabad, Mumbai, Delhi NCR, Chennai, Remote).
- **Comprehensive UI States**:
  - **Loading Skeleton**: Smooth animated pulse placeholder table rows and mobile cards.
  - **Empty State**: Informative empty view with one-click filter reset.
  - **Error & Retry Banner**: Graceful notification with instant retry button.
- **Pagination & Page Sizing**: Previous/Next navigation, page number jump controls, items-per-page selector (10, 20, 50), and record count indicator.
- **Resilient Fallback Mode**: Gracefully transitions to curated fresher datasets when the local Go backend is offline or unreachable during development.

---

## Project Structure

```
frontend/
├── app/
│   ├── globals.css          # Tailwind CSS 4 theme and design tokens
│   ├── layout.tsx           # Root metadata and font configuration
│   └── page.tsx             # Main interactive search & job board page
├── components/
│   ├── Header.tsx           # Brand header with live pipeline indicator
│   ├── SearchFilters.tsx    # Debounced search input, dropdowns, and limit picker
│   ├── JobTable.tsx         # Responsive table and mobile card listing
│   ├── JobSkeleton.tsx      # Table and card pulse loading skeletons
│   ├── EmptyState.tsx       # Zero-match state with filter reset action
│   ├── ErrorBanner.tsx      # API error banner with retry trigger
│   ├── PaginationControls.tsx # Pagination bar and page numbers
│   └── Icons.tsx            # Crisp lightweight SVG icon library
├── hooks/
│   └── useDebounce.ts       # Generic debounce hook for search inputs
├── lib/
│   └── api.ts               # Go API client and resilient fallback dataset
└── types/
    └── job.ts               # Job, Pagination, and Filter TypeScript interfaces
```

---

## Configuration

Set the backend API URL using the `NEXT_PUBLIC_API_URL` environment variable:

```bash
# Default: http://localhost:8080
NEXT_PUBLIC_API_URL=http://localhost:8080
```

---

## Verification & Commands

```bash
# Run ESLint validation
npm run lint

# Compile optimized production build
npm run build

# Start development server
npm run dev
```
