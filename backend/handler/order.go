package handler

import (
	"net/http"
	"software-metrics-project/model"
	"software-metrics-project/service"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	Service *service.OrderService
}

func (h *OrderHandler) Create(c *gin.Context) {
	var order model.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := h.Service.Create(order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}
