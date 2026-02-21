package shorten

import (
	"context"
	"errors"
	"testing"

	mockBookmarkRepo "github.com/luongtruong20201/bookmark-management/internal/repositories/bookmark/mocks"
	mockStorage "github.com/luongtruong20201/bookmark-management/internal/repositories/url/mocks"
	mockKeyGen "github.com/luongtruong20201/bookmark-management/pkg/stringutils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestShortenURL_ShortenURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		setupRepo     func(t *testing.T, ctx context.Context, url string, exp int, code string) *mockStorage.URLStorage
		setupKeyGen   func(t *testing.T) *mockKeyGen.KeyGenerator
		url           string
		exp           int
		storeCode     string
		expectedCode  string
		expectedError error
	}{
		{
			name: "success",
			setupRepo: func(t *testing.T, ctx context.Context, url string, exp int, code string) *mockStorage.URLStorage {
				repo := mockStorage.NewURLStorage(t)
				repo.On("StoreIfNotExists", ctx, code, url, exp).Return(true, nil).Once()
				return repo
			},
			setupKeyGen: func(t *testing.T) *mockKeyGen.KeyGenerator {
				keyGen := mockKeyGen.NewKeyGenerator(t)
				keyGen.On("GenerateCode", urlCodeLength).Return("1234567", nil).Once()
				return keyGen
			},
			url:           "https://truonglq.com",
			exp:           0,
			storeCode:     "1234567",
			expectedCode:  "1234567",
			expectedError: nil,
		},
		{
			name: "duplicate",
			setupRepo: func(t *testing.T, ctx context.Context, url string, exp int, code string) *mockStorage.URLStorage {
				repo := mockStorage.NewURLStorage(t)
				repo.On("StoreIfNotExists", ctx, code, url, exp).Return(false, ErrDuplicatedKey).Once()
				return repo
			},
			setupKeyGen: func(t *testing.T) *mockKeyGen.KeyGenerator {
				keyGen := mockKeyGen.NewKeyGenerator(t)
				keyGen.On("GenerateCode", urlCodeLength).Return("1234567", nil).Once()
				return keyGen
			},
			url:           "https://truonglq.com",
			exp:           0,
			storeCode:     "1234567",
			expectedCode:  "",
			expectedError: ErrDuplicatedKey,
		},
		{
			name: "key gen error",
			setupRepo: func(t *testing.T, ctx context.Context, url string, exp int, code string) *mockStorage.URLStorage {
				return mockStorage.NewURLStorage(t)
			},
			setupKeyGen: func(t *testing.T) *mockKeyGen.KeyGenerator {
				keyGen := mockKeyGen.NewKeyGenerator(t)
				keyGen.On("GenerateCode", urlCodeLength).Return("", errors.New("error")).Once()
				return keyGen
			},
			url:           "https://truonglq.com",
			exp:           0,
			storeCode:     "",
			expectedCode:  "",
			expectedError: errors.New("error"),
		},
		{
			name: "repository storage error",
			setupRepo: func(t *testing.T, ctx context.Context, url string, exp int, code string) *mockStorage.URLStorage {
				repo := mockStorage.NewURLStorage(t)
				repo.On("StoreIfNotExists", ctx, code, url, exp).Return(false, errors.New("redis connection failed")).Once()
				return repo
			},
			setupKeyGen: func(t *testing.T) *mockKeyGen.KeyGenerator {
				keyGen := mockKeyGen.NewKeyGenerator(t)
				keyGen.On("GenerateCode", urlCodeLength).Return("1234567", nil).Once()
				return keyGen
			},
			url:           "https://truonglq.com",
			exp:           0,
			storeCode:     "1234567",
			expectedCode:  "",
			expectedError: errors.New("redis connection failed"),
		},
		{
			name: "success with custom expiration",
			setupRepo: func(t *testing.T, ctx context.Context, url string, exp int, code string) *mockStorage.URLStorage {
				repo := mockStorage.NewURLStorage(t)
				repo.On("StoreIfNotExists", ctx, code, url, exp).Return(true, nil).Once()
				return repo
			},
			setupKeyGen: func(t *testing.T) *mockKeyGen.KeyGenerator {
				keyGen := mockKeyGen.NewKeyGenerator(t)
				keyGen.On("GenerateCode", urlCodeLength).Return("abcdefg", nil).Once()
				return keyGen
			},
			url:           "https://example.com",
			exp:           3600,
			storeCode:     "abcdefg",
			expectedCode:  "abcdefg",
			expectedError: nil,
		},
		{
			name: "success with zero expiration",
			setupRepo: func(t *testing.T, ctx context.Context, url string, exp int, code string) *mockStorage.URLStorage {
				repo := mockStorage.NewURLStorage(t)
				repo.On("StoreIfNotExists", ctx, code, url, exp).Return(true, nil).Once()
				return repo
			},
			setupKeyGen: func(t *testing.T) *mockKeyGen.KeyGenerator {
				keyGen := mockKeyGen.NewKeyGenerator(t)
				keyGen.On("GenerateCode", urlCodeLength).Return("xyz1234", nil).Once()
				return keyGen
			},
			url:           "https://test.com",
			exp:           0,
			storeCode:     "xyz1234",
			expectedCode:  "xyz1234",
			expectedError: nil,
		},
		{
			name: "success with maximum expiration",
			setupRepo: func(t *testing.T, ctx context.Context, url string, exp int, code string) *mockStorage.URLStorage {
				repo := mockStorage.NewURLStorage(t)
				repo.On("StoreIfNotExists", ctx, code, url, exp).Return(true, nil).Once()
				return repo
			},
			setupKeyGen: func(t *testing.T) *mockKeyGen.KeyGenerator {
				keyGen := mockKeyGen.NewKeyGenerator(t)
				keyGen.On("GenerateCode", urlCodeLength).Return("maxexp01", nil).Once()
				return keyGen
			},
			url:           "https://longurl.com",
			exp:           604800,
			storeCode:     "maxexp01",
			expectedCode:  "maxexp01",
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			repo := tc.setupRepo(t, ctx, tc.url, tc.exp, tc.storeCode)
			keyGen := tc.setupKeyGen(t)
			bookmarkRepo := mockBookmarkRepo.NewRepository(t)
			svc := NewShortenURL(keyGen, repo, bookmarkRepo)

			code, err := svc.ShortenURL(ctx, tc.url, tc.exp)

			assert.Equal(t, err, tc.expectedError)
			assert.Equal(t, tc.expectedCode, code)
		})
	}
}
