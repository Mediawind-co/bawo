package content

import (
	"context"
	"errors"

	"encore.dev/beta/errs"
	"github.com/google/uuid"
)

// ========== Public Endpoints ==========

// ListQuestions returns all questions for a lesson.
//
//encore:api auth method=GET path=/lessons/:lessonID/questions
func ListQuestions(ctx context.Context, lessonID string) (*QuestionsResponse, error) {
	uid, err := uuid.Parse(lessonID)
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid lesson ID"}
	}

	questions, err := ListQuestionsByLesson(ctx, uid)
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to list questions"}
	}

	return &QuestionsResponse{Questions: questions}, nil
}

// GetQuestion returns a specific question by ID.
//
//encore:api auth method=GET path=/questions/:id
func GetQuestion(ctx context.Context, id string) (*QuestionResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid question ID"}
	}

	question, err := FindQuestionByID(ctx, uid)
	if err != nil {
		if errors.Is(err, ErrQuestionNotFound) {
			return nil, &errs.Error{Code: errs.NotFound, Message: "question not found"}
		}
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to fetch question"}
	}

	return &QuestionResponse{Question: question}, nil
}

// ========== Admin Endpoints ==========

// CreateQuestionAdmin creates a new question (admin only).
//
//encore:api auth method=POST path=/admin/questions tag:admin
func CreateQuestionAdmin(ctx context.Context, params *CreateQuestionParams) (*QuestionResponse, error) {
	if params.LessonID == "" {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "lesson_id is required"}
	}
	if params.PromptText == "" {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "prompt_text is required"}
	}
	if params.CorrectAnswer == "" {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "correct_answer is required"}
	}
	if params.Type == "" {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "type is required"}
	}

	// Validate question type
	switch params.Type {
	case QuestionTypeListenReply, QuestionTypeMultiChoice, QuestionTypeSingleChoice:
		// valid
	default:
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid question type"}
	}

	// Validate options for choice questions
	if (params.Type == QuestionTypeMultiChoice || params.Type == QuestionTypeSingleChoice) && len(params.Options) < 2 {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "choice questions require at least 2 options"}
	}

	uid, err := uuid.Parse(params.LessonID)
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid lesson ID"}
	}

	question, err := CreateQuestion(ctx, uid, *params)
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to create question"}
	}

	return &QuestionResponse{Question: question}, nil
}

// UpdateQuestionAdmin updates an existing question (admin only).
//
//encore:api auth method=PUT path=/admin/questions/:id tag:admin
func UpdateQuestionAdmin(ctx context.Context, id string, params *UpdateQuestionParams) (*QuestionResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid question ID"}
	}

	// Validate question type if provided
	if params.Type != nil {
		switch *params.Type {
		case QuestionTypeListenReply, QuestionTypeMultiChoice, QuestionTypeSingleChoice:
			// valid
		default:
			return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid question type"}
		}
	}

	question, err := UpdateQuestion(ctx, uid, *params)
	if err != nil {
		if errors.Is(err, ErrQuestionNotFound) {
			return nil, &errs.Error{Code: errs.NotFound, Message: "question not found"}
		}
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to update question"}
	}

	return &QuestionResponse{Question: question}, nil
}

// DeleteQuestionResponse is returned after deleting a question.
type DeleteQuestionResponse struct {
	Success bool `json:"success"`
}

// DeleteQuestionAdmin deletes a question (admin only).
//
//encore:api auth method=DELETE path=/admin/questions/:id tag:admin
func DeleteQuestionAdmin(ctx context.Context, id string) (*DeleteQuestionResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid question ID"}
	}

	err = DeleteQuestion(ctx, uid)
	if err != nil {
		if errors.Is(err, ErrQuestionNotFound) {
			return nil, &errs.Error{Code: errs.NotFound, Message: "question not found"}
		}
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to delete question"}
	}

	return &DeleteQuestionResponse{Success: true}, nil
}
