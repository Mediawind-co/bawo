package enrollment

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrEnrollmentNotFound = errors.New("enrollment not found")
	ErrAlreadyEnrolled    = errors.New("already enrolled in this language")
)

// Create creates a new enrollment.
func Create(ctx context.Context, userID, languageID uuid.UUID) (*Enrollment, error) {
	var enrollment Enrollment
	err := db.QueryRow(ctx, `
		INSERT INTO enrollments (user_id, language_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, language_id)
		DO UPDATE SET is_active = true, updated_at = NOW()
		RETURNING id, user_id, language_id, is_active, enrolled_at, updated_at
	`, userID, languageID).Scan(
		&enrollment.ID, &enrollment.UserID, &enrollment.LanguageID,
		&enrollment.IsActive, &enrollment.EnrolledAt, &enrollment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &enrollment, nil
}

// FindByUserAndLanguage finds an enrollment by user and language.
func FindByUserAndLanguage(ctx context.Context, userID, languageID uuid.UUID) (*Enrollment, error) {
	var enrollment Enrollment
	err := db.QueryRow(ctx, `
		SELECT id, user_id, language_id, is_active, enrolled_at, updated_at
		FROM enrollments
		WHERE user_id = $1 AND language_id = $2
	`, userID, languageID).Scan(
		&enrollment.ID, &enrollment.UserID, &enrollment.LanguageID,
		&enrollment.IsActive, &enrollment.EnrolledAt, &enrollment.UpdatedAt,
	)
	if err != nil {
		return nil, ErrEnrollmentNotFound
	}
	return &enrollment, nil
}

// ListByUser lists all active enrollments for a user.
func ListByUser(ctx context.Context, userID uuid.UUID) ([]*Enrollment, error) {
	rows, err := db.Query(ctx, `
		SELECT id, user_id, language_id, is_active, enrolled_at, updated_at
		FROM enrollments
		WHERE user_id = $1 AND is_active = true
		ORDER BY enrolled_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enrollments []*Enrollment
	for rows.Next() {
		var e Enrollment
		if err := rows.Scan(
			&e.ID, &e.UserID, &e.LanguageID,
			&e.IsActive, &e.EnrolledAt, &e.UpdatedAt,
		); err != nil {
			return nil, err
		}
		enrollments = append(enrollments, &e)
	}

	if enrollments == nil {
		enrollments = []*Enrollment{}
	}

	return enrollments, rows.Err()
}

// Unenroll deactivates an enrollment.
func Unenroll(ctx context.Context, userID, languageID uuid.UUID) error {
	result, err := db.Exec(ctx, `
		UPDATE enrollments
		SET is_active = false, updated_at = NOW()
		WHERE user_id = $1 AND language_id = $2 AND is_active = true
	`, userID, languageID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrEnrollmentNotFound
	}
	return nil
}

// IsEnrolled checks if a user is enrolled in a language.
func IsEnrolled(ctx context.Context, userID, languageID uuid.UUID) (bool, error) {
	var exists bool
	err := db.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM enrollments
			WHERE user_id = $1 AND language_id = $2 AND is_active = true
		)
	`, userID, languageID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// CountByLanguage counts enrollments per language.
func CountByLanguage(ctx context.Context, languageID uuid.UUID) (int, error) {
	var count int
	err := db.QueryRow(ctx, `
		SELECT COUNT(*) FROM enrollments
		WHERE language_id = $1 AND is_active = true
	`, languageID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
