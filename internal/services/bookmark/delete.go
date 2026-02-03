package bookmark

import (
	"context"
)

// Delete deletes a bookmark for a specific user.
// It delegates to the repository layer to delete the bookmark.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - bookmarkID: The unique identifier of the bookmark to delete
//   - userID: The unique identifier of the user who owns the bookmark
//
// Returns:
//   - error: An error if the repository operation fails or the bookmark doesn't belong to the user
func (s bookmarkSvc) Delete(ctx context.Context, bookmarkID, userID string) error {
	err := s.repository.DeleteBookmark(ctx, bookmarkID, userID)
	if err != nil {
		return err
	}

	return nil
}

