package interfaces

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead_status"
)

// LeadStatusRepository defines the interface for lead status repository
type LeadStatusRepository interface {
	// FindByID finds a lead status by ID
	FindByID(id string) (*lead_status.LeadStatus, error)

	// FindByCode finds a lead status by code
	FindByCode(code string) (*lead_status.LeadStatus, error)

	// FindDefault finds the default lead status
	FindDefault() (*lead_status.LeadStatus, error)

	// ListAll returns a list of lead statuses
	ListAll() ([]*lead_status.LeadStatus, error)
}
