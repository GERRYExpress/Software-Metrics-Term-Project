package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWishlistRepository_GetWishlists(t *testing.T) {
	mock, cleanup := setupCollectionTest(t)
	defer cleanup()

	repo := &WishlistRepository{}

	t.Run("success", func(t *testing.T) {
		expectCollectionList(mock, "wishlists", 1, newCollectionListRows(time.Now(), "wishlist_date"))

		wishlists, err := repo.GetWishlists(1)

		assert.NoError(t, err)
		assert.Len(t, wishlists, 1) // Should group by product_id
		assert.Equal(t, 101, int(wishlists[0].ProductID))
		assert.Len(t, wishlists[0].Game.Categories, 2)
	})
}

func TestWishlistRepository_getWishlist(t *testing.T) {
	mock, cleanup := setupCollectionTest(t)
	defer cleanup()

	repo := &WishlistRepository{}

	t.Run("success", func(t *testing.T) {
		expectCollectionGet(mock, "wishlists", 1, 101, newCollectionSingleRow(time.Now(), "wishlist_date"))

		wishlist, err := repo.getWishlist(1, 101)

		assert.NoError(t, err)
		assert.NotNil(t, wishlist)
		assert.Equal(t, 1, int(wishlist.UserID))
		assert.Equal(t, 101, int(wishlist.ProductID))
	})

	t.Run("not found", func(t *testing.T) {
		expectCollectionGet(mock, "wishlists", 1, 999, emptyCollectionRows())

		_, err := repo.getWishlist(1, 999)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func TestWishlistRepository_AddWishlist(t *testing.T) {
	mock, cleanup := setupCollectionTest(t)
	defer cleanup()

	repo := &WishlistRepository{}

	t.Run("success", func(t *testing.T) {
		// 1. Expect the Insert
		expectCollectionInsert(mock, "wishlists", 1, 101)

		// 2. Expect the getWishlist Select that immediately follows
		expectCollectionGet(mock, "wishlists", 1, 101, newCollectionBareRow(time.Now(), "wishlist_date"))

		wishlist, err := repo.AddWishlist(1, 101)
		assert.NoError(t, err)
		assert.NotNil(t, wishlist)
		assert.Equal(t, 101, int(wishlist.ProductID))
	})
}

func TestWishlistRepository_DeleteWishlist(t *testing.T) {
	mock, cleanup := setupCollectionTest(t)
	defer cleanup()

	repo := &WishlistRepository{}

	t.Run("success", func(t *testing.T) {
		// 1. Expect getWishlist to verify it exists and return the deleted object
		expectCollectionGet(mock, "wishlists", 1, 101, newCollectionBareRow(time.Now(), "wishlist_date"))

		// 2. Expect the actual Delete execution
		expectCollectionDelete(mock, "wishlists", 1, 101)

		wishlist, err := repo.DeleteWishlist(1, 101)

		assert.NoError(t, err)
		assert.NotNil(t, wishlist)
		assert.Equal(t, 101, int(wishlist.ProductID))
	})
}
