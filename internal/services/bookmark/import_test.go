package bookmark

import (
	"errors"
	"testing"

	repoMocks "github.com/luongtruong20201/bookmark-management/internal/repositories/bookmark/mocks"
	mockKeyGen "github.com/luongtruong20201/bookmark-management/pkg/stringutils/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_ImportBookmarks(t *testing.T) {
	t.Parallel()

	var (
		errKeyGen = errors.New("keygen error")
		errRepo   = errors.New("repository error")
	)

	testCases := []struct {
		name        string
		userID      string
		items       []*ImportBookmarkItem
		setupKeyGen func(t *testing.T, n int) *mockKeyGen.KeyGenerator
		setupRepo   func(t *testing.T) *repoMocks.Repository
		expectedErr error
	}{
		{
			name:        "empty items returns nil",
			userID:      "user-1",
			items:       nil,
			setupKeyGen: func(t *testing.T, n int) *mockKeyGen.KeyGenerator { return mockKeyGen.NewKeyGenerator(t) },
			setupRepo:   func(t *testing.T) *repoMocks.Repository { return repoMocks.NewRepository(t) },
			expectedErr: nil,
		},
		{
			name:        "empty slice returns nil",
			userID:      "user-1",
			items:       []*ImportBookmarkItem{},
			setupKeyGen: func(t *testing.T, n int) *mockKeyGen.KeyGenerator { return mockKeyGen.NewKeyGenerator(t) },
			setupRepo:   func(t *testing.T) *repoMocks.Repository { return repoMocks.NewRepository(t) },
			expectedErr: nil,
		},
		{
			name:   "success_import_calls_CreateBookmarksBatch",
			userID: "550e8400-e29b-41d4-a716-446655440000",
			items: []*ImportBookmarkItem{
				{Description: "D1", URL: "https://a.com"},
				{Description: "D2", URL: "https://b.com"},
			},
			setupKeyGen: func(t *testing.T, n int) *mockKeyGen.KeyGenerator {
				keyGen := mockKeyGen.NewKeyGenerator(t)
				keyGen.On("GenerateCode", codeLength).Return("code001", nil).Once()
				keyGen.On("GenerateCode", codeLength).Return("code002", nil).Once()
				return keyGen
			},
			setupRepo: func(t *testing.T) *repoMocks.Repository {
				repo := repoMocks.NewRepository(t)
				repo.On("CreateBookmarksBatch", mock.Anything, mock.Anything).Return(nil).Once()
				return repo
			},
			expectedErr: nil,
		},
		{
			name:   "keyGen_error_returns_error",
			userID: "user-1",
			items:  []*ImportBookmarkItem{{Description: "D1", URL: "https://a.com"}},
			setupKeyGen: func(t *testing.T, n int) *mockKeyGen.KeyGenerator {
				keyGen := mockKeyGen.NewKeyGenerator(t)
				keyGen.On("GenerateCode", codeLength).Return("", errKeyGen).Once()
				return keyGen
			},
			setupRepo:   func(t *testing.T) *repoMocks.Repository { return repoMocks.NewRepository(t) },
			expectedErr: errKeyGen,
		},
		{
			name:   "CreateBookmarksBatch_error_returns_error",
			userID: "user-1",
			items:  []*ImportBookmarkItem{{Description: "D1", URL: "https://a.com"}},
			setupKeyGen: func(t *testing.T, n int) *mockKeyGen.KeyGenerator {
				keyGen := mockKeyGen.NewKeyGenerator(t)
				keyGen.On("GenerateCode", codeLength).Return("code1", nil).Once()
				return keyGen
			},
			setupRepo: func(t *testing.T) *repoMocks.Repository {
				repo := repoMocks.NewRepository(t)
				repo.On("CreateBookmarksBatch", mock.Anything, mock.Anything).Return(errRepo).Once()
				return repo
			},
			expectedErr: errRepo,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			keyGen := tc.setupKeyGen(t, len(tc.items))
			repo := tc.setupRepo(t)
			svc := NewBookmarkSvc(repo, keyGen)

			err := svc.ImportBookmarks(ctx, tc.userID, tc.items)

			if tc.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
				return
			}
			require.NoError(t, err)
		})
	}
}
