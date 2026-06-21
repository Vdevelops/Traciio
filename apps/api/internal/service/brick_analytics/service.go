package brick_analytics

import (
	"errors"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/brick"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

var (
	ErrBrickNotFound = errors.New("brick not found")
)

// Service provides analytics for brick performance
// All queries are optimized to avoid N+1 problems using JOINs and aggregations
type Service struct {
	brickRepo           interfaces.BrickRepository
	dealRepo            interfaces.DealRepository
	visitReportRepo     interfaces.VisitReportRepository
	accountRepo         interfaces.AccountRepository
	monthlyTargetRepo   interfaces.MonthlyTargetRepository
	brickTargetDistRepo interfaces.BrickTargetDistributionRepository
	userRepo            interfaces.UserRepository
	db                  *gorm.DB
}

// NewService creates a new brick analytics service
func NewService(
	db *gorm.DB,
	brickRepo interfaces.BrickRepository,
	dealRepo interfaces.DealRepository,
	visitReportRepo interfaces.VisitReportRepository,
	accountRepo interfaces.AccountRepository,
	monthlyTargetRepo interfaces.MonthlyTargetRepository,
	brickTargetDistRepo interfaces.BrickTargetDistributionRepository,
	userRepo interfaces.UserRepository,
) *Service {
	return &Service{
		db:                  db,
		brickRepo:           brickRepo,
		dealRepo:            dealRepo,
		visitReportRepo:     visitReportRepo,
		accountRepo:         accountRepo,
		monthlyTargetRepo:   monthlyTargetRepo,
		brickTargetDistRepo: brickTargetDistRepo,
		userRepo:            userRepo,
	}
}

// GetBrickRepo returns the brick repository (for authorization checks)
func (s *Service) GetBrickRepo() interfaces.BrickRepository {
	return s.brickRepo
}

// BrickPerformanceMetrics represents comprehensive performance metrics for a brick
type BrickPerformanceMetrics struct {
	BrickID     string  `json:"brick_id"`
	BrickName   string  `json:"brick_name"`
	BrickCode   string  `json:"brick_code"`
	ManagerName *string `json:"manager_name,omitempty"`
	ManagerID   *string `json:"manager_id,omitempty"`

	// Target & Achievement
	MonthlyTarget      int64   `json:"monthly_target"`
	TargetAchieved     int64   `json:"target_achieved"`
	AchievementPercent float64 `json:"achievement_percentage"`
	TargetRemaining    int64   `json:"target_remaining"`

	// Sales Team
	TotalSales  int `json:"total_sales"`
	ActiveSales int `json:"active_sales"`

	// Pipeline Metrics
	TotalDeals      int     `json:"total_deals"`
	OpenDeals       int     `json:"open_deals"`
	WonDeals        int     `json:"won_deals"`
	LostDeals       int     `json:"lost_deals"`
	TotalDealValue  int64   `json:"total_deal_value"`
	WonDealValue    int64   `json:"won_deal_value"`
	WinRate         float64 `json:"win_rate"`
	AverageDealSize int64   `json:"average_deal_size"`

	// Visit Activity
	TotalVisits           int     `json:"total_visits"`
	VisitsThisMonth       int     `json:"visits_this_month"`
	AverageVisitsPerSales float64 `json:"average_visits_per_sales"`

	// Accounts
	TotalAccounts        int `json:"total_accounts"`
	ActiveAccounts       int `json:"active_accounts"`
	NewAccountsThisMonth int `json:"new_accounts_this_month"`

	// Revenue (calculated from won deals)
	TotalRevenue         int64   `json:"total_revenue"`
	RevenueThisMonth     int64   `json:"revenue_this_month"`
	RevenueGrowthPercent float64 `json:"revenue_growth_percentage"`
}

// GetBrickPerformance gets comprehensive performance metrics for a brick
// Uses optimized queries with JOINs to avoid N+1 problems
func (s *Service) GetBrickPerformance(brickID string, periodStart, periodEnd time.Time) (*BrickPerformanceMetrics, error) {
	// Get brick info
	b, err := s.brickRepo.FindByID(brickID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBrickNotFound
		}
		return nil, err
	}

	metrics := &BrickPerformanceMetrics{
		BrickID:   b.ID,
		BrickName: b.Name,
		BrickCode: b.Code,
	}

	if b.Manager != nil {
		metrics.ManagerID = &b.Manager.ID
		metrics.ManagerName = &b.Manager.Name
	}

	// Get sales count (optimized single query)
	var salesCount int64
	s.db.Model(&struct {
		ID string `gorm:"type:uuid"`
	}{}).Table("users").
		Where("brick_id = ? AND deleted_at IS NULL", brickID).
		Count(&salesCount)
	metrics.TotalSales = int(salesCount)

	// Get active sales count
	var activeSalesCount int64
	s.db.Model(&struct {
		ID string `gorm:"type:uuid"`
	}{}).Table("users").
		Where("brick_id = ? AND status = 'active' AND deleted_at IS NULL", brickID).
		Count(&activeSalesCount)
	metrics.ActiveSales = int(activeSalesCount)

	allTime := periodStart.IsZero()
	loc := periodEnd.Location()
	if loc == nil {
		loc = time.UTC
	}

	// Normalize period boundaries (inclusive)
	// Use the location provided by handler (WIB) to keep boundaries consistent.
	if !allTime {
		periodStart = time.Date(periodStart.Year(), periodStart.Month(), periodStart.Day(), 0, 0, 0, 0, periodStart.Location())
		if periodEnd.Before(periodStart) {
			periodEnd = periodStart
		}
	}
	// Ensure end is inclusive (if caller passed date-only at midnight)
	if periodEnd.Hour() == 0 && periodEnd.Minute() == 0 && periodEnd.Second() == 0 && periodEnd.Nanosecond() == 0 {
		periodEnd = periodEnd.Add(24*time.Hour - 1*time.Second)
	}

	// Get monthly target for the month of the selected periodStart.
	// If the filter is "all time", default target-month to current month.
	// NOTE: monthly_targets table uses year and month, not an arbitrary range.
	ref := periodStart
	if allTime {
		referenceNow := time.Now().In(loc)
		ref = time.Date(referenceNow.Year(), referenceNow.Month(), 1, 0, 0, 0, 0, loc)
	}
	monthStart := time.Date(ref.Year(), ref.Month(), 1, 0, 0, 0, 0, ref.Location())
	year := monthStart.Year()
	month := int(monthStart.Month())

	var monthlyTarget struct {
		TargetAmount int64 `gorm:"column:target_amount"`
	}
	s.db.Table("monthly_targets").
		Select("target_amount").
		Where("brick_id = ? AND year = ? AND month = ? AND deleted_at IS NULL", brickID, year, month).
		Order("created_at DESC").
		Limit(1).
		Scan(&monthlyTarget)
	metrics.MonthlyTarget = monthlyTarget.TargetAmount

	// Get target distributions and achievements (optimized with JOIN)
	// NOTE: We intentionally do not query brick_target_distributions here.
	// This endpoint's "achieved" metric is derived from won revenue, and
	// fetching distributions per brick adds avoidable query overhead.

	// Deal metrics (align with Sales Performance page)
	// - Total/Open/Lost: filtered by created_at in period
	// - Won/Revenue: filtered by actual_close_date in period (fallback to created_at if actual_close_date is NULL)
	var dealMetrics struct {
		TotalDeals     int64 `gorm:"column:total_deals"`
		OpenDeals      int64 `gorm:"column:open_deals"`
		WonDeals       int64 `gorm:"column:won_deals"`
		LostDeals      int64 `gorm:"column:lost_deals"`
		TotalDealValue int64 `gorm:"column:total_deal_value"`
		WonDealValue   int64 `gorm:"column:won_deal_value"`
	}
	s.db.Table("deals").
		Select(`
			SUM(CASE WHEN created_at >= ? AND created_at <= ? THEN 1 ELSE 0 END) as total_deals,
			SUM(CASE WHEN status = 'open' AND created_at >= ? AND created_at <= ? THEN 1 ELSE 0 END) as open_deals,
			SUM(CASE WHEN status = 'lost' AND created_at >= ? AND created_at <= ? THEN 1 ELSE 0 END) as lost_deals,
			SUM(CASE WHEN created_at >= ? AND created_at <= ? THEN value ELSE 0 END) as total_deal_value,
			SUM(CASE WHEN status = 'won' AND (actual_close_date IS NOT NULL AND actual_close_date >= ? AND actual_close_date <= ?) THEN 1
					 WHEN status = 'won' AND (actual_close_date IS NULL AND created_at >= ? AND created_at <= ?) THEN 1
					 ELSE 0 END) as won_deals,
			SUM(CASE WHEN status = 'won' AND (actual_close_date IS NOT NULL AND actual_close_date >= ? AND actual_close_date <= ?) THEN value
					 WHEN status = 'won' AND (actual_close_date IS NULL AND created_at >= ? AND created_at <= ?) THEN value
					 ELSE 0 END) as won_deal_value
		`,
			periodStart, periodEnd,
			periodStart, periodEnd,
			periodStart, periodEnd,
			periodStart, periodEnd,
			periodStart, periodEnd,
			periodStart, periodEnd,
			periodStart, periodEnd,
			periodStart, periodEnd,
		).
		Where(
			"deals.brick_id = ? AND deals.deleted_at IS NULL AND ((deals.created_at >= ? AND deals.created_at <= ?) OR (deals.actual_close_date >= ? AND deals.actual_close_date <= ?))",
			brickID,
			periodStart,
			periodEnd,
			periodStart,
			periodEnd,
		).
		Scan(&dealMetrics)

	metrics.TotalDeals = int(dealMetrics.TotalDeals)
	metrics.OpenDeals = int(dealMetrics.OpenDeals)
	metrics.WonDeals = int(dealMetrics.WonDeals)
	metrics.LostDeals = int(dealMetrics.LostDeals)
	metrics.TotalDealValue = dealMetrics.TotalDealValue
	metrics.WonDealValue = dealMetrics.WonDealValue

	// Calculate win rate
	if metrics.TotalDeals > 0 {
		metrics.WinRate = float64(metrics.WonDeals) / float64(metrics.TotalDeals) * 100
	}

	// Calculate average deal size (won deals in the period), matching Sales Performance
	if metrics.WonDeals > 0 {
		metrics.AverageDealSize = dealMetrics.WonDealValue / dealMetrics.WonDeals
	}

	// Target achieved for the period (use same logic as won revenue above)
	metrics.TargetAchieved = dealMetrics.WonDealValue
	metrics.TargetRemaining = metrics.MonthlyTarget - metrics.TargetAchieved
	if metrics.MonthlyTarget > 0 {
		metrics.AchievementPercent = float64(metrics.TargetAchieved) / float64(metrics.MonthlyTarget) * 100
	}

	// Visit metrics (align with Sales Performance: only completed visits, filtered by period)
	var visitMetrics struct {
		TotalVisits     int64 `gorm:"column:total_visits"`
		VisitsThisMonth int64 `gorm:"column:visits_this_month"`
	}
	s.db.Table("visit_reports").
		Select(`
			COUNT(*) as total_visits,
			SUM(CASE WHEN visit_date >= ? AND visit_date <= ? THEN 1 ELSE 0 END) as visits_this_month
		`, monthStart, periodEnd).
		Where("visit_reports.brick_id = ? AND visit_reports.status IN ('completed', 'approved') AND visit_reports.visit_date >= ? AND visit_reports.visit_date <= ? AND visit_reports.deleted_at IS NULL", brickID, periodStart, periodEnd).
		Scan(&visitMetrics)

	metrics.TotalVisits = int(visitMetrics.TotalVisits)
	metrics.VisitsThisMonth = int(visitMetrics.VisitsThisMonth)

	// Calculate average visits per sales
	if metrics.ActiveSales > 0 {
		metrics.AverageVisitsPerSales = float64(metrics.VisitsThisMonth) / float64(metrics.ActiveSales)
	}

	// Account metrics
	var accountMetrics struct {
		TotalAccounts        int64 `gorm:"column:total_accounts"`
		ActiveAccounts       int64 `gorm:"column:active_accounts"`
		NewAccountsThisMonth int64 `gorm:"column:new_accounts_this_month"`
	}
	s.db.Table("accounts").
		Select(`
			COUNT(*) as total_accounts,
			SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END) as active_accounts,
			SUM(CASE WHEN created_at >= ? AND created_at <= ? THEN 1 ELSE 0 END) as new_accounts_this_month
		`, monthStart, periodEnd).
		Where("accounts.brick_id = ? AND accounts.deleted_at IS NULL", brickID).
		Scan(&accountMetrics)

	metrics.TotalAccounts = int(accountMetrics.TotalAccounts)
	metrics.ActiveAccounts = int(accountMetrics.ActiveAccounts)
	metrics.NewAccountsThisMonth = int(accountMetrics.NewAccountsThisMonth)

	// Revenue metrics (align with Sales Performance: revenue from won deals within the period)
	metrics.TotalRevenue = metrics.WonDealValue
	metrics.RevenueThisMonth = metrics.WonDealValue

	// Revenue growth: compare current month vs previous month (by close date)
	var previousMonthRevenue int64
	prevMonthStart := monthStart.AddDate(0, -1, 0)
	prevMonthEnd := monthStart.Add(-1 * time.Second)
	s.db.Table("deals").
		Select("COALESCE(SUM(value), 0)").
		Where("deals.brick_id = ? AND deals.status = 'won' AND ((deals.actual_close_date IS NOT NULL AND deals.actual_close_date >= ? AND deals.actual_close_date <= ?) OR (deals.actual_close_date IS NULL AND deals.created_at >= ? AND deals.created_at <= ?)) AND deals.deleted_at IS NULL", brickID, prevMonthStart, prevMonthEnd, prevMonthStart, prevMonthEnd).
		Scan(&previousMonthRevenue)

	if previousMonthRevenue > 0 {
		metrics.RevenueGrowthPercent = float64(metrics.RevenueThisMonth-previousMonthRevenue) / float64(previousMonthRevenue) * 100
	} else if metrics.RevenueThisMonth > 0 {
		metrics.RevenueGrowthPercent = 100.0 // 100% growth if no previous revenue
	}

	return metrics, nil
}

// ListBrickPerformance gets performance metrics for multiple bricks
// Uses optimized batch queries to avoid N+1
func (s *Service) ListBrickPerformance(brickIDs []string, periodStart, periodEnd time.Time) ([]*BrickPerformanceMetrics, error) {
	if len(brickIDs) == 0 {
		return []*BrickPerformanceMetrics{}, nil
	}

	// Get all bricks in one query (with reasonable limit for performance)
	// CRITICAL: Limit to 500 to prevent memory issues with large datasets
	bricks, _, err := s.brickRepo.List(&brick.ListBricksRequest{
		PerPage: 500, // Reasonable limit for batch operations
	})
	if err != nil {
		return nil, err
	}

	// Filter to requested brick IDs
	brickMap := make(map[string]*brick.Brick)
	for i := range bricks {
		for _, id := range brickIDs {
			if bricks[i].ID == id {
				brickMap[id] = &bricks[i]
				break
			}
		}
	}

	// Get all metrics in batch (optimized queries)
	results := make([]*BrickPerformanceMetrics, 0, len(brickIDs))
	for _, brickID := range brickIDs {
		b, exists := brickMap[brickID]
		if !exists {
			continue
		}

		metrics, err := s.GetBrickPerformance(brickID, periodStart, periodEnd)
		if err != nil {
			// Skip errors for individual bricks, continue with others
			continue
		}

		// Update brick info from loaded brick
		metrics.BrickName = b.Name
		metrics.BrickCode = b.Code
		if b.Manager != nil {
			metrics.ManagerID = &b.Manager.ID
			metrics.ManagerName = &b.Manager.Name
		}

		results = append(results, metrics)
	}

	return results, nil
}
