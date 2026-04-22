package seeders

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/permission"
	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"gorm.io/gorm/clause"
)

// AddOpportunityPermissions adds new opportunity permissions without dropping existing data
func AddOpportunityPermissions() error {
	log.Println("Adding opportunity permissions...")

	// Get Pipeline menu
	var pipelineMenu permission.Menu
	basePath := ""
	if err := database.DB.Where("url = ?", basePath+"/pipeline").First(&pipelineMenu).Error; err != nil {
		log.Printf("Error: Pipeline menu not found: %v", err)
		return err
	}

	// Define new opportunity permissions
	newPermissions := []struct {
		code   string
		name   string
		action string
	}{
		{"pipeline.opportunity-create", "Create Opportunities", "CREATE"},
		{"pipeline.opportunity-edit", "Edit Opportunities", "EDIT"},
		{"pipeline.opportunity-delete", "Delete Opportunities", "DELETE"},
	}

	// Insert or Update permissions
	for _, p := range newPermissions {
		newPerm := permission.Permission{
			Code:        p.code,
			Name:        p.name,
			Action:      p.action,
			MenuID:      &pipelineMenu.ID,
			Description: p.name,
		}

		if err := database.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "code"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "menu_id", "action", "description"}),
		}).Create(&newPerm).Error; err != nil {
			log.Printf("Error seeding permission %s: %v", p.code, err)
			return err
		}
		log.Printf("Seeding permission: %s", p.code)
	}

	// Assign to roles
	roles := []string{"admin", "sales"}
	permissions := []string{
		"pipeline.opportunity-create",
		"pipeline.opportunity-edit",
		"pipeline.opportunity-delete",
	}

	for _, roleCode := range roles {
		var r role.Role
		if err := database.DB.Where("code = ?", roleCode).First(&r).Error; err != nil {
			log.Printf("Warning: Role %s not found", roleCode)
			continue
		}

		// For Sales role, maybe exclude delete?
		// But for now, let's keep it consistent with request or give full access. 
		// Usually sales reps can't delete opportunities.
		// Let's filter: Sales gets create/edit, Admins get all.
		rolePermissions := permissions
		if roleCode == "sales" {
			rolePermissions = []string{
				"pipeline.opportunity-create",
				"pipeline.opportunity-edit",
			}
		}

		var perms []permission.Permission
		if err := database.DB.Where("code IN (?)", rolePermissions).Find(&perms).Error; err == nil {
			for _, perm := range perms {
				// Check if already assigned
				var count int64
				database.DB.Table("role_permissions").
					Where("role_id = ? AND permission_id = ?", r.ID, perm.ID).
					Count(&count)

				if count > 0 {
					continue
				}

				if err := database.DB.Exec(
					"INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)",
					r.ID, perm.ID,
				).Error; err != nil {
					log.Printf("Warning: Failed to assign permission %s to %s: %v", perm.Code, roleCode, err)
				} else {
					log.Printf("Assigned permission %s to %s", perm.Code, roleCode)
				}
			}
		}
	}

	log.Println(" Opportunity permissions added successfully")
	return nil
}
