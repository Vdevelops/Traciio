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
	Outcome          string         `gorm:"type:varchar(50);index" json:"outcome,omitempty"`               // positive, neutral, negative, very_positive
	NextSteps        string         `gorm:"type:text" json:"next_steps,omitempty"`                         // Action items after visit
	Photos           datatypes.JSON `gorm:"type:jsonb" json:"photos,omitempty"`                            // Array of photo URLs
	Status           string         `gorm:"type:varchar(20);not null;default:'draft';index" json:"status"` // draft, submitted, approved, rejected
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

// VisitReportResponse represents visit report response DTO
type VisitReportResponse struct {
	ID               string      `json:"id"`
	AccountID        *string     `json:"account_id,omitempty"`
	ContactID        *string     `json:"contact_id,omitempty"`
	DealID           *string     `json:"deal_id,omitempty"`
	LeadID           *string     `json:"lead_id,omitempty"`
	Type             string      `json:"type"` // "account", "deal", or "lead"
	SalesRepID       string      `json:"sales_rep_id"`
	BrickID          *string     `json:"brick_id,omitempty"`
	VisitDate        time.Time   `json:"visit_date"`
	CheckInTime      *time.Time  `json:"check_in_time,omitempty"`
	CheckOutTime     *time.Time  `json:"check_out_time,omitempty"`
	CheckInLocation  *Location   `json:"check_in_location,omitempty"`
	CheckOutLocation *Location   `json:"check_out_location,omitempty"`
	Purpose          string      `json:"purpose"`
	Notes            string      `json:"notes"`
	Photos           []string    `json:"photos,omitempty"`
	Status           string      `json:"status"`
	ApprovedBy       *string     `json:"approved_by,omitempty"`
	ApprovedAt       *time.Time  `json:"approved_at,omitempty"`
	RejectionReason  *string     `json:"rejection_reason,omitempty"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
	Account          interface{} `json:"account,omitempty"`
	Contact          interface{} `json:"contact,omitempty"`
	Deal             interface{} `json:"deal,omitempty"`
	Lead             interface{} `json:"lead,omitempty"`
	SalesRep         interface{} `json:"sales_rep,omitempty"`
}

// ToVisitReportResponse converts VisitReport to VisitReportResponse
func (vr *VisitReport) ToVisitReportResponse() *VisitReportResponse {
	var photos []string
	if vr.Photos != nil {
		// Parse JSON array to []string
		// This will be handled by GORM's datatypes.JSON
		// For now, we'll need to handle this in the service layer
	}

	resp := &VisitReportResponse{
		ID:           vr.ID,
		AccountID:    vr.AccountID, // Already *string
		ContactID:    vr.ContactID,
		DealID:       vr.DealID,
		LeadID:       vr.LeadID,
		SalesRepID:   vr.SalesRepID,
		BrickID:      vr.BrickID,
		VisitDate:    vr.VisitDate,
		CheckInTime:  vr.CheckInTime,
		CheckOutTime: vr.CheckOutTime,
		// CheckInLocation and CheckOutLocation will be parsed in service layer
		CheckInLocation:  nil,
		CheckOutLocation: nil,
		Purpose:          vr.Purpose,
		Notes:            vr.Notes,
		Photos:           photos,
		Status:           vr.Status,
		ApprovedBy:       vr.ApprovedBy,
		ApprovedAt:       vr.ApprovedAt,
		RejectionReason:  vr.RejectionReason,
		CreatedAt:        vr.CreatedAt,
		UpdatedAt:        vr.UpdatedAt,
		Account:          vr.Account,
		Contact:          vr.Contact,
		SalesRep:         vr.SalesRep,
	}
	return resp
}

// CreateVisitReportRequest represents create visit report request DTO
type CreateVisitReportRequest struct {
	AccountID        *string   `json:"account_id" binding:"omitempty,uuid"` // Optional: required if DealID exists
	ContactID        *string   `json:"contact_id" binding:"omitempty,uuid"`
	DealID           *string   `json:"deal_id" binding:"omitempty,uuid"`      // Optional link to deal
	LeadID           *string   `json:"lead_id" binding:"omitempty,uuid"`      // Optional link to lead
	SalesRepID       string    `json:"sales_rep_id" binding:"omitempty,uuid"` // Will be set from context
	VisitDate        string    `json:"visit_date" binding:"required"`
	Purpose          string    `json:"purpose" binding:"required,min=3"`
	Notes            string    `json:"notes" binding:"omitempty"`
	CheckInLocation  *Location `json:"check_in_location" binding:"omitempty"`
	CheckOutLocation *Location `json:"check_out_location" binding:"omitempty"`
	Photos           []string  `json:"photos" binding:"omitempty"`
}

// UpdateVisitReportRequest represents update visit report request DTO
type UpdateVisitReportRequest struct {
	AccountID        *string   `json:"account_id" binding:"omitempty,uuid"` // Optional: required if DealID exists
	ContactID        *string   `json:"contact_id" binding:"omitempty,uuid"`
	DealID           *string   `json:"deal_id" binding:"omitempty,uuid"` // Optional link to deal
	LeadID           *string   `json:"lead_id" binding:"omitempty,uuid"` // Optional link to lead
	VisitDate        string    `json:"visit_date" binding:"omitempty"`
	Purpose          string    `json:"purpose" binding:"omitempty,min=3"`
	Notes            string    `json:"notes" binding:"omitempty"`
	CheckInLocation  *Location `json:"check_in_location" binding:"omitempty"`
	CheckOutLocation *Location `json:"check_out_location" binding:"omitempty"`
	Photos           []string  `json:"photos" binding:"omitempty"`
	Status           string    `json:"status" binding:"omitempty,oneof=draft submitted"`
}

// CheckInRequest represents check-in request DTO
type CheckInRequest struct {
	Location  *Location    `json:"location" binding:"required"`
	PhotoURL  *string      `json:"photo_url,omitempty"`  // Optional: photo URL if uploaded separately
	DeviceGPS *GPSMetadata `json:"device_gps,omitempty"` // Device GPS metadata for validation
	PhotoGPS  *GPSMetadata `json:"photo_gps,omitempty"`  // Photo EXIF GPS metadata for validation
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

// CheckOutRequest represents check-out request DTO
type CheckOutRequest struct {
	Location *Location `json:"location" binding:"required"`
}

// ApproveRequest represents approve request DTO
type ApproveRequest struct {
	// No additional fields needed
}

// RejectRequest represents reject request DTO
type RejectRequest struct {
	Reason string `json:"reason" binding:"required,min=3"`
}

// SubmitRequest represents submit for approval request DTO
type SubmitRequest struct {
	Outcome   string `json:"outcome" binding:"omitempty,oneof=positive neutral negative very_positive"`
	NextSteps string `json:"next_steps" binding:"omitempty"`
}

// UploadPhotoRequest represents photo upload request DTO
type UploadPhotoRequest struct {
	PhotoURL string `json:"photo_url" binding:"required,url"`
}

// ListVisitReportsRequest represents list visit reports query parameters
type ListVisitReportsRequest struct {
	Page          int      `form:"page" binding:"omitempty,min=1"`
	PerPage       int      `form:"per_page" binding:"omitempty,min=1,max=100"`
	Offset        int      `form:"offset" binding:"omitempty,min=0"` // For offset-based pagination (infinity scroll)
	Search        string   `form:"search" binding:"omitempty"`
	Status        string   `form:"status" binding:"omitempty,oneof=draft submitted approved rejected"`
	AccountID     string   `form:"account_id" binding:"omitempty,uuid"`
	DealID        string   `form:"deal_id" binding:"omitempty,uuid"` // Filter by deal
	LeadID        string   `form:"lead_id" binding:"omitempty,uuid"` // Filter by lead
	SalesRepID    string   `form:"sales_rep_id" binding:"omitempty,uuid"`
	BrickID       string   `form:"brick_id" binding:"omitempty,uuid"`
	StartDate     string   `form:"start_date" binding:"omitempty"`
	EndDate       string   `form:"end_date" binding:"omitempty"`
	ScopedUserIDs []string `form:"-" json:"-"` // Injected by scope middleware for team-based filtering
}
