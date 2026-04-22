package interfaces

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/activity"
)

// ActivityRepository defines the interface for activity repository
type ActivityRepository interface {
	// FindByID finds an activity by ID
	FindByID(id string) (*activity.Activity, error)

	// List returns a list of activities with pagination
	List(req *activity.ListActivitiesRequest) ([]activity.Activity, int64, error)

	// Create creates a new activity
	Create(a *activity.Activity) error

	// Update updates an activity
	Update(a *activity.Activity) error

	// UpdateByLeadID updates all activities for a lead (batch update to fix N+1)
	UpdateByLeadID(leadID string, dealID, accountID *string) error

	// Delete soft deletes an activity
	Delete(id string) error

	// GetTimeline returns activity timeline for account/contact/user
	GetTimeline(req *activity.ActivityTimelineRequest) ([]activity.Activity, error)

	// FindByAccountID finds activities by account ID
	FindByAccountID(accountID string) ([]activity.Activity, error)

	// GetStatsByType returns activity statistics grouped by type (optimized aggregation)
	GetStatsByType(startDate, endDate string, accountID string) (map[string]int64, error)

	// GetStatsByTypeAndDate returns activity statistics grouped by type and date (optimized aggregation)
	GetStatsByTypeAndDate(startDate, endDate string, accountID string) (map[string]map[string]int64, error)

	// GetStatsByUser returns activity count grouped by user (optimized aggregation)
	GetStatsByUser(startDate, endDate string, accountID string) (map[string]int64, error)
}
