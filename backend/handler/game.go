package handler

import (
	"net/http"
	"software-metrics-project/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	Service *service.GameService
}

func (h *GameHandler) GetGames(c *gin.Context) {
	games, err := h.Service.GetGames()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, games)
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
