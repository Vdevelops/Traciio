package seeders

import (
	"errors"
	"log"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	lead "github.com/gilabs/crm-healthcare/api/internal/domain/lead"
	"github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
	"github.com/gilabs/crm-healthcare/api/internal/domain/visit_report"
	"gorm.io/gorm"
)

// SeedVisitReports creates visit logs that are linked directly to the seeded
// lead/deal records and carry the same product interest metadata.
func SeedVisitReports() error {
	ctx, err := loadSeedFixtureContext()
	if err != nil {
		return err
	}

	records := buildCRMSeedRecords(ctx)
	seeded := 0
	for index, record := range records {
		var sourceLead lead.Lead
		if err := database.DB.Where("email = ?", record.LeadEmail).First(&sourceLead).Error; err != nil {
			return err
		}

		var relatedDeal pipeline.Deal
		dealFound := false
		if record.DealTitle != "" {
			err := database.DB.Where("lead_id = ? AND title = ?", sourceLead.ID, record.DealTitle).First(&relatedDeal).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			dealFound = err == nil
		}

		var entity visit_report.VisitReport
		err = database.DB.Where("lead_id = ? AND purpose = ?", sourceLead.ID, record.VisitPurpose).First(&entity).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		entity.LeadID = stringPtr(sourceLead.ID)
		if record.Account != nil {
			entity.AccountID = stringPtr(record.Account.ID)
		}
		if record.Contact != nil {
			entity.ContactID = stringPtr(record.Contact.ID)
		}
		if dealFound {
			entity.DealID = stringPtr(relatedDeal.ID)
		} else {
			entity.DealID = nil
		}
		entity.SalesRepID = record.Owner.ID
		entity.BrickID = record.Owner.BrickID
		entity.VisitDate = record.VisitDate
		entity.CheckInTime = record.CheckInTime
		entity.CheckOutTime = record.CheckOutTime
		entity.CheckInLocation = mustJSON(visit_report.Location{
			Latitude:  -6.99 + (float64(index) * 0.01),
			Longitude: 110.42 + (float64(index) * 0.01),
			Address:   record.Address,
		})
		entity.CheckOutLocation = entity.CheckInLocation
		entity.Purpose = record.VisitPurpose
		entity.Notes = record.VisitNotes + " [" + crmSeedSource + "]"
		entity.Outcome = record.VisitOutcome
		entity.NextSteps = record.VisitNextSteps
		entity.Metadata = record.VisitProductMetadata
		entity.Status = record.VisitStatus

		if record.VisitStatus == "approved" {
			entity.ApprovedBy = stringPtr(ctx.AdminUser.ID)
			if record.CheckOutTime != nil {
				approvedAt := record.CheckOutTime.Add(2 * time.Hour)
				entity.ApprovedAt = &approvedAt
			}
		} else {
			entity.ApprovedBy = nil
			entity.ApprovedAt = nil
		}

		if entity.ID == "" {
			if err := database.DB.Create(&entity).Error; err != nil {
				return err
			}
		} else {
			if err := database.DB.Save(&entity).Error; err != nil {
				return err
			}
		}
		seeded++
	}

	log.Printf("Seeded %d CRM-synced visit reports", seeded)
	return nil
}
