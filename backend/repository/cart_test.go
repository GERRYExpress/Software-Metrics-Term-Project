package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"software-metrics-project/config"
	"software-metrics-project/repository/repotestutil"
)

func setupCartRepoTest(t *testing.T) (*CartRepository, sqlmock.Sqlmock, func()) {
	db, mock, cleanup := repotestutil.NewMockDB(t)
	originalDB := config.DB
	config.DB = db

	repo := &CartRepository{}
	return repo, mock, func() {
		config.DB = originalDB
		cleanup()
	}
}

func TestCartRepository_GetCarts(t *testing.T) {
	repo, mock, cleanup := setupCartRepoTest(t)
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.NewRows([]string{
			"user_id", "product_id", "cart_date", "name", "price", "poster_url", "category_id", "category_name",
		}).
			AddRow(1, 101, now, "Game 1", 29.99, "url1", 1, "Action").
			AddRow(1, 101, now, "Game 1", 29.99, "url1", 2, "RPG") // Same game, second category

		mock.ExpectQuery("^SELECT(.*)FROM carts").WithArgs(1).WillReturnRows(rows)

		carts, err := repo.GetCarts(1)

		assert.NoError(t, err)
		assert.Len(t, carts, 1) // Should group by product_id
		assert.Equal(t, 101, int(carts[0].ProductID))
		assert.Len(t, carts[0].Game.Categories, 2)
	})
}

func TestCartRepository_getCart(t *testing.T) {
	repo, mock, cleanup := setupCartRepoTest(t)
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.NewRows([]string{
			"user_id", "product_id", "cart_date", "name", "price", "poster_url", "category_id", "category_name",
		}).AddRow(1, 101, now, "Game 1", 29.99, "url1", 1, "Action")

		mock.ExpectQuery("^SELECT(.*)FROM carts").WithArgs(101, 1).WillReturnRows(rows)

		cart, err := repo.getCart(1, 101)

		assert.NoError(t, err)
		assert.NotNil(t, cart)
		assert.Equal(t, 1, int(cart.UserID))
		assert.Equal(t, 101, int(cart.ProductID))
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("^SELECT(.*)FROM carts").WithArgs(1, 999).WillReturnRows(sqlmock.NewRows([]string{}))

		_, err := repo.getCart(1, 999)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func TestCartRepository_AddCart(t *testing.T) {
	repo, mock, cleanup := setupCartRepoTest(t)
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		// 1. Expect the Insert
		mock.ExpectExec("^INSERT INTO carts").
			WithArgs(1, 101, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// 2. Expect the getCart Select that immediately follows
		now := time.Now()
		rows := sqlmock.NewRows([]string{
			"user_id", "product_id", "cart_date", "name", "price", "poster_url", "category_id", "category_name",
		}).AddRow(1, 101, now, "Game 1", 29.99, "url1", nil, nil)

		mock.ExpectQuery("^SELECT(.*)FROM carts").WithArgs(1, 101).WillReturnRows(rows)

		cart, err := repo.AddCart(1, 101)
		assert.NoError(t, err)
		assert.NotNil(t, cart)
		assert.Equal(t, 101, int(cart.ProductID))
	})
}

func TestCartRepository_DeleteCart(t *testing.T) {
	repo, mock, cleanup := setupCartRepoTest(t)
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		// 1. Expect getCart to verify it exists and return the deleted object
		now := time.Now()
		rows := sqlmock.NewRows([]string{
			"user_id", "product_id", "cart_date", "name", "price", "poster_url", "category_id", "category_name",
		}).AddRow(1, 101, now, "Game 1", 29.99, "url1", nil, nil)

		mock.ExpectQuery("^SELECT(.*)FROM carts").WithArgs(1, 101).WillReturnRows(rows)

		// 2. Expect the actual Delete execution
		// Note: The original code uses undefined variables `cartUserFilter` and `cartProductFilter`
		// inside the Delete query block. Assuming they map to user_id and product_id based on squirrel syntax.
		mock.ExpectExec("^DELETE FROM carts").WithArgs(1, 101).WillReturnResult(sqlmock.NewResult(0, 1))

		cart, err := repo.DeleteCart(1, 101)

		assert.NoError(t, err)
		assert.NotNil(t, cart)
		assert.Equal(t, 101, int(cart.ProductID))
	})

	t.Run("fails if cart does not exist", func(t *testing.T) {
		// getCart returns nothing
		mock.ExpectQuery("^SELECT(.*)FROM carts").WithArgs(1, 999).WillReturnRows(sqlmock.NewRows([]string{}))

		_, err := repo.DeleteCart(1, 999)

		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}
