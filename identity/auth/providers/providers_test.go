package providers

import (
	"context"
	"errors"
	"testing"
)

// mockProvider is a test provider implementation
type mockProvider struct {
	name      string
	claims    *TokenClaims
	err       error
	callCount int
}

func (m *mockProvider) Name() string {
	return m.name
}

func (m *mockProvider) Verify(ctx context.Context, token string) (*TokenClaims, error) {
	m.callCount++
	if m.err != nil {
		return nil, m.err
	}
	return m.claims, nil
}

func TestRegistry_Register(t *testing.T) {
	registry := NewRegistry()

	provider := &mockProvider{name: "test"}
	registry.Register(provider)

	got, ok := registry.Get("test")
	if !ok {
		t.Fatal("expected to find registered provider")
	}
	if got != provider {
		t.Error("expected to get the same provider instance")
	}
}

func TestRegistry_GetNotFound(t *testing.T) {
	registry := NewRegistry()

	_, ok := registry.Get("nonexistent")
	if ok {
		t.Error("expected not to find unregistered provider")
	}
}

func TestRegistry_VerifyAny_Success(t *testing.T) {
	registry := NewRegistry()

	expectedClaims := &TokenClaims{
		Subject:  "user-123",
		Email:    "test@example.com",
		Provider: "test",
	}

	provider := &mockProvider{
		name:   "test",
		claims: expectedClaims,
	}
	registry.Register(provider)

	claims, err := registry.VerifyAny(context.Background(), "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if claims.Subject != expectedClaims.Subject {
		t.Errorf("subject mismatch: got %s, want %s", claims.Subject, expectedClaims.Subject)
	}
	if claims.Email != expectedClaims.Email {
		t.Errorf("email mismatch: got %s, want %s", claims.Email, expectedClaims.Email)
	}
}

func TestRegistry_VerifyAny_FirstSuccessWins(t *testing.T) {
	registry := NewRegistry()

	// First provider will fail
	provider1 := &mockProvider{
		name: "provider1",
		err:  errors.New("verification failed"),
	}

	// Second provider will succeed
	provider2 := &mockProvider{
		name: "provider2",
		claims: &TokenClaims{
			Subject:  "user-456",
			Provider: "provider2",
		},
	}

	registry.Register(provider1)
	registry.Register(provider2)

	claims, err := registry.VerifyAny(context.Background(), "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if claims.Provider != "provider2" {
		t.Errorf("expected provider2 claims, got %s", claims.Provider)
	}
}

func TestRegistry_VerifyAny_AllFail(t *testing.T) {
	registry := NewRegistry()

	provider1 := &mockProvider{
		name: "provider1",
		err:  errors.New("provider1 failed"),
	}

	provider2 := &mockProvider{
		name: "provider2",
		err:  errors.New("provider2 failed"),
	}

	registry.Register(provider1)
	registry.Register(provider2)

	_, err := registry.VerifyAny(context.Background(), "test-token")
	if err == nil {
		t.Fatal("expected error when all providers fail")
	}
}

func TestRegistry_VerifyAny_NoProviders(t *testing.T) {
	registry := NewRegistry()

	_, err := registry.VerifyAny(context.Background(), "test-token")
	if err == nil {
		t.Fatal("expected error with no providers")
	}
	if !errors.Is(err, ErrInvalidToken) {
		t.Errorf("expected ErrInvalidToken, got %v", err)
	}
}
