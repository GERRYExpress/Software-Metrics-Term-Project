package handler

import (
	"net/http"
	"software-metrics-project/service"

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
