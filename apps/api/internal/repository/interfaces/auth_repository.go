package interfaces

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
)

// AuthRepository defines the interface for auth repository
type AuthRepository interface {
	// FindByEmail finds a user by email
	FindByEmail(email string) (*user.User, error)
	
	// FindByID finds a user by ID
	FindByID(id string) (*user.User, error)
	
	// Create creates a new user
	Create(user *user.User) error
	
	// Update updates a user
	Update(user *user.User) error
}

