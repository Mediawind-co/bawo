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

// ========== Admin Management Endpoints ==========

// ListAdminsResponse contains the list of admins
type ListAdminsResponse struct {
	Admins []*Admin `json:"admins"`
	Total  int      `json:"total"`
}

// ListAdmins returns all admin users
//
//encore:api auth method=GET path=/admin/admins tag:admin
func ListAdmins(ctx context.Context) (*ListAdminsResponse, error) {
	rows, err := db.Query(ctx, `
		SELECT id, username, email, name, is_active, is_superadmin, created_at
		FROM admins ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to list admins"}
	}
	defer rows.Close()

	var admins []*Admin
	for rows.Next() {
		var admin Admin
		if err := rows.Scan(&admin.ID, &admin.Username, &admin.Email, &admin.Name,
			&admin.IsActive, &admin.IsSuperadmin, &admin.CreatedAt); err != nil {
			return nil, &errs.Error{Code: errs.Internal, Message: "failed to scan admin"}
		}
		admins = append(admins, &admin)
	}

	return &ListAdminsResponse{
		Admins: admins,
		Total:  len(admins),
	}, nil
}

// GetAdminResponse contains a single admin
type GetAdminResponse struct {
	Admin *Admin `json:"admin"`
}

// GetAdmin returns a specific admin by ID
//
//encore:api auth method=GET path=/admin/admins/:id tag:admin
func GetAdmin(ctx context.Context, id string) (*GetAdminResponse, error) {
	var admin Admin
	err := db.QueryRow(ctx, `
		SELECT id, username, email, name, is_active, is_superadmin, created_at
		FROM admins WHERE id = $1
	`, id).Scan(&admin.ID, &admin.Username, &admin.Email, &admin.Name,
		&admin.IsActive, &admin.IsSuperadmin, &admin.CreatedAt)

	if err == sqldb.ErrNoRows {
		return nil, &errs.Error{Code: errs.NotFound, Message: "admin not found"}
	}
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to get admin"}
	}

	return &GetAdminResponse{Admin: &admin}, nil
}

// CreateAdminRequest contains the data for creating a new admin
type CreateAdminRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// CreateAdminResponse contains the created admin
type CreateAdminResponse struct {
	Admin *Admin `json:"admin"`
}

// CreateAdmin creates a new admin user
//
//encore:api auth method=POST path=/admin/admins tag:admin
func CreateAdmin(ctx context.Context, req *CreateAdminRequest) (*CreateAdminResponse, error) {
	if req.Username == "" || req.Email == "" || req.Password == "" || req.Name == "" {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "all fields are required"}
	}

	if len(req.Password) < 8 {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "password must be at least 8 characters"}
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to hash password"}
	}

	var admin Admin
	err = db.QueryRow(ctx, `
		INSERT INTO admins (username, email, password_hash, name)
		VALUES ($1, $2, $3, $4)
		RETURNING id, username, email, name, is_active, is_superadmin, created_at
	`, req.Username, req.Email, string(hash), req.Name).Scan(
		&admin.ID, &admin.Username, &admin.Email, &admin.Name,
		&admin.IsActive, &admin.IsSuperadmin, &admin.CreatedAt,
	)

	if err != nil {
		// Check for unique constraint violation
		if err.Error() != "" {
			return nil, &errs.Error{Code: errs.AlreadyExists, Message: "username or email already exists"}
		}
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to create admin"}
	}

	return &CreateAdminResponse{Admin: &admin}, nil
}

// UpdateAdminRequest contains the data for updating an admin
type UpdateAdminRequest struct {
	Name         *string `json:"name,omitempty"`
	Email        *string `json:"email,omitempty"`
	IsActive     *bool   `json:"is_active,omitempty"`
	IsSuperadmin *bool   `json:"is_superadmin,omitempty"`
}

// UpdateAdminResponse contains the updated admin
type UpdateAdminResponse struct {
	Admin *Admin `json:"admin"`
}

// UpdateAdmin updates an existing admin
//
//encore:api auth method=PUT path=/admin/admins/:id tag:admin
func UpdateAdmin(ctx context.Context, id string, req *UpdateAdminRequest) (*UpdateAdminResponse, error) {
	// Build dynamic update query
	var admin Admin
	err := db.QueryRow(ctx, `
		UPDATE admins SET
			name = COALESCE($2, name),
			email = COALESCE($3, email),
			is_active = COALESCE($4, is_active),
			is_superadmin = COALESCE($5, is_superadmin),
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, username, email, name, is_active, is_superadmin, created_at
	`, id, req.Name, req.Email, req.IsActive, req.IsSuperadmin).Scan(
		&admin.ID, &admin.Username, &admin.Email, &admin.Name,
		&admin.IsActive, &admin.IsSuperadmin, &admin.CreatedAt,
	)

	if err == sqldb.ErrNoRows {
		return nil, &errs.Error{Code: errs.NotFound, Message: "admin not found"}
	}
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to update admin"}
	}

	return &UpdateAdminResponse{Admin: &admin}, nil
}

// DeleteAdminResponse indicates deletion success
type DeleteAdminResponse struct {
	Success bool `json:"success"`
}

// DeleteAdmin deletes an admin user
//
//encore:api auth method=DELETE path=/admin/admins/:id tag:admin
func DeleteAdmin(ctx context.Context, id string) (*DeleteAdminResponse, error) {
	// Prevent deleting the last superadmin
	var superadminCount int
	err := db.QueryRow(ctx, `SELECT COUNT(*) FROM admins WHERE is_superadmin = true`).Scan(&superadminCount)
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to check superadmin count"}
	}

	// Check if this admin is a superadmin
	var isSuperadmin bool
	err = db.QueryRow(ctx, `SELECT is_superadmin FROM admins WHERE id = $1`, id).Scan(&isSuperadmin)
	if err == sqldb.ErrNoRows {
		return nil, &errs.Error{Code: errs.NotFound, Message: "admin not found"}
	}
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to get admin"}
	}

	if isSuperadmin && superadminCount <= 1 {
		return nil, &errs.Error{Code: errs.FailedPrecondition, Message: "cannot delete the last superadmin"}
	}

	// Delete the admin
	result, err := db.Exec(ctx, `DELETE FROM admins WHERE id = $1`, id)
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to delete admin"}
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, &errs.Error{Code: errs.NotFound, Message: "admin not found"}
	}

	return &DeleteAdminResponse{Success: true}, nil
}

// ChangePasswordRequest contains the password change data
type ChangePasswordRequest struct {
	NewPassword string `json:"new_password"`
}

// ChangePasswordResponse indicates success
type ChangePasswordResponse struct {
	Success bool `json:"success"`
}

// ChangeAdminPassword changes an admin's password
//
//encore:api auth method=PUT path=/admin/admins/:id/password tag:admin
func ChangeAdminPassword(ctx context.Context, id string, req *ChangePasswordRequest) (*ChangePasswordResponse, error) {
	if len(req.NewPassword) < 8 {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "password must be at least 8 characters"}
	}

	// Hash new password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to hash password"}
	}

	result, err := db.Exec(ctx, `
		UPDATE admins SET password_hash = $2, updated_at = NOW() WHERE id = $1
	`, id, string(hash))
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to update password"}
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, &errs.Error{Code: errs.NotFound, Message: "admin not found"}
	}

	return &ChangePasswordResponse{Success: true}, nil
}
