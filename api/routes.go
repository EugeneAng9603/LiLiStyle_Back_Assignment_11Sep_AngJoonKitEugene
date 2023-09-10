// api/routes.go

package api

import (
	"database/sql"
	//"product-like/api/handlers"

	"github.com/gin-gonic/gin"
)

// mux or gin

func SetupRoutes(router *gin.Engine, db *sql.DB) {
	// Create a new group for user-related routes that require authentication
	//userGroup := router.Group("/users")
	// userGroup.Use(auth.AuthMiddleware()) // Apply authentication middleware

	// Define user-related routes
	// userGroup.POST("/like", handlers.LikeProduct)
	// userGroup.GET("/liked", handlers.RetrieveLikedProducts)
	// userGroup.DELETE("/unlike", handlers.CancelProductLike)
}
