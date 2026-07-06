package handlers

import (
	"strings"

	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/internal/domain/ai"
	aiservice "github.com/gilabs/crm-healthcare/api/internal/service/ai"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AIHandler struct {
	aiService *aiservice.Service
}

func NewAIHandler(aiService *aiservice.Service) *AIHandler {
	return &AIHandler{
		aiService: aiService,
	}
}

// AnalyzeVisitReport handles visit report analysis request
func (h *AIHandler) AnalyzeVisitReport(c *gin.Context) {
	userCtx := middleware.GetUserContext(c)
	if userCtx == nil || !userCtx.HasPermission("ai-chatbot.view") {
		errors.ErrorResponse(c, "FORBIDDEN", map[string]interface{}{
			"required": "ai-chatbot.view",
			"message":  "You do not have permission to access AI Assistant",
		}, nil)
		return
	}

	var req ai.AnalyzeVisitReportRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

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

	insight, tokens, err := h.aiService.AnalyzeVisitReport(req.VisitReportID, userIDStr, userCtx)
	if err != nil {
		// Check for specific errors
		if err.Error() == "AI service not configured: Cerebras API key is empty" {
			errors.ErrorResponse(c, "AI_SERVICE_NOT_CONFIGURED", map[string]interface{}{
				"error": "Cerebras API key is not configured. Please set CEREBRAS_API_KEY environment variable",
			}, nil)
			return
		}
		errors.ErrorResponse(c, "AI_ANALYSIS_FAILED", map[string]interface{}{
			"error": err.Error(),
		}, nil)
		return
	}

	response.SuccessResponse(c, &ai.InsightResponse{
		Type:   ai.InsightTypeVisitReport,
		Data:   insight,
		Tokens: tokens,
	}, nil)
}

// Chat handles chat request
func (h *AIHandler) Chat(c *gin.Context) {
	userCtx := middleware.GetUserContext(c)
	if userCtx == nil || !userCtx.HasPermission("ai-chatbot.view") {
		errors.ErrorResponse(c, "FORBIDDEN", map[string]interface{}{
			"required": "ai-chatbot.view",
			"message":  "You do not have permission to access AI Assistant",
		}, nil)
		return
	}

	var req ai.ChatRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	// Use conversation history if provided, otherwise use empty slice
	history := req.ConversationHistory
	if history == nil {
		history = []ai.ChatMessage{}
	}

	// Get user ID from context (set by AuthMiddleware)
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

	chatResponse, err := h.aiService.Chat(req.Message, req.Context, req.ContextType, history, req.Model, userIDStr, req.Domain, userCtx)
	if err != nil {
		// Check for specific errors
		errMsg := err.Error()

		if strings.Contains(errMsg, "AI service not configured") || strings.Contains(errMsg, "API key is empty") {
			errors.ErrorResponse(c, "AI_SERVICE_NOT_CONFIGURED", map[string]interface{}{
				"error": "Cerebras API key is not configured. Please set CEREBRAS_API_KEY environment variable",
			}, nil)
			return
		}

		// Check for model not found errors
		if strings.Contains(errMsg, "tidak ditemukan") || strings.Contains(errMsg, "does not exist") || strings.Contains(errMsg, "model_not_found") {
			errors.ErrorResponse(c, "AI_MODEL_NOT_FOUND", map[string]interface{}{
				"error": errMsg,
			}, nil)
			return
		}

		// Check for unsupported model errors
		if strings.Contains(errMsg, "tidak didukung") || strings.Contains(errMsg, "not supported") {
			errors.ErrorResponse(c, "AI_MODEL_NOT_SUPPORTED", map[string]interface{}{
				"error": errMsg,
			}, nil)
			return
		}

		errors.ErrorResponse(c, "AI_CHAT_FAILED", map[string]interface{}{
			"error": errMsg,
		}, nil)
		return
	}

	response.SuccessResponse(c, chatResponse, nil)
}
