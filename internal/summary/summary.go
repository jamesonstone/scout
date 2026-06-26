package summary

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/jamesonstone/scout/internal/model"
)

var sentenceBoundary = regexp.MustCompile(`[.!?]+`)

type Builder struct{}

func New() Builder { return Builder{} }

func (Builder) InnovationSummary(paper model.Paper) string {
	detail := truncateWords(firstSentence(firstNonEmpty(paper.Markdown, paper.Abstract)), 20)
	if detail == "" {
		detail = "a targeted research contribution for modern AI systems"
	}
	detail = normalizeInnovationDetail(detail)
	return ensureSingleSentence(fmt.Sprintf("%s introduces %s.", cleanTitle(paper.Title), lowerFirst(detail)))
}

func (Builder) WhyItMatters(paper model.Paper, score model.ScoreBreakdown) []string {
	bullets := []string{
		fmt.Sprintf("Overall signal %d/100 driven by novelty %d and practical impact %d.", score.Overall, score.Novelty, score.PracticalImpact),
		themeBullet(paper),
		communityBullet(paper),
	}
	return compactBullets(bullets, 3)
}

func (Builder) ImplementationAngle(paper model.Paper, score model.ScoreBreakdown) []string {
	bullets := []string{
		fmt.Sprintf("Implementation potential scores %d/100; prioritize adaptation paths for internal agent, evaluation, or platform workflows.", score.ImplementationPotential),
		githubBullet(paper),
		fmt.Sprintf("Technical depth scores %d/100, so a quick skim should focus on architecture, data, and evaluation sections before full adoption work.", score.TechnicalDepth),
	}
	return compactBullets(bullets, 3)
}

func (Builder) Caveat(paper model.Paper) string {
	text := strings.ToLower(firstNonEmpty(paper.Markdown, paper.Abstract))
	switch {
	case strings.Contains(text, "benchmark"):
		return "Evidence appears benchmark-centric, so verify transfer to production workloads before acting on the claims."
	case strings.Contains(text, "simulation"):
		return "The strongest evidence comes from simulated settings, so operational impact may be less certain in live systems."
	case len(paper.Links.GitHub) == 0:
		return "No linked implementation is available yet, which raises integration cost and lowers reproducibility confidence."
	default:
		return "The abstract exposes the headline result, but deeper implementation tradeoffs still require reading the full paper."
	}
}

func (Builder) ExecutiveSummary(paper model.Paper, score model.ScoreBreakdown) string {
	parts := []string{
		fmt.Sprintf("%s scores %d/100 for a mix of novelty, practical impact, and relevance to AI-agent and software-engineering workflows.", cleanTitle(paper.Title), score.Overall),
		truncateWords(strings.TrimSpace(firstNonEmpty(paper.Abstract, paper.Markdown)), 90),
	}
	if paper.Markdown != "" {
		parts = append(parts, truncateWords(firstSentence(paper.Markdown), 45))
	}
	summary := strings.Join(compact(parts), " ")
	return truncateWords(summary, 300)
}

func (Builder) ReadingPriority(score int) string {
	switch {
	case score >= 85:
		return "Immediate — skim today and queue a deeper read this week"
	case score >= 70:
		return "High — skim this week if you work on agents, evaluation, or infrastructure"
	case score >= 55:
		return "Medium — revisit when adjacent roadmap work appears"
	default:
		return "Low — archive unless the topic is directly relevant"
	}
}

func (Builder) ExecutiveSignal(papers []model.PaperRecord, day string) string {
	if len(papers) == 0 {
		return fmt.Sprintf("No papers were available for %s.", day)
	}
	counts := map[string]int{}
	for _, paper := range papers {
		for _, category := range paper.Categories {
			counts[category]++
		}
	}
	themes := topKeys(counts, 3)
	if len(themes) == 0 {
		themes = []string{"AI systems", "evaluation", "infrastructure"}
	}
	return fmt.Sprintf("%s is led by %s, with the strongest papers skewing toward production-minded advances that pair novelty with implementation value.", day, joinThemes(themes))
}

func cleanTitle(title string) string {
	title = strings.TrimSpace(title)
	if title == "" {
		return "This paper"
	}
	return title
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func firstSentence(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}
	parts := sentenceBoundary.Split(text, -1)
	if len(parts) == 0 {
		return text
	}
	return strings.TrimSpace(parts[0])
}

func ensureSingleSentence(text string) string {
	text = strings.TrimSpace(sentenceBoundary.ReplaceAllString(text, ""))
	if text == "" {
		return "This paper introduces a research contribution for AI systems."
	}
	return text + "."
}

func truncateWords(text string, limit int) string {
	words := strings.Fields(text)
	if len(words) <= limit {
		return strings.Join(words, " ")
	}
	return strings.Join(words[:limit], " ")
}

func themeBullet(paper model.Paper) string {
	if len(paper.Categories) == 0 {
		return "It maps to cross-cutting AI systems work even without explicit category metadata."
	}
	return fmt.Sprintf("Primary categories: %s.", strings.Join(paper.Categories, ", "))
}

func communityBullet(paper model.Paper) string {
	if paper.Upvotes == 0 && paper.Comments == 0 && paper.Discussion == 0 {
		return "Community signal is still emerging, so the score leans more on technical and implementation cues than popularity."
	}
	return fmt.Sprintf("Community signal includes %d upvote(s) and %d comment(s), which helps separate durable interest from title-only curiosity.", paper.Upvotes, paper.Comments)
}

func githubBullet(paper model.Paper) string {
	if len(paper.Links.GitHub) == 0 {
		return "No linked repository is present, so expect more translation work before the ideas are production-ready."
	}
	return fmt.Sprintf("Linked implementation resources are available via %s, which lowers the cost of benchmarking or prototyping the paper's ideas.", strings.Join(paper.Links.GitHub, ", "))
}

func lowerFirst(text string) string {
	if text == "" {
		return text
	}
	runes := []rune(text)
	runes[0] = []rune(strings.ToLower(string(runes[0])))[0]
	return string(runes)
}

func normalizeInnovationDetail(text string) string {
	text = strings.TrimSpace(strings.TrimRight(text, ".!? "))
	for _, prefix := range []string{
		"this paper introduces ",
		"this paper proposes ",
		"this paper presents ",
		"this paper studies ",
		"this paper covers ",
		"this paper describes ",
		"this paper develops ",
		"we introduce ",
		"we propose ",
		"we present ",
		"we study ",
		"we describe ",
	} {
		if strings.HasPrefix(strings.ToLower(text), prefix) {
			text = strings.TrimSpace(text[len(prefix):])
			break
		}
	}
	if text == "" {
		return "a targeted research contribution for modern AI systems"
	}
	return text
}

func joinThemes(themes []string) string {
	switch len(themes) {
	case 0:
		return "cross-cutting AI systems work"
	case 1:
		return themes[0]
	case 2:
		return themes[0] + " and " + themes[1]
	default:
		return strings.Join(themes[:len(themes)-1], ", ") + ", and " + themes[len(themes)-1]
	}
}

func compactBullets(items []string, max int) []string {
	filtered := compact(items)
	if len(filtered) > max {
		filtered = filtered[:max]
	}
	return filtered
}

func compact(items []string) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item != "" {
			out = append(out, item)
		}
	}
	return out
}

func topKeys(counts map[string]int, limit int) []string {
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
	if len(pairs) > limit {
		pairs = pairs[:limit]
	}
	keys := make([]string, 0, len(pairs))
	for _, pair := range pairs {
		keys = append(keys, pair.key)
	}
	return keys
}
