package user

import (
	"context"
	"errors"
	"testing"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	mockRepo "github.com/luongtruong20201/bookmark-management/internal/repositories/user/mocks"
	mockUtils "github.com/luongtruong20201/bookmark-management/pkg/utils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestUserService_CreateUser(t *testing.T) {
	t.Parallel()

	const mockHashedPassword = "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e"

	var (
		testErrDatabase     = errors.New("database error")
		testErrDuplicateKey = errors.New("duplicate key value violates unique constraint")
	)

	testCases := []struct {
		name            string
		username        string
		password        string
		displayName     string
		email           string
		setupMockHasher func(t *testing.T, password string) *mockUtils.Hasher
		setupMockRepo   func(t *testing.T, ctx context.Context) *mockRepo.User
		expectedUser    *model.User
		expectedError   error
		verifyPassword  bool
	}{
		{
			name:        "success",
			username:    "johndoe",
			password:    "password123",
			displayName: "John Doe",
			email:       "john.doe@example.com",
			setupMockHasher: func(t *testing.T, password string) *mockUtils.Hasher {
				hasherMock := mockUtils.NewHasher(t)
				hasherMock.On("HashPassword", password).Return(mockHashedPassword).Once()
				return hasherMock
			},
			setupMockRepo: func(t *testing.T, ctx context.Context) *mockRepo.User {
				expectedUserInput := &model.User{
					Username:    "johndoe",
					Password:    mockHashedPassword,
					DisplayName: "John Doe",
					Email:       "john.doe@example.com",
				}
				repoMock := mockRepo.NewUser(t)
				repoMock.On("CreateUser", ctx, expectedUserInput).Return(&model.User{
					Base: model.Base{
						ID: "550e8400-e29b-41d4-a716-446655440000",
					},
					Username:    "johndoe",
					Password:    mockHashedPassword,
					DisplayName: "John Doe",
					Email:       "john.doe@example.com",
				}, nil).Once()
				return repoMock
			},
			expectedUser: &model.User{
				Base: model.Base{
					ID: "550e8400-e29b-41d4-a716-446655440000",
				},
				Username:    "johndoe",
				DisplayName: "John Doe",
				Email:       "john.doe@example.com",
			},
			expectedError:  nil,
			verifyPassword: true,
		},
		{
			name:        "repository error",
			username:    "johndoe",
			password:    "password123",
			displayName: "John Doe",
			email:       "john.doe@example.com",
			setupMockHasher: func(t *testing.T, password string) *mockUtils.Hasher {
				hasherMock := mockUtils.NewHasher(t)
				hasherMock.On("HashPassword", password).Return(mockHashedPassword).Once()
				return hasherMock
			},
			setupMockRepo: func(t *testing.T, ctx context.Context) *mockRepo.User {
				expectedUserInput := &model.User{
					Username:    "johndoe",
					Password:    mockHashedPassword,
					DisplayName: "John Doe",
					Email:       "john.doe@example.com",
				}
				repoMock := mockRepo.NewUser(t)
				repoMock.On("CreateUser", ctx, expectedUserInput).Return(nil, testErrDatabase).Once()
				return repoMock
			},
			expectedUser:   nil,
			expectedError:  testErrDatabase,
			verifyPassword: false,
		},
		{
			name:        "duplicate username error",
			username:    "johndoe",
			password:    "password123",
			displayName: "John Doe",
			email:       "john.doe@example.com",
			setupMockHasher: func(t *testing.T, password string) *mockUtils.Hasher {
				hasherMock := mockUtils.NewHasher(t)
				hasherMock.On("HashPassword", password).Return(mockHashedPassword).Once()
				return hasherMock
			},
			setupMockRepo: func(t *testing.T, ctx context.Context) *mockRepo.User {
				expectedUserInput := &model.User{
					Username:    "johndoe",
					Password:    mockHashedPassword,
					DisplayName: "John Doe",
					Email:       "john.doe@example.com",
				}
				repoMock := mockRepo.NewUser(t)
				repoMock.On("CreateUser", ctx, expectedUserInput).Return(nil, testErrDuplicateKey).Once()
				return repoMock
			},
			expectedUser:   nil,
			expectedError:  testErrDuplicateKey,
			verifyPassword: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			hasherMock := tc.setupMockHasher(t, tc.password)
			repoMock := tc.setupMockRepo(t, ctx)
			svc := NewUser(repoMock, hasherMock, nil)

			result, err := svc.CreateUser(ctx, tc.username, tc.password, tc.displayName, tc.email)

			if tc.expectedError != nil {
				assert.ErrorIs(t, tc.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tc.expectedUser.ID, result.ID)
				assert.Equal(t, tc.expectedUser.Username, result.Username)
				assert.Equal(t, tc.expectedUser.DisplayName, result.DisplayName)
				assert.Equal(t, tc.expectedUser.Email, result.Email)

				if tc.verifyPassword {
					assert.Equal(t, mockHashedPassword, result.Password)
				}
			}
		})
	}
}
