package brick

import (
	"errors"
	"fmt"
	"strings"
	"time"

	brickdomain "github.com/gilabs/crm-healthcare/api/internal/domain/brick"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	cachepkg "github.com/gilabs/crm-healthcare/api/pkg/cache"
	"github.com/gilabs/crm-healthcare/api/pkg/util/currency"
	"gorm.io/gorm"
)

var (
	ErrBrickNotFound       = errors.New("brick not found")
	ErrBrickAlreadyExists  = errors.New("brick already exists")
	ErrBrickInUse          = errors.New("brick is in use and cannot be deleted")
	ErrInvalidManager      = errors.New("manager must be a user with sales_manager role")
	ErrRegencyAlreadyBrick = errors.New("regency already has a brick")
	ErrMinimumSalesRequired = errors.New("brick must have at least 2 sales reps assigned")
)

type Service struct {
	brickRepo         interfaces.BrickRepository
	userRepo          interfaces.UserRepository
	monthlyTargetRepo interfaces.MonthlyTargetRepository
	db                *gorm.DB
	cacheService      *cachepkg.BrickCacheService
}

func NewService(brickRepo interfaces.BrickRepository, userRepo interfaces.UserRepository, monthlyTargetRepo interfaces.MonthlyTargetRepository, db *gorm.DB) *Service {
	return &Service{
		brickRepo:         brickRepo,
		userRepo:          userRepo,
		monthlyTargetRepo: monthlyTargetRepo,
		db:                db,
		cacheService:      cachepkg.NewBrickCacheService(nil),
	}
}

type cachedBrickListResult struct {
	Bricks     []brickdomain.BrickResponse
	Pagination *PaginationResult
}

// GetBrickRepo returns the brick repository (for authorization checks)
func (s *Service) GetBrickRepo() interfaces.BrickRepository {
	return s.brickRepo
}

// PaginationResult represents pagination result
type PaginationResult struct {
	Page       int
	PerPage    int
	Total      int
	TotalPages int
}

// List returns a list of bricks with pagination
func (s *Service) List(req *brickdomain.ListBricksRequest) ([]brickdomain.BrickResponse, *PaginationResult, error) {
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

	filterMap := map[string]interface{}{}
	if req.Search != "" {
		filterMap["search"] = req.Search
	}
	if req.Province != "" {
		filterMap["province"] = req.Province
	}
	if req.Regency != "" {
		filterMap["regency"] = req.Regency
	}
	if req.Status != "" {
		filterMap["status"] = req.Status
	}
	if req.ManagerID != nil && *req.ManagerID != "" {
		filterMap["manager_id"] = *req.ManagerID
	}

	var cachedResult cachedBrickListResult
	if found, _ := s.cacheService.GetList(page, perPage, filterMap, &cachedResult); found && cachedResult.Pagination != nil {
		return cachedResult.Bricks, cachedResult.Pagination, nil
	}

	bricks, total, err := s.brickRepo.List(req)
	if err != nil {
		return nil, nil, err
	}

	responses := make([]brickdomain.BrickResponse, len(bricks))
	for i, b := range bricks {
		resp := *b.ToBrickResponse()
		// Get sales count
		salesCount, err := s.brickRepo.CountSalesByBrickID(b.ID)
		if err == nil {
			resp.SalesCount = int(salesCount)
		}
		responses[i] = resp
	}

	pagination := &PaginationResult{
		Page:       page,
		PerPage:    perPage,
		Total:      int(total),
		TotalPages: int((total + int64(perPage) - 1) / int64(perPage)),
	}

	_ = s.cacheService.SetList(page, perPage, filterMap, cachedBrickListResult{
		Bricks:     responses,
		Pagination: pagination,
	})

	return responses, pagination, nil
}

// GetByID returns a brick by ID
func (s *Service) GetByID(id string) (*brickdomain.BrickResponse, error) {
	var cachedResponse brickdomain.BrickResponse
	if found, _ := s.cacheService.GetDetail(id, &cachedResponse); found && cachedResponse.ID != "" {
		return &cachedResponse, nil
	}

	b, err := s.brickRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBrickNotFound
		}
		return nil, err
	}

	resp := b.ToBrickResponse()
	// Get sales count
	salesCount, err := s.brickRepo.CountSalesByBrickID(id)
	if err == nil {
		resp.SalesCount = int(salesCount)
	}

	_ = s.cacheService.SetDetail(id, resp)
	return resp, nil
}

// Create creates a new brick
// If Code is empty, auto-generates a unique code based on province abbreviation
func (s *Service) Create(req *brickdomain.CreateBrickRequest) (*brickdomain.BrickResponse, error) {
	// Auto-generate code if not provided
	if req.Code == "" {
		generatedCode, err := s.generateBrickCode(req.Province, req.Regency)
		if err != nil {
			return nil, fmt.Errorf("failed to generate brick code: %w", err)
		}
		req.Code = generatedCode
	}

	// Check if code already exists
	_, err := s.brickRepo.FindByCode(req.Code)
	if err == nil {
		return nil, ErrBrickAlreadyExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Check if regency already has a brick (optional validation - can be removed if multiple bricks per regency allowed)
	_, err = s.brickRepo.FindByRegencyAndProvince(req.Regency, req.Province)
	if err == nil {
		return nil, ErrRegencyAlreadyBrick
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Validate manager if provided
	if req.ManagerID != nil && *req.ManagerID != "" {
		manager, err := s.userRepo.FindByID(*req.ManagerID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("manager not found")
			}
			return nil, err
		}
		// Check if user has sales_manager role (you may need to check role name or add role check)
		// For now, we'll just check if user exists. Role validation can be added in handler layer.
		_ = manager // Use manager to avoid unused variable
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "active"
	}

	// Create brick
	b := &brickdomain.Brick{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Province:    req.Province,
		Regency:     req.Regency,
		District:    req.District,
		ManagerID:   req.ManagerID,
		Status:      status,
	}

	if err := s.brickRepo.Create(b); err != nil {
		return nil, err
	}

	// Reload to get timestamps and relations
	createdBrick, err := s.brickRepo.FindByID(b.ID)
	if err != nil {
		return nil, err
	}

	resp := createdBrick.ToBrickResponse()
	salesCount, err := s.brickRepo.CountSalesByBrickID(b.ID)
	if err == nil {
		resp.SalesCount = int(salesCount)
	}

	_ = s.cacheService.InvalidateOnWrite(b.ID)
	_ = s.cacheService.SetDetail(b.ID, resp)

	return resp, nil
}

// Update updates a brick
func (s *Service) Update(id string, req *brickdomain.UpdateBrickRequest) (*brickdomain.BrickResponse, error) {
	// Find existing brick
	b, err := s.brickRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBrickNotFound
		}
		return nil, err
	}

	// Validate manager if being updated
	if req.ManagerID != nil && *req.ManagerID != "" {
		manager, err := s.userRepo.FindByID(*req.ManagerID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("manager not found")
			}
			return nil, err
		}
		_ = manager // Use manager to avoid unused variable
	}

	// Update fields
	if req.Name != nil {
		b.Name = *req.Name
	}
	if req.Code != nil {
		// Check if new code already exists (and is not the current brick's code)
		existingBrick, err := s.brickRepo.FindByCode(*req.Code)
		if err == nil && existingBrick.ID != id {
			return nil, ErrBrickAlreadyExists
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		b.Code = *req.Code
	}
	if req.Description != nil {
		b.Description = *req.Description
	}
	if req.District != nil {
		b.District = req.District
	}
	if req.ManagerID != nil {
		b.ManagerID = req.ManagerID
	}
	if req.Status != nil {
		b.Status = *req.Status
	}

	if err := s.brickRepo.Update(b); err != nil {
		return nil, err
	}

	// Reload to get updated timestamps and relations
	updatedBrick, err := s.brickRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	resp := updatedBrick.ToBrickResponse()
	salesCount, err := s.brickRepo.CountSalesByBrickID(id)
	if err == nil {
		resp.SalesCount = int(salesCount)
	}

	_ = s.cacheService.InvalidateOnWrite(id)
	_ = s.cacheService.SetDetail(id, resp)

	return resp, nil
}

// Delete deletes a brick (soft delete)
func (s *Service) Delete(id string) error {
	// Check if brick exists
	_, err := s.brickRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrBrickNotFound
		}
		return fmt.Errorf("failed to find brick: %w", err)
	}

	// Start transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	committed := false
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	// Soft delete brick_target_distributions
	if err := tx.Table("brick_target_distributions").
		Where("brick_id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP")).Error; err != nil {
		return fmt.Errorf("failed to soft delete target distributions: %w", err)
	}

	// Soft delete monthly_targets
	if err := tx.Table("monthly_targets").
		Where("brick_id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP")).Error; err != nil {
		return fmt.Errorf("failed to soft delete monthly targets: %w", err)
	}

	// Unassign/nullify brick_id from all related tables
	// This ensures data integrity while allowing brick deletion

	// Unassign users from brick (set brick_id to NULL)
	if err := tx.Table("users").
		Where("brick_id = ? AND deleted_at IS NULL", id).
		Update("brick_id", nil).Error; err != nil {
		return fmt.Errorf("failed to unassign users from brick: %w", err)
	}

	// Unassign accounts from brick (set brick_id to NULL)
	if err := tx.Table("accounts").
		Where("brick_id = ? AND deleted_at IS NULL", id).
		Update("brick_id", nil).Error; err != nil {
		return fmt.Errorf("failed to unassign accounts from brick: %w", err)
	}

	// Unassign deals from brick (set brick_id to NULL)
	if err := tx.Table("deals").
		Where("brick_id = ? AND deleted_at IS NULL", id).
		Update("brick_id", nil).Error; err != nil {
		return fmt.Errorf("failed to unassign deals from brick: %w", err)
	}

	// Unassign visit_reports from brick (set brick_id to NULL)
	if err := tx.Table("visit_reports").
		Where("brick_id = ? AND deleted_at IS NULL", id).
		Update("brick_id", nil).Error; err != nil {
		return fmt.Errorf("failed to unassign visit reports from brick: %w", err)
	}

	// Soft delete the brick
	result := tx.Exec("UPDATE bricks SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete brick: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrBrickNotFound
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	committed = true

	_ = s.cacheService.InvalidateOnWrite(id)

	return nil
}

// GetSalesByBrickID returns all sales users in a brick with monthly target data
func (s *Service) GetSalesByBrickID(brickID string) ([]user.UserResponse, error) {
	// Check if brick exists
	_, err := s.brickRepo.FindByID(brickID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBrickNotFound
		}
		return nil, err
	}

	users, err := s.brickRepo.GetSalesByBrickID(brickID)
	if err != nil {
		return nil, err
	}

	// Batch-load monthly targets to avoid N+1 queries
	now := time.Now()
	currentYear := now.Year()
	currentMonth := int(now.Month())

	userIDs := make([]string, len(users))
	for i, u := range users {
		userIDs[i] = u.ID
	}

	monthlyTargetsMap, _ := s.monthlyTargetRepo.BatchGetUserEffectiveTargets(userIDs, currentYear, currentMonth)

	responses := make([]user.UserResponse, len(users))
	for i, u := range users {
		resp := *u.ToUserResponse()
		if effectiveTarget, exists := monthlyTargetsMap[u.ID]; exists && effectiveTarget != nil {
			resp.MonthlyTarget = &user.UserMonthlyTarget{
				ID:                    effectiveTarget.ID,
				GroupID:               effectiveTarget.GroupID,
				UserID:                effectiveTarget.UserID,
				Year:                  effectiveTarget.Year,
				Month:                 effectiveTarget.Month,
				TargetAmount:          effectiveTarget.TargetAmount,
				TargetAmountFormatted: currency.FormatCurrency(effectiveTarget.TargetAmount),
				CreatedAt:             effectiveTarget.CreatedAt,
				UpdatedAt:             effectiveTarget.UpdatedAt,
			}
		}
		responses[i] = resp
	}

	return responses, nil
}

// AssignSales assigns sales users to a brick by setting their brick_id
func (s *Service) AssignSales(brickID string, userIDs []string) ([]user.UserResponse, error) {
	// Check if brick exists
	b, err := s.brickRepo.FindByID(brickID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBrickNotFound
		}
		return nil, err
	}
	_ = b

	// Validate all users exist and have the sales role
	for _, uid := range userIDs {
		u, uErr := s.userRepo.FindByID(uid)
		if uErr != nil {
			if errors.Is(uErr, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("user %s not found", uid)
			}
			return nil, uErr
		}
		if u.Role != nil && u.Role.Code != "sales" {
			return nil, fmt.Errorf("user %s (%s) is not a sales rep", uid, u.Name)
		}
	}

	// Assign users to brick
	if err := s.db.Table("users").
		Where("id IN ? AND deleted_at IS NULL", userIDs).
		Update("brick_id", brickID).Error; err != nil {
		return nil, fmt.Errorf("failed to assign users to brick: %w", err)
	}

	_ = s.cacheService.InvalidateOnWrite(brickID)

	return s.GetSalesByBrickID(brickID)
}

// UnassignSales removes sales users from a brick by setting their brick_id to NULL
func (s *Service) UnassignSales(brickID string, userIDs []string) ([]user.UserResponse, error) {
	// Check if brick exists
	_, err := s.brickRepo.FindByID(brickID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBrickNotFound
		}
		return nil, err
	}

	// Unassign users from brick
	if err := s.db.Table("users").
		Where("id IN ? AND brick_id = ? AND deleted_at IS NULL", userIDs, brickID).
		Update("brick_id", nil).Error; err != nil {
		return nil, fmt.Errorf("failed to unassign users from brick: %w", err)
	}

	_ = s.cacheService.InvalidateOnWrite(brickID)

	return s.GetSalesByBrickID(brickID)
}

// generateBrickCode creates an auto-generated unique code from province and regency.
// Format: BRK-{PROVINCE_ABBREV}-{SEQ} e.g. BRK-JKT-001, BRK-JBR-002
func (s *Service) generateBrickCode(province, regency string) (string, error) {
	prefix := buildCodePrefix(province)

	seq, err := s.brickRepo.GetNextCodeSequence(prefix)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s-%03d", prefix, seq), nil
}

// buildCodePrefix creates a deterministic abbreviation prefix from province name.
// Uses first 3 consonants or characters of the province, uppercased.
func buildCodePrefix(province string) string {
	cleaned := strings.ToUpper(strings.TrimSpace(province))
	// Remove common prefixes
	cleaned = strings.TrimPrefix(cleaned, "PROVINSI ")
	cleaned = strings.TrimPrefix(cleaned, "DKI ")
	cleaned = strings.TrimPrefix(cleaned, "DI ")
	cleaned = strings.ReplaceAll(cleaned, " ", "")

	// Take first 3 characters as abbreviation
	abbrev := cleaned
	if len(abbrev) > 3 {
		abbrev = abbrev[:3]
	}
	if abbrev == "" {
		abbrev = "BRK"
	}

	return fmt.Sprintf("BRK-%s", abbrev)
}
