package ai

// InsightType represents type of AI insight
type InsightType string

const (
	// Existing insight types
	InsightTypeVisitReport InsightType = "visit_report"
	InsightTypeDeal        InsightType = "deal"
	InsightTypeContact     InsightType = "contact"
	InsightTypeAccount     InsightType = "account"
	InsightTypePipeline    InsightType = "pipeline"

	// New AI module insight types (based on docs/ai-modules)
	InsightTypeSalesPerformance InsightType = "sales_performance" // Sales performance analysis

	InsightTypeBrickManagement InsightType = "brick_management" // Brick/territory management
	InsightTypeProductAnalysis InsightType = "product_analysis" // Product analysis
	InsightTypeGroups          InsightType = "groups"           // Group segmentation
	InsightTypeTarget          InsightType = "target"           // Target/quota management
	InsightTypeSchedule        InsightType = "schedule"         // Schedule planning
)

// VisitReportInsight represents AI insight for visit report
type VisitReportInsight struct {
	Summary         string   `json:"summary"`
	ActionItems     []string `json:"action_items"`
	Sentiment       string   `json:"sentiment"` // positive, neutral, negative
	KeyPoints       []string `json:"key_points"`
	Recommendations []string `json:"recommendations"`
}

// DealInsight represents AI insight for deal
type DealInsight struct {
	WinProbability  float64  `json:"win_probability"`
	NextSteps       []string `json:"next_steps"`
	RiskFactors     []string `json:"risk_factors"`
	Recommendations []string `json:"recommendations"`
	ConfidenceLevel string   `json:"confidence_level"` // high, medium, low
}

// ContactInsight represents AI insight for contact
type ContactInsight struct {
	CommunicationStyle string   `json:"communication_style"`
	Preferences        []string `json:"preferences"`
	BestContactTime    string   `json:"best_contact_time"`
	Recommendations    []string `json:"recommendations"`
}

// AccountInsight represents AI insight for account
type AccountInsight struct {
	HealthScore     int      `json:"health_score"` // 0-100
	RiskIndicators  []string `json:"risk_indicators"`
	Opportunities   []string `json:"opportunities"`
	Recommendations []string `json:"recommendations"`
}

// PipelineInsight represents AI insight for pipeline
type PipelineInsight struct {
	Forecast        float64  `json:"forecast"`
	ConfidenceLevel string   `json:"confidence_level"`
	Trends          []string `json:"trends"`
	Recommendations []string `json:"recommendations"`
}

// SalesPerformanceInsight represents AI insight for sales performance analysis
type SalesPerformanceInsight struct {
	RevenueActual      float64            `json:"revenue_actual"`
	QuotaAttainment    float64            `json:"quota_attainment"` // Percentage of quota achieved
	ConversionRates    map[string]float64 `json:"conversion_rates"` // Conversion rates per stage
	AverageDealSize    float64            `json:"average_deal_size"`
	WinRate            float64            `json:"win_rate"`
	SalesCycleDuration int                `json:"sales_cycle_duration"` // Days
	PipelineValue      float64            `json:"pipeline_value"`
	WeightedPipeline   float64            `json:"weighted_pipeline"` // Probability adjusted
	TopPerformers      []string           `json:"top_performers"`
	LowPerformers      []string           `json:"low_performers"`
	Trends             []string           `json:"trends"`
	Recommendations    []string           `json:"recommendations"`
	Alerts             []string           `json:"alerts"`
}

// BrickManagementInsight represents AI insight for brick/territory management
type BrickManagementInsight struct {
	BrickID            string   `json:"brick_id"`
	BrickName          string   `json:"brick_name"`
	RevenuePerBrick    float64  `json:"revenue_per_brick"`
	PenetrationRate    float64  `json:"penetration_rate"` // Accounts reached vs total
	VisitFrequency     float64  `json:"visit_frequency"`
	TravelCost         float64  `json:"travel_cost"`
	PotentialScore     int      `json:"potential_score"` // 0-100
	PerformanceGap     float64  `json:"performance_gap"` // Gap between potential and actual
	UnderservedAreas   []string `json:"underserved_areas"`
	ResourceAllocation []string `json:"resource_allocation"` // Recommendations
	Recommendations    []string `json:"recommendations"`
}

// ProductAnalysisInsight represents AI insight for product analysis
type ProductAnalysisInsight struct {
	ProductID           string            `json:"product_id"`
	ProductName         string            `json:"product_name"`
	RevenueContribution float64           `json:"revenue_contribution"`
	MarginContribution  float64           `json:"margin_contribution"`
	GrowthRate          float64           `json:"growth_rate"` // MoM or YoY
	VolumesSold         int64             `json:"volumes_sold"`
	InventoryTurns      float64           `json:"inventory_turns"`
	Affinity            []ProductAffinity `json:"affinity"`            // Cross-sell opportunities
	RecommendedActions  []string          `json:"recommended_actions"` // ramp-up, promo, discontinue
	Recommendations     []string          `json:"recommendations"`
}

// ProductAffinity represents product affinity for cross-sell
type ProductAffinity struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Affinity    float64 `json:"affinity"` // Correlation score
}

// GroupsInsight represents AI insight for group segmentation
type GroupsInsight struct {
	GroupID             string   `json:"group_id"`
	GroupName           string   `json:"group_name"`
	GroupType           string   `json:"group_type"` // static or dynamic
	MemberCount         int      `json:"member_count"`
	RevenuePerGroup     float64  `json:"revenue_per_group"`
	PenetrationRate     float64  `json:"penetration_rate"`
	ARPU                float64  `json:"arpu"`              // Average Revenue Per Unit
	UnvisitedMembers    []string `json:"unvisited_members"` // Members not visited in X days
	CampaignSuggestions []string `json:"campaign_suggestions"`
	Recommendations     []string `json:"recommendations"`
}

// TargetInsight represents AI insight for target/quota management
type TargetInsight struct {
	Period               string             `json:"period"` // month, quarter, year
	TargetAmount         float64            `json:"target_amount"`
	ActualAmount         float64            `json:"actual_amount"`
	AttainmentRate       float64            `json:"attainment_rate"`      // Percentage
	ProjectedAttainment  float64            `json:"projected_attainment"` // Based on current trend
	GapAmount            float64            `json:"gap_amount"`           // Target - Actual
	DaysRemaining        int                `json:"days_remaining"`
	DailyRequired        float64            `json:"daily_required"` // Daily amount needed to meet target
	DistributionByRegion map[string]float64 `json:"distribution_by_region"`
	DistributionByTeam   map[string]float64 `json:"distribution_by_team"`
	Alerts               []string           `json:"alerts"`
	Recommendations      []string           `json:"recommendations"`
}

// ScheduleInsight represents AI insight for schedule planning
type ScheduleInsight struct {
	Period              string   `json:"period"` // day, week, month
	PlannedVisits       int      `json:"planned_visits"`
	ActualVisits        int      `json:"actual_visits"`
	ComplianceRate      float64  `json:"compliance_rate"`  // Actual / Planned
	UtilizationRate     float64  `json:"utilization_rate"` // Hours used / Available hours
	AvgTravelTime       float64  `json:"avg_travel_time"`  // Minutes per visit
	VisitsPerDay        float64  `json:"visits_per_day"`
	MissedVisits        []string `json:"missed_visits"`        // Account names with missed visits
	UnderservedAccounts []string `json:"underserved_accounts"` // Accounts below required frequency
	RouteOptimization   []string `json:"route_optimization"`   // Suggestions
	Recommendations     []string `json:"recommendations"`
}

// ChatMessage represents a single chat message in conversation history
type ChatMessage struct {
	Role    string `json:"role"` // "user" or "assistant"
	Content string `json:"content"`
}
