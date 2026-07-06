package handlers

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/internal/domain/account"
	"github.com/gilabs/crm-healthcare/api/internal/domain/activity"
	"github.com/gilabs/crm-healthcare/api/internal/domain/contact"
	"github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
	"github.com/gilabs/crm-healthcare/api/internal/domain/product"
	"github.com/gilabs/crm-healthcare/api/internal/domain/visit_report"
	accountservice "github.com/gilabs/crm-healthcare/api/internal/service/account"
	activityservice "github.com/gilabs/crm-healthcare/api/internal/service/activity"
	contactservice "github.com/gilabs/crm-healthcare/api/internal/service/contact"
	pipelineservice "github.com/gilabs/crm-healthcare/api/internal/service/pipeline"
	productservice "github.com/gilabs/crm-healthcare/api/internal/service/product"
	visitreportservice "github.com/gilabs/crm-healthcare/api/internal/service/visit_report"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type DealHandler struct {
	dealService        *pipelineservice.Service
	visitReportService *visitreportservice.Service
	activityService    *activityservice.Service
	accountService     *accountservice.Service
	contactService     *contactservice.Service
	productService     *productservice.Service
}

func NewDealHandler(
	dealService *pipelineservice.Service,
	visitReportService *visitreportservice.Service,
	activityService *activityservice.Service,
	accountService *accountservice.Service,
	contactService *contactservice.Service,
	productService *productservice.Service,
) *DealHandler {
	return &DealHandler{
		dealService:        dealService,
		visitReportService: visitReportService,
		activityService:    activityService,
		accountService:     accountService,
		contactService:     contactService,
		productService:     productService,
	}
}

func ensureDealInScope(c *gin.Context, deal *pipeline.DealResponse) bool {
	userCtx := middleware.GetUserContext(c)
	if userCtx == nil {
		return true
	}

	scoped := userCtx.GetScopedUserIDs("deals")
	if scoped == nil {
		return true
	}

	for _, scopedUserID := range scoped {
		if scopedUserID == deal.AssignedTo {
			return true
		}
	}

	errors.ErrorResponse(c, "FORBIDDEN", map[string]interface{}{
		"message": "You do not have permission to access this deal",
	}, nil)
	return false
}

// List handles list deals request
func (h *DealHandler) List(c *gin.Context) {
	var req pipeline.ListDealsRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	// Apply RBAC scope filtering
	if userCtx := middleware.GetUserContext(c); userCtx != nil {
		req.ScopedUserIDs = userCtx.GetScopedUserIDs("deals")
	}

	deals, pagination, err := h.dealService.ListDeals(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{
		Pagination: &response.PaginationMeta{
			Page:       pagination.Page,
			PerPage:    pagination.PerPage,
			Total:      pagination.Total,
			TotalPages: pagination.TotalPages,
			HasNext:    pagination.Page < pagination.TotalPages,
			HasPrev:    pagination.Page > 1,
		},
		Filters: map[string]interface{}{},
	}

	if req.Search != "" {
		meta.Filters["search"] = req.Search
	}
	if req.StageID != "" {
		meta.Filters["stage_id"] = req.StageID
	}
	if req.AccountID != "" {
		meta.Filters["account_id"] = req.AccountID
	}
	if req.AssignedTo != "" {
		meta.Filters["assigned_to"] = req.AssignedTo
	}
	if req.Status != "" {
		meta.Filters["status"] = req.Status
	}
	if req.Source != "" {
		meta.Filters["source"] = req.Source
	}

	response.SuccessResponse(c, deals, meta)
}

// ListByStage handles list deals grouped by pipeline stage.
//
// @Summary List deals by stage
// @Description Returns deals grouped by stage_id
// @Tags Deals
// @Accept json
// @Produce json
// @Param stage_id query string false "Filter by stage id"
// @Param account_id query string false "Filter by account id"
// @Param assigned_to query string false "Filter by assigned user id"
// @Param search query string false "Search term"
// @Param min_value query number false "Minimum deal value"
// @Param max_value query number false "Maximum deal value"
// @Param date_from query string false "Start date (YYYY-MM-DD)"
// @Param date_to query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} response.SuccessResponse
// @Router /api/v1/deals/by-stage [get]
func (h *DealHandler) ListByStage(c *gin.Context) {
	var req pipeline.ListDealsRequest

	// We reuse ListDeals query shape. Pagination isn't meaningful for grouped result,
	// but keeping it allows existing filters and validation.
	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	// Apply RBAC scope filtering
	if userCtx := middleware.GetUserContext(c); userCtx != nil {
		req.ScopedUserIDs = userCtx.GetScopedUserIDs("deals")
	}

	// Make sure we get a large enough set to group; default per_page in service is 20.
	// For Kanban board, it's safer to request up to 100.
	if req.PerPage < 1 {
		req.PerPage = 100
	}
	if req.PerPage > 100 {
		req.PerPage = 100
	}
	if req.Page < 1 {
		req.Page = 1
	}

	deals, _, err := h.dealService.ListDeals(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	grouped := make(map[string][]pipeline.DealResponse)
	for _, d := range deals {
		grouped[d.StageID] = append(grouped[d.StageID], d)
	}

	response.SuccessResponse(c, grouped, &response.Meta{Filters: map[string]interface{}{}})
}

// GetByID handles get deal by ID request
func (h *DealHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	deal, err := h.dealService.GetDealByID(id)
	if err != nil {
		if err == pipelineservice.ErrDealNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "deal",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	if userCtx := middleware.GetUserContext(c); userCtx != nil {
		if !ensureDealInScope(c, deal) {
			return
		}
	}

	response.SuccessResponse(c, deal, nil)
}

// Create handles create deal request
func (h *DealHandler) Create(c *gin.Context) {
	var req pipeline.CreateDealRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	// Get user ID from context
	userID := ""
	if userIDVal, exists := c.Get("user_id"); exists {
		if id, ok := userIDVal.(string); ok {
			userID = id
		}
	}

	createdDeal, err := h.dealService.CreateDeal(&req, userID)
	if err != nil {
		if err == pipelineservice.ErrDealNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource": "deal",
			}, nil)
			return
		}
		if err == pipelineservice.ErrAccountNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "account",
				"resource_id": req.AccountID,
			}, nil)
			return
		}
		if err == pipelineservice.ErrInvalidStage {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "pipeline_stage",
				"resource_id": req.StageID,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if userID != "" {
		meta.CreatedBy = userID
	}

	response.SuccessResponseCreated(c, createdDeal, meta)
}

// Update handles update deal request
func (h *DealHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req pipeline.UpdateDealRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	userID := ""
	if userIDVal, exists := c.Get("user_id"); exists {
		if parsedUserID, ok := userIDVal.(string); ok {
			userID = parsedUserID
		}
	}

	existingDeal, err := h.dealService.GetDealByID(id)
	if err != nil {
		if err == pipelineservice.ErrDealNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "deal",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}
	if !ensureDealInScope(c, existingDeal) {
		return
	}

	updatedDeal, err := h.dealService.UpdateDeal(id, &req, userID)
	if err != nil {
		if err == pipelineservice.ErrDealNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "deal",
				"resource_id": id,
			}, nil)
			return
		}
		if err == pipelineservice.ErrAccountNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "account",
				"resource_id": req.AccountID,
			}, nil)
			return
		}
		if err == pipelineservice.ErrInvalidStage {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "pipeline_stage",
				"resource_id": req.StageID,
			}, nil)
			return
		}
		if err == pipelineservice.ErrCloseReasonRequired {
			errors.ErrorResponse(c, "CLOSE_REASON_REQUIRED", map[string]interface{}{
				"field": "close_reason",
			}, []response.FieldError{
				{
					Field:   "close_reason",
					Code:    "REQUIRED",
					Message: "Close reason is required when status changes to won or lost",
				},
			})
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			meta.UpdatedBy = id
		}
	}

	response.SuccessResponse(c, updatedDeal, meta)
}

// Move handles move deal request
func (h *DealHandler) Move(c *gin.Context) {
	if userCtx := middleware.GetUserContext(c); userCtx != nil {
		if !userCtx.HasPermission("pipeline.update_stage") {
			errors.ForbiddenResponse(c, "pipeline.update_stage", nil)
			return
		}
	}

	id := c.Param("id")
	var req pipeline.MoveDealRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		errors.UnauthorizedResponse(c, "")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		errors.UnauthorizedResponse(c, "")
		return
	}

	existingDeal, err := h.dealService.GetDealByID(id)
	if err != nil {
		if err == pipelineservice.ErrDealNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "deal",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}
	if !ensureDealInScope(c, existingDeal) {
		return
	}

	movedDeal, err := h.dealService.MoveStageWithValidation(id, req.StageID, userIDStr, req.Reason)
	if err != nil {
		if err == pipelineservice.ErrDealNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "deal",
				"resource_id": id,
			}, nil)
			return
		}
		if err == pipelineservice.ErrPipelineStageNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "pipeline_stage",
				"resource_id": req.StageID,
			}, nil)
			return
		}
		if err == pipelineservice.ErrStageRequirementsNotMet {
			errors.ErrorResponse(c, "STAGE_REQUIREMENTS_NOT_MET", map[string]interface{}{
				"message": "Stage transition requirements not met. Check products, deal value, and stage order.",
			}, nil)
			return
		}
		if err == pipelineservice.ErrCloseReasonRequired {
			errors.ErrorResponse(c, "CLOSE_REASON_REQUIRED", map[string]interface{}{
				"field": "reason",
			}, []response.FieldError{
				{
					Field:   "reason",
					Code:    "REQUIRED",
					Message: "Reason is required when moving a deal to won or lost",
				},
			})
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			meta.UpdatedBy = id
		}
	}

	response.SuccessResponse(c, movedDeal, meta)
}

// Delete handles delete deal request
func (h *DealHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	existingDeal, err := h.dealService.GetDealByID(id)
	if err != nil {
		if err == pipelineservice.ErrDealNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "deal",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}
	if !ensureDealInScope(c, existingDeal) {
		return
	}

	err = h.dealService.DeleteDeal(id)
	if err != nil {
		if err == pipelineservice.ErrDealNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "deal",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	// Get user ID for meta
	meta := &response.Meta{}
	if userIDVal, exists := c.Get("user_id"); exists {
		if id, ok := userIDVal.(string); ok {
			meta.DeletedBy = id
		}
	}

	response.SuccessResponseDeleted(c, "deal", id, meta)
}

// GetVisitReportsByDeal handles get visit reports by deal ID request
func (h *DealHandler) GetVisitReportsByDeal(c *gin.Context) {
	dealID := c.Param("id")

	// Verify deal exists
	_, err := h.dealService.GetDealByID(dealID)
	deal, err := h.dealService.GetDealByID(dealID)
	if err != nil {
		if err == pipelineservice.ErrDealNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "deal",
				"resource_id": dealID,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}
	if !ensureDealInScope(c, deal) {
		return
	}

	var req visit_report.ListVisitReportsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	// Set deal_id filter
	req.DealID = dealID

	visitReports, pagination, err := h.visitReportService.List(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{
		Pagination: &response.PaginationMeta{
			Page:       pagination.Page,
			PerPage:    pagination.PerPage,
			Total:      pagination.Total,
			TotalPages: pagination.TotalPages,
			HasNext:    pagination.Page < pagination.TotalPages,
			HasPrev:    pagination.Page > 1,
		},
		Filters: map[string]interface{}{
			"deal_id": dealID,
		},
	}

	if req.Status != "" {
		meta.Filters["status"] = req.Status
	}
	if req.StartDate != "" {
		meta.Filters["start_date"] = req.StartDate
	}
	if req.EndDate != "" {
		meta.Filters["end_date"] = req.EndDate
	}

	response.SuccessResponse(c, visitReports, meta)
}

// GetActivitiesByDeal handles get activities by deal ID request
func (h *DealHandler) GetActivitiesByDeal(c *gin.Context) {
	dealID := c.Param("id")

	// Verify deal exists
	deal, err := h.dealService.GetDealByID(dealID)
	if err != nil {
		if err == pipelineservice.ErrDealNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "deal",
				"resource_id": dealID,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}
	if !ensureDealInScope(c, deal) {
		return
	}

	var req activity.ListActivitiesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	// Set deal_id filter
	req.DealID = dealID

	activities, pagination, err := h.activityService.List(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{
		Pagination: &response.PaginationMeta{
			Page:       pagination.Page,
			PerPage:    pagination.PerPage,
			Total:      pagination.Total,
			TotalPages: pagination.TotalPages,
			HasNext:    pagination.Page < pagination.TotalPages,
			HasPrev:    pagination.Page > 1,
		},
		Filters: map[string]interface{}{
			"deal_id": dealID,
		},
	}

	if req.Type != "" {
		meta.Filters["type"] = req.Type
	}
	if req.StartDate != "" {
		meta.Filters["start_date"] = req.StartDate
	}
	if req.EndDate != "" {
		meta.Filters["end_date"] = req.EndDate
	}

	response.SuccessResponse(c, activities, meta)
}

// GetFormData handles get data for deal form
func (h *DealHandler) GetFormData(c *gin.Context) {
	// Verify user is authenticated
	_, exists := c.Get("user_id")
	if !exists {
		errors.UnauthorizedResponse(c, "")
		return
	}

	// 1. Get all active accounts
	accountReq := &account.ListAccountsRequest{
		Status:  "active",
		Page:    1,
		PerPage: 1000,
	}
	accounts, _, err := h.accountService.List(accountReq)
	if err != nil {
		accounts = []account.AccountResponse{}
	}

	// 2. Get all contacts
	contactReq := &contact.ListContactsRequest{
		Page:    1,
		PerPage: 1000,
	}
	contacts, _, err := h.contactService.List(contactReq)
	if err != nil {
		contacts = []contact.ContactResponse{}
	}

	// 3. Get all pipeline stages
	stages, err := h.dealService.ListStages(&pipeline.ListPipelineStagesRequest{})
	if err != nil {
		stages = []pipeline.PipelineStageResponse{}
	}

	// 4. Get active products
	productReq := &product.ListProductsRequest{
		Status:  "active",
		Page:    1,
		PerPage: 1000,
	}
	products, _, err := h.productService.ListProducts(productReq)
	if err != nil {
		products = []product.ProductResponse{}
	}

	// Build response
	formData := map[string]interface{}{
		"accounts":        accounts,
		"contacts":        contacts,
		"pipeline_stages": stages,
		"products":        products,
	}

	response.SuccessResponse(c, formData, nil)
}
