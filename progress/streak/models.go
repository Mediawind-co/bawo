package streak

import (
	"time"

	"github.com/google/uuid"
)

// Streak tracks a user's daily learning streak.
type Streak struct {
	ID             uuid.UUID  `json:"id"`
	UserID         uuid.UUID  `json:"user_id"`
	CurrentStreak  int        `json:"current_streak"`
	LongestStreak  int        `json:"longest_streak"`
	LastActiveDate *time.Time `json:"last_active_date,omitempty"`
	DailyXPGoal    int        `json:"daily_xp_goal"`
	StreakFreezes  int        `json:"streak_freezes"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// DailyActivity tracks a user's activity for a specific day.
type DailyActivity struct {
	ID               uuid.UUID `json:"id"`
	UserID           uuid.UUID `json:"user_id"`
	Date             time.Time `json:"date"`
	XPEarned         int       `json:"xp_earned"`
	LessonsCompleted int       `json:"lessons_completed"`
	GoalMet          bool      `json:"goal_met"`
	CreatedAt        time.Time `json:"created_at"`
}

// ========== Response Types ==========

// StreakResponse wraps streak data.
type StreakResponse struct {
	Streak         *Streak `json:"streak"`
	TodayXP        int     `json:"today_xp"`
	TodayLessons   int     `json:"today_lessons"`
	GoalMet        bool    `json:"goal_met"`
	XPToGoal       int     `json:"xp_to_goal"`
	StreakAtRisk   bool    `json:"streak_at_risk"`
}

// DailyGoalResponse contains today's progress.
type DailyGoalResponse struct {
	DailyXPGoal      int  `json:"daily_xp_goal"`
	TodayXP          int  `json:"today_xp"`
	TodayLessons     int  `json:"today_lessons"`
	GoalMet          bool `json:"goal_met"`
	XPRemaining      int  `json:"xp_remaining"`
	PercentComplete  int  `json:"percent_complete"`
}

// ActivityHistoryResponse contains recent activity.
type ActivityHistoryResponse struct {
	Activities []*DailyActivity `json:"activities"`
	Days       int              `json:"days"`
}

// UpdateGoalRequest contains new goal settings.
type UpdateGoalRequest struct {
	DailyXPGoal int `json:"daily_xp_goal"`
}
