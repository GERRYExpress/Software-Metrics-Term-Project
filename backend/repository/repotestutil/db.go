package repotestutil

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// NewMockDB creates a sqlmock database for testing
func NewMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	cleanup := func() { db.Close() }
	return db, mock, cleanup
}

// AssertQueryExecuted checks that a query was executed
func AssertQueryExecuted(t *testing.T, mock sqlmock.Sqlmock, expectedQuery string) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}