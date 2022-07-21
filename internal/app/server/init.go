package server

import (
	nats_streaming "github.com/BlackRRR/nats-streaming/internal/app/services/nats-streaming"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	*nats_streaming.Service
}

func NewHandler(service *nats_streaming.Service) *Handler {
	return &Handler{service}

}

func (h *Handler) InitRouters() *gin.Engine {
	router := gin.New()
	router.LoadHTMLGlob("templates/page.html")

	router.GET("/:id", h.GetOrder)

	return router
}
