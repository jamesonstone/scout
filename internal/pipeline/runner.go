package pipeline

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/jamesonstone/scout/internal/config"
	"github.com/jamesonstone/scout/internal/scoring"
	"github.com/jamesonstone/scout/internal/storage"
	"github.com/jamesonstone/scout/internal/summary"
)

type client interface {
	FetchDailyPapers(context.Context, time.Time) ([]Paper, error)
	FetchPaperDetails(context.Context, string) (Paper, error)
	FetchMarkdown(context.Context, string) (string, error)
}

type Runner struct {
	cfg     config.Config
	client  client
	store   storage.Store
	scorer  scoring.Scorer
	summary summary.Builder
}

func NewRunner(cfg config.Config, client client, store storage.Store) Runner {
	return Runner{cfg: cfg, client: client, store: store, scorer: scoring.New(), summary: summary.New()}
}

func (r Runner) Run(ctx context.Context) (RunResult, error) {
	day, err := r.cfg.ResolveRunDate(time.Now().UTC())
	if err != nil {
		return RunResult{}, err
	}
	papers, err := r.client.FetchDailyPapers(ctx, day)
	if err != nil {
		return RunResult{}, fmt.Errorf("fetch daily papers: %w", err)
	}
	papers = dedupePapers(papers)
	observedIDs := make([]string, 0, len(papers))
	records := make([]PaperRecord, 0, len(papers))
	result := RunResult{}
	for _, candidate := range papers {
		observedIDs = append(observedIDs, candidate.ID)
		record, existed, err := r.store.LoadPaper(candidate.ID)
		if err != nil {
			return RunResult{}, fmt.Errorf("load paper %s: %w", candidate.ID, err)
		}
		if existed {
			record = mergeObservedDate(record, day)
			if err := r.store.SavePaper(record); err != nil {
				return RunResult{}, err
			}
			records = append(records, record)
			result.ReusedCount++
			continue
		}
		details, err := r.client.FetchPaperDetails(ctx, candidate.ID)
		if err != nil {
			details = candidate
		}
		paper := mergePaper(candidate, details)
		markdown, err := r.client.FetchMarkdown(ctx, paper.ID)
		if err == nil {
			paper.Markdown = markdown
		}
		score := r.scorer.Score(paper)
		recommendation := scoring.Recommendation(score.Overall)
		record = PaperRecord{
			ID:                   paper.ID,
			Title:                paper.Title,
			Authors:              paper.Authors,
			Categories:           stableCategories(paper.Categories),
			PublishedAt:          formatDate(paper.PublishedAt),
			FirstSeen:            day.Format("2006-01-02"),
			ObservedDates:        []string{day.Format("2006-01-02")},
			Abstract:             paper.Abstract,
			Markdown:             paper.Markdown,
			InnovationSummary:    r.summary.InnovationSummary(paper),
			WhyItMatters:         r.summary.WhyItMatters(paper, score),
			ImplementationAngle:  r.summary.ImplementationAngle(paper, score),
			Caveat:               r.summary.Caveat(paper),
			ExecutiveSummary:     r.summary.ExecutiveSummary(paper, score),
			EstimatedPriority:    r.summary.ReadingPriority(score.Overall),
			Recommendation:       recommendation,
			Links:                paper.Links,
			Score:                score,
			ScoreHistory:         []ScoreSnapshot{{ObservedAt: day.Format("2006-01-02"), Score: score}},
			Community:            paper.RawCommunity,
			MetadataCompleteness: metadataCompleteness(paper),
		}
		if err := r.store.SavePaper(record); err != nil {
			return RunResult{}, fmt.Errorf("save paper %s: %w", paper.ID, err)
		}
		records = append(records, record)
		result.ProcessedCount++
	}
	if err := r.store.SaveObservation(day, observedIDs); err != nil {
		return RunResult{}, err
	}
	monthRecords, err := r.store.MonthRecords(day)
	if err != nil {
		return RunResult{}, err
	}
	monthlyBody := RenderMonthly(day, monthRecords)
	monthlyPath, err := r.store.SaveMonthlyReport(day, monthlyBody)
	if err != nil {
		return RunResult{}, err
	}
	dailyBody := RenderDaily(day, records, r.summary.ExecutiveSignal(records, day.Format("2006-01-02")), monthlyPath)
	dailyPath, err := r.store.SaveDailyReport(day, dailyBody)
	if err != nil {
		return RunResult{}, err
	}
	result.DailyReportPath = dailyPath
	result.MonthlyReportPath = monthlyPath
	return result, nil
}

func dedupePapers(papers []Paper) []Paper {
	seen := map[string]Paper{}
	for _, paper := range papers {
		if paper.ID == "" {
			continue
		}
		if existing, ok := seen[paper.ID]; ok {
			seen[paper.ID] = mergePaper(existing, paper)
			continue
		}
		seen[paper.ID] = paper
	}
	ids := make([]string, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	out := make([]Paper, 0, len(ids))
	for _, id := range ids {
		out = append(out, seen[id])
	}
	return out
}

func mergePaper(base, override Paper) Paper {
	paper := base
	if override.ID != "" {
		paper.ID = override.ID
	}
	if override.Title != "" {
		paper.Title = override.Title
	}
	if override.Abstract != "" {
		paper.Abstract = override.Abstract
	}
	if len(override.Authors) > 0 {
		paper.Authors = override.Authors
	}
	if len(override.Categories) > 0 {
		paper.Categories = override.Categories
	}
	if !override.PublishedAt.IsZero() {
		paper.PublishedAt = override.PublishedAt
	}
	paper.Links = mergeLinks(paper.Links, override.Links)
	if override.Upvotes > 0 {
		paper.Upvotes = override.Upvotes
	}
	if override.Comments > 0 {
		paper.Comments = override.Comments
	}
	if override.Discussion > 0 {
		paper.Discussion = override.Discussion
	}
	if override.Markdown != "" {
		paper.Markdown = override.Markdown
	}
	if paper.RawCommunity == nil {
		paper.RawCommunity = map[string]int{}
	}
	for key, value := range override.RawCommunity {
		if value > 0 {
			paper.RawCommunity[key] = value
		}
	}
	return paper
}

func mergeLinks(base, override Links) Links {
	if override.HuggingFace != "" {
		base.HuggingFace = override.HuggingFace
	}
	if override.Arxiv != "" {
		base.Arxiv = override.Arxiv
	}
	if override.Paper != "" {
		base.Paper = override.Paper
	}
	if override.PDF != "" {
		base.PDF = override.PDF
	}
	if len(override.GitHub) > 0 {
		base.GitHub = append([]string(nil), override.GitHub...)
	}
	if len(override.Project) > 0 {
		base.Project = append([]string(nil), override.Project...)
	}
	return base
}

func mergeObservedDate(record PaperRecord, day time.Time) PaperRecord {
	date := day.Format("2006-01-02")
	for _, existing := range record.ObservedDates {
		if existing == date {
			return record
		}
	}
	record.ObservedDates = append(record.ObservedDates, date)
	sort.Strings(record.ObservedDates)
	record.ScoreHistory = append(record.ScoreHistory, ScoreSnapshot{ObservedAt: date, Score: record.Score})
	return record
}

func stableCategories(categories []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(categories))
	for _, category := range categories {
		if category == "" {
			continue
		}
		if _, ok := seen[category]; ok {
			continue
		}
		seen[category] = struct{}{}
		out = append(out, category)
	}
	sort.Strings(out)
	return out
}

func formatDate(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format(time.RFC3339)
}

func metadataCompleteness(paper Paper) int {
	completeness := 0
	if paper.Title != "" {
		completeness += 1
	}
	if paper.Abstract != "" {
		completeness += 1
	}
	if len(paper.Authors) > 0 {
		completeness += 1
	}
	if len(paper.Categories) > 0 {
		completeness += 1
	}
	if paper.Links.HuggingFace != "" {
		completeness += 1
	}
	if paper.Links.Arxiv != "" {
		completeness += 1
	}
	if paper.Markdown != "" {
		completeness += 1
	}
	return completeness
}
