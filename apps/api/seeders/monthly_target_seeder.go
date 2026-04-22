package seeders

import (
	"log"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	brickdomain "github.com/gilabs/crm-healthcare/api/internal/domain/brick"
	"github.com/gilabs/crm-healthcare/api/internal/domain/group"
	monthlytargetdomain "github.com/gilabs/crm-healthcare/api/internal/domain/monthly_target"
	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
)

// SeedMonthlyTargets seeds initial monthly target data
func SeedMonthlyTargets() error {
	// Seed targets in a month-aware and type-aware way.
	// Important: brick performance needs brick-level targets, so we must not skip
	// a month just because user/group targets already exist.
	// Get users
	var users []user.User
	if err := database.DB.Find(&users).Error; err != nil {
		return err
	}
	if len(users) == 0 {
		log.Println("Warning: No users found, skipping monthly target seeding")
		return nil
	}

	// Get groups
	var groups []group.Group
	if err := database.DB.Find(&groups).Error; err != nil {
		log.Printf("Warning: Failed to get groups: %v", err)
		// Continue without groups - we can still seed user targets
	}

	now := time.Now()
	currentYear := now.Year()
	currentMonth := int(now.Month())

	// Get bricks for brick targets
	var bricks []brickdomain.Brick
	if err := database.DB.Find(&bricks).Error; err != nil {
		log.Printf("Warning: Failed to get bricks: %v", err)
	}

	// Seed monthly targets for current month and last 3 months
	// Create user targets for all active users
	// Target amount varies by role: sales users get higher targets
	targets := []monthlytargetdomain.MonthlyTarget{}

	// Get roles to determine target amounts
	var salesRole, salesManagerRole role.Role
	var salesRoleID, salesManagerRoleID string
	if err := database.DB.Where("code = ?", "sales").First(&salesRole).Error; err == nil {
		salesRoleID = salesRole.ID
	}
	if err := database.DB.Where("code = ?", "sales_manager").First(&salesManagerRole).Error; err == nil {
		salesManagerRoleID = salesManagerRole.ID
	}

	// Create targets for last 3 months and current month
	monthsToSeed := []int{currentMonth - 3, currentMonth - 2, currentMonth - 1, currentMonth}
	for _, month := range monthsToSeed {
		// Adjust year if month is negative
		targetYear := currentYear
		if month <= 0 {
			month += 12
			targetYear--
		}

		// Check if targets for this month/year already exist by type
		var userCount int64
		database.DB.Model(&monthlytargetdomain.MonthlyTarget{}).
			Where("year = ? AND month = ? AND user_id IS NOT NULL", targetYear, month).
			Count(&userCount)
		var groupCount int64
		database.DB.Model(&monthlytargetdomain.MonthlyTarget{}).
			Where("year = ? AND month = ? AND group_id IS NOT NULL", targetYear, month).
			Count(&groupCount)
		var brickCount int64
		database.DB.Model(&monthlytargetdomain.MonthlyTarget{}).
			Where("year = ? AND month = ? AND brick_id IS NOT NULL", targetYear, month).
			Count(&brickCount)

		// Create user targets for this month
		if userCount == 0 {
			for _, u := range users {
			// Only create for active users
			if u.Status != "active" {
				continue
			}

			// Determine target amount based on role
			var targetAmount int64
			if u.RoleID == salesRoleID {
				// Sales users: 100M - 150M rupiah (varied)
				targetAmount = int64(10000000000 + (len(targets)%5)*1000000000) // 100M - 150M
			} else if u.RoleID == salesManagerRoleID {
				// Sales managers: 200M - 300M rupiah
				targetAmount = int64(20000000000 + (len(targets)%5)*2000000000) // 200M - 300M
			} else {
				// Other roles: 50M rupiah
				targetAmount = int64(5000000000) // 50M
			}

				target := monthlytargetdomain.MonthlyTarget{
					UserID:       &u.ID,
					GroupID:      nil,
					BrickID:      nil,
					Year:         targetYear,
					Month:        month,
					TargetAmount: targetAmount,
				}
				targets = append(targets, target)
			}
		}

		// Create group targets for this month (optional - only if groups exist)
		if len(groups) > 0 && groupCount == 0 {
			for _, g := range groups {
			// Only create for active groups
			if g.Status != "active" {
				continue
			}

			// Group targets: sum of user targets in group (simplified: 500M per group)
			groupTargetAmount := int64(50000000000) // 500M rupiah
			groupID := g.ID
				target := monthlytargetdomain.MonthlyTarget{
					GroupID:      &groupID,
					UserID:       nil,
					BrickID:      nil,
					Year:         targetYear,
					Month:        month,
					TargetAmount: groupTargetAmount,
				}
				targets = append(targets, target)
			}
		}

		// Create brick targets for this month (if bricks exist)
		if len(bricks) > 0 && brickCount == 0 {
			for _, brick := range bricks {
			if brick.Status != "active" {
				continue
			}

			// Brick targets: 300M - 800M rupiah (varied)
			brickTargetAmount := int64(30000000000 + (len(targets)%6)*10000000000) // 300M - 800M
			brickID := brick.ID
				target := monthlytargetdomain.MonthlyTarget{
					BrickID:      &brickID,
					UserID:       nil,
					GroupID:      nil,
					Year:         targetYear,
					Month:        month,
					TargetAmount: brickTargetAmount,
				}
				targets = append(targets, target)
			}
		}
	}

	// Create targets in batches
	batchSize := 100
	for i := 0; i < len(targets); i += batchSize {
		end := i + batchSize
		if end > len(targets) {
			end = len(targets)
		}
		batch := targets[i:end]
		if err := database.DB.Create(&batch).Error; err != nil {
			return err
		}
		log.Printf("Created %d monthly targets (batch %d/%d)", len(batch), (i/batchSize)+1, (len(targets)+batchSize-1)/batchSize)
	}

	log.Printf("✅ Seeded %d monthly targets successfully (Year: %d, Month: %d)", len(targets), currentYear, currentMonth)
	return nil
}

