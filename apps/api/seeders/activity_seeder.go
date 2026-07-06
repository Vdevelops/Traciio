package seeders

import (
	"errors"
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/activity"
	lead "github.com/gilabs/crm-healthcare/api/internal/domain/lead"
	"github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
	"github.com/gilabs/crm-healthcare/api/internal/domain/visit_report"
	"gorm.io/gorm"
)

// SeedActivities creates activity timeline rows that mirror the visit/task/deal
// graph and surface the same product interest metadata used by the UI.
func SeedActivities() error {
	ctx, err := loadSeedFixtureContext()
	if err != nil {
		return err
	}

	records := buildCRMSeedRecords(ctx)
	seeded := 0
	for _, record := range records {
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

		var relatedVisit visit_report.VisitReport
		visitFound := false
		err = database.DB.Where("lead_id = ? AND purpose = ?", sourceLead.ID, record.VisitPurpose).First(&relatedVisit).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		visitFound = err == nil

		var entity activity.Activity
		err = database.DB.Where("lead_id = ? AND description = ?", sourceLead.ID, record.ActivityDescription).First(&entity).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		entity.Type = record.ActivityTypeCode
		if activityTypeID, ok := ctx.ActivityTypeIDs[record.ActivityTypeCode]; ok {
			entity.ActivityTypeID = stringPtr(activityTypeID)
		}
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
		entity.LeadID = stringPtr(sourceLead.ID)
		entity.UserID = record.Owner.ID
		entity.Description = record.ActivityDescription
		if record.CheckOutTime != nil {
			entity.Timestamp = *record.CheckOutTime
		} else {
			entity.Timestamp = record.VisitDate
		}

		metadata := map[string]any{
			"seed_source":       crmSeedSource,
			"product_interests": record.ProductInterests,
		}
		if visitFound {
			metadata["visit_report_id"] = relatedVisit.ID
			metadata["visit_date"] = relatedVisit.VisitDate.Format("2006-01-02")
		}
		if dealFound {
			metadata["deal_id"] = relatedDeal.ID
			metadata["deal_stage_code"] = record.DealStageCode
		}
		entity.Metadata = mustJSON(metadata)

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

	log.Printf("Seeded %d CRM-synced activities", seeded)
	return nil
}
