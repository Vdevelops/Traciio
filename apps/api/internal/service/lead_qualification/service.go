package lead_qualification

import (
	"encoding/json"
	"errors"

	"github.com/gilabs/crm-healthcare/api/internal/domain/lead_qualification"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var (
	ErrChecklistNotFound = errors.New("qualification checklist not found")
)

type Service struct {
	repo     interfaces.LeadQualificationRepository
	leadRepo interfaces.LeadRepository
}

func NewService(repo interfaces.LeadQualificationRepository, leadRepo interfaces.LeadRepository) *Service {
	return &Service{
		repo:     repo,
		leadRepo: leadRepo,
	}
}

func (s *Service) GetByLeadID(leadID string) (*lead_qualification.LeadQualificationResponse, error) {
	checklist, err := s.repo.FindByLeadID(leadID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If not found, return empty response instead of error
			return &lead_qualification.LeadQualificationResponse{
				LeadID:              leadID,
				QualificationStatus: "pending",
			}, nil
		}
		return nil, err
	}

	return s.MapToResponse(checklist), nil
}

func (s *Service) Upsert(leadID string, req *lead_qualification.UpsertLeadQualificationRequest) (*lead_qualification.LeadQualificationResponse, error) {
	// Verify lead exists
	_, err := s.leadRepo.FindByID(leadID)
	if err != nil {
		return nil, err
	}

	checklist, err := s.repo.FindByLeadID(leadID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			checklist = &lead_qualification.LeadQualificationChecklist{
				LeadID: leadID,
			}
		} else {
			return nil, err
		}
	}

	// Update fields
	if req.BudgetTargetAmount != nil {
		checklist.BudgetTargetAmount = *req.BudgetTargetAmount
	}
	if req.BudgetTargetCurrency != "" {
		checklist.BudgetTargetCurrency = req.BudgetTargetCurrency
	}
	if req.BudgetConfirmed != nil {
		checklist.BudgetConfirmed = *req.BudgetConfirmed
	}
	if req.BudgetNotes != "" {
		checklist.BudgetNotes = req.BudgetNotes
	}

	if req.AuthorityTargetPerson != "" {
		checklist.AuthorityTargetPerson = req.AuthorityTargetPerson
	}
	if req.AuthorityTargetRole != "" {
		checklist.AuthorityTargetRole = req.AuthorityTargetRole
	}
	if req.AuthorityConfirmed != nil {
		checklist.AuthorityConfirmed = *req.AuthorityConfirmed
	}
	if req.AuthorityNotes != "" {
		checklist.AuthorityNotes = req.AuthorityNotes
	}

	if req.NeedTargetProducts != nil {
		data, _ := json.Marshal(req.NeedTargetProducts)
		checklist.NeedTargetProducts = datatypes.JSON(data)
	}
	if req.NeedPriorityLevel != "" {
		checklist.NeedPriorityLevel = req.NeedPriorityLevel
	}
	if req.NeedConfirmed != nil {
		checklist.NeedConfirmed = *req.NeedConfirmed
	}
	if req.NeedNotes != "" {
		checklist.NeedNotes = req.NeedNotes
	}

	if req.TimelineTargetDate != nil {
		checklist.TimelineTargetDate = req.TimelineTargetDate
	}
	if req.TimelineFlexibility != "" {
		checklist.TimelineFlexibility = req.TimelineFlexibility
	}
	if req.TimelineConfirmed != nil {
		checklist.TimelineConfirmed = *req.TimelineConfirmed
	}
	if req.TimelineNotes != "" {
		checklist.TimelineNotes = req.TimelineNotes
	}

	// Score and status are calculated automatically in Domain hook BeforeSave

	if err := s.repo.Upsert(checklist); err != nil {
		return nil, err
	}

	// Reload to get calculated values
	checklist, _ = s.repo.FindByLeadID(leadID)

	return s.MapToResponse(checklist), nil
}

func (s *Service) MapToResponse(q *lead_qualification.LeadQualificationChecklist) *lead_qualification.LeadQualificationResponse {
	var products []lead_qualification.NeedProduct
	if q.NeedTargetProducts != nil {
		_ = json.Unmarshal(q.NeedTargetProducts, &products)
	}

	return &lead_qualification.LeadQualificationResponse{
		ID:                    q.ID,
		LeadID:                q.LeadID,
		BudgetTargetAmount:    q.BudgetTargetAmount,
		BudgetTargetCurrency:  q.BudgetTargetCurrency,
		BudgetConfirmed:       q.BudgetConfirmed,
		BudgetNotes:           q.BudgetNotes,
		AuthorityTargetPerson: q.AuthorityTargetPerson,
		AuthorityTargetRole:   q.AuthorityTargetRole,
		AuthorityConfirmed:    q.AuthorityConfirmed,
		AuthorityNotes:        q.AuthorityNotes,
		NeedTargetProducts:    products,
		NeedPriorityLevel:     q.NeedPriorityLevel,
		NeedConfirmed:         q.NeedConfirmed,
		NeedNotes:             q.NeedNotes,
		TimelineTargetDate:    q.TimelineTargetDate,
		TimelineFlexibility:   q.TimelineFlexibility,
		TimelineConfirmed:     q.TimelineConfirmed,
		TimelineNotes:         q.TimelineNotes,
		QualificationScore:    q.QualificationScore,
		QualificationStatus:   q.QualificationStatus,
		BANTProgress: lead_qualification.BANTProgress{
			Budget:    lead_qualification.BANTItem{Completed: q.BudgetConfirmed, Score: s.getScore(q.BudgetConfirmed), MaxScore: 25},
			Authority: lead_qualification.BANTItem{Completed: q.AuthorityConfirmed, Score: s.getScore(q.AuthorityConfirmed), MaxScore: 25},
			Need:      lead_qualification.BANTItem{Completed: q.NeedConfirmed, Score: s.getScore(q.NeedConfirmed), MaxScore: 25},
			Timeline:  lead_qualification.BANTItem{Completed: q.TimelineConfirmed, Score: s.getScore(q.TimelineConfirmed), MaxScore: 25},
		},
		CreatedAt: q.CreatedAt,
		UpdatedAt: q.UpdatedAt,
	}
}

func (s *Service) getScore(confirmed bool) int {
	if confirmed {
		return 25
	}
	return 0
}
