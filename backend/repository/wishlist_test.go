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

func setupWishlistRepoTest(t *testing.T) (*WishlistRepository, sqlmock.Sqlmock, func()) {
	db, mock, cleanup := repotestutil.NewMockDB(t)
	originalDB := config.DB
	config.DB = db

	repo := &WishlistRepository{}
	return repo, mock, func() {
		config.DB = originalDB
		cleanup()
	}
}

func TestWishlistRepository_GetWishlists(t *testing.T) {
	repo, mock, cleanup := setupWishlistRepoTest(t)
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.NewRows([]string{
			"user_id", "product_id", "wishlist_date", "name", "price", "poster_url", "category_id", "category_name",
		}).
			AddRow(1, 101, now, "Game 1", 29.99, "url1", 1, "Action").
			AddRow(1, 101, now, "Game 1", 29.99, "url1", 2, "RPG") // Same game, second category

		mock.ExpectQuery("^SELECT(.*)FROM wishlists").WithArgs(1).WillReturnRows(rows)

		wishlists, err := repo.GetWishlists(1)

		assert.NoError(t, err)
		assert.Len(t, wishlists, 1) // Should group by product_id
		assert.Equal(t, 101, int(wishlists[0].ProductID))
		assert.Len(t, wishlists[0].Game.Categories, 2)
	})
}

func TestWishlistRepository_getWishlist(t *testing.T) {
	repo, mock, cleanup := setupWishlistRepoTest(t)
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.NewRows([]string{
			"user_id", "product_id", "wishlist_date", "name", "price", "poster_url", "category_id", "category_name",
		}).AddRow(1, 101, now, "Game 1", 29.99, "url1", 1, "Action")

		mock.ExpectQuery("^SELECT(.*)FROM wishlists").WithArgs(101, 1).WillReturnRows(rows)

		wishlist, err := repo.getWishlist(1, 101)

		assert.NoError(t, err)
		assert.NotNil(t, wishlist)
		assert.Equal(t, 1, int(wishlist.UserID))
		assert.Equal(t, 101, int(wishlist.ProductID))
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("^SELECT(.*)FROM wishlists").WithArgs(999, 1).WillReturnRows(sqlmock.NewRows([]string{}))

		_, err := repo.getWishlist(1, 999)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func TestWishlistRepository_AddWishlist(t *testing.T) {
	repo, mock, cleanup := setupWishlistRepoTest(t)
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		// 1. Expect the Insert
		mock.ExpectExec("^INSERT INTO wishlists").
			WithArgs(1, 101, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// 2. Expect the getWishlist Select that immediately follows
		now := time.Now()
		rows := sqlmock.NewRows([]string{
			"user_id", "product_id", "wishlist_date", "name", "price", "poster_url", "category_id", "category_name",
		}).AddRow(1, 101, now, "Game 1", 29.99, "url1", nil, nil)

		mock.ExpectQuery("^SELECT(.*)FROM wishlists").WithArgs(101, 1).WillReturnRows(rows)

		wishlist, err := repo.AddWishlist(1, 101)
		assert.NoError(t, err)
		assert.NotNil(t, wishlist)
		assert.Equal(t, 101, int(wishlist.ProductID))
	})
}

func TestWishlistRepository_DeleteWishlist(t *testing.T) {
	repo, mock, cleanup := setupWishlistRepoTest(t)
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		// 1. Expect getWishlist to verify it exists and return the deleted object
		now := time.Now()
		rows := sqlmock.NewRows([]string{
			"user_id", "product_id", "wishlist_date", "name", "price", "poster_url", "category_id", "category_name",
		}).AddRow(1, 101, now, "Game 1", 29.99, "url1", nil, nil)

		mock.ExpectQuery("^SELECT(.*)FROM wishlists").WithArgs(101, 1).WillReturnRows(rows)

		// 2. Expect the actual Delete execution
		mock.ExpectExec("^DELETE FROM wishlists").WithArgs(101, 1).WillReturnResult(sqlmock.NewResult(0, 1))

		wishlist, err := repo.DeleteWishlist(1, 101)

		assert.NoError(t, err)
		assert.NotNil(t, wishlist)
		assert.Equal(t, 101, int(wishlist.ProductID))
	})
}
