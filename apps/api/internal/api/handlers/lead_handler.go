package handlers

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/internal/domain/activity"
	domainauth "github.com/gilabs/crm-healthcare/api/internal/domain/auth"
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead"
	lead_qual_domain "github.com/gilabs/crm-healthcare/api/internal/domain/lead_qualification"
	"github.com/gilabs/crm-healthcare/api/internal/domain/visit_report"
	activityservice "github.com/gilabs/crm-healthcare/api/internal/service/activity"
	leadservice "github.com/gilabs/crm-healthcare/api/internal/service/lead"
	leadqualificationservice "github.com/gilabs/crm-healthcare/api/internal/service/lead_qualification"
	visitreportservice "github.com/gilabs/crm-healthcare/api/internal/service/visit_report"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type LeadHandler struct {
	leadService              *leadservice.Service
	visitReportService       *visitreportservice.Service
	activityService          *activityservice.Service
	leadQualificationService *leadqualificationservice.Service
}

func NewLeadHandler(
	leadService *leadservice.Service,
	visitReportService *visitreportservice.Service,
	activityService *activityservice.Service,
	leadQualificationService *leadqualificationservice.Service,
) *LeadHandler {
	return &LeadHandler{
		leadService:              leadService,
		visitReportService:       visitReportService,
		activityService:          activityService,
		leadQualificationService: leadQualificationService,
	}
}

// List handles list leads request
func (h *LeadHandler) List(c *gin.Context) {
	var req lead.ListLeadsRequest

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
		req.ScopedUserIDs = userCtx.GetScopedUserIDs("leads")
	}

	leads, pagination, err := h.leadService.List(&req)
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

	if req.Status != "" {
		meta.Filters["status"] = req.Status
	}
	if req.Source != "" {
		meta.Filters["source"] = req.Source
	}
	if req.AssignedTo != "" {
		meta.Filters["assigned_to"] = req.AssignedTo
	}
	if req.Search != "" {
		meta.Filters["search"] = req.Search
	}
	if req.Sort != "" {
		meta.Sort = &response.SortMeta{
			Field: req.Sort,
			Order: req.Order,
		}
	}

	response.SuccessResponse(c, leads, meta)
}

// GetByID handles get lead by ID request
func (h *LeadHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	lead, err := h.leadService.GetByID(id)
	if err != nil {
		if err == leadservice.ErrLeadNotFound {
			errors.ErrorResponse(c, "LEAD_NOT_FOUND", map[string]interface{}{
				"resource":    "lead",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, lead, nil)
}

// Create handles create lead request
func (h *LeadHandler) Create(c *gin.Context) {
	var req lead.CreateLeadRequest

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
	var currentUser *domainauth.UserContext
	if userCtx := middleware.GetUserContext(c); userCtx != nil {
		currentUser = userCtx
	}
	if userIDVal, exists := c.Get("user_id"); exists {
		if id, ok := userIDVal.(string); ok {
			userID = id
		}
	}

	createdLead, err := h.leadService.Create(&req, userID, currentUser)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if userID != "" {
		meta.CreatedBy = userID
	}

	response.SuccessResponseCreated(c, createdLead, meta)
}

// Update handles update lead request
func (h *LeadHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req lead.UpdateLeadRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	updatedLead, err := h.leadService.Update(id, &req, middleware.GetUserContext(c))
	if err != nil {
		if err == leadservice.ErrLeadNotFound {
			errors.ErrorResponse(c, "LEAD_NOT_FOUND", map[string]interface{}{
				"resource":    "lead",
				"resource_id": id,
			}, nil)
			return
		}
		if err == leadservice.ErrLeadAlreadyConverted {
			errors.ErrorResponse(c, "LEAD_ALREADY_CONVERTED", map[string]interface{}{
				"lead_id": id,
			}, nil)
			return
		}
		if err == leadservice.ErrInvalidLeadStatus {
			errors.ErrorResponse(c, "INVALID_LEAD_STATUS", map[string]interface{}{
				"lead_status_id": req.LeadStatusID,
				"lead_status":    req.LeadStatus,
			}, nil)
			return
		}
		if err == leadservice.ErrLeadStatusReasonRequired {
			errors.ErrorResponse(c, "LEAD_STATUS_REASON_REQUIRED", map[string]interface{}{
				"field": "status_reason",
			}, []response.FieldError{
				{
					Field:   "status_reason",
					Code:    "REQUIRED",
					Message: "Status reason is required for won or lost lead status",
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

	response.SuccessResponse(c, updatedLead, meta)
}

// Delete handles delete lead request
func (h *LeadHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.leadService.Delete(id)
	if err != nil {
		if err == leadservice.ErrLeadNotFound {
			errors.ErrorResponse(c, "LEAD_NOT_FOUND", map[string]interface{}{
				"resource":    "lead",
				"resource_id": id,
			}, nil)
			return
		}
		if err == leadservice.ErrLeadAlreadyConverted {
			errors.ErrorResponse(c, "LEAD_ALREADY_CONVERTED", map[string]interface{}{
				"lead_id": id,
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

	response.SuccessResponseDeleted(c, "lead", id, meta)
}

// Convert handles convert lead to opportunity request.
// Requires "leads.convert" permission.
func (h *LeadHandler) Convert(c *gin.Context) {
	// Check RBAC permission
	if userCtx := middleware.GetUserContext(c); userCtx != nil {
		if !userCtx.HasPermission("leads.convert") {
			errors.ForbiddenResponse(c, "leads.convert", nil)
			return
		}
	}

	id := c.Param("id")
	var req lead.ConvertLeadRequest

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

	convertResponse, err := h.leadService.Convert(id, &req, userID)
	if err != nil {
		if err == leadservice.ErrLeadNotFound {
			errors.ErrorResponse(c, "LEAD_NOT_FOUND", map[string]interface{}{
				"resource":    "lead",
				"resource_id": id,
			}, nil)
			return
		}
		if err == leadservice.ErrLeadAlreadyConverted {
			errors.ErrorResponse(c, "LEAD_ALREADY_CONVERTED", map[string]interface{}{
				"lead_id": id,
			}, nil)
			return
		}
		if err == leadservice.ErrLeadCannotConvert {
			errors.ErrorResponse(c, "LEAD_CANNOT_CONVERT", map[string]interface{}{
				"lead_id":         id,
				"required_status": "qualified",
			}, nil)
			return
		}
		if err == leadservice.ErrStageNotFound {
			errors.ErrorResponse(c, "STAGE_NOT_FOUND", map[string]interface{}{
				"resource":    "pipeline_stage",
				"resource_id": req.StageID,
			}, nil)
			return
		}
		if err == leadservice.ErrAccountCreationFailed {
			errors.ErrorResponse(c, "ACCOUNT_CREATION_FAILED", nil, nil)
			return
		}
		if err == leadservice.ErrContactCreationFailed {
			errors.ErrorResponse(c, "CONTACT_CREATION_FAILED", nil, nil)
			return
		}
		if err == leadservice.ErrOpportunityCreationFailed {
			errors.ErrorResponse(c, "OPPORTUNITY_CREATION_FAILED", nil, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if userID != "" {
		meta.CreatedBy = userID
		meta.Additional = map[string]interface{}{
			"converted_by": userID,
		}
	}

	response.SuccessResponseCreated(c, convertResponse, meta)
}

// GetAnalytics handles get lead analytics request
func (h *LeadHandler) GetAnalytics(c *gin.Context) {
	var req lead.LeadAnalyticsRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	analytics, err := h.leadService.GetAnalytics(&req)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if req.StartDate != "" || req.EndDate != "" {
		meta.Additional = map[string]interface{}{
			"period": map[string]interface{}{
				"start_date": req.StartDate,
				"end_date":   req.EndDate,
			},
		}
	}

	response.SuccessResponse(c, analytics, meta)
}

// GetFormData handles get form data for creating a lead
func (h *LeadHandler) GetFormData(c *gin.Context) {
	formData, err := h.leadService.GetFormData()
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, formData, nil)
}

// GetMobileFormData handles get form data for creating a lead on mobile
func (h *LeadHandler) GetMobileFormData(c *gin.Context) {
	formData, err := h.leadService.GetMobileFormData()
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, formData, nil)
}

// GetVisitReportsByLead handles get visit reports by lead ID request
func (h *LeadHandler) GetVisitReportsByLead(c *gin.Context) {
	leadID := c.Param("id")

	// Verify lead exists
	_, err := h.leadService.GetByID(leadID)
	if err != nil {
		if err == leadservice.ErrLeadNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "lead",
				"resource_id": leadID,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
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

	// Set LeadID filter
	req.LeadID = leadID

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
			"lead_id": leadID,
		},
	}

	response.SuccessResponse(c, visitReports, meta)
}

// GetActivitiesByLead handles get activities by lead ID request
func (h *LeadHandler) GetActivitiesByLead(c *gin.Context) {
	leadID := c.Param("id")

	// Verify lead exists
	_, err := h.leadService.GetByID(leadID)
	if err != nil {
		if err == leadservice.ErrLeadNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "lead",
				"resource_id": leadID,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
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

	// Set LeadID filter
	req.LeadID = leadID

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
			"lead_id": leadID,
		},
	}

	response.SuccessResponse(c, activities, meta)
}

// CreateAccountFromLead handles create account from lead request (pre-convert)
func (h *LeadHandler) CreateAccountFromLead(c *gin.Context) {
	id := c.Param("id")
	var req lead.CreateAccountFromLeadRequest

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
		if uid, ok := userIDVal.(string); ok {
			userID = uid
		}
	}

	createResponse, err := h.leadService.CreateAccountFromLead(id, &req, userID)
	if err != nil {
		if err == leadservice.ErrLeadNotFound {
			errors.ErrorResponse(c, "LEAD_NOT_FOUND", map[string]interface{}{
				"resource":    "lead",
				"resource_id": id,
			}, nil)
			return
		}
		if err.Error() == "lead already has an account" {
			errors.ErrorResponse(c, "LEAD_ALREADY_HAS_ACCOUNT", map[string]interface{}{
				"lead_id": id,
			}, nil)
			return
		}
		if err.Error() == "company name is required to create account" {
			errors.ErrorResponse(c, "COMPANY_NAME_REQUIRED", map[string]interface{}{
				"lead_id": id,
			}, nil)
			return
		}
		if err == leadservice.ErrAccountCreationFailed {
			errors.ErrorResponse(c, "ACCOUNT_CREATION_FAILED", map[string]interface{}{
				"lead_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, createResponse, nil)
}

// GetQualification handles get lead qualification checklist request
func (h *LeadHandler) GetQualification(c *gin.Context) {
	id := c.Param("id")

	qualification, err := h.leadQualificationService.GetByLeadID(id)
	if err != nil {
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, qualification, nil)
}

// UpsertQualification handles create/update lead qualification checklist request
func (h *LeadHandler) UpsertQualification(c *gin.Context) {
	id := c.Param("id")
	var req lead_qual_domain.UpsertLeadQualificationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	qualification, err := h.leadQualificationService.Upsert(id, &req)
	if err != nil {
		if err == leadservice.ErrLeadNotFound || err.Error() == "record not found" {
			errors.ErrorResponse(c, "LEAD_NOT_FOUND", map[string]interface{}{
				"resource":    "lead",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, qualification, nil)
}
