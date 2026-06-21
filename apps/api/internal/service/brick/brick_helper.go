package brick

import (
	"errors"
	"fmt"
	"strings"

	brickdomain "github.com/gilabs/crm-healthcare/api/internal/domain/brick"
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

	normalizedProvince := strings.TrimSpace(province)
	for _, candidate := range buildRegencyCandidates(regency) {
		brick, err := h.brickRepo.FindByRegencyAndProvince(candidate, normalizedProvince)
		if err == nil {
			return &brick.ID, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	return nil, nil
}

// EnsureBrickIDForLocation returns an existing brick for a location,
// or creates one automatically when none exists yet.
func (h *BrickHelper) EnsureBrickIDForLocation(province, regency string) (*string, error) {
	if province == "" || regency == "" {
		return nil, nil
	}

	brickID, err := h.GetBrickIDFromLocation(province, regency)
	if err != nil || brickID != nil {
		return brickID, err
	}

	normalizedProvince := strings.TrimSpace(province)
	normalizedRegency := strings.TrimSpace(regency)

	codePrefix := buildCodePrefix(normalizedProvince)
	sequence, err := h.brickRepo.GetNextCodeSequence(codePrefix)
	if err != nil {
		return nil, err
	}

	newBrick := &brickdomain.Brick{
		Name:     normalizedRegency,
		Code:     fmt.Sprintf("%s-%03d", codePrefix, sequence),
		Province: normalizedProvince,
		Regency:  normalizedRegency,
		Status:   "active",
	}

	if err := h.brickRepo.Create(newBrick); err != nil {
		existingBrickID, lookupErr := h.GetBrickIDFromLocation(normalizedProvince, normalizedRegency)
		if lookupErr == nil && existingBrickID != nil {
			return existingBrickID, nil
		}
		if lookupErr != nil {
			return nil, lookupErr
		}
		return nil, err
	}

	return &newBrick.ID, nil
}

func buildRegencyCandidates(regency string) []string {
	normalized := strings.Join(strings.Fields(strings.TrimSpace(regency)), " ")
	if normalized == "" {
		return nil
	}

	lower := strings.ToLower(normalized)
	withoutPrefix := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(lower, "kota "), "kabupaten "))
	withoutSuffix := strings.TrimSpace(strings.TrimSuffix(strings.TrimSuffix(withoutPrefix, " kota"), " kabupaten"))

	candidates := []string{normalized}
	if withoutSuffix != "" && !strings.EqualFold(withoutSuffix, normalized) {
		candidates = append(candidates, withoutSuffix)
		candidates = append(candidates, "Kota "+withoutSuffix)
		candidates = append(candidates, "Kabupaten "+withoutSuffix)
	}

	seen := make(map[string]struct{}, len(candidates))
	deduped := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		candidate = strings.Join(strings.Fields(strings.TrimSpace(candidate)), " ")
		if candidate == "" {
			continue
		}
		key := strings.ToLower(candidate)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		deduped = append(deduped, candidate)
	}

	return deduped
}
