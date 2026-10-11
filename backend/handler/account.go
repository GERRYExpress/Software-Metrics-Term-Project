package handler

import (
	"errors"
	"net/http"
	"software-metrics-project/repository"
	"software-metrics-project/service"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	Service *service.AccountService
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AccountHandler) GetAccounts(c *gin.Context) {
	accounts, err := h.Service.GetAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func (h *AccountHandler) Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, account, err := h.Service.Login(request.Username, request.Password)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "account": account})
}

func (h *AccountHandler) Register(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	account, err := h.Service.Register(request.Username, request.Password)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "username is already registered"})
		return
	}
	c.JSON(http.StatusCreated, account)
}
