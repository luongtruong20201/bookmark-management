package bookmark

import (
	"github.com/gin-gonic/gin"
	"github.com/luongtruong20201/bookmark-management/internal/services/bookmark"
)

// Handler defines the HTTP handler interface for bookmark endpoints.
// It exposes methods used by the router to create and manage bookmarks.
type Handler interface {
	Create(c *gin.Context)
	GetBookmarks(c *gin.Context)
	UpdateBookmark(c *gin.Context)
	DeleteBookmark(c *gin.Context)
}

// bookmarkHandler implements the Handler interface and wires bookmark
// service calls to HTTP requests/responses.
type bookmarkHandler struct {
	svc bookmark.Service
}

// NewBookmarkHandler creates a new bookmark HTTP handler with the given service.
func NewBookmarkHandler(svc bookmark.Service) Handler {
	return &bookmarkHandler{
		svc: svc,
	}
}
