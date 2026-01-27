package jwt

import (
	"crypto/rsa"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// JWTGenerator defines the interface for JWT token generation.
// It provides methods to create signed JWT tokens using RSA private keys.
type JWTGenerator interface {
	GenerateToken(jwt.MapClaims) (string, error)
}

type jwtGenerator struct {
	privateKey *rsa.PrivateKey
}

// NewJWTGenerator creates a new JWT generator instance by loading and parsing
// an RSA private key from the specified file path. The private key is used
// to sign JWT tokens using RS256 algorithm.
func NewJWTGenerator(privateKeyPath string) (JWTGenerator, error) {
	privateKeyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return nil, err
	}

	return &jwtGenerator{
		privateKey: privateKey,
	}, nil
}

// GenerateToken creates a new JWT token with the provided claims and signs it
// using the RS256 algorithm with the loaded private key. Returns the signed
// token string or an error if signing fails.
func (g *jwtGenerator) GenerateToken(jwtContent jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwtContent)
	tokenString, err := token.SignedString(g.privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
