package user

import (
	"context"
	db "database/sql"
	"errors"
	"testing"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	mockRepo "github.com/luongtruong20201/bookmark-management/internal/repositories/user/mocks"
	mockJWT "github.com/luongtruong20201/bookmark-management/pkg/jwt/mocks"
	mockUtils "github.com/luongtruong20201/bookmark-management/pkg/utils/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_Login(t *testing.T) {
	t.Parallel()

	const (
		mockHashedPassword = "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e"
		mockToken          = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI1NTBlODQwMC1lMjliLTQxZDQtYTcxNi00NDY2NTU0NDAwMDAiLCJpYXQiOjE2MDAwMDAwMDAsImV4cCI6MTYwMDA4NjQwMH0.test"
		mockUserID         = "550e8400-e29b-41d4-a716-446655440000"
	)

	var (
		testErrDatabase  = errors.New("database error")
		testErrNotFound  = errors.New("not found type")
		testErrJWT       = errors.New("jwt generation error")
		testErrClientErr = ErrClientErr
	)

	testCases := []struct {
		name              string
		username          string
		password          string
		setupMockRepo     func(t *testing.T, ctx context.Context, username string) *mockRepo.User
		setupMockHasher   func(t *testing.T, password, hashedPassword string, shouldVerify bool) *mockUtils.Hasher
		setupMockJWT      func(t *testing.T, userID string) *mockJWT.JWTGenerator
		expectedToken     string
		expectedError     error
		verifyTokenClaims bool
	}{
		{
			name:     "success - valid username and password",
			username: "johndoe",
			password: "password123",
			setupMockRepo: func(t *testing.T, ctx context.Context, username string) *mockRepo.User {
				repoMock := mockRepo.NewUser(t)
				repoMock.On("GetUserByUsername", ctx, username).Return(&model.User{
					Base: model.Base{
						ID: mockUserID,
					},
					Username:    "johndoe",
					Password:    mockHashedPassword,
					DisplayName: "John Doe",
					Email:       "john.doe@example.com",
				}, nil).Once()
				return repoMock
			},
			setupMockHasher: func(t *testing.T, password, hashedPassword string, shouldVerify bool) *mockUtils.Hasher {
				hasherMock := mockUtils.NewHasher(t)
				hasherMock.On("VerifyPassword", password, hashedPassword).Return(shouldVerify).Once()
				return hasherMock
			},
			setupMockJWT: func(t *testing.T, userID string) *mockJWT.JWTGenerator {
				jwtMock := mockJWT.NewJWTGenerator(t)
				jwtMock.On("GenerateToken", mock.Anything).
					Return(mockToken, nil).Once()
				return jwtMock
			},
			expectedToken:     mockToken,
			expectedError:     nil,
			verifyTokenClaims: true,
		},
		{
			name:     "error - user not found (db.ErrNoRows mapped to client error)",
			username: "nonexistent",
			password: "password123",
			setupMockRepo: func(t *testing.T, ctx context.Context, username string) *mockRepo.User {
				repoMock := mockRepo.NewUser(t)
				repoMock.On("GetUserByUsername", ctx, username).Return(nil, db.ErrNoRows).Once()
				return repoMock
			},
			setupMockHasher: func(t *testing.T, password, hashedPassword string, shouldVerify bool) *mockUtils.Hasher {
				hasherMock := mockUtils.NewHasher(t)
				return hasherMock
			},
			setupMockJWT: func(t *testing.T, userID string) *mockJWT.JWTGenerator {
				jwtMock := mockJWT.NewJWTGenerator(t)
				return jwtMock
			},
			expectedToken:     "",
			expectedError:     testErrClientErr,
			verifyTokenClaims: false,
		},
		{
			name:     "error - database error",
			username: "johndoe",
			password: "password123",
			setupMockRepo: func(t *testing.T, ctx context.Context, username string) *mockRepo.User {
				repoMock := mockRepo.NewUser(t)
				repoMock.On("GetUserByUsername", ctx, username).Return(nil, testErrDatabase).Once()
				return repoMock
			},
			setupMockHasher: func(t *testing.T, password, hashedPassword string, shouldVerify bool) *mockUtils.Hasher {
				hasherMock := mockUtils.NewHasher(t)
				return hasherMock
			},
			setupMockJWT: func(t *testing.T, userID string) *mockJWT.JWTGenerator {
				jwtMock := mockJWT.NewJWTGenerator(t)
				return jwtMock
			},
			expectedToken:     "",
			expectedError:     testErrDatabase,
			verifyTokenClaims: false,
		},
		{
			name:     "error - invalid password",
			username: "johndoe",
			password: "wrongpassword",
			setupMockRepo: func(t *testing.T, ctx context.Context, username string) *mockRepo.User {
				repoMock := mockRepo.NewUser(t)
				repoMock.On("GetUserByUsername", ctx, username).Return(&model.User{
					Base: model.Base{
						ID: mockUserID,
					},
					Username:    "johndoe",
					Password:    mockHashedPassword,
					DisplayName: "John Doe",
					Email:       "john.doe@example.com",
				}, nil).Once()
				return repoMock
			},
			setupMockHasher: func(t *testing.T, password, hashedPassword string, shouldVerify bool) *mockUtils.Hasher {
				hasherMock := mockUtils.NewHasher(t)
				hasherMock.On("VerifyPassword", password, hashedPassword).Return(shouldVerify).Once()
				return hasherMock
			},
			setupMockJWT: func(t *testing.T, userID string) *mockJWT.JWTGenerator {
				jwtMock := mockJWT.NewJWTGenerator(t)
				return jwtMock
			},
			expectedToken:     "",
			expectedError:     testErrClientErr,
			verifyTokenClaims: false,
		},
		{
			name:     "error - JWT generation fails",
			username: "johndoe",
			password: "password123",
			setupMockRepo: func(t *testing.T, ctx context.Context, username string) *mockRepo.User {
				repoMock := mockRepo.NewUser(t)
				repoMock.On("GetUserByUsername", ctx, username).Return(&model.User{
					Base: model.Base{
						ID: mockUserID,
					},
					Username:    "johndoe",
					Password:    mockHashedPassword,
					DisplayName: "John Doe",
					Email:       "john.doe@example.com",
				}, nil).Once()
				return repoMock
			},
			setupMockHasher: func(t *testing.T, password, hashedPassword string, shouldVerify bool) *mockUtils.Hasher {
				hasherMock := mockUtils.NewHasher(t)
				hasherMock.On("VerifyPassword", password, hashedPassword).Return(true).Once()
				return hasherMock
			},
			setupMockJWT: func(t *testing.T, userID string) *mockJWT.JWTGenerator {
				jwtMock := mockJWT.NewJWTGenerator(t)
				jwtMock.On("GenerateToken", mock.Anything).Return("", testErrJWT).Once()
				return jwtMock
			},
			expectedToken:     "",
			expectedError:     testErrJWT,
			verifyTokenClaims: false,
		},
		{
			name:     "error - empty username",
			username: "",
			password: "password123",
			setupMockRepo: func(t *testing.T, ctx context.Context, username string) *mockRepo.User {
				repoMock := mockRepo.NewUser(t)
				repoMock.On("GetUserByUsername", ctx, username).Return(nil, testErrNotFound).Once()
				return repoMock
			},
			setupMockHasher: func(t *testing.T, password, hashedPassword string, shouldVerify bool) *mockUtils.Hasher {
				hasherMock := mockUtils.NewHasher(t)
				return hasherMock
			},
			setupMockJWT: func(t *testing.T, userID string) *mockJWT.JWTGenerator {
				jwtMock := mockJWT.NewJWTGenerator(t)
				return jwtMock
			},
			expectedToken:     "",
			expectedError:     testErrNotFound,
			verifyTokenClaims: false,
		},
		{
			name:     "error - empty password",
			username: "johndoe",
			password: "",
			setupMockRepo: func(t *testing.T, ctx context.Context, username string) *mockRepo.User {
				repoMock := mockRepo.NewUser(t)
				repoMock.On("GetUserByUsername", ctx, username).Return(&model.User{
					Base: model.Base{
						ID: mockUserID,
					},
					Username:    "johndoe",
					Password:    mockHashedPassword,
					DisplayName: "John Doe",
					Email:       "john.doe@example.com",
				}, nil).Once()
				return repoMock
			},
			setupMockHasher: func(t *testing.T, password, hashedPassword string, shouldVerify bool) *mockUtils.Hasher {
				hasherMock := mockUtils.NewHasher(t)
				hasherMock.On("VerifyPassword", password, hashedPassword).Return(false).Once()
				return hasherMock
			},
			setupMockJWT: func(t *testing.T, userID string) *mockJWT.JWTGenerator {
				jwtMock := mockJWT.NewJWTGenerator(t)
				return jwtMock
			},
			expectedToken:     "",
			expectedError:     testErrClientErr,
			verifyTokenClaims: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			repoMock := tc.setupMockRepo(t, ctx, tc.username)
			hasherMock := tc.setupMockHasher(t, tc.password, mockHashedPassword, tc.name == "success - valid username and password")
			jwtMock := tc.setupMockJWT(t, mockUserID)
			svc := NewUser(repoMock, hasherMock, jwtMock)

			token, err := svc.Login(ctx, tc.username, tc.password)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				assert.Equal(t, tc.expectedToken, token)
			}
		})
	}
}
