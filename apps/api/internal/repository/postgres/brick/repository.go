package brick

import (
	"fmt"
	"strings"

	"github.com/gilabs/crm-healthcare/api/internal/domain/brick"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new brick repository
func NewRepository(db *gorm.DB) interfaces.BrickRepository {
	return &repository{db: db}
}

func (r *repository) FindByID(id string) (*brick.Brick, error) {
	var b brick.Brick
	err := r.db.Preload("Manager").Where("id = ?", id).First(&b).Error
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *repository) FindByIDs(ids []string) ([]brick.Brick, error) {
	if len(ids) == 0 {
		return []brick.Brick{}, nil
	}
	var bricks []brick.Brick
	err := r.db.Preload("Manager").Where("id IN ?", ids).Find(&bricks).Error
	if err != nil {
		return nil, err
	}
	return bricks, nil
}

func (r *repository) FindByCode(code string) (*brick.Brick, error) {
	var b brick.Brick
	err := r.db.Preload("Manager").Where("code = ?", code).First(&b).Error
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *repository) List(req *brick.ListBricksRequest) ([]brick.Brick, int64, error) {
	var bricks []brick.Brick
	var total int64

	query := r.db.Model(&brick.Brick{}).Preload("Manager")

	// Apply filters
	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where(
			"LOWER(name) LIKE ? OR LOWER(code) LIKE ? OR LOWER(description) LIKE ? OR LOWER(regency) LIKE ? OR LOWER(province) LIKE ?",
			search, search, search, search, search,
		)
	}

	if req.Province != "" {
		query = query.Where("province = ?", req.Province)
	}

	if req.Regency != "" {
		query = query.Where("regency = ?", req.Regency)
	}

	if req.ManagerID != nil && *req.ManagerID != "" {
		query = query.Where("manager_id = ?", *req.ManagerID)
	}

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
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
	err := query.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&bricks).Error
	if err != nil {
		return nil, 0, err
	}

	return bricks, total, nil
}

func (r *repository) Create(b *brick.Brick) error {
	return r.db.Create(b).Error
}

func (r *repository) Update(b *brick.Brick) error {
	// Clear the Manager association to prevent GORM from overwriting ID
	b.Manager = nil
	return r.db.Save(b).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&brick.Brick{}).Error
}

func (r *repository) CountSalesByBrickID(brickID string) (int64, error) {
	var count int64
	err := r.db.Table("users").Where("brick_id = ?", brickID).Count(&count).Error
	return count, err
}

func (r *repository) GetSalesByBrickID(brickID string) ([]user.User, error) {
	var users []user.User
	err := r.db.Model(&user.User{}).Preload("Role").Where("brick_id = ?", brickID).Find(&users).Error
	return users, err
}

func (r *repository) FindByRegencyAndProvince(regency, province string) (*brick.Brick, error) {
	var b brick.Brick
	err := r.db.Preload("Manager").
		Where("regency = ? AND province = ?", regency, province).
		First(&b).Error
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// GetNextCodeSequence finds existing codes with the given prefix and returns the next sequence number.
// E.g. prefix "BRK-JKT" with existing codes "BRK-JKT-001", "BRK-JKT-002" returns 3.
func (r *repository) GetNextCodeSequence(prefix string) (int, error) {
	var maxCode string
	err := r.db.Model(&brick.Brick{}).
		Select("code").
		Where("code LIKE ?", fmt.Sprintf("%s-%%", prefix)).
		Order("code DESC").
		Limit(1).
		Pluck("code", &maxCode).Error
	if err != nil {
		return 1, nil
	}

	if maxCode == "" {
		return 1, nil
	}

	// Extract the sequence part after the last "-"
	parts := strings.Split(maxCode, "-")
	if len(parts) == 0 {
		return 1, nil
	}

	lastPart := parts[len(parts)-1]
	var seq int
	if _, err := fmt.Sscanf(lastPart, "%d", &seq); err != nil {
		return 1, nil
	}

	return seq + 1, nil
}

