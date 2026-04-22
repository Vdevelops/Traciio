package monthly_target

import (
	"errors"

	monthlytargetdomain "github.com/gilabs/crm-healthcare/api/internal/domain/monthly_target"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

var (
	ErrMonthlyTargetNotFound      = errors.New("monthly target not found")
	ErrMonthlyTargetAlreadyExists = errors.New("monthly target already exists")
	ErrInvalidTargetScope         = errors.New("either group_id, user_id, or brick_id must be provided, but only one")
)

type Service struct {
	monthlyTargetRepo interfaces.MonthlyTargetRepository
	groupRepo         interfaces.GroupRepository
	userRepo          interfaces.UserRepository
	brickRepo         interfaces.BrickRepository
}

func NewService(monthlyTargetRepo interfaces.MonthlyTargetRepository, groupRepo interfaces.GroupRepository, userRepo interfaces.UserRepository, brickRepo interfaces.BrickRepository) *Service {
	return &Service{
		monthlyTargetRepo: monthlyTargetRepo,
		groupRepo:         groupRepo,
		userRepo:          userRepo,
		brickRepo:         brickRepo,
	}
}

// PaginationResult represents pagination result
type PaginationResult struct {
	Page       int
	PerPage    int
	Total      int
	TotalPages int
	TotalAmount int64
}

// List returns a list of monthly targets with pagination
func (s *Service) List(req *monthlytargetdomain.ListMonthlyTargetsRequest) ([]monthlytargetdomain.MonthlyTargetResponse, *PaginationResult, error) {
	targets, total, totalAmount, err := s.monthlyTargetRepo.List(req)
	if err != nil {
		return nil, nil, err
	}

	responses := make([]monthlytargetdomain.MonthlyTargetResponse, len(targets))
	for i, t := range targets {
		responses[i] = *t.ToMonthlyTargetResponse()
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
		TotalAmount: totalAmount,
	}

	return responses, pagination, nil
}

// GetByID returns a monthly target by ID
func (s *Service) GetByID(id string) (*monthlytargetdomain.MonthlyTargetResponse, error) {
	t, err := s.monthlyTargetRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMonthlyTargetNotFound
		}
		return nil, err
	}
	return t.ToMonthlyTargetResponse(), nil
}

// Create creates a new monthly target
func (s *Service) Create(req *monthlytargetdomain.CreateMonthlyTargetRequest) (*monthlytargetdomain.MonthlyTargetResponse, error) {
	// Validate: either group_id, user_id, or brick_id must be provided, but only one
	scopeCount := 0
	if req.GroupID != nil && *req.GroupID != "" {
		scopeCount++
	}
	if req.UserID != nil && *req.UserID != "" {
		scopeCount++
	}
	if req.BrickID != nil && *req.BrickID != "" {
		scopeCount++
	}
	
	if scopeCount == 0 {
		return nil, ErrInvalidTargetScope
	}
	if scopeCount > 1 {
		return nil, ErrInvalidTargetScope
	}

	// Validate group exists if provided
	if req.GroupID != nil && *req.GroupID != "" {
		_, err := s.groupRepo.FindByID(*req.GroupID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("group not found")
			}
			return nil, err
		}

		// Check if group target already exists for this period
		_, err = s.monthlyTargetRepo.FindByGroupAndPeriod(*req.GroupID, req.Year, req.Month)
		if err == nil {
			return nil, ErrMonthlyTargetAlreadyExists
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	// Validate user exists if provided
	if req.UserID != nil && *req.UserID != "" {
		_, err := s.userRepo.FindByID(*req.UserID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("user not found")
			}
			return nil, err
		}

		// Check if user target already exists for this period
		_, err = s.monthlyTargetRepo.FindByUserAndPeriod(*req.UserID, req.Year, req.Month)
		if err == nil {
			return nil, ErrMonthlyTargetAlreadyExists
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	// Validate brick exists if provided
	if req.BrickID != nil && *req.BrickID != "" {
		_, err := s.brickRepo.FindByID(*req.BrickID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("brick not found")
			}
			return nil, err
		}

		// Check if brick target already exists for this period
		_, err = s.monthlyTargetRepo.FindByBrickAndPeriod(*req.BrickID, req.Year, req.Month)
		if err == nil {
			return nil, ErrMonthlyTargetAlreadyExists
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	// Create monthly target
	t := &monthlytargetdomain.MonthlyTarget{
		GroupID:   req.GroupID,
		UserID:       req.UserID,
		BrickID:      req.BrickID,
		Year:         req.Year,
		Month:        req.Month,
		TargetAmount: req.TargetAmount,
	}

	if err := s.monthlyTargetRepo.Create(t); err != nil {
		return nil, err
	}

	// Reload to get timestamps and relations
	createdTarget, err := s.monthlyTargetRepo.FindByID(t.ID)
	if err != nil {
		return nil, err
	}

	return createdTarget.ToMonthlyTargetResponse(), nil
}

// Update updates a monthly target
func (s *Service) Update(id string, req *monthlytargetdomain.UpdateMonthlyTargetRequest) (*monthlytargetdomain.MonthlyTargetResponse, error) {
	// Find existing target
	t, err := s.monthlyTargetRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMonthlyTargetNotFound
		}
		return nil, err
	}

	// Update fields
	if req.Year != nil {
		t.Year = *req.Year
	}
	if req.Month != nil {
		t.Month = *req.Month
	}
	if req.TargetAmount != nil {
		t.TargetAmount = *req.TargetAmount
	}

	// Check for duplicate if year/month changed
	if req.Year != nil || req.Month != nil {
		var existing *monthlytargetdomain.MonthlyTarget
		if t.GroupID != nil {
			existing, err = s.monthlyTargetRepo.FindByGroupAndPeriod(*t.GroupID, t.Year, t.Month)
		} else if t.UserID != nil {
			existing, err = s.monthlyTargetRepo.FindByUserAndPeriod(*t.UserID, t.Year, t.Month)
		} else if t.BrickID != nil {
			existing, err = s.monthlyTargetRepo.FindByBrickAndPeriod(*t.BrickID, t.Year, t.Month)
		}

		if err == nil && existing.ID != id {
			return nil, ErrMonthlyTargetAlreadyExists
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	if err := s.monthlyTargetRepo.Update(t); err != nil {
		return nil, err
	}

	// Reload to get updated timestamps
	updatedTarget, err := s.monthlyTargetRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return updatedTarget.ToMonthlyTargetResponse(), nil
}

// Delete deletes a monthly target
func (s *Service) Delete(id string) error {
	// Check if target exists
	_, err := s.monthlyTargetRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrMonthlyTargetNotFound
		}
		return err
	}

	// Delete target
	return s.monthlyTargetRepo.Delete(id)
}

// GetUserEffectiveTarget gets effective target for user (user target or group target fallback)
func (s *Service) GetUserEffectiveTarget(userID string, year int, month int) (*monthlytargetdomain.MonthlyTargetResponse, error) {
	target, err := s.monthlyTargetRepo.GetUserEffectiveTarget(userID, year, month)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMonthlyTargetNotFound
		}
		return nil, err
	}
	return target.ToMonthlyTargetResponse(), nil
}

// BulkCreate creates multiple monthly targets in bulk
func (s *Service) BulkCreate(req *monthlytargetdomain.BulkCreateMonthlyTargetRequest) ([]monthlytargetdomain.MonthlyTargetResponse, error) {
	// Validate: either group_ids, user_ids, or brick_ids must be provided, but only one
	scopeCount := 0
	if len(req.GroupIDs) > 0 {
		scopeCount++
	}
	if len(req.UserIDs) > 0 {
		scopeCount++
	}
	if len(req.BrickIDs) > 0 {
		scopeCount++
	}

	if scopeCount == 0 {
		return nil, ErrInvalidTargetScope
	}
	if scopeCount > 1 {
		return nil, ErrInvalidTargetScope
	}

	var createdTargets []monthlytargetdomain.MonthlyTargetResponse

	// Handle bulk create for groups
	if len(req.GroupIDs) > 0 {
		for _, groupID := range req.GroupIDs {
			// Validate group exists
			_, err := s.groupRepo.FindByID(groupID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, errors.New("group not found: " + groupID)
				}
				return nil, err
			}

			// Check if group target already exists for this period
			_, err = s.monthlyTargetRepo.FindByGroupAndPeriod(groupID, req.Year, req.Month)
			if err == nil {
				return nil, ErrMonthlyTargetAlreadyExists
			}
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}

			// Create monthly target
			t := &monthlytargetdomain.MonthlyTarget{
				GroupID:   &groupID,
				UserID:       nil,
				BrickID:      nil,
				Year:         req.Year,
				Month:        req.Month,
				TargetAmount: req.TargetAmount,
			}

			if err := s.monthlyTargetRepo.Create(t); err != nil {
				return nil, err
			}

			// Reload to get timestamps and relations
			createdTarget, err := s.monthlyTargetRepo.FindByID(t.ID)
			if err != nil {
				return nil, err
			}

			createdTargets = append(createdTargets, *createdTarget.ToMonthlyTargetResponse())
		}
	}

	// Handle bulk create for users
	if len(req.UserIDs) > 0 {
		for _, userID := range req.UserIDs {
			// Validate user exists
			_, err := s.userRepo.FindByID(userID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, errors.New("user not found: " + userID)
				}
				return nil, err
			}

			// Check if user target already exists for this period
			_, err = s.monthlyTargetRepo.FindByUserAndPeriod(userID, req.Year, req.Month)
			if err == nil {
				return nil, ErrMonthlyTargetAlreadyExists
			}
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}

			// Create monthly target
			t := &monthlytargetdomain.MonthlyTarget{
				GroupID:   nil,
				UserID:       &userID,
				BrickID:      nil,
				Year:         req.Year,
				Month:        req.Month,
				TargetAmount: req.TargetAmount,
			}

			if err := s.monthlyTargetRepo.Create(t); err != nil {
				return nil, err
			}

			// Reload to get timestamps and relations
			createdTarget, err := s.monthlyTargetRepo.FindByID(t.ID)
			if err != nil {
				return nil, err
			}

			createdTargets = append(createdTargets, *createdTarget.ToMonthlyTargetResponse())
		}
	}

	// Handle bulk create for bricks
	if len(req.BrickIDs) > 0 {
		for _, brickID := range req.BrickIDs {
			// Validate brick exists
			_, err := s.brickRepo.FindByID(brickID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, errors.New("brick not found: " + brickID)
				}
				return nil, err
			}

			// Check if brick target already exists for this period
			_, err = s.monthlyTargetRepo.FindByBrickAndPeriod(brickID, req.Year, req.Month)
			if err == nil {
				return nil, ErrMonthlyTargetAlreadyExists
			}
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}

			// Create monthly target
			t := &monthlytargetdomain.MonthlyTarget{
				GroupID:   nil,
				UserID:       nil,
				BrickID:      &brickID,
				Year:         req.Year,
				Month:        req.Month,
				TargetAmount: req.TargetAmount,
			}

			if err := s.monthlyTargetRepo.Create(t); err != nil {
				return nil, err
			}

			// Reload to get timestamps and relations
			createdTarget, err := s.monthlyTargetRepo.FindByID(t.ID)
			if err != nil {
				return nil, err
			}

			createdTargets = append(createdTargets, *createdTarget.ToMonthlyTargetResponse())
		}
	}

	return createdTargets, nil
}

// CreateGroupTargetWithUserAssignment creates a group target and automatically assigns it to all users in the group
func (s *Service) CreateGroupTargetWithUserAssignment(req *monthlytargetdomain.CreateGroupTargetWithUserAssignmentRequest) (*monthlytargetdomain.MonthlyTargetResponse, []monthlytargetdomain.MonthlyTargetResponse, error) {
	// Validate group exists
	_, err := s.groupRepo.FindByID(req.GroupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("group not found")
		}
		return nil, nil, err
	}

	// Check if group target already exists for this period
	_, err = s.monthlyTargetRepo.FindByGroupAndPeriod(req.GroupID, req.Year, req.Month)
	if err == nil {
		return nil, nil, ErrMonthlyTargetAlreadyExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, err
	}

	// Create group target
	groupTarget := &monthlytargetdomain.MonthlyTarget{
		GroupID:   &req.GroupID,
		UserID:       nil,
		Year:         req.Year,
		Month:        req.Month,
		TargetAmount: req.TargetAmount,
	}

	if err := s.monthlyTargetRepo.Create(groupTarget); err != nil {
		return nil, nil, err
	}

	// Reload group target to get timestamps and relations
	createdGroupTarget, err := s.monthlyTargetRepo.FindByID(groupTarget.ID)
	if err != nil {
		return nil, nil, err
	}

	// Get all users in group
	users, err := s.userRepo.GetUsersByGroupID(req.GroupID)
	if err != nil {
		return nil, nil, err
	}

	// Create target for each user (if not already exists)
	var createdUserTargets []monthlytargetdomain.MonthlyTargetResponse
	for _, u := range users {
		// Check if user target already exists for this period
		_, err := s.monthlyTargetRepo.FindByUserAndPeriod(u.ID, req.Year, req.Month)
		if err == nil {
			// User target already exists, skip (user can edit it later)
			continue
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, err
		}

		// Create user target with same amount as group target
		userTarget := &monthlytargetdomain.MonthlyTarget{
			GroupID:   nil,
			UserID:       &u.ID,
			Year:         req.Year,
			Month:        req.Month,
			TargetAmount: req.TargetAmount,
		}

		if err := s.monthlyTargetRepo.Create(userTarget); err != nil {
			return nil, nil, err
		}

		// Reload user target to get timestamps and relations
		createdUserTarget, err := s.monthlyTargetRepo.FindByID(userTarget.ID)
		if err != nil {
			return nil, nil, err
		}

		createdUserTargets = append(createdUserTargets, *createdUserTarget.ToMonthlyTargetResponse())
	}

	return createdGroupTarget.ToMonthlyTargetResponse(), createdUserTargets, nil
}

// BulkSetTarget sets monthly targets for a specific entity over a range of months
func (s *Service) BulkSetTarget(req *monthlytargetdomain.BulkSetTargetRequest) ([]monthlytargetdomain.MonthlyTargetResponse, error) {
	// Validate invalid months
	if req.StartMonth > req.EndMonth {
		return nil, errors.New("start_month cannot be greater than end_month")
	}

	// Validate scope: either group_id, user_id, or brick_id must be provided, but only one
	scopeCount := 0
	if req.GroupID != nil && *req.GroupID != "" {
		scopeCount++
	}
	if req.UserID != nil && *req.UserID != "" {
		scopeCount++
	}
	if req.BrickID != nil && *req.BrickID != "" {
		scopeCount++
	}
	
	if scopeCount == 0 {
		return nil, ErrInvalidTargetScope
	}
	if scopeCount > 1 {
		return nil, ErrInvalidTargetScope
	}

	// Validate entity exists
	if req.GroupID != nil && *req.GroupID != "" {
		_, err := s.groupRepo.FindByID(*req.GroupID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("group not found")
			}
			return nil, err
		}
	} else if req.UserID != nil && *req.UserID != "" {
		_, err := s.userRepo.FindByID(*req.UserID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("user not found")
			}
			return nil, err
		}
	} else if req.BrickID != nil && *req.BrickID != "" {
		_, err := s.brickRepo.FindByID(*req.BrickID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("brick not found")
			}
			return nil, err
		}
	}

	var results []monthlytargetdomain.MonthlyTargetResponse

	// Loop through months
	for m := req.StartMonth; m <= req.EndMonth; m++ {
		// potential race condition if done in parallel, but sequential loop is fine for <12 items
		
		// Check if target exists for this month
		var existing *monthlytargetdomain.MonthlyTarget
		var errFind error

		if req.GroupID != nil {
			existing, errFind = s.monthlyTargetRepo.FindByGroupAndPeriod(*req.GroupID, req.Year, m)
		} else if req.UserID != nil {
			existing, errFind = s.monthlyTargetRepo.FindByUserAndPeriod(*req.UserID, req.Year, m)
		} else if req.BrickID != nil {
			existing, errFind = s.monthlyTargetRepo.FindByBrickAndPeriod(*req.BrickID, req.Year, m)
		}

		if errFind != nil && !errors.Is(errFind, gorm.ErrRecordNotFound) {
			return nil, errFind
		}

		if existing != nil {
			// Update existing
			existing.TargetAmount = req.TargetAmount
			if err := s.monthlyTargetRepo.Update(existing); err != nil {
				return nil, err
			}
			// Reload to get updated timestamps
			updatedTarget, err := s.monthlyTargetRepo.FindByID(existing.ID)
			if err != nil {
				return nil, err
			}
			results = append(results, *updatedTarget.ToMonthlyTargetResponse())

		} else {
			// Create new
			t := &monthlytargetdomain.MonthlyTarget{
				GroupID:   req.GroupID,
				UserID:       req.UserID,
				BrickID:      req.BrickID,
				Year:         req.Year,
				Month:        m,
				TargetAmount: req.TargetAmount,
			}

			if err := s.monthlyTargetRepo.Create(t); err != nil {
				return nil, err
			}
			// Reload
			createdTarget, err := s.monthlyTargetRepo.FindByID(t.ID)
			if err != nil {
				return nil, err
			}
			results = append(results, *createdTarget.ToMonthlyTargetResponse())
		}
	}

	return results, nil
}

