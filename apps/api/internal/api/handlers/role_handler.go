package handlers

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	roleservice "github.com/gilabs/crm-healthcare/api/internal/service/role"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type RoleHandler struct {
	roleService *roleservice.Service
}

func NewRoleHandler(roleService *roleservice.Service) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

// List handles list roles request
func (h *RoleHandler) List(c *gin.Context) {
	roles, err := h.roleService.List()
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, roles, nil)
}

// GetByID handles get role by ID request
func (h *RoleHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	role, err := h.roleService.GetByID(id)
	if err != nil {
		if err == roleservice.ErrRoleNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource": "role",
				"role_id":  id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, role, nil)
}

// Create handles create role request
func (h *RoleHandler) Create(c *gin.Context) {
	var req role.CreateRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	createdRole, err := h.roleService.Create(&req)
	if err != nil {
		if err == roleservice.ErrRoleAlreadyExists {
			errors.ErrorResponse(c, "RESOURCE_ALREADY_EXISTS", map[string]interface{}{
				"resource": "role",
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

	response.SuccessResponseCreated(c, createdRole, meta)
}

// Update handles update role request
func (h *RoleHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req role.UpdateRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	updatedRole, err := h.roleService.Update(id, &req)
	if err != nil {
		if err == roleservice.ErrRoleNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource": "role",
				"role_id":  id,
			}, nil)
			return
		}
		if err == roleservice.ErrRoleAlreadyExists {
			errors.ErrorResponse(c, "RESOURCE_ALREADY_EXISTS", map[string]interface{}{
				"resource": "role",
				"field":    "code",
			}, nil)
			return
		}
		if err == roleservice.ErrRoleProtected {
			errors.ErrorResponse(c, "RESOURCE_PROTECTED", map[string]interface{}{
				"resource": "role",
				"role_id":  id,
				"message":  "This role is protected and cannot be modified",
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

	response.SuccessResponse(c, updatedRole, meta)
}

// Delete handles delete role request
func (h *RoleHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.roleService.Delete(id)
	if err != nil {
		if err == roleservice.ErrRoleNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource": "role",
				"role_id":  id,
			}, nil)
			return
		}
		if err == roleservice.ErrRoleProtected {
			errors.ErrorResponse(c, "RESOURCE_PROTECTED", map[string]interface{}{
				"resource": "role",
				"role_id":  id,
				"message":  "This role is protected and cannot be deleted",
			}, nil)
			return
		}
		if err == roleservice.ErrRoleInUse {
			errors.ErrorResponse(c, "RESOURCE_IN_USE", map[string]interface{}{
				"resource": "role",
				"role_id":  id,
				"message":  "This role is in use and cannot be deleted. Please reassign users first",
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

	response.SuccessResponseDeleted(c, "role", id, meta)
}

// GetRolePermissions handles get role permissions request
func (h *RoleHandler) GetRolePermissions(c *gin.Context) {
	id := c.Param("id")

	permissions, err := h.roleService.GetRolePermissions(id)
	if err != nil {
		if err == roleservice.ErrRoleNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource": "role",
				"role_id":  id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, permissions, nil)
}

// AssignPermissions handles assign permissions to role request
func (h *RoleHandler) AssignPermissions(c *gin.Context) {
	id := c.Param("id")
	var req role.AssignPermissionsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	err := h.roleService.AssignPermissions(id, req.PermissionIDs)
	if err != nil {
		if err == roleservice.ErrRoleNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource": "role",
				"role_id":  id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	// Return updated role
	updatedRole, err := h.roleService.GetByID(id)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, updatedRole, nil)
}

// GetRoleScopes handles get role scopes request
func (h *RoleHandler) GetRoleScopes(c *gin.Context) {
	id := c.Param("id")

	scopes, err := h.roleService.GetRoleScopes(id)
	if err != nil {
		if err == roleservice.ErrRoleNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource": "role",
				"role_id":  id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, scopes, nil)
}

// UpdateRoleScopes handles update role scopes request
func (h *RoleHandler) UpdateRoleScopes(c *gin.Context) {
	id := c.Param("id")
	var req role.UpdateRoleScopesRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	err := h.roleService.UpdateRoleScopes(id, req.Scopes)
	if err != nil {
		if err == roleservice.ErrRoleNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource": "role",
				"role_id":  id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	scopes, err := h.roleService.GetRoleScopes(id)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, scopes, nil)
}
