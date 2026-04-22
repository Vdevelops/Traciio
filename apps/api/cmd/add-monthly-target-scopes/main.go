package main

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/config"
	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/pkg/logger"
	"github.com/gilabs/crm-healthcare/api/seeders"
)

// One-time migration: adds the monthly-targets data-scope entries for all roles.
// Run this once on databases that were seeded before monthly-targets scoping was introduced.
//
//	admin        -> global  (sees all users' targets)
//	sales_manager -> team   (sees targets for members of their group)
//	sales         -> own    (sees only their own targets)
func main() {
	logger.Init()

	if err := config.Load(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	log.Println("=== Adding monthly-targets RBAC scopes ===")

	if err := seeders.AddMonthlyTargetScopes(); err != nil {
		log.Fatal("Failed to add monthly-target scopes:", err)
	}

	log.Println("=== monthly-targets scopes added successfully ===")
}
