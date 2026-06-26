package hf

import "testing"

func TestParsePaperUnwrapsDailyPaperPayload(t *testing.T) {
	paper := parsePaper(map[string]any{
		"numComments": float64(3),
		"paper": map[string]any{
			"id":          "2606.27377",
			"title":       "DanceOPD: On-Policy Generative Field Distillation",
			"summary":     "A flow-matching distillation framework for composing generative capabilities.",
			"publishedAt": "2026-06-25T00:00:00.000Z",
			"authors": []any{
				map[string]any{"name": "Wei Zhou"},
				map[string]any{"name": "Xiongwei Zhu"},
			},
			"ai_keywords": []any{"flow-matching models", "on-policy"},
			"projectPage": "https://danceopd.github.io/",
			"upvotes":     float64(51),
		},
	})

	if paper.ID != "2606.27377" {
		t.Fatalf("unexpected id: %q", paper.ID)
	}
	if paper.Title != "DanceOPD: On-Policy Generative Field Distillation" {
		t.Fatalf("unexpected title: %q", paper.Title)
	}
	if paper.Abstract == "" {
		t.Fatal("expected summary to populate abstract")
	}
	if len(paper.Authors) != 2 || paper.Authors[0] != "Wei Zhou" {
		t.Fatalf("unexpected authors: %#v", paper.Authors)
	}
	if len(paper.Categories) != 2 || paper.Categories[0] != "flow-matching models" {
		t.Fatalf("unexpected categories: %#v", paper.Categories)
	}
	if len(paper.Links.Project) != 1 || paper.Links.Project[0] != "https://danceopd.github.io/" {
		t.Fatalf("unexpected project links: %#v", paper.Links.Project)
	}
	if paper.Links.HuggingFace != "https://huggingface.co/papers/2606.27377" {
		t.Fatalf("unexpected Hugging Face link: %q", paper.Links.HuggingFace)
	}
	if paper.Links.Arxiv != "https://arxiv.org/abs/2606.27377" {
		t.Fatalf("unexpected arXiv link: %q", paper.Links.Arxiv)
	}
	if paper.Upvotes != 51 || paper.Comments != 3 {
		t.Fatalf("unexpected community signal: upvotes=%d comments=%d", paper.Upvotes, paper.Comments)
	}
	if paper.PublishedAt.IsZero() {
		t.Fatal("expected publishedAt to parse")
	}
}
