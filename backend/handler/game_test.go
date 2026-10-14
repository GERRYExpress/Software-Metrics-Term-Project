package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"software-metrics-project/config"
	"software-metrics-project/model"
	"software-metrics-project/repository"
	"software-metrics-project/repository/repotestutil"
	"software-metrics-project/service"
)

func setupGameTest(t *testing.T) (*gin.Engine, sqlmock.Sqlmock, func()) {
	t.Helper()

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

func fullDetailRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"gameID",
		"name",
		"description",
		"price",
		"posterUrl",
		"releaseDate",
		"categoryID",
		"categoryName",
		"userID",
		"username",
		"requirementID",
		"operatingSystem",
		"processor",
		"memory",
		"graphicsCard",
		"directXVersion",
		"network",
		"storage",
		"soundcard",
		"additionalNote",
		"reqType",
		"mediaID",
		"mediaUrl",
	}).AddRow(
		1,
		"Game 1",
		"Desc",
		59.99,
		"url",
		time.Now(),
		nil,
		nil,
		1,
		"user",
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
}

func gameListRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"product_id",
		"name",
		"description",
		"price",
		"poster_url",
		"category_id",
		"category_name",
		"is_hidden",
	}).AddRow(
		1,
		"Test Game",
		"Desc",
		59.99,
		"url",
		nil,
		nil,
		false,
	)
}

func TestGameHandler_GetGames(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success without keyword returns 200", func(t *testing.T) {
		router, mock, cleanup := setupGameTest(t)
		defer cleanup()

		mock.ExpectQuery("SELECT").
			WillReturnRows(gameListRows())

		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		w := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodGet,
			"/games?_page=1&_limit=10",
			nil,
		)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "1", w.Header().Get("X-Total-Count"))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success with keyword and sellerId returns 200", func(t *testing.T) {
		router, mock, cleanup := setupGameTest(t)
		defer cleanup()

		mock.ExpectQuery("SELECT").
			WillReturnRows(gameListRows())

		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		w := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodGet,
			"/games?keyword=test&sellerId=1",
			nil,
		)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("invalid limit returns 400", func(t *testing.T) {
		router, _, cleanup := setupGameTest(t)
		defer cleanup()

		w := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodGet,
			"/games?_limit=-1",
			nil,
		)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid page returns 400", func(t *testing.T) {
		router, _, cleanup := setupGameTest(t)
		defer cleanup()

		w := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodGet,
			"/games?_page=0",
			nil,
		)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid seller id returns 400", func(t *testing.T) {
		router, _, cleanup := setupGameTest(t)
		defer cleanup()

		w := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodGet,
			"/games?sellerId=abc",
			nil,
		)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("seller id less than one returns 400", func(t *testing.T) {
		router, _, cleanup := setupGameTest(t)
		defer cleanup()

		w := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodGet,
			"/games?sellerId=0",
			nil,
		)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("repository list error returns 500", func(t *testing.T) {
		router, mock, cleanup := setupGameTest(t)
		defer cleanup()

		mock.ExpectQuery("SELECT").
			WillReturnError(errors.New("database error"))

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/games", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("repository count error returns 500", func(t *testing.T) {
		router, mock, cleanup := setupGameTest(t)
		defer cleanup()

		mock.ExpectQuery("SELECT").
			WillReturnRows(gameListRows())

		mock.ExpectQuery("SELECT COUNT").
			WillReturnError(errors.New("count error"))

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/games", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGameHandler_GetGame(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success returns 200", func(t *testing.T) {
		router, mock, cleanup := setupGameTest(t)
		defer cleanup()

		mock.ExpectQuery("SELECT").
			WillReturnRows(fullDetailRows())

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/games/1", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("non numeric id returns 400", func(t *testing.T) {
		router, _, cleanup := setupGameTest(t)
		defer cleanup()

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/games/abc", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("repository error returns 500", func(t *testing.T) {
		router, mock, cleanup := setupGameTest(t)
		defer cleanup()

		mock.ExpectQuery("SELECT").
			WillReturnError(errors.New("database error"))

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/games/1", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGameHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	validJSON := `{
		"userId": 1,
		"name": "N",
		"description": "D",
		"price": 10,
		"posterUrl": "U"
	}`

	t.Run("success returns 201", func(t *testing.T) {
		router, mock, cleanup := setupGameTest(t)
		defer cleanup()

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO products").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		mock.ExpectBegin()
		mock.ExpectCommit()

		w := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodPost,
			"/games",
			bytes.NewBufferString(validJSON),
		)
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("invalid json returns 400", func(t *testing.T) {
		router, _, cleanup := setupGameTest(t)
		defer cleanup()

		w := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodPost,
			"/games",
			bytes.NewBufferString(`{bad json`),
		)
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid game details returns 400", func(t *testing.T) {
		router, _, cleanup := setupGameTest(t)
		defer cleanup()

		invalidJSON := `{
			"userId": 0,
			"name": "",
			"description": "",
			"price": -1,
			"posterUrl": ""
		}`

		w := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodPost,
			"/games",
			bytes.NewBufferString(invalidJSON),
		)
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("repository error returns 500", func(t *testing.T) {
		router, mock, cleanup := setupGameTest(t)
		defer cleanup()

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO products").
			WillReturnError(errors.New("insert failed"))
		mock.ExpectRollback()

		w := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodPost,
			"/games",
			bytes.NewBufferString(validJSON),
		)
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGameHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)

	validJSON := `{
		"userId": 1,
		"name": "N",
		"description": "D",
		"price": 10,
		"posterUrl": "U"
	}`

	t.Run("success returns 200", func(t *testing.T) {
		router, mock, cleanup := setupGameTest(t)
		defer cleanup()

		// First transaction: update product data.
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE products").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// The service retrieves the updated game afterward.
		mock.ExpectQuery("SELECT").
			WillReturnRows(fullDetailRows())

		w := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodPut,
			"/games/1",
			bytes.NewBufferString(validJSON),
		)
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("game not found returns 404", func(t *testing.T) {
		router, mock, cleanup := setupGameTest(t)
		defer cleanup()

		mock.ExpectBegin()

		mock.ExpectExec("UPDATE products").
			WillReturnResult(sqlmock.NewResult(1, 0))

		mock.ExpectQuery("SELECT EXISTS").
			WithArgs(99).
			WillReturnRows(
				sqlmock.NewRows([]string{"exists"}).
					AddRow(false),
			)

		mock.ExpectRollback()

		w := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodPut,
			"/games/99",
			bytes.NewBufferString(validJSON),
		)
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("invalid id returns 400", func(t *testing.T) {
		router, _, cleanup := setupGameTest(t)
		defer cleanup()

		w := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodPut,
			"/games/abc",
			bytes.NewBufferString(validJSON),
		)
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("zero id returns 400", func(t *testing.T) {
		router, _, cleanup := setupGameTest(t)
		defer cleanup()

		w := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodPut,
			"/games/0",
			bytes.NewBufferString(validJSON),
		)
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid json returns 400", func(t *testing.T) {
		router, _, cleanup := setupGameTest(t)
		defer cleanup()

		w := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodPut,
			"/games/1",
			bytes.NewBufferString(`invalid`),
		)
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("update query error returns 500", func(t *testing.T) {
		router, mock, cleanup := setupGameTest(t)
		defer cleanup()

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE products").
			WillReturnError(errors.New("update error"))
		mock.ExpectRollback()

		w := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodPut,
			"/games/1",
			bytes.NewBufferString(validJSON),
		)
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGameHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success returns 204", func(t *testing.T) {
		router, mock, cleanup := setupGameTest(t)
		defer cleanup()

		mock.ExpectExec("DELETE FROM products").
			WillReturnResult(sqlmock.NewResult(0, 1))

		w := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodDelete,
			"/games/1?userId=1",
			nil,
		)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("game not found returns 404", func(t *testing.T) {
		router, mock, cleanup := setupGameTest(t)
		defer cleanup()

		mock.ExpectExec("DELETE FROM products").
			WillReturnResult(sqlmock.NewResult(0, 0))

		mock.ExpectQuery("SELECT").
			WillReturnError(repository.ErrGameNotFound)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodDelete,
			"/games/99?userId=1",
			nil,
		)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("invalid game id returns 400", func(t *testing.T) {
		router, _, cleanup := setupGameTest(t)
		defer cleanup()

		w := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodDelete,
			"/games/abc?userId=1",
			nil,
		)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid user id returns 400", func(t *testing.T) {
		router, _, cleanup := setupGameTest(t)
		defer cleanup()

		w := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodDelete,
			"/games/1?userId=abc",
			nil,
		)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("missing user id returns 400", func(t *testing.T) {
		router, _, cleanup := setupGameTest(t)
		defer cleanup()

		w := httptest.NewRecorder()
		req := httptest.NewRequest(
			http.MethodDelete,
			"/games/1",
			nil,
		)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGameCreateRequest_toModel(t *testing.T) {
	request := GameCreateRequest{
		UserID:      1,
		Name:        "Test",
		Description: "Desc",
		Price:       9.99,
		PosterURL:   "poster",
		CategoryIDs: []int{1, 2},
		MediaURLs:   []string{"url1", "url2"},
		HardwareRequirements: []HardwareCreateRequest{
			{
				OS:             "Windows",
				Processor:      "Intel",
				Memory:         "8GB",
				GraphicsCard:   "GTX",
				DirectXVersion: "12",
				Network:        "Broadband",
				Storage:        "10GB",
				Soundcard:      "Any",
				AdditionalNote: "None",
				Type:           "minimum",
			},
		},
	}

	result := request.toModel()

	assert.Equal(t, model.ID(1), result.Account.ID)
	assert.Equal(t, model.Name("Test"), result.Name)
	assert.Equal(t, model.Description("Desc"), result.Description)
	assert.Equal(t, model.Price(9.99), result.Price)
	assert.Equal(t, model.PosterURL("poster"), result.PosterUrl)

	require.Len(t, result.Categories, 2)
	assert.Equal(t, model.ID(1), result.Categories[0].ID)
	assert.Equal(t, model.ID(2), result.Categories[1].ID)

	require.Len(t, result.Media, 2)
	require.Len(t, result.HardwareRequirement, 1)
}

func TestNewGameResponse(t *testing.T) {
	game := model.GameDetail{
		ID:          model.ID(1),
		Name:        model.NewName("Game"),
		Description: model.NewDescription("Description"),
		Price:       model.NewPrice(99.99),
		PosterUrl:   model.NewPosterURL("poster-url"),
		Account:     model.NewAccount(2, "seller"),
		Categories: []model.Category{
			model.NewCategory(10, "Action", false),
		},
		Media: []model.Media{
			model.NewMedia(1, "image-url"),
		},
	}

	response := newGameResponse(game)

	assert.Equal(t, game.ID, response.ID)
	assert.Equal(t, game.Name, response.Name)
	assert.Equal(t, game.Description, response.Description)
	assert.Equal(t, game.Price, response.Price)
	assert.Equal(t, game.PosterUrl, response.PosterURL)
	assert.Equal(t, game.Account, response.Account)
	assert.Equal(t, []model.Category(game.Categories), response.Categories)
	assert.Equal(t, game.Media, response.Media)
}

func TestPagination(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("uses defaults", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

		limit, page, err := pagination(c)

		assert.NoError(t, err)
		assert.Equal(t, 10, limit)
		assert.Equal(t, 1, page)
	})

	t.Run("reads values from query", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(
			http.MethodGet,
			"/?_limit=20&_page=3",
			nil,
		)

		limit, page, err := pagination(c)

		assert.NoError(t, err)
		assert.Equal(t, 20, limit)
		assert.Equal(t, 3, page)
	})

	t.Run("invalid limit returns error", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(
			http.MethodGet,
			"/?_limit=abc",
			nil,
		)

		_, _, err := pagination(c)

		assert.Error(t, err)
	})

	t.Run("invalid page returns error", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(
			http.MethodGet,
			"/?_page=abc",
			nil,
		)

		_, _, err := pagination(c)

		assert.Error(t, err)
	})
}
