package content

import (
	"time"

	"github.com/google/uuid"
)

// QuestionType defines the type of question.
type QuestionType string

const (
	QuestionTypeListenReply  QuestionType = "listen_reply"
	QuestionTypeMultiChoice  QuestionType = "multi_choice"
	QuestionTypeSingleChoice QuestionType = "single_choice"
)

// Unit represents a group of lessons within a language.
type Unit struct {
	ID          uuid.UUID `json:"id"`
	LanguageID  uuid.UUID `json:"language_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Lesson represents a single learning session within a unit.
type Lesson struct {
	ID          uuid.UUID `json:"id"`
	UnitID      uuid.UUID `json:"unit_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	XPReward    int       `json:"xp_reward"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Question represents a question within a lesson.
type Question struct {
	ID             uuid.UUID    `json:"id"`
	LessonID       uuid.UUID    `json:"lesson_id"`
	Type           QuestionType `json:"type"`
	PromptText     string       `json:"prompt_text"`
	PromptAudioKey string       `json:"prompt_audio_key,omitempty"`
	UseTTS         bool         `json:"use_tts"`
	CorrectAnswer  string       `json:"correct_answer"`
	Options        []string     `json:"options"`
	Hint           string       `json:"hint,omitempty"`
	SortOrder      int          `json:"sort_order"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

// ========== Request/Response Types ==========

// UnitResponse wraps a single unit.
type UnitResponse struct {
	Unit *Unit `json:"unit"`
}

// UnitsResponse wraps a list of units.
type UnitsResponse struct {
	Units []*Unit `json:"units"`
}

// LessonResponse wraps a single lesson.
type LessonResponse struct {
	Lesson *Lesson `json:"lesson"`
}

// LessonsResponse wraps a list of lessons.
type LessonsResponse struct {
	Lessons []*Lesson `json:"lessons"`
}

// QuestionResponse wraps a single question.
type QuestionResponse struct {
	Question *Question `json:"question"`
}

// QuestionsResponse wraps a list of questions.
type QuestionsResponse struct {
	Questions []*Question `json:"questions"`
}

// ========== Create Params ==========

// CreateUnitParams contains parameters for creating a unit.
type CreateUnitParams struct {
	LanguageID  string `json:"language_id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	SortOrder   int    `json:"sort_order,omitempty"`
}

// CreateLessonParams contains parameters for creating a lesson.
type CreateLessonParams struct {
	UnitID      string `json:"unit_id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	XPReward    int    `json:"xp_reward,omitempty"`
	SortOrder   int    `json:"sort_order,omitempty"`
}

// CreateQuestionParams contains parameters for creating a question.
type CreateQuestionParams struct {
	LessonID      string       `json:"lesson_id"`
	Type          QuestionType `json:"type"`
	PromptText    string       `json:"prompt_text"`
	UseTTS        bool         `json:"use_tts,omitempty"`
	CorrectAnswer string       `json:"correct_answer"`
	Options       []string     `json:"options,omitempty"`
	Hint          string       `json:"hint,omitempty"`
	SortOrder     int          `json:"sort_order,omitempty"`
}

// ========== Update Params ==========

// UpdateUnitParams contains parameters for updating a unit.
type UpdateUnitParams struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	SortOrder   *int    `json:"sort_order,omitempty"`
}

// UpdateLessonParams contains parameters for updating a lesson.
type UpdateLessonParams struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	XPReward    *int    `json:"xp_reward,omitempty"`
	SortOrder   *int    `json:"sort_order,omitempty"`
}

// UpdateQuestionParams contains parameters for updating a question.
type UpdateQuestionParams struct {
	Type          *QuestionType `json:"type,omitempty"`
	PromptText    *string       `json:"prompt_text,omitempty"`
	UseTTS        *bool         `json:"use_tts,omitempty"`
	CorrectAnswer *string       `json:"correct_answer,omitempty"`
	Options       []string      `json:"options,omitempty"`
	Hint          *string       `json:"hint,omitempty"`
	SortOrder     *int          `json:"sort_order,omitempty"`
}
