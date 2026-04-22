package brick_target_distribution

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/brick_target_distribution"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new brick target distribution repository
func NewRepository(db *gorm.DB) interfaces.BrickTargetDistributionRepository {
	return &repository{db: db}
}

func (r *repository) FindByID(id string) (*brick_target_distribution.BrickTargetDistribution, error) {
	var btd brick_target_distribution.BrickTargetDistribution
	err := r.db.Preload("Brick").Preload("BrickTarget").Preload("SalesUser").Preload("DistributedByUser").
		Where("id = ?", id).First(&btd).Error
	if err != nil {
		return nil, err
	}
	return &btd, nil
}

func (r *repository) List(req *brick_target_distribution.ListBrickTargetDistributionsRequest) ([]brick_target_distribution.BrickTargetDistribution, int64, error) {
	var distributions []brick_target_distribution.BrickTargetDistribution
	var total int64

	query := r.db.Model(&brick_target_distribution.BrickTargetDistribution{}).
		Preload("Brick").Preload("BrickTarget").Preload("SalesUser").Preload("DistributedByUser")

	// Apply filters
	if req.BrickID != nil && *req.BrickID != "" {
		query = query.Where("brick_id = ?", *req.BrickID)
	}

	if req.BrickTargetID != nil && *req.BrickTargetID != "" {
		query = query.Where("brick_target_id = ?", *req.BrickTargetID)
	}

	if req.SalesUserID != nil && *req.SalesUserID != "" {
		query = query.Where("sales_user_id = ?", *req.SalesUserID)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
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

	// Fetch data
	err := query.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&distributions).Error
	if err != nil {
		return nil, 0, err
	}

	return distributions, total, nil
}

func (r *repository) FindByBrickTargetID(brickTargetID string) ([]brick_target_distribution.BrickTargetDistribution, error) {
	var distributions []brick_target_distribution.BrickTargetDistribution
	err := r.db.Preload("Brick").Preload("BrickTarget").Preload("SalesUser").Preload("DistributedByUser").
		Where("brick_target_id = ?", brickTargetID).
		Find(&distributions).Error
	return distributions, err
}

func (r *repository) FindBySalesUserIDAndBrickTargetID(salesUserID, brickTargetID string) (*brick_target_distribution.BrickTargetDistribution, error) {
	var btd brick_target_distribution.BrickTargetDistribution
	err := r.db.Preload("Brick").Preload("BrickTarget").Preload("SalesUser").Preload("DistributedByUser").
		Where("sales_user_id = ? AND brick_target_id = ?", salesUserID, brickTargetID).
		First(&btd).Error
	if err != nil {
		return nil, err
	}
	return &btd, nil
}

func (r *repository) Create(btd *brick_target_distribution.BrickTargetDistribution) error {
	return r.db.Create(btd).Error
}

func (r *repository) Update(btd *brick_target_distribution.BrickTargetDistribution) error {
	// Clear associations to prevent GORM from overwriting IDs
	btd.Brick = nil
	btd.BrickTarget = nil
	btd.SalesUser = nil
	btd.DistributedByUser = nil
	return r.db.Save(btd).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&brick_target_distribution.BrickTargetDistribution{}).Error
}

func (r *repository) BulkCreate(distributions []*brick_target_distribution.BrickTargetDistribution) error {
	if len(distributions) == 0 {
		return nil
	}
	return r.db.Create(distributions).Error
}

