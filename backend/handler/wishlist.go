package handler

import (
	"net/http"
	"software-metrics-project/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WishlistHandler struct {
	Service *service.WishlistService
}

func (h *WishlistHandler) GetWishlists(c *gin.Context) {
	userID := c.Param("userID")

	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	wishlists, err := h.Service.GetWishlists(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, wishlists)
}

func (h *WishlistHandler) AddWishlist(c *gin.Context) {
	convUID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	convGID, err := strconv.Atoi(c.Param("gameID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	output, err := h.Service.AddWishlist(convUID, convGID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, output)
}

func (h *WishlistHandler) DeleteWishlist(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	gameID, err := strconv.Atoi(c.Param("gameID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.Service.DeleteWishlist(userID, gameID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, output)
}
