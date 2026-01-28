package service

import (
	"context"
	db "database/sql"
	"errors"
	"testing"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	mockRepo "github.com/luongtruong20201/bookmark-management/internal/repositories/mocks"
	mockJWT "github.com/luongtruong20201/bookmark-management/pkg/jwt/mocks"
	mockUtils "github.com/luongtruong20201/bookmark-management/pkg/utils/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
					ID:          "550e8400-e29b-41d4-a716-446655440000",
					Username:    "johndoe",
					Password:    mockHashedPassword,
					DisplayName: "John Doe",
					Email:       "john.doe@example.com",
				}, nil).Once()
				return repoMock
			},
			expectedUser: &model.User{
				ID:          "550e8400-e29b-41d4-a716-446655440000",
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
					ID:          mockUserID,
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
					ID:          mockUserID,
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
					ID:          mockUserID,
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
					ID:          mockUserID,
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
					ID:          mockUserID,
					Username:    "johndoe",
					Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e",
					DisplayName: "John Doe",
					Email:       "john.doe@example.com",
				}, nil).Once()
				return repoMock
			},
			expectedUser: &model.User{
				ID:          mockUserID,
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
					ID:          "2f6c9d14-9e42-4c31-9e6a-0f8c4b2a7c55",
					Username:    "janedoe",
					Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5f",
					DisplayName: "Jane Doe",
					Email:       "jane.doe@example.com",
				}, nil).Once()
				return repoMock
			},
			expectedUser: &model.User{
				ID:          "2f6c9d14-9e42-4c31-9e6a-0f8c4b2a7c55",
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
					ID:          userID,
					Username:    "johndoe",
					DisplayName: displayName,
					Email:       email,
				}, nil).Once()
				return repoMock
			},
			expectedUser: &model.User{
				ID:          mockUserID,
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
