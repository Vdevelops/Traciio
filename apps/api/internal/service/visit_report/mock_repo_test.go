package visit_report

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/account"
	"github.com/gilabs/crm-healthcare/api/internal/domain/activity"
	"github.com/gilabs/crm-healthcare/api/internal/domain/contact"
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead"
	"github.com/gilabs/crm-healthcare/api/internal/domain/notification"
	"github.com/gilabs/crm-healthcare/api/internal/domain/task"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"github.com/gilabs/crm-healthcare/api/internal/domain/visit_report"
)

// MockVisitReportRepository
type MockVisitReportRepository struct {
	FindByIDFunc func(id string) (*visit_report.VisitReport, error)
	ListFunc     func(req *visit_report.ListVisitReportsRequest) ([]visit_report.VisitReport, int64, error)
	CreateFunc   func(vr *visit_report.VisitReport) error
	UpdateFunc   func(vr *visit_report.VisitReport) error
	DeleteFunc   func(id string) error
}

func (m *MockVisitReportRepository) FindByID(id string) (*visit_report.VisitReport, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}
func (m *MockVisitReportRepository) List(req *visit_report.ListVisitReportsRequest) ([]visit_report.VisitReport, int64, error) {
	if m.ListFunc != nil {
		return m.ListFunc(req)
	}
	return nil, 0, nil
}
func (m *MockVisitReportRepository) Create(vr *visit_report.VisitReport) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(vr)
	}
	return nil
}
func (m *MockVisitReportRepository) Update(vr *visit_report.VisitReport) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(vr)
	}
	return nil
}
func (m *MockVisitReportRepository) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}
func (m *MockVisitReportRepository) FindByAccountID(accountID string) ([]visit_report.VisitReport, error) { return nil, nil }
func (m *MockVisitReportRepository) FindBySalesRepID(salesRepID string) ([]visit_report.VisitReport, error) { return nil, nil }
func (m *MockVisitReportRepository) GetStatsByStatus(startDate, endDate string, accountID, salesRepID, status string) (map[string]int64, error) { return nil, nil }
func (m *MockVisitReportRepository) GetStatsByDate(startDate, endDate string, accountID, salesRepID, status string) (map[string]int64, error) { return nil, nil }
func (m *MockVisitReportRepository) GetStatsByDateAndStatus(startDate, endDate string, accountID, salesRepID string) (map[string]map[string]int64, error) { return nil, nil }
func (m *MockVisitReportRepository) GetStatsByAccount(startDate, endDate string, salesRepID, status string) (map[string]int64, error) { return nil, nil }
func (m *MockVisitReportRepository) GetStatsBySalesRep(startDate, endDate string, accountID, status string) (map[string]int64, error) { return nil, nil }
func (m *MockVisitReportRepository) GetStatsBySalesRepWithAccounts(startDate, endDate string, status string) (map[string]struct {
		VisitCount   int64
		AccountCount int64
	}, error) { return nil, nil }


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

// MockActivityRepository
type MockActivityRepository struct {
	CreateFunc func(activity *activity.Activity) error
}
func (m *MockActivityRepository) Create(activity *activity.Activity) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(activity)
	}
	return nil
}
func (m *MockActivityRepository) List(req *activity.ListActivitiesRequest) ([]activity.Activity, int64, error) { return nil, 0, nil }
func (m *MockActivityRepository) GetTimeline(req *activity.ActivityTimelineRequest) ([]activity.Activity, error) { return nil, nil }
func (m *MockActivityRepository) FindByID(id string) (*activity.Activity, error) { return nil, nil }
func (m *MockActivityRepository) Update(activity *activity.Activity) error { return nil }
func (m *MockActivityRepository) Delete(id string) error { return nil }
func (m *MockActivityRepository) FindByAccountID(accountID string) ([]activity.Activity, error) { return nil, nil }
func (m *MockActivityRepository) GetStatsByType(startDate, endDate string, assignedTo string) (map[string]int64, error) { return nil, nil }
func (m *MockActivityRepository) GetStatsByTypeAndDate(startDate, endDate string, assignedTo string) (map[string]map[string]int64, error) { return nil, nil }
func (m *MockActivityRepository) GetStatsByUser(startDate, endDate string, assignedTo string) (map[string]int64, error) { return nil, nil }
func (m *MockActivityRepository) CountByDateRange(startDate, endDate interface{}) (int64, error) { return 0, nil }

// MockLeadRepository
type MockLeadRepository struct {
	FindByIDFunc func(id string) (*lead.Lead, error)
}
func (m *MockLeadRepository) FindByID(id string) (*lead.Lead, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}
func (m *MockLeadRepository) List(req *lead.ListLeadsRequest) ([]lead.Lead, int64, error) { return nil, 0, nil }
func (m *MockLeadRepository) Create(lead *lead.Lead) error { return nil }
func (m *MockLeadRepository) Update(lead *lead.Lead) error { return nil }
func (m *MockLeadRepository) Delete(id string) error { return nil }
func (m *MockLeadRepository) GetStatsByStatus() (map[string]int64, error) { return nil, nil }
func (m *MockLeadRepository) GetStatsByStatusAndDateRange(startDate, endDate interface{}) (map[string]int64, error) { return nil, nil }
func (m *MockLeadRepository) GetStatsBySource() (map[string]int64, error) { return nil, nil }
func (m *MockLeadRepository) GetAnalytics(req *lead.LeadAnalyticsRequest) (*lead.LeadAnalyticsResponse, error) { return nil, nil }
func (m *MockLeadRepository) CountByDateRange(startDate, endDate interface{}) (int64, error) { return 0, nil }
func (m *MockLeadRepository) GetStatsBySourceAndDateRange(startDate, endDate interface{}) (map[string]int64, error) { return nil, nil }
func (m *MockLeadRepository) FindByEmail(email string) (*lead.Lead, error) { return nil, nil } // Added potential missing method
func (m *MockLeadRepository) FindByPhone(phone string) (*lead.Lead, error) { return nil, nil }


// MockTaskRepository
type MockTaskRepository struct {}
func (m *MockTaskRepository) FindByID(id string) (*task.Task, error) { return nil, nil }
func (m *MockTaskRepository) List(req *task.ListTasksRequest) ([]task.Task, int64, error) { return nil, 0, nil }
func (m *MockTaskRepository) Create(task *task.Task) error { return nil }
func (m *MockTaskRepository) Update(task *task.Task) error { return nil }
func (m *MockTaskRepository) Delete(id string) error { return nil }
func (m *MockTaskRepository) GetStatsByStatus(startDate, endDate string, assignedTo string) (map[string]int64, error) { return nil, nil }
func (m *MockTaskRepository) GetStatsByPriority(startDate, endDate string, assignedTo string) (map[string]int64, error) { return nil, nil }
func (m *MockTaskRepository) CountByDateRange(startDate, endDate interface{}) (int64, error) { return 0, nil }

// MockNotificationRepository
type MockNotificationRepository struct {
	CreateFunc func(notification *notification.Notification) error
	FindByIDFunc func(id string) (*notification.Notification, error)
}
func (m *MockNotificationRepository) FindByID(id string) (*notification.Notification, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}
func (m *MockNotificationRepository) Create(notification *notification.Notification) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(notification)
	}
	return nil
}
func (m *MockNotificationRepository) List(req *notification.ListNotificationsRequest) ([]notification.Notification, int64, error) { return nil, 0, nil }
func (m *MockNotificationRepository) MarkAsRead(id string) error { return nil }
func (m *MockNotificationRepository) MarkAllAsRead(userID string) error { return nil }
func (m *MockNotificationRepository) Update(notification *notification.Notification) error { return nil }
func (m *MockNotificationRepository) CountUnread(userID string) (int64, error) { return 0, nil }
func (m *MockNotificationRepository) GetUnreadCount(userID string) (int64, error) { return 0, nil }
func (m *MockNotificationRepository) Delete(id string) error { return nil }
