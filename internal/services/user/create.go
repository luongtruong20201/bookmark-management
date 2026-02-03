package user

import (
	"context"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
)

// CreateUser creates a new user account with the provided information.
// It hashes the password using bcrypt before storing the user in the database.
// The user ID is automatically generated as a UUID by the repository layer.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - username: Unique username for the user account
//   - password: Plain text password to be hashed
//   - displayName: Display name for the user
//   - email: Unique email address for the user
//
// Returns:
//   - *model.User: The created user with generated UUID and hashed password
//   - error: Returns an error if user creation fails (e.g., duplicate username/email, database error)
func (u *user) CreateUser(ctx context.Context, username, password, displayName, email string) (*model.User, error) {
	hashPassword := u.hasher.HashPassword(password)

	newUser := model.User{
		Username:    username,
		DisplayName: displayName,
		Email:       email,
		Password:    hashPassword,
	}

	res, err := u.repo.CreateUser(ctx, &newUser)
	if err != nil {
		return nil, err
	}

	return res, nil
}
