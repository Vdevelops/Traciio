package seeders

import (
	"math/rand"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
	"github.com/gilabs/crm-healthcare/api/internal/domain/product"
	"github.com/gilabs/crm-healthcare/api/internal/domain/product_analytics"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
)

// SeedProductSales seeds sample product sales for analytics
func SeedProductSales() error {
	db := database.DB

	// Check if we already have product sales
	// Check if we already have product sales
	// We'll check per month to allow seeding new months


	// Get products
	var products []product.Product
	if err := db.Limit(20).Find(&products).Error; err != nil {
		return err
	}
	if len(products) == 0 {
		return nil // No products to seed
	}

	// Get users (sales reps)
	var users []user.User
	if err := db.Where("status = ?", "active").Find(&users).Error; err != nil {
		return err
	}
	if len(users) == 0 {
		return nil // No users to seed
	}

	// Get won deals
	var deals []pipeline.Deal
	if err := db.Where("status = ?", "won").Limit(100).Find(&deals).Error; err != nil {
		return err
	}

	now := time.Now()
	var productSales []product_analytics.ProductSales

	// Create weighted product popularity (some products sell more than others)
	productWeights := make(map[string]int)
	for i, prod := range products {
		// Top 3 products get higher weights
		if i < 3 {
			productWeights[prod.ID] = 40 // 40% chance
		} else if i < 8 {
			productWeights[prod.ID] = 20 // 20% chance
		} else {
			productWeights[prod.ID] = 5 // 5% chance
		}
	}

	// Create product sales for 3 years (36 months) for multi-year analytics
	// This ensures consistent data across years for accurate comparison
	totalMonths := 36
	for month := 0; month < totalMonths; month++ {
		soldDate := now.AddDate(0, -month, 0)

		// Check if we already have sales for this month
		startOfMonth := time.Date(soldDate.Year(), soldDate.Month(), 1, 0, 0, 0, 0, soldDate.Location())
		nextMonth := startOfMonth.AddDate(0, 1, 0)

		var count int64
		if err := db.Model(&product_analytics.ProductSales{}).
			Where("sold_at >= ? AND sold_at < ?", startOfMonth, nextMonth).
			Count(&count).Error; err != nil {
			return err
		}

		if count > 0 {
			continue
		}
		
		// Vary sales count by year to simulate business growth
		// Current year: 80-150 sales/month, previous years: gradually less
		yearsAgo := month / 12
		baseSales := 80 - (yearsAgo * 15) // Reduce base by 15 per year back
		if baseSales < 30 {
			baseSales = 30
		}
		salesCount := baseSales + rand.Intn(70)
		
		for i := 0; i < salesCount; i++ {
			// Weighted random product selection
			totalWeight := 0
			for _, weight := range productWeights {
				totalWeight += weight
			}
			
			randWeight := rand.Intn(totalWeight)
			currentWeight := 0
			var selectedProduct product.Product
			
			for _, prod := range products {
				currentWeight += productWeights[prod.ID]
				if randWeight < currentWeight {
					selectedProduct = prod
					break
				}
			}
			
			// Random sales rep
			salesRep := users[rand.Intn(len(users))]
			
			// Random day of the month
			day := 1 + rand.Intn(28)
			soldAt := time.Date(soldDate.Year(), soldDate.Month(), day, 10+rand.Intn(10), rand.Intn(60), 0, 0, soldDate.Location())
			
			// Varied quantity based on product popularity
			var quantity int
			weight := productWeights[selectedProduct.ID]
			if weight >= 40 {
				quantity = 5 + rand.Intn(95) // Top products: 5-100 units
			} else if weight >= 20 {
				quantity = 2 + rand.Intn(48) // Mid-tier: 2-50 units
			} else {
				quantity = 1 + rand.Intn(20) // Low sellers: 1-20 units
			}
			
			// Use product price as base, with some variation (±10%)
			priceVariation := float64(selectedProduct.Price) * (0.9 + rand.Float64()*0.2)
			unitPrice := int64(priceVariation)
			
			// Calculate total price
			totalPrice := unitPrice * int64(quantity)
			
			// Get or create deal ID
			dealID := ""
			if len(deals) > 0 {
				dealID = deals[rand.Intn(len(deals))].ID
			}
			
			productSale := product_analytics.ProductSales{
				DealID:     dealID,
				ProductID:  selectedProduct.ID,
				Quantity:   quantity,
				UnitPrice:  unitPrice,
				TotalPrice: totalPrice,
				SoldAt:     soldAt,
				SalesRepID: salesRep.ID,
				CreatedAt:  soldAt,
				UpdatedAt:  soldAt,
			}
			
			productSales = append(productSales, productSale)
		}
	}

	// Insert in batches
	batchSize := 100
	for i := 0; i < len(productSales); i += batchSize {
		end := i + batchSize
		if end > len(productSales) {
			end = len(productSales)
		}
		batch := productSales[i:end]
		if err := db.Create(&batch).Error; err != nil {
			return err
		}
	}

	return nil
}
