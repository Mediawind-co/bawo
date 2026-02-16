package user

import (
	"context"
	"errors"

	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
	"encore.dev/middleware"
	"github.com/google/uuid"
)

// ========== Middleware ==========

// AdminOnly is a middleware that restricts access to admin users only.
// Apply to endpoints using: tag:admin
//
//encore:middleware target=tag:admin
func AdminOnly(req middleware.Request, next middleware.Next) middleware.Response {
	data := auth.Data()
	if data == nil {
		return middleware.Response{
			Err: &errs.Error{Code: errs.Unauthenticated, Message: "authentication required"},
		}
	}

	// Type assert to get role - we expect a map or struct with Role field
	// Since auth data comes from auth service, we check the role field
	type authData interface {
		GetRole() string
	}

	// For now, we'll check the role from the auth context
	// The auth handler sets this in AuthData
	role := getRoleFromAuthData(data)
	if role != "admin" {
		return middleware.Response{
			Err: &errs.Error{Code: errs.PermissionDenied, Message: "admin access required"},
		}
	}

	return next(req)
}

// getRoleFromAuthData extracts role from auth data
func getRoleFromAuthData(data any) string {
	// AuthData is defined in the auth package, but we can use reflection or type assertion
	// For simplicity, we'll use a type assertion with the expected structure
	if m, ok := data.(interface{ Role() string }); ok {
		return m.Role()
	}
	// Try struct field access via reflection-like approach
	if ad, ok := data.(*struct {
		UserID   string
		Email    string
		Name     string
		Role     string
		Provider string
	}); ok && ad != nil {
		return ad.Role
	}
	return ""
}

// ========== Public Endpoints ==========

// GetCurrentUser returns the currently authenticated user's profile.
//
//encore:api auth method=GET path=/users/me
func GetCurrentUser(ctx context.Context) (*UserResponse, error) {
	uid, ok := auth.UserID()
	if !ok {
		return nil, &errs.Error{Code: errs.Unauthenticated, Message: "not authenticated"}
	}

	id, err := uuid.Parse(string(uid))
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid user ID"}
	}

	user, err := FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, &errs.Error{Code: errs.NotFound, Message: "user not found"}
		}
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to fetch user"}
	}

	return &UserResponse{User: user}, nil
}

// UpdateCurrentUserParams contains the update parameters
type UpdateCurrentUserParams struct {
	Name      *string `json:"name,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

// UpdateCurrentUser updates the currently authenticated user's profile.
//
//encore:api auth method=PATCH path=/users/me
func UpdateCurrentUser(ctx context.Context, params *UpdateCurrentUserParams) (*UserResponse, error) {
	uid, ok := auth.UserID()
	if !ok {
		return nil, &errs.Error{Code: errs.Unauthenticated, Message: "not authenticated"}
	}

	id, err := uuid.Parse(string(uid))
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid user ID"}
	}

	user, err := Update(ctx, id, UpdateUserParams{
		Name:      params.Name,
		AvatarURL: params.AvatarURL,
	})
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, &errs.Error{Code: errs.NotFound, Message: "user not found"}
		}
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to update user"}
	}

	return &UserResponse{User: user}, nil
}

// ========== Admin Endpoints ==========

// GetUser returns a user by ID (admin only).
//
//encore:api auth method=GET path=/admin/users/:id tag:admin
func GetUser(ctx context.Context, id string) (*UserResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid user ID format"}
	}

	user, err := FindByID(ctx, uid)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, &errs.Error{Code: errs.NotFound, Message: "user not found"}
		}
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to fetch user"}
	}

	return &UserResponse{User: user}, nil
}

// ListUsersParams contains pagination parameters
type ListUsersParams struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}

// ListUsers returns a paginated list of users (admin only).
//
//encore:api auth method=GET path=/admin/users tag:admin
func ListUsers(ctx context.Context, params *ListUsersParams) (*UsersResponse, error) {
	limit := params.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	offset := params.Offset
	if offset < 0 {
		offset = 0
	}

	users, total, err := List(ctx, limit, offset)
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to list users"}
	}

	return &UsersResponse{
		Users: users,
		Total: total,
	}, nil
}

// DeleteUserResponse is returned after deleting a user
type DeleteUserResponse struct {
	Success bool `json:"success"`
}

// DeleteUser deletes a user by ID (admin only).
//
//encore:api auth method=DELETE path=/admin/users/:id tag:admin
func DeleteUser(ctx context.Context, id string) (*DeleteUserResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid user ID format"}
	}

	err = Delete(ctx, uid)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, &errs.Error{Code: errs.NotFound, Message: "user not found"}
		}
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to delete user"}
	}

	return &DeleteUserResponse{Success: true}, nil
}

// UpdateUserRoleBody contains the role to update
type UpdateUserRoleBody struct {
	Role Role `json:"role"`
}

// UpdateUserRole updates a user's role (admin only).
//
//encore:api auth method=PATCH path=/admin/users/:id/role tag:admin
func UpdateUserRole(ctx context.Context, id string, body *UpdateUserRoleBody) (*UserResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid user ID format"}
	}

	// Validate role
	if body.Role != RoleUser && body.Role != RoleAdmin {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid role, must be 'user' or 'admin'"}
	}

	user, err := UpdateRole(ctx, uid, body.Role)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, &errs.Error{Code: errs.NotFound, Message: "user not found"}
		}
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to update role"}
	}

	return &UserResponse{User: user}, nil
}
