package brick

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/brick"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
)

// MockBrickRepository
type MockBrickRepository struct {
	FindByIDFunc                 func(id string) (*brick.Brick, error)
	FindByIDsFunc                func(ids []string) ([]brick.Brick, error)
	FindByCodeFunc               func(code string) (*brick.Brick, error)
	ListFunc                     func(req *brick.ListBricksRequest) ([]brick.Brick, int64, error)
	CreateFunc                   func(brick *brick.Brick) error
	UpdateFunc                   func(brick *brick.Brick) error
	DeleteFunc                   func(id string) error
	CountSalesByBrickIDFunc      func(brickID string) (int64, error)
	GetSalesByBrickIDFunc        func(brickID string) ([]user.User, error)
	FindByRegencyAndProvinceFunc func(regency, province string) (*brick.Brick, error)
}

func (m *MockBrickRepository) FindByID(id string) (*brick.Brick, error) {
	if m.FindByIDFunc != nil { return m.FindByIDFunc(id) }
	return nil, nil
}
func (m *MockBrickRepository) FindByIDs(ids []string) ([]brick.Brick, error) {
	if m.FindByIDsFunc != nil { return m.FindByIDsFunc(ids) }
	return nil, nil
}
func (m *MockBrickRepository) FindByCode(code string) (*brick.Brick, error) {
	if m.FindByCodeFunc != nil { return m.FindByCodeFunc(code) }
	return nil, nil
}
func (m *MockBrickRepository) List(req *brick.ListBricksRequest) ([]brick.Brick, int64, error) {
	if m.ListFunc != nil { return m.ListFunc(req) }
	return nil, 0, nil
}
func (m *MockBrickRepository) Create(brick *brick.Brick) error {
	if m.CreateFunc != nil { return m.CreateFunc(brick) }
	return nil
}
func (m *MockBrickRepository) Update(brick *brick.Brick) error {
	if m.UpdateFunc != nil { return m.UpdateFunc(brick) }
	return nil
}
func (m *MockBrickRepository) Delete(id string) error {
	if m.DeleteFunc != nil { return m.DeleteFunc(id) }
	return nil
}
func (m *MockBrickRepository) CountSalesByBrickID(brickID string) (int64, error) {
	if m.CountSalesByBrickIDFunc != nil { return m.CountSalesByBrickIDFunc(brickID) }
	return 0, nil
}
func (m *MockBrickRepository) GetSalesByBrickID(brickID string) ([]user.User, error) {
	if m.GetSalesByBrickIDFunc != nil { return m.GetSalesByBrickIDFunc(brickID) }
	return nil, nil
}
func (m *MockBrickRepository) FindByRegencyAndProvince(regency, province string) (*brick.Brick, error) {
	if m.FindByRegencyAndProvinceFunc != nil { return m.FindByRegencyAndProvinceFunc(regency, province) }
	return nil, nil
}
func (m *MockBrickRepository) GetNextCodeSequence(prefix string) (int, error) {
	return 0, nil
}

// MockUserRepository - Minimal implementation for BrickService
type MockUserRepository struct {
	FindByIDFunc func(id string) (*user.User, error)
}
func (m *MockUserRepository) FindByID(id string) (*user.User, error) {
	if m.FindByIDFunc != nil { return m.FindByIDFunc(id) }
	return nil, nil
}
func (m *MockUserRepository) FindByEmail(email string) (*user.User, error) { return nil, nil }
func (m *MockUserRepository) List(req *user.ListUsersRequest) ([]user.User, int64, error) { return nil, 0, nil }
func (m *MockUserRepository) Create(user *user.User) error { return nil }
func (m *MockUserRepository) Update(user *user.User) error { return nil }
func (m *MockUserRepository) Delete(id string) error { return nil }
func (m *MockUserRepository) CountUsersByRoleID(roleID string) (int64, error) { return 0, nil }
func (m *MockUserRepository) GetUsersByGroupID(groupID string) ([]user.User, error) { return nil, nil }
func (m *MockUserRepository) GetUsersByBrickID(brickID string) ([]user.User, error) { return nil, nil }
func (m *MockUserRepository) GetUsersByRoleID(roleID string) ([]string, error) { return nil, nil }
