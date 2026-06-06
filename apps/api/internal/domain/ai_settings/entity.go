package ai_settings

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// AISettings represents AI settings entity
type AISettings struct {
	ID          string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Enabled     bool           `gorm:"type:boolean;not null;default:true" json:"enabled"`
	Provider    string         `gorm:"type:varchar(50);not null;default:'cerebras'" json:"provider"` // cerebras, openai, anthropic, etc.
	APIKey      string         `gorm:"type:text" json:"-"`                                           // Hidden from JSON, stored encrypted
	Model       string         `gorm:"type:varchar(100);not null;default:'llama-3.1-8b'" json:"model"`
	BaseURL     string         `gorm:"type:text" json:"base_url,omitempty"`                     // Optional custom base URL
	DataPrivacy datatypes.JSON `gorm:"type:jsonb" json:"data_privacy"`                          // JSON object with data privacy settings
	Timezone    string         `gorm:"type:varchar(50);default:'Asia/Jakarta'" json:"timezone"` // Timezone for AI context (e.g., "Asia/Jakarta", "UTC", "America/New_York")
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for AISettings
func (AISettings) TableName() string {
	return "ai_settings"
}

// BeforeCreate hook to generate UUID
func (a *AISettings) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}

// DataPrivacySettings represents data privacy settings structure
// Grouped by CRM domain for modular AI context loading
type DataPrivacySettings struct {
	// Sales domain - control AI access to sales data
	AllowLeads        bool `json:"allow_leads"`
	AllowDeals        bool `json:"allow_deals"`
	AllowVisitReports bool `json:"allow_visit_reports"`
	AllowActivities   bool `json:"allow_activities"`
	AllowTasks        bool `json:"allow_tasks"`
	AllowSchedule     bool `json:"allow_schedule"`
	AllowPipelines    bool `json:"allow_pipelines"`

	// Customer domain
	AllowAccounts bool `json:"allow_accounts"`
	AllowContacts bool `json:"allow_contacts"`

	// Inventory domain
	AllowProducts bool `json:"allow_products"`

	// Analytics domain
	AllowSalesPerformance bool `json:"allow_sales_performance"`
	AllowProductAnalysis  bool `json:"allow_product_analysis"`
	AllowReports          bool `json:"allow_reports"`

	// Management domain
	AllowUsers           bool `json:"allow_users"`
	AllowRoles           bool `json:"allow_roles"`
	AllowGroups          bool `json:"allow_groups"`
	AllowBrickManagement bool `json:"allow_brick_management"`
	AllowTarget          bool `json:"allow_target"`

	// Route Optimization domain
	AllowRouteOptimization bool `json:"allow_route_optimization"`
}
