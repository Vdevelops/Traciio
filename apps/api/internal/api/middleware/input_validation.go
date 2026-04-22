package middleware

import (
	"bytes"
	"io"
	"strings"

	"github.com/gilabs/crm-healthcare/api/pkg/sanitizer"
	"github.com/gin-gonic/gin"
)

// InputSanitizationMiddleware sanitizes all input data
func InputSanitizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Sanitize query parameters
		for key, values := range c.Request.URL.Query() {
			for i, value := range values {
				// Check for SQL injection patterns
				if sanitizer.DetectSQLInjection(value) {
					c.JSON(400, gin.H{
						"success": false,
						"error": gin.H{
							"code":    "INVALID_INPUT",
							"message": "Invalid input detected in query parameter: " + key,
						},
					})
					c.Abort()
					return
				}
				
				// Sanitize for XSS
				if !sanitizer.ValidateNoScriptInjection(value) {
					c.JSON(400, gin.H{
						"success": false,
						"error": gin.H{
							"code":    "INVALID_INPUT",
							"message": "Potentially dangerous content detected in query parameter: " + key,
						},
					})
					c.Abort()
					return
				}
				
				values[i] = sanitizer.SanitizeInput(value)
			}
		}

		// Sanitize path parameters
		for _, param := range c.Params {
			if sanitizer.DetectSQLInjection(param.Value) {
				c.JSON(400, gin.H{
					"success": false,
					"error": gin.H{
						"code":    "INVALID_INPUT",
						"message": "Invalid input detected in path parameter: " + param.Key,
					},
				})
				c.Abort()
				return
			}
			
			if !sanitizer.ValidateNoScriptInjection(param.Value) {
				c.JSON(400, gin.H{
					"success": false,
					"error": gin.H{
						"code":    "INVALID_INPUT",
						"message": "Potentially dangerous content detected in path parameter: " + param.Key,
					},
				})
				c.Abort()
				return
			}
		}

		// For JSON requests, validate content type and body
		if c.Request.Method != "GET" && c.Request.Method != "DELETE" {
			contentType := c.GetHeader("Content-Type")
			if strings.Contains(contentType, "application/json") {
				// AI chat endpoints carry unrestricted conversational text and conversation
				// history from previous AI responses. Body-level pattern scanning would produce
				// false positives on natural language (semicolons, dashes, etc.).
				// These endpoints are protected by parameterized queries (GORM) at the
				// data layer, so body-level SQL injection scanning is safely skipped here.
				isAIEndpoint := strings.HasPrefix(c.Request.URL.Path, "/api/v1/ai/")
				if isAIEndpoint {
					c.Next()
					return
				}

				// Read body
				bodyBytes, err := io.ReadAll(c.Request.Body)
				if err != nil {
					c.JSON(400, gin.H{
						"success": false,
						"error": gin.H{
							"code":    "INVALID_REQUEST",
							"message": "Failed to read request body",
						},
					})
					c.Abort()
					return
				}

				// Restore body for handlers
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				// Validate body doesn't contain obvious script injection
				bodyString := string(bodyBytes)
				if !sanitizer.ValidateNoScriptInjection(bodyString) {
					c.JSON(400, gin.H{
						"success": false,
						"error": gin.H{
							"code":    "INVALID_INPUT",
							"message": "Potentially dangerous content detected in request body",
						},
					})
					c.Abort()
					return
				}

				// Check for SQL injection patterns in body
				if sanitizer.DetectSQLInjection(bodyString) {
					c.JSON(400, gin.H{
						"success": false,
						"error": gin.H{
							"code":    "INVALID_INPUT",
							"message": "Invalid input pattern detected in request body",
						},
					})
					c.Abort()
					return
				}
			}
		}

		c.Next()
	}
}

// FileUploadValidationMiddleware validates file uploads for security
func FileUploadValidationMiddleware() gin.HandlerFunc {
	// Allowed file extensions
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".pdf":  true,
		".doc":  true,
		".docx": true,
		".xls":  true,
		".xlsx": true,
		".csv":  true,
		".txt":  true,
	}

	// Maximum file size (10MB)
	maxFileSize := int64(10 * 1024 * 1024)

	return func(c *gin.Context) {
		// Only process multipart forms
		if !strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data") {
			c.Next()
			return
		}

		// Parse multipart form
		if err := c.Request.ParseMultipartForm(maxFileSize); err != nil {
			c.JSON(400, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "FILE_TOO_LARGE",
					"message": "File size exceeds maximum allowed size (10MB)",
				},
			})
			c.Abort()
			return
		}

		// Validate each file
		if c.Request.MultipartForm != nil && c.Request.MultipartForm.File != nil {
			for _, files := range c.Request.MultipartForm.File {
				for _, file := range files {
					// Check file size
					if file.Size > maxFileSize {
						c.JSON(400, gin.H{
							"success": false,
							"error": gin.H{
								"code":    "FILE_TOO_LARGE",
								"message": "File size exceeds maximum allowed size (10MB)",
							},
						})
						c.Abort()
						return
					}

					// Validate filename
					sanitizedFilename := sanitizer.SanitizeFilename(file.Filename)
					if sanitizedFilename != file.Filename {
						c.JSON(400, gin.H{
							"success": false,
							"error": gin.H{
								"code":    "INVALID_FILENAME",
								"message": "Filename contains invalid characters",
							},
						})
						c.Abort()
						return
					}

					// Check file extension
					ext := strings.ToLower(sanitizedFilename[strings.LastIndex(sanitizedFilename, "."):])
					if !allowedExtensions[ext] {
						c.JSON(400, gin.H{
							"success": false,
							"error": gin.H{
								"code":    "INVALID_FILE_TYPE",
								"message": "File type not allowed",
							},
						})
						c.Abort()
						return
					}

					// Read file header to validate MIME type
					fileContent, err := file.Open()
					if err != nil {
						c.JSON(500, gin.H{
							"success": false,
							"error": gin.H{
								"code":    "FILE_READ_ERROR",
								"message": "Failed to read file",
							},
						})
						c.Abort()
						return
					}
					defer fileContent.Close()

					// Read first 512 bytes to detect content type
					buffer := make([]byte, 512)
					_, err = fileContent.Read(buffer)
					if err != nil && err != io.EOF {
						c.JSON(500, gin.H{
							"success": false,
							"error": gin.H{
								"code":    "FILE_READ_ERROR",
								"message": "Failed to validate file content",
							},
						})
						c.Abort()
						return
					}

					// Basic validation: ensure it's not an executable
					if bytes.Contains(buffer, []byte("MZ")) || // PE executable
						bytes.Contains(buffer, []byte("\x7fELF")) || // ELF executable
						bytes.Contains(buffer, []byte("#!/")) { // Shell script
						c.JSON(400, gin.H{
							"success": false,
							"error": gin.H{
								"code":    "INVALID_FILE_TYPE",
								"message": "Executable files are not allowed",
							},
						})
						c.Abort()
						return
					}
				}
			}
		}

		c.Next()
	}
}
