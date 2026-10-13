package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"software-metrics-project/config"
	"software-metrics-project/model"
	"software-metrics-project/repository"
	"software-metrics-project/repository/repotestutil"
	"software-metrics-project/service"
)

func setupGameTest(t *testing.T) (*gin.Engine, sqlmock.Sqlmock, func()) {
	db, mock, cleanup := repotestutil.NewMockDB(t)

	originalDB := config.DB
	config.DB = db

	repo := &repository.GameRepository{}
	svc := &service.GameService{Repo: repo}
	h := &GameHandler{Service: svc}

	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/games", h.GetGames)
	router.GET("/games/:id", h.GetGame)
	router.POST("/games", h.Create)
	router.PUT("/games/:id", h.Update)
	router.DELETE("/games/:id", h.Delete)

	fullCleanup := func() {
		config.DB = originalDB
		cleanup()
	}

	return router, mock, fullCleanup
}

func TestGameHandler_GetGames(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, mock, cleanup := setupGameTest(t)
	defer cleanup()

	t.Run("success with pagination - 200 OK", func(t *testing.T) {
		// Mock the count query and the select query
		mock.ExpectQuery("SELECT COUNT(.*) FROM").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		rows := sqlmock.NewRows([]string{"id", "name", "price"}).AddRow(1, "Test Game", 59.99)
		mock.ExpectQuery("SELECT(.*)FROM").WillReturnRows(rows)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/games?_page=1&_limit=10", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "1", w.Header().Get("X-Total-Count"))
	})

	t.Run("invalid pagination parameter - 400 Bad Request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/games?_page=-1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGameHandler_GetGame(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, mock, cleanup := setupGameTest(t)
	defer cleanup()

	t.Run("success - 200 OK", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Test Game")
		mock.ExpectQuery("SELECT(.*)FROM").WithArgs(1).WillReturnRows(rows)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/games/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid id - 400 Bad Request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/games/abc", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGameHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, mock, cleanup := setupGameTest(t)
	defer cleanup()

	t.Run("success - 201 Created", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))

		validJSON := `{
			"userId": 1, 
			"name": "New Game", 
			"description": "A fun game", 
			"price": 19.99, 
			"posterUrl": "http://example.com/poster.png"
		}`
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/games", bytes.NewBufferString(validJSON))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)
		assert.NotEqual(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("invalid game details (negative price) - 400 Bad Request", func(t *testing.T) {
		invalidJSON := `{"userId": 1, "name": "New Game", "description": "Desc", "price": -5.00, "posterUrl": "url"}`
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/games", bytes.NewBufferString(invalidJSON))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid game details")
	})
}

func TestGameHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, mock, cleanup := setupGameTest(t)
	defer cleanup()

	t.Run("success - 204 No Content", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM").WithArgs(1, 42).WillReturnResult(sqlmock.NewResult(0, 1))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/games/1?userId=42", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("missing user id - 400 Bad Request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/games/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGameCreateRequest_toModel(t *testing.T) {
	req := GameCreateRequest{
		UserID:      1,
		Name:        "Test",
		Description: "Desc",
		Price:       9.99,
		CategoryIDs: []int{1, 2},
		MediaURLs:   []string{"url1", "url2"},
	}

	result := req.toModel()

	assert.Equal(t, model.ID(1), result.Account.ID)
	assert.Equal(t, model.Name("Test"), result.Name)
	assert.Len(t, result.Categories, 2)
	assert.Len(t, result.Media, 2)
}

func Test_pagination(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		queryStr      string
		expectedLimit int
		expectedPage  int
		expectError   bool
	}{
		{"default values", "", 10, 1, false},
		{"custom values", "?_limit=25&_page=3", 25, 3, false},
		{"invalid limit", "?_limit=abc", 0, 0, true},
		{"negative page", "?_page=-5", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodGet, "/"+tt.queryStr, nil)

			limit, page, err := pagination(c)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedLimit, limit)
				assert.Equal(t, tt.expectedPage, page)
			}
		})
	}
}
