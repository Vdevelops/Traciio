package seeders

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead"
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead_status"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
)

// SeedLeads seeds initial leads
func SeedLeads() error {
	// Check if leads already exist
	var count int64
	database.DB.Model(&lead.Lead{}).Count(&count)
	if count > 0 {
		log.Println("Leads already seeded, skipping...")
		return nil
	}

	// Get lead statuses
	var leadStatuses []lead_status.LeadStatus
	if err := database.DB.Where("is_active = ?", true).Find(&leadStatuses).Error; err != nil {
		log.Printf("Warning: Failed to get lead statuses: %v", err)
	}

	// Create map of lead status codes to IDs
	statusMap := make(map[string]string)
	for _, status := range leadStatuses {
		statusMap[status.Code] = status.ID
	}

	// Get users for assigned_to (prioritize sales users)
	var users []user.User
	if err := database.DB.Where("status = ?", "active").
		Joins("JOIN roles ON users.role_id = roles.id").
		Where("roles.code IN (?)", []string{"sales", "sales_manager"}).
		Find(&users).Error; err != nil {
		log.Printf("Warning: Failed to get sales users: %v", err)
		// Fallback to all users
		if err := database.DB.Find(&users).Error; err != nil {
			return err
		}
	}

	if len(users) == 0 {
		log.Println("Warning: No users found, skipping lead seeding")
		return nil
	}

	// Get admin user for created_by
	var adminUser user.User
	if err := database.DB.Where("email = ?", "admin@example.com").First(&adminUser).Error; err != nil {
		log.Printf("Warning: Admin user not found, using first user for created_by: %v", err)
		if len(users) > 0 {
			adminUser = users[0]
		}
	}

	// Assign users for leads (cycle through available users)
	userIndex := 0
	getNextUser := func() *string {
		if len(users) == 0 {
			return nil
		}
		userID := users[userIndex].ID
		userIndex = (userIndex + 1) % len(users)
		return &userID
	}

	// Helper to get lead status ID
	getLeadStatusID := func(code string) *string {
		if id, ok := statusMap[code]; ok {
			return &id
		}
		return nil
	}

	// Generate a small lead dataset for development/demo usage.
	leadTemplates := []struct {
		FirstName   string
		LastName    string
		CompanyName string
		Email       string
		Phone       string
		JobTitle    string
		Industry    string
		LeadSource  string
		StatusCode  string
		LeadScore   int
		Address     string
		City        string
		Province    string
		PostalCode  string
		Website     string
		Notes       string
	}{
		{"Budi", "Santoso", "PT Healthcare Indonesia", "budi.santoso@healthcare.id", "081234567890", "Director", "Healthcare", "website", "new", 50, "Jl. Sudirman No. 123", "Semarang", "Jawa Tengah", "50125", "https://healthcare.id", "Interested in pharmaceutical products. Requested product catalog."},
		{"Siti", "Rahayu", "Rumah Sakit Umum Daerah", "siti.rahayu@rsud.example.com", "081234567891", "Procurement Manager", "Healthcare", "referral", "contacted", 65, "Jl. Gatot Subroto No. 456", "Semarang", "Jawa Tengah", "50125", "", "Referred by existing client. Looking for medical equipment."},
		{"Ahmad", "Fauzi", "Klinik Sehat Jaya", "ahmad.fauzi@kliniksehat.com", "081234567892", "Owner", "Healthcare", "cold_call", "qualified", 75, "Jl. Merdeka No. 789", "Semarang", "Jawa Tengah", "50125", "", "Qualified lead. Budget confirmed. Ready for proposal."},
	}

	leads := []lead.Lead{}
	for _, template := range leadTemplates {
		leads = append(leads, lead.Lead{
			FirstName:    template.FirstName,
			LastName:     template.LastName,
			CompanyName:  template.CompanyName,
			Email:        template.Email,
			Phone:        template.Phone,
			JobTitle:     template.JobTitle,
			Industry:     template.Industry,
			LeadSource:   template.LeadSource,
			LeadStatus:   template.StatusCode,
			LeadStatusID: getLeadStatusID(template.StatusCode),
			LeadScore:    template.LeadScore,
			AssignedTo:   getNextUser(),
			Notes:        template.Notes,
			Address:      template.Address,
			City:         template.City,
			Province:     template.Province,
			PostalCode:   template.PostalCode,
			Country:      "Indonesia",
			Website:      template.Website,
			CreatedBy:    adminUser.ID,
		})
	}

	// Use Omit to skip empty UUID fields (AccountID, ContactID, OpportunityID, ConvertedBy)
	// This prevents PostgreSQL error: invalid input syntax for type uuid: ""
	// By omitting these fields, GORM will set them as NULL instead of empty string
	if err := database.DB.Omit("AccountID", "ContactID", "OpportunityID", "ConvertedBy", "ConvertedAt").Create(&leads).Error; err != nil {
		return err
	}

	log.Printf("✅ Seeded %d leads successfully", len(leads))
	return nil
}
