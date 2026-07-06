package permission

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Permission represents a permission entity
type Permission struct {
	ID          string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	Code        string `gorm:"type:varchar(100);uniqueIndex;not null" json:"code"` // Format: resource.action
	Description string `gorm:"type:text" json:"description"`
	Resource    string `gorm:"type:varchar(50);not null" json:"resource"` // inferred from Code if possible, but explicit is better
	Action      string `gorm:"type:varchar(50);not null" json:"action"`   // inferred from Code if possible

	// Optional UI grouping
	MenuID *string `gorm:"type:uuid;index" json:"menu_id,omitempty"`
	Menu   *Menu   `gorm:"foreignKey:MenuID" json:"menu,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Permission
func (Permission) TableName() string {
	return "permissions"
}

// BeforeCreate hook to generate UUID
func (p *Permission) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

// Menu represents a menu entity (hierarchical structure) - kept for UI navigation
type Menu struct {
	ID        string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Icon      string         `gorm:"type:varchar(100)" json:"icon"`
	URL       string         `gorm:"type:varchar(255);not null" json:"url"`
	ParentID  *string        `gorm:"type:uuid;index" json:"parent_id,omitempty"`
	Parent    *Menu          `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children  []Menu         `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Order     int            `gorm:"type:integer;default:0" json:"order"`
	Status    string         `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Menu
func (Menu) TableName() string {
	return "menus"
}

// BeforeCreate hook to generate UUID
func (m *Menu) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}
