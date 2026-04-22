package interfaces

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead_qualification"
)

type LeadQualificationRepository interface {
	FindByLeadID(leadID string) (*lead_qualification.LeadQualificationChecklist, error)
	Upsert(checklist *lead_qualification.LeadQualificationChecklist) error
	DeleteByLeadID(leadID string) error
}
