package brick_target_distribution

import (
	"errors"

	bricktargetdistributiondomain "github.com/gilabs/crm-healthcare/api/internal/domain/brick_target_distribution"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"gorm.io/gorm"
)

var (
	ErrDistributionNotFound = errors.New("brick target distribution not found")
	ErrInvalidSalesBrick    = errors.New("sales user must be in the same brick")
	ErrInvalidManager       = errors.New("only brick manager can distribute targets")
	ErrBrickTargetNotFound  = errors.New("brick target not found")
	ErrSalesNotFound        = errors.New("sales user not found")
)

type Service struct {
	distributionRepo  interfaces.BrickTargetDistributionRepository
	brickRepo         interfaces.BrickRepository
	monthlyTargetRepo interfaces.MonthlyTargetRepository
	userRepo          interfaces.UserRepository
}

func NewService(
	distributionRepo interfaces.BrickTargetDistributionRepository,
	brickRepo interfaces.BrickRepository,
	monthlyTargetRepo interfaces.MonthlyTargetRepository,
	userRepo interfaces.UserRepository,
) *Service {
	return &Service{
		distributionRepo:  distributionRepo,
		brickRepo:         brickRepo,
		monthlyTargetRepo: monthlyTargetRepo,
		userRepo:          userRepo,
	}
}

// PaginationResult represents pagination result
type PaginationResult struct {
	Page       int
	PerPage    int
	Total      int
	TotalPages int
}

// List returns a list of brick target distributions with pagination
func (s *Service) List(req *bricktargetdistributiondomain.ListBrickTargetDistributionsRequest) ([]bricktargetdistributiondomain.BrickTargetDistributionResponse, *PaginationResult, error) {
	distributions, total, err := s.distributionRepo.List(req)
	if err != nil {
		return nil, nil, err
	}

	responses := make([]bricktargetdistributiondomain.BrickTargetDistributionResponse, len(distributions))
	for i, d := range distributions {
		responses[i] = *d.ToBrickTargetDistributionResponse()
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

// GetByID returns a brick target distribution by ID
func (s *Service) GetByID(id string) (*bricktargetdistributiondomain.BrickTargetDistributionResponse, error) {
	d, err := s.distributionRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDistributionNotFound
		}
		return nil, err
	}
	return d.ToBrickTargetDistributionResponse(), nil
}

// GetByBrickTargetID returns all distributions for a brick target
func (s *Service) GetByBrickTargetID(brickTargetID string) ([]bricktargetdistributiondomain.BrickTargetDistributionResponse, error) {
	distributions, err := s.distributionRepo.FindByBrickTargetID(brickTargetID)
	if err != nil {
		return nil, err
	}

	responses := make([]bricktargetdistributiondomain.BrickTargetDistributionResponse, len(distributions))
	for i, d := range distributions {
		responses[i] = *d.ToBrickTargetDistributionResponse()
	}

	return responses, nil
}

func (s *Service) GetBrickTargetWithDistributions(brickID string, year, month int) (*BrickTargetWithDistributionsResponse, error) {
	// Get brick
	brick, err := s.brickRepo.FindByID(brickID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("brick not found")
		}
		return nil, err
	}

	// Get brick target
	target, err := s.monthlyTargetRepo.FindByBrickAndPeriod(brickID, year, month)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBrickTargetNotFound
		}
		return nil, err
	}

	// Get distributions
	distributions, err := s.distributionRepo.FindByBrickTargetID(target.ID)
	if err != nil {
		return nil, err
	}

	// Convert distributions to response
	distributionResponses := make([]bricktargetdistributiondomain.BrickTargetDistributionResponse, len(distributions))
	var totalDistributed int64 = 0
	for i, d := range distributions {
		distributionResponses[i] = *d.ToBrickTargetDistributionResponse()
		totalDistributed += d.DistributedAmount
	}

	// Calculate remaining
	remaining := target.TargetAmount - totalDistributed

	return &BrickTargetWithDistributionsResponse{
		Brick:            brick.ToBrickResponse(),
		Target:           target.ToMonthlyTargetResponse(),
		Distributions:    distributionResponses,
		TotalDistributed: totalDistributed,
		Remaining:        remaining,
	}, nil
}

// Create creates a new brick target distribution
func (s *Service) Create(brickID, brickTargetID string, req *bricktargetdistributiondomain.CreateBrickTargetDistributionRequest, distributedBy string) (*bricktargetdistributiondomain.BrickTargetDistributionResponse, error) {
	// Validate brick exists
	brick, err := s.brickRepo.FindByID(brickID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("brick not found")
		}
		return nil, err
	}

	// Validate brick target exists and belongs to this brick
	brickTarget, err := s.monthlyTargetRepo.FindByID(brickTargetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBrickTargetNotFound
		}
		return nil, err
	}
	if brickTarget.BrickID == nil || *brickTarget.BrickID != brickID {
		return nil, errors.New("brick target does not belong to this brick")
	}

	// Validate sales user exists and is in the same brick
	salesUser, err := s.userRepo.FindByID(req.SalesUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSalesNotFound
		}
		return nil, err
	}
	if salesUser.BrickID == nil || *salesUser.BrickID != brickID {
		return nil, ErrInvalidSalesBrick
	}

	// Validate distributed_by is the manager of the brick
	if brick.ManagerID == nil || *brick.ManagerID != distributedBy {
		return nil, ErrInvalidManager
	}

	// Check if distribution already exists for this sales and target
	existing, err := s.distributionRepo.FindBySalesUserIDAndBrickTargetID(req.SalesUserID, brickTargetID)
	if err == nil {
		// Update existing distribution
		existing.DistributedAmount = req.DistributedAmount
		if err := s.distributionRepo.Update(existing); err != nil {
			return nil, err
		}
		// Reload
		updated, err := s.distributionRepo.FindByID(existing.ID)
		if err != nil {
			return nil, err
		}
		return updated.ToBrickTargetDistributionResponse(), nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Create new distribution
	distribution := &bricktargetdistributiondomain.BrickTargetDistribution{
		BrickID:           brickID,
		BrickTargetID:     brickTargetID,
		SalesUserID:       req.SalesUserID,
		DistributedAmount: req.DistributedAmount,
		DistributedBy:     distributedBy,
	}

	if err := s.distributionRepo.Create(distribution); err != nil {
		return nil, err
	}

	// Reload to get timestamps and relations
	createdDistribution, err := s.distributionRepo.FindByID(distribution.ID)
	if err != nil {
		return nil, err
	}

	return createdDistribution.ToBrickTargetDistributionResponse(), nil
}

// BulkCreate creates multiple brick target distributions
func (s *Service) BulkCreate(brickID, brickTargetID string, req *bricktargetdistributiondomain.BulkCreateBrickTargetDistributionRequest, distributedBy string) ([]bricktargetdistributiondomain.BrickTargetDistributionResponse, error) {
	// Validate brick exists
	brick, err := s.brickRepo.FindByID(brickID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("brick not found")
		}
		return nil, err
	}

	// Validate brick target exists and belongs to this brick
	brickTarget, err := s.monthlyTargetRepo.FindByID(brickTargetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBrickTargetNotFound
		}
		return nil, err
	}
	if brickTarget.BrickID == nil || *brickTarget.BrickID != brickID {
		return nil, errors.New("brick target does not belong to this brick")
	}

	// Validate distributed_by is the manager of the brick
	if brick.ManagerID == nil || *brick.ManagerID != distributedBy {
		return nil, ErrInvalidManager
	}

	// Prepare distributions
	distributions := make([]*bricktargetdistributiondomain.BrickTargetDistribution, 0, len(req.Distributions))

	for _, distReq := range req.Distributions {
		// Validate sales user exists and is in the same brick
		salesUser, err := s.userRepo.FindByID(distReq.SalesUserID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("sales user not found: " + distReq.SalesUserID)
			}
			return nil, err
		}
		if salesUser.BrickID == nil || *salesUser.BrickID != brickID {
			return nil, errors.New("sales user must be in the same brick: " + distReq.SalesUserID)
		}

		// Check if distribution already exists
		existing, err := s.distributionRepo.FindBySalesUserIDAndBrickTargetID(distReq.SalesUserID, brickTargetID)
		if err == nil {
			// Update existing
			existing.DistributedAmount = distReq.DistributedAmount
			if err := s.distributionRepo.Update(existing); err != nil {
				return nil, err
			}
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		// Create new
		distribution := &bricktargetdistributiondomain.BrickTargetDistribution{
			BrickID:           brickID,
			BrickTargetID:     brickTargetID,
			SalesUserID:       distReq.SalesUserID,
			DistributedAmount: distReq.DistributedAmount,
			DistributedBy:     distributedBy,
		}
		distributions = append(distributions, distribution)
	}

	// Bulk create new distributions
	if len(distributions) > 0 {
		if err := s.distributionRepo.BulkCreate(distributions); err != nil {
			return nil, err
		}
	}

	// Reload all distributions for this target
	allDistributions, err := s.distributionRepo.FindByBrickTargetID(brickTargetID)
	if err != nil {
		return nil, err
	}

	responses := make([]bricktargetdistributiondomain.BrickTargetDistributionResponse, len(allDistributions))
	for i, d := range allDistributions {
		responses[i] = *d.ToBrickTargetDistributionResponse()
	}

	return responses, nil
}

// Update updates a brick target distribution
func (s *Service) Update(brickID, brickTargetID, distributionID string, req *bricktargetdistributiondomain.UpdateBrickTargetDistributionRequest, distributedBy string) (*bricktargetdistributiondomain.BrickTargetDistributionResponse, error) {
	// Validate brick exists
	brick, err := s.brickRepo.FindByID(brickID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("brick not found")
		}
		return nil, err
	}

	// Validate distributed_by is the manager of the brick
	if brick.ManagerID == nil || *brick.ManagerID != distributedBy {
		return nil, ErrInvalidManager
	}

	// Find existing distribution
	distribution, err := s.distributionRepo.FindByID(distributionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDistributionNotFound
		}
		return nil, err
	}

	// Validate distribution belongs to the brick and target
	if distribution.BrickID != brickID || distribution.BrickTargetID != brickTargetID {
		return nil, errors.New("distribution does not belong to this brick target")
	}

	// Update fields
	if req.DistributedAmount != nil {
		distribution.DistributedAmount = *req.DistributedAmount
	}

	if err := s.distributionRepo.Update(distribution); err != nil {
		return nil, err
	}

	// Reload to get updated timestamps and relations
	updatedDistribution, err := s.distributionRepo.FindByID(distributionID)
	if err != nil {
		return nil, err
	}

	return updatedDistribution.ToBrickTargetDistributionResponse(), nil
}

// Delete deletes a brick target distribution
func (s *Service) Delete(brickID, brickTargetID, distributionID string, distributedBy string) error {
	// Validate brick exists
	brick, err := s.brickRepo.FindByID(brickID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("brick not found")
		}
		return err
	}

	// Validate distributed_by is the manager of the brick
	if brick.ManagerID == nil || *brick.ManagerID != distributedBy {
		return ErrInvalidManager
	}

	// Find existing distribution
	distribution, err := s.distributionRepo.FindByID(distributionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrDistributionNotFound
		}
		return err
	}

	// Validate distribution belongs to the brick and target
	if distribution.BrickID != brickID || distribution.BrickTargetID != brickTargetID {
		return errors.New("distribution does not belong to this brick target")
	}

	// Delete distribution (soft delete)
	return s.distributionRepo.Delete(distributionID)
}
