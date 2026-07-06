package product

import (
	"time"

	"github.com/gilabs/crm-healthcare/api/pkg/util/currency"
)

// ProductCategoryResponse represents category response DTO.
type ProductCategoryResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Description  string    `json:"description"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ProductCount int64     `json:"product_count"`
}

// ToProductCategoryResponse converts ProductCategory to ProductCategoryResponse.
func (c *ProductCategory) ToProductCategoryResponse() *ProductCategoryResponse {
	return &ProductCategoryResponse{
		ID:           c.ID,
		Name:         c.Name,
		Slug:         c.Slug,
		Description:  c.Description,
		Status:       c.Status,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
		ProductCount: c.ProductCount,
	}
}

// ProductResponse represents product response DTO.
type ProductResponse struct {
	ID             string                   `json:"id"`
	Name           string                   `json:"name"`
	SKU            string                   `json:"sku"`
	Barcode        string                   `json:"barcode"`
	Price          int64                    `json:"price"`
	PriceFormatted string                   `json:"price_formatted,omitempty"`
	Cost           int64                    `json:"cost"`
	CategoryID     string                   `json:"category_id"`
	Category       *ProductCategoryResponse `json:"category,omitempty"`
	Description    string                   `json:"description"`
	Status         string                   `json:"status"`
	ImageURL       string                   `json:"image_url"`
	CreatedAt      time.Time                `json:"created_at"`
	UpdatedAt      time.Time                `json:"updated_at"`
}

// ToProductResponse converts Product to ProductResponse.
func (p *Product) ToProductResponse() *ProductResponse {
	resp := &ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		SKU:         p.SKU,
		Barcode:     p.Barcode,
		Price:       p.Price,
		Cost:        p.Cost,
		CategoryID:  p.CategoryID,
		Description: p.Description,
		Status:      p.Status,
		ImageURL:    p.ImageURL,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}

	if p.Category != nil {
		resp.Category = &ProductCategoryResponse{
			ID:          p.Category.ID,
			Name:        p.Category.Name,
			Slug:        p.Category.Slug,
			Description: "",
			Status:      "",
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
		}
	}

	resp.PriceFormatted = formatCurrency(p.Price)
	return resp
}

// CreateProductRequest represents create product request DTO.
type CreateProductRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=200"`
	SKU         string `json:"sku" binding:"required,min=1,max=100"`
	Barcode     string `json:"barcode" binding:"omitempty,max=100"`
	Price       int64  `json:"price" binding:"required,min=0"`
	Cost        int64  `json:"cost" binding:"omitempty,min=0"`
	CategoryID  string `json:"category_id" binding:"required,uuid"`
	Description string `json:"description" binding:"omitempty"`
	Status      string `json:"status" binding:"omitempty,oneof=active inactive"`
	ImageURL    string `json:"image_url" binding:"omitempty,max=500"`
}

// UpdateProductRequest represents update product request DTO.
type UpdateProductRequest struct {
	Name        string `json:"name" binding:"omitempty,min=3,max=200"`
	SKU         string `json:"sku" binding:"omitempty,min=1,max=100"`
	Barcode     string `json:"barcode" binding:"omitempty,max=100"`
	Price       *int64 `json:"price" binding:"omitempty,min=0"`
	Cost        *int64 `json:"cost" binding:"omitempty,min=0"`
	CategoryID  string `json:"category_id" binding:"omitempty,uuid"`
	Description string `json:"description" binding:"omitempty"`
	Status      string `json:"status" binding:"omitempty,oneof=active inactive"`
	ImageURL    string `json:"image_url" binding:"omitempty,max=500"`
}

// ListProductsRequest represents list products query parameters.
type ListProductsRequest struct {
	Page       int    `form:"page" binding:"omitempty,min=1"`
	PerPage    int    `form:"per_page" binding:"omitempty,min=1,max=100"`
	Search     string `form:"search" binding:"omitempty"`
	Status     string `form:"status" binding:"omitempty,oneof=active inactive"`
	CategoryID string `form:"category_id" binding:"omitempty,uuid"`
}

// ListProductCategoriesRequest represents list product categories query parameters.
type ListProductCategoriesRequest struct {
	Status string `form:"status" binding:"omitempty,oneof=active inactive"`
}

// CreateProductCategoryRequest represents create product category request DTO.
type CreateProductCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=150"`
	Slug        string `json:"slug" binding:"omitempty,max=150"`
	Description string `json:"description" binding:"omitempty"`
	Status      string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// UpdateProductCategoryRequest represents update product category request DTO.
type UpdateProductCategoryRequest struct {
	Name        string `json:"name" binding:"omitempty,min=3,max=150"`
	Slug        string `json:"slug" binding:"omitempty,max=150"`
	Description string `json:"description" binding:"omitempty"`
	Status      string `json:"status" binding:"omitempty,oneof=active inactive"`
}

func formatCurrency(amount int64) string {
	return currency.FormatCurrency(amount)
}
