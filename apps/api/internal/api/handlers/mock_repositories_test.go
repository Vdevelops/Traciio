package handlers

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/account"
	"github.com/gilabs/crm-healthcare/api/internal/domain/activity"
	"github.com/gilabs/crm-healthcare/api/internal/domain/brick"
	"github.com/gilabs/crm-healthcare/api/internal/domain/category"
	"github.com/gilabs/crm-healthcare/api/internal/domain/contact"
	"github.com/gilabs/crm-healthcare/api/internal/domain/contact_role"
	"github.com/gilabs/crm-healthcare/api/internal/domain/deal_history"
	"github.com/gilabs/crm-healthcare/api/internal/domain/group"
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead"
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead_source"
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead_status"
	"github.com/gilabs/crm-healthcare/api/internal/domain/monthly_target"
	"github.com/gilabs/crm-healthcare/api/internal/domain/permission"
	"github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
	"github.com/gilabs/crm-healthcare/api/internal/domain/product"
	"github.com/gilabs/crm-healthcare/api/internal/domain/refresh_token"
	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"github.com/gilabs/crm-healthcare/api/internal/domain/task"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"github.com/gilabs/crm-healthcare/api/internal/domain/visit_report"
)

// MockAuthRepository implements interfaces.AuthRepository
type MockAuthRepository struct {
	FindByEmailFunc func(email string) (*user.User, error)
	FindByIDFunc    func(id string) (*user.User, error)
	CreateFunc      func(user *user.User) error
	UpdateFunc      func(user *user.User) error
	// Helper methods that might be needed by UserService if it uses same repo interface
	// UserService uses interfaces.UserRepository which includes AuthRepository methods plus more.
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
func (m *MockAuthRepository) Create(user *user.User) error { return nil }
func (m *MockAuthRepository) Update(user *user.User) error { return nil }

// MockRefreshTokenRepository
type MockRefreshTokenRepository struct {
	CreateFunc func(token *refresh_token.RefreshToken) error
}

func (m *MockRefreshTokenRepository) Create(token *refresh_token.RefreshToken) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(token)
	}
	return nil
}
func (m *MockRefreshTokenRepository) FindByTokenID(tokenID string) (*refresh_token.RefreshToken, error) {
	return nil, nil
}
func (m *MockRefreshTokenRepository) FindByUserID(userID string) ([]*refresh_token.RefreshToken, error) {
	return nil, nil
}
func (m *MockRefreshTokenRepository) Revoke(tokenID string) error        { return nil }
func (m *MockRefreshTokenRepository) RevokeByUserID(userID string) error { return nil }
func (m *MockRefreshTokenRepository) DeleteExpired() error               { return nil }

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
func (m *MockUserRepository) Delete(id string) error                                { return nil }
func (m *MockUserRepository) CountUsersByRoleID(roleID string) (int64, error)       { return 0, nil }
func (m *MockUserRepository) GetUsersByGroupID(groupID string) ([]user.User, error) { return nil, nil }
func (m *MockUserRepository) GetUsersByBrickID(brickID string) ([]user.User, error) { return nil, nil }
func (m *MockUserRepository) GetUsersByRoleID(roleID string) ([]string, error)      { return nil, nil }

// MockRoleRepository
type MockRoleRepository struct {
	FindByIDFunc          func(id string) (*role.Role, error)
	AssignPermissionsFunc func(roleID string, permissionIDs []string) error
	FindByCodeFunc        func(code string) (*role.Role, error)
	CreateFunc            func(role *role.Role) error
	ListFunc              func() ([]role.Role, error)
}

func (m *MockRoleRepository) FindByID(id string) (*role.Role, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}
func (m *MockRoleRepository) Create(role *role.Role) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(role)
	}
	return nil
}
func (m *MockRoleRepository) Update(role *role.Role) error { return nil }
func (m *MockRoleRepository) Delete(id string) error       { return nil }
func (m *MockRoleRepository) FindByCode(code string) (*role.Role, error) {
	if m.FindByCodeFunc != nil {
		return m.FindByCodeFunc(code)
	}
	return nil, nil
}
func (m *MockRoleRepository) List() ([]role.Role, error) {
	if m.ListFunc != nil {
		return m.ListFunc()
	}
	return nil, nil
}
func (m *MockRoleRepository) AssignPermissions(roleID string, permissionIDs []string) error {
	return nil
}
func (m *MockRoleRepository) GetPermissions(roleID string) ([]string, error) { return nil, nil }
func (m *MockRoleRepository) GetScopesByRoleID(roleID string) ([]role.RoleScope, error) {
	return nil, nil
}
func (m *MockRoleRepository) UpsertScopes(roleID string, scopes []role.RoleScopeItem) error {
	return nil
}

// MockGroupRepository
type MockGroupRepository struct {
	FindByIDFunc            func(id string) (*group.Group, error)
	FindByCodeFunc          func(code string) (*group.Group, error)
	ListFunc                func(req *group.ListGroupsRequest) ([]group.Group, int64, error)
	CreateFunc              func(group *group.Group) error
	UpdateFunc              func(group *group.Group) error
	DeleteFunc              func(id string) error
	CountUsersByGroupIDFunc func(groupID string) (int64, error)
}

func (m *MockGroupRepository) FindByID(id string) (*group.Group, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}
func (m *MockGroupRepository) Create(group *group.Group) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(group)
	}
	return nil
}
func (m *MockGroupRepository) Update(group *group.Group) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(group)
	}
	return nil
}
func (m *MockGroupRepository) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}
func (m *MockGroupRepository) FindByCode(code string) (*group.Group, error) {
	if m.FindByCodeFunc != nil {
		return m.FindByCodeFunc(code)
	}
	return nil, nil
}
func (m *MockGroupRepository) List(req *group.ListGroupsRequest) ([]group.Group, int64, error) {
	if m.ListFunc != nil {
		return m.ListFunc(req)
	}
	return nil, 0, nil
}
func (m *MockGroupRepository) CountUsersByGroupID(groupID string) (int64, error) {
	if m.CountUsersByGroupIDFunc != nil {
		return m.CountUsersByGroupIDFunc(groupID)
	}
	return 0, nil
}

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
func (m *MockBrickRepository) Create(brick *brick.Brick) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(brick)
	}
	return nil
}
func (m *MockBrickRepository) Update(brick *brick.Brick) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(brick)
	}
	return nil
}
func (m *MockBrickRepository) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}
func (m *MockBrickRepository) FindByCode(code string) (*brick.Brick, error) {
	if m.FindByCodeFunc != nil {
		return m.FindByCodeFunc(code)
	}
	return nil, nil
}
func (m *MockBrickRepository) List(req *brick.ListBricksRequest) ([]brick.Brick, int64, error) {
	if m.ListFunc != nil {
		return m.ListFunc(req)
	}
	return nil, 0, nil
}
func (m *MockBrickRepository) CountSalesByBrickID(brickID string) (int64, error) {
	if m.CountSalesByBrickIDFunc != nil {
		return m.CountSalesByBrickIDFunc(brickID)
	}
	return 0, nil
}
func (m *MockBrickRepository) GetSalesByBrickID(brickID string) ([]user.User, error) {
	if m.GetSalesByBrickIDFunc != nil {
		return m.GetSalesByBrickIDFunc(brickID)
	}
	return nil, nil
}
func (m *MockBrickRepository) FindByRegencyAndProvince(regency, province string) (*brick.Brick, error) {
	if m.FindByRegencyAndProvinceFunc != nil {
		return m.FindByRegencyAndProvinceFunc(regency, province)
	}
	return nil, nil
}
func (m *MockBrickRepository) GetNextCodeSequence(prefix string) (int, error) {
	return 0, nil
}

// MockMonthlyTargetRepository
type MockMonthlyTargetRepository struct {
	GetUserEffectiveTargetFunc       func(userID string, year, month int) (*monthly_target.MonthlyTarget, error)
	BatchGetUserEffectiveTargetsFunc func(userIDs []string, year, month int) (map[string]*monthly_target.MonthlyTarget, error)
}

func (m *MockMonthlyTargetRepository) GetUserEffectiveTarget(userID string, year, month int) (*monthly_target.MonthlyTarget, error) {
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
func (m *MockMonthlyTargetRepository) Delete(id string) error                            { return nil }
func (m *MockMonthlyTargetRepository) FindByID(id string) (*monthly_target.MonthlyTarget, error) {
	return nil, nil
}
func (m *MockMonthlyTargetRepository) FindByUserAndPeriod(userID string, year, month int) (*monthly_target.MonthlyTarget, error) {
	return nil, nil
}
func (m *MockMonthlyTargetRepository) FindByGroupAndPeriod(groupID string, year, month int) (*monthly_target.MonthlyTarget, error) {
	return nil, nil
}
func (m *MockMonthlyTargetRepository) FindByBrickAndPeriod(brickID string, year, month int) (*monthly_target.MonthlyTarget, error) {
	return nil, nil
}
func (m *MockMonthlyTargetRepository) List(req *monthly_target.ListMonthlyTargetsRequest) ([]monthly_target.MonthlyTarget, int64, int64, error) {
	return nil, 0, 0, nil
}
func (m *MockMonthlyTargetRepository) BatchGetProratedTargetsForPeriod(userIDs []string, startDate, endDate string) (map[string]float64, error) {
	return nil, nil
}
func (m *MockMonthlyTargetRepository) GetTotalEffectiveTarget(year int, month int) (int64, error) {
	return 0, nil
}
func (m *MockMonthlyTargetRepository) GetProratedTargetForPeriod(userID string, startDate, endDate string) (float64, error) {
	return 0, nil
}

// MockCategoryRepository
type MockCategoryRepository struct {
	FindByIDFunc   func(id string) (*category.Category, error)
	FindByCodeFunc func(code string) (*category.Category, error)
	ListFunc       func() ([]category.Category, error)
	CreateFunc     func(cat *category.Category) error
	UpdateFunc     func(cat *category.Category) error
	DeleteFunc     func(id string) error
}

func (m *MockCategoryRepository) FindByID(id string) (*category.Category, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}
func (m *MockCategoryRepository) FindByCode(code string) (*category.Category, error) { return nil, nil }
func (m *MockCategoryRepository) List() ([]category.Category, error) {
	if m.ListFunc != nil {
		return m.ListFunc()
	}
	return nil, nil
}
func (m *MockCategoryRepository) Create(cat *category.Category) error { return nil }
func (m *MockCategoryRepository) Update(cat *category.Category) error { return nil }
func (m *MockCategoryRepository) Delete(id string) error              { return nil }

// MockAccountRepository
type MockAccountRepository struct {
	FindByIDFunc         func(id string) (*account.Account, error)
	ListFunc             func(req *account.ListAccountsRequest) ([]account.Account, int64, error)
	ListAllFunc          func(status string) ([]account.Account, error)
	CreateFunc           func(account *account.Account) error
	UpdateFunc           func(account *account.Account) error
	DeleteFunc           func(id string) error
	GetStatsByStatusFunc func() (map[string]int64, error)
	CountByDateRangeFunc func(startDate, endDate interface{}) (int64, error)
}

func (m *MockAccountRepository) FindByID(id string) (*account.Account, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}
func (m *MockAccountRepository) List(req *account.ListAccountsRequest) ([]account.Account, int64, error) {
	if m.ListFunc != nil {
		return m.ListFunc(req)
	}
	return nil, 0, nil
}
func (m *MockAccountRepository) ListAll(status string) ([]account.Account, error) { return nil, nil }
func (m *MockAccountRepository) Create(account *account.Account) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(account)
	}
	return nil
}
func (m *MockAccountRepository) Update(account *account.Account) error       { return nil }
func (m *MockAccountRepository) Delete(id string) error                      { return nil }
func (m *MockAccountRepository) GetStatsByStatus() (map[string]int64, error) { return nil, nil }
func (m *MockAccountRepository) CountByDateRange(startDate, endDate interface{}) (int64, error) {
	return 0, nil
}

// MockMenuRepository
type MockMenuRepository struct {
	FindByIDFunc     func(id string) (*permission.Menu, error)
	FindByURLFunc    func(url string) (*permission.Menu, error)
	ListFunc         func() ([]permission.Menu, error)
	GetRootMenusFunc func() ([]permission.Menu, error)
	CreateFunc       func(menu *permission.Menu) error
	UpdateFunc       func(menu *permission.Menu) error
	DeleteFunc       func(id string) error
}

func (m *MockMenuRepository) FindByID(id string) (*permission.Menu, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}
func (m *MockMenuRepository) FindByURL(url string) (*permission.Menu, error) {
	if m.FindByURLFunc != nil {
		return m.FindByURLFunc(url)
	}
	return nil, nil
}
func (m *MockMenuRepository) List() ([]permission.Menu, error) {
	if m.ListFunc != nil {
		return m.ListFunc()
	}
	return nil, nil
}
func (m *MockMenuRepository) GetRootMenus() ([]permission.Menu, error) {
	if m.GetRootMenusFunc != nil {
		return m.GetRootMenusFunc()
	}
	return nil, nil
}
func (m *MockMenuRepository) Create(menu *permission.Menu) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(menu)
	}
	return nil
}
func (m *MockMenuRepository) Update(menu *permission.Menu) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(menu)
	}
	return nil
}
func (m *MockMenuRepository) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

// MockProductRepository
type MockProductRepository struct {
	FindByIDFunc func(id string) (*product.Product, error)
	ListFunc     func(req *product.ListProductsRequest) ([]product.Product, int64, error)
	CreateFunc   func(product *product.Product) error
	UpdateFunc   func(product *product.Product) error
	DeleteFunc   func(id string) error
}

func (m *MockProductRepository) FindByID(id string) (*product.Product, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}
func (m *MockProductRepository) List(req *product.ListProductsRequest) ([]product.Product, int64, error) {
	if m.ListFunc != nil {
		return m.ListFunc(req)
	}
	return nil, 0, nil
}
func (m *MockProductRepository) Create(product *product.Product) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(product)
	}
	return nil
}
func (m *MockProductRepository) Update(product *product.Product) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(product)
	}
	return nil
}
func (m *MockProductRepository) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

// MockProductCategoryRepository
type MockProductCategoryRepository struct {
	FindByIDFunc func(id string) (*product.ProductCategory, error)
	ListFunc     func(req *product.ListProductCategoriesRequest) ([]product.ProductCategory, error)
	CreateFunc   func(category *product.ProductCategory) error
	UpdateFunc   func(category *product.ProductCategory) error
	DeleteFunc   func(id string) error
}

func (m *MockProductCategoryRepository) FindByID(id string) (*product.ProductCategory, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}
func (m *MockProductCategoryRepository) List(req *product.ListProductCategoriesRequest) ([]product.ProductCategory, error) {
	if m.ListFunc != nil {
		return m.ListFunc(req)
	}
	return nil, nil
}
func (m *MockProductCategoryRepository) Create(category *product.ProductCategory) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(category)
	}
	return nil
}
func (m *MockProductCategoryRepository) Update(category *product.ProductCategory) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(category)
	}
	return nil
}
func (m *MockProductCategoryRepository) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

// MockContactRoleRepository
type MockContactRoleRepository struct {
	FindByIDFunc   func(id string) (*contact_role.ContactRole, error)
	FindByCodeFunc func(code string) (*contact_role.ContactRole, error)
	ListFunc       func() ([]contact_role.ContactRole, error)
	CreateFunc     func(cr *contact_role.ContactRole) error
	UpdateFunc     func(cr *contact_role.ContactRole) error
	DeleteFunc     func(id string) error
}

func (m *MockContactRoleRepository) FindByID(id string) (*contact_role.ContactRole, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}
func (m *MockContactRoleRepository) FindByCode(code string) (*contact_role.ContactRole, error) {
	if m.FindByCodeFunc != nil {
		return m.FindByCodeFunc(code)
	}
	return nil, nil
}
func (m *MockContactRoleRepository) List() ([]contact_role.ContactRole, error) {
	if m.ListFunc != nil {
		return m.ListFunc()
	}
	return nil, nil
}
func (m *MockContactRoleRepository) Create(cr *contact_role.ContactRole) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(cr)
	}
	return nil
}
func (m *MockContactRoleRepository) Update(cr *contact_role.ContactRole) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(cr)
	}
	return nil
}
func (m *MockContactRoleRepository) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

// MockLeadSourceRepository
type MockLeadSourceRepository struct {
	FindByIDFunc   func(id string) (*lead_source.LeadSource, error)
	FindByCodeFunc func(code string) (*lead_source.LeadSource, error)
	ListFunc       func(req *lead_source.ListLeadSourcesRequest) ([]*lead_source.LeadSource, int64, error)
	ListAllFunc    func() ([]*lead_source.LeadSource, error)
	CreateFunc     func(ls *lead_source.LeadSource) error
	UpdateFunc     func(ls *lead_source.LeadSource) error
	DeleteFunc     func(id string) error
}

func (m *MockLeadSourceRepository) FindByID(id string) (*lead_source.LeadSource, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}
func (m *MockLeadSourceRepository) FindByCode(code string) (*lead_source.LeadSource, error) {
	if m.FindByCodeFunc != nil {
		return m.FindByCodeFunc(code)
	}
	return nil, nil
}
func (m *MockLeadSourceRepository) List(req *lead_source.ListLeadSourcesRequest) ([]*lead_source.LeadSource, int64, error) {
	if m.ListFunc != nil {
		return m.ListFunc(req)
	}
	return nil, 0, nil
}
func (m *MockLeadSourceRepository) ListAll() ([]*lead_source.LeadSource, error) {
	if m.ListAllFunc != nil {
		return m.ListAllFunc()
	}
	return nil, nil
}
func (m *MockLeadSourceRepository) Create(ls *lead_source.LeadSource) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ls)
	}
	return nil
}
func (m *MockLeadSourceRepository) Update(ls *lead_source.LeadSource) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ls)
	}
	return nil
}
func (m *MockLeadSourceRepository) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

// MockLeadRepository
type MockLeadRepository struct {
	CreateFunc                       func(l *lead.Lead) error
	UpdateFunc                       func(l *lead.Lead) error
	DeleteFunc                       func(id string) error
	FindByIDFunc                     func(id string) (*lead.Lead, error)
	ListFunc                         func(req *lead.ListLeadsRequest) ([]lead.Lead, int64, error)
	FindByEmailFunc                  func(email string) (*lead.Lead, error)
	CountByDateRangeFunc             func(startDate, endDate interface{}) (int64, error)
	GetStatsByStatusAndDateRangeFunc func(startDate, endDate interface{}) (map[string]int64, error)
	GetStatsBySourceAndDateRangeFunc func(startDate, endDate interface{}) (map[string]int64, error)
}

func (m *MockLeadRepository) Create(l *lead.Lead) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(l)
	}
	return nil
}
func (m *MockLeadRepository) Update(l *lead.Lead) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(l)
	}
	return nil
}
func (m *MockLeadRepository) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}
func (m *MockLeadRepository) FindByID(id string) (*lead.Lead, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}
func (m *MockLeadRepository) List(req *lead.ListLeadsRequest) ([]lead.Lead, int64, error) {
	if m.ListFunc != nil {
		return m.ListFunc(req)
	}
	return nil, 0, nil
}
func (m *MockLeadRepository) FindByEmail(email string) (*lead.Lead, error) {
	if m.FindByEmailFunc != nil {
		return m.FindByEmailFunc(email)
	}
	return nil, nil
}
func (m *MockLeadRepository) GetAnalytics(req *lead.LeadAnalyticsRequest) (*lead.LeadAnalyticsResponse, error) {
	return nil, nil
}
func (m *MockLeadRepository) GetStatsByStatus() (map[string]int64, error) { return nil, nil }
func (m *MockLeadRepository) GetStatsBySource() (map[string]int64, error) { return nil, nil }
func (m *MockLeadRepository) CountByDateRange(startDate, endDate interface{}) (int64, error) {
	if m.CountByDateRangeFunc != nil {
		return m.CountByDateRangeFunc(startDate, endDate)
	}
	return 0, nil
}
func (m *MockLeadRepository) GetStatsByStatusAndDateRange(startDate, endDate interface{}) (map[string]int64, error) {
	return nil, nil
}
func (m *MockLeadRepository) GetStatsBySourceAndDateRange(startDate, endDate interface{}) (map[string]int64, error) {
	if m.GetStatsBySourceAndDateRangeFunc != nil {
		return m.GetStatsBySourceAndDateRangeFunc(startDate, endDate)
	}
	return nil, nil
}

// MockDealRepository
type MockDealRepository struct {
	CreateFunc             func(d *pipeline.Deal) error
	FindByIDFunc           func(id string) (*pipeline.Deal, error)
	ListFunc               func(req *pipeline.ListDealsRequest) ([]pipeline.Deal, int64, error)
	GetSummaryInPeriodFunc func(startDate, endDate interface{}) (*pipeline.PipelineSummaryResponse, error)
}

func (m *MockDealRepository) Create(d *pipeline.Deal) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(d)
	}
	return nil
}
func (m *MockDealRepository) FindByID(id string) (*pipeline.Deal, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}
func (m *MockDealRepository) List(req *pipeline.ListDealsRequest) ([]pipeline.Deal, int64, error) {
	if m.ListFunc != nil {
		return m.ListFunc(req)
	}
	return nil, 0, nil
}
func (m *MockDealRepository) Update(d *pipeline.Deal) error                          { return nil }
func (m *MockDealRepository) Delete(id string) error                                 { return nil }
func (m *MockDealRepository) GetSummary() (*pipeline.PipelineSummaryResponse, error) { return nil, nil }
func (m *MockDealRepository) GetForecast(periodType string, start, end time.Time) (*pipeline.ForecastResponse, error) {
	return nil, nil
}
func (m *MockDealRepository) GetStatsByStatus(startDate, endDate string, assignedTo, stageID, status string) (map[string]int64, error) {
	return nil, nil
}
func (m *MockDealRepository) GetStatsByStage(startDate, endDate string, assignedTo, status string) (map[string]int64, error) {
	return nil, nil
}
func (m *MockDealRepository) CountByDateRange(startDate, endDate interface{}) (int64, error) {
	return 0, nil
}
func (m *MockDealRepository) GetWonDealsValueInPeriod(startDate, endDate interface{}) (int64, int64, error) {
	return 0, 0, nil
}
func (m *MockDealRepository) GetWonDealsValueInPeriodByUser(userID string, startDate, endDate interface{}) (int64, int64, error) {
	return 0, 0, nil
}
func (m *MockDealRepository) GetSummaryInPeriod(startDate, endDate interface{}) (*pipeline.PipelineSummaryResponse, error) {
	if m.GetSummaryInPeriodFunc != nil {
		return m.GetSummaryInPeriodFunc(startDate, endDate)
	}
	return nil, nil
}

// MockPipelineRepository
type MockPipelineRepository struct {
	FindStageByIDFunc func(id string) (*pipeline.PipelineStage, error)
	GetFirstStageFunc func() (*pipeline.PipelineStage, error)
}

func (m *MockPipelineRepository) FindStageByID(id string) (*pipeline.PipelineStage, error) {
	if m.FindStageByIDFunc != nil {
		return m.FindStageByIDFunc(id)
	}
	return nil, nil
}
func (m *MockPipelineRepository) GetFirstStage() (*pipeline.PipelineStage, error) {
	if m.GetFirstStageFunc != nil {
		return m.GetFirstStageFunc()
	}
	return nil, nil
}
func (m *MockPipelineRepository) CreateStage(s *pipeline.PipelineStage) error { return nil }
func (m *MockPipelineRepository) UpdateStage(s *pipeline.PipelineStage) error { return nil }
func (m *MockPipelineRepository) DeleteStage(id string) error                 { return nil }
func (m *MockPipelineRepository) ListStages(req *pipeline.ListPipelineStagesRequest) ([]pipeline.PipelineStage, error) {
	return nil, nil
}
func (m *MockPipelineRepository) ReorderStages(orders map[string]int) error { return nil }
func (m *MockPipelineRepository) FindStageByCode(code string) (*pipeline.PipelineStage, error) {
	return nil, nil
}

// MockContactRepository
type MockContactRepository struct {
	CreateFunc func(c *contact.Contact) error
}

func (m *MockContactRepository) Create(c *contact.Contact) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(c)
	}
	return nil
}
func (m *MockContactRepository) Update(c *contact.Contact) error              { return nil }
func (m *MockContactRepository) Delete(id string) error                       { return nil }
func (m *MockContactRepository) FindByID(id string) (*contact.Contact, error) { return nil, nil }
func (m *MockContactRepository) List(req *contact.ListContactsRequest) ([]contact.Contact, int64, error) {
	return nil, 0, nil
}
func (m *MockContactRepository) ListByAccountID(accountID string) ([]contact.Contact, error) {
	return nil, nil
}
func (m *MockContactRepository) FindByAccountID(accountID string) ([]contact.Contact, error) {
	return nil, nil
}

// MockLeadStatusRepository
type MockLeadStatusRepository struct {
	FindByIDFunc    func(id string) (*lead_status.LeadStatus, error)
	FindByCodeFunc  func(code string) (*lead_status.LeadStatus, error)
	FindDefaultFunc func() (*lead_status.LeadStatus, error)
	ListAllFunc     func() ([]*lead_status.LeadStatus, error)
}

func (m *MockLeadStatusRepository) FindByID(id string) (*lead_status.LeadStatus, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}
func (m *MockLeadStatusRepository) FindByCode(code string) (*lead_status.LeadStatus, error) {
	if m.FindByCodeFunc != nil {
		return m.FindByCodeFunc(code)
	}
	return nil, nil
}
func (m *MockLeadStatusRepository) FindDefault() (*lead_status.LeadStatus, error) {
	if m.FindDefaultFunc != nil {
		return m.FindDefaultFunc()
	}
	return nil, nil
}
func (m *MockLeadStatusRepository) ListAll() ([]*lead_status.LeadStatus, error) {
	if m.ListAllFunc != nil {
		return m.ListAllFunc()
	}
	return nil, nil
}

// MockDealHistoryRepository
type MockDealHistoryRepository struct{}

func (m *MockDealHistoryRepository) Create(h *deal_history.DealHistory) error { return nil }
func (m *MockDealHistoryRepository) ListByDealID(dealID string) ([]deal_history.DealHistory, error) {
	return nil, nil
}
func (m *MockDealHistoryRepository) Delete(id string) error { return nil }
func (m *MockDealHistoryRepository) FindByDealID(dealID string) ([]deal_history.DealHistory, error) {
	return nil, nil
}
func (m *MockDealHistoryRepository) FindByID(id string) (*deal_history.DealHistory, error) {
	return nil, nil
}
func (m *MockDealHistoryRepository) List(dealID string, limit int) ([]deal_history.DealHistory, error) {
	return nil, nil
}

// MockActivityRepository
type MockActivityRepository struct{}

func (m *MockActivityRepository) Create(a *activity.Activity) error              { return nil }
func (m *MockActivityRepository) Update(a *activity.Activity) error              { return nil }
func (m *MockActivityRepository) Delete(id string) error                         { return nil }
func (m *MockActivityRepository) FindByID(id string) (*activity.Activity, error) { return nil, nil }
func (m *MockActivityRepository) List(req *activity.ListActivitiesRequest) ([]activity.Activity, int64, error) {
	return nil, 0, nil
}
func (m *MockActivityRepository) GetStats(start, end, userID string) (map[string]int64, error) {
	return nil, nil
}
func (m *MockActivityRepository) CountActivities(userID string, start, end time.Time) (int64, error) {
	return 0, nil
}
func (m *MockActivityRepository) FindByAccountID(accountID string) ([]activity.Activity, error) {
	return nil, nil
}
func (m *MockActivityRepository) GetStatsByType(startDate, endDate string, accountID string) (map[string]int64, error) {
	return nil, nil
}
func (m *MockActivityRepository) GetStatsByTypeAndDate(startDate, endDate string, accountID string) (map[string]map[string]int64, error) {
	return nil, nil
}
func (m *MockActivityRepository) GetStatsByUser(startDate, endDate string, accountID string) (map[string]int64, error) {
	return nil, nil
}
func (m *MockActivityRepository) GetTimeline(req *activity.ActivityTimelineRequest) ([]activity.Activity, error) {
	return nil, nil
}

// MockVisitReportRepository
type MockVisitReportRepository struct{}

func (m *MockVisitReportRepository) Create(v *visit_report.VisitReport) error { return nil }
func (m *MockVisitReportRepository) Update(v *visit_report.VisitReport) error { return nil }
func (m *MockVisitReportRepository) Delete(id string) error                   { return nil }
func (m *MockVisitReportRepository) FindByID(id string) (*visit_report.VisitReport, error) {
	return nil, nil
}
func (m *MockVisitReportRepository) List(req *visit_report.ListVisitReportsRequest) ([]visit_report.VisitReport, int64, error) {
	return nil, 0, nil
}
func (m *MockVisitReportRepository) GetStatsByStatus(startDate, endDate string, accountID, salesRepID, status string) (map[string]int64, error) {
	return nil, nil
}
func (m *MockVisitReportRepository) GetStatsByBrick(start, end, userID string) (map[string]int64, error) {
	return nil, nil
}
func (m *MockVisitReportRepository) GetStatsBySalesRep(startDate, endDate string, accountID, status string) (map[string]int64, error) {
	return nil, nil
}
func (m *MockVisitReportRepository) CountVisits(userID string, start, end time.Time) (int64, error) {
	return 0, nil
}
func (m *MockVisitReportRepository) FindByAccountID(accountID string) ([]visit_report.VisitReport, error) {
	return nil, nil
}
func (m *MockVisitReportRepository) FindBySalesRepID(salesRepID string) ([]visit_report.VisitReport, error) {
	return nil, nil
}
func (m *MockVisitReportRepository) GetStatsByDate(startDate, endDate string, accountID, salesRepID, status string) (map[string]int64, error) {
	return nil, nil
}
func (m *MockVisitReportRepository) GetStatsByDateAndStatus(startDate, endDate string, accountID, salesRepID string) (map[string]map[string]int64, error) {
	return nil, nil
}
func (m *MockVisitReportRepository) GetStatsByAccount(startDate, endDate string, salesRepID, status string) (map[string]int64, error) {
	return nil, nil
}
func (m *MockVisitReportRepository) GetStatsBySalesRepWithAccounts(startDate, endDate string, status string) (map[string]struct {
	VisitCount   int64
	AccountCount int64
}, error) {
	return nil, nil
}

// MockTaskRepository
type MockTaskRepository struct{}

func (m *MockTaskRepository) Create(t *task.Task) error              { return nil }
func (m *MockTaskRepository) Update(t *task.Task) error              { return nil }
func (m *MockTaskRepository) Delete(id string) error                 { return nil }
func (m *MockTaskRepository) FindByID(id string) (*task.Task, error) { return nil, nil }
func (m *MockTaskRepository) List(req *task.ListTasksRequest) ([]task.Task, int64, error) {
	return nil, 0, nil
}
func (m *MockTaskRepository) CountTasks(userID string, start, end time.Time) (int64, error) {
	return 0, nil
}
