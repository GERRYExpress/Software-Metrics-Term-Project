package repository

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"software-metrics-project/config"
	"software-metrics-project/repository/repotestutil"
)

func setupCategoryRepoTest(t *testing.T) (*CategoryRepository, sqlmock.Sqlmock, func()) {
	db, mock, cleanup := repotestutil.NewMockDB(t)
	originalDB := config.DB
	config.DB = db

	repo := &CategoryRepository{}
	return repo, mock, func() {
		config.DB = originalDB
		cleanup()
	}
}

func TestCategoryRepository_Create(t *testing.T) {
	repo, mock, cleanup := setupCategoryRepoTest(t)
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("^INSERT INTO categories").
			WithArgs("RPG", false).
			WillReturnResult(sqlmock.NewResult(1, 1))

		cat, err := repo.Create("RPG")
		assert.NoError(t, err)
		assert.Equal(t, 1, int(cat.ID))
		assert.Equal(t, "RPG", string(cat.Name))
	})
}

func TestCategoryRepository_Update(t *testing.T) {
	repo, mock, cleanup := setupCategoryRepoTest(t)
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("^UPDATE categories").
			WithArgs("Action RPG", true, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		cat, err := repo.Update(1, "Action RPG", true)
		assert.NoError(t, err)
		assert.Equal(t, "Action RPG", string(cat.Name))
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec("^UPDATE categories").
			WithArgs("Action RPG", true, 99).
			WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected

		_, err := repo.Update(99, "Action RPG", true)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func TestCategoryRepository_Delete(t *testing.T) {
	repo, mock, cleanup := setupCategoryRepoTest(t)
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("^DELETE FROM categories").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(1)
		assert.NoError(t, err)
	})
}
