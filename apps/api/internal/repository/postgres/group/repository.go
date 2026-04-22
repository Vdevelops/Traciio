package group

import (
	"strings"

	"github.com/gilabs/crm-healthcare/api/internal/domain/group"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new group repository
func NewRepository(db *gorm.DB) interfaces.GroupRepository {
	return &repository{db: db}
}

func (r *repository) FindByID(id string) (*group.Group, error) {
	var g group.Group
	err := r.db.Where("id = ?", id).First(&g).Error
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *repository) FindByCode(code string) (*group.Group, error) {
	var g group.Group
	err := r.db.Where("code = ?", code).First(&g).Error
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *repository) List(req *group.ListGroupsRequest) ([]group.Group, int64, error) {
	var groups []group.Group
	var total int64

	query := r.db.Model(&group.Group{})

	// Apply filters
	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where(
			"LOWER(name) LIKE ? OR LOWER(code) LIKE ? OR LOWER(description) LIKE ?",
			search, search, search,
		)
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
	err := query.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&groups).Error
	if err != nil {
		return nil, 0, err
	}

	return groups, total, nil
}

func (r *repository) Create(g *group.Group) error {
	return r.db.Create(g).Error
}

func (r *repository) Update(g *group.Group) error {
	return r.db.Save(g).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&group.Group{}).Error
}

func (r *repository) CountUsersByGroupID(groupID string) (int64, error) {
	var count int64
	err := r.db.Table("users").Where("group_id = ?", groupID).Count(&count).Error
	return count, err
}

