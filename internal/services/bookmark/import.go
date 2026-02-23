package bookmark

import (
	"context"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
)

// ImportBookmarkItem is one record for batch import. Used by the worker
// handler and ImportBookmarks to create bookmarks without CSV/queue types.
type ImportBookmarkItem struct {
	Description string
	URL        string
}


// ImportBookmarks creates multiple bookmarks for the given user in one batch.
// For each item it generates a short code, builds a Bookmark model, then calls
// the repository's CreateBookmarksBatch. Returns on first error (keyGen or DB).
func (s bookmarkSvc) ImportBookmarks(ctx context.Context, userID string, items []*ImportBookmarkItem) error {
	if len(items) == 0 {
		return nil
	}
	bookmarks := make([]*model.Bookmark, 0, len(items))
	for _, item := range items {
		code, err := s.keyGen.GenerateCode(codeLength)
		if err != nil {
			return err
		}
		bookmarks = append(bookmarks, &model.Bookmark{
			Description: item.Description,
			URL:         item.URL,
			Code:        code,
			UserID:      userID,
		})
	}
	return s.repository.CreateBookmarksBatch(ctx, bookmarks)
}
