package activity_type

import (
	"errors"

	"github.com/gilabs/crm-healthcare/api/internal/domain/activity_type"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

var (
	ErrActivityTypeNotFound = errors.New("activity type not found")
)

type Service struct {
	activityTypeRepo interfaces.ActivityTypeRepository
}

func NewService(activityTypeRepo interfaces.ActivityTypeRepository) *Service {
	return &Service{
		activityTypeRepo: activityTypeRepo,
	}
}

// List returns a list of activity types
func (s *Service) List(req *activity_type.ListActivityTypesRequest) ([]activity_type.ActivityTypeResponse, error) {
	types, err := s.activityTypeRepo.List(req)
	if err != nil {
		return nil, err
	}

	responses := make([]activity_type.ActivityTypeResponse, len(types))
	for i := range types {
		responses[i] = *types[i].ToActivityTypeResponse()
	}

	return responses, nil
}

// GetByID returns an activity type by ID
func (s *Service) GetByID(id string) (*activity_type.ActivityTypeResponse, error) {
	at, err := s.activityTypeRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrActivityTypeNotFound
		}
		return nil, err
	}

	return at.ToActivityTypeResponse(), nil
}

// Create creates a new activity type
func (s *Service) Create(req *activity_type.CreateActivityTypeRequest) (*activity_type.ActivityTypeResponse, error) {
	status := req.Status
	if status == "" {
		status = "active"
	}

	badgeColor := req.BadgeColor
	if badgeColor == "" {
		badgeColor = "outline"
	}

	at := &activity_type.ActivityType{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Icon:        req.Icon,
		BadgeColor:  badgeColor,
		Status:      status,
		Order:       req.Order,
	}

	if err := s.activityTypeRepo.Create(at); err != nil {
		return nil, err
	}

	// Reload to ensure we return latest state
	at, err := s.activityTypeRepo.FindByID(at.ID)
	if err != nil {
		return nil, err
	}

	return at.ToActivityTypeResponse(), nil
}

// Update updates an existing activity type
func (s *Service) Update(id string, req *activity_type.UpdateActivityTypeRequest) (*activity_type.ActivityTypeResponse, error) {
	at, err := s.activityTypeRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrActivityTypeNotFound
		}
		return nil, err
	}

	if req.Name != "" {
		at.Name = req.Name
	}
	if req.Code != "" {
		at.Code = req.Code
	}
	if req.Description != "" {
		at.Description = req.Description
	}
	if req.Icon != "" {
		at.Icon = req.Icon
	}
	if req.BadgeColor != "" {
		at.BadgeColor = req.BadgeColor
	}
	if req.Status != "" {
		at.Status = req.Status
	}
	if req.Order != nil {
		at.Order = *req.Order
	}

	if err := s.activityTypeRepo.Update(at); err != nil {
		return nil, err
	}

	at, err = s.activityTypeRepo.FindByID(at.ID)
	if err != nil {
		return nil, err
	}

	return at.ToActivityTypeResponse(), nil
}

// Delete deletes an activity type
func (s *Service) Delete(id string) error {
	_, err := s.activityTypeRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrActivityTypeNotFound
		}
		return err
	}

	return s.activityTypeRepo.Delete(id)
}

