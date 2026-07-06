package brick_target_distribution

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/brick"
	"github.com/gilabs/crm-healthcare/api/internal/domain/monthly_target"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
)

// BrickTargetDistributionResponse represents brick target distribution response DTO.
type BrickTargetDistributionResponse struct {
	ID                         string                                `json:"id"`
	BrickID                    string                                `json:"brick_id"`
	Brick                      *brick.BrickResponse                  `json:"brick,omitempty"`
	BrickTargetID              string                                `json:"brick_target_id"`
	BrickTarget                *monthly_target.MonthlyTargetResponse `json:"brick_target,omitempty"`
	SalesUserID                string                                `json:"sales_user_id"`
	SalesUser                  *user.UserResponse                    `json:"sales_user,omitempty"`
	DistributedAmount          int64                                 `json:"distributed_amount"`
	DistributedAmountFormatted string                                `json:"distributed_amount_formatted,omitempty"`
	DistributedBy              string                                `json:"distributed_by"`
	DistributedByUser          *user.UserResponse                    `json:"distributed_by_user,omitempty"`
	DistributedAt              time.Time                             `json:"distributed_at"`
	CreatedAt                  time.Time                             `json:"created_at"`
	UpdatedAt                  time.Time                             `json:"updated_at"`
}

// ToBrickTargetDistributionResponse converts BrickTargetDistribution to DTO.
func (btd *BrickTargetDistribution) ToBrickTargetDistributionResponse() *BrickTargetDistributionResponse {
	resp := &BrickTargetDistributionResponse{
		ID:                btd.ID,
		BrickID:           btd.BrickID,
		BrickTargetID:     btd.BrickTargetID,
		SalesUserID:       btd.SalesUserID,
		DistributedAmount: btd.DistributedAmount,
		DistributedBy:     btd.DistributedBy,
		DistributedAt:     btd.DistributedAt,
		CreatedAt:         btd.CreatedAt,
		UpdatedAt:         btd.UpdatedAt,
	}
	if btd.Brick != nil {
		resp.Brick = btd.Brick.ToBrickResponse()
	}
	if btd.BrickTarget != nil {
		resp.BrickTarget = btd.BrickTarget.ToMonthlyTargetResponse()
	}
	if btd.SalesUser != nil {
		resp.SalesUser = btd.SalesUser.ToUserResponse()
	}
	if btd.DistributedByUser != nil {
		resp.DistributedByUser = btd.DistributedByUser.ToUserResponse()
	}
	return resp
}

// CreateBrickTargetDistributionRequest represents create input DTO.
type CreateBrickTargetDistributionRequest struct {
	SalesUserID       string `json:"sales_user_id" binding:"required,uuid"`
	DistributedAmount int64  `json:"distributed_amount" binding:"required,min=0"`
}

// BulkCreateBrickTargetDistributionRequest represents bulk create input DTO.
type BulkCreateBrickTargetDistributionRequest struct {
	Distributions []CreateBrickTargetDistributionRequest `json:"distributions" binding:"required,min=1,dive"`
}

// UpdateBrickTargetDistributionRequest represents update input DTO.
type UpdateBrickTargetDistributionRequest struct {
	DistributedAmount *int64 `json:"distributed_amount" binding:"omitempty,min=0"`
}

// ListBrickTargetDistributionsRequest represents list query DTO.
type ListBrickTargetDistributionsRequest struct {
	Page          int     `form:"page" binding:"omitempty,min=1"`
	PerPage       int     `form:"per_page" binding:"omitempty,min=1,max=100"`
	BrickID       *string `form:"brick_id" binding:"omitempty,uuid"`
	BrickTargetID *string `form:"brick_target_id" binding:"omitempty,uuid"`
	SalesUserID   *string `form:"sales_user_id" binding:"omitempty,uuid"`
}
