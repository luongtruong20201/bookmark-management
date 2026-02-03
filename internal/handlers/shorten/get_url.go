package shorten

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	service "github.com/luongtruong20201/bookmark-management/internal/services/shorten"
	"github.com/rs/zerolog/log"
)

// GetURL handles the request to retrieve the original URL from a short code.
// It extracts the code from the URL parameter, validates it, and redirects to the original URL.
// If the code is not found, it returns a 400 Bad Request with an error message.
// If an internal error occurs, it returns a 500 Internal Server Error.
// On success, it performs a 301 Permanent Redirect to the original URL.
// @Summary Get original URL by code
// @Description Retrieve and redirect to the original URL using the shortened code
// @Tags url
// @Accept json
// @Produce json
// @Param code path string true "Short URL code"
// @Success 301 "Permanent redirect to original URL"
// @Failure 400 {object} map[string]string "Code not found or invalid"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/links/{code} [get]
func (s *urlShortenHandler) GetURL(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "unprocessable",
		})
		return
	}

	url, err := s.svc.GetURL(c, code)
	if err != nil {
		if errors.Is(err, service.ErrCodeNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "url not found",
			})
			return
		}

		log.Error().Str("code", code).Err(err).Msg("error when get original url from code")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	c.Redirect(http.StatusMovedPermanently, url)
}
