package stringutils

import (
	"bytes"
	"crypto/rand"
	"math/big"
)

const (
	// charset contains the alphanumeric characters used for generating random codes.
	charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

// KeyGenerator defines the interface for generating random codes.
// It provides methods to generate random alphanumeric strings of a specified length.
//
//go:generate mockery --name KeyGenerator --filename key_generator.go
type KeyGenerator interface {
	GenerateCode(int) (string, error)
}

type keyGen struct{}

// NewKeyGen creates a new key generator instance.
func NewKeyGen() KeyGenerator {
	return &keyGen{}
}

// GenerateCode generates a random code of the specified length using the keyGen implementation.
func (k *keyGen) GenerateCode(length int) (string, error) {
	return GenerateCode(length)
}

// GenerateCode generates a random alphanumeric string of the specified length.
// It uses cryptographically secure random number generation to select characters from the charset.
// Returns an error if random number generation fails.
func GenerateCode(length int) (string, error) {
	var sb bytes.Buffer

	for range length {
		i, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		sb.WriteByte(charset[i.Int64()])
	}

	return sb.String(), nil
}
