package middleware

import (
	"net/http"
	"strings"
	"time"

	"test-wallet/config"
	"test-wallet/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the claims in the JWT token
type JWTClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken generates a new JWT token for the given user ID
func GenerateToken(userID string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.AppConfig.JWTConfig.Expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTConfig.Secret))
}

// AuthMiddleware is a middleware that checks for a valid JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.LogError(nil, "Authorization header missing", map[string]interface{}{
				"path": c.Request.URL.Path,
			})
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.LogError(nil, "Invalid authorization header format", map[string]interface{}{
				"path": c.Request.URL.Path,
			})
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims := &JWTClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWTConfig.Secret), nil
		})

		if err != nil || !token.Valid {
			utils.LogError(err, "Invalid token", map[string]interface{}{
				"path": c.Request.URL.Path,
			})
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set the user ID in the context for use in handlers
		c.Set("user_id", claims.UserID)
		utils.LogDebug("User authenticated", map[string]interface{}{
			"user_id": claims.UserID,
			"path":    c.Request.URL.Path,
		})
		c.Next()
	}
}
