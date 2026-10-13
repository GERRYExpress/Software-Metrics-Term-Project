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
	"software-metrics-project/repository"
	"software-metrics-project/repository/repotestutil"
	"software-metrics-project/service"
)

func setupAccountTest(t *testing.T) (*gin.Engine, sqlmock.Sqlmock, func()) {
	db, mock, cleanup := repotestutil.NewMockDB(t)

	originalDB := config.DB
	config.DB = db

	repo := &repository.AccountRepository{}
	svc := &service.AccountService{Repo: repo}
	h := &AccountHandler{Service: svc}

	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/accounts", h.GetAccounts)
	router.POST("/accounts", h.Create)
	router.POST("/accounts/login", h.Login)
	router.DELETE("/accounts/:id", h.Delete)

	fullCleanup := func() {
		config.DB = originalDB
		cleanup()
	}

	return router, mock, fullCleanup
}

func TestAccountHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, _, cleanup := setupAccountTest(t)
	defer cleanup()

	t.Run("bad json binding - 400 Bad Request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/accounts/login", bytes.NewBufferString(`{"username": "test"`)) // missing closing brace
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAccountHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, mock, cleanup := setupAccountTest(t)
	defer cleanup()

	t.Run("invalid ID parameter - 400 Bad Request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/accounts/abc", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "id must be a positive integer")
	})

	t.Run("success - 204 No Content", func(t *testing.T) {
		// Add this line to satisfy the repository's initial check
		mock.ExpectQuery("^SELECT EXISTS").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		mock.ExpectExec("DELETE FROM").WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/accounts/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}

func Test_positiveID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		raw     string
		want    int
		wantErr bool
	}{
		{"valid positive integer", "42", 42, false},
		{"zero (invalid)", "0", 0, true},
		{"negative number (invalid)", "-5", 0, true},
		{"not a number", "abc", 0, true},
		{"empty string", "", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := positiveID(tt.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("positiveID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("positiveID() = %v, want %v", got, tt.want)
			}
		})
	}
}
