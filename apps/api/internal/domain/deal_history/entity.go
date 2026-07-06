package deal_history

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DealHistory represents the history of deal stage changes
type DealHistory struct {
	ID              string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	DealID          string         `gorm:"type:uuid;not null;index" json:"deal_id"`
	FromStageID     *string        `gorm:"type:uuid" json:"from_stage_id,omitempty"` // NULL for initial creation
	FromStageName   string         `gorm:"type:varchar(255)" json:"from_stage_name,omitempty"`
	ToStageID       string         `gorm:"type:uuid;not null" json:"to_stage_id"`
	ToStageName     string         `gorm:"type:varchar(255);not null" json:"to_stage_name"`
	FromProbability int            `gorm:"type:integer;default:0" json:"from_probability"`
	ToProbability   int            `gorm:"type:integer;default:0" json:"to_probability"`
	DaysInPrevStage *int           `gorm:"type:integer" json:"days_in_prev_stage,omitempty"` // NULL for first stage
	ChangedBy       string         `gorm:"type:uuid;not null;index" json:"changed_by"`
	ChangedAt       time.Time      `gorm:"type:timestamp;not null;default:current_timestamp" json:"changed_at"`
	Reason          string         `gorm:"type:text" json:"reason,omitempty"` // Optional reason for stage change
	Notes           string         `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for DealHistory
func (DealHistory) TableName() string {
	return "deal_histories"
}

// BeforeCreate hook to generate UUID
func (dh *DealHistory) BeforeCreate(tx *gorm.DB) error {
	if dh.ID == "" {
		dh.ID = uuid.New().String()
	}
	if dh.ChangedAt.IsZero() {
		dh.ChangedAt = time.Now()
	}
	return nil
}
