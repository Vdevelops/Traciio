package dashboard

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/account"
	"github.com/gilabs/crm-healthcare/api/internal/domain/activity"
	brickdomain "github.com/gilabs/crm-healthcare/api/internal/domain/brick"
	"github.com/gilabs/crm-healthcare/api/internal/domain/dashboard"
	leaddomain "github.com/gilabs/crm-healthcare/api/internal/domain/lead"
	pipelinedomain "github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
	scheduledomain "github.com/gilabs/crm-healthcare/api/internal/domain/schedule"
	taskdomain "github.com/gilabs/crm-healthcare/api/internal/domain/task"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"github.com/gilabs/crm-healthcare/api/internal/domain/visit_report"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"github.com/gilabs/crm-healthcare/api/pkg/cache"
	"github.com/gilabs/crm-healthcare/api/pkg/util/currency"
	"gorm.io/gorm"
)

type Service struct {
	visitReportRepo   interfaces.VisitReportRepository
	accountRepo       interfaces.AccountRepository
	activityRepo      interfaces.ActivityRepository
	userRepo          interfaces.UserRepository
	dealRepo          interfaces.DealRepository
	taskRepo          interfaces.TaskRepository
	pipelineRepo      interfaces.PipelineRepository
	leadRepo          interfaces.LeadRepository
	roleRepo          interfaces.RoleRepository
	monthlyTargetRepo interfaces.MonthlyTargetRepository
	brickRepo         interfaces.BrickRepository
	scheduleRepo      interfaces.ScheduleRepository
	cacheService      *cache.DashboardCacheService
}

const dashboardListPageSize = 500

func NewService(
	visitReportRepo interfaces.VisitReportRepository,
	accountRepo interfaces.AccountRepository,
	activityRepo interfaces.ActivityRepository,
	userRepo interfaces.UserRepository,
	dealRepo interfaces.DealRepository,
	taskRepo interfaces.TaskRepository,
	pipelineRepo interfaces.PipelineRepository,
	leadRepo interfaces.LeadRepository,
	roleRepo interfaces.RoleRepository,
	monthlyTargetRepo interfaces.MonthlyTargetRepository,
	brickRepo interfaces.BrickRepository,
	scheduleRepo interfaces.ScheduleRepository,
) *Service {
	return &Service{
		visitReportRepo:   visitReportRepo,
		accountRepo:       accountRepo,
		activityRepo:      activityRepo,
		userRepo:          userRepo,
		dealRepo:          dealRepo,
		taskRepo:          taskRepo,
		pipelineRepo:      pipelineRepo,
		leadRepo:          leadRepo,
		roleRepo:          roleRepo,
		monthlyTargetRepo: monthlyTargetRepo,
		brickRepo:         brickRepo,
		scheduleRepo:      scheduleRepo,
		cacheService:      cache.NewDashboardCacheService(nil),
	}
}

// parsePeriod parses period string and returns start and end dates
func parsePeriod(period string) (time.Time, time.Time) {
	now := time.Now()
	var start, end time.Time

	switch period {
	case "today":
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
	case "week":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		start = now.AddDate(0, 0, -weekday+1)
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
		end = now
	case "month":
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		// End of current month (last day, 23:59:59.999999999)
		end = time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 999999999, now.Location())
	case "year":
		start = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		end = now
	default:
		// Default to today
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
	}

	return start, end
}

// resolveTeamUserIDs returns a deduplicated list of user IDs representing all
// sales reps whose brick is managed by the given managerID. If no team members
// are found the slice falls back to containing only the manager's own ID, so
// that the manager's own data is always included.
func (s *Service) resolveTeamUserIDs(managerID string) []string {
	seen := make(map[string]struct{})
	seen[managerID] = struct{}{} // always include the manager

	if s.brickRepo == nil {
		return []string{managerID}
	}

	// Find all bricks where manager_id = managerID
	bricks, _, err := s.brickRepo.List(&brickdomain.ListBricksRequest{
		ManagerID: &managerID,
		Page:      1,
		PerPage:   100,
	})
	if err == nil {
		for _, b := range bricks {
			salesReps, salesErr := s.brickRepo.GetSalesByBrickID(b.ID)
			if salesErr == nil {
				for _, rep := range salesReps {
					seen[rep.ID] = struct{}{}
				}
			}
		}
	}

	ids := make([]string, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	return ids
}

// isScopedUser returns true when the given userID belongs to the scoped set.
// If scopedUserIDs is nil (global access), it always returns true.
func isScopedUser(scopedUserIDs []string, userID string) bool {
	if scopedUserIDs == nil {
		return true
	}
	for _, id := range scopedUserIDs {
		if id == userID {
			return true
		}
	}
	return false
}

// scopedUserIDsSet builds a lookup set from ScopedUserIDs for O(1) checks.
// Returns nil if scopedUserIDs is nil (global access).
func scopedUserIDsSet(scopedUserIDs []string) map[string]struct{} {
	if scopedUserIDs == nil {
		return nil
	}
	set := make(map[string]struct{}, len(scopedUserIDs))
	for _, id := range scopedUserIDs {
		set[id] = struct{}{}
	}
	return set
}

func (s *Service) loadSalesUsers(scopedUserIDs []string) (map[string]user.User, []string, error) {
	salesUsers := make(map[string]user.User)
	salesUserIDs := make([]string, 0)

	for page := 1; ; page++ {
		users, total, err := s.userRepo.List(&user.ListUsersRequest{
			Page:          page,
			PerPage:       dashboardListPageSize,
			ScopedUserIDs: scopedUserIDs,
		})
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, nil, err
		}

		for _, userItem := range users {
			if userItem.Role == nil || userItem.Role.Code != "sales" {
				continue
			}
			if _, exists := salesUsers[userItem.ID]; exists {
				continue
			}
			salesUsers[userItem.ID] = userItem
			salesUserIDs = append(salesUserIDs, userItem.ID)
		}

		if int64(page*dashboardListPageSize) >= total || len(users) == 0 {
			break
		}
	}

	return salesUsers, salesUserIDs, nil
}

func (s *Service) listScopedAccounts(scopedUserIDs []string) ([]account.Account, error) {
	results := make([]account.Account, 0)

	for page := 1; ; page++ {
		items, total, err := s.accountRepo.List(&account.ListAccountsRequest{
			Page:          page,
			PerPage:       dashboardListPageSize,
			ScopedUserIDs: scopedUserIDs,
		})
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}

		results = append(results, items...)
		if int64(page*dashboardListPageSize) >= total || len(items) == 0 {
			break
		}
	}

	return results, nil
}

func (s *Service) listScopedActivities(req *activity.ListActivitiesRequest) ([]activity.Activity, error) {
	results := make([]activity.Activity, 0)

	for page := 1; ; page++ {
		pageReq := *req
		pageReq.Page = page
		pageReq.PerPage = dashboardListPageSize

		items, total, err := s.activityRepo.List(&pageReq)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}

		results = append(results, items...)
		if int64(page*dashboardListPageSize) >= total || len(items) == 0 {
			break
		}
	}

	return results, nil
}

func (s *Service) listScopedLeads(req *leaddomain.ListLeadsRequest) ([]leaddomain.Lead, error) {
	results := make([]leaddomain.Lead, 0)

	for page := 1; ; page++ {
		pageReq := *req
		pageReq.Page = page
		pageReq.PerPage = dashboardListPageSize

		items, total, err := s.leadRepo.List(&pageReq)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}

		results = append(results, items...)
		if int64(page*dashboardListPageSize) >= total || len(items) == 0 {
			break
		}
	}

	return results, nil
}

func (s *Service) listScopedDeals(req *pipelinedomain.ListDealsRequest) ([]pipelinedomain.Deal, error) {
	results := make([]pipelinedomain.Deal, 0)

	for page := 1; ; page++ {
		pageReq := *req
		pageReq.Page = page
		pageReq.PerPage = dashboardListPageSize

		items, total, err := s.dealRepo.List(&pageReq)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}

		results = append(results, items...)
		if int64(page*dashboardListPageSize) >= total || len(items) == 0 {
			break
		}
	}

	return results, nil
}

func (s *Service) listScopedVisitReports(req *visit_report.ListVisitReportsRequest) ([]visit_report.VisitReport, error) {
	results := make([]visit_report.VisitReport, 0)

	for page := 1; ; page++ {
		pageReq := *req
		pageReq.Page = page
		pageReq.PerPage = dashboardListPageSize

		items, total, err := s.visitReportRepo.List(&pageReq)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}

		results = append(results, items...)
		if int64(page*dashboardListPageSize) >= total || len(items) == 0 {
			break
		}
	}

	return results, nil
}

// GetOverview returns dashboard overview
func (s *Service) GetOverview(req *dashboard.DashboardRequest, userID string) (*dashboard.DashboardOverviewResponse, error) {
	// Parse period
	var start, end time.Time
	if req.StartDate != "" && req.EndDate != "" {
		var err error
		start, err = time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, err
		}
		end, err = time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, err
		}
		end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, end.Location())
	} else if req.Period != "" {
		start, end = parsePeriod(req.Period)
	} else {
		start, end = parsePeriod("today")
	}

	// Try to get from cache first
	var cachedResponse dashboard.DashboardOverviewResponse
	if found, _ := s.cacheService.GetOverview(userID, req.Period, req.StartDate, req.EndDate, &cachedResponse); found {
		return &cachedResponse, nil
	}

	salesUsers, salesUserIDs, err := s.loadSalesUsers(req.ScopedUserIDs)
	if err != nil {
		return nil, err
	}

	// OPTIMIZED: Use database aggregation instead of loading all records into memory
	// Get visit report stats by status using aggregation
	// FIXED: Use batch query for scoped users to prevent N+1 (was looping per user)
	visitStatsByStatus := make(map[string]int64)
	if req.ScopedUserIDs != nil && len(req.ScopedUserIDs) > 0 {
		// Use batch query for multiple users - fixes N+1 bottleneck
		var vsErr error
		visitStatsByStatus, vsErr = s.visitReportRepo.GetStatsByStatusForUsers(
			start.Format("2006-01-02"),
			end.Format("2006-01-02"),
			req.ScopedUserIDs,
		)
		if vsErr != nil && vsErr != gorm.ErrRecordNotFound {
			return nil, vsErr
		}
	} else {
		var vsErr error
		visitStatsByStatus, vsErr = s.visitReportRepo.GetStatsByStatus(
			start.Format("2006-01-02"),
			end.Format("2006-01-02"),
			"", "", "",
		)
		if vsErr != nil && vsErr != gorm.ErrRecordNotFound {
			return nil, vsErr
		}
	}

	// Calculate visit stats from aggregation
	visitStats := struct {
		Total         int
		Completed     int
		Pending       int
		Approved      int
		Rejected      int
		ChangePercent float64
	}{}

	for status, count := range visitStatsByStatus {
		visitStats.Total += int(count)
		switch status {
		case "submitted", "approved":
			visitStats.Completed += int(count)
			if status == "approved" {
				visitStats.Approved += int(count)
			}
		case "draft":
			visitStats.Pending += int(count)
		case "rejected":
			visitStats.Rejected += int(count)
		}
	}

	// OPTIMIZED: Use database aggregation for account stats
	accountStats := struct {
		Total         int
		Active        int
		Inactive      int
		ChangePercent float64
	}{}
	if req.ScopedUserIDs != nil {
		accounts, accErr := s.listScopedAccounts(salesUserIDs)
		if accErr != nil {
			return nil, accErr
		}
		for _, acc := range accounts {
			accountStats.Total++
			if acc.Status == "active" {
				accountStats.Active++
			} else {
				accountStats.Inactive++
			}
		}
	} else {
		accountStatsByStatus, err := s.accountRepo.GetStatsByStatus()
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}

		for status, count := range accountStatsByStatus {
			accountStats.Total += int(count)
			if status == "active" {
				accountStats.Active += int(count)
			} else {
				accountStats.Inactive += int(count)
			}
		}
	}

	activityStats := struct {
		Total         int
		Visits        int
		Calls         int
		Emails        int
		ChangePercent float64
	}{}
	if req.ScopedUserIDs != nil {
		activities, activityErr := s.listScopedActivities(&activity.ListActivitiesRequest{
			StartDate:     start.Format("2006-01-02"),
			EndDate:       end.Format("2006-01-02"),
			ScopedUserIDs: salesUserIDs,
		})
		if activityErr != nil {
			return nil, activityErr
		}
		for _, activityItem := range activities {
			if _, ok := salesUsers[activityItem.UserID]; !ok {
				continue
			}
			activityStats.Total++
			switch activityItem.Type {
			case "visit":
				activityStats.Visits++
			case "call":
				activityStats.Calls++
			case "email":
				activityStats.Emails++
			}
		}
	} else {
		activityStatsByType, err := s.activityRepo.GetStatsByType(
			start.Format("2006-01-02"),
			end.Format("2006-01-02"),
			"",
		)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}

		for activityType, count := range activityStatsByType {
			activityStats.Total += int(count)
			switch activityType {
			case "visit":
				activityStats.Visits += int(count)
			case "call":
				activityStats.Calls += int(count)
			case "email":
				activityStats.Emails += int(count)
			}
		}
	}

	// Get pipeline/deals summary
	var dealsSummary *pipelinedomain.PipelineSummaryResponse
	if s.dealRepo != nil {
		if req.ScopedUserIDs != nil {
			// Scoped: aggregate deal summary per user via deals list
			scopedDeals, _, listErr := s.dealRepo.List(&pipelinedomain.ListDealsRequest{
				Page:          1,
				PerPage:       1000,
				ScopedUserIDs: req.ScopedUserIDs,
			})
			if listErr == nil {
				summary := &pipelinedomain.PipelineSummaryResponse{}
				for _, d := range scopedDeals {
					summary.TotalDeals++
					summary.TotalValue += d.Value
					switch d.Status {
					case "open":
						summary.OpenDeals++
					case "won":
						summary.WonDeals++
						summary.WonValue += d.Value
					case "lost":
						summary.LostDeals++
					}
				}
				summary.TotalValueFormatted = formatCurrency(summary.TotalValue)
				summary.WonValueFormatted = formatCurrency(summary.WonValue)
				dealsSummary = summary
			}
		} else {
			dealsSummary, err = s.dealRepo.GetSummaryInPeriod(start, end)
			if err != nil && err != gorm.ErrRecordNotFound {
				return nil, err
			}
		}
	}

	// OPTIMIZED: Use database aggregation for lead statistics
	leadStats := dashboard.LeadStats{}
	if s.leadRepo != nil {
		if req.ScopedUserIDs != nil {
			// Scoped: use List with ScopedUserIDs and aggregate by status
			leads, leadErr := s.listScopedLeads(&leaddomain.ListLeadsRequest{
				ScopedUserIDs: salesUserIDs,
			})
			if leadErr == nil {
				for _, l := range leads {
					leadStats.Total++
					switch l.LeadStatus {
					case "new":
						leadStats.New++
					case "contacted", "interested":
						leadStats.Contacted++
					case "qualified", "proposal_sent":
						leadStats.Qualified++
					case "converted":
						leadStats.Converted++
					case "lost":
						leadStats.Lost++
					}
				}
			}
		} else {
			leadStatsByStatus, lsErr := s.leadRepo.GetStatsByStatusAndDateRange(start, end)
			if lsErr != nil && lsErr != gorm.ErrRecordNotFound {
				return nil, lsErr
			}
			for status, count := range leadStatsByStatus {
				leadStats.Total += count
				switch status {
				case "new":
					leadStats.New += count
				case "contacted", "interested":
					leadStats.Contacted += count
				case "qualified", "proposal_sent":
					leadStats.Qualified += count
				case "converted":
					leadStats.Converted += count
				case "lost":
					leadStats.Lost += count
				}
			}
		}
	}

	// OPTIMIZED: Use database aggregation for leads by source
	leadsBySource := dashboard.LeadsBySource{}
	if s.leadRepo != nil {
		entries := make([]dashboard.LeadsBySourceEntry, 0)
		if req.ScopedUserIDs != nil {
			leads, leadErr := s.listScopedLeads(&leaddomain.ListLeadsRequest{
				ScopedUserIDs: salesUserIDs,
			})
			if leadErr != nil {
				return nil, leadErr
			}
			sourceCounts := make(map[string]int64)
			for _, lead := range leads {
				if lead.CreatedAt.Before(start) || lead.CreatedAt.After(end) {
					continue
				}
				source := lead.LeadSource
				if source == "" {
					source = "other"
				}
				sourceCounts[source]++
				leadsBySource.Total++
			}
			for src, count := range sourceCounts {
				entries = append(entries, dashboard.LeadsBySourceEntry{
					Source: src,
					Count:  count,
				})
			}
		} else {
			sourceCounts, err := s.leadRepo.GetStatsBySourceAndDateRange(start, end)
			if err != nil && err != gorm.ErrRecordNotFound {
				return nil, err
			}
			for src, count := range sourceCounts {
				leadsBySource.Total += count
				entries = append(entries, dashboard.LeadsBySourceEntry{
					Source: src,
					Count:  count,
				})
			}
		}
		leadsBySource.BySource = entries
	}

	// Upcoming tasks (next few open tasks, sorted by due date in repository layer)
	upcomingTasks := make([]dashboard.DashboardTaskSummary, 0)
	if s.taskRepo != nil {
		taskListReq := &taskdomain.ListTasksRequest{
			Status:  "pending",
			Page:    1,
			PerPage: 10,
		}
		if req.ScopedUserIDs != nil {
			taskListReq.ScopedUserIDs = req.ScopedUserIDs
		}
		tasks, _, err := s.taskRepo.List(taskListReq)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}

		for _, t := range tasks {
			upcomingTasks = append(upcomingTasks, dashboard.DashboardTaskSummary{
				ID:       t.ID,
				Title:    t.Title,
				Priority: t.Priority,
				Status:   t.Status,
				DueDate:  t.DueDate,
			})
		}
	}

	// Target / revenue stats from settings + pipeline summary
	targetStats := dashboard.TargetStats{}
	dealsStats := dashboard.DealsStats{}
	revenueStats := dashboard.RevenueStats{}
	pipelineStages := make([]dashboard.DashboardPipelineStageSummary, 0)

	// OPTIMIZED: Calculate revenue for current user in the period using database aggregation
	var actualRevenue int64
	// Get user to check role
	currentUser, err := s.userRepo.FindByID(userID)
	if err != nil {
		// Log error but continue if possible? No, finding user is critical for role check.
		// If user not found, we can't determine revenue scope. Return 0 revenue safely.
		// But let's proceed with default behavior (0) if user not found.
	}

	var userRoleCode string
	if currentUser != nil && currentUser.Role != nil {
		userRoleCode = currentUser.Role.Code
	}

	if userID != "" && s.dealRepo != nil {
		// Use optimized method that uses database aggregation instead of loading all deals
		if req.ScopedUserIDs != nil {
			// Scoped: aggregate won deals value for all scoped users
			for _, uid := range salesUserIDs {
				_, rev, revErr := s.dealRepo.GetWonDealsValueInPeriodByUser(uid, start, end)
				if revErr == nil {
					actualRevenue += rev
				}
			}
		} else if userRoleCode == "super_admin" || userRoleCode == "admin" {
			// Check if user is admin or super_admin to show GLOBAL revenue
			_, actualRevenue, err = s.dealRepo.GetWonDealsValueInPeriod(start, end)
		} else {
			// Otherwise show PERSONAL revenue
			_, actualRevenue, err = s.dealRepo.GetWonDealsValueInPeriodByUser(userID, start, end)
		}

		if err != nil {
			// If error, set to zero (better than failing the entire request)
			actualRevenue = 0
		}
	}

	if dealsSummary != nil {
		dealsStats = dashboard.DealsStats{
			TotalDeals:          dealsSummary.TotalDeals,
			OpenDeals:           dealsSummary.OpenDeals,
			WonDeals:            dealsSummary.WonDeals,
			LostDeals:           dealsSummary.LostDeals,
			TotalValue:          dealsSummary.TotalValue,
			TotalValueFormatted: dealsSummary.TotalValueFormatted,
			ChangePercent:       0,
		}

		// Use actual revenue if calculated, otherwise use summary (for backward compatibility)
		if actualRevenue > 0 {
			revenueStats = dashboard.RevenueStats{
				TotalRevenue:          actualRevenue,
				TotalRevenueFormatted: formatCurrency(actualRevenue),
				ChangePercent:         0,
			}
		} else {
			revenueStats = dashboard.RevenueStats{
				TotalRevenue:          dealsSummary.WonValue,
				TotalRevenueFormatted: dealsSummary.WonValueFormatted,
				ChangePercent:         0,
			}
		}

		// Map pipeline stages with percentage of total deals
		if dealsSummary.TotalDeals > 0 {
			for _, st := range dealsSummary.ByStage {
				percentage := float64(st.DealCount) * 100.0 / float64(dealsSummary.TotalDeals)
				pipelineStages = append(pipelineStages, dashboard.DashboardPipelineStageSummary{
					StageID:    st.StageID,
					StageName:  st.StageName,
					StageCode:  st.StageCode,
					DealCount:  st.DealCount,
					Percentage: percentage,
				})
			}
		}
	} else if actualRevenue > 0 {
		// If no dealsSummary but we have actualRevenue, use it
		revenueStats = dashboard.RevenueStats{
			TotalRevenue:          actualRevenue,
			TotalRevenueFormatted: formatCurrency(actualRevenue),
			ChangePercent:         0,
		}
	}

	// Calculate target stats for current user if userID is provided
	if userID != "" && s.monthlyTargetRepo != nil {
		// Get effective monthly target for the period (iterate through months)
		var totalTargetAmount int64

		if req.ScopedUserIDs != nil {
			// Scoped: aggregate targets for all scoped users using batch method
			teamTargetMap, tMapErr := s.monthlyTargetRepo.BatchGetProratedTargetsForPeriod(
				req.ScopedUserIDs,
				start.Format("2006-01-02"),
				end.Format("2006-01-02"),
			)
			if tMapErr == nil {
				for _, t := range teamTargetMap {
					totalTargetAmount += int64(t)
				}
			}
		} else {

			current := start

			// Normalize current to start of month
			current = time.Date(current.Year(), current.Month(), 1, 0, 0, 0, 0, current.Location())

			for !current.After(end) {
				targetYear := current.Year()
				targetMonth := int(current.Month())

				var monthlyTargetAmt int64
				var errTarget error

				// Check if user is admin or super_admin to show GLOBAL target
				if userRoleCode == "super_admin" || userRoleCode == "admin" {
					monthlyTargetAmt, errTarget = s.monthlyTargetRepo.GetTotalEffectiveTarget(targetYear, targetMonth)
				} else {
					// Otherwise show PERSONAL target
					effectiveTarget, err := s.monthlyTargetRepo.GetUserEffectiveTarget(userID, targetYear, targetMonth)
					if err == nil && effectiveTarget != nil {
						monthlyTargetAmt = effectiveTarget.TargetAmount
					} else {
						errTarget = err
					}
				}

				if errTarget == nil && monthlyTargetAmt > 0 {
					// Calculate proration
					daysInMonth := time.Date(targetYear, time.Month(targetMonth+1), 0, 0, 0, 0, 0, time.UTC).Day()
					activeDays := daysInMonth

					// Calculate active days interaction with start/end
					monthStart := time.Date(targetYear, time.Month(targetMonth), 1, 0, 0, 0, 0, time.UTC)
					monthEnd := time.Date(targetYear, time.Month(targetMonth+1), 0, 23, 59, 59, 0, time.UTC)

					periodStart := monthStart
					if start.After(monthStart) {
						periodStart = start
					}

					periodEnd := monthEnd
					if end.Before(monthEnd) {
						periodEnd = end
					}

					if periodStart.After(periodEnd) {
						activeDays = 0
					} else {
						d1 := time.Date(periodStart.Year(), periodStart.Month(), periodStart.Day(), 0, 0, 0, 0, time.UTC)
						d2 := time.Date(periodEnd.Year(), periodEnd.Month(), periodEnd.Day(), 0, 0, 0, 0, time.UTC)
						// Add 1 because inclusive duration
						activeDays = int(d2.Sub(d1).Hours()/24) + 1
					}

					if activeDays < 0 {
						activeDays = 0
					}
					if activeDays > daysInMonth {
						activeDays = daysInMonth
					}

					proratedAmount := int64(float64(monthlyTargetAmt) * float64(activeDays) / float64(daysInMonth))
					totalTargetAmount += proratedAmount
				}

				// Move to next month
				current = time.Date(targetYear, time.Month(targetMonth+1), 1, 0, 0, 0, 0, current.Location())
			}

		} // end else (non-scoped)

		if totalTargetAmount > 0 {
			targetAmount := totalTargetAmount
			achievedAmount := revenueStats.TotalRevenue

			var progressPercent float64
			if targetAmount > 0 {
				progressPercent = float64(achievedAmount) / float64(targetAmount) * 100
			}

			targetStats = dashboard.TargetStats{
				TargetAmount:            targetAmount,
				TargetAmountFormatted:   formatCurrency(targetAmount),
				AchievedAmount:          achievedAmount,
				AchievedAmountFormatted: revenueStats.TotalRevenueFormatted,
				ProgressPercent:         progressPercent,
				ChangePercent:           0,
			}
		}
	}

	response := &dashboard.DashboardOverviewResponse{
		Period: struct {
			Type  string    `json:"type"`
			Start time.Time `json:"start"`
			End   time.Time `json:"end"`
		}{
			Type:  req.Period,
			Start: start,
			End:   end,
		},
		VisitStats: struct {
			Total         int     `json:"total"`
			Completed     int     `json:"completed"`
			Pending       int     `json:"pending"`
			Approved      int     `json:"approved"`
			Rejected      int     `json:"rejected"`
			ChangePercent float64 `json:"change_percent"`
		}{
			Total:         visitStats.Total,
			Completed:     visitStats.Completed,
			Pending:       visitStats.Pending,
			Approved:      visitStats.Approved,
			Rejected:      visitStats.Rejected,
			ChangePercent: visitStats.ChangePercent,
		},
		AccountStats: struct {
			Total         int     `json:"total"`
			Active        int     `json:"active"`
			Inactive      int     `json:"inactive"`
			ChangePercent float64 `json:"change_percent"`
		}{
			Total:         accountStats.Total,
			Active:        accountStats.Active,
			Inactive:      accountStats.Inactive,
			ChangePercent: accountStats.ChangePercent,
		},
		ActivityStats: struct {
			Total         int     `json:"total"`
			Visits        int     `json:"visits"`
			Calls         int     `json:"calls"`
			Emails        int     `json:"emails"`
			ChangePercent float64 `json:"change_percent"`
		}{
			Total:         activityStats.Total,
			Visits:        activityStats.Visits,
			Calls:         activityStats.Calls,
			Emails:        activityStats.Emails,
			ChangePercent: activityStats.ChangePercent,
		},
		Target:         targetStats,
		Deals:          dealsStats,
		Revenue:        revenueStats,
		LeadsBySource:  leadsBySource,
		UpcomingTasks:  upcomingTasks,
		PipelineStages: pipelineStages,
		LeadStats:      leadStats,
	}

	// Cache the response for future requests
	_ = s.cacheService.SetOverview(userID, req.Period, req.StartDate, req.EndDate, response)

	return response, nil
}

// GetVisitStatistics returns visit statistics
func (s *Service) GetVisitStatistics(req *dashboard.DashboardRequest) (*dashboard.VisitStatisticsResponse, error) {
	var start, end time.Time
	if req.StartDate != "" && req.EndDate != "" {
		var err error
		start, err = time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, err
		}
		end, err = time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, err
		}
		end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, end.Location())
	} else if req.Period != "" {
		start, end = parsePeriod(req.Period)
	} else {
		start, end = parsePeriod("today")
	}

	// Try to get from cache first
	var cachedResponse dashboard.VisitStatisticsResponse
	// Use a scope-aware cache key to prevent cross-scope cache pollution
	scopeKey := ""
	if len(req.ScopedUserIDs) > 0 {
		scopeKey = strings.Join(req.ScopedUserIDs, ",")
	}
	if found, _ := s.cacheService.GetVisitStats(scopeKey, req.Period, req.StartDate, req.EndDate, &cachedResponse); found {
		return &cachedResponse, nil
	}

	// OPTIMIZED: Use database aggregation instead of loading all records
	// Get visit stats by status using aggregation.
	// When scoped, filter by salesRepID (4th param) not accountID (3rd param).
	byStatusMap := make(map[string]int64)
	if req.ScopedUserIDs != nil {
		for _, uid := range req.ScopedUserIDs {
			perUser, perErr := s.visitReportRepo.GetStatsByStatus(
				start.Format("2006-01-02"),
				end.Format("2006-01-02"),
				"", uid, "", // accountID="", salesRepID=uid
			)
			if perErr == nil {
				for status, count := range perUser {
					byStatusMap[status] += count
				}
			}
		}
	} else {
		var bsErr error
		byStatusMap, bsErr = s.visitReportRepo.GetStatsByStatus(
			start.Format("2006-01-02"),
			end.Format("2006-01-02"),
			"", "", "",
		)
		if bsErr != nil && bsErr != gorm.ErrRecordNotFound {
			return nil, bsErr
		}
	}

	// Get visit stats by date AND status.
	// GetStatsByDateAndStatus accepts salesRepID as 4th param, so we can filter per scoped user.
	byDateStatusMap := make(map[string]map[string]int64)
	if req.ScopedUserIDs != nil {
		for _, uid := range req.ScopedUserIDs {
			perDateStatus, perErr := s.visitReportRepo.GetStatsByDateAndStatus(
				start.Format("2006-01-02"),
				end.Format("2006-01-02"),
				"", uid, // accountID="", salesRepID=uid
			)
			if perErr == nil {
				for date, statusCounts := range perDateStatus {
					if byDateStatusMap[date] == nil {
						byDateStatusMap[date] = make(map[string]int64)
					}
					for status, count := range statusCounts {
						byDateStatusMap[date][status] += count
					}
				}
			}
		}
	} else {
		var bdsErr error
		byDateStatusMap, bdsErr = s.visitReportRepo.GetStatsByDateAndStatus(
			start.Format("2006-01-02"),
			end.Format("2006-01-02"),
			"", "",
		)
		if bdsErr != nil && bdsErr != gorm.ErrRecordNotFound {
			return nil, bdsErr
		}
	}

	// Convert to int for compatibility
	byStatus := make(map[string]int)
	for status, count := range byStatusMap {
		byStatus[status] = int(count)
	}

	// Calculate totals from status map
	total := int(0)
	completed := 0
	pending := 0
	approved := 0
	rejected := 0

	for status, count := range byStatus {
		total += count
		switch status {
		case "submitted", "approved":
			completed += count
			if status == "approved" {
				approved += count
			}
		case "draft":
			pending += count
		case "rejected":
			rejected += count
		}
	}

	// Build byDate from the date+status aggregation for accurate per-date status breakdown
	byDate := make(map[string]struct {
		Total     int
		Completed int
		Approved  int
		Pending   int
		Rejected  int
	})
	for date, statusCounts := range byDateStatusMap {
		entry := struct {
			Total     int
			Completed int
			Approved  int
			Pending   int
			Rejected  int
		}{}
		for st, cnt := range statusCounts {
			entry.Total += int(cnt)
			switch st {
			case "submitted", "approved":
				entry.Completed += int(cnt)
				if st == "approved" {
					entry.Approved += int(cnt)
				}
			case "draft":
				entry.Pending += int(cnt)
			case "rejected":
				entry.Rejected += int(cnt)
			}
		}
		byDate[date] = entry
	}

	// Convert byDate map to slice
	dateStats := make([]dashboard.DateStat, 0, len(byDate))
	for date, stat := range byDate {
		dateStats = append(dateStats, dashboard.DateStat{
			Date:      date,
			Count:     stat.Total,
			Completed: stat.Completed,
			Approved:  stat.Approved,
			Pending:   stat.Pending,
			Rejected:  stat.Rejected,
		})
	}

	// Sort by date
	for i := 0; i < len(dateStats)-1; i++ {
		for j := i + 1; j < len(dateStats); j++ {
			if dateStats[i].Date > dateStats[j].Date {
				dateStats[i], dateStats[j] = dateStats[j], dateStats[i]
			}
		}
	}

	response := &dashboard.VisitStatisticsResponse{
		Period: struct {
			Start time.Time `json:"start"`
			End   time.Time `json:"end"`
		}{
			Start: start,
			End:   end,
		},
		Total:         total,
		Completed:     completed,
		Pending:       pending,
		Approved:      approved,
		Rejected:      rejected,
		ByStatus:      byStatus,
		ByDate:        dateStats,
		ChangePercent: 0, // Can be calculated by comparing with previous period
	}

	// Cache the response for future requests
	_ = s.cacheService.SetVisitStats(scopeKey, req.Period, req.StartDate, req.EndDate, response)

	return response, nil
}

// GetPipelineSummary returns pipeline summary with all stages including their colors
func (s *Service) GetPipelineSummary(req *dashboard.DashboardRequest) (*dashboard.PipelineSummaryResponse, error) {
	if s.dealRepo == nil || s.pipelineRepo == nil {
		return &dashboard.PipelineSummaryResponse{
			TotalDeals: 0,
			TotalValue: 0,
			WonDeals:   0,
			LostDeals:  0,
			OpenDeals:  0,
			ByStage:    []dashboard.DashboardPipelineStageSummary{},
		}, nil
	}

	// Parse period
	var start, end time.Time
	if req.StartDate != "" && req.EndDate != "" {
		var err error
		start, err = time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, err
		}
		end, err = time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, err
		}
		end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, end.Location())
	} else if req.Period != "" {
		start, end = parsePeriod(req.Period)
	} else {
		start, end = parsePeriod("today")
	}

	// Get deal summary
	var summary *pipelinedomain.PipelineSummaryResponse
	if req.ScopedUserIDs != nil {
		// Scoped: use deal list with ScopedUserIDs and build summary
		scopedDeals, _, listErr := s.dealRepo.List(&pipelinedomain.ListDealsRequest{
			Page:          1,
			PerPage:       5000,
			ScopedUserIDs: req.ScopedUserIDs,
		})
		if listErr == nil {
			psummary := &pipelinedomain.PipelineSummaryResponse{
				ByStage: []pipelinedomain.StageSummary{},
			}
			stageAgg := make(map[string]*pipelinedomain.StageSummary)
			for _, d := range scopedDeals {
				if d.CreatedAt.Before(start) || d.CreatedAt.After(end) {
					continue
				}
				psummary.TotalDeals++
				psummary.TotalValue += d.Value
				switch d.Status {
				case "open":
					psummary.OpenDeals++
				case "won":
					psummary.WonDeals++
					psummary.WonValue += d.Value
				case "lost":
					psummary.LostDeals++
				}
				if d.Stage != nil {
					if _, ok := stageAgg[d.Stage.ID]; !ok {
						stageAgg[d.Stage.ID] = &pipelinedomain.StageSummary{
							StageID:   d.Stage.ID,
							StageName: d.Stage.Name,
							StageCode: d.Stage.Code,
						}
					}
					stageAgg[d.Stage.ID].DealCount++
					stageAgg[d.Stage.ID].TotalValue += d.Value
				}
			}
			for _, s := range stageAgg {
				s.TotalValueFormatted = formatCurrency(s.TotalValue)
				psummary.ByStage = append(psummary.ByStage, *s)
			}
			psummary.TotalValueFormatted = formatCurrency(psummary.TotalValue)
			psummary.WonValueFormatted = formatCurrency(psummary.WonValue)
			summary = psummary
		}
	} else {
		var sumErr error
		summary, sumErr = s.dealRepo.GetSummaryInPeriod(start, end)
		if sumErr != nil {
			if sumErr == gorm.ErrRecordNotFound {
				summary = &pipelinedomain.PipelineSummaryResponse{
					TotalDeals: 0,
					TotalValue: 0,
					WonDeals:   0,
					LostDeals:  0,
					OpenDeals:  0,
					ByStage:    []pipelinedomain.StageSummary{},
				}
			} else {
				return nil, sumErr
			}
		}
	}
	if summary == nil {
		summary = &pipelinedomain.PipelineSummaryResponse{
			ByStage: []pipelinedomain.StageSummary{},
		}
	}

	// Get all active pipeline stages
	listReq := &pipelinedomain.ListPipelineStagesRequest{
		IsActive: func() *bool { b := true; return &b }(),
	}
	allStages, err := s.pipelineRepo.ListStages(listReq)
	if err != nil {
		return nil, err
	}

	// Create a map of stage_id to deal stats for quick lookup
	dealStatsByStageID := make(map[string]pipelinedomain.StageSummary)
	for _, st := range summary.ByStage {
		dealStatsByStageID[st.StageID] = st
	}

	// Build response with all stages (including those with 0 deals)
	byStage := make([]dashboard.DashboardPipelineStageSummary, 0, len(allStages))
	for _, stage := range allStages {
		dealStats, hasDeals := dealStatsByStageID[stage.ID]

		stageSummary := dashboard.DashboardPipelineStageSummary{
			StageID:             stage.ID,
			StageName:           stage.Name,
			StageCode:           stage.Code,
			StageColor:          stage.Color,
			DealCount:           0,
			TotalValue:          0,
			TotalValueFormatted: formatCurrency(0),
			Percentage:          0,
		}

		if hasDeals {
			stageSummary.DealCount = dealStats.DealCount
			stageSummary.TotalValue = dealStats.TotalValue
			stageSummary.TotalValueFormatted = dealStats.TotalValueFormatted
			if summary.TotalDeals > 0 {
				stageSummary.Percentage = float64(dealStats.DealCount) / float64(summary.TotalDeals) * 100
			}
		}

		byStage = append(byStage, stageSummary)
	}

	response := &dashboard.PipelineSummaryResponse{
		TotalDeals: summary.TotalDeals,
		TotalValue: summary.TotalValue,
		WonDeals:   summary.WonDeals,
		LostDeals:  summary.LostDeals,
		OpenDeals:  summary.OpenDeals,
		ByStage:    byStage,
	}
	return response, nil
}

// formatCurrency formats integer (sen) to formatted currency string
func formatCurrency(amount int64) string {
	return currency.FormatCurrency(amount)
}

// GetTopAccounts returns top accounts
func (s *Service) GetTopAccounts(req *dashboard.DashboardRequest) ([]dashboard.TopAccountResponse, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	// Parse period
	var start, end time.Time
	if req.StartDate != "" && req.EndDate != "" {
		var err error
		start, err = time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, err
		}
		end, err = time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, err
		}
		end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, end.Location())
	} else if req.Period != "" {
		start, end = parsePeriod(req.Period)
	} else {
		start, end = parsePeriod("month") // Default to month for top lists
	}

	accountVisitCount := make(map[string]int)
	if len(req.ScopedUserIDs) > 0 {
		for _, scopedUserID := range req.ScopedUserIDs {
			accountVisitCountMap, err := s.visitReportRepo.GetStatsByAccount(
				start.Format("2006-01-02"),
				end.Format("2006-01-02"),
				scopedUserID,
				"",
			)
			if err != nil && err != gorm.ErrRecordNotFound {
				return nil, err
			}
			for accountID, count := range accountVisitCountMap {
				accountVisitCount[accountID] += int(count)
			}
		}
	} else {
		accountVisitCountMap, err := s.visitReportRepo.GetStatsByAccount(
			start.Format("2006-01-02"),
			end.Format("2006-01-02"),
			"",
			"",
		)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		for accountID, count := range accountVisitCountMap {
			accountVisitCount[accountID] = int(count)
		}
	}

	// Note: Activity stats by account not yet available via aggregation
	// For now, we'll use empty map - can be enhanced with GetStatsByAccount method later
	accountActivityCount := make(map[string]int)

	// Get accounts (still need to fetch for account details, but with reasonable pagination)
	accounts, _, err := s.accountRepo.List(&account.ListAccountsRequest{
		Page:          1,
		PerPage:       100, // Reasonable limit for top accounts
		ScopedUserIDs: req.ScopedUserIDs,
	})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Build response
	results := make([]dashboard.TopAccountResponse, 0)
	for _, acc := range accounts {
		visitCount := accountVisitCount[acc.ID]
		activityCount := accountActivityCount[acc.ID]
		// Note: LastVisitDate not available via aggregation - can be added later if needed

		results = append(results, dashboard.TopAccountResponse{
			Account: struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}{
				ID:   acc.ID,
				Name: acc.Name,
			},
			VisitCount:    visitCount,
			ActivityCount: activityCount,
			LastVisitDate: nil, // Not available via aggregation - can be enhanced later
		})
	}

	// Sort by visit count (simple implementation)
	// In production, use sort.Slice or database query with ORDER BY
	if len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

// GetTopSalesRep returns top sales reps
func (s *Service) GetTopSalesRep(req *dashboard.DashboardRequest) ([]dashboard.TopSalesRepResponse, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	// Parse period
	var start, end time.Time
	if req.StartDate != "" && req.EndDate != "" {
		var err error
		start, err = time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, err
		}
		end, err = time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, err
		}
		end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, end.Location())
	} else if req.Period != "" {
		start, end = parsePeriod(req.Period)
	} else {
		start, end = parsePeriod("month") // Default to month
	}

	salesUsers, salesUserIDs, err := s.loadSalesUsers(req.ScopedUserIDs)
	if err != nil {
		return nil, err
	}

	salesRepVisitCount := make(map[string]int)
	salesRepAccountCount := make(map[string]int)
	accountSets := make(map[string]map[string]struct{})

	visitReports, err := s.listScopedVisitReports(&visit_report.ListVisitReportsRequest{
		StartDate:     start.Format("2006-01-02"),
		EndDate:       end.Format("2006-01-02"),
		ScopedUserIDs: salesUserIDs,
	})
	if err != nil {
		return nil, err
	}
	for _, visitReportItem := range visitReports {
		if _, ok := salesUsers[visitReportItem.SalesRepID]; !ok {
			continue
		}
		salesRepVisitCount[visitReportItem.SalesRepID]++
		if visitReportItem.AccountID != nil && *visitReportItem.AccountID != "" {
			if accountSets[visitReportItem.SalesRepID] == nil {
				accountSets[visitReportItem.SalesRepID] = make(map[string]struct{})
			}
			accountSets[visitReportItem.SalesRepID][*visitReportItem.AccountID] = struct{}{}
		}
	}
	for salesRepID, accountSet := range accountSets {
		salesRepAccountCount[salesRepID] = len(accountSet)
	}

	activities, err := s.listScopedActivities(&activity.ListActivitiesRequest{
		StartDate:     start.Format("2006-01-02"),
		EndDate:       end.Format("2006-01-02"),
		ScopedUserIDs: salesUserIDs,
	})
	if err != nil {
		return nil, err
	}
	salesRepActivityCount := make(map[string]int)
	for _, activityItem := range activities {
		if _, ok := salesUsers[activityItem.UserID]; !ok {
			continue
		}
		salesRepActivityCount[activityItem.UserID]++
	}

	// Build response
	results := make([]dashboard.TopSalesRepResponse, 0)
	for salesRepID, salesUser := range salesUsers {
		visitCount := salesRepVisitCount[salesRepID]
		accountCount := salesRepAccountCount[salesRepID]
		activityCount := salesRepActivityCount[salesRepID]

		if visitCount > 0 || accountCount > 0 || activityCount > 0 {
			results = append(results, dashboard.TopSalesRepResponse{
				SalesRep: struct {
					ID    string `json:"id"`
					Name  string `json:"name"`
					Email string `json:"email"`
				}{
					ID:    salesRepID,
					Name:  salesUser.Name,
					Email: salesUser.Email,
				},
				VisitCount:    visitCount,
				AccountCount:  accountCount,
				ActivityCount: activityCount,
			})
		}
	}

	// Enrich each rep with target and actual-revenue achievement data
	if len(results) > 0 {
		periodStart := start.Format("2006-01-02")
		periodEnd := end.Format("2006-01-02")

		// Collect user IDs from results for batch target fetch
		repUserIDs := make([]string, 0, len(results))
		for _, r := range results {
			repUserIDs = append(repUserIDs, r.SalesRep.ID)
		}

		targetMap, _ := s.monthlyTargetRepo.BatchGetProratedTargetsForPeriod(repUserIDs, periodStart, periodEnd)

		for i, r := range results {
			uid := r.SalesRep.ID

			// Fetch actual revenue and deals closed for this rep
			dealsClosed, actualRev, revenueErr := s.dealRepo.GetWonDealsValueInPeriodByUser(uid, start, end)
			if revenueErr != nil {
				dealsClosed = 0
				actualRev = 0
			}

			targetAmount := int64(0)
			if t, ok := targetMap[uid]; ok {
				targetAmount = int64(t)
			}

			achievementPct := float64(0)
			if targetAmount > 0 {
				achievementPct = float64(actualRev) / float64(targetAmount) * 100
			}

			results[i].DealsClosed = dealsClosed
			results[i].ActualRevenue = actualRev
			results[i].ActualRevenueFormatted = currency.FormatCurrency(actualRev)
			results[i].TargetAmount = targetAmount
			results[i].TargetAmountFormatted = currency.FormatCurrency(targetAmount)
			results[i].TargetAchievementPercent = achievementPct
		}

		// Sort by actual revenue descending so highest performers appear first
		sort.Slice(results, func(i, j int) bool {
			return results[i].ActualRevenue > results[j].ActualRevenue
		})
	}

	if len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

// GetRecentActivities returns recent activities
func (s *Service) GetRecentActivities(req *dashboard.DashboardRequest) ([]dashboard.RecentActivityResponse, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 50
	}

	// For scoped access, fetch more activities and filter by scoped user IDs
	fetchLimit := limit
	if req.ScopedUserIDs != nil {
		fetchLimit = limit * 5 // Fetch more to ensure enough results after filtering
	}

	activities, _, err := s.activityRepo.List(&activity.ListActivitiesRequest{
		Page:    1,
		PerPage: fetchLimit,
	})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	scopeSet := scopedUserIDsSet(req.ScopedUserIDs)

	results := make([]dashboard.RecentActivityResponse, 0, len(activities))
	for _, act := range activities {
		// Filter by scoped user IDs when set
		if scopeSet != nil {
			if _, ok := scopeSet[act.UserID]; !ok {
				continue
			}
		}

		response := dashboard.RecentActivityResponse{
			ID:          act.ID,
			Type:        act.Type,
			Description: act.Description,
			User: struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}{
				ID:   act.UserID,
				Name: "", // Will be populated if needed
			},
			Timestamp: act.Timestamp,
		}

		if act.AccountID != nil {
			account, err := s.accountRepo.FindByID(*act.AccountID)
			if err == nil {
				response.Account = &struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				}{
					ID:   account.ID,
					Name: account.Name,
				}
			}
		}

		results = append(results, response)
		if len(results) >= limit {
			break
		}
	}

	return results, nil
}

// GetActivityTrends returns activity trends by date
func (s *Service) GetActivityTrends(req *dashboard.DashboardRequest) (*dashboard.ActivityTrendsResponse, error) {
	var start, end time.Time
	if req.StartDate != "" && req.EndDate != "" {
		var err error
		start, err = time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, err
		}
		end, err = time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, err
		}
		end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, end.Location())
	} else if req.Period != "" {
		start, end = parsePeriod(req.Period)
	} else {
		start, end = parsePeriod("month")
	}

	byDate := make(map[string]struct {
		Visits int
		Calls  int
		Emails int
		Total  int
	})

	if req.ScopedUserIDs != nil {
		page := 1
		for {
			activities, total, err := s.activityRepo.List(&activity.ListActivitiesRequest{
				Page:          page,
				PerPage:       500,
				StartDate:     start.Format("2006-01-02"),
				EndDate:       end.Format("2006-01-02"),
				ScopedUserIDs: req.ScopedUserIDs,
			})
			if err != nil && err != gorm.ErrRecordNotFound {
				return nil, err
			}

			for _, activityItem := range activities {
				date := activityItem.Timestamp.Format("2006-01-02")
				stat := byDate[date]
				stat.Total++

				switch activityItem.Type {
				case "visit":
					stat.Visits++
				case "call":
					stat.Calls++
				case "email":
					stat.Emails++
				}

				byDate[date] = stat
			}

			if int64(page*500) >= total || len(activities) == 0 {
				break
			}
			page++
		}
	} else {
		activitiesByTypeAndDate, err := s.activityRepo.GetStatsByTypeAndDate(
			start.Format("2006-01-02"),
			end.Format("2006-01-02"),
			"",
		)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}

		for activityType, dateStats := range activitiesByTypeAndDate {
			for date, count := range dateStats {
				stat := byDate[date]
				stat.Total += int(count)

				switch activityType {
				case "visit":
					stat.Visits += int(count)
				case "call":
					stat.Calls += int(count)
				case "email":
					stat.Emails += int(count)
				}

				byDate[date] = stat
			}
		}
	}

	// Convert to slice and sort by date
	dateStats := make([]dashboard.ActivityDateStat, 0, len(byDate))
	for date, stat := range byDate {
		dateStats = append(dateStats, dashboard.ActivityDateStat{
			Date:   date,
			Visits: stat.Visits,
			Calls:  stat.Calls,
			Emails: stat.Emails,
			Total:  stat.Total,
		})
	}

	// Sort by date
	for i := 0; i < len(dateStats)-1; i++ {
		for j := i + 1; j < len(dateStats); j++ {
			if dateStats[i].Date > dateStats[j].Date {
				dateStats[i], dateStats[j] = dateStats[j], dateStats[i]
			}
		}
	}

	response := &dashboard.ActivityTrendsResponse{
		Period: struct {
			Start time.Time `json:"start"`
			End   time.Time `json:"end"`
		}{
			Start: start,
			End:   end,
		},
		ByDate: dateStats,
	}

	return response, nil
}

// ============================================================================
// Super Admin Dashboard Methods
// ============================================================================

// GetSuperAdminUsersByRole returns users grouped by role
func (s *Service) GetSuperAdminUsersByRole() (*dashboard.SuperAdminUsersByRoleResponse, error) {
	// Get all roles
	roles, err := s.roleRepo.List()
	if err != nil {
		return nil, err
	}

	usersByRole := make([]dashboard.UsersByRoleEntry, 0, len(roles))
	var totalUsers, totalActive, totalInactive int64

	for _, role := range roles {
		// Count total users for this role
		totalCount, err := s.userRepo.CountUsersByRoleID(role.ID)
		if err != nil {
			continue // Skip this role if error
		}

		// Count active users
		// OPTIMIZED: Use count query instead of loading all users
		// Get active users count for this role
		var activeCount int64
		// Note: CountUsersByRoleID doesn't filter by status, so we need to use List with count
		// But we can optimize by using a reasonable limit and just getting the count
		_, totalActive, err := s.userRepo.List(&user.ListUsersRequest{
			RoleID:  role.ID,
			Status:  "active",
			Page:    1,
			PerPage: 1, // Just need the count
		})
		if err == nil {
			activeCount = totalActive
		}

		inactiveCount := totalCount - activeCount

		usersByRole = append(usersByRole, dashboard.UsersByRoleEntry{
			RoleCode:      role.Code,
			RoleName:      role.Name,
			TotalUsers:    totalCount,
			ActiveUsers:   activeCount,
			InactiveUsers: inactiveCount,
		})

		totalUsers += totalCount
		totalActive += activeCount
		totalInactive += inactiveCount
	}

	response := &dashboard.SuperAdminUsersByRoleResponse{
		UsersByRole:   usersByRole,
		TotalUsers:    totalUsers,
		TotalActive:   totalActive,
		TotalInactive: totalInactive,
	}

	return response, nil
}

// GetSuperAdminSystemActivity returns system activity logs (placeholder - needs activity log table)
func (s *Service) GetSuperAdminSystemActivity(req *dashboard.DashboardRequest) (*dashboard.SuperAdminSystemActivityResponse, error) {
	// TODO: Implement when activity log table is available
	// For now, return empty response
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}

	response := &dashboard.SuperAdminSystemActivityResponse{
		Activities:     []dashboard.SystemActivityEntry{},
		Total:          0,
		RecentErrors:   0,
		RecentWarnings: 0,
	}

	return response, nil
}

// GetSuperAdminAIUsage returns AI usage and cost statistics (placeholder - needs AI usage log table)
func (s *Service) GetSuperAdminAIUsage(req *dashboard.DashboardRequest) (*dashboard.SuperAdminAIUsageResponse, error) {
	// TODO: Implement when AI usage log table is available
	// For now, return mock data
	response := &dashboard.SuperAdminAIUsageResponse{
		TotalRequests:     0,
		RequestsToday:     0,
		RequestsThisWeek:  0,
		RequestsThisMonth: 0,
		EstimatedCost: dashboard.AIUsageCost{
			Today:    0,
			Week:     0,
			Month:    0,
			Currency: "USD",
		},
		SuccessRate:         0,
		FallbackRate:        0,
		AverageResponseTime: 0,
	}

	return response, nil
}

// GetSuperAdminDataGrowth returns data growth statistics
func (s *Service) GetSuperAdminDataGrowth(req *dashboard.DashboardRequest) (*dashboard.SuperAdminDataGrowthResponse, error) {
	// Parse period
	var start, end time.Time
	if req.StartDate != "" && req.EndDate != "" {
		var err error
		start, err = time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, err
		}
		end, err = time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, err
		}
		end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, end.Location())
	} else if req.Period != "" {
		start, end = parsePeriod(req.Period)
	} else {
		// Default to month
		start, end = parsePeriod("month")
	}

	periodStr := req.Period
	if periodStr == "" {
		periodStr = "month"
	}

	// Get accounts count - get all and filter by created_at
	var totalAccounts int64
	var accountsLastPeriod int64
	var err error

	// OPTIMIZED: Use database count query instead of loading all records
	totalAccounts, err = s.accountRepo.CountByDateRange(start, end)
	if err != nil {
		totalAccounts = 0
	}

	// Get previous period for comparison
	prevStart := start.AddDate(0, 0, -int(end.Sub(start).Hours()/24))
	prevEnd := start.Add(-time.Second)
	accountsLastPeriod, err = s.accountRepo.CountByDateRange(prevStart, prevEnd)
	if err != nil {
		accountsLastPeriod = 0
	}

	accountsGrowthCount := totalAccounts - accountsLastPeriod
	accountsGrowthPercent := float64(0)
	if accountsLastPeriod > 0 {
		accountsGrowthPercent = float64(accountsGrowthCount) / float64(accountsLastPeriod) * 100
	}

	// OPTIMIZED: Use database count query instead of loading all records
	var totalLeads int64
	var leadsLastPeriod int64

	totalLeads, err = s.leadRepo.CountByDateRange(start, end)
	if err != nil {
		totalLeads = 0
	}

	// Reuse prevStart and prevEnd from accounts calculation
	leadsLastPeriod, err = s.leadRepo.CountByDateRange(prevStart, prevEnd)
	if err != nil {
		leadsLastPeriod = 0
	}

	leadsGrowthCount := totalLeads - leadsLastPeriod
	leadsGrowthPercent := float64(0)
	if leadsLastPeriod > 0 {
		leadsGrowthPercent = float64(leadsGrowthCount) / float64(leadsLastPeriod) * 100
	}

	// OPTIMIZED: Use database count query instead of loading all records
	var totalDeals int64
	var dealsLastPeriod int64

	totalDeals, err = s.dealRepo.CountByDateRange(start, end)
	if err != nil {
		totalDeals = 0
	}

	// Reuse prevStart and prevEnd from leads calculation
	dealsLastPeriod, err = s.dealRepo.CountByDateRange(prevStart, prevEnd)
	if err != nil {
		dealsLastPeriod = 0
	}

	dealsGrowthCount := totalDeals - dealsLastPeriod
	dealsGrowthPercent := float64(0)
	if dealsLastPeriod > 0 {
		dealsGrowthPercent = float64(dealsGrowthCount) / float64(dealsLastPeriod) * 100
	}

	response := &dashboard.SuperAdminDataGrowthResponse{
		Accounts: dashboard.DataGrowthStats{
			Total:         totalAccounts,
			GrowthPercent: accountsGrowthPercent,
			GrowthCount:   accountsGrowthCount,
			Period:        periodStr,
		},
		Leads: dashboard.DataGrowthStats{
			Total:         totalLeads,
			GrowthPercent: leadsGrowthPercent,
			GrowthCount:   leadsGrowthCount,
			Period:        periodStr,
		},
		Deals: dashboard.DataGrowthStats{
			Total:         totalDeals,
			GrowthPercent: dealsGrowthPercent,
			GrowthCount:   dealsGrowthCount,
			Period:        periodStr,
		},
	}

	return response, nil
}

// GetSuperAdminErrorSummary returns error and failed process summary (placeholder)
func (s *Service) GetSuperAdminErrorSummary(req *dashboard.DashboardRequest) (*dashboard.SuperAdminErrorSummaryResponse, error) {
	// TODO: Implement when error log table is available
	response := &dashboard.SuperAdminErrorSummaryResponse{
		TotalErrors:     0,
		ErrorsToday:     0,
		ErrorsThisWeek:  0,
		ErrorsThisMonth: 0,
		FailedProcesses: []dashboard.FailedProcessEntry{},
		ErrorTypes:      []dashboard.ErrorTypeEntry{},
	}

	return response, nil
}

// ============================================================================
// Admin Dashboard Methods
// ============================================================================

// GetAdminTotalLeads returns total leads statistics for today and this month
func (s *Service) GetAdminTotalLeads() (*dashboard.AdminTotalLeadsResponse, error) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())

	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthEnd := now

	// OPTIMIZED: Use database aggregation instead of loading all records
	// Get today stats by status using aggregation
	todayStatsMap, err := s.leadRepo.GetStatsByStatusAndDateRange(todayStart, todayEnd)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Get month stats by status using aggregation
	monthStatsMap, err := s.leadRepo.GetStatsByStatusAndDateRange(monthStart, monthEnd)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Get last month total for comparison
	lastMonthStart := monthStart.AddDate(0, -1, 0)
	lastMonthEnd := monthStart.Add(-time.Second)
	lastMonthTotal, err := s.leadRepo.CountByDateRange(lastMonthStart, lastMonthEnd)
	if err != nil {
		lastMonthTotal = 0
	}

	// Build today stats from aggregation
	todayStats := dashboard.LeadPeriodStats{}
	for status, count := range todayStatsMap {
		todayStats.Total += count
		switch status {
		case "new":
			todayStats.New += count
		case "contacted", "interested":
			todayStats.Contacted += count
		case "qualified", "proposal_sent":
			todayStats.Qualified += count
		case "converted":
			todayStats.Converted += count
		}
	}

	// Build month stats from aggregation
	monthStats := dashboard.LeadPeriodStats{}
	for status, count := range monthStatsMap {
		monthStats.Total += count
		switch status {
		case "new":
			monthStats.New += count
		case "contacted", "interested":
			monthStats.Contacted += count
		case "qualified", "proposal_sent":
			monthStats.Qualified += count
		case "converted":
			monthStats.Converted += count
		}
	}

	changePercent := float64(0)
	if lastMonthTotal > 0 {
		changePercent = float64(monthStats.Total-int64(lastMonthTotal)) / float64(lastMonthTotal) * 100
	}

	response := &dashboard.AdminTotalLeadsResponse{
		Today:         todayStats,
		ThisMonth:     monthStats,
		ChangePercent: changePercent,
	}

	return response, nil
}

// GetAdminPipelineValue returns pipeline value summary
func (s *Service) GetAdminPipelineValue(req *dashboard.DashboardRequest) (*dashboard.AdminPipelineValueResponse, error) {
	// Get deals summary
	dealsSummary, err := s.dealRepo.GetSummary()
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if dealsSummary == nil {
		dealsSummary = &pipelinedomain.PipelineSummaryResponse{}
	}

	// Format currency
	formatCurrency := func(amount int64) string {
		return fmt.Sprintf("Rp %s", strings.ReplaceAll(fmt.Sprintf("%d", amount), ",", "."))
	}

	// Calculate change percent (compare with previous period)
	// Get previous period (same period last month)
	now := time.Now()
	prevMonthStart := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location())
	prevMonthEnd := prevMonthStart.AddDate(0, 1, 0).Add(-time.Second)

	// Get previous period summary
	prevDealsSummary, err := s.dealRepo.GetForecast("month", prevMonthStart, prevMonthEnd)
	changePercent := float64(0)
	if err == nil && prevDealsSummary != nil && prevDealsSummary.ExpectedRevenue > 0 {
		currentValue := dealsSummary.TotalValue
		prevValue := prevDealsSummary.ExpectedRevenue
		changePercent = float64(currentValue-prevValue) / float64(prevValue) * 100
	}

	response := &dashboard.AdminPipelineValueResponse{
		TotalValue:          dealsSummary.TotalValue,
		TotalValueFormatted: formatCurrency(dealsSummary.TotalValue),
		OpenDealsValue:      dealsSummary.OpenValue,
		WonDealsValue:       dealsSummary.WonValue,
		LostDealsValue:      dealsSummary.LostValue,
		ChangePercent:       changePercent,
	}

	return response, nil
}

// GetAdminPendingApprovals returns pending approvals
func (s *Service) GetAdminPendingApprovals() (*dashboard.AdminPendingApprovalsResponse, error) {
	// Get visit reports with pending status
	visitReports, _, err := s.visitReportRepo.List(&visit_report.ListVisitReportsRequest{
		Status:  "submitted", // Submitted but not approved/rejected
		Page:    1,
		PerPage: 100,
	})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	items := make([]dashboard.PendingApprovalItem, 0)
	visitReportsCount := int64(0)

	for _, vr := range visitReports {
		if vr.Status == "submitted" || vr.Status == "draft" {
			visitReportsCount++
			// Get user name
			userName := "Unknown"
			if vr.SalesRepID != "" {
				user, err := s.userRepo.FindByID(vr.SalesRepID)
				if err == nil {
					userName = user.Name
				}
			}

			priority := "medium"
			if vr.VisitDate.Before(time.Now().AddDate(0, 0, -7)) {
				priority = "high"
			}

			items = append(items, dashboard.PendingApprovalItem{
				ID:          vr.ID,
				Type:        "visit_report",
				Title:       fmt.Sprintf("Visit Report - %s", vr.Purpose),
				SubmittedBy: userName,
				SubmittedAt: vr.CreatedAt,
				Priority:    priority,
			})
		}
	}

	response := &dashboard.AdminPendingApprovalsResponse{
		Total:          int64(len(items)),
		VisitReports:   visitReportsCount,
		ExpenseReports: 0, // TODO: Implement when expense reports are available
		Other:          0,
		Items:          items,
	}

	return response, nil
}

// GetAdminTaskOverdue returns global overdue tasks
func (s *Service) GetAdminTaskOverdue(req *dashboard.DashboardRequest) (*dashboard.AdminTaskOverdueResponse, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 50
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// OPTIMIZED: Use reasonable limit instead of loading all tasks
	// For overdue tasks, we only need tasks that are actually overdue
	tasks, _, err := s.taskRepo.List(&taskdomain.ListTasksRequest{
		Status:  "pending",
		Page:    1,
		PerPage: 500, // Reasonable limit for overdue tasks
	})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	overdueTasks := make([]dashboard.OverdueTaskItem, 0)
	var totalOverdue, criticalOverdue int64

	for _, task := range tasks {
		if task.DueDate != nil && task.DueDate.Before(today) {
			daysOverdue := int(today.Sub(*task.DueDate).Hours() / 24)
			totalOverdue++

			// Get assigned user name
			assignedToName := "Unassigned"
			if task.AssignedTo != nil && *task.AssignedTo != "" {
				user, err := s.userRepo.FindByID(*task.AssignedTo)
				if err == nil {
					assignedToName = user.Name
				}
			}

			priority := task.Priority
			if daysOverdue > 7 || priority == "high" {
				criticalOverdue++
			}

			overdueTasks = append(overdueTasks, dashboard.OverdueTaskItem{
				ID:          task.ID,
				Title:       task.Title,
				AssignedTo:  assignedToName,
				DueDate:     *task.DueDate,
				DaysOverdue: daysOverdue,
				Priority:    priority,
			})
		}
	}

	// Sort by days overdue (most overdue first) and limit
	if len(overdueTasks) > limit {
		// Simple sort by days overdue
		for i := 0; i < len(overdueTasks)-1; i++ {
			for j := i + 1; j < len(overdueTasks); j++ {
				if overdueTasks[i].DaysOverdue < overdueTasks[j].DaysOverdue {
					overdueTasks[i], overdueTasks[j] = overdueTasks[j], overdueTasks[i]
				}
			}
		}
		overdueTasks = overdueTasks[:limit]
	}

	response := &dashboard.AdminTaskOverdueResponse{
		TotalOverdue:    totalOverdue,
		CriticalOverdue: criticalOverdue,
		Tasks:           overdueTasks,
	}

	return response, nil
}

// ============================================================================
// Sales Manager Dashboard Methods
// ============================================================================

// GetSalesManagerPipelineFunnel returns pipeline funnel data
func (s *Service) GetSalesManagerPipelineFunnel(req *dashboard.DashboardRequest) (*dashboard.SalesManagerPipelineFunnelResponse, error) {
	// OPTIMIZED: Use database aggregation instead of loading all deals
	// When scoped, aggregate deal stats per user; otherwise fetch globally
	dealStatsByStage := make(map[string]int64)
	if req.ScopedUserIDs != nil {
		for _, uid := range req.ScopedUserIDs {
			perUser, perErr := s.dealRepo.GetStatsByStage("", "", uid, "open")
			if perErr == nil {
				for stageID, count := range perUser {
					dealStatsByStage[stageID] += count
				}
			}
		}
	} else {
		var dsErr error
		dealStatsByStage, dsErr = s.dealRepo.GetStatsByStage("", "", "", "open")
		if dsErr != nil && dsErr != gorm.ErrRecordNotFound {
			return nil, dsErr
		}
	}

	// Get all stages to build proper funnel
	stages, err := s.pipelineRepo.ListStages(&pipelinedomain.ListPipelineStagesRequest{})
	if err != nil {
		return nil, err
	}

	// Build stage ID to name map
	stageMap := make(map[string]string)
	for _, stage := range stages {
		stageMap[stage.ID] = stage.Name
	}

	// Count by stage from aggregation
	stageCounts := make(map[string]int64)
	totalLeads := int64(0)
	for stageID, count := range dealStatsByStage {
		stageName := stageMap[stageID]
		if stageName != "" {
			stageCounts[stageName] = count
			totalLeads += count
		}
	}

	// Build funnel from stages (already loaded above)
	funnel := make([]dashboard.FunnelStageEntry, 0)
	for _, stage := range stages {
		if !stage.IsWon && !stage.IsLost {
			count := stageCounts[stage.Name]
			percentage := float64(0)
			if totalLeads > 0 {
				percentage = float64(count) / float64(totalLeads) * 100
			}
			funnel = append(funnel, dashboard.FunnelStageEntry{
				Stage:      stage.Name,
				Count:      int64(count),
				Percentage: percentage,
			})
		}
	}

	// Add won stage - get won deals count using aggregation.
	// Apply scoped user filter so conversion reflects this manager's team only.
	wonCount := int64(0)
	if req.ScopedUserIDs != nil {
		for _, uid := range req.ScopedUserIDs {
			perUser, perErr := s.dealRepo.GetStatsByStatus("", "", uid, "", "won")
			if perErr == nil {
				for _, count := range perUser {
					wonCount += count
				}
			}
		}
	} else {
		wonStats, wonErr := s.dealRepo.GetStatsByStatus("", "", "", "", "won")
		if wonErr == nil {
			for _, count := range wonStats {
				wonCount += count
			}
		}
	}

	// conversionRate = won / (open + won) * 100.
	// Using (totalLeads + wonCount) as denominator prevents rate from exceeding 100%
	// because open deals cannot be less than historically-won deals.
	wonPercentage := float64(0)
	if totalLeads+wonCount > 0 {
		wonPercentage = float64(wonCount) / float64(totalLeads+wonCount) * 100
	}
	funnel = append(funnel, dashboard.FunnelStageEntry{
		Stage:      "Won",
		Count:      wonCount,
		Percentage: wonPercentage,
	})

	conversionRate := wonPercentage

	response := &dashboard.SalesManagerPipelineFunnelResponse{
		Funnel:         funnel,
		ConversionRate: conversionRate,
	}

	return response, nil
}

// GetSalesManagerTargetVsActual returns target vs actual performance
func (s *Service) GetSalesManagerTargetVsActual(userID string, req *dashboard.DashboardRequest) (*dashboard.SalesManagerTargetVsActualResponse, error) {
	period := req.Period
	if period == "" {
		period = "month"
	}

	// Get period dates
	var start, end time.Time
	if req.StartDate != "" && req.EndDate != "" {
		var err error
		parsedStart, err := time.Parse("2006-01-02", req.StartDate)
		if err == nil {
			// Normalize start date to beginning of day (00:00:00)
			start = time.Date(parsedStart.Year(), parsedStart.Month(), parsedStart.Day(), 0, 0, 0, 0, parsedStart.Location())
			parsedEnd, err := time.Parse("2006-01-02", req.EndDate)
			if err == nil {
				// Normalize end date to end of day (23:59:59.999999999)
				end = time.Date(parsedEnd.Year(), parsedEnd.Month(), parsedEnd.Day(), 23, 59, 59, 999999999, parsedEnd.Location())
			}
		}
	} else {
		start, end = parsePeriod(period)
	}

	// Resolve the full team under this manager by looking up bricks they manage,
	// then collecting all sales reps assigned to those bricks.
	// Prefer ScopedUserIDs from middleware if available (consistent with other endpoints).
	teamUserIDs := req.ScopedUserIDs
	if teamUserIDs == nil {
		teamUserIDs = s.resolveTeamUserIDs(userID)
	}

	// Get aggregate monthly targets for all team members (prorated for custom date ranges)
	var targetRevenue int64
	var targetDeals int64 = 0  // Deals target not in monthly_target table yet
	var targetVisits int64 = 0 // Visits target not in monthly_target table yet

	teamTargetMap, err := s.monthlyTargetRepo.BatchGetProratedTargetsForPeriod(
		teamUserIDs,
		start.Format("2006-01-02"),
		end.Format("2006-01-02"),
	)
	if err == nil {
		for _, t := range teamTargetMap {
			targetRevenue += int64(t)
		}
	}

	// Aggregate actual won-deal revenue and deal count for all team members
	var actualRevenue, actualDeals int64
	for _, uid := range teamUserIDs {
		deals, rev, dealErr := s.dealRepo.GetWonDealsValueInPeriodByUser(uid, start, end)
		if dealErr == nil {
			actualRevenue += rev
			actualDeals += deals
		}
	}

	// Aggregate visit counts for all team members
	var actualVisits int64
	for _, uid := range teamUserIDs {
		visitStatsByStatus, visitErr := s.visitReportRepo.GetStatsByStatus(
			start.Format("2006-01-02"),
			end.Format("2006-01-02"),
			uid, "", "",
		)
		if visitErr == nil {
			for _, count := range visitStatsByStatus {
				actualVisits += count
			}
		}
	}

	target := dashboard.TargetMetrics{
		Revenue: targetRevenue,
		Deals:   targetDeals,
		Visits:  targetVisits,
	}

	actual := dashboard.ActualMetrics{
		Revenue: actualRevenue,
		Deals:   actualDeals,
		Visits:  actualVisits,
	}

	// Calculate achievement
	achievement := dashboard.AchievementMetrics{
		RevenuePercent: 0,
		DealsPercent:   0,
		VisitsPercent:  0,
	}
	if target.Revenue > 0 {
		achievement.RevenuePercent = float64(actualRevenue) / float64(target.Revenue) * 100
	}
	if target.Deals > 0 {
		achievement.DealsPercent = float64(actualDeals) / float64(target.Deals) * 100
	}
	if target.Visits > 0 {
		achievement.VisitsPercent = float64(actualVisits) / float64(target.Visits) * 100
	}

	// Calculate gap
	gap := dashboard.GapMetrics{
		Revenue: target.Revenue - actualRevenue,
		Deals:   target.Deals - actualDeals,
		Visits:  target.Visits - actualVisits,
	}

	response := &dashboard.SalesManagerTargetVsActualResponse{
		Period:      period,
		Target:      target,
		Actual:      actual,
		Achievement: achievement,
		Gap:         gap,
	}

	return response, nil
}

// GetSalesManagerVisitCompletion returns visit completion statistics
func (s *Service) GetSalesManagerVisitCompletion(req *dashboard.DashboardRequest) (*dashboard.SalesManagerVisitCompletionResponse, error) {
	// Parse period
	var start, end time.Time
	if req.StartDate != "" && req.EndDate != "" {
		var err error
		start, err = time.Parse("2006-01-02", req.StartDate)
		if err == nil {
			end, err = time.Parse("2006-01-02", req.EndDate)
			if err == nil {
				end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, end.Location())
			}
		}
	} else if req.Period != "" {
		start, end = parsePeriod(req.Period)
	} else {
		start, end = parsePeriod("month")
	}

	// OPTIMIZED: Use database aggregation instead of loading all records.
	// When scoped, use GetStatsByStatus filtered by salesRepID to build per-rep totals,
	// because GetStatsBySalesRep takes accountID (not salesRepID) as its 3rd param.
	visitStatsBySalesRep := make(map[string]int64)
	if req.ScopedUserIDs != nil {
		for _, uid := range req.ScopedUserIDs {
			// Aggregate total visits for this sales rep across all statuses
			perUser, perErr := s.visitReportRepo.GetStatsByStatus(
				start.Format("2006-01-02"),
				end.Format("2006-01-02"),
				"", uid, "", // accountID="", salesRepID=uid
			)
			if perErr == nil {
				for _, count := range perUser {
					visitStatsBySalesRep[uid] += count
				}
			}
		}
	} else {
		var vsrErr error
		visitStatsBySalesRep, vsrErr = s.visitReportRepo.GetStatsBySalesRep(
			start.Format("2006-01-02"),
			end.Format("2006-01-02"),
			"", "",
		)
		if vsrErr != nil && vsrErr != gorm.ErrRecordNotFound {
			return nil, vsrErr
		}
	}

	// Get visit stats by status for completion calculation.
	// Filter by salesRepID (4th param) not accountID (3rd param).
	visitStatsByStatus := make(map[string]int64)
	if req.ScopedUserIDs != nil {
		for _, uid := range req.ScopedUserIDs {
			perUser, perErr := s.visitReportRepo.GetStatsByStatus(
				start.Format("2006-01-02"),
				end.Format("2006-01-02"),
				"", uid, "", // accountID="", salesRepID=uid
			)
			if perErr == nil {
				for status, count := range perUser {
					visitStatsByStatus[status] += count
				}
			}
		}
	} else {
		var vssErr error
		visitStatsByStatus, vssErr = s.visitReportRepo.GetStatsByStatus(
			start.Format("2006-01-02"),
			end.Format("2006-01-02"),
			"", "", "",
		)
		if vssErr != nil && vssErr != gorm.ErrRecordNotFound {
			return nil, vssErr
		}
	}

	// Calculate totals from aggregation
	var totalScheduled, completed, pending, missed int64
	for _, count := range visitStatsBySalesRep {
		totalScheduled += count
	}
	for status, count := range visitStatsByStatus {
		switch status {
		case "submitted", "approved":
			completed += count
		case "draft", "pending":
			pending += count
		case "rejected", "cancelled":
			missed += count
		}
	}

	completionRate := float64(0)
	if totalScheduled > 0 {
		completionRate = float64(completed) / float64(totalScheduled) * 100
	}

	// Build by sales rep list from aggregation
	// Note: For detailed completion rate per sales rep, we need GetStatsBySalesRepAndStatus method
	// For now, using simplified version with overall completion rate
	bySalesRep := make([]dashboard.VisitCompletionBySalesRep, 0)
	for salesRepID, scheduledCount := range visitStatsBySalesRep {
		user, err := s.userRepo.FindByID(salesRepID)
		userName := "Unknown"
		if err == nil {
			userName = user.Name
		}

		// Simplified completion rate (using overall completion rate)
		// Can be enhanced with GetStatsBySalesRepAndStatus method
		repCompletionRate := completionRate
		completedCount := int64(float64(scheduledCount) * completionRate / 100)

		bySalesRep = append(bySalesRep, dashboard.VisitCompletionBySalesRep{
			SalesRepID:     salesRepID,
			SalesRepName:   userName,
			Scheduled:      scheduledCount,
			Completed:      completedCount,
			CompletionRate: repCompletionRate,
		})
	}

	response := &dashboard.SalesManagerVisitCompletionResponse{
		TotalScheduled: totalScheduled,
		Completed:      completed,
		Pending:        pending,
		Missed:         missed,
		CompletionRate: completionRate,
		BySalesRep:     bySalesRep,
	}

	return response, nil
}

// GetSalesManagerDealsAtRisk returns deals at risk
func (s *Service) GetSalesManagerDealsAtRisk(req *dashboard.DashboardRequest) (*dashboard.SalesManagerDealsAtRiskResponse, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}

	now := time.Now()
	sevenDaysAgo := now.AddDate(0, 0, -7)

	// OPTIMIZED: Use reasonable limit instead of loading all deals
	// For deals at risk, we only need deals that are actually at risk
	dealsListReq := &pipelinedomain.ListDealsRequest{
		Status:  "open",
		Page:    1,
		PerPage: 500, // Reasonable limit for deals at risk
	}
	if req.ScopedUserIDs != nil {
		dealsListReq.ScopedUserIDs = req.ScopedUserIDs
	}
	deals, _, err := s.dealRepo.List(dealsListReq)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	dealsAtRisk := make([]dashboard.DealAtRiskItem, 0)

	for _, deal := range deals {
		// Check if deal has no recent activity
		// Get last activity for this account
		var lastActivity *time.Time
		if deal.AccountID != "" {
			activities, _, err := s.activityRepo.List(&activity.ListActivitiesRequest{
				AccountID: deal.AccountID,
				Page:      1,
				PerPage:   1,
			})
			if err == nil && len(activities) > 0 {
				lastActivity = &activities[0].Timestamp
			}
		}

		daysWithoutActivity := 0
		riskReason := "low_probability"
		if lastActivity != nil {
			daysWithoutActivity = int(now.Sub(*lastActivity).Hours() / 24)
			if daysWithoutActivity > 7 {
				riskReason = "no_activity"
			}
		} else if deal.CreatedAt.Before(sevenDaysAgo) {
			daysWithoutActivity = int(now.Sub(deal.CreatedAt).Hours() / 24)
			riskReason = "stale"
		}

		if deal.Probability < 30 {
			riskReason = "low_probability"
		}

		// Only include deals that are genuinely at risk (stale or low probability).
		// riskReason is excluded from the condition to prevent all open deals from being flagged.
		if daysWithoutActivity > 7 || deal.Probability < 30 {
			assignedToName := "Unassigned"
			if deal.AssignedTo != nil && *deal.AssignedTo != "" {
				user, err := s.userRepo.FindByID(*deal.AssignedTo)
				if err == nil {
					assignedToName = user.Name
				}
			}

			stageName := "Unknown"
			if deal.Stage != nil {
				stageName = deal.Stage.Name
			}

			lastAct := deal.UpdatedAt
			if lastActivity != nil {
				lastAct = *lastActivity
			}

			dealsAtRisk = append(dealsAtRisk, dashboard.DealAtRiskItem{
				ID:                  deal.ID,
				Name:                deal.Title,
				Value:               deal.Value,
				Stage:               stageName,
				DaysWithoutActivity: daysWithoutActivity,
				RiskReason:          riskReason,
				AssignedTo:          assignedToName,
				LastActivity:        lastAct,
			})
		}
	}

	// Sort by days without activity (most at risk first) and limit
	if len(dealsAtRisk) > limit {
		for i := 0; i < len(dealsAtRisk)-1; i++ {
			for j := i + 1; j < len(dealsAtRisk); j++ {
				if dealsAtRisk[i].DaysWithoutActivity < dealsAtRisk[j].DaysWithoutActivity {
					dealsAtRisk[i], dealsAtRisk[j] = dealsAtRisk[j], dealsAtRisk[i]
				}
			}
		}
		dealsAtRisk = dealsAtRisk[:limit]
	}

	response := &dashboard.SalesManagerDealsAtRiskResponse{
		TotalAtRisk: int64(len(dealsAtRisk)),
		Deals:       dealsAtRisk,
	}

	return response, nil
}

// GetSalesManagerTeamDraftApprovals aggregates submitted/pending items from all team
// members under the given manager that require review and approval.
// Approval workflow per entity:
//   - Visit Reports : draft → submitted → approved / rejected
//   - Tasks         : pending/draft → submitted → approved / rejected
//   - Schedules     : pending/draft → submitted → confirmed / cancelled
//   - Leads         : new (= just submitted for manager qualification)
//   - Pipeline      : open deals from team (needing manager oversight)
func (s *Service) GetSalesManagerTeamDraftApprovals(managerID string) (*dashboard.SalesManagerTeamDraftApprovalsResponse, error) {
	const maxItems = 10

	// Resolve team member IDs (exclude manager's own ID from the list if desired,
	// but resolveTeamUserIDs always includes the manager – keep it for completeness).
	teamIDs := s.resolveTeamUserIDs(managerID)
	teamOnlyIDs := make([]string, 0, len(teamIDs))
	for _, id := range teamIDs {
		if id != managerID {
			teamOnlyIDs = append(teamOnlyIDs, id)
		}
	}
	if len(teamOnlyIDs) == 0 {
		return &dashboard.SalesManagerTeamDraftApprovalsResponse{
			Visits:    dashboard.SalesManagerDraftVisitsData{Items: []dashboard.DraftVisitItem{}},
			Tasks:     dashboard.SalesManagerDraftTasksData{Items: []dashboard.DraftTaskItem{}},
			Schedules: dashboard.SalesManagerDraftSchedulesData{Items: []dashboard.DraftScheduleItem{}},
			Leads:     dashboard.SalesManagerDraftLeadsData{Items: []dashboard.DraftLeadItem{}},
			Pipeline:  dashboard.SalesManagerDraftPipelineData{Items: []dashboard.DraftPipelineItem{}},
		}, nil
	}

	// ── Visit Reports (submitted = awaiting manager approval) ──────────────────
	vrList, _, _ := s.visitReportRepo.List(&visit_report.ListVisitReportsRequest{
		Status:        "submitted",
		ScopedUserIDs: teamOnlyIDs,
		Page:          1,
		PerPage:       100,
	})
	vrItems := make([]dashboard.DraftVisitItem, 0, maxItems)
	for i, vr := range vrList {
		if i >= maxItems {
			break
		}
		assignedTo := vr.SalesRepID
		if u, err := s.userRepo.FindByID(vr.SalesRepID); err == nil {
			assignedTo = u.Name
		}
		vrItems = append(vrItems, dashboard.DraftVisitItem{
			ID:         vr.ID,
			Purpose:    vr.Purpose,
			Status:     vr.Status,
			VisitDate:  vr.VisitDate,
			AssignedTo: assignedTo,
			CreatedAt:  vr.CreatedAt,
		})
	}

	// ── Tasks (submitted = sales explicitly submitted task for manager review) ──
	taskList, _, _ := s.taskRepo.List(&taskdomain.ListTasksRequest{
		Status:        "pending,submitted",
		ScopedUserIDs: teamOnlyIDs,
		Page:          1,
		PerPage:       100,
	})
	taskItems := make([]dashboard.DraftTaskItem, 0, maxItems)
	for i, t := range taskList {
		if i >= maxItems {
			break
		}
		assignedTo := ""
		if t.AssignedTo != nil {
			if u, err := s.userRepo.FindByID(*t.AssignedTo); err == nil {
				assignedTo = u.Name
			}
		}
		taskItems = append(taskItems, dashboard.DraftTaskItem{
			ID:         t.ID,
			Title:      t.Title,
			Priority:   t.Priority,
			DueDate:    t.DueDate,
			AssignedTo: assignedTo,
			CreatedAt:  t.CreatedAt,
		})
	}

	// ── Schedules (submitted = sales submitted schedule for manager confirmation) ─
	scheduleItems := make([]dashboard.DraftScheduleItem, 0, maxItems)
	var scheduleTotal int64
	if s.scheduleRepo != nil {
		schedList, _, _ := s.scheduleRepo.List(&scheduledomain.ListSchedulesRequest{
			Status:        "submitted",
			ScopedUserIDs: teamOnlyIDs,
			Page:          1,
			PerPage:       100,
		})
		scheduleTotal = int64(len(schedList))
		for i, sc := range schedList {
			if i >= maxItems {
				break
			}
			assignedTo := sc.UserID
			if u, err := s.userRepo.FindByID(sc.UserID); err == nil {
				assignedTo = u.Name
			}
			scheduleItems = append(scheduleItems, dashboard.DraftScheduleItem{
				ID:          sc.ID,
				Title:       sc.Title,
				ScheduledAt: sc.ScheduledAt,
				AssignedTo:  assignedTo,
				CreatedAt:   sc.CreatedAt,
			})
		}
	}

	// ── Leads (new = just created / submitted to manager for qualification) ─────
	leadList, _, _ := s.leadRepo.List(&leaddomain.ListLeadsRequest{
		Status:        "new",
		ScopedUserIDs: teamOnlyIDs,
		Page:          1,
		PerPage:       100,
	})
	leadItems := make([]dashboard.DraftLeadItem, 0, maxItems)
	for i, l := range leadList {
		if i >= maxItems {
			break
		}
		assignedTo := ""
		if l.AssignedTo != nil {
			if u, err := s.userRepo.FindByID(*l.AssignedTo); err == nil {
				assignedTo = u.Name
			}
		}
		leadItems = append(leadItems, dashboard.DraftLeadItem{
			ID:         l.ID,
			Name:       strings.TrimSpace(l.FirstName + " " + l.LastName),
			Company:    l.CompanyName,
			Status:     l.LeadStatus,
			AssignedTo: assignedTo,
			CreatedAt:  l.CreatedAt,
		})
	}

	// ── Pipeline / Deals (open deals from team = needing manager oversight) ──────
	dealList, _, _ := s.dealRepo.List(&pipelinedomain.ListDealsRequest{
		Status:        "open",
		ScopedUserIDs: teamOnlyIDs,
		Page:          1,
		PerPage:       100,
	})
	pipelineItems := make([]dashboard.DraftPipelineItem, 0, maxItems)
	for i, d := range dealList {
		if i >= maxItems {
			break
		}
		assignedTo := ""
		if d.AssignedTo != nil {
			if u, err := s.userRepo.FindByID(*d.AssignedTo); err == nil {
				assignedTo = u.Name
			}
		}
		stageName := ""
		if d.Stage != nil {
			stageName = d.Stage.Name
		}
		pipelineItems = append(pipelineItems, dashboard.DraftPipelineItem{
			ID:         d.ID,
			Name:       d.Title,
			Stage:      stageName,
			Value:      d.Value,
			AssignedTo: assignedTo,
			CreatedAt:  d.CreatedAt,
		})
	}

	total := int64(len(vrList)) + int64(len(taskList)) + scheduleTotal + int64(len(leadList)) + int64(len(dealList))

	return &dashboard.SalesManagerTeamDraftApprovalsResponse{
		Total: total,
		Visits: dashboard.SalesManagerDraftVisitsData{
			Total: int64(len(vrList)),
			Items: vrItems,
		},
		Tasks: dashboard.SalesManagerDraftTasksData{
			Total: int64(len(taskList)),
			Items: taskItems,
		},
		Schedules: dashboard.SalesManagerDraftSchedulesData{
			Total: scheduleTotal,
			Items: scheduleItems,
		},
		Leads: dashboard.SalesManagerDraftLeadsData{
			Total: int64(len(leadList)),
			Items: leadItems,
		},
		Pipeline: dashboard.SalesManagerDraftPipelineData{
			Total: int64(len(dealList)),
			Items: pipelineItems,
		},
	}, nil
}

// ============================================================================
// Sales Dashboard Methods
// ============================================================================

// GetSalesTodayTasks returns today's tasks for sales user (requires userID from context)
func (s *Service) GetSalesTodayTasks(userID string) (*dashboard.SalesTodayTasksResponse, error) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// OPTIMIZED: Use reasonable limit instead of loading all tasks
	// For today's tasks, we only need tasks for today
	tasks, _, err := s.taskRepo.List(&taskdomain.ListTasksRequest{
		AssignedTo: userID,
		Page:       1,
		PerPage:    200, // Reasonable limit for today's tasks
	})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	var total, completed, pending, overdue int64
	taskItems := make([]dashboard.SalesTaskItem, 0)

	for _, task := range tasks {
		// Filter tasks for today
		isToday := false
		if task.DueDate != nil {
			taskDate := time.Date(task.DueDate.Year(), task.DueDate.Month(), task.DueDate.Day(), 0, 0, 0, 0, task.DueDate.Location())
			if taskDate.Equal(todayStart) {
				isToday = true
			}
		}

		if isToday || task.DueDate == nil {
			total++
			status := task.Status
			if task.DueDate != nil && task.DueDate.Before(now) && task.Status != "completed" {
				status = "overdue"
				overdue++
			} else if task.Status == "completed" {
				completed++
			} else {
				pending++
			}

			dueTime := ""
			if task.DueDate != nil {
				dueTime = task.DueDate.Format("15:04")
			}

			relatedTo := (*dashboard.TaskRelatedTo)(nil)
			if task.AccountID != nil && *task.AccountID != "" {
				account, err := s.accountRepo.FindByID(*task.AccountID)
				if err == nil {
					relatedTo = &dashboard.TaskRelatedTo{
						Type: "account",
						ID:   account.ID,
						Name: account.Name,
					}
				}
			} else if task.DealID != nil && *task.DealID != "" {
				deal, err := s.dealRepo.FindByID(*task.DealID)
				if err == nil {
					relatedTo = &dashboard.TaskRelatedTo{
						Type: "deal",
						ID:   deal.ID,
						Name: deal.Title,
					}
				}
			} else if task.ContactID != nil && *task.ContactID != "" {
				// Contact reference would need contact repo
				relatedTo = &dashboard.TaskRelatedTo{
					Type: "lead",
					ID:   *task.ContactID,
					Name: "Contact",
				}
			}

			taskItems = append(taskItems, dashboard.SalesTaskItem{
				ID:        task.ID,
				Title:     task.Title,
				DueDate:   *task.DueDate,
				DueTime:   &dueTime,
				Status:    status,
				Priority:  task.Priority,
				RelatedTo: relatedTo,
			})
		}
	}

	response := &dashboard.SalesTodayTasksResponse{
		Total:     total,
		Completed: completed,
		Pending:   pending,
		Overdue:   overdue,
		Tasks:     taskItems,
	}

	return response, nil
}

// GetSalesAssignedLeads returns assigned leads for sales user
func (s *Service) GetSalesAssignedLeads(userID string, req *dashboard.DashboardRequest) (*dashboard.SalesAssignedLeadsResponse, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	// OPTIMIZED: Use reasonable limit instead of loading all leads
	// For assigned leads, we only need a reasonable number for display
	leads, _, err := s.leadRepo.List(&leaddomain.ListLeadsRequest{
		AssignedTo: userID,
		Status:     req.Status,
		Page:       1,
		PerPage:    200, // Reasonable limit for assigned leads
	})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	var total, newCount, contactedCount, qualifiedCount, convertedCount int64
	leadItems := make([]dashboard.SalesLeadItem, 0)

	for _, lead := range leads {
		total++
		switch lead.LeadStatus {
		case "new":
			newCount++
		case "contacted", "interested":
			contactedCount++
		case "qualified", "proposal_sent":
			qualifiedCount++
		case "converted":
			convertedCount++
		}

		if len(leadItems) < limit {
			lastContact := lead.UpdatedAt
			leadItems = append(leadItems, dashboard.SalesLeadItem{
				ID:           lead.ID,
				Name:         fmt.Sprintf("%s %s", lead.FirstName, lead.LastName),
				Company:      lead.CompanyName,
				Status:       lead.LeadStatus,
				AssignedDate: lead.CreatedAt,
				LastContact:  &lastContact,
			})
		}
	}

	response := &dashboard.SalesAssignedLeadsResponse{
		Total:     total,
		New:       newCount,
		Contacted: contactedCount,
		Qualified: qualifiedCount,
		Converted: convertedCount,
		Leads:     leadItems,
	}

	return response, nil
}

// GetSalesUpcomingVisits returns upcoming visits for sales user
func (s *Service) GetSalesUpcomingVisits(userID string, req *dashboard.DashboardRequest) (*dashboard.SalesUpcomingVisitsResponse, error) {
	days := req.Days
	if days <= 0 {
		days = 7
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	now := time.Now()
	startDate := now
	endDate := now.AddDate(0, 0, days)

	// OPTIMIZED: Use reasonable limit instead of loading all visit reports
	// For upcoming visits, we only need visits in the date range
	visitReports, _, err := s.visitReportRepo.List(&visit_report.ListVisitReportsRequest{
		SalesRepID: userID,
		StartDate:  startDate.Format("2006-01-02"),
		EndDate:    endDate.Format("2006-01-02"),
		Page:       1,
		PerPage:    200, // Reasonable limit for upcoming visits
	})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	var total, todayCount, thisWeekCount, nextWeekCount int64
	visitItems := make([]dashboard.SalesVisitItem, 0)

	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekEnd := today.AddDate(0, 0, 7)
	nextWeekStart := weekEnd
	nextWeekEnd := nextWeekStart.AddDate(0, 0, 7)

	for _, vr := range visitReports {
		if vr.Status == "scheduled" || vr.Status == "draft" {
			total++
			visitDate := time.Date(vr.VisitDate.Year(), vr.VisitDate.Month(), vr.VisitDate.Day(), 0, 0, 0, 0, vr.VisitDate.Location())

			if visitDate.Equal(today) {
				todayCount++
			}
			if visitDate.After(today) && visitDate.Before(weekEnd) {
				thisWeekCount++
			}
			if visitDate.After(nextWeekStart) && visitDate.Before(nextWeekEnd) {
				nextWeekCount++
			}

			if len(visitItems) < limit {
				accountName := "Unknown Account"
				if vr.AccountID != nil && *vr.AccountID != "" {
					account, err := s.accountRepo.FindByID(*vr.AccountID)
					if err == nil {
						accountName = account.Name
					}
				}

				scheduledTime := vr.VisitDate.Format("15:04")
				visitItems = append(visitItems, dashboard.SalesVisitItem{
					ID:            vr.ID,
					AccountName:   accountName,
					ScheduledDate: vr.VisitDate,
					ScheduledTime: &scheduledTime,
					Purpose:       vr.Purpose,
					Status:        vr.Status,
				})
			}
		}
	}

	response := &dashboard.SalesUpcomingVisitsResponse{
		Total:    total,
		Today:    todayCount,
		ThisWeek: thisWeekCount,
		NextWeek: nextWeekCount,
		Visits:   visitItems,
	}

	return response, nil
}

// GetSalesReminders returns reminders for sales user
func (s *Service) GetSalesReminders(userID string, req *dashboard.DashboardRequest) (*dashboard.SalesRemindersResponse, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	// Get notifications for user (reminders are stored as notifications)
	// We need notification repo - for now, return placeholder
	// TODO: Add notification repository to service

	response := &dashboard.SalesRemindersResponse{
		Total:     0,
		Unread:    0,
		Reminders: []dashboard.ReminderItem{},
	}

	return response, nil
}

// ============================================================================
// Analyst Dashboard Methods
// ============================================================================

// GetAnalystRevenueTrend returns revenue trend over time
func (s *Service) GetAnalystRevenueTrend(req *dashboard.DashboardRequest) (*dashboard.AnalystRevenueTrendResponse, error) {
	// Parse period
	var start, end time.Time
	if req.StartDate != "" && req.EndDate != "" {
		var err error
		start, err = time.Parse("2006-01-02", req.StartDate)
		if err == nil {
			end, err = time.Parse("2006-01-02", req.EndDate)
			if err == nil {
				end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, end.Location())
			}
		}
	} else if req.Period != "" {
		start, end = parsePeriod(req.Period)
	} else {
		start, end = parsePeriod("month")
	}

	periodStr := req.Period
	if periodStr == "" {
		periodStr = "month"
	}

	_, salesUserIDs, err := s.loadSalesUsers(req.ScopedUserIDs)
	if err != nil {
		return nil, err
	}

	wonDeals, err := s.listScopedDeals(&pipelinedomain.ListDealsRequest{
		Status:        "won",
		ScopedUserIDs: salesUserIDs,
	})
	if err != nil {
		return nil, err
	}

	// Group by date for trend
	revenueByDate := make(map[string]struct {
		revenue int64
		deals   int64
	})
	var totalRevenue int64

	for _, deal := range wonDeals {
		if deal.ActualCloseDate != nil {
			closeDate := *deal.ActualCloseDate
			if (closeDate.After(start) && closeDate.Before(end)) || closeDate.Equal(start) || closeDate.Equal(end) {
				dateKey := closeDate.Format("2006-01-02")
				entry := revenueByDate[dateKey]
				entry.revenue += deal.Value
				entry.deals++
				revenueByDate[dateKey] = entry
				totalRevenue += deal.Value
			}
		}
	}

	// Build trend array
	trend := make([]dashboard.RevenueTrendEntry, 0)
	for dateKey, entry := range revenueByDate {
		trend = append(trend, dashboard.RevenueTrendEntry{
			Date:    dateKey,
			Revenue: entry.revenue,
			Deals:   entry.deals,
		})
	}

	// Calculate growth percent (compare with previous period)
	prevStart := start.AddDate(0, 0, -int(end.Sub(start).Hours()/24))
	prevEnd := start.Add(-time.Second)
	var prevRevenue int64
	for _, deal := range wonDeals {
		if deal.ActualCloseDate != nil {
			closeDate := *deal.ActualCloseDate
			if (closeDate.After(prevStart) && closeDate.Before(prevEnd)) || closeDate.Equal(prevStart) || closeDate.Equal(prevEnd) {
				prevRevenue += deal.Value
			}
		}
	}

	growthPercent := float64(0)
	if prevRevenue > 0 {
		growthPercent = float64(totalRevenue-prevRevenue) / float64(prevRevenue) * 100
	}

	daysDiff := int(end.Sub(start).Hours() / 24)
	if daysDiff == 0 {
		daysDiff = 1
	}
	averageDaily := totalRevenue / int64(daysDiff)

	response := &dashboard.AnalystRevenueTrendResponse{
		Period:        periodStr,
		TotalRevenue:  totalRevenue,
		Trend:         trend,
		GrowthPercent: growthPercent,
		AverageDaily:  averageDaily,
	}

	return response, nil
}

// GetAnalystConversionRate returns conversion rate statistics
func (s *Service) GetAnalystConversionRate(req *dashboard.DashboardRequest) (*dashboard.AnalystConversionRateResponse, error) {
	period := req.Period
	if period == "" {
		period = "month"
	}

	var start, end time.Time
	if req.StartDate != "" && req.EndDate != "" {
		var err error
		start, err = time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, err
		}
		end, err = time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, err
		}
		end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, end.Location())
	} else if req.Period != "" {
		start, end = parsePeriod(req.Period)
	} else {
		start, end = parsePeriod("month")
	}

	_, salesUserIDs, err := s.loadSalesUsers(req.ScopedUserIDs)
	if err != nil {
		return nil, err
	}

	allLeads, err := s.listScopedLeads(&leaddomain.ListLeadsRequest{
		ScopedUserIDs: salesUserIDs,
	})
	if err != nil {
		return nil, err
	}

	var totalLeads, convertedLeads int64
	sourceCounts := make(map[string]struct {
		leads     int64
		converted int64
	})

	for _, lead := range allLeads {
		if lead.CreatedAt.Before(start) || lead.CreatedAt.After(end) {
			continue
		}
		source := lead.LeadSource
		if source == "" {
			source = "other"
		}
		sourceCount := sourceCounts[source]
		sourceCount.leads++
		totalLeads++
		if lead.LeadStatus == "converted" {
			sourceCount.converted++
			convertedLeads++
		}
		sourceCounts[source] = sourceCount
	}

	trendCounts := make(map[string]struct {
		total     int64
		converted int64
	})
	for _, lead := range allLeads {
		if lead.CreatedAt.Before(start) || lead.CreatedAt.After(end) {
			continue
		}
		dayKey := lead.CreatedAt.Format("2006-01-02")
		entry := trendCounts[dayKey]
		entry.total++
		if lead.LeadStatus == "converted" {
			entry.converted++
		}
		trendCounts[dayKey] = entry
	}

	conversionRate := float64(0)
	if totalLeads > 0 {
		conversionRate = float64(convertedLeads) / float64(totalLeads) * 100
	}

	bySource := make([]dashboard.ConversionRateBySource, 0, len(sourceCounts))
	for source, counts := range sourceCounts {
		sourceConvRate := float64(0)
		if counts.leads > 0 {
			sourceConvRate = float64(counts.converted) / float64(counts.leads) * 100
		}
		bySource = append(bySource, dashboard.ConversionRateBySource{
			Source:         source,
			Leads:          counts.leads,
			Converted:      counts.converted,
			ConversionRate: sourceConvRate,
		})
	}

	trend := make([]dashboard.ConversionRateTrendEntry, 0)
	current := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	limitDate := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())
	for current.Before(limitDate) || current.Equal(limitDate) {
		dayKey := current.Format("2006-01-02")
		counts := trendCounts[dayKey]
		dayConversionRate := float64(0)
		if counts.total > 0 {
			dayConversionRate = float64(counts.converted) / float64(counts.total) * 100
		}
		trend = append(trend, dashboard.ConversionRateTrendEntry{
			Date:           dayKey,
			ConversionRate: dayConversionRate,
		})
		current = current.AddDate(0, 0, 1)
	}

	response := &dashboard.AnalystConversionRateResponse{
		Period:         period,
		TotalLeads:     totalLeads,
		ConvertedLeads: convertedLeads,
		ConversionRate: conversionRate,
		BySource:       bySource,
		Trend:          trend,
	}

	return response, nil
}

// GetAnalystSalesVelocity returns sales velocity metrics
func (s *Service) GetAnalystSalesVelocity(req *dashboard.DashboardRequest) (*dashboard.AnalystSalesVelocityResponse, error) {
	period := req.Period
	if period == "" {
		period = "month"
	}

	var start, end time.Time
	if req.StartDate != "" && req.EndDate != "" {
		var err error
		start, err = time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, err
		}
		end, err = time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, err
		}
		end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, end.Location())
	} else if req.Period != "" {
		start, end = parsePeriod(req.Period)
	} else {
		start, end = parsePeriod("month")
	}

	_, salesUserIDs, err := s.loadSalesUsers(req.ScopedUserIDs)
	if err != nil {
		return nil, err
	}

	wonDeals, err := s.listScopedDeals(&pipelinedomain.ListDealsRequest{
		Status:        "won",
		ScopedUserIDs: salesUserIDs,
	})
	if err != nil {
		return nil, err
	}

	filteredWonDeals := make([]pipelinedomain.Deal, 0, len(wonDeals))
	for _, deal := range wonDeals {
		if deal.ActualCloseDate == nil {
			continue
		}
		if deal.ActualCloseDate.Before(start) || deal.ActualCloseDate.After(end) {
			continue
		}
		filteredWonDeals = append(filteredWonDeals, deal)
	}

	var totalDays int64
	var totalValue int64
	dealCount := int64(len(filteredWonDeals))

	for _, deal := range filteredWonDeals {
		if deal.CreatedAt.Before(*deal.ActualCloseDate) {
			days := int64(deal.ActualCloseDate.Sub(deal.CreatedAt).Hours() / 24)
			totalDays += days
			totalValue += deal.Value
		}
	}

	averageSalesCycleDays := 0
	if dealCount > 0 {
		averageSalesCycleDays = int(totalDays / dealCount)
	}

	averageDealValue := int64(0)
	if dealCount > 0 {
		averageDealValue = totalValue / dealCount
	}

	salesVelocity := float64(0)
	if averageSalesCycleDays > 0 {
		salesVelocity = float64(dealCount*averageDealValue) / float64(averageSalesCycleDays)
	}

	byStage := make([]dashboard.SalesVelocityByStage, 0)
	stageDaysMap := make(map[string][]int64)
	stageDealCount := make(map[string]int64)

	for _, deal := range filteredWonDeals {
		if deal.StageID != "" && deal.CreatedAt.Before(*deal.ActualCloseDate) {
			days := int64(deal.ActualCloseDate.Sub(deal.CreatedAt).Hours() / 24)
			stageDaysMap[deal.StageID] = append(stageDaysMap[deal.StageID], days)
			stageDealCount[deal.StageID]++
		}
	}

	for stageID, daysList := range stageDaysMap {
		if len(daysList) > 0 {
			var totalStageDays int64
			for _, days := range daysList {
				totalStageDays += days
			}
			averageDays := int(totalStageDays / int64(len(daysList)))

			stage, err := s.pipelineRepo.FindStageByID(stageID)
			stageName := "Unknown"
			if err == nil && stage != nil {
				stageName = stage.Name
			}

			byStage = append(byStage, dashboard.SalesVelocityByStage{
				Stage:       stageName,
				AverageDays: averageDays,
			})
		}
	}

	response := &dashboard.AnalystSalesVelocityResponse{
		Period:                period,
		AverageSalesCycleDays: averageSalesCycleDays,
		AverageDealValue:      averageDealValue,
		SalesVelocity:         salesVelocity,
		ByStage:               byStage,
	}

	return response, nil
}

// GetAnalystAIInsights returns AI-generated insights (placeholder)
func (s *Service) GetAnalystAIInsights(req *dashboard.DashboardRequest) (*dashboard.AnalystAIInsightsResponse, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 5
	}

	// TODO: Implement when AI insights table is available
	response := &dashboard.AnalystAIInsightsResponse{
		Insights:      []dashboard.AIInsightItem{},
		TotalInsights: 0,
	}

	return response, nil
}

// ============================================================================
// Mobile Dashboard Service Methods
// ============================================================================

// GetMobileOverview returns simplified dashboard overview for mobile sales user
func (s *Service) GetMobileOverview(userID string, req *dashboard.MobileDashboardRequest) (*dashboard.MobileDashboardOverviewResponse, error) {
	// Parse period
	var start, end time.Time
	if req.StartDate != "" && req.EndDate != "" {
		var err error
		start, err = time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, err
		}
		end, err = time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, err
		}
		end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, end.Location())
	} else {
		period := req.Period
		if period == "" {
			period = "today"
		}
		start, end = parsePeriod(period)
	}

	// Get target summary
	targetSummary := dashboard.MobileTargetSummary{
		TargetAmount:            0,
		TargetAmountFormatted:   "Rp 0",
		AchievedAmount:          0,
		AchievedAmountFormatted: "Rp 0",
		ProgressPercent:         0,
		Period:                  fmt.Sprintf("%d-%02d", start.Year(), int(start.Month())),
	}

	if userID != "" && s.monthlyTargetRepo != nil {
		targetYear := start.Year()
		targetMonth := int(start.Month())

		effectiveTarget, err := s.monthlyTargetRepo.GetUserEffectiveTarget(userID, targetYear, targetMonth)
		if err == nil && effectiveTarget != nil {
			// Get achieved amount from won deals
			_, achievedAmount, err := s.dealRepo.GetWonDealsValueInPeriodByUser(userID, start, end)
			if err != nil {
				achievedAmount = 0
			}

			var progressPercent float64
			if effectiveTarget.TargetAmount > 0 {
				progressPercent = float64(achievedAmount) / float64(effectiveTarget.TargetAmount) * 100
			}

			targetSummary = dashboard.MobileTargetSummary{
				TargetAmount:            effectiveTarget.TargetAmount,
				TargetAmountFormatted:   formatCurrency(effectiveTarget.TargetAmount),
				AchievedAmount:          achievedAmount,
				AchievedAmountFormatted: formatCurrency(achievedAmount),
				ProgressPercent:         progressPercent,
				Period:                  fmt.Sprintf("%d-%02d", targetYear, targetMonth),
			}

			// Get brick name if available
			if effectiveTarget.BrickID != nil && s.accountRepo != nil {
				// Note: Brick name would need to be fetched from brick repository
				// For now, we'll leave it as nil
			}
		}
	}

	response := &dashboard.MobileDashboardOverviewResponse{
		Target: targetSummary,
	}

	return response, nil
}

// GetMobileVisits returns visits list for mobile dashboard (max 5 items)
func (s *Service) GetMobileVisits(userID string, req *dashboard.MobileDashboardRequest) (*dashboard.MobileVisitsListResponse, error) {
	// Set limit (max 5 for dashboard)
	limit := req.Limit
	if limit <= 0 {
		limit = 5
	}
	if limit > 5 {
		limit = 5
	}

	// Build visit report request
	visitReq := &visit_report.ListVisitReportsRequest{
		SalesRepID: userID,
		Page:       1,
		PerPage:    limit,
	}

	// Apply status filter
	// Status mapping: all, planned (draft), completed (approved), cancelled
	status := req.Status
	if status == "" {
		status = "all"
	}

	if status != "all" {
		switch status {
		case "planned":
			// Planned: draft status (visit reports that haven't been submitted yet)
			visitReq.Status = "draft"
		case "completed":
			// Completed: approved status
			visitReq.Status = "approved"
		case "cancelled":
			// Cancelled: cancelled status
			visitReq.Status = "cancelled"
		default:
			// For other statuses, use as-is
			visitReq.Status = status
		}
	}

	// Apply date filter
	if req.Date != "" {
		visitReq.StartDate = req.Date
		visitReq.EndDate = req.Date
	}

	// Get visits
	visits, total, err := s.visitReportRepo.List(visitReq)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Collect IDs for batch loading
	accountIDs := make([]string, 0)
	accountIDSet := make(map[string]bool)
	contactIDs := make([]string, 0)
	contactIDSet := make(map[string]bool)
	dealIDs := make([]string, 0)
	dealIDSet := make(map[string]bool)
	leadIDs := make([]string, 0)
	leadIDSet := make(map[string]bool)

	for _, vr := range visits {
		if vr.AccountID != nil && *vr.AccountID != "" && !accountIDSet[*vr.AccountID] {
			accountIDs = append(accountIDs, *vr.AccountID)
			accountIDSet[*vr.AccountID] = true
		}
		if vr.ContactID != nil && *vr.ContactID != "" && !contactIDSet[*vr.ContactID] {
			contactIDs = append(contactIDs, *vr.ContactID)
			contactIDSet[*vr.ContactID] = true
		}
		if vr.DealID != nil && *vr.DealID != "" && !dealIDSet[*vr.DealID] {
			dealIDs = append(dealIDs, *vr.DealID)
			dealIDSet[*vr.DealID] = true
		}
		if vr.LeadID != nil && *vr.LeadID != "" && !leadIDSet[*vr.LeadID] {
			leadIDs = append(leadIDs, *vr.LeadID)
			leadIDSet[*vr.LeadID] = true
		}
	}

	// Batch load accounts
	accountsMap := make(map[string]*account.Account)
	if len(accountIDs) > 0 && s.accountRepo != nil {
		for _, id := range accountIDs {
			if acc, err := s.accountRepo.FindByID(id); err == nil {
				accountsMap[id] = acc
			}
		}
	}

	// Batch load contacts
	// Note: Contact loading would need contactRepo which is not available in dashboard service
	// Contact names will be left as nil for now

	// Batch load deals
	dealsMap := make(map[string]interface{})
	if len(dealIDs) > 0 && s.dealRepo != nil {
		for _, id := range dealIDs {
			if deal, err := s.dealRepo.FindByID(id); err == nil {
				dealsMap[id] = deal
			}
		}
	}

	// Batch load leads
	leadsMap := make(map[string]interface{})
	if len(leadIDs) > 0 && s.leadRepo != nil {
		for _, id := range leadIDs {
			if lead, err := s.leadRepo.FindByID(id); err == nil {
				leadsMap[id] = lead
			}
		}
	}

	// Convert to mobile response format
	mobileVisits := make([]dashboard.MobileVisitResponse, 0, len(visits))
	for _, vr := range visits {
		// Determine type (priority: lead > deal > account)
		visitType := "account"
		if vr.LeadID != nil && *vr.LeadID != "" {
			visitType = "lead"
		} else if vr.DealID != nil && *vr.DealID != "" {
			visitType = "deal"
		}

		// Get account info from batch-loaded map
		var accountID *string
		var accountName *string
		var accountAddress *string
		if vr.AccountID != nil && *vr.AccountID != "" {
			if acc, ok := accountsMap[*vr.AccountID]; ok {
				accountID = vr.AccountID
				name := acc.Name
				accountName = &name
				if acc.Address != "" {
					accountAddress = &acc.Address
				}
			}
		}

		// Get contact info (Note: Contact loading would need contactRepo)
		var contactID *string
		var contactName *string
		if vr.ContactID != nil && *vr.ContactID != "" {
			contactID = vr.ContactID
			// Contact name loading can be added if contactRepo is available
		}

		// Get deal info
		var dealID *string
		var dealTitle *string
		if vr.DealID != nil && *vr.DealID != "" {
			if deal, ok := dealsMap[*vr.DealID]; ok {
				dealID = vr.DealID
				// Type assert to get deal title
				if d, ok := deal.(*pipelinedomain.Deal); ok {
					title := d.Title
					dealTitle = &title
				}
			}
		}

		// Get lead info
		var leadID *string
		var leadName *string
		if vr.LeadID != nil && *vr.LeadID != "" {
			if lead, ok := leadsMap[*vr.LeadID]; ok {
				leadID = vr.LeadID
				// Type assert to get lead name
				if l, ok := lead.(*leaddomain.Lead); ok {
					name := fmt.Sprintf("%s %s", l.FirstName, l.LastName)
					if l.LastName == "" {
						name = l.FirstName
					}
					leadName = &name
				}
			}
		}

		// Parse location from JSON
		var checkInLocation *dashboard.MobileVisitLocation
		var checkOutLocation *dashboard.MobileVisitLocation

		if vr.CheckInLocation != nil {
			var loc visit_report.Location
			if err := json.Unmarshal(vr.CheckInLocation, &loc); err == nil {
				checkInLocation = &dashboard.MobileVisitLocation{
					Latitude:  loc.Latitude,
					Longitude: loc.Longitude,
					Address:   &loc.Address,
				}
			}
		}

		if vr.CheckOutLocation != nil {
			var loc visit_report.Location
			if err := json.Unmarshal(vr.CheckOutLocation, &loc); err == nil {
				checkOutLocation = &dashboard.MobileVisitLocation{
					Latitude:  loc.Latitude,
					Longitude: loc.Longitude,
					Address:   &loc.Address,
				}
			}
		}

		// Format visit time
		var visitTime *string
		if vr.CheckInTime != nil {
			formatted := vr.CheckInTime.Format("15:04")
			visitTime = &formatted
		}

		mobileVisit := dashboard.MobileVisitResponse{
			ID:               vr.ID,
			Type:             visitType,
			Purpose:          vr.Purpose,
			AccountID:        accountID,
			AccountName:      accountName,
			AccountAddress:   accountAddress,
			ContactID:        contactID,
			ContactName:      contactName,
			DealID:           dealID,
			DealTitle:        dealTitle,
			LeadID:           leadID,
			LeadName:         leadName,
			VisitDate:        vr.VisitDate.Format("2006-01-02"),
			VisitTime:        visitTime,
			Status:           vr.Status,
			CheckInTime:      vr.CheckInTime,
			CheckInLocation:  checkInLocation,
			CheckOutTime:     vr.CheckOutTime,
			CheckOutLocation: checkOutLocation,
			CreatedAt:        vr.CreatedAt,
			UpdatedAt:        vr.UpdatedAt,
		}

		mobileVisits = append(mobileVisits, mobileVisit)
	}

	hasMore := int64(len(mobileVisits)) < total

	response := &dashboard.MobileVisitsListResponse{
		Visits:  mobileVisits,
		Total:   int(total),
		HasMore: hasMore,
	}

	return response, nil
}

// GetMobileTasks returns upcoming tasks list for mobile dashboard (max 5 items)
// Only returns tasks that are assigned to the user and have due date >= today
func (s *Service) GetMobileTasks(userID string, req *dashboard.MobileDashboardRequest) (*dashboard.MobileTasksListResponse, error) {
	// Set limit (max 5 for dashboard, default 5)
	limit := req.Limit
	if limit <= 0 {
		limit = 5
	}
	if limit > 5 {
		limit = 5
	}

	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Build task request - only get pending/in_progress tasks assigned to user
	taskReq := &taskdomain.ListTasksRequest{
		AssignedTo: userID,
		Status:     "pending,in_progress", // Only upcoming tasks (not completed/cancelled)
		Page:       1,
		PerPage:    100, // Get more to filter by date
	}

	// Get tasks (we'll filter and count ourselves)
	tasks, _, err := s.taskRepo.List(taskReq)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Filter: only upcoming tasks (due date >= today) or tasks without due date
	filteredTasks := make([]taskdomain.Task, 0)
	for _, t := range tasks {
		// Include tasks without due date or with due date >= today
		if t.DueDate == nil {
			filteredTasks = append(filteredTasks, t)
		} else {
			dueDate := *t.DueDate
			dueDateOnly := time.Date(dueDate.Year(), dueDate.Month(), dueDate.Day(), 0, 0, 0, 0, dueDate.Location())
			if dueDateOnly.After(todayStart) || dueDateOnly.Equal(todayStart) {
				filteredTasks = append(filteredTasks, t)
			}
		}
	}

	// Sort by due date (earliest first)
	for i := 0; i < len(filteredTasks)-1; i++ {
		for j := i + 1; j < len(filteredTasks); j++ {
			t1 := filteredTasks[i]
			t2 := filteredTasks[j]

			// Tasks without due date go to the end
			if t1.DueDate == nil && t2.DueDate != nil {
				filteredTasks[i], filteredTasks[j] = filteredTasks[j], filteredTasks[i]
				continue
			}
			if t1.DueDate != nil && t2.DueDate == nil {
				continue
			}
			if t1.DueDate != nil && t2.DueDate != nil {
				if t1.DueDate.After(*t2.DueDate) {
					filteredTasks[i], filteredTasks[j] = filteredTasks[j], filteredTasks[i]
				}
			}
		}
	}

	// Limit to max 5
	if len(filteredTasks) > limit {
		filteredTasks = filteredTasks[:limit]
	}

	// Convert to mobile response format
	mobileTasks := make([]dashboard.MobileTaskResponse, 0, len(filteredTasks))
	for _, t := range filteredTasks {
		var dueDateStr *string
		var dueTimeStr *string
		if t.DueDate != nil {
			formatted := t.DueDate.Format("2006-01-02")
			dueDateStr = &formatted
			formattedTime := t.DueDate.Format("15:04")
			dueTimeStr = &formattedTime
		}

		var assignedBy *dashboard.MobileTaskAssignee
		if t.AssignedFromUser != nil {
			assignedBy = &dashboard.MobileTaskAssignee{
				ID:   t.AssignedFromUser.ID,
				Name: t.AssignedFromUser.Name,
			}
		}

		isOverdue := false
		if t.DueDate != nil && t.DueDate.Before(todayStart) && t.Status != "completed" {
			isOverdue = true
		}

		description := t.Description
		var descriptionPtr *string
		if description != "" {
			descriptionPtr = &description
		}

		// Get task type (default to "general" if empty)
		taskType := t.Type
		if taskType == "" {
			taskType = "general"
		}

		mobileTask := dashboard.MobileTaskResponse{
			ID:          t.ID,
			Title:       t.Title,
			Description: descriptionPtr,
			Type:        taskType,
			DueDate:     dueDateStr,
			DueTime:     dueTimeStr,
			Priority:    t.Priority,
			Status:      t.Status,
			AssignedBy:  assignedBy,
			CreatedAt:   t.CreatedAt,
			IsOverdue:   isOverdue,
		}

		mobileTasks = append(mobileTasks, mobileTask)
	}

	// Calculate has_more based on filtered tasks count
	// If we have more than limit, there are more tasks available
	hasMore := len(filteredTasks) > limit
	// Total is the count of all upcoming tasks (not just the returned ones)
	totalUpcoming := int64(len(filteredTasks))

	response := &dashboard.MobileTasksListResponse{
		Tasks:   mobileTasks,
		Total:   int(totalUpcoming),
		HasMore: hasMore,
	}

	return response, nil
}
