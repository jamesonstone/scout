package summary

import (
	"strings"
	"testing"

	"github.com/jamesonstone/scout/internal/model"
)

func TestBuilderPrefersAbstractOverMarkdownHeader(t *testing.T) {
	builder := New()
	paper := model.Paper{
		Title:    "Header Heavy Paper",
		Abstract: "A concrete method improves coding agent verification with calibrated reward checks.",
		Markdown: "Title: Header Heavy Paper\n\nURL Source: https://arxiv.org/html/2606.00001\n\nMarkdown Content:\nBack to arXiv\nAbstract\n\nThe full paper repeats the same claim with more detail.",
	}

	summary := builder.InnovationSummary(paper)
	if strings.Contains(summary, "URL Source") || strings.Contains(summary, "Title:") {
		t.Fatalf("summary used markdown header: %q", summary)
	}
	if !strings.Contains(summary, "concrete method") {
		t.Fatalf("summary did not use abstract signal: %q", summary)
	}
}

func TestCleanMarkdownTextSkipsScrapedHeader(t *testing.T) {
	text := cleanMarkdownText("Title: Example\n\nURL Source: https://arxiv.org/html/1\n\nMarkdown Content:\nBack to arXiv\n###### Abstract\n\nThis is the useful abstract sentence.")
	if strings.HasPrefix(text, "Title:") || strings.Contains(text, "URL Source") {
		t.Fatalf("header was not removed: %q", text)
	}
	if !strings.HasPrefix(text, "This is the useful abstract sentence") {
		t.Fatalf("unexpected cleaned text: %q", text)
	}
}

func TestInnovationSentencePrefersContributionCue(t *testing.T) {
	text := "A classical intuition holds that verifying a solution is easier than producing one. To address this, we characterize verification signals along scalability, faithfulness, and robustness."
	got := innovationSentence(text)
	if !strings.HasPrefix(got, "To address this") {
		t.Fatalf("expected contribution sentence, got %q", got)
	}
}

func TestExecutiveSignalUsesTopRankedPaperTitles(t *testing.T) {
	builder := New()
	papers := []model.PaperRecord{
		{
			Title:      "Low Signal Agent Benchmark",
			Score:      model.ScoreBreakdown{Overall: 53},
			Categories: []string{"ALFWorld"},
		},
		{
			Title:      "SkillOpt-Lite: Better and Faster Agent Self-evolution via One Line of Vibe",
			Score:      model.ScoreBreakdown{Overall: 99},
			Categories: []string{"HarnessOpt"},
		},
		{
			Title:      "Gemma 4 Technical Report",
			Score:      model.ScoreBreakdown{Overall: 98},
			Categories: []string{"Mixture-of-Experts architectures"},
		},
		{
			Title:      "Hierarchical Sparse Attention Done Right",
			Score:      model.ScoreBreakdown{Overall: 98},
			Categories: []string{"attention mechanism"},
		},
	}

	signal := builder.ExecutiveSignal(papers, "2026-07-08")
	for _, want := range []string{"SkillOpt-Lite", "Gemma 4 Technical Report", "Hierarchical Sparse Attention"} {
		if !strings.Contains(signal, want) {
			t.Fatalf("signal %q missing top paper %q", signal, want)
		}
	}
	if strings.Contains(signal, "ALFWorld") {
		t.Fatalf("signal should not prefer low-ranked category ties: %q", signal)
	}
}
