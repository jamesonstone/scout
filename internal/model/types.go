package model

import "time"

type Recommendation string

const (
	RecommendationRead          Recommendation = "Read"
	RecommendationWorthWatching Recommendation = "Worth Watching"
	RecommendationSkip          Recommendation = "Skip"
)

type Links struct {
	HuggingFace string   `json:"hugging_face"`
	Arxiv       string   `json:"arxiv"`
	GitHub      []string `json:"github,omitempty"`
	Paper       string   `json:"paper,omitempty"`
	Project     []string `json:"project,omitempty"`
	PDF         string   `json:"pdf,omitempty"`
}

type Paper struct {
	ID           string
	Title        string
	Abstract     string
	Authors      []string
	Categories   []string
	PublishedAt  time.Time
	Links        Links
	Upvotes      int
	Comments     int
	Discussion   int
	Markdown     string
	SourceDate   time.Time
	RawCommunity map[string]int
}

type ScoreBreakdown struct {
	Novelty                 int     `json:"novelty"`
	PracticalImpact         int     `json:"practical_impact"`
	TechnicalDepth          int     `json:"technical_depth"`
	ImplementationPotential int     `json:"implementation_potential"`
	Relevance               int     `json:"relevance"`
	CommunitySignal         int     `json:"community_signal"`
	SummaryConfidence       int     `json:"summary_confidence"`
	Overall                 int     `json:"overall"`
	Weights                 Weights `json:"weights"`
}

type Weights struct {
	Novelty                 float64 `json:"novelty"`
	PracticalImpact         float64 `json:"practical_impact"`
	TechnicalDepth          float64 `json:"technical_depth"`
	ImplementationPotential float64 `json:"implementation_potential"`
	Relevance               float64 `json:"relevance"`
	CommunitySignal         float64 `json:"community_signal"`
	SummaryConfidence       float64 `json:"summary_confidence"`
}

type ScoreSnapshot struct {
	ObservedAt string         `json:"observed_at"`
	Score      ScoreBreakdown `json:"score"`
}

type PaperRecord struct {
	ID                   string          `json:"id"`
	Title                string          `json:"title"`
	Authors              []string        `json:"authors,omitempty"`
	Categories           []string        `json:"categories,omitempty"`
	PublishedAt          string          `json:"published_at,omitempty"`
	FirstSeen            string          `json:"first_seen"`
	ObservedDates        []string        `json:"observed_dates"`
	Abstract             string          `json:"abstract,omitempty"`
	Markdown             string          `json:"markdown,omitempty"`
	InnovationSummary    string          `json:"innovation_summary"`
	WhyItMatters         []string        `json:"why_it_matters"`
	ImplementationAngle  []string        `json:"implementation_angle"`
	Caveat               string          `json:"caveat"`
	ExecutiveSummary     string          `json:"executive_summary"`
	EstimatedPriority    string          `json:"estimated_reading_priority"`
	Recommendation       Recommendation  `json:"recommendation"`
	Links                Links           `json:"links"`
	Score                ScoreBreakdown  `json:"score"`
	ScoreHistory         []ScoreSnapshot `json:"score_history"`
	Community            map[string]int  `json:"community,omitempty"`
	MetadataCompleteness int             `json:"metadata_completeness"`
}

type DailyObservation struct {
	Date     string   `json:"date"`
	PaperIDs []string `json:"paper_ids"`
}

type RunResult struct {
	DailyReportPath   string
	MonthlyReportPath string
	ProcessedCount    int
	ReusedCount       int
}
