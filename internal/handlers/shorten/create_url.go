package shorten

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luongtruong20201/bookmark-management/pkg/request"
	"github.com/rs/zerolog/log"
)

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
	req, err := request.BindInputFromRequest[urlShortenReq](c)
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
