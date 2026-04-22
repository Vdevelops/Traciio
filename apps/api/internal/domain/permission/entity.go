package permission

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Permission represents a permission entity
type Permission struct {
	ID          string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Code        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"code"` // Format: resource.action
	Description string    `gorm:"type:text" json:"description"`
	Resource    string    `gorm:"type:varchar(50);not null" json:"resource"` // inferred from Code if possible, but explicit is better
	Action      string    `gorm:"type:varchar(50);not null" json:"action"`   // inferred from Code if possible
	
	// Optional UI grouping
	MenuID      *string   `gorm:"type:uuid;index" json:"menu_id,omitempty"`
	Menu        *Menu     `gorm:"foreignKey:MenuID" json:"menu,omitempty"`
	
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
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

// PermissionResponse represents permission response DTO
type PermissionResponse struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Code        string         `json:"code"`
	Resource    string         `json:"resource"`
	Action      string         `json:"action"`
	Description string         `json:"description"`
	Access      bool           `json:"access"` // For role-based permission check
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// PermissionShortResponse represents a concise permission response DTO
type PermissionShortResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Resource string `json:"resource"`
	Action   string `json:"action"`
	Access   bool   `json:"access"`
}

// PermissionSimpleResponse represents minimal permission response for list endpoint
type PermissionSimpleResponse struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Code     string  `json:"code"`
	Resource string  `json:"resource"`
	Action   string  `json:"action"`
	MenuID   *string `json:"menu_id,omitempty"`
	Menu     *Menu   `json:"menu,omitempty"`
}

// ToPermissionResponse converts Permission to PermissionResponse
func (p *Permission) ToPermissionResponse() *PermissionResponse {
	resp := &PermissionResponse{
		ID:          p.ID,
		Name:        p.Name,
		Code:        p.Code,
		Resource:    p.Resource,
		Action:      p.Action,
		Description: p.Description,
		Access:      false, // Default, will be set based on role
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
	return resp
}

// ToPermissionShortResponse converts Permission to PermissionShortResponse
func (p *Permission) ToPermissionShortResponse() *PermissionShortResponse {
	return &PermissionShortResponse{
		ID:       p.ID,
		Name:     p.Name,
		Code:     p.Code,
		Resource: p.Resource,
		Action:   p.Action,
		Access:   true, // Since this is part of role permissions list, the role has access
	}
}

// ToPermissionSimpleResponse converts Permission to PermissionSimpleResponse
func (p *Permission) ToPermissionSimpleResponse() *PermissionSimpleResponse {
	return &PermissionSimpleResponse{
		ID:       p.ID,
		Name:     p.Name,
		Code:     p.Code,
		Resource: p.Resource,
		Action:   p.Action,
		MenuID:   p.MenuID,
		Menu:     p.Menu,
	}
}

// Menu represents a menu entity (hierarchical structure) - kept for UI navigation
type Menu struct {
	ID          string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Icon        string    `gorm:"type:varchar(100)" json:"icon"`
	URL         string    `gorm:"type:varchar(255);not null" json:"url"`
	ParentID    *string   `gorm:"type:uuid;index" json:"parent_id,omitempty"`
	Parent      *Menu     `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []Menu    `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Order       int       `gorm:"type:integer;default:0" json:"order"`
	Status      string    `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
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

// MenuResponse represents menu response DTO (with nested structure)
type MenuResponse struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Icon      string          `json:"icon"`
	URL       string          `json:"url"`
	ParentID  *string         `json:"parent_id,omitempty"`
	Children  []MenuResponse  `json:"children,omitempty"`
	Order     int             `json:"order"`
	Status    string          `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// ToMenuResponse converts Menu to MenuResponse (recursive for children)
func (m *Menu) ToMenuResponse() *MenuResponse {
	resp := &MenuResponse{
		ID:        m.ID,
		Name:      m.Name,
		Icon:      m.Icon,
		URL:       m.URL,
		ParentID:  m.ParentID,
		Order:     m.Order,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
	if len(m.Children) > 0 {
		resp.Children = make([]MenuResponse, len(m.Children))
		for i, child := range m.Children {
			resp.Children[i] = *child.ToMenuResponse()
		}
	}
	return resp
}

