package handlers

import (
	"net/http"
	"product-like/db"       // Import your database package
	"product-like/pkg/auth" // Import the auth package
	"strconv"

	"github.com/gin-gonic/gin"
)

// LikeProduct allows a user to like a specific product.
func LikeProduct(c *gin.Context) {
	// Parse user ID from the JWT token
	//ctx := context.Background()

	user := auth.GetUserFromContext(c)

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := user.ID

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
	liked, err := db.CheckProductLike(c.Request.Context(), db.DB, userID, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	if liked {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already liked the product"})
		return
	}

	// Add the like to the database (insert a record into the "favorites" table)
	if _, err := db.AddProductLike(c.Request.Context(), db.DB, userID, productID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add like to the database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product liked successfully", "user_id": userID})
}

// RetrieveLikedProducts retrieves a list of products that the user liked.
func RetrieveLikedProducts(c *gin.Context) {
	// Parse user ID from the JWT token
	user := auth.GetUserFromContext(c)

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := user.ID

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
	likedProducts, totalCount, err := db.RetrieveLikedProducts(c.Request.Context(), db.DB, userID, pageInt, limitInt)
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

// func RetrieveLikedProducts(c *gin.Context) {
// 	// Parse user ID from the JWT token
// 	user := auth.GetUserFromContext(c)

// 	if user == nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	userID := user.ID

// 	// Implement pagination (you need to extract page and limit from the query parameters)
// 	page := 1   // Example: Extract from query parameter
// 	limit := 10 // Example: Extract from query parameter

// 	// Retrieve liked products and total count from the database
// 	likedProducts, totalCount, err := db.RetrieveLikedProducts(c.Request.Context(), db.DB, userID, page, limit)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve liked products"})
// 		return
// 	}

// 	// Create a response with the list of liked products and total count
// 	var response struct {
// 		Products   []models.Product `json:"products"`
// 		TotalCount int              `json:"total_count"`
// 	}

// 	response.Products = likedProducts
// 	response.TotalCount = totalCount

// 	c.JSON(http.StatusOK, response)
// }

// CancelProductLike allows a user to cancel their like for a specific product.
func CancelProductLike(c *gin.Context) {
	// Parse user ID from the JWT token
	user := auth.GetUserFromContext(c)

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := user.ID

	// Parse product ID from the request (you need to extract it from the request body)
	var request struct {
		ProductID uint64 `json:"product_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productID := request.ProductID

	// Check if the user has liked the product (you'll need to query the database)
	liked, err := db.CheckProductLike(c.Request.Context(), db.DB, userID, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	if !liked {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User has not liked the product"})
		return
	}

	// Cancel the like in the database (delete the record from the "favorites" table)
	if _, err := db.UnlikeProduct(c.Request.Context(), db.DB, userID, productID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel like in the database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product like canceled successfully", "user_id": userID})
}

// api/handlers/product_handlers.go

// package handlers

// import (
// 	"net/http"
// 	"product-like/pkg/auth"
// 	"strconv"
// 	"your-project/db/sqlc"

// 	"github.com/gorilla/mux"
// )

// // LikeProduct likes a specific product for the authenticated user.
// func LikeProduct(db *sqlc.Queries) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// Parse user ID from the JWT token
// 		userID := auth.GetUserFromContext(c).ID
// 		// Get the user ID from the JWT token (use auth package)

// 		// Parse product ID from the request URL
// 		vars := mux.Vars(r)
// 		productIDStr := vars["productID"]
// 		productID, err := strconv.ParseInt(productIDStr, 10, 64)
// 		if err != nil {
// 			http.Error(w, "Invalid product ID", http.StatusBadRequest)
// 			return
// 		}

// 		// Check if the user has liked the product
// 		liked, err := db.CheckProductLike(r.Context(), sqlc.CheckProductLikeParams{
// 			UserID:    userID,
// 			ProductID: productID,
// 		})
// 		if err != nil {
// 			http.Error(w, "Error checking if the product is liked", http.StatusInternalServerError)
// 			return
// 		}

// 		// If not liked, return an error response
// 		if !liked {
// 			http.Error(w, "Product is not liked by the user", http.StatusNotFound)
// 			return
// 		}

// 		// Return success response
// 		w.WriteHeader(http.StatusNoContent)
// 	}
// }

// // RetrieveLikedProducts retrieves a list of products that the authenticated user has liked.
// func RetrieveLikedProducts(db *sqlc.Queries) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		userID := auth.GetUserFromContext(c).ID

// 		// Query the database to retrieve liked products for the user

// 		// Serialize and return the list of liked products as JSON
// 	}
// }

// // CancelProductLike cancels the like for a specific product by the authenticated user.
// func CancelProductLike(db *sqlc.Queries) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// Parse user ID from the JWT token
// 		userID := auth.GetUserFromContext(c).ID

// 		// Parse product ID from the request URL
// 		vars := mux.Vars(r)
// 		productIDStr := vars["productID"]
// 		productID, err := strconv.ParseInt(productIDStr, 10, 64)
// 		if err != nil {
// 			http.Error(w, "Invalid product ID", http.StatusBadRequest)
// 			return
// 		}

// 		// Check if the user has liked the product (you'll need to query the database)

// 		// If liked, delete the record from the favorites table

// 		// Return success response
// 		w.WriteHeader(http.StatusNoContent)
// 	}
// }
