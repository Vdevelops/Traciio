package product

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/product"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"github.com/gilabs/crm-healthcare/api/pkg/cache"
	"gorm.io/gorm"
)

var (
	ErrProductNotFound         = errors.New("product not found")
	ErrProductCategoryNotFound = errors.New("product category not found")
	ErrInvalidImageURL         = errors.New("invalid image URL")
)

// validateImageURL validates that the image URL uses http or https protocol.
func validateImageURL(imageURL string) error {
	if imageURL == "" {
		return nil // Empty URL is allowed
	}

	parsedURL, err := url.Parse(imageURL)
	if err != nil {
		return ErrInvalidImageURL
	}

	scheme := strings.ToLower(parsedURL.Scheme)
	if scheme != "http" && scheme != "https" {
		return ErrInvalidImageURL
	}

	return nil
}

type Service struct {
	productRepo         interfaces.ProductRepository
	productCategoryRepo interfaces.ProductCategoryRepository
}

func NewService(
	productRepo interfaces.ProductRepository,
	productCategoryRepo interfaces.ProductCategoryRepository,
) *Service {
	return &Service{
		productRepo:         productRepo,
		productCategoryRepo: productCategoryRepo,
	}
}

// PaginationResult represents pagination information.
type PaginationResult struct {
	Page       int
	PerPage    int
	Total      int
	TotalPages int
}

// ListProducts returns a list of products with pagination.
func (s *Service) ListProducts(req *product.ListProductsRequest) ([]product.ProductResponse, *PaginationResult, error) {
	// Try to get from cache
	cacheKey := fmt.Sprintf("products:list:%d:%d:%s:%s:%s", req.Page, req.PerPage, req.Search, req.Status, req.CategoryID)
	var cachedResult struct {
		Products []product.ProductResponse
		Total    int
	}
	
	if cache.Client != nil && cache.Client.IsEnabled() {
		found, _ := cache.Client.Get(cacheKey, &cachedResult)
		if found {
			// Reconstruct pagination
			page := req.Page
			if page < 1 { page = 1 }
			perPage := req.PerPage
			if perPage < 1 { perPage = 20 }
			if perPage > 100 { perPage = 100 }
			
			totalPages := int((int64(cachedResult.Total) + int64(perPage) - 1) / int64(perPage))
			
			pagination := &PaginationResult{
				Page:       page,
				PerPage:    perPage,
				Total:      cachedResult.Total,
				TotalPages: totalPages,
			}
			return cachedResult.Products, pagination, nil
		}
	}

	products, total, err := s.productRepo.List(req)
	if err != nil {
		return nil, nil, err
	}

	responses := make([]product.ProductResponse, len(products))
	for i, p := range products {
		responses[i] = *p.ToProductResponse()
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}

	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	pagination := &PaginationResult{
		Page:       page,
		PerPage:    perPage,
		Total:      int(total),
		TotalPages: totalPages,
	}

	// Set cache
	if cache.Client != nil && cache.Client.IsEnabled() {
		cachedResult.Products = responses
		cachedResult.Total = int(total)
		// Cache for 5 minutes
		_ = cache.Client.Set(cacheKey, cachedResult, 5*time.Minute)
	}

	return responses, pagination, nil
}

// GetProductByID returns a product by ID.
func (s *Service) GetProductByID(id string) (*product.ProductResponse, error) {
	p, err := s.productRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	return p.ToProductResponse(), nil
}

// CreateProduct creates a new product.
func (s *Service) CreateProduct(req *product.CreateProductRequest) (*product.ProductResponse, error) {
	// Validate image URL format.
	if err := validateImageURL(req.ImageURL); err != nil {
		return nil, err
	}

	// Validate category exists.
	_, err := s.productCategoryRepo.FindByID(req.CategoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductCategoryNotFound
		}
		return nil, err
	}

	// Clear cache for product list.
	cacheKeyPattern := "products:list:*"
	
	status := req.Status
	if status == "" {
		status = "active"
	}

	p := &product.Product{
		Name:        req.Name,
		SKU:         req.SKU,
		Barcode:     req.Barcode,
		Price:       req.Price,
		Cost:        req.Cost,
		CategoryID:  req.CategoryID,
		Description: req.Description,
		Status:      status,
		ImageURL:    req.ImageURL,
	}

	if err := s.productRepo.Create(p); err != nil {
		return nil, err
	}

	// Reload to get relations.
	p, err = s.productRepo.FindByID(p.ID)
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	if cache.Client != nil && cache.Client.IsEnabled() {
		_ = cache.Client.DeletePattern(cacheKeyPattern)
	}

	return p.ToProductResponse(), nil
}

// UpdateProduct updates an existing product.
func (s *Service) UpdateProduct(id string, req *product.UpdateProductRequest) (*product.ProductResponse, error) {
	p, err := s.productRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	if req.Name != "" {
		p.Name = req.Name
	}
	if req.SKU != "" {
		p.SKU = req.SKU
	}
	if req.Barcode != "" {
		p.Barcode = req.Barcode
	}
	if req.Price != nil {
		p.Price = *req.Price
	}
	if req.Cost != nil {
		p.Cost = *req.Cost
	}
	if req.ImageURL != "" {
		// Validate image URL format.
		if err := validateImageURL(req.ImageURL); err != nil {
			return nil, err
		}
		p.ImageURL = req.ImageURL
	}
	if req.CategoryID != "" {
		// Validate category exists.
		_, err := s.productCategoryRepo.FindByID(req.CategoryID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrProductCategoryNotFound
			}
			return nil, err
		}
		p.CategoryID = req.CategoryID
	}
	if req.Description != "" {
		p.Description = req.Description
	}

	if req.Status != "" {
		p.Status = req.Status
	}

	if err := s.productRepo.Update(p); err != nil {
		return nil, err
	}

	// Reload to get relations.
	p, err = s.productRepo.FindByID(p.ID)
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	if cache.Client != nil && cache.Client.IsEnabled() {
		// Clear all product list caches
		_ = cache.Client.DeletePattern("products:list:*")
		// Clear specific product cache if we implemented GetByID caching (we haven't yet, but good practice)
		// We should also invalidate GetByID cache if we add it. 
	}

	return p.ToProductResponse(), nil
}

// DeleteProduct deletes a product.
func (s *Service) DeleteProduct(id string) error {
	_, err := s.productRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProductNotFound
		}
		return err
	}

	err = s.productRepo.Delete(id)
	if err != nil {
		return err
	}

	// Invalidate cache
	if cache.Client != nil && cache.Client.IsEnabled() {
		_ = cache.Client.DeletePattern("products:list:*")
	}

	return nil
}

// ListProductCategories returns a list of product categories.
func (s *Service) ListProductCategories(req *product.ListProductCategoriesRequest) ([]product.ProductCategoryResponse, error) {
	categories, err := s.productCategoryRepo.List(req)
	if err != nil {
		return nil, err
	}

	responses := make([]product.ProductCategoryResponse, len(categories))
	for i, c := range categories {
		responses[i] = *c.ToProductCategoryResponse()
	}

	return responses, nil
}

// GetProductCategoryByID returns a product category by ID.
func (s *Service) GetProductCategoryByID(id string) (*product.ProductCategoryResponse, error) {
	c, err := s.productCategoryRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductCategoryNotFound
		}
		return nil, err
	}

	return c.ToProductCategoryResponse(), nil
}

// CreateProductCategory creates a new product category.
func (s *Service) CreateProductCategory(req *product.CreateProductCategoryRequest) (*product.ProductCategoryResponse, error) {
	status := req.Status
	if status == "" {
		status = "active"
	}

	c := &product.ProductCategory{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Status:      status,
	}

	if err := s.productCategoryRepo.Create(c); err != nil {
		return nil, err
	}

	// Reload to ensure we return latest state.
	c, err := s.productCategoryRepo.FindByID(c.ID)
	if err != nil {
		return nil, err
	}

	return c.ToProductCategoryResponse(), nil
}

// UpdateProductCategory updates an existing product category.
func (s *Service) UpdateProductCategory(id string, req *product.UpdateProductCategoryRequest) (*product.ProductCategoryResponse, error) {
	c, err := s.productCategoryRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductCategoryNotFound
		}
		return nil, err
	}

	if req.Name != "" {
		c.Name = req.Name
	}
	if req.Slug != "" {
		c.Slug = req.Slug
	}
	if req.Description != "" {
		c.Description = req.Description
	}
	if req.Status != "" {
		c.Status = req.Status
	}

	if err := s.productCategoryRepo.Update(c); err != nil {
		return nil, err
	}

	c, err = s.productCategoryRepo.FindByID(c.ID)
	if err != nil {
		return nil, err
	}

	return c.ToProductCategoryResponse(), nil
}

// DeleteProductCategory deletes a product category.
func (s *Service) DeleteProductCategory(id string) error {
	_, err := s.productCategoryRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProductCategoryNotFound
		}
		return err
	}

	return s.productCategoryRepo.Delete(id)
}


