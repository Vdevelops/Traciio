package lead

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/account"
	domainauth "github.com/gilabs/crm-healthcare/api/internal/domain/auth"
	"github.com/gilabs/crm-healthcare/api/internal/domain/contact"
	"github.com/gilabs/crm-healthcare/api/internal/domain/deal_history"
	domainevents "github.com/gilabs/crm-healthcare/api/internal/domain/events"
	"github.com/gilabs/crm-healthcare/api/internal/domain/industry"
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead"
	leadqualification "github.com/gilabs/crm-healthcare/api/internal/domain/lead_qualification"
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead_source"
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead_status"
	"github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
	"github.com/gilabs/crm-healthcare/api/internal/domain/task"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"github.com/gilabs/crm-healthcare/api/pkg/cache"
	"gorm.io/gorm"
)

var (
	ErrLeadNotFound              = errors.New("lead not found")
	ErrLeadAlreadyConverted      = errors.New("lead already converted")
	ErrLeadCannotConvert         = errors.New("lead cannot convert")
	ErrInvalidLeadStatus         = errors.New("invalid lead status")
	ErrInvalidLeadSource         = errors.New("invalid lead source")
	ErrStageNotFound             = errors.New("stage not found")
	ErrAccountCreationFailed     = errors.New("account creation failed")
	ErrContactCreationFailed     = errors.New("contact creation failed")
	ErrOpportunityCreationFailed = errors.New("opportunity creation failed")
)

type Service struct {
	db              *gorm.DB
	leadRepo        interfaces.LeadRepository
	dealRepo        interfaces.DealRepository
	pipelineRepo    interfaces.PipelineRepository
	accountRepo     interfaces.AccountRepository
	contactRepo     interfaces.ContactRepository
	categoryRepo    interfaces.CategoryRepository
	contactRoleRepo interfaces.ContactRoleRepository
	userRepo        interfaces.UserRepository
	activityRepo    interfaces.ActivityRepository    // For auto-migrate activities
	visitReportRepo interfaces.VisitReportRepository // For auto-migrate visit reports
	taskRepo        interfaces.TaskRepository        // For auto-create tasks
	dealHistoryRepo interfaces.DealHistoryRepository // For deal history logging
	leadStatusRepo  interfaces.LeadStatusRepository  // For lead status lookup
	cacheService    *cache.LeadCacheService
	eventHelper     *domainevents.Helper // For emitting domain events
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
	eventHelper *domainevents.Helper,
) *Service {
	return &Service{
		db:              db,
		leadRepo:        leadRepo,
		dealRepo:        dealRepo,
		pipelineRepo:    pipelineRepo,
		accountRepo:     accountRepo,
		contactRepo:     contactRepo,
		categoryRepo:    categoryRepo,
		contactRoleRepo: contactRoleRepo,
		userRepo:        userRepo,
		activityRepo:    activityRepo,
		visitReportRepo: visitReportRepo,
		taskRepo:        taskRepo,
		dealHistoryRepo: dealHistoryRepo,
		leadStatusRepo:  leadStatusRepo,
		cacheService:    cache.NewLeadCacheService(nil),
		eventHelper:     eventHelper,
	}
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
		AssignedTo:         stringPtr(createdBy),
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

	if statusProvided {
		// Update canonical FK + keep legacy string in sync
		if chosenStatus != nil {
			l.LeadStatusID = stringPtr(chosenStatus.ID)
			l.LeadStatus = strings.ToLower(chosenStatus.Code)
			// If caller doesn't explicitly set lead_score, follow status default score
			if req.LeadScore == nil {
				l.LeadScore = chosenStatus.Score
			}
		} else {
			// Explicitly clear status if empty values were sent (rare). Keep legacy consistent.
			l.LeadStatusID = nil
			if req.LeadStatus != "" {
				l.LeadStatus = req.LeadStatus
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

	if err := s.leadRepo.Update(l); err != nil {
		return nil, err
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
	// Get lead
	l, err := s.leadRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLeadNotFound
		}
		return nil, err
	}

	// Validate lead status must be "qualified"
	if l.LeadStatus != "qualified" {
		return nil, ErrLeadCannotConvert
	}

	// Check if already converted
	if l.LeadStatus == "converted" || (l.OpportunityID != nil && *l.OpportunityID != "") {
		return nil, ErrLeadAlreadyConverted
	}

	// Validate stage exists
	stage, err := s.pipelineRepo.FindStageByID(req.StageID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStageNotFound
		}
		return nil, err
	}

	var accountID string
	var contactID string
	var createdAccount interface{}
	var createdContact interface{}

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
	} else if req.AccountID != "" {
		// Use account from request (fallback)
		_, err := s.accountRepo.FindByID(req.AccountID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrAccountCreationFailed
			}
			return nil, err
		}
		accountID = req.AccountID
	} else {
		// Create account automatically. Deal.account_id is required, so conversion
		// must not depend on company_name being present.
		categories, err := s.categoryRepo.List()
		if err != nil || len(categories) == 0 {
			return nil, ErrAccountCreationFailed
		}

		account := &account.Account{
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
			account.AssignedTo = l.AssignedTo
		}

		if err := s.accountRepo.Create(account); err != nil {
			return nil, ErrAccountCreationFailed
		}

		accountID = account.ID
		createdAccount = account.ToAccountResponse()
	}
	if accountID == "" {
		return nil, ErrAccountCreationFailed
	}
	if createdAccount == nil {
		updatedAccount, err := s.syncAccountFromLead(accountID, l, accountName)
		if err != nil {
			return nil, ErrAccountCreationFailed
		}
		createdAccount = updatedAccount.ToAccountResponse()
	}

	// Always create contact if account exists
	// (ignore CreateContact flag, make it automatic)
	if accountID != "" && (l.ContactID == nil || *l.ContactID == "") {
		// Find default contact role (you may need to adjust this logic)
		contactRoles, err := s.contactRoleRepo.List()
		if err != nil || len(contactRoles) == 0 {
			return nil, ErrContactCreationFailed
		}

		contactName := leadFullName
		if contactName == "" {
			contactName = accountName
		}

		contact := &contact.Contact{
			AccountID: accountID,
			Name:      contactName,
			RoleID:    contactRoles[0].ID,
			Email:     l.Email,
			Phone:     l.Phone,
			Position:  l.JobTitle,
		}

		if err := s.contactRepo.Create(contact); err != nil {
			return nil, ErrContactCreationFailed
		}

		contactID = contact.ID
		createdContact = contact.ToContactResponse()
	} else if l.ContactID != nil && *l.ContactID != "" {
		// Use existing contact from lead
		contactID = *l.ContactID
	} else if req.ContactID != "" {
		// Use contact from request (fallback)
		_, err := s.contactRepo.FindByID(req.ContactID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrContactCreationFailed
			}
			return nil, err
		}
		contactID = req.ContactID
	}

	// Create deal/opportunity
	dealValue := int64(0)
	if req.Value != nil {
		dealValue = *req.Value
	}

	dealStatus := "open"
	probability := stage.Order * 20
	conversionTime := time.Now()
	var actualCloseDate *time.Time
	if stage.Probability > 0 {
		probability = stage.Probability
	}
	if stage.IsWon {
		dealStatus = "won"
		actualCloseDate = &conversionTime
	} else if stage.IsLost {
		dealStatus = "lost"
		actualCloseDate = &conversionTime
	}

	budgetConfirmed := l.BudgetConfirmed
	authorityConfirmed := l.AuthorityConfirmed
	needConfirmed := l.NeedConfirmed
	timelineConfirmed := l.TimelineConfirmed
	qualificationSnapshot := []byte("{}")
	needProducts := make([]leadqualification.NeedProduct, 0)
	var qualification leadqualification.LeadQualificationChecklist
	if err := s.db.Where("lead_id = ?", l.ID).First(&qualification).Error; err == nil {
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

	deal := &pipeline.Deal{
		Title:           req.OpportunityTitle,
		Description:     req.OpportunityDescription,
		AccountID:       accountID,
		ContactID:       stringPtr(contactID),
		StageID:         stage.ID,
		Value:           dealValue,
		Probability:     probability,
		ActualCloseDate: actualCloseDate,
		AssignedTo:      l.AssignedTo,
		LeadID:          &l.ID, // Set LeadID to track source lead
		Status:          dealStatus,
		Source:          l.LeadSource,
		// Copy BANT qualification fields from Lead to Deal
		BudgetConfirmed:       budgetConfirmed,
		AuthorityConfirmed:    authorityConfirmed,
		NeedConfirmed:         needConfirmed,
		TimelineConfirmed:     timelineConfirmed,
		QualificationSnapshot: qualificationSnapshot,
		Notes:                 l.Notes,
		CreatedBy:             convertedBy,
	}

	if err := s.dealRepo.Create(deal); err != nil {
		return nil, ErrOpportunityCreationFailed
	}

	if len(needProducts) > 0 {
		s.createDealProductItemsFromLeadNeeds(deal.ID, needProducts)
	}

	// Reload deal to get relations
	deal, err = s.dealRepo.FindByID(deal.ID)
	if err != nil {
		return nil, ErrOpportunityCreationFailed
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
			FromStageID:     nil, // NULL for initial creation
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
		_ = s.dealHistoryRepo.Create(history) // Ignore error
	}

	// Create initial deal task
	if s.taskRepo != nil {
		dueDate := time.Now().AddDate(0, 0, 3) // Due in 3 days
		initialTask := &task.Task{
			Title:       "Conduct needs analysis meeting",
			Description: "Deep dive into customer requirements and technical specifications",
			Type:        "meeting",
			Status:      "pending",
			Priority:    "high",
			DueDate:     &dueDate,
			DealID:      &deal.ID,
			AccountID:   &accountID,
			CreatedBy:   convertedBy,
		}
		initialTask.AssignedTo = l.AssignedTo
		_ = s.taskRepo.Create(initialTask) // Ignore error
	}

	// Update lead status to converted
	now := conversionTime
	l.LeadStatus = "converted"
	if convertedStatus, err := s.findConvertedLeadStatus(); err != nil {
		return nil, err
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

	if err := s.leadRepo.Update(l); err != nil {
		return nil, err
	}

	// FIXED: Use batch UPDATE instead of N+1 loop for activities and visit reports
	// Auto-migrate Activities: Update all activities linked to this lead using batch query
	if s.activityRepo != nil {
		dealIDStr := deal.ID
		var accountIDPtr *string
		if accountID != "" {
			accountIDPtr = &accountID
		}
		// Batch update - single query instead of N updates
		_ = s.activityRepo.UpdateByLeadID(l.ID, &dealIDStr, accountIDPtr)
	}

	// Auto-migrate Visit Reports: Update all visit reports linked to this lead using batch query
	if s.visitReportRepo != nil {
		dealIDStr := deal.ID
		var accountIDPtr *string
		if accountID != "" {
			accountIDPtr = &accountID
		}
		// Batch update - single query instead of N updates
		_ = s.visitReportRepo.UpdateByLeadID(l.ID, &dealIDStr, accountIDPtr)
	}

	if s.taskRepo != nil {
		dealIDStr := deal.ID
		var accountIDPtr *string
		if accountID != "" {
			accountIDPtr = &accountID
		}
		_ = s.taskRepo.UpdateByLeadID(l.ID, &dealIDStr, accountIDPtr)
	}

	// Reload lead to get relations
	l, err = s.leadRepo.FindByID(l.ID)
	if err != nil {
		return nil, err
	}

	// Emit lead converted event
	if s.eventHelper != nil {
		s.eventHelper.EmitLeadConverted(&domainevents.LeadConvertedEvent{
			LeadID:        l.ID,
			OpportunityID: deal.ID,
			AccountID:     accountID,
			ContactID:     contactID,
			ConvertedBy:   convertedBy,
			ConvertedAt:   time.Now(),
		}, convertedBy)
	}

	response := &lead.ConvertLeadResponse{
		Lead:        l.ToLeadResponse(),
		Opportunity: deal.ToDealResponse(),
	}

	if createdAccount != nil {
		response.Account = createdAccount
	}
	if createdContact != nil {
		response.Contact = createdContact
	}

	return response, nil
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

func (s *Service) syncAccountFromLead(accountID string, l *lead.Lead, fallbackName string) (*account.Account, error) {
	accountEntity, err := s.accountRepo.FindByID(accountID)
	if err != nil {
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

	if l.AssignedTo != nil && *l.AssignedTo != "" {
		if accountEntity.AssignedTo == nil || *accountEntity.AssignedTo != *l.AssignedTo {
			accountEntity.AssignedTo = l.AssignedTo
			changed = true
		}
	}

	if changed {
		if err := s.accountRepo.Update(accountEntity); err != nil {
			return nil, err
		}
	}

	return accountEntity, nil
}

func (s *Service) createDealProductItemsFromLeadNeeds(dealID string, needProducts []leadqualification.NeedProduct) {
	if s.db == nil {
		return
	}

	for _, needProduct := range needProducts {
		if needProduct.ProductID == "" {
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

		err := s.db.Table("products AS p").
			Select("p.id, p.name, p.sku, p.price, p.cost, p.category_id, COALESCE(pc.name, '') AS category_name").
			Joins("LEFT JOIN product_categories pc ON pc.id = p.category_id AND pc.deleted_at IS NULL").
			Where("p.id = ? AND p.deleted_at IS NULL", needProduct.ProductID).
			Scan(&productSnapshot).Error
		if err != nil || productSnapshot.ID == "" {
			continue
		}

		item := &pipeline.DealProductItem{
			DealID:              dealID,
			ProductID:           productSnapshot.ID,
			ProductName:         productSnapshot.Name,
			ProductSKU:          productSnapshot.SKU,
			UnitPrice:           productSnapshot.Price,
			UnitCost:            productSnapshot.Cost,
			Quantity:            1,
			DiscountAmount:      0,
			Subtotal:            productSnapshot.Price,
			ProductCategoryID:   productSnapshot.CategoryID,
			ProductCategoryName: productSnapshot.CategoryName,
		}
		_ = s.db.Create(item).Error
	}
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
		{Value: "qualified", Label: "Qualified"},
		{Value: "unqualified", Label: "Unqualified"},
		{Value: "nurturing", Label: "Nurturing"},
		{Value: "disqualified", Label: "Disqualified"},
		{Value: "converted", Label: "Converted"},
		{Value: "lost", Label: "Lost"},
	}

	// Get active users for assigned_to
	userReq := &user.ListUsersRequest{
		Page:    1,
		PerPage: 100,
		Status:  "active",
	}
	users, _, err := s.userRepo.List(userReq)
	if err != nil {
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

// GetMobileFormData returns optimized form data for mobile lead creation
func (s *Service) GetMobileFormData() (*lead.LeadMobileFormDataResponse, error) {
	// Get Lead Sources from DB
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

	// Get Lead Statuses
	// Try to get from repository first
	var leadStatuses []lead.LeadStatusOption
	statuses, err := s.leadStatusRepo.ListAll()
	if err == nil && len(statuses) > 0 {
		for _, st := range statuses {
			leadStatuses = append(leadStatuses, lead.LeadStatusOption{
				ID:    st.ID,
				Value: st.Code,
				Label: st.Name,
			})
		}
	} else {
		// Fallback to hardcoded constants if repo fails or empty
		leadStatuses = []lead.LeadStatusOption{
			{Value: "new", Label: "New"},
			{Value: "contacted", Label: "Contacted"},
			{Value: "qualified", Label: "Qualified"},
			{Value: "unqualified", Label: "Unqualified"},
			{Value: "nurturing", Label: "Nurturing"},
			{Value: "disqualified", Label: "Disqualified"},
			{Value: "converted", Label: "Converted"},
			{Value: "lost", Label: "Lost"},
		}
	}

	// Get Industries from DB
	var dbIndustries []industry.Industry
	if err := s.db.Where("is_active = ?", true).Order("\"order\" ASC").Find(&dbIndustries).Error; err != nil {
		dbIndustries = []industry.Industry{}
	}

	industries := make([]string, len(dbIndustries))
	for i, ind := range dbIndustries {
		industries[i] = ind.Name
	}

	// Hardware Provinces (copied from GetFormData)
	provinces := []string{
		"DKI Jakarta",
		"Jawa Barat",
		"Jawa Tengah",
		"Jawa Timur",
		"Yogyakarta",
		"Banten",
		"Bali",
		"Nusa Tenggara Barat",
		"Nusa Tenggara Timur",
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

	return &lead.LeadMobileFormDataResponse{
		LeadSources:  leadSources,
		LeadStatuses: leadStatuses,
		Industries:   industries,
		Provinces:    provinces,
	}, nil
}
