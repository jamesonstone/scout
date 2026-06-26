package scoring

import (
	"testing"

	"github.com/jamesonstone/scout/internal/model"
)

func TestScoreDeterministicAndBounded(t *testing.T) {
	paper := model.Paper{
		ID:         "2501.00001",
		Title:      "Agent Memory Benchmark",
		Abstract:   "We introduce a novel benchmark for AI agents with strong evaluation, infrastructure, and open source implementation guidance.",
		Categories: []string{"agents", "evaluation"},
		Links:      model.Links{GitHub: []string{"https://github.com/example/repo"}},
		Upvotes:    7,
		Comments:   3,
		Markdown:   "The paper studies agent memory systems, orchestration tradeoffs, evaluation details, and implementation guidance.",
	}
	scorer := New()
	first := scorer.Score(paper)
	second := scorer.Score(paper)
	if first != second {
		t.Fatalf("expected deterministic scores, got %#v and %#v", first, second)
	}
	if first.Overall < 0 || first.Overall > 100 {
		t.Fatalf("overall score out of range: %d", first.Overall)
	}
	if Recommendation(first.Overall) == "" {
		t.Fatal("expected non-empty recommendation")
	}
}
