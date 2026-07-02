package seeders

import (
	"errors"
	"log"
	"strings"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead"
	leadqualification "github.com/gilabs/crm-healthcare/api/internal/domain/lead_qualification"
	"gorm.io/gorm"
)

// SeedLeads seeds CRM leads together with qualification data that becomes the
// canonical source for synced pipeline, task, visit, and activity seed records.
func SeedLeads() error {
	ctx, err := loadSeedFixtureContext()
	if err != nil {
		return err
	}

	records := buildCRMSeedRecords(ctx)
	for _, record := range records {
		var entity lead.Lead
		err := database.DB.Where("email = ?", record.LeadEmail).First(&entity).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if strings.TrimSpace(entity.ID) == "" {
			entity.CreatedBy = ctx.AdminUser.ID
		}

		entity.FirstName = record.FirstName
		entity.LastName = record.LastName
		entity.CompanyName = record.CompanyName
		entity.Email = record.LeadEmail
		entity.Phone = record.Phone
		entity.JobTitle = record.JobTitle
		entity.Industry = record.Industry
		entity.LeadSource = record.LeadSource
		entity.LeadStatus = record.LeadStatusCode
		if statusID, ok := ctx.LeadStatusIDs[record.LeadStatusCode]; ok {
			entity.LeadStatusID = stringPtr(statusID)
		}
		entity.LeadScore = record.LeadProbability
		entity.Probability = record.LeadProbability
		entity.EstimatedValue = record.EstimatedValue
		entity.BudgetConfirmed = true
		entity.BudgetAmount = int64Ptr(record.BudgetAmount)
		entity.AuthorityConfirmed = true
		entity.AuthorityPerson = record.AuthorityPerson
		entity.NeedConfirmed = len(record.NeedProducts) > 0
		entity.NeedDescription = record.NeedDescription
		entity.TimelineConfirmed = true
		entity.ExpectedCloseDate = &record.ExpectedCloseDate
		entity.AssignedTo = stringPtr(record.Owner.ID)
		entity.Notes = strings.TrimSpace(record.LeadNotes + " [" + crmSeedSource + "]")
		entity.Address = record.Address
		entity.City = record.City
		entity.Province = record.Province
		entity.PostalCode = record.PostalCode
		entity.Country = "Indonesia"
		entity.Website = record.Website

		if entity.ID == "" {
			if err := database.DB.Create(&entity).Error; err != nil {
				return err
			}
		} else {
			if err := database.DB.Save(&entity).Error; err != nil {
				return err
			}
		}

		var qualification leadqualification.LeadQualificationChecklist
		err = database.DB.Where("lead_id = ?", entity.ID).First(&qualification).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		qualification.LeadID = entity.ID
		qualification.BudgetTargetAmount = record.BudgetAmount
		qualification.BudgetTargetCurrency = "IDR"
		qualification.BudgetConfirmed = true
		qualification.BudgetNotes = "Budget tervalidasi dari skenario seeder sinkron."
		qualification.AuthorityTargetPerson = record.AuthorityPerson
		qualification.AuthorityTargetRole = record.JobTitle
		qualification.AuthorityConfirmed = true
		qualification.AuthorityNotes = "PIC pengambil keputusan sudah terverifikasi."
		qualification.NeedTargetProducts = mustJSON(record.NeedProducts)
		qualification.NeedPriorityLevel = record.NeedPriorityLevel
		qualification.NeedConfirmed = len(record.NeedProducts) > 0
		qualification.NeedNotes = record.NeedDescription
		qualification.TimelineTargetDate = &record.ExpectedCloseDate
		qualification.TimelineFlexibility = "fixed"
		qualification.TimelineConfirmed = true
		qualification.TimelineNotes = "Timeline follow-up diturunkan dari skenario CRM sync."

		if qualification.ID == "" {
			if err := database.DB.Create(&qualification).Error; err != nil {
				return err
			}
		} else {
			if err := database.DB.Save(&qualification).Error; err != nil {
				return err
			}
		}
	}

	log.Printf("Seeded %d CRM-synced leads with qualification data", len(records))
	return nil
}
