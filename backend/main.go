package main

import (
	"log"
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
	if err := config.LoadEnv(".env"); err != nil {
		log.Fatal(err)
	}
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

			wishlists := v1.Group("/wishlists")
			CreateWishlistEndpoint(wishlists)

			carts := v1.Group("/carts")
			CreateCartEndpoint(carts)

			orders := v1.Group("orders")
			CreateOrderEndpoint(orders)
		}
	}
	r.Run(":8080")
}

func CreateAccountEndpoint(r *gin.RouterGroup) {
	repo := &repository.AccountRepository{}
	service := &service.AccountService{Repo: repo}
	handler := &handler.AccountHandler{Service: service}
	r.GET("", handler.GetAccounts)
	r.POST("", handler.Register)
	r.POST("/login", handler.Login)
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
	r.POST("", handler.Create)
}

func CreateWishlistEndpoint(r *gin.RouterGroup) {
	repo := &repository.WishlistRepository{}
	service := &service.WishlistService{Repo: repo}
	handler := &handler.WishlistHandler{Service: service}
	operationPath := "/:userID/:gameID"
	r.GET("/:userID", handler.GetWishlists)
	r.POST(operationPath, handler.AddWishlist)
	r.DELETE(operationPath, handler.DeleteWishlist)
}

func CreateCartEndpoint(r *gin.RouterGroup) {
	repo := &repository.CartRepository{}
	service := &service.CartService{Repo: repo}
	handler := &handler.CartHandler{Service: service}
	operationPath := "/:userID/:gameID"
	r.GET("/:userID", handler.GetCarts)
	r.POST(operationPath, handler.AddCart)
	r.DELETE(operationPath, handler.DeleteCart)
}

func CreateOrderEndpoint(r *gin.RouterGroup) {
	repo := &repository.OrderRepository{}
	service := &service.OrderService{Repo: repo}
	handler := &handler.OrderHandler{Service: service}
	r.POST("", handler.Create)
}
