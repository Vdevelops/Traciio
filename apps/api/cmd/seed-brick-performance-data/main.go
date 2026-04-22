package main

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/config"
	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/pkg/logger"
	"github.com/gilabs/crm-healthcare/api/seeders"
)

// Seed a minimal dataset so Brick Performance (achievements/metrics) is non-empty:
// - bricks + sales user assignments
// - monthly targets (including brick targets)
// - deals (closed won/lost) in current/last month
// - visit reports (approved) in current/last month
// - backfill brick_id links where inferable
func main() {
	logger.Init()

	if err := config.Load(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	log.Println("=== Seeding brick performance data ===")

	if err := seeders.SeedAccounts(); err != nil {
		log.Printf("Warning: SeedAccounts failed: %v", err)
	}

	if err := seeders.SeedContacts(); err != nil {
		log.Printf("Warning: SeedContacts failed: %v", err)
	}

	if err := seeders.SeedPipelineStages(); err != nil {
		log.Printf("Warning: SeedPipelineStages failed: %v", err)
	}

	if err := seeders.SeedBricks(); err != nil {
		log.Printf("Warning: SeedBricks failed: %v", err)
	}

	if err := seeders.AssignBricksToUsers(); err != nil {
		log.Fatal("AssignBricksToUsers failed:", err)
	}

	if err := seeders.SeedMonthlyTargets(); err != nil {
		log.Printf("Warning: SeedMonthlyTargets failed: %v", err)
	}

	if err := seeders.SeedDeals(); err != nil {
		log.Fatal("SeedDeals failed:", err)
	}

	if err := seeders.SeedVisitReports(); err != nil {
		log.Fatal("SeedVisitReports failed:", err)
	}

	if err := seeders.BackfillBrickLinks(); err != nil {
		log.Fatal("BackfillBrickLinks failed:", err)
	}

	log.Println("=== Brick performance data seeding completed ===")
}
