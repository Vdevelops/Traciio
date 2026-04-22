package report

import "time"

// VisitReportReportResponse represents visit report report data
type VisitReportReportResponse struct {
	Period struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	} `json:"period"`
	Summary struct {
		Total        int     `json:"total"`
		Completed    int     `json:"completed"`
		Pending      int     `json:"pending"`
		Approved     int     `json:"approved"`
		Rejected     int     `json:"rejected"`
	} `json:"summary"`
	ByAccount    []AccountStat `json:"by_account"`
	BySalesRep   []SalesRepStat `json:"by_sales_rep"`
	ByDate       []DateStat     `json:"by_date"`
	ByStatus     map[string]int `json:"by_status"`
}

// AccountStat represents statistics for an account
type AccountStat struct {
	Account struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
	VisitCount int `json:"visit_count"`
}

// SalesRepStat represents statistics for a sales rep
type SalesRepStat struct {
	SalesRep struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"sales_rep"`
	VisitCount int `json:"visit_count"`
}

// DateStat represents statistics for a specific date
type DateStat struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// PipelineReportResponse represents pipeline report
type PipelineReportResponse struct {
	Period struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	} `json:"period"`
	Summary struct {
		TotalDeals        int     `json:"total_deals"`
		TotalValue        float64 `json:"total_value"`
		WonDeals          int     `json:"won_deals"`
		WonValue          float64 `json:"won_value"`
		LostDeals         int     `json:"lost_deals"`
		LostValue         float64 `json:"lost_value"`
		OpenDeals         int     `json:"open_deals"`
		OpenValue         float64 `json:"open_value"`
		ExpectedRevenue   float64 `json:"expected_revenue"`
	} `json:"summary"`
	ByStage map[string]int `json:"by_stage"`
	Deals   []DealReportItem `json:"deals,omitempty"` // Individual deals for Sales Funnel table
}

// DealReportItem represents a deal in the sales funnel report
type DealReportItem struct {
	ID                string     `json:"id"`
	CompanyName       string     `json:"company_name"`
	ContactName         string     `json:"contact_name"`
	ContactEmail      string     `json:"contact_email"`
	Stage             string     `json:"stage"`
	StageCode         string     `json:"stage_code"`
	Value             float64    `json:"value"`
	Probability       int        `json:"probability"`
	ExpectedRevenue   float64    `json:"expected_revenue"`
	CreationDate      time.Time  `json:"creation_date"`
	ExpectedCloseDate *time.Time `json:"expected_close_date"`
	TeamMember        string     `json:"team_member"`
	ProgressToWon     int        `json:"progress_to_won"` // Calculated based on stage and probability
	LastInteractedOn  *time.Time `json:"last_interacted_on"` // From activities
	NextStep          string     `json:"next_step"` // From deal notes or metadata
}

// SalesPerformanceReportResponse represents sales performance report
type SalesPerformanceReportResponse struct {
	Period struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	} `json:"period"`
	BySalesRep []SalesPerformanceStat `json:"by_sales_rep"`
	Summary    struct {
		TotalVisits      int     `json:"total_visits"`
		TotalAccounts    int     `json:"total_accounts"`
		AverageVisitsPerAccount float64 `json:"average_visits_per_account"`
	} `json:"summary"`
}

// SalesPerformanceStat represents sales performance statistics
type SalesPerformanceStat struct {
	SalesRep struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Email string `json:"email"`
	} `json:"sales_rep"`
	VisitCount    int     `json:"visit_count"`
	AccountCount  int     `json:"account_count"`
	ActivityCount int     `json:"activity_count"`
	CompletionRate float64 `json:"completion_rate"`
}

// AccountActivityReportResponse represents account activity report
type AccountActivityReportResponse struct {
	Period struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	} `json:"period"`
	AccountID   string `json:"account_id"`
	AccountName string `json:"account_name"`
	Summary     struct {
		TotalVisits   int `json:"total_visits"`
		TotalActivities int `json:"total_activities"`
		TotalContacts int `json:"total_contacts"`
	} `json:"summary"`
	Activities  []ActivityDetail `json:"activities"`
	Visits      []VisitDetail     `json:"visits"`
}

// ActivityDetail represents activity detail
type ActivityDetail struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
	User        struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"user"`
}

// VisitDetail represents visit detail
type VisitDetail struct {
	ID        string    `json:"id"`
	VisitDate time.Time `json:"visit_date"`
	Purpose   string    `json:"purpose"`
	Status    string    `json:"status"`
	SalesRep  struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"sales_rep"`
}

// ReportRequest represents request parameters for reports
type ReportRequest struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	AccountID string `form:"account_id"`
	SalesRepID string `form:"sales_rep_id"`
	Status    string `form:"status"`
	Limit     int    `form:"limit"`
}

