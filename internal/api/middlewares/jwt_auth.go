// Package middlewares provides reusable HTTP middlewares for the API layer,
// including JWT authentication for protecting authenticated routes.
package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/luongtruong20201/bookmark-management/pkg/jwt"
)

// JWTAuth defines the interface for JWT authentication middleware.
// Implementations must validate Bearer tokens from the Authorization header
// and, on success, populate the Gin context with the authenticated user ID.
type JWTAuth interface {
	JWTAuth() gin.HandlerFunc
}

type jwtAuth struct {
	jwtValidator jwt.JWTValidator
}

// NewJWTAuth creates a new JWT authentication middleware instance using the
// provided JWT validator. The returned middleware can be attached to protected
// routes to enforce authentication.
func NewJWTAuth(jwtValidator jwt.JWTValidator) JWTAuth {
	return &jwtAuth{
		jwtValidator: jwtValidator,
	}
}

// JWTAuth returns a Gin handler function that:
//   - extracts the Authorization header in "Bearer <token>" format,
//   - validates the JWT using the configured validator,
//   - reads the "sub" claim as the user ID and stores it in the context as "userID",
//   - aborts the request with 401 status if any step fails.
func (m *jwtAuth) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format is wrong"})
			c.Abort()
			return
		}

		tokenStr := parts[1]
		tokenContent, err := m.jwtValidator.ValidateToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		userId, ok := tokenContent["sub"]
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", userId)
		c.Next()
	}
}
