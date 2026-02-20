package content

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrUnitNotFound     = errors.New("unit not found")
	ErrLessonNotFound   = errors.New("lesson not found")
	ErrQuestionNotFound = errors.New("question not found")
	ErrInvalidData      = errors.New("invalid data")
)

// ========== Unit Repository ==========

// CreateUnit inserts a new unit into the database.
func CreateUnit(ctx context.Context, languageID uuid.UUID, params CreateUnitParams) (*Unit, error) {
	if params.Title == "" {
		return nil, ErrInvalidData
	}

	var unit Unit
	err := db.QueryRow(ctx, `
		INSERT INTO units (language_id, title, description, sort_order)
		VALUES ($1, $2, $3, $4)
		RETURNING id, language_id, title, description, sort_order, created_at, updated_at
	`, languageID, params.Title, params.Description, params.SortOrder).Scan(
		&unit.ID, &unit.LanguageID, &unit.Title, &unit.Description,
		&unit.SortOrder, &unit.CreatedAt, &unit.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &unit, nil
}

// FindUnitByID retrieves a unit by its ID.
func FindUnitByID(ctx context.Context, id uuid.UUID) (*Unit, error) {
	var unit Unit
	err := db.QueryRow(ctx, `
		SELECT id, language_id, title, description, sort_order, created_at, updated_at
		FROM units
		WHERE id = $1
	`, id).Scan(
		&unit.ID, &unit.LanguageID, &unit.Title, &unit.Description,
		&unit.SortOrder, &unit.CreatedAt, &unit.UpdatedAt,
	)
	if err != nil {
		return nil, ErrUnitNotFound
	}
	return &unit, nil
}

// ListUnitsByLanguage retrieves all units for a language.
func ListUnitsByLanguage(ctx context.Context, languageID uuid.UUID) ([]*Unit, error) {
	rows, err := db.Query(ctx, `
		SELECT id, language_id, title, description, sort_order, created_at, updated_at
		FROM units
		WHERE language_id = $1
		ORDER BY sort_order ASC
	`, languageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []*Unit
	for rows.Next() {
		var unit Unit
		if err := rows.Scan(
			&unit.ID, &unit.LanguageID, &unit.Title, &unit.Description,
			&unit.SortOrder, &unit.CreatedAt, &unit.UpdatedAt,
		); err != nil {
			return nil, err
		}
		units = append(units, &unit)
	}

	if units == nil {
		units = []*Unit{}
	}

	return units, rows.Err()
}

// UpdateUnit modifies an existing unit.
func UpdateUnit(ctx context.Context, id uuid.UUID, params UpdateUnitParams) (*Unit, error) {
	existing, err := FindUnitByID(ctx, id)
	if err != nil {
		return nil, err
	}

	title := existing.Title
	description := existing.Description
	sortOrder := existing.SortOrder

	if params.Title != nil {
		title = *params.Title
	}
	if params.Description != nil {
		description = *params.Description
	}
	if params.SortOrder != nil {
		sortOrder = *params.SortOrder
	}

	var unit Unit
	err = db.QueryRow(ctx, `
		UPDATE units
		SET title = $2, description = $3, sort_order = $4, updated_at = NOW()
		WHERE id = $1
		RETURNING id, language_id, title, description, sort_order, created_at, updated_at
	`, id, title, description, sortOrder).Scan(
		&unit.ID, &unit.LanguageID, &unit.Title, &unit.Description,
		&unit.SortOrder, &unit.CreatedAt, &unit.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &unit, nil
}

// DeleteUnit removes a unit by its ID.
func DeleteUnit(ctx context.Context, id uuid.UUID) error {
	result, err := db.Exec(ctx, `DELETE FROM units WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrUnitNotFound
	}
	return nil
}

// ========== Lesson Repository ==========

// CreateLesson inserts a new lesson into the database.
func CreateLesson(ctx context.Context, unitID uuid.UUID, params CreateLessonParams) (*Lesson, error) {
	if params.Title == "" {
		return nil, ErrInvalidData
	}

	xpReward := params.XPReward
	if xpReward == 0 {
		xpReward = 10 // default XP
	}

	var lesson Lesson
	err := db.QueryRow(ctx, `
		INSERT INTO lessons (unit_id, title, description, xp_reward, sort_order)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, unit_id, title, description, xp_reward, sort_order, created_at, updated_at
	`, unitID, params.Title, params.Description, xpReward, params.SortOrder).Scan(
		&lesson.ID, &lesson.UnitID, &lesson.Title, &lesson.Description,
		&lesson.XPReward, &lesson.SortOrder, &lesson.CreatedAt, &lesson.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &lesson, nil
}

// FindLessonByID retrieves a lesson by its ID.
func FindLessonByID(ctx context.Context, id uuid.UUID) (*Lesson, error) {
	var lesson Lesson
	err := db.QueryRow(ctx, `
		SELECT id, unit_id, title, description, xp_reward, sort_order, created_at, updated_at
		FROM lessons
		WHERE id = $1
	`, id).Scan(
		&lesson.ID, &lesson.UnitID, &lesson.Title, &lesson.Description,
		&lesson.XPReward, &lesson.SortOrder, &lesson.CreatedAt, &lesson.UpdatedAt,
	)
	if err != nil {
		return nil, ErrLessonNotFound
	}
	return &lesson, nil
}

// ListLessonsByUnit retrieves all lessons for a unit.
func ListLessonsByUnit(ctx context.Context, unitID uuid.UUID) ([]*Lesson, error) {
	rows, err := db.Query(ctx, `
		SELECT id, unit_id, title, description, xp_reward, sort_order, created_at, updated_at
		FROM lessons
		WHERE unit_id = $1
		ORDER BY sort_order ASC
	`, unitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []*Lesson
	for rows.Next() {
		var lesson Lesson
		if err := rows.Scan(
			&lesson.ID, &lesson.UnitID, &lesson.Title, &lesson.Description,
			&lesson.XPReward, &lesson.SortOrder, &lesson.CreatedAt, &lesson.UpdatedAt,
		); err != nil {
			return nil, err
		}
		lessons = append(lessons, &lesson)
	}

	if lessons == nil {
		lessons = []*Lesson{}
	}

	return lessons, rows.Err()
}

// UpdateLesson modifies an existing lesson.
func UpdateLesson(ctx context.Context, id uuid.UUID, params UpdateLessonParams) (*Lesson, error) {
	existing, err := FindLessonByID(ctx, id)
	if err != nil {
		return nil, err
	}

	title := existing.Title
	description := existing.Description
	xpReward := existing.XPReward
	sortOrder := existing.SortOrder

	if params.Title != nil {
		title = *params.Title
	}
	if params.Description != nil {
		description = *params.Description
	}
	if params.XPReward != nil {
		xpReward = *params.XPReward
	}
	if params.SortOrder != nil {
		sortOrder = *params.SortOrder
	}

	var lesson Lesson
	err = db.QueryRow(ctx, `
		UPDATE lessons
		SET title = $2, description = $3, xp_reward = $4, sort_order = $5, updated_at = NOW()
		WHERE id = $1
		RETURNING id, unit_id, title, description, xp_reward, sort_order, created_at, updated_at
	`, id, title, description, xpReward, sortOrder).Scan(
		&lesson.ID, &lesson.UnitID, &lesson.Title, &lesson.Description,
		&lesson.XPReward, &lesson.SortOrder, &lesson.CreatedAt, &lesson.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &lesson, nil
}

// DeleteLesson removes a lesson by its ID.
func DeleteLesson(ctx context.Context, id uuid.UUID) error {
	result, err := db.Exec(ctx, `DELETE FROM lessons WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrLessonNotFound
	}
	return nil
}

// ========== Question Repository ==========

// CreateQuestion inserts a new question into the database.
func CreateQuestion(ctx context.Context, lessonID uuid.UUID, params CreateQuestionParams) (*Question, error) {
	if params.PromptText == "" || params.CorrectAnswer == "" {
		return nil, ErrInvalidData
	}

	optionsJSON, err := json.Marshal(params.Options)
	if err != nil {
		return nil, err
	}

	var question Question
	var optionsBytes []byte
	err = db.QueryRow(ctx, `
		INSERT INTO questions (lesson_id, type, prompt_text, use_tts, correct_answer, options, hint, sort_order)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, lesson_id, type, prompt_text, prompt_audio_key, use_tts, correct_answer, options, hint, sort_order, created_at, updated_at
	`, lessonID, params.Type, params.PromptText, params.UseTTS, params.CorrectAnswer, optionsJSON, params.Hint, params.SortOrder).Scan(
		&question.ID, &question.LessonID, &question.Type, &question.PromptText,
		&question.PromptAudioKey, &question.UseTTS, &question.CorrectAnswer,
		&optionsBytes, &question.Hint, &question.SortOrder, &question.CreatedAt, &question.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(optionsBytes, &question.Options); err != nil {
		question.Options = []string{}
	}

	return &question, nil
}

// FindQuestionByID retrieves a question by its ID.
func FindQuestionByID(ctx context.Context, id uuid.UUID) (*Question, error) {
	var question Question
	var optionsBytes []byte
	err := db.QueryRow(ctx, `
		SELECT id, lesson_id, type, prompt_text, prompt_audio_key, use_tts, correct_answer, options, hint, sort_order, created_at, updated_at
		FROM questions
		WHERE id = $1
	`, id).Scan(
		&question.ID, &question.LessonID, &question.Type, &question.PromptText,
		&question.PromptAudioKey, &question.UseTTS, &question.CorrectAnswer,
		&optionsBytes, &question.Hint, &question.SortOrder, &question.CreatedAt, &question.UpdatedAt,
	)
	if err != nil {
		return nil, ErrQuestionNotFound
	}

	if err := json.Unmarshal(optionsBytes, &question.Options); err != nil {
		question.Options = []string{}
	}

	return &question, nil
}

// ListQuestionsByLesson retrieves all questions for a lesson.
func ListQuestionsByLesson(ctx context.Context, lessonID uuid.UUID) ([]*Question, error) {
	rows, err := db.Query(ctx, `
		SELECT id, lesson_id, type, prompt_text, prompt_audio_key, use_tts, correct_answer, options, hint, sort_order, created_at, updated_at
		FROM questions
		WHERE lesson_id = $1
		ORDER BY sort_order ASC
	`, lessonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*Question
	for rows.Next() {
		var question Question
		var optionsBytes []byte
		if err := rows.Scan(
			&question.ID, &question.LessonID, &question.Type, &question.PromptText,
			&question.PromptAudioKey, &question.UseTTS, &question.CorrectAnswer,
			&optionsBytes, &question.Hint, &question.SortOrder, &question.CreatedAt, &question.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(optionsBytes, &question.Options); err != nil {
			question.Options = []string{}
		}
		questions = append(questions, &question)
	}

	if questions == nil {
		questions = []*Question{}
	}

	return questions, rows.Err()
}

// UpdateQuestion modifies an existing question.
func UpdateQuestion(ctx context.Context, id uuid.UUID, params UpdateQuestionParams) (*Question, error) {
	existing, err := FindQuestionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	qType := existing.Type
	promptText := existing.PromptText
	useTTS := existing.UseTTS
	correctAnswer := existing.CorrectAnswer
	options := existing.Options
	hint := existing.Hint
	sortOrder := existing.SortOrder

	if params.Type != nil {
		qType = *params.Type
	}
	if params.PromptText != nil {
		promptText = *params.PromptText
	}
	if params.UseTTS != nil {
		useTTS = *params.UseTTS
	}
	if params.CorrectAnswer != nil {
		correctAnswer = *params.CorrectAnswer
	}
	if params.Options != nil {
		options = params.Options
	}
	if params.Hint != nil {
		hint = *params.Hint
	}
	if params.SortOrder != nil {
		sortOrder = *params.SortOrder
	}

	optionsJSON, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}

	var question Question
	var optionsBytes []byte
	err = db.QueryRow(ctx, `
		UPDATE questions
		SET type = $2, prompt_text = $3, use_tts = $4, correct_answer = $5, options = $6, hint = $7, sort_order = $8, updated_at = NOW()
		WHERE id = $1
		RETURNING id, lesson_id, type, prompt_text, prompt_audio_key, use_tts, correct_answer, options, hint, sort_order, created_at, updated_at
	`, id, qType, promptText, useTTS, correctAnswer, optionsJSON, hint, sortOrder).Scan(
		&question.ID, &question.LessonID, &question.Type, &question.PromptText,
		&question.PromptAudioKey, &question.UseTTS, &question.CorrectAnswer,
		&optionsBytes, &question.Hint, &question.SortOrder, &question.CreatedAt, &question.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(optionsBytes, &question.Options); err != nil {
		question.Options = []string{}
	}

	return &question, nil
}

// UpdateQuestionAudioKey updates the audio key for a question.
func UpdateQuestionAudioKey(ctx context.Context, id uuid.UUID, audioKey string) error {
	result, err := db.Exec(ctx, `
		UPDATE questions SET prompt_audio_key = $2, updated_at = NOW() WHERE id = $1
	`, id, audioKey)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrQuestionNotFound
	}
	return nil
}

// DeleteQuestion removes a question by its ID.
func DeleteQuestion(ctx context.Context, id uuid.UUID) error {
	result, err := db.Exec(ctx, `DELETE FROM questions WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrQuestionNotFound
	}
	return nil
}

// GetLanguageIDByQuestion retrieves the language ID for a question by joining through lesson and unit.
func GetLanguageIDByQuestion(ctx context.Context, questionID uuid.UUID) (uuid.UUID, error) {
	var languageID uuid.UUID
	err := db.QueryRow(ctx, `
		SELECT u.language_id
		FROM questions q
		JOIN lessons l ON q.lesson_id = l.id
		JOIN units u ON l.unit_id = u.id
		WHERE q.id = $1
	`, questionID).Scan(&languageID)
	if err != nil {
		return uuid.Nil, ErrQuestionNotFound
	}
	return languageID, nil
}
