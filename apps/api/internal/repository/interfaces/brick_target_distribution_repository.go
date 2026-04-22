package interfaces

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/brick_target_distribution"
)

// BrickTargetDistributionRepository defines the interface for brick target distribution repository
type BrickTargetDistributionRepository interface {
	// FindByID finds a brick target distribution by ID
	FindByID(id string) (*brick_target_distribution.BrickTargetDistribution, error)
	
	// List returns a list of brick target distributions with pagination
	List(req *brick_target_distribution.ListBrickTargetDistributionsRequest) ([]brick_target_distribution.BrickTargetDistribution, int64, error)
	
	// FindByBrickTargetID finds distributions by brick target ID
	FindByBrickTargetID(brickTargetID string) ([]brick_target_distribution.BrickTargetDistribution, error)
	
	// FindBySalesUserIDAndBrickTargetID finds distribution by sales user ID and brick target ID (for uniqueness check)
	FindBySalesUserIDAndBrickTargetID(salesUserID, brickTargetID string) (*brick_target_distribution.BrickTargetDistribution, error)
	
	// Create creates a new brick target distribution
	Create(distribution *brick_target_distribution.BrickTargetDistribution) error
	
	// Update updates a brick target distribution
	Update(distribution *brick_target_distribution.BrickTargetDistribution) error
	
	// Delete soft deletes a brick target distribution
	Delete(id string) error
	
	// BulkCreate creates multiple brick target distributions
	BulkCreate(distributions []*brick_target_distribution.BrickTargetDistribution) error
}

