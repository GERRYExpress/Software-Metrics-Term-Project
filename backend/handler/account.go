package handler

import (
	"net/http"
	"software-metrics-project/service"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	Service *service.AccountService
}

func (h *AccountHandler) GetAccounts(c *gin.Context) {
	accounts, err := h.Service.GetAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, accounts)
}
