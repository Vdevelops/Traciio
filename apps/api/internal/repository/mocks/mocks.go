package mocks

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/brick"
	"github.com/gilabs/crm-healthcare/api/internal/domain/group"
	"github.com/gilabs/crm-healthcare/api/internal/domain/monthly_target"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"

	"github.com/gilabs/crm-healthcare/api/internal/domain/account"
	"github.com/gilabs/crm-healthcare/api/internal/domain/activity"
	"github.com/gilabs/crm-healthcare/api/internal/domain/dashboard"
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead"
	"github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
	"github.com/gilabs/crm-healthcare/api/internal/domain/sales_overview"
	"github.com/gilabs/crm-healthcare/api/internal/domain/task"
	"github.com/gilabs/crm-healthcare/api/internal/domain/visit_report"
	"github.com/stretchr/testify/mock"
)

// DealRepository mock
type DealRepository struct {
	mock.Mock
}

func (m *DealRepository) FindByID(id string) (*pipeline.Deal, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pipeline.Deal), args.Error(1)
}

func (m *DealRepository) List(req *pipeline.ListDealsRequest) ([]pipeline.Deal, int64, error) {
	args := m.Called(req)
	return args.Get(0).([]pipeline.Deal), args.Get(1).(int64), args.Error(2)
}

func (m *DealRepository) Create(deal *pipeline.Deal) error {
	args := m.Called(deal)
	return args.Error(0)
}

func (m *DealRepository) Update(deal *pipeline.Deal) error {
	args := m.Called(deal)
	return args.Error(0)
}

func (m *DealRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *DealRepository) GetSummary() (*pipeline.PipelineSummaryResponse, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pipeline.PipelineSummaryResponse), args.Error(1)
}

func (m *DealRepository) GetForecast(periodType string, start, end time.Time) (*pipeline.ForecastResponse, error) {
	args := m.Called(periodType, start, end)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pipeline.ForecastResponse), args.Error(1)
}

func (m *DealRepository) GetStatsByStatus(startDate, endDate string, assignedTo, stageID, status string) (map[string]int64, error) {
	args := m.Called(startDate, endDate, assignedTo, stageID, status)
	return args.Get(0).(map[string]int64), args.Error(1)
}

func (m *DealRepository) GetStatsByStage(startDate, endDate string, assignedTo, status string) (map[string]int64, error) {
	args := m.Called(startDate, endDate, assignedTo, status)
	return args.Get(0).(map[string]int64), args.Error(1)
}

func (m *DealRepository) CountByDateRange(startDate, endDate interface{}) (int64, error) {
	args := m.Called(startDate, endDate)
	return args.Get(0).(int64), args.Error(1)
}

func (m *DealRepository) GetWonDealsValueInPeriod(startDate, endDate interface{}) (int64, int64, error) {
	args := m.Called(startDate, endDate)
	return args.Get(0).(int64), args.Get(1).(int64), args.Error(2)
}

func (m *DealRepository) GetWonDealsValueInPeriodByUser(userID string, startDate, endDate interface{}) (int64, int64, error) {
	args := m.Called(userID, startDate, endDate)
	return args.Get(0).(int64), args.Get(1).(int64), args.Error(2)
}

func (m *DealRepository) GetSummaryInPeriod(startDate, endDate interface{}) (*pipeline.PipelineSummaryResponse, error) {
	args := m.Called(startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pipeline.PipelineSummaryResponse), args.Error(1)
}

// VisitReportRepository mock
type VisitReportRepository struct {
	mock.Mock
}

func (m *VisitReportRepository) FindByID(id string) (*visit_report.VisitReport, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*visit_report.VisitReport), args.Error(1)
}

func (m *VisitReportRepository) List(req *visit_report.ListVisitReportsRequest) ([]visit_report.VisitReport, int64, error) {
	args := m.Called(req)
	return args.Get(0).([]visit_report.VisitReport), args.Get(1).(int64), args.Error(2)
}

func (m *VisitReportRepository) Create(report *visit_report.VisitReport) error {
	args := m.Called(report)
	return args.Error(0)
}

func (m *VisitReportRepository) Update(report *visit_report.VisitReport) error {
	args := m.Called(report)
	return args.Error(0)
}

func (m *VisitReportRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *VisitReportRepository) FindByAccountID(accountID string) ([]visit_report.VisitReport, error) {
    args := m.Called(accountID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]visit_report.VisitReport), args.Error(1)
}

func (m *VisitReportRepository) FindBySalesRepID(salesRepID string) ([]visit_report.VisitReport, error) {
    args := m.Called(salesRepID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]visit_report.VisitReport), args.Error(1)
}

func (m *VisitReportRepository) GetStatsByStatus(startDate, endDate string, accountID, salesRepID, status string) (map[string]int64, error) {
	args := m.Called(startDate, endDate, accountID, salesRepID, status)
	return args.Get(0).(map[string]int64), args.Error(1)
}

func (m *VisitReportRepository) GetStatsByDate(startDate, endDate string, accountID, salesRepID, status string) (map[string]int64, error) {
    args := m.Called(startDate, endDate, accountID, salesRepID, status)
    return args.Get(0).(map[string]int64), args.Error(1)
}

func (m *VisitReportRepository) GetStatsByDateAndStatus(startDate, endDate string, accountID, salesRepID string) (map[string]map[string]int64, error) {
    args := m.Called(startDate, endDate, accountID, salesRepID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(map[string]map[string]int64), args.Error(1)
}

func (m *VisitReportRepository) GetStatsByAccount(startDate, endDate string, salesRepID, status string) (map[string]int64, error) {
    args := m.Called(startDate, endDate, salesRepID, status)
    return args.Get(0).(map[string]int64), args.Error(1)
}

func (m *VisitReportRepository) GetStatsBySalesRep(startDate, endDate string, accountID, status string) (map[string]int64, error) {
    args := m.Called(startDate, endDate, accountID, status)
    return args.Get(0).(map[string]int64), args.Error(1)
}

func (m *VisitReportRepository) GetStatsBySalesRepWithAccounts(startDate, endDate string, status string) (map[string]struct {
    VisitCount   int64
    AccountCount int64
}, error) {
    args := m.Called(startDate, endDate, status)
    return args.Get(0).(map[string]struct {
        VisitCount   int64
        AccountCount int64
    }), args.Error(1)
}

// AccountRepository mock
type AccountRepository struct {
	mock.Mock
}

func (m *AccountRepository) FindByID(id string) (*account.Account, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*account.Account), args.Error(1)
}

func (m *AccountRepository) List(req *account.ListAccountsRequest) ([]account.Account, int64, error) {
	args := m.Called(req)
	return args.Get(0).([]account.Account), args.Get(1).(int64), args.Error(2)
}

func (m *AccountRepository) ListAll(status string) ([]account.Account, error) {
    args := m.Called(status)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]account.Account), args.Error(1)
}

func (m *AccountRepository) Create(acc *account.Account) error {
	args := m.Called(acc)
	return args.Error(0)
}

func (m *AccountRepository) Update(acc *account.Account) error {
	args := m.Called(acc)
	return args.Error(0)
}

func (m *AccountRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *AccountRepository) GetStatsByStatus() (map[string]int64, error) {
	args := m.Called()
	return args.Get(0).(map[string]int64), args.Error(1)
}

func (m *AccountRepository) CountTotal() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *AccountRepository) GetStatsByBrick(brickID string) (map[string]int64, error) {
	args := m.Called(brickID)
	return args.Get(0).(map[string]int64), args.Error(1)
}

func (m *AccountRepository) CountByDateRange(startDate, endDate interface{}) (int64, error) {
    args := m.Called(startDate, endDate)
    return args.Get(0).(int64), args.Error(1)
}

// ActivityRepository mock
type ActivityRepository struct {
	mock.Mock
}

func (m *ActivityRepository) List(req *activity.ListActivitiesRequest) ([]activity.Activity, int64, error) {
	args := m.Called(req)
	return args.Get(0).([]activity.Activity), args.Get(1).(int64), args.Error(2)
}

func (m *ActivityRepository) Create(act *activity.Activity) error {
	args := m.Called(act)
	return args.Error(0)
}

func (m *ActivityRepository) Update(act *activity.Activity) error {
    args := m.Called(act)
    return args.Error(0)
}

func (m *ActivityRepository) GetStatsByType(startDate, endDate, accountID string) (map[string]int64, error) {
	args := m.Called(startDate, endDate, accountID)
	return args.Get(0).(map[string]int64), args.Error(1)
}

func (m *ActivityRepository) FindByID(id string) (*activity.Activity, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*activity.Activity), args.Error(1)
}

func (m *ActivityRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *ActivityRepository) GetTimeline(req *activity.ActivityTimelineRequest) ([]activity.Activity, error) {
    args := m.Called(req)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]activity.Activity), args.Error(1)
}

func (m *ActivityRepository) FindByAccountID(accountID string) ([]activity.Activity, error) {
    args := m.Called(accountID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]activity.Activity), args.Error(1)
}

func (m *ActivityRepository) GetStatsByTypeAndDate(startDate, endDate string, accountID string) (map[string]map[string]int64, error) {
    args := m.Called(startDate, endDate, accountID)
    return args.Get(0).(map[string]map[string]int64), args.Error(1)
}

func (m *ActivityRepository) GetStatsByUser(startDate, endDate string, accountID string) (map[string]int64, error) {
    args := m.Called(startDate, endDate, accountID)
    return args.Get(0).(map[string]int64), args.Error(1)
}

// LeadRepository mock
type LeadRepository struct {
	mock.Mock
}

func (m *LeadRepository) FindByID(id string) (*lead.Lead, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*lead.Lead), args.Error(1)
}

func (m *LeadRepository) List(req *lead.ListLeadsRequest) ([]lead.Lead, int64, error) {
	args := m.Called(req)
	return args.Get(0).([]lead.Lead), args.Get(1).(int64), args.Error(2)
}

func (m *LeadRepository) Create(l *lead.Lead) error {
	args := m.Called(l)
	return args.Error(0)
}

func (m *LeadRepository) Update(l *lead.Lead) error {
	args := m.Called(l)
	return args.Error(0)
}

func (m *LeadRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *LeadRepository) GetAnalytics(req *lead.LeadAnalyticsRequest) (*lead.LeadAnalyticsResponse, error) {
    args := m.Called(req)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*lead.LeadAnalyticsResponse), args.Error(1)
}

func (m *LeadRepository) GetStatsByStatus() (map[string]int64, error) {
	args := m.Called()
	return args.Get(0).(map[string]int64), args.Error(1)
}

func (m *LeadRepository) GetStatsBySource() (map[string]int64, error) {
	args := m.Called()
	return args.Get(0).(map[string]int64), args.Error(1)
}

func (m *LeadRepository) CountByDateRange(startDate, endDate interface{}) (int64, error) {
    args := m.Called(startDate, endDate)
    return args.Get(0).(int64), args.Error(1)
}

func (m *LeadRepository) GetStatsByStatusAndDateRange(startDate, endDate interface{}) (map[string]int64, error) {
    args := m.Called(startDate, endDate)
    return args.Get(0).(map[string]int64), args.Error(1)
}

func (m *LeadRepository) GetStatsBySourceAndDateRange(startDate, endDate interface{}) (map[string]int64, error) {
    args := m.Called(startDate, endDate)
    return args.Get(0).(map[string]int64), args.Error(1)
}

// TaskRepository mock
type TaskRepository struct {
	mock.Mock
}

func (m *TaskRepository) FindByID(id string) (*task.Task, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*task.Task), args.Error(1)
}

func (m *TaskRepository) List(req *task.ListTasksRequest) ([]task.Task, int64, error) {
    args := m.Called(req)
    if args.Get(0) == nil {
        return nil, args.Get(1).(int64), args.Error(2)
    }
    return args.Get(0).([]task.Task), args.Get(1).(int64), args.Error(2)
}

func (m *TaskRepository) Create(t *task.Task) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *TaskRepository) Update(t *task.Task) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *TaskRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *TaskRepository) CountByStatus(status string, userID string) (int64, error) {
    args := m.Called(status, userID)
    return args.Get(0).(int64), args.Error(1)
}

// MonthlyTargetRepository mock (Minimal for dashboard)
type MonthlyTargetRepository struct {
	mock.Mock
}

func (m *MonthlyTargetRepository) GetTarget(userID, period string) (*dashboard.TargetStats, error) {
    args := m.Called(userID, period)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dashboard.TargetStats), args.Error(1)
}

// UserRepository mock (Minimal)
type UserRepository struct {
    mock.Mock
}
// Implement minimal methods if used

func (m *UserRepository) FindByID(id string) (*user.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *UserRepository) CountUsersByRoleID(roleID string) (int64, error) {
	args := m.Called(roleID)
	return args.Get(0).(int64), args.Error(1)
}


func (m *UserRepository) FindByEmail(email string) (*user.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *UserRepository) List(req *user.ListUsersRequest) ([]user.User, int64, error) {
	args := m.Called(req)
	return args.Get(0).([]user.User), args.Get(1).(int64), args.Error(2)
}

func (m *UserRepository) Create(u *user.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *UserRepository) Update(u *user.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *UserRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *UserRepository) GetUsersByGroupID(groupID string) ([]user.User, error) {
	args := m.Called(groupID)
	return args.Get(0).([]user.User), args.Error(1)
}

func (m *UserRepository) GetUsersByBrickID(brickID string) ([]user.User, error) {
	args := m.Called(brickID)
	return args.Get(0).([]user.User), args.Error(1)
}

func (m *UserRepository) GetUsersByRoleID(roleID string) ([]string, error) {
	args := m.Called(roleID)
	return args.Get(0).([]string), args.Error(1)
}

// RoleRepository mock (Minimal)
type RoleRepository struct {
    mock.Mock
}

// PipelineRepository mock (Minimal)
type PipelineRepository struct {
    mock.Mock
}

// SalesOverviewRepository mock
type SalesOverviewRepository struct {
    mock.Mock
}

func (m *SalesOverviewRepository) GetSalesPerformanceDetail(userID string, startDate, endDate interface{}) (*sales_overview.SalesPerformanceDetail, error) {
    args := m.Called(userID, startDate, endDate)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*sales_overview.SalesPerformanceDetail), args.Error(1)
}

func (m *SalesOverviewRepository) GetSalesRepDetail(userID string, startDate, endDate interface{}) (*sales_overview.SalesRepDetail, error) {
    args := m.Called(userID, startDate, endDate)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*sales_overview.SalesRepDetail), args.Error(1)
}

func (m *SalesOverviewRepository) ListSalesPerformance(req *sales_overview.ListSalesPerformanceRequest) ([]sales_overview.SalesPerformanceListResponse, int64, error) {
    args := m.Called(req)
    if args.Get(0) == nil {
        return nil, args.Get(1).(int64), args.Error(2)
    }
    return args.Get(0).([]sales_overview.SalesPerformanceListResponse), args.Get(1).(int64), args.Error(2)
}

func (m *SalesOverviewRepository) GetSalesRepCheckInLocations(userID string, req *sales_overview.GetSalesRepCheckInLocationsRequest, startDate, endDate interface{}) ([]sales_overview.CheckInLocation, int64, error) {
    args := m.Called(userID, req, startDate, endDate)
    if args.Get(0) == nil {
        return nil, args.Get(1).(int64), args.Error(2)
    }
	return args.Get(0).([]sales_overview.CheckInLocation), args.Get(1).(int64), args.Error(2)
}

func (m *SalesOverviewRepository) GetMonthlySalesOverview(startDate, endDate interface{}) (*sales_overview.MonthlySalesOverviewResponse, error) {
    args := m.Called(startDate, endDate)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*sales_overview.MonthlySalesOverviewResponse), args.Error(1)
}

func (m *MonthlyTargetRepository) FindByID(id string) (*monthly_target.MonthlyTarget, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*monthly_target.MonthlyTarget), args.Error(1)
}

func (m *MonthlyTargetRepository) FindByGroupAndPeriod(groupID string, year int, month int) (*monthly_target.MonthlyTarget, error) {
	args := m.Called(groupID, year, month)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*monthly_target.MonthlyTarget), args.Error(1)
}

func (m *MonthlyTargetRepository) FindByUserAndPeriod(userID string, year int, month int) (*monthly_target.MonthlyTarget, error) {
	args := m.Called(userID, year, month)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*monthly_target.MonthlyTarget), args.Error(1)
}

func (m *MonthlyTargetRepository) FindByBrickAndPeriod(brickID string, year int, month int) (*monthly_target.MonthlyTarget, error) {
	args := m.Called(brickID, year, month)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*monthly_target.MonthlyTarget), args.Error(1)
}

func (m *MonthlyTargetRepository) List(req *monthly_target.ListMonthlyTargetsRequest) ([]monthly_target.MonthlyTarget, int64, int64, error) {
	args := m.Called(req)
	return args.Get(0).([]monthly_target.MonthlyTarget), args.Get(1).(int64), args.Get(2).(int64), args.Error(3)
}

func (m *MonthlyTargetRepository) Create(target *monthly_target.MonthlyTarget) error {
	args := m.Called(target)
	return args.Error(0)
}

func (m *MonthlyTargetRepository) Update(target *monthly_target.MonthlyTarget) error {
	args := m.Called(target)
	return args.Error(0)
}

func (m *MonthlyTargetRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MonthlyTargetRepository) GetUserEffectiveTarget(userID string, year int, month int) (*monthly_target.MonthlyTarget, error) {
	args := m.Called(userID, year, month)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*monthly_target.MonthlyTarget), args.Error(1)
}

func (m *MonthlyTargetRepository) BatchGetUserEffectiveTargets(userIDs []string, year int, month int) (map[string]*monthly_target.MonthlyTarget, error) {
	args := m.Called(userIDs, year, month)
	return args.Get(0).(map[string]*monthly_target.MonthlyTarget), args.Error(1)
}

func (m *MonthlyTargetRepository) GetTotalEffectiveTarget(year int, month int) (int64, error) {
	args := m.Called(year, month)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MonthlyTargetRepository) GetProratedTargetForPeriod(userID string, startDate, endDate string) (float64, error) {
	args := m.Called(userID, startDate, endDate)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MonthlyTargetRepository) BatchGetProratedTargetsForPeriod(userIDs []string, startDate, endDate string) (map[string]float64, error) {
	args := m.Called(userIDs, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]float64), args.Error(1)
}

// GroupRepository mock
type GroupRepository struct {
	mock.Mock
}

func (m *GroupRepository) FindByID(id string) (*group.Group, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*group.Group), args.Error(1)
}

func (m *GroupRepository) List(req *group.ListGroupsRequest) ([]group.Group, int64, error) {
	args := m.Called(req)
	return args.Get(0).([]group.Group), args.Get(1).(int64), args.Error(2)
}

// BrickRepository mock
type BrickRepository struct {
	mock.Mock
}

func (m *BrickRepository) FindByID(id string) (*brick.Brick, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*brick.Brick), args.Error(1)
}

func (m *BrickRepository) List(req *brick.ListBricksRequest) ([]brick.Brick, int64, error) {
	args := m.Called(req)
	return args.Get(0).([]brick.Brick), args.Get(1).(int64), args.Error(2)
}

func (m *GroupRepository) CountUsersByGroupID(groupID string) (int64, error) {
	args := m.Called(groupID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *GroupRepository) FindByCode(code string) (*group.Group, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*group.Group), args.Error(1)
}

func (m *GroupRepository) Create(g *group.Group) error {
	args := m.Called(g)
	return args.Error(0)
}

func (m *GroupRepository) Update(g *group.Group) error {
	args := m.Called(g)
	return args.Error(0)
}

func (m *GroupRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *BrickRepository) CountSalesByBrickID(brickID string) (int64, error) {
	args := m.Called(brickID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *BrickRepository) FindByIDs(ids []string) ([]brick.Brick, error) {
	args := m.Called(ids)
	return args.Get(0).([]brick.Brick), args.Error(1)
}

func (m *BrickRepository) FindByCode(code string) (*brick.Brick, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*brick.Brick), args.Error(1)
}

func (m *BrickRepository) Create(b *brick.Brick) error {
	args := m.Called(b)
	return args.Error(0)
}

func (m *BrickRepository) Update(b *brick.Brick) error {
	args := m.Called(b)
	return args.Error(0)
}

func (m *BrickRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *BrickRepository) GetSalesByBrickID(brickID string) ([]user.User, error) {
	args := m.Called(brickID)
	return args.Get(0).([]user.User), args.Error(1)
}

func (m *BrickRepository) FindByRegencyAndProvince(regency, province string) (*brick.Brick, error) {
	args := m.Called(regency, province)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*brick.Brick), args.Error(1)
}

func (m *BrickRepository) GetNextCodeSequence(prefix string) (int, error) {
	args := m.Called(prefix)
	return args.Int(0), args.Error(1)
}
