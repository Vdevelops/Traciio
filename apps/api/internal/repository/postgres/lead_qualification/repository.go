package lead_qualification

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead_qualification"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type leadQualificationRepository struct {
	db *gorm.DB
}

func NewLeadQualificationRepository(db *gorm.DB) interfaces.LeadQualificationRepository {
	return &leadQualificationRepository{db: db}
}

func (r *leadQualificationRepository) FindByLeadID(leadID string) (*lead_qualification.LeadQualificationChecklist, error) {
	var checklist lead_qualification.LeadQualificationChecklist
	err := r.db.Where("lead_id = ?", leadID).First(&checklist).Error
	if err != nil {
		return nil, err
	}
	return &checklist, nil
}

func (r *leadQualificationRepository) Upsert(checklist *lead_qualification.LeadQualificationChecklist) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "lead_id"}},
		UpdateAll: true,
	}).Create(checklist).Error
}

func (r *leadQualificationRepository) DeleteByLeadID(leadID string) error {
	return r.db.Where("lead_id = ?", leadID).Delete(&lead_qualification.LeadQualificationChecklist{}).Error
}
