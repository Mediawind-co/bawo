package user

import (
	"testing"

	"github.com/google/uuid"
)

func TestProvider_Values(t *testing.T) {
	if ProviderGoogle != "google" {
		t.Errorf("ProviderGoogle = %q, want 'google'", ProviderGoogle)
	}
	if ProviderApple != "apple" {
		t.Errorf("ProviderApple = %q, want 'apple'", ProviderApple)
	}
}

func TestRole_Values(t *testing.T) {
	if RoleUser != "user" {
		t.Errorf("RoleUser = %q, want 'user'", RoleUser)
	}
	if RoleAdmin != "admin" {
		t.Errorf("RoleAdmin = %q, want 'admin'", RoleAdmin)
	}
}

func TestCreateUserParams(t *testing.T) {
	params := CreateUserParams{
		Provider:   ProviderGoogle,
		ProviderID: "google-123",
		Email:      "test@example.com",
		Name:       "Test User",
		AvatarURL:  "https://example.com/avatar.jpg",
	}

	if params.Provider != ProviderGoogle {
		t.Errorf("Provider = %q, want %q", params.Provider, ProviderGoogle)
	}
	if params.Email != "test@example.com" {
		t.Errorf("Email = %q, want 'test@example.com'", params.Email)
	}
}

func TestUpdateUserParams(t *testing.T) {
	name := "Updated Name"
	avatar := "https://example.com/new-avatar.jpg"

	params := UpdateUserParams{
		Name:      &name,
		AvatarURL: &avatar,
	}

	if *params.Name != name {
		t.Errorf("Name = %q, want %q", *params.Name, name)
	}
	if *params.AvatarURL != avatar {
		t.Errorf("AvatarURL = %q, want %q", *params.AvatarURL, avatar)
	}
}

func TestUser_Fields(t *testing.T) {
	id := uuid.New()
	user := User{
		ID:         id,
		Provider:   ProviderGoogle,
		ProviderID: "google-123",
		Email:      "test@example.com",
		Name:       "Test User",
		AvatarURL:  "https://example.com/avatar.jpg",
		Role:       RoleUser,
	}

	if user.ID != id {
		t.Errorf("ID mismatch")
	}
	if user.Provider != ProviderGoogle {
		t.Errorf("Provider = %q, want %q", user.Provider, ProviderGoogle)
	}
	if user.Role != RoleUser {
		t.Errorf("Role = %q, want %q", user.Role, RoleUser)
	}
}
