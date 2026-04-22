package seeders

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/permission"
	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"gorm.io/gorm/clause"
)

// AddAccountPermissions adds comprehensive Account management permissions
func AddAccountPermissions() error {
	log.Println("Adding Account permissions...")

	// Get Account menu
	var accountMenu permission.Menu
	if err := database.DB.Where("url = ?", "/accounts").First(&accountMenu).Error; err != nil {
		log.Printf("Error finding Account menu: %v", err)
		return err
	}

	// Define Account permissions
	newPermissions := []struct {
		MenuID *string
		Code   string
		Name   string
		Action string
		Menu   *permission.Menu
	}{
		{&accountMenu.ID, "accounts.view", "View Accounts", "VIEW", &accountMenu},
		{&accountMenu.ID, "accounts.create", "Create Account", "CREATE", &accountMenu},
		{&accountMenu.ID, "accounts.edit", "Edit Account", "EDIT", &accountMenu},
		{&accountMenu.ID, "accounts.delete", "Delete Account", "DELETE", &accountMenu},

		// Tabs & Category/Role management
		{&accountMenu.ID, "accounts.category", "Access Account Categories Tab", "VIEW", &accountMenu},
		{&accountMenu.ID, "accounts.category-view", "View Account Categories", "VIEW", &accountMenu},
		{&accountMenu.ID, "accounts.category-create", "Create Account Category", "CREATE", &accountMenu},
		{&accountMenu.ID, "accounts.category-edit", "Edit Account Category", "EDIT", &accountMenu},
		{&accountMenu.ID, "accounts.category-delete", "Delete Account Category", "DELETE", &accountMenu},

		{&accountMenu.ID, "accounts.role", "Access Contact Roles Tab", "VIEW", &accountMenu},
		{&accountMenu.ID, "accounts.role-view", "View Contact Roles", "VIEW", &accountMenu},
		{&accountMenu.ID, "accounts.role-create", "Create Contact Role", "CREATE", &accountMenu},
		{&accountMenu.ID, "accounts.role-edit", "Edit Contact Role", "EDIT", &accountMenu},
		{&accountMenu.ID, "accounts.role-delete", "Delete Contact Role", "DELETE", &accountMenu},
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

	// Assign to Sales (view and edit only)
	var salesRole role.Role
	if err := database.DB.Where("code = ?", "sales").First(&salesRole).Error; err == nil {
		salesPermissions := []string{
			"accounts.view",
			"accounts.create",
			"accounts.edit",
		}
		for _, code := range salesPermissions {
			var perm permission.Permission
			if err := database.DB.Where("code = ?", code).First(&perm).Error; err == nil {
				database.DB.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
					salesRole.ID, perm.ID)
			}
		}
	}

	log.Println("Account permissions added successfully")
	return nil
}
