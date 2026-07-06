package seeders

import (
	"fmt"
	"log"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/config"
	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/account"
	"github.com/gilabs/crm-healthcare/api/internal/domain/category"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"github.com/gilabs/crm-healthcare/api/internal/repository/postgres/brick"
	accountrepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/account"
	userrepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/user"
	brickservice "github.com/gilabs/crm-healthcare/api/internal/service/brick"
	"github.com/gilabs/crm-healthcare/api/pkg/geocoding"
)

// SeedAccounts seeds initial accounts
func SeedAccounts() error {
	// Check if accounts already exist (for logging purposes)
	var count int64
	database.DB.Model(&account.Account{}).Count(&count)
	log.Printf("Current account count: %d", count)

	// Get sales users that are assigned to bricks (prioritize brick-assigned sales)
	var salesUsers []user.User
	// First, try to get sales users that are assigned to bricks
	if err := database.DB.Where("users.brick_id IS NOT NULL AND users.status = ?", "active").
		Joins("JOIN roles ON users.role_id = roles.id").
		Where("roles.code = ?", "sales").
		Find(&salesUsers).Error; err != nil {
		log.Printf("Warning: Failed to get sales users with bricks: %v", err)
	}

	// If no sales users with bricks found, get all active sales users
	if len(salesUsers) == 0 {
		if err := database.DB.Where("users.status = ?", "active").
			Joins("JOIN roles ON users.role_id = roles.id").
			Where("roles.code IN (?)", []string{"sales", "sales_manager"}).
			Find(&salesUsers).Error; err != nil {
			log.Printf("Warning: Failed to get sales users: %v", err)
		}
	}

	// If still no users, use admin as fallback
	var adminUser user.User
	var adminUserID *string
	if len(salesUsers) == 0 {
		if err := database.DB.Where("email = ?", "admin@example.com").First(&adminUser).Error; err == nil {
			adminUserID = &adminUser.ID
		} else {
			// Get any active user as last resort
			var fallbackUsers []user.User
			if err := database.DB.Where("status = ?", "active").Limit(1).Find(&fallbackUsers).Error; err == nil && len(fallbackUsers) > 0 {
				adminUserID = &fallbackUsers[0].ID
			}
		}
	}

	log.Printf("Found %d sales users for account seeding", len(salesUsers))

	brickRepo := brick.NewRepository(database.DB)
	userRepo := userrepo.NewRepository(database.DB)
	accountRepo := accountrepo.NewRepository(database.DB)
	brickHelper := brickservice.NewBrickHelper(userRepo, brickRepo, accountRepo)

	// Get categories by code
	var hospitalCategory, clinicCategory, pharmacyCategory category.Category
	if err := database.DB.Where("code = ?", "HOSPITAL").First(&hospitalCategory).Error; err != nil {
		return err
	}
	if err := database.DB.Where("code = ?", "CLINIC").First(&clinicCategory).Error; err != nil {
		return err
	}
	if err := database.DB.Where("code = ?", "PHARMACY").First(&pharmacyCategory).Error; err != nil {
		return err
	}

	// Helper function to get assigned user ID (distribute across sales users)
	getAssignedUserID := func(index int) *string {
		if len(salesUsers) == 0 {
			return adminUserID
		}
		userIndex := index % len(salesUsers)
		return &salesUsers[userIndex].ID
	}

	// Initialize geocoding service if enabled
	var geocodingSvc *geocoding.GeocodingService
	geocodingEnabled := false
	if config.AppConfig != nil && config.AppConfig.Geocoding.Enabled {
		geocodingSvc = geocoding.NewGeocodingService(
			config.AppConfig.Geocoding.Provider,
			config.AppConfig.Geocoding.APIKey,
		)
		geocodingEnabled = true
		log.Println("Geocoding enabled for account seeding")
	} else {
		log.Println("Geocoding disabled, accounts will be created without coordinates")
	}

	// Helper function to geocode address with fallback
	geocodeAddress := func(addr, city, province string) (*float64, *float64) {
		if !geocodingEnabled || geocodingSvc == nil {
			return nil, nil
		}
		// Use fallback strategy for better success rate
		result, err := geocodingSvc.GeocodeAddressWithFallback(addr, city, province)
		if err != nil {
			log.Printf("Warning: Failed to geocode address for '%s, %s, %s' after fallback attempts: %v", addr, city, province, err)
			return nil, nil
		}
		return &result.Latitude, &result.Longitude
	}

	seedBaseAccounts := func() error {
		accounts := []account.Account{
			{
				Name:       "RSUP Dr Kariadi",
				CategoryID: hospitalCategory.ID,
				Address:    "Jl. Dr. Sutomo No.16, Randusari, Kec. Semarang Selatan",
				City:       "Semarang",
				Province:   "Jawa Tengah",
				Phone:      "+62248413476",
				Email:      "info@rskariadi.co.id",
				Latitude:   float64Ptr(-6.994590),
				Longitude:  float64Ptr(110.407750),
				Status:     "active",
				AssignedTo: getAssignedUserID(0),
			},
			{
				Name:       "RS Telogorejo",
				CategoryID: hospitalCategory.ID,
				Address:    "Jl. Kh Ahmad Dahlan, Pekunden, Kec. Semarang Tengah",
				City:       "Semarang",
				Province:   "Jawa Tengah",
				Phone:      "+622486466000",
				Email:      "info@rstelogorejo.com",
				Latitude:   float64Ptr(-6.987337),
				Longitude:  float64Ptr(110.426417),
				Status:     "active",
				AssignedTo: getAssignedUserID(1),
			},
			{
				Name:       "RS St Elisabeth",
				CategoryID: hospitalCategory.ID,
				Address:    "Jl. Kawi No.1, Tegalsari, Kec. Candisari",
				City:       "Semarang",
				Province:   "Jawa Tengah",
				Phone:      "+62248310076",
				Email:      "info@rs-elisabeth.com",
				Latitude:   float64Ptr(-7.008283),
				Longitude:  float64Ptr(110.419920),
				Status:     "active",
				AssignedTo: getAssignedUserID(2),
			},
			{
				Name:       "RS Bhakti Wira Tamtama",
				CategoryID: hospitalCategory.ID,
				Address:    "Jl. Dr. Sutomo No.17, Barusari, Kec. Semarang Selatan",
				City:       "Semarang",
				Province:   "Jawa Tengah",
				Phone:      "+62248315143",
				Email:      "info@rsbhaktiwiratamtama.com",
				Latitude:   float64Ptr(-6.986269),
				Longitude:  float64Ptr(110.408313),
				Status:     "active",
				AssignedTo: getAssignedUserID(3),
			},
			{
				Name:       "RS Roemani Muhammadiyah",
				CategoryID: hospitalCategory.ID,
				Address:    "Jl. Wonodri No.22, Wonodri, Kec. Semarang Selatan",
				City:       "Semarang",
				Province:   "Jawa Tengah",
				Phone:      "+62248444623",
				Email:      "info@rsroemani.com",
				Latitude:   float64Ptr(-7.001093),
				Longitude:  float64Ptr(110.425719),
				Status:     "active",
				AssignedTo: getAssignedUserID(4),
			},
		}

		// Check if accounts table exists before seeding
		if !database.DB.Migrator().HasTable(&account.Account{}) {
			return fmt.Errorf("accounts table does not exist. Migration might have failed or table was dropped mid-process")
		}

		for i, acc := range accounts {
			// If geocoding is available AND the hardcoded coordinates are nil, attempt geocoding.
			// Hardcoded coordinates take priority and skip geocoding.
			didGeocode := false
			if geocodingEnabled && acc.Latitude == nil && acc.Longitude == nil && (acc.Address != "" || acc.City != "" || acc.Province != "") {
				lat, lng := geocodeAddress(acc.Address, acc.City, acc.Province)
				if lat != nil && lng != nil {
					acc.Latitude = lat
					acc.Longitude = lng
					didGeocode = true
					log.Printf("Geocoded address for %s: lat=%.6f, lng=%.6f", acc.Name, *lat, *lng)
				} else {
					log.Printf("Warning: Could not geocode address for %s", acc.Name)
				}
			}

			// Upsert by phone so re-running the seeder is idempotent and also corrects
			// any existing rows that may have been created without a name.
			// Name is always included in Assign to ensure it is set correctly.
			assignData := map[string]interface{}{
				"name":        acc.Name,
				"category_id": acc.CategoryID,
				"address":     acc.Address,
				"city":        acc.City,
				"province":    acc.Province,
				"email":       acc.Email,
				"status":      acc.Status,
				"assigned_to": acc.AssignedTo,
			}

			brickID, brickErr := brickHelper.EnsureBrickIDForLocation(acc.Province, acc.City)
			if brickErr != nil {
				log.Printf("Warning: Failed to ensure brick for account %s (%s, %s): %v", acc.Name, acc.City, acc.Province, brickErr)
			} else if brickID != nil {
				assignData["brick_id"] = brickID
			}
			if acc.Latitude != nil {
				assignData["latitude"] = acc.Latitude
			}
			if acc.Longitude != nil {
				assignData["longitude"] = acc.Longitude
			}

			if err := database.DB.
				Where("phone = ?", acc.Phone).
				Assign(assignData).
				FirstOrCreate(&account.Account{}).Error; err != nil {
				return err
			}

			coordsInfo := ""
			if acc.Latitude != nil && acc.Longitude != nil {
				coordsInfo = fmt.Sprintf(", lat=%.6f, lng=%.6f", *acc.Latitude, *acc.Longitude)
			}
			log.Printf("Upserted account: %s (category_id: %s%s)", acc.Name, acc.CategoryID, coordsInfo)

			// Add small delay between geocoding requests to respect rate limits (Nominatim: 1 req/sec).
			if didGeocode && i < len(accounts)-1 {
				time.Sleep(1100 * time.Millisecond)
			}
		}

		log.Println("Base accounts seeded successfully")
		return nil
	}

	// Seed base accounts using upsert — safe to re-run to backfill missing coordinates.
	if err := seedBaseAccounts(); err != nil {
		return err
	}

	log.Println("Accounts seeding completed")
	return nil
}
