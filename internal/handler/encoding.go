package handler

import (
	"errors"
	"net/http"

	"github.com/bambi2/url-shorter/internal/service"
	"github.com/gin-gonic/gin"
)

type inputUrl struct {
	Url string `json:"url" binding:"required"`
}

func (h *Handler) encodeBase63(c *gin.Context) {
	var url inputUrl

	if err := c.BindJSON(&url); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid json")
		return
	}

	encodedUrl, err := h.service.Encoder.Base63(url.Url)
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
		"encodedUrl": encodedUrl,
	})
}
