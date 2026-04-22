package handlers

import (
	"strconv"

	"github.com/gilabs/crm-healthcare/api/internal/domain/brick_target_distribution"
	bricktargetdistributionservice "github.com/gilabs/crm-healthcare/api/internal/service/brick_target_distribution"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BrickTargetDistributionHandler struct {
	distributionService *bricktargetdistributionservice.Service
}

func NewBrickTargetDistributionHandler(distributionService *bricktargetdistributionservice.Service) *BrickTargetDistributionHandler {
	return &BrickTargetDistributionHandler{
		distributionService: distributionService,
	}
}

// GetBrickTargetWithDistributions handles get brick target with distributions request
func (h *BrickTargetDistributionHandler) GetBrickTargetWithDistributions(c *gin.Context) {
	brickID := c.Param("id")
	yearStr := c.Param("year")
	monthStr := c.Param("month")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		errors.InvalidQueryParamResponse(c)
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		errors.InvalidQueryParamResponse(c)
		return
	}

	// Check authorization: only admin or manager of this brick can access
	// Note: We need brickRepo to check access - this requires adding it to handler
	// For now, authorization is checked in service layer (ErrInvalidManager)
	// In production, add proper authorization check here

	result, err := h.distributionService.GetBrickTargetWithDistributions(brickID, year, month)
	if err != nil {
		if err.Error() == "brick not found" {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource": "brick",
				"brick_id": brickID,
			}, nil)
			return
		}
		if err == bricktargetdistributionservice.ErrBrickTargetNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource": "brick_target",
				"brick_id": brickID,
				"year":     year,
				"month":    month,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, result, nil)
}

// DistributeTarget handles distribute target to sales request
func (h *BrickTargetDistributionHandler) DistributeTarget(c *gin.Context) {
	brickID := c.Param("id")
	targetID := c.Param("target_id")

	// Get user ID from context (manager performing distribution)
	userID, exists := c.Get("user_id")
	if !exists {
		errors.UnauthorizedResponse(c, "user ID not found in context")
		return
	}
	distributedBy, ok := userID.(string)
	if !ok {
		errors.UnauthorizedResponse(c, "invalid user ID format")
		return
	}

	var req brick_target_distribution.BulkCreateBrickTargetDistributionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	distributions, err := h.distributionService.BulkCreate(brickID, targetID, &req, distributedBy)
	if err != nil {
		if err == bricktargetdistributionservice.ErrInvalidManager {
			errors.ErrorResponse(c, "FORBIDDEN", map[string]interface{}{
				"message": "only brick manager can distribute targets",
			}, nil)
			return
		}
		if err == bricktargetdistributionservice.ErrInvalidSalesBrick {
			errors.ErrorResponse(c, "VALIDATION_ERROR", map[string]interface{}{
				"field":   "sales_user_id",
				"message": "sales user must be in the same brick",
			}, nil)
			return
		}
		if err == bricktargetdistributionservice.ErrBrickTargetNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource": "brick_target",
				"target_id": targetID,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	meta.CreatedBy = distributedBy

	response.SuccessResponseCreated(c, map[string]interface{}{
		"distributed_count": len(distributions),
		"distributions":     distributions,
	}, meta)
}

// UpdateDistribution handles update distribution request
func (h *BrickTargetDistributionHandler) UpdateDistribution(c *gin.Context) {
	brickID := c.Param("id")
	targetID := c.Param("target_id")
	distributionID := c.Param("distribution_id")

	// Get user ID from context (manager performing update)
	userID, exists := c.Get("user_id")
	if !exists {
		errors.UnauthorizedResponse(c, "user ID not found in context")
		return
	}
	distributedBy, ok := userID.(string)
	if !ok {
		errors.UnauthorizedResponse(c, "invalid user ID format")
		return
	}

	var req brick_target_distribution.UpdateBrickTargetDistributionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	updatedDistribution, err := h.distributionService.Update(brickID, targetID, distributionID, &req, distributedBy)
	if err != nil {
		if err == bricktargetdistributionservice.ErrDistributionNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":        "brick_target_distribution",
				"distribution_id": distributionID,
			}, nil)
			return
		}
		if err == bricktargetdistributionservice.ErrInvalidManager {
			errors.ErrorResponse(c, "FORBIDDEN", map[string]interface{}{
				"message": "only brick manager can update distributions",
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	meta.UpdatedBy = distributedBy

	response.SuccessResponse(c, updatedDistribution, meta)
}

// DeleteDistribution handles delete distribution request
func (h *BrickTargetDistributionHandler) DeleteDistribution(c *gin.Context) {
	brickID := c.Param("id")
	targetID := c.Param("target_id")
	distributionID := c.Param("distribution_id")

	// Get user ID from context (manager performing delete)
	userID, exists := c.Get("user_id")
	if !exists {
		errors.UnauthorizedResponse(c, "user ID not found in context")
		return
	}
	distributedBy, ok := userID.(string)
	if !ok {
		errors.UnauthorizedResponse(c, "invalid user ID format")
		return
	}

	err := h.distributionService.Delete(brickID, targetID, distributionID, distributedBy)
	if err != nil {
		if err == bricktargetdistributionservice.ErrDistributionNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":        "brick_target_distribution",
				"distribution_id": distributionID,
			}, nil)
			return
		}
		if err == bricktargetdistributionservice.ErrInvalidManager {
			errors.ErrorResponse(c, "FORBIDDEN", map[string]interface{}{
				"message": "only brick manager can delete distributions",
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	meta.DeletedBy = distributedBy

	response.SuccessResponse(c, map[string]interface{}{
		"message": "Distribution deleted successfully",
	}, meta)
}

