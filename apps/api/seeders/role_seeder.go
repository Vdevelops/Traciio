package seeders

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
)

// SeedRoles seeds initial roles
func SeedRoles() error {
	// Check if roles already exist
	var count int64
	database.DB.Model(&role.Role{}).Count(&count)
	if count > 0 {
		log.Println("Roles already seeded, skipping...")
		return nil
	}

	roles := []role.Role{
		{
			Name:        "Admin",
			Code:        "admin",
			Description: "Business administrator with operational and configuration access",
			Status:      "active",
			MobileAccess: false,
			IsProtected: true, // Admin role is protected and cannot be deleted/disabled
		},
		{
			Name:        "Sales Manager",
			Code:        "sales_manager",
			Description: "Sales manager/supervisor focused on team performance monitoring",
			Status:      "active",
			MobileAccess: false,
		},
		{
			Name:        "Sales",
			Code:        "sales",
			Description: "Sales/Field staff for field execution and daily operations",
			Status:      "active",
			MobileAccess: true, // Sales role can access mobile app
		},
	}

	for _, r := range roles {
		if err := database.DB.Create(&r).Error; err != nil {
			return err
		}
		log.Printf("Created role: %s (%s)", r.Name, r.Code)
	}

	log.Println("Roles seeded successfully")
	return nil
}

