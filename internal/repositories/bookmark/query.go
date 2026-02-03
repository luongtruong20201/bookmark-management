package bookmark

import (
	"context"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
)

// GetBookmarks retrieves bookmarks for a specific user with pagination support.
// It queries the database for bookmarks filtered by userID, ordered by creation date (ascending),
// and applies offset and limit for pagination.
//
// Parameters:
//   - ctx: Context for database operation cancellation and timeout
//   - userID: The unique identifier of the user whose bookmarks to retrieve
//   - offset: The number of records to skip (for pagination)
//   - limit: The maximum number of records to return
//
// Returns:
//   - []*model.Bookmark: A slice of bookmarks for the user, or nil if an error occurs
//   - error: A database error if the query fails
func (r *repository) GetBookmarks(ctx context.Context, userID string, offset, limit int) ([]*model.Bookmark, error) {
	bookmarks := make([]*model.Bookmark, 0)
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at ASC").
		Offset(offset).
		Limit(limit).
		Find(&bookmarks).Error; err != nil {
		return nil, err
	}

	return bookmarks, nil
}

// CountBookmarks counts the total number of bookmarks for a specific user.
// It queries the database to get the total count of bookmarks filtered by userID.
//
// Parameters:
//   - ctx: Context for database operation cancellation and timeout
//   - userID: The unique identifier of the user whose bookmarks to count
//
// Returns:
//   - int64: The total number of bookmarks for the user, or 0 if an error occurs
//   - error: A database error if the count query fails
func (r *repository) CountBookmarks(ctx context.Context, userID string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.Bookmark{}).
		Where("user_id = ?", userID).
		Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
