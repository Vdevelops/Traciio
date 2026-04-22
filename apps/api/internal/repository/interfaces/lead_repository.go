package interfaces

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead"
)

// LeadRepository defines the interface for lead repository
type LeadRepository interface {
	// FindByID finds a lead by ID
	FindByID(id string) (*lead.Lead, error)

	// List returns a list of leads with pagination
	List(req *lead.ListLeadsRequest) ([]lead.Lead, int64, error)

	// Create creates a new lead
	Create(lead *lead.Lead) error

	// Update updates a lead
	Update(lead *lead.Lead) error

	// Delete soft deletes a lead
	Delete(id string) error

	// GetAnalytics returns lead analytics
	GetAnalytics(req *lead.LeadAnalyticsRequest) (*lead.LeadAnalyticsResponse, error)
	
	// GetStatsByStatus returns lead statistics grouped by status using database aggregation
	GetStatsByStatus() (map[string]int64, error)
	
	// GetStatsBySource returns lead statistics grouped by source using database aggregation
	GetStatsBySource() (map[string]int64, error)
	
	// CountByDateRange returns count of leads created in date range (optimized)
	CountByDateRange(startDate, endDate interface{}) (int64, error)
	
	// GetStatsByStatusAndDateRange returns lead statistics grouped by status within date range (optimized)
	GetStatsByStatusAndDateRange(startDate, endDate interface{}) (map[string]int64, error)
	
	// GetStatsBySourceAndDateRange returns lead statistics grouped by source using database aggregation with date filter
	GetStatsBySourceAndDateRange(startDate, endDate interface{}) (map[string]int64, error)
}

