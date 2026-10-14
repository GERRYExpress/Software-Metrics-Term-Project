package handler

import (
	"software-metrics-project/service"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	Service *service.CartService
}

func (h *CartHandler) GetCarts(c *gin.Context) { listByUserID(c, h.Service.GetCarts) }

func (h *CartHandler) AddCart(c *gin.Context) { modifyUserGame(c, h.Service.AddCart) }

func (h *CartHandler) DeleteCart(c *gin.Context) { modifyUserGame(c, h.Service.DeleteCart) }
