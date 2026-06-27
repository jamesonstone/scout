package site

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jamesonstone/scout/internal/model"
)

func TestBuildAndValidateStaticSite(t *testing.T) {
	root := t.TempDir()
	writeFixturePaper(t, root, model.PaperRecord{
		ID:                  "2606.00001",
		Title:               "Agent Evaluation Paper",
		Categories:          []string{"agents", "evaluation"},
		FirstSeen:           "2026-06-26",
		PublishedDate:       "2026-06-25",
		ObservedDates:       []string{"2026-06-26"},
		InnovationSummary:   "Agent Evaluation Paper: We introduce a deterministic agent evaluation benchmark.",
		WhyItMatters:        []string{"It improves production evaluation signal."},
		ImplementationAngle: []string{"Use the benchmark as a regression suite."},
		Caveat:              "Needs validation on internal workloads.",
		Recommendation:      model.RecommendationRead,
		Links:               model.Links{HuggingFace: "https://huggingface.co/papers/2606.00001", Arxiv: "https://arxiv.org/abs/2606.00001"},
		Score:               model.ScoreBreakdown{Overall: 88, Novelty: 90, PracticalImpact: 85, TechnicalDepth: 80, ImplementationPotential: 86, Relevance: 92, CommunitySignal: 70, SummaryConfidence: 95},
	})
	writeFixtureJSON(t, filepath.Join(root, "data", "daily", "2026-06", "2026-06-26.json"), model.DailyObservation{Date: "2026-06-26", PaperIDs: []string{"2606.00001"}})
	writeText(t, filepath.Join(root, "reports", "daily", "2026-06", "2026-06-26.md"), "# Scout Daily Intelligence Briefing\n\n## Executive Signal\n\nAgents are the dominant signal.\n\n## Top Papers\n")

	outDir := filepath.Join(root, "public")
	result, err := Build(Config{DataDir: root, OutDir: outDir, BasePath: "/scout/"})
	if err != nil {
		t.Fatalf("build site: %v", err)
	}
	if result.DailyPages != 1 || result.MonthlyPages != 1 || result.PaperPages != 1 {
		t.Fatalf("unexpected result: %#v", result)
	}
	home, err := os.ReadFile(filepath.Join(outDir, "index.html"))
	if err != nil {
		t.Fatalf("read home page: %v", err)
	}
	for _, token := range []string{"<h1>Scout</h1>", "Papers fetched by Scout", "Papers fetched on 2026-06-26", "Published 2026-06-25"} {
		if !strings.Contains(string(home), token) {
			t.Fatalf("home page missing %q", token)
		}
	}
	body, err := os.ReadFile(filepath.Join(outDir, "daily", "2026-06-26", "index.html"))
	if err != nil {
		t.Fatalf("read daily page: %v", err)
	}
	for _, token := range []string{"Papers fetched on 2026-06-26", "Executive Signal", "Top Papers", "Innovation Summary", "Executive Summary", "Why It Matters", "Estimated Reading Priority", "Published 2026-06-25", "/scout/papers/2606.00001/"} {
		if !strings.Contains(string(body), token) {
			t.Fatalf("daily page missing %q", token)
		}
	}
	validation, err := Validate(Config{OutDir: outDir, BasePath: "/scout/"})
	if err != nil {
		t.Fatalf("validate site: %v", err)
	}
	if validation.CheckedPages == 0 || validation.CheckedLinks == 0 {
		t.Fatalf("unexpected validation result: %#v", validation)
	}
}

func writeFixturePaper(t *testing.T, root string, record model.PaperRecord) {
	t.Helper()
	writeFixtureJSON(t, filepath.Join(root, "data", "papers", record.ID+".json"), record)
}

func writeFixtureJSON(t *testing.T, path string, value any) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := writeJSON(path, value); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func writeText(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
