package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Login authenticates a user with the provided username and password.
// It retrieves the user from the database, verifies the password hash, and generates a JWT token
// upon successful authentication. The token includes the user ID and expiration time.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - username: Username of the user attempting to log in
//   - password: Plain text password to verify against the stored hash
//
// Returns:
//   - string: JWT token string for authenticated requests (valid for 24 hours)
//   - error: Returns ErrClientErr if credentials are invalid or user doesn't exist,
//     or an error if token generation fails
func (u *user) Login(ctx context.Context, username, password string) (string, error) {
	user, err := u.repo.GetUserByUsername(ctx, username)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return "", ErrClientErr
		default:
			return "", err
		}
	}
	if check := u.hasher.VerifyPassword(password, user.Password); !check {
		return "", ErrClientErr
	}

	jwtContent := jwt.MapClaims{
		"sub": user.ID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(tokenExpiresTime).Unix(),
	}
	token, err := u.jwtGenerator.GenerateToken(jwtContent)
	if err != nil {
		return "", err
	}

	return token, nil
}
