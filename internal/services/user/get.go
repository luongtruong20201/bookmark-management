package user

import (
	"context"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
)

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
