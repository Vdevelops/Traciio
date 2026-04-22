package category

import (
	"errors"

	"github.com/gilabs/crm-healthcare/api/internal/domain/category"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

var (
	ErrCategoryNotFound      = errors.New("category not found")
	ErrCategoryAlreadyExists = errors.New("category already exists")
)

type Service struct {
	categoryRepo interfaces.CategoryRepository
}

func NewService(categoryRepo interfaces.CategoryRepository) *Service {
	return &Service{
		categoryRepo: categoryRepo,
	}
}

// List returns a list of categories
func (s *Service) List() ([]category.CategoryResponse, error) {
	categories, err := s.categoryRepo.List()
	if err != nil {
		return nil, err
	}

	responses := make([]category.CategoryResponse, len(categories))
	for i, c := range categories {
		responses[i] = *c.ToCategoryResponse()
	}

	return responses, nil
}

// GetByID returns a category by ID
func (s *Service) GetByID(id string) (*category.CategoryResponse, error) {
	c, err := s.categoryRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return c.ToCategoryResponse(), nil
}

// Create creates a new category
func (s *Service) Create(req *category.CreateCategoryRequest) (*category.CategoryResponse, error) {
	// Check if code already exists
	_, err := s.categoryRepo.FindByCode(req.Code)
	if err == nil {
		return nil, ErrCategoryAlreadyExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Set defaults
	status := req.Status
	if status == "" {
		status = "active"
	}
	badgeColor := req.BadgeColor
	if badgeColor == "" {
		badgeColor = "outline"
	}

	// Create category
	cat := &category.Category{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		BadgeColor:  badgeColor,
		Status:      status,
	}

	if err := s.categoryRepo.Create(cat); err != nil {
		return nil, err
	}

	// Reload
	createdCategory, err := s.categoryRepo.FindByID(cat.ID)
	if err != nil {
		return nil, err
	}

	return createdCategory.ToCategoryResponse(), nil
}

// Update updates a category
func (s *Service) Update(id string, req *category.UpdateCategoryRequest) (*category.CategoryResponse, error) {
	// Find category
	cat, err := s.categoryRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		cat.Name = req.Name
	}

	if req.Code != "" {
		// Check if code already exists (excluding current category)
		existingCategory, err := s.categoryRepo.FindByCode(req.Code)
		if err == nil && existingCategory.ID != id {
			return nil, ErrCategoryAlreadyExists
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		cat.Code = req.Code
	}

	if req.Description != "" {
		cat.Description = req.Description
	}

	if req.BadgeColor != "" {
		cat.BadgeColor = req.BadgeColor
	}

	if req.Status != "" {
		cat.Status = req.Status
	}

	if err := s.categoryRepo.Update(cat); err != nil {
		return nil, err
	}

	// Reload
	updatedCategory, err := s.categoryRepo.FindByID(cat.ID)
	if err != nil {
		return nil, err
	}

	return updatedCategory.ToCategoryResponse(), nil
}

// Delete deletes a category
func (s *Service) Delete(id string) error {
	// Check if category exists
	_, err := s.categoryRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCategoryNotFound
		}
		return err
	}

	return s.categoryRepo.Delete(id)
}

