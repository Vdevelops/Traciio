package seeders

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/product"
)

// SeedProductCategories seeds initial product categories.
func SeedProductCategories() error {
	db := database.DB

	now := time.Now()

	categories := []product.ProductCategory{
		{
			Name:        "Prescription Drug",
			Slug:        "prescription-drug",
			Description: "Obat resep untuk kebutuhan terapi pasien",
			Status:      "active",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Name:        "Over The Counter (OTC)",
			Slug:        "otc",
			Description: "Obat bebas dan bebas terbatas",
			Status:      "active",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Name:        "Medical Device",
			Slug:        "medical-device",
			Description: "Alat kesehatan dan medical device",
			Status:      "active",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Name:        "Supplement",
			Slug:        "supplement",
			Description: "Vitamin dan suplemen",
			Status:      "active",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	for _, c := range categories {
		var existing product.ProductCategory
		if err := db.Where("slug = ?", c.Slug).First(&existing).Error; err == nil {
			continue
		}

		if err := db.Create(&c).Error; err != nil {
			return err
		}
	}

	return nil
}


