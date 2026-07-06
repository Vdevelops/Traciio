package industry

import (
	"time"

	"gorm.io/gorm"
)

// Industry represents an industry in the system
type Industry struct {
	ID          string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	Code        string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"code"`
	Description string         `gorm:"type:text" json:"description"`
	Order       int            `gorm:"type:int;not null;default:0" json:"order"`
	IsActive    bool           `gorm:"type:boolean;not null;default:true" json:"is_active"`
	CreatedBy   string         `gorm:"type:uuid" json:"created_by"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	// Read-only fields
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	LeadCount int64          `gorm:"->" json:"lead_count"`
}

// TableName specifies the table name for Industry
func (Industry) TableName() string {
	return "industries"
}

// ToIndustryResponse converts Industry to response DTO
func (i *Industry) ToIndustryResponse() *IndustryResponse {
	return &IndustryResponse{
		ID:          i.ID,
		Name:        i.Name,
		Code:        i.Code,
		Description: i.Description,
		Order:       i.Order,
		IsActive:    i.IsActive,
		CreatedAt:   i.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   i.UpdatedAt.Format(time.RFC3339),
		LeadCount:   i.LeadCount,
	}
}
