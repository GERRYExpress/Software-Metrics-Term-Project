package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCartRepository_GetCarts(t *testing.T) {
	mock, cleanup := setupCollectionTest(t)
	defer cleanup()

	repo := &CartRepository{}

	t.Run("success", func(t *testing.T) {
		expectCollectionList(mock, "carts", 1, newCollectionListRows(time.Now(), "cart_date"))

		carts, err := repo.GetCarts(1)

		assert.NoError(t, err)
		assert.Len(t, carts, 1) // Should group by product_id
		assert.Equal(t, 101, int(carts[0].ProductID))
		assert.Len(t, carts[0].Game.Categories, 2)
	})
}

func TestCartRepository_getCart(t *testing.T) {
	mock, cleanup := setupCollectionTest(t)
	defer cleanup()

	repo := &CartRepository{}

	t.Run("success", func(t *testing.T) {
		expectCollectionGet(mock, "carts", 1, 101, newCollectionSingleRow(time.Now(), "cart_date"))

		cart, err := repo.getCart(1, 101)

		assert.NoError(t, err)
		assert.NotNil(t, cart)
		assert.Equal(t, 1, int(cart.UserID))
		assert.Equal(t, 101, int(cart.ProductID))
	})

	t.Run("not found", func(t *testing.T) {
		expectCollectionGet(mock, "carts", 1, 999, emptyCollectionRows())

		_, err := repo.getCart(1, 999)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func TestCartRepository_AddCart(t *testing.T) {
	mock, cleanup := setupCollectionTest(t)
	defer cleanup()

	repo := &CartRepository{}

	t.Run("success", func(t *testing.T) {
		// 1. Expect the Insert
		expectCollectionInsert(mock, "carts", 1, 101)

		// 2. Expect the getCart Select that immediately follows
		expectCollectionGet(mock, "carts", 1, 101, newCollectionBareRow(time.Now(), "cart_date"))

		cart, err := repo.AddCart(1, 101)
		assert.NoError(t, err)
		assert.NotNil(t, cart)
		assert.Equal(t, 101, int(cart.ProductID))
	})
}

func TestCartRepository_DeleteCart(t *testing.T) {
	mock, cleanup := setupCollectionTest(t)
	defer cleanup()

	repo := &CartRepository{}

	t.Run("success", func(t *testing.T) {
		// 1. Expect getCart to verify it exists and return the deleted object
		expectCollectionGet(mock, "carts", 1, 101, newCollectionBareRow(time.Now(), "cart_date"))

		// 2. Expect the actual Delete execution
		expectCollectionDelete(mock, "carts", 1, 101)

		cart, err := repo.DeleteCart(1, 101)

		assert.NoError(t, err)
		assert.NotNil(t, cart)
		assert.Equal(t, 101, int(cart.ProductID))
	})

	t.Run("fails if cart does not exist", func(t *testing.T) {
		// getCart returns nothing
		expectCollectionGet(mock, "carts", 1, 999, emptyCollectionRows())

		_, err := repo.DeleteCart(1, 999)

		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}
