package seeders

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/permission"
	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"gorm.io/gorm/clause"
)

// AddRouteOptimizationPermissions adds comprehensive Route Optimization permissions
func AddRouteOptimizationPermissions() error {
	log.Println("Adding Route Optimization permissions...")

	// Get Route Optimization menu
	var routeOptimizationMenu permission.Menu
	if err := database.DB.Where("url = ?", "/route-optimization").First(&routeOptimizationMenu).Error; err != nil {
		log.Printf("Error finding Route Optimization menu: %v", err)
		return err
	}

	// Define Route Optimization permissions
	newPermissions := []struct {
		MenuID *string
		Code   string
		Name   string
		Action string
		Menu   *permission.Menu
	}{
		{&routeOptimizationMenu.ID, "route-optimization.view", "View Route Optimization", "VIEW", &routeOptimizationMenu},
		{&routeOptimizationMenu.ID, "route-optimization.create", "Create Route", "CREATE", &routeOptimizationMenu},
		{&routeOptimizationMenu.ID, "route-optimization.delete", "Delete Route", "DELETE", &routeOptimizationMenu},
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

	// Assign to Sales.
	var salesRole role.Role
	if err := database.DB.Where("code = ?", "sales").First(&salesRole).Error; err == nil {
		salesPermissions := []string{
			"route-optimization.view",
			"route-optimization.create",
			"route-optimization.delete",
		}
		for _, code := range salesPermissions {
			var perm permission.Permission
			if err := database.DB.Where("code = ?", code).First(&perm).Error; err == nil {
				database.DB.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
					salesRole.ID, perm.ID)
			}
		}
	}

	log.Println("Route Optimization permissions added successfully")
	return nil
}
