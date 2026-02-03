package shorten

import (
	"github.com/gin-gonic/gin"
	service "github.com/luongtruong20201/bookmark-management/internal/services/shorten"
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

// urlShortenHandler implements the ShortenURL interface and provides HTTP handlers
// for URL shortening operations. It encapsulates the shorten URL service dependency
// for business logic execution.
type urlShortenHandler struct {
	svc service.ShortenURL
}

// NewShortenURL creates a new shorten URL handler with the provided shorten URL service.
func NewShortenURL(svc service.ShortenURL) ShortenURL {
	return &urlShortenHandler{
		svc: svc,
	}
}
