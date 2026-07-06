package seeders

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"gorm.io/gorm/clause"
)

// SeedRoleScopes seeds default data visibility scopes for each role.
// Admin -> global (sees all data)
// Sales Manager -> team (sees data from sales reps assigned to managed bricks)
// Sales -> own (sees only their own data)
// No analyst role is seeded in this project.
func SeedRoleScopes() error {
	// Fetch existing roles by code
	var adminRole, salesManagerRole, salesRole role.Role
	if err := database.DB.Where("code = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}
	if err := database.DB.Where("code = ?", "sales_manager").First(&salesManagerRole).Error; err != nil {
		return err
	}
	if err := database.DB.Where("code = ?", "sales").First(&salesRole).Error; err != nil {
		return err
	}

	// Resources that support data scoping
	resources := []string{
		"accounts",
		"contacts",
		"leads",
		"deals",
		"tasks",
		"schedules",
		"visit-reports",
		"activities",
		"monthly-targets",
		"users",
		"dashboard",
		"sales-overview",
		"route-optimization",
	}

	// Build scope assignments per role
	var scopes []role.RoleScope
	for _, resource := range resources {
		scopes = append(scopes,
			role.RoleScope{RoleID: adminRole.ID, Resource: resource, Scope: role.ScopeGlobal},
			role.RoleScope{RoleID: salesManagerRole.ID, Resource: resource, Scope: role.ScopeTeam},
			role.RoleScope{RoleID: salesRole.ID, Resource: resource, Scope: role.ScopeOwn},
		)
	}

	for _, s := range scopes {
		if err := database.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "role_id"}, {Name: "resource"}},
			DoUpdates: clause.AssignmentColumns([]string{"scope"}),
		}).Create(&s).Error; err != nil {
			return err
		}
	}

	log.Printf("Role scopes seeded successfully (%d records)", len(scopes))
	return nil
}

// AddMonthlyTargetScopes adds the missing monthly-targets and users data scope for each role.
// This is a one-time migration for databases seeded before these scopes were introduced.
// Idempotent: skips roles that already have a scope entry for the resource.
func AddMonthlyTargetScopes() error {
	var adminRole, salesManagerRole, salesRole role.Role
	if err := database.DB.Where("code = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}
	if err := database.DB.Where("code = ?", "sales_manager").First(&salesManagerRole).Error; err != nil {
		return err
	}
	if err := database.DB.Where("code = ?", "sales").First(&salesRole).Error; err != nil {
		return err
	}

	entries := []role.RoleScope{
		{RoleID: adminRole.ID, Resource: "monthly-targets", Scope: role.ScopeGlobal},
		{RoleID: salesManagerRole.ID, Resource: "monthly-targets", Scope: role.ScopeTeam},
		{RoleID: salesRole.ID, Resource: "monthly-targets", Scope: role.ScopeOwn},
		{RoleID: adminRole.ID, Resource: "users", Scope: role.ScopeGlobal},
		{RoleID: salesManagerRole.ID, Resource: "users", Scope: role.ScopeTeam},
		{RoleID: salesRole.ID, Resource: "users", Scope: role.ScopeOwn},
	}

	added := 0
	for _, entry := range entries {
		var existing role.RoleScope
		err := database.DB.
			Where("role_id = ? AND resource = ?", entry.RoleID, entry.Resource).
			First(&existing).Error
		if err == nil {
			log.Printf("Scope for role %s / resource %s already exists, skipping", entry.RoleID, entry.Resource)
			continue
		}
		if err := database.DB.Create(&entry).Error; err != nil {
			return err
		}
		added++
	}

	log.Printf("AddMonthlyTargetScopes: added %d scope records", added)
	return nil
}
