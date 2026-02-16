package dashboard

import (
	"context"
	"time"

	"encore.dev/beta/errs"

	"encore.app/identity/user"
	"encore.app/learning/content"
	"encore.app/learning/language"
	"encore.app/progress/enrollment"
)

// ========== Admin Endpoints ==========

// GetOverview returns high-level platform statistics.
//
//encore:api auth method=GET path=/admin/analytics/overview tag:admin
func GetOverview(ctx context.Context) (*OverviewResponse, error) {
	stats := &OverviewStats{}

	// Get user count
	users, _, err := user.List(ctx, 1, 0)
	if err == nil {
		// List returns limited results, we need total
		_, total, _ := user.List(ctx, 1, 0)
		stats.TotalUsers = total
		_ = users // avoid unused
	}

	// Get language count
	languages, err := language.List(ctx, true)
	if err == nil {
		stats.TotalLanguages = len(languages)
	}

	return &OverviewResponse{
		Stats:     stats,
		UpdatedAt: time.Now(),
	}, nil
}

// GetUserAnalytics returns user-related statistics.
//
//encore:api auth method=GET path=/admin/analytics/users tag:admin
func GetUserAnalytics(ctx context.Context) (*UserStatsResponse, error) {
	stats := &UserStats{
		UsersByProvider: make(map[string]int),
	}

	// Get users with pagination to count
	_, total, err := user.List(ctx, 1, 0)
	if err == nil {
		stats.TotalUsers = total
	}

	// Get all users to count by provider
	users, _, err := user.List(ctx, 1000, 0)
	if err == nil {
		for _, u := range users {
			stats.UsersByProvider[string(u.Provider)]++
		}
	}

	stats.SignupTrend = []DailyCount{} // Would need custom query

	return &UserStatsResponse{Stats: stats}, nil
}

// GetEnrollmentAnalytics returns enrollment statistics.
//
//encore:api auth method=GET path=/admin/analytics/enrollments tag:admin
func GetEnrollmentAnalytics(ctx context.Context) (*EnrollmentStatsResponse, error) {
	stats := &EnrollmentStats{
		EnrollmentsByLang: []LanguageEnrollment{},
		EnrollmentTrend:   []DailyCount{},
	}

	// Get all languages and count enrollments per language
	languages, err := language.List(ctx, true)
	if err == nil {
		for _, lang := range languages {
			count, err := enrollment.CountByLanguage(ctx, lang.ID)
			if err == nil && count > 0 {
				stats.EnrollmentsByLang = append(stats.EnrollmentsByLang, LanguageEnrollment{
					LanguageID:   lang.ID.String(),
					LanguageName: lang.Name,
					LanguageCode: lang.Code,
					Count:        count,
				})
				stats.TotalEnrollments += count
			}
		}
		stats.ActiveEnrollments = stats.TotalEnrollments
	}

	return &EnrollmentStatsResponse{Stats: stats}, nil
}

// GetCompletionAnalytics returns lesson completion statistics.
//
//encore:api auth method=GET path=/admin/analytics/completion tag:admin
func GetCompletionAnalytics(ctx context.Context) (*CompletionStatsResponse, error) {
	stats := &CompletionStats{
		CompletionTrend: []DailyCount{},
	}

	// These would need direct DB access or dedicated endpoints
	// For now return placeholder data
	stats.AverageScore = 0
	stats.CompletionRate = 0

	return &CompletionStatsResponse{Stats: stats}, nil
}

// GetStreakAnalytics returns streak-related statistics.
//
//encore:api auth method=GET path=/admin/analytics/streaks tag:admin
func GetStreakAnalytics(ctx context.Context) (*StreakStatsResponse, error) {
	stats := &StreakStats{}

	// These would need direct DB access or dedicated endpoints
	// For now return placeholder data

	return &StreakStatsResponse{Stats: stats}, nil
}

// GetContentAnalytics returns content statistics.
//
//encore:api auth method=GET path=/admin/analytics/content tag:admin
func GetContentAnalytics(ctx context.Context) (*ContentStatsResponse, error) {
	stats := &ContentStats{}

	// Get languages
	languages, err := language.List(ctx, true)
	if err == nil {
		stats.TotalLanguages = len(languages)
		for _, l := range languages {
			if l.IsActive {
				stats.ActiveLanguages++
			}
		}
	}

	// Count units and lessons per language
	for _, lang := range languages {
		units, err := content.ListUnitsByLanguage(ctx, lang.ID)
		if err == nil {
			stats.TotalUnits += len(units)

			for _, unit := range units {
				lessons, err := content.ListLessonsByUnit(ctx, unit.ID)
				if err == nil {
					stats.TotalLessons += len(lessons)

					for _, lesson := range lessons {
						questions, err := content.ListQuestionsByLesson(ctx, lesson.ID)
						if err == nil {
							stats.TotalQuestions += len(questions)
							for _, q := range questions {
								if q.PromptAudioKey != "" {
									stats.QuestionsWithAudio++
								}
							}
						}
					}
				}
			}
		}
	}

	return &ContentStatsResponse{Stats: stats}, nil
}

// ========== Health Check ==========

// Ping is a simple health check endpoint.
//
//encore:api public method=GET path=/admin/health
func Ping(ctx context.Context) (*PingResponse, error) {
	return &PingResponse{Status: "ok", Timestamp: time.Now()}, nil
}

// PingResponse contains health check result.
type PingResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

// ========== Error Handling ==========

func internalError(msg string) error {
	return &errs.Error{Code: errs.Internal, Message: msg}
}
