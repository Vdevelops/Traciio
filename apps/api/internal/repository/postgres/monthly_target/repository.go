package monthly_target

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/monthly_target"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new monthly target repository
func NewRepository(db *gorm.DB) interfaces.MonthlyTargetRepository {
	return &repository{db: db}
}

func (r *repository) FindByID(id string) (*monthly_target.MonthlyTarget, error) {
	var mt monthly_target.MonthlyTarget
	err := r.db.Preload("Group").Preload("User").Preload("Brick").Where("id = ?", id).First(&mt).Error
	if err != nil {
		return nil, err
	}
	return &mt, nil
}

func (r *repository) FindByGroupAndPeriod(groupID string, year int, month int) (*monthly_target.MonthlyTarget, error) {
	var mt monthly_target.MonthlyTarget
	err := r.db.Preload("Group").
		Where("group_id = ? AND year = ? AND month = ?", groupID, year, month).
		First(&mt).Error
	if err != nil {
		return nil, err
	}
	return &mt, nil
}

func (r *repository) FindByUserAndPeriod(userID string, year int, month int) (*monthly_target.MonthlyTarget, error) {
	var mt monthly_target.MonthlyTarget
	err := r.db.Preload("User").
		Where("user_id = ? AND year = ? AND month = ?", userID, year, month).
		First(&mt).Error
	if err != nil {
		return nil, err
	}
	return &mt, nil
}

func (r *repository) FindByBrickAndPeriod(brickID string, year int, month int) (*monthly_target.MonthlyTarget, error) {
	var mt monthly_target.MonthlyTarget
	err := r.db.Preload("Brick").
		Where("brick_id = ? AND year = ? AND month = ?", brickID, year, month).
		First(&mt).Error
	if err != nil {
		return nil, err
	}
	return &mt, nil
}

func (r *repository) List(req *monthly_target.ListMonthlyTargetsRequest) ([]monthly_target.MonthlyTarget, int64, int64, error) {
	var targets []monthly_target.MonthlyTarget
	var total int64
	var totalAmount int64

	query := r.db.Model(&monthly_target.MonthlyTarget{}).Preload("Group").Preload("User").Preload("Brick")

	// Apply filters
	if req.GroupID != nil && *req.GroupID != "" {
		query = query.Where("group_id = ?", *req.GroupID)
	}

	if req.UserID != nil && *req.UserID != "" {
		query = query.Where("user_id = ?", *req.UserID)
	}

	if req.BrickID != nil && *req.BrickID != "" {
		query = query.Where("brick_id = ?", *req.BrickID)
	}

	if req.Year != nil {
		query = query.Where("year = ?", *req.Year)
	}

	if req.Month != nil {
		query = query.Where("month = ?", *req.Month)
	}

	if req.Scope != "" && req.Scope != "all" {
		switch req.Scope {
		case "user":
			query = query.Where("user_id IS NOT NULL")
		case "group":
			query = query.Where("group_id IS NOT NULL")
		case "brick":
			query = query.Where("brick_id IS NOT NULL")
		}
	}

	// Apply RBAC scope filtering
	if len(req.ScopedUserIDs) > 0 {
		query = query.Where("user_id IN ?", req.ScopedUserIDs)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	// Calculate total amount (sum of target_amount)
	// We use a clone or re-apply logic if needed, but here simple Select works on the same query object structure *before* pagination
	// GORM's Count/Scan might modify the query object, so it's safer to use the query object for Sum
	// However, GORM usually creates a new scope. To be safe, we can use the existing query object for Sum
	if err := query.Select("COALESCE(SUM(target_amount), 0)").Scan(&totalAmount).Error; err != nil {
		return nil, 0, 0, err
	}

	// Reset Select for Find
	query = query.Select("*")

	// Apply pagination
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	if err := query.Offset(offset).Limit(perPage).Order("year desc, month desc").Find(&targets).Error; err != nil {
		return nil, 0, 0, err
	}

	return targets, total, totalAmount, nil
}

func (r *repository) Create(target *monthly_target.MonthlyTarget) error {
	return r.db.Create(target).Error
}

func (r *repository) Update(target *monthly_target.MonthlyTarget) error {
	return r.db.Save(target).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Delete(&monthly_target.MonthlyTarget{}, "id = ?", id).Error
}

func (r *repository) GetUserEffectiveTarget(userID string, year int, month int) (*monthly_target.MonthlyTarget, error) {
	// 1. Try to get user specific target
	target, err := r.FindByUserAndPeriod(userID, year, month)
	if err == nil {
		return target, nil
	}

	// 2. If not found, get user's group and check for group target
	var user struct {
		ID      string
		GroupID *string
	}
	
	if err := r.db.Table("users").Select("id, group_id").Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	if user.GroupID == nil {
		// User has no group and no specific target
		return nil, nil
	}

	return r.FindByGroupAndPeriod(*user.GroupID, year, month)
}

func (r *repository) BatchGetUserEffectiveTargets(userIDs []string, year int, month int) (map[string]*monthly_target.MonthlyTarget, error) {
	if len(userIDs) == 0 {
		return make(map[string]*monthly_target.MonthlyTarget), nil
	}

	result := make(map[string]*monthly_target.MonthlyTarget)
	
	// 1. Get user specific targets
	var userTargets []monthly_target.MonthlyTarget
	err := r.db.Preload("User").
		Where("user_id IN ? AND year = ? AND month = ?", userIDs, year, month).
		Find(&userTargets).Error
	
	if err != nil {
		return nil, err
	}

	// Map user targets
	usersWithTargets := make(map[string]bool)
	for i := range userTargets {
		userID := userTargets[i].UserID
		if userID != nil {
			result[*userID] = &userTargets[i]
			usersWithTargets[*userID] = true
		}
	}

	// Identify users needing group targets
	var usersNeedingGroupTargets []string
	for _, userID := range userIDs {
		if !usersWithTargets[userID] {
			usersNeedingGroupTargets = append(usersNeedingGroupTargets, userID)
		}
	}

	if len(usersNeedingGroupTargets) == 0 {
		return result, nil
	}

	// Get group IDs for users without user-specific targets
	type UserGroup struct {
		UserID  string
		GroupID *string
	}
	var userGroups []UserGroup
	err = r.db.Table("users").
		Select("id as user_id, group_id").
		Where("id IN ?", usersNeedingGroupTargets).
		Find(&userGroups).Error
	
	if err != nil {
		return nil, err
	}

	// Collect unique group IDs
	groupIDSet := make(map[string]bool)
	userToGroupMap := make(map[string]*string)
	for _, ug := range userGroups {
		if ug.GroupID != nil {
			groupIDSet[*ug.GroupID] = true
			userToGroupMap[ug.UserID] = ug.GroupID
		}
	}

	if len(groupIDSet) == 0 {
		return result, nil
	}

	// Get all group targets in one query
	groupIDs := make([]string, 0, len(groupIDSet))
	for gid := range groupIDSet {
		groupIDs = append(groupIDs, gid)
	}

	var groupTargets []monthly_target.MonthlyTarget
	err = r.db.Preload("Group").
		Where("group_id IN ? AND year = ? AND month = ?", groupIDs, year, month).
		Find(&groupTargets).Error
	
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Map group targets by group ID
	groupTargetMap := make(map[string]*monthly_target.MonthlyTarget)
	for i := range groupTargets {
		if groupTargets[i].GroupID != nil {
			groupTargetMap[*groupTargets[i].GroupID] = &groupTargets[i]
		}
	}

	// Assign group targets to users
	for _, userID := range usersNeedingGroupTargets {
		if groupID, exists := userToGroupMap[userID]; exists && groupID != nil {
			if groupTarget, exists := groupTargetMap[*groupID]; exists {
				result[userID] = groupTarget
			}
		}
	}

	return result, nil
}

// GetProratedTargetForPeriod calculates the time-weighted target for a specific period
func (r *repository) GetProratedTargetForPeriod(userID string, startDateStr, endDateStr string) (float64, error) {
	start, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return 0, err
	}
	end, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return 0, err
	}

	// If start date is after end date, return 0
	if start.After(end) {
		return 0, nil
	}

	var totalTarget float64
	current := start
	// Normalize to start of month for iteration
	monthStart := time.Date(current.Year(), current.Month(), 1, 0, 0, 0, 0, current.Location())
	
	for !monthStart.After(end) {
		year := monthStart.Year()
		month := int(monthStart.Month())
		
		target, err := r.GetUserEffectiveTarget(userID, year, month)
		if err == nil && target != nil {
			// Calculate active days in this month
			// Month start is max(periodStart, monthBegin)
			periodStartInMonth := monthStart
			if start.After(monthStart) {
				periodStartInMonth = start
			}
			
			// Month end is min(periodEnd, monthLastDay)
			nextMonth := monthStart.AddDate(0, 1, 0)
			monthEnd := nextMonth.Add(-time.Nanosecond)
			
			periodEndInMonth := monthEnd
			// Normalize end date time for comparison (end of day)
			endEndOfDay := time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, end.Location())
			if endEndOfDay.Before(monthEnd) {
				periodEndInMonth = endEndOfDay
			}

			// Calculate days (inclusive)
			// Truncate to days
			d1 := time.Date(periodStartInMonth.Year(), periodStartInMonth.Month(), periodStartInMonth.Day(), 0, 0, 0, 0, time.UTC)
			d2 := time.Date(periodEndInMonth.Year(), periodEndInMonth.Month(), periodEndInMonth.Day(), 0, 0, 0, 0, time.UTC)
			
			activeDays := int(d2.Sub(d1).Hours() / 24) + 1
			if activeDays < 0 {
				activeDays = 0
			}
			
			daysInMonth := daysIn(monthStart.Month(), monthStart.Year())
			
			// Prorate: (Target / DaysInMonth) * ActiveDays
			if daysInMonth > 0 {
				monthlyProrated := (float64(target.TargetAmount) / float64(daysInMonth)) * float64(activeDays)
				totalTarget += monthlyProrated
			}
		}
		
		// Move to next month
		monthStart = monthStart.AddDate(0, 1, 0)
	}

	return totalTarget, nil
}

// BatchGetProratedTargetsForPeriod calculates time-weighted targets for multiple users
func (r *repository) BatchGetProratedTargetsForPeriod(userIDs []string, startDateStr, endDateStr string) (map[string]float64, error) {
	start, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return nil, err
	}
	end, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return nil, err
	}

	result := make(map[string]float64)
	for _, uid := range userIDs {
		result[uid] = 0
	}

	if start.After(end) {
		return result, nil
	}

	// Identify all months involved
	monthStart := time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, start.Location())
	
	for !monthStart.After(end) {
		year := monthStart.Year()
		month := int(monthStart.Month())
		
		// Batch fetch targets for this month
		targets, err := r.BatchGetUserEffectiveTargets(userIDs, year, month)
		if err != nil {
			return nil, err
		}

		// Calculate proration factor for this month
		periodStartInMonth := monthStart
		if start.After(monthStart) {
			periodStartInMonth = start
		}
		
		nextMonth := monthStart.AddDate(0, 1, 0)
		monthEnd := nextMonth.Add(-time.Nanosecond)
		
		periodEndInMonth := monthEnd
		endEndOfDay := time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, end.Location())
		if endEndOfDay.Before(monthEnd) {
			periodEndInMonth = endEndOfDay
		}

		d1 := time.Date(periodStartInMonth.Year(), periodStartInMonth.Month(), periodStartInMonth.Day(), 0, 0, 0, 0, time.UTC)
		d2 := time.Date(periodEndInMonth.Year(), periodEndInMonth.Month(), periodEndInMonth.Day(), 0, 0, 0, 0, time.UTC)
		
		activeDays := int(d2.Sub(d1).Hours() / 24) + 1
		if activeDays < 0 {
			activeDays = 0
		}
		
		daysInMonth := daysIn(monthStart.Month(), monthStart.Year())
		
		factor := 0.0
		if daysInMonth > 0 {
			factor = float64(activeDays) / float64(daysInMonth)
		}

		// Apply to each user
		for userID, target := range targets {
			if target != nil {
				result[userID] += float64(target.TargetAmount) * factor
			}
		}

		monthStart = monthStart.AddDate(0, 1, 0)
	}

	return result, nil
}

func daysIn(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// GetTotalEffectiveTarget gets the sum of user monthly target records for a specific month
func (r *repository) GetTotalEffectiveTarget(year int, month int) (int64, error) {
	var total int64

	// Sum only user targets to match sales performance and target matrix views.
	err := r.db.Table("monthly_targets").
		Select("COALESCE(SUM(target_amount), 0)").
		Where("year = ? AND month = ? AND user_id IS NOT NULL AND deleted_at IS NULL", year, month).
		Scan(&total).Error
		
	if err != nil {
		return 0, err
	}
	
	return total, nil
}
