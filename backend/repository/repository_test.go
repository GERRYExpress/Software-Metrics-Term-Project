package repository

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"software-metrics-project/config"
	"software-metrics-project/model"
	"software-metrics-project/repository/repotestutil"
)

var errBoom = errors.New("boom")

// useMockDB swaps config.DB for a sqlmock DB for the lifetime of the (sub)test.
func useMockDB(t *testing.T) sqlmock.Sqlmock {
	t.Helper()
	db, mock, cleanup := repotestutil.NewMockDB(t)
	original := config.DB
	config.DB = db
	t.Cleanup(func() {
		config.DB = original
		cleanup()
	})
	return mock
}

var userGameCols = []string{"user_id", "product_id", "date", "name", "price", "poster_url", "category_id", "category_name"}

// ---------------------------------------------------------------- account

var accountCols = []string{"id", "username", "password", "created_date", "is_buyer", "is_seller", "total_game", "total_spend", "total_sale"}

func TestAccountRepository_GetAccounts_Errors(t *testing.T) {
	t.Run("query error", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectQuery("^SELECT \\* FROM accounts$").WillReturnError(errBoom)
		_, err := (&AccountRepository{}).GetAccounts()
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("scan error", func(t *testing.T) {
		mock := useMockDB(t)
		rows := sqlmock.NewRows(accountCols).AddRow("not-a-number", "u", "h", time.Now(), true, false, 0, 0, 0)
		mock.ExpectQuery("^SELECT \\* FROM accounts$").WillReturnRows(rows)
		_, err := (&AccountRepository{}).GetAccounts()
		assert.Error(t, err)
	})
	t.Run("rows error", func(t *testing.T) {
		mock := useMockDB(t)
		rows := sqlmock.NewRows(accountCols).AddRow(1, "u", "h", time.Now(), true, false, 0, 0, 0).RowError(0, errBoom)
		mock.ExpectQuery("^SELECT \\* FROM accounts$").WillReturnRows(rows)
		_, err := (&AccountRepository{}).GetAccounts()
		assert.ErrorIs(t, err, errBoom)
	})
}

func TestAccountRepository_Authenticate_Branches(t *testing.T) {
	const query = "^SELECT user_id, password FROM accounts"

	t.Run("database error is returned as is", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectQuery(query).WithArgs("u").WillReturnError(errBoom)
		_, err := (&AccountRepository{}).Authenticate("u", "pw")
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("wrong password against bcrypt hash", func(t *testing.T) {
		mock := useMockDB(t)
		hash, err := bcrypt.GenerateFromPassword([]byte("right"), bcrypt.MinCost)
		require.NoError(t, err)
		mock.ExpectQuery(query).WithArgs("u").WillReturnRows(sqlmock.NewRows([]string{"user_id", "password"}).AddRow(1, string(hash)))
		_, err = (&AccountRepository{}).Authenticate("u", "wrong")
		assert.ErrorIs(t, err, ErrInvalidCredentials)
	})
	t.Run("wrong password against legacy plaintext", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectQuery(query).WithArgs("u").WillReturnRows(sqlmock.NewRows([]string{"user_id", "password"}).AddRow(1, "plain"))
		_, err := (&AccountRepository{}).Authenticate("u", "other")
		assert.ErrorIs(t, err, ErrInvalidCredentials)
	})
	t.Run("legacy plaintext password is upgraded to bcrypt", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectQuery(query).WithArgs("u").WillReturnRows(sqlmock.NewRows([]string{"user_id", "password"}).AddRow(7, "plain"))
		mock.ExpectExec("^UPDATE accounts SET password").WithArgs(sqlmock.AnyArg(), 7).WillReturnResult(sqlmock.NewResult(0, 1))
		account, err := (&AccountRepository{}).Authenticate("u", "plain")
		assert.NoError(t, err)
		assert.Equal(t, 7, int(account.ID))
	})
	t.Run("upgrade failure is returned", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectQuery(query).WithArgs("u").WillReturnRows(sqlmock.NewRows([]string{"user_id", "password"}).AddRow(7, "plain"))
		mock.ExpectExec("^UPDATE accounts SET password").WithArgs(sqlmock.AnyArg(), 7).WillReturnError(errBoom)
		_, err := (&AccountRepository{}).Authenticate("u", "plain")
		assert.ErrorIs(t, err, errBoom)
	})
}

func TestAccountRepository_Create_Errors(t *testing.T) {
	t.Run("insert error", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectExec("^INSERT INTO accounts").WithArgs("u", sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnError(errBoom)
		_, err := (&AccountRepository{}).Create("u", "pw")
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("last insert id error", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectExec("^INSERT INTO accounts").WithArgs("u", sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewErrorResult(errBoom))
		_, err := (&AccountRepository{}).Create("u", "pw")
		assert.ErrorIs(t, err, errBoom)
	})
}

func TestAccountRepository_Delete_Errors(t *testing.T) {
	const exists = "^SELECT EXISTS"
	const del = "^DELETE FROM accounts"
	notOwner := func() *sqlmock.Rows { return sqlmock.NewRows([]string{"exists"}).AddRow(false) }

	t.Run("ownership check error", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectQuery(exists).WithArgs(1).WillReturnError(errBoom)
		assert.ErrorIs(t, (&AccountRepository{}).Delete(1), errBoom)
	})
	t.Run("delete error", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectQuery(exists).WithArgs(1).WillReturnRows(notOwner())
		mock.ExpectExec(del).WithArgs(1).WillReturnError(errBoom)
		assert.ErrorIs(t, (&AccountRepository{}).Delete(1), errBoom)
	})
	t.Run("rows affected error", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectQuery(exists).WithArgs(1).WillReturnRows(notOwner())
		mock.ExpectExec(del).WithArgs(1).WillReturnResult(sqlmock.NewErrorResult(errBoom))
		assert.ErrorIs(t, (&AccountRepository{}).Delete(1), errBoom)
	})
	t.Run("account not found", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectQuery(exists).WithArgs(1).WillReturnRows(notOwner())
		mock.ExpectExec(del).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 0))
		assert.ErrorIs(t, (&AccountRepository{}).Delete(1), sql.ErrNoRows)
	})
}

// --------------------------------------------------------------- category

var categoryCols = []string{"category_id", "name", "is_hidden"}

func TestCategoryRepository_GetCategories(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := useMockDB(t)
		rows := sqlmock.NewRows(categoryCols).AddRow(1, "RPG", false).AddRow(2, "Indie", true)
		mock.ExpectQuery("^SELECT \\* FROM categories$").WillReturnRows(rows)
		cats, err := (&CategoryRepository{}).GetCategories()
		assert.NoError(t, err)
		assert.Len(t, cats, 2)
		assert.Equal(t, "Indie", string(cats[1].Name))
	})
	t.Run("query error", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectQuery("^SELECT \\* FROM categories$").WillReturnError(errBoom)
		_, err := (&CategoryRepository{}).GetCategories()
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("scan error", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectQuery("^SELECT \\* FROM categories$").WillReturnRows(sqlmock.NewRows(categoryCols).AddRow("x", "RPG", false))
		_, err := (&CategoryRepository{}).GetCategories()
		assert.Error(t, err)
	})
	t.Run("rows error", func(t *testing.T) {
		mock := useMockDB(t)
		rows := sqlmock.NewRows(categoryCols).AddRow(1, "RPG", false).RowError(0, errBoom)
		mock.ExpectQuery("^SELECT \\* FROM categories$").WillReturnRows(rows)
		_, err := (&CategoryRepository{}).GetCategories()
		assert.ErrorIs(t, err, errBoom)
	})
}

func TestCategoryRepository_WriteErrors(t *testing.T) {
	t.Run("create: insert error", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectExec("^INSERT INTO categories").WithArgs("RPG", false).WillReturnError(errBoom)
		_, err := (&CategoryRepository{}).Create("RPG")
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("create: last insert id error", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectExec("^INSERT INTO categories").WithArgs("RPG", false).WillReturnResult(sqlmock.NewErrorResult(errBoom))
		_, err := (&CategoryRepository{}).Create("RPG")
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("update: exec error", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectExec("^UPDATE categories").WithArgs("RPG", false, 1).WillReturnError(errBoom)
		_, err := (&CategoryRepository{}).Update(1, "RPG", false)
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("update: rows affected error", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectExec("^UPDATE categories").WithArgs("RPG", false, 1).WillReturnResult(sqlmock.NewErrorResult(errBoom))
		_, err := (&CategoryRepository{}).Update(1, "RPG", false)
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("delete: exec error", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectExec("^DELETE FROM categories").WithArgs(1).WillReturnError(errBoom)
		assert.ErrorIs(t, (&CategoryRepository{}).Delete(1), errBoom)
	})
	t.Run("delete: rows affected error", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectExec("^DELETE FROM categories").WithArgs(1).WillReturnResult(sqlmock.NewErrorResult(errBoom))
		assert.ErrorIs(t, (&CategoryRepository{}).Delete(1), errBoom)
	})
	t.Run("delete: not found", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectExec("^DELETE FROM categories").WithArgs(99).WillReturnResult(sqlmock.NewResult(0, 0))
		assert.ErrorIs(t, (&CategoryRepository{}).Delete(99), sql.ErrNoRows)
	})
}

// ------------------------------------------------------------------- cart

func TestCartRepository_Errors(t *testing.T) {
	t.Run("GetCarts: query error", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectQuery("^SELECT(.*)FROM carts").WithArgs(1).WillReturnError(errBoom)
		_, err := (&CartRepository{}).GetCarts(1)
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("GetCarts: scan error", func(t *testing.T) {
		mock := useMockDB(t)
		rows := sqlmock.NewRows(userGameCols).AddRow(1, "x", time.Now(), "G", 1.0, "u", nil, nil)
		mock.ExpectQuery("^SELECT(.*)FROM carts").WithArgs(1).WillReturnRows(rows)
		_, err := (&CartRepository{}).GetCarts(1)
		assert.Error(t, err)
	})
	t.Run("GetCarts: rows error", func(t *testing.T) {
		mock := useMockDB(t)
		rows := sqlmock.NewRows(userGameCols).AddRow(1, 101, time.Now(), "G", 1.0, "u", nil, nil).RowError(0, errBoom)
		mock.ExpectQuery("^SELECT(.*)FROM carts").WithArgs(1).WillReturnRows(rows)
		_, err := (&CartRepository{}).GetCarts(1)
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("GetCarts: cart row whose product no longer exists", func(t *testing.T) {
		mock := useMockDB(t)
		rows := sqlmock.NewRows(userGameCols).AddRow(1, 101, time.Now(), nil, nil, nil, nil, nil)
		mock.ExpectQuery("^SELECT(.*)FROM carts").WithArgs(1).WillReturnRows(rows)
		carts, err := (&CartRepository{}).GetCarts(1)
		assert.NoError(t, err)
		require.Len(t, carts, 1)
		assert.Equal(t, 101, int(carts[0].ProductID))
		assert.Empty(t, carts[0].Game.Categories)
	})
	t.Run("AddCart: insert error", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectExec("^INSERT INTO carts").WithArgs(1, 101, sqlmock.AnyArg()).WillReturnError(errBoom)
		_, err := (&CartRepository{}).AddCart(1, 101)
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("getCart: query error", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectQuery("^SELECT(.*)FROM carts").WithArgs(101, 1).WillReturnError(errBoom)
		_, err := (&CartRepository{}).getCart(1, 101)
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("getCart: scan error", func(t *testing.T) {
		mock := useMockDB(t)
		rows := sqlmock.NewRows(userGameCols).AddRow(1, "x", time.Now(), "G", 1.0, "u", nil, nil)
		mock.ExpectQuery("^SELECT(.*)FROM carts").WithArgs(101, 1).WillReturnRows(rows)
		_, err := (&CartRepository{}).getCart(1, 101)
		assert.Error(t, err)
	})
	t.Run("getCart: rows error", func(t *testing.T) {
		mock := useMockDB(t)
		rows := sqlmock.NewRows(userGameCols).AddRow(1, 101, time.Now(), "G", 1.0, "u", nil, nil).RowError(0, errBoom)
		mock.ExpectQuery("^SELECT(.*)FROM carts").WithArgs(101, 1).WillReturnRows(rows)
		_, err := (&CartRepository{}).getCart(1, 101)
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("DeleteCart: delete error", func(t *testing.T) {
		mock := useMockDB(t)
		rows := sqlmock.NewRows(userGameCols).AddRow(1, 101, time.Now(), "G", 1.0, "u", nil, nil)
		mock.ExpectQuery("^SELECT(.*)FROM carts").WithArgs(101, 1).WillReturnRows(rows)
		mock.ExpectExec("^DELETE FROM carts").WithArgs(101, 1).WillReturnError(errBoom)
		_, err := (&CartRepository{}).DeleteCart(1, 101)
		assert.ErrorIs(t, err, errBoom)
	})
}

// --------------------------------------------------------------- wishlist

func TestWishlistRepository_Errors(t *testing.T) {
	t.Run("GetWishlists: query error", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectQuery("^SELECT(.*)FROM wishlists").WithArgs(1).WillReturnError(errBoom)
		_, err := (&WishlistRepository{}).GetWishlists(1)
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("GetWishlists: scan error", func(t *testing.T) {
		mock := useMockDB(t)
		rows := sqlmock.NewRows(userGameCols).AddRow(1, "x", time.Now(), "G", 1.0, "u", nil, nil)
		mock.ExpectQuery("^SELECT(.*)FROM wishlists").WithArgs(1).WillReturnRows(rows)
		_, err := (&WishlistRepository{}).GetWishlists(1)
		assert.Error(t, err)
	})
	t.Run("GetWishlists: rows error", func(t *testing.T) {
		mock := useMockDB(t)
		rows := sqlmock.NewRows(userGameCols).AddRow(1, 101, time.Now(), "G", 1.0, "u", nil, nil).RowError(0, errBoom)
		mock.ExpectQuery("^SELECT(.*)FROM wishlists").WithArgs(1).WillReturnRows(rows)
		_, err := (&WishlistRepository{}).GetWishlists(1)
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("GetWishlists: product no longer exists", func(t *testing.T) {
		mock := useMockDB(t)
		rows := sqlmock.NewRows(userGameCols).AddRow(1, 101, time.Now(), nil, nil, nil, nil, nil)
		mock.ExpectQuery("^SELECT(.*)FROM wishlists").WithArgs(1).WillReturnRows(rows)
		wishlists, err := (&WishlistRepository{}).GetWishlists(1)
		assert.NoError(t, err)
		require.Len(t, wishlists, 1)
		assert.Empty(t, wishlists[0].Game.Categories)
	})
	t.Run("AddWishlist: insert error", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectExec("^INSERT INTO wishlists").WithArgs(1, 101, sqlmock.AnyArg()).WillReturnError(errBoom)
		_, err := (&WishlistRepository{}).AddWishlist(1, 101)
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("getWishlist: query error", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectQuery("^SELECT(.*)FROM wishlists").WithArgs(101, 1).WillReturnError(errBoom)
		_, err := (&WishlistRepository{}).getWishlist(1, 101)
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("getWishlist: scan error", func(t *testing.T) {
		mock := useMockDB(t)
		rows := sqlmock.NewRows(userGameCols).AddRow(1, "x", time.Now(), "G", 1.0, "u", nil, nil)
		mock.ExpectQuery("^SELECT(.*)FROM wishlists").WithArgs(101, 1).WillReturnRows(rows)
		_, err := (&WishlistRepository{}).getWishlist(1, 101)
		assert.Error(t, err)
	})
	t.Run("getWishlist: rows error", func(t *testing.T) {
		mock := useMockDB(t)
		rows := sqlmock.NewRows(userGameCols).AddRow(1, 101, time.Now(), "G", 1.0, "u", nil, nil).RowError(0, errBoom)
		mock.ExpectQuery("^SELECT(.*)FROM wishlists").WithArgs(101, 1).WillReturnRows(rows)
		_, err := (&WishlistRepository{}).getWishlist(1, 101)
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("DeleteWishlist: not in wishlist", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectQuery("^SELECT(.*)FROM wishlists").WithArgs(999, 1).WillReturnRows(sqlmock.NewRows(userGameCols))
		_, err := (&WishlistRepository{}).DeleteWishlist(1, 999)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
	t.Run("DeleteWishlist: delete error", func(t *testing.T) {
		mock := useMockDB(t)
		rows := sqlmock.NewRows(userGameCols).AddRow(1, 101, time.Now(), "G", 1.0, "u", nil, nil)
		mock.ExpectQuery("^SELECT(.*)FROM wishlists").WithArgs(101, 1).WillReturnRows(rows)
		mock.ExpectExec("^DELETE FROM wishlists").WithArgs(101, 1).WillReturnError(errBoom)
		_, err := (&WishlistRepository{}).DeleteWishlist(1, 101)
		assert.ErrorIs(t, err, errBoom)
	})
}

// ------------------------------------------------------------------ order

var orderCols = []string{
	"order_id", "user_id", "order_date", "payment_method", "final_price",
	"product_id", "name", "description", "price", "poster_url", "release_date",
	"category_id", "category_name", "is_hidden",
}

const orderInsert = "^INSERT INTO orders"

func TestOrderRepository_Create_Branches(t *testing.T) {
	newOrder := func() model.Order {
		return model.Order{UserID: 1, PaymentMethod: "Cash", FinalPrice: 10, Games: []model.Game{{ID: 101}}}
	}
	expectInsert := func(mock sqlmock.Sqlmock) {
		mock.ExpectBegin()
		mock.ExpectExec(orderInsert).WithArgs(1, "Cash", 10.0).WillReturnResult(sqlmock.NewResult(10, 1))
	}

	t.Run("begin error", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectBegin().WillReturnError(errBoom)
		_, err := (&OrderRepository{}).Create(newOrder())
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("last insert id error rolls back", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectBegin()
		mock.ExpectExec(orderInsert).WithArgs(1, "Cash", 10.0).WillReturnResult(sqlmock.NewErrorResult(errBoom))
		mock.ExpectRollback()
		_, err := (&OrderRepository{}).Create(newOrder())
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("order item insert error rolls back", func(t *testing.T) {
		mock := useMockDB(t)
		expectInsert(mock)
		mock.ExpectExec("^INSERT INTO order_items").WithArgs(10, 101).WillReturnError(errBoom)
		mock.ExpectRollback()
		_, err := (&OrderRepository{}).Create(newOrder())
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("commit error", func(t *testing.T) {
		mock := useMockDB(t)
		expectInsert(mock)
		mock.ExpectExec("^INSERT INTO order_items").WithArgs(10, 101).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit().WillReturnError(errBoom)
		_, err := (&OrderRepository{}).Create(newOrder())
		assert.ErrorIs(t, err, errBoom)
	})

	// getOrder is reached through Create, once the transaction has committed.
	createThenSelect := func(t *testing.T) sqlmock.Sqlmock {
		mock := useMockDB(t)
		mock.ExpectBegin()
		mock.ExpectExec(orderInsert).WithArgs(1, "Cash", 10.0).WillReturnResult(sqlmock.NewResult(10, 1))
		mock.ExpectCommit()
		return mock
	}
	noItems := model.Order{UserID: 1, PaymentMethod: "Cash", FinalPrice: 10}
	now := time.Now()

	t.Run("getOrder groups categories under each game", func(t *testing.T) {
		mock := createThenSelect(t)
		rows := sqlmock.NewRows(orderCols).
			AddRow(10, 1, now, "Cash", 10.0, 101, "G1", "d", 5.0, "p", now, 1, "Action", false).
			AddRow(10, 1, now, "Cash", 10.0, 101, "G1", "d", 5.0, "p", now, 2, "RPG", true).
			AddRow(10, 1, now, "Cash", 10.0, 102, "G2", "d", 5.0, "p", now, nil, nil, nil)
		mock.ExpectQuery("^SELECT(.*)FROM orders").WithArgs(10).WillReturnRows(rows)

		order, err := (&OrderRepository{}).Create(noItems)
		require.NoError(t, err)
		require.Len(t, order.Games, 2)
		assert.Len(t, order.Games[0].Categories, 2)
		assert.Empty(t, order.Games[1].Categories)
	})
	t.Run("getOrder query error", func(t *testing.T) {
		mock := createThenSelect(t)
		mock.ExpectQuery("^SELECT(.*)FROM orders").WithArgs(10).WillReturnError(errBoom)
		_, err := (&OrderRepository{}).Create(noItems)
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("getOrder scan error", func(t *testing.T) {
		mock := createThenSelect(t)
		rows := sqlmock.NewRows(orderCols).AddRow("x", 1, now, "Cash", 10.0, 101, "G1", "d", 5.0, "p", now, nil, nil, nil)
		mock.ExpectQuery("^SELECT(.*)FROM orders").WithArgs(10).WillReturnRows(rows)
		_, err := (&OrderRepository{}).Create(noItems)
		assert.Error(t, err)
	})
	t.Run("getOrder rows error", func(t *testing.T) {
		mock := createThenSelect(t)
		rows := sqlmock.NewRows(orderCols).AddRow(10, 1, now, "Cash", 10.0, 101, "G1", "d", 5.0, "p", now, nil, nil, nil).RowError(0, errBoom)
		mock.ExpectQuery("^SELECT(.*)FROM orders").WithArgs(10).WillReturnRows(rows)
		_, err := (&OrderRepository{}).Create(noItems)
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("getOrder finds no rows", func(t *testing.T) {
		mock := createThenSelect(t)
		mock.ExpectQuery("^SELECT(.*)FROM orders").WithArgs(10).WillReturnRows(sqlmock.NewRows(orderCols))
		_, err := (&OrderRepository{}).Create(noItems)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}

// ------------------------------------------------------------------- game

var gameListCols = []string{"product_id", "name", "description", "price", "poster_url", "category_id", "category_name", "is_hidden"}

var gameDetailCols = []string{
	"product_id", "name", "description", "price", "poster_url", "release_date",
	"category_id", "category_name", "user_id", "username",
	"requirement_id", "operating_system", "processor", "memory", "graphics_card",
	"directx_version", "network", "storage", "soundcard", "additional_note", "type",
	"media_id", "media_url",
}

// detailRow builds one GetGame row; nil arguments model LEFT JOIN misses.
func detailRow(categoryID, categoryName, requirementID, mediaID, mediaURL any) []any {
	var userID, username any = 42, "seller"
	spec := func(v string) any {
		if requirementID == nil {
			return nil
		}
		return v
	}
	return []any{
		1, "Game 1", "Desc", 59.99, "poster", time.Now(),
		categoryID, categoryName, userID, username,
		requirementID, spec("Win10"), spec("i7"), spec("16GB"), spec("RTX"),
		spec("12"), spec("Broadband"), spec("50GB"), spec("Yes"), spec("note"), spec("Minimum"),
		mediaID, mediaURL,
	}
}

func addDetailRow(rows *sqlmock.Rows, values []any) *sqlmock.Rows {
	driverValues := make([]driver.Value, len(values))
	for i, v := range values {
		driverValues[i] = v
	}
	return rows.AddRow(driverValues...)
}

func TestGameRepository_GetGames(t *testing.T) {
	seller := 7
	t.Run("GetGames groups categories per game", func(t *testing.T) {
		mock := useMockDB(t)
		rows := sqlmock.NewRows(gameListCols).
			AddRow(1, "G1", "d", 9.5, "p", 1, "Action", false).
			AddRow(1, "G1", "d", 9.5, "p", 2, "RPG", true).
			AddRow(2, "G2", "d", 5.0, "p", nil, nil, nil)
		mock.ExpectQuery("^SELECT(.*)FROM \\(SELECT").WithArgs("RPG", 7).WillReturnRows(rows)

		games, err := (&GameRepository{}).GetGames("RPG", &seller, 10, 1)
		require.NoError(t, err)
		require.Len(t, games, 2)
		assert.Len(t, games[0].Categories, 2)
		assert.Empty(t, games[1].Categories)
	})
	t.Run("GetGamesByName filters by keyword", func(t *testing.T) {
		mock := useMockDB(t)
		rows := sqlmock.NewRows(gameListCols).AddRow(1, "Zelda", "d", 9.5, "p", nil, nil, nil)
		mock.ExpectQuery("^SELECT(.*)FROM \\(SELECT").WithArgs("%zelda%").WillReturnRows(rows)

		games, err := (&GameRepository{}).GetGamesByName("zelda", "", nil, 10, 1)
		require.NoError(t, err)
		require.Len(t, games, 1)
		assert.Equal(t, "Zelda", string(games[0].Name))
	})
	t.Run("query error", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectQuery("^SELECT(.*)FROM \\(SELECT").WillReturnError(errBoom)
		_, err := (&GameRepository{}).GetGames("", nil, 10, 1)
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("scan error", func(t *testing.T) {
		mock := useMockDB(t)
		rows := sqlmock.NewRows(gameListCols).AddRow("x", "G1", "d", 9.5, "p", nil, nil, nil)
		mock.ExpectQuery("^SELECT(.*)FROM \\(SELECT").WillReturnRows(rows)
		_, err := (&GameRepository{}).GetGames("", nil, 10, 1)
		assert.Error(t, err)
	})
	t.Run("rows error", func(t *testing.T) {
		mock := useMockDB(t)
		rows := sqlmock.NewRows(gameListCols).AddRow(1, "G1", "d", 9.5, "p", nil, nil, nil).RowError(0, errBoom)
		mock.ExpectQuery("^SELECT(.*)FROM \\(SELECT").WillReturnRows(rows)
		_, err := (&GameRepository{}).GetGames("", nil, 10, 1)
		assert.ErrorIs(t, err, errBoom)
	})
}

func Test_listAndDetailColumns(t *testing.T) {
	assert.Contains(t, listColumns(), "c.is_hidden")
	assert.Len(t, detailColumns(), len(gameDetailCols))
}

func TestGameRepository_GetGamesCount_Error(t *testing.T) {
	mock := useMockDB(t)
	mock.ExpectQuery("^SELECT COUNT\\(\\*\\) FROM products").WillReturnError(errBoom)
	_, err := (&GameRepository{}).GetGamesCount("", "", nil)
	assert.ErrorIs(t, err, errBoom)
}

func TestGameRepository_GetGame_Branches(t *testing.T) {
	const query = "^SELECT(.*)FROM products p"

	t.Run("duplicated join rows are merged and media order is reversed", func(t *testing.T) {
		mock := useMockDB(t)
		rows := sqlmock.NewRows(gameDetailCols)
		addDetailRow(rows, detailRow(10, "RPG", 100, 200, "http://a"))
		addDetailRow(rows, detailRow(10, "RPG", 100, 201, "http://b")) // same category and requirement, new media
		addDetailRow(rows, detailRow(11, "Indie", 100, 201, "http://b"))
		mock.ExpectQuery(query).WithArgs(false, 1).WillReturnRows(rows)

		game, err := (&GameRepository{}).GetGame(1)
		require.NoError(t, err)
		assert.Len(t, game.Categories, 2)
		assert.Len(t, game.HardwareRequirement, 1)
		require.Len(t, game.Media, 2)
		assert.Equal(t, 201, int(game.Media[0].ID)) // reversed
		assert.Equal(t, 42, int(game.Account.ID))
	})
	t.Run("left join misses leave related lists empty", func(t *testing.T) {
		mock := useMockDB(t)
		rows := sqlmock.NewRows(gameDetailCols)
		addDetailRow(rows, detailRow(nil, nil, nil, nil, nil))
		mock.ExpectQuery(query).WithArgs(false, 1).WillReturnRows(rows)

		game, err := (&GameRepository{}).GetGame(1)
		require.NoError(t, err)
		assert.Empty(t, game.Categories)
		assert.Empty(t, game.HardwareRequirement)
		assert.Empty(t, game.Media)
	})
	t.Run("media row without url is skipped", func(t *testing.T) {
		mock := useMockDB(t)
		rows := sqlmock.NewRows(gameDetailCols)
		addDetailRow(rows, detailRow(nil, nil, nil, 200, nil))
		mock.ExpectQuery(query).WithArgs(false, 1).WillReturnRows(rows)

		game, err := (&GameRepository{}).GetGame(1)
		require.NoError(t, err)
		assert.Empty(t, game.Media)
	})
	t.Run("query error", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectQuery(query).WithArgs(false, 1).WillReturnError(errBoom)
		_, err := (&GameRepository{}).GetGame(1)
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("scan error", func(t *testing.T) {
		mock := useMockDB(t)
		rows := sqlmock.NewRows(gameDetailCols)
		bad := detailRow(nil, nil, nil, nil, nil)
		bad[0] = "x"
		addDetailRow(rows, bad)
		mock.ExpectQuery(query).WithArgs(false, 1).WillReturnRows(rows)
		_, err := (&GameRepository{}).GetGame(1)
		assert.Error(t, err)
	})
	t.Run("rows error", func(t *testing.T) {
		mock := useMockDB(t)
		rows := sqlmock.NewRows(gameDetailCols)
		addDetailRow(rows, detailRow(nil, nil, nil, nil, nil))
		rows.RowError(0, errBoom)
		mock.ExpectQuery(query).WithArgs(false, 1).WillReturnRows(rows)
		_, err := (&GameRepository{}).GetGame(1)
		assert.ErrorIs(t, err, errBoom)
	})
}

func fullGame() model.GameDetail {
	game := model.NewGameDetail(0, "New Game", "Desc", 19.99, "poster", time.Now())
	game.Account = model.NewAccount(42, "seller")
	game.Categories = []model.Category{model.NewCategory(1, "RPG", false)}
	game.HardwareRequirement = []model.Hardware{model.NewHardware(model.HardwareData{
		OS: "Win10", Processor: "i7", Memory: "16GB", GraphicsCard: "RTX", DirectXVersion: "12",
		Network: "Broadband", Storage: "50GB", Soundcard: "Yes", AdditionalNote: "note", Type: "Minimum",
	})}
	game.Media = []model.Media{model.NewMedia(0, "u1"), model.NewMedia(0, "u1"), model.NewMedia(0, "u2")} // u1 is a duplicate
	return game
}

func expectRelations(mock sqlmock.Sqlmock, productID int64) {
	mock.ExpectExec("^INSERT INTO product_categories").WithArgs(productID, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("^INSERT INTO hardware_requirement").
		WithArgs("Win10", "i7", "16GB", "RTX", "12", "Broadband", "50GB", "Yes", "note", "Minimum", productID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("^INSERT INTO media").WithArgs("u1", productID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("^INSERT INTO media").WithArgs("u2", productID).WillReturnResult(sqlmock.NewResult(2, 1))
}

func TestGameRepository_Create_Branches(t *testing.T) {
	t.Run("full game with release date, requirements and duplicate media", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectBegin()
		mock.ExpectExec("^INSERT INTO products(.*)release_date").WillReturnResult(sqlmock.NewResult(10, 1))
		mock.ExpectCommit()
		mock.ExpectBegin()
		expectRelations(mock, 10)
		mock.ExpectCommit()

		created, err := (&GameRepository{}).Create(fullGame())
		require.NoError(t, err)
		assert.Equal(t, 10, int(created.ID))
		assert.Len(t, created.Media, 3) // the returned struct echoes the input
	})
	t.Run("begin error", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectBegin().WillReturnError(errBoom)
		_, err := (&GameRepository{}).Create(fullGame())
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("product insert error rolls back", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectBegin()
		mock.ExpectExec("^INSERT INTO products").WillReturnError(errBoom)
		mock.ExpectRollback()
		_, err := (&GameRepository{}).Create(fullGame())
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("last insert id error rolls back", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectBegin()
		mock.ExpectExec("^INSERT INTO products").WillReturnResult(sqlmock.NewErrorResult(errBoom))
		mock.ExpectRollback()
		_, err := (&GameRepository{}).Create(fullGame())
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("commit error", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectBegin()
		mock.ExpectExec("^INSERT INTO products").WillReturnResult(sqlmock.NewResult(10, 1))
		mock.ExpectCommit().WillReturnError(errBoom)
		_, err := (&GameRepository{}).Create(fullGame())
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("relations transaction cannot begin", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectBegin()
		mock.ExpectExec("^INSERT INTO products").WillReturnResult(sqlmock.NewResult(10, 1))
		mock.ExpectCommit()
		mock.ExpectBegin().WillReturnError(errBoom)
		_, err := (&GameRepository{}).Create(fullGame())
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("relation insert error is returned", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectBegin()
		mock.ExpectExec("^INSERT INTO products").WillReturnResult(sqlmock.NewResult(10, 1))
		mock.ExpectCommit()
		mock.ExpectBegin()
		mock.ExpectExec("^INSERT INTO product_categories").WithArgs(10, 1).WillReturnError(errBoom)
		_, err := (&GameRepository{}).Create(fullGame())
		assert.ErrorIs(t, err, errBoom)
	})
}

func TestGameRepository_Update(t *testing.T) {
	expectUpdateBody := func(mock sqlmock.Sqlmock) {
		mock.ExpectBegin()
		mock.ExpectExec("^UPDATE products").WillReturnResult(sqlmock.NewResult(0, 1))
	}
	expectExists := func(mock sqlmock.Sqlmock, exists bool) {
		mock.ExpectQuery("^SELECT EXISTS").WithArgs(5).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(exists))
	}
	expectCleanup := func(mock sqlmock.Sqlmock) {
		mock.ExpectExec("^DELETE FROM product_categories").WithArgs(5).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec("^DELETE FROM hardware_requirement").WithArgs(5).WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectExec("^DELETE FROM media").WithArgs(5).WillReturnResult(sqlmock.NewResult(0, 1))
	}

	t.Run("success replaces relations", func(t *testing.T) {
		mock := useMockDB(t)
		expectUpdateBody(mock)
		expectExists(mock, true)
		expectCleanup(mock)
		mock.ExpectCommit()
		mock.ExpectBegin()
		expectRelations(mock, 5)
		mock.ExpectCommit()

		updated, err := (&GameRepository{}).Update(5, fullGame())
		require.NoError(t, err)
		assert.Equal(t, 5, int(updated.ID))
		assert.Equal(t, "New Game", string(updated.Name))
	})
	t.Run("update error rolls back", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectBegin()
		mock.ExpectExec("^UPDATE products").WillReturnError(errBoom)
		mock.ExpectRollback()
		_, err := (&GameRepository{}).Update(5, fullGame())
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("existence check error rolls back", func(t *testing.T) {
		mock := useMockDB(t)
		expectUpdateBody(mock)
		mock.ExpectQuery("^SELECT EXISTS").WithArgs(5).WillReturnError(errBoom)
		mock.ExpectRollback()
		_, err := (&GameRepository{}).Update(5, fullGame())
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("missing game returns ErrGameNotFound", func(t *testing.T) {
		mock := useMockDB(t)
		expectUpdateBody(mock)
		expectExists(mock, false)
		mock.ExpectRollback()
		_, err := (&GameRepository{}).Update(5, fullGame())
		assert.ErrorIs(t, err, ErrGameNotFound)
	})
	t.Run("relation cleanup error rolls back", func(t *testing.T) {
		mock := useMockDB(t)
		expectUpdateBody(mock)
		expectExists(mock, true)
		mock.ExpectExec("^DELETE FROM product_categories").WithArgs(5).WillReturnError(errBoom)
		mock.ExpectRollback()
		_, err := (&GameRepository{}).Update(5, fullGame())
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("relation insert error is returned", func(t *testing.T) {
		mock := useMockDB(t)
		expectUpdateBody(mock)
		expectExists(mock, true)
		expectCleanup(mock)
		mock.ExpectCommit()
		mock.ExpectBegin()
		mock.ExpectExec("^INSERT INTO product_categories").WithArgs(5, 1).WillReturnError(errBoom)
		_, err := (&GameRepository{}).Update(5, fullGame())
		assert.ErrorIs(t, err, errBoom)
	})
}

// insertGameRelationsTx with a caller-supplied transaction: every relation type can fail.
func TestGameRepository_insertGameRelationsTx_WithTx(t *testing.T) {
	run := func(t *testing.T, arrange func(mock sqlmock.Sqlmock)) error {
		mock := useMockDB(t)
		mock.ExpectBegin()
		arrange(mock)
		tx, err := config.DB.Begin()
		require.NoError(t, err)
		defer tx.Rollback()
		return (&GameRepository{}).insertGameRelationsTx(tx, 10, fullGame())
	}

	t.Run("success", func(t *testing.T) {
		err := run(t, func(mock sqlmock.Sqlmock) { expectRelations(mock, 10) })
		assert.NoError(t, err)
	})
	t.Run("category insert error", func(t *testing.T) {
		err := run(t, func(mock sqlmock.Sqlmock) {
			mock.ExpectExec("^INSERT INTO product_categories").WithArgs(10, 1).WillReturnError(errBoom)
		})
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("requirement insert error", func(t *testing.T) {
		err := run(t, func(mock sqlmock.Sqlmock) {
			mock.ExpectExec("^INSERT INTO product_categories").WithArgs(10, 1).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec("^INSERT INTO hardware_requirement").WillReturnError(errBoom)
		})
		assert.ErrorIs(t, err, errBoom)
	})
	t.Run("media insert error", func(t *testing.T) {
		err := run(t, func(mock sqlmock.Sqlmock) {
			mock.ExpectExec("^INSERT INTO product_categories").WithArgs(10, 1).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec("^INSERT INTO hardware_requirement").WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec("^INSERT INTO media").WithArgs("u1", 10).WillReturnError(errBoom)
		})
		assert.ErrorIs(t, err, errBoom)
	})
}

func TestGameRepository_Delete_Errors(t *testing.T) {
	const del = "^DELETE FROM products"
	t.Run("exec error", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectExec(del).WithArgs(1, 42).WillReturnError(errBoom)
		assert.ErrorIs(t, (&GameRepository{}).Delete(1, 42), errBoom)
	})
	t.Run("rows affected error", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectExec(del).WithArgs(1, 42).WillReturnResult(sqlmock.NewErrorResult(errBoom))
		assert.ErrorIs(t, (&GameRepository{}).Delete(1, 42), errBoom)
	})
	t.Run("existence lookup error", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectExec(del).WithArgs(1, 42).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectQuery("^SELECT EXISTS").WithArgs(1).WillReturnError(errBoom)
		assert.ErrorIs(t, (&GameRepository{}).Delete(1, 42), errBoom)
	})
	t.Run("game does not exist", func(t *testing.T) {
		mock := useMockDB(t)
		mock.ExpectExec(del).WithArgs(1, 42).WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectQuery("^SELECT EXISTS").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
		assert.ErrorIs(t, (&GameRepository{}).Delete(1, 42), ErrGameNotFound)
	})
}
