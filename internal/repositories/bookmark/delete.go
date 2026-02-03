package bookmark

import (
	"context"
	"errors"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	"github.com/luongtruong20201/bookmark-management/pkg/dbutils"
	"gorm.io/gorm"
)

// DeleteBookmark deletes a bookmark record from the database.
// It verifies that the bookmark belongs to the specified user before deleting.
// Returns an error if the bookmark is not found or doesn't belong to the user.
func (r *repository) DeleteBookmark(ctx context.Context, bookmarkID, userID string) error {
	var bookmark model.Bookmark

	if err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", bookmarkID, userID).
		First(&bookmark).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dbutils.ErrNotFoundType
		}
		return dbutils.CatchDBErr(err)
	}

	if err := r.db.WithContext(ctx).
		Delete(&bookmark).Error; err != nil {
		return dbutils.CatchDBErr(err)
	}

	return nil
}

