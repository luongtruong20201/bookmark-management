package shorten

import (
	"context"
	"errors"
	"testing"

	mockStorage "github.com/luongtruong20201/bookmark-management/internal/repositories/url/mocks"
	mockKeyGen "github.com/luongtruong20201/bookmark-management/pkg/stringutils/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShortenURL_ShortenURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		setupRepo     func(t *testing.T, ctx context.Context, url string, exp int) *mockStorage.URLStorage
		setupKeyGen   func(t *testing.T) *mockKeyGen.KeyGenerator
		url           string
		exp           int
		expectedCode  string
		expectedError error
	}{
		{
			name: "success",
			setupRepo: func(t *testing.T, ctx context.Context, url string, exp int) *mockStorage.URLStorage {
				repo := mockStorage.NewURLStorage(t)
				repo.On("StoreIfNotExists", ctx, mock.Anything, url, exp).Return(true, nil).Once()

				return repo
			},
			setupKeyGen: func(t *testing.T) *mockKeyGen.KeyGenerator {
				keyGen := mockKeyGen.NewKeyGenerator(t)
				keyGen.On("GenerateCode", codeLength).Return("1234567", nil).Once()

				return keyGen
			},
			url:           "https://truonglq.com",
			exp:           0,
			expectedCode:  "1234567",
			expectedError: nil,
		},
		{
			name: "duplicate",
			setupRepo: func(t *testing.T, ctx context.Context, url string, exp int) *mockStorage.URLStorage {
				repo := mockStorage.NewURLStorage(t)
				repo.On("StoreIfNotExists", ctx, mock.Anything, url, exp).Return(false, ErrDuplicatedKey).Once()

				return repo
			},
			setupKeyGen: func(t *testing.T) *mockKeyGen.KeyGenerator {
				keyGen := mockKeyGen.NewKeyGenerator(t)
				keyGen.On("GenerateCode", codeLength).Return("1234567", nil).Once()

				return keyGen
			},
			url:           "https://truonglq.com",
			exp:           0,
			expectedCode:  "",
			expectedError: ErrDuplicatedKey,
		},
		{
			name: "key gen error",
			setupRepo: func(t *testing.T, ctx context.Context, url string, exp int) *mockStorage.URLStorage {
				repo := mockStorage.NewURLStorage(t)

				return repo
			},
			setupKeyGen: func(t *testing.T) *mockKeyGen.KeyGenerator {
				keyGen := mockKeyGen.NewKeyGenerator(t)
				keyGen.On("GenerateCode", codeLength).Return("", errors.New("error")).Once()

				return keyGen
			},
			url:           "https://truonglq.com",
			exp:           0,
			expectedCode:  "",
			expectedError: errors.New("error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			repo := tc.setupRepo(t, ctx, tc.url, tc.exp)
			keyGen := tc.setupKeyGen(t)
			svc := NewShortenURL(keyGen, repo)

			code, err := svc.ShortenURL(ctx, tc.url, tc.exp)

			assert.Equal(t, err, tc.expectedError)
			assert.Equal(t, tc.expectedCode, code)
		})
	}
}
