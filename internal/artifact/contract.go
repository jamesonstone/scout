package artifact

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jamesonstone/scout/internal/model"
)

const (
	MaxPaperRecordBytes   = 8 * 1024
	MaxCategories         = 6
	MaxDailyFullSummaries = 5
	MaxMonthlyTopPapers   = 10
	MaxDistilledTextChars = 360
)

var forbiddenPaperRecordFields = []string{
	"abstract",
	"authors",
	"community",
	"executive_summary",
	"markdown",
	"metadata_completeness",
	"published_at",
	"score_history",
	"estimated_reading_priority",
}

func CompactPaperRecord(record model.PaperRecord) model.PaperRecord {
	record.Categories = LimitStrings(record.Categories, MaxCategories)
	record.InnovationSummary = LimitText(record.InnovationSummary, MaxDistilledTextChars)
	record.WhyItMatters = compactWhyItMatters(record.WhyItMatters, record.Categories)
	record.ImplementationAngle = LimitTextList(LimitStrings(record.ImplementationAngle, 3), MaxDistilledTextChars)
	record.Caveat = LimitText(record.Caveat, MaxDistilledTextChars)
	record.Links = compactLinks(record.Links)
	return record
}

func LimitStrings(values []string, limit int) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, min(len(values), limit))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
		if len(out) == limit {
			break
		}
	}
	return out
}

func LimitTextList(values []string, limit int) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = LimitText(value, limit)
		if value != "" {
			out = append(out, value)
		}
	}
	return out
}

func LimitText(value string, limit int) string {
	value = strings.Join(strings.Fields(value), " ")
	if len(value) <= limit {
		return value
	}
	cut := strings.LastIndex(value[:limit], " ")
	if cut < limit/2 {
		cut = limit
	}
	return strings.TrimRight(value[:cut], " ,;:-") + "."
}

func compactWhyItMatters(values []string, categories []string) []string {
	values = LimitStrings(values, 3)
	for i, value := range values {
		if strings.HasPrefix(value, "Primary categories:") && len(categories) > 0 {
			values[i] = "Primary categories: " + strings.Join(categories, ", ") + "."
		}
	}
	return LimitTextList(values, MaxDistilledTextChars)
}

func MarshalPaperRecord(record model.PaperRecord) ([]byte, error) {
	record = CompactPaperRecord(record)
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return nil, err
	}
	data = append(data, '\n')
	if err := ValidatePaperRecordJSON("paper record", data); err != nil {
		return nil, err
	}
	return data, nil
}

func ValidatePaperRecordJSON(path string, data []byte) error {
	if len(data) > MaxPaperRecordBytes {
		return fmt.Errorf("%s exceeds %d bytes (%d bytes)", path, MaxPaperRecordBytes, len(data))
	}
	var payload map[string]json.RawMessage
	if err := json.Unmarshal(data, &payload); err != nil {
		return fmt.Errorf("decode %s: %w", path, err)
	}
	for _, field := range forbiddenPaperRecordFields {
		if _, ok := payload[field]; ok {
			return fmt.Errorf("%s contains forbidden curated-record field %q", path, field)
		}
	}
	if raw, ok := payload["categories"]; ok {
		var categories []string
		if err := json.Unmarshal(raw, &categories); err != nil {
			return fmt.Errorf("decode categories in %s: %w", path, err)
		}
		if len(categories) > MaxCategories {
			return fmt.Errorf("%s has %d categories; max is %d", path, len(categories), MaxCategories)
		}
	}
	return nil
}

func compactLinks(links model.Links) model.Links {
	links.HuggingFace = strings.TrimSpace(links.HuggingFace)
	links.Arxiv = strings.TrimSpace(links.Arxiv)
	links.Paper = strings.TrimSpace(links.Paper)
	links.PDF = strings.TrimSpace(links.PDF)
	links.GitHub = LimitStrings(links.GitHub, 3)
	links.Project = LimitStrings(links.Project, 3)
	return links
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
