package pipeline

import (
	"strings"
	"testing"
	"time"
)

func TestRenderDailyIncludesRequiredSections(t *testing.T) {
	day, _ := time.Parse("2006-01-02", "2026-01-02")
	paper := PaperRecord{
		Title:               "Scout Paper",
		Categories:          []string{"agents"},
		InnovationSummary:   "Scout Paper introduces deterministic scoring for research intelligence.",
		WhyItMatters:        []string{"High signal."},
		ImplementationAngle: []string{"Easy to integrate."},
		Caveat:              "Needs real-world validation.",
		ExecutiveSummary:    "A short executive summary.",
		Recommendation:      RecommendationRead,
		EstimatedPriority:   "Immediate",
		Score:               ScoreBreakdown{Overall: 88},
		Links:               Links{HuggingFace: "hf", Arxiv: "arxiv", GitHub: []string{"gh"}, Paper: "paper"},
	}
	body := RenderDaily(day, []PaperRecord{paper}, "Themes are converging around agents.", "/tmp/monthly.md")
	for _, token := range []string{"## Executive Signal", "## Top Papers", "## Additional Papers", "## Watchlist", "## Archive", "**Recommendation:** Read", "**Innovation Summary:** Scout Paper introduces deterministic scoring for research intelligence."} {
		if !strings.Contains(body, token) {
			t.Fatalf("expected %q in report", token)
		}
	}
}
