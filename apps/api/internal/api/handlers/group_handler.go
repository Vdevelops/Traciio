package handlers

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/group"
	groupservice "github.com/gilabs/crm-healthcare/api/internal/service/group"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type GroupHandler struct {
	groupService *groupservice.Service
}

func NewGroupHandler(groupService *groupservice.Service) *GroupHandler {
	return &GroupHandler{
		groupService: groupService,
	}
}

// List handles list groups request
func (h *GroupHandler) List(c *gin.Context) {
	var req group.ListGroupsRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	groups, pagination, err := h.groupService.List(&req)
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

	response.SuccessResponse(c, groups, meta)
}

// GetByID handles get group by ID request
func (h *GroupHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	group, err := h.groupService.GetByID(id)
	if err != nil {
		if err == groupservice.ErrGroupNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "group",
				"group_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, group, nil)
}

// Create handles create group request
func (h *GroupHandler) Create(c *gin.Context) {
	var req group.CreateGroupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	createdGroup, err := h.groupService.Create(&req)
	if err != nil {
		if err == groupservice.ErrGroupAlreadyExists {
			errors.ErrorResponse(c, "RESOURCE_ALREADY_EXISTS", map[string]interface{}{
				"resource": "group",
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

	response.SuccessResponseCreated(c, createdGroup, meta)
}

// Update handles update group request
func (h *GroupHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req group.UpdateGroupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	updatedGroup, err := h.groupService.Update(id, &req)
	if err != nil {
		if err == groupservice.ErrGroupNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "group",
				"group_id": id,
			}, nil)
			return
		}
		if err == groupservice.ErrGroupAlreadyExists {
			errors.ErrorResponse(c, "RESOURCE_ALREADY_EXISTS", map[string]interface{}{
				"resource": "group",
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
			meta.UpdatedBy = id
		}
	}

	response.SuccessResponse(c, updatedGroup, meta)
}

// Delete handles delete group request
func (h *GroupHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.groupService.Delete(id)
	if err != nil {
		if err == groupservice.ErrGroupNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "group",
				"group_id": id,
			}, nil)
			return
		}
		if err == groupservice.ErrGroupInUse {
			errors.ErrorResponse(c, "RESOURCE_IN_USE", map[string]interface{}{
				"resource":    "group",
				"group_id": id,
				"reason":      "group has associated users",
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			meta.DeletedBy = id
		}
	}

	response.SuccessResponse(c, map[string]interface{}{
		"message": "Group deleted successfully",
	}, meta)
}

