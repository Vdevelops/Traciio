package handlers

import (
	"fmt"
	"log"
	"strings"

	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/internal/domain/account"
	"github.com/gilabs/crm-healthcare/api/internal/domain/contact"
	"github.com/gilabs/crm-healthcare/api/internal/domain/lead"
	"github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
	"github.com/gilabs/crm-healthcare/api/internal/domain/visit_report"
	accountservice "github.com/gilabs/crm-healthcare/api/internal/service/account"
	contactservice "github.com/gilabs/crm-healthcare/api/internal/service/contact"
	fileservice "github.com/gilabs/crm-healthcare/api/internal/service/file"
	leadservice "github.com/gilabs/crm-healthcare/api/internal/service/lead"
	pipelineservice "github.com/gilabs/crm-healthcare/api/internal/service/pipeline"
	visitreportservice "github.com/gilabs/crm-healthcare/api/internal/service/visit_report"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type VisitReportHandler struct {
	visitReportService *visitreportservice.Service
	fileService        *fileservice.Service
	accountService     *accountservice.Service
	contactService     *contactservice.Service
	pipelineService    *pipelineservice.Service
	leadService        *leadservice.Service
}

func NewVisitReportHandler(
	visitReportService *visitreportservice.Service,
	fileService *fileservice.Service,
	accountService *accountservice.Service,
	contactService *contactservice.Service,
	pipelineService *pipelineservice.Service,
	leadService *leadservice.Service,
) *VisitReportHandler {
	return &VisitReportHandler{
		visitReportService: visitReportService,
		fileService:        fileService,
		accountService:     accountService,
		contactService:     contactService,
		pipelineService:    pipelineService,
		leadService:        leadService,
	}
}

// List handles list visit reports request
func (h *VisitReportHandler) List(c *gin.Context) {
	var req visit_report.ListVisitReportsRequest

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
		req.ScopedUserIDs = userCtx.GetScopedUserIDs("visit-reports")
	}

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
		Filters: map[string]interface{}{},
	}

	if req.Search != "" {
		meta.Filters["search"] = req.Search
	}
	if req.Status != "" {
		meta.Filters["status"] = req.Status
	}
	if req.AccountID != "" {
		meta.Filters["account_id"] = req.AccountID
	}
	if req.SalesRepID != "" {
		meta.Filters["sales_rep_id"] = req.SalesRepID
	}

	response.SuccessResponse(c, visitReports, meta)
}

// GetByID handles get visit report by ID request
func (h *VisitReportHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	visitReport, err := h.visitReportService.GetByID(id)
	if err != nil {
		if err == visitreportservice.ErrVisitReportNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "visit_report",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, visitReport, nil)
}

// Create handles create visit report request
func (h *VisitReportHandler) Create(c *gin.Context) {
	var req visit_report.CreateVisitReportRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	// Get user ID from context (set by auth middleware)
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

	// Set SalesRepID from authenticated user
	req.SalesRepID = userIDStr

	createdVisitReport, err := h.visitReportService.Create(&req)
	if err != nil {
		if err == visitreportservice.ErrAccountNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource": "account",
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			meta.CreatedBy = id
		}
	}

	response.SuccessResponseCreated(c, createdVisitReport, meta)
}

// Update handles update visit report request
func (h *VisitReportHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req visit_report.UpdateVisitReportRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	updatedVisitReport, err := h.visitReportService.Update(id, &req)
	if err != nil {
		if err == visitreportservice.ErrVisitReportNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "visit_report",
				"resource_id": id,
			}, nil)
			return
		}
		if err == visitreportservice.ErrInvalidStatus {
			errors.ErrorResponse(c, "INVALID_STATUS", map[string]interface{}{
				"message": "Cannot update visit report with current status",
			}, nil)
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

	response.SuccessResponse(c, updatedVisitReport, meta)
}

// Delete handles delete visit report request
func (h *VisitReportHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.visitReportService.Delete(id)
	if err != nil {
		if err == visitreportservice.ErrVisitReportNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "visit_report",
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

	response.SuccessResponseDeleted(c, "visit_report", id, meta)
}

// CheckIn handles check-in request
// Supports both JSON (with photo_url) and multipart/form-data (with photo file upload)
func (h *VisitReportHandler) CheckIn(c *gin.Context) {
	id := c.Param("id")
	var req visit_report.CheckInRequest

	// Check if request is multipart (file upload)
	contentType := c.GetHeader("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		// Handle multipart form data with photo upload
		// Parse location from form
		locationLat := c.PostForm("location[latitude]")
		locationLon := c.PostForm("location[longitude]")
		locationAddr := c.PostForm("location[address]")

		if locationLat == "" || locationLon == "" {
			errors.ErrorResponse(c, "INVALID_REQUEST", map[string]interface{}{
				"message": "Location latitude and longitude are required",
			}, nil)
			return
		}

		// Parse location
		var lat, lon float64
		if _, err := fmt.Sscanf(locationLat, "%f", &lat); err != nil {
			errors.ErrorResponse(c, "INVALID_REQUEST", map[string]interface{}{
				"message": "Invalid latitude format",
			}, nil)
			return
		}
		if _, err := fmt.Sscanf(locationLon, "%f", &lon); err != nil {
			errors.ErrorResponse(c, "INVALID_REQUEST", map[string]interface{}{
				"message": "Invalid longitude format",
			}, nil)
			return
		}

		req.Location = &visit_report.Location{
			Latitude:  lat,
			Longitude: lon,
			Address:   locationAddr,
		}

		// Parse device GPS metadata if provided
		if deviceLat := c.PostForm("device_gps[latitude]"); deviceLat != "" {
			var deviceGPS visit_report.GPSMetadata
			if _, err := fmt.Sscanf(deviceLat, "%f", &deviceGPS.Latitude); err == nil {
				if deviceLon := c.PostForm("device_gps[longitude]"); deviceLon != "" {
					if _, err := fmt.Sscanf(deviceLon, "%f", &deviceGPS.Longitude); err == nil {
						if acc := c.PostForm("device_gps[accuracy]"); acc != "" {
							fmt.Sscanf(acc, "%f", &deviceGPS.Accuracy)
						}
						if ts := c.PostForm("device_gps[timestamp]"); ts != "" {
							fmt.Sscanf(ts, "%d", &deviceGPS.Timestamp)
						}
						req.DeviceGPS = &deviceGPS
					}
				}
			}
		}

		// Parse photo GPS metadata if provided
		if photoLat := c.PostForm("photo_gps[latitude]"); photoLat != "" {
			var photoGPS visit_report.GPSMetadata
			if _, err := fmt.Sscanf(photoLat, "%f", &photoGPS.Latitude); err == nil {
				if photoLon := c.PostForm("photo_gps[longitude]"); photoLon != "" {
					if _, err := fmt.Sscanf(photoLon, "%f", &photoGPS.Longitude); err == nil {
						if ts := c.PostForm("photo_gps[timestamp]"); ts != "" {
							fmt.Sscanf(ts, "%d", &photoGPS.Timestamp)
						}
						req.PhotoGPS = &photoGPS
					}
				}
			}
		}

		// Handle photo file upload - REQUIRED for mobile check-in
		file, err := c.FormFile("photo")
		if err != nil {
			// Try alternative field names
			file, err = c.FormFile("file")
			if err != nil {
				file, err = c.FormFile("image")
			}
		}

		// Selfie picture is required for mobile check-in (multipart/form-data)
		if err != nil || file == nil {
			errors.ErrorResponse(c, "INVALID_REQUEST", map[string]interface{}{
				"message": "Selfie picture is required for check-in. Please provide a photo file.",
			}, nil)
			return
		}

		// Upload and compress image
		uploadedURL, err := h.fileService.UploadImage(file)
		if err != nil {
			log.Printf("Error uploading image for check-in: %v", err)
			errors.ErrorResponse(c, "UPLOAD_FAILED", map[string]interface{}{
				"message": "Failed to upload selfie picture. Please try again.",
			}, nil)
			return
		}
		req.PhotoURL = &uploadedURL
	} else {
		// Handle JSON request (for web/admin - photo is optional)
		if err := c.ShouldBindJSON(&req); err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				errors.HandleValidationError(c, validationErrors)
				return
			}
			errors.InvalidRequestBodyResponse(c)
			return
		}
		// For JSON requests, photo_url is optional (web/admin can check-in without photo)
		// Mobile should always use multipart/form-data with photo file
	}

	// Get user ID from context
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

	visitReport, err := h.visitReportService.CheckIn(id, &req, userIDStr)
	if err != nil {
		log.Printf("Error in CheckIn service for visit report %s: %v", id, err)

		if err == visitreportservice.ErrVisitReportNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "visit_report",
				"resource_id": id,
			}, nil)
			return
		}
		if err == visitreportservice.ErrInvalidGPS || strings.Contains(err.Error(), "GPS") || strings.Contains(err.Error(), "accuracy") || strings.Contains(err.Error(), "timestamp") || strings.Contains(err.Error(), "location") {
			errors.ErrorResponse(c, "INVALID_GPS", map[string]interface{}{
				"message": err.Error(),
			}, nil)
			return
		}
		// Check for "already checked in" error
		if err.Error() == "already checked in" {
			errors.ErrorResponse(c, "INVALID_OPERATION", map[string]interface{}{
				"message": "Visit report has already been checked in",
			}, nil)
			return
		}
		// Log unexpected errors and return internal server error
		// This ensures we don't expose internal errors to the client
		log.Printf("Unexpected error in CheckIn: %v", err)
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, visitReport, nil)
}

// GetByIDMobile handles get visit report by ID request for mobile (validates ownership)
func (h *VisitReportHandler) GetByIDMobile(c *gin.Context) {
	id := c.Param("id")

	// Get user ID from context
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

	// Log for debugging
	log.Printf("GetByIDMobile: User %s requesting visit report %s", userIDStr, id)

	visitReport, err := h.visitReportService.GetByID(id)
	if err != nil {
		if err == visitreportservice.ErrVisitReportNotFound {
			log.Printf("GetByIDMobile: Visit report %s not found for user %s", id, userIDStr)
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "visit_report",
				"resource_id": id,
				"message":     "Visit report not found. Please make sure the visit report exists and belongs to you.",
			}, nil)
			return
		}
		log.Printf("GetByIDMobile: Error getting visit report %s for user %s: %v", id, userIDStr, err)
		errors.InternalServerErrorResponse(c, "")
		return
	}

	// Validate ownership - check if visit report belongs to the logged-in user
	if visitReport.SalesRep == nil {
		log.Printf("GetByIDMobile: Visit report %s has no SalesRep data", id)
		errors.ErrorResponse(c, "FORBIDDEN", map[string]interface{}{
			"message": "You don't have permission to view this visit report",
		}, nil)
		return
	}

	// Type assert SalesRep to map[string]interface{}
	salesRepMap, ok := visitReport.SalesRep.(map[string]interface{})
	if !ok {
		log.Printf("GetByIDMobile: Visit report %s has invalid SalesRep data type", id)
		errors.ErrorResponse(c, "FORBIDDEN", map[string]interface{}{
			"message": "You don't have permission to view this visit report",
		}, nil)
		return
	}

	salesRepID, ok := salesRepMap["id"].(string)
	if !ok || salesRepID != userIDStr {
		log.Printf("GetByIDMobile: User %s attempted to access visit report %s owned by %s", userIDStr, id, salesRepID)
		errors.ErrorResponse(c, "FORBIDDEN", map[string]interface{}{
			"message": "You can only view your own visit reports",
		}, nil)
		return
	}

	log.Printf("GetByIDMobile: Successfully retrieved visit report %s for user %s", id, userIDStr)
	response.SuccessResponse(c, visitReport, nil)
}

// CheckInMobile handles check-in request for mobile (validates ownership)
// Processes multipart form data directly to avoid reading request body twice
func (h *VisitReportHandler) CheckInMobile(c *gin.Context) {
	id := c.Param("id")
	var req visit_report.CheckInRequest

	// Get user ID from context first
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

	// First, verify that the visit report exists and belongs to the user
	// Use GetByIDMobile which validates ownership
	vr, err := h.visitReportService.GetByID(id)
	if err != nil {
		if err == visitreportservice.ErrVisitReportNotFound {
			log.Printf("Visit report %s not found for check-in by user %s", id, userIDStr)
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "visit_report",
				"resource_id": id,
				"message":     "Visit report not found. Please make sure you're accessing your own visit report.",
			}, nil)
			return
		}
		log.Printf("Error getting visit report %s for check-in: %v", id, err)
		errors.InternalServerErrorResponse(c, "")
		return
	}

	// Validate ownership - check if visit report belongs to the logged-in user
	if vr.SalesRep == nil {
		errors.ErrorResponse(c, "FORBIDDEN", map[string]interface{}{
			"message": "You don't have permission to check-in this visit report",
		}, nil)
		return
	}

	// Type assert SalesRep to map[string]interface{}
	salesRepMap, ok := vr.SalesRep.(map[string]interface{})
	if !ok {
		errors.ErrorResponse(c, "FORBIDDEN", map[string]interface{}{
			"message": "You don't have permission to check-in this visit report",
		}, nil)
		return
	}

	salesRepID, ok := salesRepMap["id"].(string)
	if !ok || salesRepID != userIDStr {
		log.Printf("User %s attempted to check-in visit report %s owned by %s", userIDStr, id, salesRepID)
		errors.ErrorResponse(c, "FORBIDDEN", map[string]interface{}{
			"message": "You can only check-in your own visit reports",
		}, nil)
		return
	}

	// Ownership validated, now process multipart form data
	// Check if request is multipart (file upload) - REQUIRED for mobile
	contentType := c.GetHeader("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		errors.ErrorResponse(c, "INVALID_REQUEST", map[string]interface{}{
			"message": "Mobile check-in requires multipart/form-data with selfie picture",
		}, nil)
		return
	}

	// Handle multipart form data with photo upload
	// Parse location from form
	locationLat := c.PostForm("location[latitude]")
	locationLon := c.PostForm("location[longitude]")
	locationAddr := c.PostForm("location[address]")

	if locationLat == "" || locationLon == "" {
		errors.ErrorResponse(c, "INVALID_REQUEST", map[string]interface{}{
			"message": "Location latitude and longitude are required",
		}, nil)
		return
	}

	// Parse location
	var lat, lon float64
	if _, err := fmt.Sscanf(locationLat, "%f", &lat); err != nil {
		errors.ErrorResponse(c, "INVALID_REQUEST", map[string]interface{}{
			"message": "Invalid latitude format",
		}, nil)
		return
	}
	if _, err := fmt.Sscanf(locationLon, "%f", &lon); err != nil {
		errors.ErrorResponse(c, "INVALID_REQUEST", map[string]interface{}{
			"message": "Invalid longitude format",
		}, nil)
		return
	}

	req.Location = &visit_report.Location{
		Latitude:  lat,
		Longitude: lon,
		Address:   locationAddr,
	}

	// Parse device GPS metadata if provided
	if deviceLat := c.PostForm("device_gps[latitude]"); deviceLat != "" {
		var deviceGPS visit_report.GPSMetadata
		if _, err := fmt.Sscanf(deviceLat, "%f", &deviceGPS.Latitude); err == nil {
			if deviceLon := c.PostForm("device_gps[longitude]"); deviceLon != "" {
				if _, err := fmt.Sscanf(deviceLon, "%f", &deviceGPS.Longitude); err == nil {
					if acc := c.PostForm("device_gps[accuracy]"); acc != "" {
						fmt.Sscanf(acc, "%f", &deviceGPS.Accuracy)
					}
					if ts := c.PostForm("device_gps[timestamp]"); ts != "" {
						fmt.Sscanf(ts, "%d", &deviceGPS.Timestamp)
					}
					req.DeviceGPS = &deviceGPS
				}
			}
		}
	}

	// Parse photo GPS metadata if provided
	if photoLat := c.PostForm("photo_gps[latitude]"); photoLat != "" {
		var photoGPS visit_report.GPSMetadata
		if _, err := fmt.Sscanf(photoLat, "%f", &photoGPS.Latitude); err == nil {
			if photoLon := c.PostForm("photo_gps[longitude]"); photoLon != "" {
				if _, err := fmt.Sscanf(photoLon, "%f", &photoGPS.Longitude); err == nil {
					if ts := c.PostForm("photo_gps[timestamp]"); ts != "" {
						fmt.Sscanf(ts, "%d", &photoGPS.Timestamp)
					}
					req.PhotoGPS = &photoGPS
				}
			}
		}
	}

	// Handle photo file upload - REQUIRED for mobile check-in
	file, err := c.FormFile("photo")
	if err != nil {
		// Try alternative field names
		file, err = c.FormFile("file")
		if err != nil {
			file, err = c.FormFile("image")
		}
	}

	// Selfie picture is required for mobile check-in (multipart/form-data)
	if err != nil || file == nil {
		log.Printf("No photo file provided for check-in by user %s", userIDStr)
		errors.ErrorResponse(c, "INVALID_REQUEST", map[string]interface{}{
			"message": "Selfie picture is required for check-in. Please provide a photo file.",
		}, nil)
		return
	}

	// Upload and compress image
	uploadedURL, err := h.fileService.UploadImage(file)
	if err != nil {
		log.Printf("Error uploading image for check-in: %v", err)
		errors.ErrorResponse(c, "UPLOAD_FAILED", map[string]interface{}{
			"message": "Failed to upload selfie picture. Please try again.",
		}, nil)
		return
	}
	req.PhotoURL = &uploadedURL

	// Call service to perform check-in
	visitReport, err := h.visitReportService.CheckIn(id, &req, userIDStr)
	if err != nil {
		log.Printf("Error in CheckIn service for visit report %s: %v", id, err)

		if err == visitreportservice.ErrVisitReportNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "visit_report",
				"resource_id": id,
			}, nil)
			return
		}
		if err == visitreportservice.ErrInvalidGPS || strings.Contains(err.Error(), "GPS") || strings.Contains(err.Error(), "accuracy") || strings.Contains(err.Error(), "timestamp") || strings.Contains(err.Error(), "location") {
			errors.ErrorResponse(c, "INVALID_GPS", map[string]interface{}{
				"message": err.Error(),
			}, nil)
			return
		}
		// Check for "already checked in" error
		if err.Error() == "already checked in" {
			errors.ErrorResponse(c, "INVALID_OPERATION", map[string]interface{}{
				"message": "Visit report has already been checked in",
			}, nil)
			return
		}
		// Log unexpected errors and return internal server error
		log.Printf("Unexpected error in CheckInMobile: %v", err)
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, visitReport, nil)
}

// CheckOutMobile handles check-out request for mobile (validates ownership)
// Processes JSON request directly to avoid reading request body twice
func (h *VisitReportHandler) CheckOutMobile(c *gin.Context) {
	id := c.Param("id")
	var req visit_report.CheckOutRequest

	// Get user ID from context first
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

	// First, verify that the visit report exists and belongs to the user
	vr, err := h.visitReportService.GetByID(id)
	if err != nil {
		if err == visitreportservice.ErrVisitReportNotFound {
			log.Printf("Visit report %s not found for check-out by user %s", id, userIDStr)
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "visit_report",
				"resource_id": id,
				"message":     "Visit report not found. Please make sure you're accessing your own visit report.",
			}, nil)
			return
		}
		log.Printf("Error getting visit report %s for check-out: %v", id, err)
		errors.InternalServerErrorResponse(c, "")
		return
	}

	// Validate ownership - check if visit report belongs to the logged-in user
	if vr.SalesRep == nil {
		errors.ErrorResponse(c, "FORBIDDEN", map[string]interface{}{
			"message": "You don't have permission to check-out this visit report",
		}, nil)
		return
	}

	// Type assert SalesRep to map[string]interface{}
	salesRepMap, ok := vr.SalesRep.(map[string]interface{})
	if !ok {
		errors.ErrorResponse(c, "FORBIDDEN", map[string]interface{}{
			"message": "You don't have permission to check-out this visit report",
		}, nil)
		return
	}

	salesRepID, ok := salesRepMap["id"].(string)
	if !ok || salesRepID != userIDStr {
		log.Printf("User %s attempted to check-out visit report %s owned by %s", userIDStr, id, salesRepID)
		errors.ErrorResponse(c, "FORBIDDEN", map[string]interface{}{
			"message": "You can only check-out your own visit reports",
		}, nil)
		return
	}

	// Ownership validated, now process JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	// Call service to perform check-out
	visitReport, err := h.visitReportService.CheckOut(id, &req, userIDStr)
	if err != nil {
		if err == visitreportservice.ErrVisitReportNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "visit_report",
				"resource_id": id,
			}, nil)
			return
		}
		errors.ErrorResponse(c, "INVALID_OPERATION", map[string]interface{}{
			"message": err.Error(),
		}, nil)
		return
	}

	response.SuccessResponse(c, visitReport, nil)
}

// CheckOut handles check-out request
func (h *VisitReportHandler) CheckOut(c *gin.Context) {
	id := c.Param("id")
	var req visit_report.CheckOutRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	// Get user ID from context
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

	visitReport, err := h.visitReportService.CheckOut(id, &req, userIDStr)
	if err != nil {
		if err == visitreportservice.ErrVisitReportNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "visit_report",
				"resource_id": id,
			}, nil)
			return
		}
		errors.ErrorResponse(c, "INVALID_OPERATION", map[string]interface{}{
			"message": err.Error(),
		}, nil)
		return
	}

	response.SuccessResponse(c, visitReport, nil)
}

// Submit handles submit visit report request
func (h *VisitReportHandler) Submit(c *gin.Context) {
	id := c.Param("id")
	var req visit_report.SubmitRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	// Get user ID from context
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

	visitReport, err := h.visitReportService.Submit(id, &req, userIDStr)
	if err != nil {
		if err == visitreportservice.ErrVisitReportNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "visit_report",
				"resource_id": id,
			}, nil)
			return
		}
		if err == visitreportservice.ErrInvalidStatus {
			errors.ErrorResponse(c, "INVALID_STATUS", map[string]interface{}{
				"message": "Cannot submit visit report with current status",
			}, nil)
			return
		}
		if err == visitreportservice.ErrNotOwner {
			errors.ErrorResponse(c, "FORBIDDEN", map[string]interface{}{
				"message": "You can only submit your own visit reports",
			}, nil)
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

	response.SuccessResponse(c, visitReport, meta)
}

// Approve handles approve visit report request
func (h *VisitReportHandler) Approve(c *gin.Context) {
	id := c.Param("id")

	// Get user ID from context
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

	visitReport, err := h.visitReportService.Approve(id, userIDStr)
	if err != nil {
		if err == visitreportservice.ErrVisitReportNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "visit_report",
				"resource_id": id,
			}, nil)
			return
		}
		if err == visitreportservice.ErrInvalidStatus {
			errors.ErrorResponse(c, "INVALID_STATUS", map[string]interface{}{
				"message": "Cannot approve visit report with current status",
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, visitReport, nil)
}

// Reject handles reject visit report request
func (h *VisitReportHandler) Reject(c *gin.Context) {
	id := c.Param("id")
	var req visit_report.RejectRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	// Get user ID from context
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

	visitReport, err := h.visitReportService.Reject(id, &req, userIDStr)
	if err != nil {
		if err == visitreportservice.ErrVisitReportNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "visit_report",
				"resource_id": id,
			}, nil)
			return
		}
		if err == visitreportservice.ErrInvalidStatus {
			errors.ErrorResponse(c, "INVALID_STATUS", map[string]interface{}{
				"message": "Cannot reject visit report with current status",
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, visitReport, nil)
}

// UploadPhoto handles photo upload request
// Supports both multipart/form-data (file upload) and JSON (photo_url)
func (h *VisitReportHandler) UploadPhoto(c *gin.Context) {
	id := c.Param("id")
	var photoURL string

	// Check if request is multipart (file upload)
	contentType := c.GetHeader("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		// Handle multipart file upload
		file, err := c.FormFile("photo")
		if err != nil {
			// Try alternative field names
			file, err = c.FormFile("file")
			if err != nil {
				file, err = c.FormFile("image")
				if err != nil {
					errors.ErrorResponse(c, "INVALID_REQUEST", map[string]interface{}{
						"message": "No file provided. Use 'photo', 'file', or 'image' field name",
					}, nil)
					return
				}
			}
		}

		// Upload and compress image
		uploadedURL, err := h.fileService.UploadImage(file)
		if err != nil {
			errors.ErrorResponse(c, "UPLOAD_FAILED", map[string]interface{}{
				"message": err.Error(),
			}, nil)
			return
		}

		photoURL = uploadedURL
	} else {
		// Handle JSON request with photo_url
		var req visit_report.UploadPhotoRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				errors.HandleValidationError(c, validationErrors)
				return
			}
			errors.InvalidRequestBodyResponse(c)
			return
		}
		photoURL = req.PhotoURL
	}

	// Create request for service
	req := visit_report.UploadPhotoRequest{
		PhotoURL: photoURL,
	}

	visitReport, err := h.visitReportService.UploadPhoto(id, &req)
	if err != nil {
		if err == visitreportservice.ErrVisitReportNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "visit_report",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, visitReport, nil)
}

// GetMyVisitReports handles get visit reports for logged-in user request (mobile endpoint)
func (h *VisitReportHandler) GetMyVisitReports(c *gin.Context) {
	var req visit_report.ListVisitReportsRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidQueryParamResponse(c)
		return
	}

	// Get user ID from context
	userID := ""
	if userIDVal, exists := c.Get("user_id"); exists {
		if id, ok := userIDVal.(string); ok {
			userID = id
		}
	}

	if userID == "" {
		errors.ErrorResponse(c, "UNAUTHORIZED", map[string]interface{}{
			"message": "User ID not found in context",
		}, nil)
		return
	}

	// Check if this is for route optimization (allows viewing all visit reports)
	forRouteOptimization := c.Query("for_route_optimization") == "true"

	visitReports, pagination, err := h.visitReportService.GetMyVisitReports(userID, &req, forRouteOptimization)
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

	// Add offset and has_more for offset-based pagination (infinity scroll)
	if req.Offset > 0 {
		meta.Pagination.Offset = &pagination.Offset
		meta.Pagination.HasMore = &pagination.HasMore
	}

	if req.Search != "" {
		meta.Filters["search"] = req.Search
	}
	if req.Status != "" {
		meta.Filters["status"] = req.Status
	}
	if req.AccountID != "" {
		meta.Filters["account_id"] = req.AccountID
	}
	if req.DealID != "" {
		meta.Filters["deal_id"] = req.DealID
	}
	if req.LeadID != "" {
		meta.Filters["lead_id"] = req.LeadID
	}
	if req.StartDate != "" {
		meta.Filters["start_date"] = req.StartDate
	}
	if req.EndDate != "" {
		meta.Filters["end_date"] = req.EndDate
	}

	response.SuccessResponse(c, visitReports, meta)
}

// CreateMobile handles create visit report request for mobile (validates ownership)
func (h *VisitReportHandler) CreateMobile(c *gin.Context) {
	var req visit_report.CreateVisitReportRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	// Get user ID from context (set by auth middleware)
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

	// Set SalesRepID from authenticated user (mobile users can only create for themselves)
	req.SalesRepID = userIDStr

	createdVisitReport, err := h.visitReportService.Create(&req)
	if err != nil {
		if err == visitreportservice.ErrAccountNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource": "account",
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	meta := &response.Meta{}
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			meta.CreatedBy = id
		}
	}

	response.SuccessResponseCreated(c, createdVisitReport, meta)
}

// UpdateMobile handles update visit report request for mobile (validates ownership and status)
func (h *VisitReportHandler) UpdateMobile(c *gin.Context) {
	id := c.Param("id")
	var req visit_report.UpdateVisitReportRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	// Get user ID from context
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

	// Validate ownership: only the owner can update
	vr, err := h.visitReportService.GetByID(id)
	if err != nil {
		if err == visitreportservice.ErrVisitReportNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "visit_report",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	if vr.SalesRepID != userIDStr {
		errors.ErrorResponse(c, "FORBIDDEN", map[string]interface{}{
			"message": "You can only update your own visit reports",
		}, nil)
		return
	}

	// Validate status: can only update if status is draft (mobile flow: draft -> check-in -> check-out -> submit)
	if vr.Status != "draft" {
		errors.ErrorResponse(c, "INVALID_STATUS", map[string]interface{}{
			"message": "Can only update visit report with draft status",
		}, nil)
		return
	}

	updatedVisitReport, err := h.visitReportService.Update(id, &req)
	if err != nil {
		if err == visitreportservice.ErrVisitReportNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "visit_report",
				"resource_id": id,
			}, nil)
			return
		}
		if err == visitreportservice.ErrInvalidStatus {
			errors.ErrorResponse(c, "INVALID_STATUS", map[string]interface{}{
				"message": "Cannot update visit report with current status",
			}, nil)
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

	response.SuccessResponse(c, updatedVisitReport, meta)
}

// SubmitMobile handles submit visit report request for mobile (validates ownership and status)
func (h *VisitReportHandler) SubmitMobile(c *gin.Context) {
	id := c.Param("id")
	var req visit_report.SubmitRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors.HandleValidationError(c, validationErrors)
			return
		}
		errors.InvalidRequestBodyResponse(c)
		return
	}

	// Get user ID from context
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

	visitReport, err := h.visitReportService.Submit(id, &req, userIDStr)
	if err != nil {
		if err == visitreportservice.ErrVisitReportNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "visit_report",
				"resource_id": id,
			}, nil)
			return
		}
		if err == visitreportservice.ErrInvalidStatus {
			errors.ErrorResponse(c, "INVALID_STATUS", map[string]interface{}{
				"message": "Cannot submit visit report with current status",
			}, nil)
			return
		}
		if err.Error() == "unauthorized: you can only submit your own visit reports" {
			errors.ErrorResponse(c, "FORBIDDEN", map[string]interface{}{
				"message": "You can only submit your own visit reports",
			}, nil)
			return
		}
		if err.Error() == "check-in and check-out are required before submit" {
			errors.ErrorResponse(c, "INVALID_REQUEST", map[string]interface{}{
				"message": "Check-in and check-out are required before submit",
			}, nil)
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

	response.SuccessResponse(c, visitReport, meta)
}

// GetFormData handles get form data for visit report creation (accounts, contacts, deals, leads)
func (h *VisitReportHandler) GetFormData(c *gin.Context) {
	// Verify user is authenticated (required by auth middleware, but check for safety)
	_, exists := c.Get("user_id")
	if !exists {
		errors.UnauthorizedResponse(c, "")
		return
	}

	// Get all active accounts (general data, accessible to all users - consistent with web)
	accountReq := &account.ListAccountsRequest{
		Status:  "active",
		Page:    1,
		PerPage: 1000, // Get all active accounts for selection
		// Note: AssignedTo is not set, so all active accounts are returned (general data)
	}
	accounts, _, err := h.accountService.List(accountReq)
	if err != nil {
		log.Printf("Error fetching accounts: %v", err)
		accounts = []account.AccountResponse{}
	}

	// Get all contacts (general data, will be filtered by account_id on frontend)
	contactReq := &contact.ListContactsRequest{
		Page:    1,
		PerPage: 1000, // Get all contacts for selection
	}
	contacts, _, err := h.contactService.List(contactReq)
	if err != nil {
		log.Printf("Error fetching contacts: %v", err)
		contacts = []contact.ContactResponse{}
	}

	// Get all open deals (general data, accessible to all users - consistent with web)
	dealReq := &pipeline.ListDealsRequest{
		Status:  "open", // Only open deals
		Page:    1,
		PerPage: 1000, // Get all open deals for selection
		// Note: AssignedTo is not set, so all open deals are returned (general data)
	}
	deals, _, err := h.pipelineService.ListDeals(dealReq)
	if err != nil {
		log.Printf("Error fetching deals: %v", err)
		deals = []pipeline.DealResponse{}
	}

	// Get all leads (general data, accessible to all users - consistent with web)
	leadReq := &lead.ListLeadsRequest{
		Page:    1,
		PerPage: 1000, // Get all leads for selection
		// Note: AssignedTo is not set, so all leads are returned (general data)
	}
	leads, _, err := h.leadService.List(leadReq)
	if err != nil {
		log.Printf("Error fetching leads: %v", err)
		leads = []lead.LeadResponse{}
	}

	// Build response
	formData := map[string]interface{}{
		"accounts": accounts,
		"contacts": contacts,
		"deals":    deals,
		"leads":    leads,
	}

	response.SuccessResponse(c, formData, nil)
}

// DeleteMobile handles delete visit report request for mobile (validates ownership)
func (h *VisitReportHandler) DeleteMobile(c *gin.Context) {
	id := c.Param("id")

	// Get user ID from context
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

	// Validate ownership: only the owner can delete
	vr, err := h.visitReportService.GetByID(id)
	if err != nil {
		if err == visitreportservice.ErrVisitReportNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "visit_report",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	if vr.SalesRepID != userIDStr {
		errors.ErrorResponse(c, "FORBIDDEN", map[string]interface{}{
			"message": "You can only delete your own visit reports",
		}, nil)
		return
	}

	err = h.visitReportService.Delete(id)
	if err != nil {
		if err == visitreportservice.ErrVisitReportNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "visit_report",
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

	response.SuccessResponseDeleted(c, "visit_report", id, meta)
}

// UploadPhotoMobile handles photo upload request for mobile (only supports file upload, validates ownership)
func (h *VisitReportHandler) UploadPhotoMobile(c *gin.Context) {
	id := c.Param("id")

	// Get user ID from context
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

	// Validate ownership: only the owner can upload photos
	vr, err := h.visitReportService.GetByID(id)
	if err != nil {
		if err == visitreportservice.ErrVisitReportNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "visit_report",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	if vr.SalesRepID != userIDStr {
		errors.ErrorResponse(c, "FORBIDDEN", map[string]interface{}{
			"message": "You can only upload photos to your own visit reports",
		}, nil)
		return
	}

	// Mobile endpoint only supports multipart/form-data file upload (not JSON with photo_url)
	contentType := c.GetHeader("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		errors.ErrorResponse(c, "INVALID_REQUEST", map[string]interface{}{
			"message": "Mobile photo upload requires multipart/form-data with file. Use 'photo', 'file', or 'image' field name.",
		}, nil)
		return
	}

	// Handle multipart file upload
	file, err := c.FormFile("photo")
	if err != nil {
		// Try alternative field names
		file, err = c.FormFile("file")
		if err != nil {
			file, err = c.FormFile("image")
			if err != nil {
				errors.ErrorResponse(c, "INVALID_REQUEST", map[string]interface{}{
					"message": "No file provided. Use 'photo', 'file', or 'image' field name",
				}, nil)
				return
			}
		}
	}

	// Upload and compress image
	uploadedURL, err := h.fileService.UploadImage(file)
	if err != nil {
		errors.ErrorResponse(c, "UPLOAD_FAILED", map[string]interface{}{
			"message": err.Error(),
		}, nil)
		return
	}

	// Create request for service
	req := visit_report.UploadPhotoRequest{
		PhotoURL: uploadedURL,
	}

	visitReport, err := h.visitReportService.UploadPhoto(id, &req)
	if err != nil {
		if err == visitreportservice.ErrVisitReportNotFound {
			errors.ErrorResponse(c, "NOT_FOUND", map[string]interface{}{
				"resource":    "visit_report",
				"resource_id": id,
			}, nil)
			return
		}
		errors.InternalServerErrorResponse(c, "")
		return
	}

	response.SuccessResponse(c, visitReport, nil)
}
