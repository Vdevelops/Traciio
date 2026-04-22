package events

import "time"

// LeadCreatedEvent is emitted when a new lead is created
type LeadCreatedEvent struct {
	LeadID      string    `json:"lead_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Company     string    `json:"company"`
	LeadStatus  string    `json:"lead_status"`
	LeadSource  string    `json:"lead_source"`
	AssignedTo  string    `json:"assigned_to"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

// LeadConvertedEvent is emitted when a lead is converted to an opportunity
type LeadConvertedEvent struct {
	LeadID        string    `json:"lead_id"`
	OpportunityID string    `json:"opportunity_id"`
	AccountID     string    `json:"account_id"`
	ContactID     string    `json:"contact_id"`
	ConvertedBy   string    `json:"converted_by"`
	ConvertedAt   time.Time `json:"converted_at"`
}

// LeadStatusChangedEvent is emitted when a lead's status changes
type LeadStatusChangedEvent struct {
	LeadID      string    `json:"lead_id"`
	OldStatus   string    `json:"old_status"`
	NewStatus   string    `json:"new_status"`
	ChangedBy   string    `json:"changed_by"`
	ChangedAt   time.Time `json:"changed_at"`
	Reason      string    `json:"reason,omitempty"`
}
