package api

import (
	"database/sql"
	"net/http"
	"product-like/db"       // Import your database package
	"product-like/pkg/auth" // Import the auth package
	"strconv"

	"github.com/gin-gonic/gin"
)

// LikeProduct allows a user to like a specific product.
// LikeProduct allows a user to like a specific product.
func LikeProduct(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse user ID from the JWT token
		//ctx := context.Background()

		user := auth.GetUserFromContext(c)

		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		userID := uint64(user.ID)

		// Parse product ID from the request (you need to extract it from the request body)
		var request struct {
			ProductID uint64 `json:"product_id"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		productID := request.ProductID

		// Check if the user has already liked the product (you'll need to query the database)
		liked, err := db.CheckProductLike(c.Request.Context(), dbConn, userID, uint64(productID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		if liked {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already liked the product"})
			return
		}

		// Add the like to the database (insert a record into the "favorites" table)
		if _, err := db.AddProductLike(c.Request.Context(), dbConn, uint64(userID), uint64(productID)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add like to the database"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product liked successfully", "user_id": userID})
	}
}

// RetrieveLikedProducts retrieves a list of products that the user liked.
func RetrieveLikedProducts(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse user ID from the JWT token
		user := auth.GetUserFromContext(c)

		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		userID := int(user.ID)

		// Parse pagination parameters from the request query
		page := c.DefaultQuery("page", "1")
		limit := c.DefaultQuery("limit", "10")

		// Convert page and limit to integers
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
			return
		}

		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
			return
		}

		// Retrieve liked products with pagination
		likedProducts, totalCount, err := db.RetrieveLikedProducts(c.Request.Context(), dbConn, userID, pageInt, limitInt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		// Create a response object
		response := gin.H{
			"liked_products": likedProducts,
			"total_count":    totalCount,
		}

		c.JSON(http.StatusOK, response)
	}
}

// CancelProductLike allows a user to cancel their like for a specific product.
func CancelProductLike(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse user ID from the JWT token
		user := auth.GetUserFromContext(c)

		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		userID := uint64(user.ID)

		// Parse product ID from the request (you need to extract it from the request body)
		var request struct {
			ProductID uint64 `json:"product_id"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		productID := request.ProductID

		// Check if the user has liked the product
		liked, err := db.CheckProductLike(c.Request.Context(), dbConn, userID, productID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		if !liked {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User has not liked the product"})
			return
		}

		// Cancel the like in the database
		if _, err := db.UnlikeProduct(c.Request.Context(), dbConn, userID, productID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel like in the database"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product like canceled successfully", "user_id": userID})
	}
}
