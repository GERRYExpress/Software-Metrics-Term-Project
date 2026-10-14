package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

// getWishlistColumns provides the exact 8 columns expected by the repository's Scan.
func getWishlistColumns() []string {
	return []string{
		"user_id", "product_id", "wishlist_date", "name", "price", "poster_url", "category_id", "category_name",
	}
}

func TestWishlistHandler_GetWishlists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, mock, cleanup := setupWishlistTest(t)
	defer cleanup()

	t.Run("success - 200 OK", func(t *testing.T) {
		rows := sqlmock.NewRows(getWishlistColumns()).AddRow(1, 101, time.Now(), "Game 1", 29.99, "url1", nil, nil)
		mock.ExpectQuery("SELECT(.*)FROM wishlists").WithArgs(1).WillReturnRows(rows)

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
		// 1. Mock the DB insertion
		mock.ExpectExec("INSERT INTO wishlists").
			WithArgs(1, 42, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// 2. Mock the SELECT from getWishlist. Args: (gameID, userID) due to squirrel map sorting
		rows := sqlmock.NewRows(getWishlistColumns()).AddRow(1, 42, time.Now(), "Game 1", 29.99, "url1", nil, nil)
		mock.ExpectQuery("SELECT(.*)FROM wishlists").WithArgs(42, 1).WillReturnRows(rows)

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
		// 1. Mock the initial SELECT from getWishlist to verify existence
		rows := sqlmock.NewRows(getWishlistColumns()).AddRow(1, 42, time.Now(), "Game 1", 29.99, "url1", nil, nil)
		mock.ExpectQuery("SELECT(.*)FROM wishlists").WithArgs(42, 1).WillReturnRows(rows)

		// 2. Mock the DELETE execution. Args: (gameID, userID)
		mock.ExpectExec("DELETE FROM wishlists").WithArgs(42, 1).WillReturnResult(sqlmock.NewResult(0, 1))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/wishlists/1/42", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
