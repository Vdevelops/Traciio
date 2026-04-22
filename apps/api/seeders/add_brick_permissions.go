package seeders

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/permission"
	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"gorm.io/gorm/clause"
)

// AddBrickPermissions adds new Brick permissions to existing setup
func AddBrickPermissions() error {
	log.Println("Adding Brick permissions...")

	// Group menu is handled by SeedMenus

	// Get Brick menu (Master Data)
	var brickMenu permission.Menu
	if err := database.DB.Where("url = ?", "/master-data/bricks").First(&brickMenu).Error; err != nil {
		log.Printf("Error finding Brick menu: %v", err)
		return err
	}

	// Define new Brick permissions
	newPermissions := []struct {
		MenuID *string
		Code   string
		Name   string
		Action string
		Menu   *permission.Menu
	}{
		{&brickMenu.ID, "bricks.view", "View Bricks", "VIEW", &brickMenu},
		{&brickMenu.ID, "bricks.create", "Create Brick", "CREATE", &brickMenu},
		{&brickMenu.ID, "bricks.edit", "Edit Brick", "EDIT", &brickMenu},
		{&brickMenu.ID, "bricks.delete", "Delete Brick", "DELETE", &brickMenu},
		{&brickMenu.ID, "bricks.distribute-targets", "Distribute Brick Targets", "EDIT", &brickMenu},
		{&brickMenu.ID, "bricks.target-distributions-view", "View Brick Target Distributions", "VIEW", &brickMenu},
		{&brickMenu.ID, "bricks.target-distributions-edit", "Edit Brick Target Distributions", "EDIT", &brickMenu},
		{&brickMenu.ID, "bricks.target-distributions-delete", "Delete Brick Target Distributions", "DELETE", &brickMenu},
		// Kepala Area / Sales Manager permissions
		{&brickMenu.ID, "bricks.dashboard-view", "View Brick Dashboard", "VIEW", &brickMenu},
		{&brickMenu.ID, "bricks.sales-view", "View Brick Sales", "VIEW", &brickMenu},
		{&brickMenu.ID, "bricks.deals-view", "View Brick Deals", "VIEW", &brickMenu},
		{&brickMenu.ID, "bricks.visit-reports-view", "View Brick Visit Reports", "VIEW", &brickMenu},
		{&brickMenu.ID, "bricks.accounts-view", "View Brick Accounts", "VIEW", &brickMenu},
		{&brickMenu.ID, "bricks.analytics-view", "View Brick Analytics", "VIEW", &brickMenu},
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
			"bricks.view",
			"bricks.create",
			"bricks.edit",
			"bricks.delete",
			"bricks.distribute-targets",
			"bricks.target-distributions-view",
			"bricks.target-distributions-edit",
			"bricks.target-distributions-delete",
			"bricks.dashboard-view",
			"bricks.sales-view",
			"bricks.deals-view",
			"bricks.visit-reports-view",
			"bricks.accounts-view",
			"bricks.analytics-view",
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

	log.Println("Brick permissions added successfully!")
	return nil
}

