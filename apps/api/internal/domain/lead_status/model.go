package lead_status

import (
	"time"

	"gorm.io/gorm"
)

// LeadStatus represents a lead status in the system
type LeadStatus struct {
	ID          string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	Code        string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"code"`
	Description string         `gorm:"type:text" json:"description"`
	Score       int            `gorm:"type:int;not null;default:0" json:"score"` // Score percentage (0-100)
	Color       string         `gorm:"type:varchar(20);default:'#3B82F6'" json:"color"`
	Order       int            `gorm:"type:int;not null;default:0" json:"order"`
	IsActive    bool           `gorm:"type:boolean;not null;default:true" json:"is_active"`
	IsDefault   bool           `gorm:"type:boolean;not null;default:false" json:"is_default"`   // Default status for new leads
	IsConverted bool           `gorm:"type:boolean;not null;default:false" json:"is_converted"` // Mark as converted status
	CreatedBy   string         `gorm:"type:uuid" json:"created_by"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	// Read-only fields
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	LeadCount int64          `gorm:"->" json:"lead_count"`
}

// TableName specifies the table name for LeadStatus
func (LeadStatus) TableName() string {
	return "lead_statuses"
}

// ToLeadStatusResponse converts LeadStatus to response DTO
func (ls *LeadStatus) ToLeadStatusResponse() *LeadStatusResponse {
	return &LeadStatusResponse{
		ID:          ls.ID,
		Name:        ls.Name,
		Code:        ls.Code,
		Description: ls.Description,
		Score:       ls.Score,
		Color:       ls.Color,
		Order:       ls.Order,
		IsActive:    ls.IsActive,
		IsDefault:   ls.IsDefault,
		IsConverted: ls.IsConverted,
		CreatedAt:   ls.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   ls.UpdatedAt.Format(time.RFC3339),
		LeadCount:   ls.LeadCount,
	}
}
