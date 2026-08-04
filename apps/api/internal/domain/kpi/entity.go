package kpi

import "time"

// SalesRepScorecard contains the raw KPI metrics for one sales rep.
// Phase 1 keeps this layer focused on directly verifiable numbers from the
// existing CRM schema.
type SalesRepScorecard struct {
	TotalDealsClosed        int64    `json:"totalDealsClosed"`
	TotalRevenue            int64    `json:"totalRevenue"`
	DealsCreated            int64    `json:"dealsCreated"`
	ConversionRate          *float64 `json:"conversionRate"`
	AverageDealValue        *float64 `json:"averageDealValue"`
	VisitCompleted          int64    `json:"visitCompleted"`
	VisitPlanned            int64    `json:"visitPlanned"`
	VisitCompliance         *float64 `json:"visitCompliance"`
	TasksCompleted          int64    `json:"tasksCompleted"`
	OverdueTaskRate         *float64 `json:"overdueTaskRate"`
	RevenueTargetAttainment *float64 `json:"revenueTargetAttainment"`
	DealTargetAttainment    *float64 `json:"dealTargetAttainment"`
	PipelineMovementScore   int64    `json:"pipelineMovementScore"`
}

// KPITrend captures comparison against the previous period.
type KPITrend struct {
	PreviousCompositeScore *float64 `json:"previousCompositeScore,omitempty"`
	Delta                  *float64 `json:"delta,omitempty"`
	Direction              string   `json:"direction"`
}

// KPITargetGapItem represents a target vs actual comparison item.
type KPITargetGapItem struct {
	Target     int64    `json:"target"`
	Actual     int64    `json:"actual"`
	GapPercent *float64 `json:"gapPercent,omitempty"`
	Status     string   `json:"status"`
}

// KPITargetGap groups target gap comparisons.
type KPITargetGap struct {
	Revenue KPITargetGapItem `json:"revenue"`
	Deals   KPITargetGapItem `json:"deals"`
}

// KPIEvaluation captures the composite score and supporting evaluation data.
type KPIEvaluation struct {
	CompositeScore float64      `json:"compositeScore"`
	Grade          string       `json:"grade"`
	Trend          KPITrend     `json:"trend"`
	TargetGap      KPITargetGap `json:"targetGap"`
}

// SalesManagerEvaluation captures the composite score and trend for manager KPI responses.
type SalesManagerEvaluation struct {
	CompositeScore float64  `json:"compositeScore"`
	Grade          string   `json:"grade"`
	Trend          KPITrend `json:"trend"`
}

// KPIDiagnostic captures actionable insight output.
type KPIDiagnostic struct {
	Code     string `json:"code"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

// SalesRepScope identifies the scorecard context.
type SalesRepScope struct {
	UserID    string `json:"userId"`
	Role      string `json:"role"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

// SalesRepMeta contains data-quality metadata for the KPI response.
type SalesRepMeta struct {
	BrickMissingCount  int64     `json:"brickMissingCount"`
	BrickInferredCount int64     `json:"brickInferredCount"`
	GeneratedAt        time.Time `json:"generatedAt"`
}

// SalesRepScorecardResponse is the phase-1 response payload for sales rep KPI.
type SalesRepScorecardResponse struct {
	Scope       SalesRepScope     `json:"scope"`
	Scorecard   SalesRepScorecard `json:"scorecard"`
	Evaluation  KPIEvaluation     `json:"evaluation"`
	Diagnostics []KPIDiagnostic   `json:"diagnostics,omitempty"`
	Meta        SalesRepMeta      `json:"meta"`
}

// SalesManagerScope identifies the manager scorecard context.
type SalesManagerScope struct {
	ManagerID string   `json:"managerId"`
	Role      string   `json:"role"`
	StartDate string   `json:"startDate"`
	EndDate   string   `json:"endDate"`
	Bricks    []string `json:"bricks"`
}

// SalesManagerTeamSummary contains aggregate team KPI values.
type SalesManagerTeamSummary struct {
	TotalRepsCount       int64    `json:"totalRepsCount"`
	TotalDealsClosed     int64    `json:"totalDealsClosed"`
	TotalRevenue         int64    `json:"totalRevenue"`
	TeamConversionRate   *float64 `json:"teamConversionRate"`
	TeamVisitCompliance  *float64 `json:"teamVisitCompliance"`
	TeamOverdueTaskRate  *float64 `json:"teamOverdueTaskRate"`
	TeamTargetAttainment *float64 `json:"teamTargetAttainment"`
}

// SalesManagerTeamBreakdownItem represents per-rep performance in the manager view.
type SalesManagerTeamBreakdownItem struct {
	UserID         string   `json:"userId"`
	Name           string   `json:"name"`
	CompositeScore float64  `json:"compositeScore"`
	Grade          string   `json:"grade"`
	TotalRevenue   int64    `json:"totalRevenue"`
	ConversionRate *float64 `json:"conversionRate,omitempty"`
	Rank           int      `json:"rank"`
}

// SalesManagerBrickBreakdownItem represents per-brick performance.
type SalesManagerBrickBreakdownItem struct {
	BrickID             string   `json:"brickId"`
	Name                string   `json:"name"`
	CoveragePenetration *float64 `json:"coveragePenetration"`
	TotalRevenue        int64    `json:"totalRevenue"`
	RepsCount           int64    `json:"repsCount"`
	CompositeScore      float64  `json:"compositeScore"`
}

// SalesManagerTopBottomPerformers contains top and bottom performer IDs.
type SalesManagerTopBottomPerformers struct {
	Top    []string `json:"top"`
	Bottom []string `json:"bottom"`
}

// SalesManagerMeta contains data-quality metadata for the manager KPI response.
type SalesManagerMeta struct {
	BrickMissingCount  int64     `json:"brickMissingCount"`
	BrickInferredCount int64     `json:"brickInferredCount"`
	GeneratedAt        time.Time `json:"generatedAt"`
}

// SalesManagerScorecardResponse is the KPI response payload for sales managers.
type SalesManagerScorecardResponse struct {
	Scope               SalesManagerScope                `json:"scope"`
	TeamSummary         SalesManagerTeamSummary          `json:"teamSummary"`
	Evaluation          SalesManagerEvaluation           `json:"evaluation"`
	TeamBreakdown       []SalesManagerTeamBreakdownItem  `json:"teamBreakdown,omitempty"`
	BrickBreakdown      []SalesManagerBrickBreakdownItem `json:"brickBreakdown,omitempty"`
	TopBottomPerformers SalesManagerTopBottomPerformers  `json:"topBottomPerformers"`
	Diagnostics         []KPIDiagnostic                  `json:"diagnostics,omitempty"`
	Meta                SalesManagerMeta                 `json:"meta"`
}
