package seeders

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/account"
	"github.com/gilabs/crm-healthcare/api/internal/domain/activity_type"
	"github.com/gilabs/crm-healthcare/api/internal/domain/contact"
	leadqualification "github.com/gilabs/crm-healthcare/api/internal/domain/lead_qualification"
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead_status"
	"github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
	"github.com/gilabs/crm-healthcare/api/internal/domain/product"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"gorm.io/datatypes"
)

const crmSeedSource = "crm-sync-v2"

type seedProductInterest struct {
	ProductID     string `json:"product_id"`
	ProductName   string `json:"product_name"`
	InterestLevel int    `json:"interest_level"`
	Quantity      int    `json:"quantity"`
	Price         int64  `json:"price"`
}

type seedFixtureContext struct {
	SalesUsers      []user.User
	AdminUser       user.User
	Accounts        []account.Account
	Contacts        []contact.Contact
	Products        []product.Product
	LeadStatusIDs   map[string]string
	ActivityTypeIDs map[string]string
	StagesByCode    map[string]pipeline.PipelineStage
}

type crmSeedRecord struct {
	Key                  string
	LeadEmail            string
	FirstName            string
	LastName             string
	CompanyName          string
	Phone                string
	JobTitle             string
	LeadSource           string
	LeadStatusCode       string
	LeadNotes            string
	Address              string
	City                 string
	Province             string
	PostalCode           string
	Website              string
	Industry             string
	Owner                user.User
	Account              *account.Account
	Contact              *contact.Contact
	BudgetAmount         int64
	AuthorityPerson      string
	NeedDescription      string
	NeedPriorityLevel    string
	ExpectedCloseDate    time.Time
	EstimatedValue       int64
	LeadProbability      int
	DealStageCode        string
	DealTitle            string
	DealDescription      string
	DealCloseReason      string
	VisitPurpose         string
	VisitNotes           string
	VisitOutcome         string
	VisitNextSteps       string
	VisitStatus          string
	VisitDate            time.Time
	CheckInTime          *time.Time
	CheckOutTime         *time.Time
	ActivityTypeCode     string
	ActivityDescription  string
	TaskTitle            string
	TaskDescription      string
	TaskType             string
	TaskPriority         string
	TaskSource           string
	Products             []product.Product
	ProductInterests     []seedProductInterest
	NeedProducts         []leadqualification.NeedProduct
	VisitProductMetadata datatypes.JSON
}

// stringPtr creates a pointer to a string
func stringPtr(s string) *string {
	return &s
}

// int64Ptr creates a pointer to an int64
func int64Ptr(i int64) *int64 {
	return &i
}

// float64Ptr creates a pointer to a float64
func float64Ptr(f float64) *float64 {
	return &f
}

func mustJSON(value any) datatypes.JSON {
	bytes, _ := json.Marshal(value)
	return datatypes.JSON(bytes)
}

func loadSeedFixtureContext() (*seedFixtureContext, error) {
	ctx := &seedFixtureContext{
		LeadStatusIDs:   make(map[string]string),
		ActivityTypeIDs: make(map[string]string),
		StagesByCode:    make(map[string]pipeline.PipelineStage),
	}

	if err := database.DB.Where("users.status = ?", "active").
		Joins("JOIN roles ON users.role_id = roles.id").
		Where("roles.code = ?", "sales").
		Order("users.created_at ASC").
		Find(&ctx.SalesUsers).Error; err != nil {
		return nil, err
	}

	if len(ctx.SalesUsers) == 0 {
		if err := database.DB.Where("users.status = ?", "active").
			Joins("JOIN roles ON users.role_id = roles.id").
			Where("roles.code IN (?)", []string{"sales", "sales_manager"}).
			Order("users.created_at ASC").
			Find(&ctx.SalesUsers).Error; err != nil {
			return nil, err
		}
	}

	if len(ctx.SalesUsers) == 0 {
		if err := database.DB.Where("status = ?", "active").Order("created_at ASC").Find(&ctx.SalesUsers).Error; err != nil {
			return nil, err
		}
	}

	if len(ctx.SalesUsers) == 0 {
		return nil, fmt.Errorf("no active users available for CRM seed")
	}

	if err := database.DB.Where("email = ?", "admin@example.com").First(&ctx.AdminUser).Error; err != nil {
		ctx.AdminUser = ctx.SalesUsers[0]
	}

	if err := database.DB.Order("created_at ASC").Find(&ctx.Accounts).Error; err != nil {
		return nil, err
	}
	if err := database.DB.Order("created_at ASC").Find(&ctx.Contacts).Error; err != nil {
		return nil, err
	}
	if err := database.DB.Preload("Category").Where("deleted_at IS NULL").Order("created_at ASC").Find(&ctx.Products).Error; err != nil {
		return nil, err
	}

	var statuses []lead_status.LeadStatus
	if err := database.DB.Where("deleted_at IS NULL").Find(&statuses).Error; err != nil {
		return nil, err
	}
	for _, status := range statuses {
		ctx.LeadStatusIDs[status.Code] = status.ID
	}

	var stages []pipeline.PipelineStage
	if err := database.DB.Where("deleted_at IS NULL").Find(&stages).Error; err != nil {
		return nil, err
	}
	for _, stage := range stages {
		ctx.StagesByCode[stage.Code] = stage
	}

	var activityTypes []activity_type.ActivityType
	if err := database.DB.Where("deleted_at IS NULL").Find(&activityTypes).Error; err != nil {
		return nil, err
	}
	for _, activityType := range activityTypes {
		ctx.ActivityTypeIDs[activityType.Code] = activityType.ID
	}

	return ctx, nil
}

func buildCRMSeedRecords(ctx *seedFixtureContext) []crmSeedRecord {
	now := time.Now()

	type scenario struct {
		key                 string
		leadEmail           string
		firstName           string
		lastName            string
		companyName         string
		phone               string
		jobTitle            string
		leadSource          string
		leadStatusCode      string
		leadNotes           string
		address             string
		city                string
		province            string
		postalCode          string
		website             string
		industry            string
		budgetAmount        int64
		authorityPerson     string
		needDescription     string
		needPriorityLevel   string
		estimatedValue      int64
		leadProbability     int
		dealStageCode       string
		dealTitle           string
		dealDescription     string
		dealCloseReason     string
		visitPurpose        string
		visitNotes          string
		visitOutcome        string
		visitNextSteps      string
		visitStatus         string
		visitOffsetDays     int
		visitDurationHours  int
		activityTypeCode    string
		activityDescription string
		taskTitle           string
		taskDescription     string
		taskType            string
		taskPriority        string
		taskSource          string
		productIndices      []int
		interestLevels      []int
		quantities          []int
	}

	scenarios := []scenario{
		{
			key:                 "lead-new",
			leadEmail:           "budi.santoso@healthcare.id",
			firstName:           "Budi",
			lastName:            "Santoso",
			companyName:         "PT Healthcare Indonesia",
			phone:               "081234567890",
			jobTitle:            "Medical Procurement Lead",
			leadSource:          "website",
			leadStatusCode:      "new",
			leadNotes:           "Seed sinkron CRM untuk discovery awal dan minat produk antihypertensive.",
			address:             "Jl. Sudirman No. 123",
			city:                "Semarang",
			province:            "Jawa Tengah",
			postalCode:          "50125",
			website:             "https://healthcare.id",
			industry:            "Healthcare",
			budgetAmount:        850000000,
			authorityPerson:     "dr. Budi Santoso",
			needDescription:     "Mencari produk antihypertensive dan diabetes untuk pembukaan cabang baru.",
			needPriorityLevel:   "medium",
			estimatedValue:      875000000,
			leadProbability:     25,
			dealStageCode:       "",
			visitPurpose:        "Discovery visit kebutuhan formulary awal",
			visitNotes:          "Mengumpulkan kebutuhan awal, daftar SKU prioritas, dan timeline approval internal.",
			visitOutcome:        "positive",
			visitNextSteps:      "Kirim katalog dan jadwalkan follow-up minggu depan.",
			visitStatus:         "approved",
			visitOffsetDays:     -7,
			visitDurationHours:  2,
			activityTypeCode:    "call",
			activityDescription: "Call awal untuk validasi kebutuhan produk dan PIC approval.",
			taskTitle:           "Follow up katalog kebutuhan formulary PT Healthcare Indonesia",
			taskDescription:     "Kirim katalog dan rangkum kebutuhan produk hasil discovery visit.",
			taskType:            "follow_up",
			taskPriority:        "high",
			taskSource:          "lead_tab",
			productIndices:      []int{0, 1},
			interestLevels:      []int{4, 3},
			quantities:          []int{24, 18},
		},
		{
			key:                 "deal-open",
			leadEmail:           "siti.rahayu@rsud.example.com",
			firstName:           "Siti",
			lastName:            "Rahayu",
			companyName:         "Rumah Sakit Umum Daerah",
			phone:               "081234567891",
			jobTitle:            "Procurement Manager",
			leadSource:          "referral",
			leadStatusCode:      "qualified",
			leadNotes:           "Lead qualified untuk proposal open deal dengan kebutuhan produk kardiovaskular.",
			address:             "Jl. Gatot Subroto No. 456",
			city:                "Semarang",
			province:            "Jawa Tengah",
			postalCode:          "50126",
			website:             "",
			industry:            "Healthcare",
			budgetAmount:        1450000000,
			authorityPerson:     "Siti Rahayu",
			needDescription:     "Membutuhkan paket trial produk cardiovascular untuk pengadaan Q3.",
			needPriorityLevel:   "high",
			estimatedValue:      1525000000,
			leadProbability:     60,
			dealStageCode:       "desire",
			dealTitle:           "Open Pipeline - RSUD Cardiovascular Package",
			dealDescription:     "Proposal paket cardiovascular masih dalam evaluasi komite pengadaan.",
			dealCloseReason:     "",
			visitPurpose:        "Proposal review dan sampling produk cardiovascular",
			visitNotes:          "Review proposal harga, volume estimasi, dan feedback awal dari user unit.",
			visitOutcome:        "very_positive",
			visitNextSteps:      "Revisi proposal final dan tindak lanjuti komite minggu ini.",
			visitStatus:         "approved",
			visitOffsetDays:     -5,
			visitDurationHours:  2,
			activityTypeCode:    "presentation_demo_meet",
			activityDescription: "Presentasi proposal dan simulasi penggunaan produk cardiovascular.",
			taskTitle:           "Kirim revisi proposal cardiovascular RSUD",
			taskDescription:     "Siapkan revisi pricing dan volume berdasarkan hasil review proposal.",
			taskType:            "meeting",
			taskPriority:        "urgent",
			taskSource:          "pipeline_tab",
			productIndices:      []int{1, 2},
			interestLevels:      []int{5, 4},
			quantities:          []int{30, 12},
		},
		{
			key:                 "deal-won",
			leadEmail:           "ahmad.fauzi@kliniksehat.com",
			firstName:           "Ahmad",
			lastName:            "Fauzi",
			companyName:         "Klinik Sehat Jaya",
			phone:               "081234567892",
			jobTitle:            "Owner",
			leadSource:          "cold_call",
			leadStatusCode:      "qualified",
			leadNotes:           "Lead siap closing dan akan dikonversi menjadi closed won dengan item produk sinkron.",
			address:             "Jl. Merdeka No. 789",
			city:                "Semarang",
			province:            "Jawa Tengah",
			postalCode:          "50127",
			website:             "",
			industry:            "Healthcare",
			budgetAmount:        980000000,
			authorityPerson:     "Ahmad Fauzi",
			needDescription:     "Butuh paket diabetes dan vitamin untuk stok operasional 3 bulan.",
			needPriorityLevel:   "critical",
			estimatedValue:      1050000000,
			leadProbability:     90,
			dealStageCode:       "closed_won",
			dealTitle:           "Closed Won - Klinik Sehat Jaya Operational Stock",
			dealDescription:     "Kesepakatan pembelian stok operasional sudah selesai dan siap fulfilment.",
			dealCloseReason:     "Purchase order sudah diterbitkan.",
			visitPurpose:        "Final closing visit dan konfirmasi PO",
			visitNotes:          "Dokumen PO diverifikasi dan jadwal pengiriman disepakati.",
			visitOutcome:        "very_positive",
			visitNextSteps:      "Koordinasi pengiriman dan onboarding produk.",
			visitStatus:         "approved",
			visitOffsetDays:     -3,
			visitDurationHours:  1,
			activityTypeCode:    "document_proposal_sent",
			activityDescription: "PO final dan dokumen pendukung sudah diterima klinik.",
			taskTitle:           "Koordinasi fulfilment closed won Klinik Sehat Jaya",
			taskDescription:     "Pastikan fulfilment produk dan jadwal pengiriman berjalan sesuai PO.",
			taskType:            "follow_up",
			taskPriority:        "high",
			taskSource:          "pipeline_tab",
			productIndices:      []int{2, 3},
			interestLevels:      []int{5, 4},
			quantities:          []int{20, 40},
		},
		{
			key:                 "deal-lost",
			leadEmail:           "maria.wijaya@apotekprima.id",
			firstName:           "Maria",
			lastName:            "Wijaya",
			companyName:         "Apotek Prima Medika",
			phone:               "081234567893",
			jobTitle:            "Owner",
			leadSource:          "event",
			leadStatusCode:      "qualified",
			leadNotes:           "Lead evaluasi akhir namun berakhir lost karena budget ditunda.",
			address:             "Jl. Pemuda No. 18",
			city:                "Semarang",
			province:            "Jawa Tengah",
			postalCode:          "50128",
			website:             "https://apotekprima.id",
			industry:            "Pharmacy",
			budgetAmount:        675000000,
			authorityPerson:     "Maria Wijaya",
			needDescription:     "Membutuhkan paket OTC dan vitamin untuk promo semester kedua.",
			needPriorityLevel:   "medium",
			estimatedValue:      690000000,
			leadProbability:     40,
			dealStageCode:       "closed_lost",
			dealTitle:           "Closed Lost - Apotek Prima Semester Promo",
			dealDescription:     "Peluang hilang karena budget ditunda ke kuartal berikutnya.",
			dealCloseReason:     "Budget promosi ditunda ke kuartal berikutnya.",
			visitPurpose:        "Negosiasi akhir paket promosi semester dua",
			visitNotes:          "Harga disetujui sebagian, tetapi anggaran belum dibuka penuh.",
			visitOutcome:        "neutral",
			visitNextSteps:      "Jaga hubungan dan follow-up ulang saat budget dibuka.",
			visitStatus:         "approved",
			visitOffsetDays:     -2,
			visitDurationHours:  1,
			activityTypeCode:    "follow_up",
			activityDescription: "Follow-up negosiasi akhir, namun customer menunda keputusan.",
			taskTitle:           "Monitor reopening budget Apotek Prima",
			taskDescription:     "Simpan kontak hangat dan follow-up kembali saat budget semester dua dibuka.",
			taskType:            "call",
			taskPriority:        "medium",
			taskSource:          "pipeline_tab",
			productIndices:      []int{0, 3},
			interestLevels:      []int{3, 4},
			quantities:          []int{14, 22},
		},
	}

	records := make([]crmSeedRecord, 0, len(scenarios))
	for index, item := range scenarios {
		owner := ctx.SalesUsers[index%len(ctx.SalesUsers)]

		var accountRef *account.Account
		if len(ctx.Accounts) > 0 {
			accountValue := ctx.Accounts[index%len(ctx.Accounts)]
			accountRef = &accountValue
		}

		var contactRef *contact.Contact
		if accountRef != nil {
			contactRef = findContactForAccount(ctx.Contacts, accountRef.ID, index)
		}
		if contactRef == nil && len(ctx.Contacts) > 0 {
			contactValue := ctx.Contacts[index%len(ctx.Contacts)]
			contactRef = &contactValue
		}

		selectedProducts := pickSeedProducts(ctx.Products, item.productIndices)
		productInterests := make([]seedProductInterest, 0, len(selectedProducts))
		needProducts := make([]leadqualification.NeedProduct, 0, len(selectedProducts))
		needProductNames := make([]string, 0, len(selectedProducts))
		var estimatedValue int64
		for productIndex, productItem := range selectedProducts {
			quantity := 1
			if productIndex < len(item.quantities) && item.quantities[productIndex] > 0 {
				quantity = item.quantities[productIndex]
			}
			interestLevel := 3
			if productIndex < len(item.interestLevels) && item.interestLevels[productIndex] > 0 {
				interestLevel = item.interestLevels[productIndex]
			}

			productInterests = append(productInterests, seedProductInterest{
				ProductID:     productItem.ID,
				ProductName:   productItem.Name,
				InterestLevel: interestLevel,
				Quantity:      quantity,
				Price:         productItem.Price,
			})
			needProducts = append(needProducts, leadqualification.NeedProduct{
				ProductID:   productItem.ID,
				ProductName: productItem.Name,
			})
			needProductNames = append(needProductNames, productItem.Name)
			estimatedValue += productItem.Price * int64(quantity)
		}
		if estimatedValue == 0 {
			estimatedValue = item.estimatedValue
		}

		baseVisitDate := now.AddDate(0, 0, item.visitOffsetDays)
		expectedCloseDate := now.AddDate(0, 0, 14+index*5)
		if item.dealStageCode == "closed_won" || item.dealStageCode == "closed_lost" {
			expectedCloseDate = baseVisitDate.AddDate(0, 0, 2)
		}
		checkInTime := baseVisitDate
		checkInTime = time.Date(
			checkInTime.Year(),
			checkInTime.Month(),
			checkInTime.Day(),
			9+index,
			0,
			0,
			0,
			checkInTime.Location(),
		)
		checkOutTime := checkInTime.Add(time.Duration(item.visitDurationHours) * time.Hour)
		visitDate := time.Date(checkInTime.Year(), checkInTime.Month(), checkInTime.Day(), 0, 0, 0, 0, checkInTime.Location())

		records = append(records, crmSeedRecord{
			Key:               item.key,
			LeadEmail:         item.leadEmail,
			FirstName:         item.firstName,
			LastName:          item.lastName,
			CompanyName:       item.companyName,
			Phone:             item.phone,
			JobTitle:          item.jobTitle,
			LeadSource:        item.leadSource,
			LeadStatusCode:    item.leadStatusCode,
			LeadNotes:         item.leadNotes,
			Address:           item.address,
			City:              item.city,
			Province:          item.province,
			PostalCode:        item.postalCode,
			Website:           item.website,
			Industry:          item.industry,
			Owner:             owner,
			Account:           accountRef,
			Contact:           contactRef,
			BudgetAmount:      item.budgetAmount,
			AuthorityPerson:   item.authorityPerson,
			NeedDescription:   item.needDescription,
			NeedPriorityLevel: item.needPriorityLevel,
			ExpectedCloseDate: expectedCloseDate,
			EstimatedValue:    estimatedValue,
			LeadProbability:   item.leadProbability,
			DealStageCode:     item.dealStageCode,
			DealTitle:         item.dealTitle,
			DealDescription:   item.dealDescription,
			DealCloseReason:   item.dealCloseReason,
			VisitPurpose:      item.visitPurpose,
			VisitNotes:        item.visitNotes,
			VisitOutcome:      item.visitOutcome,
			VisitNextSteps:    item.visitNextSteps,
			VisitStatus:       item.visitStatus,
			VisitDate:         visitDate,
			CheckInTime:       &checkInTime,
			CheckOutTime:      &checkOutTime,
			ActivityTypeCode:  item.activityTypeCode,
			ActivityDescription: fmt.Sprintf(
				"%s [%s]",
				item.activityDescription,
				crmSeedSource,
			),
			TaskTitle:        item.taskTitle,
			TaskDescription:  item.taskDescription,
			TaskType:         item.taskType,
			TaskPriority:     item.taskPriority,
			TaskSource:       item.taskSource,
			Products:         selectedProducts,
			ProductInterests: productInterests,
			NeedProducts:     needProducts,
			VisitProductMetadata: mustJSON(map[string]any{
				"seed_source":       crmSeedSource,
				"product_interests": productInterests,
			}),
		})
	}

	return records
}

func pickSeedProducts(products []product.Product, indexes []int) []product.Product {
	if len(products) == 0 {
		return nil
	}

	selected := make([]product.Product, 0, len(indexes))
	seen := make(map[string]struct{})
	for _, index := range indexes {
		productItem := products[index%len(products)]
		if _, exists := seen[productItem.ID]; exists {
			continue
		}
		seen[productItem.ID] = struct{}{}
		selected = append(selected, productItem)
	}
	return selected
}

func findContactForAccount(contacts []contact.Contact, accountID string, fallbackIndex int) *contact.Contact {
	for _, contactItem := range contacts {
		if contactItem.AccountID == accountID {
			value := contactItem
			return &value
		}
	}
	if len(contacts) == 0 {
		return nil
	}
	value := contacts[fallbackIndex%len(contacts)]
	return &value
}
