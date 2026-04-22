package notification

import (
	"errors"

	"github.com/gilabs/crm-healthcare/api/internal/domain/notification"
	domainevents "github.com/gilabs/crm-healthcare/api/internal/domain/events"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

var (
	ErrNotificationNotFound = errors.New("notification not found")
)

type Service struct {
	notifRepo   interfaces.NotificationRepository
	hub         HubInterface
	eventHelper *domainevents.Helper
}

// HubInterface defines interface for notification hub
type HubInterface interface {
	BroadcastNotification(userID string, notification interface{})
	BroadcastNotificationUpdate(userID string, notification interface{})
	BroadcastNotificationDelete(userID string, notificationID string)
}

func NewService(notifRepo interfaces.NotificationRepository, eventHelper *domainevents.Helper) *Service {
	return &Service{
		notifRepo:   notifRepo,
		hub:         nil, // Will be set via SetHub if needed
		eventHelper: eventHelper,
	}
}

// SetHub sets the notification hub for broadcasting
func (s *Service) SetHub(hub HubInterface) {
	s.hub = hub
}

// CreateNotification creates a new notification
func (s *Service) CreateNotification(req *notification.CreateNotificationRequest) (*notification.NotificationResponse, error) {
	notifType := req.Type
	if notifType == "" {
		notifType = "reminder"
	}

	notif := &notification.Notification{
		UserID:  req.UserID,
		Title:   req.Title,
		Message: req.Message,
		Type:    notifType,
		Data:    req.Data,
		IsRead:  false,
	}

	if err := s.notifRepo.Create(notif); err != nil {
		return nil, err
	}

	// Reload to get full data
	notif, err := s.notifRepo.FindByID(notif.ID)
	if err != nil {
		return nil, err
	}

	response := notif.ToNotificationResponse()

	// Emit NotificationCreatedEvent
	if s.eventHelper != nil {
		s.eventHelper.EmitNotificationCreated(&domainevents.NotificationCreatedEvent{
			NotificationID: response.ID,
			UserID:         response.UserID,
			Title:          response.Title,
			Message:        response.Message,
			Type:           response.Type,
			Data:           req.Data,
			CreatedAt:      response.CreatedAt,
		})
	}

	// Broadcast via WebSocket if hub is set
	if s.hub != nil {
		notifMap := map[string]interface{}{
			"id":         response.ID,
			"user_id":    response.UserID,
			"title":      response.Title,
			"message":    response.Message,
			"type":       response.Type,
			"is_read":    response.IsRead,
			"created_at": response.CreatedAt,
		}
		s.hub.BroadcastNotification(response.UserID, notifMap)
	}

	return response, nil
}

// ListNotifications returns a list of notifications with pagination
func (s *Service) ListNotifications(req *notification.ListNotificationsRequest) ([]notification.NotificationResponse, *PaginationResult, error) {
	notifications, total, err := s.notifRepo.List(req)
	if err != nil {
		return nil, nil, err
	}

	responses := make([]notification.NotificationResponse, len(notifications))
	for i, notif := range notifications {
		responses[i] = *notif.ToNotificationResponse()
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}

	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	pagination := &PaginationResult{
		Page:       page,
		PerPage:    perPage,
		Total:      int(total),
		TotalPages: totalPages,
	}

	return responses, pagination, nil
}

// GetNotificationByID returns a notification by ID
func (s *Service) GetNotificationByID(id string) (*notification.NotificationResponse, error) {
	notif, err := s.notifRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotificationNotFound
		}
		return nil, err
	}

	return notif.ToNotificationResponse(), nil
}

// MarkAsRead marks a notification as read
func (s *Service) MarkAsRead(id string) error {
	_, err := s.notifRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotificationNotFound
		}
		return err
	}

	return s.notifRepo.MarkAsRead(id)
}

// MarkAllAsRead marks all notifications as read for a user
func (s *Service) MarkAllAsRead(userID string) error {
	return s.notifRepo.MarkAllAsRead(userID)
}

// DeleteNotification deletes a notification
func (s *Service) DeleteNotification(id string) error {
	_, err := s.notifRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotificationNotFound
		}
		return err
	}

	return s.notifRepo.Delete(id)
}

// GetUnreadCount returns the count of unread notifications for a user
func (s *Service) GetUnreadCount(userID string) (int64, error) {
	return s.notifRepo.GetUnreadCount(userID)
}

// PaginationResult represents pagination information
type PaginationResult struct {
	Page       int
	PerPage    int
	Total      int
	TotalPages int
}

