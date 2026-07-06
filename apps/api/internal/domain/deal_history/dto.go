package deal_history

import "time"

// DealHistoryResponse represents deal history response DTO.
type DealHistoryResponse struct {
	ID              string    `json:"id"`
	DealID          string    `json:"deal_id"`
	FromStageID     *string   `json:"from_stage_id,omitempty"`
	FromStageName   string    `json:"from_stage_name,omitempty"`
	ToStageID       string    `json:"to_stage_id"`
	ToStageName     string    `json:"to_stage_name"`
	FromProbability int       `json:"from_probability"`
	ToProbability   int       `json:"to_probability"`
	DaysInPrevStage *int      `json:"days_in_prev_stage,omitempty"`
	ChangedBy       string    `json:"changed_by"`
	ChangedAt       time.Time `json:"changed_at"`
	Reason          string    `json:"reason,omitempty"`
	Notes           string    `json:"notes,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ToDealHistoryResponse converts DealHistory to DealHistoryResponse.
func (dh *DealHistory) ToDealHistoryResponse() *DealHistoryResponse {
	return &DealHistoryResponse{
		ID:              dh.ID,
		DealID:          dh.DealID,
		FromStageID:     dh.FromStageID,
		FromStageName:   dh.FromStageName,
		ToStageID:       dh.ToStageID,
		ToStageName:     dh.ToStageName,
		FromProbability: dh.FromProbability,
		ToProbability:   dh.ToProbability,
		DaysInPrevStage: dh.DaysInPrevStage,
		ChangedBy:       dh.ChangedBy,
		ChangedAt:       dh.ChangedAt,
		Reason:          dh.Reason,
		Notes:           dh.Notes,
		CreatedAt:       dh.CreatedAt,
		UpdatedAt:       dh.UpdatedAt,
	}
}
