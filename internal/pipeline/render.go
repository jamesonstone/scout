package pipeline

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/jamesonstone/scout/internal/artifact"
)

func RenderDaily(day time.Time, papers []PaperRecord, signal string, monthlyPath string) string {
	sorted := sortRecords(papers)
	var top, additional, watchlist []PaperRecord
	for _, paper := range sorted {
		switch paper.Recommendation {
		case RecommendationRead:
			top = append(top, paper)
		case RecommendationWorthWatching:
			additional = append(additional, paper)
		default:
			watchlist = append(watchlist, paper)
		}
	}

	var b strings.Builder
	fmt.Fprintf(&b, "# Scout Daily Intelligence Briefing — %s\n\n", day.Format("2006-01-02"))
	b.WriteString("## Executive Signal\n\n")
	b.WriteString(signal)
	b.WriteString("\n\n")
	writeFullSection(&b, "Top Papers", top, artifact.MaxDailyFullSummaries)
	if len(top) > artifact.MaxDailyFullSummaries {
		writeCompactRows(&b, "Additional Read Papers", top[artifact.MaxDailyFullSummaries:])
	}
	writeCompactRows(&b, "Additional Papers", additional)
	writeCompactRows(&b, "Watchlist", watchlist)
	b.WriteString("## Archive\n\n")
	fmt.Fprintf(&b, "- Daily record count: %d\n", len(sorted))
	fmt.Fprintf(&b, "- Active monthly report: %s\n", monthlyPath)
	b.WriteString("- Persistent paper records live under `data/papers/`.\n")
	for _, paper := range sorted {
		fmt.Fprintf(&b, "- %d/100 %s — published %s — [%s](%s) — %s\n", paper.Score.Overall, paper.Recommendation, publishedDateLabel(paper), paper.Title, bestLink(paper.Links), paper.InnovationSummary)
	}
	return b.String()
}

func RenderMonthly(month time.Time, papers []PaperRecord) string {
	sorted := sortRecords(papers)
	top10 := sorted
	if len(top10) > artifact.MaxMonthlyTopPapers {
		top10 = top10[:artifact.MaxMonthlyTopPapers]
	}

	var b strings.Builder
	fmt.Fprintf(&b, "# Scout Monthly Intelligence Briefing — %s\n\n", month.Format("2006-01"))
	b.WriteString("## Top 10 Papers\n\n")
	for i, paper := range top10 {
		fmt.Fprintf(&b, "%d. **%s** — %d/100 (%s), published %s. %s\n", i+1, paper.Title, paper.Score.Overall, paper.Recommendation, publishedDateLabel(paper), paper.InnovationSummary)
	}
	b.WriteString("\n## Theme Analysis\n\n")
	b.WriteString(renderThemeAnalysis(sorted))
	b.WriteString("\n\n## Rising Papers\n\n")
	for _, paper := range risingPapers(sorted) {
		fmt.Fprintf(&b, "- **%s** — %d/100, published %s, first seen %s\n", paper.Title, paper.Score.Overall, publishedDateLabel(paper), paper.FirstSeen)
	}
	b.WriteString("\n## Historical Rankings\n\n")
	for i, paper := range sorted {
		fmt.Fprintf(&b, "%d. **%s** — current %d/100; published %s; observed %d time(s); first seen %s\n", i+1, paper.Title, paper.Score.Overall, publishedDateLabel(paper), len(paper.ObservedDates), paper.FirstSeen)
	}
	b.WriteString("\n## Complete Monthly Index\n\n")
	for i, paper := range sorted {
		fmt.Fprintf(&b, "%d. [%s](%s) — %s\n", i+1, paper.Title, firstNonEmpty(paper.Links.HuggingFace, paper.Links.Arxiv, paper.Links.Paper), categoriesLabel(paper.Categories))
	}
	return b.String()
}

func writeFullSection(b *strings.Builder, heading string, papers []PaperRecord, limit int) {
	fmt.Fprintf(b, "## %s\n\n", heading)
	if len(papers) == 0 {
		b.WriteString("No papers in this section.\n\n")
		return
	}
	if len(papers) > limit {
		papers = papers[:limit]
	}
	for i, paper := range papers {
		writePaper(b, i+1, paper)
	}
}

func writeCompactRows(b *strings.Builder, heading string, papers []PaperRecord) {
	fmt.Fprintf(b, "## %s\n\n", heading)
	if len(papers) == 0 {
		b.WriteString("No papers in this section.\n\n")
		return
	}
	for _, paper := range papers {
		fmt.Fprintf(b, "- **%s** — %d/100 (%s). %s [Source](%s)\n", paper.Title, paper.Score.Overall, paper.Recommendation, paper.InnovationSummary, bestLink(paper.Links))
	}
	b.WriteString("\n")
}

func writePaper(b *strings.Builder, rank int, paper PaperRecord) {
	fmt.Fprintf(b, "### %d. %s\n\n", rank, paper.Title)
	fmt.Fprintf(b, "- **Overall score:** %d/100\n", paper.Score.Overall)
	fmt.Fprintf(b, "- **Recommendation:** %s\n", paper.Recommendation)
	fmt.Fprintf(b, "- **Published:** %s\n", publishedDateLabel(paper))
	fmt.Fprintf(b, "- **First fetched:** %s\n", paper.FirstSeen)
	fmt.Fprintf(b, "- **Categories:** %s\n", categoriesLabel(paper.Categories))
	fmt.Fprintf(b, "- **Innovation Summary:** %s\n", paper.InnovationSummary)
	b.WriteString("- **Why It Matters:**\n")
	for _, bullet := range paper.WhyItMatters {
		fmt.Fprintf(b, "  - %s\n", bullet)
	}
	b.WriteString("- **Implementation Angle:**\n")
	for _, bullet := range paper.ImplementationAngle {
		fmt.Fprintf(b, "  - %s\n", bullet)
	}
	fmt.Fprintf(b, "- **Caveat:** %s\n", paper.Caveat)
	fmt.Fprintf(b, "- **Links:** Hugging Face: %s | arXiv: %s | GitHub: %s | Paper: %s\n", nonEmptyLink(paper.Links.HuggingFace), nonEmptyLink(paper.Links.Arxiv), strings.Join(orNA(paper.Links.GitHub), ", "), nonEmptyLink(firstNonEmpty(paper.Links.Paper, paper.Links.PDF)))
	b.WriteString("\n")
}

func sortRecords(papers []PaperRecord) []PaperRecord {
	sorted := append([]PaperRecord(nil), papers...)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Score.Overall == sorted[j].Score.Overall {
			return sorted[i].Title < sorted[j].Title
		}
		return sorted[i].Score.Overall > sorted[j].Score.Overall
	})
	return sorted
}

func renderThemeAnalysis(papers []PaperRecord) string {
	counts := map[string]int{}
	for _, paper := range papers {
		for _, category := range paper.Categories {
			counts[category]++
		}
	}
	if len(counts) == 0 {
		return "Theme coverage is still sparse, but the month already emphasizes papers with direct engineering relevance."
	}
	type pair struct {
		key   string
		count int
	}
	pairs := make([]pair, 0, len(counts))
	for key, count := range counts {
		pairs = append(pairs, pair{key: key, count: count})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count == pairs[j].count {
			return pairs[i].key < pairs[j].key
		}
		return pairs[i].count > pairs[j].count
	})
	if len(pairs) > 3 {
		pairs = pairs[:3]
	}
	labels := make([]string, 0, len(pairs))
	for _, pair := range pairs {
		labels = append(labels, fmt.Sprintf("%s (%d)", pair.key, pair.count))
	}
	return fmt.Sprintf("Dominant monthly themes: %s. Rankings are recomputed after each daily run, so this section reflects the latest cumulative signal instead of original publication order.", strings.Join(labels, ", "))
}

func risingPapers(papers []PaperRecord) []PaperRecord {
	out := make([]PaperRecord, 0, len(papers))
	for _, paper := range papers {
		if paper.Score.Overall >= 70 {
			out = append(out, paper)
		}
	}
	if len(out) > 5 {
		out = out[:5]
	}
	return out
}

func nonEmptyLink(value string) string {
	if value == "" {
		return "N/A"
	}
	return value
}

func bestLink(links Links) string {
	link := firstNonEmpty(links.HuggingFace, links.Arxiv, links.Paper, links.PDF)
	if link == "" {
		return "#"
	}
	return link
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func orNA(values []string) []string {
	if len(values) == 0 {
		return []string{"N/A"}
	}
	return values
}

func categoriesLabel(categories []string) string {
	if len(categories) == 0 {
		return "N/A"
	}
	return strings.Join(categories, ", ")
}

func publishedDateLabel(paper PaperRecord) string {
	if paper.PublishedDate == "" {
		return "unavailable"
	}
	return paper.PublishedDate
}
