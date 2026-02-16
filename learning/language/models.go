package language

import (
	"time"

	"github.com/google/uuid"
)

// Language represents a learnable language on the platform.
type Language struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`        // e.g., "Yoruba", "Igbo", "Hausa"
	Code        string    `json:"code"`        // ISO 639-1 code: "yo", "ig", "ha"
	Description string    `json:"description"`
	FlagEmoji   string    `json:"flag_emoji"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// LanguageResponse wraps a single language.
type LanguageResponse struct {
	Language *Language `json:"language"`
}

// LanguagesResponse wraps a list of languages.
type LanguagesResponse struct {
	Languages []*Language `json:"languages"`
}

// CreateLanguageParams contains parameters for creating a language.
type CreateLanguageParams struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description,omitempty"`
	FlagEmoji   string `json:"flag_emoji,omitempty"`
}

// UpdateLanguageParams contains parameters for updating a language.
type UpdateLanguageParams struct {
	Name        *string `json:"name,omitempty"`
	Code        *string `json:"code,omitempty"`
	Description *string `json:"description,omitempty"`
	FlagEmoji   *string `json:"flag_emoji,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}
