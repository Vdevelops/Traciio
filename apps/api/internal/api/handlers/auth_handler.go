package handlers

import (
	"log"

	stderrors "errors"

	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/internal/domain/auth"
	authservice "github.com/gilabs/crm-healthcare/api/internal/service/auth"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const userInactiveMsg = "User account is inactive"

type AuthHandler struct {
	authService *authservice.Service
}

func NewAuthHandler(authService *authservice.Service) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login handles login request
// @Summary Login user
// @Description Authenticate user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth.LoginRequest true "Login credentials"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req auth.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	loginResponse, err := h.authService.Login(&req)
	if err != nil {
		if stderrors.Is(err, authservice.ErrInvalidCredentials) {
			errors.ErrorResponse(c, "INVALID_CREDENTIALS", nil, nil)
			return
		}
		if stderrors.Is(err, authservice.ErrUserInactive) {
			errors.ErrorResponse(c, "ACCOUNT_DISABLED", map[string]interface{}{
				"reason": userInactiveMsg,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	// Set user_id in context for CSRF token generation
	c.Set("user_id", loginResponse.User.ID)

	// Generate CSRF token for the logged-in user
	csrfToken, err := middleware.GenerateCSRFToken(c)
	if err != nil {
		log.Printf("Failed to generate CSRF token: %v", err)
		// Don't fail the login, just log the error
	} else {
		// Add CSRF token to response
		loginResponse.CSRFToken = csrfToken
	}

	meta := &response.Meta{}
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			meta.CreatedBy = id
		}
	}

	response.SuccessResponse(c, loginResponse, meta)
}

// RefreshToken handles refresh token request
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Refresh token"
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	loginResponse, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		if stderrors.Is(err, authservice.ErrUserNotFound) {
			errors.ErrorResponse(c, "USER_NOT_FOUND", nil, nil)
			return
		}
		if stderrors.Is(err, authservice.ErrUserInactive) {
			errors.ErrorResponse(c, "ACCOUNT_DISABLED", map[string]interface{}{
				"reason": userInactiveMsg,
			}, nil)
			return
		}
		if stderrors.Is(err, authservice.ErrRefreshTokenRevoked) {
			errors.ErrorResponse(c, "REFRESH_TOKEN_REVOKED", map[string]interface{}{
				"reason": "Refresh token has been revoked",
			}, nil)
			return
		}
		if stderrors.Is(err, authservice.ErrRefreshTokenExpired) {
			errors.ErrorResponse(c, "REFRESH_TOKEN_EXPIRED", map[string]interface{}{
				"reason": "Refresh token has expired",
			}, nil)
			return
		}
		if stderrors.Is(err, authservice.ErrRefreshTokenInvalid) {
			errors.ErrorResponse(c, "REFRESH_TOKEN_INVALID", map[string]interface{}{
				"reason": "Refresh token is invalid",
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, loginResponse, nil)
}

// Logout handles logout request
// @Summary Logout user
// @Description Logout user and revoke refresh token
// @Tags auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body map[string]string false "Refresh token (optional)"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	// Try to get user_id from context (set by AuthMiddleware if authenticated)
	// If not authenticated, use "anonymous" for audit logging
	userID := "anonymous"
	if id, exists := c.Get("user_id"); exists {
		if uid, ok := id.(string); ok {
			userID = uid
		}
	}
	
	// Set user_id for audit logging (even if anonymous)
	c.Set("user_id", userID)

	// Refresh token is optional - if not provided, we still return success
	// This allows logout to work even if client doesn't have refresh token
	if err := c.ShouldBindJSON(&req); err == nil && req.RefreshToken != "" {
		// Revoke refresh token if provided
		if err := h.authService.Logout(req.RefreshToken); err != nil {
			// Log error but don't fail logout
			// This allows logout to succeed even if token is invalid
			log.Printf("Warning: Failed to revoke refresh token during logout (user: %s): %v", userID, err)
		}
	}

	// Revoke CSRF token
	middleware.RevokeCSRFToken(c)

	response.SuccessResponse(c, map[string]interface{}{
		"message": "Logged out successfully",
	}, nil)
}

// MobileLogin handles mobile login request
// @Summary Mobile login user
// @Description Authenticate user for mobile app and return JWT tokens (only sales role allowed)
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth.LoginRequest true "Login credentials"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Router /api/v1/auth/mobile/login [post]
func (h *AuthHandler) MobileLogin(c *gin.Context) {
	var req auth.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	loginResponse, err := h.authService.MobileLogin(&req)
	if err != nil {
		if stderrors.Is(err, authservice.ErrInvalidCredentials) {
			errors.ErrorResponse(c, "INVALID_CREDENTIALS", nil, nil)
			return
		}
		if stderrors.Is(err, authservice.ErrUserInactive) {
			errors.ErrorResponse(c, "ACCOUNT_DISABLED", map[string]interface{}{
				"reason": userInactiveMsg,
			}, nil)
			return
		}
		if stderrors.Is(err, authservice.ErrRoleNotAllowed) {
			errors.ErrorResponse(c, "ROLE_INSUFFICIENT", map[string]interface{}{
				"reason": "Only users with role that has mobile_access enabled can login via mobile app",
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	// Set user_id in context for CSRF token generation
	c.Set("user_id", loginResponse.User.ID)

	// Generate CSRF token for the logged-in user
	csrfToken, err := middleware.GenerateCSRFToken(c)
	if err != nil {
		log.Printf("Failed to generate CSRF token: %v", err)
		// Don't fail the login, just log the error
	} else {
		// Add CSRF token to response
		loginResponse.CSRFToken = csrfToken
	}

	meta := &response.Meta{}
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			meta.CreatedBy = id
		}
	}

	response.SuccessResponse(c, loginResponse, meta)
}

