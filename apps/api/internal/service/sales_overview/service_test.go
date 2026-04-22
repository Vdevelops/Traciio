package sales_overview_test

import (
	"errors"
	"testing"

	"github.com/gilabs/crm-healthcare/api/internal/domain/sales_overview"
	"github.com/gilabs/crm-healthcare/api/internal/repository/mocks"
	salesoverviewsvc "github.com/gilabs/crm-healthcare/api/internal/service/sales_overview"
	"github.com/stretchr/testify/assert"
)

func TestListSalesPerformance(t *testing.T) {
	// Setup mocks
	repo := new(mocks.SalesOverviewRepository)
	targetRepo := new(mocks.MonthlyTargetRepository)

	// Create service
	svc := salesoverviewsvc.NewService(repo, targetRepo)

	t.Run("Success", func(t *testing.T) {
		req := &sales_overview.ListSalesPerformanceRequest{
			Page:    1,
			PerPage: 10,
			Period:  "month",
		}
		
		expectedResp := []sales_overview.SalesPerformanceListResponse{
			{
				UserID: "user1",
				UserName:   "John Doe",
				TotalRevenue: 1000000,
			},
		}


		repo.On("ListSalesPerformance", req).Return(expectedResp, int64(1), nil)
		targetRepo.On("BatchGetProratedTargetsForPeriod", []string{"user1"}, req.StartDate, req.EndDate).Return(map[string]float64{"user1": 500000}, nil)

		resp, total, err := svc.ListSalesPerformance(req)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, resp, 1)
		assert.Equal(t, "John Doe", resp[0].UserName)
	})

	t.Run("Repository Error", func(t *testing.T) {
		req := &sales_overview.ListSalesPerformanceRequest{
			Page:    1,
			PerPage: 10,
		}

		repo.On("ListSalesPerformance", req).Return(nil, int64(0), errors.New("db error"))

		resp, total, err := svc.ListSalesPerformance(req)

		assert.Error(t, err)
		assert.Equal(t, int64(0), total)
		assert.Nil(t, resp)
		assert.Equal(t, "db error", err.Error())
	})
}
