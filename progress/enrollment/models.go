package enrollment

import (
	"time"

	"github.com/google/uuid"
)

// Enrollment represents a user's enrollment in a language.
type Enrollment struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	LanguageID uuid.UUID `json:"language_id"`
	IsActive   bool      `json:"is_active"`
	EnrolledAt time.Time `json:"enrolled_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// EnrollmentWithLanguage includes language details.
type EnrollmentWithLanguage struct {
	Enrollment
	LanguageName  string `json:"language_name"`
	LanguageCode  string `json:"language_code"`
	LanguageEmoji string `json:"language_emoji"`
}

// EnrollmentResponse wraps a single enrollment.
type EnrollmentResponse struct {
	Enrollment *Enrollment `json:"enrollment"`
}

// EnrollmentsResponse wraps a list of enrollments.
type EnrollmentsResponse struct {
	Enrollments []*EnrollmentWithLanguage `json:"enrollments"`
}

// EnrollRequest contains the language to enroll in.
type EnrollRequest struct {
	LanguageID string `json:"language_id"`
}
