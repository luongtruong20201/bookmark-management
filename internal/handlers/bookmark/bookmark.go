package bookmark

import (
	"github.com/gin-gonic/gin"
	"github.com/luongtruong20201/bookmark-management/internal/services/bookmark"
	"github.com/luongtruong20201/bookmark-management/internal/services/queue"
)

// Handler defines the HTTP handler interface for bookmark endpoints.
// It exposes methods used by the router to create and manage bookmarks.
//
//go:generate mockery --name Handler --filename bookmark_handler.go
type Handler interface {
	Create(c *gin.Context)
	GetBookmarks(c *gin.Context)
	UpdateBookmark(c *gin.Context)
	DeleteBookmark(c *gin.Context)
	ImportBookmarks(c *gin.Context)
}

// bookmarkHandler implements the Handler interface and wires bookmark
// service calls to HTTP requests/responses.
type bookmarkHandler struct {
	svc   bookmark.Service
	queue queue.Service
}

// NewBookmarkHandler creates a new bookmark HTTP handler with the given services.
// It wires the core bookmark service (with caching) and the background queue
// service used for asynchronous bookmark imports.
func NewBookmarkHandler(svc bookmark.Service, queueSvc queue.Service) Handler {
	return &bookmarkHandler{
		svc:   svc,
		queue: queueSvc,
	}
}
