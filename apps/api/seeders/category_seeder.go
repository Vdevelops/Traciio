package seeders

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/category"
)

// SeedCategories seeds the categories table
func SeedCategories() error {
	log.Println("Seeding categories...")

	categories := []category.Category{
		{
			Name:        "Hospital",
			Code:        "HOSPITAL",
			Description: "Hospital account category",
			BadgeColor:  "default",
			Status:      "active",
		},
		{
			Name:        "Clinic",
			Code:        "CLINIC",
			Description: "Clinic account category",
			BadgeColor:  "secondary",
			Status:      "active",
		},
		{
			Name:        "Pharmacy",
			Code:        "PHARMACY",
			Description: "Pharmacy account category",
			BadgeColor:  "outline",
			Status:      "active",
		},
	}

	for _, cat := range categories {
		var existing category.Category
		if err := database.DB.Where("code = ?", cat.Code).First(&existing).Error; err != nil {
			if err.Error() == "record not found" {
				if err := database.DB.Create(&cat).Error; err != nil {
					return err
				}
				log.Printf("Created category: %s", cat.Name)
			} else {
				return err
			}
		} else {
			log.Printf("Category %s already exists, skipping", cat.Code)
		}
	}

	log.Println("Categories seeded successfully")
	return nil
}

