package providers

import (
	"context"
	"sync"

	"github.com/MicahParks/keyfunc/v3"
)

// JWKSManager manages JWKS fetching and caching for multiple providers
type JWKSManager struct {
	mu       sync.RWMutex
	keyfuncs map[string]keyfunc.Keyfunc
}

// NewJWKSManager creates a new JWKS manager
func NewJWKSManager() *JWKSManager {
	return &JWKSManager{
		keyfuncs: make(map[string]keyfunc.Keyfunc),
	}
}

// GetKeyfunc returns a cached keyfunc or creates a new one
func (m *JWKSManager) GetKeyfunc(ctx context.Context, name, jwksURL string) (keyfunc.Keyfunc, error) {
	m.mu.RLock()
	kf, exists := m.keyfuncs[name]
	m.mu.RUnlock()

	if exists {
		return kf, nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Double-check after acquiring write lock
	if kf, exists = m.keyfuncs[name]; exists {
		return kf, nil
	}

	// Create new keyfunc with automatic background refresh
	kf, err := keyfunc.NewDefaultCtx(ctx, []string{jwksURL})
	if err != nil {
		return nil, err
	}

	m.keyfuncs[name] = kf
	return kf, nil
}

// Global JWKS manager instance
var globalJWKSManager = NewJWKSManager()

// GetGlobalJWKSManager returns the global JWKS manager
func GetGlobalJWKSManager() *JWKSManager {
	return globalJWKSManager
}
