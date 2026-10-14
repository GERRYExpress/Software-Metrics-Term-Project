package repository

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"software-metrics-project/config"
	"software-metrics-project/repository/repotestutil"
)

func setupCollectionTest(t *testing.T) (sqlmock.Sqlmock, func()) {
	t.Helper()
	db, mock, cleanup := repotestutil.NewMockDB(t)
	originalDB := config.DB
	config.DB = db
	return mock, func() {
		config.DB = originalDB
		cleanup()
	}
}

func collectionColumns(dateCol string) []string {
	return []string{"user_id", "product_id", dateCol, "name", "price", "poster_url", "category_id", "category_name"}
}

func newCollectionListRows(now time.Time, dateCol string) *sqlmock.Rows {
	return sqlmock.NewRows(collectionColumns(dateCol)).
		AddRow(1, 101, now, "Game 1", 29.99, "url1", 1, "Action").
		AddRow(1, 101, now, "Game 1", 29.99, "url1", 2, "RPG") // Same game, second category
}

func newCollectionSingleRow(now time.Time, dateCol string) *sqlmock.Rows {
	return sqlmock.NewRows(collectionColumns(dateCol)).
		AddRow(1, 101, now, "Game 1", 29.99, "url1", 1, "Action")
}

func newCollectionBareRow(now time.Time, dateCol string) *sqlmock.Rows {
	return sqlmock.NewRows(collectionColumns(dateCol)).
		AddRow(1, 101, now, "Game 1", 29.99, "url1", nil, nil)
}

func emptyCollectionRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{})
}

func expectCollectionList(mock sqlmock.Sqlmock, table string, userID int, rows *sqlmock.Rows) {
	mock.ExpectQuery("^SELECT(.*)FROM " + table).WithArgs(userID).WillReturnRows(rows)
}

func expectCollectionGet(mock sqlmock.Sqlmock, table string, userID, productID int, rows *sqlmock.Rows) {
	mock.ExpectQuery("^SELECT(.*)FROM " + table).WithArgs(productID, userID).WillReturnRows(rows)
}

func expectCollectionInsert(mock sqlmock.Sqlmock, table string, userID, productID int) {
	mock.ExpectExec("^INSERT INTO " + table).
		WithArgs(userID, productID, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
}

func expectCollectionDelete(mock sqlmock.Sqlmock, table string, userID, productID int) {
	// squirrel.Eq sorts keys: product_id, then user_id.
	mock.ExpectExec("^DELETE FROM " + table).WithArgs(productID, userID).WillReturnResult(sqlmock.NewResult(0, 1))
}
