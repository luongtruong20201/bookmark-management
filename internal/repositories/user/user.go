package user

import (
	"context"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	"gorm.io/gorm"
)

// User defines the interface for user repository operations.
// It provides methods to interact with user data in the database, including creation,
// retrieval by various fields, and query operations.
//
//go:generate mockery --name User --filename user.go
type User interface {
	// CreateUser persists a new user record in the database.
	// The user ID will be automatically generated as a UUID if not provided.
	CreateUser(context.Context, *model.User) (*model.User, error)

	// GetUserByUsername retrieves a user by their unique username.
	// Returns the user or an error if not found.
	GetUserByUsername(context.Context, string) (*model.User, error)

	// GetUserByID retrieves a user by their unique identifier (UUID).
	// Returns the user or an error if not found.
	GetUserByID(ctx context.Context, id string) (*model.User, error)

	// GetUserByField retrieves a user by a specified field name and value.
	// This is a generic method used by other retrieval methods.
	// Returns the user or an error if not found.
	GetUserByField(ctx context.Context, field string, value string) (*model.User, error)

	// UpdateUserProfile updates the display name and email of a user identified by ID.
	// Returns the updated user or an error if the user is not found or the update fails.
	UpdateUserProfile(ctx context.Context, id, displayName, email string) (*model.User, error)
}

// user implements the User interface and provides database operations for user entities.
// It uses GORM for database interactions and handles error translation.
type user struct {
	db *gorm.DB
}

// NewUser creates a new user repository instance with the provided database connection.
// It initializes the repository with a GORM database instance for performing user-related queries.
//
// Parameters:
//   - db: GORM database instance for database operations
//
// Returns:
//   - User: A new user repository instance implementing the User interface
func NewUser(db *gorm.DB) User {
	return &user{
		db: db,
	}
}
