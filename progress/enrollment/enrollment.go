package enrollment

import (
	"context"
	"errors"

	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
	"github.com/google/uuid"

	"encore.app/learning/language"
)

// Enroll enrolls the current user in a language.
//
//encore:api auth method=POST path=/enroll
func Enroll(ctx context.Context, req *EnrollRequest) (*EnrollmentResponse, error) {
	uid, ok := auth.UserID()
	if !ok {
		return nil, &errs.Error{Code: errs.Unauthenticated, Message: "not authenticated"}
	}

	userID, err := uuid.Parse(string(uid))
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid user ID"}
	}

	langID, err := uuid.Parse(req.LanguageID)
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid language ID"}
	}

	// Verify language exists and is active
	lang, err := language.FindByID(ctx, langID)
	if err != nil {
		return nil, &errs.Error{Code: errs.NotFound, Message: "language not found"}
	}
	if !lang.IsActive {
		return nil, &errs.Error{Code: errs.FailedPrecondition, Message: "language is not available"}
	}

	// Create enrollment
	enrollment, err := Create(ctx, userID, langID)
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to enroll"}
	}

	return &EnrollmentResponse{Enrollment: enrollment}, nil
}

// Unenroll removes the current user from a language.
//
//encore:api auth method=DELETE path=/enroll/:languageID
func UnenrollFromLanguage(ctx context.Context, languageID string) (*EnrollmentResponse, error) {
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

	// Get enrollment before unenrolling
	enrollment, err := FindByUserAndLanguage(ctx, userID, langID)
	if err != nil {
		return nil, &errs.Error{Code: errs.NotFound, Message: "not enrolled in this language"}
	}

	// Unenroll
	err = Unenroll(ctx, userID, langID)
	if err != nil {
		if errors.Is(err, ErrEnrollmentNotFound) {
			return nil, &errs.Error{Code: errs.NotFound, Message: "not enrolled in this language"}
		}
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to unenroll"}
	}

	enrollment.IsActive = false
	return &EnrollmentResponse{Enrollment: enrollment}, nil
}

// ListEnrollments lists all languages the current user is enrolled in.
//
//encore:api auth method=GET path=/enrollments
func ListEnrollments(ctx context.Context) (*EnrollmentsResponse, error) {
	uid, ok := auth.UserID()
	if !ok {
		return nil, &errs.Error{Code: errs.Unauthenticated, Message: "not authenticated"}
	}

	userID, err := uuid.Parse(string(uid))
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid user ID"}
	}

	enrollments, err := ListByUser(ctx, userID)
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to list enrollments"}
	}

	// Enrich with language details
	result := make([]*EnrollmentWithLanguage, 0, len(enrollments))
	for _, e := range enrollments {
		lang, err := language.FindByID(ctx, e.LanguageID)
		if err != nil {
			continue // Skip if language not found
		}

		result = append(result, &EnrollmentWithLanguage{
			Enrollment:    *e,
			LanguageName:  lang.Name,
			LanguageCode:  lang.Code,
			LanguageEmoji: lang.FlagEmoji,
		})
	}

	return &EnrollmentsResponse{Enrollments: result}, nil
}

// CheckEnrollment checks if the current user is enrolled in a language.
//
//encore:api auth method=GET path=/enroll/:languageID/check
func CheckEnrollment(ctx context.Context, languageID string) (*EnrollmentCheckResponse, error) {
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

	enrolled, err := IsEnrolled(ctx, userID, langID)
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to check enrollment"}
	}

	return &EnrollmentCheckResponse{Enrolled: enrolled}, nil
}

// EnrollmentCheckResponse indicates enrollment status.
type EnrollmentCheckResponse struct {
	Enrolled bool `json:"enrolled"`
}
