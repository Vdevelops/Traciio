package activity_type

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/activity_type"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new activity type repository
func NewRepository(db *gorm.DB) interfaces.ActivityTypeRepository {
	return &repository{db: db}
}

func (r *repository) FindByID(id string) (*activity_type.ActivityType, error) {
	var at activity_type.ActivityType
	err := r.db.
		Select("*, (SELECT COUNT(*) FROM activities WHERE activities.activity_type_id = activity_types.id) as activity_count").
		Where("id = ?", id).First(&at).Error
	if err != nil {
		return nil, err
	}
	return &at, nil
}

func (r *repository) FindByCode(code string) (*activity_type.ActivityType, error) {
	var at activity_type.ActivityType
	err := r.db.
		Select("*, (SELECT COUNT(*) FROM activities WHERE activities.activity_type_id = activity_types.id) as activity_count").
		Where("code = ?", code).First(&at).Error
	if err != nil {
		return nil, err
	}
	return &at, nil
}

func (r *repository) List(req *activity_type.ListActivityTypesRequest) ([]activity_type.ActivityType, error) {
	var types []activity_type.ActivityType

	query := r.db.Model(&activity_type.ActivityType{})

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	if err := query.
		Select("*, (SELECT COUNT(*) FROM activities WHERE activities.activity_type_id = activity_types.id) as activity_count").
		Order("\"order\" ASC, name ASC").Find(&types).Error; err != nil {
		return nil, err
	}

	return types, nil
}

func (r *repository) Create(at *activity_type.ActivityType) error {
	return r.db.Create(at).Error
}

func (r *repository) Update(at *activity_type.ActivityType) error {
	return r.db.Save(at).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&activity_type.ActivityType{}).Error
}

