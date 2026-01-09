package service

import (
	"context"
	"errors"

	repository "github.com/luongtruong20201/bookmark-management/internal/repositories"
	"github.com/luongtruong20201/bookmark-management/pkg/stringutils"
)

const (
	// codeLength is the length of the generated short code for URLs.
	codeLength = 7
)

var (
	// ErrDuplicatedKey is returned when a generated code already exists in storage.
	ErrDuplicatedKey = errors.New("duplicate key")
)

// ShortenURL defines the interface for shorten URL services.
// It provides methods to generate short codes for URLs and store them.
//
//go:generate mockery --name ShortenURL --filename shorten_url.go
type ShortenURL interface {
	ShortenURL(context.Context, string, int) (string, error)
}

type shortenURL struct {
	keyGen     stringutils.KeyGenerator
	repository repository.URLStorage
}

// NewShortenURL creates a new shorten URL service instance with the provided
// key generator and URL storage repository.
func NewShortenURL(keyGen stringutils.KeyGenerator, repository repository.URLStorage) ShortenURL {
	return &shortenURL{
		keyGen:     keyGen,
		repository: repository,
	}
}

// ShortenURL generates a short code for the given URL and stores it in the repository.
// It returns the generated code or an error if code generation or storage fails.
// The expire parameter specifies the expiration time in seconds (0 means default expiration).
func (s *shortenURL) ShortenURL(ctx context.Context, url string, expire int) (string, error) {
	code, err := s.keyGen.GenerateCode(codeLength)
	if err != nil {
		return "", err
	}

	ok, err := s.repository.StoreIfNotExists(ctx, code, url, expire)

	switch {
	case err != nil:
		return "", err
	case !ok:
		return "", ErrDuplicatedKey
	}

	return code, nil
}
