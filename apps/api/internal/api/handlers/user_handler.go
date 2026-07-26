package handlers

import (
	"errors"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	userservice "github.com/gilabs/crm-healthcare/api/internal/service/user"
	pkgerrors "github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	userService     *userservice.Service
	profileService  *userservice.ProfileService
	settingsService *userservice.SettingsService
}

func NewUserHandler(userService *userservice.Service, profileService *userservice.ProfileService, settingsService *userservice.SettingsService) *UserHandler {
	return &UserHandler{
		userService:     userService,
		profileService:  profileService,
		settingsService: settingsService,
	}
}

// List handles list users request
func (h *UserHandler) List(c *gin.Context) {
	var req user.ListUsersRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			pkgerrors.HandleValidationError(c, validationErrors)
			return
		}
		pkgerrors.InvalidQueryParamResponse(c)
		return
	}

	// Apply RBAC scope filtering
	if userCtx := middleware.GetUserContext(c); userCtx != nil {
		req.ScopedUserIDs = userCtx.GetScopedUserIDs("users")
	}

	users, pagination, err := h.userService.List(&req)
	if err != nil {
		pkgerrors.InternalServerErrorResponse(c, "")
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
		Filters: map[string]interface{}{},
	}

	if req.Search != "" {
		meta.Filters["search"] = req.Search
	}
	if req.Status != "" {
		meta.Filters["status"] = req.Status
	}
	if req.RoleID != "" {
		meta.Filters["role_id"] = req.RoleID
	}
	if req.GroupID != "" {
		meta.Filters["group_id"] = req.GroupID
	}

	response.SuccessResponse(c, users, meta)
}

// GetByID handles get user by ID request
func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userService.GetByID(id)
	if err != nil {
		if err == userservice.ErrUserNotFound {
			pkgerrors.ErrorResponse(c, "USER_NOT_FOUND", map[string]interface{}{
				"user_id": id,
			}, nil)
			return
		}
		pkgerrors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, user, nil)
}

// Create handles create user request
func (h *UserHandler) Create(c *gin.Context) {
	var req user.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			pkgerrors.HandleValidationError(c, validationErrors)
			return
		}
		pkgerrors.InvalidRequestBodyResponse(c)
		return
	}

	createdUser, err := h.userService.Create(&req)
	if err != nil {
		if err == userservice.ErrUserAlreadyExists {
			pkgerrors.ErrorResponse(c, "RESOURCE_ALREADY_EXISTS", map[string]interface{}{
				"resource": "user",
				"field":    "email",
				"value":    req.Email,
			}, nil)
			return
		}
		if err == userservice.ErrRoleNotFound {
			pkgerrors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource": "role",
				"role_id":  req.RoleID,
			}, nil)
			return
		}
		if err == userservice.ErrGroupNotFound {
			pkgerrors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource": "group",
			}, nil)
			return
		}
		if err == userservice.ErrBrickNotFound {
			pkgerrors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource": "brick",
			}, nil)
			return
		}
		pkgerrors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			meta.CreatedBy = id
		}
	}

	response.SuccessResponseCreated(c, createdUser, meta)
}

// Update handles update user request
func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req user.UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			pkgerrors.HandleValidationError(c, validationErrors)
			return
		}
		pkgerrors.InvalidRequestBodyResponse(c)
		return
	}

	updatedUser, err := h.userService.Update(id, &req)
	if err != nil {
		if err == userservice.ErrUserNotFound {
			pkgerrors.ErrorResponse(c, "USER_NOT_FOUND", map[string]interface{}{
				"user_id": id,
			}, nil)
			return
		}
		if err == userservice.ErrUserAlreadyExists {
			pkgerrors.ErrorResponse(c, "RESOURCE_ALREADY_EXISTS", map[string]interface{}{
				"resource": "user",
				"field":    "email",
			}, nil)
			return
		}
		if err == userservice.ErrRoleNotFound {
			pkgerrors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource": "role",
			}, nil)
			return
		}
		if err == userservice.ErrBrickNotFound {
			pkgerrors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource": "brick",
			}, nil)
			return
		}
		pkgerrors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			meta.UpdatedBy = id
		}
	}

	response.SuccessResponse(c, updatedUser, meta)
}

// Delete handles delete user request
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.userService.Delete(id)
	if err != nil {
		if err == userservice.ErrUserNotFound {
			pkgerrors.ErrorResponse(c, "USER_NOT_FOUND", map[string]interface{}{
				"user_id": id,
			}, nil)
			return
		}
		pkgerrors.InternalServerErrorResponse(c, "")
		return
	}

	// Get user ID for meta
	meta := &response.Meta{}
	if userIDVal, exists := c.Get("user_id"); exists {
		if id, ok := userIDVal.(string); ok {
			meta.DeletedBy = id
		}
	}

	response.SuccessResponseDeleted(c, "user", id, meta)
}

// GetProfile handles get user profile request
func (h *UserHandler) GetProfile(c *gin.Context) {
	id := c.Param("id")

	profile, err := h.profileService.GetProfile(id)
	if err != nil {
		if err == userservice.ErrUserNotFound {
			pkgerrors.ErrorResponse(c, "USER_NOT_FOUND", map[string]interface{}{
				"user_id": id,
			}, nil)
			return
		}

		pkgerrors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, profile, nil)
}

// UpdateProfile handles update user profile request
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	id := c.Param("id")
	var req user.UpdateProfileRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			pkgerrors.HandleValidationError(c, validationErrors)
			return
		}
		pkgerrors.InvalidRequestBodyResponse(c)
		return
	}

	updatedUser, err := h.profileService.UpdateProfile(id, &req)
	if err != nil {
		if err == userservice.ErrUserNotFound {
			pkgerrors.ErrorResponse(c, "USER_NOT_FOUND", map[string]interface{}{
				"user_id": id,
			}, nil)
			return
		}
		pkgerrors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(string); ok {
			meta.UpdatedBy = uid
		}
	}

	response.SuccessResponse(c, updatedUser, meta)
}

// ChangePassword handles change password request
func (h *UserHandler) ChangePassword(c *gin.Context) {
	id := c.Param("id")
	var req user.ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			pkgerrors.HandleValidationError(c, validationErrors)
			return
		}
		pkgerrors.InvalidRequestBodyResponse(c)
		return
	}

	err := h.profileService.ChangePassword(id, &req)
	if err != nil {
		if err == userservice.ErrUserNotFound {
			pkgerrors.ErrorResponse(c, "USER_NOT_FOUND", map[string]interface{}{
				"user_id": id,
			}, nil)
			return
		}
		if errors.Is(err, userservice.ErrIncorrectPassword) {
			pkgerrors.ErrorResponse(c, "INVALID_CREDENTIALS", map[string]interface{}{
				"field": "current_password",
			}, nil)
			return
		}
		if err.Error() == "passwords do not match" {
			pkgerrors.ErrorResponse(c, "VALIDATION_ERROR", nil, []response.FieldError{
				{
					Field:   "confirm_password",
					Code:    "INVALID_FORMAT",
					Message: "Passwords do not match",
				},
			})
			return
		}
		pkgerrors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponseNoContent(c)
}

// GetMyProfile handles get current user profile request.
// Uses user ID from JWT token, no need for :id in path
func (h *UserHandler) GetMyProfile(c *gin.Context) {
	// Get user ID from JWT token (set by AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		pkgerrors.ErrorResponse(c, "UNAUTHORIZED", map[string]interface{}{
			"reason": "User ID not found in token",
		}, nil)
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		pkgerrors.ErrorResponse(c, "UNAUTHORIZED", map[string]interface{}{
			"reason": "Invalid user ID format",
		}, nil)
		return
	}

	profile, err := h.profileService.GetProfile(userIDStr)
	if err != nil {
		if err == userservice.ErrUserNotFound {
			pkgerrors.ErrorResponse(c, "USER_NOT_FOUND", map[string]interface{}{
				"resource":    "user",
				"resource_id": userIDStr,
			}, nil)
			return
		}
		pkgerrors.InternalServerErrorResponse(c, "")
		return
	}

	// Create meta (empty but present for consistency)
	meta := &response.Meta{}

	response.SuccessResponse(c, profile, meta)
}

// UpdateMyProfile handles update current user profile request.
// Uses user ID from JWT token, no need for :id in path
func (h *UserHandler) UpdateMyProfile(c *gin.Context) {
	// Get user ID from JWT token (set by AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		pkgerrors.ErrorResponse(c, "UNAUTHORIZED", map[string]interface{}{
			"reason": "User ID not found in token",
		}, nil)
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		pkgerrors.ErrorResponse(c, "UNAUTHORIZED", map[string]interface{}{
			"reason": "Invalid user ID format",
		}, nil)
		return
	}

	var req user.UpdateProfileRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			pkgerrors.HandleValidationError(c, validationErrors)
			return
		}
		pkgerrors.InvalidRequestBodyResponse(c)
		return
	}

	updatedUser, err := h.profileService.UpdateProfile(userIDStr, &req)
	if err != nil {
		if err == userservice.ErrUserNotFound {
			pkgerrors.ErrorResponse(c, "USER_NOT_FOUND", map[string]interface{}{
				"resource":    "user",
				"resource_id": userIDStr,
			}, nil)
			return
		}
		pkgerrors.InternalServerErrorResponse(c, "")
		return
	}

	// Create meta with updated_by
	meta := &response.Meta{}
	if uid, exists := c.Get("user_id"); exists {
		if id, ok := uid.(string); ok {
			meta.UpdatedBy = id
		}
	}

	response.SuccessResponse(c, updatedUser, meta)
}

// ChangeMyPassword handles change current user password request.
// Uses user ID from JWT token, no need for :id in path
func (h *UserHandler) ChangeMyPassword(c *gin.Context) {
	// Get user ID from JWT token (set by AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		pkgerrors.ErrorResponse(c, "UNAUTHORIZED", map[string]interface{}{
			"reason": "User ID not found in token",
		}, nil)
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		pkgerrors.ErrorResponse(c, "UNAUTHORIZED", map[string]interface{}{
			"reason": "Invalid user ID format",
		}, nil)
		return
	}

	var req user.ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			pkgerrors.HandleValidationError(c, validationErrors)
			return
		}
		pkgerrors.InvalidRequestBodyResponse(c)
		return
	}

	err := h.profileService.ChangePassword(userIDStr, &req)
	if err != nil {
		if err == userservice.ErrUserNotFound {
			pkgerrors.ErrorResponse(c, "USER_NOT_FOUND", map[string]interface{}{
				"resource":    "user",
				"resource_id": userIDStr,
			}, nil)
			return
		}
		if errors.Is(err, userservice.ErrIncorrectPassword) {
			pkgerrors.ErrorResponse(c, "INVALID_CREDENTIALS", map[string]interface{}{
				"field":  "current_password",
				"reason": "Current password is incorrect",
			}, nil)
			return
		}
		if err.Error() == "passwords do not match" {
			pkgerrors.ErrorResponse(c, "VALIDATION_ERROR", nil, []response.FieldError{
				{
					Field:   "confirm_password",
					Code:    "INVALID_FORMAT",
					Message: "Passwords do not match",
				},
			})
			return
		}
		pkgerrors.InternalServerErrorResponse(c, "")
		return
	}

	// Change password returns 204 No Content (no body) as per standard
	response.SuccessResponseNoContent(c)
}

// GetMySettingsSummary handles get current user settings summary request
// Uses user ID from JWT token, returns extended stats including revenue
// Accepts optional start_date and end_date query parameters for filtering
func (h *UserHandler) GetMySettingsSummary(c *gin.Context) {
	// Get user ID from JWT token (set by AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		pkgerrors.ErrorResponse(c, "UNAUTHORIZED", map[string]interface{}{
			"reason": "User ID not found in token",
		}, nil)
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		pkgerrors.ErrorResponse(c, "UNAUTHORIZED", map[string]interface{}{
			"reason": "Invalid user ID format",
		}, nil)
		return
	}

	// Bind query parameters for date filtering
	var req user.GetSettingsSummaryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			pkgerrors.HandleValidationError(c, validationErrors)
			return
		}
		pkgerrors.InvalidQueryParamResponse(c)
		return
	}

	// Parse date range (similar to sales overview service)
	var startDate, endDate interface{}
	if req.StartDate != "" {
		parsed, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			pkgerrors.ErrorResponse(c, "VALIDATION_ERROR", map[string]interface{}{
				"message": "Invalid start_date format, expected YYYY-MM-DD",
			}, nil)
			return
		}
		startDate = parsed
	}
	if req.EndDate != "" {
		parsed, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			pkgerrors.ErrorResponse(c, "VALIDATION_ERROR", map[string]interface{}{
				"message": "Invalid end_date format, expected YYYY-MM-DD",
			}, nil)
			return
		}
		// Set to end of day for inclusive filtering
		endDate = parsed.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}

	summary, err := h.settingsService.GetSettingsSummary(userIDStr, startDate, endDate)
	if err != nil {
		if err == userservice.ErrUserNotFound {
			pkgerrors.ErrorResponse(c, "USER_NOT_FOUND", map[string]interface{}{
				"resource":    "user",
				"resource_id": userIDStr,
			}, nil)
			return
		}
		pkgerrors.InternalServerErrorResponse(c, "")
		return
	}

	// Create meta with filters
	meta := &response.Meta{
		Filters: map[string]interface{}{},
	}
	if req.StartDate != "" {
		meta.Filters["start_date"] = req.StartDate
	}
	if req.EndDate != "" {
		meta.Filters["end_date"] = req.EndDate
	}

	response.SuccessResponse(c, summary, meta)
}
