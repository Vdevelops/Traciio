package lead_status

import (
	"errors"

	"github.com/gilabs/crm-healthcare/api/internal/domain/lead_status"
	leadstatusrepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/lead_status"

	"gorm.io/gorm"
)

// Service errors
var (
	ErrLeadStatusNotFound    = errors.New("lead status not found")
	ErrLeadStatusCodeExists  = errors.New("lead status code already exists")
	ErrCannotDeleteDefault   = errors.New("cannot delete default lead status")
	ErrCannotDeleteConverted = errors.New("cannot delete converted lead status")
	ErrLeadStatusInUse       = errors.New("lead status is being used by leads")
)

// Service implements lead status use cases.
type Service struct {
	repo leadstatusrepo.Repository
	db   *gorm.DB
}

// NewService creates a new lead status service.
func NewService(repo leadstatusrepo.Repository, db *gorm.DB) *Service {
	return &Service{
		repo: repo,
		db:   db,
	}
}

func (s *Service) Create(req *lead_status.CreateLeadStatusRequest, createdBy string) (*lead_status.LeadStatusResponse, error) {
	// Check if code already exists
	existing, err := s.repo.FindByCode(req.Code)
	if err == nil && existing != nil {
		return nil, ErrLeadStatusCodeExists
	}

	status := &lead_status.LeadStatus{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Score:       *req.Score,
		Color:       req.Color,
		Order:       req.Order,
		IsActive:    true,
		IsDefault:   false,
		IsConverted: false,
		CreatedBy:   createdBy,
	}

	if req.IsActive != nil {
		status.IsActive = *req.IsActive
	}
	if req.IsDefault != nil {
		status.IsDefault = *req.IsDefault
	}
	if req.IsConverted != nil {
		status.IsConverted = *req.IsConverted
	}

	// If marked as default, unset other defaults
	if status.IsDefault {
		if err := s.repo.SetDefault(""); err != nil {
			return nil, err
		}
	}

	if err := s.repo.Create(status); err != nil {
		return nil, err
	}

	return status.ToLeadStatusResponse(), nil
}

func (s *Service) Update(id string, req *lead_status.UpdateLeadStatusRequest) (*lead_status.LeadStatusResponse, error) {
	status, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLeadStatusNotFound
		}
		return nil, err
	}

	// Check if code is being changed and already exists
	if req.Code != "" && req.Code != status.Code {
		existing, err := s.repo.FindByCode(req.Code)
		if err == nil && existing != nil && existing.ID != id {
			return nil, ErrLeadStatusCodeExists
		}
		status.Code = req.Code
	}

	if req.Name != "" {
		status.Name = req.Name
	}
	if req.Description != "" {
		status.Description = req.Description
	}
	if req.Score != nil {
		status.Score = *req.Score
	}
	if req.Color != "" {
		status.Color = req.Color
	}
	if req.Order != nil {
		status.Order = *req.Order
	}
	if req.IsActive != nil {
		status.IsActive = *req.IsActive
	}
	if req.IsDefault != nil && *req.IsDefault {
		// If setting as default, unset others first
		if err := s.repo.SetDefault(id); err != nil {
			return nil, err
		}
		status.IsDefault = true
	}
	if req.IsConverted != nil {
		status.IsConverted = *req.IsConverted
	}

	if err := s.repo.Update(status); err != nil {
		return nil, err
	}

	return status.ToLeadStatusResponse(), nil
}

func (s *Service) Delete(id string) error {
	status, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrLeadStatusNotFound
		}
		return err
	}

	// Cannot delete default status
	if status.IsDefault {
		return ErrCannotDeleteDefault
	}

	// Cannot delete converted status
	if status.IsConverted {
		return ErrCannotDeleteConverted
	}

	// Check if status is being used by any leads
	var count int64
	if err := s.db.Table("leads").Where("lead_status = ?", status.Code).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrLeadStatusInUse
	}

	return s.repo.Delete(id)
}

func (s *Service) FindByID(id string) (*lead_status.LeadStatusResponse, error) {
	status, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLeadStatusNotFound
		}
		return nil, err
	}
	return status.ToLeadStatusResponse(), nil
}

func (s *Service) List(req *lead_status.ListLeadStatusesRequest) ([]*lead_status.LeadStatusResponse, int64, error) {
	statuses, total, err := s.repo.List(req)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*lead_status.LeadStatusResponse, len(statuses))
	for i, status := range statuses {
		responses[i] = status.ToLeadStatusResponse()
	}

	return responses, total, nil
}

func (s *Service) ListAll() ([]*lead_status.LeadStatusResponse, error) {
	statuses, err := s.repo.ListAll()
	if err != nil {
		return nil, err
	}

	responses := make([]*lead_status.LeadStatusResponse, len(statuses))
	for i, status := range statuses {
		responses[i] = status.ToLeadStatusResponse()
	}

	return responses, nil
}

func (s *Service) SetDefault(id string) error {
	// Verify status exists
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrLeadStatusNotFound
		}
		return err
	}

	return s.repo.SetDefault(id)
}
