package auth

import (
	"fmt"
	"net/http"
	"product-like/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a middleware for authentication and authorization
type AuthMiddleware struct {
	jwtSecret string
}

// NewMiddleware creates a new instance of AuthMiddleware
func NewMiddleware(jwtSecret string) (*AuthMiddleware, error) {
	return &AuthMiddleware{
		jwtSecret: jwtSecret,
	}, nil
}

// Authorize is a middleware function for authorization
func (m *AuthMiddleware) Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the JWT token from the request headers
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing JWT token"})
			c.Abort()
			return
		}

		// Parse and validate the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check the signing method and set the secret key
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Invalid token signing method")
			}
			return []byte(m.jwtSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT token"})
			c.Abort()
			return
		}

		// Check if the token is valid and not expired
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT token"})
			c.Abort()
			return
		}

		// Extract claims from the token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT claims"})
			c.Abort()
			return
		}

		// Extract user ID from claims (you can customize the claims structure)
		userID, ok := claims["sub"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in JWT claims"})
			c.Abort()
			return
		}

		// Set the user context with the extracted user ID
		c.Set("user_id", uint64(userID))

		// Continue with the request
		c.Next()
	}
}

// GetUserFromContext retrieves the user ID from the request context.
//
//	func GetUserFromContext(ctx context.Context) *User {
//		// Check if the context has a user ID value set
//		if userID, ok := ctx.Value("user_id").(uint64); ok {
//			// Replace this part with your logic to retrieve user information from the database.
//			// This is just a placeholder example.
//			user := &User{
//				ID:    userID,
//				Name:  name
//				Email: email,
//				// Add other user properties as needed
//			}
//			return user
//		}
//		return nil
//	}
//
// GetUserFromContext is a helper function to retrieve the authenticated user from the context.
func GetUserFromContext(c *gin.Context) *models.User {
	user, _ := c.Get("user")
	if user, ok := user.(*models.User); ok {
		return user
	}
	return nil
}
