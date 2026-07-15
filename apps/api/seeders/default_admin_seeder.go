package seeders

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedDefaultAdmin seeds the canonical default admin user.
// This is safe to run in production as it only creates one specific user if not exists.
func SeedDefaultAdmin() error {
	// 1. Get configuration from ENV or use defaults
	email := os.Getenv("DEFAULT_ADMIN_EMAIL")
	password := os.Getenv("DEFAULT_ADMIN_PASSWORD")
	resetPassword := strings.EqualFold(os.Getenv("DEFAULT_ADMIN_RESET_PASSWORD"), "true") || os.Getenv("DEFAULT_ADMIN_RESET_PASSWORD") == "1"

	if email == "" {
		email = "admin@example.com"
		log.Printf("Using default admin email: %s", email)
	}
	if password == "" {
		password = "admin123"
		log.Println("⚠️  DEFAULT_ADMIN_PASSWORD not set. Using default password 'admin123'.")
	}

	// 2. Get Admin Role
	var adminRole role.Role
	err := database.DB.Where("code = ?", "admin").First(&adminRole).Error
	if err != nil {
		log.Printf("Error finding admin role: %v", err)
		return err
	}

	// 3. Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 4. Check if user already exists. Existing admin credentials are only reset
	// when explicitly requested, so production startup cannot silently rotate passwords.
	var existingAdmin user.User
	err = database.DB.Where("email = ?", email).First(&existingAdmin).Error
	if err == nil {
		if !resetPassword {
			log.Printf("Default admin (%s) already exists, skipping...", email)
			return nil
		}

		if err := database.DB.Model(&existingAdmin).Updates(map[string]interface{}{
			"password": string(hashedPassword),
			"role_id":  adminRole.ID,
			"status":   "active",
		}).Error; err != nil {
			return err
		}

		log.Printf("✅ Default admin password reset successfully: %s", email)
		return nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
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
