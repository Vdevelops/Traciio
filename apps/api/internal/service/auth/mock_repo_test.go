package auth

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/permission"
	"github.com/gilabs/crm-healthcare/api/internal/domain/refresh_token"
	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"github.com/redis/go-redis/v9"
)

// MockAuthRepository implements interfaces.AuthRepository
type MockAuthRepository struct {
	FindByEmailFunc func(email string) (*user.User, error)
	FindByIDFunc    func(id string) (*user.User, error)
	CreateFunc      func(user *user.User) error
	UpdateFunc      func(user *user.User) error
}

func (m *MockAuthRepository) FindByEmail(email string) (*user.User, error) {
	if m.FindByEmailFunc != nil {
		return m.FindByEmailFunc(email)
	}
	return nil, nil
}

func (m *MockAuthRepository) FindByID(id string) (*user.User, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}

func (m *MockAuthRepository) Create(user *user.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(user)
	}
	return nil
}

func (m *MockAuthRepository) Update(user *user.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(user)
	}
	return nil
}

// MockRefreshTokenRepository implements interfaces.RefreshTokenRepository
type MockRefreshTokenRepository struct {
	CreateFunc         func(token *refresh_token.RefreshToken) error
	FindByTokenIDFunc  func(tokenID string) (*refresh_token.RefreshToken, error)
	FindByUserIDFunc   func(userID string) ([]*refresh_token.RefreshToken, error)
	RevokeFunc         func(tokenID string) error
	RevokeByUserIDFunc func(userID string) error
	DeleteExpiredFunc  func() error
}

func (m *MockRefreshTokenRepository) Create(token *refresh_token.RefreshToken) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(token)
	}
	return nil
}

func (m *MockRefreshTokenRepository) FindByTokenID(tokenID string) (*refresh_token.RefreshToken, error) {
	if m.FindByTokenIDFunc != nil {
		return m.FindByTokenIDFunc(tokenID)
	}
	return nil, nil
}

func (m *MockRefreshTokenRepository) FindByUserID(userID string) ([]*refresh_token.RefreshToken, error) {
	if m.FindByUserIDFunc != nil {
		return m.FindByUserIDFunc(userID)
	}
	return nil, nil
}

func (m *MockRefreshTokenRepository) Revoke(tokenID string) error {
	if m.RevokeFunc != nil {
		return m.RevokeFunc(tokenID)
	}
	return nil
}

func (m *MockRefreshTokenRepository) RevokeByUserID(userID string) error {
	if m.RevokeByUserIDFunc != nil {
		return m.RevokeByUserIDFunc(userID)
	}
	return nil
}

func (m *MockRefreshTokenRepository) DeleteExpired() error {
	if m.DeleteExpiredFunc != nil {
		return m.DeleteExpiredFunc()
	}
	return nil
}

// MockRoleRepository implements interfaces.RoleRepository (needed for PermissionService)
type MockRoleRepository struct {
	CreateFunc           func(role *role.Role) error
	UpdateFunc           func(role *role.Role) error
	DeleteFunc           func(id string) error
	FindByIDFunc         func(id string) (*role.Role, error)
	FindByCodeFunc       func(code string) (*role.Role, error)
	ListFunc             func() ([]*role.Role, error)
	AssignPermissionFunc func(roleID string, permissionIDs []string) error
	RemovePermissionFunc func(roleID string, permissionIDs []string) error
	GetMobilePermissionsFunc func(roleID string, role *role.Role) (*role.GetMobilePermissionsResponse, error)
}

func (m *MockRoleRepository) FindByCode(code string) (*role.Role, error) {
	if m.FindByCodeFunc != nil {
		return m.FindByCodeFunc(code)
	}
	return nil, nil
}
// Implementation placeholders for other Interface methods to satisfy compiler
func (m *MockRoleRepository) Create(role *role.Role) error { return nil }
func (m *MockRoleRepository) Update(role *role.Role) error { return nil }
func (m *MockRoleRepository) Delete(id string) error { return nil }
func (m *MockRoleRepository) FindByID(id string) (*role.Role, error) { return nil, nil }
func (m *MockRoleRepository) List() ([]*role.Role, error) { return nil, nil }
func (m *MockRoleRepository) AssignPermission(roleID string, permissionIDs []string) error { return nil }
func (m *MockRoleRepository) RemovePermission(roleID string, permissionIDs []string) error { return nil }
func (m *MockRoleRepository) GetMobilePermissions(roleID string, role *role.Role) (*role.GetMobilePermissionsResponse, error) { return nil, nil }


// MockRedisClient (minimal)
type MockRedisClient struct {
	redis.UniversalClient
}
// We handle Redis by passing nil in tests if not needed, or mocking if needed. 
// Since PermissionService checks if s.redisClient != nil, we can pass nil to skip redis logic.

// MockPermissionRepository (needed for PermissionService constructor)
type MockPermissionRepository struct {
	// Add functionality if needed
}
func (m *MockPermissionRepository) Create(perm *permission.Permission) error { return nil }
func (m *MockPermissionRepository) FindByID(id string) (*permission.Permission, error) { return nil, nil }
func (m *MockPermissionRepository) List() ([]*permission.Permission, error) { return nil, nil }
func (m *MockPermissionRepository) Update(perm *permission.Permission) error { return nil }
func (m *MockPermissionRepository) Delete(id string) error { return nil }
func (m *MockPermissionRepository) FindByCode(code string) (*permission.Permission, error) { return nil, nil } // Added this if interface requires it, checking file list again might be good but let's assume standard CRUD.

// MockUserRepository (needed for PermissionService constructor)
type MockUserRepository struct {
	// Embed MockAuthRepository as they share interface? No, UserRepo interface is likely bigger.
	// Implementing minimal for compiler satiation.
	FindByIDFunc func(id string) (*user.User, error)
}
func (m *MockUserRepository) FindByID(id string) (*user.User, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}
func (m *MockUserRepository) FindByEmail(email string) (*user.User, error) { return nil, nil }
func (m *MockUserRepository) Create(user *user.User) error { return nil }
func (m *MockUserRepository) Update(user *user.User) error { return nil }
func (m *MockUserRepository) ChangePassword(userID, hashedPassword string) error { return nil }
func (m *MockUserRepository) List(req *user.ListUsersRequest) ([]user.User, int64, error) { return nil, 0, nil }
func (m *MockUserRepository) Delete(id string) error { return nil }
func (m *MockUserRepository) GetUsersByRoleID(roleID string) ([]string, error) { return nil, nil }
func (m *MockUserRepository) CountUsersByRoleID(roleID string) (int64, error) { return 0, nil }
func (m *MockUserRepository) GetUsersByGroupID(groupID string) ([]user.User, error) { return nil, nil }
func (m *MockUserRepository) GetUsersByBrickID(brickID string) ([]user.User, error) { return nil, nil }

