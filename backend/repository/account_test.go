package repository

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"software-metrics-project/config"
	"software-metrics-project/model"
	"software-metrics-project/repository/repotestutil"
)

func TestAccountRepository_GetAccounts(t *testing.T) {
	db, mock, cleanup := repotestutil.NewMockDB(t)
	defer cleanup()

	originalDB := config.DB
	config.DB = db
	defer func() { config.DB = originalDB }()

	repo := &AccountRepository{}

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"user_id", "username", "password", "created_date", "is_buyer",
			"is_seller", "total_game", "total_spend", "total_sale",
		}).
			AddRow(1, "user1", "hash1", "2024-01-01T00:00:00Z", true, false, 0, 0.0, 0.0).
			AddRow(2, "user2", "hash2", "2024-01-02T00:00:00Z", false, true, 5, 100.0, 50.0)

		mock.ExpectQuery("SELECT.*FROM.*accounts").WillReturnRows(rows)

		accounts, err := repo.GetAccounts()
		require.NoError(t, err)
		assert.Len(t, accounts, 2)
		assert.Equal(t, model.NewID(1), accounts[0].ID)
		assert.Equal(t, model.NewUsername("user1"), accounts[0].Username)
		assert.False(t, bool(accounts[0].IsSeller))
		assert.Equal(t, model.NewID(2), accounts[1].ID)
		assert.True(t, bool(accounts[1].IsSeller))
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery("SELECT.*FROM.*accounts").WillReturnError(sql.ErrConnDone)

		_, err := repo.GetAccounts()
		assert.Error(t, err)
	})
}

func TestAccountRepository_Authenticate(t *testing.T) {
	t.Run("success - plain password matches", func(t *testing.T) {
		db, mock, cleanup := repotestutil.NewMockDB(t)
		defer cleanup()

		originalDB := config.DB
		config.DB = db
		defer func() { config.DB = originalDB }()

		mock.ExpectQuery("SELECT user_id, password FROM accounts WHERE username = ?").
			WithArgs("testuser").
			WillReturnRows(sqlmock.NewRows([]string{"user_id", "password"}).
				AddRow(1, "plainpassword"))

		// The code upgrades plain passwords to bcrypt, so expect UPDATE
		mock.ExpectExec(`UPDATE accounts SET password = \? WHERE user_id = \?`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(0, 1))

		repo := &AccountRepository{}
		account, err := repo.Authenticate("testuser", "plainpassword")
		require.NoError(t, err)
		assert.Equal(t, model.NewID(1), account.ID)
		assert.Equal(t, model.NewUsername("testuser"), account.Username)
	})

	t.Run("success - bcrypt password matches", func(t *testing.T) {
		db, mock, cleanup := repotestutil.NewMockDB(t)
		defer cleanup()

		originalDB := config.DB
		config.DB = db
		defer func() { config.DB = originalDB }()

		hash, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
		mock.ExpectQuery("SELECT user_id, password FROM accounts WHERE username = ?").
			WithArgs("testuser2").
			WillReturnRows(sqlmock.NewRows([]string{"user_id", "password"}).
				AddRow(2, string(hash)))

		repo := &AccountRepository{}
		account, err := repo.Authenticate("testuser2", "secret")
		require.NoError(t, err)
		assert.Equal(t, model.NewID(2), account.ID)
	})

	t.Run("invalid credentials - user not found", func(t *testing.T) {
		db, mock, cleanup := repotestutil.NewMockDB(t)
		defer cleanup()

		originalDB := config.DB
		config.DB = db
		defer func() { config.DB = originalDB }()

		mock.ExpectQuery("SELECT user_id, password FROM accounts WHERE username = ?").
			WithArgs("nonexistent").
			WillReturnError(sql.ErrNoRows)

		repo := &AccountRepository{}
		_, err := repo.Authenticate("nonexistent", "password")
		assert.ErrorIs(t, err, ErrInvalidCredentials)
	})

	t.Run("invalid credentials - wrong password", func(t *testing.T) {
		db, mock, cleanup := repotestutil.NewMockDB(t)
		defer cleanup()

		originalDB := config.DB
		config.DB = db
		defer func() { config.DB = originalDB }()

		mock.ExpectQuery("SELECT user_id, password FROM accounts WHERE username = ?").
			WithArgs("testuser3").
			WillReturnRows(sqlmock.NewRows([]string{"user_id", "password"}).
				AddRow(3, "wronghash"))

		repo := &AccountRepository{}
		_, err := repo.Authenticate("testuser3", "wrongpassword")
		assert.ErrorIs(t, err, ErrInvalidCredentials)
	})

	t.Run("db error", func(t *testing.T) {
		db, mock, cleanup := repotestutil.NewMockDB(t)
		defer cleanup()

		originalDB := config.DB
		config.DB = db
		defer func() { config.DB = originalDB }()

		mock.ExpectQuery("SELECT user_id, password FROM accounts WHERE username = ?").
			WithArgs("testuser4").
			WillReturnError(sql.ErrConnDone)

		repo := &AccountRepository{}
		_, err := repo.Authenticate("testuser4", "password")
		assert.Error(t, err)
	})
}

func TestAccountRepository_Create(t *testing.T) {
	db, mock, cleanup := repotestutil.NewMockDB(t)
	defer cleanup()

	originalDB := config.DB
	config.DB = db
	defer func() { config.DB = originalDB }()

	repo := &AccountRepository{}

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO accounts").
			WithArgs("newuser", sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(10, 1))

		account, err := repo.Create("newuser", "password123")
		require.NoError(t, err)
		assert.Equal(t, model.NewID(10), account.ID)
		assert.Equal(t, model.NewUsername("newuser"), account.Username)
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO accounts").
			WithArgs("duplicate", sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(sql.ErrConnDone)

		_, err := repo.Create("duplicate", "password")
		assert.Error(t, err)
	})
}

func TestAccountRepository_Delete(t *testing.T) {
	db, mock, cleanup := repotestutil.NewMockDB(t)
	defer cleanup()

	originalDB := config.DB
	config.DB = db
	defer func() { config.DB = originalDB }()

	repo := &AccountRepository{}

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT EXISTS.*products.*user_id").
			WithArgs(5).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		mock.ExpectExec("DELETE FROM accounts").
			WithArgs(5).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(5)
		assert.NoError(t, err)
	})

	t.Run("account has games - error", func(t *testing.T) {
		mock.ExpectQuery("SELECT EXISTS.*products.*user_id").
			WithArgs(5).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		err := repo.Delete(5)
		assert.ErrorIs(t, err, ErrAccountHasGames)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT EXISTS.*products.*user_id").
			WithArgs(999).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		mock.ExpectExec("DELETE FROM accounts").
			WithArgs(999).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Delete(999)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("db error on exists check", func(t *testing.T) {
		mock.ExpectQuery("SELECT EXISTS.*products.*user_id").
			WithArgs(5).
			WillReturnError(sql.ErrConnDone)

		err := repo.Delete(5)
		assert.Error(t, err)
	})
}