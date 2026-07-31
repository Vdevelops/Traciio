package kpi

// GetSalesRepScorecardRequest carries request filters for a sales rep KPI.
// Query parsing will be handled in the handler layer in later phases.
type GetSalesRepScorecardRequest struct {
	UserID              string `form:"userId" binding:"omitempty,uuid"`
	StartDate           string `form:"startDate" binding:"required"`
	EndDate             string `form:"endDate" binding:"required"`
	CompareWithPrevious bool   `form:"compareWithPrevious"`
	TeamMemberIDs       []string `form:"-" json:"-"`
}

// GetSalesManagerScorecardRequest is reserved for the manager KPI shape.
// Phase 1 does not wire the endpoint yet, but the DTO is kept for symmetry.
type GetSalesManagerScorecardRequest struct {
	ManagerID            string `form:"managerId" binding:"omitempty,uuid"`
	StartDate            string `form:"startDate" binding:"required"`
	EndDate              string `form:"endDate" binding:"required"`
	BrickID              string `form:"brickId" binding:"omitempty,uuid"`
	IncludeTeamBreakdown bool   `form:"includeTeamBreakdown"`
	CompareWithPrevious  bool   `form:"compareWithPrevious"`
	TeamMemberIDs        []string `form:"-" json:"-"`
}
