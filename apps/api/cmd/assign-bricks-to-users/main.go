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

	log.Println("=== Assigning Bricks to Users ===")

	// Assign bricks to sales users that don't have brick_id yet
	if err := seeders.AssignBricksToUsers(); err != nil {
		log.Fatal("Failed to assign bricks to users:", err)
	}

	log.Println("=== Bricks assigned to users successfully ===")
}

