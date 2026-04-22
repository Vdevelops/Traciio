package handlers

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
	pipelineservice "github.com/gilabs/crm-healthcare/api/internal/service/pipeline"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PipelineHandler struct {
	pipelineService *pipelineservice.Service
}

func NewPipelineHandler(pipelineService *pipelineservice.Service) *PipelineHandler {
	return &PipelineHandler{
		pipelineService: pipelineService,
	}
}

// ListStages handles list pipeline stages request
func (h *PipelineHandler) ListStages(c *gin.Context) {
	var req pipeline.ListPipelineStagesRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	stages, err := h.pipelineService.ListStages(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{
		Filters: map[string]interface{}{},
	}

	if req.IsActive != nil {
		meta.Filters["is_active"] = *req.IsActive
	}

	response.SuccessResponse(c, stages, meta)
}

// GetStageByID handles get pipeline stage by ID request
func (h *PipelineHandler) GetStageByID(c *gin.Context) {
	id := c.Param("id")

	stage, err := h.pipelineService.GetStageByID(id)
	if err != nil {
		if err == pipelineservice.ErrPipelineStageNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "pipeline_stage",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, stage, nil)
}

// CreateStage handles create pipeline stage request
func (h *PipelineHandler) CreateStage(c *gin.Context) {
	var req pipeline.CreateStageRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	stage, err := h.pipelineService.CreateStage(&req)
	if err != nil {
		if err.Error() == "pipeline stage with this code already exists" {
			errors.ErrorResponse(c, "DUPLICATE_ENTRY", map[string]interface{}{
				"field": "code",
				"value": req.Code,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			meta.CreatedBy = id
		}
	}

	response.SuccessResponse(c, stage, meta)
}

// UpdateStage handles update pipeline stage request
func (h *PipelineHandler) UpdateStage(c *gin.Context) {
	id := c.Param("id")
	var req pipeline.UpdateStageRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	stage, err := h.pipelineService.UpdateStage(id, &req)
	if err != nil {
		if err == pipelineservice.ErrPipelineStageNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "pipeline_stage",
				"resource_id": id,
			}, nil)
			return
		}
		if err.Error() == "pipeline stage with this code already exists" {
			errors.ErrorResponse(c, "DUPLICATE_ENTRY", map[string]interface{}{
				"field": "code",
				"value": req.Code,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(string); ok {
			meta.UpdatedBy = uid
		}
	}

	response.SuccessResponse(c, stage, meta)
}

// DeleteStage handles delete pipeline stage request
func (h *PipelineHandler) DeleteStage(c *gin.Context) {
	id := c.Param("id")

	err := h.pipelineService.DeleteStage(id)
	if err != nil {
		if err == pipelineservice.ErrPipelineStageNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "pipeline_stage",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	// Get user ID for meta
	meta := &response.Meta{}
	if userIDVal, exists := c.Get("user_id"); exists {
		if id, ok := userIDVal.(string); ok {
			meta.DeletedBy = id
		}
	}

	response.SuccessResponseDeleted(c, "pipeline_stage", id, meta)
}

// UpdateStagesOrder handles update stages order request
func (h *PipelineHandler) UpdateStagesOrder(c *gin.Context) {
	var req pipeline.UpdateStagesOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	stages, err := h.pipelineService.UpdateStagesOrder(&req)
	if err != nil {
		if err == pipelineservice.ErrPipelineStageNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource": "pipeline_stage",
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(string); ok {
			meta.UpdatedBy = uid
		}
	}

	response.SuccessResponse(c, stages, meta)
}

// GetSummary handles get pipeline summary request
func (h *PipelineHandler) GetSummary(c *gin.Context) {
	summary, err := h.pipelineService.GetSummary()
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, summary, nil)
}

// GetForecast handles get forecast request
func (h *PipelineHandler) GetForecast(c *gin.Context) {
	periodType := c.DefaultQuery("period", "month")
	if periodType != "month" && periodType != "quarter" && periodType != "year" {
		periodType = "month"
	}

	forecast, err := h.pipelineService.GetForecast(periodType)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{
		Filters: map[string]interface{}{
			"period": periodType,
		},
	}

	response.SuccessResponse(c, forecast, meta)
}

// Flow Rules Handlers
// TODO: Flow Rules feature is not yet implemented. Uncomment when implementing flow rules.
/*
// ListFlowRules handles list flow rules request
func (h *PipelineHandler) ListFlowRules(c *gin.Context) {
	var req pipeline.ListFlowRulesRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	rules, err := h.pipelineService.ListFlowRules(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, rules, nil)
}

// GetFlowRuleByID handles get flow rule by ID request
func (h *PipelineHandler) GetFlowRuleByID(c *gin.Context) {
	id := c.Param("id")

	rule, err := h.pipelineService.GetFlowRuleByID(id)
	if err != nil {
		if err == pipelineservice.ErrFlowRuleNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "flow_rule",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, rule, nil)
}

// CreateFlowRule handles create flow rule request
func (h *PipelineHandler) CreateFlowRule(c *gin.Context) {
	var req pipeline.CreateFlowRuleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	rule, err := h.pipelineService.CreateFlowRule(&req)
	if err != nil {
		if err == pipelineservice.ErrInvalidStage {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "pipeline_stage",
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(string); ok {
			meta.CreatedBy = uid
		}
	}

	response.SuccessResponseCreated(c, rule, meta)
}

// UpdateFlowRule handles update flow rule request
func (h *PipelineHandler) UpdateFlowRule(c *gin.Context) {
	id := c.Param("id")
	var req pipeline.UpdateFlowRuleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	rule, err := h.pipelineService.UpdateFlowRule(id, &req)
	if err != nil {
		if err == pipelineservice.ErrFlowRuleNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "flow_rule",
				"resource_id": id,
			}, nil)
			return
		}
		if err == pipelineservice.ErrInvalidStage {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "pipeline_stage",
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(string); ok {
			meta.UpdatedBy = uid
		}
	}

	response.SuccessResponse(c, rule, meta)
}

// DeleteFlowRule handles delete flow rule request
func (h *PipelineHandler) DeleteFlowRule(c *gin.Context) {
	id := c.Param("id")

	err := h.pipelineService.DeleteFlowRule(id)
	if err != nil {
		if err == pipelineservice.ErrFlowRuleNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "flow_rule",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	// Get user ID for meta
	meta := &response.Meta{}
	if userIDVal, exists := c.Get("user_id"); exists {
		if id, ok := userIDVal.(string); ok {
			meta.DeletedBy = id
		}
	}

	response.SuccessResponseDeleted(c, "flow_rule", id, meta)
}

// ValidateTransition handles validate transition request
func (h *PipelineHandler) ValidateTransition(c *gin.Context) {
	var req pipeline.ValidateTransitionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	// Get user role from context if not provided
	if req.UserRole == "" {
		if role, exists := c.Get("user_role"); exists {
			if r, ok := role.(string); ok {
				req.UserRole = r
			}
		}
	}

	validation, err := h.pipelineService.ValidateTransition(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, validation, nil)
}
*/

// MoveStage handles move deal to another stage request.
// Requires "pipeline.update_stage" permission.
func (h *PipelineHandler) MoveStage(c *gin.Context) {
	// Check RBAC permission
	if userCtx := middleware.GetUserContext(c); userCtx != nil {
		if !userCtx.HasPermission("pipeline.update_stage") {
			errors.ForbiddenResponse(c, "pipeline.update_stage", nil)
			return
		}
	}

	id := c.Param("id")
	var req pipeline.MoveStageRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		errors.UnauthorizedResponse(c, "")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		errors.UnauthorizedResponse(c, "")
		return
	}

	deal, err := h.pipelineService.MoveStageWithValidation(id, req.ToStageID, userIDStr, req.Reason)
	if err != nil {
		if err == pipelineservice.ErrDealNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "deal",
				"resource_id": id,
			}, nil)
			return
		}
		if err == pipelineservice.ErrPipelineStageNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "pipeline_stage",
				"resource_id": req.ToStageID,
			}, nil)
			return
		}
		// Stage requirement validation errors (backward move, missing products, value = 0, etc.)
		if err == pipelineservice.ErrStageRequirementsNotMet {
			errors.ErrorResponse(c, "STAGE_REQUIREMENTS_NOT_MET", map[string]interface{}{
				"message": "Stage transition requirements not met. Check products, deal value, and stage order.",
			}, nil)
			return
		}
		// Catch-all for other validation errors from ValidateStageRequirements
		errors.ErrorResponse(c, "VALIDATION_ERROR", map[string]interface{}{
			"message": err.Error(),
		}, nil)
		return
	}

	meta := &response.Meta{}
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			meta.UpdatedBy = id
		}
	}

	response.SuccessResponse(c, deal, meta)
}

// GetDealHistory handles get deal history request
func (h *PipelineHandler) GetDealHistory(c *gin.Context) {
	id := c.Param("id")

	history, err := h.pipelineService.GetDealHistory(id)
	if err != nil {
		if err == pipelineservice.ErrDealNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "deal",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, history, nil)
}
