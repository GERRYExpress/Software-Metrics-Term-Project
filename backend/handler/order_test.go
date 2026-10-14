package handler

import (
	"bytes"
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

	t.Run("success - 201 Created", func(t *testing.T) {
		mock.ExpectBegin()

		// The repo inserts the base order
		mock.ExpectExec("INSERT INTO orders").
			WithArgs(1, "Credit Card", 59.99).
			WillReturnResult(sqlmock.NewResult(10, 1))

		// Because there are no items in our JSON payload, it skips inserting into order_items
		mock.ExpectCommit()

		// Mock the read-back query triggered by getOrder(orderID)
		orderRows := sqlmock.NewRows([]string{
			"order_id", "user_id", "order_date", "payment_method", "final_price",
			"product_id", "name", "description", "price", "poster_url", "release_date",
			"category_id", "category_name", "is_hidden",
		}).AddRow(10, 1, time.Now(), "Credit Card", 59.99, 101, "Game 1", "Desc", 59.99, "url", time.Now(), nil, nil, nil)

		mock.ExpectQuery("SELECT(.*)FROM orders o").
			WithArgs(10).
			WillReturnRows(orderRows)

		validJSON := `{"userId": 1, "paymentMethod": "Credit Card", "finalPrice": 59.99, "games": []}`
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString(validJSON))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})
}
