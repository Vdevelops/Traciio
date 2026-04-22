package lead_source

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead_source"
)

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
	if m.FindByIDFunc != nil { return m.FindByIDFunc(id) }
	return nil, nil
}
func (m *MockLeadSourceRepository) FindByCode(code string) (*lead_source.LeadSource, error) {
	if m.FindByCodeFunc != nil { return m.FindByCodeFunc(code) }
	return nil, nil
}
func (m *MockLeadSourceRepository) List(req *lead_source.ListLeadSourcesRequest) ([]*lead_source.LeadSource, int64, error) {
	if m.ListFunc != nil { return m.ListFunc(req) }
	return nil, 0, nil
}
func (m *MockLeadSourceRepository) ListAll() ([]*lead_source.LeadSource, error) {
	if m.ListAllFunc != nil { return m.ListAllFunc() }
	return nil, nil
}
func (m *MockLeadSourceRepository) Create(ls *lead_source.LeadSource) error {
	if m.CreateFunc != nil { return m.CreateFunc(ls) }
	return nil
}
func (m *MockLeadSourceRepository) Update(ls *lead_source.LeadSource) error {
	if m.UpdateFunc != nil { return m.UpdateFunc(ls) }
	return nil
}
func (m *MockLeadSourceRepository) Delete(id string) error {
	if m.DeleteFunc != nil { return m.DeleteFunc(id) }
	return nil
}
