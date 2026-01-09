package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	service "github.com/luongtruong20201/bookmark-management/internal/services"
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
	req := urlShortenReq{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "unprocessable",
		})
		return
	}

	code, err := h.svc.ShortenURL(c, req.Url, req.Exp)
	if err != nil {
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
