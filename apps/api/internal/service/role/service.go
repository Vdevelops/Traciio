package role

import (
	"errors"

	"github.com/gilabs/crm-healthcare/api/internal/domain/permission"
	roledomain "github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"github.com/gilabs/crm-healthcare/api/internal/hub"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

var (
	ErrRoleNotFound      = errors.New("role not found")
	ErrRoleAlreadyExists = errors.New("role already exists")
	ErrRoleProtected     = errors.New("role is protected and cannot be modified")
	ErrRoleInUse         = errors.New("role is in use and cannot be deleted")
)

type Service struct {
	roleRepo        interfaces.RoleRepository
	userRepo        interfaces.UserRepository
	notificationHub *hub.NotificationHub
	permService     interface{ InvalidateRoleCache(roleCode string) error }
}

func NewService(roleRepo interfaces.RoleRepository, userRepo interfaces.UserRepository) *Service {
	return &Service{
		roleRepo: roleRepo,
		userRepo: userRepo,
	}
}

// SetNotificationHub sets the notification hub for broadcasting permission updates.
func (s *Service) SetNotificationHub(h *hub.NotificationHub) {
	s.notificationHub = h
}

// SetPermissionService sets the permission service for role cache invalidation.
func (s *Service) SetPermissionService(ps interface{ InvalidateRoleCache(roleCode string) error }) {
	s.permService = ps
}

// List returns a list of roles
func (s *Service) List() ([]roledomain.RoleResponse, error) {
	roles, err := s.roleRepo.List()
	if err != nil {
		return nil, err
	}

	responses := make([]roledomain.RoleResponse, len(roles))
	for i, r := range roles {
		responses[i] = *r.ToRoleResponse()
	}

	return responses, nil
}

// GetByID returns a role by ID
func (s *Service) GetByID(id string) (*roledomain.RoleResponse, error) {
	r, err := s.roleRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}
	return r.ToRoleResponse(), nil
}

// Create creates a new role
func (s *Service) Create(req *roledomain.CreateRoleRequest) (*roledomain.RoleResponse, error) {
	// Check if code already exists
	_, err := s.roleRepo.FindByCode(req.Code)
	if err == nil {
		return nil, ErrRoleAlreadyExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "active"
	}

	// Set default mobile_access
	mobileAccess := false
	if req.MobileAccess != nil {
		mobileAccess = *req.MobileAccess
	}

	// Create role
	r := &roledomain.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      status,
		MobileAccess: mobileAccess,
	}

	if err := s.roleRepo.Create(r); err != nil {
		return nil, err
	}

	// Reload with permissions
	createdRole, err := s.roleRepo.FindByID(r.ID)
	if err != nil {
		return nil, err
	}

	return createdRole.ToRoleResponse(), nil
}

// Update updates a role
func (s *Service) Update(id string, req *roledomain.UpdateRoleRequest) (*roledomain.RoleResponse, error) {
	// Find role
	r, err := s.roleRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}

	// Check if role is protected
	if r.IsProtected {
		// Protected roles cannot have status changed to inactive or code changed
		if req.Status != "" && req.Status != "active" && req.Status != r.Status {
			return nil, ErrRoleProtected
		}
		if req.Code != "" && req.Code != r.Code {
			return nil, ErrRoleProtected
		}
	}

	// Update fields
	if req.Name != "" {
		r.Name = req.Name
	}

	if req.Code != "" {
		// Check if code already exists (excluding current role)
		existingRole, err := s.roleRepo.FindByCode(req.Code)
		if err == nil && existingRole.ID != id {
			return nil, ErrRoleAlreadyExists
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		r.Code = req.Code
	}

	if req.Description != "" {
		r.Description = req.Description
	}

	if req.Status != "" {
		r.Status = req.Status
	}

	if req.MobileAccess != nil {
		r.MobileAccess = *req.MobileAccess
	}

	if err := s.roleRepo.Update(r); err != nil {
		return nil, err
	}

	// Reload with permissions
	updatedRole, err := s.roleRepo.FindByID(r.ID)
	if err != nil {
		return nil, err
	}

	return updatedRole.ToRoleResponse(), nil
}

// Delete deletes a role
func (s *Service) Delete(id string) error {
	// Check if role exists
	r, err := s.roleRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRoleNotFound
		}
		return err
	}

	// Check if role is protected
	if r.IsProtected {
		return ErrRoleProtected
	}

	// Check if role is in use by any users
	userCount, err := s.userRepo.CountUsersByRoleID(id)
	if err != nil {
		return err
	}
	if userCount > 0 {
		return ErrRoleInUse
	}

	return s.roleRepo.Delete(id)
}

// GetRolePermissions returns permissions for a role
func (s *Service) GetRolePermissions(roleID string) ([]permission.PermissionSimpleResponse, error) {
	// Check if role exists
	role, err := s.roleRepo.FindByID(roleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}

	var resp []permission.PermissionSimpleResponse
	for _, p := range role.Permissions {
		resp = append(resp, *p.ToPermissionSimpleResponse())
	}
	return resp, nil
}

// AssignPermissions assigns permissions to a role
func (s *Service) AssignPermissions(roleID string, permissionIDs []string) error {
	// Check if role exists
	role, err := s.roleRepo.FindByID(roleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRoleNotFound
		}
		return err
	}

	if err := s.roleRepo.AssignPermissions(roleID, permissionIDs); err != nil {
		return err
	}

	// Invalidate permission cache for this role
	if s.permService != nil {
		if err := s.permService.InvalidateRoleCache(role.Code); err != nil {
			// Log error but don't fail the request
			// In production you might want structured logging here
			// fmt.Printf("Failed to invalidate role cache: %v\n", err)
		}
	}

	// Broadcast permissions update to all users with this role
	s.broadcastPermissionsUpdateToRole(roleID)

	return nil
}

// broadcastPermissionsUpdateToRole broadcasts permission updates to all users with the given role
// and invalidates their permission cache.
func (s *Service) broadcastPermissionsUpdateToRole(roleID string) {
	userIDs, err := s.userRepo.GetUsersByRoleID(roleID)
	if err != nil {
		return
	}

	if s.notificationHub != nil {
		s.notificationHub.BroadcastPermissionsUpdateToMultipleUsers(userIDs, roleID)
	}
}

// GetMobilePermissions returns mobile permissions for a role
func (s *Service) GetMobilePermissions(roleID string) (*roledomain.GetMobilePermissionsResponse, error) {
	// Check if role exists
	r, err := s.roleRepo.FindByID(roleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}

	return s.roleRepo.GetMobilePermissions(roleID, r)
}

// UpdateMobilePermissions updates mobile permissions for a role
func (s *Service) UpdateMobilePermissions(roleID string, req *roledomain.UpdateMobilePermissionsRequest) error {
	// Check if role exists
	_, err := s.roleRepo.FindByID(roleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRoleNotFound
		}
		return err
	}

	if err := s.roleRepo.UpdateMobilePermissions(roleID, req); err != nil {
		return err
	}

	// Broadcast permissions update to all users with this role
	s.broadcastPermissionsUpdateToRole(roleID)

	return nil
}

