package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GetOrder(c *gin.Context) {
	id := c.Param("id")

	order, err := h.Service.GetOrder(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "page.html", gin.H{
			"text": "failed to get order",
		})
		return
	}

	if order == nil {
		c.HTML(http.StatusNotFound, "page.html", gin.H{
			"text": "order not found",
		})
		return
	}

	c.HTML(http.StatusOK, "page.html", gin.H{
		"text": string(order),
	})
}
