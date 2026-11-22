package types

import "time"

type AnalyticsChartKind string

const (
	AnalyticsChartKind_SingleChoice   AnalyticsChartKind = "single_choice"
	AnalyticsChartKind_MultipleChoice AnalyticsChartKind = "multi_choice"
	AnalyticsChartKind_Boolean        AnalyticsChartKind = "boolean"
	AnalyticsChartKind_Rating         AnalyticsChartKind = "rating"
	AnalyticsChartKind_Ranking        AnalyticsChartKind = "ranking"
	AnalyticsChartKind_Date           AnalyticsChartKind = "date"
)

type SurveyAnalytics struct {
	SurveyUUID  string                   `json:"survey_uuid"`
	SurveyName  string                   `json:"survey_name"`
	SurveyTitle string                   `json:"survey_title"`
	GeneratedAt time.Time                `json:"generated_at"`
	Totals      AnalyticsTotals          `json:"totals"`
	Timeline    []AnalyticsTimelinePoint `json:"timeline"`
	Questions   []QuestionAnalytics      `json:"questions"`
}

type AnalyticsTotals struct {
	SessionsTotal            int     `json:"sessions_total"`
	SessionsCompleted        int     `json:"sessions_completed"`
	SessionsInProgress       int     `json:"sessions_in_progress"`
	CompletionRate           float64 `json:"completion_rate"`
	AverageCompletionSeconds *int64  `json:"average_completion_seconds,omitempty"`
}

type AnalyticsTimelinePoint struct {
	Date      string `json:"date"`
	Started   int    `json:"started"`
	Completed int    `json:"completed"`
}

type QuestionAnalytics struct {
	QuestionUUID  string                   `json:"question_uuid"`
	QuestionID    string                   `json:"question_id"`
	Label         string                   `json:"label"`
	Type          QuestionType             `json:"type"`
	ChartKind     AnalyticsChartKind       `json:"chart_kind"`
	Conditional   bool                     `json:"conditional"`
	AnsweredCount int                      `json:"answered_count"`
	Buckets       []AnalyticsBucket        `json:"buckets,omitempty"`
	Ranking       []AnalyticsRankingOption `json:"ranking,omitempty"`
	Numeric       *AnalyticsNumericSummary `json:"numeric_summary,omitempty"`
	Granularity   string                   `json:"granularity,omitempty"`
}

type AnalyticsBucket struct {
	Key        string  `json:"key"`
	Label      string  `json:"label"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

type AnalyticsRankingOption struct {
	Key         string   `json:"key"`
	Label       string   `json:"label"`
	RankCounts  []int    `json:"rank_counts"`
	AverageRank *float64 `json:"average_rank,omitempty"`
}

type AnalyticsNumericSummary struct {
	Average *float64 `json:"average,omitempty"`
	Median  *float64 `json:"median,omitempty"`
	Mode    *float64 `json:"mode,omitempty"`
}
