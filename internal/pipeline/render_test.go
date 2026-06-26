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
		Recommendation:      RecommendationRead,
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

func TestRenderDailyUsesNALabelForEmptyCategories(t *testing.T) {
	day, _ := time.Parse("2006-01-02", "2026-01-02")
	paper := PaperRecord{
		Title:               "Uncategorized Paper",
		InnovationSummary:   "Uncategorized Paper introduces a useful result.",
		WhyItMatters:        []string{"High signal."},
		ImplementationAngle: []string{"Easy to integrate."},
		Caveat:              "Needs validation.",
		Recommendation:      RecommendationRead,
		Score:               ScoreBreakdown{Overall: 88},
	}
	body := RenderDaily(day, []PaperRecord{paper}, "Themes are sparse.", "/tmp/monthly.md")
	if !strings.Contains(body, "- **Categories:** N/A") {
		t.Fatalf("expected N/A category label in report")
	}
}
