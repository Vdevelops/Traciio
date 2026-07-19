package visit_report

import (
	"encoding/json"
	"time"
)

// VisitReportResponse represents visit report response DTO.
type VisitReportResponse struct {
	ID               string      `json:"id"`
	AccountID        *string     `json:"account_id,omitempty"`
	ContactID        *string     `json:"contact_id,omitempty"`
	DealID           *string     `json:"deal_id,omitempty"`
	LeadID           *string     `json:"lead_id,omitempty"`
	Type             string      `json:"type"`
	SalesRepID       string      `json:"sales_rep_id"`
	BrickID          *string     `json:"brick_id,omitempty"`
	VisitDate        time.Time   `json:"visit_date"`
	CheckInTime      *time.Time  `json:"check_in_time,omitempty"`
	CheckOutTime     *time.Time  `json:"check_out_time,omitempty"`
	CheckInLocation  *Location   `json:"check_in_location,omitempty"`
	CheckOutLocation *Location   `json:"check_out_location,omitempty"`
	Purpose          string      `json:"purpose"`
	Notes            string      `json:"notes"`
	Outcome          string      `json:"outcome,omitempty"`
	NextSteps        string      `json:"next_steps,omitempty"`
	Photos           []string    `json:"photos,omitempty"`
	Metadata         interface{} `json:"metadata,omitempty"`
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

// ToVisitReportResponse converts VisitReport to VisitReportResponse.
func (vr *VisitReport) ToVisitReportResponse() *VisitReportResponse {
	var photos []string
	if vr.Photos != nil {
		_ = json.Unmarshal(vr.Photos, &photos)
	}
	var metadata interface{}
	if vr.Metadata != nil {
		_ = json.Unmarshal(vr.Metadata, &metadata)
	}
	var checkInLocation *Location
	if vr.CheckInLocation != nil {
		var location Location
		if err := json.Unmarshal(vr.CheckInLocation, &location); err == nil {
			checkInLocation = &location
		}
	}
	var checkOutLocation *Location
	if vr.CheckOutLocation != nil {
		var location Location
		if err := json.Unmarshal(vr.CheckOutLocation, &location); err == nil {
			checkOutLocation = &location
		}
	}

	return &VisitReportResponse{
		ID:               vr.ID,
		AccountID:        vr.AccountID,
		ContactID:        vr.ContactID,
		DealID:           vr.DealID,
		LeadID:           vr.LeadID,
		SalesRepID:       vr.SalesRepID,
		BrickID:          vr.BrickID,
		VisitDate:        vr.VisitDate,
		CheckInTime:      vr.CheckInTime,
		CheckOutTime:     vr.CheckOutTime,
		CheckInLocation:  checkInLocation,
		CheckOutLocation: checkOutLocation,
		Purpose:          vr.Purpose,
		Notes:            vr.Notes,
		Outcome:          vr.Outcome,
		NextSteps:        vr.NextSteps,
		Photos:           photos,
		Metadata:         metadata,
		Status:           NormalizeStatus(vr.Status),
		ApprovedBy:       vr.ApprovedBy,
		ApprovedAt:       vr.ApprovedAt,
		RejectionReason:  vr.RejectionReason,
		CreatedAt:        vr.CreatedAt,
		UpdatedAt:        vr.UpdatedAt,
		Account:          vr.Account,
		Contact:          vr.Contact,
		Deal:             vr.Deal,
		SalesRep:         vr.SalesRep,
	}
}

type CreateVisitReportRequest struct {
	AccountID        *string     `json:"account_id" binding:"omitempty,uuid"`
	ContactID        *string     `json:"contact_id" binding:"omitempty,uuid"`
	DealID           *string     `json:"deal_id" binding:"omitempty,uuid"`
	LeadID           *string     `json:"lead_id" binding:"omitempty,uuid"`
	SalesRepID       string      `json:"sales_rep_id" binding:"omitempty,uuid"`
	VisitDate        string      `json:"visit_date" binding:"required"`
	Purpose          string      `json:"purpose" binding:"required,min=3"`
	Notes            string      `json:"notes" binding:"omitempty"`
	CheckInLocation  *Location   `json:"check_in_location" binding:"omitempty"`
	CheckOutLocation *Location   `json:"check_out_location" binding:"omitempty"`
	Photos           []string    `json:"photos" binding:"omitempty"`
	Metadata         interface{} `json:"metadata" binding:"omitempty"`
}

type UpdateVisitReportRequest struct {
	AccountID        *string     `json:"account_id" binding:"omitempty,uuid"`
	ContactID        *string     `json:"contact_id" binding:"omitempty,uuid"`
	DealID           *string     `json:"deal_id" binding:"omitempty,uuid"`
	LeadID           *string     `json:"lead_id" binding:"omitempty,uuid"`
	VisitDate        string      `json:"visit_date" binding:"omitempty"`
	Purpose          string      `json:"purpose" binding:"omitempty,min=3"`
	Notes            string      `json:"notes" binding:"omitempty"`
	CheckInLocation  *Location   `json:"check_in_location" binding:"omitempty"`
	CheckOutLocation *Location   `json:"check_out_location" binding:"omitempty"`
	Photos           []string    `json:"photos" binding:"omitempty"`
	Metadata         interface{} `json:"metadata" binding:"omitempty"`
	Status           string      `json:"status" binding:"omitempty,oneof=pending completed draft submitted approved rejected"`
}

type CheckInRequest struct {
	Location  *Location    `json:"location" binding:"required"`
	PhotoURL  *string      `json:"photo_url,omitempty"`
	DeviceGPS *GPSMetadata `json:"device_gps,omitempty"`
	PhotoGPS  *GPSMetadata `json:"photo_gps,omitempty"`
}

type CheckOutRequest struct {
	Location *Location `json:"location" binding:"required"`
}

type ApproveRequest struct{}

type RejectRequest struct {
	Reason string `json:"reason" binding:"required,min=3"`
}

type SubmitRequest struct {
	Outcome   string `json:"outcome" binding:"omitempty,oneof=positive neutral negative very_positive"`
	NextSteps string `json:"next_steps" binding:"omitempty"`
}

type UploadPhotoRequest struct {
	PhotoURL string `json:"photo_url" binding:"required,url"`
}

type ListVisitReportsRequest struct {
	Page          int      `form:"page" binding:"omitempty,min=1"`
	PerPage       int      `form:"per_page" binding:"omitempty,min=1,max=100"`
	Offset        int      `form:"offset" binding:"omitempty,min=0"`
	Search        string   `form:"search" binding:"omitempty"`
	Status        string   `form:"status" binding:"omitempty,oneof=pending completed draft submitted approved rejected"`
	AccountID     string   `form:"account_id" binding:"omitempty,uuid"`
	DealID        string   `form:"deal_id" binding:"omitempty,uuid"`
	LeadID        string   `form:"lead_id" binding:"omitempty,uuid"`
	SalesRepID    string   `form:"sales_rep_id" binding:"omitempty,uuid"`
	BrickID       string   `form:"brick_id" binding:"omitempty,uuid"`
	StartDate     string   `form:"start_date" binding:"omitempty"`
	EndDate       string   `form:"end_date" binding:"omitempty"`
	ScopedUserIDs []string `form:"-" json:"-"`
}
