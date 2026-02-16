package auth

import (
	"encore.dev/beta/errs"
)

// Common authentication errors
var (
	ErrUnauthenticated = &errs.Error{
		Code:    errs.Unauthenticated,
		Message: "authentication required",
	}

	ErrInvalidToken = &errs.Error{
		Code:    errs.Unauthenticated,
		Message: "invalid or expired token",
	}

	ErrUserNotFound = &errs.Error{
		Code:    errs.NotFound,
		Message: "user not found",
	}

	ErrForbidden = &errs.Error{
		Code:    errs.PermissionDenied,
		Message: "insufficient permissions",
	}
)

// NewAuthError creates a new authentication error with custom message
func NewAuthError(message string) *errs.Error {
	return &errs.Error{
		Code:    errs.Unauthenticated,
		Message: message,
	}
}

// NewForbiddenError creates a new permission denied error
func NewForbiddenError(message string) *errs.Error {
	return &errs.Error{
		Code:    errs.PermissionDenied,
		Message: message,
	}
}
