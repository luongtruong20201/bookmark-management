package bookmark

import (
	"context"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	"github.com/luongtruong20201/bookmark-management/pkg/dbutils"
)

// GetBookmarkByCode retrieves a bookmark by its code from the database.
// It returns the bookmark if found, or dbutils.ErrNotFoundType if the code does not exist.
// Database errors are wrapped using dbutils.CatchDBErr for normalized error handling.
func (r *repository) GetBookmarkByCode(ctx context.Context, code string) (*model.Bookmark, error) {
	var bookmark model.Bookmark
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&bookmark).Error; err != nil {
		return nil, dbutils.CatchDBErr(err)
	}

	return &bookmark, nil
}
