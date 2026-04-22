package interfaces

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/notification"
)

// NotificationRepository defines the interface for notification repository
type NotificationRepository interface {
	// FindByID finds a notification by ID
	FindByID(id string) (*notification.Notification, error)

	// List returns a list of notifications with pagination
	List(req *notification.ListNotificationsRequest) ([]notification.Notification, int64, error)

	// Create creates a new notification
	Create(notif *notification.Notification) error

	// Update updates a notification
	Update(notif *notification.Notification) error

	// Delete soft deletes a notification
	Delete(id string) error

	// MarkAsRead marks a notification as read
	MarkAsRead(id string) error

	// MarkAllAsRead marks all notifications as read for a user
	MarkAllAsRead(userID string) error

	// GetUnreadCount returns the count of unread notifications for a user
	GetUnreadCount(userID string) (int64, error)
}

