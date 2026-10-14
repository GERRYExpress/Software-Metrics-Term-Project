package handler

import (
	"software-metrics-project/service"

	"github.com/gin-gonic/gin"
)

type WishlistHandler struct {
	Service *service.WishlistService
}

func (h *WishlistHandler) GetWishlists(c *gin.Context) { listByUserID(c, h.Service.GetWishlists) }

func (h *WishlistHandler) AddWishlist(c *gin.Context) { modifyUserGame(c, h.Service.AddWishlist) }

func (h *WishlistHandler) DeleteWishlist(c *gin.Context) { modifyUserGame(c, h.Service.DeleteWishlist) }
