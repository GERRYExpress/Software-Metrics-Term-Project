package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"software-metrics-project/repository"
	"software-metrics-project/service"
	"strconv"

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
	h.Create(c)
}

func (h *AccountHandler) Create(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	account, err := h.Service.Create(request.Username, request.Password)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "username is already registered"})
		return
	}
	c.JSON(http.StatusCreated, account)
}

func (h *AccountHandler) Delete(c *gin.Context) {
	userID, err := positiveID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Service.Delete(userID); err != nil {
		switch {
		case errors.Is(err, repository.ErrAccountHasGames):
			c.JSON(http.StatusConflict, gin.H{"error": "delete your games before deleting your account"})
		case errors.Is(err, sql.ErrNoRows):
			c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.Status(http.StatusNoContent)
}

func positiveID(raw string) (int, error) {
	id, err := strconv.Atoi(raw)
	if err != nil || id < 1 {
		return 0, errors.New("id must be a positive integer")
	}
	return id, nil
}
