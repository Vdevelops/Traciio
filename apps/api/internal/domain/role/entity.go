package role

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/permission"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Role represents a role entity
type Role struct {
	ID           string                  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name         string                  `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Code         string                  `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Description  string                  `gorm:"type:text" json:"description"`
	Status       string                  `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	MobileAccess bool                    `gorm:"type:boolean;default:false" json:"mobile_access"`
	IsProtected  bool                    `gorm:"type:boolean;default:false" json:"is_protected"` // Protected roles cannot be deleted or disabled
	Permissions  []permission.Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	CreatedAt    time.Time               `json:"created_at"`
	UpdatedAt    time.Time               `json:"updated_at"`
	DeletedAt    gorm.DeletedAt          `gorm:"index" json:"-"`
	// Read-only fields
	UserCount int64 `gorm:"->" json:"user_count"`
}

// TableName specifies the table name for Role
func (Role) TableName() string {
	return "roles"
}

// BeforeCreate hook to generate UUID
func (r *Role) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}

// ScopeType defines the data visibility level for a role on a specific resource
type ScopeType string

const (
	ScopeGlobal ScopeType = "global" // Can see all data across the organization
	ScopeTeam   ScopeType = "team"   // Can see data belonging to the same group/team
	ScopeOwn    ScopeType = "own"    // Can only see data they own or are assigned to
)

// RoleScope defines data visibility rules per role per resource.
// This is the core table for dynamic, policy-driven data scoping.
type RoleScope struct {
	ID        string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	RoleID    string    `gorm:"type:uuid;not null;uniqueIndex:idx_role_scopes_role_resource" json:"role_id"`
	Role      *Role     `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Resource  string    `gorm:"type:varchar(50);not null;uniqueIndex:idx_role_scopes_role_resource" json:"resource"` // leads, deals, tasks, schedules, visit-reports
	Scope     ScopeType `gorm:"type:varchar(20);not null;default:'own'" json:"scope"`                                // global, team, own
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for RoleScope
func (RoleScope) TableName() string {
	return "role_scopes"
}

// BeforeCreate hook to generate UUID
func (rs *RoleScope) BeforeCreate(tx *gorm.DB) error {
	if rs.ID == "" {
		rs.ID = uuid.New().String()
	}
	return nil
}
