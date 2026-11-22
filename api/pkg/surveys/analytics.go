package surveys

import (
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/plutov/formulosity/api/pkg/services"
	"github.com/plutov/formulosity/api/pkg/types"
)

const analyticsMaxSessions = 1000000
const maxRatingBuckets = 1000

// GetSurveyAnalytics loads sessions and builds the aggregated analytics payload.
func GetSurveyAnalytics(svc services.Services, survey types.Survey) (*types.SurveyAnalytics, error) {
	logCtx := svc.Logger.With("survey_uuid", survey.UUID)
	logCtx.Info("getting survey analytics")

	sessions, _, err := svc.Storage.GetSurveySessionsWithAnswers(
		survey.UUID,
		&types.SurveySessionsFilter{
			Limit:  analyticsMaxSessions,
			Offset: 0,
			SortBy: "created_at",
			Order:  "desc",
		},
	)
	if err != nil {
		logCtx.Error("unable to get survey sessions", "err", err)
		return nil, err
	}

	for i := range sessions {
		if sessions[i].Status != types.SurveySessionStatus_Completed {
			continue
		}

		sessions[i].QuestionAnswers = convertAnswerBytesToAnswerType(
			svc,
			&survey,
			sessions[i].QuestionAnswers,
		)
		if err := filterVisibleAnswers(&survey, &sessions[i]); err != nil {
			return nil, err
		}
	}

	return BuildSurveyAnalytics(&survey, sessions), nil
}

// BuildSurveyAnalytics aggregates a survey and its sessions into chart-ready data.
func BuildSurveyAnalytics(survey *types.Survey, sessions []types.SurveySession) *types.SurveyAnalytics {
	analytics := &types.SurveyAnalytics{
		SurveyUUID:  survey.UUID,
		SurveyName:  survey.Name,
		GeneratedAt: time.Now().UTC(),
		Timeline:    []types.AnalyticsTimelinePoint{},
		Questions:   []types.QuestionAnalytics{},
	}
	if survey.Config != nil {
		analytics.SurveyTitle = survey.Config.Title
	}

	completedSessions := make([]types.SurveySession, 0, len(sessions))
	timelineByDate := map[string]*types.AnalyticsTimelinePoint{}
	completedWithDuration := 0
	var totalCompletionSeconds int64

	getTimelinePoint := func(date string) *types.AnalyticsTimelinePoint {
		point, ok := timelineByDate[date]
		if !ok {
			point = &types.AnalyticsTimelinePoint{Date: date}
			timelineByDate[date] = point
		}
		return point
	}

	for i := range sessions {
		session := sessions[i]
		analytics.Totals.SessionsTotal++

		getTimelinePoint(session.CreatedAt.Format(types.DATE_FORMAT)).Started++

		if session.Status == types.SurveySessionStatus_Completed {
			analytics.Totals.SessionsCompleted++
			completedSessions = append(completedSessions, session)

			if session.CompletedAt != nil {
				getTimelinePoint(session.CompletedAt.Format(types.DATE_FORMAT)).Completed++

				duration := session.CompletedAt.Sub(session.CreatedAt)
				if duration > 0 {
					totalCompletionSeconds += int64(duration.Seconds())
					completedWithDuration++
				}
			}
		} else {
			analytics.Totals.SessionsInProgress++
		}
	}

	if analytics.Totals.SessionsTotal > 0 {
		analytics.Totals.CompletionRate = roundTo(
			float64(analytics.Totals.SessionsCompleted)/float64(analytics.Totals.SessionsTotal)*100,
			1,
		)
	}

	if completedWithDuration > 0 {
		average := totalCompletionSeconds / int64(completedWithDuration)
		analytics.Totals.AverageCompletionSeconds = &average
	}

	dates := make([]string, 0, len(timelineByDate))
	for date := range timelineByDate {
		dates = append(dates, date)
	}
	sort.Strings(dates)
	for _, date := range dates {
		analytics.Timeline = append(analytics.Timeline, *timelineByDate[date])
	}

	if survey.Config == nil || survey.Config.Questions == nil {
		return analytics
	}

	for _, question := range survey.Config.Questions.Questions {
		chartKind, supported := chartKindForQuestionType(question.Type)
		if !supported {
			continue
		}

		questionAnalytics := types.QuestionAnalytics{
			QuestionUUID: question.UUID,
			QuestionID:   question.ID,
			Label:        question.Label,
			Type:         question.Type,
			ChartKind:    chartKind,
			Conditional:  question.VisibleIf != nil,
		}

		switch question.Type {
		case types.QuestionType_DropdownSingle:
			analyzeSingleChoice(question, completedSessions, &questionAnalytics)
		case types.QuestionType_DropdownMultiple:
			analyzeMultipleChoice(question, completedSessions, &questionAnalytics)
		case types.QuestionType_YesNo:
			analyzeBoolean(question, completedSessions, &questionAnalytics)
		case types.QuestionType_Rating:
			analyzeRating(question, completedSessions, &questionAnalytics)
		case types.QuestionType_Ranking:
			analyzeRanking(question, completedSessions, &questionAnalytics)
		case types.QuestionType_Date:
			analyzeDate(question, completedSessions, &questionAnalytics)
		}

		analytics.Questions = append(analytics.Questions, questionAnalytics)
	}

	return analytics
}

func chartKindForQuestionType(questionType types.QuestionType) (types.AnalyticsChartKind, bool) {
	switch questionType {
	case types.QuestionType_DropdownSingle:
		return types.AnalyticsChartKind_SingleChoice, true
	case types.QuestionType_DropdownMultiple:
		return types.AnalyticsChartKind_MultipleChoice, true
	case types.QuestionType_YesNo:
		return types.AnalyticsChartKind_Boolean, true
	case types.QuestionType_Rating:
		return types.AnalyticsChartKind_Rating, true
	case types.QuestionType_Ranking:
		return types.AnalyticsChartKind_Ranking, true
	case types.QuestionType_Date:
		return types.AnalyticsChartKind_Date, true
	default:
		return "", false
	}
}

func analyzeSingleChoice(
	question types.Question,
	sessions []types.SurveySession,
	analytics *types.QuestionAnalytics,
) {
	order, seen := choiceOrder(question.Options)
	counts := make(map[string]int, len(question.Options))

	for i := range sessions {
		value, ok := answerStringValue(sessions[i], question.UUID)
		if !ok {
			continue
		}

		analytics.AnsweredCount++
		if !seen[value] {
			order = append(order, value)
			seen[value] = true
		}
		counts[value]++
	}

	analytics.Buckets = buildChoiceBuckets(order, counts, analytics.AnsweredCount)
}

func analyzeMultipleChoice(
	question types.Question,
	sessions []types.SurveySession,
	analytics *types.QuestionAnalytics,
) {
	order, seen := choiceOrder(question.Options)
	counts := make(map[string]int, len(question.Options))

	for i := range sessions {
		values, ok := answerStringSliceValue(sessions[i], question.UUID)
		if !ok {
			continue
		}

		analytics.AnsweredCount++
		unique := make(map[string]bool, len(values))
		for _, value := range values {
			if unique[value] {
				continue
			}
			unique[value] = true

			if !seen[value] {
				order = append(order, value)
				seen[value] = true
			}
			counts[value]++
		}
	}

	analytics.Buckets = buildChoiceBuckets(order, counts, analytics.AnsweredCount)
}

func analyzeBoolean(
	question types.Question,
	sessions []types.SurveySession,
	analytics *types.QuestionAnalytics,
) {
	counts := map[bool]int{}
	for i := range sessions {
		value, ok := answerBoolValue(sessions[i], question.UUID)
		if !ok {
			continue
		}

		analytics.AnsweredCount++
		counts[value]++
	}

	analytics.Buckets = []types.AnalyticsBucket{
		newBucket("true", "true", counts[true], analytics.AnsweredCount),
		newBucket("false", "false", counts[false], analytics.AnsweredCount),
	}
}

func analyzeRating(
	question types.Question,
	sessions []types.SurveySession,
	analytics *types.QuestionAnalytics,
) {
	counts := map[int]int{}
	values := []int{}
	var sum float64

	for i := range sessions {
		value, ok := answerNumberValue(sessions[i], question.UUID)
		if !ok {
			continue
		}

		rating := int(math.Round(value))
		counts[rating]++
		values = append(values, rating)
		sum += float64(rating)
	}

	analytics.AnsweredCount = len(values)
	if len(values) == 0 {
		return
	}

	minValue, maxValue := ratingRange(question, values)
	span := int64(maxValue) - int64(minValue) + 1
	useObservedOnly := span <= 0 || span > maxRatingBuckets
	bucketCount := int(span)
	if useObservedOnly {
		bucketCount = len(counts)
	}
	buckets := make([]types.AnalyticsBucket, 0, bucketCount)
	if useObservedOnly {
		observed := make([]int, 0, len(counts))
		for value := range counts {
			observed = append(observed, value)
		}
		sort.Ints(observed)

		for _, value := range observed {
			label := strconv.Itoa(value)
			buckets = append(buckets, newBucket(label, label, counts[value], analytics.AnsweredCount))
		}
	} else {
		for value := int64(minValue); value <= int64(maxValue); value++ {
			label := strconv.FormatInt(value, 10)
			buckets = append(buckets, newBucket(label, label, counts[int(value)], analytics.AnsweredCount))
		}
	}

	sort.Ints(values)
	average := roundTo(float64(sum)/float64(len(values)), 2)
	median := medianOfSortedInts(values)
	mode := modeOfInts(counts)

	analytics.Buckets = buckets
	analytics.Numeric = &types.AnalyticsNumericSummary{
		Average: &average,
		Median:  &median,
		Mode:    &mode,
	}
}

func analyzeRanking(
	question types.Question,
	sessions []types.SurveySession,
	analytics *types.QuestionAnalytics,
) {
	order, seen := choiceOrder(question.Options)
	rankCounts := map[string][]int{}
	rankTotals := map[string]int{}
	rankedCounts := map[string]int{}

	for i := range sessions {
		values, ok := answerStringSliceValue(sessions[i], question.UUID)
		if !ok {
			continue
		}

		analytics.AnsweredCount++
		for position, value := range values {
			if !seen[value] {
				order = append(order, value)
				seen[value] = true
			}

			for len(rankCounts[value]) <= position {
				rankCounts[value] = append(rankCounts[value], 0)
			}
			rankCounts[value][position]++
			rankTotals[value] += position + 1
			rankedCounts[value]++
		}
	}

	maxRanks := len(question.Options)
	for _, counts := range rankCounts {
		if len(counts) > maxRanks {
			maxRanks = len(counts)
		}
	}

	ranking := make([]types.AnalyticsRankingOption, 0, len(order))
	for _, option := range order {
		counts := append([]int(nil), rankCounts[option]...)
		for len(counts) < maxRanks {
			counts = append(counts, 0)
		}

		var averageRank *float64
		if rankedCounts[option] > 0 {
			average := roundTo(
				float64(rankTotals[option])/float64(rankedCounts[option]),
				2,
			)
			averageRank = &average
		}

		ranking = append(ranking, types.AnalyticsRankingOption{
			Key:         option,
			Label:       option,
			RankCounts:  counts,
			AverageRank: averageRank,
		})
	}

	analytics.Ranking = ranking
}

func analyzeDate(
	question types.Question,
	sessions []types.SurveySession,
	analytics *types.QuestionAnalytics,
) {
	dates := []time.Time{}
	for i := range sessions {
		value, ok := answerStringValue(sessions[i], question.UUID)
		if !ok {
			continue
		}

		date, err := time.Parse(types.DATE_FORMAT, value)
		if err != nil {
			continue
		}
		dates = append(dates, date)
	}

	analytics.AnsweredCount = len(dates)
	if len(dates) == 0 {
		return
	}

	granularity := dateGranularity(dates)
	analytics.Granularity = granularity

	counts := map[string]int{}
	for _, date := range dates {
		counts[formatDateBucket(date, granularity)]++
	}

	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	buckets := make([]types.AnalyticsBucket, 0, len(keys))
	for _, key := range keys {
		buckets = append(buckets, newBucket(key, key, counts[key], analytics.AnsweredCount))
	}
	analytics.Buckets = buckets
}

func choiceOrder(options []string) ([]string, map[string]bool) {
	order := append([]string(nil), options...)
	seen := make(map[string]bool, len(options))
	for _, option := range options {
		seen[option] = true
	}
	return order, seen
}

func buildChoiceBuckets(
	order []string,
	counts map[string]int,
	answeredCount int,
) []types.AnalyticsBucket {
	buckets := make([]types.AnalyticsBucket, 0, len(order))
	for _, option := range order {
		buckets = append(buckets, newBucket(option, option, counts[option], answeredCount))
	}
	return buckets
}

func newBucket(key string, label string, count int, answeredCount int) types.AnalyticsBucket {
	return types.AnalyticsBucket{
		Key:        key,
		Label:      label,
		Count:      count,
		Percentage: percentage(count, answeredCount),
	}
}

func percentage(count int, total int) float64 {
	if total == 0 {
		return 0
	}
	return roundTo(float64(count)/float64(total)*100, 1)
}

func findAnswer(answers []types.QuestionAnswer, questionUUID string) (types.QuestionAnswer, bool) {
	for i := range answers {
		if answers[i].QuestionUUID == questionUUID {
			return answers[i], true
		}
	}
	return types.QuestionAnswer{}, false
}

func answerValueFor(session types.SurveySession, questionUUID string) (any, bool) {
	answer, ok := findAnswer(session.QuestionAnswers, questionUUID)
	if !ok || answer.Answer == nil {
		return nil, false
	}

	value, err := types.AnswerValue(answer.Answer)
	if err != nil || isEmptyAnswer(value) {
		return nil, false
	}
	return value, true
}

func answerStringValue(session types.SurveySession, questionUUID string) (string, bool) {
	value, ok := answerValueFor(session, questionUUID)
	if !ok {
		return "", false
	}

	text, ok := value.(string)
	if !ok || strings.TrimSpace(text) == "" {
		return "", false
	}
	return text, true
}

func answerStringSliceValue(session types.SurveySession, questionUUID string) ([]string, bool) {
	value, ok := answerValueFor(session, questionUUID)
	if !ok {
		return nil, false
	}

	switch typed := value.(type) {
	case []string:
		values := make([]string, 0, len(typed))
		for _, item := range typed {
			text := strings.TrimSpace(item)
			if text != "" {
				values = append(values, text)
			}
		}
		return values, len(values) > 0
	case []any:
		values := make([]string, 0, len(typed))
		for _, item := range typed {
			if item == nil {
				continue
			}
			text := strings.TrimSpace(stringValue(item))
			if text == "" {
				continue
			}
			values = append(values, text)
		}
		return values, len(values) > 0
	default:
		return nil, false
	}
}

func answerBoolValue(session types.SurveySession, questionUUID string) (bool, bool) {
	value, ok := answerValueFor(session, questionUUID)
	if !ok {
		return false, false
	}

	switch typed := value.(type) {
	case bool:
		return typed, true
	case string:
		parsed, err := strconv.ParseBool(strings.TrimSpace(typed))
		if err != nil {
			return false, false
		}
		return parsed, true
	default:
		return false, false
	}
}

func answerNumberValue(session types.SurveySession, questionUUID string) (float64, bool) {
	value, ok := answerValueFor(session, questionUUID)
	if !ok {
		return 0, false
	}

	switch typed := value.(type) {
	case float64:
		return typed, true
	case float32:
		return float64(typed), true
	case int:
		return float64(typed), true
	case int64:
		return float64(typed), true
	case string:
		parsed, err := strconv.ParseFloat(strings.TrimSpace(typed), 64)
		if err != nil {
			return 0, false
		}
		return parsed, true
	default:
		return 0, false
	}
}

func stringValue(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case float64:
		return strconv.FormatFloat(typed, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(typed)
	default:
		return ""
	}
}

func isEmptyAnswer(value any) bool {
	switch typed := value.(type) {
	case nil:
		return true
	case string:
		return strings.TrimSpace(typed) == ""
	case []string:
		return len(typed) == 0
	case []any:
		return len(typed) == 0
	default:
		return false
	}
}

func ratingRange(question types.Question, values []int) (int, int) {
	minValue, maxValue := values[0], values[0]
	for _, value := range values[1:] {
		if value < minValue {
			minValue = value
		}
		if value > maxValue {
			maxValue = value
		}
	}

	if question.Min != nil && *question.Min < minValue {
		minValue = *question.Min
	}
	if question.Max != nil && *question.Max > maxValue {
		maxValue = *question.Max
	}

	return minValue, maxValue
}

func medianOfSortedInts(values []int) float64 {
	length := len(values)
	if length == 0 {
		return 0
	}
	middle := length / 2
	if length%2 == 1 {
		return float64(values[middle])
	}
	return float64(values[middle-1]+values[middle]) / 2
}

func modeOfInts(counts map[int]int) float64 {
	mode := 0
	modeCount := -1
	for value, count := range counts {
		if count > modeCount || (count == modeCount && value < mode) {
			mode = value
			modeCount = count
		}
	}
	return float64(mode)
}

func dateGranularity(dates []time.Time) string {
	minDate, maxDate := dates[0], dates[0]
	for _, date := range dates[1:] {
		if date.Before(minDate) {
			minDate = date
		}
		if date.After(maxDate) {
			maxDate = date
		}
	}

	spanDays := maxDate.Sub(minDate).Hours() / 24
	if spanDays <= 45 {
		return "day"
	}

	months := (maxDate.Year()-minDate.Year())*12 + int(maxDate.Month()) - int(minDate.Month())
	if months <= 36 {
		return "month"
	}

	return "year"
}

func formatDateBucket(date time.Time, granularity string) string {
	switch granularity {
	case "year":
		return date.Format("2006")
	case "month":
		return date.Format("2006-01")
	default:
		return date.Format(types.DATE_FORMAT)
	}
}

func roundTo(value float64, places int) float64 {
	factor := math.Pow(10, float64(places))
	return math.Round(value*factor) / factor
}
