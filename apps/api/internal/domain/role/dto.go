package role

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/permission"
)

// RoleResponse represents role response DTO.
type RoleResponse struct {
	ID           string                                `json:"id"`
	Name         string                                `json:"name"`
	Code         string                                `json:"code"`
	Description  string                                `json:"description"`
	Status       string                                `json:"status"`
	MobileAccess bool                                  `json:"mobile_access"`
	IsProtected  bool                                  `json:"is_protected"`
	Permissions  []permission.PermissionSimpleResponse `json:"permissions,omitempty"`
	CreatedAt    time.Time                             `json:"created_at"`
	UpdatedAt    time.Time                             `json:"updated_at"`
	UserCount    int64                                 `json:"user_count"`
}

// ToRoleResponse converts Role to RoleResponse.
func (r *Role) ToRoleResponse() *RoleResponse {
	var permissions []permission.PermissionSimpleResponse
	if len(r.Permissions) > 0 {
		permissions = make([]permission.PermissionSimpleResponse, len(r.Permissions))
		for i, p := range r.Permissions {
			permissions[i] = *p.ToPermissionSimpleResponse()
		}
	}

	return &RoleResponse{
		ID:           r.ID,
		Name:         r.Name,
		Code:         r.Code,
		Description:  r.Description,
		Status:       r.Status,
		MobileAccess: r.MobileAccess,
		IsProtected:  r.IsProtected,
		Permissions:  permissions,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
		UserCount:    r.UserCount,
	}
}

// CreateRoleRequest represents create role request DTO.
type CreateRoleRequest struct {
	Name         string `json:"name" binding:"required,min=3"`
	Code         string `json:"code" binding:"required,min=3"`
	Description  string `json:"description"`
	Status       string `json:"status" binding:"omitempty,oneof=active inactive"`
	MobileAccess *bool  `json:"mobile_access"`
}

// UpdateRoleRequest represents update role request DTO.
type UpdateRoleRequest struct {
	Name         string `json:"name" binding:"omitempty,min=3"`
	Code         string `json:"code" binding:"omitempty,min=3"`
	Description  string `json:"description"`
	Status       string `json:"status" binding:"omitempty,oneof=active inactive"`
	MobileAccess *bool  `json:"mobile_access"`
}

// RoleScopeResponse represents role scope response DTO.
type RoleScopeResponse struct {
	ID       string    `json:"id"`
	RoleID   string    `json:"role_id"`
	Resource string    `json:"resource"`
	Scope    ScopeType `json:"scope"`
}

// ToResponse converts RoleScope to RoleScopeResponse.
func (rs *RoleScope) ToResponse() *RoleScopeResponse {
	return &RoleScopeResponse{
		ID:       rs.ID,
		RoleID:   rs.RoleID,
		Resource: rs.Resource,
		Scope:    rs.Scope,
	}
}

// UpdateRoleScopesRequest represents update role scopes request DTO.
type UpdateRoleScopesRequest struct {
	Scopes []RoleScopeItem `json:"scopes" binding:"required,min=1,dive"`
}

// RoleScopeItem represents a single scope assignment in the update request.
type RoleScopeItem struct {
	Resource string    `json:"resource" binding:"required"`
	Scope    ScopeType `json:"scope" binding:"required,oneof=global team own"`
}

// AssignPermissionsRequest represents assign permissions to role request DTO.
type AssignPermissionsRequest struct {
	PermissionIDs []string `json:"permission_ids" binding:"required,min=1,dive,uuid"`
}

// GetMobilePermissionsResponse represents mobile permissions response for a role.
type GetMobilePermissionsResponse struct {
	Menus []MobileMenuPermission `json:"menus"`
}

// MobileMenuPermission represents a mobile menu with CRUD permissions.
type MobileMenuPermission struct {
	Menu    string   `json:"menu"`
	Actions []string `json:"actions"`
}

// UpdateMobilePermissionsRequest represents update mobile permissions request DTO.
type UpdateMobilePermissionsRequest struct {
	Menus []MobileMenuPermission `json:"menus" binding:"required,min=1,dive"`
}
