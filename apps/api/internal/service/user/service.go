package user

import (
	"errors"
	"net/url"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"github.com/gilabs/crm-healthcare/api/pkg/cache"
	"github.com/gilabs/crm-healthcare/api/pkg/util/currency"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	userRepo          interfaces.UserRepository
	roleRepo          interfaces.RoleRepository
	groupRepo         interfaces.GroupRepository
	brickRepo         interfaces.BrickRepository
	monthlyTargetRepo interfaces.MonthlyTargetRepository
	cache             *cache.Cache
}

func NewService(userRepo interfaces.UserRepository, roleRepo interfaces.RoleRepository, groupRepo interfaces.GroupRepository, brickRepo interfaces.BrickRepository, monthlyTargetRepo interfaces.MonthlyTargetRepository, cache *cache.Cache) *Service {
	return &Service{
		userRepo:          userRepo,
		roleRepo:          roleRepo,
		groupRepo:         groupRepo,
		brickRepo:         brickRepo,
		monthlyTargetRepo: monthlyTargetRepo,
		cache:             cache,
	}
}

// List returns a list of users with pagination
func (s *Service) List(req *user.ListUsersRequest) ([]user.UserResponse, *PaginationResult, error) {
	users, total, err := s.userRepo.List(req)
	if err != nil {
		return nil, nil, err
	}

	// Get current year and month for monthly target lookup
	now := time.Now()
	currentYear := now.Year()
	currentMonth := int(now.Month())

	// OPTIMIZED: Batch load monthly targets for all users in one query (prevents N+1)
	userIDs := make([]string, len(users))
	for i, u := range users {
		userIDs[i] = u.ID
	}

	monthlyTargetsMap, _ := s.monthlyTargetRepo.BatchGetUserEffectiveTargets(userIDs, currentYear, currentMonth)

	// OPTIMIZED: Batch load bricks for all users in one query (prevents N+1)
	// Collect unique brick IDs from users
	brickIDSet := make(map[string]bool)
	for _, u := range users {
		if u.BrickID != nil && *u.BrickID != "" {
			brickIDSet[*u.BrickID] = true
		}
	}

	// Convert set to slice
	brickIDs := make([]string, 0, len(brickIDSet))
	for brickID := range brickIDSet {
		brickIDs = append(brickIDs, brickID)
	}

	// Batch load all bricks using WHERE IN query (prevents N+1)
	bricksMap := make(map[string]interface{})
	if len(brickIDs) > 0 {
		bricks, err := s.brickRepo.FindByIDs(brickIDs)
		if err == nil {
			for _, brickEntity := range bricks {
				bricksMap[brickEntity.ID] = brickEntity.ToBrickResponse()
			}
		} else {
			// Log error but don't fail the entire request
			// This allows users without bricks to still be returned
		}
	}

	responses := make([]user.UserResponse, len(users))
	for i, u := range users {
		userResp := *u.ToUserResponse()

		// Get brick from batch-loaded map
		if u.BrickID != nil && *u.BrickID != "" {
			if brickResp, exists := bricksMap[*u.BrickID]; exists {
				userResp.Brick = brickResp
			}
		}

		// Get effective monthly target from batch-loaded map
		if effectiveTarget, exists := monthlyTargetsMap[u.ID]; exists && effectiveTarget != nil {
			// Convert to simplified format to avoid circular dependency
			userResp.MonthlyTarget = &user.UserMonthlyTarget{
				ID:                    effectiveTarget.ID,
				GroupID:               effectiveTarget.GroupID,
				UserID:                effectiveTarget.UserID,
				Year:                  effectiveTarget.Year,
				Month:                 effectiveTarget.Month,
				TargetAmount:          effectiveTarget.TargetAmount,
				TargetAmountFormatted: formatCurrency(effectiveTarget.TargetAmount),
				CreatedAt:             effectiveTarget.CreatedAt,
				UpdatedAt:             effectiveTarget.UpdatedAt,
			}
		}
		// If not found, monthly_target will be nil (omitempty)

		responses[i] = userResp
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

	pagination := &PaginationResult{
		Page:       page,
		PerPage:    perPage,
		Total:      int(total),
		TotalPages: int((total + int64(perPage) - 1) / int64(perPage)),
	}

	return responses, pagination, nil
}

// GetByID returns a user by ID
func (s *Service) GetByID(id string) (*user.UserResponse, error) {
	// 1. Try Cache
	if s.cache != nil && s.cache.IsEnabled() {
		var cachedUser user.UserResponse
		found, err := s.cache.Get("user:"+id, &cachedUser)
		if err == nil && found {
			return &cachedUser, nil
		}
	}

	u, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	userResp := *u.ToUserResponse()

	// Load brick if brick_id exists
	if u.BrickID != nil && *u.BrickID != "" {
		brickEntity, err := s.brickRepo.FindByID(*u.BrickID)
		if err == nil && brickEntity != nil {
			userResp.Brick = brickEntity.ToBrickResponse()
		}
		// If error (e.g., not found), brick will be nil (omitempty)
	}

	// Get effective monthly target for current month/year (user target > group target)
	now := time.Now()
	currentYear := now.Year()
	currentMonth := int(now.Month())

	effectiveTarget, err := s.monthlyTargetRepo.GetUserEffectiveTarget(u.ID, currentYear, currentMonth)
	if err == nil && effectiveTarget != nil {
		// Convert to simplified format to avoid circular dependency
		userResp.MonthlyTarget = &user.UserMonthlyTarget{
			ID:                    effectiveTarget.ID,
			GroupID:               effectiveTarget.GroupID,
			UserID:                effectiveTarget.UserID,
			Year:                  effectiveTarget.Year,
			Month:                 effectiveTarget.Month,
			TargetAmount:          effectiveTarget.TargetAmount,
			TargetAmountFormatted: formatCurrency(effectiveTarget.TargetAmount),
			CreatedAt:             effectiveTarget.CreatedAt,
			UpdatedAt:             effectiveTarget.UpdatedAt,
		}
	}
	// If error (e.g., not found), monthly_target will be nil (omitempty)

	// If error (e.g., not found), monthly_target will be nil (omitempty)

	// 3. Set Cache
	if s.cache != nil && s.cache.IsEnabled() {
		// Use 5 minutes TTL as per standard
		_ = s.cache.Set("user:"+id, &userResp, 5*time.Minute)
	}

	return &userResp, nil
}

// Create creates a new user
func (s *Service) Create(req *user.CreateUserRequest) (*user.UserResponse, error) {
	// Check if role exists
	_, err := s.roleRepo.FindByID(req.RoleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}

	// Check if email already exists
	_, err = s.userRepo.FindByEmail(req.Email)
	if err == nil {
		return nil, ErrUserAlreadyExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "active"
	}

	// Generate avatar URL using dicebear lorelei
	avatarURL := "https://api.dicebear.com/7.x/lorelei/svg?seed=" + url.QueryEscape(req.Email)

	// Check if group exists (if provided)
	if req.GroupID != nil && *req.GroupID != "" {
		_, err = s.groupRepo.FindByID(*req.GroupID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrGroupNotFound
			}
			return nil, err
		}
	}

	// Check if brick exists (if provided)
	if req.BrickID != nil && *req.BrickID != "" {
		_, err = s.brickRepo.FindByID(*req.BrickID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrBrickNotFound
			}
			return nil, err
		}
	}

	// Create user
	u := &user.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		Name:      req.Name,
		AvatarURL: avatarURL,
		RoleID:    req.RoleID,
		GroupID:   req.GroupID,
		BrickID:   req.BrickID,
		Status:    status,
	}

	if err := s.userRepo.Create(u); err != nil {
		return nil, err
	}

	// Reload with role
	createdUser, err := s.userRepo.FindByID(u.ID)
	if err != nil {
		return nil, err
	}

	return createdUser.ToUserResponse(), nil
}

// Update updates a user
func (s *Service) Update(id string, req *user.UpdateUserRequest) (*user.UserResponse, error) {
	// Find user
	u, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Update fields
	if req.Email != "" {
		// Check if email already exists (excluding current user)
		existingUser, err := s.userRepo.FindByEmail(req.Email)
		if err == nil && existingUser.ID != id {
			return nil, ErrUserAlreadyExists
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		u.Email = req.Email
	}

	if req.Name != "" {
		u.Name = req.Name
	}

	if req.RoleID != "" {
		// Check if role exists
		_, err := s.roleRepo.FindByID(req.RoleID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrRoleNotFound
			}
			return nil, err
		}
		u.RoleID = req.RoleID
	}

	if req.Status != "" {
		u.Status = req.Status
	}

	if req.GroupID != nil {
		if *req.GroupID != "" {
			// Check if group exists
			_, err := s.groupRepo.FindByID(*req.GroupID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, ErrGroupNotFound
				}
				return nil, err
			}
			u.GroupID = req.GroupID
		} else {
			// Set to nil if empty string
			u.GroupID = nil
		}
	}

	if req.BrickID != nil {
		if *req.BrickID != "" {
			// Check if brick exists
			_, err := s.brickRepo.FindByID(*req.BrickID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, ErrBrickNotFound
				}
				return nil, err
			}
			u.BrickID = req.BrickID
		} else {
			// Set to nil if empty string
			u.BrickID = nil
		}
	}

	if err := s.userRepo.Update(u); err != nil {
		return nil, err
	}

	// Reload with role
	updatedUser, err := s.userRepo.FindByID(u.ID)
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	if s.cache != nil && s.cache.IsEnabled() {
		_ = s.cache.Delete("user:" + id)
	}

	return updatedUser.ToUserResponse(), nil
}

// Delete deletes a user
func (s *Service) Delete(id string) error {
	// Check if user exists
	_, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	if err := s.userRepo.Delete(id); err != nil {
		return err
	}

	// Invalidate cache
	if s.cache != nil && s.cache.IsEnabled() {
		_ = s.cache.Delete("user:" + id)
	}

	return nil
}

// PaginationResult represents pagination result
type PaginationResult struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// formatCurrency formats integer (sen) to formatted currency string
func formatCurrency(amount int64) string {
	return currency.FormatCurrency(amount)
}
