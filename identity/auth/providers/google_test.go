package providers

import (
	"testing"
)

func TestGoogleProvider_Name(t *testing.T) {
	provider := NewGoogleProvider("test-client-id")

	if provider.Name() != "google" {
		t.Errorf("expected provider name 'google', got %q", provider.Name())
	}
}

func TestGoogleProvider_IsValidIssuer(t *testing.T) {
	provider := NewGoogleProvider("test-client-id")

	tests := []struct {
		issuer string
		valid  bool
	}{
		{"accounts.google.com", true},
		{"https://accounts.google.com", true},
		{"https://invalid.google.com", false},
		{"accounts.apple.com", false},
		{"", false},
		{"google.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.issuer, func(t *testing.T) {
			result := provider.isValidIssuer(tt.issuer)
			if result != tt.valid {
				t.Errorf("isValidIssuer(%q) = %v, want %v", tt.issuer, result, tt.valid)
			}
		})
	}
}

func TestGoogleProvider_IsValidAudience(t *testing.T) {
	clientID := "my-client-id.apps.googleusercontent.com"
	provider := NewGoogleProvider(clientID)

	tests := []struct {
		name      string
		audiences []string
		valid     bool
	}{
		{
			name:      "correct client ID",
			audiences: []string{clientID},
			valid:     true,
		},
		{
			name:      "wrong client ID",
			audiences: []string{"wrong-client-id"},
			valid:     false,
		},
		{
			name:      "empty audience",
			audiences: []string{},
			valid:     false,
		},
		{
			name:      "multiple audiences with correct one",
			audiences: []string{"other-id", clientID},
			valid:     true,
		},
		{
			name:      "multiple wrong audiences",
			audiences: []string{"wrong1", "wrong2"},
			valid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := provider.isValidAudience(tt.audiences)
			if result != tt.valid {
				t.Errorf("isValidAudience(%v) = %v, want %v", tt.audiences, result, tt.valid)
			}
		})
	}
}
