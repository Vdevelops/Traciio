package handlers

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/category"
	categoryservice "github.com/gilabs/crm-healthcare/api/internal/service/category"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CategoryHandler struct {
	categoryService *categoryservice.Service
}

func NewCategoryHandler(categoryService *categoryservice.Service) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// List handles list categories request
func (h *CategoryHandler) List(c *gin.Context) {
	categories, err := h.categoryService.List()
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, categories, nil)
}

// GetByID handles get category by ID request
func (h *CategoryHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	category, err := h.categoryService.GetByID(id)
	if err != nil {
		if err == categoryservice.ErrCategoryNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "category",
				"category_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, category, nil)
}

// Create handles create category request
func (h *CategoryHandler) Create(c *gin.Context) {
	var req category.CreateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	createdCategory, err := h.categoryService.Create(&req)
	if err != nil {
		if err == categoryservice.ErrCategoryAlreadyExists {
			errors.ErrorResponse(c, "RESOURCE_ALREADY_EXISTS", map[string]interface{}{
				"resource": "category",
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

	response.SuccessResponseCreated(c, createdCategory, meta)
}

// Update handles update category request
func (h *CategoryHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req category.UpdateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	updatedCategory, err := h.categoryService.Update(id, &req)
	if err != nil {
		if err == categoryservice.ErrCategoryNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "category",
				"category_id": id,
			}, nil)
			return
		}
		if err == categoryservice.ErrCategoryAlreadyExists {
			errors.ErrorResponse(c, "RESOURCE_ALREADY_EXISTS", map[string]interface{}{
				"resource": "category",
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

	response.SuccessResponse(c, updatedCategory, meta)
}

// Delete handles delete category request
func (h *CategoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.categoryService.Delete(id)
	if err != nil {
		if err == categoryservice.ErrCategoryNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "category",
				"category_id": id,
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

	response.SuccessResponseDeleted(c, "category", id, meta)
}

