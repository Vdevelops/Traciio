package user

import (
	"errors"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
	"github.com/gilabs/crm-healthcare/api/internal/domain/task"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"github.com/gilabs/crm-healthcare/api/internal/domain/visit_report"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"github.com/gilabs/crm-healthcare/api/pkg/util/currency"
	"gorm.io/gorm"
)

type SettingsService struct {
	userRepo        interfaces.UserRepository
	dealRepo        interfaces.DealRepository
	visitReportRepo interfaces.VisitReportRepository
	taskRepo        interfaces.TaskRepository
	profileService  *ProfileService // Reuse for activities/transactions
}

func NewSettingsService(
	userRepo interfaces.UserRepository,
	dealRepo interfaces.DealRepository,
	visitReportRepo interfaces.VisitReportRepository,
	taskRepo interfaces.TaskRepository,
	profileService *ProfileService,
) *SettingsService {
	return &SettingsService{
		userRepo:        userRepo,
		dealRepo:        dealRepo,
		visitReportRepo: visitReportRepo,
		taskRepo:        taskRepo,
		profileService:  profileService,
	}
}

// GetSettingsSummary returns complete settings summary for authenticated user
// This is user-scoped and should ONLY be called with authenticated user's ID
// Accepts optional date filtering via startDate and endDate parameters
func (s *SettingsService) GetSettingsSummary(userID string, startDate, endDate interface{}) (*user.SettingsSummaryResponse, error) {
	// Get user
	u, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Convert date parameters to *time.Time for consistent handling
	var dueDateFrom, dueDateTo *time.Time
	if startDate != nil {
		if t, ok := startDate.(time.Time); ok {
			dueDateFrom = &t
		}
	}
	if endDate != nil {
		if t, ok := endDate.(time.Time); ok {
			dueDateTo = &t
		}
	}

	// === OPTIMIZED: Use direct SQL aggregation instead of fetching all deals ===
	
	// 1. Get won deals count and revenue using optimized repository method
	// This method uses actual_close_date for won deals (when they were actually closed)
	dealsWonCount, totalRevenue, err := s.dealRepo.GetWonDealsValueInPeriodByUser(userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// 2. Count total deals in period
	dealReq := &pipeline.ListDealsRequest{
		AssignedTo: userID,
		Page:       1,
		PerPage:    1, // Only need count, not data
	}
	if startDate != nil {
		if t, ok := startDate.(time.Time); ok {
			dealReq.DateFrom = t.Format("2006-01-02")
		}
	}
	if endDate != nil {
		if t, ok := endDate.(time.Time); ok {
			dealReq.DateTo = t.Format("2006-01-02")
		}
	}
	_, totalDeals, err := s.dealRepo.List(dealReq)
	if err != nil {
		return nil, err
	}

	// 3. Count visits COMPLETED - MATCH SALES OVERVIEW
	// Filter by visit_date and completed visit status only
	var visitsCompleted int64
	visitReq := &visit_report.ListVisitReportsRequest{
		SalesRepID: userID,
		Status:     "completed",
		Page:       1,
		PerPage:    1,
	}
	if startDate != nil {
		if t, ok := startDate.(time.Time); ok {
			dateStr := t.Format("2006-01-02")
			visitReq.StartDate = dateStr // Filters by visit_date
		}
	}
	if endDate != nil {
		if t, ok := endDate.(time.Time); ok {
			dateStr := t.Format("2006-01-02")
			visitReq.EndDate = dateStr // Filters by visit_date
		}
	}
	_, visitsCompleted, err = s.visitReportRepo.List(visitReq)
	if err != nil {
		return nil, err
	}

	// 4. Count tasks COMPLETED - MATCH SALES OVERVIEW
	// Filter by created_at (not due_date!) and status='completed'
	var tasksCompleted int64
	taskReq := &task.ListTasksRequest{
		AssignedTo: userID,
		Status:     "completed", // Only count completed tasks like Sales Overview
		Page:       1,
		PerPage:    1,
	}
	// Sales Overview uses created_at for tasks
	if dueDateFrom != nil {
		taskReq.CreatedFrom = dueDateFrom
	}
	if dueDateTo != nil {
		taskReq.CreatedTo = dueDateTo
	}
	_, tasksCompleted, err = s.taskRepo.List(taskReq)
	if err != nil {
		return nil, err
	}

	// Get activities (reuse from ProfileService)
	activities, err := s.profileService.getActivities(userID)
	if err != nil {
		return nil, err
	}

	// Get transactions (reuse from ProfileService)
	transactions, err := s.profileService.getTransactions(userID)
	if err != nil {
		return nil, err
	}

	// Format currency (convert sen to rupiah)
	totalRevenueFormatted := formatCurrency(totalRevenue)

	// Calculate conversion rate
	conversionRate := 0.0
	if totalDeals > 0 {
		conversionRate = (float64(dealsWonCount) / float64(totalDeals)) * 100.0
	}

	// Build response
	summary := &user.SettingsSummaryResponse{
		User: u.ToUserResponse(),
		Stats: &user.SettingsStats{
			TotalRevenue:          totalRevenue,
			TotalRevenueFormatted: totalRevenueFormatted,
			Deals:                 int(totalDeals),        // Total deals
			DealsWon:              int(dealsWonCount),      // Won deals
			Visits:                int(visitsCompleted),    // Approved visits only
			Tasks:                 int(tasksCompleted),     // Completed tasks only
			ConversionRate:        conversionRate,
		},
		Activities:   activities,
		Transactions: transactions,
	}

	return summary, nil
}

// getExtendedStats calculates extended statistics including revenue metrics
// Accepts optional date filtering via startDate and endDate parameters
func (s *SettingsService) getExtendedStats(userID string, startDate, endDate interface{}) (*user.SettingsStats, error) {
	// Build date filter strings for repository queries
	var startDateStr, endDateStr string
	if startDate != nil {
		if t, ok := startDate.(time.Time); ok {
			startDateStr = t.Format("2006-01-02")
		}
	}
	if endDate != nil {
		if t, ok := endDate.(time.Time); ok {
			endDateStr = t.Format("2006-01-02")
		}
	}

	// Count visit reports with date filter
	visitReq := &visit_report.ListVisitReportsRequest{
		SalesRepID: userID,
		StartDate:  startDateStr,
		EndDate:    endDateStr,
		Page:       1,
		PerPage:    1,
	}
	_, totalVisits, err := s.visitReportRepo.List(visitReq)
	if err != nil {
		return nil, err
	}

	// Get all deals for user (for aggregation) with date filter
	// Note: pipeline.ListDealsRequest uses DateFrom/DateTo, not StartDate/EndDate
	dealReq := &pipeline.ListDealsRequest{
		AssignedTo: userID,
		DateFrom:   startDateStr,
		DateTo:     endDateStr,
		Page:       1,
		PerPage:    10000, // Get all deals for aggregation (optimize with dedicated query later)
	}
	deals, totalDeals, err := s.dealRepo.List(dealReq)
	if err != nil {
		return nil, err
	}

	// Parse dates for task filtering
	// Note: task.ListTasksRequest uses DueDateFrom/DueDateTo as *time.Time, not strings
	var dueDateFrom, dueDateTo *time.Time
	if startDate != nil {
		if t, ok := startDate.(time.Time); ok {
			dueDateFrom = &t
		}
	}
	if endDate != nil {
		if t, ok := endDate.(time.Time); ok {
			dueDateTo = &t
		}
	}

	// Count tasks with date filter
	taskReq := &task.ListTasksRequest{
		AssignedTo:  userID,
		DueDateFrom: dueDateFrom,
		DueDateTo:   dueDateTo,
		Page:        1,
		PerPage:     1,
	}
	_, totalTasks, err := s.taskRepo.List(taskReq)
	if err != nil {
		return nil, err
	}

	// Aggregate deal stats
	// IMPORTANT: Revenue should ONLY include WON deals, not all deals
	var totalRevenue int64
	var dealsWon, dealsLost, dealsOpen int

	for _, deal := range deals {
		switch deal.Status {
		case "won":
			dealsWon++
			totalRevenue += deal.Value // Only add won deal values to revenue
		case "lost":
			dealsLost++
		default:
			dealsOpen++
		}
	}

	// Calculate conversion rate (won deals / total deals * 100)
	var conversionRate float64
	if totalDeals > 0 {
		conversionRate = (float64(dealsWon) / float64(totalDeals)) * 100
	}

	// Calculate average deal value
	var avgDealValue int64
	if dealsWon > 0 {
		avgDealValue = totalRevenue / int64(dealsWon)
	}

	// Format currency values
	revenueFormatted := currency.FormatCurrency(totalRevenue)
	avgDealValueFormatted := currency.FormatCurrency(avgDealValue)

	return &user.SettingsStats{
		Visits:                    int(totalVisits),
		Deals:                     int(totalDeals),
		Tasks:                     int(totalTasks),
		TotalRevenue:              totalRevenue,
		DealsWon:                  dealsWon,
		DealsLost:                 dealsLost,
		DealsOpen:                 dealsOpen,
		TotalRevenueFormatted:     revenueFormatted,
		ConversionRate:            conversionRate,
		AverageDealValue:          avgDealValue,
		AverageDealValueFormatted: avgDealValueFormatted,
	}, nil
}
