package reminder

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/reminder"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new reminder repository
func NewRepository(db *gorm.DB) interfaces.ReminderRepository {
	return &repository{db: db}
}

func (r *repository) FindByID(id string) (*reminder.Reminder, error) {
	var rem reminder.Reminder
	err := r.db.
		Preload("Task").
		Where("id = ?", id).
		First(&rem).Error
	if err != nil {
		return nil, err
	}
	return &rem, nil
}

func (r *repository) FindByTaskID(taskID string) ([]reminder.Reminder, error) {
	var reminders []reminder.Reminder
	err := r.db.
		Preload("Task").
		Where("task_id = ?", taskID).
		Order("remind_at ASC").
		Find(&reminders).Error
	if err != nil {
		return nil, err
	}
	return reminders, nil
}

func (r *repository) List(req *reminder.ListRemindersRequest) ([]reminder.Reminder, int64, error) {
	var reminders []reminder.Reminder
	var total int64

	query := r.db.Model(&reminder.Reminder{})

	// Apply filters
	if req.TaskID != "" {
		query = query.Where("task_id = ?", req.TaskID)
	}

	if req.ReminderType != "" {
		query = query.Where("reminder_type = ?", req.ReminderType)
	}

	if req.IsSent != nil {
		query = query.Where("is_sent = ?", *req.IsSent)
	}

	if req.RemindAtFrom != nil {
		query = query.Where("remind_at >= ?", *req.RemindAtFrom)
	}

	if req.RemindAtTo != nil {
		query = query.Where("remind_at <= ?", *req.RemindAtTo)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}

	offset := (page - 1) * perPage

	// Fetch data with preload
	err := query.
		Preload("Task").
		Order("remind_at ASC").
		Offset(offset).
		Limit(perPage).
		Find(&reminders).Error
	if err != nil {
		return nil, 0, err
	}

	return reminders, total, nil
}

func (r *repository) Create(rem *reminder.Reminder) error {
	return r.db.Create(rem).Error
}

func (r *repository) Update(rem *reminder.Reminder) error {
	// Clear relations to avoid updating them
	rem.Task = nil

	return r.db.Model(rem).Omit("Task").Updates(rem).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&reminder.Reminder{}).Error
}

func (r *repository) FindPendingReminders(beforeTime time.Time) ([]reminder.Reminder, error) {
	var reminders []reminder.Reminder
	err := r.db.
		Preload("Task").
		Where("remind_at <= ?", beforeTime).
		Where("is_sent = ?", false).
		Order("remind_at ASC").
		Find(&reminders).Error
	if err != nil {
		return nil, err
	}
	return reminders, nil
}

func (r *repository) MarkAsSent(id string, sentAt time.Time) error {
	return r.db.Model(&reminder.Reminder{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_sent": true,
			"sent_at": sentAt,
		}).Error
}



