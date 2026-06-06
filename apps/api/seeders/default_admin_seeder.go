package seeders

import (
	"log"
	"os"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

// SeedDefaultAdmin seeds the canonical default admin user.
// This is safe to run in production as it only creates one specific user if not exists.
func SeedDefaultAdmin() error {
	// 1. Get configuration from ENV or use defaults
	email := os.Getenv("DEFAULT_ADMIN_EMAIL")
	password := os.Getenv("DEFAULT_ADMIN_PASSWORD")

	if email == "" {
		email = "admin@example.com"
		log.Printf("Using default admin email: %s", email)
	}
	if password == "" {
		password = "admin123"
		log.Println("⚠️  DEFAULT_ADMIN_PASSWORD not set. Using default password 'admin123'.")
	}

	// 2. Check if user already exists
	var count int64
	database.DB.Model(&user.User{}).Where("email = ?", email).Count(&count)
	if count > 0 {
		log.Printf("Default admin (%s) already exists, skipping...", email)
		return nil
	}

	// 3. Get Admin Role
	var adminRole role.Role
	err := database.DB.Where("code = ?", "admin").First(&adminRole).Error
	if err != nil {
		log.Printf("Error finding admin role: %v", err)
		return err
	}

	// 4. Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 5. Create User
	adminUser := user.User{
		Email:     email,
		Password:  string(hashedPassword),
		Name:      "Admin",
		AvatarURL: "https://api.dicebear.com/7.x/lorelei/svg?seed=" + email,
		RoleID:    adminRole.ID,
		Status:    "active",
	}

	if err := database.DB.Create(&adminUser).Error; err != nil {
		return err
	}

	log.Printf("✅ Default admin created successfully: %s", email)
	return nil
}
