package seeders

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/permission"
	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"gorm.io/gorm/clause"
)

// AddIndustryLeadSourcePermissions adds new Industry and Lead Source permissions to existing setup
func AddIndustryLeadSourcePermissions() error {
	log.Println("Adding Industry and Lead Source permissions...")

	// Get Leads menu (Lead Management)
	var leadsMenu permission.Menu
	if err := database.DB.Where("name = ?", "Lead Management").First(&leadsMenu).Error; err != nil {
		log.Printf("Error finding Lead Management menu: %v", err)
		return err
	}

	// Define new Industry permissions
	industryPermissions := []struct {
		MenuID *string
		Code   string
		Name   string
		Action string
		Menu   *permission.Menu
	}{
		{&leadsMenu.ID, "leads.industries-view", "View Industries", "leads.industries-view", &leadsMenu},
		{&leadsMenu.ID, "leads.industries-create", "Create Industries", "leads.industries-create", &leadsMenu},
		{&leadsMenu.ID, "leads.industries-edit", "Edit Industries", "leads.industries-edit", &leadsMenu},
		{&leadsMenu.ID, "leads.industries-delete", "Delete Industries", "leads.industries-delete", &leadsMenu},
	}

	// Define new Lead Source permissions
	leadSourcePermissions := []struct {
		MenuID *string
		Code   string
		Name   string
		Action string
		Menu   *permission.Menu
	}{
		{&leadsMenu.ID, "leads.sources-view", "View Lead Sources", "leads.sources-view", &leadsMenu},
		{&leadsMenu.ID, "leads.sources-create", "Create Lead Sources", "leads.sources-create", &leadsMenu},
		{&leadsMenu.ID, "leads.sources-edit", "Edit Lead Sources", "leads.sources-edit", &leadsMenu},
		{&leadsMenu.ID, "leads.sources-delete", "Delete Lead Sources", "leads.sources-delete", &leadsMenu},
	}

	// Insert or Update Industry permissions
	for _, p := range industryPermissions {
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

	// Insert or Update Lead Source permissions
	for _, p := range leadSourcePermissions {
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



	// Assign to Admin role (business admin)
	var adminRole role.Role
	if err := database.DB.Where("code = ?", "admin").First(&adminRole).Error; err == nil {
		var adminPermissions []permission.Permission
		if err := database.DB.Where("code IN (?)", []string{
			"leads.industries-view",
			"leads.industries-create",
			"leads.industries-edit",
			"leads.industries-delete",
			"leads.sources-view",
			"leads.sources-create",
			"leads.sources-edit",
			"leads.sources-delete",
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

	// Assign to Sales Manager role (full access)
	var salesManagerRole role.Role
	if err := database.DB.Where("code = ?", "sales_manager").First(&salesManagerRole).Error; err == nil {
		var managerPermissions []permission.Permission
		if err := database.DB.Where("code IN (?)", []string{
			"leads.industries-view",
			"leads.industries-create",
			"leads.industries-edit",
			"leads.industries-delete",
			"leads.sources-view",
			"leads.sources-create",
			"leads.sources-edit",
			"leads.sources-delete",
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

	// Assign VIEW_INDUSTRIES and VIEW_LEAD_SOURCES to Sales role (read-only)
	var salesRole role.Role
	if err := database.DB.Where("code = ?", "sales").First(&salesRole).Error; err == nil {
		var salesPermissions []permission.Permission
		if err := database.DB.Where("code IN (?)", []string{
			"leads.industries-view",
			"leads.sources-view",
		}).Find(&salesPermissions).Error; err == nil {
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

	log.Println("Industry and Lead Source permissions added successfully!")
	return nil
}

