package notification

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/notification"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new notification repository
func NewRepository(db *gorm.DB) interfaces.NotificationRepository {
	return &repository{db: db}
}

func (r *repository) FindByID(id string) (*notification.Notification, error) {
	var notif notification.Notification
	err := r.db.Where("id = ?", id).First(&notif).Error
	if err != nil {
		return nil, err
	}
	return &notif, nil
}

func (r *repository) List(req *notification.ListNotificationsRequest) ([]notification.Notification, int64, error) {
	var notifications []notification.Notification
	var total int64

	query := r.db.Model(&notification.Notification{})

	// Apply filters
	if req.UserID != "" {
		query = query.Where("user_id = ?", req.UserID)
	}

	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	if req.IsRead != nil {
		query = query.Where("is_read = ?", *req.IsRead)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
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

	offset := (page - 1) * perPage

	// Fetch data
	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(perPage).
		Find(&notifications).Error
	if err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

func (r *repository) Create(notif *notification.Notification) error {
	return r.db.Create(notif).Error
}

func (r *repository) Update(notif *notification.Notification) error {
	return r.db.Model(notif).Updates(notif).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&notification.Notification{}).Error
}

func (r *repository) MarkAsRead(id string) error {
	now := time.Now()
	return r.db.Model(&notification.Notification{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": &now,
		}).Error
}

func (r *repository) MarkAllAsRead(userID string) error {
	now := time.Now()
	return r.db.Model(&notification.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": &now,
		}).Error
}

func (r *repository) GetUnreadCount(userID string) (int64, error) {
	var count int64
	err := r.db.Model(&notification.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

