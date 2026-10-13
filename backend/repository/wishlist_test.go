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

func TestWishlistRepository_GetWishlists(t *testing.T) {
	db, mock, cleanup := repotestutil.NewMockDB(t)
	defer cleanup()

	originalDB := config.DB
	config.DB = db
	defer func() { config.DB = originalDB }()

	repo := &WishlistRepository{}

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"w.user_id", "w.product_id", "w.wishlist_date",
			"p.name", "p.price", "p.poster_url", "c.category_id", "c.name",
		}).
			AddRow(1, 1, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				"Game 1", 29.99, "poster1.jpg", 1, "Action").
			AddRow(1, 1, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				"Game 1", 29.99, "poster1.jpg", 2, "RPG").
			AddRow(1, 2, time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC),
				"Game 2", 19.99, "poster2.jpg", 3, "Strategy")

		mock.ExpectQuery("SELECT.*FROM.*wishlists.*").
			WillReturnRows(rows)

		wishlists, err := repo.GetWishlists(1)
		require.NoError(t, err)
		assert.Len(t, wishlists, 2)
		assert.Equal(t, model.NewID(1), wishlists[0].ProductID)
		assert.Equal(t, model.NewName("Game 1"), wishlists[0].Game.Name)
		assert.Equal(t, model.NewPrice(29.99), wishlists[0].Game.Price)
		assert.Len(t, wishlists[0].Game.Categories, 2)
		assert.Equal(t, model.NewName("Action"), wishlists[0].Game.Categories[0].Name)
		assert.Equal(t, model.NewName("RPG"), wishlists[0].Game.Categories[1].Name)
		assert.Equal(t, model.NewID(2), wishlists[1].ProductID)
		assert.Equal(t, model.NewName("Game 2"), wishlists[1].Game.Name)
	})

	t.Run("empty result", func(t *testing.T) {
		mock.ExpectQuery("SELECT.*FROM.*wishlists.*").
			WillReturnRows(sqlmock.NewRows([]string{
				"w.user_id", "w.product_id", "w.wishlist_date",
				"p.name", "p.price", "p.poster_url", "c.category_id", "c.name",
			}))

		wishlists, err := repo.GetWishlists(999)
		require.NoError(t, err)
		assert.Len(t, wishlists, 0)
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery("SELECT.*FROM.*wishlists.*").
			WillReturnError(sql.ErrConnDone)

		_, err := repo.GetWishlists(1)
		assert.Error(t, err)
	})
}

func TestWishlistRepository_AddWishlist(t *testing.T) {
	db, mock, cleanup := repotestutil.NewMockDB(t)
	defer cleanup()

	originalDB := config.DB
	config.DB = db
	defer func() { config.DB = originalDB }()

	repo := &WishlistRepository{}

	t.Run("success - new wishlist item", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO.*wishlists.*ON DUPLICATE KEY UPDATE").
			WithArgs(1, 1, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectQuery("SELECT.*FROM.*wishlists.*").
			WillReturnRows(sqlmock.NewRows([]string{
				"w.user_id", "w.product_id", "w.wishlist_date",
				"p.name", "p.price", "p.poster_url", "c.category_id", "c.name",
			}).AddRow(1, 1, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), "Game 1", 29.99, "poster.jpg", 1, "Action"))

		wishlist, err := repo.AddWishlist(1, 1)
		require.NoError(t, err)
		assert.NotNil(t, wishlist)
		assert.Equal(t, model.NewID(1), wishlist.ProductID)
	})

	t.Run("success - duplicate wishlist item (upsert)", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO.*wishlists.*ON DUPLICATE KEY UPDATE").
			WithArgs(1, 2, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectQuery("SELECT.*FROM.*wishlists.*").
			WillReturnRows(sqlmock.NewRows([]string{
				"w.user_id", "w.product_id", "w.wishlist_date",
				"p.name", "p.price", "p.poster_url", "c.category_id", "c.name",
			}).AddRow(1, 2, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), "Game 2", 19.99, "poster2.jpg", 1, "RPG"))

		wishlist, err := repo.AddWishlist(1, 2)
		require.NoError(t, err)
		assert.NotNil(t, wishlist)
		assert.Equal(t, model.NewID(2), wishlist.ProductID)
	})

	t.Run("db error on insert", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO.*wishlists.*ON DUPLICATE KEY UPDATE").
			WithArgs(1, 3, sqlmock.AnyArg()).
			WillReturnError(sql.ErrConnDone)

		_, err := repo.AddWishlist(1, 3)
		assert.Error(t, err)
	})

	t.Run("db error on getWishlist after insert", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO.*wishlists.*ON DUPLICATE KEY UPDATE").
			WithArgs(1, 4, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectQuery("SELECT.*FROM.*wishlists.*").
			WillReturnError(sql.ErrConnDone)

		_, err := repo.AddWishlist(1, 4)
		assert.Error(t, err)
	})
}

func TestWishlistRepository_DeleteWishlist(t *testing.T) {
	db, mock, cleanup := repotestutil.NewMockDB(t)
	defer cleanup()

	originalDB := config.DB
	config.DB = db
	defer func() { config.DB = originalDB }()

	repo := &WishlistRepository{}

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT.*FROM.*wishlists.*").
			WillReturnRows(sqlmock.NewRows([]string{
				"w.user_id", "w.product_id", "w.wishlist_date",
				"p.name", "p.price", "p.poster_url", "c.category_id", "c.name",
			}).AddRow(1, 1, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), "Game 1", 29.99, "poster.jpg", 1, "Action"))

		mock.ExpectExec("DELETE FROM.*wishlists.*").
			WithArgs(1, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		wishlist, err := repo.DeleteWishlist(1, 1)
		require.NoError(t, err)
		assert.NotNil(t, wishlist)
		assert.Equal(t, model.NewID(1), wishlist.ProductID)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT.*FROM.*wishlists.*").
			WillReturnError(sql.ErrNoRows)

		_, err := repo.DeleteWishlist(999, 999)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("db error on getWishlist", func(t *testing.T) {
		mock.ExpectQuery("SELECT.*FROM.*wishlists.*").
			WillReturnError(sql.ErrConnDone)

		_, err := repo.DeleteWishlist(1, 1)
		assert.Error(t, err)
	})

	t.Run("db error on delete", func(t *testing.T) {
		mock.ExpectQuery("SELECT.*FROM.*wishlists.*").
			WillReturnRows(sqlmock.NewRows([]string{
				"w.user_id", "w.product_id", "w.wishlist_date",
				"p.name", "p.price", "p.poster_url", "c.category_id", "c.name",
			}).AddRow(1, 1, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), "Game 1", 29.99, "poster.jpg", 1, "Action"))

		mock.ExpectExec("DELETE FROM.*wishlists.*").
			WithArgs(1, 1).
			WillReturnError(sql.ErrConnDone)

		_, err := repo.DeleteWishlist(1, 1)
		assert.Error(t, err)
	})
}