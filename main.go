package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"product-like/api"
	"product-like/pkg/auth"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	gin.SetMode(gin.DebugMode)
	dbConn, err := sql.Open("postgres", "postgres://postgres:postgrespass@db:5432/product_like_db?sslmode=disable")
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

	authMiddleware, err := auth.NewMiddleware("my-secret-key")
	if err != nil {
		log.Fatalf("Failed to initialize auth middleware: %v", err)
	}

	// Initialize API routes
	apiRoutes := router.Group("/api")

	apiRoutes.POST("/like-product", authMiddleware.Authorize(), api.LikeProduct(dbConn))
	apiRoutes.GET("/liked-products", authMiddleware.Authorize(), api.RetrieveLikedProducts(dbConn))
	apiRoutes.POST("/cancel-like", authMiddleware.Authorize(), api.CancelProductLike(dbConn))

	apiRoutes.GET("/products", api.GetProducts(dbConn))
	apiRoutes.POST("/products", authMiddleware.Authorize(), api.CreateProduct(dbConn))
	apiRoutes.PUT("/products/:id", authMiddleware.Authorize(), api.UpdateProduct(dbConn))
	apiRoutes.DELETE("/products/:id", authMiddleware.Authorize(), api.DeleteProduct(dbConn))

	// Example: Generate a JWT token for user with ID 1234
	userID := int64(1234)
	token, err := GenerateToken(userID)
	if err != nil {
		log.Fatal("Failed to generate token:", err)
	}

	fmt.Printf("Generated token: %s\n", token)

	// Start the server on localhost:8080
	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}

func GenerateToken(userID int64) (string, error) {
	// expired in 1 hour
	expirationTime := time.Now().Add(1 * time.Hour)

	// Create the JWT claims
	claims := jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		Issuer:    "product-like-issuer",
		Subject:   strconv.FormatInt(userID, 10),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte("my-secret-key")
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
