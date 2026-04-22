package interfaces

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/visit_report"
)

// VisitReportRepository defines the interface for visit report repository
type VisitReportRepository interface {
	// FindByID finds a visit report by ID
	FindByID(id string) (*visit_report.VisitReport, error)

	// List returns a list of visit reports with pagination
	List(req *visit_report.ListVisitReportsRequest) ([]visit_report.VisitReport, int64, error)

	// Create creates a new visit report
	Create(vr *visit_report.VisitReport) error

	// Update updates a visit report
	Update(vr *visit_report.VisitReport) error

	// UpdateByLeadID updates all visit reports for a lead (batch update to fix N+1)
	UpdateByLeadID(leadID string, dealID, accountID *string) error

	// Delete soft deletes a visit report
	Delete(id string) error

	// FindByAccountID finds visit reports by account ID
	FindByAccountID(accountID string) ([]visit_report.VisitReport, error)

	// FindBySalesRepID finds visit reports by sales rep ID
	FindBySalesRepID(salesRepID string) ([]visit_report.VisitReport, error)

	// GetStatsByStatus returns visit report statistics grouped by status (optimized aggregation)
	GetStatsByStatus(startDate, endDate string, accountID, salesRepID, status string) (map[string]int64, error)

	// GetStatsByStatusForUsers returns visit report statistics for multiple users (batch query to fix N+1)
	GetStatsByStatusForUsers(startDate, endDate string, userIDs []string) (map[string]int64, error)

	// GetStatsByDate returns visit report count grouped by date (optimized aggregation)
	GetStatsByDate(startDate, endDate string, accountID, salesRepID, status string) (map[string]int64, error)

	// GetStatsByDateAndStatus returns visit report count grouped by date and status (optimized aggregation)
	// Returns map[date]map[status]count
	GetStatsByDateAndStatus(startDate, endDate string, accountID, salesRepID string) (map[string]map[string]int64, error)

	// GetStatsByAccount returns visit report count grouped by account (optimized aggregation)
	GetStatsByAccount(startDate, endDate string, salesRepID, status string) (map[string]int64, error)

	// GetStatsBySalesRep returns visit report count grouped by sales rep (optimized aggregation)
	GetStatsBySalesRep(startDate, endDate string, accountID, status string) (map[string]int64, error)

	// GetStatsBySalesRepWithAccounts returns visit count and unique account count per sales rep (optimized aggregation)
	GetStatsBySalesRepWithAccounts(startDate, endDate string, status string) (map[string]struct {
		VisitCount   int64
		AccountCount int64
	}, error)
}
