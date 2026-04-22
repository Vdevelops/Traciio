package notification

import (
	"testing"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/notification"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Define local mock to avoid dependency cycles if reusing from other package
type IntMockNotifRepo struct {
	mock.Mock
	// interfaces.NotificationRepository // Embedded interface
	// Redefine methods needed
}

func (m *IntMockNotifRepo) Create(n *notification.Notification) error {
	args := m.Called(n)
	return args.Error(0)
}

func (m *IntMockNotifRepo) FindByID(id string) (*notification.Notification, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*notification.Notification), args.Error(1)
}

func (m *IntMockNotifRepo) List(req *notification.ListNotificationsRequest) ([]notification.Notification, int64, error) {
	args := m.Called(req)
	return args.Get(0).([]notification.Notification), args.Get(1).(int64), args.Error(2)
}

func (m *IntMockNotifRepo) MarkAsRead(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *IntMockNotifRepo) GetUnreadCount(userID string) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *IntMockNotifRepo) MarkAllAsRead(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *IntMockNotifRepo) Update(notif *notification.Notification) error {
	args := m.Called(notif)
	return args.Error(0)
}

func (m *IntMockNotifRepo) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestNotification_Integration_CreateAndRead(t *testing.T) {
	// Setup Mock Repo
	mockRepo := new(IntMockNotifRepo)
	
	// Create Service
	service := NewService(mockRepo, nil) // eventHelper is nil

	// Test Case 1: Create Notification
	req := &notification.CreateNotificationRequest{
		UserID:  "user-1",
		Title:   "New Task",
		Message: "You have a new task assigned",
		Type:    "task_assigned",
		Data:    `{"ref_id": "task-123", "ref_type": "task"}`,
	}

	mockRepo.On("Create", mock.AnythingOfType("*notification.Notification")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*notification.Notification)
		arg.ID = "generated-id"
		arg.CreatedAt = time.Now()
	})

	// Add expectation for FindByID which is called inside CreateNotification
	mockRepo.On("FindByID", "generated-id").Return(&notification.Notification{
		ID: "generated-id",
		UserID: "user-1",
		Title: "New Task",
		Message: "You have a new task assigned",
		Type: "task_assigned",
		Data: `{"ref_id": "task-123", "ref_type": "task"}`,
		IsRead: false,
		CreatedAt: time.Now(),
	}, nil).Once() // Use Once() because we set it up again for MarkAsRead later or we can rely on .Return sequence if we are careful

	created, err := service.CreateNotification(req)
	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, "New Task", created.Title)
	assert.Equal(t, "generated-id", created.ID)
	
	// Test Case 2: Mark As Read
	mockRepo.On("FindByID", "generated-id").Return(&notification.Notification{
		ID: "generated-id",
		UserID: "user-1",
		IsRead: false,
	}, nil)
	
	mockRepo.On("MarkAsRead", "generated-id").Return(nil)

	err = service.MarkAsRead("generated-id")
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}
