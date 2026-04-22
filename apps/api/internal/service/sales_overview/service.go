package sales_overview

import (
	"errors"
	"fmt"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/sales_overview"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	cachepkg "github.com/gilabs/crm-healthcare/api/pkg/cache"
	"github.com/gilabs/crm-healthcare/api/pkg/util/currency"
	// For currency formatting
)

var (
	ErrSalesPerformanceNotFound = errors.New("sales performance not found")
	ErrInvalidDateRange         = errors.New("invalid date range")
)

type Service struct {
	salesOverviewRepo interfaces.SalesOverviewRepository
	monthlyTargetRepo interfaces.MonthlyTargetRepository
	ac                *cachepkg.AdvancedCache
}

func NewService(
	salesOverviewRepo interfaces.SalesOverviewRepository,
	monthlyTargetRepo interfaces.MonthlyTargetRepository,
) *Service {
	return &Service{
		salesOverviewRepo: salesOverviewRepo,
		monthlyTargetRepo: monthlyTargetRepo,
		ac:                cachepkg.Advanced(),
	}
}

type cachedSalesPerformanceList struct {
	Items    []sales_overview.SalesPerformanceListResponse `msgpack:"items"`
	Total    int64                                         `msgpack:"total"`
	CachedAt time.Time                                     `msgpack:"cached_at"`
}

// GetMonthlySalesOverview returns aggregated monthly sales data for the chart
func (s *Service) GetMonthlySalesOverview(startDate, endDate interface{}, scopedUserIDs []string) (*sales_overview.MonthlySalesOverviewResponse, error) {
	// If startDate and endDate are strings, try to parse them if needed, 
	// but the handler usually handles parsing query params to interface{} or specific types.
	// The repository interface expects interface{} which usually means time.Time or string that GORM can handle.
	// For consistency with other methods, let's assume the handler passes proper types or compatible strings.
	
	// Check cache
	startKey := ""
	endKey := ""
	if t, ok := startDate.(time.Time); ok {
		startKey = t.Format("2006-01-02")
	} else if s, ok := startDate.(string); ok {
		startKey = s
	}
	if t, ok := endDate.(time.Time); ok {
		endKey = t.Format("2006-01-02")
	} else if s, ok := endDate.(string); ok {
		endKey = s
	}

	if s.ac != nil && s.ac.IsEnabled() {
		// Include scope in cache key to prevent cross-scope pollution
		scopeKey := "global"
		if len(scopedUserIDs) > 0 {
			scopeKey = fmt.Sprintf("%v", scopedUserIDs)
		}
		cacheKey := fmt.Sprintf("sales_overview:monthly:scope:%s:start:%s:end:%s", scopeKey, startKey, endKey)
		var cached sales_overview.MonthlySalesOverviewResponse
		if found, _ := s.ac.Get(cacheKey, &cached); found {
			return &cached, nil
		}
	}

	result, err := s.salesOverviewRepo.GetMonthlySalesOverview(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Calculate and populate targets for each month
	if result != nil && len(result.MonthlyData) > 0 {
		var start, end time.Time
		if t, ok := startDate.(time.Time); ok {
			start = t
		} else if s, ok := startDate.(string); ok {
			start, _ = time.Parse("2006-01-02", s)
		}
		if t, ok := endDate.(time.Time); ok {
			end = t
		} else if s, ok := endDate.(string); ok {
			end, _ = time.Parse("2006-01-02", s)
		}

		for i := range result.MonthlyData {
			year := result.MonthlyData[i].Year
			month := result.MonthlyData[i].Month

			var totalTarget int64
			if len(scopedUserIDs) > 0 {
				// Scoped: aggregate targets for scoped users only
				for _, uid := range scopedUserIDs {
					et, etErr := s.monthlyTargetRepo.GetUserEffectiveTarget(uid, year, month)
					if etErr == nil && et != nil {
						totalTarget += et.TargetAmount
					}
				}
			} else {
				// Global: get total effective target for all active users
				var getErr error
				totalTarget, getErr = s.monthlyTargetRepo.GetTotalEffectiveTarget(year, month)
				if getErr != nil {
					continue
				}
			}

			// Calculate proration
			daysInMonth := time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC).Day()
			activeDays := daysInMonth

			if !start.IsZero() && !end.IsZero() {
				// Calculate overlap
				monthStart := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
				monthEnd := time.Date(year, time.Month(month+1), 0, 23, 59, 59, 0, time.UTC)
				
				periodStart := monthStart
				if start.After(monthStart) {
					periodStart = start
				}
				
				periodEnd := monthEnd
				endEndOfDay := time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, time.UTC)
				if endEndOfDay.Before(monthEnd) {
					periodEnd = endEndOfDay
				}

				if periodStart.After(periodEnd) {
					activeDays = 0
				} else {
					// activeDays should be inclusive
					// Truncate to days
					d1 := time.Date(periodStart.Year(), periodStart.Month(), periodStart.Day(), 0, 0, 0, 0, time.UTC)
					d2 := time.Date(periodEnd.Year(), periodEnd.Month(), periodEnd.Day(), 0, 0, 0, 0, time.UTC)
					activeDays = int(d2.Sub(d1).Hours() / 24) + 1
				}
			}
			
			if activeDays < 0 { activeDays = 0 }

			if daysInMonth > 0 {
				result.MonthlyData[i].TargetAmount = int64(float64(totalTarget) * float64(activeDays) / float64(daysInMonth))
			}
		}
	}

	if s.ac != nil && s.ac.IsEnabled() && result != nil {
		scopeKey := "global"
		if len(scopedUserIDs) > 0 {
			scopeKey = fmt.Sprintf("%v", scopedUserIDs)
		}
		cacheKey := fmt.Sprintf("sales_overview:monthly:scope:%s:start:%s:end:%s", scopeKey, startKey, endKey)
		_ = s.ac.Set(cacheKey, result, cachepkg.TTLStatsShort)
	}

	return result, nil
}

// GetSalesPerformanceDetail gets detailed performance metrics for a user
func (s *Service) GetSalesPerformanceDetail(userID string, req *sales_overview.GetSalesPerformanceDetailRequest) (*sales_overview.SalesPerformanceDetail, error) {
	var startDate, endDate interface{}

	// Parse date range
	if req.StartDate != "" {
		parsed, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, ErrInvalidDateRange
		}
		startDate = parsed
	}
	if req.EndDate != "" {
		parsed, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, ErrInvalidDateRange
		}
		// Set to end of day
		endDate = parsed.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}

	// If period is specified but dates are not, calculate dates based on period
	if req.Period != "" && startDate == nil && endDate == nil {
		now := time.Now()
		switch req.Period {
		case "today":
			startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			endDate = startDate.(time.Time).Add(24*time.Hour - 1*time.Second)
		case "week":
			weekday := int(now.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			startDate = startDate.(time.Time).AddDate(0, 0, -(weekday - 1))
			endDate = startDate.(time.Time).AddDate(0, 0, 7).Add(-1 * time.Second)
		case "month":
			startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			endDate = startDate.(time.Time).AddDate(0, 1, 0).Add(-1 * time.Second)
		case "year":
			startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
			endDate = startDate.(time.Time).AddDate(1, 0, 0).Add(-1 * time.Second)
		}
	}

	startKey := ""
	endKey := ""
	if startDate != nil {
		startKey = startDate.(time.Time).Format(time.RFC3339)
	}
	if endDate != nil {
		endKey = endDate.(time.Time).Format(time.RFC3339)
	}

	if s.ac != nil && s.ac.IsEnabled() {
		cacheKey := fmt.Sprintf("sales_overview:perf_detail:user:%s:start:%s:end:%s", userID, startKey, endKey)
		var cached sales_overview.SalesPerformanceDetail
		if found, _ := s.ac.Get(cacheKey, &cached); found {
			return &cached, nil
		}
	}

	detail, err := s.salesOverviewRepo.GetSalesPerformanceDetail(userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	if s.ac != nil && s.ac.IsEnabled() {
		cacheKey := fmt.Sprintf("sales_overview:perf_detail:user:%s:start:%s:end:%s", userID, startKey, endKey)
		_ = s.ac.Set(cacheKey, detail, cachepkg.TTLStatsShort)
	}

	return detail, nil
}

// GetSalesRepDetail gets comprehensive detail for sales rep detail page
func (s *Service) GetSalesRepDetail(userID string, req *sales_overview.GetSalesRepDetailRequest) (*sales_overview.SalesRepDetail, error) {
	var startDate, endDate interface{}

	// Parse date range
	if req.StartDate != "" {
		parsed, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, ErrInvalidDateRange
		}
		startDate = parsed
	}
	if req.EndDate != "" {
		parsed, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, ErrInvalidDateRange
		}
		endDate = parsed.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}

	// If period is specified but dates are not, calculate dates based on period
	if req.Period != "" && startDate == nil && endDate == nil {
		now := time.Now()
		switch req.Period {
		case "today":
			startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			endDate = startDate.(time.Time).Add(24*time.Hour - 1*time.Second)
		case "week":
			weekday := int(now.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			startDate = startDate.(time.Time).AddDate(0, 0, -(weekday - 1))
			endDate = startDate.(time.Time).AddDate(0, 0, 7).Add(-1 * time.Second)
		case "month":
			startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			endDate = startDate.(time.Time).AddDate(0, 1, 0).Add(-1 * time.Second)
		case "year":
			startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
			endDate = startDate.(time.Time).AddDate(1, 0, 0).Add(-1 * time.Second)
		}
	}

	startKey := ""
	endKey := ""
	if startDate != nil {
		startKey = startDate.(time.Time).Format(time.RFC3339)
	}
	if endDate != nil {
		endKey = endDate.(time.Time).Format(time.RFC3339)
	}

	if s.ac != nil && s.ac.IsEnabled() {
		cacheKey := fmt.Sprintf("sales_overview:sales_rep_detail:user:%s:start:%s:end:%s", userID, startKey, endKey)
		var cached sales_overview.SalesRepDetail
		if found, _ := s.ac.Get(cacheKey, &cached); found {
			return &cached, nil
		}
	}

	detail, err := s.salesOverviewRepo.GetSalesRepDetail(userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	if s.ac != nil && s.ac.IsEnabled() {
		cacheKey := fmt.Sprintf("sales_overview:sales_rep_detail:user:%s:start:%s:end:%s", userID, startKey, endKey)
		_ = s.ac.Set(cacheKey, detail, cachepkg.TTLStatsShort)
	}

	return detail, nil
}

// ListSalesPerformance lists all sales reps with performance summary
func (s *Service) ListSalesPerformance(req *sales_overview.ListSalesPerformanceRequest) ([]sales_overview.SalesPerformanceListResponse, int64, error) {
	cacheKey := ""
	if s.ac != nil && s.ac.IsEnabled() {
		scopeKey := "global"
		if len(req.ScopedUserIDs) > 0 {
			scopeKey = fmt.Sprintf("%v", req.ScopedUserIDs)
		}
		cacheKey = fmt.Sprintf(
			"sales_overview:list:scope:%s:page:%d:per:%d:period:%s:start:%s:end:%s:search:%s:brick:%s:sort:%s:order:%s",
			scopeKey,
			req.Page,
			req.PerPage,
			req.Period,
			req.StartDate,
			req.EndDate,
			req.Search,
			req.BrickID,
			req.SortBy,
			req.Order,
		)
		var cached cachedSalesPerformanceList
		if found, _ := s.ac.Get(cacheKey, &cached); found {
			return cached.Items, cached.Total, nil
		}
	}

	results, total, err := s.salesOverviewRepo.ListSalesPerformance(req)
	if err != nil {
		return nil, 0, err
	}

	// Calculate and attach targets
	
	// Identify months involved
	// We iterate from start month to end month
	userIDs := make([]string, len(results))
    for i, r := range results {
        userIDs[i] = r.UserID
    }

    // We need to fetch targets for each month in the range
    // Since BatchGetUserEffectiveTargets takes year/month, we loop through months
    
	// Calculate Prorated Targets for the period
	proratedTargets, err := s.monthlyTargetRepo.BatchGetProratedTargetsForPeriod(userIDs, req.StartDate, req.EndDate)
	if err != nil {
		// Log error but continue
		// fmt.Printf("Error fetching prorated targets: %v\n", err)
	}

	// Update results
	for i := range results {
		target := int64(0)
		if val, ok := proratedTargets[results[i].UserID]; ok {
			target = int64(val)
		}
		
		results[i].TargetAmount = target
		results[i].TargetAmountFormatted = currency.FormatCurrency(target)
		
		if target > 0 {
			results[i].TargetAchievementPercentage = (float64(results[i].TotalRevenue) / float64(target)) * 100
		} else {
			results[i].TargetAchievementPercentage = 0
		}
	}



	if cacheKey != "" {
		_ = s.ac.Set(cacheKey, &cachedSalesPerformanceList{Items: results, Total: total, CachedAt: time.Now()}, cachepkg.TTLStatsShort)
	}

	return results, total, nil
}

// GetSalesRepCheckInLocations gets check-in locations for sales rep with pagination
func (s *Service) GetSalesRepCheckInLocations(userID string, req *sales_overview.GetSalesRepCheckInLocationsRequest) (*sales_overview.SalesRepCheckInLocationsResponse, error) {
	var startDate, endDate interface{}

	// Parse date range
	if req.StartDate != "" {
		parsed, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, ErrInvalidDateRange
		}
		startDate = parsed
	}
	if req.EndDate != "" {
		parsed, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, ErrInvalidDateRange
		}
		endDate = parsed.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}

	// If no date range provided, use last 30 days as default to show recent data
	if startDate == nil && endDate == nil {
		now := time.Now()
		endDate = now
		startDate = now.AddDate(0, 0, -30)
	}

	startKey := ""
	endKey := ""
	if startDate != nil {
		startKey = startDate.(time.Time).Format(time.RFC3339)
	}
	if endDate != nil {
		endKey = endDate.(time.Time).Format(time.RFC3339)
	}

	cacheKey := ""
	if s.ac != nil && s.ac.IsEnabled() {
		cacheKey = fmt.Sprintf("sales_overview:checkins:user:%s:start:%s:end:%s:page:%d:per:%d", userID, startKey, endKey, req.Page, req.PerPage)
		var cached sales_overview.SalesRepCheckInLocationsResponse
		if found, _ := s.ac.Get(cacheKey, &cached); found {
			return &cached, nil
		}
	}

	locations, total, err := s.salesOverviewRepo.GetSalesRepCheckInLocations(userID, req, startDate, endDate)
	if err != nil {
		return nil, err
	}

	var period *sales_overview.PeriodRange
	if startDate != nil && endDate != nil {
		period = &sales_overview.PeriodRange{
			Start: startDate.(time.Time),
			End:   endDate.(time.Time),
		}
	}

	response := &sales_overview.SalesRepCheckInLocationsResponse{
		CheckInLocations: locations,
		TotalVisits:      total,
		Period:           period,
	}

	if cacheKey != "" {
		_ = s.ac.Set(cacheKey, response, cachepkg.TTLStatsShort)
	}

	return response, nil
}
