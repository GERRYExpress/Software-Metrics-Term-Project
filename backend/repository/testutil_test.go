package repository

import (
	"database/sql"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"software-metrics-project/config"
	"software-metrics-project/repository/repotestutil"
)

// MockRowsForSuccess returns standard rows for successful GetCarts/GetWishlists
func MockRowsForSuccess() *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"ca.user_id", "ca.product_id", "ca.cart_date",
		"p.name", "p.price", "p.poster_url", "c.category_id", "c.name",
	}).
		AddRow(1, 1, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			"Game 1", 29.99, "poster1.jpg", 1, "Action").
		AddRow(1, 1, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			"Game 1", 29.99, "poster1.jpg", 2, "RPG").
		AddRow(1, 2, time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC),
			"Game 2", 19.99, "poster2.jpg", 3, "Strategy")
}

// MockRowsForEmpty returns empty rows with standard columns
func MockRowsForEmpty() *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"ca.user_id", "ca.product_id", "ca.cart_date",
		"p.name", "p.price", "p.poster_url", "c.category_id", "c.name",
	})
}

// SetupTestDB sets up test database with common configuration
func SetupTestDB() (*sql.DB, interface{}, func(), *sql.DB) {
	db, mock, cleanup := repotestutil.NewMockDB(nil)

	originalDB := config.DB
	config.DB = db

	cleanupFn := func() {
		_ = cleanup
		config.DB = originalDB
	}

	return db, mock, cleanupFn, nil
}