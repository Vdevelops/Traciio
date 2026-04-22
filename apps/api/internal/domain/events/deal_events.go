package events

import "time"

// DealCreatedEvent is emitted when a new deal is created
type DealCreatedEvent struct {
	DealID            string    `json:"deal_id"`
	Title             string    `json:"title"`
	Value             int64     `json:"value"`
	AccountID         string    `json:"account_id"`
	ContactID         string    `json:"contact_id"`
	StageID           string    `json:"stage_id"`
	StageName         string    `json:"stage_name"`
	PipelineID        string    `json:"pipeline_id"`
	AssignedTo        string    `json:"assigned_to"`
	ExpectedCloseDate *time.Time `json:"expected_close_date,omitempty"`
	CreatedBy         string    `json:"created_by"`
	CreatedAt         time.Time `json:"created_at"`
}

// DealStageChangedEvent is emitted when a deal moves to a different stage
type DealStageChangedEvent struct {
	DealID       string    `json:"deal_id"`
	OldStageID   string    `json:"old_stage_id"`
	OldStageName string    `json:"old_stage_name"`
	NewStageID   string    `json:"new_stage_id"`
	NewStageName string    `json:"new_stage_name"`
	ChangedBy    string    `json:"changed_by"`
	ChangedAt    time.Time `json:"changed_at"`
}

// DealWonEvent is emitted when a deal is marked as won
type DealWonEvent struct {
	DealID          string    `json:"deal_id"`
	Title           string    `json:"title"`
	Value           int64     `json:"value"`
	AccountID       string    `json:"account_id"`
	AssignedTo      string    `json:"assigned_to"`
	ActualCloseDate time.Time `json:"actual_close_date"`
	WonBy           string    `json:"won_by"`
	WonAt           time.Time `json:"won_at"`
}

// DealLostEvent is emitted when a deal is marked as lost
type DealLostEvent struct {
	DealID          string    `json:"deal_id"`
	Title           string    `json:"title"`
	Value           int64     `json:"value"`
	AccountID       string    `json:"account_id"`
	AssignedTo      string    `json:"assigned_to"`
	ActualCloseDate time.Time `json:"actual_close_date"`
	LostReason      string    `json:"lost_reason,omitempty"`
	LostBy          string    `json:"lost_by"`
	LostAt          time.Time `json:"lost_at"`
}
