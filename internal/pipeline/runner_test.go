package pipeline

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jamesonstone/scout/internal/config"
	"github.com/jamesonstone/scout/internal/storage"
)

type fakeClient struct {
	dailyCalls    int
	detailCalls   int
	markdownCalls int
}

func (f *fakeClient) FetchDailyPapers(_ context.Context, day time.Time) ([]Paper, error) {
	f.dailyCalls++
	return []Paper{
		{ID: "2501.00001", Title: "Agent System Paper", Abstract: "A novel agent orchestration framework with open source code.", Categories: []string{"agents", "orchestration"}, Links: Links{GitHub: []string{"https://github.com/example/agent"}, HuggingFace: "https://huggingface.co/papers/2501.00001"}, Upvotes: 4, Comments: 1, SourceDate: day},
		{ID: "2501.00001", Title: "Agent System Paper", Abstract: "duplicate", Categories: []string{"agents"}, SourceDate: day},
		{ID: "2501.00002", Title: "Infra Eval Paper", Abstract: "Evaluation for infrastructure and llm memory systems.", Categories: []string{"evaluation", "infrastructure"}, Links: Links{HuggingFace: "https://huggingface.co/papers/2501.00002"}, Upvotes: 2, SourceDate: day},
	}, nil
}

func (f *fakeClient) FetchPaperDetails(_ context.Context, id string) (Paper, error) {
	f.detailCalls++
	switch id {
	case "2501.00001":
		return Paper{ID: id, Title: "Agent System Paper", Abstract: "A novel agent orchestration framework with open source code and evaluation details.", Categories: []string{"agents", "orchestration"}, Links: Links{Arxiv: "https://arxiv.org/abs/2501.00001", HuggingFace: "https://huggingface.co/papers/2501.00001", GitHub: []string{"https://github.com/example/agent"}}}, nil
	case "2501.00002":
		return Paper{ID: id, Title: "Infra Eval Paper", Abstract: "Evaluation for infrastructure and llm memory systems.", Categories: []string{"evaluation", "infrastructure"}, Links: Links{Arxiv: "https://arxiv.org/abs/2501.00002", HuggingFace: "https://huggingface.co/papers/2501.00002"}}, nil
	default:
		return Paper{}, errors.New("unknown id")
	}
}

func (f *fakeClient) FetchMarkdown(_ context.Context, id string) (string, error) {
	f.markdownCalls++
	return "Markdown summary for " + id + ".", nil
}

func TestRunnerProducesDailyAndMonthlyReportsAndAvoidsReprocessing(t *testing.T) {
	root := t.TempDir()
	cfg := config.Config{DataDir: root, RunDate: "2026-01-02"}
	runner := NewRunner(cfg, &fakeClient{}, storage.New(root))
	result, err := runner.Run(context.Background())
	if err != nil {
		t.Fatalf("first run: %v", err)
	}
	if result.ProcessedCount != 2 || result.ReusedCount != 0 {
		t.Fatalf("unexpected counts on first run: %#v", result)
	}
	for _, path := range []string{result.DailyReportPath, result.MonthlyReportPath} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected output %s: %v", path, err)
		}
	}
	body, err := os.ReadFile(result.DailyReportPath)
	if err != nil {
		t.Fatalf("read daily report: %v", err)
	}
	text := string(body)
	for _, token := range []string{"## Executive Signal", "### 1.", "**Recommendation:**", "**Implementation Angle:**", "**Caveat:**", "**Estimated reading priority:**"} {
		if !strings.Contains(text, token) {
			t.Fatalf("missing %q in daily report", token)
		}
	}

	runner.cfg.RunDate = "2026-01-03"
	second, err := runner.Run(context.Background())
	if err != nil {
		t.Fatalf("second run: %v", err)
	}
	if second.ProcessedCount != 0 || second.ReusedCount != 2 {
		t.Fatalf("unexpected counts on second run: %#v", second)
	}
	paperPath := filepath.Join(root, "data", "papers", "2501.00001.json")
	data, err := os.ReadFile(paperPath)
	if err != nil {
		t.Fatalf("read paper record: %v", err)
	}
	if !strings.Contains(string(data), "2026-01-03") {
		t.Fatalf("expected observed date update in %s", paperPath)
	}
}
