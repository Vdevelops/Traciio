package interfaces

import "time"

// KPIRepository defines the KPI data-access contract.
// The repository exposes focused aggregations so the service can orchestrate
// scorecard assembly without one giant query.
type KPIRepository interface {
	CountDealsCreated(userID string, startDate, endDate time.Time) (int64, error)
	CountWonDeals(userID string, startDate, endDate time.Time) (int64, error)
	SumWonRevenue(userID string, startDate, endDate time.Time) (int64, error)
	CountVisitCompleted(userID string, startDate, endDate time.Time) (int64, error)
	CountVisitPlanned(userID string, startDate, endDate time.Time) (int64, error)
	CountTasksCreated(userID string, startDate, endDate time.Time) (int64, error)
	CountTasksCompleted(userID string, startDate, endDate time.Time) (int64, error)
	CountOverdueTasks(userID string, startDate, endDate time.Time) (int64, error)
	GetPipelineMovementScore(userID string, startDate, endDate time.Time) (int64, error)
	CountDealsWithoutBrick(userID string, startDate, endDate time.Time) (int64, error)
	CountVisitReportsWithoutBrick(userID string, startDate, endDate time.Time) (int64, error)
	CountCustomersWithInteractionInBrick(brickID string, startDate, endDate time.Time) (int64, error)
	CountRegisteredCustomersInBrick(brickID string) (int64, error)
}
