package main

import (
	"database/sql"
	"fmt"
	"log"
	"product-like/api/handlers"

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

	// API routes
	apiRoutes := router.Group("/api")
	{
		apiRoutes.POST("/like-product", handlers.LikeProduct(dbConn))
		apiRoutes.GET("/liked-products", handlers.RetrieveLikedProducts(dbConn))
		apiRoutes.DELETE("/cancel-like", handlers.CancelProductLike(dbConn))
	}

	// Start the server on localhost:8080
	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}
