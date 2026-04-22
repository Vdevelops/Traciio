package seeders

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/group"
	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

// SeedUsers seeds initial users
func SeedUsers() error {
	// Check if users already exist
	var count int64
	database.DB.Model(&user.User{}).Count(&count)
	if count > 0 {
		log.Println("Users already seeded, skipping...")
		return nil
	}

	// Get roles
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

	// Get groups
	var salesGroup, itGroup group.Group
	var salesGroupID, itGroupID *string
	
	if err := database.DB.Where("code = ?", "SALES").First(&salesGroup).Error; err == nil {
		salesGroupID = &salesGroup.ID
	} else {
		log.Printf("Warning: Sales group not found, users will be created without group")
	}
	// opsGroup removed as it was only for Analyst role
	if err := database.DB.Where("code = ?", "IT").First(&itGroup).Error; err == nil {
		itGroupID = &itGroup.ID
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	users := []user.User{
		// Admin users
		{
			Email:      "admin@example.com",
			Password:   string(hashedPassword),
			Name:       "Admin User",
			AvatarURL:  "https://api.dicebear.com/7.x/lorelei/svg?seed=admin@example.com",
			RoleID:     adminRole.ID,
			GroupID: itGroupID, // IT Group
			Status:     "active",
		},
		// Sales Managers
		{
			Email:      "salesmanager@example.com",
			Password:   string(hashedPassword),
			Name:       "Sales Manager",
			AvatarURL:  "https://api.dicebear.com/7.x/lorelei/svg?seed=salesmanager@example.com",
			RoleID:     salesManagerRole.ID,
			GroupID: salesGroupID, // Sales Group
			Status:     "active",
		},
		{
			Email:      "salesmanager2@example.com",
			Password:   string(hashedPassword),
			Name:       "Sales Manager 2",
			AvatarURL:  "https://api.dicebear.com/7.x/lorelei/svg?seed=salesmanager2@example.com",
			RoleID:     salesManagerRole.ID,
			GroupID: salesGroupID,
			Status:     "active",
		},
		{
			Email:      "salesmanager3@example.com",
			Password:   string(hashedPassword),
			Name:       "Sales Manager 3",
			AvatarURL:  "https://api.dicebear.com/7.x/lorelei/svg?seed=salesmanager3@example.com",
			RoleID:     salesManagerRole.ID,
			GroupID: salesGroupID,
			Status:     "active",
		},
		// Sales Users (many for testing)
		{
			Email:      "sales@example.com",
			Password:   string(hashedPassword),
			Name:       "Sales User",
			AvatarURL:  "https://api.dicebear.com/7.x/lorelei/svg?seed=sales@example.com",
			RoleID:     salesRole.ID,
			GroupID: salesGroupID, // Sales Group
			Status:     "active",
		},
		{
			Email:      "sales1@example.com",
			Password:   string(hashedPassword),
			Name:       "Sales Rep 1",
			AvatarURL:  "https://api.dicebear.com/7.x/lorelei/svg?seed=sales1@example.com",
			RoleID:     salesRole.ID,
			GroupID: salesGroupID,
			Status:     "active",
		},
		{
			Email:      "sales2@example.com",
			Password:   string(hashedPassword),
			Name:       "Sales Rep 2",
			AvatarURL:  "https://api.dicebear.com/7.x/lorelei/svg?seed=sales2@example.com",
			RoleID:     salesRole.ID,
			GroupID: salesGroupID,
			Status:     "active",
		},
		{
			Email:      "sales3@example.com",
			Password:   string(hashedPassword),
			Name:       "Sales Rep 3",
			AvatarURL:  "https://api.dicebear.com/7.x/lorelei/svg?seed=sales3@example.com",
			RoleID:     salesRole.ID,
			GroupID: salesGroupID,
			Status:     "active",
		},
		{
			Email:      "sales4@example.com",
			Password:   string(hashedPassword),
			Name:       "Sales Rep 4",
			AvatarURL:  "https://api.dicebear.com/7.x/lorelei/svg?seed=sales4@example.com",
			RoleID:     salesRole.ID,
			GroupID: salesGroupID,
			Status:     "active",
		},
		{
			Email:      "sales5@example.com",
			Password:   string(hashedPassword),
			Name:       "Sales Rep 5",
			AvatarURL:  "https://api.dicebear.com/7.x/lorelei/svg?seed=sales5@example.com",
			RoleID:     salesRole.ID,
			GroupID: salesGroupID,
			Status:     "active",
		},
		{
			Email:      "sales6@example.com",
			Password:   string(hashedPassword),
			Name:       "Sales Rep 6",
			AvatarURL:  "https://api.dicebear.com/7.x/lorelei/svg?seed=sales6@example.com",
			RoleID:     salesRole.ID,
			GroupID: salesGroupID,
			Status:     "active",
		},
		{
			Email:      "sales7@example.com",
			Password:   string(hashedPassword),
			Name:       "Sales Rep 7",
			AvatarURL:  "https://api.dicebear.com/7.x/lorelei/svg?seed=sales7@example.com",
			RoleID:     salesRole.ID,
			GroupID: salesGroupID,
			Status:     "active",
		},
		{
			Email:      "sales8@example.com",
			Password:   string(hashedPassword),
			Name:       "Sales Rep 8",
			AvatarURL:  "https://api.dicebear.com/7.x/lorelei/svg?seed=sales8@example.com",
			RoleID:     salesRole.ID,
			GroupID: salesGroupID,
			Status:     "active",
		},
		{
			Email:      "sales9@example.com",
			Password:   string(hashedPassword),
			Name:       "Sales Rep 9",
			AvatarURL:  "https://api.dicebear.com/7.x/lorelei/svg?seed=sales9@example.com",
			RoleID:     salesRole.ID,
			GroupID: salesGroupID,
			Status:     "active",
		},
		{
			Email:      "sales10@example.com",
			Password:   string(hashedPassword),
			Name:       "Sales Rep 10",
			AvatarURL:  "https://api.dicebear.com/7.x/lorelei/svg?seed=sales10@example.com",
			RoleID:     salesRole.ID,
			GroupID: salesGroupID,
			Status:     "active",
		},
	}

	for _, u := range users {
		if err := database.DB.Create(&u).Error; err != nil {
			return err
		}
		groupInfo := "no group"
		if u.GroupID != nil {
			groupInfo = "group_id: " + *u.GroupID
		}
		log.Printf("Created user: %s (role_id: %s, %s)", u.Email, u.RoleID, groupInfo)
	}

	log.Println("Users seeded successfully")
	return nil
}
