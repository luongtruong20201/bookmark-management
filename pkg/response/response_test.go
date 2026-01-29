package response

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestInputFieldError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		err             error
		expectedMsg     string
		expectedDetails interface{}
		hasDetails      bool
	}{
		{
			name:            "error - non-validation error returns InputErrResponse",
			err:             errors.New("some error"),
			expectedMsg:     "Input Error",
			expectedDetails: nil,
			hasDetails:      false,
		},
		{
			name:            "error - nil error returns InputErrResponse",
			err:             nil,
			expectedMsg:     "Input Error",
			expectedDetails: nil,
			hasDetails:      false,
		},
		{
			name:        "success - validation errors are formatted",
			err:         createValidationError(),
			expectedMsg: "Input error",
			hasDetails:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := InputFieldError(tc.err)

			assert.Equal(t, tc.expectedMsg, result.Message)

			if tc.hasDetails {
				assert.NotNil(t, result.Details)
				details, ok := result.Details.([]string)
				assert.True(t, ok, "details should be a slice of strings")
				assert.NotEmpty(t, details)
				for _, detail := range details {
					assert.Contains(t, detail, "is invalid")
					assert.Contains(t, detail, "(")
					assert.Contains(t, detail, ")")
				}
			} else {
				assert.Equal(t, tc.expectedDetails, result.Details)
			}
		})
	}
}

func TestInputFieldError_ValidationErrors(t *testing.T) {
	t.Parallel()

	type TestStruct struct {
		Email    string `validate:"required,email"`
		Username string `validate:"required,min=3"`
		Age      int    `validate:"required,min=18"`
	}

	validate := validator.New()

	testCases := []struct {
		name            string
		data            TestStruct
		expectedFields  []string
		expectedMessage string
	}{
		{
			name: "success - multiple validation errors",
			data: TestStruct{
				Email:    "invalid-email",
				Username: "ab",
				Age:      15,
			},
			expectedFields:  []string{"Email", "Username", "Age"},
			expectedMessage: "Input error",
		},
		{
			name: "success - single validation error",
			data: TestStruct{
				Email:    "",
				Username: "validuser",
				Age:      25,
			},
			expectedFields:  []string{"Email"},
			expectedMessage: "Input error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := validate.Struct(tc.data)
			assert.Error(t, err)

			result := InputFieldError(err)

			assert.Equal(t, tc.expectedMessage, result.Message)
			assert.NotNil(t, result.Details)

			details, ok := result.Details.([]string)
			assert.True(t, ok)
			assert.Len(t, details, len(tc.expectedFields))

			for _, field := range tc.expectedFields {
				found := false
				for _, detail := range details {
					if containsField(detail, field) {
						found = true
						break
					}
				}
				assert.True(t, found, "field %s should be in details", field)
			}
		})
	}
}

func TestMessage_Structure(t *testing.T) {
	t.Parallel()

	t.Run("InternalErrResponse has correct structure", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "Processing Error", InternalErrResponse.Message)
		assert.Nil(t, InternalErrResponse.Details)
	})

	t.Run("InputErrResponse has correct structure", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "Input Error", InputErrResponse.Message)
		assert.Nil(t, InputErrResponse.Details)
	})

	t.Run("Message with details", func(t *testing.T) {
		t.Parallel()

		msg := Message{
			Message: "Test message",
			Details: []string{"detail1", "detail2"},
		}

		assert.Equal(t, "Test message", msg.Message)
		assert.NotNil(t, msg.Details)
		details, ok := msg.Details.([]string)
		assert.True(t, ok)
		assert.Len(t, details, 2)
	})
}

func createValidationError() error {
	type TestStruct struct {
		Email string `validate:"required,email"`
		Name  string `validate:"required"`
	}

	validate := validator.New()
	testData := TestStruct{
		Email: "invalid-email",
		Name:  "",
	}

	return validate.Struct(testData)
}

func containsField(detail, field string) bool {
	return len(detail) > len(field) && detail[:len(field)] == field
}
