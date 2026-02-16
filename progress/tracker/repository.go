package tracker

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrProgressNotFound = errors.New("progress not found")
	ErrMistakeNotFound  = errors.New("mistake not found")
)

// ========== Progress Repository ==========

// GetOrCreateProgress gets existing progress or creates new.
func GetOrCreateProgress(ctx context.Context, userID, lessonID uuid.UUID) (*LessonProgress, error) {
	var progress LessonProgress
	err := db.QueryRow(ctx, `
		INSERT INTO lesson_progress (user_id, lesson_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, lesson_id) DO UPDATE SET updated_at = NOW()
		RETURNING id, user_id, lesson_id, status, xp_earned, attempts, best_score, completed_at, updated_at
	`, userID, lessonID).Scan(
		&progress.ID, &progress.UserID, &progress.LessonID, &progress.Status,
		&progress.XPEarned, &progress.Attempts, &progress.BestScore,
		&progress.CompletedAt, &progress.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &progress, nil
}

// FindProgressByUserAndLesson finds progress for a user and lesson.
func FindProgressByUserAndLesson(ctx context.Context, userID, lessonID uuid.UUID) (*LessonProgress, error) {
	var progress LessonProgress
	err := db.QueryRow(ctx, `
		SELECT id, user_id, lesson_id, status, xp_earned, attempts, best_score, completed_at, updated_at
		FROM lesson_progress
		WHERE user_id = $1 AND lesson_id = $2
	`, userID, lessonID).Scan(
		&progress.ID, &progress.UserID, &progress.LessonID, &progress.Status,
		&progress.XPEarned, &progress.Attempts, &progress.BestScore,
		&progress.CompletedAt, &progress.UpdatedAt,
	)
	if err != nil {
		return nil, ErrProgressNotFound
	}
	return &progress, nil
}

// ListProgressByUser lists all progress for a user.
func ListProgressByUser(ctx context.Context, userID uuid.UUID) ([]*LessonProgress, error) {
	rows, err := db.Query(ctx, `
		SELECT id, user_id, lesson_id, status, xp_earned, attempts, best_score, completed_at, updated_at
		FROM lesson_progress
		WHERE user_id = $1
		ORDER BY updated_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var progress []*LessonProgress
	for rows.Next() {
		var p LessonProgress
		if err := rows.Scan(
			&p.ID, &p.UserID, &p.LessonID, &p.Status,
			&p.XPEarned, &p.Attempts, &p.BestScore,
			&p.CompletedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		progress = append(progress, &p)
	}

	if progress == nil {
		progress = []*LessonProgress{}
	}

	return progress, rows.Err()
}

// UpdateProgress updates lesson progress after completion.
func UpdateProgress(ctx context.Context, userID, lessonID uuid.UUID, xpEarned int, score float64) (*LessonProgress, error) {
	now := time.Now()
	var progress LessonProgress
	err := db.QueryRow(ctx, `
		INSERT INTO lesson_progress (user_id, lesson_id, status, xp_earned, attempts, best_score, completed_at)
		VALUES ($1, $2, 'completed', $3, 1, $4, $5)
		ON CONFLICT (user_id, lesson_id) DO UPDATE SET
			status = 'completed',
			xp_earned = lesson_progress.xp_earned + $3,
			attempts = lesson_progress.attempts + 1,
			best_score = GREATEST(lesson_progress.best_score, $4),
			completed_at = COALESCE(lesson_progress.completed_at, $5),
			updated_at = NOW()
		RETURNING id, user_id, lesson_id, status, xp_earned, attempts, best_score, completed_at, updated_at
	`, userID, lessonID, xpEarned, score, now).Scan(
		&progress.ID, &progress.UserID, &progress.LessonID, &progress.Status,
		&progress.XPEarned, &progress.Attempts, &progress.BestScore,
		&progress.CompletedAt, &progress.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &progress, nil
}

// GetUserStats gets aggregate statistics for a user.
func GetUserStats(ctx context.Context, userID uuid.UUID) (*UserStatsResponse, error) {
	var stats UserStatsResponse

	// Get XP and lesson counts
	err := db.QueryRow(ctx, `
		SELECT
			COALESCE(SUM(xp_earned), 0),
			COUNT(*) FILTER (WHERE status = 'completed'),
			COUNT(*) FILTER (WHERE status = 'in_progress'),
			COALESCE(AVG(best_score) FILTER (WHERE status = 'completed'), 0)
		FROM lesson_progress
		WHERE user_id = $1
	`, userID).Scan(&stats.TotalXP, &stats.LessonsCompleted, &stats.LessonsInProgress, &stats.AverageScore)
	if err != nil {
		return nil, err
	}

	// Get mistake counts
	err = db.QueryRow(ctx, `
		SELECT
			COUNT(*),
			COUNT(*) FILTER (WHERE reviewed = false)
		FROM mistakes
		WHERE user_id = $1
	`, userID).Scan(&stats.TotalMistakes, &stats.UnreviewedMistakes)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// ========== Mistake Repository ==========

// RecordMistake records a mistake.
func RecordMistake(ctx context.Context, userID, questionID, languageID uuid.UUID, userAnswer, correctAnswer string) (*Mistake, error) {
	var mistake Mistake
	err := db.QueryRow(ctx, `
		INSERT INTO mistakes (user_id, question_id, language_id, user_answer, correct_answer)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, question_id, language_id, user_answer, correct_answer, reviewed, created_at
	`, userID, questionID, languageID, userAnswer, correctAnswer).Scan(
		&mistake.ID, &mistake.UserID, &mistake.QuestionID, &mistake.LanguageID,
		&mistake.UserAnswer, &mistake.CorrectAnswer, &mistake.Reviewed, &mistake.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &mistake, nil
}

// ListMistakesByUser lists all mistakes for a user.
func ListMistakesByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Mistake, int, error) {
	var total int
	err := db.QueryRow(ctx, `SELECT COUNT(*) FROM mistakes WHERE user_id = $1`, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := db.Query(ctx, `
		SELECT id, user_id, question_id, language_id, user_answer, correct_answer, reviewed, created_at
		FROM mistakes
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var mistakes []*Mistake
	for rows.Next() {
		var m Mistake
		if err := rows.Scan(
			&m.ID, &m.UserID, &m.QuestionID, &m.LanguageID,
			&m.UserAnswer, &m.CorrectAnswer, &m.Reviewed, &m.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		mistakes = append(mistakes, &m)
	}

	if mistakes == nil {
		mistakes = []*Mistake{}
	}

	return mistakes, total, rows.Err()
}

// ListMistakesByLanguage lists mistakes for a user filtered by language.
func ListMistakesByLanguage(ctx context.Context, userID, languageID uuid.UUID, limit, offset int) ([]*Mistake, int, error) {
	var total int
	err := db.QueryRow(ctx, `
		SELECT COUNT(*) FROM mistakes WHERE user_id = $1 AND language_id = $2
	`, userID, languageID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := db.Query(ctx, `
		SELECT id, user_id, question_id, language_id, user_answer, correct_answer, reviewed, created_at
		FROM mistakes
		WHERE user_id = $1 AND language_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`, userID, languageID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var mistakes []*Mistake
	for rows.Next() {
		var m Mistake
		if err := rows.Scan(
			&m.ID, &m.UserID, &m.QuestionID, &m.LanguageID,
			&m.UserAnswer, &m.CorrectAnswer, &m.Reviewed, &m.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		mistakes = append(mistakes, &m)
	}

	if mistakes == nil {
		mistakes = []*Mistake{}
	}

	return mistakes, total, rows.Err()
}

// MarkMistakeReviewed marks a mistake as reviewed.
func MarkMistakeReviewed(ctx context.Context, mistakeID uuid.UUID) error {
	result, err := db.Exec(ctx, `
		UPDATE mistakes SET reviewed = true WHERE id = $1
	`, mistakeID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrMistakeNotFound
	}
	return nil
}
