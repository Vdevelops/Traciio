package account

import (
	"errors"
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/config"
	"github.com/gilabs/crm-healthcare/api/internal/domain/account"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	brick "github.com/gilabs/crm-healthcare/api/internal/service/brick"
	"github.com/gilabs/crm-healthcare/api/pkg/cache"
	"github.com/gilabs/crm-healthcare/api/pkg/geocoding"
	"gorm.io/gorm"
)

var (
	ErrAccountNotFound  = errors.New("account not found")
	ErrCategoryNotFound = errors.New("category not found")
)

type Service struct {
	accountRepo      interfaces.AccountRepository
	categoryRepo     interfaces.CategoryRepository
	brickHelper      *brick.BrickHelper
	geocodingSvc     *geocoding.GeocodingService
	geocodingEnabled bool
	cacheService     *cache.AccountCacheService
}

func NewService(accountRepo interfaces.AccountRepository, categoryRepo interfaces.CategoryRepository, brickHelper *brick.BrickHelper) *Service {
	var geocodingSvc *geocoding.GeocodingService
	geocodingEnabled := false

	if config.AppConfig != nil && config.AppConfig.Geocoding.Enabled {
		geocodingSvc = geocoding.NewGeocodingService(
			config.AppConfig.Geocoding.Provider,
			config.AppConfig.Geocoding.APIKey,
		)
		geocodingEnabled = true
	}

	return &Service{
		accountRepo:      accountRepo,
		categoryRepo:     categoryRepo,
		brickHelper:      brickHelper,
		geocodingSvc:     geocodingSvc,
		geocodingEnabled: geocodingEnabled,
		cacheService:     cache.NewAccountCacheService(nil),
	}
}

// PaginationResult represents pagination information
type PaginationResult struct {
	Page       int
	PerPage    int
	Total      int
	TotalPages int
}

// cachedListResult represents cached list data
type cachedListResult struct {
	Accounts []account.AccountResponse `msgpack:"accounts"`
	Total    int64                     `msgpack:"total"`
}

// List returns a list of accounts with pagination
func (s *Service) List(req *account.ListAccountsRequest) ([]account.AccountResponse, *PaginationResult, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 1000 {
		perPage = 1000
	}

	// Build filter map for cache key
	filters := map[string]interface{}{
		"search":      req.Search,
		"status":      req.Status,
		"category_id": req.CategoryID,
		"assigned_to": req.AssignedTo,
		"brick_id":    req.BrickID,
	}

	// Try to get from cache first
	var cachedResult cachedListResult
	if found, _ := s.cacheService.GetList(page, perPage, filters, &cachedResult); found {
		totalPages := int((cachedResult.Total + int64(perPage) - 1) / int64(perPage))
		pagination := &PaginationResult{
			Page:       page,
			PerPage:    perPage,
			Total:      int(cachedResult.Total),
			TotalPages: totalPages,
		}
		return cachedResult.Accounts, pagination, nil
	}

	accounts, total, err := s.accountRepo.List(req)
	if err != nil {
		return nil, nil, err
	}

	responses := make([]account.AccountResponse, len(accounts))
	for i, a := range accounts {
		responses[i] = *a.ToAccountResponse()
	}

	pagination := &PaginationResult{
		Page:       page,
		PerPage:    perPage,
		Total:      int(total),
		TotalPages: int((total + int64(perPage) - 1) / int64(perPage)),
	}

	// Cache the result
	_ = s.cacheService.SetList(page, perPage, filters, cachedListResult{
		Accounts: responses,
		Total:    total,
	})

	return responses, pagination, nil
}

// GetByID returns an account by ID
func (s *Service) GetByID(id string) (*account.AccountResponse, error) {
	// Try to get from cache first
	var cachedAccount account.AccountResponse
	if found, _ := s.cacheService.GetDetail(id, &cachedAccount); found {
		return &cachedAccount, nil
	}

	a, err := s.accountRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}

	response := a.ToAccountResponse()

	// Cache the result
	_ = s.cacheService.SetDetail(id, response)

	return response, nil
}

// Create creates a new account
func (s *Service) Create(req *account.CreateAccountRequest) (*account.AccountResponse, error) {
	// Validate category exists
	_, err := s.categoryRepo.FindByID(req.CategoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	var assignedTo *string
	if req.AssignedTo != "" {
		assignedTo = &req.AssignedTo
	}

	// Auto-populate brick_id if not provided
	var brickID *string
	if req.BrickID != "" {
		brickID = &req.BrickID
	} else if s.brickHelper != nil {
		// Try to get brick_id from assigned_to user
		if assignedTo != nil {
			brickID, _ = s.brickHelper.GetBrickIDFromUser(*assignedTo)
		}
		// If still nil, try to get from location (province + city/regency)
		if brickID == nil && req.Province != "" && req.City != "" {
			brickID, _ = s.brickHelper.GetBrickIDFromLocation(req.Province, req.City)
		}
	}

	a := &account.Account{
		Name:       req.Name,
		CategoryID: req.CategoryID,
		Address:    req.Address,
		City:       req.City,
		Province:   req.Province,
		Phone:      req.Phone,
		Email:      req.Email,
		Latitude:   req.Latitude,
		Longitude:  req.Longitude,
		PostalCode: req.PostalCode,
		Country:    req.Country,
		Website:    req.Website,
		Industry:   req.Industry,
		AssignedTo: assignedTo,
		BrickID:    brickID,
	}

	// Set default country if not provided
	if a.Country == "" {
		a.Country = "Indonesia"
	}

	if req.Status != "" {
		a.Status = req.Status
	} else {
		a.Status = "active"
	}

	// Geocode address if enabled, ONLY when coords are not provided
	if (a.Latitude == nil || a.Longitude == nil) && s.geocodingEnabled && s.geocodingSvc != nil {
		if req.Address != "" || req.City != "" || req.Province != "" {
			result, err := s.geocodingSvc.GeocodeAddressWithFallback(req.Address, req.City, req.Province)
			if err != nil {
				// Log error but don't fail the account creation
				log.Printf("Warning: Failed to geocode address for account %s: %v", req.Name, err)
			} else {
				a.Latitude = &result.Latitude
				a.Longitude = &result.Longitude
			}
		}
	}

	if err := s.accountRepo.Create(a); err != nil {
		return nil, err
	}

	// Invalidate cache after create
	_ = s.cacheService.InvalidateOnWrite("")

	// Reload with category
	createdAccount, err := s.accountRepo.FindByID(a.ID)
	if err != nil {
		return nil, err
	}

	return createdAccount.ToAccountResponse(), nil
}

// Update updates an account
func (s *Service) Update(id string, req *account.UpdateAccountRequest) (*account.AccountResponse, error) {
	a, err := s.accountRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}

	// Track if fields that affect brick_id are being updated
	assignedToChanged := false
	locationChanged := false

	// Update fields if provided
	if req.Name != "" {
		a.Name = req.Name
	}
	if req.CategoryID != "" {
		// Validate category exists
		_, err := s.categoryRepo.FindByID(req.CategoryID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrCategoryNotFound
			}
			return nil, err
		}
		a.CategoryID = req.CategoryID
	}
	addressChanged := false
	if req.Address != "" {
		if a.Address != req.Address {
			addressChanged = true
		}
		a.Address = req.Address
	}
	if req.City != "" {
		if a.City != req.City {
			locationChanged = true
			addressChanged = true
		}
		a.City = req.City
	}
	if req.Province != "" {
		if a.Province != req.Province {
			locationChanged = true
			addressChanged = true
		}
		a.Province = req.Province
	}

	// Re-geocode if address changed
	if addressChanged && (req.Latitude == nil || req.Longitude == nil) && s.geocodingEnabled && s.geocodingSvc != nil {
		result, err := s.geocodingSvc.GeocodeAddressWithFallback(a.Address, a.City, a.Province)
		if err != nil {
			// Log error but don't fail the update
			log.Printf("Warning: Failed to geocode address for account %s: %v", a.ID, err)
		} else {
			a.Latitude = &result.Latitude
			a.Longitude = &result.Longitude
		}
	}

	// Manual coordinates take precedence when provided
	if req.Latitude != nil {
		a.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		a.Longitude = req.Longitude
	}
	if req.Phone != "" {
		a.Phone = req.Phone
	}
	if req.Email != "" {
		a.Email = req.Email
	}
	if req.Status != "" {
		a.Status = req.Status
	}
	if req.AssignedTo != "" {
		assignedTo := req.AssignedTo
		// Check if assigned_to is actually changing
		if a.AssignedTo == nil || *a.AssignedTo != assignedTo {
			assignedToChanged = true
		}
		a.AssignedTo = &assignedTo
	} else if req.AssignedTo == "" && a.AssignedTo != nil {
		// Explicitly set to nil (clearing assignment)
		assignedToChanged = true
		a.AssignedTo = nil
	}
	if req.PostalCode != "" {
		a.PostalCode = req.PostalCode
	}
	if req.Country != "" {
		a.Country = req.Country
	}
	if req.Website != "" {
		a.Website = req.Website
	}
	if req.Industry != "" {
		a.Industry = req.Industry
	}

	// Auto-update brick_id if assigned_to or location changed
	if (assignedToChanged || locationChanged) && s.brickHelper != nil {
		var brickID *string
		// Try to get brick_id from assigned_to user first
		if a.AssignedTo != nil && *a.AssignedTo != "" {
			brickID, _ = s.brickHelper.GetBrickIDFromUser(*a.AssignedTo)
		}
		// If still nil, try to get from location (province + city/regency)
		if brickID == nil && a.Province != "" && a.City != "" {
			brickID, _ = s.brickHelper.GetBrickIDFromLocation(a.Province, a.City)
		}
		a.BrickID = brickID
	}

	if err := s.accountRepo.Update(a); err != nil {
		return nil, err
	}

	// Invalidate cache after update
	_ = s.cacheService.InvalidateOnWrite(id)

	// Reload with category
	updatedAccount, err := s.accountRepo.FindByID(a.ID)
	if err != nil {
		return nil, err
	}

	return updatedAccount.ToAccountResponse(), nil
}

// Delete deletes an account
func (s *Service) Delete(id string) error {
	_, err := s.accountRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAccountNotFound
		}
		return err
	}

	err = s.accountRepo.Delete(id)
	if err != nil {
		return err
	}

	// Invalidate cache after delete
	_ = s.cacheService.InvalidateOnWrite(id)

	return nil
}

// ListAllForMap returns all accounts for map display (read-only; no geocoding/backfill)
func (s *Service) ListAllForMap(status string) ([]account.AccountResponse, error) {
	// Fetch all accounts
	accounts, err := s.accountRepo.ListAll(status)
	if err != nil {
		return nil, err
	}

	// Convert to response
	responses := make([]account.AccountResponse, len(accounts))
	for i, a := range accounts {
		responses[i] = *a.ToAccountResponse()
	}

	return responses, nil
}

// ListByBBox returns accounts within a geographic bounding box for viewport-based map loading
func (s *Service) ListByBBox(req *account.BBoxRequest) ([]account.AccountResponse, error) {
	accounts, err := s.accountRepo.ListByBBox(req)
	if err != nil {
		return nil, err
	}

	responses := make([]account.AccountResponse, len(accounts))
	for i, a := range accounts {
		responses[i] = *a.ToAccountResponse()
	}

	return responses, nil
}
