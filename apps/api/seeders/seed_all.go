package seeders

import (
	"log"

	"gorm.io/gorm"
)

// SeedAll runs all seeders
func SeedAll(db *gorm.DB) error {
	// Seed in order: roles -> menus -> permissions -> users -> accounts -> contacts
	if err := SeedRoles(); err != nil {
		return err
	}

	if err := SeedMenus(); err != nil {
		return err
	}

	// Update menu structure for existing menus (migration)
	if err := UpdateMenuStructure(); err != nil {
		return err
	}

	if err := SeedPermissions(); err != nil {
		return err
	}

	// Seed role scopes (data visibility: global/team/own per resource)
	if err := SeedRoleScopes(); err != nil {
		log.Printf("Warning: Failed to seed role scopes: %v", err)
	}

	// Add Group permissions (if not already added in SeedPermissions)
	if err := AddGroupPermissions(); err != nil {
		log.Printf("Warning: Failed to add Group permissions: %v", err)
	}

	// Add Brick permissions (if not already added in SeedPermissions)
	if err := AddBrickPermissions(); err != nil {
		log.Printf("Warning: Failed to add Brick permissions: %v", err)
	}

	// Add Product permissions (including categories)
	if err := AddProductPermissions(); err != nil {
		log.Printf("Warning: Failed to add Product permissions: %v", err)
	}

	// Add Account permissions (including delete)
	if err := AddAccountPermissions(); err != nil {
		log.Printf("Warning: Failed to add Account permissions: %v", err)
	}

	// Add Task permissions (including complete, start, cancel, delete)
	if err := AddTaskPermissions(); err != nil {
		log.Printf("Warning: Failed to add Task permissions: %v", err)
	}

	// Add Route Optimization permissions (including create, delete)
	if err := AddRouteOptimizationPermissions(); err != nil {
		log.Printf("Warning: Failed to add Route Optimization permissions: %v", err)
	}

	// Add Visit Report and Schedule permissions (including delete)
	if err := AddVisitReportAndSchedulePermissions(); err != nil {
		log.Printf("Warning: Failed to add Visit Report and Schedule permissions: %v", err)
	}

	// Add Area Mapping permissions
	if err := AddAreaMappingPermissions(); err != nil {
		log.Printf("Warning: Failed to add Area Mapping permissions: %v", err)
	}

	// Add Profile and Notification permissions
	if err := AddProfileAndNotificationPermissions(); err != nil {
		log.Printf("Warning: Failed to add Profile and Notification permissions: %v", err)
	}

	// Seed groups (required for users)
	if err := SeedGroups(); err != nil {
		log.Printf("Warning: Failed to seed groups: %v", err)
	}

	if err := SeedUsers(); err != nil {
		return err
	}

	// Seed bricks early so subsequent seeders (accounts/deals/visit reports)
	// can attach to brick-assigned sales users.
	// This prevents brick-scoped analytics from appearing empty in dev.
	if err := SeedBricks(); err != nil {
		log.Printf("Warning: Failed to seed bricks: %v", err)
	}

	// Ensure sales users from SeedUsers() get assigned to bricks
	if err := AssignBricksToUsers(); err != nil {
		log.Printf("Warning: Failed to assign bricks to users: %v", err)
	}

	// Seed categories (required for accounts)
	if err := SeedCategories(); err != nil {
		return err
	}

	// Seed product categories (required for products)
	if err := SeedProductCategories(); err != nil {
		return err
	}

	// Seed activity types (required for activities)
	if err := SeedActivityTypes(); err != nil {
		return err
	}

	// Seed products (requires product categories)
	if err := SeedProducts(); err != nil {
		return err
	}

	// Seed contact roles (required for contacts)
	if err := SeedContactRoles(); err != nil {
		return err
	}

	// Seed accounts (requires users for assigned_to and categories for category_id)
	if err := SeedAccounts(); err != nil {
		return err
	}

	// Seed contacts (requires accounts for account_id and contact_roles for role_id)
	if err := SeedContacts(); err != nil {
		return err
	}

	// =====================================================
	// SEEDER UTAMA: Lead Management, Pipeline, Visit Reports
	// =====================================================

	// Seed lead statuses (required for leads)
	if err := SeedLeadStatuses(db); err != nil {
		return err
	}

	// Seed industries (required for leads)
	if err := SeedIndustries(); err != nil {
		log.Printf("Warning: Failed to seed industries: %v", err)
	}

	// Seed lead sources (required for leads)
	if err := SeedLeadSources(); err != nil {
		log.Printf("Warning: Failed to seed lead sources: %v", err)
	}

	// Add Industry and Lead Source permissions
	if err := AddIndustryLeadSourcePermissions(); err != nil {
		log.Printf("Warning: Failed to add Industry and Lead Source permissions: %v", err)
	}

	// Add Opportunity permissions
	if err := AddOpportunityPermissions(); err != nil {
		log.Printf("Warning: Failed to add Opportunity permissions: %v", err)
	}

	// Seed pipeline stages as master data before creating sample deals.
	if err := SeedPipelineStages(); err != nil {
		return err
	}

	// Seed a synchronized CRM dataset for development/demo usage:
	// - leads with qualification + product interest
	// - deals derived from those leads
	// - visit reports, activities, and tasks bound to the same lead/deal graph
	if err := SeedLeads(); err != nil {
		return err
	}

	if err := SeedDeals(); err != nil {
		return err
	}

	if err := SeedVisitReports(); err != nil {
		return err
	}

	if err := SeedActivities(); err != nil {
		return err
	}

	if err := SeedTasks(); err != nil {
		return err
	}

	// Seed monthly targets (requires users, groups, bricks)
	// Needed for achievement/target calculation
	if err := SeedMonthlyTargets(); err != nil {
		log.Printf("Warning: Failed to seed monthly targets: %v", err)
	}

	// Best-effort: backfill brick_id links for records created manually before brick assignment.
	if err := BackfillBrickLinks(); err != nil {
		log.Printf("Warning: Failed to backfill brick links: %v", err)
	}

	// Note: Schedules are now auto-created in SeedTasks, so SeedSchedules is no longer needed
	// But we keep it commented for reference
	// if err := SeedSchedules(); err != nil {
	// 	return err
	// }

	// // Seed reminders (requires tasks)
	// // COMMENT: Berlebihan, tidak perlu untuk testing sales performance
	// if err := SeedReminders(); err != nil {
	// 	return err
	// }

	// // Seed notifications (requires reminders)
	// // COMMENT: Berlebihan, tidak perlu untuk testing sales performance
	// if err := SeedNotifications(); err != nil {
	// 	return err
	// }

	// // Seed route optimization (requires users, accounts, and visit reports)
	// // COMMENT: Berlebihan, tidak critical untuk data utama
	// if err := SeedRouteOptimization(); err != nil {
	// 	return err
	// }

	// // Seed schedules (requires users, accounts, contacts, and deals)
	// // COMMENT: Berlebihan, tidak critical untuk data utama
	// if err := SeedSchedules(); err != nil {
	// 	return err
	// }

	// // Seed area mapping (requires users and visit reports)
	// // COMMENT: Berlebihan, tidak critical untuk data utama
	// if err := SeedAreaMapping(db); err != nil {
	// 	return err
	// }

	// Seed product sales (requires products, users, deals)
	if err := SeedProductSales(); err != nil {
		return err
	}

	// // Seed leaderboard (requires users, deals, visit_reports, tasks)
	// // COMMENT: Leaderboard sekarang dihitung real-time, tidak perlu seeder fake data
	// if err := SeedLeaderboard(); err != nil {
	// 	return err
	// }

	// // Seed achievements (requires users)
	// // COMMENT: Achievements sekarang dihitung real-time, tidak perlu seeder fake data
	// if err := SeedAchievements(); err != nil {
	// 	return err
	// }

	// // Seed KPI settings (no dependencies)
	// // COMMENT: Bisa di-comment jika tidak critical
	// if err := SeedKPISettings(); err != nil {
	// 	return err
	// }

	return nil
}
