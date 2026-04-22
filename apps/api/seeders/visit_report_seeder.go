package seeders

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/account"
	"github.com/gilabs/crm-healthcare/api/internal/domain/contact"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"github.com/gilabs/crm-healthcare/api/internal/domain/visit_report"
	"gorm.io/datatypes"
)

// SeedVisitReports seeds initial visit reports
func SeedVisitReports() error {
	// Check if visit reports already exist for the current month.
	// We intentionally do NOT use a global count, because that would prevent
	// current-month metrics from ever being seeded after the calendar moves.
	nowCheck := time.Now()
	currentMonthStartCheck := time.Date(nowCheck.Year(), nowCheck.Month(), 1, 0, 0, 0, 0, nowCheck.Location())
	var count int64
	database.DB.Model(&visit_report.VisitReport{}).Where("visit_date >= ?", currentMonthStartCheck).Count(&count)
	if count > 0 {
		log.Println("Visit reports already exist for current month, skipping main seeding...")
		// Removed extra seeding logic
		return nil
	}

	// Get accounts
	var accounts []account.Account
	if err := database.DB.Find(&accounts).Error; err != nil {
		return err
	}
	if len(accounts) == 0 {
		log.Println("Warning: No accounts found, skipping visit report seeding")
		return nil
	}

	// Get contacts
	var contacts []contact.Contact
	if err := database.DB.Find(&contacts).Error; err != nil {
		return err
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
		log.Println("Warning: No users found, skipping visit report seeding")
		return nil
	}

	log.Printf("Found %d sales users for visit report seeding", len(users))

	// Helper function to marshal location
	marshalLocation := func(lat, lng float64, address string) datatypes.JSON {
		loc := visit_report.Location{
			Latitude:  lat,
			Longitude: lng,
			Address:   address,
		}
		bytes, _ := json.Marshal(loc)
		return bytes
	}

	// Helper function to marshal photos
	marshalPhotos := func(urls []string) datatypes.JSON {
		bytes, _ := json.Marshal(urls)
		return bytes
	}

	// Get current time for check-in/out
	now := time.Now()
	twoDaysAgo := now.AddDate(0, 0, -2)
	oneDayAgo := now.AddDate(0, 0, -1)
	yesterday := now.AddDate(0, 0, -1)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	visitReports := []visit_report.VisitReport{}

	// Visit Report 1: Draft visit to first account
	if len(accounts) > 0 && len(users) > 0 {
		var contactID *string
		if len(contacts) > 0 {
			contactID = &contacts[0].ID
		}
		accountID := accounts[0].ID
		visitReports = append(visitReports, visit_report.VisitReport{
			AccountID:  &accountID, // Convert string to *string
			ContactID:  contactID,
			SalesRepID: users[0].ID,
			VisitDate:  today,
			Purpose:    "Product presentation for new cardiovascular medications",
			Notes:      "Scheduled meeting to discuss new product line. Need to prepare samples and pricing information.",
			Status:     "draft",
		})
	}

	// Visit Report 2: Submitted visit with check-in
	if len(accounts) > 0 && len(users) > 0 {
		var contactID *string
		if len(contacts) > 1 {
			contactID = &contacts[1].ID
		}
		checkInTime := yesterday.Add(9 * time.Hour) // 9 AM yesterday
		accountID := accounts[0].ID
		visitReports = append(visitReports, visit_report.VisitReport{
			AccountID:       &accountID, // Convert string to *string
			ContactID:       contactID,
			SalesRepID:      users[0].ID,
			VisitDate:       yesterday,
			CheckInTime:     &checkInTime,
			CheckInLocation: marshalLocation(-6.2088, 106.8456, "Jl. Salemba Raya No. 6, Jakarta Pusat"),
			Purpose:         "Follow-up meeting regarding previous product inquiry",
			Notes:           "Met with procurement manager. Discussed pricing and delivery schedule. Positive response.",
			Status:          "submitted",
		})
	}

	// Visit Report 3: Approved visit with check-in and check-out
	if len(accounts) > 1 && len(users) > 0 {
		var contactID *string
		if len(contacts) > 2 {
			contactID = &contacts[2].ID
		}
		checkInTime := twoDaysAgo.Add(10 * time.Hour)  // 10 AM two days ago
		checkOutTime := twoDaysAgo.Add(11 * time.Hour) // 11 AM two days ago
		approvedAt := twoDaysAgo.Add(14 * time.Hour)   // 2 PM two days ago
		var approvedBy *string
		if len(users) > 0 {
			approvedBy = &users[0].ID
		}
		accountID := accounts[1].ID
		visitReports = append(visitReports, visit_report.VisitReport{
			AccountID:        &accountID, // Convert string to *string
			ContactID:        contactID,
			SalesRepID:       users[0].ID,
			VisitDate:        twoDaysAgo,
			CheckInTime:      &checkInTime,
			CheckOutTime:     &checkOutTime,
			CheckInLocation:  marshalLocation(-6.1944, 106.8229, "Jl. Diponegoro No. 71, Jakarta Pusat"),
			CheckOutLocation: marshalLocation(-6.1944, 106.8229, "Jl. Diponegoro No. 71, Jakarta Pusat"),
			Purpose:          "Product demonstration and training session",
			Notes:            "Conducted product demo for medical staff. Training session went well. Received positive feedback.",
			Photos:           marshalPhotos([]string{"https://example.com/photos/visit-001.jpg", "https://example.com/photos/visit-002.jpg"}),
			Status:           "approved",
			ApprovedBy:       approvedBy,
			ApprovedAt:       &approvedAt,
		})
	}

	// Visit Report 4: Rejected visit
	if len(accounts) > 2 && len(users) > 0 {
		var contactID *string
		if len(contacts) > 3 {
			contactID = &contacts[3].ID
		}
		checkInTime := oneDayAgo.Add(13 * time.Hour) // 1 PM one day ago
		rejectedAt := oneDayAgo.Add(16 * time.Hour)  // 4 PM one day ago
		var rejectedBy *string
		if len(users) > 0 {
			rejectedBy = &users[0].ID
		}
		rejectionReason := "Incomplete documentation. Missing required photos and detailed notes."
		accountID := accounts[2].ID
		visitReports = append(visitReports, visit_report.VisitReport{
			AccountID:       &accountID, // Convert string to *string
			ContactID:       contactID,
			SalesRepID:      users[0].ID,
			VisitDate:       oneDayAgo,
			CheckInTime:     &checkInTime,
			CheckInLocation: marshalLocation(-6.2297, 106.7986, "Jl. Metro Duta Kav. UE, Jakarta Selatan"),
			Purpose:         "Initial contact and product introduction",
			Notes:           "Brief meeting with contact person. Discussed basic product information.",
			Status:          "rejected",
			ApprovedBy:      rejectedBy,
			ApprovedAt:      &rejectedAt,
			RejectionReason: &rejectionReason,
		})
	}

	// Visit Report 5: Submitted visit with photos
	if len(accounts) > 3 && len(users) > 0 {
		var contactID *string
		if len(contacts) > 4 {
			contactID = &contacts[4].ID
		}
		checkInTime := yesterday.Add(14 * time.Hour) // 2 PM yesterday
		accountID := accounts[3].ID
		visitReports = append(visitReports, visit_report.VisitReport{
			AccountID:       &accountID, // Convert string to *string
			ContactID:       contactID,
			SalesRepID:      users[0].ID,
			VisitDate:       yesterday,
			CheckInTime:     &checkInTime,
			CheckInLocation: marshalLocation(-6.2297, 106.7986, "Jl. Sudirman No. 123, Jakarta Selatan"),
			Purpose:         "Quarterly review meeting",
			Notes:           "Quarterly business review. Discussed sales performance and upcoming promotions. Took photos of product display.",
			Photos:          marshalPhotos([]string{"https://example.com/photos/visit-003.jpg"}),
			Status:          "submitted",
		})
	}

	// Save the exactly 5 visit reports to the database
	for _, vr := range visitReports {
		if err := database.DB.Create(&vr).Error; err != nil {
			return err
		}
		accountIDStr := "nil"
		if vr.AccountID != nil {
			accountIDStr = *vr.AccountID
		}
		log.Printf("Created visit report: %s (id: %s, account_id: %s, status: %s)", vr.Purpose, vr.ID, accountIDStr, vr.Status)
	}

	log.Printf("Visit reports seeded successfully (%d visit reports created)", len(visitReports))

	// Removed extra seeding functions to keep the exact limit of 5 records

	return nil
}

func marshalLocationJSON(lat, lng float64, address string) datatypes.JSON {
	loc := visit_report.Location{
		Latitude:  lat,
		Longitude: lng,
		Address:   address,
	}
	bytes, _ := json.Marshal(loc)
	return bytes
}

func seedRouteOptimizationTestVisitReports() error {
	// Create deterministic, approved visit reports in the current month, with check-in locations
	// exactly matching the ROUTE TEST accounts.

	// Find route-test accounts
	var routeAccounts []account.Account
	if err := database.DB.Where("name LIKE ?", "ROUTE TEST - %").Order("name ASC").Find(&routeAccounts).Error; err != nil {
		return err
	}
	if len(routeAccounts) == 0 {
		return nil
	}

	// Pick a sales rep (prefer active sales/sales_manager)
	var salesRep user.User
	if err := database.DB.Where("users.status = ?", "active").
		Joins("JOIN roles ON users.role_id = roles.id").
		Where("roles.code IN (?)", []string{"sales", "sales_manager"}).
		First(&salesRep).Error; err != nil {
		// Fallback: any user
		if err2 := database.DB.First(&salesRep).Error; err2 != nil {
			return err
		}
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	createdOrUpdated := 0
	for i, acc := range routeAccounts {
		if acc.Latitude == nil || acc.Longitude == nil {
			continue
		}
		accountID := acc.ID
		checkInTime := today.Add(9*time.Hour + time.Duration(i)*10*time.Minute)
		approvedAt := checkInTime.Add(2 * time.Hour)
		approvedBy := salesRep.ID

		notesMarker := fmt.Sprintf("ROUTE_TEST_VISIT::%s", acc.ID)
		purpose := fmt.Sprintf("Route optimization test visit: %s", acc.Name)
		address := acc.Address
		if address == "" {
			address = acc.Name
		}

		vr := visit_report.VisitReport{
			AccountID:       &accountID,
			SalesRepID:      salesRep.ID,
			VisitDate:       today,
			CheckInTime:     &checkInTime,
			CheckInLocation: marshalLocationJSON(*acc.Latitude, *acc.Longitude, address),
			Purpose:         purpose,
			Notes:           notesMarker,
			Status:          "approved",
			ApprovedBy:      &approvedBy,
			ApprovedAt:      &approvedAt,
		}

		if err := database.DB.
			Where("notes = ?", notesMarker).
			Assign(map[string]interface{}{
				"account_id":        vr.AccountID,
				"sales_rep_id":      vr.SalesRepID,
				"visit_date":        vr.VisitDate,
				"check_in_time":     vr.CheckInTime,
				"check_in_location": vr.CheckInLocation,
				"purpose":           vr.Purpose,
				"status":            vr.Status,
				"approved_by":       vr.ApprovedBy,
				"approved_at":       vr.ApprovedAt,
			}).
			FirstOrCreate(&visit_report.VisitReport{}).Error; err != nil {
			return err
		}
		createdOrUpdated++
	}

	if createdOrUpdated > 0 {
		log.Printf("Route-test visit reports ensured: %d", createdOrUpdated)
	}
	return nil
}

// jakartaLocations contains real Jakarta area coordinates for seeding
var jakartaLocations = []struct {
	lat     float64
	lng     float64
	address string
}{
	// Jakarta Pusat
	{-6.1944, 106.8229, "Jl. Diponegoro No. 71, Jakarta Pusat"},
	{-6.2088, 106.8456, "Jl. Salemba Raya No. 6, Jakarta Pusat"},
	{-6.2146, 106.8451, "Jl. Cikini Raya No. 73, Jakarta Pusat"},
	{-6.1811, 106.8285, "Jl. MH Thamrin No. 1, Jakarta Pusat"},
	{-6.1914, 106.8227, "Jl. Menteng Raya No. 31, Jakarta Pusat"},

	// Jakarta Selatan
	{-6.2297, 106.7986, "Jl. Sudirman No. 123, Jakarta Selatan"},
	{-6.2415, 106.7970, "Jl. Gatot Subroto No. 45, Jakarta Selatan"},
	{-6.2608, 106.8106, "Jl. Metro Duta Kav. UE, Jakarta Selatan"},
	{-6.2489, 106.7931, "Jl. Kebayoran Baru No. 78, Jakarta Selatan"},
	{-6.2350, 106.8074, "Jl. Rasuna Said Kav. C-22, Jakarta Selatan"},

	// Jakarta Utara
	{-6.1247, 106.9130, "Jl. Pluit Raya No. 2, Jakarta Utara"},
	{-6.1356, 106.8954, "Jl. Yos Sudarso No. 1, Jakarta Utara"},
	{-6.1412, 106.8998, "Jl. Tanjung Priok No. 15, Jakarta Utara"},

	// Jakarta Barat
	{-6.1683, 106.7625, "Jl. Letjen S. Parman No. 1, Jakarta Barat"},
	{-6.1819, 106.7791, "Jl. Kebon Jeruk No. 27, Jakarta Barat"},
	{-6.1750, 106.7498, "Jl. Daan Mogot No. 59, Jakarta Barat"},

	// Jakarta Timur
	{-6.2250, 106.9004, "Jl. Raya Bekasi Km. 19, Jakarta Timur"},
	{-6.2425, 106.8825, "Jl. Raya Bogor Km. 24, Jakarta Timur"},
	{-6.2089, 106.9101, "Jl. Pemuda No. 3, Jakarta Timur"},
}

// seedJakartaCheckInLocations seeds additional check-in locations with Jakarta coordinates
// This is called separately to ensure check-in locations are always available
func seedJakartaCheckInLocations() error {
	// Get sales users that are assigned to bricks (prioritize brick-assigned sales)
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
		log.Println("Warning: No users found, skipping Jakarta check-in locations seeding")
		return nil
	}

	// Get accounts
	var accounts []account.Account
	if err := database.DB.Find(&accounts).Error; err != nil {
		return err
	}
	if len(accounts) == 0 {
		log.Println("Warning: No accounts found, skipping Jakarta check-in locations seeding")
		return nil
	}

	// Get contacts
	var contacts []contact.Contact
	if err := database.DB.Find(&contacts).Error; err != nil {
		return err
	}

	// Helper function to marshal location
	marshalLocation := func(lat, lng float64, address string) datatypes.JSON {
		loc := visit_report.Location{
			Latitude:  lat,
			Longitude: lng,
			Address:   address,
		}
		bytes, _ := json.Marshal(loc)
		return bytes
	}

	// Use first sales user (preferably one assigned to a brick)
	userID := users[0].ID
	log.Printf("Using sales user %s for Jakarta check-in locations seeding", userID)

	now := time.Now()

	// Check if Jakarta check-in locations already exist for this user in last 30 days
	var existingCount int64
	last30DaysStart := now.AddDate(0, 0, -30)

	database.DB.Model(&visit_report.VisitReport{}).
		Where("sales_rep_id = ? AND check_in_location IS NOT NULL AND visit_date >= ? AND visit_date <= ?",
			userID, last30DaysStart, now).
		Count(&existingCount)

	// Only add if there are less than 10 check-in locations in last 30 days
	if existingCount >= 10 {
		log.Printf("Jakarta check-in locations already exist for user %s in last 30 days (%d found), skipping...", userID, existingCount)
		return nil
	}

	log.Printf("Found %d existing check-in locations for user %s in last 30 days, will add more...", existingCount, userID)

	visitReports := []visit_report.VisitReport{}

	for i, loc := range jakartaLocations {
		// Spread visits across last 30 days to ensure they're within the default query range
		// Use last 30 days instead of current month to ensure data is always available
		daysAgo := i % 30 // 0-29 days ago
		visitDate := now.AddDate(0, 0, -daysAgo)
		// Ensure visit date is not in the future
		if visitDate.After(now) {
			visitDate = now.AddDate(0, 0, -daysAgo)
		}

		// Vary check-in time throughout the day (8 AM to 5 PM)
		hour := 8 + (i % 10) // 8 AM to 5 PM
		checkInTime := visitDate.Add(time.Duration(hour) * time.Hour)
		// Ensure check-in time is not in the future
		if checkInTime.After(now) {
			checkInTime = visitDate.Add(time.Duration(8) * time.Hour) // Default to 8 AM
		}

		approvedAt := checkInTime.Add(2 * time.Hour) // Approved 2 hours after check-in
		var approvedBy *string
		if len(users) > 0 {
			approvedBy = &users[0].ID
		}

		var contactID *string
		if len(contacts) > 0 {
			contactIndex := i % len(contacts)
			contactID = &contacts[contactIndex].ID
		}

		accountIndex := i % len(accounts)
		accountID := accounts[accountIndex].ID

		visitReports = append(visitReports, visit_report.VisitReport{
			AccountID:       &accountID,
			ContactID:       contactID,
			SalesRepID:      userID,
			VisitDate:       visitDate,
			CheckInTime:     &checkInTime,
			CheckInLocation: marshalLocation(loc.lat, loc.lng, loc.address),
			Purpose:         fmt.Sprintf("Sales visit to %s", loc.address),
			Notes:           fmt.Sprintf("Visit to location in Jakarta. Check-in location recorded at %s", loc.address),
			Status:          "approved",
			ApprovedBy:      approvedBy,
			ApprovedAt:      &approvedAt,
		})
	}

	// Create visit reports
	createdCount := 0
	for _, vr := range visitReports {
		if err := database.DB.Create(&vr).Error; err != nil {
			log.Printf("Warning: Failed to create Jakarta check-in location visit report: %v", err)
			continue // Continue with next visit report instead of failing
		}
		createdCount++
		log.Printf("Created Jakarta check-in location visit report: ID=%s, UserID=%s, VisitDate=%s, Purpose=%s",
			vr.ID, vr.SalesRepID, vr.VisitDate.Format("2006-01-02"), vr.Purpose)
	}

	log.Printf("Jakarta check-in locations seeded successfully (%d/%d visit reports created)", createdCount, len(visitReports))
	return nil
}
