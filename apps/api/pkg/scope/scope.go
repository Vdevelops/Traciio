package scope

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/auth"
	roledomain "github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"gorm.io/gorm"
)

// OwnerColumn maps resources to their ownership column in the database.
// This determines which column is used for "own" scope filtering.
var OwnerColumn = map[string]string{
	"leads":         "assigned_to",
	"deals":         "assigned_to",
	"tasks":         "assigned_to",
	"schedules":     "user_id",
	"visit-reports": "sales_rep_id",
}

// TeamColumn maps resources to their team/group column in the database.
// Used for "team" scope filtering based on the user's group.
var TeamColumn = map[string]string{
	"leads":         "assigned_to",
	"deals":         "assigned_to",
	"tasks":         "assigned_to",
	"schedules":     "user_id",
	"visit-reports": "sales_rep_id",
}

// ApplyScope applies dynamic data scoping to a GORM query based on the user's role scope.
// This replaces hardcoded role checks (e.g., if role == "admin") with policy-driven filtering.
//
// Scope behavior:
//   - global: No filter applied, user sees all data
//   - team:   Filters by group members (users in the same group_id)
//   - own:    Filters by the user's own ID on the ownership column
func ApplyScope(db *gorm.DB, ctx *auth.UserContext, resource string) *gorm.DB {
	if ctx == nil {
		return db.Where("1 = 0") // Deny all if no context
	}

	scopeType := ctx.GetScope(resource)

	switch scopeType {
	case roledomain.ScopeGlobal:
		return db
	case roledomain.ScopeTeam:
		if ctx.GroupID == "" {
			// No group assigned, fall back to own scope
			return applyOwnScope(db, ctx, resource)
		}
		return applyTeamScope(db, ctx, resource)
	case roledomain.ScopeOwn:
		return applyOwnScope(db, ctx, resource)
	default:
		// Default to own scope (principle of least privilege)
		return applyOwnScope(db, ctx, resource)
	}
}

// applyOwnScope filters data to only the user's own records
func applyOwnScope(db *gorm.DB, ctx *auth.UserContext, resource string) *gorm.DB {
	column, ok := OwnerColumn[resource]
	if !ok {
		column = "created_by" // Default fallback
	}
	return db.Where(column+" = ?", ctx.UserID)
}

// applyTeamScope filters data to records belonging to users in the same group.
// Uses a subquery to find all user IDs within the same group.
func applyTeamScope(db *gorm.DB, ctx *auth.UserContext, resource string) *gorm.DB {
	column, ok := TeamColumn[resource]
	if !ok {
		column = "created_by"
	}
	// Subquery: SELECT id FROM users WHERE group_id = ?
	return db.Where(column+" IN (SELECT id FROM users WHERE group_id = ?)", ctx.GroupID)
}
