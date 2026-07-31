package kpi

import (
	"errors"
	"fmt"
	"sort"
	"time"

	domainkpi "github.com/gilabs/crm-healthcare/api/internal/domain/kpi"
	brickdomain "github.com/gilabs/crm-healthcare/api/internal/domain/brick"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"gorm.io/gorm"
)

var ErrKPIUserNotFound = errors.New("kpi user not found")

type Service struct {
	kpiRepo           interfaces.KPIRepository
	userRepo          interfaces.UserRepository
	monthlyTargetRepo interfaces.MonthlyTargetRepository
	brickRepo         interfaces.BrickRepository
}

func NewService(
	kpiRepo interfaces.KPIRepository,
	userRepo interfaces.UserRepository,
	monthlyTargetRepo interfaces.MonthlyTargetRepository,
	brickRepo interfaces.BrickRepository,
) *Service {
	return &Service{
		kpiRepo:           kpiRepo,
		userRepo:          userRepo,
		monthlyTargetRepo: monthlyTargetRepo,
		brickRepo:         brickRepo,
	}
}

func parseKPIPeriod(startDate, endDate string) (time.Time, time.Time, error) {
	loc := response.GetTimezoneWIB()

	start, err := time.ParseInLocation("2006-01-02", startDate, loc)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("parse start date: %w", err)
	}
	end, err := time.ParseInLocation("2006-01-02", endDate, loc)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("parse end date: %w", err)
	}

	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, loc)
	end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, loc)

	if end.Before(start) {
		return time.Time{}, time.Time{}, fmt.Errorf("end date before start date")
	}

	return start, end, nil
}

func pct(numerator, denominator int64) *float64 {
	if denominator == 0 {
		return nil
	}
	value := (float64(numerator) / float64(denominator)) * 100
	return &value
}

func normalizeRoleCode(roleCode string) string {
	switch roleCode {
	case "sales":
		return "sales_rep"
	case "sales_manager":
		return "sales_manager"
	default:
		return roleCode
	}
}

func (s *Service) GetSalesRepScorecard(userID, startDate, endDate string) (*domainkpi.SalesRepScorecardResponse, error) {
	start, end, err := parseKPIPeriod(startDate, endDate)
	if err != nil {
		return nil, err
	}

	currentUser, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrKPIUserNotFound
		}
		return nil, err
	}

	roleCode := ""
	if currentUser.Role != nil {
		roleCode = currentUser.Role.Code
	}

	currentScorecard, meta, err := s.buildRawSalesRepScorecard(currentUser, start, end)
	if err != nil {
		return nil, err
	}

	var previousScorecard *domainkpi.SalesRepScorecard
	if true {
		prevStart, prevEnd := previousPeriodRange(start, end)
		if prevScorecard, _, prevErr := s.buildRawSalesRepScorecard(currentUser, prevStart, prevEnd); prevErr == nil {
			previousScorecard = &prevScorecard
		}
	}

	revenueTarget, err := s.monthlyTargetRepo.GetProratedTargetForPeriod(userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	if revenueTarget > 0 {
		attainment := (float64(currentScorecard.TotalRevenue) / revenueTarget) * 100
		currentScorecard.RevenueTargetAttainment = &attainment
	}

	evaluation, diagnostics := s.buildSalesRepEvaluation(currentScorecard, previousScorecard)

	responseScorecard := currentScorecard

	return &domainkpi.SalesRepScorecardResponse{
		Scope: domainkpi.SalesRepScope{
			UserID:    userID,
			Role:      normalizeRoleCode(roleCode),
			StartDate: startDate,
			EndDate:   endDate,
		},
		Scorecard: responseScorecard,
		Evaluation: evaluation,
		Diagnostics: diagnostics,
		Meta: domainkpi.SalesRepMeta{
			BrickMissingCount:  meta.BrickMissingCount,
			BrickInferredCount: meta.BrickInferredCount,
			GeneratedAt:        meta.GeneratedAt,
		},
	}, nil
}

func (s *Service) GetSalesManagerScorecard(req *domainkpi.GetSalesManagerScorecardRequest) (*domainkpi.SalesManagerScorecardResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("missing manager request")
	}
	start, end, err := parseKPIPeriod(req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}
	managerID := req.ManagerID
	if managerID == "" {
		return nil, fmt.Errorf("missing manager id")
	}

	currentUser, err := s.userRepo.FindByID(managerID)
	if err != nil {
		return nil, err
	}

	bricks, _, err := s.brickRepo.List(&brickdomain.ListBricksRequest{ManagerID: &managerID, Page: 1, PerPage: 100})
	if err != nil {
		return nil, err
	}
	if req.BrickID != "" {
		filtered := make([]brickdomain.Brick, 0, 1)
		for _, brick := range bricks {
			if brick.ID == req.BrickID {
				filtered = append(filtered, brick)
				break
			}
		}
		bricks = filtered
	}

	teamUsers := make([]user.User, 0)
	brickIDs := make([]string, 0, len(bricks))
	for _, brick := range bricks {
		brickIDs = append(brickIDs, brick.ID)
		members, memberErr := s.brickRepo.GetSalesByBrickID(brick.ID)
		if memberErr != nil {
			return nil, memberErr
		}
		teamUsers = append(teamUsers, members...)
	}
	teamUsers = uniqueUsers(teamUsers)
	brickIDs = uniqueStrings(brickIDs)

	teamSummary := domainkpi.SalesManagerTeamSummary{TotalRepsCount: int64(len(teamUsers))}
	teamBreakdown := make([]domainkpi.SalesManagerTeamBreakdownItem, 0, len(teamUsers))
	brickBreakdown := make([]domainkpi.SalesManagerBrickBreakdownItem, 0, len(brickIDs))
	brickMissingCount := int64(0)
	brickInferredCount := int64(0)
	totalDealsCreated := int64(0)
	totalVisitCompleted := int64(0)
	totalVisitPlanned := int64(0)
	overdueRateSum := 0.0
	overdueRateCount := 0
	type repResult struct {
		user       user.User
		scorecard  domainkpi.SalesRepScorecard
		evaluation domainkpi.KPIEvaluation
	}
	repResults := make([]repResult, 0, len(teamUsers))

	for _, rep := range teamUsers {
		raw, meta, rawErr := s.buildRawSalesRepScorecard(&rep, start, end)
		if rawErr != nil {
			return nil, rawErr
		}
		prevStart, prevEnd := previousPeriodRange(start, end)
		var previous *domainkpi.SalesRepScorecard
		if prevRaw, _, prevErr := s.buildRawSalesRepScorecard(&rep, prevStart, prevEnd); prevErr == nil {
			previous = &prevRaw
		}
		raw.RevenueTargetAttainment = nil
		if target, targetErr := s.monthlyTargetRepo.GetProratedTargetForPeriod(rep.ID, start.Format("2006-01-02"), end.Format("2006-01-02")); targetErr == nil && target > 0 {
			attainment := (float64(raw.TotalRevenue) / target) * 100
			raw.RevenueTargetAttainment = &attainment
		}
		eval, _ := s.buildSalesRepEvaluation(raw, previous)
		repResults = append(repResults, repResult{user: rep, scorecard: raw, evaluation: eval})

		teamSummary.TotalDealsClosed += raw.TotalDealsClosed
		teamSummary.TotalRevenue += raw.TotalRevenue
		totalDealsCreated += raw.DealsCreated
		totalVisitCompleted += raw.VisitCompleted
		totalVisitPlanned += raw.VisitPlanned
		if raw.OverdueTaskRate != nil {
			overdueRateSum += *raw.OverdueTaskRate
			overdueRateCount++
		}
		if raw.RevenueTargetAttainment != nil {
			brickInferredCount += meta.BrickInferredCount
			brickMissingCount += meta.BrickMissingCount
		}
	}

	teamTargetTotal := int64(0)
	for _, rep := range teamUsers {
		if target, targetErr := s.monthlyTargetRepo.GetProratedTargetForPeriod(rep.ID, start.Format("2006-01-02"), end.Format("2006-01-02")); targetErr == nil {
			teamTargetTotal += int64(target + 0.5)
		}
	}
	if teamTargetTotal > 0 {
		attainment := (float64(teamSummary.TotalRevenue) / float64(teamTargetTotal)) * 100
		teamSummary.TeamTargetAttainment = &attainment
	}
	if totalDealsCreated > 0 {
		teamSummary.TeamConversionRate = pct(teamSummary.TotalDealsClosed, totalDealsCreated)
	}
	if totalVisitPlanned > 0 {
		teamSummary.TeamVisitCompliance = pct(totalVisitCompleted, totalVisitPlanned)
	}
	if overdueRateCount > 0 {
		avgOverdue := overdueRateSum / float64(overdueRateCount)
		teamSummary.TeamOverdueTaskRate = &avgOverdue
	}

	for _, brick := range bricks {
		customersWithInteraction, custErr := s.kpiRepo.CountCustomersWithInteractionInBrick(brick.ID, start, end)
		if custErr != nil {
			return nil, custErr
		}
		registeredCustomers, regErr := s.kpiRepo.CountRegisteredCustomersInBrick(brick.ID)
		if regErr != nil {
			return nil, regErr
		}
		var coverage *float64
		if registeredCustomers > 0 {
			v := (float64(customersWithInteraction) / float64(registeredCustomers)) * 100
			coverage = &v
		}
		brickBreakdown = append(brickBreakdown, domainkpi.SalesManagerBrickBreakdownItem{
			BrickID:             brick.ID,
			Name:                brick.Name,
			CoveragePenetration: coverage,
			TotalRevenue:        teamSummary.TotalRevenue,
			RepsCount:           int64(len(teamUsers)),
			CompositeScore:      0,
		})
	}

	weights := s.salesManagerWeights()
	values := map[string]*float64{
		"team_target_attainment": nil,
		"team_conversion_rate":   teamSummary.TeamConversionRate,
		"territory_coverage":     nil,
		"team_visit_compliance":  teamSummary.TeamVisitCompliance,
		"team_overdue_task_rate":  teamSummary.TeamOverdueTaskRate,
		"brick_pipeline_movement": nil,
	}
	if teamSummary.TeamTargetAttainment != nil {
		values["team_target_attainment"] = teamSummary.TeamTargetAttainment
	}
	currentScore := computeWeightedScore(values, weights)
	var previousScore *float64
	if req.CompareWithPrevious {
		previousScore = &currentScore
	}

	diagnostics := []domainkpi.KPIDiagnostic{{Code: "DATA_QUALITY_ISSUE", Severity: "info", Message: "Brick inference and missing data were evaluated"}}
	if brickMissingCount > 0 {
		diagnostics = append(diagnostics, domainkpi.KPIDiagnostic{Code: "DATA_QUALITY_ISSUE", Severity: "info", Message: "Ada record tanpa brick_id"})
	}

	teamBreakdown = make([]domainkpi.SalesManagerTeamBreakdownItem, 0, len(repResults))
	for i, item := range repResults {
		teamBreakdown = append(teamBreakdown, domainkpi.SalesManagerTeamBreakdownItem{
			UserID:         item.user.ID,
			Name:           item.user.Name,
			CompositeScore: item.evaluation.CompositeScore,
			Grade:          item.evaluation.Grade,
			TotalRevenue:   item.scorecard.TotalRevenue,
			ConversionRate: item.scorecard.ConversionRate,
			Rank:           i + 1,
		})
	}

	sort.Slice(teamBreakdown, func(i, j int) bool { return teamBreakdown[i].CompositeScore > teamBreakdown[j].CompositeScore })
	for idx := range teamBreakdown {
		teamBreakdown[idx].Rank = idx + 1
	}

	topBottom := domainkpi.SalesManagerTopBottomPerformers{}
	if len(teamBreakdown) > 0 {
		topBottom.Top = []string{teamBreakdown[0].UserID}
		topBottom.Bottom = []string{teamBreakdown[len(teamBreakdown)-1].UserID}
	}

	return &domainkpi.SalesManagerScorecardResponse{
		Scope: domainkpi.SalesManagerScope{
			ManagerID: managerID,
			Role:      normalizeRoleCode(currentUser.Role.Code),
			StartDate: req.StartDate,
			EndDate:   req.EndDate,
			Bricks:    brickIDs,
		},
		TeamSummary: teamSummary,
		Evaluation:  s.buildManagerEvaluation(currentScore, previousScore),
		TeamBreakdown: teamBreakdown,
		BrickBreakdown: brickBreakdown,
		TopBottomPerformers: topBottom,
		Diagnostics: diagnostics,
		Meta: domainkpi.SalesManagerMeta{
			BrickMissingCount:  brickMissingCount,
			BrickInferredCount: brickInferredCount,
			GeneratedAt:        time.Now().In(response.GetTimezoneWIB()),
		},
	}, nil
}
