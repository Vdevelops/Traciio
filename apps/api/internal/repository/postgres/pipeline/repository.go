package pipeline

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new pipeline repository
func NewRepository(db *gorm.DB) interfaces.PipelineRepository {
	return &repository{db: db}
}

func (r *repository) FindStageByID(id string) (*pipeline.PipelineStage, error) {
	var stage pipeline.PipelineStage
	err := r.db.Where("id = ?", id).First(&stage).Error
	if err != nil {
		return nil, err
	}
	return &stage, nil
}

func (r *repository) FindStageByCode(code string) (*pipeline.PipelineStage, error) {
	var stage pipeline.PipelineStage
	err := r.db.Where("code = ?", code).First(&stage).Error
	if err != nil {
		return nil, err
	}
	return &stage, nil
}

func (r *repository) ListStages(req *pipeline.ListPipelineStagesRequest) ([]pipeline.PipelineStage, error) {
	var stages []pipeline.PipelineStage

	query := r.db.Model(&pipeline.PipelineStage{})

	// Apply filters
	if req.IsActive != nil {
		query = query.Where("is_active = ?", *req.IsActive)
	}

	// Order by order field (use double quotes for PostgreSQL reserved keyword)
	err := query.Order("\"order\" ASC").Find(&stages).Error
	if err != nil {
		return nil, err
	}

	return stages, nil
}

func (r *repository) CreateStage(stage *pipeline.PipelineStage) error {
	return r.db.Create(stage).Error
}

func (r *repository) UpdateStage(stage *pipeline.PipelineStage) error {
	return r.db.Save(stage).Error
}

func (r *repository) DeleteStage(id string) error {
	return r.db.Where("id = ?", id).Delete(&pipeline.PipelineStage{}).Error
}
