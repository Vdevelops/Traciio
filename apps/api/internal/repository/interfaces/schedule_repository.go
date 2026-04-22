package interfaces

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/schedule"
)

// ScheduleRepository defines the interface for schedule repository
type ScheduleRepository interface {
	// FindByID finds a schedule by ID
	FindByID(id string) (*schedule.Schedule, error)

	// List returns a list of schedules with pagination
	List(req *schedule.ListSchedulesRequest) ([]schedule.Schedule, int64, error)

	// Create creates a new schedule
	Create(s *schedule.Schedule) error

	// Update updates a schedule
	Update(s *schedule.Schedule) error

	// Delete soft deletes a schedule
	Delete(id string) error

	// FindConflicts finds schedules that conflict with the given time range
	FindConflicts(userID string, startTime, endTime time.Time, excludeScheduleID string) ([]schedule.Schedule, error)

	// FindByUserID finds schedules assigned to a user
	FindByUserID(userID string, startDate, endDate *time.Time) ([]schedule.Schedule, error)

	// FindByDateRange finds schedules within a date range
	FindByDateRange(startDate, endDate time.Time, assignedTo *string) ([]schedule.Schedule, error)
}

// ScheduleAssignmentRepository defines the interface for schedule assignment repository
type ScheduleAssignmentRepository interface {
	// FindByID finds a schedule assignment by ID
	FindByID(id string) (*schedule.ScheduleAssignment, error)

	// FindByScheduleID finds assignments by schedule ID
	FindByScheduleID(scheduleID string) ([]schedule.ScheduleAssignment, error)

	// FindByUserID finds assignments by user ID
	FindByUserID(userID string) ([]schedule.ScheduleAssignment, error)

	// Create creates a new schedule assignment
	Create(sa *schedule.ScheduleAssignment) error

	// CreateBatch creates multiple schedule assignments
	CreateBatch(assignments []schedule.ScheduleAssignment) error

	// Update updates a schedule assignment
	Update(sa *schedule.ScheduleAssignment) error

	// Delete soft deletes a schedule assignment
	Delete(id string) error

	// DeleteByScheduleID deletes all assignments for a schedule
	DeleteByScheduleID(scheduleID string) error
}

