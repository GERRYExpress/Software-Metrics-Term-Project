package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"software-metrics-project/config"
	"software-metrics-project/model"
	"software-metrics-project/repository/repotestutil"
)

func setupOrderRepoTest(t *testing.T) (*OrderRepository, sqlmock.Sqlmock, func()) {
	db, mock, cleanup := repotestutil.NewMockDB(t)
	originalDB := config.DB
	config.DB = db

	repo := &OrderRepository{}
	return repo, mock, func() {
		config.DB = originalDB
		cleanup()
	}
}

func TestOrderRepository_Create(t *testing.T) {
	repo, mock, cleanup := setupOrderRepoTest(t)
	defer cleanup()

	t.Run("success with transaction", func(t *testing.T) {
		order := model.Order{
			UserID:        1,
			PaymentMethod: "Credit Card",
			FinalPrice:    59.99,
			Games: []model.Game{
				{ID: 101},
				{ID: 102},
			},
		}

		// 1. Expect Transaction Begin
		mock.ExpectBegin()

		// 2. Expect Order Insert
		mock.ExpectExec("^INSERT INTO orders").
			WithArgs(1, "Credit Card", 59.99).
			WillReturnResult(sqlmock.NewResult(10, 1)) // Returns OrderID 10

		// 3. Expect Game Items Insert (Loop)
		mock.ExpectExec("^INSERT INTO order_items").WithArgs(10, 101).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("^INSERT INTO order_items").WithArgs(10, 102).WillReturnResult(sqlmock.NewResult(2, 1))

		// 4. Expect Transaction Commit
		mock.ExpectCommit()

		// 5. Expect the final SELECT query from getOrder
		rows := sqlmock.NewRows([]string{
			"order_id", "user_id", "order_date", "payment_method", "final_price",
			"product_id", "name", "description", "price", "poster_url", "release_date",
			"category_id", "category_name", "is_hidden",
		}).AddRow(10, 1, time.Now(), "Credit Card", 59.99, 101, "Game 1", "Desc", 29.99, "url", time.Now(), nil, nil, nil)

		mock.ExpectQuery("^SELECT(.*)FROM orders").WithArgs(10).WillReturnRows(rows)

		createdOrder, err := repo.Create(order)
		assert.NoError(t, err)
		assert.NotNil(t, createdOrder)
		assert.Equal(t, 10, int(createdOrder.ID))
	})

	t.Run("rollback on insert failure", func(t *testing.T) {
		order := model.Order{UserID: 1}

		mock.ExpectBegin()
		mock.ExpectExec("^INSERT INTO orders").WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		_, err := repo.Create(order)
		assert.Error(t, err)
	})
}
