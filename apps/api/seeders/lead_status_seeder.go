package seeders

import (
	"log"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/lead_status"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"gorm.io/gorm"
)

// SeedLeadStatuses seeds default lead statuses with scoring
func SeedLeadStatuses(db *gorm.DB) error {
	// Check if lead statuses already exist
	var count int64
	if err := db.Model(&lead_status.LeadStatus{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil // Already seeded
	}

	// Get admin user for CreatedBy (or use first user if admin not found)
	var createdBy string
	var adminUser user.User
	if err := db.Where("email = ?", "admin@example.com").First(&adminUser).Error; err != nil {
		// Try to get any user
		var firstUser user.User
		if err := db.First(&firstUser).Error; err == nil {
			createdBy = firstUser.ID
			log.Printf("Warning: Admin user not found, using first user for lead status created_by: %s", firstUser.Email)
		} else {
			// No users found, use empty string (will be NULL in database)
			createdBy = ""
			log.Printf("Warning: No users found, lead statuses will be created with NULL created_by")
		}
	} else {
		createdBy = adminUser.ID
	}

	statuses := []lead_status.LeadStatus{
		{
			Name:        "New",
			Code:        "new",
			Description: "Brand new lead, not yet contacted",
			Score:       5,
			Color:       "#94A3B8", // Slate
			Order:       1,
			IsActive:    true,
			IsDefault:   true, // Default status for new leads
			IsConverted: false,
			CreatedBy:   createdBy,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Contacted",
			Code:        "contacted",
			Description: "Initial contact made with the lead",
			Score:       15,
			Color:       "#60A5FA", // Blue
			Order:       2,
			IsActive:    true,
			IsDefault:   false,
			IsConverted: false,
			CreatedBy:   createdBy,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Engaged",
			Code:        "engaged",
			Description: "Lead is actively engaging with us",
			Score:       30,
			Color:       "#A78BFA", // Purple
			Order:       3,
			IsActive:    true,
			IsDefault:   false,
			IsConverted: false,
			CreatedBy:   createdBy,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Interested",
			Code:        "interested",
			Description: "Lead has shown strong interest in our products/services",
			Score:       50,
			Color:       "#34D399", // Green
			Order:       4,
			IsActive:    true,
			IsDefault:   false,
			IsConverted: false,
			CreatedBy:   createdBy,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Qualified",
			Code:        "qualified",
			Description: "Lead meets qualification criteria and ready for conversion",
			Score:       75,
			Color:       "#10B981", // Emerald
			Order:       5,
			IsActive:    true,
			IsDefault:   false,
			IsConverted: false,
			CreatedBy:   createdBy,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Converted",
			Code:        "converted",
			Description: "Lead has been converted to opportunity",
			Score:       100,
			Color:       "#3B82F6", // Blue
			Order:       6,
			IsActive:    true,
			IsDefault:   false,
			IsConverted: true, // Mark as converted status
			CreatedBy:   createdBy,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Nurturing",
			Code:        "nurturing",
			Description: "Lead needs more time and education",
			Score:       25,
			Color:       "#FBBF24", // Amber
			Order:       7,
			IsActive:    true,
			IsDefault:   false,
			IsConverted: false,
			CreatedBy:   createdBy,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Unqualified",
			Code:        "unqualified",
			Description: "Lead does not meet qualification criteria",
			Score:       0,
			Color:       "#F97316", // Orange
			Order:       8,
			IsActive:    true,
			IsDefault:   false,
			IsConverted: false,
			CreatedBy:   createdBy,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Disqualified",
			Code:        "disqualified",
			Description: "Lead is not a good fit for our business",
			Score:       0,
			Color:       "#EF4444", // Red
			Order:       9,
			IsActive:    true,
			IsDefault:   false,
			IsConverted: false,
			CreatedBy:   createdBy,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Lost",
			Code:        "lost",
			Description: "Lead was lost to competitor or not interested",
			Score:       0,
			Color:       "#DC2626", // Red-600
			Order:       10,
			IsActive:    true,
			IsDefault:   false,
			IsConverted: false,
			CreatedBy:   createdBy,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	return db.Create(&statuses).Error
}
