package service

import (
	"bytes"
	"crypto/rand"
	"math/big"
)

const (
	charset    = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	passLength = 10
)

type passwordService struct{}

// Password interface represents password service
//
//go:generate mockery --name Password --filename password_service.go
type Password interface {
	GeneratePassword() (string, error)
}

// NewPassword return a new instance of the password service
func NewPassword() Password {
	return &passwordService{}
}

// GeneratePassword generates a random password of the specified length using
// alphanumeric characters. Returns an error if random number generation fails.
func (s *passwordService) GeneratePassword() (string, error) {
	var sb bytes.Buffer

	for range passLength {
		i, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		sb.WriteByte(charset[i.Int64()])
	}

	return sb.String(), nil
}
