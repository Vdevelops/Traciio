package visit_report

import (
	"strings"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/visit_report"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new visit report repository
func NewRepository(db *gorm.DB) interfaces.VisitReportRepository {
	return &repository{db: db}
}

func (r *repository) FindByID(id string) (*visit_report.VisitReport, error) {
	var vr visit_report.VisitReport
	err := r.db.Where("id = ?", id).First(&vr).Error
	if err != nil {
		return nil, err
	}
	return &vr, nil
}

func (r *repository) List(req *visit_report.ListVisitReportsRequest) ([]visit_report.VisitReport, int64, error) {
	var visitReports []visit_report.VisitReport
	var total int64

	query := r.db.Model(&visit_report.VisitReport{})

	// Apply RBAC scope filtering
	if len(req.ScopedUserIDs) > 0 {
		query = query.Where("sales_rep_id IN ?", req.ScopedUserIDs)
	}

	// Apply filters
	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where(
			"LOWER(purpose) LIKE ? OR LOWER(notes) LIKE ?",
			search, search,
		)
	}

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	if req.AccountID != "" {
		query = query.Where("account_id = ?", req.AccountID)
	} else {
		// If AccountID is empty, we need to handle NULL values correctly
		// This allows filtering for visit reports without account (qualification phase)
	}

	if req.DealID != "" {
		query = query.Where("deal_id = ?", req.DealID)
	}

	if req.LeadID != "" {
		query = query.Where("lead_id = ?", req.LeadID)
	}

	if req.SalesRepID != "" {
		query = query.Where("sales_rep_id = ?", req.SalesRepID)
	}

	if req.BrickID != "" {
		query = query.Where("brick_id = ?", req.BrickID)
	}

	// Date range filter
	if req.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err == nil {
			query = query.Where("visit_date >= ?", startDate)
		}
	}

	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err == nil {
			// Add one day to include the end date
			endDate = endDate.Add(24 * time.Hour)
			query = query.Where("visit_date < ?", endDate)
		}
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination - support both page-based and offset-based
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}

	var offset int
	// If offset is provided, use offset-based pagination (for infinity scroll)
	if req.Offset > 0 {
		offset = req.Offset
	} else {
		// Otherwise, use page-based pagination
		page := req.Page
		if page < 1 {
			page = 1
		}
		offset = (page - 1) * perPage
	}

	// Fetch data
	err := query.Order("visit_date DESC, created_at DESC").Offset(offset).Limit(perPage).Find(&visitReports).Error
	if err != nil {
		return nil, 0, err
	}

	return visitReports, total, nil
}

func (r *repository) Create(vr *visit_report.VisitReport) error {
	return r.db.Create(vr).Error
}

func (r *repository) Update(vr *visit_report.VisitReport) error {
	return r.db.Save(vr).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&visit_report.VisitReport{}).Error
}

// UpdateByLeadID updates all visit reports for a lead using batch UPDATE (fixes N+1 bottleneck)
func (r *repository) UpdateByLeadID(leadID string, dealID, accountID *string) error {
	updates := make(map[string]interface{})

	if dealID != nil {
		updates["deal_id"] = *dealID
	}
	if accountID != nil {
		updates["account_id"] = *accountID
	}

	if len(updates) == 0 {
		return nil
	}

	return r.db.Model(&visit_report.VisitReport{}).
		Where("lead_id = ?", leadID).
		Updates(updates).Error
}

func (r *repository) FindByAccountID(accountID string) ([]visit_report.VisitReport, error) {
	var visitReports []visit_report.VisitReport
	err := r.db.Where("account_id = ?", accountID).Order("visit_date DESC").Find(&visitReports).Error
	if err != nil {
		return nil, err
	}
	return visitReports, nil
}

func (r *repository) FindBySalesRepID(salesRepID string) ([]visit_report.VisitReport, error) {
	var visitReports []visit_report.VisitReport
	err := r.db.Where("sales_rep_id = ?", salesRepID).Order("visit_date DESC").Find(&visitReports).Error
	if err != nil {
		return nil, err
	}
	return visitReports, nil
}

// GetStatsByStatus returns visit report statistics grouped by status using database aggregation
func (r *repository) GetStatsByStatus(startDate, endDate string, accountID, salesRepID, status string) (map[string]int64, error) {
	query := r.db.Table("visit_reports")

	// Apply date filters
	if startDate != "" {
		if start, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("visit_date >= ?", start)
		}
	}
	if endDate != "" {
		if end, err := time.Parse("2006-01-02", endDate); err == nil {
			end = end.Add(24 * time.Hour)
			query = query.Where("visit_date < ?", end)
		}
	}

	// Apply other filters
	if accountID != "" {
		query = query.Where("account_id = ?", accountID)
	}
	if salesRepID != "" {
		query = query.Where("sales_rep_id = ?", salesRepID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
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

// GetStatsByStatusForUsers returns visit report statistics for multiple users using batch query (fixes N+1)
func (r *repository) GetStatsByStatusForUsers(startDate, endDate string, userIDs []string) (map[string]int64, error) {
	if len(userIDs) == 0 {
		return make(map[string]int64), nil
	}

	query := r.db.Table("visit_reports")

	// Apply date filters
	if startDate != "" {
		if start, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("visit_date >= ?", start)
		}
	}
	if endDate != "" {
		if end, err := time.Parse("2006-01-02", endDate); err == nil {
			end = end.Add(24 * time.Hour)
			query = query.Where("visit_date < ?", end)
		}
	}

	// Batch filter by user IDs using IN clause
	query = query.Where("sales_rep_id IN ?", userIDs)

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
	for _, res := range results {
		stats[res.Status] = res.Count
	}

	return stats, nil
}

// GetStatsByDate returns visit report count grouped by date using database aggregation
func (r *repository) GetStatsByDate(startDate, endDate string, accountID, salesRepID, status string) (map[string]int64, error) {
	query := r.db.Table("visit_reports")

	// Apply date filters
	if startDate != "" {
		if start, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("visit_date >= ?", start)
		}
	}
	if endDate != "" {
		if end, err := time.Parse("2006-01-02", endDate); err == nil {
			end = end.Add(24 * time.Hour)
			query = query.Where("visit_date < ?", end)
		}
	}

	// Apply other filters
	if accountID != "" {
		query = query.Where("account_id = ?", accountID)
	}
	if salesRepID != "" {
		query = query.Where("sales_rep_id = ?", salesRepID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Aggregate by date
	var results []struct {
		Date  string
		Count int64
	}
	err := query.
		Select("DATE(visit_date) as date, COUNT(*) as count").
		Group("DATE(visit_date)").
		Order("DATE(visit_date) ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, r := range results {
		stats[r.Date] = r.Count
	}

	return stats, nil
}

// GetStatsByDateAndStatus returns visit report count grouped by both date and status using a single aggregation query
func (r *repository) GetStatsByDateAndStatus(startDate, endDate string, accountID, salesRepID string) (map[string]map[string]int64, error) {
	query := r.db.Table("visit_reports")

	if startDate != "" {
		if start, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("visit_date >= ?", start)
		}
	}
	if endDate != "" {
		if end, err := time.Parse("2006-01-02", endDate); err == nil {
			end = end.Add(24 * time.Hour)
			query = query.Where("visit_date < ?", end)
		}
	}
	if accountID != "" {
		query = query.Where("account_id = ?", accountID)
	}
	if salesRepID != "" {
		query = query.Where("sales_rep_id = ?", salesRepID)
	}

	var results []struct {
		Date   string
		Status string
		Count  int64
	}
	err := query.
		Select("DATE(visit_date) as date, status, COUNT(*) as count").
		Group("DATE(visit_date), status").
		Order("DATE(visit_date) ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	stats := make(map[string]map[string]int64)
	for _, row := range results {
		if stats[row.Date] == nil {
			stats[row.Date] = make(map[string]int64)
		}
		stats[row.Date][row.Status] = row.Count
	}

	return stats, nil
}

// GetStatsByAccount returns visit report count grouped by account using database aggregation
func (r *repository) GetStatsByAccount(startDate, endDate string, salesRepID, status string) (map[string]int64, error) {
	query := r.db.Table("visit_reports").Where("account_id IS NOT NULL")

	// Apply date filters
	if startDate != "" {
		if start, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("visit_date >= ?", start)
		}
	}
	if endDate != "" {
		if end, err := time.Parse("2006-01-02", endDate); err == nil {
			end = end.Add(24 * time.Hour)
			query = query.Where("visit_date < ?", end)
		}
	}

	// Apply other filters
	if salesRepID != "" {
		query = query.Where("sales_rep_id = ?", salesRepID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Aggregate by account
	var results []struct {
		AccountID string
		Count     int64
	}
	err := query.
		Select("account_id, COUNT(*) as count").
		Group("account_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, r := range results {
		stats[r.AccountID] = r.Count
	}

	return stats, nil
}

// GetStatsBySalesRep returns visit report count grouped by sales rep using database aggregation
func (r *repository) GetStatsBySalesRep(startDate, endDate string, accountID, status string) (map[string]int64, error) {
	query := r.db.Table("visit_reports")

	// Apply date filters
	if startDate != "" {
		if start, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("visit_date >= ?", start)
		}
	}
	if endDate != "" {
		if end, err := time.Parse("2006-01-02", endDate); err == nil {
			end = end.Add(24 * time.Hour)
			query = query.Where("visit_date < ?", end)
		}
	}

	// Apply other filters
	if accountID != "" {
		query = query.Where("account_id = ?", accountID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Aggregate by sales rep
	var results []struct {
		SalesRepID string
		Count      int64
	}
	err := query.
		Select("sales_rep_id, COUNT(*) as count").
		Group("sales_rep_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, r := range results {
		stats[r.SalesRepID] = r.Count
	}

	return stats, nil
}

// GetStatsBySalesRepWithAccounts returns visit count and unique account count per sales rep using database aggregation
func (r *repository) GetStatsBySalesRepWithAccounts(startDate, endDate string, status string) (map[string]struct {
	VisitCount   int64
	AccountCount int64
}, error) {
	query := r.db.Table("visit_reports")

	// Apply date filters
	if startDate != "" {
		if start, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("visit_date >= ?", start)
		}
	}
	if endDate != "" {
		if end, err := time.Parse("2006-01-02", endDate); err == nil {
			end = end.Add(24 * time.Hour)
			query = query.Where("visit_date < ?", end)
		}
	}

	// Apply status filter
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Aggregate by sales rep with visit count and unique account count
	var results []struct {
		SalesRepID   string
		VisitCount   int64
		AccountCount int64
	}
	err := query.
		Select("sales_rep_id, COUNT(*) as visit_count, COUNT(DISTINCT account_id) as account_count").
		Where("account_id IS NOT NULL").
		Group("sales_rep_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	stats := make(map[string]struct {
		VisitCount   int64
		AccountCount int64
	})
	for _, r := range results {
		stats[r.SalesRepID] = struct {
			VisitCount   int64
			AccountCount int64
		}{
			VisitCount:   r.VisitCount,
			AccountCount: r.AccountCount,
		}
	}

	return stats, nil
}
