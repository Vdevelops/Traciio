package schedule

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/schedule"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

type assignmentRepository struct {
	db *gorm.DB
}

// NewAssignmentRepository creates a new schedule assignment repository
func NewAssignmentRepository(db *gorm.DB) interfaces.ScheduleAssignmentRepository {
	return &assignmentRepository{db: db}
}

func (r *assignmentRepository) FindByID(id string) (*schedule.ScheduleAssignment, error) {
	var sa schedule.ScheduleAssignment
	err := r.db.
		Preload("Schedule").
		Preload("User").
		Where("id = ?", id).
		First(&sa).Error
	if err != nil {
		return nil, err
	}
	return &sa, nil
}

func (r *assignmentRepository) FindByScheduleID(scheduleID string) ([]schedule.ScheduleAssignment, error) {
	var assignments []schedule.ScheduleAssignment
	err := r.db.
		Preload("Schedule").
		Preload("User").
		Where("schedule_id = ?", scheduleID).
		Order("assigned_at DESC").
		Find(&assignments).Error
	return assignments, err
}

func (r *assignmentRepository) FindByUserID(userID string) ([]schedule.ScheduleAssignment, error) {
	var assignments []schedule.ScheduleAssignment
	err := r.db.
		Preload("Schedule").
		Preload("User").
		Where("user_id = ?", userID).
		Order("assigned_at DESC").
		Find(&assignments).Error
	return assignments, err
}

func (r *assignmentRepository) Create(sa *schedule.ScheduleAssignment) error {
	return r.db.Create(sa).Error
}

func (r *assignmentRepository) CreateBatch(assignments []schedule.ScheduleAssignment) error {
	if len(assignments) == 0 {
		return nil
	}
	return r.db.Create(&assignments).Error
}

func (r *assignmentRepository) Update(sa *schedule.ScheduleAssignment) error {
	// Clear relations to avoid updating them
	sa.Schedule = nil
	sa.User = nil

	return r.db.Model(sa).Omit("Schedule", "User").Updates(sa).Error
}

func (r *assignmentRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&schedule.ScheduleAssignment{}).Error
}

func (r *assignmentRepository) DeleteByScheduleID(scheduleID string) error {
	return r.db.Where("schedule_id = ?", scheduleID).Delete(&schedule.ScheduleAssignment{}).Error
}

