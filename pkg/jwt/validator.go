package jwt

import (
	"crypto/rsa"
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// errInvalidToken is returned when a JWT token cannot be validated or is malformed.
	errInvalidToken = errors.New("invalid token")
)

// JWTValidator defines the interface for JWT token validation.
// It provides methods to verify and parse JWT tokens using RSA public keys.
//
//go:generate mockery --name JWTValidator --filename jwt_validator.go
type JWTValidator interface {
	ValidateToken(string) (jwt.MapClaims, error)
}

type jwtValidator struct {
	publicKey *rsa.PublicKey
}

// NewJWTValidator creates a new JWT validator instance by loading and parsing
// an RSA public key from the specified file path. The public key is used
// to verify JWT token signatures.
func NewJWTValidator(publicKeyPath string) (JWTValidator, error) {
	publicKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return nil, err
	}

	return &jwtValidator{
		publicKey: publicKey,
	}, nil
}

// ValidateToken verifies the signature and validity of a JWT token string.
// It parses the token, validates it against the loaded public key, and returns
// the token claims if valid. Returns an error if the token is invalid or malformed.
func (v *jwtValidator) ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		return v.publicKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errInvalidToken
	}

	return token.Claims.(jwt.MapClaims), nil
}
