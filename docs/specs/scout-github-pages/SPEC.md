---
workflow_version: 2
phase: implemented
title: Scout GitHub Pages static site
issue: GH-4
skills: []
---

# Scout GitHub Pages Static Site

## Thesis

Scout remains a standalone repository and publishes a read-only GitHub Pages project site at the project path, expected as `/scout/`. The site is generated from Scout's committed curated `data/` and `reports/` artifacts and does not fetch Hugging Face, run Go code, write storage in the browser, or serve full paper text.

## Current Code Map

- CLI root: `pkg/cli/root.go`, `rootCmd`, `Execute`, `Version`.
- Daily pipeline command: `pkg/cli/run.go`, `runPipeline`, `scout run`.
- Configuration: `internal/config/config.go`, `Config`, `FromEnv`, `ResolveRunDate`; default data directory is `.scout`.
- Hugging Face fetching: `internal/hf/client.go`, `Client`, `FetchDailyPapers`, `FetchPaperDetails`, `FetchMarkdown`.
- Storage: `internal/storage/store.go`, `Store`, `SavePaper`, `LoadPaper`, `SaveObservation`, `MonthRecords`, `SaveDailyReport`, `SaveMonthlyReport`.
- Data model: `internal/model/types.go`, `PaperRecord`, `DailyObservation`, `ScoreBreakdown`, `Recommendation`.
- Scoring: `internal/scoring/scorer.go`, `Scorer.Score`, `Recommendation`.
- Summary generation: `internal/summary/summary.go`, `Builder` methods.
- Markdown report rendering: `internal/pipeline/render.go`, `RenderDaily`, `RenderMonthly`.
- Pipeline orchestration: `internal/pipeline/runner.go`, `Runner.Run`.
- Validation commands: `make fmt`, `make lint`, `make test`, `make build`.
- Existing committed artifact paths: `data/papers/*.json`, `data/daily/2026-06/2026-06-26.json`, `reports/daily/2026-06/2026-06-26.md`, `reports/monthly/2026-06.md`.
- Existing workflow: `.github/workflows/daily-scout.yml` runs a manual smoke check with `SCOUT_DATA_DIR=.scout`; it does not publish Pages or commit artifacts.

## Curated Storage Contract

Scout is a curated intelligence artifact, not an archive of papers.

- Product question: Scout answers "Should I care?" Official paper pages answer "What are all the details?"
- Transient inputs: the run may fetch markdown, abstract, author metadata, and community signals for scoring and summarization.
- Committed paper records: store only compact distilled output: title, publish date when available, first/observed dates, score, recommendation, capped categories, innovation summary, why it matters, implementation angle, caveat, and source links.
- Forbidden committed fields: `markdown`, `abstract`, `authors`, `community`, `score_history`, `metadata_completeness`, rendered executive prose, and reading-priority prose.
- Size limits: paper JSON must be 8KB or smaller; category lists are capped; daily reports cap full summaries to the highest-signal papers; archive/lower-signal rows stay one-line.
- Retention stance: raw paper text belongs in ignored `.scout/cache/` or temp storage only. Month-close compaction should preserve this contract for old low-signal records.

## Reconciliation

- Live behavior matches the desired durable artifact layout: `data/papers/<paper-id>.json`, `data/daily/YYYY-MM/YYYY-MM-DD.json`, `reports/daily/YYYY-MM/YYYY-MM-DD.md`, and `reports/monthly/YYYY-MM.md`.
- RESOLVED: the repository has a manual GitHub Action smoke workflow, but it writes to `.scout` with read-only contents permission and is not the durable publishing/update lane. Codex Scheduled Automation is the research-update delivery mechanism; GitHub Actions Pages only deploys committed artifacts after merge.
- CONFLICT: using `docs/` as the Pages source would mix generated site output with repo-local agent documentation. The selected publish directory is `public/`, uploaded through GitHub Actions Pages.
- UNKNOWN: repository Pages settings are not represented in source. Jameson Stone owns enabling GitHub Pages with GitHub Actions as the source if it is not already enabled.

## Static Site Architecture

- Build command: `go run ./cmd/scout site build --data-dir . --out-dir public --base-path /scout/`.
- Validate command: `go run ./cmd/scout site validate --out-dir public --base-path /scout/`.
- Publish directory: `public/`, generated deterministically and suitable for `actions/upload-pages-artifact`.
- Base path: `/scout/`, configurable with `--base-path` for local previews or repository renames.
- Data strategy: copy compact JSON artifacts from `data/` to `public/data/` and write `public/data/index.json` as a static manifest. `site build` fails if source paper records violate the curated storage contract.
- Asset strategy: write `public/assets/styles.css` and `public/.nojekyll`; no JavaScript is required.
- Routes:
  - `/scout/` -> homepage focused on Scout-branded latest fetched papers, grouped by date.
  - `/scout/daily/` -> daily archive organized by Scout fetch date.
  - `/scout/daily/YYYY-MM-DD/` -> daily briefing page for papers fetched on that date.
  - `/scout/monthly/` -> monthly archive.
  - `/scout/monthly/YYYY-MM/` -> monthly ranking page.
  - `/scout/papers/<paper-id>/` -> paper detail page.
  - `/scout/data/index.json`, `/scout/data/papers/*.json`, `/scout/data/daily/YYYY-MM/*.json` -> durable JSON.

## File Edits

- Add `pkg/cli/site.go` for `scout site build` and `scout site validate`.
- Add `internal/artifact/` for the curated committed-artifact contract, paper JSON size limit, forbidden raw fields, and compact-list limits.
- Add `internal/site/` for static site loading, rendering, styles, validation, and tests.
- Update `internal/storage/` and `internal/pipeline/` so committed paper records are compact curated summaries, while raw fetched text remains transient.
- Update `Makefile` with `site-build` and `site-validate`.
- Add `.github/workflows/pages.yml` for GitHub Actions Pages artifact deployment from `public/`.
- Update `README.md` with site generation, Pages deployment, and Codex automation flow.
- Add this `docs/specs/scout-github-pages/SPEC.md` plan.
- Generate `public/` from committed `data/` and `reports/` artifacts for validation and PR review.

## Acceptance Criteria Mapping

- Current file/symbol map is recorded in this spec.
- Static site architecture, route structure, data strategy, base path, and publish directory are recorded in this spec.
- `scout site build` consumes existing artifacts and writes a static output directory.
- `scout site validate` checks required pages, required daily/monthly sections, internal base-path links, JSON score objects, paper JSON size, and forbidden raw fields.
- Homepage and daily pages show Scout as the product, not repository documentation.
- Daily pages show Executive Signal, Top Papers, Additional Papers, Watchlist, Archive, score, recommendation, publish date when available, innovation summary, why it matters, implementation angle, caveat, and source links.
- Paper detail pages show publish date, first fetched date, observed dates, source links, and score breakdowns.
- Monthly pages show Top 10 papers, Rising Papers, Themes, and Complete Monthly Index.
- Site validation fails if the homepage contains README-style project documentation instead of the generated Scout paper briefing.
- Committed paper records and generated `public/data/papers/*.json` omit raw markdown, full abstracts, authors, raw community blobs, score history arrays, metadata completeness counters, and duplicated rendered prose.
- The Pages workflow publishes static output only; it does not fetch Hugging Face at page view time.

## Validation Plan

- `make fmt`: formats all Go source.
- `make lint`: verifies lint cleanliness, including the new site package.
- `make test`: covers CLI-adjacent package behavior, rendering, storage, scoring, summary, Hugging Face parsing, and site build/validation.
- `make build`: ensures the Scout binary still builds.
- `go run ./cmd/scout run --date 2026-06-26 --data-dir .`: validates the daily pipeline against the committed durable artifact directory.
- `go run ./cmd/scout site build --data-dir . --out-dir public --base-path /scout/`: writes the Pages output.
- `go run ./cmd/scout site validate --out-dir public --base-path /scout/`: validates generated pages, sections, links, and JSON score data.

## Codex Automation Flow

1. Run Scout for the target date with an explicit data directory.
2. Persist daily observation, paper records, daily report, and monthly report.
3. Build the static site into `public/`.
4. Validate the generated site.
5. Deliver changed `data/`, `reports/`, and `public/` artifacts through the Kit issue, `GH-123` branch, commit, push, and ready PR workflow.
6. After merge, GitHub Actions Pages deploys the static `public/` artifact. The repository Pages source must be GitHub Actions, not legacy `main`/root.

## Risks And Owners

- GitHub Pages source setting: owner Jameson Stone; mitigation is `.github/workflows/pages.yml` and this spec documenting GitHub Actions Pages as the expected source.
- Project-site base path drift: owner Codex for code; mitigation is configurable `--base-path` plus link validation.
- Generated artifact growth: owner Codex for repo hygiene; mitigation is deterministic generated paths and future pruning policy if history grows too large.
- Hugging Face API payload drift: owner external API plus Codex tests; mitigation is existing parser coverage and daily run validation.
- Report quality regressions: owner Codex; mitigation is site validation plus existing summary/render tests.
- Personal site integration: owner Jameson Stone; mitigation is link-only guidance in `README.md`, no generated Scout artifacts copied into `jamesonstone.github.io`.
