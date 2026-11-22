package surveys

import (
	"math"
	"testing"
	"time"

	"github.com/plutov/formulosity/api/pkg/services"
	"github.com/plutov/formulosity/api/pkg/types"
)

func intPtr(value int) *int {
	return &value
}

func timePtr(value time.Time) *time.Time {
	return &value
}

func testSurvey() types.Survey {
	questions := []types.Question{
		{
			UUID:    "q-employment",
			ID:      "employment",
			Type:    types.QuestionType_DropdownSingle,
			Label:   "Employment",
			Options: []string{"yes", "no"},
		},
		{
			UUID:    "q-interests",
			ID:      "interests",
			Type:    types.QuestionType_DropdownMultiple,
			Label:   "Interests",
			Options: []string{"Go", "Rust", "TypeScript"},
		},
		{
			UUID:  "q-contact",
			ID:    "contact",
			Type:  types.QuestionType_YesNo,
			Label: "Contact",
		},
		{
			UUID:  "q-satisfaction",
			ID:    "satisfaction",
			Type:  types.QuestionType_Rating,
			Label: "Satisfaction",
			Min:   intPtr(1),
			Max:   intPtr(5),
		},
		{
			UUID:    "q-cities",
			ID:      "cities",
			Type:    types.QuestionType_Ranking,
			Label:   "Cities",
			Options: []string{"Berlin", "Munich", "Paris"},
		},
		{
			UUID:  "q-wall",
			ID:    "wall",
			Type:  types.QuestionType_Date,
			Label: "Wall",
		},
		{
			UUID:  "q-name",
			ID:    "name",
			Type:  types.QuestionType_ShortText,
			Label: "Name",
		},
	}

	return types.Survey{
		UUID: "survey-1",
		Name: "test",
		Config: &types.SurveyConfig{
			Title: "Test survey",
			Questions: &types.Questions{
				Questions: questions,
			},
		},
	}
}

func testSessions() []types.SurveySession {
	base := time.Date(2026, time.November, 1, 10, 0, 0, 0, time.UTC)

	return []types.SurveySession{
		{
			UUID:        "s1",
			Status:      types.SurveySessionStatus_Completed,
			CreatedAt:   base,
			CompletedAt: timePtr(base.Add(100 * time.Second)),
			QuestionAnswers: []types.QuestionAnswer{
				{QuestionID: "employment", QuestionUUID: "q-employment", Answer: &types.SingleOptionAnswer{AnswerValue: "yes"}},
				{QuestionID: "interests", QuestionUUID: "q-interests", Answer: &types.MultiOptionsAnswer{AnswerValue: []string{"Go", "Rust"}}},
				{QuestionID: "contact", QuestionUUID: "q-contact", Answer: &types.BoolAnswer{AnswerValue: true}},
				{QuestionID: "satisfaction", QuestionUUID: "q-satisfaction", Answer: &types.NumberAnswer{AnswerValue: 5}},
				{QuestionID: "cities", QuestionUUID: "q-cities", Answer: &types.MultiOptionsAnswer{AnswerValue: []string{"Berlin", "Munich", "Paris"}}},
				{QuestionID: "wall", QuestionUUID: "q-wall", Answer: &types.DateAnswer{AnswerValue: "1961-08-13"}},
				{QuestionID: "name", QuestionUUID: "q-name", Answer: &types.TextAnswer{AnswerValue: "Alice"}},
			},
		},
		{
			UUID:        "s2",
			Status:      types.SurveySessionStatus_Completed,
			CreatedAt:   base.Add(24 * time.Hour),
			CompletedAt: timePtr(base.Add(24*time.Hour + 200*time.Second)),
			QuestionAnswers: []types.QuestionAnswer{
				{QuestionID: "employment", QuestionUUID: "q-employment", Answer: &types.SingleOptionAnswer{AnswerValue: "no"}},
				{QuestionID: "interests", QuestionUUID: "q-interests", Answer: &types.MultiOptionsAnswer{AnswerValue: []string{"Go"}}},
				{QuestionID: "contact", QuestionUUID: "q-contact", Answer: &types.BoolAnswer{AnswerValue: false}},
				{QuestionID: "satisfaction", QuestionUUID: "q-satisfaction", Answer: &types.NumberAnswer{AnswerValue: 3}},
				{QuestionID: "cities", QuestionUUID: "q-cities", Answer: &types.MultiOptionsAnswer{AnswerValue: []string{"Munich", "Berlin", "Paris"}}},
				{QuestionID: "wall", QuestionUUID: "q-wall", Answer: &types.DateAnswer{AnswerValue: "1961-08-13"}},
				{QuestionID: "name", QuestionUUID: "q-name", Answer: &types.TextAnswer{AnswerValue: "Bob"}},
			},
		},
		{
			UUID:        "s3",
			Status:      types.SurveySessionStatus_Completed,
			CreatedAt:   base.Add(48 * time.Hour),
			CompletedAt: timePtr(base.Add(48*time.Hour + 300*time.Second)),
			QuestionAnswers: []types.QuestionAnswer{
				{QuestionID: "employment", QuestionUUID: "q-employment", Answer: &types.SingleOptionAnswer{AnswerValue: "yes"}},
				{QuestionID: "interests", QuestionUUID: "q-interests", Answer: &types.MultiOptionsAnswer{AnswerValue: []string{}}},
				{QuestionID: "contact", QuestionUUID: "q-contact", Answer: &types.BoolAnswer{AnswerValue: true}},
				{QuestionID: "satisfaction", QuestionUUID: "q-satisfaction", Answer: &types.NumberAnswer{AnswerValue: 3}},
				{QuestionID: "cities", QuestionUUID: "q-cities", Answer: &types.MultiOptionsAnswer{AnswerValue: []string{"Berlin", "Paris", "Munich"}}},
				{QuestionID: "wall", QuestionUUID: "q-wall", Answer: &types.DateAnswer{AnswerValue: "1989-11-09"}},
				{QuestionID: "name", QuestionUUID: "q-name", Answer: &types.TextAnswer{AnswerValue: "Carol"}},
			},
		},
		{
			UUID:      "s4",
			Status:    types.SurveySessionStatus_InProgress,
			CreatedAt: base.Add(72 * time.Hour),
			QuestionAnswers: []types.QuestionAnswer{
				{QuestionID: "employment", QuestionUUID: "q-employment", Answer: &types.SingleOptionAnswer{AnswerValue: "yes"}},
			},
		},
	}
}

func TestBuildSurveyAnalytics(t *testing.T) {
	survey := testSurvey()
	analytics := BuildSurveyAnalytics(&survey, testSessions())

	if analytics.Totals.SessionsTotal != 4 {
		t.Fatalf("expected 4 total sessions, got %d", analytics.Totals.SessionsTotal)
	}
	if analytics.Totals.SessionsCompleted != 3 {
		t.Fatalf("expected 3 completed sessions, got %d", analytics.Totals.SessionsCompleted)
	}
	if analytics.Totals.SessionsInProgress != 1 {
		t.Fatalf("expected 1 in-progress session, got %d", analytics.Totals.SessionsInProgress)
	}
	assertFloat(t, "completion rate", analytics.Totals.CompletionRate, 75)
	if analytics.Totals.AverageCompletionSeconds == nil {
		t.Fatal("expected average completion seconds")
	}
	if *analytics.Totals.AverageCompletionSeconds != 200 {
		t.Fatalf(
			"expected average completion seconds 200, got %d",
			*analytics.Totals.AverageCompletionSeconds,
		)
	}

	if len(analytics.Timeline) != 4 {
		t.Fatalf("expected 4 timeline points, got %d", len(analytics.Timeline))
	}
	for i, expectedDate := range []string{
		"2026-11-01",
		"2026-11-02",
		"2026-11-03",
		"2026-11-04",
	} {
		if analytics.Timeline[i].Date != expectedDate {
			t.Fatalf("expected timeline date %s, got %s", expectedDate, analytics.Timeline[i].Date)
		}
		if analytics.Timeline[i].Started != 1 {
			t.Fatalf("expected 1 started session on %s, got %d", expectedDate, analytics.Timeline[i].Started)
		}
	}
	if analytics.Timeline[3].Completed != 0 {
		t.Fatalf("expected no completed sessions on 2026-11-04, got %d", analytics.Timeline[3].Completed)
	}

	if len(analytics.Questions) != 6 {
		t.Fatalf("expected 6 chartable questions, got %d", len(analytics.Questions))
	}

	employment := findQuestionAnalytics(t, analytics, "employment")
	if employment.AnsweredCount != 3 {
		t.Fatalf("expected 3 employment answers, got %d", employment.AnsweredCount)
	}
	assertBucket(t, employment.Buckets, "yes", 2, 66.7)
	assertBucket(t, employment.Buckets, "no", 1, 33.3)

	interests := findQuestionAnalytics(t, analytics, "interests")
	if interests.AnsweredCount != 2 {
		t.Fatalf("expected 2 interest answers, got %d", interests.AnsweredCount)
	}
	assertBucket(t, interests.Buckets, "Go", 2, 100)
	assertBucket(t, interests.Buckets, "Rust", 1, 50)
	assertBucket(t, interests.Buckets, "TypeScript", 0, 0)

	contact := findQuestionAnalytics(t, analytics, "contact")
	if contact.AnsweredCount != 3 {
		t.Fatalf("expected 3 contact answers, got %d", contact.AnsweredCount)
	}
	assertBucket(t, contact.Buckets, "true", 2, 66.7)
	assertBucket(t, contact.Buckets, "false", 1, 33.3)

	satisfaction := findQuestionAnalytics(t, analytics, "satisfaction")
	if satisfaction.AnsweredCount != 3 {
		t.Fatalf("expected 3 satisfaction answers, got %d", satisfaction.AnsweredCount)
	}
	if len(satisfaction.Buckets) != 5 {
		t.Fatalf("expected rating buckets 1-5, got %d", len(satisfaction.Buckets))
	}
	assertBucket(t, satisfaction.Buckets, "3", 2, 66.7)
	assertBucket(t, satisfaction.Buckets, "5", 1, 33.3)
	if satisfaction.Numeric == nil {
		t.Fatal("expected numeric summary")
	}
	assertFloatPtr(t, "rating average", satisfaction.Numeric.Average, 3.67)
	assertFloatPtr(t, "rating median", satisfaction.Numeric.Median, 3)
	assertFloatPtr(t, "rating mode", satisfaction.Numeric.Mode, 3)

	cities := findQuestionAnalytics(t, analytics, "cities")
	if cities.AnsweredCount != 3 {
		t.Fatalf("expected 3 ranking answers, got %d", cities.AnsweredCount)
	}
	assertRanking(t, cities.Ranking, "Berlin", []int{2, 1, 0}, 1.33)
	assertRanking(t, cities.Ranking, "Munich", []int{1, 1, 1}, 2)
	assertRanking(t, cities.Ranking, "Paris", []int{0, 1, 2}, 2.67)

	wall := findQuestionAnalytics(t, analytics, "wall")
	if wall.AnsweredCount != 3 {
		t.Fatalf("expected 3 date answers, got %d", wall.AnsweredCount)
	}
	if wall.Granularity != "year" {
		t.Fatalf("expected year granularity, got %s", wall.Granularity)
	}
	assertBucket(t, wall.Buckets, "1961", 2, 66.7)
	assertBucket(t, wall.Buckets, "1989", 1, 33.3)
}

func TestBuildSurveyAnalyticsEmpty(t *testing.T) {
	survey := testSurvey()
	analytics := BuildSurveyAnalytics(&survey, []types.SurveySession{})

	if analytics.Totals.SessionsTotal != 0 {
		t.Fatalf("expected 0 total sessions, got %d", analytics.Totals.SessionsTotal)
	}
	if analytics.Totals.CompletionRate != 0 {
		t.Fatalf("expected 0 completion rate, got %f", analytics.Totals.CompletionRate)
	}
	if len(analytics.Timeline) != 0 {
		t.Fatalf("expected empty timeline, got %d", len(analytics.Timeline))
	}
	if len(analytics.Questions) != 6 {
		t.Fatalf("expected 6 chartable questions, got %d", len(analytics.Questions))
	}
	for _, question := range analytics.Questions {
		if question.AnsweredCount != 0 {
			t.Fatalf("expected 0 answers for %s, got %d", question.QuestionID, question.AnsweredCount)
		}
	}
}

func TestBuildSurveyAnalyticsNilConfig(t *testing.T) {
	survey := types.Survey{UUID: "survey-empty", Name: "empty"}

	analytics := BuildSurveyAnalytics(&survey, []types.SurveySession{})
	if analytics.SurveyUUID != survey.UUID {
		t.Fatalf("expected survey UUID %s, got %s", survey.UUID, analytics.SurveyUUID)
	}
	if len(analytics.Questions) != 0 {
		t.Fatalf("expected no questions for nil config, got %d", len(analytics.Questions))
	}
}

func TestAnswerConversionHandlesNilConfig(t *testing.T) {
	survey := types.Survey{UUID: "survey-empty", Name: "empty"}
	answers := []types.QuestionAnswer{
		{QuestionID: "name", QuestionUUID: "q-name", AnswerBytes: []byte(`{"value":"Alice"}`)},
	}

	converted := convertAnswerBytesToAnswerType(services.Services{}, &survey, answers)
	if len(converted) != len(answers) {
		t.Fatalf("expected %d answers, got %d", len(answers), len(converted))
	}

	session := types.SurveySession{QuestionAnswers: converted}
	if err := filterVisibleAnswers(&survey, &session); err != nil {
		t.Fatalf("expected nil config filtering to succeed, got %v", err)
	}
}

func TestBuildSurveyAnalyticsExcludesHiddenAnswers(t *testing.T) {
	condition := true
	survey := types.Survey{
		UUID: "survey-conditional",
		Name: "conditional",
		Config: &types.SurveyConfig{
			Title: "Conditional",
			Questions: &types.Questions{
				Questions: []types.Question{
					{
						UUID:  "q-contact",
						ID:    "contact",
						Type:  types.QuestionType_YesNo,
						Label: "Contact",
					},
					{
						UUID:    "q-followup",
						ID:      "followup",
						Type:    types.QuestionType_DropdownSingle,
						Label:   "Follow up",
						Options: []string{"a", "b"},
						VisibleIf: &types.QuestionCondition{
							Question: "contact",
							Equals:   condition,
						},
					},
				},
			},
		},
	}

	base := time.Date(2026, time.November, 1, 10, 0, 0, 0, time.UTC)
	sessions := []types.SurveySession{
		{
			UUID:        "s1",
			Status:      types.SurveySessionStatus_Completed,
			CreatedAt:   base,
			CompletedAt: timePtr(base.Add(10 * time.Second)),
			QuestionAnswers: []types.QuestionAnswer{
				{QuestionID: "contact", QuestionUUID: "q-contact", Answer: &types.BoolAnswer{AnswerValue: true}},
				{QuestionID: "followup", QuestionUUID: "q-followup", Answer: &types.SingleOptionAnswer{AnswerValue: "a"}},
			},
		},
		{
			UUID:        "s2",
			Status:      types.SurveySessionStatus_Completed,
			CreatedAt:   base.Add(24 * time.Hour),
			CompletedAt: timePtr(base.Add(24*time.Hour + 10*time.Second)),
			QuestionAnswers: []types.QuestionAnswer{
				{QuestionID: "contact", QuestionUUID: "q-contact", Answer: &types.BoolAnswer{AnswerValue: false}},
				{QuestionID: "followup", QuestionUUID: "q-followup", Answer: &types.SingleOptionAnswer{AnswerValue: "b"}},
			},
		},
	}

	for i := range sessions {
		if err := filterVisibleAnswers(&survey, &sessions[i]); err != nil {
			t.Fatalf("unable to filter visible answers: %v", err)
		}
	}

	analytics := BuildSurveyAnalytics(&survey, sessions)
	followup := findQuestionAnalytics(t, analytics, "followup")
	if !followup.Conditional {
		t.Fatal("expected followup question to be marked conditional")
	}
	if followup.AnsweredCount != 1 {
		t.Fatalf("expected 1 followup answer, got %d", followup.AnsweredCount)
	}
	assertBucket(t, followup.Buckets, "a", 1, 100)
	assertBucket(t, followup.Buckets, "b", 0, 0)
}

func TestBuildSurveyAnalyticsIncludesUnknownOptions(t *testing.T) {
	survey := types.Survey{
		UUID: "survey-unknown",
		Name: "unknown",
		Config: &types.SurveyConfig{
			Title: "Unknown option",
			Questions: &types.Questions{
				Questions: []types.Question{
					{
						UUID:    "q-choice",
						ID:      "choice",
						Type:    types.QuestionType_DropdownSingle,
						Label:   "Choice",
						Options: []string{"a", "b"},
					},
				},
			},
		},
	}

	base := time.Date(2026, time.November, 1, 10, 0, 0, 0, time.UTC)
	sessions := []types.SurveySession{
		{
			UUID:        "s1",
			Status:      types.SurveySessionStatus_Completed,
			CreatedAt:   base,
			CompletedAt: timePtr(base.Add(10 * time.Second)),
			QuestionAnswers: []types.QuestionAnswer{
				{QuestionID: "choice", QuestionUUID: "q-choice", Answer: &types.SingleOptionAnswer{AnswerValue: "c"}},
			},
		},
	}

	analytics := BuildSurveyAnalytics(&survey, sessions)
	choice := findQuestionAnalytics(t, analytics, "choice")
	if len(choice.Buckets) != 3 {
		t.Fatalf("expected 3 buckets including unknown option, got %d", len(choice.Buckets))
	}
	assertBucket(t, choice.Buckets, "c", 1, 100)
}

func findQuestionAnalytics(
	t *testing.T,
	analytics *types.SurveyAnalytics,
	questionID string,
) types.QuestionAnalytics {
	t.Helper()

	for _, question := range analytics.Questions {
		if question.QuestionID == questionID {
			return question
		}
	}

	t.Fatalf("question %s not found", questionID)
	return types.QuestionAnalytics{}
}

func assertBucket(
	t *testing.T,
	buckets []types.AnalyticsBucket,
	key string,
	count int,
	percentage float64,
) {
	t.Helper()

	for _, bucket := range buckets {
		if bucket.Key == key {
			if bucket.Count != count {
				t.Fatalf("expected bucket %s count %d, got %d", key, count, bucket.Count)
			}
			assertFloat(t, "bucket "+key+" percentage", bucket.Percentage, percentage)
			return
		}
	}

	t.Fatalf("bucket %s not found", key)
}

func assertRanking(
	t *testing.T,
	ranking []types.AnalyticsRankingOption,
	key string,
	rankCounts []int,
	averageRank float64,
) {
	t.Helper()

	for _, option := range ranking {
		if option.Key != key {
			continue
		}

		if len(option.RankCounts) != len(rankCounts) {
			t.Fatalf(
				"expected %d rank counts for %s, got %d",
				len(rankCounts),
				key,
				len(option.RankCounts),
			)
		}
		for i, count := range rankCounts {
			if option.RankCounts[i] != count {
				t.Fatalf(
					"expected rank count %d for %s at position %d, got %d",
					count,
					key,
					i,
					option.RankCounts[i],
				)
			}
		}
		if option.AverageRank == nil {
			t.Fatalf("expected average rank for %s", key)
		}
		assertFloat(t, "average rank "+key, *option.AverageRank, averageRank)
		return
	}

	t.Fatalf("ranking option %s not found", key)
}

func assertFloat(t *testing.T, label string, actual float64, expected float64) {
	t.Helper()

	if math.Abs(actual-expected) > 0.01 {
		t.Fatalf("expected %s %.2f, got %.2f", label, expected, actual)
	}
}

func assertFloatPtr(t *testing.T, label string, actual *float64, expected float64) {
	t.Helper()

	if actual == nil {
		t.Fatalf("expected %s to be set", label)
	}
	assertFloat(t, label, *actual, expected)
}
