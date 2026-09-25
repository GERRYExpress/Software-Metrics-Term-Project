package main

import (
	"software-metrics-project/config"
	"software-metrics-project/handler"
	"software-metrics-project/repository"
	"software-metrics-project/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Endpoint interface {
	CreateGameEndpoint(r *gin.RouterGroup)
	CreateAccountEndpoint(r *gin.RouterGroup)
	CreateCategoryEndpoint(r *gin.RouterGroup)
}

func main() {
	config.ConnectDB()
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "X-Total-Count"},
		AllowCredentials: false,
	}))
	api := r.Group("/api")
	{
		// /api/v1
		v1 := api.Group("/v1")
		{
			accounts := v1.Group("/accounts")
			CreateAccountEndpoint(accounts)

			games := v1.Group("/games")
			CreateGameEndpoint(games)

			categories := v1.Group("/categories")
			CreateCategoryEndpoint(categories)
		}
	}
	r.Run(":8080")
}

func CreateAccountEndpoint(r *gin.RouterGroup) {
	repo := &repository.AccountRepository{}
	service := &service.AccountService{Repo: repo}
	handler := &handler.AccountHandler{Service: service}
	r.GET("", handler.GetAccounts)
}

func CreateGameEndpoint(r *gin.RouterGroup) {
	repo := &repository.GameRepository{}
	service := &service.GameService{Repo: repo}
	handler := &handler.GameHandler{Service: service}
	r.GET("", handler.GetGames)
	r.GET("/:id", handler.GetGame)
}

func CreateCategoryEndpoint(r *gin.RouterGroup) {
	repo := &repository.CategoryRepository{}
	service := &service.CategoryService{Repo: repo}
	handler := &handler.CategoryHandler{Service: service}
	r.GET("", handler.GetCategories)
}
