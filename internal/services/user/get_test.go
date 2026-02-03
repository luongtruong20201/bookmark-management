package user

import (
	"context"
	"errors"
	"testing"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	mockRepo "github.com/luongtruong20201/bookmark-management/internal/repositories/user/mocks"
	"github.com/stretchr/testify/assert"
)

func TestUserService_GetUserByID(t *testing.T) {
	t.Parallel()

	const mockUserID = "550e8400-e29b-41d4-a716-446655440000"

	var (
		testErrDatabase = errors.New("database error")
		testErrNotFound = errors.New("not found type")
	)

	testCases := []struct {
		name          string
		userID        string
		setupMockRepo func(t *testing.T, ctx context.Context, userID string) *mockRepo.User
		expectedUser  *model.User
		expectedError error
	}{
		{
			name:   "success - get existing user by ID",
			userID: mockUserID,
			setupMockRepo: func(t *testing.T, ctx context.Context, userID string) *mockRepo.User {
				repoMock := mockRepo.NewUser(t)
				repoMock.On("GetUserByID", ctx, userID).Return(&model.User{
					Base: model.Base{
						ID: mockUserID,
					},
					Username:    "johndoe",
					Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e",
					DisplayName: "John Doe",
					Email:       "john.doe@example.com",
				}, nil).Once()
				return repoMock
			},
			expectedUser: &model.User{
				Base: model.Base{
					ID: mockUserID,
				},
				Username:    "johndoe",
				Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e",
				DisplayName: "John Doe",
				Email:       "john.doe@example.com",
			},
			expectedError: nil,
		},
		{
			name:   "success - get another existing user by ID",
			userID: "2f6c9d14-9e42-4c31-9e6a-0f8c4b2a7c55",
			setupMockRepo: func(t *testing.T, ctx context.Context, userID string) *mockRepo.User {
				repoMock := mockRepo.NewUser(t)
				repoMock.On("GetUserByID", ctx, userID).Return(&model.User{
					Base: model.Base{
						ID: "2f6c9d14-9e42-4c31-9e6a-0f8c4b2a7c55",
					},
					Username:    "janedoe",
					Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5f",
					DisplayName: "Jane Doe",
					Email:       "jane.doe@example.com",
				}, nil).Once()
				return repoMock
			},
			expectedUser: &model.User{
				Base: model.Base{
					ID: "2f6c9d14-9e42-4c31-9e6a-0f8c4b2a7c55",
				},
				Username:    "janedoe",
				Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5f",
				DisplayName: "Jane Doe",
				Email:       "jane.doe@example.com",
			},
			expectedError: nil,
		},
		{
			name:   "error - user not found",
			userID: "00000000-0000-0000-0000-000000000000",
			setupMockRepo: func(t *testing.T, ctx context.Context, userID string) *mockRepo.User {
				repoMock := mockRepo.NewUser(t)
				repoMock.On("GetUserByID", ctx, userID).Return(nil, testErrNotFound).Once()
				return repoMock
			},
			expectedUser:  nil,
			expectedError: testErrNotFound,
		},
		{
			name:   "error - database error",
			userID: mockUserID,
			setupMockRepo: func(t *testing.T, ctx context.Context, userID string) *mockRepo.User {
				repoMock := mockRepo.NewUser(t)
				repoMock.On("GetUserByID", ctx, userID).Return(nil, testErrDatabase).Once()
				return repoMock
			},
			expectedUser:  nil,
			expectedError: testErrDatabase,
		},
		{
			name:   "error - empty user ID",
			userID: "",
			setupMockRepo: func(t *testing.T, ctx context.Context, userID string) *mockRepo.User {
				repoMock := mockRepo.NewUser(t)
				repoMock.On("GetUserByID", ctx, userID).Return(nil, testErrNotFound).Once()
				return repoMock
			},
			expectedUser:  nil,
			expectedError: testErrNotFound,
		},
		{
			name:   "error - invalid UUID format",
			userID: "invalid-uuid",
			setupMockRepo: func(t *testing.T, ctx context.Context, userID string) *mockRepo.User {
				repoMock := mockRepo.NewUser(t)
				repoMock.On("GetUserByID", ctx, userID).Return(nil, testErrNotFound).Once()
				return repoMock
			},
			expectedUser:  nil,
			expectedError: testErrNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			repoMock := tc.setupMockRepo(t, ctx, tc.userID)
			svc := NewUser(repoMock, nil, nil)

			result, err := svc.GetUserByID(ctx, tc.userID)

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
				assert.Equal(t, tc.expectedUser.Password, result.Password)
			}
		})
	}
}
