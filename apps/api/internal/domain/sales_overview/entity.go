package sales_overview

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
)

// SalesPerformanceDetail represents detailed performance metrics for a user
type SalesPerformanceDetail struct {
	UserID                string  `json:"user_id"`
	User                  *user.UserResponse `json:"user,omitempty"`
	PeriodStart           time.Time `json:"period_start"`
	PeriodEnd             time.Time `json:"period_end"`
	TotalRevenue          int64   `json:"total_revenue"` // in smallest currency unit (sen)
	TotalRevenueFormatted string  `json:"total_revenue_formatted"`
	WonDeals              int     `json:"won_deals"`
	TotalDeals            int     `json:"total_deals"`
	LostDeals             int     `json:"lost_deals"`
	OpenDeals             int     `json:"open_deals"`
	ConversionRate        float64 `json:"conversion_rate"` // percentage
	AverageDealValue      float64 `json:"average_deal_value"`
	AverageDealValueFormatted string `json:"average_deal_value_formatted"`
	VisitsCompleted       int     `json:"visits_completed"`
	TasksCompleted        int     `json:"tasks_completed"`
	TotalTasks            int     `json:"total_tasks"`
	TaskCompletionRate    float64 `json:"task_completion_rate"` // percentage
}

// SalesRepDetail represents comprehensive detail for sales rep detail page
type SalesRepDetail struct {
	UserID                string  `json:"user_id"`
	User                  *user.UserResponse `json:"user,omitempty"`
	PeriodStart           *time.Time `json:"period_start,omitempty"`
	PeriodEnd             *time.Time `json:"period_end,omitempty"`
	Statistics            *SalesRepStatistics `json:"statistics,omitempty"`
}

// SalesRepStatistics represents statistics for sales rep
type SalesRepStatistics struct {
	TotalRevenue          int64   `json:"total_revenue"`
	TotalRevenueFormatted string  `json:"total_revenue_formatted"`
	DealsClosed           int     `json:"deals_closed"`
	VisitsCompleted       int     `json:"visits_completed"`
	TasksCompleted        int     `json:"tasks_completed"`
	ConversionRate        float64 `json:"conversion_rate"`
	AverageDealValue      float64 `json:"average_deal_value"`
	AverageDealValueFormatted string `json:"average_deal_value_formatted"`
	PeriodComparison      *PeriodComparison `json:"period_comparison,omitempty"`
}

// PeriodComparison represents comparison with previous period
type PeriodComparison struct {
	RevenueChange        float64 `json:"revenue_change"` // percentage
	RevenueChangeDirection string `json:"revenue_change_direction"` // up, down, same
	DealsChange          float64 `json:"deals_change"` // percentage
	DealsChangeDirection string  `json:"deals_change_direction"` // up, down, same
}

// CheckInLocation represents a check-in location for sales rep
type CheckInLocation struct {
	VisitNumber   int       `json:"visit_number"`
	VisitReportID string    `json:"visit_report_id"`
	VisitDate     time.Time `json:"visit_date"`
	CheckInTime   time.Time `json:"check_in_time"`
	Location      *Location `json:"location"`
	Account       *AccountRef `json:"account,omitempty"`
	Purpose       string    `json:"purpose"`
}

// Location represents GPS location
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address,omitempty"`
}

// AccountRef represents account reference
type AccountRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// SalesPerformanceListResponse represents list of sales performance (management overview)
type SalesPerformanceListResponse struct {
	UserID                string  `json:"user_id"`
	UserName              string  `json:"user_name"`
	UserEmail             string  `json:"user_email"`
	AvatarURL             string  `json:"avatar_url,omitempty"`
	TotalRevenue          int64   `json:"total_revenue"`
	TotalRevenueFormatted string  `json:"total_revenue_formatted"`
	DealsClosed           int     `json:"deals_closed"`
	VisitsCompleted       int     `json:"visits_completed"`
	TasksCompleted              int     `json:"tasks_completed"`
	ConversionRate              float64 `json:"conversion_rate"`
	TargetAmount                int64   `json:"target_amount"`
	TargetAmountFormatted       string  `json:"target_amount_formatted"`
	TargetAchievementPercentage float64 `json:"target_achievement_percentage"`
}

// SalesRepCheckInLocationsResponse represents check-in locations response
type SalesRepCheckInLocationsResponse struct {
	SalesRep       *user.UserResponse `json:"sales_rep,omitempty"`
	CheckInLocations []CheckInLocation `json:"check_in_locations"`
	TotalVisits    int64              `json:"total_visits"`
	Period         *PeriodRange       `json:"period,omitempty"`
}

// PeriodRange represents date range
type PeriodRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// MonthlySalesData represents data for a single month in the chart
type MonthlySalesData struct {
	Month        int    `json:"month"`
	MonthName    string `json:"month_name"`
	Year         int    `json:"year"`
	TotalRevenue int64  `json:"total_revenue"`
	TotalDeals   int    `json:"total_deals"`
	TotalVisits  int    `json:"total_visits"`
	TotalTasks   int    `json:"total_tasks"`
	TargetAmount int64  `json:"target_amount"`
}

// MonthlySalesOverviewResponse represents the chart data
type MonthlySalesOverviewResponse struct {
	MonthlyData  []MonthlySalesData `json:"monthly_data"`
	TotalRevenue int64              `json:"total_revenue"`
	TotalDeals   int                `json:"total_deals"`
	TotalVisits  int                `json:"total_visits"`
	TotalTasks   int                `json:"total_tasks"`
}

// ListSalesPerformanceRequest represents list sales performance query parameters
type ListSalesPerformanceRequest struct {
	Page          int      `form:"page" binding:"omitempty,min=1"`
	PerPage       int      `form:"per_page" binding:"omitempty,min=1,max=100"`
	Period        string   `form:"period" binding:"omitempty"`     // YYYY-MM-DD format
	StartDate     string   `form:"start_date" binding:"omitempty"` // YYYY-MM-DD format
	EndDate       string   `form:"end_date" binding:"omitempty"`   // YYYY-MM-DD format
	Search        string   `form:"search" binding:"omitempty"`
	BrickID       string   `form:"brick_id" binding:"omitempty,uuid"` // Filter by brick_id
	SortBy        string   `form:"sort_by" binding:"omitempty,oneof=revenue deals visits tasks name target achievement"`
	Order         string   `form:"order" binding:"omitempty,oneof=asc desc"`
	SortOrder     string   `form:"sort_order" binding:"omitempty,oneof=asc desc"`
	ScopedUserIDs []string `form:"-" json:"-"` // Injected by scope middleware for team-based filtering
}

// GetSalesPerformanceDetailRequest represents get sales performance detail query parameters
type GetSalesPerformanceDetailRequest struct {
	Period    string `form:"period" binding:"omitempty"` // YYYY-MM-DD format
	StartDate string `form:"start_date" binding:"omitempty"` // YYYY-MM-DD format
	EndDate   string `form:"end_date" binding:"omitempty"` // YYYY-MM-DD format
}

// GetSalesRepDetailRequest represents get sales rep detail query parameters
type GetSalesRepDetailRequest struct {
	Period    string `form:"period" binding:"omitempty"` // YYYY-MM-DD format
	StartDate string `form:"start_date" binding:"omitempty"` // YYYY-MM-DD format
	EndDate   string `form:"end_date" binding:"omitempty"` // YYYY-MM-DD format
}

// GetSalesRepCheckInLocationsRequest represents get sales rep check-in locations query parameters
type GetSalesRepCheckInLocationsRequest struct {
	StartDate string `form:"start_date" binding:"omitempty"` // YYYY-MM-DD format
	EndDate   string `form:"end_date" binding:"omitempty"` // YYYY-MM-DD format
	Page      int    `form:"page" binding:"omitempty,min=1"`
	PerPage   int    `form:"per_page" binding:"omitempty,min=1,max=100"`
}