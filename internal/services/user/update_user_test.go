package user

import (
	"context"
	"errors"
	"testing"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	mockRepo "github.com/luongtruong20201/bookmark-management/internal/repositories/user/mocks"
	"github.com/stretchr/testify/assert"
)

func TestUserService_UpdateUserProfile(t *testing.T) {
	t.Parallel()

	const mockUserID = "550e8400-e29b-41d4-a716-446655440000"

	var (
		testErrDatabase = errors.New("database error")
		testErrNotFound = errors.New("not found type")
	)

	testCases := []struct {
		name          string
		userID        string
		displayName   string
		email         string
		setupMockRepo func(t *testing.T, ctx context.Context, userID, displayName, email string) *mockRepo.User
		expectedUser  *model.User
		expectedError error
	}{
		{
			name:        "success - update existing user profile",
			userID:      mockUserID,
			displayName: "John Updated",
			email:       "john.updated@example.com",
			setupMockRepo: func(t *testing.T, ctx context.Context, userID, displayName, email string) *mockRepo.User {
				repoMock := mockRepo.NewUser(t)
				repoMock.On("UpdateUserProfile", ctx, userID, displayName, email).Return(&model.User{
					Base: model.Base{
						ID: userID,
					},
					Username:    "johndoe",
					DisplayName: displayName,
					Email:       email,
				}, nil).Once()
				return repoMock
			},
			expectedUser: &model.User{
				Base: model.Base{
					ID: mockUserID,
				},
				Username:    "johndoe",
				DisplayName: "John Updated",
				Email:       "john.updated@example.com",
			},
			expectedError: nil,
		},
		{
			name:        "error - user not found",
			userID:      "00000000-0000-0000-0000-000000000000",
			displayName: "John Updated",
			email:       "john.updated@example.com",
			setupMockRepo: func(t *testing.T, ctx context.Context, userID, displayName, email string) *mockRepo.User {
				repoMock := mockRepo.NewUser(t)
				repoMock.On("UpdateUserProfile", ctx, userID, displayName, email).Return(nil, testErrNotFound).Once()
				return repoMock
			},
			expectedUser:  nil,
			expectedError: testErrNotFound,
		},
		{
			name:        "error - database error",
			userID:      mockUserID,
			displayName: "John Updated",
			email:       "john.updated@example.com",
			setupMockRepo: func(t *testing.T, ctx context.Context, userID, displayName, email string) *mockRepo.User {
				repoMock := mockRepo.NewUser(t)
				repoMock.On("UpdateUserProfile", ctx, userID, displayName, email).Return(nil, testErrDatabase).Once()
				return repoMock
			},
			expectedUser:  nil,
			expectedError: testErrDatabase,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			repoMock := tc.setupMockRepo(t, ctx, tc.userID, tc.displayName, tc.email)
			svc := NewUser(repoMock, nil, nil)

			result, err := svc.UpdateUserProfile(ctx, tc.userID, tc.displayName, tc.email)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.expectedUser.ID, result.ID)
				assert.Equal(t, tc.expectedUser.Username, result.Username)
				assert.Equal(t, tc.expectedUser.DisplayName, result.DisplayName)
				assert.Equal(t, tc.expectedUser.Email, result.Email)
			}
		})
	}
}
