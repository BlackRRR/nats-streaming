package server

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

func (h *Handler) GetModel(c *gin.Context) {
	id := c.Param("id")

	order, err := h.Service.GetOrder(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	if order == nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.New("order null"))
		return
	}

	c.HTML(http.StatusAccepted, "page.html", gin.H{
		"text": string(order),
	})
}
