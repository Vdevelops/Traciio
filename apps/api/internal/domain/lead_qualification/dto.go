package lead_qualification

import "time"

// NeedProduct represents a product in the need tracking
type NeedProduct struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
}

// BANTProgress represents the BANT progress for frontend display
type BANTProgress struct {
	Budget    BANTItem `json:"budget"`
	Authority BANTItem `json:"authority"`
	Need      BANTItem `json:"need"`
	Timeline  BANTItem `json:"timeline"`
}

// BANTItem represents a single BANT item status
type BANTItem struct {
	Completed bool `json:"completed"`
	Score     int  `json:"score"`
	MaxScore  int  `json:"max_score"`
}

// LeadQualificationResponse represents qualification response DTO
type LeadQualificationResponse struct {
	ID     string `json:"id"`
	LeadID string `json:"lead_id"`
	BudgetTargetAmount   int64         `json:"budget_target_amount"`
	BudgetTargetCurrency string        `json:"budget_target_currency"`
	BudgetConfirmed      bool          `json:"budget_confirmed"`
	BudgetNotes          string        `json:"budget_notes"`
	AuthorityTargetPerson string       `json:"authority_target_person"`
	AuthorityTargetRole   string       `json:"authority_target_role"`
	AuthorityConfirmed    bool         `json:"authority_confirmed"`
	AuthorityNotes        string       `json:"authority_notes"`
	NeedTargetProducts    []NeedProduct `json:"need_target_products"`
	NeedPriorityLevel     string       `json:"need_priority_level"`
	NeedConfirmed         bool         `json:"need_confirmed"`
	NeedNotes             string       `json:"need_notes"`
	TimelineTargetDate    *time.Time   `json:"timeline_target_date"`
	TimelineFlexibility   string       `json:"timeline_flexibility"`
	TimelineConfirmed     bool         `json:"timeline_confirmed"`
	TimelineNotes         string       `json:"timeline_notes"`
	QualificationScore    int          `json:"qualification_score"`
	QualificationStatus   string       `json:"qualification_status"`
	BANTProgress          BANTProgress `json:"bant_progress"`
	CreatedAt             time.Time    `json:"created_at"`
	UpdatedAt             time.Time    `json:"updated_at"`
}

// UpsertLeadQualificationRequest represents create/update qualification request
type UpsertLeadQualificationRequest struct {
	BudgetTargetAmount    *int64        `json:"budget_target_amount" binding:"omitempty,min=0"`
	BudgetTargetCurrency  string        `json:"budget_target_currency" binding:"omitempty,len=3"`
	BudgetConfirmed       *bool         `json:"budget_confirmed" binding:"omitempty"`
	BudgetNotes           string        `json:"budget_notes" binding:"omitempty"`
	AuthorityTargetPerson string        `json:"authority_target_person" binding:"omitempty,max=255"`
	AuthorityTargetRole   string        `json:"authority_target_role" binding:"omitempty,max=100"`
	AuthorityConfirmed    *bool         `json:"authority_confirmed" binding:"omitempty"`
	AuthorityNotes        string        `json:"authority_notes" binding:"omitempty"`
	NeedTargetProducts    []NeedProduct `json:"need_target_products" binding:"omitempty"`
	NeedPriorityLevel     string        `json:"need_priority_level" binding:"omitempty,oneof=low medium high critical"`
	NeedConfirmed         *bool         `json:"need_confirmed" binding:"omitempty"`
	NeedNotes             string        `json:"need_notes" binding:"omitempty"`
	TimelineTargetDate    *time.Time    `json:"timeline_target_date" binding:"omitempty"`
	TimelineFlexibility   string        `json:"timeline_flexibility" binding:"omitempty,oneof=fixed flexible urgent"`
	TimelineConfirmed     *bool         `json:"timeline_confirmed" binding:"omitempty"`
	TimelineNotes         string        `json:"timeline_notes" binding:"omitempty"`
}
