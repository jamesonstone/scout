# ScoutResearchDistillation Scheduled Prompt

Run the weekday Scout automation for the `scout` repository and deliver the result through the repository's Kit-managed GitHub workflow.

Repository: `/Users/jamesonstone/go/src/github.com/jamesonstone/scout`

## Scheduled Task Contract

This automation is intended to run Monday through Friday only.

The scheduled task should use this cron RRULE:

```text
RRULE:FREQ=WEEKLY;BYHOUR=8;BYMINUTE=0;BYDAY=MO,TU,WE,TH,FR
```

Use the local timezone configured for the scheduled task.

Before doing any repository, file, git, or GitHub mutation, determine the current local date and day of week.

If the current local day is Saturday or Sunday:

1. Stop before any checkout, pull, branch, issue, artifact, commit, push, PR, or merge mutation.
2. Report that Scout is weekday-only and no run was performed.
3. Do not create an issue, branch, commit, pull request, or artifact update.

If the current local day is Monday through Friday, continue with the run.

## Purpose

Scout is a scheduled research intelligence system. It fetches Hugging Face Daily Papers, detects new papers, scores them, summarizes the highest-signal papers, writes a concise Markdown research brief, persists paper records, updates the current month's Top 10 papers, publishes static site artifacts, opens a pull request, and self-merges that pull request when the required merge gate passes.

Primary goal:

Generate and deliver an AI research briefing that is more useful than the default Hugging Face Daily Papers email. The report must not be a title-only digest. Each selected paper must explain the paper's core innovation, why it matters, implementation implications, and whether it deserves deeper reading.

Optimize for:

- skimmability
- signal density
- deterministic output
- historical tracking
- engineering usefulness
- low cognitive load
- minimal, production-ready changes
- autonomous weekday delivery

## Transparency Rule

Explain decisions, assumptions, failure modes, and tradeoffs briefly. Do not expose hidden chain-of-thought. Provide concise rationale and observable evidence only.

## Authoritative Source

Use Hugging Face Daily Papers as the source of truth.

Known resources:

- Daily Papers page: `https://huggingface.co/papers`
- OpenAPI definition: `https://huggingface.co/.well-known/openapi.md`
- Hub API docs: `https://huggingface.co/docs/hub/api`

Prefer official API endpoints over scraping. Inspect the OpenAPI definition or existing project code before changing fetch behavior.

## Repo-Local Instructions

1. Start by reading `AGENTS.md`.
2. Then read `docs/agents/README.md`.
3. Load only the additional docs needed for the decision at hand.
4. Treat repo-local docs under `docs/` as the system of record.
5. Reuse existing runtime, storage, scoring, report rendering, CLI, test, lint, formatting, and site conventions.
6. Do not invent new paths when equivalent project paths already exist.
7. Repo-local Kit rules control issue, branch, commit, PR, and merge shape.
8. This scheduled prompt adds the automation-specific rule that required checks are the merge gate. Optional checks such as CodeRabbit review-in-progress are reported but must not block self-merge unless GitHub marks them required.

## GitHub Delivery Requirement

This run must use the repository's Kit-managed delivery conventions for any issue, branch, commit, push, pull request, or merge work.

In this automation, "merge readiness gate" means:

- required GitHub checks pass, or
- GitHub reports no required checks for the branch.

Optional advisory checks, including CodeRabbit review-in-progress, must be reported in the final response but are not merge blockers unless they appear in the required-check set.

Use `gh pr checks --required` for the merge gate.

Use plain `gh pr checks` as an observation/reporting command only.

## Startup Main-Sync Requirement

Before creating or mutating any issue, issue branch, commit, push, pull request, merge state, or Scout artifact, first return the local checkout to the latest base branch:

1. Load the repo-local delivery docs and run the initial delivery recon below.
2. Discover the default/base branch with:
   - `gh repo view --json defaultBranchRef -q .defaultBranchRef.name`
3. If the working tree has uncommitted changes before syncing, stop and report them. Do not stash, reset, clean, or overwrite local work.
4. If the current branch is not the base branch, check out the base branch:
   - `git checkout $BASE_BRANCH`
5. Pull the latest base branch with fast-forward only:
   - `git pull --ff-only origin $BASE_BRANCH`
6. Verify the local base branch is current:
   - `git status --short --branch`
   - `git rev-parse HEAD`
   - `git rev-parse origin/$BASE_BRANCH`
7. If checkout, pull, or verification fails, stop and report the exact state. Do not continue from a stale local base branch.
8. Run Scout, create or reuse the issue, create the `GH-123` branch, commit, push, create the PR, and self-merge the PR on top of this latest synced base when the required merge gate passes.

## Required Pre-Mutation Context

Before any GitHub or git mutation, load:

- `.kit.yaml`
- `docs/agents/GUARDRAILS.md`
- `docs/agents/TOOLING.md`
- `docs/references/rules/safety-guardrails.md`
- `docs/references/rules/work-lane-gating.md`
- `docs/references/rules/github-pr-delivery.md`
- `.github/pull_request_template.md`

Then run and show delivery recon:

- `pwd`
- `git status --short --branch`
- `git remote -v`
- `git rev-parse --abbrev-ref HEAD`
- default branch via `gh repo view --json defaultBranchRef -q .defaultBranchRef.name`
- active PRs for the current branch
- git author and committer identity
- `gh auth status`

Confirm the authenticated GitHub user is the human user, Jameson Stone. If the GitHub login cannot be confirmed as Jameson, stop and ask.

## Delivery Contract

Before creating or mutating issue, branch, staging, commit, push, PR, or merge state, present this Delivery Contract and stop if any field is unknown:

```text
Delivery Contract:
- Repository:
- Base branch:
- Base branch local sync:
- Issue source:
- Issue number/link:
- Issue assignee:
- Branch name:
- Branch base:
- Branch/status/staleness check:
- Staging method:
- Commit format:
- PR title format:
- PR template:
- Draft or ready:
- PR assignee:
- Merge method:
- Merge readiness gate:
- Required checks:
- Optional checks:
- Cross-repo dependencies:
- Unknowns/blockers:
```

The merge readiness gate must resolve to `gh pr checks --required`: merge only if all required checks pass or there are no required checks. Do not use optional CodeRabbit pending status as a blocker unless it is required.

## Issue Rules

1. Search for an existing open issue with a title matching this weekday Scout run before creating a new one.
2. Use a same-day title such as `ScoutResearchDistillation YYYY-MM-DD`.
3. If exactly one strong matching issue exists, reuse it.
4. If multiple plausible matches exist, stop and ask which issue to use.
5. If no matching issue exists, create one assigned to the human user with `--assignee @me` or the confirmed GitHub login.
6. If reusing an existing issue, confirm it is assigned to the human user. If not, add the human user as assignee before branching.
7. Do not assign any bot, agent, assistant, or tool account to the issue.

## Branch Rules

- Branch name must be exactly the GitHub issue number in `GH-123` form.
- Do not use `codex/*`.
- Do not add a slug or suffix.
- Start every weekday run by checking out the base branch and pulling the latest changes with `git pull --ff-only origin $BASE_BRANCH`.
- Fetch the remote base again before branching.
- Create the branch from `origin/$BASE_BRANCH`, not the local base branch.
- Do not commit directly to `main`, `master`, or any protected or assumed-protected branch.
- Do not create or use a git worktree.

## Implementation Behavior

1. Inspect the repository.
2. Identify the package manager, runtime, source layout, config conventions, test runner, lint command, formatter command, HTTP client, storage conventions, report-generation patterns, site-generation patterns, and automation entrypoints.
3. Reuse existing conventions.
4. Record exact file paths and symbols used.
5. Run the existing Scout automation before changing code unless the repository cannot currently run.
6. If the existing implementation already satisfies the weekday automation contract, do not make unnecessary code changes.
7. If it does not satisfy the contract, make the smallest production-ready change needed.

## Weekday Run Behavior

Use the current local date for `YYYY-MM-DD` and `YYYY-MM`.

Run Scout with an explicit data directory appropriate to delivery:

- If the weekday artifacts should be committed in the PR, use `--data-dir .` so outputs land under tracked repo-root `data/`, `reports/`, and derived `public/` paths.
- If repo-local docs or gitignore indicate generated artifacts should remain local, use `.scout` and do not force an empty PR. Report that no PR was created unless source/docs/config changes were required.
- Never stage `.scout/`, `.env`, `.envrc`, `bin/`, `.kit/runs/`, `.kit/loops/`, `.kit/cache/`, `.kit/tmp/`, or other local machine artifacts.

Expected Scout artifact layout follows the live code and configured data directory:

- `data/papers/<paper-id>.json`
- `data/daily/YYYY-MM/YYYY-MM-DD.json`
- `reports/daily/YYYY-MM/YYYY-MM-DD.md`
- `reports/monthly/YYYY-MM.md`
- `public/**` generated site files

Do not use stale prompt paths if they conflict with the current repo implementation. Prefer live repo behavior unless intentionally changing it.

## Paper Processing Requirements

1. Fetch today's Hugging Face Daily Papers.
2. Persist durable daily observation data.
3. Normalize each paper into stable internal records.
4. Deduplicate by paper ID, with normalized arXiv ID as a secondary key if supported.
5. Do not reprocess an already summarized paper.
6. Update observed dates or last-seen equivalent for previously known papers.
7. Preserve existing score and summary history unless source content changed.
8. Fetch paper details and Markdown when available.
9. If full text is unavailable, summarize from metadata and abstract only and lower summary confidence.
10. Continue processing other papers when one paper has partial metadata.

If Hugging Face returns an empty same-date response, treat it as an authoritative zero-paper weekday. Still publish durable zero-paper artifacts and self-merge them after the required merge gate passes.

## Scoring Requirements

Score each paper from 0 to 100 using the repository's existing scoring model unless explicitly fixing that model.

The intended rubric is:

- novelty
- practical impact
- technical depth
- relevance to AI agents, LLMs, memory, orchestration, tool use, evaluation, MoE, inference, training, ML systems, or software engineering
- implementation potential
- community signal
- summary confidence

Recommendations:

- Read: highest-signal papers, generally score >= 80
- Worth Watching: useful but not urgent
- Archive or Skip: lower signal or low confidence

If the live repo uses slightly different thresholds or names, preserve live behavior unless changing it is necessary and justified.

## Report Quality Rules

- Do not produce a title-only digest.
- Do not summarize every paper in depth.
- Do not bury the best papers below low-value items.
- Do not use generic filler such as "this paper is important because it advances the field."
- Do not overstate weak papers.
- Do not include a paper in Top Papers unless the summary explains a concrete innovation.
- Prefer direct engineering implications over broad academic claims.
- Prefer fewer, better summaries over comprehensive coverage.
- Make the report useful even when the reader only scans headings, scores, and innovation summaries.
- On zero-paper weekdays, still publish a concise report that clearly states no papers were available and includes the standard sections.

Daily report should include, at minimum:

- Executive Signal
- Top Papers
- Additional Papers
- Watchlist
- Archive

Each full paper summary should include:

- rank
- title
- final score
- recommendation
- categories
- one-sentence innovation summary
- why it matters
- implementation angle
- one-sentence caveat
- executive summary under 300 words
- links
- estimated reading priority

Monthly report should:

- recompute rankings from current month paper records
- avoid stale ordering from previous daily reports
- include Top 10 papers, rising papers, themes, and a full monthly index or the closest equivalent supported by the current renderer

## Failure Handling

- If Hugging Face is unavailable, use same-date cached data if present.
- If no same-date cache exists, use the most recent cached response only if the report clearly marks the source date.
- If paper Markdown cannot be fetched, continue with metadata and abstract.
- If scoring cannot be completed for a paper, put it in Watchlist or the closest supported low-confidence section with a reason.
- If report generation fails, do not overwrite the previous valid report.
- If monthly reranking fails, preserve the previous monthly report and report the error.
- On git, GitHub, auth, push, PR creation, merge, lint, test, or state mismatch failure, stop and report exact state.
- Do not retry with mutation after a failure.
- Do not force-push, rebase, reset, clean, bypass branch protection, use admin merge, enable auto-merge, or delete branches unless explicitly requested.

## Validation

Run the smallest complete validation set for this repo:

1. Scout weekday run:
   - `go run ./cmd/scout run --date YYYY-MM-DD --data-dir .`
2. Static site build:
   - `go run ./cmd/scout site build --data-dir . --out-dir public --base-path /`
3. Static site validation:
   - `go run ./cmd/scout site validate --out-dir public --base-path /`
4. Formatter:
   - `make fmt`
5. Linter:
   - `make lint`
6. Tests:
   - `make test`
7. Build:
   - `make build`
8. Diff checks:
   - `git diff --check`
   - `git diff --staged --check` after staging
9. Artifact and safety checks:
   - daily observation file exists
   - daily report exists
   - monthly report exists
   - public daily page exists
   - public daily JSON exists
   - at least one paper record exists in the repository unless this is the first ever run
   - daily report contains Executive Signal, Top Papers, Additional Papers, Watchlist, and Archive
   - full summaries, when present, contain Innovation Summary, Why It Matters, Implementation Angle, Caveat, Executive Summary under 300 words, score, and recommendation
   - underlying JSON records include score breakdowns or the repo's equivalent score object
   - focused secret scan of staged artifact paths

If `golangci-lint` or another local tool is unavailable, report that specifically and continue only with checks that can run. Do not claim a skipped or unavailable check passed.

## Staging And Commit Rules

1. Self-review the diff before staging.
2. Check for secrets and local-only files before staging.
3. Stage files explicitly with `git add <file>`.
4. Never use `git add .` or `git add -A`.
5. Review `git diff --staged`.
6. Confirm git author and committer are the human user.
7. Commit title must be:
   - `<type>(GH-123): <gitmoji> <short title message>`
8. Use deterministic type/gitmoji mapping:
   - `feat` -> `:sparkles:`
   - `fix` -> `:bug:`
   - `docs` -> `:memo:`
   - `test` -> `:white_check_mark:`
   - `refactor` -> `:recycle:`
   - `chore` -> `:wrench:`
   - `ci` -> `:green_heart:`
9. Commit body must include:
   - Original ask
   - Implementation summary
   - Verification performed
   - Reference to the GitHub issue
10. Do not add agent attribution trailers.

For routine Scout artifact publication, use:

```text
chore(GH-123): :wrench: publish Scout briefing
```

## Push, PR, And Self-Merge Rules

1. Before pushing, rerun branch, remote, and active PR recon.
2. Push only the `GH-123` branch.
3. If a PR already exists for the branch, update it. Do not create a duplicate.
4. Create the PR ready for review, not draft, unless explicitly asked otherwise.
5. Use `.github/pull_request_template.md` exactly:
   - `## Description`
   - `## How to Test`
   - `## Ticket`
6. Preserve template headings.
7. PR title must match the commit title format:
   - `<type>(GH-123): <gitmoji> <short title message>`
8. Prefix concrete bullets in the Description section with descriptive gitmoji codes.
9. Use `Closes #123` if the PR fully resolves the issue. Otherwise use `Refs #123`.
10. Assign the PR to the human user.
11. Confirm the GitHub issue and PR both show the human user as assignee.
12. Run:
    - `gh issue view 123 --json number,url,assignees,state`
    - `gh pr view 123 --json number,url,author,state,isDraft,assignees,headRefName,baseRefName`
    - `gh pr checks 123 --required`
    - `gh pr checks 123`
13. Do not claim CI passed unless required checks were observed passing or GitHub reported no required checks.
14. Merge the pull request after creation only after confirming:
    - the PR is open and ready for review, not draft
    - the PR head is the exact `GH-123` branch
    - the PR base is the resolved base branch
    - the GitHub issue and PR are assigned to the human user
    - `gh pr checks 123 --required` reports all required checks passing, or reports no required checks for the branch
15. If required checks are pending, failing, unavailable due to auth or network failure, or ambiguous, stop and report the exact required-check state.
16. If plain `gh pr checks 123` shows optional CodeRabbit pending, report it but do not block merge unless `gh pr checks 123 --required` includes CodeRabbit as required.
17. Merge with a normal merge commit using GitHub CLI, preserving the PR title as the merge subject:
    - `gh pr merge 123 --merge --subject "<PR title>"`
18. Do not squash, rebase, use admin bypass, enable auto-merge, force-push, or delete the branch unless explicitly requested.
19. After merging, verify:
    - `gh pr view 123 --json number,url,state,mergedAt,mergeCommit,assignees,headRefName,baseRefName`
    - `gh issue view 123 --json number,url,state,assignees`
    - `git status --short --branch`
20. After a successful merge, return the local checkout to the updated base branch:
    - `git checkout $BASE_BRANCH`
    - `git pull --ff-only origin $BASE_BRANCH`
    - `git status --short --branch`
21. If post-merge checkout or pull fails, report the exact state. Do not reset, clean, force-pull, or delete branches unless explicitly requested.

## Final Response Requirements

Return:

1. Weekday guard result, including local date and day of week.
2. Issue number and URL.
3. Issue assignee.
4. Branch name.
5. Commit hash.
6. PR URL.
7. PR assignee.
8. Daily report path.
9. Monthly report path.
10. Daily observation or raw data path.
11. Public daily page path.
12. Public daily JSON path.
13. Number of papers fetched.
14. Number of new papers processed.
15. Number of duplicate or reused papers.
16. Top 5 paper titles with scores, or `none - zero-paper day`.
17. Validation commands executed.
18. Validation results, including exact skipped or unavailable checks.
19. Required-check state from `gh pr checks --required`.
20. Optional check state from plain `gh pr checks`.
21. Merge state and merge commit hash, if available.
22. Observed GitHub issue state after merge.
23. Initial base-branch sync state before the run.
24. Final local branch and sync state after merge.
25. Files changed.
26. Failures encountered.
27. Assumptions made.
28. Open questions.
