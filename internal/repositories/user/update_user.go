package user

import (
	"context"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	"github.com/luongtruong20201/bookmark-management/pkg/dbutils"
)

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
