package ai_model_usage

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ModelUsage represents token usage per model
type ModelUsage struct {
	ID          string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Model       string    `gorm:"type:varchar(100);not null;index" json:"model"` // Model name (e.g., "llama-3.1-8b")
	Tokens      int64     `gorm:"type:bigint;not null;default:0" json:"tokens"`   // Total tokens used for this model
	RequestCount int64    `gorm:"type:bigint;not null;default:0" json:"request_count"` // Number of requests
	LastUsedAt  time.Time `gorm:"type:timestamp" json:"last_used_at"`           // Last time this model was used
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for ModelUsage
func (ModelUsage) TableName() string {
	return "ai_model_usage"
}

// BeforeCreate hook to generate UUID
func (m *ModelUsage) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}

// ModelUsageResponse represents model usage response DTO
type ModelUsageResponse struct {
	Model        string    `json:"model"`
	Tokens       int64     `json:"tokens"`
	RequestCount int64     `json:"request_count"`
	LastUsedAt   time.Time `json:"last_used_at"`
}

// ToModelUsageResponse converts ModelUsage to ModelUsageResponse
func (m *ModelUsage) ToModelUsageResponse() *ModelUsageResponse {
	return &ModelUsageResponse{
		Model:        m.Model,
		Tokens:       m.Tokens,
		RequestCount: m.RequestCount,
		LastUsedAt:   m.LastUsedAt,
	}
}
