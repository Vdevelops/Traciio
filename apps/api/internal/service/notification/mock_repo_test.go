package notification

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/notification"
)

// MockNotificationRepository
type MockNotificationRepository struct {
	FindByIDFunc        func(id string) (*notification.Notification, error)
	CreateFunc          func(notif *notification.Notification) error
	ListFunc            func(req *notification.ListNotificationsRequest) ([]notification.Notification, int64, error)
	MarkAsReadFunc      func(id string) error
	MarkAllAsReadFunc   func(userID string) error
	DeleteFunc          func(id string) error
	GetUnreadCountFunc  func(userID string) (int64, error)
}

func (m *MockNotificationRepository) FindByID(id string) (*notification.Notification, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}
func (m *MockNotificationRepository) Create(notif *notification.Notification) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(notif)
	}
	notif.ID = "notif-1" // Simulate ID assignment
	return nil
}
func (m *MockNotificationRepository) List(req *notification.ListNotificationsRequest) ([]notification.Notification, int64, error) {
	if m.ListFunc != nil {
		return m.ListFunc(req)
	}
	return nil, 0, nil
}
func (m *MockNotificationRepository) MarkAsRead(id string) error {
	if m.MarkAsReadFunc != nil {
		return m.MarkAsReadFunc(id)
	}
	return nil
}
func (m *MockNotificationRepository) MarkAllAsRead(userID string) error {
	if m.MarkAllAsReadFunc != nil {
		return m.MarkAllAsReadFunc(userID)
	}
	return nil
}
func (m *MockNotificationRepository) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}
func (m *MockNotificationRepository) GetUnreadCount(userID string) (int64, error) {
	if m.GetUnreadCountFunc != nil {
		return m.GetUnreadCountFunc(userID)
	}
	return 0, nil
}
func (m *MockNotificationRepository) Update(notif *notification.Notification) error { return nil }
func (m *MockNotificationRepository) CountUnread(userID string) (int64, error) { return 0, nil }
