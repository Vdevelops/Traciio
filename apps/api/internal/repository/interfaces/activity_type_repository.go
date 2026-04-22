package interfaces

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/activity_type"
)

// ActivityTypeRepository defines the interface for activity type repository
type ActivityTypeRepository interface {
	// FindByID finds an activity type by ID
	FindByID(id string) (*activity_type.ActivityType, error)

	// FindByCode finds an activity type by code
	FindByCode(code string) (*activity_type.ActivityType, error)

	// List returns a list of activity types
	List(req *activity_type.ListActivityTypesRequest) ([]activity_type.ActivityType, error)

	// Create creates a new activity type
	Create(at *activity_type.ActivityType) error

	// Update updates an activity type
	Update(at *activity_type.ActivityType) error

	// Delete soft deletes an activity type
	Delete(id string) error
}

