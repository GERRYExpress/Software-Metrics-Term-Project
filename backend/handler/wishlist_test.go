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

func setupWishlistTest(t *testing.T) (*gin.Engine, sqlmock.Sqlmock, func()) {
	repo := &repository.WishlistRepository{}
	svc := &service.WishlistService{Repo: repo}
	h := &WishlistHandler{Service: svc}

	return setupHandlerTest(t, func(router *gin.Engine) {
		router.GET("/wishlists/:userID", h.GetWishlists)
		router.POST("/wishlists/:userID/:gameID", h.AddWishlist)
		router.DELETE("/wishlists/:userID/:gameID", h.DeleteWishlist)
	})
}

func TestWishlistHandler_GetWishlists(t *testing.T) {
	router, mock, cleanup := setupWishlistTest(t)
	defer cleanup()

	t.Run("success - 200 OK", func(t *testing.T) {
		expectCollectionListQuery(mock, "wishlists", 1, newCollectionHandlerRow("wishlist_date", 1, 101))

		w := serve(router, http.MethodGet, "/wishlists/1")

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad userID - 400 Bad Request", func(t *testing.T) {
		w := serve(router, http.MethodGet, "/wishlists/abc")

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestWishlistHandler_AddWishlist(t *testing.T) {
	router, mock, cleanup := setupWishlistTest(t)
	defer cleanup()

	t.Run("success - 200 OK", func(t *testing.T) {
		// 1. Mock the DB insertion
		expectCollectionInsert(mock, "wishlists", 1, 42)

		// 2. Mock the SELECT from getWishlist
		expectCollectionGetQuery(mock, "wishlists", 42, 1, newCollectionHandlerRow("wishlist_date", 1, 42))

		w := serve(router, http.MethodPost, "/wishlists/1/42")

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad params - 400 Bad Request", func(t *testing.T) {
		w := serve(router, http.MethodPost, "/wishlists/1/xyz")

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestWishlistHandler_DeleteWishlist(t *testing.T) {
	router, mock, cleanup := setupWishlistTest(t)
	defer cleanup()

	t.Run("success - 200 OK", func(t *testing.T) {
		// 1. Mock the initial SELECT from getWishlist to verify existence
		expectCollectionGetQuery(mock, "wishlists", 42, 1, newCollectionHandlerRow("wishlist_date", 1, 42))

		// 2. Mock the DELETE execution
		expectCollectionDelete(mock, "wishlists", 42, 1)

		w := serve(router, http.MethodDelete, "/wishlists/1/42")

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
