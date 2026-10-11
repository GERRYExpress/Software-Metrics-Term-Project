package handler

import (
	"fmt"
	"net/http"
	"software-metrics-project/model"
	"software-metrics-project/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	Service *service.GameService
}

func (h *GameHandler) GetGames(c *gin.Context) {
	var (
		games []model.Game
		err   error
	)
	limit, page, err := pagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	keyword := c.Query("keyword")
	category := c.Query("category")
	var sellerID *int
	if rawSellerID := c.Query("sellerId"); rawSellerID != "" {
		parsedSellerID, parseErr := strconv.Atoi(rawSellerID)
		if parseErr != nil || parsedSellerID < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "sellerId must be a positive integer"})
			return
		}
		sellerID = &parsedSellerID
	}
	if keyword != "" {
		games, err = h.Service.GetGamesByName(keyword, category, sellerID, limit, page)
	} else {
		games, err = h.Service.GetGames(category, sellerID, limit, page)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	totalCount, err := h.Service.GetGamesCount(keyword, category, sellerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("X-Total-Count", strconv.Itoa(totalCount))
	c.JSON(http.StatusOK, games)
}

func pagination(c *gin.Context) (int, int, error) {
	limit, err := strconv.Atoi(c.DefaultQuery("_limit", "10"))
	if err != nil || limit < 1 {
		return 0, 0, fmt.Errorf("_limit must be a positive integer")
	}

	page, err := strconv.Atoi(c.DefaultQuery("_page", "1"))
	if err != nil || page < 1 {
		return 0, 0, fmt.Errorf("_page must be a positive integer")
	}

	return limit, page, nil
}

func (h *GameHandler) GetGame(c *gin.Context) {
	gameID := c.Param("id")

	id, err := strconv.Atoi(gameID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	game, err := h.Service.GetGame(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, game)
}
