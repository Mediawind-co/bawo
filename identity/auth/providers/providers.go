package providers

import (
	"context"
	"errors"
)

var (
	ErrInvalidToken    = errors.New("invalid token")
	ErrTokenExpired    = errors.New("token expired")
	ErrUnknownProvider = errors.New("unknown provider")
	ErrInvalidIssuer   = errors.New("invalid token issuer")
	ErrInvalidAudience = errors.New("invalid token audience")
)

// TokenClaims contains the verified claims from a JWT
type TokenClaims struct {
	Subject   string // Provider's unique user ID (sub claim)
	Email     string
	Name      string
	AvatarURL string
	Provider  string // "google" or "apple"
}

// Provider interface for JWT verification
type Provider interface {
	// Verify validates a JWT and returns the extracted claims
	Verify(ctx context.Context, token string) (*TokenClaims, error)
	// Name returns the provider identifier
	Name() string
}

// Registry holds all registered providers
type Registry struct {
	providers map[string]Provider
}

// NewRegistry creates a new provider registry
func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]Provider),
	}
}

// Register adds a provider to the registry
func (r *Registry) Register(p Provider) {
	r.providers[p.Name()] = p
}

// Get retrieves a provider by name
func (r *Registry) Get(name string) (Provider, bool) {
	p, ok := r.providers[name]
	return p, ok
}

// VerifyAny attempts to verify the token with each provider
// Returns on first successful verification
func (r *Registry) VerifyAny(ctx context.Context, token string) (*TokenClaims, error) {
	var lastErr error
	for _, p := range r.providers {
		claims, err := p.Verify(ctx, token)
		if err == nil {
			return claims, nil
		}
		lastErr = err
	}
	if lastErr != nil {
		return nil, lastErr
	}
	return nil, ErrInvalidToken
}
