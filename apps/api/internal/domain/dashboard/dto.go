package dashboard

import "time"

// DashboardOverviewResponse represents dashboard overview data
type DashboardOverviewResponse struct {
	Period struct {
		Type  string    `json:"type"` // today, week, month, year
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	} `json:"period"`
	VisitStats struct {
		Total         int     `json:"total"`
		Completed     int     `json:"completed"`
		Pending       int     `json:"pending"`
		Approved      int     `json:"approved"`
		Rejected      int     `json:"rejected"`
		ChangePercent float64 `json:"change_percent"`
	} `json:"visit_stats"`
	AccountStats struct {
		Total         int     `json:"total"`
		Active        int     `json:"active"`
		Inactive      int     `json:"inactive"`
		ChangePercent float64 `json:"change_percent"`
	} `json:"account_stats"`
	ActivityStats struct {
		Total         int     `json:"total"`
		Visits        int     `json:"visits"`
		Calls         int     `json:"calls"`
		Emails        int     `json:"emails"`
		ChangePercent float64 `json:"change_percent"`
	} `json:"activity_stats"`

	// Target represents sales target vs achievement for the selected period.
	// This is used by the "Your target is incomplete" card in the dashboard UI.
	Target TargetStats `json:"target"`

	// Deals contains high‑level deal statistics for the selected period.
	Deals DealsStats `json:"deals"`

	// Revenue contains revenue statistics (typically based on won deals).
	Revenue RevenueStats `json:"revenue"`

	// LeadsBySource aggregates open deals/leads grouped by their source
	// (e.g. social, email, call, others) for the donut chart.
	LeadsBySource LeadsBySource `json:"leads_by_source"`

	// UpcomingTasks contains a small snapshot of upcoming tasks used by the
	// dashboard tasks widget.
	UpcomingTasks []DashboardTaskSummary `json:"upcoming_tasks"`

	// PipelineStages contains simplified pipeline stage distribution used by
	// the sales pipeline progress bar in the dashboard.
	PipelineStages []DashboardPipelineStageSummary `json:"pipeline_stages"`

	// LeadStats contains lead statistics for the selected period.
	LeadStats LeadStats `json:"lead_stats"`
}

// TargetStats represents sales target vs achievement metrics.
type TargetStats struct {
	TargetAmount          int64   `json:"target_amount"`
	TargetAmountFormatted string  `json:"target_amount_formatted"`
	AchievedAmount        int64   `json:"achieved_amount"`
	AchievedAmountFormatted string `json:"achieved_amount_formatted"`
	ProgressPercent       float64 `json:"progress_percent"`
	ChangePercent         float64 `json:"change_percent"`
}

// DealsStats represents high‑level deal metrics.
type DealsStats struct {
	TotalDeals int64  `json:"total_deals"`
	OpenDeals  int64  `json:"open_deals"`
	WonDeals   int64  `json:"won_deals"`
	LostDeals  int64  `json:"lost_deals"`
	TotalValue int64  `json:"total_value"`
	TotalValueFormatted string `json:"total_value_formatted"`
	ChangePercent float64 `json:"change_percent"`
}

// RevenueStats represents revenue metrics (derived from won deals).
type RevenueStats struct {
	TotalRevenue          int64   `json:"total_revenue"`
	TotalRevenueFormatted string  `json:"total_revenue_formatted"`
	ChangePercent         float64 `json:"change_percent"`
}

// LeadsBySourceEntry represents a single lead source bucket.
type LeadsBySourceEntry struct {
	Source string `json:"source"`
	Count  int64  `json:"count"`
}

// LeadsBySource aggregates total leads and distribution by source.
type LeadsBySource struct {
	Total    int64               `json:"total"`
	BySource []LeadsBySourceEntry `json:"by_source"`
}

// DashboardTaskSummary is a lightweight projection of a task for the dashboard.
type DashboardTaskSummary struct {
	ID       string     `json:"id"`
	Title    string     `json:"title"`
	Priority string     `json:"priority"`
	Status   string     `json:"status"`
	DueDate  *time.Time `json:"due_date,omitempty"`
}

// DashboardPipelineStageSummary represents simplified stage stats for dashboard.
type DashboardPipelineStageSummary struct {
	StageID             string  `json:"stage_id"`
	StageName           string  `json:"stage_name"`
	StageCode           string  `json:"stage_code"`
	StageColor          string  `json:"stage_color"`
	DealCount           int64   `json:"deal_count"`
	TotalValue          int64   `json:"total_value"`
	TotalValueFormatted string  `json:"total_value_formatted"`
	Percentage          float64 `json:"percentage"`
}

// LeadStats represents lead statistics.
type LeadStats struct {
	Total         int64   `json:"total"`
	New           int64   `json:"new"`
	Contacted     int64   `json:"contacted"`
	Qualified     int64   `json:"qualified"`
	Converted     int64   `json:"converted"`
	Lost          int64   `json:"lost"`
	ChangePercent float64 `json:"change_percent"`
}

// VisitStatisticsResponse represents visit statistics
type VisitStatisticsResponse struct {
	Period struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	} `json:"period"`
	Total        int     `json:"total"`
	Completed    int     `json:"completed"`
	Pending      int     `json:"pending"`
	Approved     int     `json:"approved"`
	Rejected     int     `json:"rejected"`
	ByStatus     map[string]int `json:"by_status"`
	ByDate       []DateStat     `json:"by_date"`
	ChangePercent float64 `json:"change_percent"`
}

// DateStat represents statistics for a specific date
type DateStat struct {
	Date      string `json:"date"`
	Count     int    `json:"count"`
	Completed int    `json:"completed"`
	Approved  int    `json:"approved"`
	Pending   int    `json:"pending"`
	Rejected  int    `json:"rejected"`
}

// PipelineSummaryResponse represents pipeline summary (placeholder for future)
type PipelineSummaryResponse struct {
	TotalDeals int64                          `json:"total_deals"`
	TotalValue int64                          `json:"total_value"`
	WonDeals   int64                          `json:"won_deals"`
	LostDeals  int64                          `json:"lost_deals"`
	OpenDeals  int64                          `json:"open_deals"`
	// ByStage contains all stages with their stats (including stages with 0 deals)
	ByStage []DashboardPipelineStageSummary `json:"by_stage"`
}

// TopAccountResponse represents top account data
type TopAccountResponse struct {
	Account struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
	VisitCount    int     `json:"visit_count"`
	ActivityCount int     `json:"activity_count"`
	LastVisitDate *time.Time `json:"last_visit_date,omitempty"`
}

// TopSalesRepResponse represents top sales rep data
type TopSalesRepResponse struct {
	SalesRep struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"sales_rep"`
	VisitCount                int     `json:"visit_count"`
	AccountCount              int     `json:"account_count"`
	ActivityCount             int     `json:"activity_count"`
	DealsClosed               int64   `json:"deals_closed"`
	ActualRevenue             int64   `json:"actual_revenue"`
	ActualRevenueFormatted    string  `json:"actual_revenue_formatted"`
	TargetAmount              int64   `json:"target_amount"`
	TargetAmountFormatted     string  `json:"target_amount_formatted"`
	TargetAchievementPercent  float64 `json:"target_achievement_percent"`
}

// RecentActivityResponse represents recent activity data
type RecentActivityResponse struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Account     *struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"account,omitempty"`
	Contact     *struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"contact,omitempty"`
	User        struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"user"`
	Timestamp   time.Time `json:"timestamp"`
}

// ActivityTrendsResponse represents activity trends by date
type ActivityTrendsResponse struct {
	Period struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	} `json:"period"`
	ByDate []ActivityDateStat `json:"by_date"`
}

// ActivityDateStat represents activity statistics for a specific date
type ActivityDateStat struct {
	Date   string `json:"date"`
	Visits int    `json:"visits"`
	Calls  int    `json:"calls"`
	Emails int    `json:"emails"`
	Total  int    `json:"total"`
}

// DashboardRequest represents request parameters for dashboard
type DashboardRequest struct {
	StartDate    string   `form:"start_date"`
	EndDate      string   `form:"end_date"`
	Period       string   `form:"period"` // today, week, month, year
	Limit        int      `form:"limit"`
	Offset       int      `form:"offset"`
	Days         int      `form:"days"`    // For upcoming visits
	Status       string   `form:"status"`  // For leads filtering
	TeamID       string   `form:"team_id"` // For sales manager filtering
	ScopedUserIDs []string `form:"-" json:"-"` // RBAC-resolved user IDs for data scoping
}

// ============================================================================
// Super Admin Dashboard Responses
// ============================================================================

// SuperAdminUsersByRoleResponse represents users by role statistics
type SuperAdminUsersByRoleResponse struct {
	UsersByRole []UsersByRoleEntry `json:"users_by_role"`
	TotalUsers  int64              `json:"total_users"`
	TotalActive int64              `json:"total_active"`
	TotalInactive int64            `json:"total_inactive"`
}

// UsersByRoleEntry represents user count for a specific role
type UsersByRoleEntry struct {
	RoleCode     string `json:"role_code"`
	RoleName     string `json:"role_name"`
	TotalUsers   int64  `json:"total_users"`
	ActiveUsers  int64  `json:"active_users"`
	InactiveUsers int64 `json:"inactive_users"`
}

// SuperAdminSystemActivityResponse represents system activity logs
type SuperAdminSystemActivityResponse struct {
	Activities   []SystemActivityEntry `json:"activities"`
	Total        int64                 `json:"total"`
	RecentErrors int                   `json:"recent_errors"`
	RecentWarnings int                 `json:"recent_warnings"`
}

// SystemActivityEntry represents a single system activity log entry
type SystemActivityEntry struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"` // user_login, user_created, error, warning
	Description string    `json:"description"`
	UserID      *string   `json:"user_id,omitempty"`
	UserName    *string   `json:"user_name,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// SuperAdminAIUsageResponse represents AI usage and cost statistics
type SuperAdminAIUsageResponse struct {
	TotalRequests        int64              `json:"total_requests"`
	RequestsToday        int64              `json:"requests_today"`
	RequestsThisWeek     int64              `json:"requests_this_week"`
	RequestsThisMonth    int64              `json:"requests_this_month"`
	EstimatedCost        AIUsageCost        `json:"estimated_cost"`
	SuccessRate          float64            `json:"success_rate"`
	FallbackRate         float64            `json:"fallback_rate"`
	AverageResponseTime  float64            `json:"average_response_time"`
}

// AIUsageCost represents cost estimation for different periods
type AIUsageCost struct {
	Today    float64 `json:"today"`
	Week     float64 `json:"week"`
	Month    float64 `json:"month"`
	Currency string  `json:"currency"`
}

// SuperAdminDataGrowthResponse represents data growth statistics
type SuperAdminDataGrowthResponse struct {
	Accounts DataGrowthStats `json:"accounts"`
	Leads    DataGrowthStats `json:"leads"`
	Deals    DataGrowthStats `json:"deals"`
}

// DataGrowthStats represents growth statistics for a data type
type DataGrowthStats struct {
	Total         int64   `json:"total"`
	GrowthPercent float64 `json:"growth_percent"`
	GrowthCount   int64   `json:"growth_count"`
	Period        string  `json:"period"`
}

// SuperAdminErrorSummaryResponse represents error and failed process summary
type SuperAdminErrorSummaryResponse struct {
	TotalErrors      int64                `json:"total_errors"`
	ErrorsToday      int64                `json:"errors_today"`
	ErrorsThisWeek   int64                `json:"errors_this_week"`
	ErrorsThisMonth  int64                `json:"errors_this_month"`
	FailedProcesses  []FailedProcessEntry `json:"failed_processes"`
	ErrorTypes       []ErrorTypeEntry     `json:"error_types"`
}

// FailedProcessEntry represents a failed process
type FailedProcessEntry struct {
	ProcessName string    `json:"process_name"`
	FailureCount int      `json:"failure_count"`
	LastFailure  time.Time `json:"last_failure"`
}

// ErrorTypeEntry represents error count by type
type ErrorTypeEntry struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

// ============================================================================
// Admin Dashboard Responses
// ============================================================================

// AdminTotalLeadsResponse represents total leads statistics
type AdminTotalLeadsResponse struct {
	Today       LeadPeriodStats `json:"today"`
	ThisMonth   LeadPeriodStats `json:"this_month"`
	ChangePercent float64       `json:"change_percent"`
}

// LeadPeriodStats represents lead statistics for a period
type LeadPeriodStats struct {
	Total     int64 `json:"total"`
	New       int64 `json:"new"`
	Contacted int64 `json:"contacted"`
	Qualified int64 `json:"qualified"`
	Converted int64 `json:"converted,omitempty"`
}

// AdminPipelineValueResponse represents pipeline value summary
type AdminPipelineValueResponse struct {
	TotalValue          int64  `json:"total_value"`
	TotalValueFormatted string `json:"total_value_formatted"`
	OpenDealsValue      int64  `json:"open_deals_value"`
	WonDealsValue       int64  `json:"won_deals_value"`
	LostDealsValue      int64  `json:"lost_deals_value"`
	ChangePercent       float64 `json:"change_percent"`
}

// AdminPendingApprovalsResponse represents pending approvals
type AdminPendingApprovalsResponse struct {
	Total         int64                `json:"total"`
	VisitReports  int64                `json:"visit_reports"`
	ExpenseReports int64               `json:"expense_reports"`
	Other         int64                `json:"other"`
	Items         []PendingApprovalItem `json:"items"`
}

// PendingApprovalItem represents a single pending approval item
type PendingApprovalItem struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"` // visit_report, expense_report, other
	Title       string    `json:"title"`
	SubmittedBy string    `json:"submitted_by"`
	SubmittedAt time.Time `json:"submitted_at"`
	Priority    string    `json:"priority"` // high, medium, low
}

// AdminTaskOverdueResponse represents global overdue tasks
type AdminTaskOverdueResponse struct {
	TotalOverdue    int64            `json:"total_overdue"`
	CriticalOverdue int64            `json:"critical_overdue"`
	Tasks           []OverdueTaskItem `json:"tasks"`
}

// OverdueTaskItem represents a single overdue task
type OverdueTaskItem struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	AssignedTo   string    `json:"assigned_to"`
	DueDate      time.Time `json:"due_date"`
	DaysOverdue  int       `json:"days_overdue"`
	Priority     string    `json:"priority"` // high, medium, low
}

// ============================================================================
// Sales Manager Dashboard Responses
// ============================================================================

// SalesManagerPipelineFunnelResponse represents pipeline funnel data
type SalesManagerPipelineFunnelResponse struct {
	Funnel        []FunnelStageEntry `json:"funnel"`
	ConversionRate float64           `json:"conversion_rate"`
}

// FunnelStageEntry represents a stage in the conversion funnel
type FunnelStageEntry struct {
	Stage     string  `json:"stage"`
	Count     int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

// SalesManagerTargetVsActualResponse represents target vs actual performance
type SalesManagerTargetVsActualResponse struct {
	Period      string              `json:"period"`
	Target      TargetMetrics        `json:"target"`
	Actual      ActualMetrics        `json:"actual"`
	Achievement AchievementMetrics   `json:"achievement"`
	Gap         GapMetrics          `json:"gap"`
}

// TargetMetrics represents target metrics
type TargetMetrics struct {
	Revenue int64 `json:"revenue"`
	Deals   int64 `json:"deals"`
	Visits  int64 `json:"visits"`
}

// ActualMetrics represents actual metrics
type ActualMetrics struct {
	Revenue int64 `json:"revenue"`
	Deals   int64 `json:"deals"`
	Visits  int64 `json:"visits"`
}

// AchievementMetrics represents achievement percentages
type AchievementMetrics struct {
	RevenuePercent float64 `json:"revenue_percent"`
	DealsPercent   float64 `json:"deals_percent"`
	VisitsPercent  float64 `json:"visits_percent"`
}

// GapMetrics represents gap between target and actual
type GapMetrics struct {
	Revenue int64 `json:"revenue"`
	Deals   int64 `json:"deals"`
	Visits  int64 `json:"visits"`
}

// SalesManagerVisitCompletionResponse represents visit completion statistics
type SalesManagerVisitCompletionResponse struct {
	TotalScheduled int64                    `json:"total_scheduled"`
	Completed      int64                    `json:"completed"`
	Pending        int64                    `json:"pending"`
	Missed         int64                    `json:"missed"`
	CompletionRate float64                  `json:"completion_rate"`
	BySalesRep     []VisitCompletionBySalesRep `json:"by_sales_rep"`
}

// VisitCompletionBySalesRep represents visit completion by sales rep
type VisitCompletionBySalesRep struct {
	SalesRepID     string  `json:"sales_rep_id"`
	SalesRepName   string  `json:"sales_rep_name"`
	Scheduled      int64   `json:"scheduled"`
	Completed      int64   `json:"completed"`
	CompletionRate float64 `json:"completion_rate"`
}

// SalesManagerDealsAtRiskResponse represents deals at risk
type SalesManagerDealsAtRiskResponse struct {
	TotalAtRisk int64            `json:"total_at_risk"`
	Deals        []DealAtRiskItem `json:"deals"`
}

// DealAtRiskItem represents a deal at risk
type DealAtRiskItem struct {
	ID                  string    `json:"id"`
	Name                string    `json:"name"`
	Value               int64     `json:"value"`
	Stage               string    `json:"stage"`
	DaysWithoutActivity int       `json:"days_without_activity"`
	RiskReason          string    `json:"risk_reason"` // no_activity, stale, low_probability
	AssignedTo          string    `json:"assigned_to"`
	LastActivity        time.Time `json:"last_activity"`
}

// ============================================================================
// Sales Manager Team Draft Approvals Responses
// ============================================================================

// SalesManagerTeamDraftApprovalsResponse aggregates pending/draft items from all team
// members managed by the requesting sales manager.
type SalesManagerTeamDraftApprovalsResponse struct {
	Total     int64                          `json:"total"`
	Leads     SalesManagerDraftLeadsData     `json:"leads"`
	Pipeline  SalesManagerDraftPipelineData  `json:"pipeline"`
	Schedules SalesManagerDraftSchedulesData `json:"schedules"`
	Visits    SalesManagerDraftVisitsData    `json:"visits"`
	Tasks     SalesManagerDraftTasksData     `json:"tasks"`
}

// SalesManagerDraftLeadsData represents new/uncontacted leads from team members.
type SalesManagerDraftLeadsData struct {
	Total int64           `json:"total"`
	Items []DraftLeadItem `json:"items"`
}

// DraftLeadItem is a single lead item needing manager attention.
type DraftLeadItem struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Company    string    `json:"company"`
	Status     string    `json:"status"`
	AssignedTo string    `json:"assigned_to"`
	CreatedAt  time.Time `json:"created_at"`
}

// SalesManagerDraftPipelineData represents open pipeline deals from team members.
type SalesManagerDraftPipelineData struct {
	Total int64              `json:"total"`
	Items []DraftPipelineItem `json:"items"`
}

// DraftPipelineItem is a single open deal item needing manager attention.
type DraftPipelineItem struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Stage      string    `json:"stage"`
	Value      int64     `json:"value"`
	AssignedTo string    `json:"assigned_to"`
	CreatedAt  time.Time `json:"created_at"`
}

// SalesManagerDraftSchedulesData represents pending schedules from team members.
type SalesManagerDraftSchedulesData struct {
	Total int64               `json:"total"`
	Items []DraftScheduleItem `json:"items"`
}

// DraftScheduleItem is a single pending schedule item needing manager confirmation.
type DraftScheduleItem struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	ScheduledAt time.Time `json:"scheduled_at"`
	AssignedTo  string    `json:"assigned_to"`
	CreatedAt   time.Time `json:"created_at"`
}

// SalesManagerDraftVisitsData represents submitted/draft visit reports from team members.
type SalesManagerDraftVisitsData struct {
	Total int64              `json:"total"`
	Items []DraftVisitItem   `json:"items"`
}

// DraftVisitItem is a single draft or submitted visit report awaiting approval.
type DraftVisitItem struct {
	ID          string    `json:"id"`
	Purpose     string    `json:"purpose"`
	Status      string    `json:"status"` // draft, submitted
	VisitDate   time.Time `json:"visit_date"`
	AssignedTo  string    `json:"assigned_to"`
	CreatedAt   time.Time `json:"created_at"`
}

// SalesManagerDraftTasksData represents pending tasks from team members.
type SalesManagerDraftTasksData struct {
	Total int64           `json:"total"`
	Items []DraftTaskItem `json:"items"`
}

// DraftTaskItem is a single pending task needing manager attention.
type DraftTaskItem struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	Priority   string     `json:"priority"`
	DueDate    *time.Time `json:"due_date"`
	AssignedTo string     `json:"assigned_to"`
	CreatedAt  time.Time  `json:"created_at"`
}

// ============================================================================
// Sales Dashboard Responses
// ============================================================================

// SalesTodayTasksResponse represents today's tasks for sales user
type SalesTodayTasksResponse struct {
	Total    int64              `json:"total"`
	Completed int64             `json:"completed"`
	Pending  int64              `json:"pending"`
	Overdue  int64              `json:"overdue"`
	Tasks    []SalesTaskItem    `json:"tasks"`
}

// SalesTaskItem represents a task item for sales dashboard
type SalesTaskItem struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	DueDate   time.Time `json:"due_date"`
	DueTime   *string   `json:"due_time,omitempty"`
	Status    string    `json:"status"` // pending, completed, overdue
	Priority  string    `json:"priority"` // high, medium, low
	RelatedTo *TaskRelatedTo `json:"related_to,omitempty"`
}

// TaskRelatedTo represents related entity for a task
type TaskRelatedTo struct {
	Type string `json:"type"` // account, deal, lead
	ID   string `json:"id"`
	Name string `json:"name"`
}

// SalesAssignedLeadsResponse represents assigned leads for sales user
type SalesAssignedLeadsResponse struct {
	Total     int64           `json:"total"`
	New       int64           `json:"new"`
	Contacted int64           `json:"contacted"`
	Qualified int64           `json:"qualified"`
	Converted int64           `json:"converted"`
	Leads     []SalesLeadItem `json:"leads"`
}

// SalesLeadItem represents a lead item for sales dashboard
type SalesLeadItem struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Company     string     `json:"company"`
	Status      string     `json:"status"`
	AssignedDate time.Time `json:"assigned_date"`
	LastContact *time.Time `json:"last_contact,omitempty"`
}

// SalesUpcomingVisitsResponse represents upcoming visits for sales user
type SalesUpcomingVisitsResponse struct {
	Total     int64              `json:"total"`
	Today     int64              `json:"today"`
	ThisWeek  int64              `json:"this_week"`
	NextWeek  int64              `json:"next_week"`
	Visits    []SalesVisitItem   `json:"visits"`
}

// SalesVisitItem represents a visit item for sales dashboard
type SalesVisitItem struct {
	ID            string    `json:"id"`
	AccountName   string    `json:"account_name"`
	ScheduledDate time.Time `json:"scheduled_date"`
	ScheduledTime *string   `json:"scheduled_time,omitempty"`
	Purpose       string    `json:"purpose"`
	Status        string    `json:"status"` // scheduled, completed, cancelled
}

// SalesRemindersResponse represents reminders for sales user
type SalesRemindersResponse struct {
	Total    int64              `json:"total"`
	Unread   int64              `json:"unread"`
	Reminders []ReminderItem    `json:"reminders"`
}

// ReminderItem represents a reminder item
type ReminderItem struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"` // task_due, visit_upcoming, deal_reminder, system
	Title     string          `json:"title"`
	Message   string          `json:"message"`
	RelatedTo *ReminderRelatedTo `json:"related_to,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	Read      bool            `json:"read"`
}

// ReminderRelatedTo represents related entity for a reminder
type ReminderRelatedTo struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// ============================================================================
// Analyst Dashboard Responses
// ============================================================================

// AnalystRevenueTrendResponse represents revenue trend over time
type AnalystRevenueTrendResponse struct {
	Period        string              `json:"period"`
	TotalRevenue  int64               `json:"total_revenue"`
	Trend         []RevenueTrendEntry `json:"trend"`
	GrowthPercent float64             `json:"growth_percent"`
	AverageDaily  int64               `json:"average_daily"`
}

// RevenueTrendEntry represents revenue for a specific date
type RevenueTrendEntry struct {
	Date    string `json:"date"`
	Revenue int64  `json:"revenue"`
	Deals   int64  `json:"deals"`
}

// AnalystConversionRateResponse represents conversion rate statistics
type AnalystConversionRateResponse struct {
	Period         string                    `json:"period"`
	TotalLeads     int64                     `json:"total_leads"`
	ConvertedLeads int64                     `json:"converted_leads"`
	ConversionRate float64                   `json:"conversion_rate"`
	BySource       []ConversionRateBySource  `json:"by_source"`
	Trend          []ConversionRateTrendEntry `json:"trend"`
}

// ConversionRateBySource represents conversion rate by source
type ConversionRateBySource struct {
	Source         string  `json:"source"`
	Leads          int64   `json:"leads"`
	Converted      int64   `json:"converted"`
	ConversionRate float64 `json:"conversion_rate"`
}

// ConversionRateTrendEntry represents conversion rate trend
type ConversionRateTrendEntry struct {
	Date          string  `json:"date"`
	ConversionRate float64 `json:"conversion_rate"`
}

// AnalystSalesVelocityResponse represents sales velocity metrics
type AnalystSalesVelocityResponse struct {
	Period              string              `json:"period"`
	AverageSalesCycleDays int               `json:"average_sales_cycle_days"`
	AverageDealValue    int64               `json:"average_deal_value"`
	SalesVelocity       float64             `json:"sales_velocity"`
	ByStage             []SalesVelocityByStage `json:"by_stage"`
}

// SalesVelocityByStage represents average days for a stage
type SalesVelocityByStage struct {
	Stage       string `json:"stage"`
	AverageDays int    `json:"average_days"`
}

// AnalystAIInsightsResponse represents AI-generated insights
type AnalystAIInsightsResponse struct {
	Insights     []AIInsightItem `json:"insights"`
	TotalInsights int64          `json:"total_insights"`
}

// AIInsightItem represents a single AI insight
type AIInsightItem struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"` // trend, anomaly, recommendation
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Impact      string    `json:"impact"` // high, medium, low
	GeneratedAt time.Time `json:"generated_at"`
}

// ============================================================================
// Mobile Dashboard Responses
// ============================================================================

// MobileDashboardOverviewResponse represents simplified dashboard data for mobile sales user
type MobileDashboardOverviewResponse struct {
	Target MobileTargetSummary `json:"target"`
}

// MobileTargetSummary represents target information for mobile dashboard
type MobileTargetSummary struct {
	TargetAmount          int64   `json:"target_amount"`
	TargetAmountFormatted string  `json:"target_amount_formatted"`
	AchievedAmount        int64   `json:"achieved_amount"`
	AchievedAmountFormatted string `json:"achieved_amount_formatted"`
	ProgressPercent       float64 `json:"progress_percent"`
	Period                string  `json:"period"` // Format: "YYYY-MM"
	BrickName             *string `json:"brick_name,omitempty"`
}

// MobileVisitSummary represents visit summary for mobile dashboard
type MobileVisitSummary struct {
	TotalToday      int `json:"total_today"`
	Active          int `json:"active"`
	Completed       int `json:"completed"`
	PendingApproval int `json:"pending_approval"`
}

// MobileTaskSummary represents task summary for mobile dashboard
type MobileTaskSummary struct {
	Total    int `json:"total"`
	Today    int `json:"today"`
	Overdue  int `json:"overdue"`
	Upcoming int `json:"upcoming"`
}

// MobileVisitResponse represents a visit for mobile dashboard
type MobileVisitResponse struct {
	ID              string                 `json:"id"`
	Type            string                 `json:"type"` // "account", "deal", or "lead"
	Purpose         string                 `json:"purpose"`
	AccountID       *string                `json:"account_id,omitempty"`
	AccountName     *string                `json:"account_name,omitempty"`
	AccountAddress  *string                `json:"account_address,omitempty"`
	ContactID       *string                `json:"contact_id,omitempty"`
	ContactName     *string                `json:"contact_name,omitempty"`
	DealID          *string                `json:"deal_id,omitempty"`
	DealTitle       *string                `json:"deal_title,omitempty"`
	LeadID          *string                `json:"lead_id,omitempty"`
	LeadName        *string                `json:"lead_name,omitempty"`
	VisitDate       string                 `json:"visit_date"` // Format: "YYYY-MM-DD"
	VisitTime       *string                `json:"visit_time,omitempty"` // Format: "HH:MM"
	Status          string                 `json:"status"`
	CheckInTime     *time.Time             `json:"check_in_time,omitempty"`
	CheckInLocation *MobileVisitLocation   `json:"check_in_location,omitempty"`
	CheckOutTime    *time.Time             `json:"check_out_time,omitempty"`
	CheckOutLocation *MobileVisitLocation `json:"check_out_location,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// MobileVisitLocation represents location information for mobile visit
type MobileVisitLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   *string `json:"address,omitempty"`
}

// MobileVisitsListResponse represents list of visits for mobile dashboard
type MobileVisitsListResponse struct {
	Visits  []MobileVisitResponse `json:"visits"`
	Total   int                   `json:"total"`
	HasMore bool                  `json:"has_more"`
}

// MobileTaskResponse represents a task for mobile dashboard
type MobileTaskResponse struct {
	ID          string                 `json:"id"`
	Title       string                 `json:"title"`
	Description *string                `json:"description,omitempty"`
	Type        string                 `json:"type"` // general, call, email, meeting, follow_up
	DueDate     *string                `json:"due_date,omitempty"` // Format: "YYYY-MM-DD"
	DueTime     *string                `json:"due_time,omitempty"` // Format: "HH:MM"
	Priority    string                 `json:"priority"`
	Status      string                 `json:"status"`
	AssignedBy  *MobileTaskAssignee    `json:"assigned_by,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	IsOverdue   bool                   `json:"is_overdue"`
}

// MobileTaskAssignee represents task assignee information
type MobileTaskAssignee struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// MobileTasksListResponse represents list of tasks for mobile dashboard
type MobileTasksListResponse struct {
	Tasks   []MobileTaskResponse `json:"tasks"`
	Total   int                  `json:"total"`
	HasMore bool                 `json:"has_more"`
}

// MobileDashboardRequest represents request parameters for mobile dashboard
type MobileDashboardRequest struct {
	Period    string `form:"period"`     // today, week, month (default: today)
	StartDate string `form:"start_date"` // ISO 8601 date string
	EndDate   string `form:"end_date"`   // ISO 8601 date string
	Status    string `form:"status"`      // For visits: active, completed, all (default: all)
	Filter    string `form:"filter"`     // For tasks: today, week, overdue, all (default: all)
	Date      string `form:"date"`       // Filter visits by date (ISO 8601)
	Limit     int    `form:"limit"`       // Max items (default: 5 for visits, 3 for tasks)
}
