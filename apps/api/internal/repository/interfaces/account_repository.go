package interfaces

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/account"
)

// AccountRepository defines the interface for account repository
type AccountRepository interface {
	// FindByID finds an account by ID
	FindByID(id string) (*account.Account, error)
	
	// List returns a list of accounts with pagination
	List(req *account.ListAccountsRequest) ([]account.Account, int64, error)
	
	// ListAll returns all accounts without pagination (for map display)
	ListAll(status string) ([]account.Account, error)

	// ListByBBox returns accounts within a geographic bounding box for viewport-based map loading
	ListByBBox(req *account.BBoxRequest) ([]account.Account, error)
	
	// Create creates a new account
	Create(account *account.Account) error
	
	// Update updates an account
	Update(account *account.Account) error
	
	// Delete soft deletes an account
	Delete(id string) error
	
	// GetStatsByStatus returns account statistics grouped by status (optimized aggregation)
	GetStatsByStatus() (map[string]int64, error)
	
	// CountByDateRange returns count of accounts created in date range (optimized)
	CountByDateRange(startDate, endDate interface{}) (int64, error)
}

