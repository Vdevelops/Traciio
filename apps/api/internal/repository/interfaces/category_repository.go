package interfaces

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/category"
)

// CategoryRepository defines the interface for category repository
type CategoryRepository interface {
	// FindByID finds a category by ID
	FindByID(id string) (*category.Category, error)
	
	// FindByCode finds a category by code
	FindByCode(code string) (*category.Category, error)
	
	// List returns a list of categories
	List() ([]category.Category, error)
	
	// Create creates a new category
	Create(cat *category.Category) error
	
	// Update updates a category
	Update(cat *category.Category) error
	
	// Delete soft deletes a category
	Delete(id string) error
}

