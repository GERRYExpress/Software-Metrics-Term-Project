package handler

import (
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"software-metrics-project/repository"
	"software-metrics-project/service"
)

// setupCartTest wires up the repository, service, and routes for testing.
// Mock DB swap and router boilerplate live in setupHandlerTest.
func setupCartTest(t *testing.T) (*gin.Engine, sqlmock.Sqlmock, func()) {
	repo := &repository.CartRepository{}
	svc := &service.CartService{Repo: repo}
	h := &CartHandler{Service: svc}

	return setupHandlerTest(t, func(router *gin.Engine) {
		router.GET("/carts/:userID", h.GetCarts)
		router.POST("/carts/:userID/:gameID", h.AddCart)
		router.DELETE("/carts/:userID/:gameID", h.DeleteCart)
	})
}

func TestCartHandler_GetCarts(t *testing.T) {
	router, mock, cleanup := setupCartTest(t)
	defer cleanup()

	t.Run("success - 200 OK", func(t *testing.T) {
		expectCollectionListQuery(mock, "carts", 1, newCollectionHandlerRow("cart_date", 1, 101))

		w := serve(router, http.MethodGet, "/carts/1")

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad userID - 400 Bad Request", func(t *testing.T) {
		w := serve(router, http.MethodGet, "/carts/invalid")

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid syntax")
	})
}

func TestCartHandler_AddCart(t *testing.T) {
	router, mock, cleanup := setupCartTest(t)
	defer cleanup()

	t.Run("success - 200 OK", func(t *testing.T) {
		// 1. Mock the DB insertion expected by service.AddCart
		expectCollectionInsert(mock, "carts", 1, 42)

		// 2. Mock the SELECT triggered by the internal getCart call
		expectCollectionGetQuery(mock, "carts", 42, 1, newCollectionHandlerRow("cart_date", 1, 42))

		w := serve(router, http.MethodPost, "/carts/1/42")

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad gameID - 400 Bad Request", func(t *testing.T) {
		w := serve(router, http.MethodPost, "/carts/1/abc")

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCartHandler_DeleteCart(t *testing.T) {
	router, mock, cleanup := setupCartTest(t)
	defer cleanup()

	t.Run("success - 200 OK", func(t *testing.T) {
		// 1. Mock the initial SELECT triggered by getCart to verify existence
		expectCollectionGetQuery(mock, "carts", 42, 1, newCollectionHandlerRow("cart_date", 1, 42))

		// 2. Mock the actual DELETE execution
		expectCollectionDelete(mock, "carts", 42, 1)

		w := serve(router, http.MethodDelete, "/carts/1/42")

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad userID - 400 Bad Request", func(t *testing.T) {
		w := serve(router, http.MethodDelete, "/carts/xyz/42")

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
