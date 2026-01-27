// Package utils provides utility functions for common operations such as password hashing and verification.
package utils

import "golang.org/x/crypto/bcrypt"

//go:generate mockery --name Hasher --filename hash.go
// Hasher defines the interface for password hashing operations.
// It provides methods to hash passwords and verify passwords against their hashes.
type Hasher interface {
	// HashPassword generates a hash from the given plain text password.
	// Returns the hashed password string.
	HashPassword(password string) string
	// VerifyPassword checks if the given plain text password matches the provided hash.
	// Returns true if the password matches the hash, false otherwise.
	VerifyPassword(password, hash string) bool
}

// bcryptHasher implements the Hasher interface using bcrypt algorithm.
type bcryptHasher struct{}

// HashPassword generates a bcrypt hash from the given plain text password.
// Uses bcrypt.DefaultCost for hashing. Returns the hashed password string.
func (b *bcryptHasher) HashPassword(password string) string {
	hashBytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashBytes)
}

// VerifyPassword verifies if the given plain text password matches the bcrypt hash.
// Returns true if the password matches the hash, false otherwise.
func (b *bcryptHasher) VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// NewHasher creates a new instance of Hasher using bcrypt implementation.
// Returns a Hasher interface that can be used for password hashing and verification.
func NewHasher() Hasher {
	return &bcryptHasher{}
}

// HashPassword is a convenience function that hashes a password using the default Hasher.
// It creates a new Hasher instance and hashes the given password.
// Returns the hashed password string.
func HashPassword(password string) string {
	hasher := NewHasher()
	return hasher.HashPassword(password)
}

// VerifyPassword is a convenience function that verifies a password against a hash using the default Hasher.
// It creates a new Hasher instance and verifies if the password matches the hash.
// Returns true if the password matches the hash, false otherwise.
func VerifyPassword(password, hash string) bool {
	hasher := NewHasher()
	return hasher.VerifyPassword(password, hash)
}
