package service

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"software-metrics-project/config"
	"software-metrics-project/repository"
	"software-metrics-project/repository/repotestutil"
)

func setupCategoryServiceTest(t *testing.T) (*CategoryService, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, cleanup := repotestutil.NewMockDB(t)

	originalDB := config.DB
	config.DB = db

	service := &CategoryService{
		Repo: &repository.CategoryRepository{},
	}

	return service, mock, func() {
		config.DB = originalDB
		cleanup()
	}
}

func TestCategoryService_GetCategories(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mock, cleanup := setupCategoryServiceTest(t)
		defer cleanup()

		rows := sqlmock.NewRows([]string{
			"category_id",
			"name",
			"is_hidden",
		}).AddRow(1, "Action", false)

		mock.ExpectQuery("SELECT").
			WillReturnRows(rows)

		categories, err := svc.GetCategories()

		assert.NoError(t, err)
		assert.Len(t, categories, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("repository error", func(t *testing.T) {
		svc, mock, cleanup := setupCategoryServiceTest(t)
		defer cleanup()

		mock.ExpectQuery("SELECT").
			WillReturnError(errors.New("database error"))

		_, err := svc.GetCategories()

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCategoryService_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mock, cleanup := setupCategoryServiceTest(t)
		defer cleanup()

		mock.ExpectExec("INSERT INTO categories").
			WithArgs("Action", false).
			WillReturnResult(sqlmock.NewResult(1, 1))

		category, err := svc.Create("Action")

		assert.NoError(t, err)
		assert.Equal(t, 1, int(category.ID))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("repository error", func(t *testing.T) {
		svc, mock, cleanup := setupCategoryServiceTest(t)
		defer cleanup()

		mock.ExpectExec("INSERT INTO").
			WithArgs("Action").
			WillReturnError(errors.New("insert error"))

		_, err := svc.Create("Action")

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCategoryService_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mock, cleanup := setupCategoryServiceTest(t)
		defer cleanup()

		mock.ExpectExec("UPDATE").
			WithArgs("Action", true, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		category, err := svc.Update(1, "Action", true)

		assert.NoError(t, err)
		assert.Equal(t, 1, int(category.ID))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("repository error", func(t *testing.T) {
		svc, mock, cleanup := setupCategoryServiceTest(t)
		defer cleanup()

		mock.ExpectExec("UPDATE").
			WithArgs("Action", true, 1).
			WillReturnError(errors.New("update error"))

		_, err := svc.Update(1, "Action", true)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCategoryService_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, mock, cleanup := setupCategoryServiceTest(t)
		defer cleanup()

		mock.ExpectExec("DELETE FROM").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := svc.Delete(1)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("repository error", func(t *testing.T) {
		svc, mock, cleanup := setupCategoryServiceTest(t)
		defer cleanup()

		mock.ExpectExec("DELETE FROM").
			WithArgs(1).
			WillReturnError(errors.New("delete error"))

		err := svc.Delete(1)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
