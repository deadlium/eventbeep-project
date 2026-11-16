## Pull Request Checklist

> **Target Branch Notice**: All feature and fix PRs must target **`build`**, NOT `master`.
> The CI/CD God Flow automatically promotes `build` into `master` upon successful test verification.
> **Branch Retention Notice**: Do NOT delete the source branch upon merge. All branches are preserved permanently.

### Phase & Type
- [ ] Phase 1: Database & Go Backend Setup
- [ ] Phase 2: Scraper Engine
- [ ] Phase 3: LLM Integration & Scheduler
- [ ] Phase 4: Next.js Frontend
- [ ] Infrastructure / Docs / CI

### Verification
- [ ] Code compiles and passes tests locally
- [ ] Follows Conventional Commits standard
- [ ] No hardcoded API keys or secrets
