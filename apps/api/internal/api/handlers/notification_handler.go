package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gilabs/crm-healthcare/api/internal/domain/notification"
	notificationservice "github.com/gilabs/crm-healthcare/api/internal/service/notification"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
)

type NotificationHandler struct {
	notificationService *notificationservice.Service
}

func NewNotificationHandler(notificationService *notificationservice.Service) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// List returns a list of notifications
func (h *NotificationHandler) List(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		errors.UnauthorizedResponse(c, "user ID not found in context")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		errors.UnauthorizedResponse(c, "invalid user ID format")
		return
	}

	// Parse query parameters
	req := &notification.ListNotificationsRequest{
		UserID: userIDStr,
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			req.Page = page
		}
	}

	if perPageStr := c.Query("per_page"); perPageStr != "" {
		if perPage, err := strconv.Atoi(perPageStr); err == nil {
			req.PerPage = perPage
		}
	}

	if typeFilter := c.Query("type"); typeFilter != "" {
		req.Type = typeFilter
	}

	if isReadStr := c.Query("is_read"); isReadStr != "" {
		if isRead, err := strconv.ParseBool(isReadStr); err == nil {
			req.IsRead = &isRead
		}
	}

	notifications, pagination, err := h.notificationService.ListNotifications(req)
	if err != nil {
		errors.ErrorResponse(c, "NOTIFICATION_LIST_FAILED", map[string]interface{}{
			"error": err.Error(),
		}, nil)
		return
	}

	meta := &response.Meta{
		Pagination: &response.PaginationMeta{
			Page:       pagination.Page,
			PerPage:    pagination.PerPage,
			Total:      pagination.Total,
			TotalPages: pagination.TotalPages,
			HasNext:    pagination.Page < pagination.TotalPages,
			HasPrev:    pagination.Page > 1,
		},
	}

	response.SuccessResponse(c, notifications, meta)
}

// GetUnreadCount returns the unread notification count
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		errors.UnauthorizedResponse(c, "user ID not found in context")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		errors.UnauthorizedResponse(c, "invalid user ID format")
		return
	}

	count, err := h.notificationService.GetUnreadCount(userIDStr)
	if err != nil {
		errors.ErrorResponse(c, "NOTIFICATION_COUNT_FAILED", map[string]interface{}{
			"error": err.Error(),
		}, nil)
		return
	}

	response.SuccessResponse(c, gin.H{
		"unread_count": count,
	}, nil)
}

// MarkAsRead marks a notification as read
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		errors.InvalidQueryParamResponse(c)
		return
	}

	err := h.notificationService.MarkAsRead(id)
	if err != nil {
		if err.Error() == "notification not found" {
			errors.NotFoundResponse(c, "NOTIFICATION_NOT_FOUND", "notification not found")
			return
		}
		errors.ErrorResponse(c, "NOTIFICATION_MARK_READ_FAILED", map[string]interface{}{
			"error": err.Error(),
		}, nil)
		return
	}

	response.SuccessResponse(c, gin.H{
		"message": "Notification marked as read",
	}, nil)
}

// MarkAllAsRead marks all notifications as read for the current user
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		errors.UnauthorizedResponse(c, "user ID not found in context")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		errors.UnauthorizedResponse(c, "invalid user ID format")
		return
	}

	err := h.notificationService.MarkAllAsRead(userIDStr)
	if err != nil {
		errors.ErrorResponse(c, "NOTIFICATION_MARK_ALL_READ_FAILED", map[string]interface{}{
			"error": err.Error(),
		}, nil)
		return
	}

	response.SuccessResponse(c, gin.H{
		"message": "All notifications marked as read",
	}, nil)
}

// Delete deletes a notification
func (h *NotificationHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		errors.InvalidQueryParamResponse(c)
		return
	}

	err := h.notificationService.DeleteNotification(id)
	if err != nil {
		if err.Error() == "notification not found" {
			errors.NotFoundResponse(c, "NOTIFICATION_NOT_FOUND", "notification not found")
			return
		}
		errors.ErrorResponse(c, "NOTIFICATION_DELETE_FAILED", map[string]interface{}{
			"error": err.Error(),
		}, nil)
		return
	}

	response.SuccessResponse(c, gin.H{
		"message": "Notification deleted",
	}, nil)
}

