package bookmark

import (
	"context"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	bookmarkRepo "github.com/luongtruong20201/bookmark-management/internal/repositories/bookmark"
	"github.com/luongtruong20201/bookmark-management/pkg/stringutils"
)

// Service defines the interface for bookmark business operations.
// It is responsible for generating bookmark codes and delegating
// persistence to the repository layer.
//
//go:generate mockery --name Service --filename bookmark.go
type Service interface {
	Create(ctx context.Context, description, url, userId string) (*model.Bookmark, error)
	GetBookmarks(ctx context.Context, userID string, offset, limit int) (*GetBookmarksResponse, error)
	CountBookmarks(ctx context.Context, userID string) (int64, error)
	Update(ctx context.Context, bookmarkID, userID, description, url string) (*model.Bookmark, error)
	Delete(ctx context.Context, bookmarkID, userID string) error
}

// bookmarkSvc is the concrete implementation of the Service interface.
// It composes a bookmark repository and a key generator used to create
// short, unique codes for each bookmark.
type bookmarkSvc struct {
	repository bookmarkRepo.Repository
	keyGen     stringutils.KeyGenerator
}

// NewBookmarkSvc constructs a new bookmark service with the provided
// repository and key generator dependencies.
func NewBookmarkSvc(repo bookmarkRepo.Repository, keyGen stringutils.KeyGenerator) Service {
	return &bookmarkSvc{
		repository: repo,
		keyGen:     keyGen,
	}
}
