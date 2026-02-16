package auth

import (
	"testing"
)

func TestExtractBearerToken(t *testing.T) {
	tests := []struct {
		name     string
		header   string
		expected string
	}{
		{
			name:     "valid bearer token",
			header:   "Bearer abc123xyz",
			expected: "abc123xyz",
		},
		{
			name:     "bearer lowercase",
			header:   "bearer abc123xyz",
			expected: "abc123xyz",
		},
		{
			name:     "empty header",
			header:   "",
			expected: "",
		},
		{
			name:     "no bearer prefix",
			header:   "abc123xyz",
			expected: "",
		},
		{
			name:     "only bearer",
			header:   "Bearer",
			expected: "",
		},
		{
			name:     "bearer with extra spaces",
			header:   "Bearer   abc123xyz  ",
			expected: "abc123xyz",
		},
		{
			name:     "basic auth (not bearer)",
			header:   "Basic dXNlcjpwYXNz",
			expected: "",
		},
		{
			name:     "mixed case bearer",
			header:   "BEARER token123",
			expected: "token123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractBearerToken(tt.header)
			if result != tt.expected {
				t.Errorf("extractBearerToken(%q) = %q, want %q", tt.header, result, tt.expected)
			}
		})
	}
}
