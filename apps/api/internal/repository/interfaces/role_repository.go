package interfaces

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
)

// RoleRepository defines the interface for role repository
type RoleRepository interface {
	// FindByID finds a role by ID
	FindByID(id string) (*role.Role, error)
	
	// FindByCode finds a role by code
	FindByCode(code string) (*role.Role, error)
	
	// List returns a list of roles
	List() ([]role.Role, error)
	
	// Create creates a new role
	Create(role *role.Role) error
	
	// Update updates a role
	Update(role *role.Role) error
	
	// Delete soft deletes a role
	Delete(id string) error
	
	// AssignPermissions assigns permissions to a role
	AssignPermissions(roleID string, permissionIDs []string) error
	
	// GetPermissions gets all permissions for a role
	GetPermissions(roleID string) ([]string, error)
	
	// GetMobilePermissions gets mobile permissions for a role
	GetMobilePermissions(roleID string, role *role.Role) (*role.GetMobilePermissionsResponse, error)
	
	// UpdateMobilePermissions updates mobile permissions for a role
	UpdateMobilePermissions(roleID string, req *role.UpdateMobilePermissionsRequest) error

	// GetScopesByRoleID returns all data scopes for a role
	GetScopesByRoleID(roleID string) ([]role.RoleScope, error)

	// UpsertScopes creates or updates data scopes for a role
	UpsertScopes(roleID string, scopes []role.RoleScopeItem) error
}

