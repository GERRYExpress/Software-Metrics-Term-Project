package handler

import (
	"net/http"
	"software-metrics-project/service"
	"time"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	Service *service.OrderService
}

type OrderRequest struct {
	UserID        int       `json:"userId"`
	OrderDate     time.Time `json:"orderDate"`
	PaymentMethod string    `json:"paymentMethod"`
	FinalPrice    float64   `json:"finalPrice"`
}

func (h *OrderHandler) Create(c *gin.Context) {
	var order OrderRequest
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := h.Service.Create(order.UserID, order.PaymentMethod, order.FinalPrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}
