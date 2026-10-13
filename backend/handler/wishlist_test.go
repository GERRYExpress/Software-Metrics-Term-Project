package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"software-metrics-project/config"
	"software-metrics-project/repository"
	"software-metrics-project/repository/repotestutil"
	"software-metrics-project/service"
)

func setupWishlistTest(t *testing.T) (*gin.Engine, sqlmock.Sqlmock, func()) {
	db, mock, cleanup := repotestutil.NewMockDB(t)

	originalDB := config.DB
	config.DB = db

	repo := &repository.WishlistRepository{}
	svc := &service.WishlistService{Repo: repo}
	h := &WishlistHandler{Service: svc}

	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/wishlists/:userID", h.GetWishlists)
	router.POST("/wishlists/:userID/:gameID", h.AddWishlist)
	router.DELETE("/wishlists/:userID/:gameID", h.DeleteWishlist)

	fullCleanup := func() {
		config.DB = originalDB
		cleanup()
	}

	return router, mock, fullCleanup
}

func TestWishlistHandler_GetWishlists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, mock, cleanup := setupWishlistTest(t)
	defer cleanup()

	t.Run("success - 200 OK", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "game_id"}).AddRow(1, 101)
		mock.ExpectQuery("SELECT(.*)FROM").WithArgs(1).WillReturnRows(rows)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/wishlists/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad userID - 400 Bad Request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/wishlists/abc", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestWishlistHandler_AddWishlist(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, mock, cleanup := setupWishlistTest(t)
	defer cleanup()

	t.Run("success - 200 OK", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO").WithArgs(1, 42).WillReturnResult(sqlmock.NewResult(1, 1))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/wishlists/1/42", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad params - 400 Bad Request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/wishlists/1/xyz", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestWishlistHandler_DeleteWishlist(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, mock, cleanup := setupWishlistTest(t)
	defer cleanup()

	t.Run("success - 200 OK", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM").WithArgs(1, 42).WillReturnResult(sqlmock.NewResult(0, 1))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/wishlists/1/42", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
