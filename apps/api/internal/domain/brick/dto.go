package brick

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
)

// BrickResponse represents brick response DTO
type BrickResponse struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Code        string             `json:"code"`
	Description string             `json:"description"`
	Province    string             `json:"province"`
	Regency     string             `json:"regency"`
	District    *string            `json:"district"`
	ManagerID   *string            `json:"manager_id"`
	Manager     *user.UserResponse `json:"manager,omitempty"`
	Status      string             `json:"status"`
	SalesCount  int                `json:"sales_count,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// CreateBrickRequest represents create brick request DTO
type CreateBrickRequest struct {
	Name        string  `json:"name" binding:"required,min=3"`
	Code        string  `json:"code" binding:"omitempty"`
	Description string  `json:"description"`
	Province    string  `json:"province" binding:"required"`
	Regency     string  `json:"regency" binding:"required"`
	District    *string `json:"district"`
	ManagerID   *string `json:"manager_id" binding:"omitempty,uuid"`
	Status      string  `json:"status" binding:"omitempty,oneof=active inactive"`
}

// UpdateBrickRequest represents update brick request DTO
type UpdateBrickRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=3"`
	Code        *string `json:"code" binding:"omitempty,min=3"`
	Description *string `json:"description"`
	Province    *string `json:"province" binding:"omitempty,min=1"`
	Regency     *string `json:"regency" binding:"omitempty,min=1"`
	District    *string `json:"district"`
	ManagerID   *string `json:"manager_id" binding:"omitempty,uuid"`
	Status      *string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// ListBricksRequest represents list bricks query parameters
type ListBricksRequest struct {
	Page      int     `form:"page" binding:"omitempty,min=1"`
	PerPage   int     `form:"per_page" binding:"omitempty,min=1,max=100"`
	Search    string  `form:"search" binding:"omitempty"`
	Province  string  `form:"province" binding:"omitempty"`
	Regency   string  `form:"regency" binding:"omitempty"`
	ManagerID *string `form:"manager_id" binding:"omitempty,uuid"`
	Status    string  `form:"status" binding:"omitempty,oneof=active inactive"`
}

// AssignSalesRequest represents assign/unassign sales users to a brick
type AssignSalesRequest struct {
	UserIDs []string `json:"user_ids" binding:"required,min=1,dive,uuid"`
}
