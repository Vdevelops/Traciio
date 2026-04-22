package role

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/permission"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Role represents a role entity
type Role struct {
	ID          string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Code        string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Description string    `gorm:"type:text" json:"description"`
	Status      string    `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	MobileAccess bool     `gorm:"type:boolean;default:false" json:"mobile_access"`
	IsProtected bool      `gorm:"type:boolean;default:false" json:"is_protected"` // Protected roles cannot be deleted or disabled
	Permissions []permission.Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
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

// RoleResponse represents role response DTO
type RoleResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	MobileAccess bool     `json:"mobile_access"`
	IsProtected bool      `json:"is_protected"`
	Permissions []permission.PermissionSimpleResponse `json:"permissions,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserCount   int64     `json:"user_count"`
}

// ToRoleResponse converts Role to RoleResponse
func (r *Role) ToRoleResponse() *RoleResponse {
	var permissions []permission.PermissionSimpleResponse
	if len(r.Permissions) > 0 {
		permissions = make([]permission.PermissionSimpleResponse, len(r.Permissions))
		for i, p := range r.Permissions {
			permissions[i] = *p.ToPermissionSimpleResponse()
		}
	}

	return &RoleResponse{
		ID:          r.ID,
		Name:        r.Name,
		Code:        r.Code,
		Description: r.Description,
		Status:      r.Status,
		MobileAccess: r.MobileAccess,
		IsProtected: r.IsProtected,
		Permissions: permissions,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
		UserCount:   r.UserCount,
	}
}

// CreateRoleRequest represents create role request DTO
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required,min=3"`
	Code        string `json:"code" binding:"required,min=3"`
	Description string `json:"description"`
	Status      string `json:"status" binding:"omitempty,oneof=active inactive"`
	MobileAccess *bool `json:"mobile_access"`
}

// UpdateRoleRequest represents update role request DTO
type UpdateRoleRequest struct {
	Name        string `json:"name" binding:"omitempty,min=3"`
	Code        string `json:"code" binding:"omitempty,min=3"`
	Description string `json:"description"`
	Status      string `json:"status" binding:"omitempty,oneof=active inactive"`
	MobileAccess *bool `json:"mobile_access"`
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
	ID       string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	RoleID   string    `gorm:"type:uuid;not null;uniqueIndex:idx_role_scopes_role_resource" json:"role_id"`
	Role     *Role     `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Resource string    `gorm:"type:varchar(50);not null;uniqueIndex:idx_role_scopes_role_resource" json:"resource"` // leads, deals, tasks, schedules, visit-reports
	Scope    ScopeType `gorm:"type:varchar(20);not null;default:'own'" json:"scope"`                               // global, team, own
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

// RoleScopeResponse represents role scope response DTO
type RoleScopeResponse struct {
	ID       string    `json:"id"`
	RoleID   string    `json:"role_id"`
	Resource string    `json:"resource"`
	Scope    ScopeType `json:"scope"`
}

// ToResponse converts RoleScope to RoleScopeResponse
func (rs *RoleScope) ToResponse() *RoleScopeResponse {
	return &RoleScopeResponse{
		ID:       rs.ID,
		RoleID:   rs.RoleID,
		Resource: rs.Resource,
		Scope:    rs.Scope,
	}
}

// UpdateRoleScopesRequest represents update role scopes request DTO
type UpdateRoleScopesRequest struct {
	Scopes []RoleScopeItem `json:"scopes" binding:"required,min=1,dive"`
}

// RoleScopeItem represents a single scope assignment in the update request
type RoleScopeItem struct {
	Resource string    `json:"resource" binding:"required"`
	Scope    ScopeType `json:"scope" binding:"required,oneof=global team own"`
}

// AssignPermissionsRequest represents assign permissions to role request DTO
type AssignPermissionsRequest struct {
	PermissionIDs []string `json:"permission_ids" binding:"required,min=1,dive,uuid"`
}

// GetMobilePermissionsResponse represents mobile permissions response for a role
type GetMobilePermissionsResponse struct {
	Menus []MobileMenuPermission `json:"menus"`
}

// MobileMenuPermission represents a mobile menu with CRUD permissions
type MobileMenuPermission struct {
	Menu    string   `json:"menu"`    // dashboard, task, accounts, contacts, visit_reports
	Actions []string `json:"actions"` // VIEW, CREATE, EDIT, DELETE
}

// UpdateMobilePermissionsRequest represents update mobile permissions request DTO
type UpdateMobilePermissionsRequest struct {
	Menus []MobileMenuPermission `json:"menus" binding:"required,min=1,dive"`
}

