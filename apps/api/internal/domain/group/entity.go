package group

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Group represents a group entity
type Group struct {
	ID          string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Code        string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Description string    `gorm:"type:text" json:"description"`
	Status      string    `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Group
func (Group) TableName() string {
	return "groups"
}

// BeforeCreate hook to generate UUID
func (g *Group) BeforeCreate(tx *gorm.DB) error {
	if g.ID == "" {
		g.ID = uuid.New().String()
	}
	return nil
}

// ToGroupResponse converts Group to GroupResponse
func (g *Group) ToGroupResponse() *GroupResponse {
	return &GroupResponse{
		ID:          g.ID,
		Name:        g.Name,
		Code:        g.Code,
		Description: g.Description,
		Status:      g.Status,
		CreatedAt:   g.CreatedAt,
		UpdatedAt:   g.UpdatedAt,
	}
}
