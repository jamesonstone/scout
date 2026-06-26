```text
███████╗ ██████╗ ██████╗ ██╗   ██╗████████╗
██╔════╝██╔════╝██╔═══██╗██║   ██║╚══██╔══╝
███████╗██║     ██║   ██║██║   ██║   ██║
╚════██║██║     ██║   ██║██║   ██║   ██║
███████║╚██████╗╚██████╔╝╚██████╔╝   ██║
╚══════╝ ╚═════╝ ╚═════╝  ╚═════╝    ╚═╝   daily papers intelligence distillation
```

**`scout` is a research-intelligence pipeline for engineers who want the highest-value AI papers without reading dozens of titles every day.**

Scout ingests the Hugging Face Daily Papers feed, enriches each paper with metadata and markdown during the run, scores it across explicit weighted dimensions, persists compact curated paper records, and publishes executive-style daily and monthly intelligence briefings.

## Why Scout Exists

- 🧠 Replace title-only discovery with concise signal-dense paper intelligence.
- 📈 Maintain a continuously reranked monthly leaderboard instead of static daily ordering.
- 🗃️ Preserve compact paper records, score breakdowns, and observation history without mirroring full papers.
- 🔁 Keep execution deterministic, incremental, and easy to audit.
- 🛠️ Bias summaries toward implementation implications for AI agents, LLM systems, evaluation, memory, and ML infrastructure work.

## Install From Source

```sh
git clone https://github.com/jamesonstone/scout.git
cd scout
make build
./bin/scout --help
```

For a local install:

```sh
make install
scout version
```

## Quick Start

```sh
# run today's deterministic pipeline
scout run

# run a specific day into a custom workspace
scout run --date 2026-01-02 --data-dir /tmp/scout-output

# build the static GitHub Pages site from committed artifacts
scout site build --data-dir . --out-dir public --base-path /scout/

# validate generated pages, links, sections, and JSON score records
scout site validate --out-dir public --base-path /scout/
```

## Output Layout

Scout writes durable artifacts beneath the configured data directory:

- `data/papers/<paper-id>.json` — compact curated per-paper records with score, recommendation, distilled summaries, links, and observation dates.
- `data/daily/YYYY-MM/YYYY-MM-DD.json` — observed paper IDs for each daily run.
- `reports/daily/YYYY-MM/YYYY-MM-DD.md` — daily executive briefing.
- `reports/monthly/YYYY-MM.md` — continuously reranked monthly briefing.

## Storage Contract

Scout answers "Should I care?" The official paper answers "What are all the details?"

Committed artifacts are curated summaries only:

- raw `markdown`, full `abstract`, full `authors`, raw community blobs, score history arrays, and duplicated rendered prose are not committed in paper JSON;
- raw paper text may be fetched during a run for scoring and summarization, but stays transient in memory or ignored `.scout/cache/` paths;
- paper JSON records are capped at 8KB, with category and link lists capped before persistence;
- daily reports give full treatment only to the highest-signal papers, while lower-signal and archive entries stay compact;
- the static site optimizes for title, score, recommendation, concrete innovation, why it matters, implementation angle, caveat, and links out.

## Deterministic Scoring Model

Each paper receives a 0-100 component score and weighted final score across:

- Novelty
- Practical impact
- Technical depth
- Implementation potential
- Relevance to agents, LLM systems, orchestration, evaluation, memory, ML infrastructure, and software engineering
- Community signal
- Summary confidence

## Configuration

Environment variables:

- `SCOUT_DATA_DIR` — output root (default `.scout`)
- `SCOUT_BASE_URL` — Hugging Face base URL (default `https://huggingface.co`)
- `SCOUT_RUN_DATE` — override run date (`YYYY-MM-DD`)
- `SCOUT_HTTP_TIMEOUT` — request timeout (default `30s`)
- `SCOUT_HTTP_RETRIES` — retry count (default `3`)
- `SCOUT_HTTP_RETRY_WAIT` — wait between retries (default `2s`)
- `SCOUT_HTTP_USER_AGENT` — request user-agent (default `scout/1`)

CLI flags on `scout run` can override `data-dir`, `date`, `base-url`, `timeout`, and `retries`.

## Validation

```sh
make test
make build
go run ./cmd/scout run --date 2026-01-02 --data-dir /tmp/scout-output --base-url http://127.0.0.1:8080
```

Expected observable results:

- tests pass for deterministic scoring, persistence, duplicate handling, and report rendering;
- `make build` produces `bin/scout`;
- `scout run` writes a daily report, a monthly report, and per-paper JSON records under the chosen data directory.

## Scheduled Automation

`.github/workflows/daily-scout.yml` runs the pipeline on a daily cron schedule and supports manual dispatch.

Codex Scheduled Automation is the intended durable research-update path: run Scout for a date, commit the updated `data/` and `reports/` artifacts, build `public/`, validate it, and deliver the result through the Kit-managed issue/branch/PR workflow.

## Static GitHub Pages Site

Scout publishes as its own project site instead of being copied into `jamesonstone.github.io`. The personal site can link to the Scout project URL, but generated Scout artifacts stay in this repository.

The static site build is deterministic and backend-free:

- `go run ./cmd/scout site build --data-dir . --out-dir public --base-path /scout/`
- `go run ./cmd/scout site validate --out-dir public --base-path /scout/`

Generated output includes:

- `public/index.html` — homepage.
- `public/daily/index.html` and `public/daily/YYYY-MM-DD/index.html` — daily archive and briefings.
- `public/monthly/index.html` and `public/monthly/YYYY-MM/index.html` — monthly rankings.
- `public/papers/<paper-id>/index.html` — paper detail pages.
- `public/data/**` — copied JSON records and a generated manifest.
- `public/assets/styles.css` and `public/.nojekyll` — static assets for GitHub Pages.

`.github/workflows/pages.yml` builds and validates `public/`, then deploys it with GitHub Actions Pages. Repository settings should use GitHub Actions as the Pages source. GitHub Pages only serves static files; it does not fetch Hugging Face, execute Go, or mutate Scout storage at page view time.

## Repository Notes

Current repository inspection found no pre-existing implementation for scheduling, automation, report generation, storage, HTTP clients, markdown rendering, prompt management, or configuration. Scout introduces those subsystems under:

- `cmd/scout/main.go`
- `pkg/cli/root.go`
- `pkg/cli/run.go`
- `internal/config/config.go`
- `internal/hf/client.go`
- `internal/pipeline/*.go`
- `internal/scoring/scorer.go`
- `internal/summary/summary.go`
- `internal/pipeline/render.go`
- `internal/storage/store.go`
- `internal/prompt/*.tmpl`

## License

MIT.

## Maintainer

❤️ Lovingly overthought by [Jameson Stone](https://github.com/jamesonstone)
🪖 Field notes at [jamesonstone.io](https://jamesonstone.io)
