package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func SetupAIRoutes(v1 *gin.RouterGroup, aiHandler *handlers.AIHandler, aiSettingsHandler *handlers.AISettingsHandler, jwtManager *jwt.JWTManager) {
	ai := v1.Group("/ai")
	ai.Use(middleware.AuthMiddleware(jwtManager))

	{
		// Visit Report Insights
		ai.POST("/analyze/visit-report", aiHandler.AnalyzeVisitReport)

		// Chat
		ai.POST("/chat", aiHandler.Chat)

		// Settings
		ai.GET("/settings", aiSettingsHandler.GetSettings)
		ai.PUT("/settings", aiSettingsHandler.UpdateSettings)
	}
}

