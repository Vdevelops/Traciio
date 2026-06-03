package visit_report

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/account"
	"github.com/gilabs/crm-healthcare/api/internal/domain/activity"
	"github.com/gilabs/crm-healthcare/api/internal/domain/contact"
	"github.com/gilabs/crm-healthcare/api/internal/domain/task"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"github.com/gilabs/crm-healthcare/api/internal/domain/visit_report"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	brick "github.com/gilabs/crm-healthcare/api/internal/service/brick"
	"github.com/gilabs/crm-healthcare/api/pkg/cache"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var (
	ErrVisitReportNotFound = errors.New("visit report not found")
	ErrAccountNotFound     = errors.New("account not found")
	ErrInvalidStatus       = errors.New("invalid status transition")
	ErrLeadNotFound        = errors.New("lead not found")
	ErrNotOwner            = errors.New("not the owner of this visit report")
	ErrInvalidGPS          = errors.New("invalid GPS data or GPS spoofing detected")
	ErrSubmitPrerequisite  = errors.New("submit prerequisites not met")
)

type Service struct {
	visitReportRepo  interfaces.VisitReportRepository
	accountRepo      interfaces.AccountRepository
	contactRepo      interfaces.ContactRepository
	userRepo         interfaces.UserRepository
	activityRepo     interfaces.ActivityRepository
	activityTypeRepo interfaces.ActivityTypeRepository
	leadRepo         interfaces.LeadRepository
	taskRepo         interfaces.TaskRepository
	notificationRepo interfaces.NotificationRepository
	brickHelper      *brick.BrickHelper
	cacheService     *cache.VisitReportCacheService
	db               *gorm.DB
}

func NewService(
	visitReportRepo interfaces.VisitReportRepository,
	accountRepo interfaces.AccountRepository,
	contactRepo interfaces.ContactRepository,
	userRepo interfaces.UserRepository,
	activityRepo interfaces.ActivityRepository,
	activityTypeRepo interfaces.ActivityTypeRepository,
	leadRepo interfaces.LeadRepository,
	taskRepo interfaces.TaskRepository,
	notificationRepo interfaces.NotificationRepository,
	brickHelper *brick.BrickHelper,
	db *gorm.DB,
) *Service {
	return &Service{
		visitReportRepo:  visitReportRepo,
		accountRepo:      accountRepo,
		contactRepo:      contactRepo,
		userRepo:         userRepo,
		activityRepo:     activityRepo,
		activityTypeRepo: activityTypeRepo,
		leadRepo:         leadRepo,
		taskRepo:         taskRepo,
		notificationRepo: notificationRepo,
		brickHelper:      brickHelper,
		cacheService:     cache.NewVisitReportCacheService(nil),
		db:               db,
	}
}

// cachedVisitReportListResult for msgpack serialization
type cachedVisitReportListResult struct {
	VisitReports []visit_report.VisitReportResponse `msgpack:"visit_reports"`
	Pagination   *PaginationResult                  `msgpack:"pagination"`
}

// PaginationResult represents pagination information
type PaginationResult struct {
	Page       int
	PerPage    int
	Total      int
	TotalPages int
	Offset     int  // For offset-based pagination (infinity scroll)
	HasMore    bool // Indicates if there are more items to load (for infinity scroll)
}

// loadRelations loads Account, Contact, Deal, Lead, and SalesRep relations into response
func (s *Service) loadRelations(response *visit_report.VisitReportResponse, vr *visit_report.VisitReport) {
	// Load Account (if AccountID is provided)
	if vr.AccountID != nil && *vr.AccountID != "" {
		if account, err := s.accountRepo.FindByID(*vr.AccountID); err == nil {
			response.Account = map[string]interface{}{
				"id":   account.ID,
				"name": account.Name,
			}
		}
	}
	// Load Contact
	if vr.ContactID != nil && *vr.ContactID != "" {
		if contact, err := s.contactRepo.FindByID(*vr.ContactID); err == nil {
			response.Contact = map[string]interface{}{
				"id":   contact.ID,
				"name": contact.Name,
			}
		}
	}
	// Load Lead (if LeadID is provided)
	if vr.LeadID != nil && *vr.LeadID != "" && s.leadRepo != nil {
		if lead, err := s.leadRepo.FindByID(*vr.LeadID); err == nil {
			response.Lead = map[string]interface{}{
				"id":           lead.ID,
				"first_name":   lead.FirstName,
				"last_name":    lead.LastName,
				"company_name": lead.CompanyName,
			}
		}
	}
	// Load SalesRep (User)
	if user, err := s.userRepo.FindByID(vr.SalesRepID); err == nil {
		response.SalesRep = map[string]interface{}{
			"id":   user.ID,
			"name": user.Name,
		}
	}
}

// List returns a list of visit reports with pagination
func (s *Service) List(req *visit_report.ListVisitReportsRequest) ([]visit_report.VisitReportResponse, *PaginationResult, error) {
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}

	// Calculate page and offset
	var page int
	var offset int
	if req.Offset > 0 {
		// Offset-based pagination (for infinity scroll)
		offset = req.Offset
		page = (offset / perPage) + 1
	} else {
		// Page-based pagination (default)
		page = req.Page
		if page < 1 {
			page = 1
		}
		offset = (page - 1) * perPage
	}

	// Try cache first (only for page-based pagination, not offset-based)
	if req.Offset == 0 {
		filterMap := map[string]interface{}{
			"search":       req.Search,
			"status":       req.Status,
			"account_id":   req.AccountID,
			"deal_id":      req.DealID,
			"lead_id":      req.LeadID,
			"sales_rep_id": req.SalesRepID,
			"brick_id":     req.BrickID,
			"start_date":   req.StartDate,
			"end_date":     req.EndDate,
		}
		if len(req.ScopedUserIDs) > 0 {
			filterMap["scoped_user_ids"] = strings.Join(req.ScopedUserIDs, ",")
		}
		var cachedResult cachedVisitReportListResult
		if found, _ := s.cacheService.GetList(page, perPage, filterMap, &cachedResult); found && cachedResult.Pagination != nil {
			return cachedResult.VisitReports, cachedResult.Pagination, nil
		}
	}

	visitReports, total, err := s.visitReportRepo.List(req)
	if err != nil {
		return nil, nil, err
	}

	// CRITICAL: Batch load relations to avoid N+1 queries
	// Collect all unique IDs first
	accountIDs := make(map[string]bool)
	contactIDs := make(map[string]bool)
	dealIDs := make(map[string]bool)
	userIDs := make(map[string]bool)

	for _, vr := range visitReports {
		if vr.AccountID != nil && *vr.AccountID != "" {
			accountIDs[*vr.AccountID] = true
		}
		if vr.ContactID != nil && *vr.ContactID != "" {
			contactIDs[*vr.ContactID] = true
		}
		if vr.DealID != nil && *vr.DealID != "" {
			dealIDs[*vr.DealID] = true
		}
		if vr.SalesRepID != "" {
			userIDs[vr.SalesRepID] = true
		}
	}

	// Batch load all accounts
	accountsMap := make(map[string]*account.Account)
	if len(accountIDs) > 0 {
		accountIDList := make([]string, 0, len(accountIDs))
		for id := range accountIDs {
			accountIDList = append(accountIDList, id)
		}
		// Load accounts in batch (using repository if batch method exists, otherwise fallback to individual)
		for _, id := range accountIDList {
			if acc, err := s.accountRepo.FindByID(id); err == nil {
				accountsMap[id] = acc
			}
		}
	}

	// Batch load all contacts
	contactsMap := make(map[string]*contact.Contact)
	if len(contactIDs) > 0 {
		contactIDList := make([]string, 0, len(contactIDs))
		for id := range contactIDs {
			contactIDList = append(contactIDList, id)
		}
		for _, id := range contactIDList {
			if c, err := s.contactRepo.FindByID(id); err == nil {
				contactsMap[id] = c
			}
		}
	}

	// Batch load all deals
	dealsMap := make(map[string]interface{})
	if len(dealIDs) > 0 {
		dealIDList := make([]string, 0, len(dealIDs))
		for id := range dealIDs {
			dealIDList = append(dealIDList, id)
		}

		// Query deals in batch
		var deals []struct {
			ID    string `gorm:"column:id"`
			Title string `gorm:"column:title"`
		}
		if err := s.db.Table("deals").
			Select("id, title").
			Where("id IN ? AND deleted_at IS NULL", dealIDList).
			Scan(&deals).Error; err == nil {
			// Map deals by ID
			for _, deal := range deals {
				dealsMap[deal.ID] = map[string]interface{}{
					"id":    deal.ID,
					"title": deal.Title,
				}
			}
		}
	}

	// Batch load all users
	usersMap := make(map[string]*user.User)
	if len(userIDs) > 0 {
		userIDList := make([]string, 0, len(userIDs))
		for id := range userIDs {
			userIDList = append(userIDList, id)
		}
		for _, id := range userIDList {
			if u, err := s.userRepo.FindByID(id); err == nil {
				usersMap[id] = u
			}
		}
	}

	// Build responses with pre-loaded relations
	responses := make([]visit_report.VisitReportResponse, len(visitReports))
	for i, vr := range visitReports {
		response := *vr.ToVisitReportResponse()
		// Parse photos JSON
		if vr.Photos != nil {
			var photos []string
			if err := json.Unmarshal(vr.Photos, &photos); err == nil {
				response.Photos = photos
			}
		}
		if vr.Metadata != nil {
			var metadata interface{}
			if err := json.Unmarshal(vr.Metadata, &metadata); err == nil {
				response.Metadata = metadata
			}
		}
		// Parse check-in location JSON
		if vr.CheckInLocation != nil {
			var location visit_report.Location
			if err := json.Unmarshal(vr.CheckInLocation, &location); err == nil {
				response.CheckInLocation = &location
			}
		}
		// Parse check-out location JSON
		if vr.CheckOutLocation != nil {
			var location visit_report.Location
			if err := json.Unmarshal(vr.CheckOutLocation, &location); err == nil {
				response.CheckOutLocation = &location
			}
		}
		// Determine type based on lead_id, deal_id, or account_id (priority: lead > deal > account)
		if vr.LeadID != nil && *vr.LeadID != "" {
			response.Type = "lead"
		} else if vr.DealID != nil && *vr.DealID != "" {
			response.Type = "deal"
		} else {
			response.Type = "account"
		}

		// Load relations from pre-loaded maps (no additional queries)
		if vr.AccountID != nil && *vr.AccountID != "" {
			if acc, ok := accountsMap[*vr.AccountID]; ok {
				response.Account = map[string]interface{}{
					"id":   acc.ID,
					"name": acc.Name,
				}
			}
		}
		if vr.ContactID != nil && *vr.ContactID != "" {
			if c, ok := contactsMap[*vr.ContactID]; ok {
				response.Contact = map[string]interface{}{
					"id":   c.ID,
					"name": c.Name,
				}
			}
		}
		if vr.SalesRepID != "" {
			if u, ok := usersMap[vr.SalesRepID]; ok {
				response.SalesRep = map[string]interface{}{
					"id":   u.ID,
					"name": u.Name,
				}
			}
		}
		// Load deal relation from pre-loaded map
		if vr.DealID != nil && *vr.DealID != "" {
			if deal, ok := dealsMap[*vr.DealID]; ok {
				response.Deal = deal
			}
		}
		responses[i] = response
	}

	// Calculate has_more for offset-based pagination (infinity scroll)
	hasMore := (offset + perPage) < int(total)

	pagination := &PaginationResult{
		Page:       page,
		PerPage:    perPage,
		Total:      int(total),
		TotalPages: int((total + int64(perPage) - 1) / int64(perPage)),
		Offset:     offset,
		HasMore:    hasMore,
	}

	// Only cache page-based pagination, not offset-based
	if req.Offset == 0 {
		filterMap := map[string]interface{}{
			"search":       req.Search,
			"status":       req.Status,
			"account_id":   req.AccountID,
			"deal_id":      req.DealID,
			"lead_id":      req.LeadID,
			"sales_rep_id": req.SalesRepID,
			"brick_id":     req.BrickID,
			"start_date":   req.StartDate,
			"end_date":     req.EndDate,
		}
		if len(req.ScopedUserIDs) > 0 {
			filterMap["scoped_user_ids"] = strings.Join(req.ScopedUserIDs, ",")
		}
		_ = s.cacheService.SetList(page, perPage, filterMap, cachedVisitReportListResult{
			VisitReports: responses,
			Pagination:   pagination,
		})
	}

	return responses, pagination, nil
}

// GetByID returns a visit report by ID
func (s *Service) GetByID(id string) (*visit_report.VisitReportResponse, error) {
	// Try cache first
	var cachedResponse visit_report.VisitReportResponse
	if found, _ := s.cacheService.GetDetail(id, &cachedResponse); found && cachedResponse.ID != "" {
		return &cachedResponse, nil
	}

	vr, err := s.visitReportRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVisitReportNotFound
		}
		return nil, err
	}

	response := *vr.ToVisitReportResponse()

	// Determine type based on lead_id, deal_id, or account_id (priority: lead > deal > account)
	if vr.LeadID != nil && *vr.LeadID != "" {
		response.Type = "lead"
	} else if vr.DealID != nil && *vr.DealID != "" {
		response.Type = "deal"
	} else {
		response.Type = "account"
	}

	// Parse photos JSON
	if vr.Photos != nil {
		var photos []string
		if err := json.Unmarshal(vr.Photos, &photos); err == nil {
			response.Photos = photos
		}
	}
	if vr.Metadata != nil {
		var metadata interface{}
		if err := json.Unmarshal(vr.Metadata, &metadata); err == nil {
			response.Metadata = metadata
		}
	}
	// Parse check-in location JSON
	if vr.CheckInLocation != nil {
		var location visit_report.Location
		if err := json.Unmarshal(vr.CheckInLocation, &location); err == nil {
			response.CheckInLocation = &location
		}
	}
	// Parse check-out location JSON
	if vr.CheckOutLocation != nil {
		var location visit_report.Location
		if err := json.Unmarshal(vr.CheckOutLocation, &location); err == nil {
			response.CheckOutLocation = &location
		}
	}
	// Load relations
	s.loadRelations(&response, vr)
	_ = s.cacheService.SetDetail(id, &response)
	return &response, nil
}

// Create creates a new visit report
func (s *Service) Create(req *visit_report.CreateVisitReportRequest) (*visit_report.VisitReportResponse, error) {
	// Validate SalesRepID
	if req.SalesRepID == "" {
		return nil, errors.New("sales_rep_id is required")
	}

	// Business rule validation: Either LeadID or AccountID is required
	hasLeadID := req.LeadID != nil && *req.LeadID != ""
	hasAccountID := req.AccountID != nil && *req.AccountID != ""
	hasDealID := req.DealID != nil && *req.DealID != ""

	if !hasLeadID && !hasAccountID {
		return nil, errors.New("either lead_id or account_id is required")
	}

	// If DealID is provided, AccountID must be provided (post-conversion phase)
	if hasDealID && !hasAccountID {
		return nil, errors.New("account_id is required when deal_id is provided")
	}

	// Verify account exists if provided
	var err error
	if hasAccountID {
		_, err = s.accountRepo.FindByID(*req.AccountID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrAccountNotFound
			}
			return nil, err
		}
	}

	// Parse visit date (support both "YYYY-MM-DD" and "YYYY-MM-DD HH:mm" formats)
	var visitDate time.Time
	if len(req.VisitDate) > 10 {
		// Format with time: "YYYY-MM-DD HH:mm"
		visitDate, err = time.Parse("2006-01-02 15:04", req.VisitDate)
		if err != nil {
			// Try alternative format "2006-01-02T15:04:05"
			visitDate, err = time.Parse("2006-01-02T15:04:05", req.VisitDate)
		}
		if err != nil {
			// Try ISO format
			visitDate, err = time.Parse(time.RFC3339, req.VisitDate)
		}
	} else {
		// Format without time: "YYYY-MM-DD"
		visitDate, err = time.Parse("2006-01-02", req.VisitDate)
	}
	if err != nil {
		return nil, errors.New("invalid visit_date format, expected YYYY-MM-DD or YYYY-MM-DD HH:mm")
	}

	// Marshal photos to JSON
	var photosJSON datatypes.JSON
	if len(req.Photos) > 0 {
		photosBytes, err := json.Marshal(req.Photos)
		if err != nil {
			return nil, err
		}
		photosJSON = photosBytes
	}

	metadataJSON := datatypes.JSON([]byte("{}"))
	if req.Metadata != nil {
		metadataBytes, err := json.Marshal(req.Metadata)
		if err != nil {
			return nil, err
		}
		metadataJSON = metadataBytes
	}

	// Marshal check-in location to JSON
	var checkInLocationJSON datatypes.JSON
	if req.CheckInLocation != nil {
		locationBytes, err := json.Marshal(req.CheckInLocation)
		if err != nil {
			return nil, err
		}
		checkInLocationJSON = locationBytes
	}

	// Marshal check-out location to JSON
	var checkOutLocationJSON datatypes.JSON
	if req.CheckOutLocation != nil {
		locationBytes, err := json.Marshal(req.CheckOutLocation)
		if err != nil {
			return nil, err
		}
		checkOutLocationJSON = locationBytes
	}

	// Auto-populate brick_id if not provided
	var brickID *string
	if s.brickHelper != nil {
		// Try to get brick_id from sales_rep_id user
		if req.SalesRepID != "" {
			brickID, _ = s.brickHelper.GetBrickIDFromUser(req.SalesRepID)
		}
		// If still nil, try to get from account
		if brickID == nil && hasAccountID {
			brickID, _ = s.brickHelper.GetBrickIDFromAccount(*req.AccountID)
		}
	}

	// Validate deal exists if provided
	if req.DealID != nil && *req.DealID != "" {
		// Note: We'll need to inject dealRepo if we want to validate
		// For now, we'll just set it and let database foreign key handle validation
	}

	vr := &visit_report.VisitReport{
		AccountID:        req.AccountID,
		ContactID:        req.ContactID,
		DealID:           req.DealID,
		LeadID:           req.LeadID,
		SalesRepID:       req.SalesRepID,
		BrickID:          brickID,
		VisitDate:        visitDate,
		Purpose:          req.Purpose,
		Notes:            req.Notes,
		CheckInLocation:  checkInLocationJSON,
		CheckOutLocation: checkOutLocationJSON,
		Photos:           photosJSON,
		Metadata:         metadataJSON,
		Status:           "pending",
	}

	if err := s.visitReportRepo.Create(vr); err != nil {
		return nil, err
	}
	_ = s.cacheService.InvalidateOnWrite(vr.ID)

	// Reload
	createdVR, err := s.visitReportRepo.FindByID(vr.ID)
	if err != nil {
		return nil, err
	}

	response := *createdVR.ToVisitReportResponse()
	if createdVR.Photos != nil {
		var photos []string
		if err := json.Unmarshal(createdVR.Photos, &photos); err == nil {
			response.Photos = photos
		}
	}
	if createdVR.Metadata != nil {
		var metadata interface{}
		if err := json.Unmarshal(createdVR.Metadata, &metadata); err == nil {
			response.Metadata = metadata
		}
	}
	// Parse check-in location JSON
	if createdVR.CheckInLocation != nil {
		var location visit_report.Location
		if err := json.Unmarshal(createdVR.CheckInLocation, &location); err == nil {
			response.CheckInLocation = &location
		}
	}
	// Parse check-out location JSON
	if createdVR.CheckOutLocation != nil {
		var location visit_report.Location
		if err := json.Unmarshal(createdVR.CheckOutLocation, &location); err == nil {
			response.CheckOutLocation = &location
		}
	}
	// Load relations
	s.loadRelations(&response, createdVR)
	_ = s.cacheService.SetDetail(vr.ID, &response)

	return &response, nil
}

// Update updates a visit report
func (s *Service) Update(id string, req *visit_report.UpdateVisitReportRequest) (*visit_report.VisitReportResponse, error) {
	vr, err := s.visitReportRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVisitReportNotFound
		}
		return nil, err
	}

	// Business rule validation for update
	hasLeadID := req.LeadID != nil && *req.LeadID != ""
	hasAccountID := req.AccountID != nil && *req.AccountID != ""
	hasDealID := req.DealID != nil && *req.DealID != ""

	// Determine current state
	currentHasLeadID := vr.LeadID != nil && *vr.LeadID != ""
	currentHasAccountID := vr.AccountID != nil && *vr.AccountID != ""

	// After update, must have either LeadID or AccountID
	finalHasLeadID := hasLeadID || (!hasLeadID && currentHasLeadID)
	finalHasAccountID := hasAccountID || (!hasAccountID && currentHasAccountID)

	if !finalHasLeadID && !finalHasAccountID {
		return nil, errors.New("either lead_id or account_id is required")
	}

	// If DealID is provided, AccountID must be provided
	if hasDealID && !finalHasAccountID {
		return nil, errors.New("account_id is required when deal_id is provided")
	}

	// Track if fields that affect brick_id are being updated
	accountIDChanged := false

	// Update fields if provided
	if req.AccountID != nil {
		// Verify account exists
		_, err := s.accountRepo.FindByID(*req.AccountID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrAccountNotFound
			}
			return nil, err
		}
		// Check if account_id is actually changing
		if vr.AccountID == nil || *vr.AccountID != *req.AccountID {
			accountIDChanged = true
		}
		vr.AccountID = req.AccountID
	}

	if req.ContactID != nil {
		vr.ContactID = req.ContactID
	}

	if req.DealID != nil {
		vr.DealID = req.DealID
	}

	if req.LeadID != nil {
		vr.LeadID = req.LeadID
	}

	if req.VisitDate != "" {
		var visitDate time.Time
		var err error
		if len(req.VisitDate) > 10 {
			// Format with time: "YYYY-MM-DD HH:mm"
			visitDate, err = time.Parse("2006-01-02 15:04", req.VisitDate)
			if err != nil {
				// Try alternative format "2006-01-02T15:04:05"
				visitDate, err = time.Parse("2006-01-02T15:04:05", req.VisitDate)
			}
			if err != nil {
				// Try ISO format
				visitDate, err = time.Parse(time.RFC3339, req.VisitDate)
			}
		} else {
			// Format without time: "YYYY-MM-DD"
			visitDate, err = time.Parse("2006-01-02", req.VisitDate)
		}
		if err != nil {
			return nil, errors.New("invalid visit_date format, expected YYYY-MM-DD or YYYY-MM-DD HH:mm")
		}
		vr.VisitDate = visitDate
	}

	if req.Purpose != "" {
		vr.Purpose = req.Purpose
	}

	if req.Notes != "" {
		vr.Notes = req.Notes
	}

	if req.CheckInLocation != nil {
		locationBytes, err := json.Marshal(req.CheckInLocation)
		if err != nil {
			return nil, err
		}
		vr.CheckInLocation = locationBytes
	}

	if req.CheckOutLocation != nil {
		locationBytes, err := json.Marshal(req.CheckOutLocation)
		if err != nil {
			return nil, err
		}
		vr.CheckOutLocation = locationBytes
	}

	if req.Photos != nil {
		photosBytes, err := json.Marshal(req.Photos)
		if err != nil {
			return nil, err
		}
		vr.Photos = photosBytes
	}

	if req.Metadata != nil {
		metadataBytes, err := json.Marshal(req.Metadata)
		if err != nil {
			return nil, err
		}
		vr.Metadata = metadataBytes
	}

	if req.Status != "" {
		normalizedStatus := visit_report.NormalizeStatus(req.Status)
		if normalizedStatus == "completed" && vr.CheckInTime == nil {
			return nil, ErrSubmitPrerequisite
		}
		vr.Status = normalizedStatus
	}

	// Auto-update brick_id if account_id changed
	// Note: sales_rep_id typically doesn't change after creation, but if it did, we'd also need to update brick_id
	if accountIDChanged && s.brickHelper != nil {
		var brickID *string
		// Try to get brick_id from sales_rep_id user first (sales rep doesn't change, but we check for consistency)
		if vr.SalesRepID != "" {
			brickID, _ = s.brickHelper.GetBrickIDFromUser(vr.SalesRepID)
		}
		// If still nil, try to get from account
		if brickID == nil && vr.AccountID != nil && *vr.AccountID != "" {
			brickID, _ = s.brickHelper.GetBrickIDFromAccount(*vr.AccountID)
		}
		vr.BrickID = brickID
	}

	if err := s.visitReportRepo.Update(vr); err != nil {
		return nil, err
	}
	_ = s.cacheService.InvalidateOnWrite(vr.ID)

	// Reload
	updatedVR, err := s.visitReportRepo.FindByID(vr.ID)
	if err != nil {
		return nil, err
	}

	response := *updatedVR.ToVisitReportResponse()
	if updatedVR.Photos != nil {
		var photos []string
		if err := json.Unmarshal(updatedVR.Photos, &photos); err == nil {
			response.Photos = photos
		}
	}
	if updatedVR.Metadata != nil {
		var metadata interface{}
		if err := json.Unmarshal(updatedVR.Metadata, &metadata); err == nil {
			response.Metadata = metadata
		}
	}
	// Parse check-in location JSON
	if updatedVR.CheckInLocation != nil {
		var location visit_report.Location
		if err := json.Unmarshal(updatedVR.CheckInLocation, &location); err == nil {
			response.CheckInLocation = &location
		}
	}
	// Parse check-out location JSON
	if updatedVR.CheckOutLocation != nil {
		var location visit_report.Location
		if err := json.Unmarshal(updatedVR.CheckOutLocation, &location); err == nil {
			response.CheckOutLocation = &location
		}
	}
	// Load relations
	s.loadRelations(&response, updatedVR)
	_ = s.cacheService.SetDetail(vr.ID, &response)

	return &response, nil
}

// Delete deletes a visit report
func (s *Service) Delete(id string) error {
	_, err := s.visitReportRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrVisitReportNotFound
		}
		return err
	}

	if err := s.visitReportRepo.Delete(id); err != nil {
		return err
	}
	_ = s.cacheService.InvalidateOnWrite(id)

	return nil
}

// CheckIn performs check-in for a visit report
func (s *Service) CheckIn(id string, req *visit_report.CheckInRequest, userID string) (*visit_report.VisitReportResponse, error) {
	vr, err := s.visitReportRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVisitReportNotFound
		}
		return nil, err
	}

	if vr.CheckInTime != nil {
		return nil, errors.New("already checked in")
	}

	// Validate GPS anti-spoofing
	if err := s.validateGPS(req); err != nil {
		// Wrap error to ensure it's recognized as GPS error
		if err.Error() != "" && (strings.Contains(err.Error(), "GPS") || strings.Contains(err.Error(), "accuracy") || strings.Contains(err.Error(), "timestamp") || strings.Contains(err.Error(), "location")) {
			return nil, err
		}
		return nil, ErrInvalidGPS
	}

	now := time.Now()
	vr.CheckInTime = &now

	// Marshal check-in location to JSON
	if req.Location != nil {
		locationBytes, err := json.Marshal(req.Location)
		if err != nil {
			return nil, err
		}
		vr.CheckInLocation = locationBytes
	}

	// Handle photo upload if provided
	if req.PhotoURL != nil && *req.PhotoURL != "" {
		var photos []string
		if vr.Photos != nil {
			if err := json.Unmarshal(vr.Photos, &photos); err != nil {
				photos = []string{}
			}
		}
		photos = append(photos, *req.PhotoURL)
		photosBytes, err := json.Marshal(photos)
		if err != nil {
			return nil, err
		}
		vr.Photos = photosBytes
	}

	if err := s.visitReportRepo.Update(vr); err != nil {
		return nil, err
	}
	_ = s.cacheService.InvalidateOnWrite(vr.ID)

	// Create activity
	s.createActivity(vr, "visit", "Checked in to visit")

	// Reload
	updatedVR, err := s.visitReportRepo.FindByID(vr.ID)
	if err != nil {
		return nil, err
	}

	response := *updatedVR.ToVisitReportResponse()
	if updatedVR.Photos != nil {
		var photos []string
		if err := json.Unmarshal(updatedVR.Photos, &photos); err == nil {
			response.Photos = photos
		}
	}
	// Parse check-in location JSON
	if updatedVR.CheckInLocation != nil {
		var location visit_report.Location
		if err := json.Unmarshal(updatedVR.CheckInLocation, &location); err == nil {
			response.CheckInLocation = &location
		}
	}
	// Parse check-out location JSON
	if updatedVR.CheckOutLocation != nil {
		var location visit_report.Location
		if err := json.Unmarshal(updatedVR.CheckOutLocation, &location); err == nil {
			response.CheckOutLocation = &location
		}
	}

	return &response, nil
}

// validateGPS validates GPS data to prevent spoofing
func (s *Service) validateGPS(req *visit_report.CheckInRequest) error {
	if req.Location == nil {
		return errors.New("location is required")
	}

	// Maximum allowed distance between device GPS and photo GPS (in meters)
	const maxDistanceMeters = 200.0 // 200 meters tolerance (increased for better UX)
	// Maximum allowed GPS accuracy (in meters) - reject if accuracy is too poor
	const maxAccuracyMeters = 100.0 // 100 meters max accuracy (increased for indoor/weak signal)
	// Maximum time difference between GPS capture and check-in (in seconds)
	const maxTimeDifferenceSeconds = 300 // 5 minutes tolerance (increased for better UX)

	now := time.Now().Unix()

	// Validate device GPS accuracy if provided
	if req.DeviceGPS != nil {
		// Check GPS accuracy - reject if accuracy is too poor (indicates fake GPS)
		if req.DeviceGPS.Accuracy > maxAccuracyMeters {
			return errors.New("GPS accuracy too poor (exceeds " + fmt.Sprintf("%.0f", maxAccuracyMeters) + "m). Accuracy: " + fmt.Sprintf("%.0f", req.DeviceGPS.Accuracy) + "m")
		}

		// Check timestamp - GPS should be recent (not from hours ago)
		if req.DeviceGPS.Timestamp > 0 {
			timeDiff := now - req.DeviceGPS.Timestamp
			if timeDiff < 0 {
				return errors.New("GPS timestamp is in the future")
			}
			if timeDiff > maxTimeDifferenceSeconds {
				return errors.New("GPS timestamp too old (exceeds " + fmt.Sprintf("%.0f", float64(maxTimeDifferenceSeconds)/60) + " minutes). Time difference: " + fmt.Sprintf("%.0f", float64(timeDiff)/60) + " minutes")
			}
		}

		// Validate that device GPS matches check-in location (within tolerance)
		distance := calculateDistance(
			req.Location.Latitude, req.Location.Longitude,
			req.DeviceGPS.Latitude, req.DeviceGPS.Longitude,
		)
		if distance > maxDistanceMeters {
			return errors.New("Device GPS location doesn't match check-in location. Distance: " + fmt.Sprintf("%.0f", distance) + "m (max: " + fmt.Sprintf("%.0f", maxDistanceMeters) + "m)")
		}
	}

	// Validate photo GPS if provided (EXIF GPS from photo)
	if req.PhotoGPS != nil {
		// Check timestamp - photo should be recent
		if req.PhotoGPS.Timestamp > 0 {
			timeDiff := now - req.PhotoGPS.Timestamp
			if timeDiff < 0 {
				return errors.New("Photo GPS timestamp is in the future")
			}
			if timeDiff > maxTimeDifferenceSeconds {
				return errors.New("Photo GPS timestamp too old (exceeds " + fmt.Sprintf("%.0f", float64(maxTimeDifferenceSeconds)/60) + " minutes)")
			}
		}

		// Validate that photo GPS matches check-in location (within tolerance)
		distance := calculateDistance(
			req.Location.Latitude, req.Location.Longitude,
			req.PhotoGPS.Latitude, req.PhotoGPS.Longitude,
		)
		if distance > maxDistanceMeters {
			return errors.New("Photo GPS location doesn't match check-in location. Distance: " + fmt.Sprintf("%.0f", distance) + "m (max: " + fmt.Sprintf("%.0f", maxDistanceMeters) + "m)")
		}

		// If both device GPS and photo GPS are provided, validate they match
		if req.DeviceGPS != nil {
			distance := calculateDistance(
				req.DeviceGPS.Latitude, req.DeviceGPS.Longitude,
				req.PhotoGPS.Latitude, req.PhotoGPS.Longitude,
			)
			if distance > maxDistanceMeters {
				return errors.New("Device GPS and Photo GPS don't match. Distance: " + fmt.Sprintf("%.0f", distance) + "m (max: " + fmt.Sprintf("%.0f", maxDistanceMeters) + "m)")
			}
		}
	}

	// If photo is provided but no GPS metadata, warn but don't reject
	// (some devices may not embed GPS in photos)
	if req.PhotoURL != nil && *req.PhotoURL != "" {
		if req.DeviceGPS == nil && req.PhotoGPS == nil {
			// Photo provided but no GPS metadata - this is suspicious but not blocking
			// In production, you might want to log this for review
		}
	}

	return nil
}

// calculateDistance calculates the distance between two GPS coordinates using Haversine formula
// Returns distance in meters
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusMeters = 6371000 // Earth radius in meters

	// Convert degrees to radians
	lat1Rad := lat1 * (3.141592653589793 / 180.0)
	lon1Rad := lon1 * (3.141592653589793 / 180.0)
	lat2Rad := lat2 * (3.141592653589793 / 180.0)
	lon2Rad := lon2 * (3.141592653589793 / 180.0)

	// Haversine formula
	dlat := lat2Rad - lat1Rad
	dlon := lon2Rad - lon1Rad

	a := (dlat/2)*(dlat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			(dlon/2)*(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusMeters * c
}

// CheckOut performs check-out for a visit report
func (s *Service) CheckOut(id string, req *visit_report.CheckOutRequest, userID string) (*visit_report.VisitReportResponse, error) {
	vr, err := s.visitReportRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVisitReportNotFound
		}
		return nil, err
	}

	if vr.CheckInTime == nil {
		return nil, errors.New("must check in first")
	}

	if vr.CheckOutTime != nil {
		return nil, errors.New("already checked out")
	}

	now := time.Now()
	vr.CheckOutTime = &now

	// Marshal check-out location to JSON
	if req.Location != nil {
		locationBytes, err := json.Marshal(req.Location)
		if err != nil {
			return nil, err
		}
		vr.CheckOutLocation = locationBytes
	}

	if err := s.visitReportRepo.Update(vr); err != nil {
		return nil, err
	}

	// Create activity
	s.createActivity(vr, "visit", "Checked out from visit")

	// Reload
	updatedVR, err := s.visitReportRepo.FindByID(vr.ID)
	if err != nil {
		return nil, err
	}

	response := *updatedVR.ToVisitReportResponse()
	if updatedVR.Photos != nil {
		var photos []string
		if err := json.Unmarshal(updatedVR.Photos, &photos); err == nil {
			response.Photos = photos
		}
	}
	// Parse check-in location JSON
	if updatedVR.CheckInLocation != nil {
		var location visit_report.Location
		if err := json.Unmarshal(updatedVR.CheckInLocation, &location); err == nil {
			response.CheckInLocation = &location
		}
	}
	// Parse check-out location JSON
	if updatedVR.CheckOutLocation != nil {
		var location visit_report.Location
		if err := json.Unmarshal(updatedVR.CheckOutLocation, &location); err == nil {
			response.CheckOutLocation = &location
		}
	}
	// Load relations
	s.loadRelations(&response, updatedVR)

	return &response, nil
}

// Approve approves a visit report
func (s *Service) Approve(id string, userID string) (*visit_report.VisitReportResponse, error) {
	vr, err := s.visitReportRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVisitReportNotFound
		}
		return nil, err
	}

	if visit_report.NormalizeStatus(vr.Status) != "pending" {
		return nil, ErrInvalidStatus
	}

	now := time.Now()
	vr.Status = "completed"
	vr.ApprovedBy = &userID
	vr.ApprovedAt = &now

	if err := s.visitReportRepo.Update(vr); err != nil {
		return nil, err
	}

	// Create activity
	s.createActivity(vr, "visit", "Visit report approved")

	// Reload
	updatedVR, err := s.visitReportRepo.FindByID(vr.ID)
	if err != nil {
		return nil, err
	}

	response := *updatedVR.ToVisitReportResponse()
	if updatedVR.Photos != nil {
		var photos []string
		if err := json.Unmarshal(updatedVR.Photos, &photos); err == nil {
			response.Photos = photos
		}
	}
	// Parse check-in location JSON
	if updatedVR.CheckInLocation != nil {
		var location visit_report.Location
		if err := json.Unmarshal(updatedVR.CheckInLocation, &location); err == nil {
			response.CheckInLocation = &location
		}
	}
	// Parse check-out location JSON
	if updatedVR.CheckOutLocation != nil {
		var location visit_report.Location
		if err := json.Unmarshal(updatedVR.CheckOutLocation, &location); err == nil {
			response.CheckOutLocation = &location
		}
	}
	// Load relations
	s.loadRelations(&response, updatedVR)

	return &response, nil
}

// Reject rejects a visit report
func (s *Service) Reject(id string, req *visit_report.RejectRequest, userID string) (*visit_report.VisitReportResponse, error) {
	vr, err := s.visitReportRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVisitReportNotFound
		}
		return nil, err
	}

	if visit_report.NormalizeStatus(vr.Status) != "pending" {
		return nil, ErrInvalidStatus
	}

	vr.Status = "completed"
	vr.RejectionReason = &req.Reason

	if err := s.visitReportRepo.Update(vr); err != nil {
		return nil, err
	}

	// Create activity
	s.createActivity(vr, "visit", "Visit report rejected: "+req.Reason)

	// Reload
	updatedVR, err := s.visitReportRepo.FindByID(vr.ID)
	if err != nil {
		return nil, err
	}

	response := *updatedVR.ToVisitReportResponse()
	if updatedVR.Photos != nil {
		var photos []string
		if err := json.Unmarshal(updatedVR.Photos, &photos); err == nil {
			response.Photos = photos
		}
	}
	// Parse check-in location JSON
	if updatedVR.CheckInLocation != nil {
		var location visit_report.Location
		if err := json.Unmarshal(updatedVR.CheckInLocation, &location); err == nil {
			response.CheckInLocation = &location
		}
	}
	// Parse check-out location JSON
	if updatedVR.CheckOutLocation != nil {
		var location visit_report.Location
		if err := json.Unmarshal(updatedVR.CheckOutLocation, &location); err == nil {
			response.CheckOutLocation = &location
		}
	}
	// Load relations
	s.loadRelations(&response, updatedVR)

	return &response, nil
}

// UploadPhoto adds a photo to a visit report
func (s *Service) UploadPhoto(id string, req *visit_report.UploadPhotoRequest) (*visit_report.VisitReportResponse, error) {
	vr, err := s.visitReportRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVisitReportNotFound
		}
		return nil, err
	}

	// Get existing photos
	var photos []string
	if vr.Photos != nil {
		if err := json.Unmarshal(vr.Photos, &photos); err != nil {
			photos = []string{}
		}
	}

	// Add new photo
	photos = append(photos, req.PhotoURL)

	// Marshal back to JSON
	photosBytes, err := json.Marshal(photos)
	if err != nil {
		return nil, err
	}
	vr.Photos = photosBytes

	if err := s.visitReportRepo.Update(vr); err != nil {
		return nil, err
	}

	// Reload
	updatedVR, err := s.visitReportRepo.FindByID(vr.ID)
	if err != nil {
		return nil, err
	}

	response := *updatedVR.ToVisitReportResponse()
	if updatedVR.Photos != nil {
		var photos []string
		if err := json.Unmarshal(updatedVR.Photos, &photos); err == nil {
			response.Photos = photos
		}
	}
	// Parse check-in location JSON
	if updatedVR.CheckInLocation != nil {
		var location visit_report.Location
		if err := json.Unmarshal(updatedVR.CheckInLocation, &location); err == nil {
			response.CheckInLocation = &location
		}
	}
	// Parse check-out location JSON
	if updatedVR.CheckOutLocation != nil {
		var location visit_report.Location
		if err := json.Unmarshal(updatedVR.CheckOutLocation, &location); err == nil {
			response.CheckOutLocation = &location
		}
	}
	// Load relations
	s.loadRelations(&response, updatedVR)

	return &response, nil
}

// GetMyVisitReports returns visit reports for the logged-in user (sales rep)
// If forRouteOptimization is true, returns all visit reports (not filtered by user)
func (s *Service) GetMyVisitReports(userID string, req *visit_report.ListVisitReportsRequest, forRouteOptimization bool) ([]visit_report.VisitReportResponse, *PaginationResult, error) {
	// For route optimization, don't filter by sales rep (allow viewing all visit reports)
	// Otherwise, override SalesRepID filter to only show visit reports for the user
	if !forRouteOptimization {
		req.SalesRepID = userID
	}

	return s.List(req)
}

// createActivity creates an activity record for a visit report
func (s *Service) createActivity(vr *visit_report.VisitReport, activityType, description string) {
	if s.activityRepo == nil {
		return // Skip if activity repo not available
	}

	var activityTypeID *string
	if s.activityTypeRepo != nil {
		if activityTypeEntity, err := s.activityTypeRepo.FindByCode(activityType); err == nil && activityTypeEntity != nil {
			activityTypeID = &activityTypeEntity.ID
		}
	}

	activity := &activity.Activity{
		Type:           activityType,
		ActivityTypeID: activityTypeID,
		AccountID:      vr.AccountID, // Already *string, can be nil
		ContactID:      vr.ContactID,
		DealID:         vr.DealID,
		LeadID:         vr.LeadID,
		UserID:         vr.SalesRepID,
		Description:    description,
		Timestamp:      time.Now(),
	}

	// Add metadata with visit report ID
	metadata := map[string]interface{}{
		"visit_report_id": vr.ID,
		"visit_date":      vr.VisitDate.Format("2006-01-02"),
	}
	if vr.Metadata != nil {
		var visitMetadata map[string]interface{}
		if err := json.Unmarshal(vr.Metadata, &visitMetadata); err == nil {
			if productInterests, ok := visitMetadata["product_interests"]; ok {
				metadata["product_interests"] = productInterests
			}
		}
	}
	if metadataBytes, err := json.Marshal(metadata); err == nil {
		activity.Metadata = datatypes.JSON(metadataBytes)
	}

	_ = s.activityRepo.Create(activity) // Ignore error for now
}

// Submit finalizes a visit report after the visit is completed.
func (s *Service) Submit(id string, req *visit_report.SubmitRequest, userID string) (*visit_report.VisitReportResponse, error) {
	vr, err := s.visitReportRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVisitReportNotFound
		}
		return nil, err
	}

	// Validate that only the owner can submit
	if vr.SalesRepID != userID {
		return nil, ErrNotOwner
	}

	if visit_report.NormalizeStatus(vr.Status) != "pending" {
		return nil, ErrInvalidStatus
	}

	if vr.CheckInTime == nil {
		return nil, ErrSubmitPrerequisite
	}

	vr.Status = "completed"

	// Update outcome and next_steps if provided
	if req.Outcome != "" {
		vr.Outcome = req.Outcome
	}
	if req.NextSteps != "" {
		vr.NextSteps = req.NextSteps
	}

	if err := s.visitReportRepo.Update(vr); err != nil {
		return nil, err
	}

	// AUTO-TRIGGERS AFTER SUBMIT

	// 1. Update Lead status if linked to a lead
	if vr.LeadID != nil && *vr.LeadID != "" && s.leadRepo != nil {
		lead, err := s.leadRepo.FindByID(*vr.LeadID)
		if err == nil && lead.LeadStatus == "new" {
			// Update lead status from "new" to "contacted"
			lead.LeadStatus = "contacted"
			_ = s.leadRepo.Update(lead) // Ignore error for now
		}
	}

	// 2. Auto-create tasks based on next_steps or default tasks
	if s.taskRepo != nil {
		s.createAutoTasks(vr)
	}

	// 3. Create activity record
	s.createActivity(vr, "visit", "Visit completed")

	// Reload and return response
	return s.GetByID(id)
}

// createAutoTasks creates automatic tasks after visit report submission
func (s *Service) createAutoTasks(vr *visit_report.VisitReport) {
	if s.taskRepo == nil {
		return
	}

	// Define default tasks to create
	tasksToCreate := []struct {
		title       string
		description string
		daysOffset  int
		priority    string
		taskType    string
	}{
		{
			title:       "Send product catalog",
			description: "Follow-up: Send detailed product catalog to customer",
			daysOffset:  1, // Due tomorrow
			priority:    "high",
			taskType:    "follow_up",
		},
		{
			title:       "Schedule product demo",
			description: "Schedule product demonstration with customer",
			daysOffset:  3, // Due in 3 days
			priority:    "medium",
			taskType:    "meeting",
		},
		{
			title:       "Follow-up call",
			description: "Make follow-up call to discuss visit outcomes",
			daysOffset:  2, // Due in 2 days
			priority:    "high",
			taskType:    "call",
		},
	}

	// Create tasks
	for _, taskDef := range tasksToCreate {
		dueDate := time.Now().AddDate(0, 0, taskDef.daysOffset)

		// Check if task already exists to avoid duplicates (idempotency)
		listReq := &task.ListTasksRequest{
			Page:    1,
			PerPage: 10,
		}

		// Set filters based on what's available
		if vr.DealID != nil && *vr.DealID != "" {
			listReq.DealID = *vr.DealID
		}
		if vr.AccountID != nil && *vr.AccountID != "" {
			listReq.AccountID = *vr.AccountID
		}

		existingTasks, _, err := s.taskRepo.List(listReq)

		// Check if task with same title already exists
		taskExists := false
		if err == nil {
			for _, t := range existingTasks {
				if t.Title == taskDef.title && task.NormalizeStatus(t.Status) != "completed" {
					taskExists = true
					break
				}
			}
		}

		// Skip if task already exists
		if taskExists {
			continue
		}

		// Create new task
		newTask := &task.Task{
			Title:       taskDef.title,
			Description: taskDef.description,
			Type:        taskDef.taskType,
			Status:      "pending",
			Priority:    taskDef.priority,
			DueDate:     &dueDate,
			CreatedBy:   vr.SalesRepID,
		}

		// Set assigned_to to sales rep
		newTask.AssignedTo = &vr.SalesRepID

		// Link to deal/account/contact/lead if available
		if vr.DealID != nil {
			newTask.DealID = vr.DealID
		}
		if vr.AccountID != nil {
			newTask.AccountID = vr.AccountID
		}
		if vr.ContactID != nil {
			newTask.ContactID = vr.ContactID
		}

		// Create task
		_ = s.taskRepo.Create(newTask) // Ignore error for now, log in production
	}
}

// notifyManager creates a notification for the manager to approve visit report
func (s *Service) notifyManager(vr *visit_report.VisitReport) {
	if s.notificationRepo == nil {
		return
	}

	// Get sales rep's manager (this would require manager relationship in user table)
	// For now, we'll create a general notification

	// Find manager user (placeholder logic - adjust based on actual org structure)
	// In production, query user.manager_id or use role-based notification

	// Create notification
	// NOTE: This requires notification domain to be properly imported
	// Placeholder for actual implementation
}

// Helper function to get string value from pointer
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
