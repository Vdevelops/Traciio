package contact

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/account"
	"github.com/gilabs/crm-healthcare/api/internal/domain/contact"
	"github.com/gilabs/crm-healthcare/api/internal/domain/contact_role"
)

// MockContactRepository
type MockContactRepository struct {
	CreateFunc     func(c *contact.Contact) error
	UpdateFunc     func(c *contact.Contact) error
	DeleteFunc     func(id string) error
	FindByIDFunc   func(id string) (*contact.Contact, error)
	ListFunc       func(req *contact.ListContactsRequest) ([]contact.Contact, int64, error)
	ListByAccountIDFunc func(accountID string) ([]contact.Contact, error)
	FindByAccountIDFunc func(accountID string) ([]contact.Contact, error)
}
func (m *MockContactRepository) Create(c *contact.Contact) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(c)
	}
	return nil
}
func (m *MockContactRepository) Update(c *contact.Contact) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(c)
	}
	return nil
}
func (m *MockContactRepository) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}
func (m *MockContactRepository) FindByID(id string) (*contact.Contact, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}
func (m *MockContactRepository) List(req *contact.ListContactsRequest) ([]contact.Contact, int64, error) {
	if m.ListFunc != nil {
		return m.ListFunc(req)
	}
	return nil, 0, nil
}
func (m *MockContactRepository) ListByAccountID(accountID string) ([]contact.Contact, error) {
	if m.ListByAccountIDFunc != nil {
		return m.ListByAccountIDFunc(accountID)
	}
	return nil, nil
}
func (m *MockContactRepository) FindByAccountID(accountID string) ([]contact.Contact, error) {
	if m.FindByAccountIDFunc != nil {
		return m.FindByAccountIDFunc(accountID)
	}
	return nil, nil
}
func (m *MockContactRepository) GetStatsByStatusAndDateRange(startDate, endDate interface{}) (map[string]int64, error) { return nil, nil }
func (m *MockContactRepository) CountByDateRange(startDate, endDate interface{}) (int64, error) { return 0, nil }
func (m *MockContactRepository) GetStatsByUser(startDate, endDate interface{}, accountID string) (map[string]int64, error) { return nil, nil }
func (m *MockContactRepository) CountContacts(userID string, start, end interface{}) (int64, error) { return 0, nil }


// MockAccountRepository
type MockAccountRepository struct {
	FindByIDFunc func(id string) (*account.Account, error)
}
func (m *MockAccountRepository) FindByID(id string) (*account.Account, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}
func (m *MockAccountRepository) Create(a *account.Account) error { return nil }
func (m *MockAccountRepository) Update(a *account.Account) error { return nil }
func (m *MockAccountRepository) Delete(id string) error { return nil }
func (m *MockAccountRepository) List(req *account.ListAccountsRequest) ([]account.Account, int64, error) { return nil, 0, nil }
func (m *MockAccountRepository) ListAll(status string) ([]account.Account, error) { return nil, nil }
func (m *MockAccountRepository) FindByName(name string) (*account.Account, error) { return nil, nil }
func (m *MockAccountRepository) GetStatsByStatus() (map[string]int64, error) { return nil, nil }
func (m *MockAccountRepository) CountByDateRange(startDate, endDate interface{}) (int64, error) { return 0, nil }

// MockContactRoleRepository
type MockContactRoleRepository struct {
	FindByIDFunc func(id string) (*contact_role.ContactRole, error)
}
func (m *MockContactRoleRepository) FindByID(id string) (*contact_role.ContactRole, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}
func (m *MockContactRoleRepository) List() ([]contact_role.ContactRole, error) { return nil, nil }
func (m *MockContactRoleRepository) Create(c *contact_role.ContactRole) error { return nil }
func (m *MockContactRoleRepository) Update(c *contact_role.ContactRole) error { return nil }
func (m *MockContactRoleRepository) Delete(id string) error { return nil }
func (m *MockContactRoleRepository) FindByCode(code string) (*contact_role.ContactRole, error) { return nil, nil }
