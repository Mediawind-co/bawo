package streak

import (
	"context"
	"time"

	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
	"encore.dev/cron"
	"encore.dev/pubsub"
	"github.com/google/uuid"

	"encore.app/learning/lesson"
)

// ========== Pub/Sub Subscriber ==========

// Subscribe to lesson completed events to update streaks.
var _ = pubsub.NewSubscription(
	lesson.LessonCompleted, "update-streak",
	pubsub.SubscriptionConfig[*lesson.LessonCompletedEvent]{
		Handler: HandleStreakUpdate,
	},
)

// HandleStreakUpdate processes lesson completion for streak updates.
func HandleStreakUpdate(ctx context.Context, event *lesson.LessonCompletedEvent) error {
	userID, err := uuid.Parse(event.UserID)
	if err != nil {
		return nil
	}

	// Get or create streak
	userStreak, err := GetOrCreateStreak(ctx, userID)
	if err != nil {
		return err
	}

	today := time.Now().Truncate(24 * time.Hour)

	// Update daily activity
	activity, err := UpdateDailyActivity(ctx, userID, today, event.XPEarned, userStreak.DailyXPGoal)
	if err != nil {
		return err
	}

	// If goal met and this is a new streak day, increment streak
	if activity.GoalMet {
		// Check if this is a new day for the streak
		isNewStreakDay := userStreak.LastActiveDate == nil ||
			userStreak.LastActiveDate.Truncate(24*time.Hour).Before(today)

		if isNewStreakDay {
			_, err = IncrementStreak(ctx, userID, today)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ========== Cron Job ==========

// Check and reset broken streaks daily at 2 AM.
var _ = cron.NewJob("check-streaks", cron.JobConfig{
	Title:    "Check and reset broken streaks",
	Schedule: "0 2 * * *",
	Endpoint: CheckStreaks,
})

// CheckStreaks checks for broken streaks and resets them.
//
//encore:api private
func CheckStreaks(ctx context.Context) error {
	yesterday := time.Now().AddDate(0, 0, -1)

	// Find users with broken streaks
	userIDs, err := GetUsersWithBrokenStreaks(ctx, yesterday)
	if err != nil {
		return err
	}

	for _, userID := range userIDs {
		// Try to use a streak freeze first
		used, err := UseStreakFreeze(ctx, userID)
		if err != nil {
			continue
		}

		if !used {
			// No freeze available, reset streak
			_ = ResetStreak(ctx, userID)
		}
	}

	return nil
}

// ========== API Endpoints ==========

// GetMyStreak returns the current user's streak information.
//
//encore:api auth method=GET path=/streak
func GetMyStreak(ctx context.Context) (*StreakResponse, error) {
	uid, ok := auth.UserID()
	if !ok {
		return nil, &errs.Error{Code: errs.Unauthenticated, Message: "not authenticated"}
	}

	userID, err := uuid.Parse(string(uid))
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid user ID"}
	}

	// Get or create streak
	userStreak, err := GetOrCreateStreak(ctx, userID)
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to get streak"}
	}

	// Get today's activity
	today := time.Now().Truncate(24 * time.Hour)
	activity, _ := GetDailyActivity(ctx, userID, today)

	todayXP := 0
	todayLessons := 0
	goalMet := false
	if activity != nil {
		todayXP = activity.XPEarned
		todayLessons = activity.LessonsCompleted
		goalMet = activity.GoalMet
	}

	xpToGoal := userStreak.DailyXPGoal - todayXP
	if xpToGoal < 0 {
		xpToGoal = 0
	}

	// Check if streak is at risk (close to end of day without meeting goal)
	streakAtRisk := !goalMet && userStreak.CurrentStreak > 0

	return &StreakResponse{
		Streak:       userStreak,
		TodayXP:      todayXP,
		TodayLessons: todayLessons,
		GoalMet:      goalMet,
		XPToGoal:     xpToGoal,
		StreakAtRisk: streakAtRisk,
	}, nil
}

// GetDailyGoal returns today's progress toward the daily goal.
//
//encore:api auth method=GET path=/daily-goal
func GetDailyGoal(ctx context.Context) (*DailyGoalResponse, error) {
	uid, ok := auth.UserID()
	if !ok {
		return nil, &errs.Error{Code: errs.Unauthenticated, Message: "not authenticated"}
	}

	userID, err := uuid.Parse(string(uid))
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid user ID"}
	}

	userStreak, err := GetOrCreateStreak(ctx, userID)
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to get streak"}
	}

	today := time.Now().Truncate(24 * time.Hour)
	activity, _ := GetDailyActivity(ctx, userID, today)

	todayXP := 0
	todayLessons := 0
	goalMet := false
	if activity != nil {
		todayXP = activity.XPEarned
		todayLessons = activity.LessonsCompleted
		goalMet = activity.GoalMet
	}

	xpRemaining := userStreak.DailyXPGoal - todayXP
	if xpRemaining < 0 {
		xpRemaining = 0
	}

	percentComplete := 0
	if userStreak.DailyXPGoal > 0 {
		percentComplete = (todayXP * 100) / userStreak.DailyXPGoal
		if percentComplete > 100 {
			percentComplete = 100
		}
	}

	return &DailyGoalResponse{
		DailyXPGoal:     userStreak.DailyXPGoal,
		TodayXP:         todayXP,
		TodayLessons:    todayLessons,
		GoalMet:         goalMet,
		XPRemaining:     xpRemaining,
		PercentComplete: percentComplete,
	}, nil
}

// UpdateDailyGoal updates the user's daily XP goal.
//
//encore:api auth method=PUT path=/daily-goal
func UpdateDailyGoal(ctx context.Context, req *UpdateGoalRequest) (*StreakResponse, error) {
	uid, ok := auth.UserID()
	if !ok {
		return nil, &errs.Error{Code: errs.Unauthenticated, Message: "not authenticated"}
	}

	userID, err := uuid.Parse(string(uid))
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid user ID"}
	}

	if req.DailyXPGoal < 10 || req.DailyXPGoal > 500 {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "daily goal must be between 10 and 500 XP"}
	}

	// Ensure streak exists
	_, err = GetOrCreateStreak(ctx, userID)
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to get streak"}
	}

	// Update goal
	_, err = UpdateStreakGoal(ctx, userID, req.DailyXPGoal)
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to update goal"}
	}

	// Return updated streak info
	return GetMyStreak(ctx)
}

// GetActivityHistory returns recent daily activity.
//
//encore:api auth method=GET path=/activity
func GetActivityHistory(ctx context.Context) (*ActivityHistoryResponse, error) {
	uid, ok := auth.UserID()
	if !ok {
		return nil, &errs.Error{Code: errs.Unauthenticated, Message: "not authenticated"}
	}

	userID, err := uuid.Parse(string(uid))
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid user ID"}
	}

	// Get last 30 days of activity
	activities, err := ListRecentActivity(ctx, userID, 30)
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to get activity"}
	}

	return &ActivityHistoryResponse{
		Activities: activities,
		Days:       30,
	}, nil
}
