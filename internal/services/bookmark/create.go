package bookmark

import (
	"context"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
)

const codeLength = 8

// Create generates a new short code for the given URL and persists the bookmark.
// It returns the created bookmark with its generated code and database identifier.
func (s bookmarkSvc) Create(ctx context.Context, description, url, userId string) (*model.Bookmark, error) {
	code, err := s.keyGen.GenerateCode(codeLength)
	if err != nil {
		return nil, err
	}

	bookmark := &model.Bookmark{
		Description: description,
		URL:         url,
		Code:        code,
		UserID:      userId,
	}

	bookmark, err = s.repository.CreateBookmark(ctx, bookmark)
	if err != nil {
		return nil, err
	}

	return bookmark, nil
}
