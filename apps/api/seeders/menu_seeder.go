package seeders

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/permission"
)

// SeedMenus seeds initial menus based on the provided structure using UPSERT logic
func SeedMenus() error {
	basePath := ""

	// Define all menus (Root and Children)
	type menuDef struct {
		Name      string
		Icon      string
		URL       string
		ParentURL string
		Order     int
	}

	menus := []menuDef{
		// Dashboard items
		{Name: "Dashboard", Icon: "layout-dashboard", URL: basePath + "/dashboard", Order: 1},
		{Name: "Route Optimization", Icon: "map", URL: basePath + "/route-optimization", Order: 2},

		// Sales Group Items
		{Name: "Leads", Icon: "users", URL: basePath + "/leads", Order: 10},
		{Name: "Pipeline", Icon: "kanban", URL: basePath + "/pipeline", Order: 11},
		{Name: "Tasks & Schedules", Icon: "calendar-check", URL: basePath + "/tasks", Order: 12},
		{Name: "Visits", Icon: "file-text", URL: basePath + "/visit-reports", Order: 13},

		// Inventory Group Items
		{Name: "Products", Icon: "package", URL: basePath + "/products", Order: 20},

		// Customers Group Items
		{Name: "Accounts", Icon: "building", URL: basePath + "/accounts", Order: 30},

		// Analytics Group Items
		{Name: "Sales Performance", Icon: "bar-chart-3", URL: basePath + "/sales-overview", Order: 40},
		{Name: "Product Analytics", Icon: "pie-chart", URL: basePath + "/product-analytics", Order: 41},
		// KPI menu entry (analytics)
		{Name: "KPI", Icon: "leaderboard", URL: basePath + "/kpi", Order: 42},

		{Name: "Reports", Icon: "file-bar-chart", URL: basePath + "/reports", Order: 43},

		// Management Group Items
		{Name: "Users", Icon: "users-2", URL: basePath + "/master-data/users", Order: 50},
		{Name: "Groups", Icon: "users-round", URL: basePath + "/master-data/groups", Order: 51},
		{Name: "Bricks", Icon: "map-pin", URL: basePath + "/master-data/bricks", Order: 53},
		{Name: "Targets", Icon: "target", URL: basePath + "/master-data/monthly-targets", Order: 54},

		// Area Mapping Group Items
		{Name: "Area Mapping", Icon: "map", URL: basePath + "/area-mapping", Order: 55},

		// AI Group Items
		{Name: "Chatbot", Icon: "bot", URL: basePath + "/ai-chatbot", Order: 60},

		// User Items (personal)
		{Name: "Profile", Icon: "user", URL: basePath + "/profile", Order: 70},
		{Name: "Notifications", Icon: "bell", URL: basePath + "/notifications", Order: 71},
	}

	// 1. Upsert all menus by URL
	for _, m := range menus {
		var existing permission.Menu
		result := database.DB.Where("url = ?", m.URL).First(&existing)
		if result.Error == nil {
			// Update existing
			existing.Name = m.Name
			existing.Icon = m.Icon
			existing.Order = m.Order
			existing.Status = "active"
			if err := database.DB.Save(&existing).Error; err != nil {
				log.Printf("Failed to update menu %s: %v", m.Name, err)
			}
		} else {
			// Create new
			newMenu := permission.Menu{
				Name:   m.Name,
				Icon:   m.Icon,
				URL:    m.URL,
				Order:  m.Order,
				Status: "active",
			}
			if err := database.DB.Create(&newMenu).Error; err != nil {
				log.Printf("Failed to create menu %s: %v", m.Name, err)
			}
		}
	}

	log.Println("Menus standardized and UPSERTed successfully")
	return nil
}

// UpdateMenuStructure updates existing menu structure (kept for backward compatibility)
func UpdateMenuStructure() error {
	return SeedMenus()
}
