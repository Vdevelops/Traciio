package pipeline

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/customer_purchase"
	"github.com/gilabs/crm-healthcare/api/internal/domain/deal_history"
	domainevents "github.com/gilabs/crm-healthcare/api/internal/domain/events"
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead_status"
	"github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
	"github.com/gilabs/crm-healthcare/api/internal/domain/task"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	brick "github.com/gilabs/crm-healthcare/api/internal/service/brick"
	"github.com/gilabs/crm-healthcare/api/pkg/cache"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var (
	ErrPipelineStageNotFound   = errors.New("pipeline stage not found")
	ErrDealNotFound            = errors.New("deal not found")
	ErrAccountNotFound         = errors.New("account not found")
	ErrInvalidStage            = errors.New("invalid pipeline stage")
	ErrCloseReasonRequired     = errors.New("close reason is required")
	ErrStageRequirementsNotMet = errors.New("stage requirements not met")
)

type Service struct {
	pipelineRepo         interfaces.PipelineRepository
	dealRepo             interfaces.DealRepository
	accountRepo          interfaces.AccountRepository
	productRepo          interfaces.ProductRepository
	dealProductItemRepo  interfaces.DealProductItemRepository
	dealHistoryRepo      interfaces.DealHistoryRepository
	taskRepo             interfaces.TaskRepository
	visitReportRepo      interfaces.VisitReportRepository
	leadRepo             interfaces.LeadRepository
	customerPurchaseRepo interfaces.CustomerPurchaseHistoryRepository
	brickHelper          *brick.BrickHelper
	db                   *gorm.DB
	eventHelper          *domainevents.Helper
	cacheService         *cache.DealCacheService
}

func NewService(
	db *gorm.DB,
	pipelineRepo interfaces.PipelineRepository,
	dealRepo interfaces.DealRepository,
	accountRepo interfaces.AccountRepository,
	productRepo interfaces.ProductRepository,
	dealProductItemRepo interfaces.DealProductItemRepository,
	dealHistoryRepo interfaces.DealHistoryRepository,
	taskRepo interfaces.TaskRepository,
	visitReportRepo interfaces.VisitReportRepository,
	leadRepo interfaces.LeadRepository,
	customerPurchaseRepo interfaces.CustomerPurchaseHistoryRepository,
	brickHelper *brick.BrickHelper,
	eventHelper *domainevents.Helper,
) *Service {
	return &Service{
		db:                   db,
		pipelineRepo:         pipelineRepo,
		dealRepo:             dealRepo,
		accountRepo:          accountRepo,
		productRepo:          productRepo,
		dealProductItemRepo:  dealProductItemRepo,
		dealHistoryRepo:      dealHistoryRepo,
		taskRepo:             taskRepo,
		visitReportRepo:      visitReportRepo,
		leadRepo:             leadRepo,
		customerPurchaseRepo: customerPurchaseRepo,
		brickHelper:          brickHelper,
		eventHelper:          eventHelper,
		cacheService:         cache.NewDealCacheService(nil),
	}
}

type cachedDealListResult struct {
	Deals      []pipeline.DealResponse
	Pagination *PaginationResult
}

// ListStages returns a list of pipeline stages
func (s *Service) ListStages(req *pipeline.ListPipelineStagesRequest) ([]pipeline.PipelineStageResponse, error) {
	stages, err := s.pipelineRepo.ListStages(req)
	if err != nil {
		return nil, err
	}

	responses := make([]pipeline.PipelineStageResponse, len(stages))
	for i, stage := range stages {
		responses[i] = *stage.ToPipelineStageResponse()
	}

	return responses, nil
}

// GetStageByID returns a pipeline stage by ID
func (s *Service) GetStageByID(id string) (*pipeline.PipelineStageResponse, error) {
	stage, err := s.pipelineRepo.FindStageByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPipelineStageNotFound
		}
		return nil, err
	}

	return stage.ToPipelineStageResponse(), nil
}

// ListDeals returns a list of deals with pagination
func (s *Service) ListDeals(req *pipeline.ListDealsRequest) ([]pipeline.DealResponse, *PaginationResult, error) {
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

	filterMap := map[string]interface{}{}
	if req.Search != "" {
		filterMap["search"] = req.Search
	}
	if req.StageID != "" {
		filterMap["stage_id"] = req.StageID
	}
	if req.AccountID != "" {
		filterMap["account_id"] = req.AccountID
	}
	if req.AssignedTo != "" {
		filterMap["assigned_to"] = req.AssignedTo
	}
	if req.BrickID != "" {
		filterMap["brick_id"] = req.BrickID
	}
	if req.Status != "" {
		filterMap["status"] = req.Status
	}
	if req.Source != "" {
		filterMap["source"] = req.Source
	}
	if req.MinValue != nil {
		filterMap["min_value"] = *req.MinValue
	}
	if req.MaxValue != nil {
		filterMap["max_value"] = *req.MaxValue
	}
	if req.DateFrom != "" {
		filterMap["date_from"] = req.DateFrom
	}
	if req.DateTo != "" {
		filterMap["date_to"] = req.DateTo
	}
	if len(req.ScopedUserIDs) > 0 {
		filterMap["scoped_user_ids"] = fmt.Sprintf("%v", req.ScopedUserIDs)
	}

	var cachedResult cachedDealListResult
	if found, _ := s.cacheService.GetList(page, perPage, filterMap, &cachedResult); found && cachedResult.Pagination != nil {
		return cachedResult.Deals, cachedResult.Pagination, nil
	}

	deals, total, err := s.dealRepo.List(req)
	if err != nil {
		return nil, nil, err
	}

	responses := make([]pipeline.DealResponse, len(deals))
	for i, deal := range deals {
		responses[i] = *deal.ToDealResponse()
	}

	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	pagination := &PaginationResult{
		Page:       page,
		PerPage:    perPage,
		Total:      int(total),
		TotalPages: totalPages,
	}

	_ = s.cacheService.SetList(page, perPage, filterMap, cachedDealListResult{
		Deals:      responses,
		Pagination: pagination,
	})

	return responses, pagination, nil
}

// GetDealByID returns a deal by ID
func (s *Service) GetDealByID(id string) (*pipeline.DealResponse, error) {
	var cachedResponse pipeline.DealResponse
	if found, _ := s.cacheService.GetDetail(id, &cachedResponse); found && cachedResponse.ID != "" {
		return &cachedResponse, nil
	}

	deal, err := s.dealRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDealNotFound
		}
		return nil, err
	}

	resp := deal.ToDealResponse()
	_ = s.cacheService.SetDetail(id, resp)
	return resp, nil
}

// CreateDeal creates a new deal
func (s *Service) CreateDeal(req *pipeline.CreateDealRequest, createdBy string) (*pipeline.DealResponse, error) {
	if s.db == nil {
		return nil, errors.New("database connection not available")
	}

	var createdDeal *pipeline.Deal
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Validate account exists
		accountRecord, err := s.accountRepo.FindByID(req.AccountID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrAccountNotFound
			}
			return err
		}

		// Validate stage exists
		stage, err := s.pipelineRepo.FindStageByID(req.StageID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrInvalidStage
			}
			return err
		}

		// Set default status based on stage
		status := "open"
		var actualCloseDate *time.Time
		if stage.IsWon {
			status = "won"
			now := time.Now()
			actualCloseDate = &now
		} else if stage.IsLost {
			status = "lost"
			now := time.Now()
			actualCloseDate = &now
		}

		// Probability always follows the stage percent (fallback by stage order)
		probability := stage.Probability
		if probability == 0 {
			probability = stage.Order * 20
		}

		// Auto-populate brick_id if not provided
		var brickID *string
		if s.brickHelper != nil {
			// Account territory is the source of truth for deal brick assignment.
			if req.AccountID != "" {
				brickID, _ = s.brickHelper.GetBrickIDFromAccount(req.AccountID)
			}
			if brickID == nil && createdBy != "" {
				brickID, _ = s.brickHelper.GetBrickIDFromUser(createdBy)
			}
		}

		assignedTo := strings.TrimSpace(req.AssignedTo)
		if assignedTo == "" && req.LeadID != nil && *req.LeadID != "" && s.leadRepo != nil {
			if leadRecord, err := s.leadRepo.FindByID(*req.LeadID); err == nil && leadRecord != nil && leadRecord.AssignedTo != nil {
				assignedTo = strings.TrimSpace(*leadRecord.AssignedTo)
			}
		}
		if assignedTo == "" && accountRecord.AssignedTo != nil {
			assignedTo = strings.TrimSpace(*accountRecord.AssignedTo)
		}
		if assignedTo == "" {
			assignedTo = createdBy
		}

		var assignedToPtr *string
		if assignedTo != "" {
			assignedToPtr = &assignedTo
		}

		var contactIDPtr *string
		if req.ContactID != "" {
			contactIDPtr = &req.ContactID
		}

		deal := &pipeline.Deal{
			Title:             req.Title,
			Description:       req.Description,
			AccountID:         req.AccountID,
			ContactID:         contactIDPtr,
			StageID:           req.StageID,
			Value:             req.Value,
			Probability:       probability,
			ExpectedCloseDate: req.ExpectedCloseDate,
			ActualCloseDate:   actualCloseDate,
			AssignedTo:        assignedToPtr,
			LeadID:            req.LeadID,
			BrickID:           brickID,
			Status:            status,
			Source:            req.Source,
			Notes:             req.Notes,
			CreatedBy:         createdBy,
		}

		// If product_items provided, compute deal total from items and persist items.
		normalizedProductItems := normalizeCreateDealProductItemRequests(req.ProductItems)
		if len(normalizedProductItems) > 0 {
			total := int64(0)
			items := make([]pipeline.DealProductItem, 0, len(normalizedProductItems))

			for _, itemReq := range normalizedProductItems {
				if s.productRepo == nil {
					return errors.New("product repository not available")
				}
				p, err := s.productRepo.FindByID(itemReq.ProductID)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return errors.New("product not found")
					}
					return err
				}

				unitPrice := p.Price
				if itemReq.UnitPrice != nil {
					unitPrice = *itemReq.UnitPrice
				}
				discount := int64(0)
				if itemReq.DiscountAmount != nil {
					discount = *itemReq.DiscountAmount
				}
				subtotal := (unitPrice * int64(itemReq.Quantity)) - discount
				if subtotal < 0 {
					subtotal = 0
				}
				total += subtotal

				items = append(items, pipeline.DealProductItem{
					ProductID:      p.ID,
					ProductName:    p.Name,
					ProductSKU:     p.SKU,
					UnitPrice:      unitPrice,
					Quantity:       itemReq.Quantity,
					DiscountAmount: discount,
					Subtotal:       subtotal,
					Notes:          itemReq.Notes,
				})
			}

			deal.Value = total
		}

		if err := tx.Create(deal).Error; err != nil {
			return err
		}

		if len(normalizedProductItems) > 0 {
			items := make([]pipeline.DealProductItem, 0, len(normalizedProductItems))
			for _, itemReq := range normalizedProductItems {
				p, err := s.productRepo.FindByID(itemReq.ProductID)
				if err != nil {
					return err
				}
				unitPrice := p.Price
				if itemReq.UnitPrice != nil {
					unitPrice = *itemReq.UnitPrice
				}
				discount := int64(0)
				if itemReq.DiscountAmount != nil {
					discount = *itemReq.DiscountAmount
				}
				subtotal := (unitPrice * int64(itemReq.Quantity)) - discount
				if subtotal < 0 {
					subtotal = 0
				}
				items = append(items, pipeline.DealProductItem{
					DealID:         deal.ID,
					ProductID:      p.ID,
					ProductName:    p.Name,
					ProductSKU:     p.SKU,
					UnitPrice:      unitPrice,
					Quantity:       itemReq.Quantity,
					DiscountAmount: discount,
					Subtotal:       subtotal,
					Notes:          itemReq.Notes,
				})
			}

			if len(items) > 0 {
				if err := tx.Create(&items).Error; err != nil {
					return err
				}
			}
		}

		// Create initial deal history entry
		notes := "Deal created"
		if req.LeadID != nil && *req.LeadID != "" {
			notes = "Deal created from qualified lead"
			// Convert lead to "converted" status if it's qualified
			if s.leadRepo != nil {
				lead, err := s.leadRepo.FindByID(*req.LeadID)
				if err == nil && lead != nil && strings.EqualFold(lead.LeadStatus, "qualified") {
					now := time.Now()
					lead.LeadStatus = "converted"
					convertedStatus, statusErr := s.findConvertedLeadStatus()
					if statusErr != nil {
						return statusErr
					}
					if convertedStatus != nil {
						lead.LeadStatusID = &convertedStatus.ID
						lead.LeadScore = convertedStatus.Score
					}
					dealID := deal.ID
					lead.OpportunityID = &dealID
					if req.AccountID != "" {
						accountID := req.AccountID
						lead.AccountID = &accountID
					}
					if req.ContactID != "" {
						contactID := req.ContactID
						lead.ContactID = &contactID
					}
					lead.ConvertedAt = &now
					convertedBy := createdBy
					lead.ConvertedBy = &convertedBy
					if err := s.leadRepo.Update(lead); err != nil {
						return err
					}
				}
			}
			if err := s.migrateLeadAssociationsToDeal(tx, *req.LeadID, deal.ID, deal.AccountID); err != nil {
				return err
			}
		}
		if err := s.createDealHistory(deal, nil, "", 0, createdBy, "Deal creation", notes); err != nil {
			_ = err
		}

		// Reload with relations
		var reloaded pipeline.Deal
		if err := tx.
			Preload("Account").
			Preload("Contact").
			Preload("Stage").
			Preload("ProductItems").
			Preload("AssignedUser").
			Where("id = ?", deal.ID).
			First(&reloaded).Error; err != nil {
			return err
		}
		createdDeal = &reloaded
		return nil
	})
	if err != nil {
		// ERROR: Deal creation failed - logged via structured logging
		// REMOVED: Debug file I/O that caused blocking under load
		// Use centralized logging instead for production debugging
		return nil, err
	}
	if createdDeal == nil {
		return nil, errors.New("failed to create deal")
	}
	resp := createdDeal.ToDealResponse()
	_ = s.cacheService.InvalidateOnWrite(createdDeal.ID)
	_ = s.cacheService.SetDetail(createdDeal.ID, resp)
	return resp, nil
}

func (s *Service) findConvertedLeadStatus() (*lead_status.LeadStatus, error) {
	if s.db == nil {
		return nil, nil
	}

	var status lead_status.LeadStatus
	err := s.db.
		Where("code IN ? OR is_converted = ?", []string{"CONVERTED", "converted"}, true).
		Order("is_converted DESC").
		First(&status).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &status, nil
}

func (s *Service) migrateLeadAssociationsToDeal(tx *gorm.DB, leadID, dealID, accountID string) error {
	if tx == nil || leadID == "" || dealID == "" {
		return nil
	}

	updates := map[string]interface{}{
		"deal_id": dealID,
	}
	if accountID != "" {
		updates["account_id"] = accountID
	}

	for _, table := range []string{"activities", "visit_reports", "tasks"} {
		if err := tx.Table(table).
			Where("lead_id = ? AND deleted_at IS NULL", leadID).
			Updates(updates).Error; err != nil {
			return err
		}
	}

	return nil
}

// UpdateDeal updates a deal
func (s *Service) UpdateDeal(id string, req *pipeline.UpdateDealRequest, changedBy string) (*pipeline.DealResponse, error) {
	deal, err := s.dealRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDealNotFound
		}
		return nil, err
	}

	// Track if fields that affect brick_id are being updated
	accountIDChanged := false
	oldStatus := deal.Status
	oldStageID := deal.StageID
	oldProbability := deal.Probability
	oldStageName := ""
	if deal.Stage != nil {
		oldStageName = deal.Stage.Name
	}
	statusChanged := false
	stageChanged := false
	closeReason := strings.TrimSpace(req.CloseReason)

	// Update fields if provided
	if req.Title != "" {
		deal.Title = req.Title
	}
	if req.Description != "" {
		deal.Description = req.Description
	}
	if req.AccountID != "" {
		// Validate account exists
		_, err := s.accountRepo.FindByID(req.AccountID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrAccountNotFound
			}
			return nil, err
		}
		// Check if account_id is actually changing
		if deal.AccountID != req.AccountID {
			accountIDChanged = true
		}
		deal.AccountID = req.AccountID
	}
	if req.ContactID != "" {
		deal.ContactID = &req.ContactID
	}
	if req.StageID != "" {
		// Validate stage exists
		stage, err := s.pipelineRepo.FindStageByID(req.StageID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrInvalidStage
			}
			return nil, err
		}
		deal.StageID = req.StageID
		stageChanged = oldStageID != req.StageID
		// Update status based on stage
		if stage.IsWon {
			deal.Status = "won"
			now := time.Now()
			deal.ActualCloseDate = &now
			// CRM Enhancement Phase 1: Convert deal to customer purchase history when won
			s.convertDealToPurchaseHistory(deal)
		} else if stage.IsLost {
			deal.Status = "lost"
			now := time.Now()
			deal.ActualCloseDate = &now
		} else {
			deal.Status = "open"
			deal.CloseReason = ""
			deal.ActualCloseDate = nil
		}
		// Probability always follows stage percent (fallback by stage order)
		if stage.Probability > 0 {
			deal.Probability = stage.Probability
		} else {
			deal.Probability = stage.Order * 20
		}
	}
	// Update value if provided. Closed won revenue can still be overridden by
	// product items later in this method.
	if req.Value != nil {
		deal.Value = *req.Value
	}
	// Probability is derived from stage.
	if req.ExpectedCloseDate != nil {
		deal.ExpectedCloseDate = req.ExpectedCloseDate
	}
	if req.LeadID != nil {
		deal.LeadID = req.LeadID
	}
	if req.Status != "" {
		deal.Status = req.Status
		if req.Status == "open" {
			deal.ActualCloseDate = nil
			deal.CloseReason = ""
		} else if deal.ActualCloseDate == nil {
			now := time.Now()
			deal.ActualCloseDate = &now
		}
	}
	if req.Source != "" {
		deal.Source = req.Source
	}
	if req.Notes != "" {
		deal.Notes = req.Notes
	}
	if req.BudgetConfirmed != nil {
		deal.BudgetConfirmed = *req.BudgetConfirmed
	}
	if req.AuthorityConfirmed != nil {
		deal.AuthorityConfirmed = *req.AuthorityConfirmed
	}
	if req.NeedConfirmed != nil {
		deal.NeedConfirmed = *req.NeedConfirmed
	}
	if req.TimelineConfirmed != nil {
		deal.TimelineConfirmed = *req.TimelineConfirmed
	}
	if req.QualificationSnapshot != nil {
		snapshot, marshalErr := json.Marshal(req.QualificationSnapshot)
		if marshalErr != nil {
			return nil, marshalErr
		}
		deal.QualificationSnapshot = datatypes.JSON(snapshot)
	}

	statusChanged = oldStatus != deal.Status
	if statusChanged && isClosedDealStatus(deal.Status) {
		if closeReason == "" {
			return nil, ErrCloseReasonRequired
		}
		deal.CloseReason = closeReason
	} else if closeReason != "" {
		deal.CloseReason = closeReason
	}

	// Auto-update brick_id if account_id changed. Assignee is managed by create flow.
	if accountIDChanged && s.brickHelper != nil {
		var brickID *string
		// Account territory is the source of truth for deal brick assignment.
		if deal.AccountID != "" {
			brickID, _ = s.brickHelper.GetBrickIDFromAccount(deal.AccountID)
		}
		if brickID == nil && deal.AssignedTo != nil && *deal.AssignedTo != "" {
			brickID, _ = s.brickHelper.GetBrickIDFromUser(*deal.AssignedTo)
		}
		deal.BrickID = brickID
	}

	// Handle product items update if provided
	normalizedProductItems := normalizeCreateDealProductItemRequests(req.ProductItems)
	if len(normalizedProductItems) > 0 {
		// Delete existing product items for this deal
		if err := s.dealProductItemRepo.DeleteByDealID(deal.ID); err != nil {
			return nil, err
		}

		// Validate and create new product items
		total := int64(0)
		items := make([]pipeline.DealProductItem, 0, len(normalizedProductItems))

		for _, itemReq := range normalizedProductItems {
			if s.productRepo == nil {
				return nil, errors.New("product repository not available")
			}
			p, err := s.productRepo.FindByID(itemReq.ProductID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, errors.New("product not found")
				}
				return nil, err
			}

			unitPrice := p.Price
			if itemReq.UnitPrice != nil {
				unitPrice = *itemReq.UnitPrice
			}
			discount := int64(0)
			if itemReq.DiscountAmount != nil {
				discount = *itemReq.DiscountAmount
			}
			subtotal := (unitPrice * int64(itemReq.Quantity)) - discount
			if subtotal < 0 {
				subtotal = 0
			}
			total += subtotal

			items = append(items, pipeline.DealProductItem{
				DealID:         deal.ID,
				ProductID:      p.ID,
				ProductName:    p.Name,
				ProductSKU:     p.SKU,
				UnitPrice:      unitPrice,
				Quantity:       itemReq.Quantity,
				DiscountAmount: discount,
				Subtotal:       subtotal,
				Notes:          itemReq.Notes,
			})
		}

		// Create new product items in batch
		if err := s.dealProductItemRepo.CreateMany(items); err != nil {
			return nil, err
		}

		// Update deal value with calculated total
		deal.Value = total
	}

	if err := s.dealRepo.Update(deal); err != nil {
		return nil, err
	}

	if stageChanged && s.dealHistoryRepo != nil && changedBy != "" {
		fromStageID := oldStageID
		notes := "Stage moved from " + oldStageName
		if deal.StageID != "" {
			if stageEntity, stageErr := s.pipelineRepo.FindStageByID(deal.StageID); stageErr == nil {
				notes += " to " + stageEntity.Name
			}
		}
		_ = s.createDealHistory(deal, &fromStageID, oldStageName, oldProbability, changedBy, closeReason, notes)
	}

	s.emitDealStatusEvents(deal, oldStageID, oldStageName, oldStatus, changedBy)

	// Reload to get relations
	deal, err = s.dealRepo.FindByID(deal.ID)
	if err != nil {
		return nil, err
	}

	resp := deal.ToDealResponse()
	_ = s.cacheService.InvalidateOnWrite(deal.ID)
	_ = s.cacheService.SetDetail(deal.ID, resp)
	return resp, nil
}

// MoveDeal moves a deal to a different stage
func (s *Service) MoveDeal(id string, req *pipeline.MoveDealRequest) (*pipeline.DealResponse, error) {
	deal, err := s.dealRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDealNotFound
		}
		return nil, err
	}

	// Validate new stage exists
	stage, err := s.pipelineRepo.FindStageByID(req.StageID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidStage
		}
		return nil, err
	}

	deal.StageID = req.StageID
	// Update probability from stage if set, otherwise fallback to formula
	if stage.Probability > 0 {
		deal.Probability = stage.Probability
	} else {
		deal.Probability = stage.Order * 20 // Fallback formula: each stage = +20%
	}
	// Update status based on stage
	if stage.IsWon {
		deal.Status = "won"
		now := time.Now()
		deal.ActualCloseDate = &now
		// CRM Enhancement Phase 1: Convert deal to customer purchase history when won
		s.convertDealToPurchaseHistory(deal)
	} else if stage.IsLost {
		deal.Status = "lost"
		now := time.Now()
		deal.ActualCloseDate = &now
	} else {
		deal.Status = "open"
		deal.ActualCloseDate = nil
	}

	if err := s.dealRepo.Update(deal); err != nil {
		return nil, err
	}

	// Reload to get relations
	deal, err = s.dealRepo.FindByID(deal.ID)
	if err != nil {
		return nil, err
	}

	resp := deal.ToDealResponse()
	_ = s.cacheService.InvalidateOnWrite(deal.ID)
	_ = s.cacheService.SetDetail(deal.ID, resp)
	return resp, nil
}

// DeleteDeal deletes a deal
func (s *Service) DeleteDeal(id string) error {
	deal, err := s.dealRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrDealNotFound
		}
		return err
	}

	err = s.dealRepo.Delete(deal.ID)
	if err != nil {
		return err
	}
	_ = s.cacheService.InvalidateOnWrite(deal.ID)
	return nil
}

// GetSummary returns pipeline summary
func (s *Service) GetSummary() (*pipeline.PipelineSummaryResponse, error) {
	var cachedSummary pipeline.PipelineSummaryResponse
	if found, _ := s.cacheService.GetSummary("", &cachedSummary); found {
		return &cachedSummary, nil
	}

	summary, err := s.dealRepo.GetSummary()
	if err != nil {
		return nil, err
	}
	_ = s.cacheService.SetSummary("", summary)
	return summary, nil
}

// GetForecast returns forecast data
func (s *Service) GetForecast(periodType string) (*pipeline.ForecastResponse, error) {
	now := time.Now()
	var start, end time.Time

	switch periodType {
	case "month":
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 1, 0).Add(-time.Second)
	case "quarter":
		quarter := (now.Month() - 1) / 3
		start = time.Date(now.Year(), quarter*3+1, 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 3, 0).Add(-time.Second)
	case "year":
		start = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(1, 0, 0).Add(-time.Second)
	default:
		// Default to current month
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 1, 0).Add(-time.Second)
		periodType = "month"
	}

	return s.dealRepo.GetForecast(periodType, start, end)
}

// CreateStage creates a new pipeline stage
func (s *Service) CreateStage(req *pipeline.CreateStageRequest) (*pipeline.PipelineStageResponse, error) {
	// Check if code already exists
	existing, err := s.pipelineRepo.FindStageByCode(req.Code)
	if err == nil && existing != nil {
		return nil, errors.New("pipeline stage with this code already exists")
	}

	// Set default color if not provided
	color := req.Color
	if color == "" {
		color = "#3B82F6"
	}

	stage := &pipeline.PipelineStage{
		Name:        req.Name,
		Code:        req.Code,
		Order:       req.Order,
		Color:       color,
		IsActive:    req.IsActive,
		IsWon:       req.IsWon,
		IsLost:      req.IsLost,
		Probability: req.Probability,
		Description: req.Description,
	}

	if err := s.pipelineRepo.CreateStage(stage); err != nil {
		return nil, err
	}

	// Reload to get relations
	stage, err = s.pipelineRepo.FindStageByID(stage.ID)
	if err != nil {
		return nil, err
	}
	_ = s.cacheService.InvalidateOnWrite("")
	return stage.ToPipelineStageResponse(), nil
}

// UpdateStage updates a pipeline stage
func (s *Service) UpdateStage(id string, req *pipeline.UpdateStageRequest) (*pipeline.PipelineStageResponse, error) {
	stage, err := s.pipelineRepo.FindStageByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPipelineStageNotFound
		}
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		stage.Name = req.Name
	}
	if req.Code != "" {
		// Check if new code already exists (excluding current stage)
		existing, err := s.pipelineRepo.FindStageByCode(req.Code)
		if err == nil && existing != nil && existing.ID != id {
			return nil, errors.New("pipeline stage with this code already exists")
		}
		stage.Code = req.Code
	}
	if req.Order != nil {
		stage.Order = *req.Order
	}
	if req.Color != "" {
		stage.Color = req.Color
	}
	if req.IsActive != nil {
		stage.IsActive = *req.IsActive
	}
	if req.IsWon != nil {
		stage.IsWon = *req.IsWon
	}
	if req.IsLost != nil {
		stage.IsLost = *req.IsLost
	}
	if req.Probability != nil {
		stage.Probability = *req.Probability
	}
	if req.Description != "" {
		stage.Description = req.Description
	}

	if err := s.pipelineRepo.UpdateStage(stage); err != nil {
		return nil, err
	}

	// Reload to get relations
	stage, err = s.pipelineRepo.FindStageByID(stage.ID)
	if err != nil {
		return nil, err
	}
	_ = s.cacheService.InvalidateOnWrite("")
	return stage.ToPipelineStageResponse(), nil
}

// DeleteStage deletes a pipeline stage
func (s *Service) DeleteStage(id string) error {
	stage, err := s.pipelineRepo.FindStageByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPipelineStageNotFound
		}
		return err
	}

	// Check if stage is being used by any deals
	// Get deals count for this stage
	dealsWithStage, _, err := s.dealRepo.List(&pipeline.ListDealsRequest{
		StageID: stage.ID,
		Page:    1,
		PerPage: 1, // Just need to check if any exists
	})
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if len(dealsWithStage) > 0 {
		return errors.New("cannot delete stage: stage is being used by existing deals")
	}

	err = s.pipelineRepo.DeleteStage(stage.ID)
	if err != nil {
		return err
	}
	_ = s.cacheService.InvalidateOnWrite("")
	return nil
}

// UpdateStagesOrder updates the order of multiple stages
func (s *Service) UpdateStagesOrder(req *pipeline.UpdateStagesOrderRequest) ([]pipeline.PipelineStageResponse, error) {
	// Update each stage's order
	for _, item := range req.Stages {
		stage, err := s.pipelineRepo.FindStageByID(item.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrPipelineStageNotFound
			}
			return nil, err
		}

		stage.Order = item.Order
		if err := s.pipelineRepo.UpdateStage(stage); err != nil {
			return nil, err
		}
	}

	// Return updated list of stages
	listReq := &pipeline.ListPipelineStagesRequest{}
	stages, err := s.pipelineRepo.ListStages(listReq)
	if err != nil {
		return nil, err
	}

	responses := make([]pipeline.PipelineStageResponse, len(stages))
	for i, stage := range stages {
		responses[i] = *stage.ToPipelineStageResponse()
	}
	_ = s.cacheService.InvalidateOnWrite("")
	return responses, nil
}

// PaginationResult represents pagination information
type PaginationResult struct {
	Page       int
	PerPage    int
	Total      int
	TotalPages int
}

// createDealHistory creates a deal history entry
func (s *Service) createDealHistory(deal *pipeline.Deal, fromStageID *string, fromStageName string, fromProbability int, changedBy string, reason string, notes string) error {
	if s.dealHistoryRepo == nil {
		return nil // Skip if repository not available
	}

	// Get current stage info
	currentStage, err := s.pipelineRepo.FindStageByID(deal.StageID)
	if err != nil {
		return err
	}

	// Calculate days in previous stage if fromStageID is provided
	var daysInPrevStage *int
	if fromStageID != nil {
		// Query deal histories to find when deal entered previous stage
		histories, err := s.dealHistoryRepo.FindByDealID(deal.ID)
		if err == nil && len(histories) > 0 {
			// Find the entry when deal moved to fromStage
			for i := len(histories) - 1; i >= 0; i-- {
				if histories[i].ToStageID == *fromStageID {
					days := int(time.Since(histories[i].ChangedAt).Hours() / 24)
					daysInPrevStage = &days
					break
				}
			}
		}
	}

	history := &deal_history.DealHistory{
		DealID:          deal.ID,
		FromStageID:     fromStageID,
		FromStageName:   fromStageName,
		ToStageID:       deal.StageID,
		ToStageName:     currentStage.Name,
		FromProbability: fromProbability,
		ToProbability:   deal.Probability,
		DaysInPrevStage: daysInPrevStage,
		ChangedBy:       changedBy,
		ChangedAt:       time.Now(),
		Reason:          reason,
		Notes:           notes,
	}

	return s.dealHistoryRepo.Create(history)
}

// ValidateStageRequirements validates if a deal can move to another stage.
// Business rules:
//  1. Backward movement is allowed as a correction in the Kanban pipeline.
//  2. A deal MUST have at least one product item to move to a "won" stage.
//  3. Deal value must be > 0 before moving past proposal.
func (s *Service) ValidateStageRequirements(dealID string, toStageID string, incomingProductItems ...[]pipeline.CreateDealProductItemRequest) error {
	// Get deal with products preloaded
	deal, err := s.dealRepo.FindByID(dealID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrDealNotFound
		}
		return err
	}

	// Get target stage
	toStage, err := s.pipelineRepo.FindStageByID(toStageID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPipelineStageNotFound
		}
		return err
	}

	// Get current stage
	currentStage, err := s.pipelineRepo.FindStageByID(deal.StageID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPipelineStageNotFound
		}
		return err
	}

	// Same-stage move is a no-op, always allowed
	if currentStage.ID == toStage.ID {
		return nil
	}

	// Rule: Moving to a lost stage is always permitted (cancel scenario)
	if toStage.IsLost {
		return nil
	}

	// Rule: Closed won revenue must come from explicit sold products, while
	// intermediate stages can still carry a manual estimated deal value.
	if toStage.IsWon {
		validProductCount := 0
		if len(incomingProductItems) > 0 {
			for _, item := range incomingProductItems[0] {
				if item.ProductID != "" && item.Quantity > 0 {
					validProductCount++
				}
			}
		}
		for _, item := range deal.ProductItems {
			if item.ProductID != "" && item.Quantity > 0 {
				validProductCount++
			}
		}
		if validProductCount == 0 {
			return errors.New("at least one sold product is required before moving to closed won")
		}
	}
	if toStage.Order >= 3 && !toStage.IsWon {
		computedValue := deal.Value
		itemValue := int64(0)
		if len(incomingProductItems) > 0 {
			for _, item := range incomingProductItems[0] {
				if item.Quantity <= 0 {
					continue
				}
				unitPrice := int64(0)
				if item.UnitPrice != nil {
					unitPrice = *item.UnitPrice
				}
				discount := int64(0)
				if item.DiscountAmount != nil {
					discount = *item.DiscountAmount
				}
				subtotal := unitPrice*int64(item.Quantity) - discount
				if subtotal > 0 {
					itemValue += subtotal
				}
			}
		}
		for _, item := range deal.ProductItems {
			itemValue += item.Subtotal
		}
		if itemValue > 0 {
			computedValue = itemValue
		}
		if computedValue == 0 {
			return errors.New("deal value must be greater than 0 before moving to this stage")
		}
	}

	return nil
}

// MoveStageWithValidation moves a deal to a new stage with validation and history logging
func (s *Service) MoveStageWithValidation(dealID string, toStageID string, changedBy string, reason string, productItems ...[]pipeline.CreateDealProductItemRequest) (*pipeline.DealResponse, error) {
	// Validate stage requirements
	var incomingProductItems []pipeline.CreateDealProductItemRequest
	if len(productItems) > 0 {
		incomingProductItems = productItems[0]
	}
	if err := s.ValidateStageRequirements(dealID, toStageID, incomingProductItems); err != nil {
		if err == ErrDealNotFound {
			return nil, ErrDealNotFound
		}
		if err == ErrPipelineStageNotFound {
			return nil, ErrPipelineStageNotFound
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPipelineStageNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrStageRequirementsNotMet, err)
	}

	// Get current deal state
	deal, err := s.dealRepo.FindByID(dealID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDealNotFound
		}
		return nil, err
	}

	// Get current stage info for history
	currentStage, err := s.pipelineRepo.FindStageByID(deal.StageID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPipelineStageNotFound
		}
		return nil, err
	}

	fromStageID := deal.StageID
	fromStageName := currentStage.Name
	fromProbability := deal.Probability

	// Get new stage info
	newStage, err := s.pipelineRepo.FindStageByID(toStageID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPipelineStageNotFound
		}
		return nil, err
	}

	// Update deal
	deal.StageID = toStageID
	// Use probability from stage if set, otherwise fallback to formula based on order
	if newStage.Probability > 0 {
		deal.Probability = newStage.Probability
	} else {
		deal.Probability = newStage.Order * 20 // Fallback formula: each stage = +20%
	}

	// Update status based on stage flags
	if newStage.IsWon {
		if strings.TrimSpace(reason) == "" {
			return nil, ErrCloseReasonRequired
		}
		if len(incomingProductItems) > 0 {
			if err := s.replaceDealProductItems(deal, incomingProductItems); err != nil {
				return nil, err
			}
		}
		deal.Status = "won"
		now := time.Now()
		deal.ActualCloseDate = &now
		deal.CloseReason = strings.TrimSpace(reason)
	} else if newStage.IsLost {
		if strings.TrimSpace(reason) == "" {
			return nil, ErrCloseReasonRequired
		}
		deal.Status = "lost"
		now := time.Now()
		deal.ActualCloseDate = &now
		deal.CloseReason = strings.TrimSpace(reason)
	} else {
		deal.Status = "open"
		deal.ActualCloseDate = nil
		deal.CloseReason = ""
	}

	if err := s.dealRepo.Update(deal); err != nil {
		return nil, err
	}

	// Create deal history
	fromStageIDPtr := &fromStageID
	notes := "Stage moved from " + fromStageName + " to " + newStage.Name
	if err := s.createDealHistory(deal, fromStageIDPtr, fromStageName, fromProbability, changedBy, reason, notes); err != nil {
		// Log error but don't fail the operation
		_ = err
	}

	// Auto-create stage-specific tasks
	if s.taskRepo != nil {
		s.createStageTransitionTasks(deal, newStage, changedBy)
	}

	// If deal is now won, generate customer purchase history
	if newStage.IsWon {
		s.convertDealToPurchaseHistory(deal)
	}

	s.emitDealStatusEvents(deal, fromStageID, fromStageName, currentStageStatus(currentStage), changedBy)

	// Reload and return
	_ = s.cacheService.InvalidateOnWrite(dealID)
	return s.GetDealByID(dealID)
}

func isClosedDealStatus(status string) bool {
	switch status {
	case "won", "lost":
		return true
	default:
		return false
	}
}

func currentStageStatus(stage *pipeline.PipelineStage) string {
	if stage == nil {
		return "open"
	}
	if stage.IsWon {
		return "won"
	}
	if stage.IsLost {
		return "lost"
	}
	return "open"
}

func (s *Service) replaceDealProductItems(deal *pipeline.Deal, productItems []pipeline.CreateDealProductItemRequest) error {
	productItems = normalizeCreateDealProductItemRequests(productItems)
	if deal == nil || len(productItems) == 0 {
		return nil
	}
	if s.productRepo == nil || s.dealProductItemRepo == nil {
		return errors.New("product repository not available")
	}

	if err := s.dealProductItemRepo.DeleteByDealID(deal.ID); err != nil {
		return err
	}

	total := int64(0)
	items := make([]pipeline.DealProductItem, 0, len(productItems))
	for _, itemReq := range productItems {
		if itemReq.ProductID == "" || itemReq.Quantity < 1 {
			continue
		}

		p, err := s.productRepo.FindByID(itemReq.ProductID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("product not found")
			}
			return err
		}

		unitPrice := p.Price
		if itemReq.UnitPrice != nil {
			unitPrice = *itemReq.UnitPrice
		}
		discount := int64(0)
		if itemReq.DiscountAmount != nil {
			discount = *itemReq.DiscountAmount
		}
		subtotal := unitPrice*int64(itemReq.Quantity) - discount
		if subtotal < 0 {
			subtotal = 0
		}
		total += subtotal

		items = append(items, pipeline.DealProductItem{
			DealID:         deal.ID,
			ProductID:      p.ID,
			ProductName:    p.Name,
			ProductSKU:     p.SKU,
			UnitPrice:      unitPrice,
			Quantity:       itemReq.Quantity,
			DiscountAmount: discount,
			Subtotal:       subtotal,
			Notes:          itemReq.Notes,
		})
	}

	if len(items) > 0 {
		if err := s.dealProductItemRepo.CreateMany(items); err != nil {
			return err
		}
		deal.ProductItems = items
		deal.Value = total
	}
	return nil
}

func normalizeCreateDealProductItemRequests(productItems []pipeline.CreateDealProductItemRequest) []pipeline.CreateDealProductItemRequest {
	normalized := make([]pipeline.CreateDealProductItemRequest, 0, len(productItems))
	indexByProductID := make(map[string]int, len(productItems))

	for _, item := range productItems {
		item.ProductID = strings.TrimSpace(item.ProductID)
		if item.ProductID == "" || item.Quantity < 1 {
			continue
		}

		if existingIndex, exists := indexByProductID[item.ProductID]; exists {
			existing := &normalized[existingIndex]
			existing.Quantity += item.Quantity
			if item.UnitPrice != nil {
				existing.UnitPrice = item.UnitPrice
			}
			if item.DiscountAmount != nil {
				if existing.DiscountAmount == nil {
					discount := int64(0)
					existing.DiscountAmount = &discount
				}
				*existing.DiscountAmount += *item.DiscountAmount
			}
			if strings.TrimSpace(item.Notes) != "" {
				if strings.TrimSpace(existing.Notes) == "" {
					existing.Notes = strings.TrimSpace(item.Notes)
				} else {
					existing.Notes = existing.Notes + "; " + strings.TrimSpace(item.Notes)
				}
			}
			continue
		}

		item.Notes = strings.TrimSpace(item.Notes)
		indexByProductID[item.ProductID] = len(normalized)
		normalized = append(normalized, item)
	}

	return normalized
}

func (s *Service) emitDealStatusEvents(deal *pipeline.Deal, oldStageID, oldStageName, oldStatus, changedBy string) {
	if s.eventHelper == nil || changedBy == "" || deal == nil {
		return
	}

	if oldStageID != "" && oldStageID != deal.StageID {
		newStageName := oldStageName
		if deal.Stage != nil && deal.Stage.Name != "" {
			newStageName = deal.Stage.Name
		} else if stage, err := s.pipelineRepo.FindStageByID(deal.StageID); err == nil {
			newStageName = stage.Name
		}

		s.eventHelper.EmitDealStageChanged(&domainevents.DealStageChangedEvent{
			DealID:       deal.ID,
			OldStageID:   oldStageID,
			OldStageName: oldStageName,
			NewStageID:   deal.StageID,
			NewStageName: newStageName,
			ChangedBy:    changedBy,
			ChangedAt:    time.Now(),
		}, changedBy)
	}

	assignedTo := ""
	if deal.AssignedTo != nil {
		assignedTo = *deal.AssignedTo
	}

	if oldStatus != "won" && deal.Status == "won" && deal.ActualCloseDate != nil {
		s.eventHelper.EmitDealWon(&domainevents.DealWonEvent{
			DealID:          deal.ID,
			Title:           deal.Title,
			Value:           deal.Value,
			AccountID:       deal.AccountID,
			AssignedTo:      assignedTo,
			ActualCloseDate: *deal.ActualCloseDate,
			WonBy:           changedBy,
			WonAt:           time.Now(),
		}, changedBy)
	}

	if oldStatus != "lost" && deal.Status == "lost" && deal.ActualCloseDate != nil {
		s.eventHelper.EmitDealLost(&domainevents.DealLostEvent{
			DealID:          deal.ID,
			Title:           deal.Title,
			Value:           deal.Value,
			AccountID:       deal.AccountID,
			AssignedTo:      assignedTo,
			ActualCloseDate: *deal.ActualCloseDate,
			LostReason:      deal.CloseReason,
			LostBy:          changedBy,
			LostAt:          time.Now(),
		}, changedBy)
	}
}

// createStageTransitionTasks creates automatic tasks when deal moves to a new stage
func (s *Service) createStageTransitionTasks(deal *pipeline.Deal, stage *pipeline.PipelineStage, createdBy string) {
	if s.taskRepo == nil {
		return
	}

	var tasksToCreate []struct {
		title       string
		description string
		daysOffset  int
		priority    string
		taskType    string
	}

	// Define tasks based on stage
	switch stage.Code {
	case "proposal":
		tasksToCreate = []struct {
			title       string
			description string
			daysOffset  int
			priority    string
			taskType    string
		}{
			{
				title:       "Prepare detailed quotation",
				description: "Create comprehensive quotation with all line items",
				daysOffset:  1,
				priority:    "high",
				taskType:    "general",
			},
			{
				title:       "Schedule proposal presentation",
				description: "Arrange meeting with decision makers to present proposal",
				daysOffset:  3,
				priority:    "high",
				taskType:    "meeting",
			},
			{
				title:       "Get pricing approval from manager",
				description: "Internal pricing review and approval",
				daysOffset:  2,
				priority:    "medium",
				taskType:    "general",
			},
		}
	case "negotiation":
		tasksToCreate = []struct {
			title       string
			description string
			daysOffset  int
			priority    string
			taskType    string
		}{
			{
				title:       "Send contract for legal review",
				description: "Submit contract to legal team for review",
				daysOffset:  1,
				priority:    "high",
				taskType:    "general",
			},
			{
				title:       "Schedule contract negotiation meeting",
				description: "Arrange meeting to discuss and finalize contract terms",
				daysOffset:  3,
				priority:    "high",
				taskType:    "meeting",
			},
		}
	}

	// Create tasks
	for _, taskDef := range tasksToCreate {
		dueDate := time.Now().AddDate(0, 0, taskDef.daysOffset)

		newTask := &task.Task{
			Title:       taskDef.title,
			Description: taskDef.description,
			Type:        taskDef.taskType,
			Status:      "pending",
			Priority:    taskDef.priority,
			DueDate:     &dueDate,
			DealID:      &deal.ID,
			AccountID:   &deal.AccountID,
			CreatedBy:   createdBy,
		}

		newTask.AssignedTo = deal.AssignedTo

		_ = s.taskRepo.Create(newTask) // Ignore error for now
	}
}

// GetDealHistory returns the history of stage changes for a deal
func (s *Service) GetDealHistory(dealID string) ([]deal_history.DealHistoryResponse, error) {
	if s.dealHistoryRepo == nil {
		return []deal_history.DealHistoryResponse{}, nil
	}

	histories, err := s.dealHistoryRepo.FindByDealID(dealID)
	if err != nil {
		return nil, err
	}

	responses := make([]deal_history.DealHistoryResponse, len(histories))
	for i, h := range histories {
		responses[i] = *h.ToDealHistoryResponse()
	}

	return responses, nil
}

// convertDealToPurchaseHistory creates a purchase record from a won deal
func (s *Service) convertDealToPurchaseHistory(deal *pipeline.Deal) {
	if s.customerPurchaseRepo == nil {
		return
	}

	// Create history record
	history := &customer_purchase.CustomerPurchaseHistory{
		AccountID:    deal.AccountID,
		DealID:       deal.ID,
		PurchaseDate: time.Now(),
		TotalAmount:  deal.Value,
		TotalItems:   len(deal.ProductItems),
		SourceLeadID: deal.LeadID,
		SourceType:   "pipeline",
	}

	// Set Sales Rep info if assigned
	if deal.AssignedTo != nil {
		history.SalesRepID = deal.AssignedTo
		if deal.AssignedUser != nil {
			history.SalesRepName = deal.AssignedUser.Name
		}
	}

	// Populate products from deal items
	items := make([]customer_purchase.PurchaseProduct, 0, len(deal.ProductItems))
	for _, item := range deal.ProductItems {
		items = append(items, customer_purchase.PurchaseProduct{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			ProductSKU:  item.ProductSKU,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Subtotal:    item.Subtotal,
		})
	}

	data, _ := json.Marshal(items)
	history.Products = datatypes.JSON(data)

	_ = s.customerPurchaseRepo.Create(history)
}
