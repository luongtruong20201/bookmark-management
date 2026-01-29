// Package request provides helpers for binding and validating HTTP request inputs.
// It centralizes JSON, URI, query, and header binding plus validation logic.
package request

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
