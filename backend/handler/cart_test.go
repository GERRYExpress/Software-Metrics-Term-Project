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

// setupCartTest wires up the mock DB, repository, service, and a Gin router for testing.
func setupCartTest(t *testing.T) (*gin.Engine, sqlmock.Sqlmock, func()) {
	db, mock, cleanup := repotestutil.NewMockDB(t)

	// Override global DB for the repository layer
	originalDB := config.DB
	config.DB = db

	repo := &repository.CartRepository{}
	svc := &service.CartService{Repo: repo}
	h := &CartHandler{Service: svc}

	router := gin.New()
	router.Use(gin.Recovery())

	// Map the routes exactly as expected by the handler
	router.GET("/carts/:userID", h.GetCarts)
	router.POST("/carts/:userID/:gameID", h.AddCart)
	router.DELETE("/carts/:userID/:gameID", h.DeleteCart)

	fullCleanup := func() {
		config.DB = originalDB
		cleanup()
	}

	return router, mock, fullCleanup
}

func TestCartHandler_GetCarts(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router, mock, cleanup := setupCartTest(t)
	defer cleanup()

	t.Run("success - 200 OK", func(t *testing.T) {
		// Mock the DB query expected by service.GetCarts(id)
		rows := sqlmock.NewRows([]string{"id", "game_id"}).AddRow(1, 101)
		mock.ExpectQuery("SELECT(.*)FROM").WithArgs(1).WillReturnRows(rows)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/carts/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad userID - 400 Bad Request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/carts/invalid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid syntax")
	})
}

func TestCartHandler_AddCart(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router, mock, cleanup := setupCartTest(t)
	defer cleanup()

	t.Run("success - 200 OK", func(t *testing.T) {
		// Mock the DB insertion expected by service.AddCart
		mock.ExpectExec("INSERT INTO").WithArgs(1, 42).WillReturnResult(sqlmock.NewResult(1, 1))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/carts/1/42", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad gameID - 400 Bad Request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/carts/1/abc", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCartHandler_DeleteCart(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router, mock, cleanup := setupCartTest(t)
	defer cleanup()

	t.Run("success - 200 OK", func(t *testing.T) {
		// Mock the DB deletion expected by service.DeleteCart
		mock.ExpectExec("DELETE FROM").WithArgs(1, 42).WillReturnResult(sqlmock.NewResult(0, 1))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/carts/1/42", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad userID - 400 Bad Request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/carts/xyz/42", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
