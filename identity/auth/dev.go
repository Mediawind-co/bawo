package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"os"

	"encore.dev/beta/errs"

	"encore.app/identity/user"
)

// DevLoginRequest contains dev login credentials
type DevLoginRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// DevLoginResponse contains the auth token
type DevLoginResponse struct {
	Token string     `json:"token"`
	User  *user.User `json:"user"`
}

// devTokens stores valid dev tokens (in-memory, cleared on restart)
var devTokens = make(map[string]*DevTokenData)

// DevTokenData holds the user info for a dev token
type DevTokenData struct {
	UserID string
	Email  string
	Name   string
}

// DevLogin creates a dev user and returns a token (development only)
//
//encore:api public method=POST path=/auth/dev-login
func (s *Service) DevLogin(ctx context.Context, req *DevLoginRequest) (*DevLoginResponse, error) {
	// Only allow in development
	if os.Getenv("ENCORE_ENVIRONMENT") != "" {
		return nil, &errs.Error{
			Code:    errs.PermissionDenied,
			Message: "dev login only available in local development",
		}
	}

	if req.Email == "" {
		req.Email = "dev@bawo.test"
	}
	if req.Name == "" {
		req.Name = "Dev User"
	}

	// Find or create dev user
	u, _, err := user.FindOrCreate(ctx, user.CreateUserParams{
		Provider:   user.ProviderGoogle, // Use google as placeholder
		ProviderID: "dev-" + req.Email,
		Email:      req.Email,
		Name:       req.Name,
		AvatarURL:  "",
	})
	if err != nil {
		return nil, &errs.Error{
			Code:    errs.Internal,
			Message: "failed to create dev user",
		}
	}

	// Generate a random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, &errs.Error{
			Code:    errs.Internal,
			Message: "failed to generate token",
		}
	}
	token := "dev_" + hex.EncodeToString(tokenBytes)

	// Store token
	devTokens[token] = &DevTokenData{
		UserID: u.ID.String(),
		Email:  u.Email,
		Name:   u.Name,
	}

	return &DevLoginResponse{
		Token: token,
		User:  u,
	}, nil
}

// ValidateDevToken checks if a dev token is valid
func ValidateDevToken(token string) (*DevTokenData, bool) {
	data, ok := devTokens[token]
	return data, ok
}
