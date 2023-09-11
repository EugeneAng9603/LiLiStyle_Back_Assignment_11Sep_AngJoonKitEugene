package api

import (
	"database/sql"
	"log"
	"net/http"
	"product-like/db"
	"product-like/pkg/auth"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// LikeProduct allows user to like a specific product.
func LikeProduct(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		user := auth.GetUserFromContext(c)

		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		userID := uint64(user.ID)

		var request struct {
			ProductID uint64 `json:"product_id"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		productID := request.ProductID

		// Check if the user has already liked the product
		liked, err := db.CheckProductLike(c.Request.Context(), dbConn, userID, uint64(productID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		if liked {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already liked the product"})
			return
		}

		// Update like to the database
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

		user := auth.GetUserFromContext(c)

		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		userID := int(user.ID)

		// Parse pagination parameters from the request query
		page := c.DefaultQuery("page", "1")
		limit := c.DefaultQuery("limit", "10")

		// Convert page and limit from string to integers
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

type Product struct {
	ID                uint64    `db:"id"`
	ShopID            uint64    `db:"shop_id"`
	Name              string    `db:"name"`
	Description       string    `db:"description"`
	ThumbnailURL      string    `db:"thumbnail_url"`
	OriginPrice       int64     `db:"origin_price"`
	DiscountedPrice   int64     `db:"discounted_price"`
	DiscountedRate    float64   `db:"discounted_rate"`
	Status            string    `db:"status"`
	InStock           bool      `db:"in_stock"`
	IsPreorder        bool      `db:"is_preorder"`
	IsPurchasable     bool      `db:"is_purchasable"`
	DeliveryCondition string    `db:"delivery_condition"`
	DeliveryDisplay   string    `db:"delivery_display"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}

// GetProducts retrieves a list of products.
func GetProducts(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		rows, err := dbConn.QueryContext(c.Request.Context(), "SELECT * FROM products")
		if err != nil {
			log.Println("Error fetching products:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
			return
		}
		defer rows.Close()

		// Iterate through the rows and build a slice of products.
		var products []Product
		for rows.Next() {
			var product Product
			err := rows.Scan(
				&product.ID,
				&product.ShopID,
				&product.Name,
				&product.Description,
				&product.ThumbnailURL,
				&product.OriginPrice,
				&product.DiscountedPrice,
				&product.DiscountedRate,
				&product.Status,
				&product.InStock,
				&product.IsPreorder,
				&product.IsPurchasable,
				&product.DeliveryCondition,
				&product.DeliveryDisplay,
				&product.CreatedAt,
				&product.UpdatedAt,
			)
			if err != nil {
				log.Println("Error fetching products:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
				return
			}
			products = append(products, product)
		}

		// Return the list of products as a JSON response.
		c.JSON(http.StatusOK, gin.H{"products": products})
	}
}

// CreateProduct creates a new product.
func CreateProduct(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Define a struct to represent the request body data.
		var request struct {
			Name        string  `json:"name" binding:"required"`
			Description string  `json:"description" binding:"required"`
			Price       float64 `json:"origin_price" binding:"required"`
		}

		// Bind the request body to the request struct.
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Insert the new product into the database.
		query := `
			INSERT INTO products (name, description, price)
			VALUES ($1, $2, $3)
			RETURNING id
		`

		var productID int64
		err := dbConn.QueryRowContext(
			c.Request.Context(),
			query,
			request.Name,
			request.Description,
			request.Price,
		).Scan(&productID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create the product"})
			return
		}

		// Return a success response with the created product ID.
		c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully", "product_id": productID})
	}
}

// UpdateProduct updates an existing product.
func UpdateProduct(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the product ID from the URL parameters
		productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		// Define a struct to represent the request body data for updates.
		var request struct {
			Name        string  `json:"name"`
			Description string  `json:"description"`
			Price       float64 `json:"price"`
		}

		// Bind the request body to the request struct.
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update the product in the database.
		query := `
			UPDATE products
			SET name = COALESCE($1, name), description = COALESCE($2, description), price = COALESCE($3, price)
			WHERE id = $4
		`

		result, err := dbConn.ExecContext(
			c.Request.Context(),
			query,
			request.Name,
			request.Description,
			request.Price,
			productID,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the product"})
			return
		}

		// Check if the product was updated successfully.
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		// Return a success response.
		c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
	}
}

// DeleteProduct deletes a product by its ID.
func DeleteProduct(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the product ID from the URL parameters
		productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		// Delete the product from the database
		query := `
			DELETE FROM products
			WHERE id = $1
		`

		result, err := dbConn.ExecContext(c.Request.Context(), query, productID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the product"})
			return
		}

		// Check if the product was deleted successfully
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		// Return a success response
		c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
	}
}

// CreateUser creates a new user and returns the user's ID.
func CreateUser(dbConn *sql.DB, username, password string) (int64, error) {
	// Hash the user's password before storing it in the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	// Insert the user into the database
	query := `
        INSERT INTO users (username, password)
        VALUES ($1, $2)
        RETURNING id
    `
	var userID int64
	err = dbConn.QueryRow(query, username, hashedPassword).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
