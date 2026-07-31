package handlers

import (
	domainkpi "github.com/gilabs/crm-healthcare/api/internal/domain/kpi"
	kpiservice "github.com/gilabs/crm-healthcare/api/internal/service/kpi"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type KPIHandler struct {
	kpiService *kpiservice.Service
}

func NewKPIHandler(kpiService *kpiservice.Service) *KPIHandler {
	return &KPIHandler{kpiService: kpiService}
}

func (h *KPIHandler) GetSalesRepScorecard(c *gin.Context) {
	var req domainkpi.GetSalesRepScorecardRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}
	if req.UserID == "" {
		if userID, exists := c.Get("user_id"); exists {
			if id, ok := userID.(string); ok {
				req.UserID = id
			}
		}
	}
	result, err := h.kpiService.GetSalesRepScorecard(req.UserID, req.StartDate, req.EndDate)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}
	response.SuccessResponse(c, result, nil)
}

func (h *KPIHandler) GetSalesManagerScorecard(c *gin.Context) {
	var req domainkpi.GetSalesManagerScorecardRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}
	if req.ManagerID == "" {
		if userID, exists := c.Get("user_id"); exists {
			if id, ok := userID.(string); ok {
				req.ManagerID = id
			}
		}
	}
	result, err := h.kpiService.GetSalesManagerScorecard(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}
	response.SuccessResponse(c, result, nil)
}