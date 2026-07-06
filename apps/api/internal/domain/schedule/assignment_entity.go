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
