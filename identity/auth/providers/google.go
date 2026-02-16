package providers

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Google JWKS endpoint
const googleJWKSURL = "https://www.googleapis.com/oauth2/v3/certs"

// Valid Google issuers
var googleIssuers = []string{
	"accounts.google.com",
	"https://accounts.google.com",
}

// GoogleProvider implements JWT verification for Google Sign-In
type GoogleProvider struct {
	clientID    string
	jwksManager *JWKSManager
}

// NewGoogleProvider creates a new Google JWT verification provider
func NewGoogleProvider(clientID string) *GoogleProvider {
	return &GoogleProvider{
		clientID:    clientID,
		jwksManager: GetGlobalJWKSManager(),
	}
}

// Name returns the provider identifier
func (p *GoogleProvider) Name() string {
	return "google"
}

// GoogleClaims represents the claims in a Google ID token
type GoogleClaims struct {
	jwt.RegisteredClaims
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
}

// Verify validates a Google ID token and extracts claims
func (p *GoogleProvider) Verify(ctx context.Context, tokenString string) (*TokenClaims, error) {
	// Get keyfunc for Google JWKS
	kf, err := p.jwksManager.GetKeyfunc(ctx, "google", googleJWKSURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get Google JWKS: %w", err)
	}

	// Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenString, &GoogleClaims{}, kf.Keyfunc,
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*GoogleClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// Verify issuer
	if !p.isValidIssuer(claims.Issuer) {
		return nil, ErrInvalidIssuer
	}

	// Verify audience
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
		Name:      claims.Name,
		AvatarURL: claims.Picture,
		Provider:  "google",
	}, nil
}

func (p *GoogleProvider) isValidIssuer(issuer string) bool {
	for _, valid := range googleIssuers {
		if issuer == valid {
			return true
		}
	}
	return false
}

func (p *GoogleProvider) isValidAudience(audiences []string) bool {
	for _, aud := range audiences {
		if aud == p.clientID {
			return true
		}
	}
	return false
}
