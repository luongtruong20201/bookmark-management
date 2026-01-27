package response

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

// Message represents a standardized API response message structure.
// It contains a message string and optional details for additional information.
type Message struct {
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

var (
	// InternalErrResponse represents a standardized internal server error response.
	// It is used when an unexpected error occurs during request processing.
	InternalErrResponse = Message{
		Message: "Processing Error",
		Details: nil,
	}
	// InputErrResponse represents a standardized input validation error response.
	// It is used when the request input fails basic validation checks.
	InputErrResponse = Message{
		Message: "Input Error",
		Details: nil,
	}
)

// InputFieldError processes validation errors and returns a formatted Message response.
// If the error is a validator.ValidationErrors, it extracts field-specific error messages.
// Otherwise, it returns the default InputErrResponse.
func InputFieldError(err error) Message {
	if ok := errors.As(err, &validator.ValidationErrors{}); !ok {
		return InputErrResponse
	}

	var errs []string
	for _, err := range err.(validator.ValidationErrors) {
		errs = append(errs, err.Field()+" is invalid ("+err.Tag()+")")
	}

	return Message{
		Message: "Input error",
		Details: errs,
	}
}
