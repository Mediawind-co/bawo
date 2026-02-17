package adminauth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"encore.dev/beta/errs"
	"encore.dev/storage/sqldb"
	"golang.org/x/crypto/bcrypt"
)

// Database connection
var db = sqldb.NewDatabase("adminauth", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})

// Admin represents an admin user
type Admin struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	IsActive     bool      `json:"is_active"`
	IsSuperadmin bool      `json:"is_superadmin"`
	CreatedAt    time.Time `json:"created_at"`
}

// LoginRequest contains admin login credentials
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse contains the admin token and info
type LoginResponse struct {
	Token string `json:"token"`
	Admin *Admin `json:"admin"`
}

// Admin session tokens (in-memory for simplicity)
var adminTokens = make(map[string]*Admin)

// Login authenticates an admin with username and password
//
//encore:api public method=POST path=/admin/auth/login
func Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "username and password required"}
	}

	// Find admin by username
	var admin Admin
	var passwordHash string
	err := db.QueryRow(ctx, `
		SELECT id, username, email, name, password_hash, is_active, is_superadmin, created_at
		FROM admins WHERE username = $1
	`, req.Username).Scan(
		&admin.ID, &admin.Username, &admin.Email, &admin.Name,
		&passwordHash, &admin.IsActive, &admin.IsSuperadmin, &admin.CreatedAt,
	)
	if err == sqldb.ErrNoRows {
		return nil, &errs.Error{Code: errs.Unauthenticated, Message: "invalid username or password"}
	}
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "database error"}
	}

	// Check if active
	if !admin.IsActive {
		return nil, &errs.Error{Code: errs.PermissionDenied, Message: "account is disabled"}
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		return nil, &errs.Error{Code: errs.Unauthenticated, Message: "invalid username or password"}
	}

	// Generate token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to generate token"}
	}
	token := "admin_" + hex.EncodeToString(tokenBytes)

	// Store token
	adminTokens[token] = &admin

	return &LoginResponse{
		Token: token,
		Admin: &admin,
	}, nil
}

// ValidateTokenRequest contains the token to validate
type ValidateTokenRequest struct {
	Token string `json:"token"`
}

// ValidateTokenResponse contains the admin info if valid
type ValidateTokenResponse struct {
	Valid bool   `json:"valid"`
	Admin *Admin `json:"admin,omitempty"`
}

// ValidateToken checks if an admin token is valid
//
//encore:api public method=POST path=/admin/auth/validate
func ValidateToken(ctx context.Context, req *ValidateTokenRequest) (*ValidateTokenResponse, error) {
	admin, ok := adminTokens[req.Token]
	if !ok {
		return &ValidateTokenResponse{Valid: false}, nil
	}
	return &ValidateTokenResponse{Valid: true, Admin: admin}, nil
}

// LogoutRequest contains the token to invalidate
type LogoutRequest struct {
	Token string `json:"token"`
}

// LogoutResponse indicates success
type LogoutResponse struct {
	Success bool `json:"success"`
}

// Logout invalidates an admin token
//
//encore:api public method=POST path=/admin/auth/logout
func Logout(ctx context.Context, req *LogoutRequest) (*LogoutResponse, error) {
	delete(adminTokens, req.Token)
	return &LogoutResponse{Success: true}, nil
}

// GetAdminByToken returns the admin for a valid token (internal use)
func GetAdminByToken(token string) (*Admin, bool) {
	admin, ok := adminTokens[token]
	return admin, ok
}
