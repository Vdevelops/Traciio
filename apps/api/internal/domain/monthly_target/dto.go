package monthly_target

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/brick"
	"github.com/gilabs/crm-healthcare/api/internal/domain/group"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
)

// MonthlyTargetResponse represents monthly target response DTO
type MonthlyTargetResponse struct {
	ID           string               `json:"id"`
	GroupID      *string              `json:"group_id"`
	Group        *group.GroupResponse `json:"group,omitempty"`
	UserID       *string              `json:"user_id"`
	User         *user.UserResponse   `json:"user,omitempty"`
	BrickID      *string              `json:"brick_id"`
	Brick        *brick.BrickResponse `json:"brick,omitempty"`
	Year         int                  `json:"year"`
	Month        int                  `json:"month"`
	TargetAmount int64                `json:"target_amount"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
}

// CreateMonthlyTargetRequest represents create monthly target request DTO
type CreateMonthlyTargetRequest struct {
	GroupID      *string `json:"group_id" binding:"omitempty,uuid"`
	UserID       *string `json:"user_id" binding:"omitempty,uuid"`
	BrickID      *string `json:"brick_id" binding:"omitempty,uuid"`
	Year         int     `json:"year" binding:"required,min=2000,max=2100"`
	Month        int     `json:"month" binding:"required,min=1,max=12"`
	TargetAmount int64   `json:"target_amount" binding:"required,min=0"`
}

// BulkCreateMonthlyTargetRequest represents bulk create monthly target request DTO
type BulkCreateMonthlyTargetRequest struct {
	GroupIDs     []string `json:"group_ids" binding:"omitempty,dive,uuid"`
	UserIDs      []string `json:"user_ids" binding:"omitempty,dive,uuid"`
	BrickIDs     []string `json:"brick_ids" binding:"omitempty,dive,uuid"`
	Year         int      `json:"year" binding:"required,min=2000,max=2100"`
	Month        int      `json:"month" binding:"required,min=1,max=12"`
	TargetAmount int64    `json:"target_amount" binding:"required,min=0"`
}

// UpdateMonthlyTargetRequest represents update monthly target request DTO
type UpdateMonthlyTargetRequest struct {
	Year         *int   `json:"year" binding:"omitempty,min=2000,max=2100"`
	Month        *int   `json:"month" binding:"omitempty,min=1,max=12"`
	TargetAmount *int64 `json:"target_amount" binding:"omitempty,min=0"`
}

// ListMonthlyTargetsRequest represents list monthly targets query parameters
type ListMonthlyTargetsRequest struct {
	Page          int      `form:"page" binding:"omitempty,min=1"`
	PerPage       int      `form:"per_page" binding:"omitempty,min=1,max=1000"`
	GroupID       *string  `form:"group_id" binding:"omitempty,uuid"`
	UserID        *string  `form:"user_id" binding:"omitempty,uuid"`
	BrickID       *string  `form:"brick_id" binding:"omitempty,uuid"`
	Year          *int     `form:"year" binding:"omitempty,min=2000,max=2100"`
	Month         *int     `form:"month" binding:"omitempty,min=1,max=12"`
	Scope         string   `form:"scope" binding:"omitempty,oneof=all user group brick"`
	ScopedUserIDs []string `form:"-" json:"-"`
}

// GetUserTargetRequest represents get user target request (with fallback to group)
type GetUserTargetRequest struct {
	UserID string `form:"user_id" binding:"required,uuid"`
	Year   int    `form:"year" binding:"required,min=2000,max=2100"`
	Month  int    `form:"month" binding:"required,min=1,max=12"`
}

// CreateGroupTargetWithUserAssignmentRequest represents create group target with auto-assign to all users request DTO
type CreateGroupTargetWithUserAssignmentRequest struct {
	GroupID      string `json:"group_id" binding:"required,uuid"`
	Year         int    `json:"year" binding:"required,min=2000,max=2100"`
	Month        int    `json:"month" binding:"required,min=1,max=12"`
	TargetAmount int64  `json:"target_amount" binding:"required,min=0"`
}

// BulkSetTargetRequest represents bulk set target request DTO (for a range of months)
type BulkSetTargetRequest struct {
	GroupID      *string `json:"group_id" binding:"omitempty,uuid"`
	UserID       *string `json:"user_id" binding:"omitempty,uuid"`
	BrickID      *string `json:"brick_id" binding:"omitempty,uuid"`
	Year         int     `json:"year" binding:"required,min=2000,max=2100"`
	StartMonth   int     `json:"start_month" binding:"required,min=1,max=12"`
	EndMonth     int     `json:"end_month" binding:"required,min=1,max=12"`
	TargetAmount int64   `json:"target_amount" binding:"required,min=0"`
}
