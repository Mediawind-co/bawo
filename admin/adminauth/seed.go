package adminauth

import (
	"context"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// SeedSuperadmin ensures the superadmin exists with correct password
// Called on service initialization
func seedSuperadmin(ctx context.Context) error {
	// Default superadmin credentials
	username := "superadmin"
	email := "admin@bawo.app"
	password := "BawoAdmin2024!"
	name := "Super Admin"

	// Generate password hash
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Upsert superadmin
	_, err = db.Exec(ctx, `
		INSERT INTO admins (username, email, password_hash, name, is_superadmin)
		VALUES ($1, $2, $3, $4, true)
		ON CONFLICT (username) DO UPDATE SET
			password_hash = $3,
			updated_at = NOW()
	`, username, email, string(hash), name)

	if err != nil {
		log.Printf("Failed to seed superadmin: %v", err)
		return err
	}

	log.Println("Superadmin seeded successfully")
	return nil
}

//encore:service
type Service struct{}

func initService() (*Service, error) {
	ctx := context.Background()
	if err := seedSuperadmin(ctx); err != nil {
		log.Printf("Warning: failed to seed superadmin: %v", err)
	}
	return &Service{}, nil
}
