package content

import (
	"context"
	"errors"

	"encore.dev/beta/errs"
	"github.com/google/uuid"
)

// ========== Public Endpoints ==========

// ListLessons returns all lessons for a unit.
//
//encore:api auth method=GET path=/units/:unitID/lessons
func ListLessons(ctx context.Context, unitID string) (*LessonsResponse, error) {
	uid, err := uuid.Parse(unitID)
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid unit ID"}
	}

	lessons, err := ListLessonsByUnit(ctx, uid)
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to list lessons"}
	}

	return &LessonsResponse{Lessons: lessons}, nil
}

// GetLesson returns a specific lesson by ID.
//
//encore:api auth method=GET path=/lessons/:id
func GetLesson(ctx context.Context, id string) (*LessonResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid lesson ID"}
	}

	lesson, err := FindLessonByID(ctx, uid)
	if err != nil {
		if errors.Is(err, ErrLessonNotFound) {
			return nil, &errs.Error{Code: errs.NotFound, Message: "lesson not found"}
		}
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to fetch lesson"}
	}

	return &LessonResponse{Lesson: lesson}, nil
}

// ========== Admin Endpoints ==========

// CreateLessonAdmin creates a new lesson (admin only).
//
//encore:api auth method=POST path=/admin/lessons tag:admin
func CreateLessonAdmin(ctx context.Context, params *CreateLessonParams) (*LessonResponse, error) {
	if params.UnitID == "" {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "unit_id is required"}
	}
	if params.Title == "" {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "title is required"}
	}

	uid, err := uuid.Parse(params.UnitID)
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid unit ID"}
	}

	lesson, err := CreateLesson(ctx, uid, *params)
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to create lesson"}
	}

	return &LessonResponse{Lesson: lesson}, nil
}

// UpdateLessonAdmin updates an existing lesson (admin only).
//
//encore:api auth method=PUT path=/admin/lessons/:id tag:admin
func UpdateLessonAdmin(ctx context.Context, id string, params *UpdateLessonParams) (*LessonResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid lesson ID"}
	}

	lesson, err := UpdateLesson(ctx, uid, *params)
	if err != nil {
		if errors.Is(err, ErrLessonNotFound) {
			return nil, &errs.Error{Code: errs.NotFound, Message: "lesson not found"}
		}
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to update lesson"}
	}

	return &LessonResponse{Lesson: lesson}, nil
}

// DeleteLessonResponse is returned after deleting a lesson.
type DeleteLessonResponse struct {
	Success bool `json:"success"`
}

// DeleteLessonAdmin deletes a lesson (admin only).
//
//encore:api auth method=DELETE path=/admin/lessons/:id tag:admin
func DeleteLessonAdmin(ctx context.Context, id string) (*DeleteLessonResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid lesson ID"}
	}

	err = DeleteLesson(ctx, uid)
	if err != nil {
		if errors.Is(err, ErrLessonNotFound) {
			return nil, &errs.Error{Code: errs.NotFound, Message: "lesson not found"}
		}
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to delete lesson"}
	}

	return &DeleteLessonResponse{Success: true}, nil
}
