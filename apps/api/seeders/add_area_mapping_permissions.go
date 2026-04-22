package seeders

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/permission"
	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"gorm.io/gorm/clause"
)

// AddAreaMappingPermissions adds Area Mapping module permissions
func AddAreaMappingPermissions() error {
	log.Println("Adding Area Mapping permissions...")

	// Get Area Mapping menu
	var areaMenu permission.Menu
	if err := database.DB.Where("url = ?", "/area-mapping").First(&areaMenu).Error; err != nil {
		log.Printf("Error finding Area Mapping menu: %v", err)
		return err
	}

	// Define Area Mapping permissions
	newPermissions := []struct {
		MenuID *string
		Code   string
		Name   string
		Action string
		Menu   *permission.Menu
	}{
		{&areaMenu.ID, "area-mapping.view", "View Area Mapping", "VIEW", &areaMenu},
		{&areaMenu.ID, "area-mapping.territories-view", "View Territories", "VIEW", &areaMenu},
		{&areaMenu.ID, "area-mapping.territories-create", "Create Territory", "CREATE", &areaMenu},
		{&areaMenu.ID, "area-mapping.territories-edit", "Edit Territory", "EDIT", &areaMenu},
		{&areaMenu.ID, "area-mapping.territories-delete", "Delete Territory", "DELETE", &areaMenu},
		{&areaMenu.ID, "area-mapping.captures-view", "View Area Captures", "VIEW", &areaMenu},
		{&areaMenu.ID, "area-mapping.captures-create", "Create Area Capture", "CREATE", &areaMenu},
		{&areaMenu.ID, "area-mapping.coverage-view", "View Coverage Analysis", "VIEW", &areaMenu},
		{&areaMenu.ID, "area-mapping.heatmap-view", "View Heatmap", "VIEW", &areaMenu},
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

	// Assign to Sales Manager (selective permissions)
	var salesManagerRole role.Role
	if err := database.DB.Where("code = ?", "sales_manager").First(&salesManagerRole).Error; err == nil {
		viewPermissions := []string{
			"area-mapping.view",
			"area-mapping.territories-view",
			"area-mapping.captures-view",
			"area-mapping.coverage-view",
			"area-mapping.heatmap-view",
		}
		for _, code := range viewPermissions {
			var perm permission.Permission
			if err := database.DB.Where("code = ?", code).First(&perm).Error; err == nil {
				database.DB.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
					salesManagerRole.ID, perm.ID)
			}
		}
	}

	log.Println("Area Mapping permissions added successfully")
	return nil
}
