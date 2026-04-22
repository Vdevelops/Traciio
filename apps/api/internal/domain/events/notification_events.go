package events

import "time"

// NotificationCreatedEvent is emitted when a new notification is created
type NotificationCreatedEvent struct {
	NotificationID string    `json:"notification_id"`
	UserID         string    `json:"user_id"`
	Title          string    `json:"title"`
	Message        string    `json:"message"`
	Type           string    `json:"type"`
	Data           string    `json:"data,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}
