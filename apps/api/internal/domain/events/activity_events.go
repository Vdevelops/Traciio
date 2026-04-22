package events

import "time"

// ActivityLoggedEvent is emitted when a new activity is logged
type ActivityLoggedEvent struct {
	ActivityID  string    `json:"activity_id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	AccountID   string    `json:"account_id,omitempty"`
	ContactID   string    `json:"contact_id,omitempty"`
	DealID      string    `json:"deal_id,omitempty"`
	LeadID      string    `json:"lead_id,omitempty"`
	UserID      string    `json:"user_id"`
	Timestamp   time.Time `json:"timestamp"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

// VisitCompletedEvent is emitted when a visit report is completed
type VisitCompletedEvent struct {
	VisitReportID string    `json:"visit_report_id"`
	AccountID     string    `json:"account_id"`
	AccountName   string    `json:"account_name"`
	SalesRepID    string    `json:"sales_rep_id"`
	SalesRepName  string    `json:"sales_rep_name"`
	VisitDate     time.Time `json:"visit_date"`
	Status        string    `json:"status"`
	Notes         string    `json:"notes,omitempty"`
	DealID        string    `json:"deal_id,omitempty"`
	LeadID        string    `json:"lead_id,omitempty"`
	CreatedBy     string    `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
}
