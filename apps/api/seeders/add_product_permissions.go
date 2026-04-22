package seeders

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/permission"
	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"gorm.io/gorm/clause"
)

// AddProductPermissions adds comprehensive Product permissions including categories
func AddProductPermissions() error {
	log.Println("Adding Product permissions...")

	// Get Product menu
	var productMenu permission.Menu
	if err := database.DB.Where("url = ?", "/products").First(&productMenu).Error; err != nil {
		log.Printf("Error finding Product menu: %v", err)
		return err
	}

	// Define Product permissions
	newPermissions := []struct {
		MenuID *string
		Code   string
		Name   string
		Action string
		Menu   *permission.Menu
	}{
		{&productMenu.ID, "products.view", "View Products", "VIEW", &productMenu},
		{&productMenu.ID, "products.create", "Create Product", "CREATE", &productMenu},
		{&productMenu.ID, "products.edit", "Edit Product", "EDIT", &productMenu},
		{&productMenu.ID, "products.delete", "Delete Product", "DELETE", &productMenu},
		{&productMenu.ID, "products.category-view", "View Product Categories", "VIEW", &productMenu},
		{&productMenu.ID, "products.category-create", "Create Product Category", "CREATE", &productMenu},
		{&productMenu.ID, "products.category-edit", "Edit Product Category", "EDIT", &productMenu},
		{&productMenu.ID, "products.category-delete", "Delete Product Category", "DELETE", &productMenu},
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

	log.Println("Product permissions added successfully")
	return nil
}
