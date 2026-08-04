package kpi

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a KPI repository backed by PostgreSQL.
func NewRepository(db *gorm.DB) interfaces.KPIRepository {
	return &repository{db: db}
}

func (r *repository) count(query *gorm.DB) (int64, error) {
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *repository) CountDealsCreated(userID string, startDate, endDate time.Time) (int64, error) {
	return r.count(
		r.db.Table("deals").
			Where("deleted_at IS NULL").
			Where("assigned_to = ?", userID).
			Where("created_at >= ? AND created_at <= ?", startDate, endDate),
	)
}

func (r *repository) CountWonDeals(userID string, startDate, endDate time.Time) (int64, error) {
	return r.count(
		r.db.Table("deals").
			Where("deleted_at IS NULL").
			Where("assigned_to = ? AND status = ?", userID, "won").
			Where("COALESCE(actual_close_date, created_at) >= ? AND COALESCE(actual_close_date, created_at) <= ?", startDate, endDate),
	)
}

func (r *repository) SumWonRevenue(userID string, startDate, endDate time.Time) (int64, error) {
	var total int64
	err := r.db.Table("deals").
		Select("COALESCE(SUM(value), 0)").
		Where("deleted_at IS NULL").
		Where("assigned_to = ? AND status = ?", userID, "won").
		Where("COALESCE(actual_close_date, created_at) >= ? AND COALESCE(actual_close_date, created_at) <= ?", startDate, endDate).
		Scan(&total).Error
	return total, err
}

func (r *repository) CountVisitCompleted(userID string, startDate, endDate time.Time) (int64, error) {
	return r.count(
		r.db.Table("visit_reports").
			Where("deleted_at IS NULL").
			Where("sales_rep_id = ?", userID).
			Where("status IN ?", []string{"completed", "approved"}).
			Where("visit_date >= ? AND visit_date <= ?", startDate, endDate),
	)
}

func (r *repository) CountVisitPlanned(userID string, startDate, endDate time.Time) (int64, error) {
	return r.count(
		r.db.Table("schedules").
			Where("deleted_at IS NULL").
			Where("user_id = ?", userID).
			Where("scheduled_at >= ? AND scheduled_at <= ?", startDate, endDate).
			Where("status NOT IN ?", []string{"cancelled", "rejected"}),
	)
}

func (r *repository) CountTasksCreated(userID string, startDate, endDate time.Time) (int64, error) {
	return r.count(
		r.db.Table("tasks").
			Where("deleted_at IS NULL").
			Where("assigned_to = ?", userID).
			Where("created_at >= ? AND created_at <= ?", startDate, endDate),
	)
}

func (r *repository) CountTasksCompleted(userID string, startDate, endDate time.Time) (int64, error) {
	return r.count(
		r.db.Table("tasks").
			Where("deleted_at IS NULL").
			Where("assigned_to = ?", userID).
			Where("status = ?", "completed").
			Where("COALESCE(completed_at, updated_at, created_at) >= ? AND COALESCE(completed_at, updated_at, created_at) <= ?", startDate, endDate),
	)
}

func (r *repository) CountOverdueTasks(userID string, startDate, endDate time.Time) (int64, error) {
	return r.count(
		r.db.Table("tasks").
			Where("deleted_at IS NULL").
			Where("assigned_to = ?", userID).
			Where("created_at >= ? AND created_at <= ?", startDate, endDate).
			Where("due_date IS NOT NULL AND due_date < ? AND status != ?", endDate, "completed"),
	)
}

func (r *repository) GetPipelineMovementScore(userID string, startDate, endDate time.Time) (int64, error) {
	var score int64
	err := r.db.Table("deal_histories dh").
		Select(`
			COALESCE(SUM(
				CASE
					WHEN COALESCE(ts.code, '') = 'closed_lost' THEN 0
					WHEN COALESCE(ts."order", 0) > COALESCE(fs."order", 0) THEN COALESCE(ts."order", 0) - COALESCE(fs."order", 0)
					ELSE 0
				END
			), 0) AS score
		`).
		Joins("INNER JOIN deals d ON d.id = dh.deal_id AND d.deleted_at IS NULL AND d.assigned_to = ?", userID).
		Joins("LEFT JOIN pipeline_stages fs ON fs.id = dh.from_stage_id AND fs.deleted_at IS NULL").
		Joins("LEFT JOIN pipeline_stages ts ON ts.id = dh.to_stage_id AND ts.deleted_at IS NULL").
		Where("dh.deleted_at IS NULL").
		Where("dh.changed_at >= ? AND dh.changed_at <= ?", startDate, endDate).
		Scan(&score).Error
	return score, err
}

func (r *repository) CountDealsWithoutBrick(userID string, startDate, endDate time.Time) (int64, error) {
	return r.count(
		r.db.Table("deals").
			Where("deleted_at IS NULL").
			Where("assigned_to = ?", userID).
			Where("created_at >= ? AND created_at <= ?", startDate, endDate).
			Where("brick_id IS NULL"),
	)
}

func (r *repository) CountVisitReportsWithoutBrick(userID string, startDate, endDate time.Time) (int64, error) {
	return r.count(
		r.db.Table("visit_reports").
			Where("deleted_at IS NULL").
			Where("sales_rep_id = ?", userID).
			Where("visit_date >= ? AND visit_date <= ?", startDate, endDate).
			Where("brick_id IS NULL"),
	)
}

func (r *repository) CountCustomersWithInteractionInBrick(brickID string, startDate, endDate time.Time) (int64, error) {
	var total int64
	err := r.db.Raw(`
		SELECT COUNT(DISTINCT account_id) AS total
		FROM (
			SELECT d.account_id
			FROM deals d
			WHERE d.deleted_at IS NULL
			  AND d.brick_id = ?
			  AND d.created_at >= ? AND d.created_at <= ?
			UNION ALL
			SELECT vr.account_id
			FROM visit_reports vr
			WHERE vr.deleted_at IS NULL
			  AND vr.brick_id = ?
			  AND vr.visit_date >= ? AND vr.visit_date <= ?
		) interactions
	`, brickID, startDate, endDate, brickID, startDate, endDate).Scan(&total).Error
	return total, err
}

func (r *repository) CountRegisteredCustomersInBrick(brickID string) (int64, error) {
	return r.count(
		r.db.Table("accounts").
			Where("deleted_at IS NULL").
			Where("brick_id = ?", brickID),
	)
}
