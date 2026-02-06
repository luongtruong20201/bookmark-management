package bookmark

import (
	"context"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
)

// GetBookmarkByCode retrieves a bookmark by its code from the repository.
// It delegates to the repository layer to fetch the bookmark.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - code: The short code of the bookmark to retrieve
//
// Returns:
//   - *model.Bookmark: The bookmark with the given code, or nil if not found
//   - error: An error if the repository operation fails
func (s bookmarkSvc) GetBookmarkByCode(ctx context.Context, code string) (*model.Bookmark, error) {
	return s.repository.GetBookmarkByCode(ctx, code)
}

