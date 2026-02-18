package auth

import (
	"context"
	"strings"

	"encore.dev/beta/auth"
	"encore.dev/beta/errs"

	"encore.app/admin/adminauth"
	"encore.app/identity/auth/providers"
	"encore.app/identity/user"
)

// Secrets for OAuth providers
var secrets struct {
	GoogleClientID string
	AppleBundleID  string
	AppleServiceID string
}

// AuthData contains the authenticated user information.
// This is accessible via auth.Data() in all authenticated endpoints.
type AuthData struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	Provider string `json:"provider"`
}

// AuthParams defines the authentication parameters from the request
type AuthParams struct {
	Authorization string `header:"Authorization"`
}

// Service holds dependencies for the auth service
//
//encore:service
type Service struct {
	registry *providers.Registry
}

// initService initializes the auth service with providers
func initService() (*Service, error) {
	registry := providers.NewRegistry()

	// Register Google provider
	if secrets.GoogleClientID != "" {
		googleProvider := providers.NewGoogleProvider(secrets.GoogleClientID)
		registry.Register(googleProvider)
	}

	// Register Apple provider
	if secrets.AppleBundleID != "" {
		appleProvider := providers.NewAppleProvider(secrets.AppleBundleID, secrets.AppleServiceID)
		registry.Register(appleProvider)
	}

	return &Service{
		registry: registry,
	}, nil
}

// AuthHandler validates JWT tokens from Google or Apple Sign-In
// and returns the authenticated user's data.
//
//encore:authhandler
func (s *Service) AuthHandler(ctx context.Context, params *AuthParams) (auth.UID, *AuthData, error) {
	// Extract token from Authorization header
	token := extractBearerToken(params.Authorization)
	if token == "" {
		return "", nil, ErrUnauthenticated
	}

	// Check for dev token first (only works in local development)
	if strings.HasPrefix(token, "dev_") {
		return s.handleDevToken(ctx, token)
	}

	// Check for admin token
	if strings.HasPrefix(token, "admin_") {
		return s.handleAdminToken(ctx, token)
	}

	// Verify token with providers (tries each registered provider)
	claims, err := s.registry.VerifyAny(ctx, token)
	if err != nil {
		return "", nil, &errs.Error{
			Code:    errs.Unauthenticated,
			Message: err.Error(),
		}
	}

	// Find or create user in database
	u, _, err := user.FindOrCreate(ctx, user.CreateUserParams{
		Provider:   user.Provider(claims.Provider),
		ProviderID: claims.Subject,
		Email:      claims.Email,
		Name:       claims.Name,
		AvatarURL:  claims.AvatarURL,
	})
	if err != nil {
		return "", nil, &errs.Error{
			Code:    errs.Internal,
			Message: "failed to process user",
		}
	}

	// Build auth data
	authData := &AuthData{
		UserID:   u.ID.String(),
		Email:    u.Email,
		Name:     u.Name,
		Role:     string(u.Role),
		Provider: claims.Provider,
	}

	return auth.UID(u.ID.String()), authData, nil
}

// extractBearerToken extracts the token from "Bearer <token>" format
func extractBearerToken(header string) string {
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

// handleDevToken validates dev tokens for local development
func (s *Service) handleDevToken(ctx context.Context, token string) (auth.UID, *AuthData, error) {
	devData, ok := ValidateDevToken(token)
	if !ok {
		return "", nil, ErrInvalidToken
	}

	// Return auth data directly from dev token (user was created during dev login)
	authData := &AuthData{
		UserID:   devData.UserID,
		Email:    devData.Email,
		Name:     devData.Name,
		Role:     "user", // Default role for dev users
		Provider: "dev",
	}

	return auth.UID(devData.UserID), authData, nil
}

// handleAdminToken validates admin tokens
func (s *Service) handleAdminToken(ctx context.Context, token string) (auth.UID, *AuthData, error) {
	admin, ok := adminauth.GetAdminByToken(token)
	if !ok {
		return "", nil, ErrInvalidToken
	}

	// Return auth data from admin token
	authData := &AuthData{
		UserID:   admin.ID,
		Email:    admin.Email,
		Name:     admin.Name,
		Role:     "admin",
		Provider: "admin",
	}

	return auth.UID(admin.ID), authData, nil
}
