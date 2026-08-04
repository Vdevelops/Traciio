package kpi

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/kpi"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
)

func (s *Service) previousPeriod(start, end time.Time) (time.Time, time.Time) {
	loc := start.Location()
	durationDays := int(end.Sub(start).Hours()/24) + 1
	prevEndDate := start.AddDate(0, 0, -1)
	prevStartDate := prevEndDate.AddDate(0, 0, -durationDays+1)

	prevStart := time.Date(prevStartDate.Year(), prevStartDate.Month(), prevStartDate.Day(), 0, 0, 0, 0, loc)
	prevEnd := time.Date(prevEndDate.Year(), prevEndDate.Month(), prevEndDate.Day(), 23, 59, 59, 999999999, loc)
	return prevStart, prevEnd
}

func (s *Service) buildRawSalesRepScorecard(currentUser *user.User, start, end time.Time) (kpi.SalesRepScorecard, kpi.SalesRepMeta, error) {
	userID := currentUser.ID

	dealsCreated, err := s.kpiRepo.CountDealsCreated(userID, start, end)
	if err != nil {
		return kpi.SalesRepScorecard{}, kpi.SalesRepMeta{}, err
	}
	dealsClosed, err := s.kpiRepo.CountWonDeals(userID, start, end)
	if err != nil {
		return kpi.SalesRepScorecard{}, kpi.SalesRepMeta{}, err
	}
	totalRevenue, err := s.kpiRepo.SumWonRevenue(userID, start, end)
	if err != nil {
		return kpi.SalesRepScorecard{}, kpi.SalesRepMeta{}, err
	}
	visitsCompleted, err := s.kpiRepo.CountVisitCompleted(userID, start, end)
	if err != nil {
		return kpi.SalesRepScorecard{}, kpi.SalesRepMeta{}, err
	}
	visitsPlanned, err := s.kpiRepo.CountVisitPlanned(userID, start, end)
	if err != nil {
		return kpi.SalesRepScorecard{}, kpi.SalesRepMeta{}, err
	}
	tasksCreated, err := s.kpiRepo.CountTasksCreated(userID, start, end)
	if err != nil {
		return kpi.SalesRepScorecard{}, kpi.SalesRepMeta{}, err
	}
	tasksCompleted, err := s.kpiRepo.CountTasksCompleted(userID, start, end)
	if err != nil {
		return kpi.SalesRepScorecard{}, kpi.SalesRepMeta{}, err
	}
	overdueTasks, err := s.kpiRepo.CountOverdueTasks(userID, start, end)
	if err != nil {
		return kpi.SalesRepScorecard{}, kpi.SalesRepMeta{}, err
	}
	pipelineMovementScore, err := s.kpiRepo.GetPipelineMovementScore(userID, start, end)
	if err != nil {
		return kpi.SalesRepScorecard{}, kpi.SalesRepMeta{}, err
	}
	dealsWithoutBrick, err := s.kpiRepo.CountDealsWithoutBrick(userID, start, end)
	if err != nil {
		return kpi.SalesRepScorecard{}, kpi.SalesRepMeta{}, err
	}
	visitReportsWithoutBrick, err := s.kpiRepo.CountVisitReportsWithoutBrick(userID, start, end)
	if err != nil {
		return kpi.SalesRepScorecard{}, kpi.SalesRepMeta{}, err
	}

	conversionRate := pct(dealsClosed, dealsCreated)
	averageDealValue := (*float64)(nil)
	if dealsClosed > 0 {
		value := float64(totalRevenue) / float64(dealsClosed)
		averageDealValue = &value
	}
	visitCompliance := pct(visitsCompleted, visitsPlanned)
	overdueTaskRate := pct(overdueTasks, tasksCreated)

	revenueTarget, err := s.monthlyTargetRepo.GetProratedTargetForPeriod(userID, start.Format("2006-01-02"), end.Format("2006-01-02"))
	if err != nil {
		return kpi.SalesRepScorecard{}, kpi.SalesRepMeta{}, err
	}
	var revenueTargetAttainment *float64
	if revenueTarget > 0 {
		v := (float64(totalRevenue) / revenueTarget) * 100
		revenueTargetAttainment = &v
	}

	scorecard := kpi.SalesRepScorecard{
		TotalDealsClosed:        dealsClosed,
		TotalRevenue:            totalRevenue,
		DealsCreated:            dealsCreated,
		ConversionRate:          conversionRate,
		AverageDealValue:        averageDealValue,
		VisitCompleted:          visitsCompleted,
		VisitPlanned:            visitsPlanned,
		VisitCompliance:         visitCompliance,
		TasksCompleted:          tasksCompleted,
		OverdueTaskRate:         overdueTaskRate,
		RevenueTargetAttainment: revenueTargetAttainment,
		DealTargetAttainment:    nil,
		PipelineMovementScore:   pipelineMovementScore,
	}

	brickMissingCount := dealsWithoutBrick + visitReportsWithoutBrick
	brickInferredCount := int64(0)
	if currentUser.BrickID != nil && *currentUser.BrickID != "" {
		brickInferredCount = brickMissingCount
		brickMissingCount = 0
	}

	meta := kpi.SalesRepMeta{
		BrickMissingCount:  brickMissingCount,
		BrickInferredCount: brickInferredCount,
		GeneratedAt:        time.Now().In(response.GetTimezoneWIB()),
	}

	return scorecard, meta, nil
}

func (s *Service) peerAverageDealValue(currentUser *user.User, start, end time.Time) (*float64, error) {
	if s.brickRepo == nil || currentUser == nil || currentUser.BrickID == nil || *currentUser.BrickID == "" {
		return nil, nil
	}

	peers, err := s.brickRepo.GetSalesByBrickID(*currentUser.BrickID)
	if err != nil {
		return nil, err
	}

	var totalRevenue int64
	var totalWon int64
	for _, peer := range peers {
		if peer.ID == currentUser.ID {
			continue
		}
		wonCount, err := s.kpiRepo.CountWonDeals(peer.ID, start, end)
		if err != nil {
			return nil, err
		}
		if wonCount == 0 {
			continue
		}
		revenue, err := s.kpiRepo.SumWonRevenue(peer.ID, start, end)
		if err != nil {
			return nil, err
		}
		totalWon += wonCount
		totalRevenue += revenue
	}

	if totalWon == 0 {
		return nil, nil
	}

	avg := float64(totalRevenue) / float64(totalWon)
	return &avg, nil
}

func (s *Service) revenueTargetForPeriod(userID string, start, end time.Time) (int64, error) {
	target, err := s.monthlyTargetRepo.GetProratedTargetForPeriod(userID, start.Format("2006-01-02"), end.Format("2006-01-02"))
	if err != nil {
		return 0, err
	}
	return int64(target + 0.5), nil
}

func (s *Service) loadCurrentUser(userID string) (*user.User, error) {
	currentUser, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	return currentUser, nil
}

func (s *Service) validateCurrentUserRole(userID string) (*user.User, string, error) {
	currentUser, err := s.loadCurrentUser(userID)
	if err != nil {
		return nil, "", err
	}
	roleCode := ""
	if currentUser.Role != nil {
		roleCode = currentUser.Role.Code
	}
	return currentUser, roleCode, nil
}

func previousPeriodRange(start, end time.Time) (time.Time, time.Time) {
	loc := start.Location()
	durationDays := int(end.Sub(start).Hours()/24) + 1
	prevEndDate := start.AddDate(0, 0, -1)
	prevStartDate := prevEndDate.AddDate(0, 0, -durationDays+1)

	prevStart := time.Date(prevStartDate.Year(), prevStartDate.Month(), prevStartDate.Day(), 0, 0, 0, 0, loc)
	prevEnd := time.Date(prevEndDate.Year(), prevEndDate.Month(), prevEndDate.Day(), 23, 59, 59, 999999999, loc)
	return prevStart, prevEnd
}
