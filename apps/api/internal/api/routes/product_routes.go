package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gilabs/crm-healthcare/api/internal/api/handlers"
	"github.com/gilabs/crm-healthcare/api/internal/api/middleware"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
)

// SetupProductRoutes sets up product routes.
func SetupProductRoutes(router *gin.RouterGroup, productHandler *handlers.ProductHandler, jwtManager *jwt.JWTManager) {
	products := router.Group("/products")
	products.Use(middleware.AuthMiddleware(jwtManager))
	{
		products.GET("", productHandler.List)
		products.GET("/:id", productHandler.GetByID)
		products.POST("", productHandler.Create)
		products.PUT("/:id", productHandler.Update)
		products.DELETE("/:id", productHandler.Delete)
	}

	productCategories := router.Group("/product-categories")
	productCategories.Use(middleware.AuthMiddleware(jwtManager))
	{
		productCategories.GET("", productHandler.ListCategories)
		productCategories.GET("/:id", productHandler.GetCategoryByID)
		productCategories.POST("", productHandler.CreateCategory)
		productCategories.PUT("/:id", productHandler.UpdateCategory)
		productCategories.DELETE("/:id", productHandler.DeleteCategory)
	}
}


