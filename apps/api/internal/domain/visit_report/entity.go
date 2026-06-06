package visit_report

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// VisitReport represents a visit report entity
type VisitReport struct {
	ID               string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AccountID        *string        `gorm:"type:uuid;index" json:"account_id,omitempty"` // Optional: required if DealID exists, optional if LeadID exists
	ContactID        *string        `gorm:"type:uuid;index" json:"contact_id,omitempty"`
	DealID           *string        `gorm:"type:uuid;index" json:"deal_id,omitempty"` // Optional link to deal
	LeadID           *string        `gorm:"type:uuid;index" json:"lead_id,omitempty"` // Optional link to lead
	SalesRepID       string         `gorm:"type:uuid;not null;index" json:"sales_rep_id"`
	BrickID          *string        `gorm:"type:uuid;index" json:"brick_id,omitempty"` // Brick/Area assignment
	VisitDate        time.Time      `gorm:"type:timestamp;not null;index" json:"visit_date"`
	CheckInTime      *time.Time     `gorm:"type:timestamp;index" json:"check_in_time,omitempty"`
	CheckOutTime     *time.Time     `gorm:"type:timestamp;index" json:"check_out_time,omitempty"`
	CheckInLocation  datatypes.JSON `gorm:"type:jsonb" json:"check_in_location,omitempty"`
	CheckOutLocation datatypes.JSON `gorm:"type:jsonb" json:"check_out_location,omitempty"`
	Purpose          string         `gorm:"type:text;not null" json:"purpose"`
	Notes            string         `gorm:"type:text" json:"notes"`
	Outcome          string         `gorm:"type:varchar(50);index" json:"outcome,omitempty"`                 // positive, neutral, negative, very_positive
	NextSteps        string         `gorm:"type:text" json:"next_steps,omitempty"`                           // Action items after visit
	Photos           datatypes.JSON `gorm:"type:jsonb" json:"photos,omitempty"`                              // Array of photo URLs
	Metadata         datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'" json:"metadata,omitempty"`      // Product interests and structured visit context
	Status           string         `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"` // pending, completed
	ApprovedBy       *string        `gorm:"type:uuid;index" json:"approved_by,omitempty"`
	ApprovedAt       *time.Time     `gorm:"type:timestamp;index" json:"approved_at,omitempty"`
	RejectionReason  *string        `gorm:"type:text" json:"rejection_reason,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations (for preloading)
	Account  interface{} `gorm:"-" json:"account,omitempty"`
	Contact  interface{} `gorm:"-" json:"contact,omitempty"`
	SalesRep interface{} `gorm:"-" json:"sales_rep,omitempty"`
}

// Location represents GPS location
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address,omitempty"`
}

// TableName specifies the table name for VisitReport
func (VisitReport) TableName() string {
	return "visit_reports"
}

// BeforeCreate hook to generate UUID
func (vr *VisitReport) BeforeCreate(tx *gorm.DB) error {
	if vr.ID == "" {
		vr.ID = uuid.New().String()
	}
	return nil
}

// GPSMetadata represents GPS metadata for validation
type GPSMetadata struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Accuracy  float64 `json:"accuracy,omitempty"`  // GPS accuracy in meters
	Timestamp int64   `json:"timestamp,omitempty"` // Unix timestamp when GPS was captured
	Altitude  float64 `json:"altitude,omitempty"`  // Altitude in meters
	Heading   float64 `json:"heading,omitempty"`   // Heading in degrees
	Speed     float64 `json:"speed,omitempty"`     // Speed in m/s
}

func NormalizeStatus(status string) string {
	switch status {
	case "completed", "approved", "rejected", "cancelled":
		return "completed"
	default:
		return "pending"
	}
}
