package dashboard

import "time"

// OverviewStats contains high-level platform statistics.
type OverviewStats struct {
	TotalUsers        int `json:"total_users"`
	TotalLanguages    int `json:"total_languages"`
	ActiveUsersToday  int `json:"active_users_today"`
	ActiveUsersWeek   int `json:"active_users_week"`
	ActiveUsersMonth  int `json:"active_users_month"`
	TotalEnrollments  int `json:"total_enrollments"`
	TotalLessons      int `json:"total_lessons"`
	TotalQuestions    int `json:"total_questions"`
	LessonsCompleted  int `json:"lessons_completed_today"`
}

// UserStats contains user-related statistics.
type UserStats struct {
	TotalUsers       int              `json:"total_users"`
	NewUsersToday    int              `json:"new_users_today"`
	NewUsersWeek     int              `json:"new_users_week"`
	NewUsersMonth    int              `json:"new_users_month"`
	UsersByProvider  map[string]int   `json:"users_by_provider"`
	SignupTrend      []DailyCount     `json:"signup_trend"`
}

// EnrollmentStats contains enrollment statistics.
type EnrollmentStats struct {
	TotalEnrollments    int                   `json:"total_enrollments"`
	ActiveEnrollments   int                   `json:"active_enrollments"`
	EnrollmentsByLang   []LanguageEnrollment  `json:"enrollments_by_language"`
	EnrollmentTrend     []DailyCount          `json:"enrollment_trend"`
}

// LanguageEnrollment represents enrollments for a language.
type LanguageEnrollment struct {
	LanguageID   string `json:"language_id"`
	LanguageName string `json:"language_name"`
	LanguageCode string `json:"language_code"`
	Count        int    `json:"count"`
}

// CompletionStats contains lesson completion statistics.
type CompletionStats struct {
	TotalCompletions   int              `json:"total_completions"`
	CompletionsToday   int              `json:"completions_today"`
	CompletionsWeek    int              `json:"completions_week"`
	AverageScore       float64          `json:"average_score"`
	CompletionRate     float64          `json:"completion_rate"`
	CompletionTrend    []DailyCount     `json:"completion_trend"`
}

// StreakStats contains streak-related statistics.
type StreakStats struct {
	UsersWithStreak    int     `json:"users_with_streak"`
	AverageStreak      float64 `json:"average_streak"`
	LongestStreak      int     `json:"longest_streak"`
	GoalCompletionRate float64 `json:"goal_completion_rate"`
}

// ContentStats contains content statistics.
type ContentStats struct {
	TotalLanguages  int `json:"total_languages"`
	ActiveLanguages int `json:"active_languages"`
	TotalUnits      int `json:"total_units"`
	TotalLessons    int `json:"total_lessons"`
	TotalQuestions  int `json:"total_questions"`
	QuestionsWithAudio int `json:"questions_with_audio"`
}

// DailyCount represents a count for a specific day.
type DailyCount struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}

// ========== Response Types ==========

// OverviewResponse wraps overview stats.
type OverviewResponse struct {
	Stats     *OverviewStats `json:"stats"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// UserStatsResponse wraps user stats.
type UserStatsResponse struct {
	Stats *UserStats `json:"stats"`
}

// EnrollmentStatsResponse wraps enrollment stats.
type EnrollmentStatsResponse struct {
	Stats *EnrollmentStats `json:"stats"`
}

// CompletionStatsResponse wraps completion stats.
type CompletionStatsResponse struct {
	Stats *CompletionStats `json:"stats"`
}

// StreakStatsResponse wraps streak stats.
type StreakStatsResponse struct {
	Stats *StreakStats `json:"stats"`
}

// ContentStatsResponse wraps content stats.
type ContentStatsResponse struct {
	Stats *ContentStats `json:"stats"`
}
