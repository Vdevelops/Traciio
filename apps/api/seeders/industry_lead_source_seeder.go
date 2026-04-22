package seeders

import (
	"log"
	"strings"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/industry"
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead_source"
)

// SeedIndustries seeds initial industries
func SeedIndustries() error {
	log.Println("Seeding industries...")

	// Check if industries already exist
	var count int64
	database.DB.Model(&industry.Industry{}).Count(&count)
	if count > 0 {
		log.Println("Industries already seeded, skipping...")
		return nil
	}

	// Get admin user for created_by
	var adminUser struct {
		ID string
	}
	if err := database.DB.Table("users").Where("email = ?", "admin@example.com").First(&adminUser).Error; err != nil {
		log.Printf("Warning: Admin user not found, using empty string for created_by: %v", err)
	}

	createdBy := adminUser.ID
	if createdBy == "" {
		// Try to get any admin user
		if err := database.DB.Table("users").Where("role_id IN (SELECT id FROM roles WHERE code = 'admin')").First(&adminUser).Error; err == nil {
			createdBy = adminUser.ID
		}
	}

	// Common industries in Indonesia healthcare
	industries := []struct {
		Name        string
		Code        string
		Description string
		Order       int
	}{
		{"Healthcare", "HEALTHCARE", "General healthcare industry", 1},
		{"Pharmaceutical", "PHARMACEUTICAL", "Pharmaceutical and medicine industry", 2},
		{"Hospital", "HOSPITAL", "Hospital and medical facility", 3},
		{"Clinic", "CLINIC", "Clinic and medical practice", 4},
		{"Medical Device", "MEDICAL_DEVICE", "Medical device and equipment", 5},
		{"Biotechnology", "BIOTECHNOLOGY", "Biotechnology and research", 6},
		{"Research & Development", "RND", "Research and development", 7},
		{"Medical Equipment", "MEDICAL_EQUIPMENT", "Medical equipment manufacturing", 8},
		{"Health Insurance", "HEALTH_INSURANCE", "Health insurance provider", 9},
		{"Telemedicine", "TELEMEDICINE", "Telemedicine and digital health", 10},
		{"Other", "OTHER", "Other healthcare related industry", 99},
	}

	for _, ind := range industries {
		industry := &industry.Industry{
			Name:        ind.Name,
			Code:        ind.Code,
			Description: ind.Description,
			Order:       ind.Order,
			IsActive:    true,
			CreatedBy:   createdBy,
		}

		if err := database.DB.Create(industry).Error; err != nil {
			log.Printf("Error creating industry %s: %v", ind.Name, err)
			return err
		}
		log.Printf("Created industry: %s", ind.Name)
	}

	log.Println("Industries seeded successfully!")
	return nil
}

// SeedLeadSources seeds initial lead sources
func SeedLeadSources() error {
	log.Println("Seeding lead sources...")

	// Check if lead sources already exist
	var count int64
	database.DB.Model(&lead_source.LeadSource{}).Count(&count)
	if count > 0 {
		log.Println("Lead sources already seeded, skipping...")
		return nil
	}

	// Get admin user for created_by
	var adminUser struct {
		ID string
	}
	if err := database.DB.Table("users").Where("email = ?", "admin@example.com").First(&adminUser).Error; err != nil {
		log.Printf("Warning: Admin user not found, using empty string for created_by: %v", err)
	}

	createdBy := adminUser.ID
	if createdBy == "" {
		// Try to get any admin user
		if err := database.DB.Table("users").Where("role_id IN (SELECT id FROM roles WHERE code = 'admin')").First(&adminUser).Error; err == nil {
			createdBy = adminUser.ID
		}
	}

	// Common lead sources
	leadSources := []struct {
		Name        string
		Code        string
		Description string
		Order       int
	}{
		{"Website", "WEBSITE", "Lead from website form or inquiry", 1},
		{"Referral", "REFERRAL", "Lead from referral or recommendation", 2},
		{"Cold Call", "COLD_CALL", "Lead from cold calling campaign", 3},
		{"Event", "EVENT", "Lead from trade show or event", 4},
		{"Social Media", "SOCIAL_MEDIA", "Lead from social media platforms", 5},
		{"Email Campaign", "EMAIL_CAMPAIGN", "Lead from email marketing campaign", 6},
		{"Partner", "PARTNER", "Lead from business partner", 7},
		{"Other", "OTHER", "Other lead source", 99},
	}

	for _, ls := range leadSources {
		leadSource := &lead_source.LeadSource{
			Name:        ls.Name,
			Code:        strings.ToUpper(ls.Code),
			Description: ls.Description,
			Order:       ls.Order,
			IsActive:    true,
			CreatedBy:   createdBy,
		}

		if err := database.DB.Create(leadSource).Error; err != nil {
			log.Printf("Error creating lead source %s: %v", ls.Name, err)
			return err
		}
		log.Printf("Created lead source: %s", ls.Name)
	}

	log.Println("Lead sources seeded successfully!")
	return nil
}

