package user

import (
	"context"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
)

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
