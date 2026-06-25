package sales_overview

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/account"
	"github.com/gilabs/crm-healthcare/api/internal/domain/sales_overview"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

const (
	defaultProspectReason = "no_reason_provided"
	prospectTypeDeal      = "deal"
	prospectTypeLead      = "lead"
)

// NewRepository creates a new sales overview repository
func NewRepository(db *gorm.DB) interfaces.SalesOverviewRepository {
	return &repository{db: db}
}

// formatCurrency formats integer (sen) to formatted currency string
func formatCurrency(amount int64) string {
	rupiah := float64(amount) / 100.0
	if rupiah < 0 {
		return "-Rp " + formatNumber(-rupiah)
	}
	return "Rp " + formatNumber(rupiah)
}

func formatNumber(n float64) string {
	// ... (no change to formatNumber implementation)
	amount := int64(n)
	if amount == 0 {
		return "0"
	}
	str := ""
	if amount >= 1000 {
		parts := make([]string, 0)
		for amount > 0 {
			part := amount % 1000
			if amount >= 1000 {
				parts = append([]string{fmt.Sprintf("%03d", part)}, parts...)
			} else {
				parts = append([]string{fmt.Sprintf("%d", part)}, parts...)
			}
			amount = amount / 1000
		}
		str = strings.Join(parts, ".")
	} else {
		str = fmt.Sprintf("%d", amount)
	}
	return str
}

func normalizeProspectReason(reason string) string {
	normalized := strings.TrimSpace(strings.ToLower(reason))
	if normalized == "" || normalized == defaultProspectReason {
		return "unknown"
	}

	switch {
	case strings.Contains(normalized, "price"), strings.Contains(normalized, "harga"), strings.Contains(normalized, "cost"), strings.Contains(normalized, "budget"):
		return "price"
	case strings.Contains(normalized, "compet"), strings.Contains(normalized, "rival"), strings.Contains(normalized, "vendor lain"), strings.Contains(normalized, "pesaing"):
		return "competition"
	case strings.Contains(normalized, "timing"), strings.Contains(normalized, "later"), strings.Contains(normalized, "delay"), strings.Contains(normalized, "tunda"), strings.Contains(normalized, "follow up nanti"):
		return "timing"
	case strings.Contains(normalized, "fit"), strings.Contains(normalized, "feature"), strings.Contains(normalized, "need"), strings.Contains(normalized, "kebutuhan"), strings.Contains(normalized, "spec"), strings.Contains(normalized, "produk"):
		return "fit"
	case strings.Contains(normalized, "decision"), strings.Contains(normalized, "approval"), strings.Contains(normalized, "approv"), strings.Contains(normalized, "otor"), strings.Contains(normalized, "authority"):
		return "decision"
	case strings.Contains(normalized, "relationship"), strings.Contains(normalized, "trust"), strings.Contains(normalized, "relasi"), strings.Contains(normalized, "service"):
		return "relationship"
	default:
		return "other"
	}
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func endOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
}

func startOfQuarter(t time.Time) time.Time {
	month := ((int(t.Month())-1)/3)*3 + 1
	return time.Date(t.Year(), time.Month(month), 1, 0, 0, 0, 0, t.Location())
}

func endOfQuarter(t time.Time) time.Time {
	return startOfQuarter(t).AddDate(0, 3, 0).Add(-time.Second)
}

func normalizeTrendRange(startDate, endDate interface{}, trendMode string) (time.Time, time.Time) {
	now := time.Now()
	var start time.Time
	var end time.Time

	if t, ok := startDate.(time.Time); ok {
		start = t
	}
	if t, ok := endDate.(time.Time); ok {
		end = t
	}

	switch trendMode {
	case "rolling_30d":
		end = endOfDay(now)
		start = startOfDay(end.AddDate(0, 0, -29))
	case "rolling_90d":
		end = endOfDay(now)
		start = startOfDay(end.AddDate(0, 0, -89))
	default:
		if start.IsZero() {
			start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		}
		if end.IsZero() {
			end = endOfDay(now)
		}
	}

	return start, end
}

func buildTrendPeriods(start, end time.Time, trendMode string) []sales_overview.MonthlySalesData {
	periods := make([]sales_overview.MonthlySalesData, 0)
	current := start

	switch trendMode {
	case "qoq":
		current = startOfQuarter(start)
		for !current.After(end) {
			periodEnd := endOfQuarter(current)
			quarter := ((int(current.Month()) - 1) / 3) + 1
			periods = append(periods, sales_overview.MonthlySalesData{
				Month:       int(current.Month()),
				MonthName:   current.Month().String(),
				Year:        current.Year(),
				PeriodKey:   fmt.Sprintf("%d-Q%d", current.Year(), quarter),
				PeriodLabel: fmt.Sprintf("Q%d %d", quarter, current.Year()),
				PeriodStart: current,
				PeriodEnd:   periodEnd,
			})
			current = current.AddDate(0, 3, 0)
		}
	case "rolling_30d", "rolling_90d":
		current = startOfDay(start)
		for !current.After(end) {
			periods = append(periods, sales_overview.MonthlySalesData{
				Month:       int(current.Month()),
				MonthName:   current.Month().String(),
				Year:        current.Year(),
				PeriodKey:   current.Format("2006-01-02"),
				PeriodLabel: current.Format("02 Jan"),
				PeriodStart: current,
				PeriodEnd:   endOfDay(current),
			})
			current = current.AddDate(0, 0, 1)
		}
	default:
		current = time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, start.Location())
		for !current.After(end) {
			periodEnd := time.Date(current.Year(), current.Month()+1, 1, 0, 0, 0, 0, current.Location()).Add(-time.Second)
			periods = append(periods, sales_overview.MonthlySalesData{
				Month:       int(current.Month()),
				MonthName:   current.Month().String(),
				Year:        current.Year(),
				PeriodKey:   current.Format("2006-01"),
				PeriodLabel: current.Format("Jan 2006"),
				PeriodStart: current,
				PeriodEnd:   periodEnd,
			})
			current = current.AddDate(0, 1, 0)
		}
	}

	return periods
}

func applyDealOutcomeDateFilter(query *gorm.DB, startDate, endDate interface{}, tablePrefix string) *gorm.DB {
	if tablePrefix != "" && !strings.HasSuffix(tablePrefix, ".") {
		tablePrefix += "."
	}

	actualCloseDateColumn := tablePrefix + "actual_close_date"
	createdAtColumn := tablePrefix + "created_at"

	if startDate != nil {
		query = query.Where(
			fmt.Sprintf("(%s >= ? OR (%s IS NULL AND %s >= ?))", actualCloseDateColumn, actualCloseDateColumn, createdAtColumn),
			startDate,
			startDate,
		)
	}
	if endDate != nil {
		query = query.Where(
			fmt.Sprintf("(%s <= ? OR (%s IS NULL AND %s <= ?))", actualCloseDateColumn, actualCloseDateColumn, createdAtColumn),
			endDate,
			endDate,
		)
	}

	return query
}

func applyProspectOutcomeDateFilter(query *gorm.DB, startDate, endDate interface{}) *gorm.DB {
	if startDate != nil {
		query = query.Where("outcome_at >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("outcome_at <= ?", endDate)
	}

	return query
}

func restrictDealAssignedToSalesRole(query *gorm.DB, dealAlias string) *gorm.DB {
	sanitizedAlias := strings.NewReplacer(".", "_").Replace(dealAlias)
	userAlias := sanitizedAlias + "_sales_scope_user"
	roleAlias := sanitizedAlias + "_sales_scope_role"

	return query.
		Joins(fmt.Sprintf(
			"INNER JOIN users %s ON %s.id = %s.assigned_to AND %s.deleted_at IS NULL",
			userAlias, userAlias, dealAlias, userAlias,
		)).
		Joins(fmt.Sprintf(
			"INNER JOIN roles %s ON %s.id = %s.role_id AND %s.deleted_at IS NULL AND %s.code = ?",
			roleAlias, roleAlias, userAlias, roleAlias, roleAlias,
		), "sales")
}

func restrictSalesRepToSalesRole(query *gorm.DB, userColumn string) *gorm.DB {
	sanitizedColumn := strings.NewReplacer(".", "_").Replace(userColumn)
	userAlias := sanitizedColumn + "_sales_scope_user"
	roleAlias := sanitizedColumn + "_sales_scope_role"

	return query.
		Joins(fmt.Sprintf(
			"INNER JOIN users %s ON %s.id = %s AND %s.deleted_at IS NULL",
			userAlias, userAlias, userColumn, userAlias,
		)).
		Joins(fmt.Sprintf(
			"INNER JOIN roles %s ON %s.id = %s.role_id AND %s.deleted_at IS NULL AND %s.code = ?",
			roleAlias, roleAlias, userAlias, roleAlias, roleAlias,
		), "sales")
}

func prospectOutcomeDatasetSQL() string {
	leadClosedAtExpression := "COALESCE(NULLIF(l.conversion_metadata->>'latest_status_changed_at', '')::timestamptz, l.updated_at, l.created_at)"

	return fmt.Sprintf(`
		(
			SELECT
				d.id::text AS id,
				'%s' AS type,
				d.title AS title,
				COALESCE(a.name, '') AS account_name,
				COALESCE(d.assigned_to::text, '') AS sales_rep_id,
				COALESCE(u.name, '') AS sales_rep_name,
				COALESCE(u.email, '') AS sales_rep_email,
				COALESCE(u.avatar_url, '') AS sales_rep_avatar_url,
				LOWER(COALESCE(d.status, 'open')) AS status,
				COALESCE(d.value, 0) AS value,
				COALESCE(NULLIF(TRIM(d.close_reason), ''), '') AS reason,
				COALESCE(d.source, '') AS source,
				d.created_at AS created_at,
				d.actual_close_date AS closed_at,
				COALESCE(d.actual_close_date, d.created_at) AS outcome_at
			FROM deals d
			LEFT JOIN accounts a ON d.account_id = a.id AND a.deleted_at IS NULL
			INNER JOIN users u ON d.assigned_to = u.id AND u.deleted_at IS NULL
			INNER JOIN roles ur ON u.role_id = ur.id AND ur.deleted_at IS NULL AND ur.code = 'sales'
			WHERE d.deleted_at IS NULL

			UNION ALL

			SELECT
				l.id::text AS id,
				'%s' AS type,
				COALESCE(
					NULLIF(TRIM(l.company_name), ''),
					NULLIF(TRIM(CONCAT(l.first_name, ' ', COALESCE(l.last_name, ''))), ''),
					l.email
				) AS title,
				COALESCE(NULLIF(TRIM(l.company_name), ''), '') AS account_name,
				COALESCE(l.assigned_to::text, '') AS sales_rep_id,
				COALESCE(u.name, '') AS sales_rep_name,
				COALESCE(u.email, '') AS sales_rep_email,
				COALESCE(u.avatar_url, '') AS sales_rep_avatar_url,
				CASE
					WHEN LOWER(COALESCE(l.lead_status, '')) = 'converted' THEN 'won'
					WHEN LOWER(COALESCE(l.lead_status, '')) = 'lost' THEN 'lost'
					ELSE LOWER(COALESCE(l.lead_status, 'open'))
				END AS status,
				COALESCE(l.estimated_value, 0) AS value,
				COALESCE(NULLIF(TRIM(l.conversion_metadata->>'latest_status_reason'), ''), '') AS reason,
				COALESCE(l.lead_source, '') AS source,
				l.created_at AS created_at,
				CASE
					WHEN LOWER(COALESCE(l.lead_status, '')) IN ('converted', 'lost') THEN %s
					ELSE NULL
				END AS closed_at,
				CASE
					WHEN LOWER(COALESCE(l.lead_status, '')) IN ('converted', 'lost') THEN %s
					ELSE l.created_at
				END AS outcome_at
			FROM leads l
			INNER JOIN users u ON l.assigned_to = u.id AND u.deleted_at IS NULL
			INNER JOIN roles ur ON u.role_id = ur.id AND ur.deleted_at IS NULL AND ur.code = 'sales'
			WHERE l.deleted_at IS NULL
		) AS prospects
	`, prospectTypeDeal, prospectTypeLead, leadClosedAtExpression, leadClosedAtExpression)
}

func (r *repository) getProspectReasonBreakdown(userID, status string, total int64, startDate, endDate interface{}) ([]sales_overview.ProspectReasonBreakdown, error) {
	if total == 0 {
		return []sales_overview.ProspectReasonBreakdown{}, nil
	}

	reasonExpression := "COALESCE(NULLIF(TRIM(reason), ''), '" + defaultProspectReason + "')"
	var rows []struct {
		Reason string
		Count  int
	}

	query := r.db.Table(prospectOutcomeDatasetSQL()).
		Select(reasonExpression+" as reason, COUNT(*) as count").
		Where("sales_rep_id = ? AND status = ?", userID, status)
	query = applyProspectOutcomeDateFilter(query, startDate, endDate)

	if err := query.Group(reasonExpression).Order("count DESC, reason ASC").Scan(&rows).Error; err != nil {
		return nil, err
	}

	breakdown := make([]sales_overview.ProspectReasonBreakdown, len(rows))
	for i, row := range rows {
		breakdown[i] = sales_overview.ProspectReasonBreakdown{
			Reason:     row.Reason,
			Category:   normalizeProspectReason(row.Reason),
			Count:      row.Count,
			Percentage: (float64(row.Count) / float64(total)) * 100.0,
		}
	}

	return breakdown, nil
}

func (r *repository) getProspectOutcomeSummary(userID string, startDate, endDate interface{}, includeRecent bool) (*sales_overview.ProspectOutcomeSummary, error) {
	var counts struct {
		TotalProspects int64
		WonProspects   int64
		LostProspects  int64
		OpenProspects  int64
	}

	countQuery := r.db.Table(prospectOutcomeDatasetSQL()).
		Where("sales_rep_id = ?", userID)
	countQuery = applyProspectOutcomeDateFilter(countQuery, startDate, endDate)

	if err := countQuery.Select(`
			COUNT(*) as total_prospects,
			COALESCE(SUM(CASE WHEN status = 'won' THEN 1 ELSE 0 END), 0) as won_prospects,
			COALESCE(SUM(CASE WHEN status = 'lost' THEN 1 ELSE 0 END), 0) as lost_prospects,
			COALESCE(SUM(CASE WHEN status NOT IN ('won', 'lost') THEN 1 ELSE 0 END), 0) as open_prospects
		`).Scan(&counts).Error; err != nil {
		return nil, err
	}

	wonReasons, err := r.getProspectReasonBreakdown(userID, "won", counts.WonProspects, startDate, endDate)
	if err != nil {
		return nil, err
	}

	lostReasons, err := r.getProspectReasonBreakdown(userID, "lost", counts.LostProspects, startDate, endDate)
	if err != nil {
		return nil, err
	}

	conversionRate := 0.0
	if counts.TotalProspects > 0 {
		conversionRate = (float64(counts.WonProspects) / float64(counts.TotalProspects)) * 100.0
	}

	summary := &sales_overview.ProspectOutcomeSummary{
		TotalProspects:         int(counts.TotalProspects),
		WonProspects:           int(counts.WonProspects),
		LostProspects:          int(counts.LostProspects),
		OpenProspects:          int(counts.OpenProspects),
		ProspectConversionRate: conversionRate,
		WonReasons:             wonReasons,
		LostReasons:            lostReasons,
	}

	if !includeRecent {
		return summary, nil
	}

	var rows []struct {
		ID              string
		Type            string
		Title           string
		AccountName     string
		Status          string
		Value           int64
		Reason          string
		Source          string
		CreatedAt       time.Time
		ActualCloseDate *time.Time
	}

	recentQuery := r.db.Table(prospectOutcomeDatasetSQL()).
		Select(`
			id,
			type,
			title,
			account_name,
			status,
			value,
			COALESCE(NULLIF(TRIM(reason), ''), '') as reason,
			source,
			created_at,
			closed_at as actual_close_date
		`).
		Where("sales_rep_id = ?", userID)
	recentQuery = applyProspectOutcomeDateFilter(recentQuery, startDate, endDate)

	if err := recentQuery.
		Order("outcome_at DESC").
		Limit(10).
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	summary.RecentProspects = make([]sales_overview.ProspectOutcomeItem, len(rows))
	for i, row := range rows {
		reason := row.Reason
		if reason == "" && (row.Status == "won" || row.Status == "lost") {
			reason = defaultProspectReason
		}

		summary.RecentProspects[i] = sales_overview.ProspectOutcomeItem{
			ID:             row.ID,
			Type:           row.Type,
			Title:          row.Title,
			AccountName:    row.AccountName,
			Status:         row.Status,
			Value:          row.Value,
			ValueFormatted: formatCurrency(row.Value),
			Reason:         reason,
			ReasonCategory: normalizeProspectReason(reason),
			Source:         row.Source,
			CreatedAt:      row.CreatedAt,
			ClosedAt:       row.ActualCloseDate,
		}
	}

	return summary, nil
}

// ListProspectOutcomes lists prospect outcomes across scoped sales reps
func (r *repository) ListProspectOutcomes(req *sales_overview.ListProspectOutcomesRequest, startDate, endDate interface{}) ([]sales_overview.ProspectOutcomeListItem, int64, error) {
	var total int64
	var rows []struct {
		ID                string
		Type              string
		Title             string
		AccountName       string
		SalesRepID        string
		SalesRepName      string
		SalesRepEmail     string
		SalesRepAvatarURL string
		Status            string
		Value             int64
		Reason            string
		Source            string
		CreatedAt         time.Time
		ActualCloseDate   *time.Time
	}

	query := r.db.Table(prospectOutcomeDatasetSQL()).
		Where("sales_rep_id <> ''")
	query = applyProspectOutcomeDateFilter(query, startDate, endDate)

	if len(req.ScopedUserIDs) > 0 {
		query = query.Where("sales_rep_id IN ?", req.ScopedUserIDs)
	}
	if req.SalesUserID != "" {
		query = query.Where("sales_rep_id = ?", req.SalesUserID)
	}
	if req.Status != "" {
		normalizedStatus := strings.ToLower(strings.TrimSpace(req.Status))
		if normalizedStatus == "open" {
			query = query.Where("status NOT IN ('won', 'lost')")
		} else {
			query = query.Where("status = ?", normalizedStatus)
		}
	}
	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where(`
			LOWER(COALESCE(title, '')) LIKE ?
			OR LOWER(COALESCE(account_name, '')) LIKE ?
			OR LOWER(COALESCE(sales_rep_name, '')) LIKE ?
			OR LOWER(COALESCE(reason, '')) LIKE ?
			OR LOWER(COALESCE(source, '')) LIKE ?
		`, search, search, search, search, search)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}
	offset := (page - 1) * perPage

	if err := query.Select(`
				id,
				type,
				title,
				account_name,
				sales_rep_id,
				sales_rep_name,
				sales_rep_email,
				sales_rep_avatar_url,
				status,
				value,
				COALESCE(NULLIF(TRIM(reason), ''), '') as reason,
				source,
				created_at,
				closed_at as actual_close_date
			`).
		Order("outcome_at DESC").
		Limit(perPage).
		Offset(offset).
		Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	results := make([]sales_overview.ProspectOutcomeListItem, len(rows))
	for i, row := range rows {
		reason := row.Reason
		if reason == "" && (row.Status == "won" || row.Status == "lost") {
			reason = defaultProspectReason
		}

		results[i] = sales_overview.ProspectOutcomeListItem{
			ID:                row.ID,
			Type:              row.Type,
			Title:             row.Title,
			AccountName:       row.AccountName,
			SalesRepID:        row.SalesRepID,
			SalesRepName:      row.SalesRepName,
			SalesRepEmail:     row.SalesRepEmail,
			SalesRepAvatarURL: row.SalesRepAvatarURL,
			Status:            row.Status,
			Value:             row.Value,
			ValueFormatted:    formatCurrency(row.Value),
			Reason:            reason,
			ReasonCategory:    normalizeProspectReason(reason),
			Source:            row.Source,
			CreatedAt:         row.CreatedAt,
			ClosedAt:          row.ActualCloseDate,
		}
	}

	return results, total, nil
}

// GetSalesPerformanceDetail gets detailed performance metrics for a user
func (r *repository) GetSalesPerformanceDetail(userID string, startDate, endDate interface{}) (*sales_overview.SalesPerformanceDetail, error) {
	var totalDeals, wonDeals, lostDeals, openDeals int64
	var totalRevenue, wonRevenue int64
	var visitsCompleted, tasksCompleted, totalTasks int64

	// Calculate deals metrics
	// Total deals: count all deals created in period
	totalDealsQuery := r.db.Table("deals").
		Where("assigned_to = ?", userID)
	if startDate != nil {
		totalDealsQuery = totalDealsQuery.Where("created_at >= ?", startDate)
	}
	if endDate != nil {
		totalDealsQuery = totalDealsQuery.Where("created_at <= ?", endDate)
	}

	if err := totalDealsQuery.Count(&totalDeals).Error; err != nil {
		return nil, err
	}

	// Won deals and revenue: filter by actual_close_date (when deal was actually closed/won)
	// Use actual_close_date for revenue calculation, fallback to created_at if NULL
	wonDealsQuery := r.db.Table("deals").
		Where("assigned_to = ? AND status = ?", userID, "won")
	if startDate != nil {
		wonDealsQuery = wonDealsQuery.Where("(actual_close_date >= ? OR (actual_close_date IS NULL AND created_at >= ?))", startDate, startDate)
	}
	if endDate != nil {
		wonDealsQuery = wonDealsQuery.Where("(actual_close_date <= ? OR (actual_close_date IS NULL AND created_at <= ?))", endDate, endDate)
	}

	// Count won deals
	if err := wonDealsQuery.Count(&wonDeals).Error; err != nil {
		return nil, err
	}

	// Sum revenue from won deals
	var revenueResult struct {
		Total int64
	}
	if err := wonDealsQuery.Select("COALESCE(SUM(value), 0) as total").Scan(&revenueResult).Error; err != nil {
		return nil, err
	}
	wonRevenue = revenueResult.Total
	totalRevenue = wonRevenue // Total revenue is from won deals only

	// Count lost deals (filter by created_at for all deals)
	lostDealsQuery := r.db.Table("deals").
		Where("assigned_to = ? AND status = ?", userID, "lost")
	if startDate != nil {
		lostDealsQuery = lostDealsQuery.Where("created_at >= ?", startDate)
	}
	if endDate != nil {
		lostDealsQuery = lostDealsQuery.Where("created_at <= ?", endDate)
	}
	if err := lostDealsQuery.Count(&lostDeals).Error; err != nil {
		return nil, err
	}

	// Count open deals (filter by created_at for all deals)
	openDealsQuery := r.db.Table("deals").
		Where("assigned_to = ? AND status = ?", userID, "open")
	if startDate != nil {
		openDealsQuery = openDealsQuery.Where("created_at >= ?", startDate)
	}
	if endDate != nil {
		openDealsQuery = openDealsQuery.Where("created_at <= ?", endDate)
	}
	if err := openDealsQuery.Count(&openDeals).Error; err != nil {
		return nil, err
	}

	// Calculate visits completed
	visitsQuery := r.db.Table("visit_reports").
		Where("sales_rep_id = ?", userID)
	if startDate != nil {
		visitsQuery = visitsQuery.Where("visit_date >= ?", startDate)
	}
	if endDate != nil {
		visitsQuery = visitsQuery.Where("visit_date <= ?", endDate)
	}
	if err := visitsQuery.Where("status IN ?", []string{"completed", "approved"}).Count(&visitsCompleted).Error; err != nil {
		return nil, err
	}

	// Calculate tasks completed
	tasksQuery := r.db.Table("tasks").
		Where("assigned_to = ?", userID)
	if startDate != nil {
		tasksQuery = tasksQuery.Where("created_at >= ?", startDate)
	}
	if endDate != nil {
		tasksQuery = tasksQuery.Where("created_at <= ?", endDate)
	}
	if err := tasksQuery.Count(&totalTasks).Error; err != nil {
		return nil, err
	}
	if err := tasksQuery.Where("status = ?", "completed").Count(&tasksCompleted).Error; err != nil {
		return nil, err
	}

	// Calculate conversion rate
	conversionRate := 0.0
	if totalDeals > 0 {
		conversionRate = (float64(wonDeals) / float64(totalDeals)) * 100.0
	}

	// Calculate average deal value
	avgDealValue := 0.0
	if wonDeals > 0 {
		avgDealValue = float64(wonRevenue) / float64(wonDeals)
	}

	// Calculate task completion rate
	taskCompletionRate := 0.0
	if totalTasks > 0 {
		taskCompletionRate = (float64(tasksCompleted) / float64(totalTasks)) * 100.0
	}

	// Get user info
	var userObj user.User
	if err := r.db.
		Joins("INNER JOIN roles user_roles ON users.role_id = user_roles.id AND user_roles.deleted_at IS NULL").
		Where("users.id = ? AND user_roles.code = ?", userID, "sales").
		First(&userObj).Error; err != nil {
		return nil, err
	}

	var startTime, endTime time.Time
	if startDate != nil {
		if t, ok := startDate.(time.Time); ok {
			startTime = t
		}
	}
	if endDate != nil {
		if t, ok := endDate.(time.Time); ok {
			endTime = t
		}
	}

	prospectOutcome, err := r.getProspectOutcomeSummary(userID, startDate, endDate, true)
	if err != nil {
		return nil, err
	}

	detail := &sales_overview.SalesPerformanceDetail{
		UserID:                    userID,
		User:                      userObj.ToUserResponse(),
		PeriodStart:               startTime,
		PeriodEnd:                 endTime,
		TotalRevenue:              totalRevenue,
		TotalRevenueFormatted:     formatCurrency(totalRevenue),
		WonDeals:                  int(wonDeals),
		TotalDeals:                int(totalDeals),
		LostDeals:                 int(lostDeals),
		OpenDeals:                 int(openDeals),
		ConversionRate:            conversionRate,
		AverageDealValue:          avgDealValue,
		AverageDealValueFormatted: formatCurrency(int64(avgDealValue)),
		VisitsCompleted:           int(visitsCompleted),
		TasksCompleted:            int(tasksCompleted),
		TotalTasks:                int(totalTasks),
		TaskCompletionRate:        taskCompletionRate,
		ProspectOutcome:           prospectOutcome,
	}

	return detail, nil
}

// GetSalesRepDetail gets comprehensive detail for sales rep detail page
func (r *repository) GetSalesRepDetail(userID string, startDate, endDate interface{}) (*sales_overview.SalesRepDetail, error) {
	// Get user info
	var userObj user.User
	if err := r.db.
		Joins("INNER JOIN roles user_roles ON users.role_id = user_roles.id AND user_roles.deleted_at IS NULL").
		Where("users.id = ? AND user_roles.code = ?", userID, "sales").
		First(&userObj).Error; err != nil {
		return nil, err
	}

	// Get statistics
	stats, err := r.getSalesRepStatistics(userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	var startTime, endTime *time.Time
	if startDate != nil {
		if t, ok := startDate.(time.Time); ok {
			startTime = &t
		}
	}
	if endDate != nil {
		if t, ok := endDate.(time.Time); ok {
			endTime = &t
		}
	}

	detail := &sales_overview.SalesRepDetail{
		UserID:      userID,
		User:        userObj.ToUserResponse(),
		PeriodStart: startTime,
		PeriodEnd:   endTime,
		Statistics:  stats,
	}

	return detail, nil
}

// getSalesRepStatistics calculates statistics for sales rep
func (r *repository) getSalesRepStatistics(userID string, startDate, endDate interface{}) (*sales_overview.SalesRepStatistics, error) {
	var totalRevenue int64
	var dealsClosed, visitsCompleted, tasksCompleted int64

	// Calculate revenue from won deals
	dealsQuery := r.db.Table("deals").
		Where("assigned_to = ? AND status = ?", userID, "won")
	if startDate != nil {
		dealsQuery = dealsQuery.Where("actual_close_date >= ? OR (actual_close_date IS NULL AND created_at >= ?)", startDate, startDate)
	}
	if endDate != nil {
		dealsQuery = dealsQuery.Where("actual_close_date <= ? OR (actual_close_date IS NULL AND created_at <= ?)", endDate, endDate)
	}

	var revenueResult struct {
		Total int64
	}
	if err := dealsQuery.Select("COALESCE(SUM(value), 0) as total").Scan(&revenueResult).Error; err != nil {
		return nil, err
	}
	totalRevenue = revenueResult.Total

	// Count deals closed
	if err := dealsQuery.Count(&dealsClosed).Error; err != nil {
		return nil, err
	}

	// Count visits completed
	visitsQuery := r.db.Table("visit_reports").
		Where("sales_rep_id = ? AND status IN ?", userID, []string{"completed", "approved"})
	if startDate != nil {
		visitsQuery = visitsQuery.Where("visit_date >= ?", startDate)
	}
	if endDate != nil {
		visitsQuery = visitsQuery.Where("visit_date <= ?", endDate)
	}
	if err := visitsQuery.Count(&visitsCompleted).Error; err != nil {
		return nil, err
	}

	// Count tasks completed
	tasksQuery := r.db.Table("tasks").
		Where("assigned_to = ? AND (completed_at IS NOT NULL OR status = ?)", userID, "completed")
	if startDate != nil {
		tasksQuery = tasksQuery.Where("COALESCE(completed_at, updated_at, created_at) >= ?", startDate)
	}
	if endDate != nil {
		tasksQuery = tasksQuery.Where("COALESCE(completed_at, updated_at, created_at) <= ?", endDate)
	}
	if err := tasksQuery.Count(&tasksCompleted).Error; err != nil {
		return nil, err
	}

	// Calculate conversion rate (needs total deals)
	var totalDeals int64
	totalDealsQuery := r.db.Table("deals").Where("assigned_to = ?", userID)
	if startDate != nil {
		totalDealsQuery = totalDealsQuery.Where("created_at >= ?", startDate)
	}
	if endDate != nil {
		totalDealsQuery = totalDealsQuery.Where("created_at <= ?", endDate)
	}
	if err := totalDealsQuery.Count(&totalDeals).Error; err != nil {
		return nil, err
	}

	conversionRate := 0.0
	if totalDeals > 0 {
		conversionRate = (float64(dealsClosed) / float64(totalDeals)) * 100.0
	}

	avgDealValue := 0.0
	if dealsClosed > 0 {
		avgDealValue = float64(totalRevenue) / float64(dealsClosed)
	}

	prospectOutcome, err := r.getProspectOutcomeSummary(userID, startDate, endDate, true)
	if err != nil {
		return nil, err
	}

	stats := &sales_overview.SalesRepStatistics{
		TotalRevenue:              totalRevenue,
		TotalRevenueFormatted:     formatCurrency(totalRevenue),
		DealsClosed:               int(dealsClosed),
		VisitsCompleted:           int(visitsCompleted),
		TasksCompleted:            int(tasksCompleted),
		ConversionRate:            conversionRate,
		AverageDealValue:          avgDealValue,
		AverageDealValueFormatted: formatCurrency(int64(avgDealValue)),
		ProspectOutcome:           prospectOutcome,
	}

	// TODO: Add period comparison if needed
	// stats.PeriodComparison = r.calculatePeriodComparison(userID, startDate, endDate)

	return stats, nil
}

// ListSalesPerformance lists all sales reps with performance summary (management overview)
// Optimized to use JOINs and Aggregations instead of N+1 queries
func (r *repository) ListSalesPerformance(req *sales_overview.ListSalesPerformanceRequest) ([]sales_overview.SalesPerformanceListResponse, int64, error) {
	var total int64
	var rawResults []struct {
		UserID          string
		UserName        string
		UserEmail       string
		AvatarURL       string
		TotalRevenue    int64
		DealsClosed     int
		VisitsCompleted int
		TasksCompleted  int
		TotalDeals      int
	}

	// Parse date range
	var startDate, endDate interface{}
	if req.StartDate != "" {
		parsed, err := time.Parse("2006-01-02", req.StartDate)
		if err == nil {
			startDate = parsed
		}
	}
	if req.EndDate != "" {
		parsed, err := time.Parse("2006-01-02", req.EndDate)
		if err == nil {
			endDate = parsed.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		}
	}

	// Build Joins Dynamically to avoid "could not determine data type of parameter $1" error regarding ? IS NULL

	// 1. Deals Join
	dealsJoin := `
		LEFT JOIN deals d ON users.id = d.assigned_to 
		AND d.deleted_at IS NULL
	`
	var dealsArgs []interface{}
	if startDate != nil {
		dealsJoin += " AND (d.actual_close_date >= ? OR (d.actual_close_date IS NULL AND d.created_at >= ?))"
		dealsArgs = append(dealsArgs, startDate, startDate)
	}
	if endDate != nil {
		dealsJoin += " AND (d.actual_close_date <= ? OR (d.actual_close_date IS NULL AND d.created_at <= ?))"
		dealsArgs = append(dealsArgs, endDate, endDate)
	}

	// 2. Visits Join
	visitsCondition := "status IN ('completed', 'approved')"
	var visitsArgs []interface{}
	if startDate != nil {
		visitsCondition += " AND visit_date >= ?"
		visitsArgs = append(visitsArgs, startDate)
	}
	if endDate != nil {
		visitsCondition += " AND visit_date <= ?"
		visitsArgs = append(visitsArgs, endDate)
	}
	visitsJoin := fmt.Sprintf(`
		LEFT JOIN (
			SELECT sales_rep_id, COUNT(*) as completed_count
			FROM visit_reports
			WHERE %s
			GROUP BY sales_rep_id
		) visits ON users.id = visits.sales_rep_id
	`, visitsCondition)

	// 3. Tasks Join
	tasksCondition := "status = 'completed'"
	var tasksArgs []interface{}
	if startDate != nil {
		tasksCondition += " AND created_at >= ?"
		tasksArgs = append(tasksArgs, startDate)
	}
	if endDate != nil {
		tasksCondition += " AND created_at <= ?"
		tasksArgs = append(tasksArgs, endDate)
	}
	tasksJoin := fmt.Sprintf(`
		LEFT JOIN (
			SELECT assigned_to, COUNT(*) as completed_count
			FROM tasks
			WHERE %s
			GROUP BY assigned_to
		) tasks ON users.id = tasks.assigned_to
	`, tasksCondition)

	targetSelect := "0 as target_amount"
	targetJoin := ""
	var targetArgs []interface{}
	if startDate != nil && endDate != nil {
		targetSelect = "COALESCE(targets.target_amount, 0) as target_amount"
		targetJoin = `
			LEFT JOIN (
				SELECT u.id as user_id,
					COALESCE(SUM(
						COALESCE(ut.target_amount, gt.target_amount, 0) *
						GREATEST(
							(LEAST(?::date, (date_trunc('month', ms.month_start) + INTERVAL '1 month' - INTERVAL '1 day')::date) -
							GREATEST(?::date, date_trunc('month', ms.month_start)::date) + 1),
							0
						)::numeric
						/ EXTRACT(day FROM (date_trunc('month', ms.month_start) + INTERVAL '1 month' - INTERVAL '1 day'))
					), 0) as target_amount
				FROM users u
				CROSS JOIN (
					SELECT generate_series(date_trunc('month', ?::date), date_trunc('month', ?::date), interval '1 month') as month_start
				) ms
				LEFT JOIN monthly_targets ut ON ut.user_id = u.id
					AND ut.year = EXTRACT(year FROM ms.month_start)
					AND ut.month = EXTRACT(month FROM ms.month_start)
					AND ut.deleted_at IS NULL
				LEFT JOIN monthly_targets gt ON gt.group_id = u.group_id
					AND gt.year = EXTRACT(year FROM ms.month_start)
					AND gt.month = EXTRACT(month FROM ms.month_start)
					AND gt.deleted_at IS NULL
				WHERE u.deleted_at IS NULL
				GROUP BY u.id
			) targets ON users.id = targets.user_id
		`
		targetArgs = append(targetArgs, endDate, startDate, startDate, endDate)
	}

	selectFields := fmt.Sprintf(`
			users.id as user_id,
			users.name as user_name,
			users.email as user_email,
			users.avatar_url,
			COALESCE(SUM(CASE WHEN d.status = 'won' THEN d.value ELSE 0 END), 0) as total_revenue,
			COALESCE(SUM(CASE WHEN d.status = 'won' THEN 1 ELSE 0 END), 0) as deals_closed,
			COALESCE(COUNT(d.id), 0) as total_deals,
			COALESCE(visits.completed_count, 0) as visits_completed,
			COALESCE(tasks.completed_count, 0) as tasks_completed,
			%s
		`, targetSelect)

	// Build the main query
	query := r.db.Table("users").
		Select(selectFields).
		Joins("INNER JOIN roles user_roles ON users.role_id = user_roles.id AND user_roles.deleted_at IS NULL").
		Joins(dealsJoin, dealsArgs...).
		Joins(visitsJoin, visitsArgs...).
		Joins(tasksJoin, tasksArgs...).
		Where("users.deleted_at IS NULL").
		Where("user_roles.code = ?", "sales")

	if targetJoin != "" {
		query = query.Joins(targetJoin, targetArgs...)
	}

	// Apply Filters
	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where("LOWER(users.name) LIKE ? OR LOWER(users.email) LIKE ?", search, search)
	}

	if req.BrickID != "" {
		query = query.Where("users.brick_id = ?", req.BrickID)
	}

	// Apply RBAC scope filtering
	if len(req.ScopedUserIDs) > 0 {
		query = query.Where("users.id IN ?", req.ScopedUserIDs)
	}

	// Group By
	if targetJoin != "" {
		query = query.Group("users.id, visits.completed_count, tasks.completed_count, targets.target_amount")
	} else {
		query = query.Group("users.id, visits.completed_count, tasks.completed_count")
	}

	// Count Total (Using a subquery or separate count is safer with Group By)
	// GORM's Count() with Group() can be tricky, let's use a cleaner approach for total count:
	// We count the number of users matching the filter criteria.
	countQuery := r.db.Model(&user.User{}).
		Joins("INNER JOIN roles user_roles ON users.role_id = user_roles.id AND user_roles.deleted_at IS NULL").
		Where("users.deleted_at IS NULL").
		Where("user_roles.code = ?", "sales")
	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"
		countQuery = countQuery.Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ?", search, search)
	}
	if req.BrickID != "" {
		countQuery = countQuery.Where("brick_id = ?", req.BrickID)
	}
	if len(req.ScopedUserIDs) > 0 {
		countQuery = countQuery.Where("id IN ?", req.ScopedUserIDs)
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply Sorting
	order := "desc"
	if req.Order == "asc" {
		order = "asc"
	}

	switch req.SortBy {
	case "revenue":
		query = query.Order(fmt.Sprintf("total_revenue %s", order))
	case "deals":
		query = query.Order(fmt.Sprintf("deals_closed %s", order))
	case "visits":
		query = query.Order(fmt.Sprintf("visits_completed %s", order))
	case "tasks":
		query = query.Order(fmt.Sprintf("tasks_completed %s", order))
	case "name":
		query = query.Order(fmt.Sprintf("users.name %s", order))
	case "target":
		query = query.Order(fmt.Sprintf("COALESCE(targets.target_amount, 0) %s", order))
	case "achievement":
		query = query.Order(fmt.Sprintf("CASE WHEN COALESCE(targets.target_amount, 0) > 0 THEN (COALESCE(SUM(CASE WHEN d.status = 'won' THEN d.value ELSE 0 END), 0)::numeric / COALESCE(targets.target_amount, 0)::numeric) * 100 ELSE 0 END %s", order))
	default:
		// Default sort by revenue desc
		query = query.Order("total_revenue desc")
	}

	// Apply Pagination
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	if err := query.Limit(perPage).Offset(offset).Scan(&rawResults).Error; err != nil {
		return nil, 0, err
	}

	// Map results
	results := make([]sales_overview.SalesPerformanceListResponse, len(rawResults))
	for i, rResult := range rawResults {
		conversionRate := 0.0
		if rResult.TotalDeals > 0 {
			conversionRate = (float64(rResult.DealsClosed) / float64(rResult.TotalDeals)) * 100.0
		}

		prospectOutcome, err := r.getProspectOutcomeSummary(rResult.UserID, startDate, endDate, false)
		if err != nil {
			return nil, 0, err
		}

		topWonReason := ""
		if len(prospectOutcome.WonReasons) > 0 {
			topWonReason = prospectOutcome.WonReasons[0].Reason
		}

		topLostReason := ""
		if len(prospectOutcome.LostReasons) > 0 {
			topLostReason = prospectOutcome.LostReasons[0].Reason
		}

		results[i] = sales_overview.SalesPerformanceListResponse{
			UserID:                 rResult.UserID,
			UserName:               rResult.UserName,
			UserEmail:              rResult.UserEmail,
			AvatarURL:              rResult.AvatarURL,
			TotalRevenue:           rResult.TotalRevenue,
			TotalRevenueFormatted:  formatCurrency(rResult.TotalRevenue),
			DealsClosed:            rResult.DealsClosed,
			VisitsCompleted:        rResult.VisitsCompleted,
			TasksCompleted:         rResult.TasksCompleted,
			ConversionRate:         conversionRate,
			TotalProspects:         prospectOutcome.TotalProspects,
			WonProspects:           prospectOutcome.WonProspects,
			LostProspects:          prospectOutcome.LostProspects,
			OpenProspects:          prospectOutcome.OpenProspects,
			ProspectConversionRate: prospectOutcome.ProspectConversionRate,
			TopWonReason:           topWonReason,
			TopLostReason:          topLostReason,
		}
	}

	return results, total, nil
}

// GetMonthlySalesOverview returns aggregated trend data for the chart.
func (r *repository) GetMonthlySalesOverview(startDate, endDate interface{}, trendMode string, scopedUserIDs []string) (*sales_overview.MonthlySalesOverviewResponse, error) {
	var truncateUnit string
	switch trendMode {
	case "qoq":
		truncateUnit = "quarter"
	case "rolling_30d", "rolling_90d":
		truncateUnit = "day"
	default:
		truncateUnit = "month"
	}

	normalizedStart, normalizedEnd := normalizeTrendRange(startDate, endDate, trendMode)
	periods := buildTrendPeriods(normalizedStart, normalizedEnd, trendMode)
	periodMap := make(map[string]*sales_overview.MonthlySalesData, len(periods))
	for i := range periods {
		period := periods[i]
		periodMap[period.PeriodKey] = &periods[i]
	}

	type aggregateRow struct {
		PeriodStart  time.Time
		TotalRevenue int64
		Count        int
	}
	dealPeriodExpr := fmt.Sprintf("DATE_TRUNC('%s', COALESCE(deals.actual_close_date, deals.created_at))", truncateUnit)
	visitPeriodExpr := fmt.Sprintf("DATE_TRUNC('%s', visit_reports.visit_date)", truncateUnit)
	taskPeriodExpr := fmt.Sprintf("DATE_TRUNC('%s', tasks.created_at)", truncateUnit)

	dealRows := make([]aggregateRow, 0)
	dealQuery := r.db.Table("deals").
		Select(dealPeriodExpr + " as period_start, COALESCE(SUM(deals.value), 0) as total_revenue, COUNT(deals.id) as count").
		Where("deals.status = 'won' AND deals.deleted_at IS NULL")
	dealQuery = restrictDealAssignedToSalesRole(dealQuery, "deals")
	if len(scopedUserIDs) > 0 {
		dealQuery = dealQuery.Where("deals.assigned_to IN ?", scopedUserIDs)
	}
	dealQuery = dealQuery.Where(
		"(COALESCE(deals.actual_close_date, deals.created_at) >= ? AND COALESCE(deals.actual_close_date, deals.created_at) <= ?)",
		normalizedStart,
		normalizedEnd,
	)
	if err := dealQuery.Group(dealPeriodExpr).Order("period_start ASC").Scan(&dealRows).Error; err != nil {
		return nil, err
	}

	visitRows := make([]aggregateRow, 0)
	visitQuery := r.db.Table("visit_reports").
		Select(visitPeriodExpr+" as period_start, 0 as total_revenue, COUNT(visit_reports.id) as count").
		Where("visit_reports.status IN ?", []string{"completed", "approved"}).
		Where("visit_reports.visit_date >= ? AND visit_reports.visit_date <= ?", normalizedStart, normalizedEnd)
	visitQuery = restrictSalesRepToSalesRole(visitQuery, "visit_reports.sales_rep_id")
	if len(scopedUserIDs) > 0 {
		visitQuery = visitQuery.Where("visit_reports.sales_rep_id IN ?", scopedUserIDs)
	}
	if err := visitQuery.Group(visitPeriodExpr).Order("period_start ASC").Scan(&visitRows).Error; err != nil {
		return nil, err
	}

	taskRows := make([]aggregateRow, 0)
	taskQuery := r.db.Table("tasks").
		Select(taskPeriodExpr+" as period_start, 0 as total_revenue, COUNT(tasks.id) as count").
		Where("tasks.status = 'completed'").
		Where("tasks.created_at >= ? AND tasks.created_at <= ?", normalizedStart, normalizedEnd)
	taskQuery = restrictSalesRepToSalesRole(taskQuery, "tasks.assigned_to")
	if len(scopedUserIDs) > 0 {
		taskQuery = taskQuery.Where("tasks.assigned_to IN ?", scopedUserIDs)
	}
	if err := taskQuery.Group(taskPeriodExpr).Order("period_start ASC").Scan(&taskRows).Error; err != nil {
		return nil, err
	}

	resolvePeriodKey := func(t time.Time) string {
		switch trendMode {
		case "qoq":
			quarter := ((int(t.Month()) - 1) / 3) + 1
			return fmt.Sprintf("%d-Q%d", t.Year(), quarter)
		case "rolling_30d", "rolling_90d":
			return t.Format("2006-01-02")
		default:
			return t.Format("2006-01")
		}
	}

	for _, row := range dealRows {
		if period, ok := periodMap[resolvePeriodKey(row.PeriodStart)]; ok {
			period.TotalRevenue = row.TotalRevenue
			period.TotalDeals = row.Count
		}
	}
	for _, row := range visitRows {
		if period, ok := periodMap[resolvePeriodKey(row.PeriodStart)]; ok {
			period.TotalVisits = row.Count
		}
	}
	for _, row := range taskRows {
		if period, ok := periodMap[resolvePeriodKey(row.PeriodStart)]; ok {
			period.TotalTasks = row.Count
		}
	}

	sort.Slice(periods, func(i, j int) bool {
		return periods[i].PeriodStart.Before(periods[j].PeriodStart)
	})

	var totalRevenue int64
	var totalDeals, totalVisits, totalTasks int
	for i := range periods {
		totalRevenue += periods[i].TotalRevenue
		totalDeals += periods[i].TotalDeals
		totalVisits += periods[i].TotalVisits
		totalTasks += periods[i].TotalTasks
		if i == 0 {
			continue
		}

		var previousValue float64
		switch {
		case periods[i-1].TotalRevenue > 0:
			previousValue = float64(periods[i-1].TotalRevenue)
			periods[i].ChangeRate = ((float64(periods[i].TotalRevenue) - previousValue) / previousValue) * 100
		case periods[i-1].TotalDeals > 0:
			previousValue = float64(periods[i-1].TotalDeals)
			periods[i].ChangeRate = ((float64(periods[i].TotalDeals) - previousValue) / previousValue) * 100
		default:
			periods[i].ChangeRate = 0
		}
	}

	return &sales_overview.MonthlySalesOverviewResponse{
		TrendMode:    trendMode,
		MonthlyData:  periods,
		TotalRevenue: totalRevenue,
		TotalDeals:   totalDeals,
		TotalVisits:  totalVisits,
		TotalTasks:   totalTasks,
	}, nil
}

func (r *repository) GetFunnelDiagnostics(req *sales_overview.GetFunnelDiagnosticsRequest, scopedUserIDs []string, stalledThresholdDays, noActivityThresholdDays, limit int) (*sales_overview.FunnelDiagnosticsResponse, error) {
	if limit <= 0 {
		limit = 20
	}

	now := time.Now()
	selectedSalesUserID := ""
	selectedStageID := ""
	if req != nil {
		selectedSalesUserID = req.SalesUserID
		selectedStageID = req.StageID
	}

	type stalledRow struct {
		ID                string
		Title             string
		AccountName       string
		AssignedToID      string
		AssignedToName    string
		StageID           string
		StageName         string
		Value             int64
		Probability       int
		ExpectedCloseDate *time.Time
		LastStageChangeAt time.Time
		DaysInStage       int
	}

	stalledRows := make([]stalledRow, 0)
	stalledQuery := r.db.Table("deals d").
		Select(`
			d.id,
			d.title,
			COALESCE(a.name, '') as account_name,
			COALESCE(d.assigned_to::text, '') as assigned_to_id,
			COALESCE(u.name, '') as assigned_to_name,
			d.stage_id,
			COALESCE(ps.name, '') as stage_name,
			COALESCE(d.value, 0) as value,
			COALESCE(d.probability, 0) as probability,
			d.expected_close_date,
			COALESCE(last_history.changed_at, d.created_at) as last_stage_change_at,
			GREATEST(0, FLOOR(EXTRACT(EPOCH FROM (NOW() - COALESCE(last_history.changed_at, d.created_at))) / 86400))::int as days_in_stage
		`).
		Joins("LEFT JOIN accounts a ON a.id = d.account_id AND a.deleted_at IS NULL").
		Joins("LEFT JOIN pipeline_stages ps ON ps.id = d.stage_id AND ps.deleted_at IS NULL").
		Joins("LEFT JOIN users u ON u.id = d.assigned_to AND u.deleted_at IS NULL").
		Joins(`
			LEFT JOIN (
				SELECT DISTINCT ON (deal_id) deal_id, changed_at
				FROM deal_histories
				WHERE deleted_at IS NULL
				ORDER BY deal_id, changed_at DESC
			) last_history ON last_history.deal_id = d.id
		`).
		Where("d.deleted_at IS NULL AND d.status = 'open'")
	stalledQuery = restrictDealAssignedToSalesRole(stalledQuery, "d")
	if len(scopedUserIDs) > 0 {
		stalledQuery = stalledQuery.Where("d.assigned_to IN ?", scopedUserIDs)
	}
	if selectedSalesUserID != "" {
		stalledQuery = stalledQuery.Where("d.assigned_to = ?", selectedSalesUserID)
	}
	if selectedStageID != "" {
		stalledQuery = stalledQuery.Where("d.stage_id = ?", selectedStageID)
	}
	if err := stalledQuery.
		Where("COALESCE(last_history.changed_at, d.created_at) <= ?", now.AddDate(0, 0, -stalledThresholdDays)).
		Order("days_in_stage DESC, d.updated_at DESC").
		Limit(limit).
		Scan(&stalledRows).Error; err != nil {
		return nil, err
	}

	type noActivityRow struct {
		ID                  string
		Title               string
		AccountName         string
		AssignedToID        string
		AssignedToName      string
		StageID             string
		StageName           string
		Value               int64
		Probability         int
		ExpectedCloseDate   *time.Time
		LastActivityAt      *time.Time
		DaysWithoutActivity int
	}

	noActivityRows := make([]noActivityRow, 0)
	noActivityQuery := r.db.Table("deals d").
		Select(`
			d.id,
			d.title,
			COALESCE(a.name, '') as account_name,
			COALESCE(d.assigned_to::text, '') as assigned_to_id,
			COALESCE(u.name, '') as assigned_to_name,
			d.stage_id,
			COALESCE(ps.name, '') as stage_name,
			COALESCE(d.value, 0) as value,
			COALESCE(d.probability, 0) as probability,
			d.expected_close_date,
			COALESCE(deal_last_activity.last_activity_at, account_last_activity.last_activity_at) as last_activity_at,
			GREATEST(0, FLOOR(EXTRACT(EPOCH FROM (NOW() - COALESCE(deal_last_activity.last_activity_at, account_last_activity.last_activity_at, d.created_at))) / 86400))::int as days_without_activity
		`).
		Joins("LEFT JOIN accounts a ON a.id = d.account_id AND a.deleted_at IS NULL").
		Joins("LEFT JOIN pipeline_stages ps ON ps.id = d.stage_id AND ps.deleted_at IS NULL").
		Joins("LEFT JOIN users u ON u.id = d.assigned_to AND u.deleted_at IS NULL").
		Joins(`
			LEFT JOIN (
				SELECT deal_id, MAX(timestamp) as last_activity_at
				FROM activities
				WHERE deleted_at IS NULL AND deal_id IS NOT NULL
				GROUP BY deal_id
			) deal_last_activity ON deal_last_activity.deal_id = d.id
		`).
		Joins(`
			LEFT JOIN (
				SELECT account_id, MAX(timestamp) as last_activity_at
				FROM activities
				WHERE deleted_at IS NULL AND account_id IS NOT NULL
				GROUP BY account_id
			) account_last_activity ON account_last_activity.account_id = d.account_id
		`).
		Where("d.deleted_at IS NULL AND d.status = 'open'")
	noActivityQuery = restrictDealAssignedToSalesRole(noActivityQuery, "d")
	if len(scopedUserIDs) > 0 {
		noActivityQuery = noActivityQuery.Where("d.assigned_to IN ?", scopedUserIDs)
	}
	if selectedSalesUserID != "" {
		noActivityQuery = noActivityQuery.Where("d.assigned_to = ?", selectedSalesUserID)
	}
	if selectedStageID != "" {
		noActivityQuery = noActivityQuery.Where("d.stage_id = ?", selectedStageID)
	}
	if err := noActivityQuery.
		Where("COALESCE(deal_last_activity.last_activity_at, account_last_activity.last_activity_at, d.created_at) <= ?", now.AddDate(0, 0, -noActivityThresholdDays)).
		Order("days_without_activity DESC, d.updated_at DESC").
		Limit(limit).
		Scan(&noActivityRows).Error; err != nil {
		return nil, err
	}

	type stageAgingRow struct {
		FromStageName string
		ToStageName   string
		AverageDays   float64
		MedianDays    float64
		Transitions   int
	}

	stageAgingRows := make([]stageAgingRow, 0)
	stageAgingQuery := r.db.Table("deal_histories dh").
		Select(`
			COALESCE(NULLIF(TRIM(dh.from_stage_name), ''), 'Created') as from_stage_name,
			COALESCE(NULLIF(TRIM(dh.to_stage_name), ''), 'Unknown') as to_stage_name,
			COALESCE(AVG(dh.days_in_prev_stage), 0) as average_days,
			COALESCE(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY dh.days_in_prev_stage), 0) as median_days,
			COUNT(*) as transitions
		`).
		Joins("INNER JOIN deals d ON d.id = dh.deal_id AND d.deleted_at IS NULL").
		Where("dh.deleted_at IS NULL AND dh.days_in_prev_stage IS NOT NULL")
	stageAgingQuery = restrictDealAssignedToSalesRole(stageAgingQuery, "d")
	if len(scopedUserIDs) > 0 {
		stageAgingQuery = stageAgingQuery.Where("d.assigned_to IN ?", scopedUserIDs)
	}
	if selectedSalesUserID != "" {
		stageAgingQuery = stageAgingQuery.Where("d.assigned_to = ?", selectedSalesUserID)
	}
	if selectedStageID != "" {
		stageAgingQuery = stageAgingQuery.Where("(d.stage_id = ? OR dh.to_stage_id = ?)", selectedStageID, selectedStageID)
	}
	if err := stageAgingQuery.
		Group("COALESCE(NULLIF(TRIM(dh.from_stage_name), ''), 'Created'), COALESCE(NULLIF(TRIM(dh.to_stage_name), ''), 'Unknown')").
		Order("average_days DESC, transitions DESC").
		Limit(12).
		Scan(&stageAgingRows).Error; err != nil {
		return nil, err
	}

	type optionRow struct {
		ID   string
		Name string
	}

	salesRepOptionRows := make([]optionRow, 0)
	salesRepOptionsQuery := r.db.Table("deals d").
		Select("DISTINCT u.id::text as id, u.name").
		Joins("INNER JOIN users u ON u.id = d.assigned_to AND u.deleted_at IS NULL").
		Where("d.deleted_at IS NULL AND d.status = 'open'")
	salesRepOptionsQuery = restrictDealAssignedToSalesRole(salesRepOptionsQuery, "d")
	if len(scopedUserIDs) > 0 {
		salesRepOptionsQuery = salesRepOptionsQuery.Where("d.assigned_to IN ?", scopedUserIDs)
	}
	if err := salesRepOptionsQuery.Order("u.name ASC").Scan(&salesRepOptionRows).Error; err != nil {
		return nil, err
	}

	stageOptionRows := make([]optionRow, 0)
	stageOptionsQuery := r.db.Table("deals d").
		Select("DISTINCT ps.id::text as id, ps.name").
		Joins("INNER JOIN pipeline_stages ps ON ps.id = d.stage_id AND ps.deleted_at IS NULL").
		Where("d.deleted_at IS NULL AND d.status = 'open'")
	stageOptionsQuery = restrictDealAssignedToSalesRole(stageOptionsQuery, "d")
	if len(scopedUserIDs) > 0 {
		stageOptionsQuery = stageOptionsQuery.Where("d.assigned_to IN ?", scopedUserIDs)
	}
	if selectedSalesUserID != "" {
		stageOptionsQuery = stageOptionsQuery.Where("d.assigned_to = ?", selectedSalesUserID)
	}
	if err := stageOptionsQuery.Order("ps.name ASC").Scan(&stageOptionRows).Error; err != nil {
		return nil, err
	}

	response := &sales_overview.FunnelDiagnosticsResponse{
		GeneratedAt:             now,
		StalledThresholdDays:    stalledThresholdDays,
		NoActivityThresholdDays: noActivityThresholdDays,
		SelectedSalesUserID:     selectedSalesUserID,
		SelectedStageID:         selectedStageID,
		AvailableSalesReps:      make([]sales_overview.FunnelDiagnosticsSalesRepOption, len(salesRepOptionRows)),
		AvailableStages:         make([]sales_overview.FunnelDiagnosticsStageOption, len(stageOptionRows)),
		StalledDeals:            make([]sales_overview.StalledDealItem, len(stalledRows)),
		NoActivityDeals:         make([]sales_overview.NoActivityDealItem, len(noActivityRows)),
		StageAging:              make([]sales_overview.StageAgingItem, len(stageAgingRows)),
	}

	for i, row := range salesRepOptionRows {
		response.AvailableSalesReps[i] = sales_overview.FunnelDiagnosticsSalesRepOption{
			ID:   row.ID,
			Name: row.Name,
		}
	}

	for i, row := range stageOptionRows {
		response.AvailableStages[i] = sales_overview.FunnelDiagnosticsStageOption{
			ID:   row.ID,
			Name: row.Name,
		}
	}

	for i, row := range stalledRows {
		response.StalledDeals[i] = sales_overview.StalledDealItem{
			ID:                row.ID,
			Title:             row.Title,
			AccountName:       row.AccountName,
			AssignedToID:      row.AssignedToID,
			AssignedToName:    row.AssignedToName,
			StageID:           row.StageID,
			StageName:         row.StageName,
			Value:             row.Value,
			ValueFormatted:    formatCurrency(row.Value),
			Probability:       row.Probability,
			ExpectedCloseDate: row.ExpectedCloseDate,
			LastStageChangeAt: row.LastStageChangeAt,
			DaysInStage:       row.DaysInStage,
		}
	}

	for i, row := range noActivityRows {
		response.NoActivityDeals[i] = sales_overview.NoActivityDealItem{
			ID:                  row.ID,
			Title:               row.Title,
			AccountName:         row.AccountName,
			AssignedToID:        row.AssignedToID,
			AssignedToName:      row.AssignedToName,
			StageID:             row.StageID,
			StageName:           row.StageName,
			Value:               row.Value,
			ValueFormatted:      formatCurrency(row.Value),
			Probability:         row.Probability,
			ExpectedCloseDate:   row.ExpectedCloseDate,
			LastActivityAt:      row.LastActivityAt,
			DaysWithoutActivity: row.DaysWithoutActivity,
		}
	}

	for i, row := range stageAgingRows {
		response.StageAging[i] = sales_overview.StageAgingItem{
			FromStageName: row.FromStageName,
			ToStageName:   row.ToStageName,
			TransitionKey: strings.ToLower(strings.ReplaceAll(row.FromStageName+"-"+row.ToStageName, " ", "_")),
			AverageDays:   row.AverageDays,
			MedianDays:    row.MedianDays,
			Transitions:   row.Transitions,
		}
	}

	response.Summary = sales_overview.FunnelDiagnosticsSummary{
		StalledDeals:          len(response.StalledDeals),
		NoActivityDeals:       len(response.NoActivityDeals),
		StageAgingTransitions: len(response.StageAging),
	}

	return response, nil
}

// getPerformanceSummary gets performance summary for a user
func (r *repository) getPerformanceSummary(userID string, startDate, endDate interface{}) (*sales_overview.SalesPerformanceListResponse, error) {
	var totalRevenue int64
	var dealsClosed, visitsCompleted, tasksCompleted int64

	// Calculate revenue and deals closed
	dealsQuery := r.db.Table("deals").
		Where("assigned_to = ? AND status = ?", userID, "won")
	if startDate != nil {
		dealsQuery = dealsQuery.Where("actual_close_date >= ? OR (actual_close_date IS NULL AND created_at >= ?)", startDate, startDate)
	}
	if endDate != nil {
		dealsQuery = dealsQuery.Where("actual_close_date <= ? OR (actual_close_date IS NULL AND created_at <= ?)", endDate, endDate)
	}

	var revenueResult struct {
		Total int64
	}
	if err := dealsQuery.Select("COALESCE(SUM(value), 0) as total").Scan(&revenueResult).Error; err != nil {
		return nil, err
	}
	totalRevenue = revenueResult.Total

	if err := dealsQuery.Count(&dealsClosed).Error; err != nil {
		return nil, err
	}

	// Count visits completed
	visitsQuery := r.db.Table("visit_reports").
		Where("sales_rep_id = ? AND status IN ?", userID, []string{"completed", "approved"})
	if startDate != nil {
		visitsQuery = visitsQuery.Where("visit_date >= ?", startDate)
	}
	if endDate != nil {
		visitsQuery = visitsQuery.Where("visit_date <= ?", endDate)
	}
	if err := visitsQuery.Count(&visitsCompleted).Error; err != nil {
		return nil, err
	}

	// Count tasks completed
	tasksQuery := r.db.Table("tasks").
		Where("assigned_to = ? AND status = ?", userID, "completed")
	if startDate != nil {
		tasksQuery = tasksQuery.Where("created_at >= ?", startDate)
	}
	if endDate != nil {
		tasksQuery = tasksQuery.Where("created_at <= ?", endDate)
	}
	if err := tasksQuery.Count(&tasksCompleted).Error; err != nil {
		return nil, err
	}

	// Calculate conversion rate
	var totalDeals int64
	totalDealsQuery := r.db.Table("deals").Where("assigned_to = ?", userID)
	if startDate != nil {
		totalDealsQuery = totalDealsQuery.Where("created_at >= ?", startDate)
	}
	if endDate != nil {
		totalDealsQuery = totalDealsQuery.Where("created_at <= ?", endDate)
	}
	if err := totalDealsQuery.Count(&totalDeals).Error; err != nil {
		return nil, err
	}

	conversionRate := 0.0
	if totalDeals > 0 {
		conversionRate = (float64(dealsClosed) / float64(totalDeals)) * 100.0
	}

	return &sales_overview.SalesPerformanceListResponse{
		TotalRevenue:          totalRevenue,
		TotalRevenueFormatted: formatCurrency(totalRevenue),
		DealsClosed:           int(dealsClosed),
		VisitsCompleted:       int(visitsCompleted),
		TasksCompleted:        int(tasksCompleted),
		ConversionRate:        conversionRate,
	}, nil
}

// GetSalesRepCheckInLocations gets check-in locations for sales rep (ordered by visit number) with pagination
// Optimized for enterprise-scale: uses database-level pagination and batch account loading
func (r *repository) GetSalesRepCheckInLocations(userID string, req *sales_overview.GetSalesRepCheckInLocationsRequest, startDate, endDate interface{}) ([]sales_overview.CheckInLocation, int64, error) {
	// Build base query for filtering
	baseQuery := r.db.Table("visit_reports").
		Where("sales_rep_id = ? AND check_in_location IS NOT NULL", userID)

	// If no date range provided, use last 30 days as default
	if startDate == nil && endDate == nil {
		now := time.Now()
		endDate = now
		startDate = now.AddDate(0, 0, -30)
	}

	if startDate != nil {
		baseQuery = baseQuery.Where("visit_date >= ?", startDate)
	}
	if endDate != nil {
		baseQuery = baseQuery.Where("visit_date <= ?", endDate)
	}

	// Get total count first (efficient COUNT query)
	var totalCount int64
	if err := baseQuery.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination parameters
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 50 // Default to 50 for locations (larger than activities)
	}
	if perPage > 100 {
		perPage = 100
	}
	offset := (page - 1) * perPage

	// Query visit reports with pagination at database level (optimized for enterprise-scale)
	var visits []struct {
		ID              string
		VisitDate       time.Time
		CheckInTime     *time.Time
		CheckInLocation datatypes.JSON
		Purpose         string
		AccountID       *string
	}

	query := baseQuery.
		Select("id, visit_date, check_in_time, check_in_location, purpose, account_id").
		Order("visit_date ASC, check_in_time ASC").
		Limit(perPage).
		Offset(offset)

	if err := query.Find(&visits).Error; err != nil {
		return nil, 0, err
	}

	// Collect unique account IDs for batch loading (optimize N+1 query problem)
	accountIDSet := make(map[string]bool)
	accountIDs := make([]string, 0)
	for _, v := range visits {
		if v.AccountID != nil && !accountIDSet[*v.AccountID] {
			accountIDSet[*v.AccountID] = true
			accountIDs = append(accountIDs, *v.AccountID)
		}
	}

	// Batch load all accounts in one query
	accountsMap := make(map[string]*account.Account)
	if len(accountIDs) > 0 {
		var accounts []account.Account
		if err := r.db.Where("id IN ?", accountIDs).Find(&accounts).Error; err == nil {
			for i := range accounts {
				accountsMap[accounts[i].ID] = &accounts[i]
			}
		}
	}

	// Calculate base visit number (accounting for pagination offset)
	baseVisitNumber := int64(offset + 1)

	locations := make([]sales_overview.CheckInLocation, 0, len(visits))
	for i, v := range visits {
		// Parse check-in location JSON
		if len(v.CheckInLocation) > 0 {
			var locationData map[string]interface{}

			// datatypes.JSON is []byte, so we can unmarshal directly
			if err := json.Unmarshal(v.CheckInLocation, &locationData); err != nil {
				continue
			}

			if locationData != nil {
				var loc *sales_overview.Location
				if lat, ok := locationData["latitude"].(float64); ok {
					if lng, ok := locationData["longitude"].(float64); ok {
						loc = &sales_overview.Location{
							Latitude:  lat,
							Longitude: lng,
						}
						if addr, ok := locationData["address"].(string); ok {
							loc.Address = addr
						}
					}
				}

				// Only append if location was successfully parsed
				if loc != nil {
					checkInTime := v.VisitDate
					if v.CheckInTime != nil {
						checkInTime = *v.CheckInTime
					}

					location := sales_overview.CheckInLocation{
						VisitNumber:   int(baseVisitNumber + int64(i)), // 1-based indexing with pagination offset
						VisitReportID: v.ID,
						VisitDate:     v.VisitDate,
						CheckInTime:   checkInTime,
						Location:      loc,
						Purpose:       v.Purpose,
					}

					// Set account from batch-loaded map
					if v.AccountID != nil {
						if acc, exists := accountsMap[*v.AccountID]; exists {
							location.Account = &sales_overview.AccountRef{
								ID:   acc.ID,
								Name: acc.Name,
							}
						}
					}

					locations = append(locations, location)
				}
			}
		}
	}

	return locations, totalCount, nil
}
