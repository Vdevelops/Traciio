package handlers

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/product"
	productservice "github.com/gilabs/crm-healthcare/api/internal/service/product"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ProductHandler struct {
	productService *productservice.Service
}

func NewProductHandler(productService *productservice.Service) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// List handles list products request.
func (h *ProductHandler) List(c *gin.Context) {
	var req product.ListProductsRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	products, pagination, err := h.productService.ListProducts(&req)
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

	response.SuccessResponse(c, products, meta)
}

// GetByID handles get product by ID request.
func (h *ProductHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	p, err := h.productService.GetProductByID(id)
	if err != nil {
		if err == productservice.ErrProductNotFound {
			errors.ErrorResponse(c, "PRODUCT_NOT_FOUND", map[string]interface{}{
				"resource":    "product",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, p, nil)
}

// Create handles create product request.
func (h *ProductHandler) Create(c *gin.Context) {
	var req product.CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	createdProduct, err := h.productService.CreateProduct(&req)
	if err != nil {
		if err == productservice.ErrProductCategoryNotFound {
			errors.ErrorResponse(c, "CATEGORY_NOT_FOUND", map[string]interface{}{
				"category_id": req.CategoryID,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponseCreated(c, createdProduct, nil)
}

// Update handles update product request.
func (h *ProductHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req product.UpdateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	updatedProduct, err := h.productService.UpdateProduct(id, &req)
	if err != nil {
		if err == productservice.ErrProductNotFound {
			errors.ErrorResponse(c, "PRODUCT_NOT_FOUND", map[string]interface{}{
				"resource":    "product",
				"resource_id": id,
			}, nil)
			return
		}
		if err == productservice.ErrProductCategoryNotFound {
			errors.ErrorResponse(c, "CATEGORY_NOT_FOUND", map[string]interface{}{
				"category_id": req.CategoryID,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, updatedProduct, nil)
}

// Delete handles delete product request.
func (h *ProductHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.productService.DeleteProduct(id)
	if err != nil {
		if err == productservice.ErrProductNotFound {
			errors.ErrorResponse(c, "PRODUCT_NOT_FOUND", map[string]interface{}{
				"resource":    "product",
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

	response.SuccessResponseDeleted(c, "product", id, meta)
}

// ListCategories handles list product categories request.
func (h *ProductHandler) ListCategories(c *gin.Context) {
	var req product.ListProductCategoriesRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	categories, err := h.productService.ListProductCategories(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, categories, nil)
}

// GetCategoryByID handles get product category by ID request.
func (h *ProductHandler) GetCategoryByID(c *gin.Context) {
	id := c.Param("id")

	category, err := h.productService.GetProductCategoryByID(id)
	if err != nil {
		if err == productservice.ErrProductCategoryNotFound {
			errors.ErrorResponse(c, "CATEGORY_NOT_FOUND", map[string]interface{}{
				"resource":    "product_category",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, category, nil)
}

// CreateCategory handles create product category request.
func (h *ProductHandler) CreateCategory(c *gin.Context) {
	var req product.CreateProductCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	category, err := h.productService.CreateProductCategory(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponseCreated(c, category, nil)
}

// UpdateCategory handles update product category request.
func (h *ProductHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var req product.UpdateProductCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	category, err := h.productService.UpdateProductCategory(id, &req)
	if err != nil {
		if err == productservice.ErrProductCategoryNotFound {
			errors.ErrorResponse(c, "CATEGORY_NOT_FOUND", map[string]interface{}{
				"resource":    "product_category",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, category, nil)
}

// DeleteCategory handles delete product category request.
func (h *ProductHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	err := h.productService.DeleteProductCategory(id)
	if err != nil {
		if err == productservice.ErrProductCategoryNotFound {
			errors.ErrorResponse(c, "CATEGORY_NOT_FOUND", map[string]interface{}{
				"resource":    "product_category",
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

	response.SuccessResponseDeleted(c, "product_category", id, meta)
}


