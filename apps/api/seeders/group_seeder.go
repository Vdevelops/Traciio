package seeders

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/group"
)

// SeedGroups seeds initial groups
func SeedGroups() error {
	// Check if groups already exist
	var count int64
	database.DB.Model(&group.Group{}).Count(&count)
	if count > 0 {
		log.Println("Groups already seeded, skipping...")
		return nil
	}

	groups := []group.Group{
		{
			Name:        "Sales Group",
			Code:        "SALES",
			Description: "Sales and marketing group",
			Status:      "active",
		},
		{
			Name:        "Marketing Group",
			Code:        "MARKETING",
			Description: "Marketing and promotion group",
			Status:      "active",
		},
		{
			Name:        "Operations Group",
			Code:        "OPS",
			Description: "Operations and logistics group",
			Status:      "active",
		},
		{
			Name:        "IT Group",
			Code:        "IT",
			Description: "Information Technology group",
			Status:      "active",
		},
	}

	for _, g := range groups {
		if err := database.DB.Create(&g).Error; err != nil {
			return err
		}
		log.Printf("Created group: %s (%s)", g.Name, g.Code)
	}

	log.Println("Groups seeded successfully")
	return nil
}

