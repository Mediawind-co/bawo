package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

// FindByID retrieves a user by their UUID
func FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var u User
	err := db.QueryRow(ctx, `
		SELECT id, provider, provider_id, email, name, avatar_url, role, created_at, updated_at
		FROM users WHERE id = $1
	`, id).Scan(&u.ID, &u.Provider, &u.ProviderID, &u.Email, &u.Name, &u.AvatarURL, &u.Role, &u.CreatedAt, &u.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// FindByProviderID retrieves a user by provider and provider_id
func FindByProviderID(ctx context.Context, provider Provider, providerID string) (*User, error) {
	var u User
	err := db.QueryRow(ctx, `
		SELECT id, provider, provider_id, email, name, avatar_url, role, created_at, updated_at
		FROM users WHERE provider = $1 AND provider_id = $2
	`, provider, providerID).Scan(&u.ID, &u.Provider, &u.ProviderID, &u.Email, &u.Name, &u.AvatarURL, &u.Role, &u.CreatedAt, &u.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Create creates a new user and returns it
func Create(ctx context.Context, params CreateUserParams) (*User, error) {
	var u User
	err := db.QueryRow(ctx, `
		INSERT INTO users (provider, provider_id, email, name, avatar_url)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, provider, provider_id, email, name, avatar_url, role, created_at, updated_at
	`, params.Provider, params.ProviderID, params.Email, params.Name, params.AvatarURL).Scan(
		&u.ID, &u.Provider, &u.ProviderID, &u.Email, &u.Name, &u.AvatarURL, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Update updates an existing user
func Update(ctx context.Context, id uuid.UUID, params UpdateUserParams) (*User, error) {
	var u User
	err := db.QueryRow(ctx, `
		UPDATE users
		SET name = COALESCE($2, name),
		    avatar_url = COALESCE($3, avatar_url)
		WHERE id = $1
		RETURNING id, provider, provider_id, email, name, avatar_url, role, created_at, updated_at
	`, id, params.Name, params.AvatarURL).Scan(
		&u.ID, &u.Provider, &u.ProviderID, &u.Email, &u.Name, &u.AvatarURL, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// UpdateRole updates a user's role
func UpdateRole(ctx context.Context, id uuid.UUID, role Role) (*User, error) {
	var u User
	err := db.QueryRow(ctx, `
		UPDATE users SET role = $2
		WHERE id = $1
		RETURNING id, provider, provider_id, email, name, avatar_url, role, created_at, updated_at
	`, id, role).Scan(
		&u.ID, &u.Provider, &u.ProviderID, &u.Email, &u.Name, &u.AvatarURL, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Delete removes a user by ID
func Delete(ctx context.Context, id uuid.UUID) error {
	result, err := db.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}

// FindOrCreate finds an existing user or creates a new one
// Returns the user, a boolean indicating if it was created, and any error
func FindOrCreate(ctx context.Context, params CreateUserParams) (*User, bool, error) {
	// Try to find existing user first
	user, err := FindByProviderID(ctx, params.Provider, params.ProviderID)
	if err == nil {
		return user, false, nil // found, not created
	}
	if !errors.Is(err, ErrUserNotFound) {
		return nil, false, err
	}

	// Create new user
	user, err = Create(ctx, params)
	if err != nil {
		return nil, false, err
	}
	return user, true, nil // created
}

// List returns all users with pagination
func List(ctx context.Context, limit, offset int) ([]*User, int, error) {
	// Get total count
	var total int
	err := db.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get users
	rows, err := db.Query(ctx, `
		SELECT id, provider, provider_id, email, name, avatar_url, role, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Provider, &u.ProviderID, &u.Email, &u.Name, &u.AvatarURL, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
