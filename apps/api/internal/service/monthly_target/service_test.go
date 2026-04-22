package monthly_target_test

import (
	"testing"

	"github.com/gilabs/crm-healthcare/api/internal/domain/monthly_target"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"github.com/gilabs/crm-healthcare/api/internal/repository/mocks"
	targetsvc "github.com/gilabs/crm-healthcare/api/internal/service/monthly_target"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

const newMonthlyTargetID = "new-id"

func TestCreateMonthlyTarget(t *testing.T) {
	t.Run("Success - User Target", func(t *testing.T) {
		repo := new(mocks.MonthlyTargetRepository)
		groupRepo := new(mocks.GroupRepository) // Assume GroupRepository mock exists or unused
		userRepo := new(mocks.UserRepository)
		brickRepo := new(mocks.BrickRepository) // Assume BrickRepository mock exists or unused

		svc := targetsvc.NewService(repo, groupRepo, userRepo, brickRepo)

		userID := "user1"
		req := &monthly_target.CreateMonthlyTargetRequest{
			UserID: &userID,
			Month:  1,
			Year:   2026,
			TargetAmount: 10000000,
		}

		userRepo.On("FindByID", "user1").Return(&user.User{ID: "user1"}, nil)
		repo.On("FindByUserAndPeriod", "user1", 2026, 1).Return(nil, gorm.ErrRecordNotFound) // Should not exist
		
		// Mock Create to set ID
		repo.On("Create", mock.AnythingOfType("*monthly_target.MonthlyTarget")).Run(func(args mock.Arguments) {
			mt := args.Get(0).(*monthly_target.MonthlyTarget)
			mt.ID = newMonthlyTargetID
		}).Return(nil)
		
		// Mock FindByID for reload
		repo.On("FindByID", newMonthlyTargetID).Return(&monthly_target.MonthlyTarget{
			ID: newMonthlyTargetID,
			UserID: &userID,
			Year: 2026,
			Month: 1,
			TargetAmount: 10000000,
		}, nil)

		resp, err := svc.Create(req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, int64(10000000), resp.TargetAmount)
	})

	t.Run("Failure - Duplicate Target", func(t *testing.T) {
		repo := new(mocks.MonthlyTargetRepository)
		groupRepo := new(mocks.GroupRepository)
		userRepo := new(mocks.UserRepository)
		brickRepo := new(mocks.BrickRepository)

		svc := targetsvc.NewService(repo, groupRepo, userRepo, brickRepo)

		userID := "user1"
		req := &monthly_target.CreateMonthlyTargetRequest{
			UserID: &userID,
			Month:  1,
			Year:   2026,
			TargetAmount: 10000000,
		}

		userRepo.On("FindByID", "user1").Return(&user.User{ID: "user1"}, nil)
		// Return existing target
		repo.On("FindByUserAndPeriod", "user1", 2026, 1).Return(&monthly_target.MonthlyTarget{ID: "existing"}, nil)

		resp, err := svc.Create(req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "already exists")
	})
}
