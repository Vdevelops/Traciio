package user

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/brick"
	"github.com/gilabs/crm-healthcare/api/internal/domain/group"
	"github.com/gilabs/crm-healthcare/api/internal/domain/monthly_target"
	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
)

// MockUserRepository implements interfaces.UserRepository
type MockUserRepository struct {
	FindByIDFunc           func(id string) (*user.User, error)
	FindByEmailFunc        func(email string) (*user.User, error)
	ListFunc               func(req *user.ListUsersRequest) ([]user.User, int64, error)
	CreateFunc             func(user *user.User) error
	UpdateFunc             func(user *user.User) error
	DeleteFunc             func(id string) error
	CountUsersByRoleIDFunc func(roleID string) (int64, error)
	GetUsersByGroupIDFunc  func(groupID string) ([]user.User, error)
	GetUsersByBrickIDFunc  func(brickID string) ([]user.User, error)
	GetUsersByRoleIDFunc   func(roleID string) ([]string, error)
}

func (m *MockUserRepository) FindByID(id string) (*user.User, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}
func (m *MockUserRepository) FindByEmail(email string) (*user.User, error) {
	if m.FindByEmailFunc != nil {
		return m.FindByEmailFunc(email)
	}
	return nil, nil
}
func (m *MockUserRepository) List(req *user.ListUsersRequest) ([]user.User, int64, error) {
	if m.ListFunc != nil {
		return m.ListFunc(req)
	}
	return nil, 0, nil
}
func (m *MockUserRepository) Create(user *user.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(user)
	}
	return nil
}
func (m *MockUserRepository) Update(user *user.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(user)
	}
	return nil
}
func (m *MockUserRepository) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}
func (m *MockUserRepository) CountUsersByRoleID(roleID string) (int64, error) {
	if m.CountUsersByRoleIDFunc != nil {
		return m.CountUsersByRoleIDFunc(roleID)
	}
	return 0, nil
}
func (m *MockUserRepository) GetUsersByGroupID(groupID string) ([]user.User, error) {
	if m.GetUsersByGroupIDFunc != nil {
		return m.GetUsersByGroupIDFunc(groupID)
	}
	return nil, nil
}
func (m *MockUserRepository) GetUsersByBrickID(brickID string) ([]user.User, error) {
	if m.GetUsersByBrickIDFunc != nil {
		return m.GetUsersByBrickIDFunc(brickID)
	}
	return nil, nil
}
func (m *MockUserRepository) GetUsersByRoleID(roleID string) ([]string, error) {
	if m.GetUsersByRoleIDFunc != nil {
		return m.GetUsersByRoleIDFunc(roleID)
	}
	return nil, nil
}

// MockRoleRepository implements interfaces.RoleRepository
type MockRoleRepository struct {
	FindByIDFunc func(id string) (*role.Role, error)
	GetMobilePermissionsFunc func(roleID string, role *role.Role) (*role.GetMobilePermissionsResponse, error)
	AssignPermissionsFunc func(roleID string, permissionIDs []string) error
}
func (m *MockRoleRepository) FindByID(id string) (*role.Role, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}
func (m *MockRoleRepository) Create(role *role.Role) error { return nil }
func (m *MockRoleRepository) Update(role *role.Role) error { return nil }
func (m *MockRoleRepository) Delete(id string) error { return nil }
func (m *MockRoleRepository) FindByCode(code string) (*role.Role, error) { return nil, nil }
func (m *MockRoleRepository) List() ([]role.Role, error) { return nil, nil }
func (m *MockRoleRepository) AssignPermissions(roleID string, permissionIDs []string) error {
	if m.AssignPermissionsFunc != nil {
		return m.AssignPermissionsFunc(roleID, permissionIDs)
	}
	return nil
}
func (m *MockRoleRepository) GetPermissions(roleID string) ([]string, error) { return nil, nil }
func (m *MockRoleRepository) GetMobilePermissions(roleID string, role *role.Role) (*role.GetMobilePermissionsResponse, error) {
	if m.GetMobilePermissionsFunc != nil {
		return m.GetMobilePermissionsFunc(roleID, role)
	}
	return nil, nil
}
func (m *MockRoleRepository) UpdateMobilePermissions(roleID string, req *role.UpdateMobilePermissionsRequest) error { return nil }
func (m *MockRoleRepository) GetScopesByRoleID(roleID string) ([]role.RoleScope, error) {
	return nil, nil
}
func (m *MockRoleRepository) UpsertScopes(roleID string, scopes []role.RoleScopeItem) error {
	return nil
}

// MockGroupRepository implements interfaces.GroupRepository
type MockGroupRepository struct {
	FindByIDFunc func(id string) (*group.Group, error)
	CountUsersByGroupIDFunc func(groupID string) (int64, error)
}
func (m *MockGroupRepository) FindByID(id string) (*group.Group, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}
func (m *MockGroupRepository) Create(group *group.Group) error { return nil }
func (m *MockGroupRepository) Update(group *group.Group) error { return nil }
func (m *MockGroupRepository) Delete(id string) error { return nil }
func (m *MockGroupRepository) FindByCode(code string) (*group.Group, error) { return nil, nil }
func (m *MockGroupRepository) List(req *group.ListGroupsRequest) ([]group.Group, int64, error) { return nil, 0, nil }
func (m *MockGroupRepository) CountUsersByGroupID(groupID string) (int64, error) {
	if m.CountUsersByGroupIDFunc != nil {
		return m.CountUsersByGroupIDFunc(groupID)
	}
	return 0, nil
}

// MockBrickRepository implements interfaces.BrickRepository
type MockBrickRepository struct {
	FindByIDFunc func(id string) (*brick.Brick, error)
	FindByIDsFunc func(ids []string) ([]brick.Brick, error)
	CountSalesByBrickIDFunc func(brickID string) (int64, error)
}
func (m *MockBrickRepository) FindByID(id string) (*brick.Brick, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}
func (m *MockBrickRepository) FindByIDs(ids []string) ([]brick.Brick, error) {
	if m.FindByIDsFunc != nil {
		return m.FindByIDsFunc(ids)
	}
	return nil, nil
}
func (m *MockBrickRepository) Create(brick *brick.Brick) error { return nil }
func (m *MockBrickRepository) Update(brick *brick.Brick) error { return nil }
func (m *MockBrickRepository) Delete(id string) error { return nil }
func (m *MockBrickRepository) FindByCode(code string) (*brick.Brick, error) { return nil, nil }
func (m *MockBrickRepository) List(req *brick.ListBricksRequest) ([]brick.Brick, int64, error) { return nil, 0, nil }
func (m *MockBrickRepository) CountSalesByBrickID(brickID string) (int64, error) {
	if m.CountSalesByBrickIDFunc != nil {
		return m.CountSalesByBrickIDFunc(brickID)
	}
	return 0, nil
}
func (m *MockBrickRepository) GetSalesByBrickID(brickID string) ([]user.User, error) { return nil, nil }
func (m *MockBrickRepository) FindByRegencyAndProvince(regency, province string) (*brick.Brick, error) { return nil, nil }
func (m *MockBrickRepository) GetNextCodeSequence(prefix string) (int, error) { return 0, nil }


// MockMonthlyTargetRepository implements interfaces.MonthlyTargetRepository
type MockMonthlyTargetRepository struct {
	GetUserEffectiveTargetFunc func(userID string, year, month int) (*monthly_target.MonthlyTarget, error)
	BatchGetUserEffectiveTargetsFunc func(userIDs []string, year, month int) (map[string]*monthly_target.MonthlyTarget, error)
}
func (m *MockMonthlyTargetRepository) GetUserEffectiveTarget(userID string, year, month int) (*monthly_target.MonthlyTarget, error) {
	if m.GetUserEffectiveTargetFunc != nil {
		return m.GetUserEffectiveTargetFunc(userID, year, month)
	}
	return nil, nil
}
func (m *MockMonthlyTargetRepository) BatchGetUserEffectiveTargets(userIDs []string, year, month int) (map[string]*monthly_target.MonthlyTarget, error) {
	if m.BatchGetUserEffectiveTargetsFunc != nil {
		return m.BatchGetUserEffectiveTargetsFunc(userIDs, year, month)
	}
	return nil, nil
}
func (m *MockMonthlyTargetRepository) Create(target *monthly_target.MonthlyTarget) error { return nil }
func (m *MockMonthlyTargetRepository) Update(target *monthly_target.MonthlyTarget) error { return nil }
func (m *MockMonthlyTargetRepository) Delete(id string) error { return nil }
func (m *MockMonthlyTargetRepository) FindByID(id string) (*monthly_target.MonthlyTarget, error) { return nil, nil }
func (m *MockMonthlyTargetRepository) FindByUserAndPeriod(userID string, year int, month int) (*monthly_target.MonthlyTarget, error) { return nil, nil }
func (m *MockMonthlyTargetRepository) FindByGroupAndPeriod(groupID string, year int, month int) (*monthly_target.MonthlyTarget, error) { return nil, nil }
func (m *MockMonthlyTargetRepository) FindByBrickAndPeriod(brickID string, year int, month int) (*monthly_target.MonthlyTarget, error) { return nil, nil }
func (m *MockMonthlyTargetRepository) List(req *monthly_target.ListMonthlyTargetsRequest) ([]monthly_target.MonthlyTarget, int64, int64, error) { return nil, 0, 0, nil }
func (m *MockMonthlyTargetRepository) BatchGetProratedTargetsForPeriod(userIDs []string, startDate, endDate string) (map[string]float64, error) { return nil, nil }
func (m *MockMonthlyTargetRepository) GetTotalEffectiveTarget(year int, month int) (int64, error) { return 0, nil }
func (m *MockMonthlyTargetRepository) GetProratedTargetForPeriod(userID string, startDate, endDate string) (float64, error) { return 0, nil }

