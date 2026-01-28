package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBcryptHasher_HashPassword(t *testing.T) {
	t.Parallel()

	hasher := NewHasher()

	testCases := []struct {
		name     string
		password string
		verify   func(t *testing.T, hash string)
	}{
		{
			name:     "success - hash password",
			password: "password123",
			verify: func(t *testing.T, hash string) {
				assert.NotEmpty(t, hash)
				assert.NotEqual(t, "password123", hash)
				assert.Contains(t, hash, "$2a$")
			},
		},
		{
			name:     "success - hash empty password",
			password: "",
			verify: func(t *testing.T, hash string) {
				assert.NotEmpty(t, hash)
				assert.Contains(t, hash, "$2a$")
			},
		},
		{
			name:     "success - hash long password (within bcrypt limit)",
			password: "this is a long password that is still within bcrypt 72 byte limit",
			verify: func(t *testing.T, hash string) {
				assert.NotEmpty(t, hash)
				assert.Contains(t, hash, "$2a$")
			},
		},
		{
			name:     "success - hash special characters",
			password: "p@ssw0rd!#$%^&*()",
			verify: func(t *testing.T, hash string) {
				assert.NotEmpty(t, hash)
				assert.Contains(t, hash, "$2a$")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			hash := hasher.HashPassword(tc.password)
			tc.verify(t, hash)
		})
	}
}

func TestBcryptHasher_VerifyPassword(t *testing.T) {
	t.Parallel()

	hasher := NewHasher()

	testCases := []struct {
		name          string
		password      string
		hash          string
		expectedMatch bool
		setupHash     bool
	}{
		{
			name:          "success - verify correct password",
			password:      "password123",
			setupHash:     true,
			expectedMatch: true,
		},
		{
			name:          "error - verify incorrect password",
			password:      "wrongpassword",
			setupHash:     true,
			expectedMatch: false,
		},
		{
			name:          "success - verify with empty password (bcrypt supports empty passwords)",
			password:      "",
			setupHash:     true,
			expectedMatch: true,
		},
		{
			name:          "error - verify with invalid hash",
			password:      "password123",
			hash:          "invalid-hash",
			setupHash:     false,
			expectedMatch: false,
		},
		{
			name:          "error - verify with empty hash",
			password:      "password123",
			hash:          "",
			setupHash:     false,
			expectedMatch: false,
		},
		{
			name:          "error - verify with malformed hash",
			password:      "password123",
			hash:          "$2a$10$short",
			setupHash:     false,
			expectedMatch: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var hash string
			if tc.setupHash {
				hash = hasher.HashPassword(tc.password)
				if !tc.expectedMatch && tc.name == "error - verify incorrect password" {
					hash = hasher.HashPassword("differentpassword")
				}
			} else {
				hash = tc.hash
			}

			result := hasher.VerifyPassword(tc.password, hash)
			assert.Equal(t, tc.expectedMatch, result)
		})
	}
}

func TestHashPassword(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		password string
		verify   func(t *testing.T, hash string)
	}{
		{
			name:     "success - convenience function works",
			password: "testpassword",
			verify: func(t *testing.T, hash string) {
				assert.NotEmpty(t, hash)
				assert.Contains(t, hash, "$2a$")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			hash := HashPassword(tc.password)
			tc.verify(t, hash)
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		password      string
		setupHash     bool
		expectedMatch bool
	}{
		{
			name:          "success - convenience function works with correct password",
			password:      "testpassword",
			setupHash:     true,
			expectedMatch: true,
		},
		{
			name:          "error - convenience function works with incorrect password",
			password:      "wrongpassword",
			setupHash:     true,
			expectedMatch: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var hash string
			if tc.setupHash {
				if tc.expectedMatch {
					hash = HashPassword(tc.password)
				} else {
					hash = HashPassword("differentpassword")
				}
			}

			result := VerifyPassword(tc.password, hash)
			assert.Equal(t, tc.expectedMatch, result)
		})
	}
}
