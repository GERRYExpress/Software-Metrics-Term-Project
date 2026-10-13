package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"software-metrics-project/config"
	"software-metrics-project/model"
	"software-metrics-project/repository/repotestutil"
)

func TestCartRepository_GetCarts(t *testing.T) {
	db, mock, cleanup := repotestutil.NewMockDB(t)
	defer cleanup()

	originalDB := config.DB
	config.DB = db
	defer func() { config.DB = originalDB }()

	repo := &CartRepository{}

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"ca.user_id", "ca.product_id", "ca.cart_date",
			"p.name", "p.price", "p.poster_url", "c.category_id", "c.name",
		}).
			AddRow(1, 1, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				"Game 1", 29.99, "poster1.jpg", 1, "Action").
			AddRow(1, 1, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				"Game 1", 29.99, "poster1.jpg", 2, "RPG").
			AddRow(1, 2, time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC),
				"Game 2", 19.99, "poster2.jpg", 3, "Strategy")

		mock.ExpectQuery("SELECT.*FROM.*carts.*").
			WillReturnRows(rows)

		carts, err := repo.GetCarts(1)
		require.NoError(t, err)
		assert.Len(t, carts, 2)
		assert.Equal(t, model.NewID(1), carts[0].ProductID)
		assert.Equal(t, model.NewName("Game 1"), carts[0].Game.Name)
		assert.Equal(t, model.NewPrice(29.99), carts[0].Game.Price)
		assert.Len(t, carts[0].Game.Categories, 2)
		assert.Equal(t, model.NewName("Action"), carts[0].Game.Categories[0].Name)
		assert.Equal(t, model.NewName("RPG"), carts[0].Game.Categories[1].Name)
		assert.Equal(t, model.NewID(2), carts[1].ProductID)
		assert.Equal(t, model.NewName("Game 2"), carts[1].Game.Name)
	})

	t.Run("empty result", func(t *testing.T) {
		mock.ExpectQuery("SELECT.*FROM.*carts.*").
			WillReturnRows(sqlmock.NewRows([]string{
				"ca.user_id", "ca.product_id", "ca.cart_date",
				"p.name", "p.price", "p.poster_url", "c.category_id", "c.name",
			}))

		carts, err := repo.GetCarts(999)
		require.NoError(t, err)
		assert.Len(t, carts, 0)
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery("SELECT.*FROM.*carts.*").
			WillReturnError(sql.ErrConnDone)

		_, err := repo.GetCarts(1)
		assert.Error(t, err)
	})
}

func TestCartRepository_AddCart(t *testing.T) {
	db, mock, cleanup := repotestutil.NewMockDB(t)
	defer cleanup()

	originalDB := config.DB
	config.DB = db
	defer func() { config.DB = originalDB }()

	repo := &CartRepository{}

	t.Run("success - new cart item", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO.*carts.*ON DUPLICATE KEY UPDATE").
			WithArgs(1, 1, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// getCart is called after insert
		mock.ExpectQuery("SELECT.*FROM.*carts.*").
			WillReturnRows(sqlmock.NewRows([]string{
				"ca.user_id", "ca.product_id", "ca.cart_date",
				"p.name", "p.price", "p.poster_url", "c.category_id", "c.name",
			}).AddRow(1, 1, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), "Game 1", 29.99, "poster.jpg", 1, "Action"))

		cart, err := repo.AddCart(1, 1)
		require.NoError(t, err)
		assert.NotNil(t, cart)
		assert.Equal(t, model.NewID(1), cart.ProductID)
	})

	t.Run("success - duplicate cart item (upsert)", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO.*carts.*ON DUPLICATE KEY UPDATE").
			WithArgs(1, 2, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectQuery("SELECT.*FROM.*carts.*").
			WillReturnRows(sqlmock.NewRows([]string{
				"ca.user_id", "ca.product_id", "ca.cart_date",
				"p.name", "p.price", "p.poster_url", "c.category_id", "c.name",
			}).AddRow(1, 2, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), "Game 2", 19.99, "poster2.jpg", 1, "RPG"))

		cart, err := repo.AddCart(1, 2)
		require.NoError(t, err)
		assert.NotNil(t, cart)
		assert.Equal(t, model.NewID(2), cart.ProductID)
	})

	t.Run("db error on insert", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO.*carts.*ON DUPLICATE KEY UPDATE").
			WithArgs(1, 3, sqlmock.AnyArg()).
			WillReturnError(sql.ErrConnDone)

		_, err := repo.AddCart(1, 3)
		assert.Error(t, err)
	})

	t.Run("db error on getCart after insert", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO.*carts.*ON DUPLICATE KEY UPDATE").
			WithArgs(1, 4, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectQuery("SELECT.*FROM.*carts.*").
			WillReturnError(sql.ErrConnDone)

		_, err := repo.AddCart(1, 4)
		assert.Error(t, err)
	})
}

func TestCartRepository_DeleteCart(t *testing.T) {
	db, mock, cleanup := repotestutil.NewMockDB(t)
	defer cleanup()

	originalDB := config.DB
	config.DB = db
	defer func() { config.DB = originalDB }()

	repo := &CartRepository{}

	t.Run("success", func(t *testing.T) {
		// getCart is called first to get the cart item
		mock.ExpectQuery("SELECT.*FROM.*carts.*").
			WillReturnRows(sqlmock.NewRows([]string{
				"ca.user_id", "ca.product_id", "ca.cart_date",
				"p.name", "p.price", "p.poster_url", "c.category_id", "c.name",
			}).AddRow(1, 1, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), "Game 1", 29.99, "poster.jpg", 1, "Action"))

		mock.ExpectExec("DELETE FROM.*carts.*").
			WithArgs(1, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		cart, err := repo.DeleteCart(1, 1)
		require.NoError(t, err)
		assert.NotNil(t, cart)
		assert.Equal(t, model.NewID(1), cart.ProductID)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT.*FROM.*carts.*").
			WillReturnError(sql.ErrNoRows)

		_, err := repo.DeleteCart(999, 999)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("db error on getCart", func(t *testing.T) {
		mock.ExpectQuery("SELECT.*FROM.*carts.*").
			WillReturnError(sql.ErrConnDone)

		_, err := repo.DeleteCart(1, 1)
		assert.Error(t, err)
	})

	t.Run("db error on delete", func(t *testing.T) {
		mock.ExpectQuery("SELECT.*FROM.*carts.*").
			WillReturnRows(sqlmock.NewRows([]string{
				"ca.user_id", "ca.product_id", "ca.cart_date",
				"p.name", "p.price", "p.poster_url", "c.category_id", "c.name",
			}).AddRow(1, 1, "2024-01-15T00:00:00Z", "Game 1", 29.99, "poster.jpg", 1, "Action"))

		mock.ExpectExec("DELETE FROM.*carts.*").
			WithArgs(1, 1).
			WillReturnError(sql.ErrConnDone)

		_, err := repo.DeleteCart(1, 1)
		assert.Error(t, err)
	})
}