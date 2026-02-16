package providers

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Apple JWKS endpoint
const appleJWKSURL = "https://appleid.apple.com/auth/keys"

// Apple issuer
const appleIssuer = "https://appleid.apple.com"

// AppleProvider implements JWT verification for Apple Sign-In
type AppleProvider struct {
	bundleID    string // iOS App Bundle ID
	serviceID   string // Web Service ID (optional)
	jwksManager *JWKSManager
}

// NewAppleProvider creates a new Apple JWT verification provider
func NewAppleProvider(bundleID, serviceID string) *AppleProvider {
	return &AppleProvider{
		bundleID:    bundleID,
		serviceID:   serviceID,
		jwksManager: GetGlobalJWKSManager(),
	}
}

// Name returns the provider identifier
func (p *AppleProvider) Name() string {
	return "apple"
}

// AppleClaims represents the claims in an Apple ID token
type AppleClaims struct {
	jwt.RegisteredClaims
	Email          string `json:"email"`
	EmailVerified  any    `json:"email_verified"` // Can be bool or string "true"/"false"
	IsPrivateEmail any    `json:"is_private_email"`
	RealUserStatus int    `json:"real_user_status"`
	NonceSupported bool   `json:"nonce_supported"`
	AuthTime       int64  `json:"auth_time"`
}

// Verify validates an Apple ID token and extracts claims
func (p *AppleProvider) Verify(ctx context.Context, tokenString string) (*TokenClaims, error) {
	// Get keyfunc for Apple JWKS
	kf, err := p.jwksManager.GetKeyfunc(ctx, "apple", appleJWKSURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get Apple JWKS: %w", err)
	}

	// Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenString, &AppleClaims{}, kf.Keyfunc,
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*AppleClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// Verify issuer
	if claims.Issuer != appleIssuer {
		return nil, ErrInvalidIssuer
	}

	// Verify audience (can be Bundle ID or Service ID)
	if !p.isValidAudience(claims.Audience) {
		return nil, ErrInvalidAudience
	}

	// Verify expiration
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, ErrTokenExpired
	}

	return &TokenClaims{
		Subject:   claims.Subject,
		Email:     claims.Email,
		Name:      "", // Apple only provides name on first sign-in via a separate parameter
		AvatarURL: "", // Apple doesn't provide avatar
		Provider:  "apple",
	}, nil
}

func (p *AppleProvider) isValidAudience(audiences []string) bool {
	for _, aud := range audiences {
		if aud == p.bundleID || (p.serviceID != "" && aud == p.serviceID) {
			return true
		}
	}
	return false
}
