package bookmark

import (
	"context"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	"github.com/luongtruong20201/bookmark-management/pkg/dbutils"
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
