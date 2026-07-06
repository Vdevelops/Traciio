package interfaces

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/sales_overview"
)

// SalesOverviewRepository defines the interface for sales overview repository
type SalesOverviewRepository interface {
	// GetSalesPerformanceDetail gets detailed performance metrics for a user
	GetSalesPerformanceDetail(userID string, startDate, endDate interface{}) (*sales_overview.SalesPerformanceDetail, error)

	// GetSalesRepDetail gets detailed information about a sales representative
	GetSalesRepDetail(userID string, startDate, endDate interface{}) (*sales_overview.SalesRepDetail, error)

	// ListSalesPerformance returns a list of sales performance metrics with pagination
	ListSalesPerformance(req *sales_overview.ListSalesPerformanceRequest) ([]sales_overview.SalesPerformanceListResponse, int64, error)

	// ListProspectOutcomes returns prospect outcome rows across sales reps with pagination
	ListProspectOutcomes(req *sales_overview.ListProspectOutcomesRequest, startDate, endDate interface{}) ([]sales_overview.ProspectOutcomeListItem, int64, error)

	// GetSalesRepCheckInLocations returns a list of check-in locations for a sales representative
	GetSalesRepCheckInLocations(userID string, req *sales_overview.GetSalesRepCheckInLocationsRequest, startDate, endDate interface{}) ([]sales_overview.CheckInLocation, int64, error)

	// GetMonthlySalesOverview returns sales trend data for the chart
	GetMonthlySalesOverview(startDate, endDate interface{}, trendMode string, scopedUserIDs []string) (*sales_overview.MonthlySalesOverviewResponse, error)

	// GetFunnelDiagnostics returns stalled deals, no-activity deals, and stage aging summary
	GetFunnelDiagnostics(req *sales_overview.GetFunnelDiagnosticsRequest, scopedUserIDs []string, stalledThresholdDays, noActivityThresholdDays, limit int) (*sales_overview.FunnelDiagnosticsResponse, error)
}
