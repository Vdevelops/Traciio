package kpi

import (
	"errors"
	"fmt"
	"sort"
	"time"

	brickdomain "github.com/gilabs/crm-healthcare/api/internal/domain/brick"
	domainkpi "github.com/gilabs/crm-healthcare/api/internal/domain/kpi"
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

func (s *Service) GetSalesRepScorecard(userID, startDate, endDate string, compareWithPrevious bool) (*domainkpi.SalesRepScorecardResponse, error) {
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
	averageDealBenchmark, err := s.peerAverageDealValue(currentUser, start, end)
	if err != nil {
		return nil, err
	}

	var previousScorecard *domainkpi.SalesRepScorecard
	var previousAverageDealBenchmark *float64
	if compareWithPrevious {
		prevStart, prevEnd := previousPeriodRange(start, end)
		if prevScorecard, _, prevErr := s.buildRawSalesRepScorecard(currentUser, prevStart, prevEnd); prevErr == nil {
			previousScorecard = &prevScorecard
		}
		if benchmark, benchmarkErr := s.peerAverageDealValue(currentUser, prevStart, prevEnd); benchmarkErr == nil {
			previousAverageDealBenchmark = benchmark
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

	evaluation, diagnostics := s.buildSalesRepEvaluation(currentScorecard, previousScorecard, averageDealBenchmark, previousAverageDealBenchmark, meta)

	responseScorecard := roundSalesRepScorecard(currentScorecard)

	return &domainkpi.SalesRepScorecardResponse{
		Scope: domainkpi.SalesRepScope{
			UserID:    userID,
			Role:      normalizeRoleCode(roleCode),
			StartDate: startDate,
			EndDate:   endDate,
		},
		Scorecard:   responseScorecard,
		Evaluation:  evaluation,
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
	roleCode := ""
	if currentUser.Role != nil {
		roleCode = currentUser.Role.Code
	}

	bricks, err := s.resolveManagerBricks(currentUser)
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
	brickMembers := make(map[string][]user.User, len(bricks))
	for _, brick := range bricks {
		brickIDs = append(brickIDs, brick.ID)
		members, memberErr := s.brickRepo.GetSalesByBrickID(brick.ID)
		if memberErr != nil {
			return nil, memberErr
		}
		for _, member := range members {
			if member.ID == managerID {
				continue
			}
			if member.Role == nil || member.Role.Code != "sales" {
				continue
			}
			brickMembers[brick.ID] = append(brickMembers[brick.ID], member)
			teamUsers = append(teamUsers, member)
		}
	}
	teamUsers = uniqueUsers(teamUsers)
	brickIDs = uniqueStrings(brickIDs)

	teamSummary := domainkpi.SalesManagerTeamSummary{TotalRepsCount: int64(len(teamUsers))}
	teamBreakdown := make([]domainkpi.SalesManagerTeamBreakdownItem, 0, len(teamUsers))
	brickBreakdown := make([]domainkpi.SalesManagerBrickBreakdownItem, 0, len(brickIDs))
	brickMissingCount := int64(0)
	brickInferredCount := int64(0)
	totalDealsCreated := int64(0)
	totalPipelineMovement := int64(0)
	totalVisitCompleted := int64(0)
	totalVisitPlanned := int64(0)
	totalTasksCreated := int64(0)
	totalOverdueTasks := int64(0)
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
		var previousAverageDealBenchmark *float64
		if prevRaw, _, prevErr := s.buildRawSalesRepScorecard(&rep, prevStart, prevEnd); prevErr == nil {
			previous = &prevRaw
		}
		averageDealBenchmark, benchmarkErr := s.peerAverageDealValue(&rep, start, end)
		if benchmarkErr != nil {
			return nil, benchmarkErr
		}
		if benchmark, benchmarkErr := s.peerAverageDealValue(&rep, prevStart, prevEnd); benchmarkErr == nil {
			previousAverageDealBenchmark = benchmark
		}
		raw.RevenueTargetAttainment = nil
		if target, targetErr := s.monthlyTargetRepo.GetProratedTargetForPeriod(rep.ID, start.Format("2006-01-02"), end.Format("2006-01-02")); targetErr == nil && target > 0 {
			attainment := (float64(raw.TotalRevenue) / target) * 100
			raw.RevenueTargetAttainment = &attainment
		}
		eval, _ := s.buildSalesRepEvaluation(raw, previous, averageDealBenchmark, previousAverageDealBenchmark, meta)
		repResults = append(repResults, repResult{user: rep, scorecard: raw, evaluation: eval})

		teamSummary.TotalDealsClosed += raw.TotalDealsClosed
		teamSummary.TotalRevenue += raw.TotalRevenue
		totalDealsCreated += raw.DealsCreated
		totalPipelineMovement += raw.PipelineMovementScore
		totalVisitCompleted += raw.VisitCompleted
		totalVisitPlanned += raw.VisitPlanned
		tasksCreated, taskErr := s.kpiRepo.CountTasksCreated(rep.ID, start, end)
		if taskErr != nil {
			return nil, taskErr
		}
		overdueTasks, overdueErr := s.kpiRepo.CountOverdueTasks(rep.ID, start, end)
		if overdueErr != nil {
			return nil, overdueErr
		}
		totalTasksCreated += tasksCreated
		totalOverdueTasks += overdueTasks
		brickInferredCount += meta.BrickInferredCount
		brickMissingCount += meta.BrickMissingCount
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
	teamSummary.TeamOverdueTaskRate = pct(totalOverdueTasks, totalTasksCreated)

	coverageValues := make([]*float64, 0, len(bricks))
	repScoreByID := make(map[string]float64, len(repResults))
	for _, item := range repResults {
		repScoreByID[item.user.ID] = item.evaluation.CompositeScore
	}
	for _, brick := range bricks {
		members := brickMembers[brick.ID]
		var brickRevenue int64
		var brickScoreTotal float64
		brickScoreCount := 0
		for _, member := range members {
			revenue, revenueErr := s.kpiRepo.SumWonRevenue(member.ID, start, end)
			if revenueErr != nil {
				return nil, revenueErr
			}
			brickRevenue += revenue
			if score, exists := repScoreByID[member.ID]; exists {
				brickScoreTotal += score
				brickScoreCount++
			}
		}
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
		brickCompositeScore := 0.0
		if brickScoreCount > 0 {
			brickCompositeScore = brickScoreTotal / float64(brickScoreCount)
		}
		coverageValues = append(coverageValues, coverage)
		brickBreakdown = append(brickBreakdown, domainkpi.SalesManagerBrickBreakdownItem{
			BrickID:             brick.ID,
			Name:                brick.Name,
			CoveragePenetration: coverage,
			TotalRevenue:        brickRevenue,
			RepsCount:           int64(len(members)),
			CompositeScore:      brickCompositeScore,
		})
	}

	weights := s.salesManagerWeights()
	values := map[string]*float64{
		"team_target_attainment":  nil,
		"team_conversion_rate":    normalizePercentage(teamSummary.TeamConversionRate),
		"territory_coverage":      normalizePercentage(averageFloatPointers(coverageValues)),
		"team_visit_compliance":   normalizePercentage(teamSummary.TeamVisitCompliance),
		"team_overdue_task_rate":  normalizeInverseRate(teamSummary.TeamOverdueTaskRate),
		"brick_pipeline_movement": normalizePipelineMovement(totalPipelineMovement, totalDealsCreated),
	}
	if teamSummary.TeamTargetAttainment != nil {
		values["team_target_attainment"] = normalizePercentage(teamSummary.TeamTargetAttainment)
	}
	currentScore := computeWeightedScore(values, weights)
	var previousScore *float64
	if req.CompareWithPrevious {
		prevStart, prevEnd := previousPeriodRange(start, end)
		if score, scoreErr := s.computeManagerCompositeScore(teamUsers, bricks, prevStart, prevEnd); scoreErr == nil {
			previousScore = &score
		}
	}

	diagnostics := make([]domainkpi.KPIDiagnostic, 0)
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
			Role:      normalizeRoleCode(roleCode),
			StartDate: req.StartDate,
			EndDate:   req.EndDate,
			Bricks:    brickIDs,
		},
		TeamSummary:         roundSalesManagerTeamSummary(teamSummary),
		Evaluation:          s.buildManagerEvaluation(currentScore, previousScore),
		TeamBreakdown:       roundSalesManagerTeamBreakdown(teamBreakdown),
		BrickBreakdown:      roundSalesManagerBrickBreakdown(brickBreakdown),
		TopBottomPerformers: topBottom,
		Diagnostics:         diagnostics,
		Meta: domainkpi.SalesManagerMeta{
			BrickMissingCount:  brickMissingCount,
			BrickInferredCount: brickInferredCount,
			GeneratedAt:        time.Now().In(response.GetTimezoneWIB()),
		},
	}, nil
}

func (s *Service) resolveManagerBricks(manager *user.User) ([]brickdomain.Brick, error) {
	if manager == nil {
		return nil, fmt.Errorf("missing manager user")
	}
	bricks, _, err := s.brickRepo.List(&brickdomain.ListBricksRequest{ManagerID: &manager.ID, Page: 1, PerPage: 100})
	if err != nil {
		return nil, err
	}
	if manager.BrickID != nil && *manager.BrickID != "" {
		ownBrick, ownBrickErr := s.brickRepo.FindByID(*manager.BrickID)
		if ownBrickErr != nil {
			return nil, ownBrickErr
		}
		bricks = append(bricks, *ownBrick)
	}
	return uniqueBricks(bricks), nil
}

func averageFloatPointers(values []*float64) *float64 {
	total := 0.0
	count := 0
	for _, value := range values {
		if value == nil {
			continue
		}
		total += *value
		count++
	}
	if count == 0 {
		return nil
	}
	avg := total / float64(count)
	return &avg
}

func roundSalesRepScorecard(scorecard domainkpi.SalesRepScorecard) domainkpi.SalesRepScorecard {
	scorecard.ConversionRate = roundFloatPtr(scorecard.ConversionRate)
	scorecard.AverageDealValue = roundFloatPtr(scorecard.AverageDealValue)
	scorecard.VisitCompliance = roundFloatPtr(scorecard.VisitCompliance)
	scorecard.OverdueTaskRate = roundFloatPtr(scorecard.OverdueTaskRate)
	scorecard.RevenueTargetAttainment = roundFloatPtr(scorecard.RevenueTargetAttainment)
	scorecard.DealTargetAttainment = roundFloatPtr(scorecard.DealTargetAttainment)
	return scorecard
}

func roundSalesManagerTeamSummary(summary domainkpi.SalesManagerTeamSummary) domainkpi.SalesManagerTeamSummary {
	summary.TeamConversionRate = roundFloatPtr(summary.TeamConversionRate)
	summary.TeamVisitCompliance = roundFloatPtr(summary.TeamVisitCompliance)
	summary.TeamOverdueTaskRate = roundFloatPtr(summary.TeamOverdueTaskRate)
	summary.TeamTargetAttainment = roundFloatPtr(summary.TeamTargetAttainment)
	return summary
}

func roundSalesManagerTeamBreakdown(items []domainkpi.SalesManagerTeamBreakdownItem) []domainkpi.SalesManagerTeamBreakdownItem {
	for idx := range items {
		items[idx].CompositeScore = roundFloat(items[idx].CompositeScore)
		items[idx].ConversionRate = roundFloatPtr(items[idx].ConversionRate)
	}
	return items
}

func roundSalesManagerBrickBreakdown(items []domainkpi.SalesManagerBrickBreakdownItem) []domainkpi.SalesManagerBrickBreakdownItem {
	for idx := range items {
		items[idx].CoveragePenetration = roundFloatPtr(items[idx].CoveragePenetration)
		items[idx].CompositeScore = roundFloat(items[idx].CompositeScore)
	}
	return items
}

func uniqueBricks(bricks []brickdomain.Brick) []brickdomain.Brick {
	seen := make(map[string]struct{}, len(bricks))
	result := make([]brickdomain.Brick, 0, len(bricks))
	for _, brick := range bricks {
		if brick.ID == "" {
			continue
		}
		if _, exists := seen[brick.ID]; exists {
			continue
		}
		seen[brick.ID] = struct{}{}
		result = append(result, brick)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}

func (s *Service) computeManagerCompositeScore(teamUsers []user.User, bricks []brickdomain.Brick, start, end time.Time) (float64, error) {
	summary := domainkpi.SalesManagerTeamSummary{}
	totalDealsCreated := int64(0)
	totalPipelineMovement := int64(0)
	totalVisitCompleted := int64(0)
	totalVisitPlanned := int64(0)
	totalTasksCreated := int64(0)
	totalOverdueTasks := int64(0)
	teamTargetTotal := int64(0)

	for _, rep := range teamUsers {
		raw, _, err := s.buildRawSalesRepScorecard(&rep, start, end)
		if err != nil {
			return 0, err
		}
		target, err := s.revenueTargetForPeriod(rep.ID, start, end)
		if err != nil {
			return 0, err
		}
		teamTargetTotal += target
		summary.TotalDealsClosed += raw.TotalDealsClosed
		summary.TotalRevenue += raw.TotalRevenue
		totalDealsCreated += raw.DealsCreated
		totalPipelineMovement += raw.PipelineMovementScore
		totalVisitCompleted += raw.VisitCompleted
		totalVisitPlanned += raw.VisitPlanned
		tasksCreated, taskErr := s.kpiRepo.CountTasksCreated(rep.ID, start, end)
		if taskErr != nil {
			return 0, taskErr
		}
		overdueTasks, overdueErr := s.kpiRepo.CountOverdueTasks(rep.ID, start, end)
		if overdueErr != nil {
			return 0, overdueErr
		}
		totalTasksCreated += tasksCreated
		totalOverdueTasks += overdueTasks
	}

	if teamTargetTotal > 0 {
		attainment := (float64(summary.TotalRevenue) / float64(teamTargetTotal)) * 100
		summary.TeamTargetAttainment = &attainment
	}
	summary.TeamConversionRate = pct(summary.TotalDealsClosed, totalDealsCreated)
	summary.TeamVisitCompliance = pct(totalVisitCompleted, totalVisitPlanned)
	summary.TeamOverdueTaskRate = pct(totalOverdueTasks, totalTasksCreated)

	coverageValues := make([]*float64, 0, len(bricks))
	for _, brick := range bricks {
		customersWithInteraction, err := s.kpiRepo.CountCustomersWithInteractionInBrick(brick.ID, start, end)
		if err != nil {
			return 0, err
		}
		registeredCustomers, err := s.kpiRepo.CountRegisteredCustomersInBrick(brick.ID)
		if err != nil {
			return 0, err
		}
		if registeredCustomers == 0 {
			coverageValues = append(coverageValues, nil)
			continue
		}
		coverage := (float64(customersWithInteraction) / float64(registeredCustomers)) * 100
		coverageValues = append(coverageValues, &coverage)
	}

	values := map[string]*float64{
		"team_target_attainment":  normalizePercentage(summary.TeamTargetAttainment),
		"team_conversion_rate":    normalizePercentage(summary.TeamConversionRate),
		"territory_coverage":      normalizePercentage(averageFloatPointers(coverageValues)),
		"team_visit_compliance":   normalizePercentage(summary.TeamVisitCompliance),
		"team_overdue_task_rate":  normalizeInverseRate(summary.TeamOverdueTaskRate),
		"brick_pipeline_movement": normalizePipelineMovement(totalPipelineMovement, totalDealsCreated),
	}

	return computeWeightedScore(values, s.salesManagerWeights()), nil
}
