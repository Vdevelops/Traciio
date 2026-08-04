package seeders

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/permission"
	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"gorm.io/gorm/clause"
)

// SeedPermissions seeds initial permissions based on the provided structure
func SeedPermissions() error {
	// Base path must match menu_seeder
	basePath := ""

	// Get menus by their standardized URLs
	var dashboardMenu, usersMenu, groupsMenu, targetsMenu permission.Menu
	var leadsMenu, pipelineMenu, tasksMenu, visitReportsMenu, routeOptimizationMenu, scheduleMenu permission.Menu
	var productsMenu, accountsMenu, reportsMenu permission.Menu
	var aiChatbotMenu, salesOverviewMenu, productAnalyticsMenu permission.Menu
	var areaMappingMenu, profileMenu, notificationsMenu permission.Menu

	database.DB.Where("url = ?", basePath+"/dashboard").First(&dashboardMenu)
	database.DB.Where("url = ?", basePath+"/master-data/users").First(&usersMenu)
	database.DB.Where("url = ?", basePath+"/master-data/groups").First(&groupsMenu)
	database.DB.Where("url = ?", basePath+"/master-data/monthly-targets").First(&targetsMenu)
	database.DB.Where("url = ?", basePath+"/leads").First(&leadsMenu)
	database.DB.Where("url = ?", basePath+"/pipeline").First(&pipelineMenu)
	database.DB.Where("url = ?", basePath+"/tasks").First(&tasksMenu)
	database.DB.Where("url = ?", basePath+"/visit-reports").First(&visitReportsMenu)
	database.DB.Where("url = ?", basePath+"/route-optimization").First(&routeOptimizationMenu)
	database.DB.Where("url = ?", basePath+"/schedules").First(&scheduleMenu)
	database.DB.Where("url = ?", basePath+"/products").First(&productsMenu)
	database.DB.Where("url = ?", basePath+"/accounts").First(&accountsMenu)
	database.DB.Where("url = ?", basePath+"/reports").First(&reportsMenu)
	database.DB.Where("url = ?", basePath+"/ai-chatbot").First(&aiChatbotMenu)

	database.DB.Where("url = ?", basePath+"/sales-overview").First(&salesOverviewMenu)
	database.DB.Where("url = ?", basePath+"/product-analytics").First(&productAnalyticsMenu)
	database.DB.Where("url = ?", basePath+"/area-mapping").First(&areaMappingMenu)
	database.DB.Where("url = ?", basePath+"/profile").First(&profileMenu)
	database.DB.Where("url = ?", basePath+"/notifications").First(&notificationsMenu)

	// Define actions for each menu
	actions := []struct {
		menuID string
		code   string
		name   string
		action string
	}{
		// Dashboard
		{dashboardMenu.ID, "dashboard.view", "View Dashboard", "VIEW"},

		// Management - Users
		{usersMenu.ID, "users.view", "View Users", "VIEW"},
		{usersMenu.ID, "users.create", "Create Users", "CREATE"},
		{usersMenu.ID, "users.edit", "Edit Users", "EDIT"},
		{usersMenu.ID, "users.delete", "Delete Users", "DELETE"},
		{usersMenu.ID, "users.roles", "View Roles", "VIEW"},
		{usersMenu.ID, "users.permissions", "View Permissions", "VIEW"},

		// Management - Groups
		{groupsMenu.ID, "groups.view", "View Groups", "VIEW"},
		{groupsMenu.ID, "groups.create", "Create Group", "CREATE"},
		{groupsMenu.ID, "groups.edit", "Edit Group", "EDIT"},
		{groupsMenu.ID, "groups.delete", "Delete Group", "DELETE"},

		// Management - Targets
		{targetsMenu.ID, "monthly-targets.view", "View Targets", "VIEW"},
		{targetsMenu.ID, "monthly-targets.create", "Create Target", "CREATE"},
		{targetsMenu.ID, "monthly-targets.edit", "Edit Target", "EDIT"},
		{targetsMenu.ID, "monthly-targets.delete", "Delete Target", "DELETE"},

		// Sales - Leads
		{leadsMenu.ID, "leads.view", "View Leads", "VIEW"},
		{leadsMenu.ID, "leads.create", "Create Leads", "CREATE"},
		{leadsMenu.ID, "leads.edit", "Edit Leads", "EDIT"},
		{leadsMenu.ID, "leads.delete", "Delete Leads", "DELETE"},
		{leadsMenu.ID, "leads.convert", "Convert Leads", "CONVERT"},
		{leadsMenu.ID, "leads.status-view", "View Lead Status", "VIEW_STATUS"},
		{leadsMenu.ID, "leads.industries-view", "View Industries", "VIEW"},
		{leadsMenu.ID, "leads.industries-create", "Create Industry", "CREATE"},
		{leadsMenu.ID, "leads.industries-edit", "Edit Industry", "EDIT"},
		{leadsMenu.ID, "leads.industries-delete", "Delete Industry", "DELETE"},
		{leadsMenu.ID, "leads.sources-view", "View Lead Sources", "VIEW"},
		{leadsMenu.ID, "leads.sources-create", "Create Lead Source", "CREATE"},
		{leadsMenu.ID, "leads.sources-edit", "Edit Lead Source", "EDIT"},
		{leadsMenu.ID, "leads.sources-delete", "Delete Lead Source", "DELETE"},

		// Sales - Pipeline
		{pipelineMenu.ID, "pipeline.view", "View Pipeline", "VIEW"},
		{pipelineMenu.ID, "pipeline.create", "Create Deals", "CREATE"},
		{pipelineMenu.ID, "pipeline.edit", "Edit Deals", "EDIT"},
		{pipelineMenu.ID, "pipeline.delete", "Delete Deals", "DELETE"},
		{pipelineMenu.ID, "pipeline.move", "Move Deals", "MOVE"},
		{pipelineMenu.ID, "pipeline.update_stage", "Update Deal Stage (Validated)", "UPDATE_STAGE"},
		{pipelineMenu.ID, "pipeline.convert_quotation", "Convert Deal to Quotation", "CONVERT_QUOTATION"},
		{pipelineMenu.ID, "pipeline.convert_sales_order", "Convert Deal to Sales Order", "CONVERT_SALES_ORDER"},
		{pipelineMenu.ID, "pipeline.stages-view", "View Pipeline Stages", "VIEW"},
		{pipelineMenu.ID, "pipeline.stages-create", "Create Pipeline Stages", "CREATE"},
		{pipelineMenu.ID, "pipeline.stages-edit", "Edit Pipeline Stages", "EDIT"},
		{pipelineMenu.ID, "pipeline.stages-delete", "Delete Pipeline Stages", "DELETE"},
		{pipelineMenu.ID, "pipeline.stages-order", "Reorder Pipeline Stages", "EDIT"},

		// Sales - Visits
		{visitReportsMenu.ID, "visit-reports.view", "View Visits", "VIEW"},
		{visitReportsMenu.ID, "visit-reports.create", "Create Visit", "CREATE"},
		{visitReportsMenu.ID, "visit-reports.edit", "Edit Visit", "EDIT"},
		{visitReportsMenu.ID, "visit-reports.delete", "Delete Visit", "DELETE"},
		{visitReportsMenu.ID, "visit-reports.approve", "Approve Visit", "APPROVE"},
		{visitReportsMenu.ID, "visit-reports.activity-type", "Manage Activity Types", "MANAGE"},

		// Sales - Schedules
		{scheduleMenu.ID, "schedules.view", "View Schedules", "VIEW"},
		{scheduleMenu.ID, "schedules.create", "Create Schedule", "CREATE"},
		{scheduleMenu.ID, "schedules.edit", "Edit Schedule", "EDIT"},
		{scheduleMenu.ID, "schedules.delete", "Delete Schedule", "DELETE"},
		{scheduleMenu.ID, "schedules.assign", "Assign Schedule", "ASSIGN"},

		// Sales - Tasks
		{tasksMenu.ID, "tasks.view", "View Tasks", "VIEW"},
		{tasksMenu.ID, "tasks.create", "Create Task", "CREATE"},
		{tasksMenu.ID, "tasks.edit", "Edit Task", "EDIT"},
		{tasksMenu.ID, "tasks.delete", "Delete Task", "DELETE"},
		{tasksMenu.ID, "tasks.complete", "Complete Task", "COMPLETE"},
		{tasksMenu.ID, "tasks.start", "Start Task", "START"},
		{tasksMenu.ID, "tasks.cancel", "Cancel Task", "CANCEL"},
		{tasksMenu.ID, "tasks.create_lead", "Create & Link Lead from Task", "CREATE_LEAD"},

		// Main - Route Optimization
		{routeOptimizationMenu.ID, "route-optimization.view", "View Route Optimization", "VIEW"},
		{routeOptimizationMenu.ID, "route-optimization.create", "Create Route", "CREATE"},
		{routeOptimizationMenu.ID, "route-optimization.delete", "Delete Route", "DELETE"},

		// Inventory - Products
		{productsMenu.ID, "products.view", "View Products", "VIEW"},
		{productsMenu.ID, "products.create", "Create Product", "CREATE"},
		{productsMenu.ID, "products.edit", "Edit Product", "EDIT"},
		{productsMenu.ID, "products.delete", "Delete Product", "DELETE"},
		{productsMenu.ID, "products.category-view", "View Product Categories", "VIEW"},
		{productsMenu.ID, "products.category-create", "Create Product Category", "CREATE"},
		{productsMenu.ID, "products.category-edit", "Edit Product Category", "EDIT"},
		{productsMenu.ID, "products.category-delete", "Delete Product Category", "DELETE"},

		// Customers - Accounts
		{accountsMenu.ID, "accounts.view", "View Accounts", "VIEW"},
		{accountsMenu.ID, "accounts.create", "Create Account", "CREATE"},
		{accountsMenu.ID, "accounts.edit", "Edit Account", "EDIT"},
		{accountsMenu.ID, "accounts.delete", "Delete Account", "DELETE"},
		{accountsMenu.ID, "accounts.category", "Manage Account Categories", "MANAGE"},
		{accountsMenu.ID, "accounts.role", "Manage Contact Roles", "MANAGE"},

		// Analytics
		{reportsMenu.ID, "reports.view", "View Reports", "VIEW"},
		{reportsMenu.ID, "reports.generate", "Generate Report", "CREATE"},

		{salesOverviewMenu.ID, "sales-overview.view", "View Sales Performance", "VIEW"},
		{productAnalyticsMenu.ID, "product-analytics.view", "View Product Analytics", "VIEW"},
		// KPI permissions (view only)
		{"", "kpi.view", "View KPI", "VIEW"},

		// AI
		{aiChatbotMenu.ID, "ai-chatbot.view", "View Chatbot", "VIEW"},
		{"", "ai-settings.view", "View AI Settings", "VIEW"},
		{"", "ai-settings.edit", "Edit AI Settings", "EDIT"},

		// Area Mapping
		{areaMappingMenu.ID, "area-mapping.view", "View Area Mapping", "VIEW"},
		{areaMappingMenu.ID, "area-mapping.territories-view", "View Territories", "VIEW"},
		{areaMappingMenu.ID, "area-mapping.territories-create", "Create Territory", "CREATE"},
		{areaMappingMenu.ID, "area-mapping.territories-edit", "Edit Territory", "EDIT"},
		{areaMappingMenu.ID, "area-mapping.territories-delete", "Delete Territory", "DELETE"},
		{areaMappingMenu.ID, "area-mapping.captures-view", "View Area Captures", "VIEW"},
		{areaMappingMenu.ID, "area-mapping.captures-create", "Create Area Capture", "CREATE"},
		{areaMappingMenu.ID, "area-mapping.coverage-view", "View Coverage Analysis", "VIEW"},
		{areaMappingMenu.ID, "area-mapping.heatmap-view", "View Heatmap", "VIEW"},

		// Profile & User Settings
		{profileMenu.ID, "profile.view", "View Profile", "VIEW"},
		{profileMenu.ID, "profile.edit", "Edit Profile", "EDIT"},
		{profileMenu.ID, "profile.change-password", "Change Password", "EDIT"},

		// Notifications
		{notificationsMenu.ID, "notifications.view", "View Notifications", "VIEW"},
		{notificationsMenu.ID, "notifications.mark-read", "Mark Notification as Read", "EDIT"},
		{notificationsMenu.ID, "notifications.delete", "Delete Notification", "DELETE"},
	}

	// Create permissions using UPSERT
	for _, act := range actions {
		perm := permission.Permission{
			Name:   act.name,
			Code:   act.code,
			Action: act.action,
		}
		if act.menuID != "" {
			perm.MenuID = &act.menuID
		}

		if err := database.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "code"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "menu_id", "action"}),
		}).Create(&perm).Error; err != nil {
			log.Printf("Warning: Failed to seed permission %s: %v", perm.Code, err)
		}
	}

	// Re-assign permissions to roles...

	// Admin Role
	var adminRole role.Role
	if err := database.DB.Where("code = ?", "admin").First(&adminRole).Error; err == nil {
		SyncAdminPermissions()
	}

	// Sales Role
	var salesRole role.Role
	if err := database.DB.Where("code = ?", "sales").First(&salesRole).Error; err == nil {
		// Assign detailed permissions
		permissions := []string{
			"dashboard.view",
			"monthly-targets.view",
			// Accounts
			"accounts.view", "accounts.create", "accounts.edit",
			// Leads
			"leads.view", "leads.create", "leads.edit", "leads.convert",
			// Pipeline
			"pipeline.view", "pipeline.create", "pipeline.edit", "pipeline.move", "pipeline.update_stage",
			// Tasks
			"tasks.view", "tasks.edit", "tasks.complete", "tasks.create_lead",
			// Visits
			"visit-reports.view", "visit-reports.create", "visit-reports.edit",
			// Route Optimization
			"route-optimization.view", "route-optimization.create", "route-optimization.delete",
			// Schedules
			"schedules.view", "schedules.create", "schedules.edit",
			// Products (View Only)
			"products.view",
			// Reports
			"reports.view", "reports.generate",
			// Analytics
			"sales-overview.view", "product-analytics.view",
			// AI
			"ai-chatbot.view",
		}

		database.DB.Exec("INSERT INTO role_permissions (role_id, permission_id) SELECT ?, id FROM permissions WHERE code IN (?) ON CONFLICT DO NOTHING",
			salesRole.ID, permissions)
	}

	// Sales Manager Role: limited management access; users can view/create,
	// while roles and permissions are view-only.
	var salesManagerRole role.Role
	if err := database.DB.Where("code = ?", "sales_manager").First(&salesManagerRole).Error; err == nil {
		salesManagerPermissions := []string{
			"dashboard.view",
			"monthly-targets.view",
			"monthly-targets.edit",
			"users.view",
			"users.create",
			"users.roles",
			"users.permissions",
			"leads.view",
			"leads.status-view",
			"leads.industries-view",
			"leads.sources-view",
			"pipeline.view",
			"tasks.view",
			"visit-reports.view",
			"schedules.view",
			"route-optimization.view",
			"products.view",
			"accounts.view",
			"reports.view",
			"sales-overview.view",
			"product-analytics.view",
			"ai-chatbot.view",
			"area-mapping.view",
			"area-mapping.territories-view",
			"area-mapping.captures-view",
			"area-mapping.coverage-view",
			"area-mapping.heatmap-view",
			"profile.view",
			"notifications.view",
		}

		database.DB.Exec("INSERT INTO role_permissions (role_id, permission_id) SELECT ?, id FROM permissions WHERE code IN (?) ON CONFLICT DO NOTHING",
			salesManagerRole.ID, salesManagerPermissions)
	}

	log.Println("Permissions seeded and standardized successfully")
	return nil
}

// SyncAdminPermissions syncs all existing permissions to admin role
func SyncAdminPermissions() error {
	var adminRole role.Role
	if err := database.DB.Where("code = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}

	var allPermissions []permission.Permission
	if err := database.DB.Find(&allPermissions).Error; err != nil {
		return err
	}

	for _, perm := range allPermissions {
		database.DB.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
			adminRole.ID, perm.ID)
	}

	log.Printf("Synced %d permissions to admin role", len(allPermissions))
	return nil
}
