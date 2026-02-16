package language

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrLanguageNotFound   = errors.New("language not found")
	ErrLanguageExists     = errors.New("language already exists")
	ErrInvalidLanguage    = errors.New("invalid language data")
)

// Create inserts a new language into the database.
func Create(ctx context.Context, params CreateLanguageParams) (*Language, error) {
	if params.Name == "" || params.Code == "" {
		return nil, ErrInvalidLanguage
	}

	var lang Language
	err := db.QueryRow(ctx, `
		INSERT INTO languages (name, code, description, flag_emoji)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, code, description, flag_emoji, is_active, created_at, updated_at
	`, params.Name, params.Code, params.Description, params.FlagEmoji).Scan(
		&lang.ID, &lang.Name, &lang.Code, &lang.Description,
		&lang.FlagEmoji, &lang.IsActive, &lang.CreatedAt, &lang.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &lang, nil
}

// FindByID retrieves a language by its ID.
func FindByID(ctx context.Context, id uuid.UUID) (*Language, error) {
	var lang Language
	err := db.QueryRow(ctx, `
		SELECT id, name, code, description, flag_emoji, is_active, created_at, updated_at
		FROM languages
		WHERE id = $1
	`, id).Scan(
		&lang.ID, &lang.Name, &lang.Code, &lang.Description,
		&lang.FlagEmoji, &lang.IsActive, &lang.CreatedAt, &lang.UpdatedAt,
	)
	if err != nil {
		return nil, ErrLanguageNotFound
	}
	return &lang, nil
}

// FindByCode retrieves a language by its code.
func FindByCode(ctx context.Context, code string) (*Language, error) {
	var lang Language
	err := db.QueryRow(ctx, `
		SELECT id, name, code, description, flag_emoji, is_active, created_at, updated_at
		FROM languages
		WHERE code = $1
	`, code).Scan(
		&lang.ID, &lang.Name, &lang.Code, &lang.Description,
		&lang.FlagEmoji, &lang.IsActive, &lang.CreatedAt, &lang.UpdatedAt,
	)
	if err != nil {
		return nil, ErrLanguageNotFound
	}
	return &lang, nil
}

// List retrieves all active languages.
func List(ctx context.Context, includeInactive bool) ([]*Language, error) {
	query := `
		SELECT id, name, code, description, flag_emoji, is_active, created_at, updated_at
		FROM languages
	`
	if !includeInactive {
		query += ` WHERE is_active = true`
	}
	query += ` ORDER BY name ASC`

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var languages []*Language
	for rows.Next() {
		var lang Language
		if err := rows.Scan(
			&lang.ID, &lang.Name, &lang.Code, &lang.Description,
			&lang.FlagEmoji, &lang.IsActive, &lang.CreatedAt, &lang.UpdatedAt,
		); err != nil {
			return nil, err
		}
		languages = append(languages, &lang)
	}

	if languages == nil {
		languages = []*Language{}
	}

	return languages, rows.Err()
}

// Update modifies an existing language.
func Update(ctx context.Context, id uuid.UUID, params UpdateLanguageParams) (*Language, error) {
	// First check if language exists
	existing, err := FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	name := existing.Name
	code := existing.Code
	description := existing.Description
	flagEmoji := existing.FlagEmoji
	isActive := existing.IsActive

	if params.Name != nil {
		name = *params.Name
	}
	if params.Code != nil {
		code = *params.Code
	}
	if params.Description != nil {
		description = *params.Description
	}
	if params.FlagEmoji != nil {
		flagEmoji = *params.FlagEmoji
	}
	if params.IsActive != nil {
		isActive = *params.IsActive
	}

	var lang Language
	err = db.QueryRow(ctx, `
		UPDATE languages
		SET name = $2, code = $3, description = $4, flag_emoji = $5, is_active = $6, updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, code, description, flag_emoji, is_active, created_at, updated_at
	`, id, name, code, description, flagEmoji, isActive).Scan(
		&lang.ID, &lang.Name, &lang.Code, &lang.Description,
		&lang.FlagEmoji, &lang.IsActive, &lang.CreatedAt, &lang.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &lang, nil
}

// Delete removes a language by its ID.
func Delete(ctx context.Context, id uuid.UUID) error {
	result, err := db.Exec(ctx, `DELETE FROM languages WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrLanguageNotFound
	}
	return nil
}
