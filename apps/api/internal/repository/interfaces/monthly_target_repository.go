package interfaces

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/monthly_target"
)

// MonthlyTargetRepository defines the interface for monthly target repository
type MonthlyTargetRepository interface {
	// FindByID finds a monthly target by ID
	FindByID(id string) (*monthly_target.MonthlyTarget, error)
	
	// FindByGroupAndPeriod finds group target by group ID, year, and month
	FindByGroupAndPeriod(groupID string, year int, month int) (*monthly_target.MonthlyTarget, error)
	
	// FindByUserAndPeriod finds user target by user ID, year, and month
	FindByUserAndPeriod(userID string, year int, month int) (*monthly_target.MonthlyTarget, error)
	
	// FindByBrickAndPeriod finds brick target by brick ID, year, and month
	FindByBrickAndPeriod(brickID string, year int, month int) (*monthly_target.MonthlyTarget, error)
	
	// List returns a list of monthly targets with pagination
	List(req *monthly_target.ListMonthlyTargetsRequest) ([]monthly_target.MonthlyTarget, int64, int64, error)
	
	// Create creates a new monthly target
	Create(target *monthly_target.MonthlyTarget) error
	
	// Update updates a monthly target
	Update(target *monthly_target.MonthlyTarget) error
	
	// Delete soft deletes a monthly target
	Delete(id string) error
	
	// GetUserEffectiveTarget gets effective target for user (user target or group target fallback)
	GetUserEffectiveTarget(userID string, year int, month int) (*monthly_target.MonthlyTarget, error)
	
	// BatchGetUserEffectiveTargets gets effective targets for multiple users in one query (optimized for N+1 prevention)
	BatchGetUserEffectiveTargets(userIDs []string, year int, month int) (map[string]*monthly_target.MonthlyTarget, error)

	// GetProratedTargetForPeriod calculates the time-weighted target for a specific period
	GetProratedTargetForPeriod(userID string, startDate, endDate string) (float64, error)

	// BatchGetProratedTargetsForPeriod calculates time-weighted targets for multiple users
	BatchGetProratedTargetsForPeriod(userIDs []string, startDate, endDate string) (map[string]float64, error)

	// GetTotalEffectiveTarget gets the sum of user monthly target records for a specific month
	GetTotalEffectiveTarget(year int, month int) (int64, error)
}

