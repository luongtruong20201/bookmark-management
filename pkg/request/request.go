// Package request provides helpers for binding and validating HTTP request inputs.
// It centralizes JSON, URI, query, and header binding plus validation logic.
package request

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/luongtruong20201/bookmark-management/internal/utils"
	"github.com/luongtruong20201/bookmark-management/pkg/response"
)

// BindInputFromRequest binds and validates all supported input sources (JSON body,
// URI params, query params, and headers) into a typed struct T. On validation
// failure it writes a 400 response using response.InputFieldError and aborts
// the Gin context, returning the encountered error.
func BindInputFromRequest[T any](c *gin.Context) (*T, error) {
	reqInput := new(T)

	if err := c.ShouldBindJSON(reqInput); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		c.Abort()
		return nil, err
	}

	if err := c.ShouldBindUri(reqInput); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		c.Abort()
		return nil, err
	}

	if err := c.ShouldBindQuery(reqInput); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		c.Abort()
		return nil, err
	}

	if err := c.ShouldBindHeader(reqInput); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		c.Abort()
		return nil, err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(reqInput); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		c.Abort()
		return nil, err
	}

	return reqInput, nil
}

// BindInputFromRequestWithAuth binds and validates all supported input sources (JSON body,
// URI params, query params, and headers) and extracts the user ID from JWT claims.
// It returns the bound input, user ID, and any error.
// On validation failure it writes a 400 response using response.InputFieldError and aborts
// the Gin context. If user ID extraction fails, it writes a 401 response and aborts the context.
func BindInputFromRequestWithAuth[T any](c *gin.Context) (*T, string, error) {
	input, err := BindInputFromRequest[T](c)
	if err != nil {
		return nil, "", err
	}

	userId, err := utils.GetUserIDFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &response.Message{
			Message: "Invalid jwt token",
		})
		c.Abort()
		return nil, "", err
	}

	return input, userId, nil
}

// BindInputFromQueryWithAuth binds and validates query parameters from the request
// and extracts the user ID from JWT claims. It returns the bound input, user ID, and any error.
// On validation failure it writes a 400 response using response.InputFieldError and aborts
// the Gin context. If user ID extraction fails, it writes a 401 response and aborts the context.
func BindInputFromQueryWithAuth[T any](c *gin.Context) (*T, string, error) {
	reqInput := new(T)

	if err := c.ShouldBindQuery(reqInput); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		c.Abort()
		return nil, "", err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(reqInput); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		c.Abort()
		return nil, "", err
	}

	userId, err := utils.GetUserIDFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &response.Message{
			Message: "Invalid jwt token",
		})
		c.Abort()
		return nil, "", err
	}

	return reqInput, userId, nil
}

// BindInputFromUriWithAuth binds and validates URI parameters from the request
// and extracts the user ID from JWT claims. It returns the bound input, user ID, and any error.
// On validation failure it writes a 400 response using response.InputFieldError and aborts
// the Gin context. If user ID extraction fails, it writes a 401 response and aborts the context.
func BindInputFromUriWithAuth[T any](c *gin.Context) (*T, string, error) {
	reqInput := new(T)

	if err := c.ShouldBindUri(reqInput); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		c.Abort()
		return nil, "", err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(reqInput); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		c.Abort()
		return nil, "", err
	}

	userId, err := utils.GetUserIDFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &response.Message{
			Message: "Invalid jwt token",
		})
		c.Abort()
		return nil, "", err
	}

	return reqInput, userId, nil
}
