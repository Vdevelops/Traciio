package group

import (
	"errors"

	groupdomain "github.com/gilabs/crm-healthcare/api/internal/domain/group"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

var (
	ErrGroupNotFound      = errors.New("group not found")
	ErrGroupAlreadyExists = errors.New("group already exists")
	ErrGroupInUse         = errors.New("group is in use and cannot be deleted")
)

type Service struct {
	groupRepo interfaces.GroupRepository
	userRepo  interfaces.UserRepository
}

func NewService(groupRepo interfaces.GroupRepository, userRepo interfaces.UserRepository) *Service {
	return &Service{
		groupRepo: groupRepo,
		userRepo:  userRepo,
	}
}

// PaginationResult represents pagination result
type PaginationResult struct {
	Page       int
	PerPage    int
	Total      int
	TotalPages int
}

// List returns a list of groups with pagination
func (s *Service) List(req *groupdomain.ListGroupsRequest) ([]groupdomain.GroupResponse, *PaginationResult, error) {
	groups, total, err := s.groupRepo.List(req)
	if err != nil {
		return nil, nil, err
	}

	responses := make([]groupdomain.GroupResponse, len(groups))
	for i, g := range groups {
		responses[i] = *g.ToGroupResponse()
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}

	pagination := &PaginationResult{
		Page:       page,
		PerPage:    perPage,
		Total:      int(total),
		TotalPages: int((total + int64(perPage) - 1) / int64(perPage)),
	}

	return responses, pagination, nil
}

// GetByID returns a group by ID
func (s *Service) GetByID(id string) (*groupdomain.GroupResponse, error) {
	g, err := s.groupRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrGroupNotFound
		}
		return nil, err
	}
	return g.ToGroupResponse(), nil
}

// Create creates a new group
func (s *Service) Create(req *groupdomain.CreateGroupRequest) (*groupdomain.GroupResponse, error) {
	// Check if code already exists
	_, err := s.groupRepo.FindByCode(req.Code)
	if err == nil {
		return nil, ErrGroupAlreadyExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "active"
	}

	// Create group
	g := &groupdomain.Group{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      status,
	}

	if err := s.groupRepo.Create(g); err != nil {
		return nil, err
	}

	// Reload to get timestamps
	createdGroup, err := s.groupRepo.FindByID(g.ID)
	if err != nil {
		return nil, err
	}

	return createdGroup.ToGroupResponse(), nil
}

// Update updates a group
func (s *Service) Update(id string, req *groupdomain.UpdateGroupRequest) (*groupdomain.GroupResponse, error) {
	// Find existing group
	g, err := s.groupRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrGroupNotFound
		}
		return nil, err
	}

	// Check if code is being changed and if new code already exists
	if req.Code != "" && req.Code != g.Code {
		existing, err := s.groupRepo.FindByCode(req.Code)
		if err == nil && existing.ID != id {
			return nil, ErrGroupAlreadyExists
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	// Update fields
	if req.Name != "" {
		g.Name = req.Name
	}
	if req.Code != "" {
		g.Code = req.Code
	}
	if req.Description != "" {
		g.Description = req.Description
	}
	if req.Status != "" {
		g.Status = req.Status
	}

	if err := s.groupRepo.Update(g); err != nil {
		return nil, err
	}

	// Reload to get updated timestamps
	updatedGroup, err := s.groupRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return updatedGroup.ToGroupResponse(), nil
}

// Delete deletes a group
func (s *Service) Delete(id string) error {
	// Check if group exists
	_, err := s.groupRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrGroupNotFound
		}
		return err
	}

	// Check if group is in use
	count, err := s.groupRepo.CountUsersByGroupID(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrGroupInUse
	}

	// Delete group
	return s.groupRepo.Delete(id)
}

