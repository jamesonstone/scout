package scoring

import (
	"math"
	"strings"

	"github.com/jamesonstone/scout/internal/model"
)

var defaultWeights = model.Weights{
	Novelty:                 0.20,
	PracticalImpact:         0.20,
	TechnicalDepth:          0.15,
	ImplementationPotential: 0.15,
	Relevance:               0.15,
	CommunitySignal:         0.10,
	SummaryConfidence:       0.05,
}

type Scorer struct{}

func New() Scorer { return Scorer{} }

func (Scorer) Score(paper model.Paper) model.ScoreBreakdown {
	text := strings.ToLower(strings.Join([]string{paper.Title, paper.Abstract, paper.Markdown, strings.Join(paper.Categories, " ")}, " "))
	novelty := keywordScore(text, 45, map[string]int{
		"novel": 12, "new": 8, "first": 10, "state-of-the-art": 12, "sota": 12, "breakthrough": 14, "agent": 8, "reasoning": 8,
		"multimodal": 10, "benchmark": 6, "scaling": 10, "memory": 10, "orchestration": 12,
	})
	practical := keywordScore(text, 40, map[string]int{
		"deploy": 12, "production": 12, "real-world": 10, "inference": 10, "latency": 10, "cost": 8, "efficiency": 10,
		"evaluation": 8, "tool": 8, "framework": 8, "infrastructure": 12, "software": 10,
	})
	technical := keywordScore(text, min(35+len(strings.Fields(paper.Abstract))/6, 55), map[string]int{
		"algorithm": 10, "architecture": 10, "theorem": 12, "optimization": 8, "dataset": 6, "benchmark": 8,
		"ablation": 10, "training": 8, "alignment": 10, "diffusion": 8, "transformer": 8,
	})
	implementation := keywordScore(text, 35+len(paper.Links.GitHub)*12, map[string]int{
		"open source": 16, "github": 16, "code": 10, "implementation": 12, "reproducible": 10, "library": 10,
		"sdk": 10, "api": 8, "agent": 8,
	})
	relevance := keywordScore(text, 30, map[string]int{
		"ai agent": 18, "agents": 16, "llm": 16, "language model": 16, "orchestration": 16, "evaluation": 14,
		"memory": 14, "infrastructure": 14, "software engineering": 16, "developer": 12, "benchmark": 8,
	})
	community := clamp(20+paper.Upvotes*5+paper.Comments*3+paper.Discussion*2+len(paper.Links.GitHub)*8, 0, 100)
	confidence := metadataConfidence(paper)
	overall := int(math.Round(
		float64(novelty)*defaultWeights.Novelty +
			float64(practical)*defaultWeights.PracticalImpact +
			float64(technical)*defaultWeights.TechnicalDepth +
			float64(implementation)*defaultWeights.ImplementationPotential +
			float64(relevance)*defaultWeights.Relevance +
			float64(community)*defaultWeights.CommunitySignal +
			float64(confidence)*defaultWeights.SummaryConfidence,
	))
	return model.ScoreBreakdown{
		Novelty:                 novelty,
		PracticalImpact:         practical,
		TechnicalDepth:          technical,
		ImplementationPotential: implementation,
		Relevance:               relevance,
		CommunitySignal:         community,
		SummaryConfidence:       confidence,
		Overall:                 clamp(overall, 0, 100),
		Weights:                 defaultWeights,
	}
}

func Recommendation(score int) model.Recommendation {
	switch {
	case score >= 80:
		return model.RecommendationRead
	case score >= 60:
		return model.RecommendationWorthWatching
	default:
		return model.RecommendationSkip
	}
}

func metadataConfidence(paper model.Paper) int {
	score := 20
	if paper.Abstract != "" {
		score += 20
	}
	if paper.Markdown != "" {
		score += 25
	}
	if len(paper.Authors) > 0 {
		score += 10
	}
	if len(paper.Categories) > 0 {
		score += 10
	}
	if paper.Links.Arxiv != "" || paper.Links.HuggingFace != "" {
		score += 10
	}
	if len(paper.Links.GitHub) > 0 {
		score += 5
	}
	return clamp(score, 0, 100)
}

func keywordScore(text string, base int, weights map[string]int) int {
	score := base
	for keyword, weight := range weights {
		if strings.Contains(text, keyword) {
			score += weight
		}
	}
	return clamp(score, 0, 100)
}

func clamp(value, low, high int) int {
	if value < low {
		return low
	}
	if value > high {
		return high
	}
	return value
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
