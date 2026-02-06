package shorten

import (
	"context"
	"errors"

	"github.com/luongtruong20201/bookmark-management/pkg/dbutils"
	"github.com/redis/go-redis/v9"
)

// GetURL retrieves the original URL associated with the given short code.
// It checks the code length to determine the source:
// - Code length 7: retrieves from Redis (shortened URLs)
// - Code length 8: retrieves from database (bookmarks)
// It returns the original URL if found, or ErrCodeNotFound if the code does not exist.
// Any other error from the repository is returned as-is.
func (s *shortenURL) GetURL(ctx context.Context, code string) (string, error) {
	codeLen := len(code)

	if codeLen == urlCodeLength {
		url, err := s.repository.Get(ctx, code)
		if errors.Is(err, redis.Nil) {
			return "", ErrCodeNotFound
		}
		return url, err
	}

	if codeLen == bookmarkCodeLength {
		bookmark, err := s.bookmarkRepo.GetBookmarkByCode(ctx, code)
		if errors.Is(err, dbutils.ErrNotFoundType) {
			return "", ErrCodeNotFound
		}
		if err != nil {
			return "", err
		}
		return bookmark.URL, nil
	}

	return "", ErrCodeNotFound
}
