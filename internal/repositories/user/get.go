package user

import (
	"context"
	"fmt"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	"github.com/luongtruong20201/bookmark-management/pkg/dbutils"
)

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
