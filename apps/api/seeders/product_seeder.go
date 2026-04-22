package seeders

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/product"
)

// SeedProducts seeds sample products for Sales CRM.
func SeedProducts() error {
	db := database.DB

	// Get some categories to attach products to.
	var categories []product.ProductCategory
	if err := db.Limit(3).Find(&categories).Error; err != nil {
		return err
	}

	if len(categories) == 0 {
		// Categories might depend on other seeders or failed to seed
		// Let's try to seed them explicitly if they don't exist
		if err := SeedProductCategories(); err != nil {
			return err
		}
		
		// Fetch again
		if err := db.Limit(3).Find(&categories).Error; err != nil {
			return err
		}
		
		// If still empty, return
		if len(categories) == 0 {
			return nil
		}
	}

	now := time.Now()

	samples := []product.Product{
		{
			Name:        "Amoxicillin 500mg Capsule",
			SKU:         "AMOX-500-CAP",
			Barcode:     "8991234567001",
			Price:       750000,
			Cost:        500000,
			CategoryID:  categories[0].ID,
			Description: "Antibiotik spektrum luas untuk infeksi bakteri.",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Name:        "Paracetamol 500mg Tablet",
			SKU:         "PARA-500-TAB",
			Barcode:     "8991234567002",
			Price:       300000,
			Cost:        150000,
			CategoryID:  categories[1%len(categories)].ID,
			Description: "Analgetik dan antipiretik untuk menurunkan demam dan nyeri.",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Name:        "Blood Pressure Monitor",
			SKU:         "BP-MON-001",
			Barcode:     "8991234567003",
			Price:       35000000,
			Cost:        25000000,
			CategoryID:  categories[2%len(categories)].ID,
			Description: "Alat pengukur tekanan darah digital untuk klinik dan rumah sakit.",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Name:        "Ibuprofen 400mg Tablet",
			SKU:         "IBU-400-TAB",
			Barcode:     "8991234567004",
			Price:       450000,
			Cost:        250000,
			CategoryID:  categories[1%len(categories)].ID,
			Description: "Obat anti inflamasi non-steroid (OAINS) untuk nyeri.",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Name:        "Vitamin C 500mg",
			SKU:         "VIT-C-500",
			Barcode:     "8991234567005",
			Price:       150000,
			Cost:        80000,
			CategoryID:  categories[3%len(categories)].ID,
			Description: "Suplemen Vitamin C untuk daya tahan tubuh.",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Name:        "Omeprazole 20mg Capsule",
			SKU:         "OMEP-20-CAP",
			Barcode:     "8991234567006",
			Price:       1200000,
			Cost:        800000,
			CategoryID:  categories[0].ID,
			Description: "Obat untuk masalah lambung dan GERD.",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Name:        "Stethoscope Classic III",
			SKU:         "STETH-003",
			Barcode:     "8991234567007",
			Price:       150000000, // Rp 1.500.000,00
			Cost:        120000000,
			CategoryID:  categories[2%len(categories)].ID,
			Description: "Stetoskop kualitas tinggi untuk diagnosa fisik.",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Name:        "Cetirizine 10mg Tablet",
			SKU:         "CET-10-TAB",
			Barcode:     "8991234567008",
			Price:       600000,
			Cost:        350000,
			CategoryID:  categories[1%len(categories)].ID,
			Description: "Antihistamin untuk mengatasi alergi.",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Name:        "Thermometer Digital",
			SKU:         "THERM-DIG-009",
			Barcode:     "8991234567009",
			Price:       7500000, // Rp 75.000,00
			Cost:        4000000,
			CategoryID:  categories[2%len(categories)].ID,
			Description: "Termometer digital akurat.",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Name:        "Azithromycin 500mg Tablet",
			SKU:         "AZITH-500-TAB",
			Barcode:     "8991234567010",
			Price:       1800000,
			Cost:        1200000,
			CategoryID:  categories[0].ID,
			Description: "Antibiotik makrolida untuk infeksi pernapasan.",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	for _, p := range samples {
		var existing product.Product
		if err := db.Where("sku = ?", p.SKU).First(&existing).Error; err == nil {
			continue
		}

		if err := db.Create(&p).Error; err != nil {
			return err
		}
	}

	return nil
}


