package bookmark

import (
	"context"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	"github.com/luongtruong20201/bookmark-management/pkg/dbutils"
)

// DeleteBookmark deletes a bookmark record from the database.
// It verifies that the bookmark belongs to the specified user before deleting.
// Returns an error if the bookmark is not found or doesn't belong to the user.
func (r *repository) DeleteBookmark(ctx context.Context, bookmarkID, userID string) error {
	var bookmark model.Bookmark

	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", bookmarkID, userID).First(&bookmark).Error
	if err != nil {
		return dbutils.CatchDBErr(err)
	}

	if err = r.db.WithContext(ctx).Delete(&bookmark).Error; err != nil {
		return dbutils.CatchDBErr(err)
	}

	return nil
}
