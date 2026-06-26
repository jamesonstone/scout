package hf

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jamesonstone/scout/internal/config"
	"github.com/jamesonstone/scout/internal/model"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	retries    int
	retryWait  time.Duration
	userAgent  string
}

func NewClient(cfg config.Config) Client {
	return Client{
		baseURL:    strings.TrimRight(cfg.BaseURL, "/"),
		httpClient: &http.Client{Timeout: cfg.Timeout},
		retries:    cfg.Retries,
		retryWait:  cfg.RetryWait,
		userAgent:  cfg.UserAgent,
	}
}

func (c Client) FetchDailyPapers(ctx context.Context, day time.Time) ([]model.Paper, error) {
	endpoint := fmt.Sprintf("%s/api/daily_papers?date=%s", c.baseURL, day.Format("2006-01-02"))
	var payload any
	if err := c.getJSON(ctx, endpoint, &payload); err != nil {
		return nil, err
	}
	items := extractItems(payload)
	papers := make([]model.Paper, 0, len(items))
	for _, item := range items {
		paper := parsePaper(item)
		if paper.ID == "" || paper.Title == "" {
			continue
		}
		paper.SourceDate = day
		papers = append(papers, paper)
	}
	sort.Slice(papers, func(i, j int) bool { return papers[i].ID < papers[j].ID })
	return papers, nil
}

func (c Client) FetchPaperDetails(ctx context.Context, id string) (model.Paper, error) {
	endpoint := fmt.Sprintf("%s/api/papers/%s", c.baseURL, url.PathEscape(id))
	var payload any
	if err := c.getJSON(ctx, endpoint, &payload); err != nil {
		return model.Paper{}, err
	}
	paper := parsePaper(payload)
	if paper.ID == "" {
		paper.ID = id
	}
	return paper, nil
}

func (c Client) FetchMarkdown(ctx context.Context, id string) (string, error) {
	endpoint := fmt.Sprintf("%s/papers/%s.md", c.baseURL, path.Clean(id))
	body, err := c.getText(ctx, endpoint)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(body), nil
}

func (c Client) getJSON(ctx context.Context, endpoint string, target any) error {
	body, err := c.getText(ctx, endpoint)
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(body), target); err != nil {
		return fmt.Errorf("decode %s: %w", endpoint, err)
	}
	return nil
}

func (c Client) getText(ctx context.Context, endpoint string) (string, error) {
	var lastErr error
	for attempt := 0; attempt < c.retries; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		if err != nil {
			return "", err
		}
		req.Header.Set("User-Agent", c.userAgent)
		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
		} else {
			data, readErr := io.ReadAll(resp.Body)
			closeErr := resp.Body.Close()
			if readErr != nil {
				lastErr = readErr
			} else if closeErr != nil {
				lastErr = closeErr
			} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				return string(data), nil
			} else {
				lastErr = fmt.Errorf("GET %s: status %d", endpoint, resp.StatusCode)
			}
		}
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(c.retryWait):
		}
	}
	return "", lastErr
}

func extractItems(payload any) []map[string]any {
	switch value := payload.(type) {
	case []any:
		return mapsFromSlice(value)
	case map[string]any:
		for _, key := range []string{"papers", "items", "results", "data"} {
			if raw, ok := value[key]; ok {
				if items, ok := raw.([]any); ok {
					return mapsFromSlice(items)
				}
			}
		}
		return []map[string]any{value}
	default:
		return nil
	}
}

func mapsFromSlice(items []any) []map[string]any {
	out := make([]map[string]any, 0, len(items))
	for _, item := range items {
		if mapped, ok := item.(map[string]any); ok {
			out = append(out, mapped)
		}
	}
	return out
}

func parsePaper(value any) model.Paper {
	mapped, ok := value.(map[string]any)
	if !ok {
		return model.Paper{}
	}
	paper := model.Paper{
		ID:         firstString(mapped, "arxiv_id", "id", "slug"),
		Title:      firstString(mapped, "title"),
		Abstract:   firstString(mapped, "abstract", "summary", "description"),
		Authors:    stringSlice(mapped["authors"]),
		Categories: stringSlice(firstValue(mapped, "categories", "tags")),
		Links: model.Links{
			HuggingFace: firstString(mapped, "hf_paper_url", "url", "paper_url"),
			Arxiv:       firstString(mapped, "arxiv_url"),
			Paper:       firstString(mapped, "pdf_url", "paper_pdf_url", "paper_link"),
			PDF:         firstString(mapped, "pdf_url"),
			GitHub:      stringSlice(firstValue(mapped, "github_urls", "github", "code_urls")),
			Project:     stringSlice(firstValue(mapped, "project_urls", "resources")),
		},
		Upvotes:    intValue(firstValue(mapped, "upvotes", "likes")),
		Comments:   intValue(firstValue(mapped, "comments_count", "comments")),
		Discussion: intValue(firstValue(mapped, "discussion_count", "discussions")),
	}
	if paper.ID == "" && paper.Links.Arxiv != "" {
		paper.ID = trailingID(paper.Links.Arxiv)
	}
	if paper.ID == "" && paper.Links.HuggingFace != "" {
		paper.ID = trailingID(paper.Links.HuggingFace)
	}
	paper.PublishedAt = parseDate(firstString(mapped, "publication_date", "published_at", "created_at"))
	if paper.Links.HuggingFace == "" && paper.ID != "" {
		paper.Links.HuggingFace = strings.TrimRight("https://huggingface.co/papers/"+paper.ID, "/")
	}
	if paper.Links.Arxiv == "" && paper.ID != "" {
		paper.Links.Arxiv = "https://arxiv.org/abs/" + paper.ID
	}
	paper.RawCommunity = map[string]int{"upvotes": paper.Upvotes, "comments": paper.Comments, "discussion": paper.Discussion}
	return paper
}

func firstValue(mapped map[string]any, keys ...string) any {
	for _, key := range keys {
		if value, ok := mapped[key]; ok {
			return value
		}
	}
	return nil
}

func firstString(mapped map[string]any, keys ...string) string {
	for _, key := range keys {
		if value, ok := mapped[key]; ok {
			switch typed := value.(type) {
			case string:
				return strings.TrimSpace(typed)
			}
		}
	}
	return ""
}

func stringSlice(value any) []string {
	switch typed := value.(type) {
	case []string:
		return typed
	case []any:
		out := make([]string, 0, len(typed))
		for _, item := range typed {
			switch v := item.(type) {
			case string:
				if strings.TrimSpace(v) != "" {
					out = append(out, strings.TrimSpace(v))
				}
			case map[string]any:
				if label := firstString(v, "name", "label", "title", "url"); label != "" {
					out = append(out, label)
				}
			}
		}
		return out
	case string:
		if typed == "" {
			return nil
		}
		parts := strings.Split(typed, ",")
		out := make([]string, 0, len(parts))
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				out = append(out, part)
			}
		}
		return out
	default:
		return nil
	}
}

func intValue(value any) int {
	switch typed := value.(type) {
	case float64:
		return int(typed)
	case int:
		return typed
	case string:
		i, _ := strconv.Atoi(strings.TrimSpace(typed))
		return i
	default:
		return 0
	}
}

func parseDate(value string) time.Time {
	for _, layout := range []string{time.RFC3339, "2006-01-02", "2006-01-02 15:04:05"} {
		if value == "" {
			return time.Time{}
		}
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed.UTC()
		}
	}
	return time.Time{}
}

func trailingID(value string) string {
	value = strings.TrimRight(value, "/")
	parts := strings.Split(value, "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}
