package contact_role

import (
	"errors"

	"github.com/gilabs/crm-healthcare/api/internal/domain/contact_role"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

var (
	ErrContactRoleNotFound      = errors.New("contact role not found")
	ErrContactRoleAlreadyExists = errors.New("contact role already exists")
)

type Service struct {
	contactRoleRepo interfaces.ContactRoleRepository
}

func NewService(contactRoleRepo interfaces.ContactRoleRepository) *Service {
	return &Service{
		contactRoleRepo: contactRoleRepo,
	}
}

// List returns a list of contact roles
func (s *Service) List() ([]contact_role.ContactRoleResponse, error) {
	contactRoles, err := s.contactRoleRepo.List()
	if err != nil {
		return nil, err
	}

	responses := make([]contact_role.ContactRoleResponse, len(contactRoles))
	for i, cr := range contactRoles {
		responses[i] = *cr.ToContactRoleResponse()
	}

	return responses, nil
}

// GetByID returns a contact role by ID
func (s *Service) GetByID(id string) (*contact_role.ContactRoleResponse, error) {
	cr, err := s.contactRoleRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrContactRoleNotFound
		}
		return nil, err
	}
	return cr.ToContactRoleResponse(), nil
}

// Create creates a new contact role
func (s *Service) Create(req *contact_role.CreateContactRoleRequest) (*contact_role.ContactRoleResponse, error) {
	// Check if code already exists
	_, err := s.contactRoleRepo.FindByCode(req.Code)
	if err == nil {
		return nil, ErrContactRoleAlreadyExists
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

	// Create contact role
	cr := &contact_role.ContactRole{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		BadgeColor:  badgeColor,
		Status:      status,
	}

	if err := s.contactRoleRepo.Create(cr); err != nil {
		return nil, err
	}

	// Reload
	createdContactRole, err := s.contactRoleRepo.FindByID(cr.ID)
	if err != nil {
		return nil, err
	}

	return createdContactRole.ToContactRoleResponse(), nil
}

// Update updates a contact role
func (s *Service) Update(id string, req *contact_role.UpdateContactRoleRequest) (*contact_role.ContactRoleResponse, error) {
	// Find contact role
	cr, err := s.contactRoleRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrContactRoleNotFound
		}
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		cr.Name = req.Name
	}

	if req.Code != "" {
		// Check if code already exists (excluding current contact role)
		existingContactRole, err := s.contactRoleRepo.FindByCode(req.Code)
		if err == nil && existingContactRole.ID != id {
			return nil, ErrContactRoleAlreadyExists
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		cr.Code = req.Code
	}

	if req.Description != "" {
		cr.Description = req.Description
	}

	if req.BadgeColor != "" {
		cr.BadgeColor = req.BadgeColor
	}

	if req.Status != "" {
		cr.Status = req.Status
	}

	if err := s.contactRoleRepo.Update(cr); err != nil {
		return nil, err
	}

	// Reload
	updatedContactRole, err := s.contactRoleRepo.FindByID(cr.ID)
	if err != nil {
		return nil, err
	}

	return updatedContactRole.ToContactRoleResponse(), nil
}

// Delete deletes a contact role
func (s *Service) Delete(id string) error {
	// Check if contact role exists
	_, err := s.contactRoleRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrContactRoleNotFound
		}
		return err
	}

	return s.contactRoleRepo.Delete(id)
}

