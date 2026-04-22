package seeders

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/permission"
	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"gorm.io/gorm/clause"
)

// AddTaskPermissions adds comprehensive Task management permissions
func AddTaskPermissions() error {
	log.Println("Adding Task permissions...")

	// Get Task menu
	var taskMenu permission.Menu
	if err := database.DB.Where("url = ?", "/tasks").First(&taskMenu).Error; err != nil {
		log.Printf("Error finding Task menu: %v", err)
		return err
	}

	// Define Task permissions
	newPermissions := []struct {
		MenuID *string
		Code   string
		Name   string
		Action string
		Menu   *permission.Menu
	}{
		{&taskMenu.ID, "tasks.view", "View Tasks", "VIEW", &taskMenu},
		{&taskMenu.ID, "tasks.create", "Create Task", "CREATE", &taskMenu},
		{&taskMenu.ID, "tasks.edit", "Edit Task", "EDIT", &taskMenu},
		{&taskMenu.ID, "tasks.delete", "Delete Task", "DELETE", &taskMenu},
		{&taskMenu.ID, "tasks.complete", "Complete Task", "COMPLETE", &taskMenu},
		{&taskMenu.ID, "tasks.start", "Start Task", "START", &taskMenu},
		{&taskMenu.ID, "tasks.cancel", "Cancel Task", "CANCEL", &taskMenu},
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



	// Assign to Admin
	var adminRole role.Role
	if err := database.DB.Where("code = ?", "admin").First(&adminRole).Error; err == nil {
		for _, p := range newPermissions {
			var perm permission.Permission
			if err := database.DB.Where("code = ?", p.Code).First(&perm).Error; err == nil {
				database.DB.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
					adminRole.ID, perm.ID)
			}
		}
	}

	// Assign to Sales (own tasks only)
	var salesRole role.Role
	if err := database.DB.Where("code = ?", "sales").First(&salesRole).Error; err == nil {
		salesPermissions := []string{
			"tasks.view",
			"tasks.create",
			"tasks.edit",
			"tasks.complete",
			"tasks.start",
			"tasks.cancel",
		}
		for _, code := range salesPermissions {
			var perm permission.Permission
			if err := database.DB.Where("code = ?", code).First(&perm).Error; err == nil {
				database.DB.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
					salesRole.ID, perm.ID)
			}
		}
	}

	log.Println("Task permissions added successfully")
	return nil
}
