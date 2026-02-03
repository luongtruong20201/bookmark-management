package user

import (
	"context"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	"github.com/luongtruong20201/bookmark-management/pkg/dbutils"
)

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
