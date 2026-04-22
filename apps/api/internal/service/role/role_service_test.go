package role

import (
	"errors"
	"testing"

	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"gorm.io/gorm"
)

func setupTestService() (*Service, *MockRoleRepository, *MockUserRepository) {
	mockRoleRepo := &MockRoleRepository{}
	mockUserRepo := &MockUserRepository{}
	
	svc := NewService(mockRoleRepo, mockUserRepo)
	return svc, mockRoleRepo, mockUserRepo
}

func TestService_Create_Success(t *testing.T) {
	svc, mockRoleRepo, _ := setupTestService()
	
	req := &role.CreateRoleRequest{
		Name: "New Role",
		Code: "new_role",
	}
	
	mockRoleRepo.FindByCodeFunc = func(code string) (*role.Role, error) {
		return nil, gorm.ErrRecordNotFound
	}
	
	mockRoleRepo.CreateFunc = func(r *role.Role) error {
		r.ID = "generated-id"
		return nil
	}
	
	mockRoleRepo.FindByIDFunc = func(id string) (*role.Role, error) {
		return &role.Role{ID: id, Code: req.Code, Name: req.Name, Status: "active"}, nil
	}
	
	resp, err := svc.Create(req)
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Code != "new_role" {
		t.Errorf("expected code new_role, got %s", resp.Code)
	}
}

func TestService_Create_AlreadyExists(t *testing.T) {
	svc, mockRoleRepo, _ := setupTestService()
	
	mockRoleRepo.FindByCodeFunc = func(code string) (*role.Role, error) {
		return &role.Role{Code: code}, nil
	}
	
	req := &role.CreateRoleRequest{Code: "existing"}
	
	_, err := svc.Create(req)
	
	if !errors.Is(err, ErrRoleAlreadyExists) {
		t.Errorf("expected ErrRoleAlreadyExists, got %v", err)
	}
}

func TestService_Delete_Success(t *testing.T) {
	svc, mockRoleRepo, mockUserRepo := setupTestService()
	
	mockRoleRepo.FindByIDFunc = func(id string) (*role.Role, error) {
		return &role.Role{ID: id, IsProtected: false}, nil
	}
	
	mockUserRepo.CountUsersByRoleIDFunc = func(roleID string) (int64, error) {
		return 0, nil
	}
	
	mockRoleRepo.DeleteFunc = func(id string) error {
		return nil
	}
	
	err := svc.Delete("role-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestService_Delete_Protected(t *testing.T) {
	svc, mockRoleRepo, _ := setupTestService()
	
	mockRoleRepo.FindByIDFunc = func(id string) (*role.Role, error) {
		return &role.Role{ID: id, IsProtected: true}, nil
	}
	
	err := svc.Delete("role-1")
	if !errors.Is(err, ErrRoleProtected) {
		t.Errorf("expected ErrRoleProtected, got %v", err)
	}
}

func TestService_Delete_InUse(t *testing.T) {
	svc, mockRoleRepo, mockUserRepo := setupTestService()
	
	mockRoleRepo.FindByIDFunc = func(id string) (*role.Role, error) {
		return &role.Role{ID: id, IsProtected: false}, nil
	}
	
	mockUserRepo.CountUsersByRoleIDFunc = func(roleID string) (int64, error) {
		return 5, nil
	}
	
	err := svc.Delete("role-1")
	if !errors.Is(err, ErrRoleInUse) {
		t.Errorf("expected ErrRoleInUse, got %v", err)
	}
}
