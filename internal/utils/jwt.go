package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	// ErrInvaidToken is returned when JWT claims cannot be extracted from the request context
	// or when the claims are not of the expected type (jwt.MapClaims).
	ErrInvaidToken = errors.New("invalid token")
	// ErrEmptyUID is returned when the user ID (sub claim) is missing or empty in the JWT claims.
	ErrEmptyUID = errors.New("empty uid")
)

// GetJWTClaimsFromRequest extracts JWT claims from the Gin request context.
// The claims are expected to be set by the JWT authentication middleware with the key "claims".
// Returns the JWT claims as jwt.MapClaims or an error if the claims are not found or invalid.
//
// Parameters:
//   - c: Gin context containing the JWT claims set by the authentication middleware
//
// Returns:
//   - jwt.MapClaims: The JWT claims extracted from the context
//   - error: ErrInvaidToken if claims are not found or not of the expected type
func GetJWTClaimsFromRequest(c *gin.Context) (jwt.MapClaims, error) {
	tokenInfo, _ := c.Get("claims")
	claims, ok := tokenInfo.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvaidToken
	}

	return claims, nil
}

// GetUserIDFromRequest extracts the user ID from JWT claims in the Gin request context.
// It first retrieves the JWT claims, then extracts the "sub" (subject) claim which contains
// the user ID. Returns the user ID as a string or an error if extraction fails.
//
// Parameters:
//   - c: Gin context containing the JWT claims set by the authentication middleware
//
// Returns:
//   - string: The user ID extracted from the "sub" claim
//   - error: ErrInvaidToken if claims cannot be retrieved, or ErrEmptyUID if the "sub" claim
//            is missing, not a string, or empty
func GetUserIDFromRequest(c *gin.Context) (string, error) {
	claims, err := GetJWTClaimsFromRequest(c)
	if err != nil {
		return "", err
	}

	uid, ok := claims["sub"].(string)
	if !ok || uid == "" {
		return "", ErrEmptyUID
	}

	return uid, nil
}
