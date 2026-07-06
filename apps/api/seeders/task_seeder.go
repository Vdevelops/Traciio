package seeders

import (
	"errors"
	"log"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	lead "github.com/gilabs/crm-healthcare/api/internal/domain/lead"
	"github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
	"github.com/gilabs/crm-healthcare/api/internal/domain/schedule"
	"github.com/gilabs/crm-healthcare/api/internal/domain/task"
	"github.com/gilabs/crm-healthcare/api/internal/domain/visit_report"
	"gorm.io/gorm"
)

// SeedTasks creates follow-up tasks derived from the same lead/deal product
// interests, then mirrors them into schedule rows for calendar-based surfaces.
func SeedTasks() error {
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

		var entity task.Task
		err = database.DB.Where("lead_id = ? AND title = ?", sourceLead.ID, record.TaskTitle).First(&entity).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		dueDate := record.ExpectedCloseDate.Add(-72 * time.Hour)
		entity.Title = record.TaskTitle
		entity.Description = record.TaskDescription + " [" + crmSeedSource + "]"
		entity.Type = record.TaskType
		entity.Status = "pending"
		entity.Priority = record.TaskPriority
		entity.DueDate = &dueDate
		entity.AssignedTo = stringPtr(record.Owner.ID)
		entity.AssignedFrom = stringPtr(ctx.AdminUser.ID)
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
		entity.TaskSource = record.TaskSource
		entity.QuickActionType = "product_interest_follow_up"
		payload := map[string]any{
			"seed_source":       crmSeedSource,
			"product_interests": record.ProductInterests,
			"lead_id":           sourceLead.ID,
		}
		if dealFound {
			payload["deal_id"] = relatedDeal.ID
		}
		if visitFound {
			payload["visit_report_id"] = relatedVisit.ID
		}
		entity.QuickActionPayload = mustJSON(payload)
		entity.CreatedBy = ctx.AdminUser.ID

		if entity.ID == "" {
			if err := database.DB.Create(&entity).Error; err != nil {
				return err
			}
		} else {
			if err := database.DB.Save(&entity).Error; err != nil {
				return err
			}
		}

		var taskSchedule schedule.Schedule
		err = database.DB.Where("task_id = ?", entity.ID).First(&taskSchedule).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		scheduledAt := dueDate.Add(-1 * time.Hour)
		reminderMinutesBefore := 60
		description := entity.Description
		taskSchedule.TaskID = stringPtr(entity.ID)
		taskSchedule.UserID = record.Owner.ID
		taskSchedule.Title = "Reminder: " + entity.Title
		taskSchedule.Description = &description
		taskSchedule.ScheduledAt = scheduledAt
		taskSchedule.Status = "pending"
		taskSchedule.ReminderMinutesBefore = &reminderMinutesBefore
		taskSchedule.CreatedBy = ctx.AdminUser.ID

		if taskSchedule.ID == "" {
			if err := database.DB.Create(&taskSchedule).Error; err != nil {
				return err
			}
		} else {
			if err := database.DB.Save(&taskSchedule).Error; err != nil {
				return err
			}
		}

		seeded++
	}

	log.Printf("Seeded %d CRM-synced tasks with schedules", seeded)
	return nil
}
