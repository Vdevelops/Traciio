package schedule

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/schedule"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

// Field selections for optimized queries (reduce memory usage)
const (
	userFields        = "id, name, email, avatar_url"
	taskFields        = "id, title, description, status, priority, due_date, assigned_to"
	taskFieldsMinimal = "id, title, status, assigned_to"
)

// Preload relation names
const (
	preloadTask           = "Task"
	preloadTaskAssignedUser = "Task.AssignedUser"
)

// Query conditions
const (
	whereUserID        = "user_id = ?"
	whereScheduledAtGE = "scheduled_at >= ?"
	orderScheduledAtASC = "scheduled_at ASC"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new schedule repository
func NewRepository(db *gorm.DB) interfaces.ScheduleRepository {
	return &repository{db: db}
}

func (r *repository) FindByID(id string) (*schedule.Schedule, error) {
	var s schedule.Schedule
	err := r.db.
		Preload(preloadTask, func(db *gorm.DB) *gorm.DB {
			return db.Select(taskFields)
		}).
		Preload(preloadTaskAssignedUser, func(db *gorm.DB) *gorm.DB {
			return db.Select(userFields)
		}).
		Where("id = ?", id).
		First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *repository) List(req *schedule.ListSchedulesRequest) ([]schedule.Schedule, int64, error) {
	var schedules []schedule.Schedule
	var total int64

	query := r.db.Model(&schedule.Schedule{})

	// Apply RBAC scope filtering
	if len(req.ScopedUserIDs) > 0 {
		query = query.Where("user_id IN ?", req.ScopedUserIDs)
	}

	// Apply filters
	if req.Search != "" {
		// Optimized: Use Full Text Search instead of LIKE %...%
		// Uses GIN index on title, description
		query = query.Where(
			"to_tsvector('english', title || ' ' || COALESCE(description, '')) @@ plainto_tsquery('english', ?)",
			req.Search,
		)
	}

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	if req.TaskID != "" {
		query = query.Where("task_id = ?", req.TaskID)
	}

	if req.UserID != "" {
		query = query.Where(whereUserID, req.UserID)
	}

	if req.ScheduledAtFrom != nil {
		query = query.Where(whereScheduledAtGE, *req.ScheduledAtFrom)
	}

	if req.ScheduledAtTo != nil {
		query = query.Where("scheduled_at <= ?", *req.ScheduledAtTo)
	}

	if req.GoogleCalendarSyncStatus != "" {
		query = query.Where("google_calendar_sync_status = ?", req.GoogleCalendarSyncStatus)
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
		Preload(preloadTask, func(db *gorm.DB) *gorm.DB {
			return db.Select(taskFields)
		}).
		Preload(preloadTaskAssignedUser, func(db *gorm.DB) *gorm.DB {
			return db.Select(userFields)
		}).
		Order(orderScheduledAtASC).
		Offset(offset).
		Limit(perPage).
		Find(&schedules).Error
	if err != nil {
		return nil, 0, err
	}

	return schedules, total, nil
}

func (r *repository) Create(s *schedule.Schedule) error {
	return r.db.Create(s).Error
}

func (r *repository) Update(s *schedule.Schedule) error {
	// Clear relations to avoid updating them
	s.Task = nil

	return r.db.Model(s).Omit("Task").Updates(s).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&schedule.Schedule{}).Error
}

func (r *repository) FindConflicts(userID string, startTime, endTime time.Time, excludeScheduleID string) ([]schedule.Schedule, error) {
	var conflicts []schedule.Schedule

	query := r.db.Model(&schedule.Schedule{}).
		Where(whereUserID, userID).
		Where("status IN ?", []string{"pending", "confirmed"}).
		Where(
			// Check for overlap: new schedule starts before existing ends AND new schedule ends after existing starts
			whereScheduledAtGE+" AND scheduled_at <= ?",
			startTime, endTime,
		)

	if excludeScheduleID != "" {
		query = query.Where("id != ?", excludeScheduleID)
	}

	err := query.
		Preload(preloadTask, func(db *gorm.DB) *gorm.DB {
			return db.Select(taskFieldsMinimal)
		}).
		Preload(preloadTaskAssignedUser, func(db *gorm.DB) *gorm.DB {
			return db.Select(userFields)
		}).
		Order(orderScheduledAtASC).
		Find(&conflicts).Error

	return conflicts, err
}

func (r *repository) FindByUserID(userID string, startDate, endDate *time.Time) ([]schedule.Schedule, error) {
	var schedules []schedule.Schedule

	query := r.db.Model(&schedule.Schedule{}).
		Where(whereUserID, userID)

	if startDate != nil {
		query = query.Where(whereScheduledAtGE, *startDate)
	}

	if endDate != nil {
		// Add one day to include the entire end date
		endDateWithTime := endDate.Add(24 * time.Hour)
		query = query.Where("scheduled_at < ?", endDateWithTime)
	}

	err := query.
		Preload(preloadTask, func(db *gorm.DB) *gorm.DB {
			return db.Select(taskFields)
		}).
		Preload(preloadTaskAssignedUser, func(db *gorm.DB) *gorm.DB {
			return db.Select(userFields)
		}).
		Order(orderScheduledAtASC).
		Find(&schedules).Error

	return schedules, err
}

func (r *repository) FindByDateRange(startDate, endDate time.Time, assignedTo *string) ([]schedule.Schedule, error) {
	var schedules []schedule.Schedule

	query := r.db.Model(&schedule.Schedule{}).
		Where(whereScheduledAtGE, startDate).
		Where("scheduled_at < ?", endDate.Add(24*time.Hour))

	if assignedTo != nil {
		query = query.Where(whereUserID, *assignedTo)
	}

	err := query.
		Preload(preloadTask, func(db *gorm.DB) *gorm.DB {
			return db.Select(taskFields)
		}).
		Preload(preloadTaskAssignedUser, func(db *gorm.DB) *gorm.DB {
			return db.Select(userFields)
		}).
		Order(orderScheduledAtASC).
		Find(&schedules).Error

	return schedules, err
}

