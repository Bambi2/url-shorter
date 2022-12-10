package handler

import (
	"errors"
	"net/http"

	"github.com/bambi2/url-shorter/internal/service"
	"github.com/gin-gonic/gin"
)

func (h *Handler) decodeBase63(c *gin.Context) {
	encodedUrl := c.Param("encoded_url")

	url, err := h.service.Decoder.Base63(encodedUrl)
	if err != nil {
		var serviceErr *service.ServiceError
		switch {
		case errors.As(err, &serviceErr):
			newErrorResponse(c, serviceErr.StatusCode, serviceErr.Msg)
		default:
			newErrorResponse(c, http.StatusInternalServerError, "unknown server error")
		}
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"url": url,
	})
}
