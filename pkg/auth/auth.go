package auth

import (
	"fmt"
	"net/http"
	"product-like/models"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// middleware for authentication and authorization
type AuthMiddleware struct {
	jwtSecret string
}

// creates a new instance of AuthMiddleware
func NewMiddleware(jwtSecret string) (*AuthMiddleware, error) {
	return &AuthMiddleware{
		jwtSecret: jwtSecret,
	}, nil
}

// function for authorization
func (m *AuthMiddleware) Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the JWT token from the request headers
		tokenString := c.GetHeader("Authorization")
		// Remove the "Bearer " prefix if it exists
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		fmt.Print("tokenString: ", tokenString)
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
			fmt.Printf("\n JWT Parse Error: %s\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT token"})
			c.Abort()
			return
		}

		fmt.Print("\n token: ", token)
		// Check if the token is valid and not expired
		if !token.Valid {
			fmt.Printf("\n JWT Token Invalid!")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT token"})
			c.Abort()
			return
		}

		// Extract claims from the token
		claims, ok := token.Claims.(jwt.MapClaims)
		fmt.Println("\n Claims: ", claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT claims"})
			c.Abort()
			return
		}

		// Extract user ID from claims as a string
		userID, ok := claims["sub"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in JWT claims"})
			c.Abort()
			return
		}

		// Parse the user ID as needed (e.g., convert it to an integer)
		parsedUserID, err := strconv.Atoi(userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format in JWT claims"})
			c.Abort()
			return
		}

		// Convert the parsed user ID to uint64
		userUint64 := uint64(parsedUserID)

		// Set the user context with the extracted user ID
		c.Set("user_id", userUint64)

		// Continue with the request
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
