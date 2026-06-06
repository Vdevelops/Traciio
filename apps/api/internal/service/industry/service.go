package industry

import (
	"errors"

	"github.com/gilabs/crm-healthcare/api/internal/domain/industry"
	industryrepo "github.com/gilabs/crm-healthcare/api/internal/repository/postgres/industry"

	"gorm.io/gorm"
)

// Service errors
var (
	ErrIndustryNotFound   = errors.New("industry not found")
	ErrIndustryCodeExists = errors.New("industry code already exists")
	ErrIndustryInUse      = errors.New("industry is being used by leads")
)

// Service implements industry use cases.
type Service struct {
	repo industryrepo.Repository
	db   *gorm.DB
}

// NewService creates a new industry service.
func NewService(repo industryrepo.Repository, db *gorm.DB) *Service {
	return &Service{
		repo: repo,
		db:   db,
	}
}

func (s *Service) Create(req *industry.CreateIndustryRequest, createdBy string) (*industry.IndustryResponse, error) {
	// Check if code already exists
	existing, err := s.repo.FindByCode(req.Code)
	if err == nil && existing != nil {
		return nil, ErrIndustryCodeExists
	}

	ind := &industry.Industry{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Order:       req.Order,
		IsActive:    true,
		CreatedBy:   createdBy,
	}

	if req.IsActive != nil {
		ind.IsActive = *req.IsActive
	}

	if err := s.repo.Create(ind); err != nil {
		return nil, err
	}

	return ind.ToIndustryResponse(), nil
}

func (s *Service) Update(id string, req *industry.UpdateIndustryRequest) (*industry.IndustryResponse, error) {
	ind, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrIndustryNotFound
		}
		return nil, err
	}

	// Check if code is being changed and already exists
	if req.Code != "" && req.Code != ind.Code {
		existing, err := s.repo.FindByCode(req.Code)
		if err == nil && existing != nil && existing.ID != id {
			return nil, ErrIndustryCodeExists
		}
		ind.Code = req.Code
	}

	if req.Name != "" {
		ind.Name = req.Name
	}
	if req.Description != "" {
		ind.Description = req.Description
	}
	if req.Order != nil {
		ind.Order = *req.Order
	}
	if req.IsActive != nil {
		ind.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ind); err != nil {
		return nil, err
	}

	return ind.ToIndustryResponse(), nil
}

func (s *Service) Delete(id string) error {
	ind, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrIndustryNotFound
		}
		return err
	}

	// Check if industry is being used by any leads
	var count int64
	if err := s.db.Table("leads").Where("industry = ?", ind.Name).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrIndustryInUse
	}

	return s.repo.Delete(id)
}

func (s *Service) FindByID(id string) (*industry.IndustryResponse, error) {
	ind, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrIndustryNotFound
		}
		return nil, err
	}
	return ind.ToIndustryResponse(), nil
}

func (s *Service) List(req *industry.ListIndustriesRequest) ([]*industry.IndustryResponse, int64, error) {
	industries, total, err := s.repo.List(req)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*industry.IndustryResponse, len(industries))
	for i, ind := range industries {
		responses[i] = ind.ToIndustryResponse()
	}

	return responses, total, nil
}

func (s *Service) ListAll() ([]*industry.IndustryResponse, error) {
	industries, err := s.repo.ListAll()
	if err != nil {
		return nil, err
	}

	responses := make([]*industry.IndustryResponse, len(industries))
	for i, ind := range industries {
		responses[i] = ind.ToIndustryResponse()
	}

	return responses, nil
}
