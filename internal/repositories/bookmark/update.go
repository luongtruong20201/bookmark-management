package bookmark

import (
	"context"
	"errors"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	"github.com/luongtruong20201/bookmark-management/pkg/dbutils"
	"gorm.io/gorm"
)

// UpdateBookmark updates an existing bookmark record in the database.
// It verifies that the bookmark belongs to the specified user before updating.
// Returns an error if the bookmark is not found or doesn't belong to the user.
func (r *repository) UpdateBookmark(ctx context.Context, bookmarkID, userID string, updates *model.Bookmark) (*model.Bookmark, error) {
	var bookmark model.Bookmark

	if err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", bookmarkID, userID).
		First(&bookmark).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dbutils.ErrNotFoundType
		}
		return nil, dbutils.CatchDBErr(err)
	}

	updates.ID = bookmarkID
	updates.UserID = userID

	if err := r.db.WithContext(ctx).
		Model(&bookmark).
		Updates(updates).Error; err != nil {
		return nil, dbutils.CatchDBErr(err)
	}

	return &bookmark, nil
}

