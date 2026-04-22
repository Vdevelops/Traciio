package handlers

import (
	permissionservice "github.com/gilabs/crm-healthcare/api/internal/service/permission"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	permissionService *permissionservice.Service
}

func NewPermissionHandler(permissionService *permissionservice.Service) *PermissionHandler {
	return &PermissionHandler{
		permissionService: permissionService,
	}
}

// List handles list permissions request
func (h *PermissionHandler) List(c *gin.Context) {
	permissions, err := h.permissionService.List()
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, permissions, nil)
}

// GetByID handles get permission by ID request
func (h *PermissionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	permission, err := h.permissionService.GetByID(id)
	if err != nil {
		if err == permissionservice.ErrPermissionNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource": "permission",
				"permission_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, permission, nil)
}

// GetUserPermissions handles get user permissions request
func (h *PermissionHandler) GetUserPermissions(c *gin.Context) {
	userID := c.Param("id")

	permissions, err := h.permissionService.GetUserPermissions(userID)
	if err != nil {
		if err == permissionservice.ErrUserNotFound {
			errors.ErrorResponse(c, "USER_NOT_FOUND", map[string]interface{}{
				"user_id": userID,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, permissions, nil)
}

// GetMobilePermissions handles get mobile permissions request
func (h *PermissionHandler) GetMobilePermissions(c *gin.Context) {
	// Get user ID from JWT token (set by AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		errors.ErrorResponse(c, "UNAUTHORIZED", map[string]interface{}{
			"reason": "User ID not found in token",
		}, nil)
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		errors.ErrorResponse(c, "UNAUTHORIZED", map[string]interface{}{
			"reason": "Invalid user ID format",
		}, nil)
		return
	}

	permissions, err := h.permissionService.GetMobilePermissions(userIDStr)
	if err != nil {
		if err == permissionservice.ErrUserNotFound {
			errors.ErrorResponse(c, "USER_NOT_FOUND", map[string]interface{}{
				"user_id": userIDStr,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, permissions, nil)
}

