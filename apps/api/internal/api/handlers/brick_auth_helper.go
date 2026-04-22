package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/gilabs/crm-healthcare/api/internal/domain/brick"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
)

// checkBrickAccess checks if user has access to a brick
// Returns:
// - true, nil if user has access (admin or manager of the brick)
// - false, error if user doesn't have access
func checkBrickAccess(c *gin.Context, brickRepo interfaces.BrickRepository, brickID string) (bool, error) {
	// Get user ID and role from context
	userID, exists := c.Get("user_id")
	if !exists {
		return false, errors.New("user ID not found in context")
	}
	userIDStr, ok := userID.(string)
	if !ok {
		return false, errors.New("invalid user ID format")
	}

	userRole, exists := c.Get("user_role")
	if !exists {
		return false, errors.New("user role not found in context")
	}
	roleStr, ok := userRole.(string)
	if !ok {
		return false, errors.New("invalid user role format")
	}

	// Admin and super_admin have full access
	if roleStr == "admin" || roleStr == "super_admin" {
		return true, nil
	}

	// For sales_manager, check if they are the manager of this brick
	if roleStr == "sales_manager" {
		brick, err := brickRepo.FindByID(brickID)
		if err != nil {
			return false, err
		}

		// Check if user is the manager of this brick
		if brick.ManagerID != nil && *brick.ManagerID == userIDStr {
			return true, nil
		}

		// Sales manager doesn't have access to this brick
		return false, nil
	}

	// For sales, check if they are assigned to this brick
	if roleStr == "sales" {
		brick, err := brickRepo.FindByID(brickID)
		if err != nil {
			return false, err
		}

		// Sales can only view their own brick
		// We need to check if user's brick_id matches
		// This requires user repository, so we'll handle it in the handler
		// For now, return false - sales should not access brick analytics directly
		_ = brick
		return false, nil
	}

	// Other roles don't have access
	return false, nil
}

// filterBricksByManager filters bricks list to only show bricks managed by the user
// Returns filtered brick IDs or nil if user is admin (no filter needed)
func filterBricksByManager(c *gin.Context, brickRepo interfaces.BrickRepository) ([]string, error) {
	userRole, exists := c.Get("user_role")
	if !exists {
		return nil, errors.New("user role not found in context")
	}
	roleStr, ok := userRole.(string)
	if !ok {
		return nil, errors.New("invalid user role format")
	}

	// Admin and super_admin can see all bricks
	if roleStr == "admin" || roleStr == "super_admin" {
		return nil, nil // nil means no filter (all bricks)
	}

	// For sales_manager, get only bricks they manage
	if roleStr == "sales_manager" {
		userID, exists := c.Get("user_id")
		if !exists {
			return nil, errors.New("user ID not found in context")
		}
		userIDStr, ok := userID.(string)
		if !ok {
			return nil, errors.New("invalid user ID format")
		}

		// Get all bricks managed by this user
		bricks, _, err := brickRepo.List(&brick.ListBricksRequest{
			ManagerID: &userIDStr,
			PerPage:   1000, // Get all bricks managed by this user
		})
		if err != nil {
			return nil, err
		}

		brickIDs := make([]string, len(bricks))
		for i, brick := range bricks {
			brickIDs[i] = brick.ID
		}

		return brickIDs, nil
	}

	// For sales, they can only see their own brick (via brick_id in users table)
	// This is handled differently - they access via their own brick_id
	return []string{}, nil // Empty array means no access
}

