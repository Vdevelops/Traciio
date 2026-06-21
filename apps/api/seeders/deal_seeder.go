package seeders

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/account"
	brickdomain "github.com/gilabs/crm-healthcare/api/internal/domain/brick"
	"github.com/gilabs/crm-healthcare/api/internal/domain/contact"
	"github.com/gilabs/crm-healthcare/api/internal/domain/customer_purchase"
	"github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
	"github.com/gilabs/crm-healthcare/api/internal/domain/product"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"gorm.io/datatypes"
)

// SeedDeals seeds initial deals data for pipeline and dashboard widgets.
// This seeder creates sample deals with proper relationships for testing sales performance.
func SeedDeals() error {
	// Check if deals already exist
	// We'll check for current month existence later to allow partial seeding

	// Check if deals exist for the current month
	nowCheck := time.Now()
	currentMonthStartCheck := time.Date(nowCheck.Year(), nowCheck.Month(), 1, 0, 0, 0, 0, nowCheck.Location())

	var count int64
	database.DB.Model(&pipeline.Deal{}).Where("created_at >= ?", currentMonthStartCheck).Count(&count)
	if count > 0 {
		log.Println("Deals already exist for current month, skipping...")
		return nil
	}

	// Get users (sales reps) - prioritize sales users that are assigned to bricks
	var users []user.User
	// First, try to get sales users that are assigned to bricks
	if err := database.DB.Where("users.brick_id IS NOT NULL AND users.status = ?", "active").
		Joins("JOIN roles ON users.role_id = roles.id").
		Where("roles.code = ?", "sales").
		Find(&users).Error; err != nil {
		log.Printf("Warning: Failed to get sales users with bricks: %v", err)
	}

	// If no sales users with bricks found, get all active sales users
	if len(users) == 0 {
		if err := database.DB.Where("users.status = ?", "active").
			Joins("JOIN roles ON users.role_id = roles.id").
			Where("roles.code IN (?)", []string{"sales", "sales_manager"}).
			Find(&users).Error; err != nil {
			log.Printf("Warning: Failed to get sales users: %v", err)
		}
	}

	// If still no users, get all users as fallback
	if len(users) == 0 {
		if err := database.DB.Find(&users).Error; err != nil {
			return err
		}
	}

	if len(users) == 0 {
		log.Println("Warning: No users found, skipping deal seeding")
		return nil
	}

	log.Printf("Found %d sales users for deal seeding", len(users))

	userMap := make(map[string]user.User)
	for _, u := range users {
		userMap[u.ID] = u
	}

	// Get accounts
	var accounts []account.Account
	if err := database.DB.Find(&accounts).Error; err != nil {
		return err
	}
	if len(accounts) == 0 {
		log.Println("Warning: No accounts found, skipping deal seeding")
		return nil
	}

	// Get contacts
	var contacts []contact.Contact
	if err := database.DB.Find(&contacts).Error; err != nil {
		return err
	}

	// Get pipeline stages
	var stages []pipeline.PipelineStage
	if err := database.DB.Find(&stages).Error; err != nil {
		return err
	}
	if len(stages) == 0 {
		log.Println("Warning: No pipeline stages found, skipping deal seeding")
		return nil
	}

	// Find closed pipeline stages used by this simple transactional seeder.
	var wonStage, lostStage *pipeline.PipelineStage
	stageProbability := func(stage *pipeline.PipelineStage) int {
		if stage == nil {
			return 0
		}
		return stage.Probability
	}
	stageStatus := func(stage *pipeline.PipelineStage) string {
		if stage == nil {
			return "open"
		}
		if stage.IsWon {
			return "won"
		}
		if stage.IsLost {
			return "lost"
		}
		return "open"
	}
	for i := range stages {
		switch stages[i].Code {
		case "closed_won":
			wonStage = &stages[i]
		case "closed_lost":
			lostStage = &stages[i]
		}
	}
	if wonStage == nil {
		log.Println("Warning: Closed won stage not found, skipping deal seeding")
		return nil
	}

	// Get products for product items
	var products []product.Product
	if err := database.DB.Where("deleted_at IS NULL").Limit(20).Find(&products).Error; err != nil {
		log.Printf("Warning: Failed to get products: %v", err)
	}
	if len(products) == 0 {
		log.Println("Warning: No products found, attempting to seed products...")
		if err := SeedProducts(); err != nil {
			log.Printf("Warning: Failed to seed products fallback: %v", err)
		} else {
			// Refresh products
			if err := database.DB.Where("deleted_at IS NULL").Limit(20).Find(&products).Error; err != nil {
				log.Printf("Warning: Failed to get new products: %v", err)
			}
		}

		if len(products) == 0 {
			log.Println("Warning: Still no products found, deals will be created without product items")
		}
	}

	// Get admin user for CreatedBy
	var adminUser user.User
	if err := database.DB.Where("email = ?", "admin@example.com").First(&adminUser).Error; err != nil {
		if len(users) > 0 {
			adminUser = users[0]
		}
	}

	now := time.Now()
	// Create deals for last month and current month to ensure data appears
	lastMonthStart := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location())
	currentMonthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// Create sample deals for sales performance data
	// Assign deals to different users to generate performance metrics
	deals := []pipeline.Deal{}

	// Get bricks for deals
	var bricks []brickdomain.Brick
	if err := database.DB.Find(&bricks).Error; err != nil {
		log.Printf("Warning: Failed to get bricks: %v", err)
	}

	// Create deals for more users (up to 10 users)
	maxUsers := 1
	if len(users) < maxUsers {
		maxUsers = len(users)
	}

	for i := 0; i < maxUsers; i++ {
		userID := users[i].ID
		accountIndex := i % len(accounts)

		// Get user's brick if available
		var userBrickID *string
		if i < len(users) {
			user := users[i]
			if user.BrickID != nil {
				userBrickID = user.BrickID
			}
		}

		// Create multiple won deals per user with different values
		// FOKUS: Closed Won deals (for revenue calculation)
		wonDealValues := []int64{
			1000000000, // Rp 10,000,000
			2000000000, // Rp 20,000,000
		}

		for j, value := range wonDealValues {
			if accountIndex+j >= len(accounts) {
				accountIndex = 0
			}

			var contactID string
			if len(contacts) > 0 {
				contactID = contacts[(accountIndex+j)%len(contacts)].ID
			}

			// Distribute deals across last month and current month
			var closeDate time.Time
			if j%2 == 0 {
				// Last month deals
				daysAgo := (j/2)*7 + 5 // Stagger: 5, 12, 19, 26 days ago
				closeDate = lastMonthStart.AddDate(0, 0, daysAgo)
			} else {
				// Current month deals
				daysAgo := ((j-1)/2)*7 + 3 // Stagger: 3, 10, 17, 24 days ago
				closeDate = currentMonthStart.AddDate(0, 0, daysAgo)
				if closeDate.After(now) {
					closeDate = now.AddDate(0, 0, -daysAgo)
				}
			}

			var contactIDPtr *string
			if contactID != "" {
				contactIDPtr = &contactID
			}
			deal := pipeline.Deal{
				Title:             "Pharmaceutical Supply Agreement",
				Description:       "Annual pharmaceutical supply contract - Closed Won",
				AccountID:         accounts[(accountIndex+j)%len(accounts)].ID,
				ContactID:         contactIDPtr,
				StageID:           wonStage.ID,
				Value:             value,
				Probability:       stageProbability(wonStage),
				ExpectedCloseDate: &closeDate,
				ActualCloseDate:   &closeDate,
				AssignedTo:        &userID,
				BrickID:           userBrickID,
				Status:            stageStatus(wonStage),
				Source:            "referral",
				CloseReason:       []string{"Annual supply contract approved", "Product fit confirmed by procurement"}[j%2],
				Notes:             "Sample closed won deal for sales performance testing",
				CreatedBy:         adminUser.ID,
			}

			deals = append(deals, deal)
		}

		// FOKUS: Create exactly two closed lost deals (for conversion rate calculation)
		lostDealValues := []int64{
			500000000, // Rp 5,000,000
			750000000, // Rp 7,500,000
		}
		if lostStage != nil {
			for j, value := range lostDealValues {
				if accountIndex+j >= len(accounts) {
					accountIndex = 0
				}

				var contactID string
				if len(contacts) > 0 {
					contactID = contacts[(accountIndex+j+len(wonDealValues))%len(contacts)].ID
				}

				// Distribute lost deals across last month and current month
				var closeDate time.Time
				if j%2 == 0 {
					// Last month deals
					daysAgo := (j/2)*7 + 3 // Stagger: 3, 10 days ago
					closeDate = lastMonthStart.AddDate(0, 0, daysAgo)
				} else {
					// Current month deals
					daysAgo := ((j-1)/2)*7 + 2 // Stagger: 2, 9 days ago
					closeDate = currentMonthStart.AddDate(0, 0, daysAgo)
					if closeDate.After(now) {
						closeDate = now.AddDate(0, 0, -daysAgo)
					}
				}

				var contactIDPtr *string
				if contactID != "" {
					contactIDPtr = &contactID
				}
				deal := pipeline.Deal{
					Title:             "Pharmaceutical Supply Agreement",
					Description:       "Annual pharmaceutical supply contract - Closed Lost",
					AccountID:         accounts[(accountIndex+j+len(wonDealValues))%len(accounts)].ID,
					ContactID:         contactIDPtr,
					StageID:           lostStage.ID,
					Value:             value,
					Probability:       stageProbability(lostStage),
					ExpectedCloseDate: &closeDate,
					ActualCloseDate:   &closeDate,
					AssignedTo:        &userID,
					BrickID:           userBrickID,
					Status:            stageStatus(lostStage),
					Source:            "referral",
					CloseReason:       []string{"Price higher than competitor", "Budget postponed to next quarter"}[j%2],
					Notes:             "Sample closed lost deal for conversion rate calculation",
					CreatedBy:         adminUser.ID,
				}

				deals = append(deals, deal)
			}
		}
	}

	// Create deals with product items
	createdDeals := []pipeline.Deal{}
	for _, deal := range deals {
		// Calculate product items and total value for won deals
		var productItems []pipeline.DealProductItem
		var calculatedValue int64

		if deal.Status == "won" && len(products) > 0 {
			// Add 2-4 products per won deal with variety
			numProducts := 2 + rand.Intn(3)           // 2-4 products
			selectedProducts := make(map[string]bool) // Track selected products to avoid duplicates

			for k := 0; k < numProducts && k < len(products); k++ {
				// Select random product (avoid duplicates)
				var selectedProduct product.Product
				attempts := 0
				for {
					selectedProduct = products[rand.Intn(len(products))]
					if !selectedProducts[selectedProduct.ID] {
						selectedProducts[selectedProduct.ID] = true
						break
					}
					attempts++
					if attempts > 10 {
						// Fallback: use first available product
						for _, p := range products {
							if !selectedProducts[p.ID] {
								selectedProduct = p
								selectedProducts[p.ID] = true
								break
							}
						}
						break
					}
				}

				// Varied quantity based on product price
				var quantity int
				if selectedProduct.Price > 10000000 { // Expensive products (> Rp 100k)
					quantity = 1 + rand.Intn(5) // 1-5 units
				} else if selectedProduct.Price > 1000000 { // Mid-range products (> Rp 10k)
					quantity = 5 + rand.Intn(20) // 5-25 units
				} else { // Cheap products
					quantity = 10 + rand.Intn(50) // 10-60 units
				}

				// Use product price with some variation (±5%)
				priceVariation := float64(selectedProduct.Price) * (0.95 + rand.Float64()*0.1)
				unitPrice := int64(priceVariation)

				// Random discount (0-10% of subtotal)
				subtotalBeforeDiscount := unitPrice * int64(quantity)
				discountPercent := rand.Float64() * 0.1 // 0-10%
				discountAmount := int64(float64(subtotalBeforeDiscount) * discountPercent)
				subtotal := subtotalBeforeDiscount - discountAmount

				productItems = append(productItems, pipeline.DealProductItem{
					ProductID:      selectedProduct.ID,
					ProductName:    selectedProduct.Name,
					ProductSKU:     selectedProduct.SKU,
					UnitPrice:      unitPrice,
					Quantity:       quantity,
					DiscountAmount: discountAmount,
					Subtotal:       subtotal,
				})

				calculatedValue += subtotal
			}

			// Update deal value to match product items total
			if calculatedValue > 0 {
				deal.Value = calculatedValue
			}
		}

		// Create deal
		if err := database.DB.Create(&deal).Error; err != nil {
			log.Printf("Warning: Failed to create deal: %v", err)
			continue
		}

		// Create product items for won deals
		if deal.Status == "won" && len(productItems) > 0 {
			for i := range productItems {
				productItems[i].DealID = deal.ID
			}

			if err := database.DB.Create(&productItems).Error; err != nil {
				log.Printf("Warning: Failed to create product items for deal %s: %v", deal.ID, err)
			} else {
				assignedToVal := ""
				if deal.AssignedTo != nil {
					assignedToVal = *deal.AssignedTo
				}
				log.Printf("Created deal: %s (Value: %d, AssignedTo: %s) with %d product items",
					deal.Title, deal.Value, assignedToVal, len(productItems))

				// Create CustomerPurchaseHistory record for won deals
				items := make([]customer_purchase.PurchaseProduct, 0, len(productItems))
				for _, item := range productItems {
					items = append(items, customer_purchase.PurchaseProduct{
						ProductID:   item.ProductID,
						ProductName: item.ProductName,
						ProductSKU:  item.ProductSKU,
						Quantity:    item.Quantity,
						UnitPrice:   item.UnitPrice,
						Subtotal:    item.Subtotal,
					})
				}

				data, _ := json.Marshal(items)
				history := customer_purchase.CustomerPurchaseHistory{
					AccountID:    deal.AccountID,
					DealID:       deal.ID,
					PurchaseDate: *deal.ActualCloseDate,
					TotalAmount:  deal.Value,
					TotalItems:   len(productItems),
					Products:     datatypes.JSON(data),
					SalesRepID:   deal.AssignedTo,
					SourceType:   "pipeline",
				}

				// Set SalesRepName from user map
				if deal.AssignedTo != nil {
					if u, ok := userMap[*deal.AssignedTo]; ok {
						history.SalesRepName = u.Name
					}
				}

				if err := database.DB.Create(&history).Error; err != nil {
					log.Printf("Warning: Failed to create purchase history for deal %s: %v", deal.ID, err)
				}
			}
		} else {
			assignedToVal := ""
			if deal.AssignedTo != nil {
				assignedToVal = *deal.AssignedTo
			}
			log.Printf("Created deal: %s (Value: %d, AssignedTo: %s)", deal.Title, deal.Value, assignedToVal)
		}

		createdDeals = append(createdDeals, deal)
	}

	log.Printf("Seeded %d deals successfully", len(createdDeals))
	return nil
}
