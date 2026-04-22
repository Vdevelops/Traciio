package main

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/config"
	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/pkg/logger"
	"github.com/gilabs/crm-healthcare/api/seeders"
)

func main() {
	// Initialize logger
	logger.Init()

	// Load configuration
	if err := config.Load(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	log.Println("=== Adding Brick Menu and Permissions ===")

	// Update all menus first
	if err := seeders.SeedMenus(); err != nil {
		log.Printf("Warning: Failed to seed menus: %v", err)
	}

	// Add new Brick permissions
	if err := seeders.AddBrickPermissions(); err != nil {
		log.Fatal("Failed to add Brick permissions:", err)
	}

	log.Println("=== Brick menu and permissions added successfully ===")
}
