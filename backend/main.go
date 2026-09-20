package main

import (
	"software-metrics-project/config"
	"software-metrics-project/handler"
	"software-metrics-project/repository"
	"software-metrics-project/service"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()
	r := gin.Default()
	repo := &repository.AccountRepository{}
	service := &service.AccountService{Repo: repo}
	handler := &handler.AccountHandler{Service: service}
	r.GET("/accounts", handler.GetAccounts)
	r.Run(":8080")
}
