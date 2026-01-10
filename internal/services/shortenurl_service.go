package service

import (
	"context"
	"errors"

	repository "github.com/luongtruong20201/bookmark-management/internal/repositories"
	"github.com/luongtruong20201/bookmark-management/pkg/stringutils"
	"github.com/redis/go-redis/v9"
)

const (
	// codeLength is the length of the generated short code for URLs.
	codeLength = 7
)

var (
	// ErrDuplicatedKey is returned when a generated code already exists in storage.
	ErrDuplicatedKey = errors.New("duplicate key")
	ErrCodeNotFound  = errors.New("code not found")
)

// ShortenURL defines the interface for shorten URL services.
// It provides methods to generate short codes for URLs and store them.
//
//go:generate mockery --name ShortenURL --filename shorten_url.go
type ShortenURL interface {
	ShortenURL(context.Context, string, int) (string, error)
	// GetURL retrieves the original URL associated with the given short code.
	// It returns the original URL if found, or ErrCodeNotFound if the code does not exist.
	GetURL(context.Context, string) (string, error)
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

// GetURL retrieves the original URL associated with the given short code from the repository.
// It returns the original URL if found, or ErrCodeNotFound if the code does not exist in storage.
// Any other error from the repository is returned as-is.
func (s *shortenURL) GetURL(ctx context.Context, code string) (string, error) {
	url, err := s.repository.Get(ctx, code)
	if errors.Is(err, redis.Nil) {
		return "", ErrCodeNotFound
	}

	return url, err
}
