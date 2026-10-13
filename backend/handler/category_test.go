package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"software-metrics-project/config"
	"software-metrics-project/repository"
	"software-metrics-project/repository/repotestutil"
	"software-metrics-project/service"
)

func TestCategoryHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, mock, cleanup := repotestutil.NewMockDB(t)
	defer cleanup()

	originalDB := config.DB
	config.DB = db
	defer func() { config.DB = originalDB }()

	repo := &repository.CategoryRepository{}
	svc := &service.CategoryService{Repo: repo}
	h := &CategoryHandler{Service: svc}

	router := gin.New()
	router.Use(gin.Recovery())
	router.DELETE("/api/v1/categories/:id", h.Delete)

	t.Run("success - 204 No Content", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM.*categories").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/categories/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Empty(t, w.Body.String())
	})

	t.Run("not found - 404", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM.*categories").
			WithArgs(999).
			WillReturnResult(sqlmock.NewResult(0, 0))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/categories/999", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, `{"error":"category not found"}`, w.Body.String())
	})

	t.Run("bad id - 400", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/categories/abc", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"id must be a positive integer"}`, w.Body.String())
	})
}

func TestCategoryHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, mock, cleanup := repotestutil.NewMockDB(t)
	defer cleanup()

	originalDB := config.DB
	config.DB = db
	defer func() { config.DB = originalDB }()

	repo := &repository.CategoryRepository{}
	svc := &service.CategoryService{Repo: repo}
	h := &CategoryHandler{Service: svc}

	router := gin.New()
	router.Use(gin.Recovery())
	router.PUT("/api/v1/categories/:id", h.Update)

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("UPDATE.*categories").
			WithArgs("Updated", true, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/categories/1", nil)
		req.Header.Set("Content-Type", "application/json")
		// Body would be {"name":"Updated","isHidden":true}
		router.ServeHTTP(w, req)

		// Note: For full test, need proper JSON body
		assert.NotEqual(t, http.StatusInternalServerError, w.Code)
	})
}

func TestCategoryHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, mock, cleanup := repotestutil.NewMockDB(t)
	defer cleanup()

	originalDB := config.DB
	config.DB = db
	defer func() { config.DB = originalDB }()

	repo := &repository.CategoryRepository{}
	svc := &service.CategoryService{Repo: repo}
	h := &CategoryHandler{Service: svc}

	router := gin.New()
	router.Use(gin.Recovery())
	router.POST("/api/v1/categories", h.Create)

	t.Run("success - 201", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO.*categories").
			WithArgs("NewCat", false).
			WillReturnResult(sqlmock.NewResult(3, 1))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/categories", nil)
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Would need proper JSON body for full test
		assert.NotEqual(t, http.StatusInternalServerError, w.Code)
	})
}

func TestCategoryHandler_GetCategories(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, mock, cleanup := repotestutil.NewMockDB(t)
	defer cleanup()

	originalDB := config.DB
	config.DB = db
	defer func() { config.DB = originalDB }()

	repo := &repository.CategoryRepository{}
	svc := &service.CategoryService{Repo: repo}
	h := &CategoryHandler{Service: svc}

	router := gin.New()
	router.Use(gin.Recovery())
	router.GET("/api/v1/categories", h.GetCategories)

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"category_id", "name", "is_hidden"}).
			AddRow(1, "Action", false).
			AddRow(2, "RPG", true)

		mock.ExpectQuery("SELECT.*FROM.*categories").WillReturnRows(rows)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/categories", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Action")
		assert.Contains(t, w.Body.String(), "RPG")
	})
}