package lesson

import (
	"time"

	"github.com/google/uuid"
)

// SessionStatus represents the status of a lesson session.
type SessionStatus string

const (
	SessionStatusInProgress SessionStatus = "in_progress"
	SessionStatusCompleted  SessionStatus = "completed"
	SessionStatusAbandoned  SessionStatus = "abandoned"
)

// LessonSession represents an active or completed lesson attempt.
type LessonSession struct {
	ID             uuid.UUID     `json:"id"`
	UserID         uuid.UUID     `json:"user_id"`
	LessonID       uuid.UUID     `json:"lesson_id"`
	Status         SessionStatus `json:"status"`
	CurrentIndex   int           `json:"current_index"`
	CorrectCount   int           `json:"correct_count"`
	IncorrectCount int           `json:"incorrect_count"`
	XPEarned       int           `json:"xp_earned"`
	StartedAt      time.Time     `json:"started_at"`
	CompletedAt    *time.Time    `json:"completed_at,omitempty"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

// SessionAnswer represents a user's answer to a question within a session.
type SessionAnswer struct {
	ID         uuid.UUID `json:"id"`
	SessionID  uuid.UUID `json:"session_id"`
	QuestionID uuid.UUID `json:"question_id"`
	UserAnswer string    `json:"user_answer"`
	IsCorrect  bool      `json:"is_correct"`
	Similarity float64   `json:"similarity"`
	AnsweredAt time.Time `json:"answered_at"`
}

// ========== Request/Response Types ==========

// StartSessionResponse is returned when starting a new lesson session.
type StartSessionResponse struct {
	Session   *LessonSession `json:"session"`
	Questions []QuestionInfo `json:"questions"`
}

// QuestionInfo contains question data for the client (without correct answer).
type QuestionInfo struct {
	ID         uuid.UUID `json:"id"`
	Type       string    `json:"type"`
	PromptText string    `json:"prompt_text"`
	HasAudio   bool      `json:"has_audio"`
	UseTTS     bool      `json:"use_tts"`
	Options    []string  `json:"options,omitempty"`
	Hint       string    `json:"hint,omitempty"`
	SortOrder  int       `json:"sort_order"`
}

// SubmitAnswerRequest contains the user's answer to a question.
type SubmitAnswerRequest struct {
	QuestionID string `json:"question_id"`
	Answer     string `json:"answer"`
}

// SubmitAnswerResponse is returned after evaluating an answer.
type SubmitAnswerResponse struct {
	IsCorrect     bool    `json:"is_correct"`
	CorrectAnswer string  `json:"correct_answer"`
	Similarity    float64 `json:"similarity,omitempty"`
	Hint          string  `json:"hint,omitempty"`
	XPEarned      int     `json:"xp_earned"`
}

// CompleteSessionResponse is returned when completing a session.
type CompleteSessionResponse struct {
	Session      *LessonSession `json:"session"`
	TotalXP      int            `json:"total_xp"`
	Accuracy     float64        `json:"accuracy"`
	CorrectCount int            `json:"correct_count"`
	TotalCount   int            `json:"total_count"`
}

// SessionResponse wraps a single session.
type SessionResponse struct {
	Session *LessonSession `json:"session"`
}

// SessionsResponse wraps a list of sessions.
type SessionsResponse struct {
	Sessions []*LessonSession `json:"sessions"`
}
