package handler

import (
	"github.com/bambi2/url-shorter/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		url := api.Group("/url")
		{
			base63 := url.Group("base63")
			{
				base63.POST("/", h.encodeBase63)
				base63.GET("/:encoded_url", h.decodeBase63)
			}
		}
	}

	return router
}
