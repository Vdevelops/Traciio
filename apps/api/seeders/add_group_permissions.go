package seeders

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/permission"
	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"gorm.io/gorm/clause"
)

// AddGroupPermissions adds new Group permissions to existing setup
func AddGroupPermissions() error {
	log.Println("Adding Group permissions...")

	// Group menu is handled by SeedMenus

	// Get Group menu (Master Data)
	var groupMenu permission.Menu
	if err := database.DB.Where("url = ?", "/master-data/groups").First(&groupMenu).Error; err != nil {
		log.Printf("Error finding Group menu: %v", err)
		return err
	}

	// Define new Group permissions
	newPermissions := []struct {
		MenuID *string
		Code   string
		Name   string
		Action string
		Menu   *permission.Menu
	}{
		{&groupMenu.ID, "groups.view", "View Groups", "VIEW", &groupMenu},
		{&groupMenu.ID, "groups.create", "Create Group", "CREATE", &groupMenu},
		{&groupMenu.ID, "groups.edit", "Edit Group", "EDIT", &groupMenu},
		{&groupMenu.ID, "groups.delete", "Delete Group", "DELETE", &groupMenu},
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
			log.Printf("Error seeding permission %s: %v", p.Code, err)
			return err
		}
		log.Printf("Seeded permission: %s", p.Code)
	}



	// Assign to Admin role
	var adminRole role.Role
	if err := database.DB.Where("code = ?", "admin").First(&adminRole).Error; err == nil {
		var adminPermissions []permission.Permission
		if err := database.DB.Where("code IN (?)", []string{
			"groups.view",
			"groups.create",
			"groups.edit",
			"groups.delete",
		}).Find(&adminPermissions).Error; err == nil {
			for _, perm := range adminPermissions {
				// Check if already assigned
				var count int64
				database.DB.Table("role_permissions").
					Where("role_id = ? AND permission_id = ?", adminRole.ID, perm.ID).
					Count(&count)

				if count > 0 {
					log.Printf("Permission %s already assigned to Admin", perm.Code)
					continue
				}

				if err := database.DB.Exec(
					"INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)",
					adminRole.ID, perm.ID,
				).Error; err != nil {
					log.Printf("Warning: Failed to assign permission %s to Admin: %v", perm.Code, err)
				} else {
					log.Printf("Assigned permission %s to Admin", perm.Code)
				}
			}
		}
	}

	log.Println("Group permissions added successfully!")
	return nil
}

