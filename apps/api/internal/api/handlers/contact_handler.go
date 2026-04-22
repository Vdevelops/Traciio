package handlers

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/contact"
	contactservice "github.com/gilabs/crm-healthcare/api/internal/service/contact"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ContactHandler struct {
	contactService *contactservice.Service
}

func NewContactHandler(contactService *contactservice.Service) *ContactHandler {
	return &ContactHandler{
		contactService: contactService,
	}
}

// List handles list contacts request
func (h *ContactHandler) List(c *gin.Context) {
	var req contact.ListContactsRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	contacts, pagination, err := h.contactService.List(&req)
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
	if req.AccountID != "" {
		meta.Filters["account_id"] = req.AccountID
	}
	if req.RoleID != "" {
		meta.Filters["role_id"] = req.RoleID
	}

	response.SuccessResponse(c, contacts, meta)
}

// GetByID handles get contact by ID request
func (h *ContactHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	contact, err := h.contactService.GetByID(id)
	if err != nil {
		if err == contactservice.ErrContactNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "contact",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, contact, nil)
}

// Create handles create contact request
func (h *ContactHandler) Create(c *gin.Context) {
	var req contact.CreateContactRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	createdContact, err := h.contactService.Create(&req)
	if err != nil {
		if err == contactservice.ErrAccountNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "account",
				"resource_id": req.AccountID,
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

	response.SuccessResponseCreated(c, createdContact, meta)
}

// Update handles update contact request
func (h *ContactHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req contact.UpdateContactRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	updatedContact, err := h.contactService.Update(id, &req)
	if err != nil {
		if err == contactservice.ErrContactNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "contact",
				"resource_id": id,
			}, nil)
			return
		}
		if err == contactservice.ErrAccountNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "account",
				"resource_id": req.AccountID,
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

	response.SuccessResponse(c, updatedContact, meta)
}

// Delete handles delete contact request
func (h *ContactHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.contactService.Delete(id)
	if err != nil {
		if err == contactservice.ErrContactNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "contact",
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

	response.SuccessResponseDeleted(c, "contact", id, meta)
}

