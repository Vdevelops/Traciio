package activity

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/activity"
	domainevents "github.com/gilabs/crm-healthcare/api/internal/domain/events"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"github.com/gilabs/crm-healthcare/api/pkg/cache"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var (
	ErrActivityNotFound = errors.New("activity not found")
)

type Service struct {
	activityRepo     interfaces.ActivityRepository
	activityTypeRepo interfaces.ActivityTypeRepository
	accountRepo      interfaces.AccountRepository
	contactRepo      interfaces.ContactRepository
	userRepo         interfaces.UserRepository
	db               *gorm.DB
	eventHelper      *domainevents.Helper
	cacheService     *cache.ActivityCacheService
}

func NewService(activityRepo interfaces.ActivityRepository, activityTypeRepo interfaces.ActivityTypeRepository, accountRepo interfaces.AccountRepository, contactRepo interfaces.ContactRepository, userRepo interfaces.UserRepository, db *gorm.DB, eventHelper *domainevents.Helper) *Service {
	return &Service{
		activityRepo:     activityRepo,
		activityTypeRepo: activityTypeRepo,
		accountRepo:      accountRepo,
		contactRepo:      contactRepo,
		userRepo:         userRepo,
		db:               db,
		eventHelper:      eventHelper,
		cacheService:     cache.NewActivityCacheService(nil),
	}
}

func parseActivityTimestamp(value string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, value); err == nil {
		return t, nil
	}
	if t, err := time.Parse("2006-01-02T15:04:05Z07:00", value); err == nil {
		return t, nil
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.Local
	}
	for _, layout := range []string{"2006-01-02T15:04:05", "2006-01-02T15:04", "2006-01-02 15:04"} {
		if t, parseErr := time.ParseInLocation(layout, value, loc); parseErr == nil {
			return t, nil
		}
	}

	return time.Time{}, errors.New("invalid timestamp format")
}

type cachedActivityListResult struct {
	Activities []activity.ActivityResponse
	Pagination *PaginationResult
}

// PaginationResult represents pagination information
type PaginationResult struct {
	Page       int
	PerPage    int
	Total      int
	TotalPages int
}

// List returns a list of activities with pagination
func (s *Service) List(req *activity.ListActivitiesRequest) ([]activity.ActivityResponse, *PaginationResult, error) {
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
	if req.Type != "" {
		filterMap["type"] = req.Type
	}
	if req.AccountID != "" {
		filterMap["account_id"] = req.AccountID
	}
	if req.ContactID != "" {
		filterMap["contact_id"] = req.ContactID
	}
	if req.DealID != "" {
		filterMap["deal_id"] = req.DealID
	}
	if req.LeadID != "" {
		filterMap["lead_id"] = req.LeadID
	}
	if req.UserID != "" {
		filterMap["user_id"] = req.UserID
	}
	if req.StartDate != "" {
		filterMap["start_date"] = req.StartDate
	}
	if req.EndDate != "" {
		filterMap["end_date"] = req.EndDate
	}

	// Try cache first
	var cachedResult cachedActivityListResult
	if found, _ := s.cacheService.GetList(page, perPage, filterMap, &cachedResult); found && cachedResult.Pagination != nil {
		return cachedResult.Activities, cachedResult.Pagination, nil
	}

	activities, total, err := s.activityRepo.List(req)
	if err != nil {
		return nil, nil, err
	}

	responses := make([]activity.ActivityResponse, len(activities))
	for i, a := range activities {
		response := *a.ToActivityResponse()
		// Parse metadata JSON
		if a.Metadata != nil {
			var metadata interface{}
			if err := json.Unmarshal(a.Metadata, &metadata); err == nil {
				response.Metadata = metadata
			}
		}
		// Load Account
		if a.AccountID != nil && *a.AccountID != "" {
			if account, err := s.accountRepo.FindByID(*a.AccountID); err == nil {
				response.Account = map[string]interface{}{
					"id":   account.ID,
					"name": account.Name,
				}
			}
		}
		// Load Contact
		if a.ContactID != nil && *a.ContactID != "" {
			if contact, err := s.contactRepo.FindByID(*a.ContactID); err == nil {
				response.Contact = map[string]interface{}{
					"id":   contact.ID,
					"name": contact.Name,
				}
			}
		}
		// Load User
		if user, err := s.userRepo.FindByID(a.UserID); err == nil {
			response.User = map[string]interface{}{
				"id":   user.ID,
				"name": user.Name,
			}
		}
		// Load ActivityType
		if a.ActivityTypeID != nil && *a.ActivityTypeID != "" {
			if activityType, err := s.activityTypeRepo.FindByID(*a.ActivityTypeID); err == nil {
				response.ActivityType = activityType.ToActivityTypeResponse()
			}
		}
		responses[i] = response
	}

	pagination := &PaginationResult{
		Page:       page,
		PerPage:    perPage,
		Total:      int(total),
		TotalPages: int((total + int64(perPage) - 1) / int64(perPage)),
	}

	_ = s.cacheService.SetList(page, perPage, filterMap, cachedActivityListResult{
		Activities: responses,
		Pagination: pagination,
	})

	return responses, pagination, nil
}

// GetByID returns an activity by ID
func (s *Service) GetByID(id string) (*activity.ActivityResponse, error) {
	// Try cache first
	var cachedResponse activity.ActivityResponse
	if found, _ := s.cacheService.GetDetail(id, &cachedResponse); found && cachedResponse.ID != "" {
		return &cachedResponse, nil
	}

	a, err := s.activityRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrActivityNotFound
		}
		return nil, err
	}

	response := *a.ToActivityResponse()
	// Parse metadata JSON
	if a.Metadata != nil {
		var metadata interface{}
		if err := json.Unmarshal(a.Metadata, &metadata); err == nil {
			response.Metadata = metadata
		}
	}
	// Load ActivityType
	if a.ActivityTypeID != nil && *a.ActivityTypeID != "" {
		if activityType, err := s.activityTypeRepo.FindByID(*a.ActivityTypeID); err == nil {
			response.ActivityType = activityType.ToActivityTypeResponse()
		}
	}

	_ = s.cacheService.SetDetail(id, &response)

	return &response, nil
}

// Create creates a new activity
func (s *Service) Create(req *activity.CreateActivityRequest) (*activity.ActivityResponse, error) {
	// Validate UserID
	if req.UserID == "" {
		return nil, errors.New("user_id is required")
	}

	// Validate that either activity_type_id or type is provided
	if (req.ActivityTypeID == nil || *req.ActivityTypeID == "") && req.Type == "" {
		return nil, errors.New("either activity_type_id or type is required")
	}

	// Business rule validation: Activity must be linked to at least one entity (Lead, Account, or Deal)
	hasLeadID := req.LeadID != nil && *req.LeadID != ""
	hasAccountID := req.AccountID != nil && *req.AccountID != ""
	hasDealID := req.DealID != nil && *req.DealID != ""

	if !hasLeadID && !hasAccountID && !hasDealID {
		return nil, errors.New("activity must be linked to lead, account, or deal")
	}

	// Parse timestamp
	timestamp, err := parseActivityTimestamp(req.Timestamp)
	if err != nil {
		return nil, errors.New("invalid timestamp format")
	}

	// Marshal metadata to JSON
	var metadataJSON datatypes.JSON
	if req.Metadata != nil {
		metadataBytes, err := json.Marshal(req.Metadata)
		if err != nil {
			return nil, err
		}
		metadataJSON = metadataBytes
	}

	// Determine type: use ActivityTypeID if provided, otherwise use Type
	var activityType string
	if req.ActivityTypeID != nil && *req.ActivityTypeID != "" {
		// Load ActivityType to get the code/type
		activityTypeEntity, err := s.activityTypeRepo.FindByID(*req.ActivityTypeID)
		if err == nil && activityTypeEntity != nil {
			activityType = activityTypeEntity.Code
		} else {
			// Fallback to provided Type if ActivityType not found
			activityType = req.Type
		}
	} else {
		// Use provided Type (backward compatibility)
		activityType = req.Type
	}

	a := &activity.Activity{
		Type:           activityType,
		ActivityTypeID: req.ActivityTypeID,
		AccountID:      req.AccountID,
		ContactID:      req.ContactID,
		DealID:         req.DealID,
		LeadID:         req.LeadID,
		UserID:         req.UserID,
		Description:    req.Description,
		Timestamp:      timestamp,
		Metadata:       metadataJSON,
	}

	if err := s.activityRepo.Create(a); err != nil {
		return nil, err
	}

	// Reload
	createdActivity, err := s.activityRepo.FindByID(a.ID)
	if err != nil {
		return nil, err
	}

	// Emit activity logged event
	if s.eventHelper != nil {
		accountID := ""
		if createdActivity.AccountID != nil {
			accountID = *createdActivity.AccountID
		}
		contactID := ""
		if createdActivity.ContactID != nil {
			contactID = *createdActivity.ContactID
		}
		dealID := ""
		if createdActivity.DealID != nil {
			dealID = *createdActivity.DealID
		}
		leadID := ""
		if createdActivity.LeadID != nil {
			leadID = *createdActivity.LeadID
		}

		s.eventHelper.EmitActivityLogged(&domainevents.ActivityLoggedEvent{
			ActivityID:  createdActivity.ID,
			Type:        createdActivity.Type,
			Description: createdActivity.Description,
			UserID:      createdActivity.UserID,
			AccountID:   accountID,
			ContactID:   contactID,
			DealID:      dealID,
			LeadID:      leadID,
			Timestamp:   createdActivity.Timestamp,
			CreatedBy:   req.UserID,
			CreatedAt:   createdActivity.CreatedAt,
		}, req.UserID)
	}

	response := *createdActivity.ToActivityResponse()
	if createdActivity.Metadata != nil {
		var metadata interface{}
		if err := json.Unmarshal(createdActivity.Metadata, &metadata); err == nil {
			response.Metadata = metadata
		}
	}

	_ = s.cacheService.InvalidateOnWrite(createdActivity.ID)
	_ = s.cacheService.SetDetail(createdActivity.ID, &response)

	return &response, nil
}

// Update updates an existing activity owned by the authenticated user.
func (s *Service) Update(id string, req *activity.UpdateActivityRequest, userID string) (*activity.ActivityResponse, error) {
	a, err := s.activityRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrActivityNotFound
		}
		return nil, err
	}

	if a.UserID != userID {
		return nil, errors.New("forbidden")
	}

	if req.ActivityTypeID != nil && *req.ActivityTypeID != "" {
		a.ActivityTypeID = req.ActivityTypeID
		if activityTypeEntity, typeErr := s.activityTypeRepo.FindByID(*req.ActivityTypeID); typeErr == nil && activityTypeEntity != nil {
			a.Type = activityTypeEntity.Code
		}
	}

	if req.Description != "" {
		a.Description = req.Description
	}

	if req.Timestamp != "" {
		timestamp, parseErr := parseActivityTimestamp(req.Timestamp)
		if parseErr != nil {
			return nil, errors.New("invalid timestamp format")
		}
		a.Timestamp = timestamp
	}

	if req.Metadata != nil {
		metadataBytes, marshalErr := json.Marshal(req.Metadata)
		if marshalErr != nil {
			return nil, marshalErr
		}
		a.Metadata = datatypes.JSON(metadataBytes)
	}

	if err := s.activityRepo.Update(a); err != nil {
		return nil, err
	}

	_ = s.cacheService.InvalidateOnWrite(a.ID)

	return s.GetByID(a.ID)
}

// GetTimeline returns activity timeline
func (s *Service) GetTimeline(req *activity.ActivityTimelineRequest) ([]activity.ActivityResponse, error) {
	activities, err := s.activityRepo.GetTimeline(req)
	if err != nil {
		return nil, err
	}

	// Batch load deals to avoid N+1 queries
	dealIDs := make(map[string]bool)
	for _, a := range activities {
		if a.DealID != nil && *a.DealID != "" {
			dealIDs[*a.DealID] = true
		}
	}

	dealsMap := make(map[string]interface{})
	if len(dealIDs) > 0 && s.db != nil {
		dealIDList := make([]string, 0, len(dealIDs))
		for id := range dealIDs {
			dealIDList = append(dealIDList, id)
		}

		var deals []struct {
			ID    string `gorm:"column:id"`
			Title string `gorm:"column:title"`
		}
		if err := s.db.Table("deals").
			Select("id, title").
			Where("id IN ? AND deleted_at IS NULL", dealIDList).
			Scan(&deals).Error; err == nil {
			for _, deal := range deals {
				dealsMap[deal.ID] = map[string]interface{}{
					"id":    deal.ID,
					"title": deal.Title,
				}
			}
		}
	}

	responses := make([]activity.ActivityResponse, len(activities))
	for i, a := range activities {
		response := *a.ToActivityResponse()
		// Parse metadata JSON
		if a.Metadata != nil {
			var metadata interface{}
			if err := json.Unmarshal(a.Metadata, &metadata); err == nil {
				response.Metadata = metadata
			}
		}
		// Load Account
		if a.AccountID != nil && *a.AccountID != "" {
			if account, err := s.accountRepo.FindByID(*a.AccountID); err == nil {
				response.Account = map[string]interface{}{
					"id":   account.ID,
					"name": account.Name,
				}
			}
		}
		// Load Contact
		if a.ContactID != nil && *a.ContactID != "" {
			if contact, err := s.contactRepo.FindByID(*a.ContactID); err == nil {
				response.Contact = map[string]interface{}{
					"id":   contact.ID,
					"name": contact.Name,
				}
			}
		}
		// Load User
		if user, err := s.userRepo.FindByID(a.UserID); err == nil {
			response.User = map[string]interface{}{
				"id":   user.ID,
				"name": user.Name,
			}
		}
		// Load ActivityType
		if a.ActivityTypeID != nil && *a.ActivityTypeID != "" {
			if activityType, err := s.activityTypeRepo.FindByID(*a.ActivityTypeID); err == nil {
				response.ActivityType = activityType.ToActivityTypeResponse()
			}
		}
		// Load Deal from batch-loaded map
		if a.DealID != nil && *a.DealID != "" {
			if deal, ok := dealsMap[*a.DealID]; ok {
				response.Deal = deal
			}
		}
		responses[i] = response
	}

	return responses, nil
}
