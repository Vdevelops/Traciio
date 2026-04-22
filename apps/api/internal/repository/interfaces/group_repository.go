package interfaces

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/group"
)

// GroupRepository defines the interface for group repository
type GroupRepository interface {
	// FindByID finds a group by ID
	FindByID(id string) (*group.Group, error)
	
	// FindByCode finds a group by code
	FindByCode(code string) (*group.Group, error)
	
	// List returns a list of groups with pagination
	List(req *group.ListGroupsRequest) ([]group.Group, int64, error)
	
	// Create creates a new group
	Create(group *group.Group) error
	
	// Update updates a group
	Update(group *group.Group) error
	
	// Delete soft deletes a group
	Delete(id string) error
	
	// CountUsersByGroupID counts users with a specific group ID
	CountUsersByGroupID(groupID string) (int64, error)
}

