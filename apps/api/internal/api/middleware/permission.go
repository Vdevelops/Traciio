package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gilabs/crm-healthcare/api/internal/service/permission"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
)

// PermissionMiddleware creates a gin middleware for permission checking
func PermissionMiddleware(permService *permission.Service, requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get user role from context (set by AuthMiddleware)
		roleInterface, exists := c.Get("user_role")
		if !exists {
			errors.UnauthorizedResponse(c, "Role not found in context")
			c.Abort()
			return
		}
		roleCode, ok := roleInterface.(string)
		if !ok {
			errors.InternalServerErrorResponse(c, "Invalid role format in context")
			c.Abort()
			return
		}

		// 2. Super Admin bypass (optional, dependent on requirements, but safe to add)
		if roleCode == "super_admin" {
			c.Next()
			return
		}

		// 3. Get effective permissions for the role (from Cache/DB)
		perms, err := permService.GetPermissionsByRole(roleCode)
		if err != nil {
			// Failed to load permissions means we can't verify access
			errors.InternalServerErrorResponse(c, "Failed to load permissions")
			c.Abort()
			return
		}

		// 4. Check if user has the required permission
		hasPermission := false
		for _, p := range perms {
			if p == requiredPermission {
				hasPermission = true
				break
			}
			// Optional: Support wildcard like "user.*"
			if p == "*" {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			errors.ErrorResponse(c, "FORBIDDEN", map[string]interface{}{
				"required": requiredPermission,
				"message":  "You do not have permission to access this resource",
			}, nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
