package shorten

import (
	"context"
	"errors"

	bookmarkRepo "github.com/luongtruong20201/bookmark-management/internal/repositories/bookmark"
	repository "github.com/luongtruong20201/bookmark-management/internal/repositories/url"
	"github.com/luongtruong20201/bookmark-management/pkg/stringutils"
)

const (
	// urlCodeLength is the length of the generated short code for URLs.
	urlCodeLength = 7
	// bookmarkCodeLength is the length of the generated code for bookmarks.
	bookmarkCodeLength = 8
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

// shortenURL implements the ShortenURL interface and provides business logic
// for URL shortening operations. It uses a key generator to create short codes
// and a repository to store and retrieve URL mappings.
type shortenURL struct {
	keyGen       stringutils.KeyGenerator
	repository   repository.URLStorage
	bookmarkRepo bookmarkRepo.Repository
}

// NewShortenURL creates a new shorten URL service instance with the provided
// key generator, URL storage repository, and bookmark repository.
func NewShortenURL(keyGen stringutils.KeyGenerator, repository repository.URLStorage, bookmarkRepo bookmarkRepo.Repository) ShortenURL {
	return &shortenURL{
		keyGen:       keyGen,
		repository:   repository,
		bookmarkRepo: bookmarkRepo,
	}
}
