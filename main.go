package main

import (
	"database/sql"
	"fmt"
	"log"

	"product-like/pkg/auth"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	dbConn, err := sql.Open("postgres", "postgres://postgres:postgrespass@localhost:5432/product_like_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	// Test the database connection
	err = dbConn.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	fmt.Println("Connected to the PostgreSQL database!")

	// Create a new router
	router := gin.Default()

	// // Set up authentication middleware
	// router.Use(auth.Middleware())

	authMiddleware, err := auth.NewMiddleware("my-secret-key")
	if err != nil {
		log.Fatalf("Failed to initialize auth middleware: %v", err)
	}

	// Initialize API routes
	apiRoutes := router.Group("/api")

	// Setup routes
	// setupRoutes(apiRoutes, authMiddleware)
	apiRoutes.GET("/products", api.GetProducts)
	apiRoutes.POST("/products", authMiddleware.Authorize(), api.CreateProduct)
	apiRoutes.PUT("/products/:id", authMiddleware.Authorize(), api.UpdateProduct)
	apiRoutes.DELETE("/products/:id", authMiddleware.Authorize(), api.DeleteProduct)

	// Start the server on localhost:8080
	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}

/*
func setupRoutes(apiRoutes *gin.RouterGroup, authMiddleware *auth.Middleware) {
	// Set up your routes here using apiRoutes
	apiRoutes.GET("/products", api.GetProducts)
	apiRoutes.POST("/products", authMiddleware.Authorize(), api.CreateProduct)
	apiRoutes.PUT("/products/:id", authMiddleware.Authorize(), api.UpdateProduct)
	apiRoutes.DELETE("/products/:id", authMiddleware.Authorize(), api.DeleteProduct)

	// Add more routes as needed
}
*/
