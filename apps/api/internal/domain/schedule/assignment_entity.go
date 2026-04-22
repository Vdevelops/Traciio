package schedule

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ScheduleAssignment represents a schedule assignment (for bulk assignment tracking)
type ScheduleAssignment struct {
	ID         string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ScheduleID string         `gorm:"type:uuid;not null;index" json:"schedule_id"`
	Schedule   *Schedule      `gorm:"foreignKey:ScheduleID" json:"schedule,omitempty"`
	UserID     string         `gorm:"type:uuid;not null;index" json:"user_id"`
	User       *UserRef       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	AssignedAt time.Time      `gorm:"type:timestamp;default:now()" json:"assigned_at"`
	Status     string         `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"` // pending, accepted, rejected
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for ScheduleAssignment
func (ScheduleAssignment) TableName() string {
	return "schedule_assignments"
}

// BeforeCreate hook to generate UUID
func (sa *ScheduleAssignment) BeforeCreate(tx *gorm.DB) error {
	if sa.ID == "" {
		sa.ID = uuid.New().String()
	}
	return nil
}

// ScheduleAssignmentResponse represents schedule assignment response DTO
type ScheduleAssignmentResponse struct {
	ID         string           `json:"id"`
	ScheduleID string           `json:"schedule_id"`
	Schedule   *ScheduleResponse `json:"schedule,omitempty"`
	UserID     string           `json:"user_id"`
	User       *UserRefResponse `json:"user,omitempty"`
	AssignedAt time.Time        `json:"assigned_at"`
	Status     string           `json:"status"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}

// ToScheduleAssignmentResponse converts ScheduleAssignment to ScheduleAssignmentResponse
func (sa *ScheduleAssignment) ToScheduleAssignmentResponse() *ScheduleAssignmentResponse {
	resp := &ScheduleAssignmentResponse{
		ID:         sa.ID,
		ScheduleID: sa.ScheduleID,
		UserID:     sa.UserID,
		AssignedAt: sa.AssignedAt,
		Status:     sa.Status,
		CreatedAt:  sa.CreatedAt,
		UpdatedAt:  sa.UpdatedAt,
	}

	if sa.Schedule != nil {
		resp.Schedule = sa.Schedule.ToScheduleResponse()
	}

	if sa.User != nil {
		resp.User = &UserRefResponse{
			ID:    sa.User.ID,
			Name:  sa.User.Name,
			Email: sa.User.Email,
		}
	}

	return resp
}

