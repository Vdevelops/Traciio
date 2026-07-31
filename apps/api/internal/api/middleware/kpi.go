package middleware

import (
	brickdomain "github.com/gilabs/crm-healthcare/api/internal/domain/brick"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gin-gonic/gin"
)

func KPIRepScopeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userCtx := GetUserContext(c)
		if userCtx == nil {
			c.Next()
			return
		}
		if userCtx.RoleCode == "admin" || userCtx.RoleCode == "super_admin" {
			c.Next()
			return
		}

		requestedUserID := c.Query("userId")
		if requestedUserID == "" || requestedUserID == userCtx.UserID {
			c.Next()
			return
		}

		if userCtx.RoleCode == "sales_manager" {
			for _, teamMemberID := range userCtx.TeamMemberIDs {
				if teamMemberID == requestedUserID {
					c.Next()
					return
				}
			}
		}

		errors.ErrorResponse(c, "FORBIDDEN", map[string]interface{}{
			"message": "You do not have permission to access this KPI scope",
		}, nil)
		c.Abort()
	}
}

func KPIManagerScopeMiddleware(brickRepo interfaces.BrickRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userCtx := GetUserContext(c)
		if userCtx == nil {
			c.Next()
			return
		}
		if userCtx.RoleCode == "admin" || userCtx.RoleCode == "super_admin" {
			c.Next()
			return
		}

		requestedManagerID := c.Query("managerId")
		if requestedManagerID != "" && requestedManagerID != userCtx.UserID {
			errors.ErrorResponse(c, "FORBIDDEN", map[string]interface{}{
				"message": "Manager KPI can only be accessed for the authenticated manager",
			}, nil)
			c.Abort()
			return
		}

		brickID := c.Query("brickId")
		if brickID != "" && brickRepo != nil {
			bricks, _, err := brickRepo.List(&brickdomain.ListBricksRequest{Page: 1, PerPage: 100, ManagerID: &userCtx.UserID})
			if err == nil {
				allowed := false
				for _, brick := range bricks {
					if brick.ID == brickID {
						allowed = true
						break
					}
				}
				if !allowed {
					errors.ErrorResponse(c, "FORBIDDEN", map[string]interface{}{
						"message": "You do not have permission to access this brick KPI scope",
					}, nil)
					c.Abort()
					return
				}
			}
		}

		c.Next()
	}
}