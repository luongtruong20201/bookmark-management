package bookmark

import (
	"context"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
)

// GetBookmarksResponse represents the response structure for GetBookmarks service method.
type GetBookmarksResponse struct {
	Data  []*model.Bookmark `json:"data"`
	Total int64             `json:"total"`
}

// GetBookmarks retrieves bookmarks for a specific user with pagination support
// and returns the total count of bookmarks for that user.
// It delegates to the repository layer to fetch bookmarks ordered by creation date (ascending)
// and to count the total number of bookmarks.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - userID: The unique identifier of the user whose bookmarks to retrieve
//   - offset: The number of records to skip (for pagination)
//   - limit: The maximum number of records to return
//
// Returns:
//   - *GetBookmarksResponse: A response containing bookmarks and total count, or nil if an error occurs
//   - error: An error if the repository operation fails
func (s bookmarkSvc) GetBookmarks(ctx context.Context, userID string, offset, limit int) (*GetBookmarksResponse, error) {
	bookmarks, err := s.repository.GetBookmarks(ctx, userID, offset, limit)
	if err != nil {
		return nil, err
	}

	total, err := s.repository.CountBookmarks(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &GetBookmarksResponse{
		Data:  bookmarks,
		Total: total,
	}, nil
}

// CountBookmarks counts the total number of bookmarks for a specific user.
// It delegates to the repository layer to get the total count of bookmarks.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - userID: The unique identifier of the user whose bookmarks to count
//
// Returns:
//   - int64: The total number of bookmarks for the user, or 0 if an error occurs
//   - error: An error if the repository operation fails
func (s bookmarkSvc) CountBookmarks(ctx context.Context, userID string) (int64, error) {
	total, err := s.repository.CountBookmarks(ctx, userID)
	if err != nil {
		return 0, err
	}

	return total, nil
}

