package brick

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Brick represents a brick/area entity
type Brick struct {
	ID          string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Code        string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Description string    `gorm:"type:text" json:"description"`
	Province    string    `gorm:"type:varchar(100);not null" json:"province"`
	Regency     string    `gorm:"type:varchar(100);not null" json:"regency"`
	District    *string   `gorm:"type:varchar(100)" json:"district"`
	ManagerID   *string   `gorm:"type:uuid;index" json:"manager_id"`
	Manager     *user.User `gorm:"foreignKey:ManagerID" json:"manager,omitempty"`
	Status      string    `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Brick
func (Brick) TableName() string {
	return "bricks"
}

// BeforeCreate hook to generate UUID
func (b *Brick) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return nil
}

// BrickResponse represents brick response DTO
type BrickResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Province    string    `json:"province"`
	Regency     string    `json:"regency"`
	District    *string   `json:"district"`
	ManagerID   *string   `json:"manager_id"`
	Manager     *user.UserResponse `json:"manager,omitempty"`
	Status      string    `json:"status"`
	SalesCount  int       `json:"sales_count,omitempty"` // Count of sales in this brick
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToBrickResponse converts Brick to BrickResponse
func (b *Brick) ToBrickResponse() *BrickResponse {
	resp := &BrickResponse{
		ID:          b.ID,
		Name:        b.Name,
		Code:        b.Code,
		Description: b.Description,
		Province:    b.Province,
		Regency:     b.Regency,
		District:    b.District,
		ManagerID:   b.ManagerID,
		Status:      b.Status,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}
	if b.Manager != nil {
		resp.Manager = b.Manager.ToUserResponse()
	}
	return resp
}

// CreateBrickRequest represents create brick request DTO
// Code is optional: if empty, the server auto-generates a unique code
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
	Page       int     `form:"page" binding:"omitempty,min=1"`
	PerPage    int     `form:"per_page" binding:"omitempty,min=1,max=100"`
	Search     string  `form:"search" binding:"omitempty"`
	Province   string  `form:"province" binding:"omitempty"`
	Regency    string  `form:"regency" binding:"omitempty"`
	ManagerID  *string `form:"manager_id" binding:"omitempty,uuid"`
	Status     string  `form:"status" binding:"omitempty,oneof=active inactive"`
}

// AssignSalesRequest represents assign/unassign sales users to a brick
type AssignSalesRequest struct {
	UserIDs []string `json:"user_ids" binding:"required,min=1,dive,uuid"`
}

