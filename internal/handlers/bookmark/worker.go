package bookmark

import (
	"context"
	"encoding/json"

	"github.com/luongtruong20201/bookmark-management/internal/services/bookmark"
	"github.com/luongtruong20201/bookmark-management/internal/services/queue"
	"github.com/luongtruong20201/bookmark-management/internal/worker"
)

// workerHandler implements worker.Handler for bookmark import jobs. It
// unmarshals queue.ImportMessage (uid + bookmarks) and calls the bookmark
// service's ImportBookmarks to batch-insert records.
type workerHandler struct {
	bookmarkSvc bookmark.Service
}

// NewWorkerHandler returns a worker.Handler that processes bookmark import
// messages using the given bookmark service (e.g. for batch insert).
func NewWorkerHandler(bookmarkSvc bookmark.Service) worker.Handler {
	return &workerHandler{bookmarkSvc: bookmarkSvc}
}

// Handle decodes the message as JSON ImportMessage, maps bookmarks to
// ImportBookmarkItem, and calls ImportBookmarks on the bookmark service.
func (h *workerHandler) Handle(ctx context.Context, message []byte) error {
	var msg queue.ImportMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		return err
	}
	items := make([]*bookmark.ImportBookmarkItem, 0, len(msg.Bookmarks))
	for _, b := range msg.Bookmarks {
		items = append(items, &bookmark.ImportBookmarkItem{
			Description: b.Description,
			URL:         b.URL,
		})
	}
	return h.bookmarkSvc.ImportBookmarks(ctx, msg.UID, items)
}
