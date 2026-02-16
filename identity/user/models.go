package user

import (
	"time"

	"github.com/google/uuid"
)

// Provider represents the authentication provider
type Provider string

const (
	ProviderGoogle Provider = "google"
	ProviderApple  Provider = "apple"
)

// Role represents user authorization level
type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// User represents a user in the system
type User struct {
	ID         uuid.UUID `json:"id"`
	Provider   Provider  `json:"provider"`
	ProviderID string    `json:"provider_id"`
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	AvatarURL  string    `json:"avatar_url,omitempty"`
	Role       Role      `json:"role"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CreateUserParams contains parameters for creating a new user
type CreateUserParams struct {
	Provider   Provider
	ProviderID string
	Email      string
	Name       string
	AvatarURL  string
}

// UpdateUserParams contains parameters for updating a user
type UpdateUserParams struct {
	Name      *string
	AvatarURL *string
}

// UserResponse is the API response for user data
type UserResponse struct {
	User *User `json:"user"`
}

// UsersResponse is the API response for multiple users
type UsersResponse struct {
	Users []*User `json:"users"`
	Total int     `json:"total"`
}
