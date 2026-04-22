package seeders

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/permission"
	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"gorm.io/gorm/clause"
)

// AddLeadStatusPermissions adds new Lead Status permissions to existing setup
func AddLeadStatusPermissions() error {
	log.Println("Adding Lead Status permissions...")

	// Get Leads menu (Lead Management)
	var leadsMenu permission.Menu
	if err := database.DB.Where("name = ?", "Lead Management").First(&leadsMenu).Error; err != nil {
		log.Printf("Error finding Lead Management menu: %v", err)
		return err
	}

	// Define new Lead Status permissions
	newPermissions := []struct {
		MenuID *string
		Code   string
		Name   string
		Action string
		Menu   *permission.Menu
	}{
		{&leadsMenu.ID, "leads.status-view", "View Lead Status", "VIEW_STATUS", &leadsMenu},
		{&leadsMenu.ID, "leads.status-create", "Create Lead Status", "CREATE_STATUS", &leadsMenu},
		{&leadsMenu.ID, "leads.status-edit", "Edit Lead Status", "EDIT_STATUS", &leadsMenu},
		{&leadsMenu.ID, "leads.status-delete", "Delete Lead Status", "DELETE_STATUS", &leadsMenu},
		{&leadsMenu.ID, "leads.status-default", "Set Default Lead Status", "SET_DEFAULT_STATUS", &leadsMenu},
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



	// Assign to Sales Manager role (full access)
	var salesManagerRole role.Role
	if err := database.DB.Where("code = ?", "sales_manager").First(&salesManagerRole).Error; err == nil {
		var managerPermissions []permission.Permission
		if err := database.DB.Where("code IN (?)", []string{
			"leads.status-view",
			"leads.status-create",
			"leads.status-edit",
			"leads.status-delete",
			"leads.status-default",
		}).Find(&managerPermissions).Error; err == nil {
			for _, perm := range managerPermissions {
				// Check if already assigned
				var count int64
				database.DB.Table("role_permissions").
					Where("role_id = ? AND permission_id = ?", salesManagerRole.ID, perm.ID).
					Count(&count)

				if count > 0 {
					log.Printf("Permission %s already assigned to Sales Manager", perm.Code)
					continue
				}

				if err := database.DB.Exec(
					"INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)",
					salesManagerRole.ID, perm.ID,
				).Error; err != nil {
					log.Printf("Warning: Failed to assign permission %s to Sales Manager: %v", perm.Code, err)
				} else {
					log.Printf("Assigned permission %s to Sales Manager", perm.Code)
				}
			}
		}
	}

	// Assign VIEW_LEAD_STATUS to Sales role (read-only)
	var salesRole role.Role
	if err := database.DB.Where("code = ?", "sales").First(&salesRole).Error; err == nil {
		var salesPermissions []permission.Permission
		if err := database.DB.Where("code = ?", "leads.status-view").Find(&salesPermissions).Error; err == nil {
			for _, perm := range salesPermissions {
				// Check if already assigned
				var count int64
				database.DB.Table("role_permissions").
					Where("role_id = ? AND permission_id = ?", salesRole.ID, perm.ID).
					Count(&count)

				if count > 0 {
					log.Printf("Permission %s already assigned to Sales", perm.Code)
					continue
				}

				if err := database.DB.Exec(
					"INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)",
					salesRole.ID, perm.ID,
				).Error; err != nil {
					log.Printf("Warning: Failed to assign permission %s to Sales: %v", perm.Code, err)
				} else {
					log.Printf("Assigned permission %s to Sales", perm.Code)
				}
			}
		}
	}

	log.Println("Lead Status permissions added successfully!")
	return nil
}
