package seeders

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"gorm.io/gorm/clause"
)

// SeedRoles seeds the canonical roles used by the application.
func SeedRoles() error {
	roles := []role.Role{
		{
			Name:        "Admin",
			Code:        "admin",
			Description: "Administrator with full menu, action, and global data access",
			Status:      "active",
			IsProtected: true,
		},
		{
			Name:        "Sales Manager",
			Code:        "sales_manager",
			Description: "Sales manager with full actions scoped to managed bricks",
			Status:      "active",
			IsProtected: true,
		},
		{
			Name:        "Sales",
			Code:        "sales",
			Description: "Sales representative for field execution and daily operations",
			Status:      "active",
			IsProtected: true,
		},
	}

	canonicalCodes := []string{"admin", "sales_manager", "sales"}
	if err := cleanupLegacyAnalystRole(); err != nil {
		return err
	}
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

func cleanupLegacyAnalystRole() error {
	var analystRole role.Role
	if err := database.DB.Where("code = ?", "analyst").First(&analystRole).Error; err != nil {
		return nil
	}

	if err := database.DB.Where("role_id = ?", analystRole.ID).Delete(&user.User{}).Error; err != nil {
		return err
	}
	if err := database.DB.Exec("DELETE FROM role_permissions WHERE role_id = ?", analystRole.ID).Error; err != nil {
		return err
	}
	if err := database.DB.Exec("DELETE FROM role_scopes WHERE role_id = ?", analystRole.ID).Error; err != nil {
		return err
	}
	if err := database.DB.Delete(&analystRole).Error; err != nil {
		return err
	}

	log.Println("Removed legacy analyst role and users")
	return nil
}
