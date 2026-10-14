package service

import (
	"testing"
	"time"

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

// getCartColumns returns the exact 8 columns expected by the repository's Scan.
func getCartColumns() []string {
	return []string{
		"user_id", "product_id", "cart_date", "name", "price", "poster_url", "category_id", "category_name",
	}
}

func TestCartService_GetCarts(t *testing.T) {
	svc, mock, cleanup := setupCartServiceTest(t)
	defer cleanup()

	rows := sqlmock.NewRows(getCartColumns()).AddRow(1, 101, time.Now(), "Game 1", 29.99, "url", nil, nil)
	mock.ExpectQuery("SELECT(.*)FROM").WithArgs(1).WillReturnRows(rows)

	carts, err := svc.GetCarts(1)
	if err != nil || len(carts) == 0 {
		t.Errorf("failed to get carts: %v", err)
	}
}

func TestCartService_AddCart(t *testing.T) {
	svc, mock, cleanup := setupCartServiceTest(t)
	defer cleanup()

	// The repo inserts, then immediately queries to return the updated object.
	mock.ExpectExec("INSERT INTO carts").WithArgs(1, 101, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))

	rows := sqlmock.NewRows(getCartColumns()).AddRow(1, 101, time.Now(), "Game 1", 29.99, "url", nil, nil)
	mock.ExpectQuery("SELECT(.*)FROM carts").WithArgs(101, 1).WillReturnRows(rows)

	cart, err := svc.AddCart(1, 101)
	if err != nil {
		t.Errorf("AddCart executed with error: %v", err)
	} else if cart == nil {
		t.Errorf("expected cart to be returned")
	}
}

func TestCartService_DeleteCart(t *testing.T) {
	svc, mock, cleanup := setupCartServiceTest(t)
	defer cleanup()

	// The repo queries first to verify, then deletes.
	rows := sqlmock.NewRows(getCartColumns()).AddRow(1, 101, time.Now(), "Game 1", 29.99, "url", nil, nil)
	mock.ExpectQuery("SELECT(.*)FROM carts").WithArgs(101, 1).WillReturnRows(rows)

	mock.ExpectExec("DELETE FROM carts").WithArgs(101, 1).WillReturnResult(sqlmock.NewResult(0, 1))

	cart, err := svc.DeleteCart(1, 101)
	if err != nil {
		t.Errorf("DeleteCart executed with error: %v", err)
	} else if cart == nil {
		t.Errorf("expected cart to be returned")
	}
}
