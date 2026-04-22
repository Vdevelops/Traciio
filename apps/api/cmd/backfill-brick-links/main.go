package main

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/config"
	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/pkg/logger"
	"github.com/gilabs/crm-healthcare/api/seeders"
)

func main() {
	logger.Init()

	if err := config.Load(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	log.Println("=== Backfilling brick links ===")

	// Ensure sales users have brick_id.
	if err := seeders.AssignBricksToUsers(); err != nil {
		log.Fatal("Failed to assign bricks to users:", err)
	}

	// Backfill brick_id across other tables that can infer it.
	if err := seeders.BackfillBrickLinks(); err != nil {
		log.Fatal("Failed to backfill brick links:", err)
	}

	log.Println("=== Backfill completed successfully ===")
}
