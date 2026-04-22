package seeders

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/activity_type"
)

// SeedActivityTypes seeds initial activity types
func SeedActivityTypes() error {
	// Check if activity types already exist
	var count int64
	database.DB.Model(&activity_type.ActivityType{}).Count(&count)
	if count > 0 {
		log.Println("Activity types already seeded, skipping...")
		return nil
	}

	activityTypes := []activity_type.ActivityType{
		{
			Name:        "Visit",
			Code:        "visit",
			Description: "Visit to account or contact",
			Icon:        "activity",
			BadgeColor:  "default",
			Status:      "active",
			Order:       1,
		},
		{
			Name:        "Call",
			Code:        "call",
			Description: "Phone call with account or contact",
			Icon:        "phone",
			BadgeColor:  "secondary",
			Status:      "active",
			Order:       2,
		},
		{
			Name:        "Email",
			Code:        "email",
			Description: "Email communication",
			Icon:        "mail",
			BadgeColor:  "secondary",
			Status:      "active",
			Order:       3,
		},
		{
			Name:        "Task",
			Code:        "task",
			Description: "Task or reminder",
			Icon:        "circle-check",
			BadgeColor:  "outline",
			Status:      "active",
			Order:       4,
		},
		{
			Name:        "Deal",
			Code:        "deal",
			Description: "Deal or opportunity related activity",
			Icon:        "handshake",
			BadgeColor:  "default",
			Status:      "active",
			Order:       5,
		},
	}

	// Create activity types
	for _, at := range activityTypes {
		if err := database.DB.Create(&at).Error; err != nil {
			return err
		}
	}

	log.Printf("Seeded %d activity types", len(activityTypes))
	return nil
}

