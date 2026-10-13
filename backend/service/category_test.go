package service

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"software-metrics-project/config"
	"software-metrics-project/model"
	"software-metrics-project/repository"
	"software-metrics-project/repository/repotestutil"
)

// Integration-style test using real repository with mocked DB
func TestCategoryService_GetCategories(t *testing.T) {
	db, mock, cleanup := repotestutil.NewMockDB(t)
	defer cleanup()

	originalDB := config.DB
	config.DB = db
	defer func() { config.DB = originalDB }()

	repo := &repository.CategoryRepository{}
	svc := &CategoryService{Repo: repo}

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"category_id", "name", "is_hidden"}).
			AddRow(1, "Action", false).
			AddRow(2, "RPG", true)

		mock.ExpectQuery("SELECT.*FROM.*categories").WillReturnRows(rows)

		cats, err := svc.GetCategories()
		require.NoError(t, err)
		assert.Len(t, cats, 2)
		assert.Equal(t, model.NewName("Action"), cats[0].Name)
		assert.False(t, bool(cats[0].IsHidden))
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery("SELECT.*FROM.*categories").WillReturnError(sql.ErrConnDone)

		_, err := svc.GetCategories()
		assert.Error(t, err)
	})
}

func TestCategoryService_Create(t *testing.T) {
	db, mock, cleanup := repotestutil.NewMockDB(t)
	defer cleanup()

	originalDB := config.DB
	config.DB = db
	defer func() { config.DB = originalDB }()

	repo := &repository.CategoryRepository{}
	svc := &CategoryService{Repo: repo}

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO.*categories").
			WithArgs("Strategy", false).
			WillReturnResult(sqlmock.NewResult(3, 1))

		cat, err := svc.Create("Strategy")
		require.NoError(t, err)
		assert.Equal(t, model.NewID(3), cat.ID)
		assert.Equal(t, model.NewName("Strategy"), cat.Name)
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO.*categories").
			WithArgs("Duplicate", false).
			WillReturnError(sql.ErrConnDone)

		_, err := svc.Create("Duplicate")
		assert.Error(t, err)
	})
}

func TestCategoryService_Update(t *testing.T) {
	db, mock, cleanup := repotestutil.NewMockDB(t)
	defer cleanup()

	originalDB := config.DB
	config.DB = db
	defer func() { config.DB = originalDB }()

	repo := &repository.CategoryRepository{}
	svc := &CategoryService{Repo: repo}

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("UPDATE.*categories").
			WithArgs("Updated Action", true, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		cat, err := svc.Update(1, "Updated Action", true)
		require.NoError(t, err)
		assert.Equal(t, model.NewID(1), cat.ID)
		assert.Equal(t, model.NewName("Updated Action"), cat.Name)
		assert.True(t, bool(cat.IsHidden))
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec("UPDATE.*categories").
			WithArgs("Name", false, 999).
			WillReturnResult(sqlmock.NewResult(0, 0))

		_, err := svc.Update(999, "Name", false)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func TestCategoryService_Delete(t *testing.T) {
	db, mock, cleanup := repotestutil.NewMockDB(t)
	defer cleanup()

	originalDB := config.DB
	config.DB = db
	defer func() { config.DB = originalDB }()

	repo := &repository.CategoryRepository{}
	svc := &CategoryService{Repo: repo}

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM.*categories").
			WithArgs(5).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := svc.Delete(5)
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM.*categories").
			WithArgs(999).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := svc.Delete(999)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}