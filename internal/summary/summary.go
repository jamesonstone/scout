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
	detail := truncateWords(innovationSentence(paperSignalText(paper)), 26)
	if detail == "" {
		detail = "A targeted research contribution for modern AI systems"
	}
	return ensureSingleSentence(fmt.Sprintf("%s: %s.", cleanTitle(paper.Title), upperFirst(detail)))
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
	text := strings.ToLower(paperSignalText(paper))
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

func paperSignalText(paper model.Paper) string {
	return firstNonEmpty(paper.Abstract, cleanMarkdownText(paper.Markdown))
}

func cleanMarkdownText(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}

	lines := strings.Split(text, "\n")
	for i, line := range lines {
		heading := strings.Trim(strings.TrimSpace(line), "#* ")
		if strings.EqualFold(heading, "Abstract") && i+1 < len(lines) {
			return strings.TrimSpace(strings.Join(lines[i+1:], "\n"))
		}
	}

	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		lower := strings.ToLower(trimmed)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(lower, "title:") ||
			strings.HasPrefix(lower, "url source:") ||
			strings.HasPrefix(lower, "markdown content:") ||
			lower == "back to arxiv" ||
			lower == "why html?" ||
			lower == "report issue" ||
			lower == "back to abstract" {
			continue
		}
		filtered = append(filtered, trimmed)
	}
	return strings.Join(filtered, " ")
}

func innovationSentence(text string) string {
	sentences := sentenceBoundary.Split(strings.TrimSpace(text), -1)
	best := ""
	bestScore := 0
	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if sentence == "" {
			continue
		}
		score := innovationCueScore(strings.ToLower(sentence))
		if score > bestScore {
			best = sentence
			bestScore = score
		}
	}
	if best != "" {
		return best
	}
	return firstSentence(text)
}

func innovationCueScore(sentence string) int {
	score := 0
	for _, cue := range []string{
		"we introduce",
		"we propose",
		"we present",
		"we develop",
		"we show",
		"we demonstrate",
		"we identify",
		"we characterize",
		"we study",
		"this paper introduces",
		"this paper proposes",
		"this paper presents",
		"to address this",
		"to tackle this",
		"our contributions",
	} {
		if strings.Contains(sentence, cue) {
			score += 10
		}
	}
	for _, cue := range []string{"framework", "benchmark", "method", "model", "algorithm", "pipeline", "system"} {
		if strings.Contains(sentence, cue) {
			score++
		}
	}
	return score
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

func upperFirst(text string) string {
	if text == "" {
		return text
	}
	runes := []rune(text)
	runes[0] = []rune(strings.ToUpper(string(runes[0])))[0]
	return string(runes)
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
