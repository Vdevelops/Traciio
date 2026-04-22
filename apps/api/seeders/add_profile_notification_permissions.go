package seeders

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/permission"
	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"gorm.io/gorm/clause"
)

// AddProfileAndNotificationPermissions adds Profile and Notification permissions
func AddProfileAndNotificationPermissions() error {
	log.Println("Adding Profile and Notification permissions...")

	// Get Profile menu
	var profileMenu permission.Menu
	if err := database.DB.Where("url = ?", "/profile").First(&profileMenu).Error; err != nil {
		log.Printf("Warning: Profile menu not found: %v", err)
	}

	// Get Notification menu
	var notificationMenu permission.Menu
	if err := database.DB.Where("url = ?", "/notifications").First(&notificationMenu).Error; err != nil {
		log.Printf("Warning: Notification menu not found: %v", err)
	}

	// Define permissions
	newPermissions := []struct {
		MenuID *string
		Code   string
		Name   string
		Action string
	}{
		// Profile permissions
		{&profileMenu.ID, "profile.view", "View Profile", "VIEW"},
		{&profileMenu.ID, "profile.edit", "Edit Profile", "EDIT"},
		{&profileMenu.ID, "profile.change-password", "Change Password", "EDIT"},

		// Notification permissions
		{&notificationMenu.ID, "notifications.view", "View Notifications", "VIEW"},
		{&notificationMenu.ID, "notifications.mark-read", "Mark Notification as Read", "EDIT"},
		{&notificationMenu.ID, "notifications.delete", "Delete Notification", "DELETE"},
	}

	// Insert or Update permissions
	for _, p := range newPermissions {
		perm := permission.Permission{
			MenuID: p.MenuID,
			Code:   p.Code,
			Name:   p.Name,
			Action: p.Action,
		}

		if err := database.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "code"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "menu_id", "action"}),
		}).Create(&perm).Error; err != nil {
			log.Printf("Warning: Failed to seed permission %s: %v", p.Code, err)
		}
	}

	// Assign ALL these permissions to ALL roles (including admin, sales, etc.)
	// Because every user should be able to view/edit their own profile and notifications
	var allRoles []role.Role
	if err := database.DB.Find(&allRoles).Error; err == nil {
		for _, r := range allRoles {
			for _, p := range newPermissions {
				var perm permission.Permission
				if err := database.DB.Where("code = ?", p.Code).First(&perm).Error; err == nil {
					database.DB.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
						r.ID, perm.ID)
				}
			}
		}
		log.Printf("Assigned profile & notification permissions to %d roles", len(allRoles))
	}

	log.Println("Profile and Notification permissions added successfully")
	return nil
}
