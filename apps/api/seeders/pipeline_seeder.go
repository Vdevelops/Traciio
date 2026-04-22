package seeders

import (
	"errors"
	"log"
	"strings"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
	"gorm.io/gorm"
)

// SeedPipelineStages seeds pipeline stages
// Note: "Lead" stage is NOT included because leads don't go into pipeline.
// Leads are managed in Lead Management module, and only converted leads (deals) enter the pipeline.
// Pipeline stages are for deals/opportunities only.
// This function ensures stages exist - it will restore soft-deleted stages or create new ones if missing.
func SeedPipelineStages() error {
	// Check if any active stages exist
	var activeCount int64
	database.DB.Model(&pipeline.PipelineStage{}).Where("is_active = ?", true).Count(&activeCount)
	if activeCount >= 6 {
		// At least 6 stages exist (expected: Needs Analysis, Qualification, Proposal, Negotiation, Closed Won, Closed Lost)
		log.Println("Pipeline stages already exist, skipping...")
		return nil
	}
	stages := []pipeline.PipelineStage{
		{
			Name:        "Awareness",
			Code:        "awareness",
			Order:       1,
			Color:       "#3B82F6",
			IsActive:    true,
			IsWon:       false,
			IsLost:      false,
			Probability: 0,
			Description: "Initial stage after lead conversion. Prospect is aware of the solution.",
		},
		{
			Name:        "Interest",
			Code:        "interest",
			Order:       2,
			Color:       "#6366F1",
			IsActive:    true,
			IsWon:       false,
			IsLost:      false,
			Probability: 25,
			Description: "Prospect has shown interest. Follow-up meetings and discovery held.",
		},
		{
			Name:        "Desire",
			Code:        "desire",
			Order:       3,
			Color:       "#8B5CF6",
			IsActive:    true,
			IsWon:       false,
			IsLost:      false,
			Probability: 50,
			Description: "Prospect desires the solution. Proposal and pricing discussed.",
		},
		{
			Name:        "Negotiation",
			Code:        "negotiation",
			Order:       4,
			Color:       "#F59E0B",
			IsActive:    true,
			IsWon:       false,
			IsLost:      false,
			Probability: 75,
			Description: "Negotiating terms, pricing, and contract details. Final adjustments before closing.",
		},
		{
			Name:        "Closed Won",
			Code:        "closed_won",
			Order:       5,
			Color:       "#10B981",
			IsActive:    true,
			IsWon:       true,
			IsLost:      false,
			Probability: 100,
			Description: "Deal successfully closed and won. Contract signed.",
		},
		{
			Name:        "Closed Lost",
			Code:        "closed_lost",
			Order:       6,
			Color:       "#EF4444",
			IsActive:    true,
			IsWon:       false,
			IsLost:      true,
			Probability: 0,
			Description: "Deal lost or cancelled. Opportunity did not convert.",
		},
	}

	for _, stage := range stages {
		var existing pipeline.PipelineStage
		// Check if stage exists (not soft-deleted)
		err := database.DB.Where("code = ?", stage.Code).First(&existing).Error
		if err == nil {
			// Stage already exists and is active, skip
			log.Printf("Pipeline stage %s already exists, skipping...", stage.Code)
			continue
		}

		// Check if error is "record not found" (expected) or something else
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Error checking pipeline stage %s: %v", stage.Code, err)
			return err
		}

		// Check if soft-deleted record exists
		var softDeleted pipeline.PipelineStage
		err = database.DB.Unscoped().Where("code = ? AND deleted_at IS NOT NULL", stage.Code).First(&softDeleted).Error
		if err == nil {
			// Restore soft-deleted stage
			softDeleted.DeletedAt = gorm.DeletedAt{}
			softDeleted.Name = stage.Name
			softDeleted.Order = stage.Order
			softDeleted.Color = stage.Color
			softDeleted.IsActive = stage.IsActive
			softDeleted.IsWon = stage.IsWon
			softDeleted.IsLost = stage.IsLost
			softDeleted.Probability = stage.Probability
			softDeleted.Description = stage.Description
			if err := database.DB.Unscoped().Save(&softDeleted).Error; err != nil {
				log.Printf("Error restoring pipeline stage %s: %v", stage.Code, err)
				return err
			}
			log.Printf("Restored pipeline stage: %s", stage.Name)
			continue
		}

		// Stage doesn't exist, create new one
		if err := database.DB.Create(&stage).Error; err != nil {
			// Handle duplicate key error gracefully (in case of race condition)
			if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "SQLSTATE 23505") {
				log.Printf("Pipeline stage %s already exists (duplicate key), skipping...", stage.Code)
				continue
			}
			log.Printf("Error seeding pipeline stage %s: %v", stage.Code, err)
			return err
		}
		log.Printf("Seeded pipeline stage: %s", stage.Name)
	}

	return nil
}
