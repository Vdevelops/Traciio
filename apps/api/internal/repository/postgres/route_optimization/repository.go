package route_optimization

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/route_optimization"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new route optimization repository
func NewRepository(db *gorm.DB) interfaces.RouteOptimizationRepository {
	return &repository{db: db}
}

func (r *repository) FindByID(id string) (*route_optimization.OptimizedRoute, error) {
	var route route_optimization.OptimizedRoute
	err := r.db.Where("id = ?", id).First(&route).Error
	if err != nil {
		return nil, err
	}
	return &route, nil
}

func (r *repository) List(req *route_optimization.ListRoutesRequest) ([]route_optimization.OptimizedRoute, int64, error) {
	var routes []route_optimization.OptimizedRoute
	var total int64

	query := r.db.Model(&route_optimization.OptimizedRoute{})

	// Apply filters
	if req.UserID != "" {
		query = query.Where("user_id = ?", req.UserID)
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
	err := query.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&routes).Error
	if err != nil {
		return nil, 0, err
	}

	return routes, total, nil
}

func (r *repository) Create(route *route_optimization.OptimizedRoute) error {
	return r.db.Create(route).Error
}

func (r *repository) Update(route *route_optimization.OptimizedRoute) error {
	return r.db.Save(route).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&route_optimization.OptimizedRoute{}).Error
}

func (r *repository) FindByUserID(userID string) ([]route_optimization.OptimizedRoute, error) {
	var routes []route_optimization.OptimizedRoute
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&routes).Error
	if err != nil {
		return nil, err
	}
	return routes, nil
}


