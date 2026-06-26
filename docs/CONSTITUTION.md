# CONSTITUTION

Scout is a deterministic research-intelligence pipeline that turns the Hugging Face Daily Papers feed into auditable daily and monthly executive briefings. This document is the canonical contract for how Scout is designed, built, and changed. When code and this document disagree, update one or the other in the same change so they stay aligned.

## MISSION

- Replace title-only paper discovery with concise, signal-dense intelligence for engineers.
- Bias every summary and score toward implementation implications for AI agents, LLM systems, orchestration, evaluation, memory, ML infrastructure, and software engineering.
- Maintain a continuously reranked monthly leaderboard instead of static daily ordering.
- Preserve durable, inspectable paper records, score breakdowns, and observation history.
- Keep execution deterministic, incremental, and easy to audit.

## PRINCIPLES

- **Determinism first.** Identical inputs must produce byte-identical outputs. No randomness, no unsorted map iteration in output, no wall-clock or locale dependence in scoring, summarization, or rendering. Sort before emitting and use UTC for all dates.
- **Incremental and idempotent.** Re-running any day must be safe and convergent. Reuse existing per-paper records, append (never overwrite) score history, and merge observation dates without duplication.
- **Auditable artifacts.** Every run writes durable, human-readable JSON and Markdown under the data directory. Scores expose their per-dimension breakdown and weights so any result can be explained.
- **Standard library first.** Reach for the Go stdlib before adding dependencies. New direct dependencies require a clear, documented justification.
- **Tolerant ingestion, strict output.** Upstream feeds are messy: parse defensively with key fallbacks and graceful degradation. Internal records and reports are strict and well-typed.
- **Layered separation of concerns.** Presentation (CLI) → orchestration (pipeline) → domain logic (scoring, summary) → I/O (hf client, storage). Business logic never leaks into the CLI or the storage layer.
- **Dependency injection at boundaries.** Collaborators are injected and depend on consumer-defined interfaces, so external systems (HTTP, filesystem) are swappable and testable.
- **Small, single-purpose files.** One concept per file; co-locate a type with its helpers, constants, and methods. Prefer source files at or under 300 lines.
- **Repo-local docs are the system of record.** Route through `AGENTS.md` → `docs/agents/README.md`; load only what the current decision needs (see `docs/agents/RLM.md`).

## ARCHITECTURE

Go module `github.com/jamesonstone/scout`, currently targeting the Go version pinned in `go.mod`. A single `scout` binary built from `cmd/scout`.

Layer map:

- `cmd/scout/main.go` — thin entrypoint; delegates to `pkg/cli.Execute()`.
- `pkg/cli/` — Cobra command surface (`root.go`, `run.go`). Owns flag/env resolution, exit codes via `cliExitError`, and wiring of dependencies. The only intentionally public package.
- `internal/config/` — configuration from `SCOUT_*` environment variables with sensible defaults, overridable by CLI flags; resolves and validates the run date and data directory.
- `internal/hf/` — Hugging Face HTTP client: retry/backoff, configurable timeout and user agent, and defensive JSON/Markdown parsing.
- `internal/model/` — domain types and JSON contracts (`Paper`, `PaperRecord`, `ScoreBreakdown`, `Weights`, `ScoreSnapshot`, `DailyObservation`, `RunResult`, `Recommendation`). Single source of truth for the persisted schema.
- `internal/pipeline/` — orchestration `Runner`, dedup/merge logic, and report rendering (`render.go`); re-exports `model` types via aliases.
- `internal/scoring/` — deterministic weighted keyword scorer and recommendation thresholds.
- `internal/summary/` — deterministic, template-free text builders for innovation summary, why-it-matters, implementation angle, caveat, executive summary, and signal.
- `internal/storage/` — filesystem persistence of per-paper JSON, daily observations, and Markdown reports.
- `internal/prompt/` — embedded `*.tmpl` assets via `go:embed`.

Data flow: `cli.run` builds `config`, `hf.Client`, and `storage.Store`, constructs a `pipeline.Runner`, and calls `Run(ctx)`. The runner fetches the daily feed, dedups, reuses or enriches each paper, scores and summarizes new papers, persists records and the daily observation, then renders the monthly and daily reports.

### Persistence Layout

Durable artifacts live beneath the configured data directory (default `.scout`):

- `data/papers/<paper-id>.json` — per-paper record, summaries, links, score, and score history.
- `data/daily/YYYY-MM/YYYY-MM-DD.json` — observed paper IDs for a daily run.
- `reports/daily/YYYY-MM/YYYY-MM-DD.md` — daily executive briefing.
- `reports/monthly/YYYY-MM.md` — continuously reranked monthly briefing.

File paths are derived deterministically; paper IDs are sanitized for filesystem safety. JSON is written pretty-printed with a trailing newline; directories are created on demand.

## DOMAIN MODEL AND SCORING

- **Temporal provenance is mandatory on records.** Every `PaperRecord` carries `FirstSeen`, `ObservedDates`, and `ScoreHistory`. New data models must preserve equivalent creation/observation/update provenance; do not introduce records without durable temporal tracking.
- **Scoring weights must sum to 1.0** and every component and overall score is clamped to `0–100`. Current weights: novelty `0.20`, practical impact `0.20`, technical depth `0.15`, implementation potential `0.15`, relevance `0.15`, community signal `0.10`, summary confidence `0.05`.
- **Recommendation thresholds:** `>= 80` → Read, `>= 60` → Worth Watching, otherwise Skip. Reading-priority bands: `>= 85`, `>= 70`, `>= 55`, else low.
- Scoring is pure keyword/metadata heuristics over lowercased text. Changing weights, keywords, or thresholds is a behavioral change: update tests and any affected docs in the same change.
- Summaries are deterministic string construction, not generated prose. No external model calls.

## CODE STYLE

- Format with `gofmt`; keep `go vet` and `golangci-lint` clean (`make fmt vet lint`).
- Comments are sparse and begin lowercase; prefer self-explaining code and descriptive names over narration.
- Constructors return value types (`New(...) T`) consistent with existing packages; avoid pointers unless mutation or identity requires them.
- Wrap errors with context using `fmt.Errorf("...: %w", err)`; prefer explicit handling over silent failure.
- Resolve config as env-defaults first, then apply CLI flags only when `Changed`.
- Keep packages cohesive and files focused; split when a file grows past ~300 lines and the split improves clarity.

## DEPENDENCIES

- `github.com/spf13/cobra` — CLI command framework (direct).
- `github.com/spf13/pflag` — POSIX flags (indirect, via cobra).
- `github.com/inconshreveable/mousetrap` — Windows entrypoint guard (indirect, via cobra).
- Everything else (HTTP, JSON, regexp, embed, filesystem, time) is Go standard library. Keep it that way unless a dependency is clearly justified and recorded.

## CONSTRAINTS

- Determinism, incrementality, and idempotency are non-negotiable invariants.
- Runs must be reproducible and offline-testable: `--base-url` / `SCOUT_BASE_URL` must allow pointing at a local fixture server with no hard dependency on the live Hugging Face API.
- No secrets, credentials, or PII in code, config, logs, or artifacts. Hugging Face endpoints used are public; the scheduled workflow runs with `contents: read` only.
- Preserve the persisted JSON schema and on-disk layout; schema changes must remain backward-compatible for existing records or include a migration path.
- Keep `cmd` and `cli` thin; do not place scoring, summarization, or persistence logic in the command layer.

### Kit-Managed Baseline Rules

<!-- BEGIN KIT-MANAGED BASELINE RULES -->
- Treat `docs/CONSTITUTION.md` as the canonical project contract.
- Keep `AGENTS.md`, `CLAUDE.md`, and `.github/copilot-instructions.md` aligned with the repo-local docs tree.
- Prefer implementation/source code files around 300 lines or less when splitting improves clarity and ownership.
- Do not apply the code-file size guideline to documentation files, all `docs/**`, all `.kit/**`, or `.kit.yaml`.
- Do not split or rewrite docs, generated state, or Kit config artifacts solely because they exceed 300 lines.
<!-- END KIT-MANAGED BASELINE RULES -->

## VALIDATION

- `make test` — run the full test suite; never claim tests passed unless they ran.
- `make build` — must produce `bin/scout`.
- End-to-end smoke against a local fixture: `go run ./cmd/scout run --date 2026-01-02 --data-dir /tmp/scout-output --base-url http://127.0.0.1:8080`.
- Expected observable results: deterministic scoring/persistence/dedup/render tests pass; a daily report, a monthly report, and per-paper JSON records are written under the chosen data directory.
- Fix relevant lint and test failures before calling work complete.

## CHANGE CLASSIFICATION

All work falls into one of two tracks — classify before acting. See `docs/agents/WORKFLOWS.md` for the full flow and `docs/agents/GUARDRAILS.md` for hard rules.

### Spec-Driven (Formal)

- Use when: new features, substantial architectural or behavioral changes, cross-component changes, or work that already has feature docs.
- Workflow: `kit spec <feature>` → `docs/specs/<feature>/SPEC.md` phases: clarify → ready → implement → validate → reflect → deliver.
- `SPEC.md` is the single durable feature artifact; front matter carries `workflow_version: 2` and the current `phase`.
- Legacy staged documents (`BRAINSTORM.md`, legacy `SPEC.md`, `PLAN.md`, `TASKS.md`) are historical context only, used when a legacy staged command is explicitly chosen.
- Never mix multiple features in one `docs/specs/<feature>/` directory.

### Ad Hoc (Lightweight)

- Use when: bug fixes, security reviews, refactors, dependency updates, config changes, small refinements.
- Workflow: understand → implement → verify with the smallest relevant checks.
- Docs: update only the practical docs that changed (README, inline docs, this constitution when invariants shift).
- Do not create a feature `SPEC.md` or legacy staged artifacts for ad hoc work.

### Ad Hoc with Existing Specs

- If a change touches code with existing spec docs, default to updating them.
- Skip spec updates only for purely mechanical changes (formatting, typo, dependency bump).

## PROGRESS TRACKING

- `docs/PROJECT_PROGRESS_SUMMARY.md` must reflect the highest completed artifact per feature at all times, and is the shortlist index for prior feature discovery (see `docs/agents/RLM.md`).
- Do not generate periodic standalone progress reports or separate summary documents; update the canonical artifact instead.

## DELIVERY

- This is a Kit-managed project: issue, branch, staging, commit, push, and PR actions are mutation boundaries gated by `docs/agents/GUARDRAILS.md` and `docs/references/rules/*`.
- Never run `git add`, `git commit`, `git push`, or history-rewriting commands without explicit, current approval.
- PR titles use Conventional Commits with the GitHub issue as scope: `<type>(<issue_number>): <gitmoji> <short title message>`.
- PRs use the repo template at `.github/pull_request_template.md`; do not use generic branches, ad hoc bodies, or default drafts when Kit rules define the contract.
- CodeRabbit review scope excludes `docs/**`, `AGENTS.md`, and `CLAUDE.md` per `.coderabbit.yaml`.
- Markdown is formatted with Prettier; do not emit compact or unpadded Markdown tables.

## NON-GOALS

- No external LLM or AI inference for scoring or summarization; all enrichment is deterministic local heuristics.
- No database, queue, or external service dependency; persistence is the local filesystem only.
- No web UI, API server, or long-running daemon; Scout is a batch CLI invoked on demand or by cron.
- No authentication, accounts, or handling of private/credentialed data.
- Not a general-purpose web scraper; ingestion is scoped to the Hugging Face Daily Papers feed and its paper endpoints.
- No real-time or streaming processing.

## DEFINITIONS

- **Paper** — an ingested feed item (`model.Paper`): transient fetch/enrichment shape before persistence.
- **PaperRecord** — the durable persisted artifact (`model.PaperRecord`) with summaries, links, score, and history.
- **Daily Observation** — the set of paper IDs seen on a given day (`data/daily/...`).
- **Daily / Monthly Briefing** — rendered Markdown reports; the monthly briefing is reranked after each daily run.
- **Score Breakdown** — per-dimension scores plus weights producing a `0–100` overall score.
- **Recommendation** — `Read`, `Worth Watching`, or `Skip`, derived from the overall score.
- **Kit-managed project** — a repo containing `.kit.yaml`, `docs/CONSTITUTION.md`, or `docs/agents/README.md`; delivery follows repo-local Kit rules over generic defaults.
- **RLM** — Kit's just-in-time, progressive-disclosure context routing pattern (`docs/agents/RLM.md`).
