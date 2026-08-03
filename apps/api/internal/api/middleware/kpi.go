package middleware

import (
	brickdomain "github.com/gilabs/crm-healthcare/api/internal/domain/brick"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gin-gonic/gin"
)

func KPIRepScopeMiddleware(brickRepo interfaces.BrickRepository) gin.HandlerFunc {
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
			if brickRepo != nil {
				bricks, _, err := brickRepo.List(&brickdomain.ListBricksRequest{Page: 1, PerPage: 100, ManagerID: &userCtx.UserID})
				if err == nil {
					for _, brick := range bricks {
						members, memberErr := brickRepo.GetSalesByBrickID(brick.ID)
						if memberErr != nil {
							continue
						}
						for _, member := range members {
							if member.ID == requestedUserID && member.Role != nil && member.Role.Code == "sales" {
								c.Next()
								return
							}
						}
					}
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
		if userCtx.RoleCode != "sales_manager" {
			errors.ErrorResponse(c, "FORBIDDEN", map[string]interface{}{
				"message": "Only sales managers can access manager KPI scope",
			}, nil)
			c.Abort()
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
			if err != nil {
				errors.ErrorResponse(c, "FORBIDDEN", map[string]interface{}{
					"message": "Unable to validate brick KPI scope",
				}, nil)
				c.Abort()
				return
			}
			allowed := false
			for _, brick := range bricks {
				if brick.ID == brickID {
					allowed = true
					break
				}
			}
			if !allowed {
				members, memberErr := brickRepo.GetSalesByBrickID(brickID)
				if memberErr == nil {
					teamMemberSet := make(map[string]struct{}, len(userCtx.TeamMemberIDs))
					for _, teamMemberID := range userCtx.TeamMemberIDs {
						teamMemberSet[teamMemberID] = struct{}{}
					}
					for _, member := range members {
						if _, exists := teamMemberSet[member.ID]; exists {
							allowed = true
							break
						}
					}
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

		c.Next()
	}
}
