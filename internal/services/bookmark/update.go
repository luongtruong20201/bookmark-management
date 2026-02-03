package bookmark

import (
	"context"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
)

// Update updates an existing bookmark for a specific user.
// It delegates to the repository layer to update the bookmark.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - bookmarkID: The unique identifier of the bookmark to update
//   - userID: The unique identifier of the user who owns the bookmark
//   - description: Optional new description for the bookmark
//   - url: Optional new URL for the bookmark
//
// Returns:
//   - *model.Bookmark: The updated bookmark, or nil if an error occurs
//   - error: An error if the repository operation fails or the bookmark doesn't belong to the user
func (s bookmarkSvc) Update(ctx context.Context, bookmarkID, userID, description, url string) (*model.Bookmark, error) {
	updates := &model.Bookmark{}

	if description != "" {
		updates.Description = description
	}

	if url != "" {
		updates.URL = url
	}

	bookmark, err := s.repository.UpdateBookmark(ctx, bookmarkID, userID, updates)
	if err != nil {
		return nil, err
	}

	return bookmark, nil
}

