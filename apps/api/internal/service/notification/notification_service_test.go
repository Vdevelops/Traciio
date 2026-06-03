package notification

import (
	"testing"

	"github.com/gilabs/crm-healthcare/api/internal/domain/notification"
)

func setupTest(t *testing.T) (*Service, *MockNotificationRepository) {
	notifRepo := &MockNotificationRepository{}
	service := NewService(notifRepo, nil) // nil eventHelper for testing
	return service, notifRepo
}

func TestService_CreateNotification_Success(t *testing.T) {
	service, notifRepo := setupTest(t)

	notifRepo.CreateFunc = func(notif *notification.Notification) error {
		notif.ID = "notif-1"
		return nil
	}
	notifRepo.FindByIDFunc = func(id string) (*notification.Notification, error) {
		return &notification.Notification{
			ID:      id,
			UserID:  "user-1",
			Title:   "Test Notification",
			Message: "Test message",
			Type:    "reminder",
			IsRead:  false,
		}, nil
	}

	req := &notification.CreateNotificationRequest{
		UserID:  "user-1",
		Title:   "Test Notification",
		Message: "Test message",
		Type:    "reminder",
	}

	resp, err := service.CreateNotification(req)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if resp == nil {
		t.Error("expected response, got nil")
	}
	if resp.Title != "Test Notification" {
		t.Errorf("expected title 'Test Notification', got %s", resp.Title)
	}
}

func TestService_MarkAsRead_Success(t *testing.T) {
	service, notifRepo := setupTest(t)

	notifRepo.FindByIDFunc = func(id string) (*notification.Notification, error) {
		return &notification.Notification{ID: id, UserID: "user-1"}, nil
	}
	notifRepo.MarkAsReadFunc = func(id string) error {
		return nil
	}

	err := service.MarkAsRead("notif-1", "user-1")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestService_MarkAsRead_Forbidden(t *testing.T) {
	service, notifRepo := setupTest(t)

	notifRepo.FindByIDFunc = func(id string) (*notification.Notification, error) {
		return &notification.Notification{ID: id, UserID: "user-2"}, nil
	}

	err := service.MarkAsRead("notif-1", "user-1")
	if err != ErrNotificationForbidden {
		t.Fatalf("expected ErrNotificationForbidden, got: %v", err)
	}
}

func TestService_GetUnreadCount_Success(t *testing.T) {
	service, notifRepo := setupTest(t)

	notifRepo.GetUnreadCountFunc = func(userID string) (int64, error) {
		return 5, nil
	}

	count, err := service.GetUnreadCount("user-1")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if count != 5 {
		t.Errorf("expected count 5, got %d", count)
	}
}
