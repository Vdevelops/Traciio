package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.RouterGroup, authHandler *handlers.AuthHandler, permissionHandler *handlers.PermissionHandler, userHandler *handlers.UserHandler, jwtManager *jwt.JWTManager) {
	auth := router.Group("/auth")
	{
		// Login endpoint with rate limiting (5 requests per 15 minutes)
		auth.POST("/login", middleware.RateLimitMiddleware("login"), authHandler.Login)
		// Mobile login endpoint with rate limiting (5 requests per 15 minutes)
		// Only users with role that has mobile_access enabled can login via this endpoint
		auth.POST("/mobile/login", middleware.RateLimitMiddleware("mobile_login"), authHandler.MobileLogin)
		// Refresh token endpoint with rate limiting (10 requests per hour)
		auth.POST("/refresh", middleware.RateLimitMiddleware("refresh"), authHandler.RefreshToken)
		// Logout endpoint (no auth required - should work even with expired token)
		auth.POST("/logout", authHandler.Logout)
		
		// Mobile permissions endpoint (authenticated)
		auth.GET("/mobile/permissions", middleware.AuthMiddleware(jwtManager), permissionHandler.GetMobilePermissions)
		
		// Mobile profile endpoints (authenticated)
		mobile := auth.Group("/mobile")
		mobile.Use(middleware.AuthMiddleware(jwtManager))
		{
			mobile.GET("/profile", userHandler.GetMyProfile)
			mobile.PUT("/profile", userHandler.UpdateMyProfile)
			mobile.PUT("/password", userHandler.ChangeMyPassword)
		}
	}
}

