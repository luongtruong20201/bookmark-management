package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	model "github.com/luongtruong20201/bookmark-management/internal/models"
	"github.com/luongtruong20201/bookmark-management/internal/services/mocks"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_RegisterUser(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	var (
		testErrDatabase = errors.New("database error")
	)

	testCases := []struct {
		name           string
		requestBody    interface{}
		setupMockSvc   func(t *testing.T, ctx context.Context) *mocks.User
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "success",
			requestBody: createUserInputBody{
				Username:    "johndoe",
				Password:    "password123",
				DisplayName: "John Doe",
				Email:       "john.doe@example.com",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				svcMock.On("CreateUser", ctx, "johndoe", "password123", "John Doe", "john.doe@example.com").
					Return(&model.User{
						ID:          "550e8400-e29b-41d4-a716-446655440000",
						Username:    "johndoe",
						DisplayName: "John Doe",
						Email:       "john.doe@example.com",
					}, nil).Once()
				return svcMock
			},
			expectedStatus: http.StatusOK,
			expectedBody: model.User{
				ID:          "550e8400-e29b-41d4-a716-446655440000",
				Username:    "johndoe",
				DisplayName: "John Doe",
				Email:       "john.doe@example.com",
			},
		},
		{
			name: "invalid request body - missing required fields",
			requestBody: createUserInputBody{
				Username: "",
				Password: "password123",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				return svcMock
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name: "invalid request body - invalid email format",
			requestBody: createUserInputBody{
				Username:    "johndoe",
				Password:    "password123",
				DisplayName: "John Doe",
				Email:       "invalid-email",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				return svcMock
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name: "invalid request body - password too short",
			requestBody: createUserInputBody{
				Username:    "johndoe",
				Password:    "12345",
				DisplayName: "John Doe",
				Email:       "john.doe@example.com",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				return svcMock
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name: "internal server error",
			requestBody: createUserInputBody{
				Username:    "johndoe",
				Password:    "password123",
				DisplayName: "John Doe",
				Email:       "john.doe@example.com",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				svcMock.On("CreateUser", ctx, "johndoe", "password123", "John Doe", "john.doe@example.com").
					Return(nil, testErrDatabase).Once()
				return svcMock
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
		{
			name:        "invalid JSON",
			requestBody: "invalid json string",
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				return svcMock
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)

			var reqBody []byte
			var err error
			if str, ok := tc.requestBody.(string); ok {
				reqBody = []byte(str)
			} else {
				reqBody, err = json.Marshal(tc.requestBody)
				assert.NoError(t, err)
			}

			ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/register", bytes.NewBuffer(reqBody))
			ctx.Request.Header.Set("Content-Type", "application/json")
			mockSvc := tc.setupMockSvc(t, ctx)
			testHandler := NewUser(mockSvc)

			testHandler.RegisterUser(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectedBody != nil {
				var responseBody map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
				assert.NoError(t, err)

				assert.Equal(t, "Register an user successfully!", responseBody["message"])

				data, ok := responseBody["data"].(map[string]interface{})
				assert.True(t, ok, "data should be a map")

				expectedUser := tc.expectedBody.(model.User)
				assert.Equal(t, expectedUser.ID, data["id"])
				assert.Equal(t, expectedUser.Username, data["username"])
				assert.Equal(t, expectedUser.DisplayName, data["display_name"])
				assert.Equal(t, expectedUser.Email, data["email"])
			}
		})
	}
}
