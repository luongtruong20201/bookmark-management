package bookmark

import (
	"context"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	"gorm.io/gorm"
)

// Repository defines the interface for bookmark persistence operations.
// It provides methods to create and work with bookmark records in the database.
//
//go:generate mockery --name Repository --filename bookmark.go
type Repository interface {
	CreateBookmark(ctx context.Context, bookmark *model.Bookmark) (*model.Bookmark, error)
	GetBookmarks(ctx context.Context, userID string, offset, limit int) ([]*model.Bookmark, error)
	CountBookmarks(ctx context.Context, userID string) (int64, error)
	UpdateBookmark(ctx context.Context, bookmarkID, userID string, updates *model.Bookmark) (*model.Bookmark, error)
	DeleteBookmark(ctx context.Context, bookmarkID, userID string) error
}

// repository is the concrete implementation of the Repository interface.
// It wraps a GORM database connection and translates low‑level errors
// into domain friendly errors in the corresponding methods.
type repository struct {
	db *gorm.DB
}

// NewBookmark creates a new bookmark repository using the given database connection.
// The returned Repository can be used by services to persist bookmark entities.
func NewBookmark(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}
