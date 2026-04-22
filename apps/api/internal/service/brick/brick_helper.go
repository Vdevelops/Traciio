package brick

import (
	"errors"

	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

// BrickHelper provides helper functions for brick-related operations
// This helps avoid N+1 queries by providing optimized helper methods
type BrickHelper struct {
	userRepo    interfaces.UserRepository
	brickRepo   interfaces.BrickRepository
	accountRepo interfaces.AccountRepository
}

// NewBrickHelper creates a new brick helper
func NewBrickHelper(userRepo interfaces.UserRepository, brickRepo interfaces.BrickRepository, accountRepo interfaces.AccountRepository) *BrickHelper {
	return &BrickHelper{
		userRepo:    userRepo,
		brickRepo:   brickRepo,
		accountRepo: accountRepo,
	}
}

// GetBrickIDFromUser gets brick_id from user's brick_id field
// Returns nil if user doesn't have brick_id or user not found
func (h *BrickHelper) GetBrickIDFromUser(userID string) (*string, error) {
	if userID == "" {
		return nil, nil
	}

	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return user.BrickID, nil
}

// GetBrickIDFromAccount gets brick_id from account's brick_id field
// If account doesn't have brick_id, try to get from assigned_to user's brick_id
// Returns nil if no brick_id found
func (h *BrickHelper) GetBrickIDFromAccount(accountID string) (*string, error) {
	if accountID == "" {
		return nil, nil
	}

	account, err := h.accountRepo.FindByID(accountID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// If account has brick_id, use it
	if account.BrickID != nil && *account.BrickID != "" {
		return account.BrickID, nil
	}

	// Otherwise, try to get from assigned_to user's brick_id
	if account.AssignedTo != nil && *account.AssignedTo != "" {
		return h.GetBrickIDFromUser(*account.AssignedTo)
	}

	return nil, nil
}

// GetBrickIDFromLocation tries to find brick_id based on province and regency/city
// This is used for auto-populating brick_id when creating accounts
// Returns nil if no matching brick found
func (h *BrickHelper) GetBrickIDFromLocation(province, regency string) (*string, error) {
	if province == "" || regency == "" {
		return nil, nil
	}

	brick, err := h.brickRepo.FindByRegencyAndProvince(regency, province)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &brick.ID, nil
}

