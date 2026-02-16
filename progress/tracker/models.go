package tracker

import (
	"time"

	"github.com/google/uuid"
)

// ProgressStatus represents the status of lesson progress.
type ProgressStatus string

const (
	ProgressStatusNotStarted ProgressStatus = "not_started"
	ProgressStatusInProgress ProgressStatus = "in_progress"
	ProgressStatusCompleted  ProgressStatus = "completed"
)

// LessonProgress tracks a user's progress on a specific lesson.
type LessonProgress struct {
	ID          uuid.UUID      `json:"id"`
	UserID      uuid.UUID      `json:"user_id"`
	LessonID    uuid.UUID      `json:"lesson_id"`
	Status      ProgressStatus `json:"status"`
	XPEarned    int            `json:"xp_earned"`
	Attempts    int            `json:"attempts"`
	BestScore   float64        `json:"best_score"`
	CompletedAt *time.Time     `json:"completed_at,omitempty"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// Mistake represents a question the user got wrong.
type Mistake struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	QuestionID    uuid.UUID `json:"question_id"`
	LanguageID    uuid.UUID `json:"language_id"`
	UserAnswer    string    `json:"user_answer"`
	CorrectAnswer string    `json:"correct_answer"`
	Reviewed      bool      `json:"reviewed"`
	CreatedAt     time.Time `json:"created_at"`
}

// MistakeWithQuestion includes question details.
type MistakeWithQuestion struct {
	Mistake
	PromptText string `json:"prompt_text"`
	Hint       string `json:"hint,omitempty"`
}

// ========== Response Types ==========

// ProgressResponse wraps lesson progress.
type ProgressResponse struct {
	Progress *LessonProgress `json:"progress"`
}

// ProgressListResponse wraps a list of progress entries.
type ProgressListResponse struct {
	Progress []*LessonProgress `json:"progress"`
}

// MistakesResponse wraps a list of mistakes.
type MistakesResponse struct {
	Mistakes []*MistakeWithQuestion `json:"mistakes"`
	Total    int                    `json:"total"`
}

// UserStatsResponse contains user statistics.
type UserStatsResponse struct {
	TotalXP            int     `json:"total_xp"`
	LessonsCompleted   int     `json:"lessons_completed"`
	LessonsInProgress  int     `json:"lessons_in_progress"`
	TotalMistakes      int     `json:"total_mistakes"`
	UnreviewedMistakes int     `json:"unreviewed_mistakes"`
	AverageScore       float64 `json:"average_score"`
}
