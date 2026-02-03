package user

import (
	"context"
	"errors"
	"time"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	repository "github.com/luongtruong20201/bookmark-management/internal/repositories/user"
	jwtPkg "github.com/luongtruong20201/bookmark-management/pkg/jwt"
	"github.com/luongtruong20201/bookmark-management/pkg/utils"
)

const (
	// tokenExpiresTime defines the expiration duration for JWT tokens.
	// Tokens are valid for 24 hours from the time of generation.
	tokenExpiresTime = 24 * time.Hour
)

// ErrClientErr is returned when user authentication fails due to invalid credentials.
// It indicates that the provided username or password is incorrect.
var ErrClientErr = errors.New("invalid username or password")

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
