package interfaces

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/permission"
)

// PermissionRepository defines the interface for permission repository
type PermissionRepository interface {
	// FindByID finds a permission by ID
	FindByID(id string) (*permission.Permission, error)
	
	// FindByCode finds a permission by code
	FindByCode(code string) (*permission.Permission, error)
	
	// List returns a list of permissions
	List() ([]permission.Permission, error)
	
	// GetByMenuID returns permissions by menu ID
	GetByMenuID(menuID string) ([]permission.Permission, error)
	
	// GetByRoleID returns permissions for a role
	GetByRoleID(roleID string) ([]permission.Permission, error)
}

// MenuRepository defines the interface for menu repository
type MenuRepository interface {
	// FindByID finds a menu by ID
	FindByID(id string) (*permission.Menu, error)
	
	// FindByURL finds a menu by URL
	FindByURL(url string) (*permission.Menu, error)
	
	// List returns a list of menus (hierarchical)
	List() ([]permission.Menu, error)
	
	// GetRootMenus returns root level menus
	GetRootMenus() ([]permission.Menu, error)
	
	// Create creates a new menu
	Create(menu *permission.Menu) error
	
	// Update updates a menu
	Update(menu *permission.Menu) error
	
	// Delete soft deletes a menu
	Delete(id string) error
}

