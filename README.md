```text
███████╗ ██████╗ ██████╗ ██╗   ██╗████████╗
██╔════╝██╔════╝██╔═══██╗██║   ██║╚══██╔══╝
███████╗██║     ██║   ██║██║   ██║   ██║
╚════██║██║     ██║   ██║██║   ██║   ██║
███████║╚██████╗╚██████╔╝╚██████╔╝   ██║
╚══════╝ ╚═════╝ ╚═════╝  ╚═════╝    ╚═╝   daily papers intelligence distillation
```

**`scout` is a research-intelligence pipeline for engineers who want the highest-value AI papers without reading dozens of titles every day.**

Scout ingests the Hugging Face Daily Papers feed, enriches each paper with metadata and markdown, scores it across explicit weighted dimensions, persists durable paper records, and publishes executive-style daily and monthly intelligence briefings.

## Why Scout Exists

- 🧠 Replace title-only discovery with concise signal-dense paper intelligence.
- 📈 Maintain a continuously reranked monthly leaderboard instead of static daily ordering.
- 🗃️ Preserve historical paper records, score breakdowns, and observation history.
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
```

## Output Layout

Scout writes durable artifacts beneath the configured data directory:

- `data/papers/<paper-id>.json` — persistent per-paper records, summaries, links, and score history.
- `data/daily/YYYY-MM/YYYY-MM-DD.json` — observed paper IDs for each daily run.
- `reports/daily/YYYY-MM/YYYY-MM-DD.md` — daily executive briefing.
- `reports/monthly/YYYY-MM.md` — continuously reranked monthly briefing.

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
