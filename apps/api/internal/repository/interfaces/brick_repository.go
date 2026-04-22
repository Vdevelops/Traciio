package interfaces

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/brick"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
)

// BrickRepository defines the interface for brick repository
type BrickRepository interface {
	// FindByID finds a brick by ID
	FindByID(id string) (*brick.Brick, error)
	
	// FindByIDs finds multiple bricks by IDs (batch load)
	FindByIDs(ids []string) ([]brick.Brick, error)
	
	// FindByCode finds a brick by code
	FindByCode(code string) (*brick.Brick, error)
	
	// List returns a list of bricks with pagination
	List(req *brick.ListBricksRequest) ([]brick.Brick, int64, error)
	
	// Create creates a new brick
	Create(brick *brick.Brick) error
	
	// Update updates a brick
	Update(brick *brick.Brick) error
	
	// Delete soft deletes a brick
	Delete(id string) error
	
	// CountSalesByBrickID counts sales users with a specific brick ID
	CountSalesByBrickID(brickID string) (int64, error)
	
	// GetSalesByBrickID gets all sales users by brick ID
	GetSalesByBrickID(brickID string) ([]user.User, error)
	
	// FindByRegencyAndProvince finds brick by regency and province (for map-based creation validation)
	FindByRegencyAndProvince(regency, province string) (*brick.Brick, error)

	// GetNextCodeSequence returns the next available sequence number for a given code prefix
	GetNextCodeSequence(prefix string) (int, error)
}

