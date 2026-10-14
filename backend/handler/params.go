package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func paramUserID(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 0, false
	}
	return id, true
}

func paramUserGameIDs(c *gin.Context) (int, int, bool) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 0, 0, false
	}
	gameID, err := strconv.Atoi(c.Param("gameID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 0, 0, false
	}
	return userID, gameID, true
}

func listByUserID[T any](c *gin.Context, fetch func(int) (T, error)) {
	id, ok := paramUserID(c)
	if !ok {
		return
	}
	result, err := fetch(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func modifyUserGame[T any](c *gin.Context, action func(int, int) (T, error)) {
	userID, gameID, ok := paramUserGameIDs(c)
	if !ok {
		return
	}
	result, err := action(userID, gameID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
