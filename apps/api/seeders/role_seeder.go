package seeders

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"gorm.io/gorm/clause"
)

// SeedRoles seeds the canonical roles used by the application.
func SeedRoles() error {
	roles := []role.Role{
		{
			Name:         "Admin",
			Code:         "admin",
			Description:  "Administrator with full menu, action, and global data access",
			Status:       "active",
			MobileAccess: false,
			IsProtected:  true,
		},
		{
			Name:         "Sales Manager",
			Code:         "sales_manager",
			Description:  "Sales manager with full actions scoped to managed bricks",
			Status:       "active",
			MobileAccess: false,
			IsProtected:  true,
		},
		{
			Name:         "Sales",
			Code:         "sales",
			Description:  "Sales representative for field execution and daily operations",
			Status:       "active",
			MobileAccess: true,
			IsProtected:  true,
		},
		{
			Name:         "Analyst",
			Code:         "analyst",
			Description:  "Analyst with reporting and analytics access",
			Status:       "active",
			MobileAccess: false,
			IsProtected:  true,
		},
	}

	canonicalCodes := []string{"admin", "sales_manager", "sales", "analyst"}
	if err := database.DB.Where("code NOT IN ?", canonicalCodes).Delete(&role.Role{}).Error; err != nil {
		return err
	}

	for _, r := range roles {
		if err := database.DB.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "code"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"name",
				"description",
				"status",
				"mobile_access",
				"is_protected",
			}),
		}).Create(&r).Error; err != nil {
			return err
		}
		log.Printf("Seeded role: %s (%s)", r.Name, r.Code)
	}

	log.Println("Canonical roles seeded successfully")
	return nil
}
