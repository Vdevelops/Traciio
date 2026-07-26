package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	domainauth "github.com/gilabs/crm-healthcare/api/internal/domain/auth"
	brickdomain "github.com/gilabs/crm-healthcare/api/internal/domain/brick"
	roledomain "github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	permissionservice "github.com/gilabs/crm-healthcare/api/internal/service/permission"
	"github.com/gilabs/crm-healthcare/api/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const (
	// UserContextKey is the gin context key for the resolved UserContext
	UserContextKey = "user_context"

	// Cache key prefix and TTL for role scopes
	roleScopeCacheKeyPrefix = "rbac:role:%s:scopes"
	roleScopeCacheTTL       = 20 * time.Minute
)

// ScopeMiddleware resolves the full UserContext (permissions + data scopes)
// from the authenticated user's JWT claims. The resolved context is stored
// in gin.Context under the "user_context" key for downstream use.
//
// This middleware must be placed AFTER AuthMiddleware in the chain.
func ScopeMiddleware(
	permService *permissionservice.Service,
	roleRepo interfaces.RoleRepository,
	userRepo interfaces.UserRepository,
	brickRepo interfaces.BrickRepository,
	redisClient redis.UniversalClient,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract identity from auth middleware
		userIDVal, exists := c.Get("user_id")
		if !exists {
			errors.UnauthorizedResponse(c, "user not authenticated")
			c.Abort()
			return
		}
		userID, ok := userIDVal.(string)
		if !ok || userID == "" {
			errors.UnauthorizedResponse(c, "invalid user identity")
			c.Abort()
			return
		}

		roleCodeVal, _ := c.Get("user_role")
		jwtRoleCode, _ := roleCodeVal.(string)
		emailVal, _ := c.Get("user_email")
		email, _ := emailVal.(string)

		// Resolve current role info and scopes from the database. JWT role claims can be stale
		// after an admin changes a user's role, so DB-backed role is the authorization source.
		roleID, groupID, roleCode, scopes, err := resolveRoleAndScopes(userID, jwtRoleCode, roleRepo, userRepo, redisClient)
		if err != nil {
			errors.InternalServerErrorResponse(c, "Failed to resolve role scopes")
			c.Abort()
			return
		}

		// Resolve permissions from cache/DB
		perms, err := permService.GetPermissionsByRole(roleCode)
		if err != nil {
			errors.InternalServerErrorResponse(c, "Failed to resolve permissions")
			c.Abort()
			return
		}
		permMap := make(map[string]bool, len(perms))
		for _, p := range perms {
			permMap[p] = true
		}

		// Build UserContext
		userCtx := &domainauth.UserContext{
			UserID:      userID,
			Email:       email,
			RoleID:      roleID,
			RoleCode:    roleCode,
			GroupID:     groupID,
			Permissions: permMap,
			Scopes:      scopes,
		}

		// Pre-resolve team member IDs from bricks managed by this user.
		// A sales manager's "team" consists of all sales reps assigned to bricks
		// where this user is the brick manager.
		teamIDs := resolveBrickTeamMemberIDs(userID, brickRepo)
		if len(teamIDs) > 0 {
			userCtx.TeamMemberIDs = teamIDs
		}

		c.Set(UserContextKey, userCtx)
		c.Next()
	}
}

// resolveBrickTeamMemberIDs returns a deduplicated list of user IDs representing
// all sales reps assigned to bricks managed by the given userID, plus the user
// themselves. This is the authoritative team resolution for RBAC "team" scope.
func resolveBrickTeamMemberIDs(userID string, brickRepo interfaces.BrickRepository) []string {
	seen := make(map[string]struct{})
	seen[userID] = struct{}{} // always include the user themselves

	if brickRepo == nil {
		return []string{userID}
	}

	// Find all bricks where this user is the manager
	managerID := userID
	bricks, _, err := brickRepo.List(&brickdomain.ListBricksRequest{
		ManagerID: &managerID,
		Page:      1,
		PerPage:   100,
	})
	if err != nil {
		return []string{userID}
	}

	for _, b := range bricks {
		salesReps, salesErr := brickRepo.GetSalesByBrickID(b.ID)
		if salesErr == nil {
			for _, rep := range salesReps {
				seen[rep.ID] = struct{}{}
			}
		}
	}

	ids := make([]string, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	return ids
}

// resolveRoleAndScopes fetches role ID, group ID, and data scopes.
// Uses Redis cache for scope resolution to minimize DB load.
func resolveRoleAndScopes(
	userID, fallbackRoleCode string,
	roleRepo interfaces.RoleRepository,
	userRepo interfaces.UserRepository,
	redisClient redis.UniversalClient,
) (roleID, groupID, roleCode string, scopes map[string]roledomain.ScopeType, err error) {
	// Fetch user for role_id and group_id
	user, err := userRepo.FindByID(userID)
	if err != nil {
		return "", "", "", nil, fmt.Errorf("user not found: %w", err)
	}
	roleID = user.RoleID
	roleCode = fallbackRoleCode
	if user.Role != nil && user.Role.Code != "" {
		roleCode = user.Role.Code
	}
	if user.GroupID != nil {
		groupID = *user.GroupID
	}

	// Try cache for scopes
	scopes, err = getScopesFromCache(roleID, redisClient)
	if err == nil && scopes != nil {
		return roleID, groupID, roleCode, scopes, nil
	}

	// Cache miss: fetch from DB
	dbScopes, err := roleRepo.GetScopesByRoleID(roleID)
	if err != nil {
		return roleID, groupID, roleCode, nil, fmt.Errorf("failed to get scopes: %w", err)
	}

	scopes = make(map[string]roledomain.ScopeType, len(dbScopes))
	for _, s := range dbScopes {
		scopes[s.Resource] = s.Scope
	}

	// Cache the scopes
	cacheScopeData(roleID, scopes, redisClient)

	return roleID, groupID, roleCode, scopes, nil
}

// getScopesFromCache retrieves role scopes from Redis
func getScopesFromCache(roleCode string, redisClient redis.UniversalClient) (map[string]roledomain.ScopeType, error) {
	if redisClient == nil {
		return nil, fmt.Errorf("redis not available")
	}
	ctx := context.Background()
	cacheKey := fmt.Sprintf(roleScopeCacheKeyPrefix, roleCode)
	val, err := redisClient.Get(ctx, cacheKey).Result()
	if err != nil {
		return nil, err
	}
	var scopes map[string]roledomain.ScopeType
	if err := json.Unmarshal([]byte(val), &scopes); err != nil {
		return nil, err
	}
	return scopes, nil
}

// cacheScopeData stores role scopes in Redis
func cacheScopeData(roleCode string, scopes map[string]roledomain.ScopeType, redisClient redis.UniversalClient) {
	if redisClient == nil {
		return
	}
	ctx := context.Background()
	cacheKey := fmt.Sprintf(roleScopeCacheKeyPrefix, roleCode)
	data, err := json.Marshal(scopes)
	if err != nil {
		return
	}
	_ = redisClient.Set(ctx, cacheKey, string(data), roleScopeCacheTTL).Err()
}

// InvalidateScopeCache removes cached scopes for a role ID (call when scopes change)
func InvalidateScopeCache(roleID string, redisClient redis.UniversalClient) {
	if redisClient == nil {
		return
	}
	ctx := context.Background()
	cacheKey := fmt.Sprintf(roleScopeCacheKeyPrefix, roleID)
	_ = redisClient.Del(ctx, cacheKey).Err()
}

// GetUserContext extracts the resolved UserContext from gin.Context.
// Returns nil if the context has not been resolved (middleware not applied).
func GetUserContext(c *gin.Context) *domainauth.UserContext {
	val, exists := c.Get(UserContextKey)
	if !exists {
		return nil
	}
	ctx, ok := val.(*domainauth.UserContext)
	if !ok {
		return nil
	}
	return ctx
}
