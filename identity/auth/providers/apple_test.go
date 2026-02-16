package providers

import (
	"testing"
)

func TestAppleProvider_Name(t *testing.T) {
	provider := NewAppleProvider("com.example.app", "com.example.service")

	if provider.Name() != "apple" {
		t.Errorf("expected provider name 'apple', got %q", provider.Name())
	}
}

func TestAppleProvider_IsValidAudience(t *testing.T) {
	bundleID := "com.example.app"
	serviceID := "com.example.service"
	provider := NewAppleProvider(bundleID, serviceID)

	tests := []struct {
		name      string
		audiences []string
		valid     bool
	}{
		{
			name:      "bundle ID",
			audiences: []string{bundleID},
			valid:     true,
		},
		{
			name:      "service ID",
			audiences: []string{serviceID},
			valid:     true,
		},
		{
			name:      "wrong bundle ID",
			audiences: []string{"com.other.app"},
			valid:     false,
		},
		{
			name:      "empty audience",
			audiences: []string{},
			valid:     false,
		},
		{
			name:      "multiple with valid",
			audiences: []string{"other", bundleID},
			valid:     true,
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

func TestAppleProvider_NoServiceID(t *testing.T) {
	bundleID := "com.example.app"
	provider := NewAppleProvider(bundleID, "") // No service ID

	// Bundle ID should still work
	if !provider.isValidAudience([]string{bundleID}) {
		t.Error("expected bundle ID to be valid even without service ID")
	}

	// Empty service ID should not match anything
	if provider.isValidAudience([]string{"com.example.service"}) {
		t.Error("expected random service ID to be invalid when no service ID configured")
	}
}
