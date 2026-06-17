package seeders

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/group"
	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
)

// SeedUsers seeds the canonical application users.
func SeedUsers() error {
	var adminRole, salesManagerRole, salesRole role.Role
	if err := database.DB.Where("code = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}
	if err := database.DB.Where("code = ?", "sales_manager").First(&salesManagerRole).Error; err != nil {
		return err
	}
	if err := database.DB.Where("code = ?", "sales").First(&salesRole).Error; err != nil {
		return err
	}

	var salesGroup, itGroup group.Group
	var salesGroupID, itGroupID *string
	if err := database.DB.Where("code = ?", "SALES").First(&salesGroup).Error; err == nil {
		salesGroupID = &salesGroup.ID
	} else {
		log.Printf("Warning: Sales group not found, sales users will be created without group")
	}
	if err := database.DB.Where("code = ?", "IT").First(&itGroup).Error; err == nil {
		itGroupID = &itGroup.ID
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	users := []user.User{
		{
			Email:     "admin@example.com",
			Password:  string(hashedPassword),
			Name:      "Admin",
			AvatarURL: "https://api.dicebear.com/7.x/lorelei/svg?seed=admin@example.com",
			RoleID:    adminRole.ID,
			GroupID:   itGroupID,
			Status:    "active",
		},
		{
			Email:     "salesmanager@example.com",
			Password:  string(hashedPassword),
			Name:      "Sales Manager",
			AvatarURL: "https://api.dicebear.com/7.x/lorelei/svg?seed=salesmanager@example.com",
			RoleID:    salesManagerRole.ID,
			GroupID:   salesGroupID,
			Status:    "active",
		},
		{
			Email:     "sales@example.com",
			Password:  string(hashedPassword),
			Name:      "Sales Representative",
			AvatarURL: "https://api.dicebear.com/7.x/lorelei/svg?seed=sales@example.com",
			RoleID:    salesRole.ID,
			GroupID:   salesGroupID,
			Status:    "active",
		},
	}

	canonicalEmails := []string{
		"admin@example.com",
		"salesmanager@example.com",
		"sales@example.com",
	}
	if err := database.DB.Where("email NOT IN ?", canonicalEmails).Delete(&user.User{}).Error; err != nil {
		return err
	}

	for _, u := range users {
		if err := database.DB.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "email"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"name",
				"avatar_url",
				"role_id",
				"group_id",
				"status",
			}),
		}).Create(&u).Error; err != nil {
			return err
		}
		log.Printf("Seeded user: %s", u.Email)
	}

	log.Println("Canonical users seeded successfully")
	return nil
}
