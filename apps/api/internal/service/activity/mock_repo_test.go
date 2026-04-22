package activity

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/account"
	"github.com/gilabs/crm-healthcare/api/internal/domain/activity"
	"github.com/gilabs/crm-healthcare/api/internal/domain/activity_type"
	"github.com/gilabs/crm-healthcare/api/internal/domain/contact"
	"github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
)

// MockActivityRepository
type MockActivityRepository struct {
	FindByIDFunc     func(id string) (*activity.Activity, error)
	ListFunc         func(req *activity.ListActivitiesRequest) ([]activity.Activity, int64, error)
	CreateFunc       func(activity *activity.Activity) error
	UpdateFunc       func(activity *activity.Activity) error
	DeleteFunc       func(id string) error
	FindByDealIDFunc func(dealID string) ([]activity.Activity, error)
}

func (m *MockActivityRepository) FindByID(id string) (*activity.Activity, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}
func (m *MockActivityRepository) List(req *activity.ListActivitiesRequest) ([]activity.Activity, int64, error) {
	if m.ListFunc != nil {
		return m.ListFunc(req)
	}
	return nil, 0, nil
}
func (m *MockActivityRepository) Create(activity *activity.Activity) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(activity)
	}
	return nil
}
func (m *MockActivityRepository) Update(activity *activity.Activity) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(activity)
	}
	return nil
}
func (m *MockActivityRepository) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}
func (m *MockActivityRepository) FindByDealID(dealID string) ([]activity.Activity, error) {
	if m.FindByDealIDFunc != nil {
		return m.FindByDealIDFunc(dealID)
	}
	return nil, nil
}
func (m *MockActivityRepository) GetTimeline(req *activity.ActivityTimelineRequest) ([]activity.Activity, error) {
	return nil, nil
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
func (m *MockActivityRepository) GetStatsByStatus(startDate, endDate string, assignedTo string) (map[string]int64, error) {
	return nil, nil
}
func (m *MockActivityRepository) CountByDateRange(startDate, endDate interface{}) (int64, error) {
	return 0, nil
}

// MockActivityTypeRepository
type MockActivityTypeRepository struct {
	FindByIDFunc func(id string) (*activity_type.ActivityType, error)
}

func (m *MockActivityTypeRepository) FindByID(id string) (*activity_type.ActivityType, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}
func (m *MockActivityTypeRepository) List(req *activity_type.ListActivityTypesRequest) ([]activity_type.ActivityType, error)      { return nil, nil }
func (m *MockActivityTypeRepository) FindByCode(code string) (*activity_type.ActivityType, error) { return nil, nil }
func (m *MockActivityTypeRepository) Create(at *activity_type.ActivityType) error     { return nil }
func (m *MockActivityTypeRepository) Update(at *activity_type.ActivityType) error     { return nil }
func (m *MockActivityTypeRepository) Delete(id string) error                     { return nil }


// MockDealRepository
type MockDealRepository struct {
	FindByIDFunc func(id string) (*pipeline.Deal, error)
}
func (m *MockDealRepository) FindByID(id string) (*pipeline.Deal, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}
func (m *MockDealRepository) List(req *pipeline.ListDealsRequest) ([]pipeline.Deal, int64, error) { return nil, 0, nil }
func (m *MockDealRepository) Create(deal *pipeline.Deal) error { return nil }
func (m *MockDealRepository) Update(deal *pipeline.Deal) error { return nil }
func (m *MockDealRepository) Delete(id string) error { return nil }
func (m *MockDealRepository) GetSummary() (*pipeline.PipelineSummaryResponse, error) { return nil, nil }
func (m *MockDealRepository) GetForecast(periodType string, start, end time.Time) (*pipeline.ForecastResponse, error) { return nil, nil }
func (m *MockDealRepository) GetStatsByStatus(startDate, endDate string, assignedTo, stageID, status string) (map[string]int64, error) { return nil, nil }
func (m *MockDealRepository) GetStatsByStage(startDate, endDate string, assignedTo, status string) (map[string]int64, error) { return nil, nil }
func (m *MockDealRepository) CountByDateRange(startDate, endDate interface{}) (int64, error) { return 0, nil }
func (m *MockDealRepository) GetWonDealsValueInPeriod(startDate, endDate interface{}) (int64, int64, error) { return 0, 0, nil }
func (m *MockDealRepository) GetWonDealsValueInPeriodByUser(userID string, startDate, endDate interface{}) (int64, int64, error) { return 0, 0, nil }


// MockContactRepository
type MockContactRepository struct {
	FindByIDFunc func(id string) (*contact.Contact, error)
}
func (m *MockContactRepository) FindByID(id string) (*contact.Contact, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}
func (m *MockContactRepository) List(req *contact.ListContactsRequest) ([]contact.Contact, int64, error) { return nil, 0, nil }
func (m *MockContactRepository) Create(contact *contact.Contact) error { return nil }
func (m *MockContactRepository) Update(contact *contact.Contact) error { return nil }
func (m *MockContactRepository) Delete(id string) error { return nil }
func (m *MockContactRepository) FindByEmail(email string) (*contact.Contact, error) { return nil, nil }
func (m *MockContactRepository) FindByAccountID(accountID string) ([]contact.Contact, error) { return nil, nil }
func (m *MockContactRepository) GetStatsBySource() (map[string]int64, error) { return nil, nil }
func (m *MockContactRepository) CountByDateRange(startDate, endDate interface{}) (int64, error) { return 0, nil }


// MockUserRepository
type MockUserRepository struct {
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
func (m *MockUserRepository) Delete(id string) error { return nil }
func (m *MockUserRepository) List(req *user.ListUsersRequest) ([]user.User, int64, error) { return nil, 0, nil }
func (m *MockUserRepository) CountUsersByRoleID(roleID string) (int64, error) { return 0, nil }
func (m *MockUserRepository) GetUsersByGroupID(groupID string) ([]user.User, error) { return nil, nil }
func (m *MockUserRepository) GetUsersByBrickID(brickID string) ([]user.User, error) { return nil, nil }
func (m *MockUserRepository) GetUsersByRoleID(roleID string) ([]string, error) { return nil, nil }


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
func (m *MockAccountRepository) List(req *account.ListAccountsRequest) ([]account.Account, int64, error) { return nil, 0, nil }
func (m *MockAccountRepository) ListAll(status string) ([]account.Account, error) { return nil, nil }
func (m *MockAccountRepository) Create(account *account.Account) error { return nil }
func (m *MockAccountRepository) Update(account *account.Account) error { return nil }
func (m *MockAccountRepository) Delete(id string) error { return nil }
func (m *MockAccountRepository) GetStatsByStatus() (map[string]int64, error) { return nil, nil }
func (m *MockAccountRepository) CountByDateRange(startDate, endDate interface{}) (int64, error) { return 0, nil }
