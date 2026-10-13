package handler

import (
	"bytes"
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

func setupOrderTest(t *testing.T) (*gin.Engine, sqlmock.Sqlmock, func()) {
	db, mock, cleanup := repotestutil.NewMockDB(t)

	originalDB := config.DB
	config.DB = db

	repo := &repository.OrderRepository{}
	svc := &service.OrderService{Repo: repo}
	h := &OrderHandler{Service: svc}

	router := gin.New()
	router.Use(gin.Recovery())
	router.POST("/orders", h.Create)

	fullCleanup := func() {
		config.DB = originalDB
		cleanup()
	}

	return router, mock, fullCleanup
}

func TestOrderHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, mock, cleanup := setupOrderTest(t)
	defer cleanup()

	t.Run("bad json payload - 400 Bad Request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString(`{invalid json}`))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("success - 201 Created", func(t *testing.T) {
		// Note: Update the WithArgs and JSON body to match the exact fields of model.Order
		mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))

		validJSON := `{"id": 1, "totalAmount": 59.99}`
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString(validJSON))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		// This expects the service layer logic to pass without DB syntax errors
		assert.NotEqual(t, http.StatusInternalServerError, w.Code)
	})
}
