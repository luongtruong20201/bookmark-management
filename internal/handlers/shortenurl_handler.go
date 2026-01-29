package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	service "github.com/luongtruong20201/bookmark-management/internal/services"
	requtil "github.com/luongtruong20201/bookmark-management/pkg/request"
	"github.com/rs/zerolog/log"
)

// urlShortenReq represents the request body for shortening a URL.
type urlShortenReq struct {
	Url string `json:"url" binding:"required,url" example:"https://fb.com"`
	Exp int    `json:"exp" binding:"required,gte=0,lte=604800" example:"3600"`
}

// urlShortenRes represents the response body after shortening a URL.
type urlShortenRes struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

// ShortenURL defines the interface for shorten URL handlers.
// It provides methods to handle URL shortening requests.
type ShortenURL interface {
	ShortenURL(*gin.Context)
	// GetURL handles the request to retrieve and redirect to the original URL from a short code.
	GetURL(*gin.Context)
}

type urlShortenHandler struct {
	svc service.ShortenURL
}

// NewShortenURL creates a new shorten URL handler with the provided shorten URL service.
func NewShortenURL(svc service.ShortenURL) ShortenURL {
	return &urlShortenHandler{
		svc: svc,
	}
}

// ShortenURL handles the URL shortening endpoint request. It validates the input,
// generates a short code for the URL, and returns the shortened URL code.
// @Summary Shorten URL
// @Description Create a shortened URL with an optional expiration time (in seconds, max 604800)
// @Tags url
// @Accept json
// @Produce json
// @Param request body urlShortenReq true "URL shortening request"
// @Success 200 {object} urlShortenRes "Successfully shortened URL"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /links/shorten [post]
func (h *urlShortenHandler) ShortenURL(c *gin.Context) {
	req, err := requtil.BindInputFromRequest[urlShortenReq](c)
	if err != nil {
		return
	}

	code, err := h.svc.ShortenURL(c, req.Url, req.Exp)
	if err != nil {
		log.Error().Str("url", req.Url).Err(err).Msg("error when create shorten url")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, urlShortenRes{
		Message: "OK",
		Code:    code,
	})
}

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
