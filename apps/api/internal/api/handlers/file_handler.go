package handlers

import (
	"strings"

	"github.com/gilabs/crm-healthcare/api/internal/service/file"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	fileService *file.Service
}

func NewFileHandler(fileService *file.Service) *FileHandler {
	return &FileHandler{
		fileService: fileService,
	}
}

// UploadImage handles image upload request
func (h *FileHandler) UploadImage(c *gin.Context) {
	// Get file from form data
	file, err := c.FormFile("image")
	if err != nil {
		// Try alternative field names
		file, err = c.FormFile("file")
		if err != nil {
			file, err = c.FormFile("photo")
			if err != nil {
				errors.ErrorResponse(c, "INVALID_REQUEST", map[string]interface{}{
					"message": "No file provided. Use 'image', 'file', or 'photo' field name",
				}, nil)
				return
			}
		}
	}

	// Validate file type
	contentType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		errors.ErrorResponse(c, "INVALID_FILE_TYPE", map[string]interface{}{
			"message": "Only image files are allowed",
			"allowed": "image/jpeg, image/png, image/gif, image/webp",
		}, nil)
		return
	}

	// Validate file size (max 10MB)
	maxSize := int64(10 * 1024 * 1024) // 10MB
	if file.Size > maxSize {
		errors.ErrorResponse(c, "FILE_TOO_LARGE", map[string]interface{}{
			"message":  "File size exceeds maximum allowed",
			"max_size": "10MB",
			"size":     file.Size,
		}, nil)
		return
	}

	// Upload and compress image
	uploadedURL, err := h.fileService.UploadImage(file)
	if err != nil {
		errors.ErrorResponse(c, "UPLOAD_FAILED", map[string]interface{}{
			"message": "Failed to upload image",
			"error":   err.Error(),
		}, nil)
		return
	}

	response.SuccessResponse(c, gin.H{
		"url": uploadedURL,
	}, nil)
}

// DeleteFile handles file deletion request
func (h *FileHandler) DeleteFile(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		errors.ErrorResponse(c, "INVALID_REQUEST", map[string]interface{}{
			"message": "Filename is required",
		}, nil)
		return
	}

	err := h.fileService.DeleteFile(filename)
	if err != nil {
		errors.ErrorResponse(c, "DELETE_FAILED", map[string]interface{}{
			"message": "Failed to delete file",
			"error":   err.Error(),
		}, nil)
		return
	}

	response.SuccessResponse(c, gin.H{
		"message": "File deleted successfully",
	}, nil)
}
