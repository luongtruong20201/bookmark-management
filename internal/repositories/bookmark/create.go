package bookmark

import (
	"context"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	"github.com/luongtruong20201/bookmark-management/pkg/dbutils"
)

const (
	batchSize = 100
)

// CreateBookmark persists a new bookmark record into the database.
// It wraps GORM errors using dbutils.CatchDBErr so callers receive
// normalized error types (e.g. duplicate key, not found, etc).
func (b *repository) CreateBookmark(ctx context.Context, bookmark *model.Bookmark) (*model.Bookmark, error) {
	if err := b.db.WithContext(ctx).Create(bookmark).Error; err != nil {
		return nil, dbutils.CatchDBErr(err)
	}

	return bookmark, nil
}

// CreateBookmarksBatch inserts multiple bookmarks in one batch using GORM Create.
// If bookmarks is nil or empty, it returns nil without calling the DB.
func (b *repository) CreateBookmarksBatch(ctx context.Context, bookmarks []*model.Bookmark) error {
	if len(bookmarks) == 0 {
		return nil
	}
	if err := b.db.WithContext(ctx).CreateInBatches(bookmarks, batchSize).Error; err != nil {
		return dbutils.CatchDBErr(err)
	}
	return nil
}
