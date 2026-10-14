package handler

import (
	"bytes"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"software-metrics-project/config"
	"software-metrics-project/repository/repotestutil"
)

// setupHandlerTest swaps in the mock DB and wires routes via wire.
// Table/column names, routes, and service/handler constructors stay
// with the callers; only the mock DB swap and router boilerplate is shared.
func setupHandlerTest(t *testing.T, wire func(*gin.Engine)) (*gin.Engine, sqlmock.Sqlmock, func()) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	db, mock, cleanup := repotestutil.NewMockDB(t)

	originalDB := config.DB
	config.DB = db

	router := gin.New()
	router.Use(gin.Recovery())
	wire(router)

	return router, mock, func() {
		config.DB = originalDB
		cleanup()
	}
}

// collectionHandlerColumns returns the exact 8 columns expected by the
// cart/wishlist repository Scan. The date column differs per table.
func collectionHandlerColumns(dateCol string) []string {
	return []string{
		"user_id", "product_id", dateCol, "name", "price", "poster_url", "category_id", "category_name",
	}
}

func newCollectionHandlerRow(dateCol string, userID, productID int) *sqlmock.Rows {
	return sqlmock.NewRows(collectionHandlerColumns(dateCol)).
		AddRow(userID, productID, time.Now(), "Game 1", 29.99, "url1", nil, nil)
}

func expectCollectionListQuery(mock sqlmock.Sqlmock, table string, userID int, rows *sqlmock.Rows) {
	mock.ExpectQuery("SELECT(.*)FROM " + table).WithArgs(userID).WillReturnRows(rows)
}

// Squirrel sorts Eq keys alphabetically: product_id < user_id,
// so args are (gameID, userID).
func expectCollectionGetQuery(mock sqlmock.Sqlmock, table string, gameID, userID int, rows *sqlmock.Rows) {
	mock.ExpectQuery("SELECT(.*)FROM " + table).WithArgs(gameID, userID).WillReturnRows(rows)
}

func expectCollectionInsert(mock sqlmock.Sqlmock, table string, userID, gameID int) {
	mock.ExpectExec("INSERT INTO " + table).
		WithArgs(userID, gameID, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
}

func expectCollectionDelete(mock sqlmock.Sqlmock, table string, gameID, userID int) {
	mock.ExpectExec("DELETE FROM " + table).
		WithArgs(gameID, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

// serve sends a bodyless request to the router.
func serve(router *gin.Engine, method, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, nil)
	router.ServeHTTP(w, req)
	return w
}

// serveJSON sends a request with a JSON body to the router.
func serveJSON(router *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	return w
}

// expectGameListOK stubs the game list SELECT with the given rows.
func expectGameListOK(mock sqlmock.Sqlmock, rows *sqlmock.Rows) {
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
}

// expectGameCountOK stubs the game count SELECT with the given total.
func expectGameCountOK(mock sqlmock.Sqlmock, n int) {
	mock.ExpectQuery("SELECT COUNT").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(n))
}

func expectSelectError(mock sqlmock.Sqlmock, err error) {
	mock.ExpectQuery("SELECT").WillReturnError(err)
}

func expectCountError(mock sqlmock.Sqlmock, err error) {
	mock.ExpectQuery("SELECT COUNT").WillReturnError(err)
}

func expectDeleteProducts(mock sqlmock.Sqlmock, rowsAffected int64) {
	mock.ExpectExec("DELETE FROM products").
		WillReturnResult(sqlmock.NewResult(0, rowsAffected))
}

// assertStatusAndMet checks the status code and that all mock expectations met.
func assertStatusAndMet(t *testing.T, mock sqlmock.Sqlmock, w *httptest.ResponseRecorder, want int) {
	t.Helper()
	assert.Equal(t, want, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}
