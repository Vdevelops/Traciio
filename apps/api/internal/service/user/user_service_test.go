package user

import (
	"errors"
	"testing"

	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"gorm.io/gorm"
)

func setupTestService() (*Service, *MockUserRepository, *MockRoleRepository, *MockGroupRepository, *MockBrickRepository) {
	mockUserRepo := &MockUserRepository{}
	mockRoleRepo := &MockRoleRepository{}
	mockGroupRepo := &MockGroupRepository{}
	mockBrickRepo := &MockBrickRepository{}
	mockTargetRepo := &MockMonthlyTargetRepository{}
	
	svc := NewService(mockUserRepo, mockRoleRepo, mockGroupRepo, mockBrickRepo, mockTargetRepo, nil)
	return svc, mockUserRepo, mockRoleRepo, mockGroupRepo, mockBrickRepo
}

func TestService_Create_Success(t *testing.T) {
	svc, mockUserRepo, mockRoleRepo, _, _ := setupTestService()
	
	req := &user.CreateUserRequest{
		Email:    "new@example.com",
		Password: "password123",
		Name:     "New User",
		RoleID:   "role-1",
	}
	
	// Expectations
	mockRoleRepo.FindByIDFunc = func(id string) (*role.Role, error) {
		if id == "role-1" {
			return &role.Role{ID: "role-1", Code: "user", Name: "User"}, nil
		}
		return nil, gorm.ErrRecordNotFound
	}
	
	mockUserRepo.FindByEmailFunc = func(email string) (*user.User, error) {
		return nil, gorm.ErrRecordNotFound // User does not exist
	}
	
	mockUserRepo.CreateFunc = func(u *user.User) error {
		if u.Email != "new@example.com" {
			return errors.New("wrong email")
		}
		u.ID = "generated-id"
		return nil
	}
	
	mockUserRepo.FindByIDFunc = func(id string) (*user.User, error) {
		// Return user with Preloaded Role
		return &user.User{
			ID:     id,
			Email:  "new@example.com",
			Name:   "New User",
			RoleID: "role-1",
			Role:   &role.Role{ID: "role-1", Code: "user", Name: "User"},
			Status: "active",
		}, nil
	}
	
	// Execute
	resp, err := svc.Create(req)
	
	// Verify
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("response is nil")
	}
	if resp.Email != "new@example.com" {
		t.Errorf("expected email %s, got %s", "new@example.com", resp.Email)
	}
	if resp.Role == nil {
		t.Error("expected role to be present in response")
	}
}

func TestService_Create_RoleNotFound(t *testing.T) {
	svc, _, mockRoleRepo, _, _ := setupTestService()
	
	mockRoleRepo.FindByIDFunc = func(id string) (*role.Role, error) {
		return nil, gorm.ErrRecordNotFound
	}
	
	req := &user.CreateUserRequest{
		Email:  "new@example.com",
		RoleID: "invalid-role",
	}
	
	resp, err := svc.Create(req)
	
	if err == nil {
		t.Error("expected error, got nil")
	}
	if !errors.Is(err, ErrRoleNotFound) {
		t.Errorf("expected ErrRoleNotFound, got %v", err)
	}
	if resp != nil {
		t.Error("expected nil response")
	}
}

func TestService_Create_UserAlreadyExists(t *testing.T) {
	svc, mockUserRepo, mockRoleRepo, _, _ := setupTestService()
	
	mockRoleRepo.FindByIDFunc = func(id string) (*role.Role, error) {
		return &role.Role{ID: id}, nil
	}
	
	mockUserRepo.FindByEmailFunc = func(email string) (*user.User, error) {
		return &user.User{Email: email}, nil // Found
	}
	
	req := &user.CreateUserRequest{
		Email:  "exists@example.com",
		RoleID: "role-1",
	}
	
	_, err := svc.Create(req)
	
	if err == nil {
		t.Error("expected error, got nil")
	}
	if !errors.Is(err, ErrUserAlreadyExists) {
		t.Errorf("expected ErrUserAlreadyExists, got %v", err)
	}
}
