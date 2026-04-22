package interfaces

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/route_optimization"
)

// RouteOptimizationRepository defines the interface for route optimization repository
type RouteOptimizationRepository interface {
	// FindByID finds an optimized route by ID
	FindByID(id string) (*route_optimization.OptimizedRoute, error)

	// List returns a list of optimized routes with pagination
	List(req *route_optimization.ListRoutesRequest) ([]route_optimization.OptimizedRoute, int64, error)

	// Create creates a new optimized route
	Create(route *route_optimization.OptimizedRoute) error

	// Update updates an optimized route
	Update(route *route_optimization.OptimizedRoute) error

	// Delete soft deletes an optimized route
	Delete(id string) error

	// FindByUserID finds optimized routes by user ID
	FindByUserID(userID string) ([]route_optimization.OptimizedRoute, error)
}


