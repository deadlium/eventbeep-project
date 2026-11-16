# Contributing to EventBeep

Thank you for contributing to EventBeep! We maintain high engineering standards across Go, Next.js, and automated data pipelines.

---

## 1. The GitHub Flow ("God Flow")

This repository strictly adheres to a three-tier branch lifecycle:

```
feature/<name> / fix/<name>
             │
             │  (Pull Request into build)
             ▼
           build  (Integration branch where all sub-branches merge)
             │
             │  (Automated CI pass -> Auto promotion)
             ▼
           master (Stable production release branch)
```

### Core Branching Rules:
1. **Branch Targets**:
   - All feature, bugfix, and documentation branches originate from `build`.
   - All Pull Requests must target **`build`**, never `master`.
2. **Automated Promotion**:
   - The CI/CD engine validates all PRs and pushes to `build`.
   - Once tests pass on `build`, the automated workflow promotes `build` to `master`.
3. **Branch Retention Policy (STRICT)**:
   - **DO NOT DELETE ANY BRANCH**.
   - Every feature branch, experiment branch, and release branch is preserved permanently for full auditability and historical tracking. Do not check "Delete branch after merge" on GitHub.

---

## 2. Commit Message Guidelines

All commits must follow the [Conventional Commits](https://www.conventionalcommits.org/) standard:
```
<type>(<scope>): <short imperative summary>

- <bullet point details>
- <additional context or rationale>
```

Allowed types: `feat`, `fix`, `refactor`, `perf`, `test`, `docs`, `style`, `chore`, `build`, `ci`.
Common scopes: `backend`, `frontend`, `db`, `scraper`, `llm`, `api`, `workflow`.

---

## 3. Local Verification

- **Backend**: `cd backend && go test ./...`
- **Frontend**: `cd frontend && npm run build`
