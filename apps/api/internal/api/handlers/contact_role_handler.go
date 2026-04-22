package handlers

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/contact_role"
	contactroleservice "github.com/gilabs/crm-healthcare/api/internal/service/contact_role"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ContactRoleHandler struct {
	contactRoleService *contactroleservice.Service
}

func NewContactRoleHandler(contactRoleService *contactroleservice.Service) *ContactRoleHandler {
	return &ContactRoleHandler{
		contactRoleService: contactRoleService,
	}
}

// List handles list contact roles request
func (h *ContactRoleHandler) List(c *gin.Context) {
	contactRoles, err := h.contactRoleService.List()
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, contactRoles, nil)
}

// GetByID handles get contact role by ID request
func (h *ContactRoleHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	contactRole, err := h.contactRoleService.GetByID(id)
	if err != nil {
		if err == contactroleservice.ErrContactRoleNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":        "contact_role",
				"contact_role_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, contactRole, nil)
}

// Create handles create contact role request
func (h *ContactRoleHandler) Create(c *gin.Context) {
	var req contact_role.CreateContactRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	createdContactRole, err := h.contactRoleService.Create(&req)
	if err != nil {
		if err == contactroleservice.ErrContactRoleAlreadyExists {
			errors.ErrorResponse(c, "RESOURCE_ALREADY_EXISTS", map[string]interface{}{
				"resource": "contact_role",
				"field":    "code",
				"value":    req.Code,
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

	response.SuccessResponseCreated(c, createdContactRole, meta)
}

// Update handles update contact role request
func (h *ContactRoleHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req contact_role.UpdateContactRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	updatedContactRole, err := h.contactRoleService.Update(id, &req)
	if err != nil {
		if err == contactroleservice.ErrContactRoleNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":        "contact_role",
				"contact_role_id": id,
			}, nil)
			return
		}
		if err == contactroleservice.ErrContactRoleAlreadyExists {
			errors.ErrorResponse(c, "RESOURCE_ALREADY_EXISTS", map[string]interface{}{
				"resource": "contact_role",
				"field":    "code",
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

	response.SuccessResponse(c, updatedContactRole, meta)
}

// Delete handles delete contact role request
func (h *ContactRoleHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.contactRoleService.Delete(id)
	if err != nil {
		if err == contactroleservice.ErrContactRoleNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":        "contact_role",
				"contact_role_id": id,
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

	response.SuccessResponseDeleted(c, "contact_role", id, meta)
}

