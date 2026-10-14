package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func newTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestCreateAccountEndpoint(t *testing.T) {
	router := newTestRouter()
	group := router.Group("/accounts")

	CreateAccountEndpoint(group)

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{"get accounts", http.MethodGet, "/accounts"},
		{"create account", http.MethodPost, "/accounts"},
		{"delete account", http.MethodDelete, "/accounts/1"},
		{"login", http.MethodPost, "/accounts/login"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.path, nil)

			router.ServeHTTP(w, req)

			// A registered route must not return 404.
			// It may return 400/500 because no request body/database mock exists.
			assert.NotEqual(t, http.StatusNotFound, w.Code)
		})
	}
}

func TestCreateGameEndpoint(t *testing.T) {
	router := newTestRouter()
	group := router.Group("/games")

	CreateGameEndpoint(group)

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{"get games", http.MethodGet, "/games"},
		{"create game", http.MethodPost, "/games"},
		{"update game", http.MethodPut, "/games/1"},
		{"delete game", http.MethodDelete, "/games/1"},
		{"get game", http.MethodGet, "/games/1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.path, nil)

			router.ServeHTTP(w, req)

			assert.NotEqual(t, http.StatusNotFound, w.Code)
		})
	}
}

func TestCreateCategoryEndpoint(t *testing.T) {
	router := newTestRouter()
	group := router.Group("/categories")

	CreateCategoryEndpoint(group)

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{"get categories", http.MethodGet, "/categories"},
		{"create category", http.MethodPost, "/categories"},
		{"update category", http.MethodPut, "/categories/1"},
		{"delete category", http.MethodDelete, "/categories/1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.path, nil)

			router.ServeHTTP(w, req)

			assert.NotEqual(t, http.StatusNotFound, w.Code)
		})
	}
}

func TestCreateWishlistEndpoint(t *testing.T) {
	router := newTestRouter()
	group := router.Group("/wishlists")

	CreateWishlistEndpoint(group)

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{"get wishlists", http.MethodGet, "/wishlists/1"},
		{"add wishlist", http.MethodPost, "/wishlists/1/10"},
		{"delete wishlist", http.MethodDelete, "/wishlists/1/10"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.path, nil)

			router.ServeHTTP(w, req)

			assert.NotEqual(t, http.StatusNotFound, w.Code)
		})
	}
}

func TestCreateCartEndpoint(t *testing.T) {
	router := newTestRouter()
	group := router.Group("/carts")

	CreateCartEndpoint(group)

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{"get carts", http.MethodGet, "/carts/1"},
		{"add cart", http.MethodPost, "/carts/1/10"},
		{"delete cart", http.MethodDelete, "/carts/1/10"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.path, nil)

			router.ServeHTTP(w, req)

			assert.NotEqual(t, http.StatusNotFound, w.Code)
		})
	}
}

func TestCreateOrderEndpoint(t *testing.T) {
	router := newTestRouter()
	group := router.Group("/orders")

	CreateOrderEndpoint(group)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/orders", nil)

	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusNotFound, w.Code)
}

func TestAllEndpointsCanBeRegisteredTogether(t *testing.T) {
	router := newTestRouter()

	api := router.Group("/api")
	v1 := api.Group("/v1")

	CreateAccountEndpoint(v1.Group("/accounts"))
	CreateGameEndpoint(v1.Group("/games"))
	CreateCategoryEndpoint(v1.Group("/categories"))
	CreateWishlistEndpoint(v1.Group("/wishlists"))
	CreateCartEndpoint(v1.Group("/carts"))
	CreateOrderEndpoint(v1.Group("/orders"))

	tests := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/v1/accounts"},
		{http.MethodPost, "/api/v1/accounts"},
		{http.MethodPost, "/api/v1/accounts/login"},
		{http.MethodDelete, "/api/v1/accounts/1"},

		{http.MethodGet, "/api/v1/games"},
		{http.MethodPost, "/api/v1/games"},
		{http.MethodGet, "/api/v1/games/1"},
		{http.MethodPut, "/api/v1/games/1"},
		{http.MethodDelete, "/api/v1/games/1"},

		{http.MethodGet, "/api/v1/categories"},
		{http.MethodPost, "/api/v1/categories"},
		{http.MethodPut, "/api/v1/categories/1"},
		{http.MethodDelete, "/api/v1/categories/1"},

		{http.MethodGet, "/api/v1/wishlists/1"},
		{http.MethodPost, "/api/v1/wishlists/1/10"},
		{http.MethodDelete, "/api/v1/wishlists/1/10"},

		{http.MethodGet, "/api/v1/carts/1"},
		{http.MethodPost, "/api/v1/carts/1/10"},
		{http.MethodDelete, "/api/v1/carts/1/10"},

		{http.MethodPost, "/api/v1/orders"},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.path, nil)

			router.ServeHTTP(w, req)

			assert.NotEqualf(
				t,
				http.StatusNotFound,
				w.Code,
				"route %s %s was not registered",
				tt.method,
				tt.path,
			)
		})
	}
}
