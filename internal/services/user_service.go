// Package service provides business logic implementations for user-related operations.
// It handles user creation, authentication, and profile retrieval with password hashing
// and JWT token generation.
package service

import (
	"context"
	db "database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	model "github.com/luongtruong20201/bookmark-management/internal/models"
	repository "github.com/luongtruong20201/bookmark-management/internal/repositories"
	jwtPkg "github.com/luongtruong20201/bookmark-management/pkg/jwt"

	"github.com/luongtruong20201/bookmark-management/pkg/utils"
)

const (
	// tokenExpiresTime defines the expiration duration for JWT tokens.
	// Tokens are valid for 24 hours from the time of generation.
	tokenExpiresTime = 24 * time.Hour
)

// User defines the interface for user service operations.
// It provides methods to handle user-related business logic including user creation,
// authentication, and profile retrieval.
//
//go:generate mockery --name User --filename user_service.go
type User interface {
	// CreateUser creates a new user account with hashed password.
	// It validates input, hashes the password, and persists the user to the database.
	// Returns the created user with generated UUID or an error if creation fails.
	CreateUser(ctx context.Context, username, password, displayName, email string) (*model.User, error)

	// Login authenticates a user with username and password.
	// It verifies credentials, and upon successful authentication, generates and returns a JWT token.
	// Returns the JWT token string or an error if authentication fails.
	Login(ctx context.Context, username, password string) (string, error)

	// GetUserByID retrieves a user by their unique identifier.
	// Returns the user information or an error if the user is not found.
	GetUserByID(ctx context.Context, id string) (*model.User, error)

	// UpdateUserProfile updates the display name and email of a user identified by ID.
	// Returns the updated user information or an error if the update fails.
	UpdateUserProfile(ctx context.Context, id, displayName, email string) (*model.User, error)
}

// user implements the User interface and provides business logic for user operations.
// It encapsulates dependencies for repository access, password hashing, and JWT token generation.
type user struct {
	repo         repository.User
	hasher       utils.Hasher
	jwtGenerator jwtPkg.JWTGenerator
}

// NewUser creates a new user service instance with the provided dependencies.
// It initializes the service with a user repository, password hasher, and JWT generator.
//
// Parameters:
//   - repo: Repository interface for database operations
//   - hasher: Hasher interface for password hashing and verification
//   - jwtGenerator: JWT generator for creating authentication tokens
//
// Returns:
//   - User: A new user service instance implementing the User interface
func NewUser(
	repo repository.User,
	hasher utils.Hasher,
	jwtGenerator jwtPkg.JWTGenerator,
) User {
	return &user{
		repo:         repo,
		hasher:       hasher,
		jwtGenerator: jwtGenerator,
	}
}

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

// ErrClientErr is returned when user authentication fails due to invalid credentials.
// It indicates that the provided username or password is incorrect.
var ErrClientErr = errors.New("invalid username or password")

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
		case errors.Is(err, db.ErrNoRows):
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

// GetUserByID retrieves a user's information by their unique identifier.
// It delegates to the repository layer to fetch the user from the database.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - id: UUID string identifying the user
//
// Returns:
//   - *model.User: User information including username, display name, and email (password is excluded)
//   - error: Returns ErrNotFoundType if user doesn't exist, or an error if database query fails
func (u *user) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return u.repo.GetUserByID(ctx, id)
}

// UpdateUserProfile updates a user's display name and email.
// It delegates the update operation to the repository layer.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - id: UUID string identifying the user to update
//   - displayName: New display name
//   - email: New email address
//
// Returns:
//   - *model.User: Updated user information
//   - error: Returns ErrNotFoundType if user doesn't exist, or an error if database update fails
func (u *user) UpdateUserProfile(ctx context.Context, id, displayName, email string) (*model.User, error) {
	return u.repo.UpdateUserProfile(ctx, id, displayName, email)
}
