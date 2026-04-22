package lead_source

import (
	"time"

	"gorm.io/gorm"
)

// LeadSource represents a lead source in the system
type LeadSource struct {
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

// TableName specifies the table name for LeadSource
func (LeadSource) TableName() string {
	return "lead_sources"
}

// ToLeadSourceResponse converts LeadSource to response DTO
func (ls *LeadSource) ToLeadSourceResponse() *LeadSourceResponse {
	return &LeadSourceResponse{
		ID:          ls.ID,
		Name:        ls.Name,
		Code:        ls.Code,
		Description: ls.Description,
		Order:       ls.Order,
		IsActive:    ls.IsActive,
		CreatedAt:   ls.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   ls.UpdatedAt.Format(time.RFC3339),
		LeadCount:   ls.LeadCount,
	}
}

