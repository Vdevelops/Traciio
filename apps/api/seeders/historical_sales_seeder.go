package seeders

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	accountdomain "github.com/gilabs/crm-healthcare/api/internal/domain/account"
	activitydomain "github.com/gilabs/crm-healthcare/api/internal/domain/activity"
	"github.com/gilabs/crm-healthcare/api/internal/domain/activity_type"
	"github.com/gilabs/crm-healthcare/api/internal/domain/contact"
	"github.com/gilabs/crm-healthcare/api/internal/domain/customer_purchase"
	leaddomain "github.com/gilabs/crm-healthcare/api/internal/domain/lead"
	leadqualification "github.com/gilabs/crm-healthcare/api/internal/domain/lead_qualification"
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead_status"
	monthlytargetdomain "github.com/gilabs/crm-healthcare/api/internal/domain/monthly_target"
	"github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
	"github.com/gilabs/crm-healthcare/api/internal/domain/product"
	"github.com/gilabs/crm-healthcare/api/internal/domain/task"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"github.com/gilabs/crm-healthcare/api/internal/domain/visit_report"
	"gorm.io/gorm"
)

const historicalSeedSource = "historical-sales-2023"

// SeedHistoricalSalesData creates a deterministic, idempotent sales history from
// January 2023 to the current month. It keeps revenue, product analytics, visits,
// tasks, and targets aligned with the same sales/deal graph.
func SeedHistoricalSalesData() error {
	db := database.DB
	now := time.Now()
	start := time.Date(2023, time.January, 1, 0, 0, 0, 0, now.Location())

	var salesUsers []user.User
	if err := db.Where("users.status = ?", "active").
		Joins("JOIN roles ON roles.id = users.role_id AND roles.deleted_at IS NULL").
		Where("roles.code = ?", "sales").
		Order("users.created_at ASC").
		Find(&salesUsers).Error; err != nil {
		return err
	}
	if len(salesUsers) == 0 {
		log.Println("Warning: no active sales users found, skipping historical sales seed")
		return nil
	}

	var accounts []accountdomain.Account
	if err := db.Where("deleted_at IS NULL").Order("created_at ASC").Find(&accounts).Error; err != nil {
		return err
	}
	if len(accounts) == 0 {
		log.Println("Warning: no accounts found, skipping historical sales seed")
		return nil
	}

	var contacts []contact.Contact
	if err := db.Where("deleted_at IS NULL").Order("created_at ASC").Find(&contacts).Error; err != nil {
		return err
	}
	contactsByAccount := map[string][]contact.Contact{}
	for _, c := range contacts {
		contactsByAccount[c.AccountID] = append(contactsByAccount[c.AccountID], c)
	}

	var products []product.Product
	if err := db.Preload("Category").
		Where("deleted_at IS NULL").
		Order("created_at ASC").
		Find(&products).Error; err != nil {
		return err
	}
	if len(products) == 0 {
		log.Println("Warning: no products found, skipping historical sales seed")
		return nil
	}

	leadStatusIDs, err := loadHistoricalLeadStatusIDs(db)
	if err != nil {
		return err
	}
	activityTypeIDs, err := loadHistoricalActivityTypeIDs(db)
	if err != nil {
		return err
	}

	stages := map[string]pipeline.PipelineStage{}
	var stageRows []pipeline.PipelineStage
	if err := db.Where("deleted_at IS NULL").Find(&stageRows).Error; err != nil {
		return err
	}
	for _, stage := range stageRows {
		stages[stage.Code] = stage
	}
	for _, code := range []string{"closed_won", "closed_lost", "negotiation"} {
		if _, ok := stages[code]; !ok {
			return fmt.Errorf("missing pipeline stage %s for historical sales seed", code)
		}
	}

	adminID := salesUsers[0].ID
	var admin user.User
	if err := db.Where("email = ?", "admin@example.com").First(&admin).Error; err == nil {
		adminID = admin.ID
	}

	rng := rand.New(rand.NewSource(20240727))
	leadsSeeded := 0
	activitiesSeeded := 0
	dealsSeeded := 0
	visitsSeeded := 0
	tasksSeeded := 0

	for cursor := start; !cursor.After(now); cursor = cursor.AddDate(0, 1, 0) {
		for userIndex, salesUser := range salesUsers {
			dealsThisMonth := 3 + ((cursor.Year() + int(cursor.Month()) + userIndex) % 2)
			for sequence := 0; sequence < dealsThisMonth; sequence++ {
				seedDate := historicalSeedDate(cursor, now, userIndex, sequence)
				accountEntity := accounts[(userIndex+sequence+int(cursor.Month()))%len(accounts)]
				var contactID *string
				if accountContacts := contactsByAccount[accountEntity.ID]; len(accountContacts) > 0 {
					selected := accountContacts[(sequence+userIndex)%len(accountContacts)]
					contactID = stringPtr(selected.ID)
				}

				stageCode := "closed_won"
				status := "won"
				closeReason := "repeat order after product fit validation"
				if (cursor.Year()+int(cursor.Month())+sequence+userIndex)%7 == 0 {
					stageCode = "closed_lost"
					status = "lost"
					closeReason = "budget delayed by customer"
				} else if cursor.Year() == now.Year() && cursor.Month() == now.Month() && sequence == dealsThisMonth-1 {
					stageCode = "negotiation"
					status = "open"
					closeReason = ""
				}

				stage := stages[stageCode]
				dealTitle := fmt.Sprintf("Historical Sales %04d-%02d %s #%02d", cursor.Year(), cursor.Month(), salesUser.Name, sequence+1)
				leadEntity, err := upsertHistoricalLead(db, cursor, seedDate, accountEntity, contactID, salesUser, adminID, leadStatusIDs, "converted", products, rng, sequence, true)
				if err != nil {
					return err
				}
				leadsSeeded++
				if err := upsertHistoricalLeadActivity(db, leadEntity, accountEntity, contactID, salesUser, activityTypeIDs, seedDate, sequence); err != nil {
					return err
				}
				activitiesSeeded++

				deal, items, err := upsertHistoricalDeal(db, dealTitle, accountEntity, contactID, &leadEntity, salesUser, adminID, stage, status, closeReason, seedDate, products, rng, sequence)
				if err != nil {
					return err
				}
				dealsSeeded++

				if err := markHistoricalLeadConverted(db, leadEntity, deal, adminID, seedDate); err != nil {
					return err
				}
				if err := upsertHistoricalOpportunityActivity(db, leadEntity, deal, accountEntity, contactID, salesUser, activityTypeIDs, seedDate, sequence); err != nil {
					return err
				}
				activitiesSeeded++

				if err := upsertHistoricalVisit(db, deal, accountEntity, contactID, &leadEntity, salesUser, adminID, seedDate, sequence); err != nil {
					return err
				}
				visitsSeeded++

				if err := upsertHistoricalTask(db, deal, accountEntity, contactID, &leadEntity, salesUser, adminID, seedDate, sequence); err != nil {
					return err
				}
				tasksSeeded++

				if status == "won" {
					if err := upsertHistoricalPurchase(db, deal, salesUser, items); err != nil {
						return err
					}
				} else if deal.ID != "" {
					if err := db.Where("deal_id = ?", deal.ID).Delete(&customer_purchase.CustomerPurchaseHistory{}).Error; err != nil {
						return err
					}
				}
			}

			standaloneLeads := 2 + ((cursor.Year() + int(cursor.Month()) + userIndex) % 2)
			for sequence := 0; sequence < standaloneLeads; sequence++ {
				seedDate := historicalSeedDate(cursor, now, userIndex, sequence+dealsThisMonth).Add(2 * time.Hour)
				accountEntity := accounts[(userIndex+sequence+dealsThisMonth+int(cursor.Month()))%len(accounts)]
				var contactID *string
				if accountContacts := contactsByAccount[accountEntity.ID]; len(accountContacts) > 0 {
					selected := accountContacts[(sequence+dealsThisMonth+userIndex)%len(accountContacts)]
					contactID = stringPtr(selected.ID)
				}

				statusCode := historicalStandaloneLeadStatus(cursor, now, userIndex, sequence)
				leadEntity, err := upsertHistoricalLead(db, cursor, seedDate, accountEntity, contactID, salesUser, adminID, leadStatusIDs, statusCode, products, rng, sequence+dealsThisMonth, false)
				if err != nil {
					return err
				}
				leadsSeeded++
				if err := upsertHistoricalLeadActivity(db, leadEntity, accountEntity, contactID, salesUser, activityTypeIDs, seedDate, sequence+dealsThisMonth); err != nil {
					return err
				}
				activitiesSeeded++
			}
		}
	}

	if err := seedHistoricalMonthlyTargets(start, now, salesUsers); err != nil {
		return err
	}

	log.Printf("Seeded historical sales data from %s to %s: %d leads, %d deals, %d activities, %d visits, %d tasks", start.Format("2006-01"), now.Format("2006-01"), leadsSeeded, dealsSeeded, activitiesSeeded, visitsSeeded, tasksSeeded)
	return nil
}

func historicalSeedDate(monthStart time.Time, now time.Time, userIndex int, sequence int) time.Time {
	lastDay := time.Date(monthStart.Year(), monthStart.Month()+1, 0, 0, 0, 0, 0, monthStart.Location()).Day()
	if monthStart.Year() == now.Year() && monthStart.Month() == now.Month() && now.Day() < lastDay {
		lastDay = now.Day()
	}
	if lastDay < 1 {
		lastDay = 1
	}
	day := 4 + ((sequence*6 + userIndex*3) % lastDay)
	if day > lastDay {
		day = lastDay
	}
	return time.Date(monthStart.Year(), monthStart.Month(), day, 10+sequence%6, 15, 0, 0, monthStart.Location())
}

func loadHistoricalLeadStatusIDs(db *gorm.DB) (map[string]string, error) {
	var statuses []lead_status.LeadStatus
	if err := db.Where("deleted_at IS NULL").Find(&statuses).Error; err != nil {
		return nil, err
	}
	statusIDs := map[string]string{}
	for _, status := range statuses {
		statusIDs[status.Code] = status.ID
	}
	for _, code := range []string{"new", "contacted", "interested", "qualified", "proposal_sent", "converted", "lost"} {
		if _, ok := statusIDs[code]; !ok {
			return nil, fmt.Errorf("missing lead status %s for historical sales seed", code)
		}
	}
	return statusIDs, nil
}

func loadHistoricalActivityTypeIDs(db *gorm.DB) (map[string]string, error) {
	var types []activity_type.ActivityType
	if err := db.Where("status = ?", "active").Find(&types).Error; err != nil {
		return nil, err
	}
	typeIDs := map[string]string{}
	for _, item := range types {
		typeIDs[item.Code] = item.ID
	}
	for _, code := range []string{"call", "whatsapp_chat", "email", "note", "follow_up", "presentation_demo_meet", "document_proposal_sent"} {
		if _, ok := typeIDs[code]; !ok {
			return nil, fmt.Errorf("missing activity type %s for historical sales seed", code)
		}
	}
	return typeIDs, nil
}

func historicalStandaloneLeadStatus(monthStart time.Time, now time.Time, userIndex int, sequence int) string {
	if monthStart.Year() == now.Year() && monthStart.Month() == now.Month() {
		return []string{"new", "contacted", "interested"}[(userIndex+sequence)%3]
	}
	return []string{"new", "contacted", "interested", "qualified", "proposal_sent", "lost"}[(monthStart.Year()+int(monthStart.Month())+userIndex+sequence)%6]
}

func upsertHistoricalLead(
	db *gorm.DB,
	monthStart time.Time,
	seedDate time.Time,
	accountEntity accountdomain.Account,
	contactID *string,
	salesUser user.User,
	adminID string,
	leadStatusIDs map[string]string,
	statusCode string,
	products []product.Product,
	rng *rand.Rand,
	sequence int,
	willConvert bool,
) (leaddomain.Lead, error) {
	firstNames := []string{"Andi", "Rina", "Dewi", "Fajar", "Sari", "Hendra", "Maya", "Rizky", "Nadia", "Agus", "Lestari", "Dimas"}
	lastNames := []string{"Pratama", "Wijaya", "Utami", "Saputra", "Permata", "Nugroho", "Siregar", "Halim", "Susanto", "Kurniawan", "Mahendra", "Wibowo"}
	leadSources := []string{"website", "referral", "event", "cold_call", "field_visit", "partner"}
	jobTitles := []string{"Procurement Manager", "Head Pharmacist", "Clinic Owner", "Medical Director", "Purchasing Staff", "Hospital Administrator"}

	nameIndex := (monthStart.Year() + int(monthStart.Month()) + sequence + rng.Intn(len(firstNames))) % len(firstNames)
	firstName := firstNames[nameIndex]
	lastName := lastNames[(nameIndex+sequence)%len(lastNames)]
	leadKind := "pipeline"
	if willConvert {
		leadKind = "converted"
	}
	email := fmt.Sprintf("hist.%s.%s.%04d%02d.%s.%02d@seed.traciio.local",
		strings.ToLower(firstName),
		strings.ToLower(lastName),
		monthStart.Year(),
		monthStart.Month(),
		salesUser.ID[:8],
		sequence+1,
	)

	var entity leaddomain.Lead
	err := db.Where("email = ?", email).First(&entity).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return entity, err
	}

	scoreByStatus := map[string]int{
		"new":           10,
		"contacted":     25,
		"interested":    45,
		"qualified":     70,
		"proposal_sent": 85,
		"converted":     100,
		"lost":          0,
	}
	estimatedValue := historicalEstimatedValue(products, rng, sequence)
	expectedClose := seedDate.AddDate(0, 0, 21+(sequence%20))
	budgetConfirmed := scoreByStatus[statusCode] >= 45
	authorityConfirmed := scoreByStatus[statusCode] >= 70 || statusCode == "converted"
	needConfirmed := scoreByStatus[statusCode] >= 25
	timelineConfirmed := scoreByStatus[statusCode] >= 70 || statusCode == "converted"

	entity.FirstName = firstName
	entity.LastName = lastName
	entity.CompanyName = accountEntity.Name
	entity.Email = email
	entity.Phone = fmt.Sprintf("08%010d", 2100000000+monthStart.Year()+int(monthStart.Month())*100+sequence*13)
	entity.JobTitle = jobTitles[(sequence+nameIndex)%len(jobTitles)]
	entity.Industry = "Healthcare"
	entity.LeadSource = leadSources[(sequence+historicalUserShard(salesUser.ID))%len(leadSources)]
	entity.LeadStatus = statusCode
	entity.LeadStatusID = stringPtr(leadStatusIDs[statusCode])
	entity.LeadScore = scoreByStatus[statusCode]
	entity.Probability = scoreByStatus[statusCode]
	entity.EstimatedValue = estimatedValue
	entity.BudgetConfirmed = budgetConfirmed
	if budgetConfirmed {
		entity.BudgetAmount = int64Ptr(estimatedValue)
	} else {
		entity.BudgetAmount = nil
	}
	entity.AuthorityConfirmed = authorityConfirmed
	if authorityConfirmed {
		entity.AuthorityPerson = firstName + " " + lastName
	} else {
		entity.AuthorityPerson = ""
	}
	entity.NeedConfirmed = needConfirmed
	entity.NeedDescription = historicalNeedDescription(products, sequence)
	entity.TimelineConfirmed = timelineConfirmed
	if timelineConfirmed {
		entity.ExpectedCloseDate = &expectedClose
	} else {
		entity.ExpectedCloseDate = nil
	}
	entity.AssignedTo = stringPtr(salesUser.ID)
	entity.AccountID = stringPtr(accountEntity.ID)
	entity.ContactID = contactID
	entity.Notes = fmt.Sprintf("Historical %s lead generated for CRM funnel coverage [%s]", leadKind, historicalSeedSource)
	entity.Address = accountEntity.Address
	entity.City = accountEntity.City
	entity.Province = accountEntity.Province
	entity.PostalCode = accountEntity.PostalCode
	entity.Country = "Indonesia"
	entity.Website = accountEntity.Website
	entity.CreatedBy = adminID
	entity.CreatedAt = seedDate.AddDate(0, 0, -14)
	entity.UpdatedAt = seedDate
	if !willConvert {
		entity.OpportunityID = nil
		entity.ConvertedPipelineID = nil
		entity.ConvertedAt = nil
		entity.ConvertedBy = nil
		entity.ConversionMetadata = mustJSON(map[string]any{
			"seed_source": historicalSeedSource,
			"status":      statusCode,
		})
	}

	if entity.ID == "" {
		if err := db.Create(&entity).Error; err != nil {
			return entity, err
		}
	} else if err := db.Save(&entity).Error; err != nil {
		return entity, err
	}

	if err := upsertHistoricalQualification(db, entity, statusCode, products, estimatedValue, expectedClose, sequence); err != nil {
		return entity, err
	}

	return entity, nil
}

func markHistoricalLeadConverted(db *gorm.DB, entity leaddomain.Lead, deal pipeline.Deal, adminID string, convertedAt time.Time) error {
	entity.LeadStatus = "converted"
	entity.LeadScore = 100
	entity.Probability = 100
	entity.AccountID = stringPtr(deal.AccountID)
	entity.OpportunityID = stringPtr(deal.ID)
	entity.ConvertedPipelineID = stringPtr(deal.ID)
	entity.ConvertedAt = &convertedAt
	entity.ConvertedBy = stringPtr(adminID)
	entity.ConversionMetadata = mustJSON(map[string]any{
		"seed_source": historicalSeedSource,
		"deal_id":     deal.ID,
		"deal_value":  deal.Value,
	})

	var status lead_status.LeadStatus
	if err := db.Where("code = ?", "converted").First(&status).Error; err == nil {
		entity.LeadStatusID = stringPtr(status.ID)
	}
	return db.Save(&entity).Error
}

func upsertHistoricalQualification(db *gorm.DB, entity leaddomain.Lead, statusCode string, products []product.Product, estimatedValue int64, expectedClose time.Time, sequence int) error {
	var qualification leadqualification.LeadQualificationChecklist
	err := db.Where("lead_id = ?", entity.ID).First(&qualification).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	score := entity.LeadScore
	qualification.LeadID = entity.ID
	qualification.BudgetTargetAmount = estimatedValue
	qualification.BudgetTargetCurrency = "IDR"
	qualification.BudgetConfirmed = score >= 45
	qualification.BudgetNotes = "Historical seed: budget follows lead maturity."
	qualification.AuthorityTargetPerson = entity.AuthorityPerson
	qualification.AuthorityTargetRole = entity.JobTitle
	qualification.AuthorityConfirmed = score >= 70 || statusCode == "converted"
	qualification.AuthorityNotes = "Historical seed: authority identified for qualified leads."
	qualification.NeedTargetProducts = mustJSON(historicalNeedProducts(products, sequence))
	qualification.NeedPriorityLevel = []string{"medium", "high", "critical"}[sequence%3]
	qualification.NeedConfirmed = score >= 25
	qualification.NeedNotes = entity.NeedDescription
	qualification.TimelineTargetDate = &expectedClose
	qualification.TimelineFlexibility = []string{"flexible", "fixed", "urgent"}[sequence%3]
	qualification.TimelineConfirmed = score >= 70 || statusCode == "converted"
	qualification.TimelineNotes = "Historical seed: timeline follows funnel stage."

	if qualification.ID == "" {
		return db.Create(&qualification).Error
	}
	return db.Save(&qualification).Error
}

func historicalEstimatedValue(products []product.Product, rng *rand.Rand, sequence int) int64 {
	items := buildHistoricalDealItems(products, rng, sequence)
	total := int64(0)
	for _, item := range items {
		total += item.Subtotal
	}
	return total
}

func historicalNeedProducts(products []product.Product, sequence int) []leadqualification.NeedProduct {
	if len(products) == 0 {
		return nil
	}
	count := 1 + sequence%2
	result := make([]leadqualification.NeedProduct, 0, count)
	for i := 0; i < count; i++ {
		item := products[(sequence+i)%len(products)]
		result = append(result, leadqualification.NeedProduct{
			ProductID:   item.ID,
			ProductName: item.Name,
		})
	}
	return result
}

func historicalNeedDescription(products []product.Product, sequence int) string {
	needs := historicalNeedProducts(products, sequence)
	if len(needs) == 0 {
		return "General product inquiry for healthcare procurement."
	}
	names := make([]string, 0, len(needs))
	for _, need := range needs {
		names = append(names, need.ProductName)
	}
	return "Customer is evaluating " + strings.Join(names, ", ") + " for recurring healthcare procurement."
}

func historicalUserShard(userID string) int {
	total := 0
	for _, r := range userID {
		total += int(r)
	}
	return total
}

func upsertHistoricalLeadActivity(
	db *gorm.DB,
	leadEntity leaddomain.Lead,
	accountEntity accountdomain.Account,
	contactID *string,
	salesUser user.User,
	activityTypeIDs map[string]string,
	seedDate time.Time,
	sequence int,
) error {
	activityCode := []string{"call", "whatsapp_chat", "email", "note"}[sequence%4]
	description := fmt.Sprintf("Historical lead %s activity for %s [%s]", activityCode, leadEntity.Email, historicalSeedSource)

	var entity activitydomain.Activity
	err := db.Where("lead_id = ? AND deal_id IS NULL AND description = ?", leadEntity.ID, description).First(&entity).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	entity.Type = activityCode
	entity.ActivityTypeID = stringPtr(activityTypeIDs[activityCode])
	entity.AccountID = stringPtr(accountEntity.ID)
	entity.ContactID = contactID
	entity.DealID = nil
	entity.LeadID = stringPtr(leadEntity.ID)
	entity.UserID = salesUser.ID
	entity.Description = description
	entity.Timestamp = seedDate.Add(-48 * time.Hour)
	entity.Metadata = mustJSON(map[string]any{
		"seed_source": historicalSeedSource,
		"entity":      "lead",
		"lead_status": leadEntity.LeadStatus,
		"account_id":  accountEntity.ID,
	})
	entity.CreatedAt = entity.Timestamp
	entity.UpdatedAt = seedDate

	if entity.ID == "" {
		return db.Create(&entity).Error
	}
	return db.Save(&entity).Error
}

func upsertHistoricalOpportunityActivity(
	db *gorm.DB,
	leadEntity leaddomain.Lead,
	deal pipeline.Deal,
	accountEntity accountdomain.Account,
	contactID *string,
	salesUser user.User,
	activityTypeIDs map[string]string,
	seedDate time.Time,
	sequence int,
) error {
	activityCode := "presentation_demo_meet"
	if deal.Status == "won" {
		activityCode = "document_proposal_sent"
	} else if deal.Status == "lost" {
		activityCode = "follow_up"
	}
	description := fmt.Sprintf("Historical opportunity %s activity for %s [%s]", activityCode, deal.Title, historicalSeedSource)

	var entity activitydomain.Activity
	err := db.Where("deal_id = ? AND description = ?", deal.ID, description).First(&entity).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	entity.Type = activityCode
	entity.ActivityTypeID = stringPtr(activityTypeIDs[activityCode])
	entity.AccountID = stringPtr(accountEntity.ID)
	entity.ContactID = contactID
	entity.DealID = stringPtr(deal.ID)
	entity.LeadID = stringPtr(leadEntity.ID)
	entity.UserID = salesUser.ID
	entity.Description = description
	entity.Timestamp = seedDate.Add(-6 * time.Hour).Add(time.Duration(sequence) * time.Minute)
	entity.Metadata = mustJSON(map[string]any{
		"seed_source": historicalSeedSource,
		"entity":      "opportunity",
		"deal_status": deal.Status,
		"deal_value":  deal.Value,
		"lead_id":     leadEntity.ID,
	})
	entity.CreatedAt = entity.Timestamp
	entity.UpdatedAt = seedDate

	if entity.ID == "" {
		return db.Create(&entity).Error
	}
	return db.Save(&entity).Error
}

func upsertHistoricalDeal(
	db *gorm.DB,
	title string,
	accountEntity accountdomain.Account,
	contactID *string,
	leadEntity *leaddomain.Lead,
	salesUser user.User,
	adminID string,
	stage pipeline.PipelineStage,
	status string,
	closeReason string,
	seedDate time.Time,
	products []product.Product,
	rng *rand.Rand,
	sequence int,
) (pipeline.Deal, []pipeline.DealProductItem, error) {
	var deal pipeline.Deal
	err := db.Where("title = ? AND notes LIKE ?", title, "%"+historicalSeedSource+"%").First(&deal).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return deal, nil, err
	}

	actualClose := seedDate
	expectedClose := seedDate.AddDate(0, 0, 7)
	if status == "open" {
		actualClose = time.Time{}
		expectedClose = seedDate.AddDate(0, 0, 21)
	}

	items := buildHistoricalDealItems(products, rng, sequence)
	value := int64(0)
	for _, item := range items {
		value += item.Subtotal
	}

	deal.Title = title
	deal.Description = "Historical seeded sales opportunity for analytics coverage from 2023."
	deal.AccountID = accountEntity.ID
	deal.ContactID = contactID
	if leadEntity != nil {
		deal.LeadID = stringPtr(leadEntity.ID)
	}
	deal.StageID = stage.ID
	deal.Value = value
	deal.Probability = stage.Probability
	deal.ExpectedCloseDate = &expectedClose
	deal.AssignedTo = stringPtr(salesUser.ID)
	deal.BrickID = salesUser.BrickID
	if deal.BrickID == nil {
		deal.BrickID = accountEntity.BrickID
	}
	deal.Status = status
	deal.Source = "historical_seed"
	deal.BudgetConfirmed = true
	deal.AuthorityConfirmed = true
	deal.NeedConfirmed = true
	deal.TimelineConfirmed = true
	deal.CloseReason = closeReason
	deal.Notes = "Generated historical sales fixture [" + historicalSeedSource + "]"
	deal.CreatedBy = adminID
	deal.CreatedAt = seedDate.AddDate(0, 0, -10)
	deal.UpdatedAt = seedDate
	if status == "won" || status == "lost" {
		deal.ActualCloseDate = &actualClose
	} else {
		deal.ActualCloseDate = nil
	}

	if deal.ID == "" {
		if err := db.Create(&deal).Error; err != nil {
			return deal, nil, err
		}
	} else {
		if err := db.Save(&deal).Error; err != nil {
			return deal, nil, err
		}
		if err := db.Where("deal_id = ?", deal.ID).Delete(&pipeline.DealProductItem{}).Error; err != nil {
			return deal, nil, err
		}
	}

	for i := range items {
		items[i].DealID = deal.ID
		items[i].CreatedAt = seedDate
		items[i].UpdatedAt = seedDate
	}
	if len(items) > 0 {
		if err := db.Create(&items).Error; err != nil {
			return deal, nil, err
		}
	}

	return deal, items, nil
}

func buildHistoricalDealItems(products []product.Product, rng *rand.Rand, sequence int) []pipeline.DealProductItem {
	itemCount := 1 + (sequence % 2)
	items := make([]pipeline.DealProductItem, 0, itemCount)
	used := map[string]bool{}
	for i := 0; i < itemCount; i++ {
		productEntity := products[(sequence+i+rng.Intn(len(products)))%len(products)]
		if used[productEntity.ID] {
			continue
		}
		used[productEntity.ID] = true

		quantity := 8 + rng.Intn(28)
		if productEntity.Price >= 5000000 {
			quantity = 1 + rng.Intn(5)
		}

		categoryID := stringPtr(productEntity.CategoryID)
		categoryName := ""
		if productEntity.Category != nil {
			categoryName = productEntity.Category.Name
		}

		subtotal := productEntity.Price * int64(quantity)
		items = append(items, pipeline.DealProductItem{
			ProductID:           productEntity.ID,
			ProductName:         productEntity.Name,
			ProductSKU:          productEntity.SKU,
			UnitPrice:           productEntity.Price,
			UnitCost:            productEntity.Cost,
			Quantity:            quantity,
			DiscountAmount:      0,
			Subtotal:            subtotal,
			ProductCategoryID:   categoryID,
			ProductCategoryName: categoryName,
			Notes:               "Historical sales seed [" + historicalSeedSource + "]",
		})
	}
	return items
}

func upsertHistoricalVisit(
	db *gorm.DB,
	deal pipeline.Deal,
	accountEntity accountdomain.Account,
	contactID *string,
	leadEntity *leaddomain.Lead,
	salesUser user.User,
	adminID string,
	seedDate time.Time,
	sequence int,
) error {
	purpose := fmt.Sprintf("Historical follow-up visit %s", deal.Title)
	var visit visit_report.VisitReport
	err := db.Where("deal_id = ? AND purpose = ?", deal.ID, purpose).First(&visit).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	checkIn := seedDate.Add(-90 * time.Minute)
	checkOut := seedDate.Add(-25 * time.Minute)
	visit.AccountID = stringPtr(accountEntity.ID)
	visit.ContactID = contactID
	visit.DealID = stringPtr(deal.ID)
	if leadEntity != nil {
		visit.LeadID = stringPtr(leadEntity.ID)
	} else {
		visit.LeadID = nil
	}
	visit.SalesRepID = salesUser.ID
	visit.BrickID = deal.BrickID
	visit.VisitDate = seedDate.Add(-2 * time.Hour)
	visit.CheckInTime = &checkIn
	visit.CheckOutTime = &checkOut
	visit.CheckInLocation = mustJSON(visit_report.Location{
		Latitude:  -6.98 + float64(sequence)*0.003,
		Longitude: 110.41 + float64(sequence)*0.003,
		Address:   accountEntity.Address,
	})
	visit.CheckOutLocation = visit.CheckInLocation
	visit.Purpose = purpose
	visit.Notes = "Historical customer visit generated for sales trend analytics [" + historicalSeedSource + "]"
	visit.Outcome = "positive"
	visit.NextSteps = "Continue procurement follow-up and product availability check."
	visit.Metadata = mustJSON(map[string]any{
		"seed_source": historicalSeedSource,
		"deal_id":     deal.ID,
	})
	visit.Status = "completed"
	if sequence%3 == 0 {
		visit.Status = "approved"
		visit.ApprovedBy = stringPtr(adminID)
		approvedAt := checkOut.Add(2 * time.Hour)
		visit.ApprovedAt = &approvedAt
	} else {
		visit.ApprovedBy = nil
		visit.ApprovedAt = nil
	}
	visit.CreatedAt = seedDate.Add(-3 * time.Hour)
	visit.UpdatedAt = seedDate

	if visit.ID == "" {
		return db.Create(&visit).Error
	}
	return db.Save(&visit).Error
}

func upsertHistoricalTask(
	db *gorm.DB,
	deal pipeline.Deal,
	accountEntity accountdomain.Account,
	contactID *string,
	leadEntity *leaddomain.Lead,
	salesUser user.User,
	adminID string,
	seedDate time.Time,
	sequence int,
) error {
	title := fmt.Sprintf("Historical follow-up task %s", deal.Title)
	var entity task.Task
	err := db.Where("deal_id = ? AND title = ?", deal.ID, title).First(&entity).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	dueDate := seedDate.AddDate(0, 0, 2)
	completedAt := seedDate.Add(4 * time.Hour)
	entity.Title = title
	entity.Description = "Historical completed task generated for sales performance analytics [" + historicalSeedSource + "]"
	entity.Type = "follow_up"
	entity.Status = "completed"
	entity.Priority = []string{"medium", "high", "urgent"}[sequence%3]
	entity.DueDate = &dueDate
	entity.CompletedAt = &completedAt
	entity.AssignedTo = stringPtr(salesUser.ID)
	entity.AssignedFrom = stringPtr(adminID)
	entity.AccountID = stringPtr(accountEntity.ID)
	entity.ContactID = contactID
	entity.DealID = stringPtr(deal.ID)
	if leadEntity != nil {
		entity.LeadID = stringPtr(leadEntity.ID)
	} else {
		entity.LeadID = nil
	}
	entity.TaskSource = "historical_seed"
	entity.QuickActionType = "sales_follow_up"
	entity.QuickActionPayload = mustJSON(map[string]any{
		"seed_source": historicalSeedSource,
		"deal_id":     deal.ID,
	})
	entity.CreatedBy = adminID
	entity.CreatedAt = seedDate
	entity.UpdatedAt = completedAt

	if entity.ID == "" {
		return db.Create(&entity).Error
	}
	return db.Save(&entity).Error
}

func upsertHistoricalPurchase(db *gorm.DB, deal pipeline.Deal, salesUser user.User, items []pipeline.DealProductItem) error {
	var purchase customer_purchase.CustomerPurchaseHistory
	err := db.Where("deal_id = ?", deal.ID).First(&purchase).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	purchaseItems := make([]customer_purchase.PurchaseProduct, 0, len(items))
	for _, item := range items {
		categoryID := ""
		if item.ProductCategoryID != nil {
			categoryID = *item.ProductCategoryID
		}
		purchaseItems = append(purchaseItems, customer_purchase.PurchaseProduct{
			ProductID:           item.ProductID,
			ProductName:         item.ProductName,
			ProductSKU:          item.ProductSKU,
			Quantity:            item.Quantity,
			UnitPrice:           item.UnitPrice,
			Subtotal:            item.Subtotal,
			ProductCategoryID:   categoryID,
			ProductCategoryName: item.ProductCategoryName,
		})
	}

	purchase.AccountID = deal.AccountID
	purchase.DealID = deal.ID
	if deal.ActualCloseDate != nil {
		purchase.PurchaseDate = *deal.ActualCloseDate
	} else {
		purchase.PurchaseDate = deal.CreatedAt
	}
	purchase.TotalAmount = deal.Value
	purchase.TotalItems = len(purchaseItems)
	purchase.Products = mustJSON(purchaseItems)
	purchase.SalesRepID = stringPtr(salesUser.ID)
	purchase.SalesRepName = salesUser.Name
	purchase.SourceLeadID = nil
	purchase.SourceType = "pipeline"
	purchase.CustomerLifetimeValue = deal.Value
	purchase.CreatedAt = purchase.PurchaseDate
	purchase.UpdatedAt = purchase.PurchaseDate

	if purchase.ID == "" {
		return db.Create(&purchase).Error
	}
	return database.DB.Save(&purchase).Error
}

func seedHistoricalMonthlyTargets(start time.Time, now time.Time, salesUsers []user.User) error {
	created := 0
	for cursor := time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, start.Location()); !cursor.After(now); cursor = cursor.AddDate(0, 1, 0) {
		for index, salesUser := range salesUsers {
			year := cursor.Year()
			month := int(cursor.Month())
			var existing monthlytargetdomain.MonthlyTarget
			err := database.DB.Where("user_id = ? AND year = ? AND month = ? AND deleted_at IS NULL", salesUser.ID, year, month).
				First(&existing).Error
			if err == nil {
				continue
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			seasonalBoost := int64((month%4 + index%3) * 1000000000)
			target := monthlytargetdomain.MonthlyTarget{
				UserID:       stringPtr(salesUser.ID),
				GroupID:      nil,
				BrickID:      nil,
				Year:         year,
				Month:        month,
				TargetAmount: 9000000000 + seasonalBoost,
				CreatedAt:    cursor,
				UpdatedAt:    cursor,
			}
			if err := database.DB.Create(&target).Error; err != nil {
				return err
			}
			created++
		}
	}
	log.Printf("Seeded %d missing historical monthly user targets", created)
	return nil
}
