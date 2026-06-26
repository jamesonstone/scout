package pipeline

import "github.com/jamesonstone/scout/internal/model"

type Recommendation = model.Recommendation

const (
	RecommendationRead          = model.RecommendationRead
	RecommendationWorthWatching = model.RecommendationWorthWatching
	RecommendationSkip          = model.RecommendationSkip
)

type Links = model.Links
type Paper = model.Paper
type ScoreBreakdown = model.ScoreBreakdown
type Weights = model.Weights
type PaperRecord = model.PaperRecord
type DailyObservation = model.DailyObservation
type RunResult = model.RunResult
