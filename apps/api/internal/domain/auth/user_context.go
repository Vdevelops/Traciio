package auth

import (
	roledomain "github.com/gilabs/crm-healthcare/api/internal/domain/role"
)

// UserContext carries the authenticated user's identity, permissions, and data scope.
// It is resolved once per request by the ScopeMiddleware and stored in gin.Context.
// Services and repositories use this to dynamically filter data.
type UserContext struct {
	UserID        string                          `json:"user_id"`
	Email         string                          `json:"email"`
	RoleID        string                          `json:"role_id"`
	RoleCode      string                          `json:"role_code"`
	GroupID       string                          `json:"group_id"`
	TeamMemberIDs []string                        `json:"team_member_ids,omitempty"` // Pre-resolved team member user IDs for team scope
	Permissions   map[string]bool                 `json:"permissions"`
	Scopes        map[string]roledomain.ScopeType `json:"scopes"` // resource -> scope_type
}

// HasPermission checks if the user has a specific permission
func (uc *UserContext) HasPermission(permission string) bool {
	if uc.Permissions == nil {
		return false
	}
	return uc.Permissions[permission]
}

// GetScope returns the data scope for a given resource.
// Falls back to ScopeOwn if no explicit scope is defined (principle of least privilege).
func (uc *UserContext) GetScope(resource string) roledomain.ScopeType {
	if uc.Scopes == nil {
		return roledomain.ScopeOwn
	}
	scope, ok := uc.Scopes[resource]
	if !ok {
		return roledomain.ScopeOwn
	}
	return scope
}

// IsGlobalScope checks if the user has global scope for a resource
func (uc *UserContext) IsGlobalScope(resource string) bool {
	return uc.GetScope(resource) == roledomain.ScopeGlobal
}

// IsTeamScope checks if the user has team scope for a resource
func (uc *UserContext) IsTeamScope(resource string) bool {
	return uc.GetScope(resource) == roledomain.ScopeTeam
}

// IsOwnScope checks if the user has own scope for a resource
func (uc *UserContext) IsOwnScope(resource string) bool {
	return uc.GetScope(resource) == roledomain.ScopeOwn
}

// GetScopedUserIDs returns the list of user IDs that should be used for data filtering.
// Returns nil for global scope (no filtering needed).
// Returns [userID] for own scope.
// Returns team member IDs for team scope.
func (uc *UserContext) GetScopedUserIDs(resource string) []string {
	scope := uc.GetScope(resource)
	switch scope {
	case roledomain.ScopeGlobal:
		return nil // No filtering
	case roledomain.ScopeTeam:
		if len(uc.TeamMemberIDs) > 0 {
			return uc.TeamMemberIDs
		}
		// Fall back to own if no team members resolved
		return []string{uc.UserID}
	case roledomain.ScopeOwn:
		return []string{uc.UserID}
	default:
		return []string{uc.UserID}
	}
}
