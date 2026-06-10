package deal

import (
	"fmt"
	"strings"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

const (
	queryWhereAssignedTo   = "assigned_to = ?"
	queryWhereStatus       = "status = ?"
	queryWhereID           = "id = ?"
	queryWhereCreatedAtGte = "created_at >= ?"
	queryWhereCreatedAtLte = "created_at <= ?"
	queryWhereDateGte      = "(actual_close_date >= ? OR (actual_close_date IS NULL AND created_at >= ?))"
	queryWhereDateLte      = "(actual_close_date <= ? OR (actual_close_date IS NULL AND created_at <= ?))"
	dateFormatISO          = "2006-01-02"
	exprSumValue           = "COALESCE(SUM(value), 0)"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new deal repository
func NewRepository(db *gorm.DB) interfaces.DealRepository {
	return &repository{db: db}
}

func (r *repository) FindByID(id string) (*pipeline.Deal, error) {
	var deal pipeline.Deal
	err := r.db.
		Preload("Account").
		Preload("Contact").
		Preload("Stage").
		Preload("ProductItems").
		Preload("AssignedUser").
		Where(queryWhereID, id).
		First(&deal).Error
	if err != nil {
		return nil, err
	}
	return &deal, nil
}

func (r *repository) applyListFilters(query *gorm.DB, req *pipeline.ListDealsRequest) *gorm.DB {
	// Apply RBAC scope filtering
	if len(req.ScopedUserIDs) > 0 {
		query = query.Where("assigned_to IN ?", req.ScopedUserIDs)
	}

	if req.Search != "" {
		searchQuery := strings.Join(strings.Fields(req.Search), " & ") + ":*"
		query = query.Where("to_tsvector('english', title || ' ' || COALESCE(description, '') || ' ' || COALESCE(notes, '')) @@ to_tsquery('english', ?)", searchQuery)
	}
	if req.StageID != "" {
		query = query.Where("stage_id = ?", req.StageID)
	}
	if req.AccountID != "" {
		query = query.Where("account_id = ?", req.AccountID)
	}
	if req.AssignedTo != "" {
		query = query.Where(queryWhereAssignedTo, req.AssignedTo)
	}
	if req.BrickID != "" {
		query = query.Where("brick_id = ?", req.BrickID)
	}
	if req.Status != "" {
		query = query.Where(queryWhereStatus, req.Status)
	}
	if req.Source != "" {
		query = query.Where("source = ?", req.Source)
	}
	if req.MinValue != nil {
		query = query.Where("value >= ?", *req.MinValue)
	}
	if req.MaxValue != nil {
		query = query.Where("value <= ?", *req.MaxValue)
	}
	if req.DateFrom != "" {
		if from, err := time.Parse(dateFormatISO, req.DateFrom); err == nil {
			query = query.Where(queryWhereCreatedAtGte, from)
		}
	}
	if req.DateTo != "" {
		if to, err := time.Parse(dateFormatISO, req.DateTo); err == nil {
			to = time.Date(to.Year(), to.Month(), to.Day(), 23, 59, 59, 999999999, to.Location())
			query = query.Where(queryWhereCreatedAtLte, to)
		}
	}
	return query
}

func (r *repository) List(req *pipeline.ListDealsRequest) ([]pipeline.Deal, int64, error) {
	var deals []pipeline.Deal
	var total int64

	query := r.db.Model(&pipeline.Deal{})
	query = r.applyListFilters(query, req)

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

	err := query.
		Preload("Account").
		Preload("Contact").
		Preload("Stage").
		Preload("ProductItems").
		Preload("AssignedUser").
		Order("created_at DESC").
		Offset(offset).
		Limit(perPage).
		Find(&deals).Error
	if err != nil {
		return nil, 0, err
	}

	return deals, total, nil
}

func (r *repository) Create(deal *pipeline.Deal) error {
	return r.db.Create(deal).Error
}

func (r *repository) Update(deal *pipeline.Deal) error {
	updates := map[string]interface{}{
		"title":                  deal.Title,
		"description":            deal.Description,
		"account_id":             deal.AccountID,
		"contact_id":             deal.ContactID,
		"stage_id":               deal.StageID,
		"value":                  deal.Value,
		"probability":            deal.Probability,
		"expected_close_date":    deal.ExpectedCloseDate,
		"actual_close_date":      deal.ActualCloseDate,
		"assigned_to":            deal.AssignedTo,
		"lead_id":                deal.LeadID,
		"brick_id":               deal.BrickID,
		"status":                 deal.Status,
		"source":                 deal.Source,
		"budget_confirmed":       deal.BudgetConfirmed,
		"authority_confirmed":    deal.AuthorityConfirmed,
		"need_confirmed":         deal.NeedConfirmed,
		"timeline_confirmed":     deal.TimelineConfirmed,
		"qualification_snapshot": deal.QualificationSnapshot,
		"close_reason":           deal.CloseReason,
		"notes":                  deal.Notes,
		"created_by":             deal.CreatedBy,
	}

	return r.db.Model(&pipeline.Deal{}).Where(queryWhereID, deal.ID).Updates(updates).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Where(queryWhereID, id).Delete(&pipeline.Deal{}).Error
}

func (r *repository) getStatusStats(baseQuery *gorm.DB, status string) (count, value int64, err error) {
	if err = baseQuery.Where(queryWhereStatus, status).Count(&count).Error; err != nil {
		return 0, 0, err
	}
	if err = baseQuery.Where(queryWhereStatus, status).Select(exprSumValue).Scan(&value).Error; err != nil {
		return 0, 0, err
	}
	return count, value, nil
}

func (r *repository) getStatsForQuery(baseQuery *gorm.DB) (dealsCount, dealsValue int64, wonDeals, wonValue int64, lostDeals, lostValue int64, openDeals, openValue int64, err error) {
	if err = baseQuery.Count(&dealsCount).Error; err != nil {
		return
	}
	if err = baseQuery.Select(exprSumValue).Scan(&dealsValue).Error; err != nil {
		return
	}

	if wonDeals, wonValue, err = r.getStatusStats(baseQuery, "won"); err != nil {
		return
	}
	if lostDeals, lostValue, err = r.getStatusStats(baseQuery, "lost"); err != nil {
		return
	}
	if openDeals, openValue, err = r.getStatusStats(baseQuery, "open"); err != nil {
		return
	}

	return
}

func (r *repository) GetSummary() (*pipeline.PipelineSummaryResponse, error) {
	baseQuery := r.db.Model(&pipeline.Deal{})
	totalDeals, totalValue, wonDeals, wonValue, lostDeals, lostValue, openDeals, openValue, err := r.getStatsForQuery(baseQuery)
	if err != nil {
		return nil, err
	}

	// Get summary by stage
	var stageSummaries []pipeline.StageSummary
	err = r.db.Model(&pipeline.Deal{}).
		Select(`
			stage_id,
			COUNT(*) as deal_count,
			COALESCE(SUM(value), 0) as total_value
		`).
		Group("stage_id").
		Scan(&stageSummaries).Error
	if err != nil {
		return nil, err
	}

	// OPTIMIZED: Batch load all stages in one query instead of N+1 queries
	if len(stageSummaries) > 0 {
		// Collect all stage IDs
		stageIDs := make([]string, 0, len(stageSummaries))
		stageIDMap := make(map[string]int) // Map stage ID to index in stageSummaries
		for i, ss := range stageSummaries {
			if ss.StageID != "" {
				stageIDs = append(stageIDs, ss.StageID)
				stageIDMap[ss.StageID] = i
			}
		}

		// Batch load all stages in one query
		var stages []pipeline.PipelineStage
		if len(stageIDs) > 0 {
			if err := r.db.Where("id IN ?", stageIDs).Find(&stages).Error; err == nil {
				// Map stages to summaries
				stageMap := make(map[string]*pipeline.PipelineStage)
				for i := range stages {
					stageMap[stages[i].ID] = &stages[i]
				}

				// Populate stage names and codes
				for i := range stageSummaries {
					if stage, exists := stageMap[stageSummaries[i].StageID]; exists {
						stageSummaries[i].StageName = stage.Name
						stageSummaries[i].StageCode = stage.Code
					}
				}
			}
		}
	}

	// Format values
	for i := range stageSummaries {
		stageSummaries[i].TotalValueFormatted = formatCurrency(stageSummaries[i].TotalValue)
	}

	summary := &pipeline.PipelineSummaryResponse{
		TotalDeals:          totalDeals,
		TotalValue:          totalValue,
		TotalValueFormatted: formatCurrency(totalValue),
		WonDeals:            wonDeals,
		WonValue:            wonValue,
		WonValueFormatted:   formatCurrency(wonValue),
		LostDeals:           lostDeals,
		LostValue:           lostValue,
		LostValueFormatted:  formatCurrency(lostValue),
		OpenDeals:           openDeals,
		OpenValue:           openValue,
		OpenValueFormatted:  formatCurrency(openValue),
		ByStage:             stageSummaries,
	}

	return summary, nil
}

func (r *repository) mapToForecastDeal(deal pipeline.Deal) pipeline.ForecastDeal {
	weightedValue := deal.Value * int64(deal.Probability) / 100

	accountName := ""
	if deal.Account != nil {
		accountName = deal.Account.Name
	}

	contactName := ""
	if deal.Contact != nil {
		contactName = deal.Contact.Name
	}

	stageName := ""
	if deal.Stage != nil {
		stageName = deal.Stage.Name
	}

	contactID := ""
	if deal.ContactID != nil {
		contactID = *deal.ContactID
	}

	return pipeline.ForecastDeal{
		ID:                     deal.ID,
		Title:                  deal.Title,
		AccountID:              deal.AccountID,
		AccountName:            accountName,
		ContactID:              contactID,
		ContactName:            contactName,
		StageName:              stageName,
		Value:                  deal.Value,
		ValueFormatted:         formatCurrency(deal.Value),
		Probability:            deal.Probability,
		WeightedValue:          weightedValue,
		WeightedValueFormatted: formatCurrency(weightedValue),
		ExpectedCloseDate:      deal.ExpectedCloseDate,
	}
}

func (r *repository) GetForecast(periodType string, start, end time.Time) (*pipeline.ForecastResponse, error) {
	var deals []pipeline.Deal
	err := r.db.
		Preload("Account").
		Preload("Contact").
		Preload("Stage").
		Where("expected_close_date >= ? AND expected_close_date <= ?", start, end).
		Where(queryWhereStatus, "open").
		Find(&deals).Error
	if err != nil {
		return nil, err
	}

	var expectedRevenue, weightedRevenue int64
	forecastDeals := make([]pipeline.ForecastDeal, 0, len(deals))

	for _, deal := range deals {
		mapped := r.mapToForecastDeal(deal)
		expectedRevenue += mapped.Value
		weightedRevenue += mapped.WeightedValue
		forecastDeals = append(forecastDeals, mapped)
	}

	forecast := &pipeline.ForecastResponse{
		Period: pipeline.ForecastPeriod{
			Type:  periodType,
			Start: start,
			End:   end,
		},
		ExpectedRevenue:          expectedRevenue,
		ExpectedRevenueFormatted: formatCurrency(expectedRevenue),
		WeightedRevenue:          weightedRevenue,
		WeightedRevenueFormatted: formatCurrency(weightedRevenue),
		Deals:                    forecastDeals,
	}

	return forecast, nil
}

// formatCurrency formats integer (sen) to formatted currency string
func formatCurrency(amount int64) string {
	// Convert to Rupiah (divide by 100 if stored in sen)
	rupiah := float64(amount) / 100.0
	// Format with thousand separator
	formatted := formatNumber(rupiah)
	return "Rp " + formatted
}

// formatNumber formats number with thousand separator
func formatNumber(n float64) string {
	// Convert to int64 to remove decimal places
	amount := int64(n)

	// Handle zero case
	if amount == 0 {
		return "0"
	}

	// Handle negative numbers
	negative := false
	if amount < 0 {
		negative = true
		amount = -amount
	}

	// Convert to string
	str := fmt.Sprintf("%d", amount)
	length := len(str)

	// Add thousand separators (dot for Indonesian format)
	// Split into chunks of 3 digits from right
	var parts []string
	for i := length; i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		parts = append([]string{str[start:i]}, parts...)
	}

	result := strings.Join(parts, ".")
	if negative {
		result = "-" + result
	}

	return result
}

// GetStatsByStatus returns deal statistics grouped by status using database aggregation
func (r *repository) GetStatsByStatus(startDate, endDate string, assignedTo, stageID, status string) (map[string]int64, error) {
	query := r.db.Table("deals")

	// Apply date filters
	if startDate != "" {
		if start, err := time.Parse(dateFormatISO, startDate); err == nil {
			query = query.Where(queryWhereCreatedAtGte, start)
		}
	}
	if endDate != "" {
		if end, err := time.Parse(dateFormatISO, endDate); err == nil {
			end = end.Add(24 * time.Hour)
			query = query.Where(queryWhereCreatedAtLte, end)
		}
	}

	// Apply other filters
	if assignedTo != "" {
		query = query.Where(queryWhereAssignedTo, assignedTo)
	}
	if stageID != "" {
		query = query.Where("stage_id = ?", stageID)
	}
	if status != "" {
		query = query.Where(queryWhereStatus, status)
	}

	// Aggregate by status
	var results []struct {
		Status string
		Count  int64
	}
	err := query.
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, r := range results {
		stats[r.Status] = r.Count
	}

	return stats, nil
}

// GetStatsByStage returns deal statistics grouped by stage using database aggregation
func (r *repository) GetStatsByStage(startDate, endDate string, assignedTo, status string) (map[string]int64, error) {
	query := r.db.Table("deals")

	// Apply date filters
	if startDate != "" {
		if start, err := time.Parse(dateFormatISO, startDate); err == nil {
			query = query.Where(queryWhereCreatedAtGte, start)
		}
	}
	if endDate != "" {
		if end, err := time.Parse(dateFormatISO, endDate); err == nil {
			end = end.Add(24 * time.Hour)
			query = query.Where(queryWhereCreatedAtLte, end)
		}
	}

	// Apply other filters
	if assignedTo != "" {
		query = query.Where(queryWhereAssignedTo, assignedTo)
	}
	if status != "" {
		query = query.Where(queryWhereStatus, status)
	}

	// Aggregate by stage
	var results []struct {
		StageID string
		Count   int64
	}
	err := query.
		Select("stage_id, COUNT(*) as count").
		Group("stage_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, r := range results {
		stats[r.StageID] = r.Count
	}

	return stats, nil
}

// CountByDateRange returns count of deals created in date range using database aggregation
func (r *repository) CountByDateRange(startDate, endDate interface{}) (int64, error) {
	var count int64
	query := r.db.Table("deals")

	if startDate != nil {
		if start, ok := startDate.(time.Time); ok {
			query = query.Where(queryWhereCreatedAtGte, start)
		}
	}
	if endDate != nil {
		if end, ok := endDate.(time.Time); ok {
			query = query.Where(queryWhereCreatedAtLte, end)
		}
	}

	err := query.Count(&count).Error
	return count, err
}

// GetWonDealsValueInPeriod returns count and total value of won deals closed in date range using database aggregation
func (r *repository) GetWonDealsValueInPeriod(startDate, endDate interface{}) (int64, int64, error) {
	query := r.db.Table("deals").Where(queryWhereStatus, "won")

	if startDate != nil {
		if start, ok := startDate.(time.Time); ok {
			query = query.Where(queryWhereDateGte, start, start)
		}
	}
	if endDate != nil {
		if end, ok := endDate.(time.Time); ok {
			query = query.Where(queryWhereDateLte, end, end)
		}
	}

	var result struct {
		Count int64
		Value int64
	}
	err := query.
		Select("COUNT(*) as count, COALESCE(SUM(value), 0) as value").
		Scan(&result).Error

	if err != nil {
		return 0, 0, err
	}

	return result.Count, result.Value, nil
}

// GetWonDealsValueInPeriodByUser returns count and total value of won deals closed in date range for a specific user using database aggregation
func (r *repository) GetWonDealsValueInPeriodByUser(userID string, startDate, endDate interface{}) (int64, int64, error) {
	query := r.db.Table("deals").Where(queryWhereStatus+" AND "+queryWhereAssignedTo, "won", userID)

	if startDate != nil {
		if start, ok := startDate.(time.Time); ok {
			query = query.Where(queryWhereDateGte, start, start)
		}
	}
	if endDate != nil {
		if end, ok := endDate.(time.Time); ok {
			query = query.Where(queryWhereDateLte, end, end)
		}
	}

	var result struct {
		Count int64
		Value int64
	}
	err := query.
		Select("COUNT(*) as count, COALESCE(SUM(value), 0) as value").
		Scan(&result).Error

	if err != nil {
		return 0, 0, err
	}

	return result.Count, result.Value, nil
}

func (r *repository) loadStagesByID(stageIDs []string) (map[string]*pipeline.PipelineStage, error) {
	var stages []pipeline.PipelineStage
	if err := r.db.Where("id IN ?", stageIDs).Find(&stages).Error; err != nil {
		return nil, err
	}
	stageMap := make(map[string]*pipeline.PipelineStage)
	for i := range stages {
		stageMap[stages[i].ID] = &stages[i]
	}
	return stageMap, nil
}

func (r *repository) populateStageSummaries(stageSummaries []pipeline.StageSummary) {
	if len(stageSummaries) == 0 {
		return
	}

	stageIDs := make([]string, 0, len(stageSummaries))
	for _, ss := range stageSummaries {
		if ss.StageID != "" {
			stageIDs = append(stageIDs, ss.StageID)
		}
	}

	if len(stageIDs) > 0 {
		if stageMap, err := r.loadStagesByID(stageIDs); err == nil {
			for i := range stageSummaries {
				if stage, exists := stageMap[stageSummaries[i].StageID]; exists {
					stageSummaries[i].StageName = stage.Name
					stageSummaries[i].StageCode = stage.Code
				}
			}
		}
	}

	for i := range stageSummaries {
		stageSummaries[i].TotalValueFormatted = formatCurrency(stageSummaries[i].TotalValue)
	}
}

func (r *repository) GetSummaryInPeriod(startDate, endDate interface{}) (*pipeline.PipelineSummaryResponse, error) {
	baseQuery := r.db.Model(&pipeline.Deal{})
	if startDate != nil {
		if start, ok := startDate.(time.Time); ok {
			baseQuery = baseQuery.Where(queryWhereDateGte, start, start)
		}
	}
	if endDate != nil {
		if end, ok := endDate.(time.Time); ok {
			baseQuery = baseQuery.Where(queryWhereDateLte, end, end)
		}
	}

	// Use a clone for stats calculation to avoid modifying baseQuery
	totalDeals, totalValue, wonDeals, wonValue, lostDeals, lostValue, openDeals, openValue, err := r.getStatsForQuery(baseQuery.Session(&gorm.Session{}))
	if err != nil {
		return nil, err
	}

	var stageSummaries []pipeline.StageSummary
	// Use a clone for grouping query
	err = baseQuery.Session(&gorm.Session{}).
		Select(`
			stage_id,
			COUNT(*) as deal_count,
			COALESCE(SUM(value), 0) as total_value
		`).
		Group("stage_id").
		Scan(&stageSummaries).Error
	if err != nil {
		return nil, err
	}

	r.populateStageSummaries(stageSummaries)

	return &pipeline.PipelineSummaryResponse{
		TotalDeals:          totalDeals,
		TotalValue:          totalValue,
		TotalValueFormatted: formatCurrency(totalValue),
		WonDeals:            wonDeals,
		WonValue:            wonValue,
		WonValueFormatted:   formatCurrency(wonValue),
		LostDeals:           lostDeals,
		LostValue:           lostValue,
		LostValueFormatted:  formatCurrency(lostValue),
		OpenDeals:           openDeals,
		OpenValue:           openValue,
		OpenValueFormatted:  formatCurrency(openValue),
		ByStage:             stageSummaries,
	}, nil
}
