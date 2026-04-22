package role

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
)

// MockRoleRepository implements interfaces.RoleRepository
type MockRoleRepository struct {
	FindByIDFunc func(id string) (*role.Role, error)
	FindByCodeFunc func(code string) (*role.Role, error)
	ListFunc func() ([]role.Role, error) // Changed based on service List() call
	CreateFunc func(role *role.Role) error
	UpdateFunc func(role *role.Role) error
	DeleteFunc func(id string) error
	AssignPermissionsFunc func(roleID string, permissionIDs []string) error
	GetPermissionsFunc func(roleID string) ([]string, error)
	GetMobilePermissionsFunc func(roleID string, role *role.Role) (*role.GetMobilePermissionsResponse, error)
	UpdateMobilePermissionsFunc func(roleID string, req *role.UpdateMobilePermissionsRequest) error
}

func (m *MockRoleRepository) FindByID(id string) (*role.Role, error) {
	if m.FindByIDFunc != nil { return m.FindByIDFunc(id) }
	return nil, nil
}
func (m *MockRoleRepository) FindByCode(code string) (*role.Role, error) {
	if m.FindByCodeFunc != nil { return m.FindByCodeFunc(code) }
	return nil, nil
}
func (m *MockRoleRepository) List() ([]role.Role, error) {  // Changed return type to value not pointer to match interface if that was the case
	// Wait, checking service.go usage: `roles, err := s.roleRepo.List()` then `for i, r := range roles` and `r.ToRoleResponse()`.
	// Usually repo returns `[]*role.Role`. Let's check interface definition in step 151.
	// Interface says: `List() ([]role.Role, error)` (Slice of values, not pointers? Let's re-verify Step 151)
	// Step 151: `List() ([]role.Role, error)` -> YES, slice of values.
	// But service.go lines 46-53: `for i, r := range roles` -> `r` is `role.Role`. `r.ToRoleResponse()` is method on `*Role`?
	// If `ToRoleResponse` is pointer receiver, we might need `&r`.
	// Let's stick to interface definition for now.
	if m.ListFunc != nil { return m.ListFunc() }
	return nil, nil
}
func (m *MockRoleRepository) Create(role *role.Role) error {
	if m.CreateFunc != nil { return m.CreateFunc(role) }
	return nil
}
func (m *MockRoleRepository) Update(role *role.Role) error {
	if m.UpdateFunc != nil { return m.UpdateFunc(role) }
	return nil
}
func (m *MockRoleRepository) Delete(id string) error {
	if m.DeleteFunc != nil { return m.DeleteFunc(id) }
	return nil
}
func (m *MockRoleRepository) AssignPermissions(roleID string, permissionIDs []string) error {
	if m.AssignPermissionsFunc != nil { return m.AssignPermissionsFunc(roleID, permissionIDs) }
	return nil
}
func (m *MockRoleRepository) GetPermissions(roleID string) ([]string, error) {
	if m.GetPermissionsFunc != nil { return m.GetPermissionsFunc(roleID) }
	return nil, nil
}
func (m *MockRoleRepository) GetMobilePermissions(roleID string, role *role.Role) (*role.GetMobilePermissionsResponse, error) {
	if m.GetMobilePermissionsFunc != nil { return m.GetMobilePermissionsFunc(roleID, role) }
	return nil, nil
}
func (m *MockRoleRepository) UpdateMobilePermissions(roleID string, req *role.UpdateMobilePermissionsRequest) error {
	if m.UpdateMobilePermissionsFunc != nil { return m.UpdateMobilePermissionsFunc(roleID, req) }
	return nil
}
func (m *MockRoleRepository) GetScopesByRoleID(roleID string) ([]role.RoleScope, error) {
	return nil, nil
}
func (m *MockRoleRepository) UpsertScopes(roleID string, scopes []role.RoleScopeItem) error {
	return nil
}

// MockUserRepository implements interfaces.UserRepository (subset needed for RoleService)
type MockUserRepository struct {
	CountUsersByRoleIDFunc func(roleID string) (int64, error)
	GetUsersByRoleIDFunc func(roleID string) ([]string, error)
	// Other methods needed for interface compliance
	FindByIDFunc           func(id string) (*user.User, error)
	FindByEmailFunc        func(email string) (*user.User, error)
	ListFunc               func(req *user.ListUsersRequest) ([]user.User, int64, error)
	CreateFunc             func(user *user.User) error
	UpdateFunc             func(user *user.User) error
	DeleteFunc             func(id string) error
	GetUsersByGroupIDFunc  func(groupID string) ([]user.User, error)
	GetUsersByBrickIDFunc  func(brickID string) ([]user.User, error)
}
func (m *MockUserRepository) CountUsersByRoleID(roleID string) (int64, error) {
	if m.CountUsersByRoleIDFunc != nil { return m.CountUsersByRoleIDFunc(roleID) }
	return 0, nil
}
func (m *MockUserRepository) GetUsersByRoleID(roleID string) ([]string, error) {
	if m.GetUsersByRoleIDFunc != nil { return m.GetUsersByRoleIDFunc(roleID) }
	return nil, nil
}
// Stubs for interface compliance
func (m *MockUserRepository) FindByID(id string) (*user.User, error) { return nil, nil }
func (m *MockUserRepository) FindByEmail(email string) (*user.User, error) { return nil, nil }
func (m *MockUserRepository) List(req *user.ListUsersRequest) ([]user.User, int64, error) { return nil, 0, nil }
func (m *MockUserRepository) Create(user *user.User) error { return nil }
func (m *MockUserRepository) Update(user *user.User) error { return nil }
func (m *MockUserRepository) Delete(id string) error { return nil }
func (m *MockUserRepository) GetUsersByGroupID(groupID string) ([]user.User, error) { return nil, nil }
func (m *MockUserRepository) GetUsersByBrickID(brickID string) ([]user.User, error) { return nil, nil }

