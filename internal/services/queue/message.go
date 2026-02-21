package queue

import (
	"context"
	"encoding/json"

	"github.com/luongtruong20201/bookmark-management/pkg/array"
)

const (
	// chunkSize is the maximum number of bookmarks bundled into a single
	// ImportMessage before being pushed to the queue.
	chunkSize = 100
)

// ImportMessage represents the payload pushed to the queue for a bookmark
// import job. It contains the ID of the user who triggered the import and
// the list of bookmarks to be imported.
type ImportMessage struct {
	// UID is the ID of the user who initiated the import.
	UID string `json:"uid"`
	// Bookmarks is the list of bookmark records parsed from the CSV file.
	Bookmarks []*ImportBookmarkInput `json:"bookmarks"`
}

// ImportBookmarkInput describes a single bookmark entry as it appears in the
// imported CSV file. The struct tags are mapped to CSV column names.
type ImportBookmarkInput struct {
	// Description is a human-readable description of the bookmark.
	// It is required when importing from CSV.
	Description string `csv:"description" validate:"required"`
	// URL is the target URL to be stored as a bookmark.
	// It is required and must be a valid URL.
	URL string `csv:"url" validate:"required,url"`
}

// SendBookmarkJob groups the given bookmarks into chunks and enqueues one or
// more ImportMessage jobs for the specified user. Each chunk is delegated to
// sendJob, which performs the actual serialization and push to the queue.
func (s *queueService) SendBookmarkJob(ctx context.Context, uid string, bookmarks []*ImportBookmarkInput) error {
	chunks := array.SplitIntoChunks[*ImportBookmarkInput](bookmarks, chunkSize)
	for _, chunk := range chunks {
		if err := s.sendJob(ctx, uid, chunk); err != nil {
			return err
		}
	}
	return nil
}

// sendJob builds a single ImportMessage from the provided user ID and chunk of
// bookmarks, serializes it to JSON, and pushes it to the underlying queue
// repository.
func (s *queueService) sendJob(ctx context.Context, uid string, bookmark []*ImportBookmarkInput) error {
	msg := &ImportMessage{
		UID:       uid,
		Bookmarks: bookmark,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return s.repo.PushMessage(ctx, msgBytes)

}
