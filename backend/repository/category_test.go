package repository

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"software-metrics-project/config"
	"software-metrics-project/model"
	"software-metrics-project/repository/repotestutil"
)

func TestCategoryRepository_GetCategories(t *testing.T) {
	db, mock, cleanup := repotestutil.NewMockDB(t)
	defer cleanup()

	// Override global config.DB for testing
	originalDB := config.DB
	config.DB = db
	defer func() { config.DB = originalDB }()

	repo := &CategoryRepository{}

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"category_id", "name", "is_hidden"}).
			AddRow(1, "Action", false).
			AddRow(2, "RPG", true)

		mock.ExpectQuery("SELECT.*FROM.*categories").WillReturnRows(rows)

		cats, err := repo.GetCategories()
		require.NoError(t, err)
		assert.Len(t, cats, 2)
		assert.Equal(t, model.NewName("Action"), cats[0].Name)
		assert.False(t, bool(cats[0].IsHidden))
		assert.Equal(t, model.NewName("RPG"), cats[1].Name)
		assert.True(t, bool(cats[1].IsHidden))
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery("SELECT.*FROM.*categories").WillReturnError(sql.ErrConnDone)

		_, err := repo.GetCategories()
		assert.Error(t, err)
	})
}

func TestCategoryRepository_Create(t *testing.T) {
	db, mock, cleanup := repotestutil.NewMockDB(t)
	defer cleanup()

	originalDB := config.DB
	config.DB = db
	defer func() { config.DB = originalDB }()

	repo := &CategoryRepository{}

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO.*categories").
			WithArgs("Strategy", false).
			WillReturnResult(sqlmock.NewResult(3, 1))

		cat, err := repo.Create("Strategy")
		require.NoError(t, err)
		assert.Equal(t, model.NewID(3), cat.ID)
		assert.Equal(t, model.NewName("Strategy"), cat.Name)
		assert.False(t, bool(cat.IsHidden))
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO.*categories").
			WithArgs("Duplicate", false).
			WillReturnError(sql.ErrConnDone)

		_, err := repo.Create("Duplicate")
		assert.Error(t, err)
	})
}

func TestCategoryRepository_Update(t *testing.T) {
	db, mock, cleanup := repotestutil.NewMockDB(t)
	defer cleanup()

	originalDB := config.DB
	config.DB = db
	defer func() { config.DB = originalDB }()

	repo := &CategoryRepository{}

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("UPDATE.*categories").
			WithArgs("Updated Action", true, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		cat, err := repo.Update(1, "Updated Action", true)
		require.NoError(t, err)
		assert.Equal(t, model.NewID(1), cat.ID)
		assert.Equal(t, model.NewName("Updated Action"), cat.Name)
		assert.True(t, bool(cat.IsHidden))
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec("UPDATE.*categories").
			WithArgs("Name", false, 999).
			WillReturnResult(sqlmock.NewResult(0, 0))

		_, err := repo.Update(999, "Name", false)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func TestCategoryRepository_Delete(t *testing.T) {
	db, mock, cleanup := repotestutil.NewMockDB(t)
	defer cleanup()

	originalDB := config.DB
	config.DB = db
	defer func() { config.DB = originalDB }()

	repo := &CategoryRepository{}

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM.*categories").
			WithArgs(5).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(5)
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM.*categories").
			WithArgs(999).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Delete(999)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}