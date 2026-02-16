package streak

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrStreakNotFound = errors.New("streak not found")
)

// ========== Streak Repository ==========

// GetOrCreateStreak gets or creates a streak for a user.
func GetOrCreateStreak(ctx context.Context, userID uuid.UUID) (*Streak, error) {
	var streak Streak
	err := db.QueryRow(ctx, `
		INSERT INTO streaks (user_id)
		VALUES ($1)
		ON CONFLICT (user_id) DO UPDATE SET updated_at = NOW()
		RETURNING id, user_id, current_streak, longest_streak, last_active_date, daily_xp_goal, streak_freezes, updated_at
	`, userID).Scan(
		&streak.ID, &streak.UserID, &streak.CurrentStreak, &streak.LongestStreak,
		&streak.LastActiveDate, &streak.DailyXPGoal, &streak.StreakFreezes, &streak.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &streak, nil
}

// FindStreakByUser finds a streak by user ID.
func FindStreakByUser(ctx context.Context, userID uuid.UUID) (*Streak, error) {
	var streak Streak
	err := db.QueryRow(ctx, `
		SELECT id, user_id, current_streak, longest_streak, last_active_date, daily_xp_goal, streak_freezes, updated_at
		FROM streaks
		WHERE user_id = $1
	`, userID).Scan(
		&streak.ID, &streak.UserID, &streak.CurrentStreak, &streak.LongestStreak,
		&streak.LastActiveDate, &streak.DailyXPGoal, &streak.StreakFreezes, &streak.UpdatedAt,
	)
	if err != nil {
		return nil, ErrStreakNotFound
	}
	return &streak, nil
}

// UpdateStreakGoal updates the daily XP goal.
func UpdateStreakGoal(ctx context.Context, userID uuid.UUID, goal int) (*Streak, error) {
	var streak Streak
	err := db.QueryRow(ctx, `
		UPDATE streaks
		SET daily_xp_goal = $2, updated_at = NOW()
		WHERE user_id = $1
		RETURNING id, user_id, current_streak, longest_streak, last_active_date, daily_xp_goal, streak_freezes, updated_at
	`, userID, goal).Scan(
		&streak.ID, &streak.UserID, &streak.CurrentStreak, &streak.LongestStreak,
		&streak.LastActiveDate, &streak.DailyXPGoal, &streak.StreakFreezes, &streak.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &streak, nil
}

// IncrementStreak increments the streak and updates last active date.
func IncrementStreak(ctx context.Context, userID uuid.UUID, date time.Time) (*Streak, error) {
	var streak Streak
	err := db.QueryRow(ctx, `
		UPDATE streaks
		SET
			current_streak = current_streak + 1,
			longest_streak = GREATEST(longest_streak, current_streak + 1),
			last_active_date = $2,
			updated_at = NOW()
		WHERE user_id = $1
		RETURNING id, user_id, current_streak, longest_streak, last_active_date, daily_xp_goal, streak_freezes, updated_at
	`, userID, date).Scan(
		&streak.ID, &streak.UserID, &streak.CurrentStreak, &streak.LongestStreak,
		&streak.LastActiveDate, &streak.DailyXPGoal, &streak.StreakFreezes, &streak.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &streak, nil
}

// ResetStreak resets a user's streak to 0.
func ResetStreak(ctx context.Context, userID uuid.UUID) error {
	_, err := db.Exec(ctx, `
		UPDATE streaks
		SET current_streak = 0, updated_at = NOW()
		WHERE user_id = $1
	`, userID)
	return err
}

// UseStreakFreeze uses one streak freeze to prevent reset.
func UseStreakFreeze(ctx context.Context, userID uuid.UUID) (bool, error) {
	result, err := db.Exec(ctx, `
		UPDATE streaks
		SET streak_freezes = streak_freezes - 1, updated_at = NOW()
		WHERE user_id = $1 AND streak_freezes > 0
	`, userID)
	if err != nil {
		return false, err
	}
	return result.RowsAffected() > 0, nil
}

// AddStreakFreeze adds streak freezes to a user.
func AddStreakFreeze(ctx context.Context, userID uuid.UUID, count int) error {
	_, err := db.Exec(ctx, `
		UPDATE streaks
		SET streak_freezes = streak_freezes + $2, updated_at = NOW()
		WHERE user_id = $1
	`, userID, count)
	return err
}

// ========== Daily Activity Repository ==========

// GetOrCreateDailyActivity gets or creates activity for a date.
func GetOrCreateDailyActivity(ctx context.Context, userID uuid.UUID, date time.Time) (*DailyActivity, error) {
	dateOnly := date.Truncate(24 * time.Hour)
	var activity DailyActivity
	err := db.QueryRow(ctx, `
		INSERT INTO daily_activity (user_id, date)
		VALUES ($1, $2)
		ON CONFLICT (user_id, date) DO UPDATE SET user_id = daily_activity.user_id
		RETURNING id, user_id, date, xp_earned, lessons_completed, goal_met, created_at
	`, userID, dateOnly).Scan(
		&activity.ID, &activity.UserID, &activity.Date,
		&activity.XPEarned, &activity.LessonsCompleted, &activity.GoalMet, &activity.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &activity, nil
}

// GetDailyActivity gets activity for a specific date.
func GetDailyActivity(ctx context.Context, userID uuid.UUID, date time.Time) (*DailyActivity, error) {
	dateOnly := date.Truncate(24 * time.Hour)
	var activity DailyActivity
	err := db.QueryRow(ctx, `
		SELECT id, user_id, date, xp_earned, lessons_completed, goal_met, created_at
		FROM daily_activity
		WHERE user_id = $1 AND date = $2
	`, userID, dateOnly).Scan(
		&activity.ID, &activity.UserID, &activity.Date,
		&activity.XPEarned, &activity.LessonsCompleted, &activity.GoalMet, &activity.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &activity, nil
}

// UpdateDailyActivity updates activity and checks goal.
func UpdateDailyActivity(ctx context.Context, userID uuid.UUID, date time.Time, xp int, dailyGoal int) (*DailyActivity, error) {
	dateOnly := date.Truncate(24 * time.Hour)
	var activity DailyActivity
	err := db.QueryRow(ctx, `
		INSERT INTO daily_activity (user_id, date, xp_earned, lessons_completed, goal_met)
		VALUES ($1, $2, $3, 1, $3 >= $4)
		ON CONFLICT (user_id, date) DO UPDATE SET
			xp_earned = daily_activity.xp_earned + $3,
			lessons_completed = daily_activity.lessons_completed + 1,
			goal_met = (daily_activity.xp_earned + $3) >= $4
		RETURNING id, user_id, date, xp_earned, lessons_completed, goal_met, created_at
	`, userID, dateOnly, xp, dailyGoal).Scan(
		&activity.ID, &activity.UserID, &activity.Date,
		&activity.XPEarned, &activity.LessonsCompleted, &activity.GoalMet, &activity.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &activity, nil
}

// ListRecentActivity lists recent daily activity.
func ListRecentActivity(ctx context.Context, userID uuid.UUID, days int) ([]*DailyActivity, error) {
	rows, err := db.Query(ctx, `
		SELECT id, user_id, date, xp_earned, lessons_completed, goal_met, created_at
		FROM daily_activity
		WHERE user_id = $1
		ORDER BY date DESC
		LIMIT $2
	`, userID, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []*DailyActivity
	for rows.Next() {
		var a DailyActivity
		if err := rows.Scan(
			&a.ID, &a.UserID, &a.Date,
			&a.XPEarned, &a.LessonsCompleted, &a.GoalMet, &a.CreatedAt,
		); err != nil {
			return nil, err
		}
		activities = append(activities, &a)
	}

	if activities == nil {
		activities = []*DailyActivity{}
	}

	return activities, rows.Err()
}

// GetUsersWithBrokenStreaks finds users who missed their goal yesterday.
func GetUsersWithBrokenStreaks(ctx context.Context, yesterday time.Time) ([]uuid.UUID, error) {
	yesterdayDate := yesterday.Truncate(24 * time.Hour)

	rows, err := db.Query(ctx, `
		SELECT s.user_id
		FROM streaks s
		WHERE s.current_streak > 0
		AND s.last_active_date < $1
		AND NOT EXISTS (
			SELECT 1 FROM daily_activity da
			WHERE da.user_id = s.user_id
			AND da.date = $1
			AND da.goal_met = true
		)
	`, yesterdayDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, id)
	}

	return userIDs, rows.Err()
}
