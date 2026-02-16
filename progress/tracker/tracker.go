package tracker

import (
	"context"

	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
	"encore.dev/pubsub"
	"github.com/google/uuid"

	"encore.app/learning/content"
	"encore.app/learning/lesson"
)

// ========== Pub/Sub Subscriber ==========

// Subscribe to lesson completed events.
var _ = pubsub.NewSubscription(
	lesson.LessonCompleted, "track-progress",
	pubsub.SubscriptionConfig[*lesson.LessonCompletedEvent]{
		Handler: HandleLessonCompleted,
	},
)

// HandleLessonCompleted processes lesson completion events.
func HandleLessonCompleted(ctx context.Context, event *lesson.LessonCompletedEvent) error {
	userID, err := uuid.Parse(event.UserID)
	if err != nil {
		return nil // Skip invalid events
	}
	lessonID, err := uuid.Parse(event.LessonID)
	if err != nil {
		return nil
	}

	// Calculate score
	total := event.Correct + event.Incorrect
	score := 0.0
	if total > 0 {
		score = float64(event.Correct) / float64(total) * 100
	}

	// Update progress
	_, err = UpdateProgress(ctx, userID, lessonID, event.XPEarned, score)
	return err
}

// ========== API Endpoints ==========

// GetMyStats returns statistics for the current user.
//
//encore:api auth method=GET path=/stats
func GetMyStats(ctx context.Context) (*UserStatsResponse, error) {
	uid, ok := auth.UserID()
	if !ok {
		return nil, &errs.Error{Code: errs.Unauthenticated, Message: "not authenticated"}
	}

	userID, err := uuid.Parse(string(uid))
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid user ID"}
	}

	stats, err := GetUserStats(ctx, userID)
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to get stats"}
	}

	return stats, nil
}

// GetLessonProgress returns progress for a specific lesson.
//
//encore:api auth method=GET path=/progress/lessons/:lessonID
func GetLessonProgress(ctx context.Context, lessonID string) (*ProgressResponse, error) {
	uid, ok := auth.UserID()
	if !ok {
		return nil, &errs.Error{Code: errs.Unauthenticated, Message: "not authenticated"}
	}

	userID, err := uuid.Parse(string(uid))
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid user ID"}
	}

	lID, err := uuid.Parse(lessonID)
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid lesson ID"}
	}

	progress, err := FindProgressByUserAndLesson(ctx, userID, lID)
	if err != nil {
		// Return empty progress if not found
		return &ProgressResponse{
			Progress: &LessonProgress{
				UserID:   userID,
				LessonID: lID,
				Status:   ProgressStatusNotStarted,
			},
		}, nil
	}

	return &ProgressResponse{Progress: progress}, nil
}

// ListMyProgress lists all progress for the current user.
//
//encore:api auth method=GET path=/progress
func ListMyProgress(ctx context.Context) (*ProgressListResponse, error) {
	uid, ok := auth.UserID()
	if !ok {
		return nil, &errs.Error{Code: errs.Unauthenticated, Message: "not authenticated"}
	}

	userID, err := uuid.Parse(string(uid))
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid user ID"}
	}

	progress, err := ListProgressByUser(ctx, userID)
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to list progress"}
	}

	return &ProgressListResponse{Progress: progress}, nil
}

// ========== Mistakes Endpoints ==========

// ListMistakesParams contains pagination parameters.
type ListMistakesParams struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}

// ListMyMistakes lists all mistakes for the current user.
//
//encore:api auth method=GET path=/mistakes
func ListMyMistakes(ctx context.Context, params *ListMistakesParams) (*MistakesResponse, error) {
	uid, ok := auth.UserID()
	if !ok {
		return nil, &errs.Error{Code: errs.Unauthenticated, Message: "not authenticated"}
	}

	userID, err := uuid.Parse(string(uid))
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid user ID"}
	}

	limit := params.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := params.Offset
	if offset < 0 {
		offset = 0
	}

	mistakes, total, err := ListMistakesByUser(ctx, userID, limit, offset)
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to list mistakes"}
	}

	// Enrich with question details
	result := make([]*MistakeWithQuestion, 0, len(mistakes))
	for _, m := range mistakes {
		mwq := &MistakeWithQuestion{Mistake: *m}

		// Try to get question details
		q, err := content.FindQuestionByID(ctx, m.QuestionID)
		if err == nil {
			mwq.PromptText = q.PromptText
			mwq.Hint = q.Hint
		}

		result = append(result, mwq)
	}

	return &MistakesResponse{Mistakes: result, Total: total}, nil
}

// ListMistakesByLanguageEndpoint lists mistakes for a specific language.
//
//encore:api auth method=GET path=/mistakes/:languageID
func ListMistakesByLanguageEndpoint(ctx context.Context, languageID string, params *ListMistakesParams) (*MistakesResponse, error) {
	uid, ok := auth.UserID()
	if !ok {
		return nil, &errs.Error{Code: errs.Unauthenticated, Message: "not authenticated"}
	}

	userID, err := uuid.Parse(string(uid))
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid user ID"}
	}

	langID, err := uuid.Parse(languageID)
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid language ID"}
	}

	limit := params.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := params.Offset
	if offset < 0 {
		offset = 0
	}

	mistakes, total, err := ListMistakesByLanguage(ctx, userID, langID, limit, offset)
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to list mistakes"}
	}

	// Enrich with question details
	result := make([]*MistakeWithQuestion, 0, len(mistakes))
	for _, m := range mistakes {
		mwq := &MistakeWithQuestion{Mistake: *m}

		q, err := content.FindQuestionByID(ctx, m.QuestionID)
		if err == nil {
			mwq.PromptText = q.PromptText
			mwq.Hint = q.Hint
		}

		result = append(result, mwq)
	}

	return &MistakesResponse{Mistakes: result, Total: total}, nil
}

// MarkReviewedRequest contains the mistake ID to mark.
type MarkReviewedRequest struct {
	MistakeID string `json:"mistake_id"`
}

// MarkReviewedResponse indicates success.
type MarkReviewedResponse struct {
	Success bool `json:"success"`
}

// MarkMistakeAsReviewed marks a mistake as reviewed.
//
//encore:api auth method=POST path=/mistakes/review
func MarkMistakeAsReviewed(ctx context.Context, req *MarkReviewedRequest) (*MarkReviewedResponse, error) {
	_, ok := auth.UserID()
	if !ok {
		return nil, &errs.Error{Code: errs.Unauthenticated, Message: "not authenticated"}
	}

	mistakeID, err := uuid.Parse(req.MistakeID)
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid mistake ID"}
	}

	err = MarkMistakeReviewed(ctx, mistakeID)
	if err != nil {
		return nil, &errs.Error{Code: errs.NotFound, Message: "mistake not found"}
	}

	return &MarkReviewedResponse{Success: true}, nil
}
