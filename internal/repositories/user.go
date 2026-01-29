// Package repository provides data access layer implementations for user-related database operations.
// It handles CRUD operations for user entities using GORM.
package repository

import (
	"context"
	"fmt"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	"github.com/luongtruong20201/bookmark-management/pkg/dbutils"
	"gorm.io/gorm"
)

//go:generate mockery --name User --filename user.go

// User defines the interface for user repository operations.
// It provides methods to interact with user data in the database, including creation,
// retrieval by various fields, and query operations.
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

// CreateUser persists a new user record in the database.
// The user ID will be automatically generated as a UUID if not provided (via BeforeCreate hook).
// This method handles database errors and translates them to application-specific error types.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - user: User model instance to be created (ID will be auto-generated if empty)
//
// Returns:
//   - *model.User: The created user with generated UUID and all persisted fields
//   - error: Returns ErrDuplicationType if username/email already exists, or other database errors
func (u *user) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	if err := u.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, dbutils.CatchDBErr(err)
	}

	return user, nil
}

// GetUserByUsername retrieves a user from the database by their unique username.
// This is a convenience method that delegates to GetUserByField with the "username" field.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - username: Username string to search for
//
// Returns:
//   - *model.User: User information if found
//   - error: Returns ErrNotFoundType if user doesn't exist, or other database errors
func (u *user) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return u.GetUserByField(ctx, "username", username)
}

// GetUserByID retrieves a user from the database by their unique identifier (UUID).
// This is a convenience method that delegates to GetUserByField with the "id" field.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - id: UUID string identifying the user
//
// Returns:
//   - *model.User: User information if found
//   - error: Returns ErrNotFoundType if user doesn't exist, or other database errors
func (u *user) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return u.GetUserByField(ctx, "id", id)
}

// GetUserByField retrieves a user from the database by a specified field name and value.
// This is a generic method used internally by other retrieval methods (GetUserByUsername, GetUserByID).
// It performs a WHERE clause query on the specified field and returns the first matching user.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - field: Database column name to search (e.g., "username", "id", "email")
//   - value: Value to match against the specified field
//
// Returns:
//   - *model.User: User information if found
//   - error: Returns ErrNotFoundType if no user matches the criteria, or other database errors
//
// Note: This method uses parameterized queries to prevent SQL injection. The field name should
// be a trusted value (typically a constant) and not user input.
func (u *user) GetUserByField(ctx context.Context, field string, value string) (*model.User, error) {
	user := &model.User{}
	if err := u.db.WithContext(ctx).Where(fmt.Sprintf("%s = ?", field), value).First(user).Error; err != nil {
		return nil, dbutils.CatchDBErr(err)
	}
	return user, nil
}

// UpdateUserProfile updates the display name and email of a user identified by their ID.
// It returns the updated user or an error if the user does not exist or the update fails.
func (u *user) UpdateUserProfile(ctx context.Context, id, displayName, email string) (*model.User, error) {
	updates := map[string]interface{}{
		"display_name": displayName,
		"email":        email,
	}

	tx := u.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(updates)
	if tx.Error != nil {
		return nil, dbutils.CatchDBErr(tx.Error)
	}

	if tx.RowsAffected == 0 {
		return nil, dbutils.ErrNotFoundType
	}

	return u.GetUserByID(ctx, id)
}
