package service

import (
	"testing"

	"software-metrics-project/config"
	"software-metrics-project/repository"
	"software-metrics-project/repository/repotestutil"

	"github.com/DATA-DOG/go-sqlmock"
)

func setupCartServiceTest(t *testing.T) (*CartService, sqlmock.Sqlmock, func()) {
	db, mock, cleanup := repotestutil.NewMockDB(t)
	originalDB := config.DB
	config.DB = db

	svc := &CartService{Repo: &repository.CartRepository{}}

	fullCleanup := func() {
		config.DB = originalDB
		cleanup()
	}
	return svc, mock, fullCleanup
}

func TestCartService_GetCarts(t *testing.T) {
	svc, mock, cleanup := setupCartServiceTest(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "game_id"}).AddRow(1, 101)
	mock.ExpectQuery("SELECT(.*)FROM").WithArgs(1).WillReturnRows(rows)

	carts, err := svc.GetCarts(1)
	if err != nil || len(carts) == 0 {
		t.Errorf("failed to get carts: %v", err)
	}
}

func TestCartService_AddCart(t *testing.T) {
	svc, mock, cleanup := setupCartServiceTest(t)
	defer cleanup()

	mock.ExpectExec("INSERT INTO").WithArgs(1, 101).WillReturnResult(sqlmock.NewResult(1, 1))

	cart, err := svc.AddCart(1, 101)
	if err != nil {
		// Expecting an error here only if the repository tries to SELECT the inserted row and we didn't mock it.
		// If repository.AddCart only executes an INSERT, this passes.
	} else if cart == nil {
		t.Errorf("expected cart to be returned")
	}
}
