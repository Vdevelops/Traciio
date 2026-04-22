package seeders

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/permission"
	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"gorm.io/gorm/clause"
)

// AddVisitReportAndSchedulePermissions adds missing Visit Report and Schedule permissions
func AddVisitReportAndSchedulePermissions() error {
	log.Println("Adding Visit Report and Schedule permissions...")

	// Get Visit Report menu
	var visitMenu permission.Menu
	if err := database.DB.Where("url = ?", "/visit-reports").First(&visitMenu).Error; err != nil {
		log.Printf("Error finding Visit Report menu: %v", err)
		return err
	}

	// Get Schedule menu
	var scheduleMenu permission.Menu
	if err := database.DB.Where("url = ?", "/schedules").First(&scheduleMenu).Error; err != nil {
		log.Printf("Error finding Schedule menu: %v", err)
		return err
	}

	// Define permissions
	newPermissions := []struct {
		MenuID *string
		Code   string
		Name   string
		Action string
		Menu   *permission.Menu
	}{
		// Visit Report - adding delete permission
		{&visitMenu.ID, "visit-reports.view", "View Visits", "VIEW", &visitMenu},
		{&visitMenu.ID, "visit-reports.create", "Create Visit", "CREATE", &visitMenu},
		{&visitMenu.ID, "visit-reports.edit", "Edit Visit", "EDIT", &visitMenu},
		{&visitMenu.ID, "visit-reports.delete", "Delete Visit", "DELETE", &visitMenu},
		{&visitMenu.ID, "visit-reports.approve", "Approve Visit", "APPROVE", &visitMenu},

		// Schedule - adding delete permission
		{&scheduleMenu.ID, "schedules.view", "View Schedules", "VIEW", &scheduleMenu},
		{&scheduleMenu.ID, "schedules.create", "Create Schedule", "CREATE", &scheduleMenu},
		{&scheduleMenu.ID, "schedules.edit", "Edit Schedule", "EDIT", &scheduleMenu},
		{&scheduleMenu.ID, "schedules.delete", "Delete Schedule", "DELETE", &scheduleMenu},
		{&scheduleMenu.ID, "schedules.assign", "Assign Schedule", "ASSIGN", &scheduleMenu},
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

	// Assign to Sales (selective - no delete)
	var salesRole role.Role
	if err := database.DB.Where("code = ?", "sales").First(&salesRole).Error; err == nil {
		salesPermissions := []string{
			"visit-reports.view",
			"visit-reports.create",
			"visit-reports.edit",
			"schedules.view",
			"schedules.create",
			"schedules.edit",
		}
		for _, code := range salesPermissions {
			var perm permission.Permission
			if err := database.DB.Where("code = ?", code).First(&perm).Error; err == nil {
				database.DB.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
					salesRole.ID, perm.ID)
			}
		}
	}

	log.Println("Visit Report and Schedule permissions added successfully")
	return nil
}
