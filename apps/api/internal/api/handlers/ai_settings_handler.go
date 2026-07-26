package handlers

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/internal/domain/ai_settings"
	aisettingsservice "github.com/gilabs/crm-healthcare/api/internal/service/ai_settings"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AISettingsHandler struct {
	settingsService *aisettingsservice.Service
}

func NewAISettingsHandler(settingsService *aisettingsservice.Service) *AISettingsHandler {
	return &AISettingsHandler{
		settingsService: settingsService,
	}
}

// GetSettings handles get AI settings request.
// Chatbot users need read-only runtime settings (enabled/model/privacy) to render the chat UI.
// Editing settings remains restricted to ai-settings.edit.
func (h *AISettingsHandler) GetSettings(c *gin.Context) {
	userCtx := middleware.GetUserContext(c)
	if userCtx == nil || (!userCtx.HasPermission("ai-settings.view") && !userCtx.HasPermission("ai-chatbot.view")) {
		errors.ForbiddenResponse(c, "ai-chatbot.view", []string{"ai-settings.view"})
		return
	}

	settings, err := h.settingsService.GetSettings()
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, settings, nil)
}

// UpdateSettings handles update AI settings request
// AI Settings is universal and only modifiable by admin
func (h *AISettingsHandler) UpdateSettings(c *gin.Context) {
	userCtx := middleware.GetUserContext(c)
	if userCtx == nil || !userCtx.HasPermission("ai-settings.edit") {
		errors.ForbiddenResponse(c, "ai-settings.edit", []string{})
		return
	}

	var req ai_settings.UpdateAISettingsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	settings, err := h.settingsService.UpdateSettings(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, settings, nil)
}
