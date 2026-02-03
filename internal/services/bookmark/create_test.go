package bookmark

import (
	"context"
	"errors"
	"testing"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	repoMocks "github.com/luongtruong20201/bookmark-management/internal/repositories/bookmark/mocks"
	mockKeyGen "github.com/luongtruong20201/bookmark-management/pkg/stringutils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestBookmarkService_Create(t *testing.T) {
	t.Parallel()

	var (
		testErrKeyGen   = errors.New("keygen error")
		testErrDatabase = errors.New("database error")
	)

	testCases := []struct {
		name           string
		setupRepo      func(t *testing.T, ctx context.Context, description, url, userID, code string) *repoMocks.Repository
		setupKeyGen    func(t *testing.T, code string, err error) *mockKeyGen.KeyGenerator
		description    string
		url            string
		userID         string
		expectedCode   string
		expectedError  error
		verifyBookmark func(t *testing.T, bookmark *model.Bookmark)
	}{
		{
			name:         "success - create bookmark",
			description:  "My blog",
			url:          "https://truonglq.com",
			userID:       "550e8400-e29b-41d4-a716-446655440000",
			expectedCode: "abcd1234",
			setupKeyGen: func(t *testing.T, code string, err error) *mockKeyGen.KeyGenerator {
				keyGen := mockKeyGen.NewKeyGenerator(t)
				keyGen.On("GenerateCode", codeLength).Return(code, err).Once()
				return keyGen
			},
			setupRepo: func(t *testing.T, ctx context.Context, description, url, userID, code string) *repoMocks.Repository {
				repo := repoMocks.NewRepository(t)
				input := &model.Bookmark{
					Description: description,
					URL:         url,
					Code:        code,
					UserID:      userID,
				}
				output := &model.Bookmark{
					Base: model.Base{
						ID: "11111111-2222-3333-4444-555555555555",
					},
					Description: description,
					URL:         url,
					Code:        code,
					UserID:      userID,
				}
				repo.On("CreateBookmark", ctx, input).Return(output, nil).Once()
				return repo
			},
			expectedError: nil,
			verifyBookmark: func(t *testing.T, bookmark *model.Bookmark) {
				assert.Equal(t, "11111111-2222-3333-4444-555555555555", bookmark.ID)
				assert.Equal(t, "My blog", bookmark.Description)
				assert.Equal(t, "https://truonglq.com", bookmark.URL)
				assert.Equal(t, "abcd1234", bookmark.Code)
				assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", bookmark.UserID)
			},
		},
		{
			name:        "error - key generator error",
			description: "My blog",
			url:         "https://truonglq.com",
			userID:      "550e8400-e29b-41d4-a716-446655440000",
			setupKeyGen: func(t *testing.T, code string, err error) *mockKeyGen.KeyGenerator {
				keyGen := mockKeyGen.NewKeyGenerator(t)
				keyGen.On("GenerateCode", codeLength).Return("", err).Once()
				return keyGen
			},
			setupRepo: func(t *testing.T, ctx context.Context, description, url, userID, code string) *repoMocks.Repository {
				return repoMocks.NewRepository(t)
			},
			expectedError:  testErrKeyGen,
			expectedCode:   "",
			verifyBookmark: nil,
		},
		{
			name:         "error - repository error",
			description:  "My blog",
			url:          "https://truonglq.com",
			userID:       "550e8400-e29b-41d4-a716-446655440000",
			expectedCode: "xyz98765",
			setupKeyGen: func(t *testing.T, code string, err error) *mockKeyGen.KeyGenerator {
				keyGen := mockKeyGen.NewKeyGenerator(t)
				keyGen.On("GenerateCode", codeLength).Return(code, nil).Once()
				return keyGen
			},
			setupRepo: func(t *testing.T, ctx context.Context, description, url, userID, code string) *repoMocks.Repository {
				repo := repoMocks.NewRepository(t)
				input := &model.Bookmark{
					Description: description,
					URL:         url,
					Code:        code,
					UserID:      userID,
				}
				repo.On("CreateBookmark", ctx, input).Return(nil, testErrDatabase).Once()
				return repo
			},
			expectedError:  testErrDatabase,
			verifyBookmark: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			keyGen := tc.setupKeyGen(t, tc.expectedCode, tc.expectedError)
			repo := tc.setupRepo(t, ctx, tc.description, tc.url, tc.userID, tc.expectedCode)

			svc := NewBookmarkSvc(repo, keyGen)

			result, err := svc.Create(ctx, tc.description, tc.url, tc.userID)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tc.verifyBookmark != nil {
					tc.verifyBookmark(t, result)
				}
			}
		})
	}
}
