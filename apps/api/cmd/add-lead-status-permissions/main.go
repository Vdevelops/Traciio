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

	log.Println("=== Adding Lead Status Permissions ===")

	// Add new Lead Status permissions
	if err := seeders.AddLeadStatusPermissions(); err != nil {
		log.Fatal("Failed to add Lead Status permissions:", err)
	}

	log.Println("=== Lead Status permissions added successfully ===")
}
