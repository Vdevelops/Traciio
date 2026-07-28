package lead

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/config"
	"github.com/gilabs/crm-healthcare/api/internal/domain/account"
	"github.com/gilabs/crm-healthcare/api/internal/domain/activity"
	domainauth "github.com/gilabs/crm-healthcare/api/internal/domain/auth"
	"github.com/gilabs/crm-healthcare/api/internal/domain/contact"
	"github.com/gilabs/crm-healthcare/api/internal/domain/deal_history"
	domainevents "github.com/gilabs/crm-healthcare/api/internal/domain/events"
	"github.com/gilabs/crm-healthcare/api/internal/domain/industry"
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead"
	"github.com/gilabs/crm-healthcare/api/internal/domain/task"
	leadqualification "github.com/gilabs/crm-healthcare/api/internal/domain/lead_qualification"
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead_source"
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead_status"
	"github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	brickservice "github.com/gilabs/crm-healthcare/api/internal/service/brick"
	"github.com/gilabs/crm-healthcare/api/pkg/cache"
	"github.com/gilabs/crm-healthcare/api/pkg/geocoding"
	"gorm.io/gorm"
)

var (
	ErrLeadNotFound              = errors.New("lead not found")
	ErrLeadAlreadyConverted      = errors.New("lead already converted")
	ErrLeadCannotConvert         = errors.New("lead cannot convert")
	ErrInvalidLeadStatus         = errors.New("invalid lead status")
	ErrLeadStatusReasonRequired  = errors.New("lead status reason is required")
	ErrInvalidLeadSource         = errors.New("invalid lead source")
	ErrStageNotFound             = errors.New("stage not found")
	ErrInvalidConversionStage    = errors.New("conversion stage must be closed won")
	ErrSoldProductsRequired      = errors.New("at least one sold product is required")
	ErrAccountCreationFailed     = errors.New("account creation failed")
	ErrContactCreationFailed     = errors.New("contact creation failed")
	ErrOpportunityCreationFailed = errors.New("opportunity creation failed")
)

type Service struct {
	db               *gorm.DB
	leadRepo         interfaces.LeadRepository
	dealRepo         interfaces.DealRepository
	pipelineRepo     interfaces.PipelineRepository
	accountRepo      interfaces.AccountRepository
	contactRepo      interfaces.ContactRepository
	categoryRepo     interfaces.CategoryRepository
	contactRoleRepo  interfaces.ContactRoleRepository
	userRepo         interfaces.UserRepository
	activityRepo     interfaces.ActivityRepository    // For auto-migrate activities
	visitReportRepo  interfaces.VisitReportRepository // For auto-migrate visit reports
	taskRepo         interfaces.TaskRepository        // For auto-migrate tasks
	dealHistoryRepo  interfaces.DealHistoryRepository // For deal history logging
	leadStatusRepo   interfaces.LeadStatusRepository  // For lead status lookup
	brickHelper      *brickservice.BrickHelper
	geocodingSvc     *geocoding.GeocodingService
	geocodingEnabled bool
	cacheService     *cache.LeadCacheService
	eventHelper      *domainevents.Helper // For emitting domain events
}

func NewService(
	db *gorm.DB,
	leadRepo interfaces.LeadRepository,
	dealRepo interfaces.DealRepository,
	pipelineRepo interfaces.PipelineRepository,
	accountRepo interfaces.AccountRepository,
	contactRepo interfaces.ContactRepository,
	categoryRepo interfaces.CategoryRepository,
	contactRoleRepo interfaces.ContactRoleRepository,
	userRepo interfaces.UserRepository,
	activityRepo interfaces.ActivityRepository,
	visitReportRepo interfaces.VisitReportRepository,
	taskRepo interfaces.TaskRepository,
	dealHistoryRepo interfaces.DealHistoryRepository,
	leadStatusRepo interfaces.LeadStatusRepository,
	brickHelper *brickservice.BrickHelper,
	eventHelper *domainevents.Helper,
) *Service {
	var geocodingSvc *geocoding.GeocodingService
	geocodingEnabled := false

	if config.AppConfig != nil && config.AppConfig.Geocoding.Enabled {
		geocodingSvc = geocoding.NewGeocodingService(
			config.AppConfig.Geocoding.Provider,
			config.AppConfig.Geocoding.APIKey,
		)
		geocodingEnabled = true
	}

	return &Service{
		db:               db,
		leadRepo:         leadRepo,
		dealRepo:         dealRepo,
		pipelineRepo:     pipelineRepo,
		accountRepo:      accountRepo,
		contactRepo:      contactRepo,
		categoryRepo:     categoryRepo,
		contactRoleRepo:  contactRoleRepo,
		userRepo:         userRepo,
		activityRepo:     activityRepo,
		visitReportRepo:  visitReportRepo,
		taskRepo:         taskRepo,
		dealHistoryRepo:  dealHistoryRepo,
		leadStatusRepo:   leadStatusRepo,
		brickHelper:      brickHelper,
		geocodingSvc:     geocodingSvc,
		geocodingEnabled: geocodingEnabled,
		cacheService:     cache.NewLeadCacheService(nil),
		eventHelper:      eventHelper,
	}
}

func (s *Service) resolveBrickIDFromLead(l *lead.Lead) (*string, error) {
	if s.brickHelper == nil || l == nil {
		return nil, nil
	}

	if strings.TrimSpace(l.Province) != "" && strings.TrimSpace(l.City) != "" {
		brickID, err := s.brickHelper.EnsureBrickIDForLocation(l.Province, l.City)
		if err != nil {
			log.Printf("Warning: Failed to ensure brick ID for location (%s, %s): %v", l.Province, l.City, err)
		} else if brickID != nil {
			return brickID, nil
		}
	}

	if l.AssignedTo != nil && *l.AssignedTo != "" {
		if brickID, err := s.brickHelper.GetBrickIDFromUser(*l.AssignedTo); err == nil && brickID != nil {
			return brickID, nil
		} else if err != nil {
			log.Printf("Warning: Failed to get brick ID from assigned user %s: %v", *l.AssignedTo, err)
		}
	}

	return nil, nil
}

// PaginationResult represents pagination information
type PaginationResult struct {
	Page       int
	PerPage    int
	Total      int
	TotalPages int
}

// cachedLeadListResult represents cached lead list data
type cachedLeadListResult struct {
	Leads []lead.LeadResponse `msgpack:"leads"`
	Total int64               `msgpack:"total"`
}

// List returns a list of leads with pagination
func (s *Service) List(req *lead.ListLeadsRequest) ([]lead.LeadResponse, *PaginationResult, error) {
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

	// Build filter map for cache key (includes ScopedUserIDs to prevent cross-scope cache pollution)
	filters := map[string]interface{}{
		"search":      req.Search,
		"status":      req.Status,
		"assigned_to": req.AssignedTo,
		"source":      req.Source,
		"start_date":  req.StartDate,
		"end_date":    req.EndDate,
	}
	if len(req.ScopedUserIDs) > 0 {
		filters["scoped_user_ids"] = strings.Join(req.ScopedUserIDs, ",")
	}

	// Try to get from cache first
	var cachedResult cachedLeadListResult
	if found, _ := s.cacheService.GetList(page, perPage, filters, &cachedResult); found {
		totalPages := int((cachedResult.Total + int64(perPage) - 1) / int64(perPage))
		pagination := &PaginationResult{
			Page:       page,
			PerPage:    perPage,
			Total:      int(cachedResult.Total),
			TotalPages: totalPages,
		}
		return cachedResult.Leads, pagination, nil
	}

	leads, total, err := s.leadRepo.List(req)
	if err != nil {
		return nil, nil, err
	}

	responses := make([]lead.LeadResponse, len(leads))
	for i, l := range leads {
		responses[i] = *l.ToLeadResponse()
	}

	pagination := &PaginationResult{
		Page:       page,
		PerPage:    perPage,
		Total:      int(total),
		TotalPages: int((total + int64(perPage) - 1) / int64(perPage)),
	}

	// Cache the result
	_ = s.cacheService.SetList(page, perPage, filters, cachedLeadListResult{
		Leads: responses,
		Total: total,
	})

	return responses, pagination, nil
}

// GetByID returns a lead by ID
func (s *Service) GetByID(id string) (*lead.LeadResponse, error) {
	// Try to get from cache first
	var cachedLead lead.LeadResponse
	if found, _ := s.cacheService.GetDetail(id, &cachedLead); found {
		return &cachedLead, nil
	}

	l, err := s.leadRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLeadNotFound
		}
		return nil, err
	}

	response := l.ToLeadResponse()

	// Cache the result
	_ = s.cacheService.SetDetail(id, response)

	return response, nil
}

// Create creates a new lead
func (s *Service) Create(req *lead.CreateLeadRequest, createdBy string, currentUser *domainauth.UserContext) (*lead.LeadResponse, error) {
	// Resolve lead status from lead_statuses table to keep create lead consistent with Lead Status management.
	// Preferred input: lead_status_id. Backward compatible: lead_status (string).
	var resolvedLeadStatusID *string
	resolvedLegacyLeadStatus := ""
	resolvedLeadScore := req.LeadScore

	resolveByID := func(id string) (*lead_status.LeadStatus, error) {
		return s.leadStatusRepo.FindByID(id)
	}

	resolveByCodeOrLegacy := func(codeOrLegacy string) (*lead_status.LeadStatus, error) {
		if codeOrLegacy == "" {
			return nil, nil
		}
		return s.leadStatusRepo.FindByCode(codeOrLegacy)
	}

	getDefault := func() (*lead_status.LeadStatus, error) {
		return s.leadStatusRepo.FindDefault()
	}

	var chosen *lead_status.LeadStatus
	var err error
	if req.LeadStatusID != "" {
		chosen, err = resolveByID(req.LeadStatusID)
		if err != nil {
			return nil, err
		}
	}
	if chosen == nil && req.LeadStatus != "" {
		chosen, err = resolveByCodeOrLegacy(req.LeadStatus)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	if chosen == nil {
		chosen, err = getDefault()
		if err != nil || chosen == nil {
			// Fallback if no default exists
			resolvedLegacyLeadStatus = "new"
		} else {
			resolvedLeadStatusID = &chosen.ID
			resolvedLegacyLeadStatus = strings.ToLower(chosen.Code)
			if resolvedLeadScore == 0 {
				resolvedLeadScore = chosen.Score
			}
		}
	} else {
		resolvedLeadStatusID = &chosen.ID
		resolvedLegacyLeadStatus = strings.ToLower(chosen.Code)
		if resolvedLeadScore == 0 {
			resolvedLeadScore = chosen.Score
		}
	}

	// Helper function to convert empty string to nil pointer
	stringPtr := func(s string) *string {
		if s == "" {
			return nil
		}
		return &s
	}

	assignedTo := s.resolveSalesAssignee(req.AssignedTo, createdBy)

	l := &lead.Lead{
		FirstName:          req.FirstName,
		LastName:           req.LastName,
		CompanyName:        req.CompanyName,
		Email:              req.Email,
		Phone:              req.Phone,
		JobTitle:           req.JobTitle,
		Industry:           req.Industry,
		LeadSource:         req.LeadSource,
		LeadStatus:         resolvedLegacyLeadStatus,
		LeadStatusID:       resolvedLeadStatusID,
		LeadScore:          resolvedLeadScore,
		Probability:        req.Probability,
		EstimatedValue:     req.EstimatedValue,
		BudgetConfirmed:    req.BudgetConfirmed,
		BudgetAmount:       req.BudgetAmount,
		AuthorityConfirmed: req.AuthorityConfirmed,
		AuthorityPerson:    req.AuthorityPerson,
		NeedConfirmed:      req.NeedConfirmed,
		NeedDescription:    req.NeedDescription,
		TimelineConfirmed:  req.TimelineConfirmed,
		AssignedTo:         stringPtr(assignedTo),
		Notes:              req.Notes,
		Address:            req.Address,
		City:               req.City,
		Province:           req.Province,
		PostalCode:         req.PostalCode,
		Country:            req.Country,
		Latitude:           req.Latitude,
		Longitude:          req.Longitude,
		Website:            req.Website,
		CreatedBy:          createdBy,
	}

	if err := s.leadRepo.Create(l); err != nil {
		return nil, err
	}

	// Emit lead created event
	if s.eventHelper != nil {
		assignedTo := ""
		if l.AssignedTo != nil {
			assignedTo = *l.AssignedTo
		}
		s.eventHelper.EmitLeadCreated(&domainevents.LeadCreatedEvent{
			LeadID:     l.ID,
			Company:    l.CompanyName,
			FirstName:  l.FirstName,
			LastName:   l.LastName,
			Email:      l.Email,
			Phone:      l.Phone,
			LeadSource: l.LeadSource,
			LeadStatus: l.LeadStatus,
			AssignedTo: assignedTo,
			CreatedBy:  createdBy,
			CreatedAt:  l.CreatedAt,
		}, createdBy)
	}

	// Invalidate cache after create
	_ = s.cacheService.InvalidateOnWrite("")

	// Reload to get relations
	l, err = s.leadRepo.FindByID(l.ID)
	if err != nil {
		return nil, err
	}

	return l.ToLeadResponse(), nil
}

// Update updates a lead
func (s *Service) Update(id string, req *lead.UpdateLeadRequest, currentUser *domainauth.UserContext) (*lead.LeadResponse, error) {
	l, err := s.leadRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLeadNotFound
		}
		return nil, err
	}

	// Helper function to convert empty string to nil pointer
	stringPtr := func(s string) *string {
		if s == "" {
			return nil
		}
		return &s
	}

	// Resolve lead status from lead_statuses table.
	// Preferred input: lead_status_id. Backward compatible: lead_status (string).
	// Keep legacy l.LeadStatus in sync for existing UI paths.
	resolveLeadStatus := func() (*lead_status.LeadStatus, bool, error) {
		// If nothing provided, don't change.
		if req.LeadStatusID == "" && req.LeadStatus == "" {
			return nil, false, nil
		}

		resolveByID := func(id string) (*lead_status.LeadStatus, error) {
			if id == "" {
				return nil, nil
			}
			var ls lead_status.LeadStatus
			err := s.db.Where("id = ? AND deleted_at IS NULL", id).First(&ls).Error
			if err != nil {
				return nil, err
			}
			return &ls, nil
		}

		resolveByCodeOrLegacy := func(codeOrLegacy string) (*lead_status.LeadStatus, error) {
			if codeOrLegacy == "" {
				return nil, nil
			}
			upper := strings.ToUpper(codeOrLegacy)
			var ls lead_status.LeadStatus
			err := s.db.Where("code = ? AND deleted_at IS NULL", upper).First(&ls).Error
			if err != nil {
				return nil, err
			}
			return &ls, nil
		}

		if req.LeadStatusID != "" {
			ls, err := resolveByID(req.LeadStatusID)
			if err != nil {
				return nil, false, err
			}
			return ls, true, nil
		}

		ls, err := resolveByCodeOrLegacy(req.LeadStatus)
		if err != nil {
			return nil, false, err
		}
		return ls, true, nil
	}

	chosenStatus, statusProvided, err := resolveLeadStatus()
	if err != nil {
		// Invalid lead status provided (id/code not found)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidLeadStatus
		}
		return nil, err
	}

	// Check if lead is already converted
	if l.LeadStatus == "converted" {
		return nil, ErrLeadAlreadyConverted
	}

	oldStatus := l.LeadStatus
	newStatus := oldStatus
	statusChanged := false

	if statusProvided {
		// Update canonical FK + keep legacy string in sync
		if chosenStatus != nil {
			newStatus = strings.ToLower(chosenStatus.Code)
			statusChanged = oldStatus != newStatus
			if statusChanged && isLeadTerminalStatus(newStatus) && strings.TrimSpace(req.StatusReason) == "" {
				return nil, ErrLeadStatusReasonRequired
			}
			l.LeadStatusID = stringPtr(chosenStatus.ID)
			l.LeadStatus = newStatus
			// If caller doesn't explicitly set lead_score, follow status default score
			if req.LeadScore == nil {
				l.LeadScore = chosenStatus.Score
			}
		} else {
			// Explicitly clear status if empty values were sent (rare). Keep legacy consistent.
			l.LeadStatusID = nil
			if req.LeadStatus != "" {
				newStatus = strings.ToLower(req.LeadStatus)
				statusChanged = oldStatus != newStatus
				if statusChanged && isLeadTerminalStatus(newStatus) && strings.TrimSpace(req.StatusReason) == "" {
					return nil, ErrLeadStatusReasonRequired
				}
				l.LeadStatus = newStatus
			}
		}
	}

	// Update fields if provided
	if req.FirstName != "" {
		l.FirstName = req.FirstName
	}
	if req.LastName != "" {
		l.LastName = req.LastName
	}
	if req.CompanyName != "" {
		l.CompanyName = req.CompanyName
	}
	if req.Email != "" {
		l.Email = req.Email
	}
	if req.Phone != "" {
		l.Phone = req.Phone
	}
	if req.JobTitle != "" {
		l.JobTitle = req.JobTitle
	}
	if req.Industry != "" {
		l.Industry = req.Industry
	}
	if req.LeadSource != "" {
		l.LeadSource = req.LeadSource
	}
	if req.LeadScore != nil {
		l.LeadScore = *req.LeadScore
	}
	if req.Probability != nil {
		l.Probability = *req.Probability
	}
	if req.EstimatedValue != nil {
		l.EstimatedValue = *req.EstimatedValue
	}
	if req.BudgetConfirmed != nil {
		l.BudgetConfirmed = *req.BudgetConfirmed
	}
	if req.BudgetAmount != nil {
		l.BudgetAmount = req.BudgetAmount
	}
	if req.AuthorityConfirmed != nil {
		l.AuthorityConfirmed = *req.AuthorityConfirmed
	}
	if req.AuthorityPerson != "" {
		l.AuthorityPerson = req.AuthorityPerson
	}
	if req.NeedConfirmed != nil {
		l.NeedConfirmed = *req.NeedConfirmed
	}
	if req.NeedDescription != "" {
		l.NeedDescription = req.NeedDescription
	}
	if req.TimelineConfirmed != nil {
		l.TimelineConfirmed = *req.TimelineConfirmed
	}

	if req.Notes != "" {
		l.Notes = req.Notes
	}
	if req.Address != "" {
		l.Address = req.Address
	}
	if req.City != "" {
		l.City = req.City
	}
	if req.Province != "" {
		l.Province = req.Province
	}
	if req.PostalCode != "" {
		l.PostalCode = req.PostalCode
	}
	if req.Country != "" {
		l.Country = req.Country
	}
	if req.Latitude != nil {
		l.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		l.Longitude = req.Longitude
	}
	if req.Website != "" {
		l.Website = req.Website
	}

	if statusChanged {
		appendLeadStatusHistory(l, oldStatus, newStatus, req.StatusReason, currentUser)
	}

	if err := s.leadRepo.Update(l); err != nil {
		return nil, err
	}

	if statusChanged && s.eventHelper != nil && currentUser != nil {
		s.eventHelper.EmitLeadStatusChanged(&domainevents.LeadStatusChangedEvent{
			LeadID:    l.ID,
			OldStatus: oldStatus,
			NewStatus: newStatus,
			ChangedBy: currentUser.UserID,
			ChangedAt: time.Now(),
			Reason:    strings.TrimSpace(req.StatusReason),
		}, currentUser.UserID)
	}

	// Invalidate cache after update
	_ = s.cacheService.InvalidateOnWrite(id)

	// Reload to get relations
	l, err = s.leadRepo.FindByID(l.ID)
	if err != nil {
		return nil, err
	}

	return l.ToLeadResponse(), nil
}

func isLeadTerminalStatus(status string) bool {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "converted", "won", "lost":
		return true
	default:
		return false
	}
}

func appendLeadStatusHistory(l *lead.Lead, oldStatus, newStatus, reason string, currentUser *domainauth.UserContext) {
	var metadata map[string]interface{}
	if len(l.ConversionMetadata) > 0 {
		_ = json.Unmarshal(l.ConversionMetadata, &metadata)
	}
	if metadata == nil {
		metadata = map[string]interface{}{}
	}

	history, _ := metadata["status_history"].([]interface{})
	entry := map[string]interface{}{
		"from_status": oldStatus,
		"to_status":   newStatus,
		"changed_at":  time.Now().Format(time.RFC3339),
	}
	if currentUser != nil && currentUser.UserID != "" {
		entry["changed_by"] = currentUser.UserID
	}
	if trimmedReason := strings.TrimSpace(reason); trimmedReason != "" {
		entry["reason"] = trimmedReason
		metadata["latest_status_reason"] = trimmedReason
	} else {
		delete(metadata, "latest_status_reason")
	}
	metadata["latest_status"] = newStatus
	metadata["latest_status_changed_at"] = entry["changed_at"]

	metadata["status_history"] = append(history, entry)
	encoded, err := json.Marshal(metadata)
	if err == nil {
		l.ConversionMetadata = encoded
	}
}

// Delete deletes a lead
func (s *Service) Delete(id string) error {
	l, err := s.leadRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrLeadNotFound
		}
		return err
	}

	// Check if lead is already converted - cannot delete converted leads
	if l.LeadStatus == "converted" {
		return ErrLeadAlreadyConverted
	}

	err = s.leadRepo.Delete(l.ID)
	if err != nil {
		return err
	}

	// Invalidate cache after delete
	_ = s.cacheService.InvalidateOnWrite(id)

	return nil
}

// Convert converts a qualified lead to opportunity/deal
func (s *Service) Convert(id string, req *lead.ConvertLeadRequest, convertedBy string) (*lead.ConvertLeadResponse, error) {
	trimmedReason := strings.TrimSpace(req.StatusReason)
	if trimmedReason == "" {
		return nil, ErrLeadStatusReasonRequired
	}

	var response *lead.ConvertLeadResponse
	var contactID string
	var deal *pipeline.Deal
	var createdAccount interface{}
	var createdContact interface{}

	l := &lead.Lead{}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Get lead inside transaction
		if err := tx.Where("id = ? AND deleted_at IS NULL", id).First(l).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrLeadNotFound
			}
			return err
		}

		// Validate lead status must be "qualified"
		if l.LeadStatus != "qualified" {
			return ErrLeadCannotConvert
		}

		// Check if already converted
		if l.LeadStatus == "converted" || (l.OpportunityID != nil && *l.OpportunityID != "") {
			return ErrLeadAlreadyConverted
		}

		// Validate stage exists
		var stage pipeline.PipelineStage
		if err := tx.Where("id = ? AND deleted_at IS NULL", req.StageID).First(&stage).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrStageNotFound
			}
			return err
		}
		if !stage.IsWon {
			return ErrInvalidConversionStage
		}
		if len(req.ProductItems) == 0 {
			return ErrSoldProductsRequired
		}

		var accountID string

		leadFullName := strings.TrimSpace(strings.Join([]string{l.FirstName, l.LastName}, " "))
		accountName := strings.TrimSpace(l.CompanyName)
		if accountName == "" {
			accountName = leadFullName
		}
		if accountName == "" {
			accountName = strings.TrimSpace(l.Email)
		}
		if accountName == "" {
			accountName = "Lead " + l.ID
		}

		if l.AccountID != nil && *l.AccountID != "" {
			// Use existing account from lead
			accountID = *l.AccountID
		} else {
			// Create account automatically from lead data
			var categories []account.Category
			if err := tx.Order("code ASC").Find(&categories).Error; err != nil || len(categories) == 0 {
				return ErrAccountCreationFailed
			}

			acc := &account.Account{
				Name:       accountName,
				CategoryID: categories[0].ID,
				Email:      l.Email,
				Phone:      l.Phone,
				Address:    l.Address,
				City:       l.City,
				Province:   l.Province,
				PostalCode: l.PostalCode,
				Country:    l.Country,
				Website:    l.Website,
				Industry:   l.Industry,
				Latitude:   l.Latitude,
				Longitude:  l.Longitude,
				Status:     "active",
			}
			if l.AssignedTo != nil && *l.AssignedTo != "" {
				acc.AssignedTo = l.AssignedTo
			}
			s.populateAccountCoordinatesFromLead(acc, l)
			var brickErr error
			acc.BrickID, brickErr = s.resolveBrickIDFromLead(l)
			if brickErr != nil {
				// resolvedBrickID already logs warnings instead of returning errors
			}

			if err := tx.Create(acc).Error; err != nil {
				return ErrAccountCreationFailed
			}

			accountID = acc.ID
			createdAccount = acc.ToAccountResponse()
		}

		if accountID == "" {
			return ErrAccountCreationFailed
		}
		if createdAccount == nil {
			updatedAccount, err := s.syncAccountFromLead(tx, accountID, l, accountName)
			if err != nil {
				return ErrAccountCreationFailed
			}
			createdAccount = updatedAccount.ToAccountResponse()
		}

		// Always ensure a contact is attached to the converted account using lead data.
		if accountID != "" && l.ContactID != nil && *l.ContactID != "" {
			updatedContact, err := s.syncContactFromLead(tx, *l.ContactID, accountID, l, leadFullName, accountName)
			if err != nil {
				return ErrContactCreationFailed
			}
			contactID = updatedContact.ID
			createdContact = updatedContact.ToContactResponse()
		} else if accountID != "" {
			var contactRoles []contact.ContactRole
			if err := tx.Order("code ASC").Find(&contactRoles).Error; err != nil || len(contactRoles) == 0 {
				return ErrContactCreationFailed
			}

			contactName := leadFullName
			if contactName == "" {
				contactName = accountName
			}

			c := &contact.Contact{
				AccountID: accountID,
				Name:      contactName,
				RoleID:    contactRoles[0].ID,
				Email:     l.Email,
				Phone:     l.Phone,
				Position:  l.JobTitle,
			}

			if err := tx.Create(c).Error; err != nil {
				return ErrContactCreationFailed
			}

			contactID = c.ID
			createdContact = c.ToContactResponse()
		}

		// Create deal/opportunity
		dealValue := l.EstimatedValue
		if req.Value != nil {
			dealValue = *req.Value
		}
		if len(req.ProductItems) > 0 {
			dealValue = 0
			for _, itemReq := range req.ProductItems {
				if itemReq.Quantity <= 0 {
					continue
				}
				unitPrice := int64(0)
				if itemReq.UnitPrice != nil {
					unitPrice = *itemReq.UnitPrice
				}
				discount := int64(0)
				if itemReq.DiscountAmount != nil {
					discount = *itemReq.DiscountAmount
				}
				subtotal := unitPrice*int64(itemReq.Quantity) - discount
				if subtotal > 0 {
					dealValue += subtotal
				}
			}
		}

		dealStatus := "won" // Won since stage.IsWon is verified above
		conversionTime := time.Now()
		var actualCloseDate *time.Time = &conversionTime

		budgetConfirmed := l.BudgetConfirmed
		authorityConfirmed := l.AuthorityConfirmed
		needConfirmed := l.NeedConfirmed
		timelineConfirmed := l.TimelineConfirmed
		qualificationSnapshot := []byte("{}")
		needProducts := make([]leadqualification.NeedProduct, 0)
		var qualification leadqualification.LeadQualificationChecklist
		if err := tx.Where("lead_id = ?", l.ID).First(&qualification).Error; err == nil {
			budgetConfirmed = qualification.BudgetConfirmed
			authorityConfirmed = qualification.AuthorityConfirmed
			needConfirmed = qualification.NeedConfirmed
			timelineConfirmed = qualification.TimelineConfirmed
			_ = json.Unmarshal(qualification.NeedTargetProducts, &needProducts)
			qualificationSnapshot, _ = json.Marshal(map[string]interface{}{
				"budget_target_amount":    qualification.BudgetTargetAmount,
				"budget_target_currency":  qualification.BudgetTargetCurrency,
				"budget_confirmed":        qualification.BudgetConfirmed,
				"budget_notes":            qualification.BudgetNotes,
				"authority_target_person": qualification.AuthorityTargetPerson,
				"authority_target_role":   qualification.AuthorityTargetRole,
				"authority_confirmed":     qualification.AuthorityConfirmed,
				"authority_notes":         qualification.AuthorityNotes,
				"need_target_products":    needProducts,
				"need_priority_level":     qualification.NeedPriorityLevel,
				"need_confirmed":          qualification.NeedConfirmed,
				"need_notes":              qualification.NeedNotes,
				"timeline_target_date":    qualification.TimelineTargetDate,
				"timeline_flexibility":    qualification.TimelineFlexibility,
				"timeline_confirmed":      qualification.TimelineConfirmed,
				"timeline_notes":          qualification.TimelineNotes,
				"qualification_score":     qualification.QualificationScore,
				"qualification_status":    qualification.QualificationStatus,
			})
		}

		// Helper function to convert empty string to nil pointer
		stringPtr := func(s string) *string {
			if s == "" {
				return nil
			}
			return &s
		}

		dealAssignedTo := ""
		if l.AssignedTo != nil {
			dealAssignedTo = *l.AssignedTo
		}
		dealAssignedTo = s.resolveSalesAssignee(dealAssignedTo, convertedBy)

		deal = &pipeline.Deal{
			Title:                 req.OpportunityTitle,
			Description:           req.OpportunityDescription,
			AccountID:             accountID,
			ContactID:             stringPtr(contactID),
			StageID:               stage.ID,
			Value:                 dealValue,
			Probability:           100, // Closed Won
			ActualCloseDate:       actualCloseDate,
			AssignedTo:            stringPtr(dealAssignedTo),
			LeadID:                &l.ID,
			Status:                dealStatus,
			Source:                l.LeadSource,
			BudgetConfirmed:       budgetConfirmed,
			AuthorityConfirmed:    authorityConfirmed,
			NeedConfirmed:         needConfirmed,
			TimelineConfirmed:     timelineConfirmed,
			QualificationSnapshot: qualificationSnapshot,
			Notes:                 l.Notes,
			CreatedBy:             convertedBy,
		}

		if err := tx.Create(deal).Error; err != nil {
			return ErrOpportunityCreationFailed
		}

		if len(req.ProductItems) > 0 {
			productItemsTotal, err := s.createDealProductItemsFromConvertItems(tx, deal.ID, req.ProductItems)
			if err != nil {
				if errors.Is(err, ErrSoldProductsRequired) {
					return ErrSoldProductsRequired
				}
				return ErrOpportunityCreationFailed
			}
			if productItemsTotal > 0 && productItemsTotal != deal.Value {
				deal.Value = productItemsTotal
				if err := tx.Model(deal).Update("value", productItemsTotal).Error; err != nil {
					return ErrOpportunityCreationFailed
				}
			}
		}

		if err := tx.
			Preload("Account").
			Preload("Contact").
			Preload("Stage").
			Preload("ProductItems").
			Preload("AssignedUser").
			Where("id = ?", deal.ID).
			First(deal).Error; err != nil {
			return ErrOpportunityCreationFailed
		}

		// Create initial deal history entry
		if s.dealHistoryRepo != nil {
			leadJourney := "Lead journey: New → Contacted → Qualified"
			if l.CreatedAt.IsZero() == false {
				daysTotal := int(time.Since(l.CreatedAt).Hours() / 24)
				leadJourney += " (" + strconv.Itoa(daysTotal) + " days total)"
			}

			history := &deal_history.DealHistory{
				DealID:          deal.ID,
				FromStageID:     nil,
				FromStageName:   "",
				ToStageID:       deal.StageID,
				ToStageName:     stage.Name,
				FromProbability: 0,
				ToProbability:   deal.Probability,
				DaysInPrevStage: nil,
				ChangedBy:       convertedBy,
				ChangedAt:       time.Now(),
				Reason:          "Deal created from qualified lead",
				Notes:           leadJourney,
			}
			if err := tx.Create(history).Error; err != nil {
				log.Printf("Warning: Failed to create deal history: %v", err)
			}
		}

		// Update lead status to converted
		now := conversionTime
		oldLeadStatus := l.LeadStatus
		l.LeadStatus = "converted"
		if convertedStatus, err := s.findConvertedLeadStatus(); err != nil {
			return err
		} else if convertedStatus != nil {
			l.LeadStatusID = &convertedStatus.ID
			l.LeadScore = convertedStatus.Score
		} else {
			l.LeadStatusID = nil
		}
		dealID := deal.ID
		l.OpportunityID = &dealID
		accountIDPtr := accountID
		l.AccountID = &accountIDPtr
		if contactID != "" {
			contactIDPtr := contactID
			l.ContactID = &contactIDPtr
		}
		l.ConvertedAt = &now
		convertedByPtr := convertedBy
		l.ConvertedBy = &convertedByPtr
		appendLeadStatusHistory(l, oldLeadStatus, l.LeadStatus, trimmedReason, &domainauth.UserContext{UserID: convertedBy})

		if err := tx.Save(l).Error; err != nil {
			return err
		}

		// Auto-migrate Activities
		if s.activityRepo != nil {
			updates := make(map[string]interface{})
			updates["deal_id"] = deal.ID
			if accountID != "" {
				updates["account_id"] = accountID
			}
			if err := tx.Model(&activity.Activity{}).Where("lead_id = ?", l.ID).Updates(updates).Error; err != nil {
				return err
			}
		}

		// Auto-migrate Visit Reports
		if s.visitReportRepo != nil {
			updates := make(map[string]interface{})
			updates["deal_id"] = deal.ID
			if accountID != "" {
				updates["account_id"] = accountID
			}
			if err := tx.Table("visit_reports").Where("lead_id = ?", l.ID).Updates(updates).Error; err != nil {
				return err
			}
		}

		// Auto-migrate Tasks
		if s.taskRepo != nil {
			updates := make(map[string]interface{})
			updates["deal_id"] = deal.ID
			updates["lead_id"] = nil
			if accountID != "" {
				updates["account_id"] = accountID
			}
			if err := tx.Model(&task.Task{}).Where("lead_id = ?", l.ID).Updates(updates).Error; err != nil {
				return err
			}
		}

		s.logLeadConversionActivity(l, deal, accountID, contactID, convertedBy, conversionTime)

		response = &lead.ConvertLeadResponse{
			Lead:        l.ToLeadResponse(),
			Opportunity: deal.ToDealResponse(),
		}
		if createdAccount != nil {
			response.Account = createdAccount
		}
		if createdContact != nil {
			response.Contact = createdContact
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Reload lead outside transaction to get fully populated relations (e.g. Creator, Status, etc.)
	finalLead, err := s.leadRepo.FindByID(id)
	if err == nil && finalLead != nil {
		response.Lead = finalLead.ToLeadResponse()
	}

	_ = s.cacheService.InvalidateOnWrite(id)

	// Emit events outside transaction after it succeeds
	if s.eventHelper != nil {
		assignedTo := ""
		if deal.AssignedTo != nil {
			assignedTo = *deal.AssignedTo
		}
		s.eventHelper.EmitDealCreated(&domainevents.DealCreatedEvent{
			DealID:            deal.ID,
			Title:             deal.Title,
			Value:             deal.Value,
			AccountID:         deal.AccountID,
			ContactID:         contactID,
			StageID:           deal.StageID,
			StageName:         deal.Stage.Name,
			PipelineID:        "",
			AssignedTo:        assignedTo,
			ExpectedCloseDate: deal.ExpectedCloseDate,
			CreatedBy:         convertedBy,
			CreatedAt:         deal.CreatedAt,
		}, convertedBy)
		s.eventHelper.EmitLeadConverted(&domainevents.LeadConvertedEvent{
			LeadID:        id,
			OpportunityID: deal.ID,
			AccountID:     deal.AccountID,
			ContactID:     contactID,
			ConvertedBy:   convertedBy,
			ConvertedAt:   time.Now(),
		}, convertedBy)
		s.eventHelper.EmitLeadStatusChanged(&domainevents.LeadStatusChangedEvent{
			LeadID:    id,
			OldStatus: "qualified",
			NewStatus: "converted",
			ChangedBy: convertedBy,
			ChangedAt: time.Now(),
			Reason:    "Lead converted to deal",
		}, convertedBy)

		if deal.Status == "won" && deal.ActualCloseDate != nil {
			s.eventHelper.EmitDealWon(&domainevents.DealWonEvent{
				DealID:          deal.ID,
				Title:           deal.Title,
				Value:           deal.Value,
				AccountID:       deal.AccountID,
				AssignedTo:      assignedTo,
				ActualCloseDate: *deal.ActualCloseDate,
				WonBy:           convertedBy,
				WonAt:           time.Now(),
			}, convertedBy)
		}
	}

	return response, nil
}

func (s *Service) logLeadConversionActivity(l *lead.Lead, deal *pipeline.Deal, accountID, contactID, convertedBy string, convertedAt time.Time) {
	if s.activityRepo == nil || l == nil || deal == nil || convertedBy == "" {
		return
	}

	description := "Lead converted to deal"
	if strings.TrimSpace(deal.Title) != "" {
		description = "Lead converted to deal: " + strings.TrimSpace(deal.Title)
	}

	metadata, err := json.Marshal(map[string]interface{}{
		"event":        "lead_converted",
		"lead_id":      l.ID,
		"deal_id":      deal.ID,
		"deal_status":  deal.Status,
		"stage_id":     deal.StageID,
		"converted_at": convertedAt.Format(time.RFC3339),
	})
	if err != nil {
		metadata = nil
	}

	activityRecord := &activity.Activity{
		Type:        "note",
		LeadID:      &l.ID,
		DealID:      &deal.ID,
		UserID:      convertedBy,
		Description: description,
		Timestamp:   convertedAt,
		Metadata:    metadata,
	}
	if accountID != "" {
		activityRecord.AccountID = &accountID
	}
	if contactID != "" {
		activityRecord.ContactID = &contactID
	}

	_ = s.activityRepo.Create(activityRecord)
}

func (s *Service) findConvertedLeadStatus() (*lead_status.LeadStatus, error) {
	if s.leadStatusRepo == nil {
		return nil, nil
	}

	for _, code := range []string{"CONVERTED", "converted"} {
		status, err := s.leadStatusRepo.FindByCode(code)
		if err == nil {
			return status, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	statuses, err := s.leadStatusRepo.ListAll()
	if err != nil {
		return nil, err
	}
	for _, status := range statuses {
		if status.IsConverted {
			return status, nil
		}
	}

	return nil, nil
}

func (s *Service) syncAccountFromLead(tx *gorm.DB, accountID string, l *lead.Lead, fallbackName string) (*account.Account, error) {
	var accountEntity account.Account
	if err := tx.Preload("Category").Where("id = ? AND deleted_at IS NULL", accountID).First(&accountEntity).Error; err != nil {
		return nil, err
	}

	changed := false
	setString := func(target *string, value string) {
		value = strings.TrimSpace(value)
		if value != "" && *target != value {
			*target = value
			changed = true
		}
	}
	setFloatPtr := func(target **float64, value *float64) {
		if value == nil {
			return
		}
		if *target == nil || **target != *value {
			*target = value
			changed = true
		}
	}

	if strings.TrimSpace(l.CompanyName) != "" {
		setString(&accountEntity.Name, l.CompanyName)
	} else if strings.TrimSpace(accountEntity.Name) == "" {
		setString(&accountEntity.Name, fallbackName)
	}
	setString(&accountEntity.Email, l.Email)
	setString(&accountEntity.Phone, l.Phone)
	setString(&accountEntity.Address, l.Address)
	setString(&accountEntity.City, l.City)
	setString(&accountEntity.Province, l.Province)
	setString(&accountEntity.PostalCode, l.PostalCode)
	setString(&accountEntity.Country, l.Country)
	setString(&accountEntity.Website, l.Website)
	setString(&accountEntity.Industry, l.Industry)
	setFloatPtr(&accountEntity.Latitude, l.Latitude)
	setFloatPtr(&accountEntity.Longitude, l.Longitude)
	if s.populateAccountCoordinatesFromLead(&accountEntity, l) {
		changed = true
	}

	if l.AssignedTo != nil && *l.AssignedTo != "" {
		if accountEntity.AssignedTo == nil || *accountEntity.AssignedTo != *l.AssignedTo {
			accountEntity.AssignedTo = l.AssignedTo
			changed = true
		}
	}

	resolvedBrickID, err := s.resolveBrickIDFromLead(l)
	if err != nil {
		return nil, err
	}
	if !sameStringPointer(accountEntity.BrickID, resolvedBrickID) {
		accountEntity.BrickID = resolvedBrickID
		changed = true
	}

	if changed {
		if err := tx.Save(&accountEntity).Error; err != nil {
			return nil, err
		}
	}

	return &accountEntity, nil
}

func (s *Service) syncContactFromLead(tx *gorm.DB, contactID, accountID string, l *lead.Lead, fullName, fallbackAccountName string) (*contact.Contact, error) {
	var contactEntity contact.Contact
	if err := tx.Preload("Role").Where("id = ? AND deleted_at IS NULL", contactID).First(&contactEntity).Error; err != nil {
		return nil, err
	}

	changed := false
	setString := func(target *string, value string) {
		value = strings.TrimSpace(value)
		if value != "" && *target != value {
			*target = value
			changed = true
		}
	}

	if contactEntity.AccountID != accountID {
		contactEntity.AccountID = accountID
		changed = true
	}

	contactName := strings.TrimSpace(fullName)
	if contactName == "" {
		contactName = strings.TrimSpace(fallbackAccountName)
	}
	if contactName != "" && strings.TrimSpace(contactEntity.Name) != contactName {
		contactEntity.Name = contactName
		changed = true
	}

	setString(&contactEntity.Email, l.Email)
	setString(&contactEntity.Phone, l.Phone)
	setString(&contactEntity.Position, l.JobTitle)

	if changed {
		if err := tx.Save(&contactEntity).Error; err != nil {
			return nil, err
		}
	}

	return &contactEntity, nil
}

func (s *Service) populateAccountCoordinatesFromLead(accountEntity *account.Account, l *lead.Lead) bool {
	if accountEntity == nil || l == nil {
		return false
	}

	if l.Latitude != nil && l.Longitude != nil {
		return false
	}
	if accountEntity.Latitude != nil && accountEntity.Longitude != nil {
		return false
	}
	if !s.geocodingEnabled || s.geocodingSvc == nil {
		return false
	}
	if strings.TrimSpace(l.Address) == "" && strings.TrimSpace(l.City) == "" && strings.TrimSpace(l.Province) == "" {
		return false
	}

	result, err := s.geocodingSvc.GeocodeAddressWithFallback(l.Address, l.City, l.Province)
	if err != nil {
		log.Printf("Warning: Failed to geocode account from lead %s during convert: %v", l.ID, err)
		return false
	}

	changed := false
	if accountEntity.Latitude == nil || *accountEntity.Latitude != result.Latitude {
		accountEntity.Latitude = &result.Latitude
		changed = true
	}
	if accountEntity.Longitude == nil || *accountEntity.Longitude != result.Longitude {
		accountEntity.Longitude = &result.Longitude
		changed = true
	}

	return changed
}

func (s *Service) createDealProductItemsFromConvertItems(tx *gorm.DB, dealID string, productItems []lead.ConvertLeadProductItemRequest) (int64, error) {
	productItems = normalizeConvertLeadProductItemRequests(productItems)
	total := int64(0)
	items := make([]pipeline.DealProductItem, 0, len(productItems))
	for _, itemReq := range productItems {
		if itemReq.ProductID == "" || itemReq.Quantity < 1 {
			continue
		}

		var productSnapshot struct {
			ID           string
			Name         string
			SKU          string
			Price        int64
			Cost         int64
			CategoryID   *string
			CategoryName string
		}
		err := tx.Table("products AS p").
			Select("p.id, p.name, p.sku, p.price, p.cost, p.category_id, COALESCE(pc.name, '') AS category_name").
			Joins("LEFT JOIN product_categories pc ON pc.id = p.category_id").
			Where("p.id = ? AND p.deleted_at IS NULL", itemReq.ProductID).
			Scan(&productSnapshot).Error
		if err != nil || productSnapshot.ID == "" {
			return 0, errors.New("product not found")
		}

		unitPrice := productSnapshot.Price
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
			DealID:              dealID,
			ProductID:           productSnapshot.ID,
			ProductName:         productSnapshot.Name,
			ProductSKU:          productSnapshot.SKU,
			UnitPrice:           unitPrice,
			UnitCost:            productSnapshot.Cost,
			Quantity:            itemReq.Quantity,
			DiscountAmount:      discount,
			Subtotal:            subtotal,
			ProductCategoryID:   productSnapshot.CategoryID,
			ProductCategoryName: productSnapshot.CategoryName,
			Notes:               itemReq.Notes,
		})
	}

	if len(items) == 0 {
		return 0, ErrSoldProductsRequired
	}
	if err := tx.Create(&items).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func normalizeConvertLeadProductItemRequests(productItems []lead.ConvertLeadProductItemRequest) []lead.ConvertLeadProductItemRequest {
	normalized := make([]lead.ConvertLeadProductItemRequest, 0, len(productItems))
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

func (s *Service) resolveSalesAssignee(preferredUserID string, fallbackUserID string) string {
	preferredUserID = strings.TrimSpace(preferredUserID)
	fallbackUserID = strings.TrimSpace(fallbackUserID)

	if s.isActiveSalesUser(preferredUserID) {
		return preferredUserID
	}
	if s.isActiveSalesUser(fallbackUserID) {
		return fallbackUserID
	}

	var salesUser struct {
		ID string
	}
	if s.db != nil {
		err := s.db.Table("users u").
			Select("u.id").
			Joins("INNER JOIN roles r ON r.id = u.role_id AND r.deleted_at IS NULL AND r.code = ?", "sales").
			Where("u.deleted_at IS NULL AND u.status = ?", "active").
			Order("u.created_at ASC").
			Limit(1).
			Scan(&salesUser).Error
		if err == nil && salesUser.ID != "" {
			return salesUser.ID
		}
	}

	if preferredUserID != "" {
		return preferredUserID
	}
	return fallbackUserID
}

func (s *Service) isActiveSalesUser(userID string) bool {
	if strings.TrimSpace(userID) == "" || s.db == nil {
		return false
	}

	var count int64
	err := s.db.Table("users u").
		Joins("INNER JOIN roles r ON r.id = u.role_id AND r.deleted_at IS NULL AND r.code = ?", "sales").
		Where("u.id = ? AND u.deleted_at IS NULL AND u.status = ?", userID, "active").
		Count(&count).Error
	return err == nil && count > 0
}

// GetAnalytics returns lead analytics
func (s *Service) GetAnalytics(req *lead.LeadAnalyticsRequest) (*lead.LeadAnalyticsResponse, error) {
	return s.leadRepo.GetAnalytics(req)
}

// CreateAccountFromLead creates an account from a lead (pre-convert)
func (s *Service) CreateAccountFromLead(leadID string, req *lead.CreateAccountFromLeadRequest, createdBy string) (*lead.CreateAccountFromLeadResponse, error) {
	// Get lead
	l, err := s.leadRepo.FindByID(leadID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLeadNotFound
		}
		return nil, err
	}

	// Validate lead status must be "qualified" to create account
	if l.LeadStatus != "qualified" {
		return nil, errors.New("only qualified leads can be used to create accounts")
	}

	// Check if lead already has an account
	if l.AccountID != nil && *l.AccountID != "" {
		return nil, errors.New("lead already has an account")
	}

	// Check if company name exists
	if l.CompanyName == "" {
		return nil, errors.New("company name is required to create account")
	}

	// Get category
	var categoryID string
	if req.CategoryID != "" {
		// Verify category exists
		_, err := s.categoryRepo.FindByID(req.CategoryID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("category not found")
			}
			return nil, err
		}
		categoryID = req.CategoryID
	} else {
		// Use first available category
		categories, err := s.categoryRepo.List()
		if err != nil || len(categories) == 0 {
			return nil, ErrAccountCreationFailed
		}
		categoryID = categories[0].ID
	}

	// Create account
	account := &account.Account{
		Name:       l.CompanyName,
		CategoryID: categoryID,
		Email:      l.Email,
		Phone:      l.Phone,
		Address:    l.Address,
		City:       l.City,
		Province:   l.Province,
		Status:     "active",
	}
	if l.AssignedTo != nil && *l.AssignedTo != "" {
		account.AssignedTo = l.AssignedTo
	}
	account.BrickID, err = s.resolveBrickIDFromLead(l)
	if err != nil {
		return nil, ErrAccountCreationFailed
	}

	if err := s.accountRepo.Create(account); err != nil {
		return nil, ErrAccountCreationFailed
	}

	// Update lead with account ID
	accountIDPtr := account.ID
	l.AccountID = &accountIDPtr
	if err := s.leadRepo.Update(l); err != nil {
		return nil, err
	}

	// Create contact if requested
	var createdContact interface{}
	if req.CreateContact {
		// Find default contact role
		contactRoles, err := s.contactRoleRepo.List()
		if err != nil || len(contactRoles) == 0 {
			// Contact creation is optional, continue without it
		} else {
			contactName := l.FirstName
			if l.LastName != "" {
				contactName += " " + l.LastName
			}

			contact := &contact.Contact{
				AccountID: account.ID,
				Name:      contactName,
				RoleID:    contactRoles[0].ID,
				Email:     l.Email,
				Phone:     l.Phone,
				Position:  l.JobTitle,
			}

			if err := s.contactRepo.Create(contact); err == nil {
				contactIDPtr := contact.ID
				l.ContactID = &contactIDPtr
				_ = s.leadRepo.Update(l) // Ignore error
				createdContact = contact.ToContactResponse()
			}
		}
	}

	// Reload lead to get relations
	l, err = s.leadRepo.FindByID(l.ID)
	if err != nil {
		return nil, err
	}

	return &lead.CreateAccountFromLeadResponse{
		Lead:    l.ToLeadResponse(),
		Account: account.ToAccountResponse(),
		Contact: createdContact,
	}, nil
}

func sameStringPointer(left, right *string) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil || right == nil {
		return false
	}
	return *left == *right
}

// GetFormData returns form data for creating a lead
func (s *Service) GetFormData() (*lead.LeadFormDataResponse, error) {
	// Get lead sources from database
	var dbLeadSources []lead_source.LeadSource
	if err := s.db.Where("is_active = ?", true).Order("\"order\" ASC").Find(&dbLeadSources).Error; err != nil {
		// Fallback to empty array if error
		dbLeadSources = []lead_source.LeadSource{}
	}

	leadSources := make([]lead.LeadSourceOption, len(dbLeadSources))
	for i, ls := range dbLeadSources {
		leadSources[i] = lead.LeadSourceOption{
			ID:    ls.ID,
			Value: ls.Code,
			Label: ls.Name,
		}
	}

	// Lead statuses (with intermediate statuses) - legacy, should use lead_statuses API
	leadStatuses := []lead.LeadStatusOption{
		{Value: "new", Label: "New"},
		{Value: "contacted", Label: "Contacted"},
		{Value: "interested", Label: "Interested"},
		{Value: "qualified", Label: "Qualified"},
		{Value: "proposal_sent", Label: "Proposal Sent"},
		{Value: "converted", Label: "Converted"},
		{Value: "lost", Label: "Lost"},
	}

	// Get active sales users for assigned_to. Opportunities are scoped by sales
	// owner, so lead assignment options must not include admin/non-sales users.
	var users []user.User
	if err := s.db.Model(&user.User{}).
		Joins("INNER JOIN roles r ON r.id = users.role_id AND r.deleted_at IS NULL AND r.code = ?", "sales").
		Where("users.deleted_at IS NULL AND users.status = ?", "active").
		Order("users.name ASC").
		Limit(100).
		Find(&users).Error; err != nil {
		return nil, err
	}

	userOptions := make([]lead.UserOption, len(users))
	for i, u := range users {
		userOptions[i] = lead.UserOption{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		}
	}

	// Get industries from database
	var dbIndustries []industry.Industry
	if err := s.db.Where("is_active = ?", true).Order("\"order\" ASC").Find(&dbIndustries).Error; err != nil {
		// Fallback to empty array if error
		dbIndustries = []industry.Industry{}
	}

	industries := make([]string, len(dbIndustries))
	for i, ind := range dbIndustries {
		industries[i] = ind.Name
	}

	// Provinces in Indonesia
	provinces := []string{
		"DKI Jakarta",
		"Jawa Barat",
		"Jawa Tengah",
		"Jawa Timur",
		"Yogyakarta",
		"Banten",
		"Bali",
		"Sumatera Utara",
		"Sumatera Barat",
		"Sumatera Selatan",
		"Lampung",
		"Riau",
		"Kepulauan Riau",
		"Jambi",
		"Aceh",
		"Kalimantan Barat",
		"Kalimantan Tengah",
		"Kalimantan Selatan",
		"Kalimantan Timur",
		"Kalimantan Utara",
		"Sulawesi Utara",
		"Sulawesi Tengah",
		"Sulawesi Selatan",
		"Sulawesi Tenggara",
		"Gorontalo",
		"Maluku",
		"Maluku Utara",
		"Papua",
		"Papua Barat",
		"Papua Selatan",
		"Papua Tengah",
		"Papua Pegunungan",
	}

	// Default values
	defaults := lead.LeadFormDefaults{
		Country:    "Indonesia",
		LeadStatus: "new",
		LeadScore:  0,
	}

	return &lead.LeadFormDataResponse{
		LeadSources:  leadSources,
		LeadStatuses: leadStatuses,
		Users:        userOptions,
		Industries:   industries,
		Provinces:    provinces,
		Defaults:     defaults,
	}, nil
}
