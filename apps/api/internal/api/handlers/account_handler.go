package handlers

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/account"
	accountservice "github.com/gilabs/crm-healthcare/api/internal/service/account"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AccountHandler struct {
	accountService *accountservice.Service
}

func NewAccountHandler(accountService *accountservice.Service) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
	}
}

// List handles list accounts request
func (h *AccountHandler) List(c *gin.Context) {
	var req account.ListAccountsRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	accounts, pagination, err := h.accountService.List(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
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
	if req.CategoryID != "" {
		meta.Filters["category_id"] = req.CategoryID
	}
	if req.AssignedTo != "" {
		meta.Filters["assigned_to"] = req.AssignedTo
	}

	response.SuccessResponse(c, accounts, meta)
}

// GetByID handles get account by ID request
func (h *AccountHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	account, err := h.accountService.GetByID(id)
	if err != nil {
		if err == accountservice.ErrAccountNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "account",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, account, nil)
}

// Create handles create account request
func (h *AccountHandler) Create(c *gin.Context) {
	var req account.CreateAccountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	createdAccount, err := h.accountService.Create(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			meta.CreatedBy = id
		}
	}

	response.SuccessResponseCreated(c, createdAccount, meta)
}

// Update handles update account request
func (h *AccountHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req account.UpdateAccountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	updatedAccount, err := h.accountService.Update(id, &req)
	if err != nil {
		if err == accountservice.ErrAccountNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "account",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			meta.UpdatedBy = id
		}
	}

	response.SuccessResponse(c, updatedAccount, meta)
}

// Delete handles delete account request
func (h *AccountHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.accountService.Delete(id)
	if err != nil {
		if err == accountservice.ErrAccountNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "account",
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

	response.SuccessResponseDeleted(c, "account", id, meta)
}

// ListAllForMap handles list all accounts for map display request (no pagination, with geocoding)
func (h *AccountHandler) ListAllForMap(c *gin.Context) {
	status := c.Query("status") // Optional status filter

	accounts, err := h.accountService.ListAllForMap(status)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{
		Filters: map[string]interface{}{},
	}

	if status != "" {
		meta.Filters["status"] = status
	}

	response.SuccessResponse(c, accounts, meta)
}

// ListByBBox handles viewport-based map loading using bounding box coordinates
func (h *AccountHandler) ListByBBox(c *gin.Context) {
	var req account.BBoxRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	accounts, err := h.accountService.ListByBBox(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{
		Filters: map[string]interface{}{
			"bbox": map[string]float64{
				"north": req.North,
				"south": req.South,
				"east":  req.East,
				"west":  req.West,
			},
			"count": len(accounts),
		},
	}

	if req.Status != "" {
		meta.Filters["status"] = req.Status
	}
	if req.CategoryID != "" {
		meta.Filters["category_id"] = req.CategoryID
	}

	response.SuccessResponse(c, accounts, meta)
}

