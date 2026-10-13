package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"software-metrics-project/config"
	"software-metrics-project/repository/repotestutil"
)

func setupAccountRepoTest(t *testing.T) (*AccountRepository, sqlmock.Sqlmock, func()) {
	db, mock, cleanup := repotestutil.NewMockDB(t)
	originalDB := config.DB
	config.DB = db

	repo := &AccountRepository{}
	return repo, mock, func() {
		config.DB = originalDB
		cleanup()
	}
}

func TestAccountRepository_GetAccounts(t *testing.T) {
	repo, mock, cleanup := setupAccountRepoTest(t)
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.NewRows([]string{"id", "username", "password", "created_date", "is_buyer", "is_seller", "total_game", "total_spend", "total_sale"}).
			AddRow(1, "testuser", "hash", now, true, false, 0, 0, 0)

		mock.ExpectQuery("^SELECT \\* FROM accounts$").WillReturnRows(rows)

		accounts, err := repo.GetAccounts()
		assert.NoError(t, err)
		assert.Len(t, accounts, 1)
		assert.Equal(t, "testuser", string(accounts[0].Username))
	})
}

func TestAccountRepository_Authenticate(t *testing.T) {
	repo, mock, cleanup := setupAccountRepoTest(t)
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		rows := sqlmock.NewRows([]string{"user_id", "password"}).AddRow(1, string(hash))

		mock.ExpectQuery("^SELECT user_id, password FROM accounts").
			WithArgs("testuser").
			WillReturnRows(rows)

		account, err := repo.Authenticate("testuser", "password123")
		assert.NoError(t, err)
		assert.Equal(t, 1, int(account.ID))
	})

	t.Run("invalid credentials", func(t *testing.T) {
		mock.ExpectQuery("^SELECT user_id, password FROM accounts").
			WithArgs("testuser").
			WillReturnError(sql.ErrNoRows)

		_, err := repo.Authenticate("testuser", "wrongpass")
		assert.ErrorIs(t, err, ErrInvalidCredentials)
	})
}

func TestAccountRepository_Create(t *testing.T) {
	repo, mock, cleanup := setupAccountRepoTest(t)
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("^INSERT INTO accounts").
			WithArgs("newuser", sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		account, err := repo.Create("newuser", "password123")
		assert.NoError(t, err)
		assert.Equal(t, 1, int(account.ID))
		assert.Equal(t, "newuser", string(account.Username))
	})
}

func TestAccountRepository_Delete(t *testing.T) {
	repo, mock, cleanup := setupAccountRepoTest(t)
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		// First query checks if account owns games
		mock.ExpectQuery("^SELECT EXISTS").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
		// Second query deletes the account
		mock.ExpectExec("^DELETE FROM accounts").WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(1)
		assert.NoError(t, err)
	})

	t.Run("fails if account owns games", func(t *testing.T) {
		mock.ExpectQuery("^SELECT EXISTS").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		err := repo.Delete(1)
		assert.ErrorIs(t, err, ErrAccountHasGames)
	})
}
