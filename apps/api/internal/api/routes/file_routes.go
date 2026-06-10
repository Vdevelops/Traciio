package routes

import (
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func SetupFileRoutes(router *gin.RouterGroup, fileHandler *handlers.FileHandler, jwtManager *jwt.JWTManager) {
	files := router.Group("/files")
	{
		files.GET("/image/*filepath", fileHandler.ServeImage)
	}

	upload := router.Group("/upload")
	upload.Use(middleware.AuthMiddleware(jwtManager))
	{
		upload.POST("/image", fileHandler.UploadImage)
		upload.DELETE("/file/:filename", fileHandler.DeleteFile)
	}
}
