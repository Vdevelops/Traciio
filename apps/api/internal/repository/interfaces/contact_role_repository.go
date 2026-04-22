package interfaces

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/contact_role"
)

// ContactRoleRepository defines the interface for contact role repository
type ContactRoleRepository interface {
	// FindByID finds a contact role by ID
	FindByID(id string) (*contact_role.ContactRole, error)
	
	// FindByCode finds a contact role by code
	FindByCode(code string) (*contact_role.ContactRole, error)
	
	// List returns a list of contact roles
	List() ([]contact_role.ContactRole, error)
	
	// Create creates a new contact role
	Create(cr *contact_role.ContactRole) error
	
	// Update updates a contact role
	Update(cr *contact_role.ContactRole) error
	
	// Delete soft deletes a contact role
	Delete(id string) error
}

