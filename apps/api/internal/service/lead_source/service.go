package lead_source

import (
	"errors"

	"github.com/gilabs/crm-healthcare/api/internal/domain/lead_source"
	lead_source_repo "github.com/gilabs/crm-healthcare/api/internal/repository/lead_source"

	"gorm.io/gorm"
)

// Service errors
var (
	ErrLeadSourceNotFound   = errors.New("lead source not found")
	ErrLeadSourceCodeExists = errors.New("lead source code already exists")
	ErrLeadSourceInUse      = errors.New("lead source is being used by leads")
)

// Service defines lead source service interface
type Service interface {
	Create(req *lead_source.CreateLeadSourceRequest, createdBy string) (*lead_source.LeadSourceResponse, error)
	Update(id string, req *lead_source.UpdateLeadSourceRequest) (*lead_source.LeadSourceResponse, error)
	Delete(id string) error
	FindByID(id string) (*lead_source.LeadSourceResponse, error)
	List(req *lead_source.ListLeadSourcesRequest) ([]*lead_source.LeadSourceResponse, int64, error)
	ListAll() ([]*lead_source.LeadSourceResponse, error)
}

type service struct {
	repo lead_source_repo.Repository
	db   *gorm.DB
}

// NewService creates a new lead source service
func NewService(repo lead_source_repo.Repository, db *gorm.DB) Service {
	return &service{
		repo: repo,
		db:   db,
	}
}

func (s *service) Create(req *lead_source.CreateLeadSourceRequest, createdBy string) (*lead_source.LeadSourceResponse, error) {
	// Check if code already exists
	existing, err := s.repo.FindByCode(req.Code)
	if err == nil && existing != nil {
		return nil, ErrLeadSourceCodeExists
	}

	ls := &lead_source.LeadSource{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Order:       req.Order,
		IsActive:    true,
		CreatedBy:   createdBy,
	}

	if req.IsActive != nil {
		ls.IsActive = *req.IsActive
	}

	if err := s.repo.Create(ls); err != nil {
		return nil, err
	}

	return ls.ToLeadSourceResponse(), nil
}

func (s *service) Update(id string, req *lead_source.UpdateLeadSourceRequest) (*lead_source.LeadSourceResponse, error) {
	ls, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLeadSourceNotFound
		}
		return nil, err
	}

	// Check if code is being changed and already exists
	if req.Code != "" && req.Code != ls.Code {
		existing, err := s.repo.FindByCode(req.Code)
		if err == nil && existing != nil && existing.ID != id {
			return nil, ErrLeadSourceCodeExists
		}
		ls.Code = req.Code
	}

	if req.Name != "" {
		ls.Name = req.Name
	}
	if req.Description != "" {
		ls.Description = req.Description
	}
	if req.Order != nil {
		ls.Order = *req.Order
	}
	if req.IsActive != nil {
		ls.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ls); err != nil {
		return nil, err
	}

	return ls.ToLeadSourceResponse(), nil
}

func (s *service) Delete(id string) error {
	ls, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrLeadSourceNotFound
		}
		return err
	}

	// Check if lead source is being used by any leads
	var count int64
	if err := s.db.Table("leads").Where("lead_source = ?", ls.Code).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrLeadSourceInUse
	}

	return s.repo.Delete(id)
}

func (s *service) FindByID(id string) (*lead_source.LeadSourceResponse, error) {
	ls, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLeadSourceNotFound
		}
		return nil, err
	}
	return ls.ToLeadSourceResponse(), nil
}

func (s *service) List(req *lead_source.ListLeadSourcesRequest) ([]*lead_source.LeadSourceResponse, int64, error) {
	leadSources, total, err := s.repo.List(req)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*lead_source.LeadSourceResponse, len(leadSources))
	for i, ls := range leadSources {
		responses[i] = ls.ToLeadSourceResponse()
	}

	return responses, total, nil
}

func (s *service) ListAll() ([]*lead_source.LeadSourceResponse, error) {
	leadSources, err := s.repo.ListAll()
	if err != nil {
		return nil, err
	}

	responses := make([]*lead_source.LeadSourceResponse, len(leadSources))
	for i, ls := range leadSources {
		responses[i] = ls.ToLeadSourceResponse()
	}

	return responses, nil
}

