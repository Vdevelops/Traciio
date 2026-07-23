package permission

import (
	"context"
	"fmt"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/permission"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"github.com/redis/go-redis/v9"
)

// Cache key prefixes and TTL following enterprise best practices
const (
	// Role -> Permissions cache (CRITICAL: checked on every request)
	PermissionCacheKeyPrefix = "rbac:role:%s:permissions"
	PermissionCacheTTL       = 20 * time.Minute // Role->Permissions: 10-30 min recommended

	// User -> Role cache (CRITICAL: resolves user to role on every request)
	UserRoleCacheKeyPrefix = "rbac:user:%s:role"
	UserRoleCacheTTL       = 10 * time.Minute // User->Role: shorter TTL for quicker role change propagation

	// Permission list cache (master data, rarely changes)
	PermissionListCacheKey = "rbac:permissions:list"
	PermissionListCacheTTL = 30 * time.Minute
)

var (
	ErrPermissionNotFound = fmt.Errorf("permission not found")
	ErrUserNotFound       = fmt.Errorf("user not found")
)

type Service struct {
	repo        interfaces.PermissionRepository
	roleRepo    interfaces.RoleRepository
	userRepo    interfaces.UserRepository
	redisClient redis.UniversalClient
}

func NewService(repo interfaces.PermissionRepository, roleRepo interfaces.RoleRepository, userRepo interfaces.UserRepository, redisClient redis.UniversalClient) *Service {
	return &Service{
		repo:        repo,
		roleRepo:    roleRepo,
		userRepo:    userRepo,
		redisClient: redisClient,
	}
}

// GetPermissionsByRole returns list of permissions for a specific role, using cache
func (s *Service) GetPermissionsByRole(roleCode string) ([]string, error) {
	// Admin roles should always have all permissions.
	// This also prevents stale Redis role-permissions cache from hiding newly added menus.
	if roleCode == "admin" || roleCode == "super_admin" {
		allPerms, err := s.repo.List()
		if err != nil {
			return nil, err
		}

		codes := make([]string, 0, len(allPerms))
		for _, p := range allPerms {
			if p.Code != "" {
				codes = append(codes, p.Code)
			} else if p.Resource != "" && p.Action != "" {
				codes = append(codes, fmt.Sprintf("%s.%s", p.Resource, p.Action))
			}
		}

		// Best-effort cache for admin too (optional). If Redis is unavailable, just return.
		ctx := context.Background()
		cacheKey := fmt.Sprintf(PermissionCacheKeyPrefix, roleCode)
		if s.redisClient != nil && len(codes) > 0 {
			pipe := s.redisClient.Pipeline()
			// Replace set contents to avoid mixing stale/partial data
			pipe.Del(ctx, cacheKey)
			pipe.SAdd(ctx, cacheKey, codes)
			pipe.Expire(ctx, cacheKey, PermissionCacheTTL)
			_, _ = pipe.Exec(ctx)
		}

		return codes, nil
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf(PermissionCacheKeyPrefix, roleCode)

	// Try to get from redis
	if s.redisClient != nil {
		val, err := s.redisClient.SMembers(ctx, cacheKey).Result()
		if err == nil && len(val) > 0 {
			return val, nil
		}
	}

	// If not in cache, get from DB
	r, err := s.roleRepo.FindByCode(roleCode)
	if err != nil {
		return nil, err
	}

	// Assuming FindByCode preloads permissions, if not we need to fetch them
	perms := []string{}
	for _, p := range r.Permissions {
		// Use Code (resource.action)
		if p.Code != "" {
			perms = append(perms, p.Code)
		} else {
			// Fallback: construct from Resource + Action if Code is empty
			if p.Resource != "" && p.Action != "" {
				perms = append(perms, fmt.Sprintf("%s.%s", p.Resource, p.Action))
			}
		}
	}

	// Cache in Redis
	if s.redisClient != nil && len(perms) > 0 {
		pipe := s.redisClient.Pipeline()
		pipe.SAdd(ctx, cacheKey, perms)
		pipe.Expire(ctx, cacheKey, PermissionCacheTTL)
		_, _ = pipe.Exec(ctx)
	}

	return perms, nil
}

// InvalidateRoleCache invalidates cache for a specific role
func (s *Service) InvalidateRoleCache(roleCode string) error {
	if s.redisClient == nil {
		return nil
	}
	ctx := context.Background()
	cacheKey := fmt.Sprintf(PermissionCacheKeyPrefix, roleCode)
	return s.redisClient.Del(ctx, cacheKey).Err()
}

// List returns all available permissions with simplified response
func (s *Service) List() ([]permission.PermissionSimpleResponse, error) {
	perms, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	var resp []permission.PermissionSimpleResponse
	for _, p := range perms {
		resp = append(resp, *p.ToPermissionSimpleResponse())
	}
	return resp, nil
}

// GetByID returns a permission by ID
func (s *Service) GetByID(id string) (*permission.PermissionResponse, error) {
	perm, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return perm.ToPermissionResponse(), nil
}

// GetUserPermissions returns permissions for a user with caching
func (s *Service) GetUserPermissions(userID string) ([]string, error) {
	// Try to get user's role from cache first
	roleCode, err := s.getUserRoleFromCache(userID)
	if err != nil || roleCode == "" {
		// Cache miss or error - get from DB
		user, err := s.userRepo.FindByID(userID)
		if err != nil {
			return nil, ErrUserNotFound
		}

		if user.Role == nil {
			return []string{}, nil
		}

		roleCode = user.Role.Code
		// Cache the user->role mapping
		s.cacheUserRole(userID, roleCode)
	}

	return s.GetPermissionsByRole(roleCode)
}

// getUserRoleFromCache retrieves user's role code from Redis cache
func (s *Service) getUserRoleFromCache(userID string) (string, error) {
	if s.redisClient == nil {
		return "", nil
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf(UserRoleCacheKeyPrefix, userID)
	return s.redisClient.Get(ctx, cacheKey).Result()
}

// cacheUserRole caches the user->role mapping in Redis
func (s *Service) cacheUserRole(userID, roleCode string) {
	if s.redisClient == nil || roleCode == "" {
		return
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf(UserRoleCacheKeyPrefix, userID)
	_ = s.redisClient.Set(ctx, cacheKey, roleCode, UserRoleCacheTTL).Err()
}

// InvalidateUserCache invalidates cache for a specific user (call when user's role changes)
func (s *Service) InvalidateUserCache(userID string) error {
	if s.redisClient == nil {
		return nil
	}
	ctx := context.Background()
	cacheKey := fmt.Sprintf(UserRoleCacheKeyPrefix, userID)
	return s.redisClient.Del(ctx, cacheKey).Err()
}

// InvalidateAllUserCachesForRole invalidates cache for all users with a specific role
// Call this when role permissions are updated
func (s *Service) InvalidateAllUserCachesForRole(roleID string) error {
	if s.redisClient == nil {
		return nil
	}

	// Get all users with this role
	userIDs, err := s.userRepo.GetUsersByRoleID(roleID)
	if err != nil {
		return err
	}

	// Invalidate cache for each user
	ctx := context.Background()
	for _, userID := range userIDs {
		cacheKey := fmt.Sprintf(UserRoleCacheKeyPrefix, userID)
		_ = s.redisClient.Del(ctx, cacheKey).Err()
	}

	return nil
}
