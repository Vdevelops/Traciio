package main

import (
	"log"
	"os"

	"github.com/gilabs/crm-healthcare/api/internal/config"
	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/pkg/logger"
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

	log.Println("=== Running Lead Status Migration ===")

	// Read migration file
	migrationSQL, err := os.ReadFile("internal/database/migrations/20250128_create_lead_statuses_table.sql")
	if err != nil {
		log.Fatal("Failed to read migration file:", err)
	}

	// Execute migration
	if err := database.DB.Exec(string(migrationSQL)).Error; err != nil {
		log.Fatal("Failed to execute migration:", err)
	}

	log.Println("=== Migration completed successfully ===")
	log.Println("Lead statuses table created!")
}
