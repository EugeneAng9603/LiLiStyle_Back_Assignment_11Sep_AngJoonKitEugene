package auth

import (
	"database/sql" // Import your database package
	"net/http"
	"product-like/models" // Import your user model

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a middleware function to handle JWT authentication.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}

		claims, err := ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Check if the user exists in the database (you need to implement this function)
		user, err := GetUserByID(uint64(claims.UserID))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Set the user in the context for further request handling
		c.Set("user", user)
		c.Next()
	}
}

// GetUserFromContext is a helper function to retrieve the authenticated user from the context.
func GetUserFromContext(c *gin.Context) *models.User {
	user, _ := c.Get("user")
	if user, ok := user.(*models.User); ok {
		return user
	}
	return nil
}

var db *sql.DB

// GetUserByID retrieves a user from the database by their ID.
func GetUserByID(userID uint64) (*models.User, error) {
	// Define your SQL query to retrieve the user by ID
	query := "SELECT id, name, email, password, phone, status, created_at, updated_at, deleted_at FROM users WHERE id = $1"

	var user models.User

	// Execute the query and scan the result into the user struct
	err := db.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Phone,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// Handle the case where the user is not found
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
