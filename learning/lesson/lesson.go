package lesson

import (
	"context"
	"errors"

	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
	"github.com/google/uuid"

	"encore.app/learning/content"
)

// ========== Lesson Session Endpoints ==========

// StartLesson starts a new lesson session or resumes an existing one.
//
//encore:api auth method=POST path=/lessons/:lessonID/start
func StartLesson(ctx context.Context, lessonID string) (*StartSessionResponse, error) {
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

	// Check if lesson exists
	lessonData, err := content.FindLessonByID(ctx, lID)
	if err != nil {
		return nil, &errs.Error{Code: errs.NotFound, Message: "lesson not found"}
	}

	// Check for existing active session
	session, err := FindActiveSession(ctx, userID, lID)
	if err != nil {
		// No active session, create a new one
		session, err = CreateSession(ctx, userID, lID)
		if err != nil {
			return nil, &errs.Error{Code: errs.Internal, Message: "failed to create session"}
		}
	}

	// Get questions for this lesson
	questions, err := content.ListQuestionsByLesson(ctx, lID)
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to fetch questions"}
	}

	// Convert to QuestionInfo (hide correct answers)
	questionInfos := make([]QuestionInfo, len(questions))
	for i, q := range questions {
		questionInfos[i] = QuestionInfo{
			ID:         q.ID,
			Type:       string(q.Type),
			PromptText: q.PromptText,
			HasAudio:   q.PromptAudioKey != "",
			UseTTS:     q.UseTTS,
			Options:    q.Options,
			Hint:       q.Hint,
			SortOrder:  q.SortOrder,
		}
	}

	// Update session with lesson XP (store for later use)
	_ = lessonData.XPReward

	return &StartSessionResponse{
		Session:   session,
		Questions: questionInfos,
	}, nil
}

// SubmitAnswer submits an answer for the current question in a session.
//
//encore:api auth method=POST path=/sessions/:sessionID/answer
func SubmitAnswer(ctx context.Context, sessionID string, req *SubmitAnswerRequest) (*SubmitAnswerResponse, error) {
	uid, ok := auth.UserID()
	if !ok {
		return nil, &errs.Error{Code: errs.Unauthenticated, Message: "not authenticated"}
	}

	userID, err := uuid.Parse(string(uid))
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid user ID"}
	}

	sID, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid session ID"}
	}

	qID, err := uuid.Parse(req.QuestionID)
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid question ID"}
	}

	// Get session
	session, err := FindSessionByID(ctx, sID)
	if err != nil {
		return nil, &errs.Error{Code: errs.NotFound, Message: "session not found"}
	}

	// Verify ownership
	if session.UserID != userID {
		return nil, &errs.Error{Code: errs.PermissionDenied, Message: "not your session"}
	}

	// Check session is in progress
	if session.Status != SessionStatusInProgress {
		return nil, &errs.Error{Code: errs.FailedPrecondition, Message: "session is not in progress"}
	}

	// Get question
	question, err := content.FindQuestionByID(ctx, qID)
	if err != nil {
		return nil, &errs.Error{Code: errs.NotFound, Message: "question not found"}
	}

	// Get lesson for XP info
	lesson, err := content.FindLessonByID(ctx, session.LessonID)
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to fetch lesson"}
	}

	// Evaluate answer
	result := EvaluateAnswer(question, req.Answer, lesson.XPReward)

	// Record answer
	_, err = RecordAnswer(ctx, sID, qID, req.Answer, result.IsCorrect, result.Similarity)
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to record answer"}
	}

	// Update session progress
	_, err = UpdateSessionProgress(ctx, sID, result.IsCorrect, result.XPEarned)
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to update progress"}
	}

	response := &SubmitAnswerResponse{
		IsCorrect:     result.IsCorrect,
		CorrectAnswer: question.CorrectAnswer,
		Similarity:    result.Similarity,
		XPEarned:      result.XPEarned,
	}

	// Add hint for incorrect answers and record the mistake
	if !result.IsCorrect {
		response.Hint = GenerateHint(question, req.Answer)

		// Get language ID and publish mistake event
		languageID, err := content.GetLanguageIDByQuestion(ctx, qID)
		if err == nil {
			_ = PublishMistakeRecorded(
				userID.String(),
				qID.String(),
				languageID.String(),
				req.Answer,
				question.CorrectAnswer,
			)
		}
	}

	return response, nil
}

// CompleteLesson completes a lesson session.
//
//encore:api auth method=POST path=/sessions/:sessionID/complete
func CompleteLesson(ctx context.Context, sessionID string) (*CompleteSessionResponse, error) {
	uid, ok := auth.UserID()
	if !ok {
		return nil, &errs.Error{Code: errs.Unauthenticated, Message: "not authenticated"}
	}

	userID, err := uuid.Parse(string(uid))
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid user ID"}
	}

	sID, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid session ID"}
	}

	// Get session
	session, err := FindSessionByID(ctx, sID)
	if err != nil {
		return nil, &errs.Error{Code: errs.NotFound, Message: "session not found"}
	}

	// Verify ownership
	if session.UserID != userID {
		return nil, &errs.Error{Code: errs.PermissionDenied, Message: "not your session"}
	}

	// Check session is in progress
	if session.Status != SessionStatusInProgress {
		return nil, &errs.Error{Code: errs.FailedPrecondition, Message: "session is not in progress"}
	}

	// Complete the session
	session, err = CompleteSession(ctx, sID)
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to complete session"}
	}

	// Publish completion event
	_ = PublishLessonCompleted(session)

	// Calculate accuracy
	totalCount := session.CorrectCount + session.IncorrectCount
	accuracy := 0.0
	if totalCount > 0 {
		accuracy = float64(session.CorrectCount) / float64(totalCount) * 100
	}

	return &CompleteSessionResponse{
		Session:      session,
		TotalXP:      session.XPEarned,
		Accuracy:     accuracy,
		CorrectCount: session.CorrectCount,
		TotalCount:   totalCount,
	}, nil
}

// AbandonLesson abandons an in-progress lesson session.
//
//encore:api auth method=POST path=/sessions/:sessionID/abandon
func AbandonLesson(ctx context.Context, sessionID string) (*SessionResponse, error) {
	uid, ok := auth.UserID()
	if !ok {
		return nil, &errs.Error{Code: errs.Unauthenticated, Message: "not authenticated"}
	}

	userID, err := uuid.Parse(string(uid))
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid user ID"}
	}

	sID, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid session ID"}
	}

	// Get session
	session, err := FindSessionByID(ctx, sID)
	if err != nil {
		return nil, &errs.Error{Code: errs.NotFound, Message: "session not found"}
	}

	// Verify ownership
	if session.UserID != userID {
		return nil, &errs.Error{Code: errs.PermissionDenied, Message: "not your session"}
	}

	// Abandon the session
	err = AbandonSession(ctx, sID)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) {
			return nil, &errs.Error{Code: errs.FailedPrecondition, Message: "session is not in progress"}
		}
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to abandon session"}
	}

	// Fetch updated session
	session, _ = FindSessionByID(ctx, sID)

	return &SessionResponse{Session: session}, nil
}

// GetSession retrieves a session by ID.
//
//encore:api auth method=GET path=/sessions/:sessionID
func GetSession(ctx context.Context, sessionID string) (*SessionResponse, error) {
	uid, ok := auth.UserID()
	if !ok {
		return nil, &errs.Error{Code: errs.Unauthenticated, Message: "not authenticated"}
	}

	userID, err := uuid.Parse(string(uid))
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid user ID"}
	}

	sID, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid session ID"}
	}

	// Get session
	session, err := FindSessionByID(ctx, sID)
	if err != nil {
		return nil, &errs.Error{Code: errs.NotFound, Message: "session not found"}
	}

	// Verify ownership
	if session.UserID != userID {
		return nil, &errs.Error{Code: errs.PermissionDenied, Message: "not your session"}
	}

	return &SessionResponse{Session: session}, nil
}

// ListSessionsParams contains pagination parameters.
type ListSessionsParams struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}

// ListMySessions lists the current user's lesson sessions.
//
//encore:api auth method=GET path=/sessions
func ListMySessions(ctx context.Context, params *ListSessionsParams) (*SessionsResponse, error) {
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

	sessions, err := ListUserSessions(ctx, userID, limit, offset)
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to list sessions"}
	}

	return &SessionsResponse{Sessions: sessions}, nil
}
