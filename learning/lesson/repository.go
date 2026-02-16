package lesson

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrAnswerNotFound  = errors.New("answer not found")
	ErrInvalidData     = errors.New("invalid data")
)

// ========== Session Repository ==========

// CreateSession creates a new lesson session.
func CreateSession(ctx context.Context, userID, lessonID uuid.UUID) (*LessonSession, error) {
	var session LessonSession
	err := db.QueryRow(ctx, `
		INSERT INTO lesson_sessions (user_id, lesson_id)
		VALUES ($1, $2)
		RETURNING id, user_id, lesson_id, status, current_index, correct_count, incorrect_count, xp_earned, started_at, completed_at, updated_at
	`, userID, lessonID).Scan(
		&session.ID, &session.UserID, &session.LessonID, &session.Status,
		&session.CurrentIndex, &session.CorrectCount, &session.IncorrectCount,
		&session.XPEarned, &session.StartedAt, &session.CompletedAt, &session.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// FindSessionByID retrieves a session by its ID.
func FindSessionByID(ctx context.Context, id uuid.UUID) (*LessonSession, error) {
	var session LessonSession
	err := db.QueryRow(ctx, `
		SELECT id, user_id, lesson_id, status, current_index, correct_count, incorrect_count, xp_earned, started_at, completed_at, updated_at
		FROM lesson_sessions
		WHERE id = $1
	`, id).Scan(
		&session.ID, &session.UserID, &session.LessonID, &session.Status,
		&session.CurrentIndex, &session.CorrectCount, &session.IncorrectCount,
		&session.XPEarned, &session.StartedAt, &session.CompletedAt, &session.UpdatedAt,
	)
	if err != nil {
		return nil, ErrSessionNotFound
	}
	return &session, nil
}

// FindActiveSession finds an in-progress session for a user and lesson.
func FindActiveSession(ctx context.Context, userID, lessonID uuid.UUID) (*LessonSession, error) {
	var session LessonSession
	err := db.QueryRow(ctx, `
		SELECT id, user_id, lesson_id, status, current_index, correct_count, incorrect_count, xp_earned, started_at, completed_at, updated_at
		FROM lesson_sessions
		WHERE user_id = $1 AND lesson_id = $2 AND status = 'in_progress'
		ORDER BY started_at DESC
		LIMIT 1
	`, userID, lessonID).Scan(
		&session.ID, &session.UserID, &session.LessonID, &session.Status,
		&session.CurrentIndex, &session.CorrectCount, &session.IncorrectCount,
		&session.XPEarned, &session.StartedAt, &session.CompletedAt, &session.UpdatedAt,
	)
	if err != nil {
		return nil, ErrSessionNotFound
	}
	return &session, nil
}

// ListUserSessions retrieves all sessions for a user.
func ListUserSessions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*LessonSession, error) {
	rows, err := db.Query(ctx, `
		SELECT id, user_id, lesson_id, status, current_index, correct_count, incorrect_count, xp_earned, started_at, completed_at, updated_at
		FROM lesson_sessions
		WHERE user_id = $1
		ORDER BY started_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*LessonSession
	for rows.Next() {
		var session LessonSession
		if err := rows.Scan(
			&session.ID, &session.UserID, &session.LessonID, &session.Status,
			&session.CurrentIndex, &session.CorrectCount, &session.IncorrectCount,
			&session.XPEarned, &session.StartedAt, &session.CompletedAt, &session.UpdatedAt,
		); err != nil {
			return nil, err
		}
		sessions = append(sessions, &session)
	}

	if sessions == nil {
		sessions = []*LessonSession{}
	}

	return sessions, rows.Err()
}

// UpdateSessionProgress updates session progress after an answer.
func UpdateSessionProgress(ctx context.Context, id uuid.UUID, correct bool, xp int) (*LessonSession, error) {
	var query string
	if correct {
		query = `
			UPDATE lesson_sessions
			SET current_index = current_index + 1, correct_count = correct_count + 1, xp_earned = xp_earned + $2, updated_at = NOW()
			WHERE id = $1
			RETURNING id, user_id, lesson_id, status, current_index, correct_count, incorrect_count, xp_earned, started_at, completed_at, updated_at
		`
	} else {
		query = `
			UPDATE lesson_sessions
			SET current_index = current_index + 1, incorrect_count = incorrect_count + 1, updated_at = NOW()
			WHERE id = $1
			RETURNING id, user_id, lesson_id, status, current_index, correct_count, incorrect_count, xp_earned, started_at, completed_at, updated_at
		`
	}

	var session LessonSession
	err := db.QueryRow(ctx, query, id, xp).Scan(
		&session.ID, &session.UserID, &session.LessonID, &session.Status,
		&session.CurrentIndex, &session.CorrectCount, &session.IncorrectCount,
		&session.XPEarned, &session.StartedAt, &session.CompletedAt, &session.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// CompleteSession marks a session as completed.
func CompleteSession(ctx context.Context, id uuid.UUID) (*LessonSession, error) {
	now := time.Now()
	var session LessonSession
	err := db.QueryRow(ctx, `
		UPDATE lesson_sessions
		SET status = 'completed', completed_at = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING id, user_id, lesson_id, status, current_index, correct_count, incorrect_count, xp_earned, started_at, completed_at, updated_at
	`, id, now).Scan(
		&session.ID, &session.UserID, &session.LessonID, &session.Status,
		&session.CurrentIndex, &session.CorrectCount, &session.IncorrectCount,
		&session.XPEarned, &session.StartedAt, &session.CompletedAt, &session.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// AbandonSession marks a session as abandoned.
func AbandonSession(ctx context.Context, id uuid.UUID) error {
	result, err := db.Exec(ctx, `
		UPDATE lesson_sessions
		SET status = 'abandoned', updated_at = NOW()
		WHERE id = $1 AND status = 'in_progress'
	`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrSessionNotFound
	}
	return nil
}

// ========== Answer Repository ==========

// RecordAnswer saves a user's answer to a question.
func RecordAnswer(ctx context.Context, sessionID, questionID uuid.UUID, userAnswer string, isCorrect bool, similarity float64) (*SessionAnswer, error) {
	var answer SessionAnswer
	err := db.QueryRow(ctx, `
		INSERT INTO session_answers (session_id, question_id, user_answer, is_correct, similarity)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, session_id, question_id, user_answer, is_correct, similarity, answered_at
	`, sessionID, questionID, userAnswer, isCorrect, similarity).Scan(
		&answer.ID, &answer.SessionID, &answer.QuestionID,
		&answer.UserAnswer, &answer.IsCorrect, &answer.Similarity, &answer.AnsweredAt,
	)
	if err != nil {
		return nil, err
	}
	return &answer, nil
}

// ListSessionAnswers retrieves all answers for a session.
func ListSessionAnswers(ctx context.Context, sessionID uuid.UUID) ([]*SessionAnswer, error) {
	rows, err := db.Query(ctx, `
		SELECT id, session_id, question_id, user_answer, is_correct, similarity, answered_at
		FROM session_answers
		WHERE session_id = $1
		ORDER BY answered_at ASC
	`, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []*SessionAnswer
	for rows.Next() {
		var answer SessionAnswer
		if err := rows.Scan(
			&answer.ID, &answer.SessionID, &answer.QuestionID,
			&answer.UserAnswer, &answer.IsCorrect, &answer.Similarity, &answer.AnsweredAt,
		); err != nil {
			return nil, err
		}
		answers = append(answers, &answer)
	}

	if answers == nil {
		answers = []*SessionAnswer{}
	}

	return answers, rows.Err()
}
