package brick_target_distribution

import (
	brickdomain "github.com/gilabs/crm-healthcare/api/internal/domain/brick"
	bricktargetdistributiondomain "github.com/gilabs/crm-healthcare/api/internal/domain/brick_target_distribution"
	monthlytargetdomain "github.com/gilabs/crm-healthcare/api/internal/domain/monthly_target"
)

// BrickTargetWithDistributionsResponse represents a service-level aggregate response.
type BrickTargetWithDistributionsResponse struct {
	Brick            *brickdomain.BrickResponse                                      `json:"brick"`
	Target           *monthlytargetdomain.MonthlyTargetResponse                      `json:"target"`
	Distributions    []bricktargetdistributiondomain.BrickTargetDistributionResponse `json:"distributions"`
	TotalDistributed int64                                                           `json:"total_distributed"`
	Remaining        int64                                                           `json:"remaining"`
}
