package handlertestutil

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

// NewTestRouter creates a gin router with recovery middleware for testing
func NewTestRouter(setup func(*gin.Engine)) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery())
	if setup != nil {
		setup(r)
	}
	return r
}

// PerformRequest executes a request against the router and returns the response recorder
func PerformRequest(r *gin.Engine, method, path string, body string, headers map[string]string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	if body != "" {
		req.Body = http.NoBody // will be set by test if needed
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// AuthHeader returns Authorization header with Bearer token
func AuthHeader(token string) map[string]string {
	return map[string]string{"Authorization": "Bearer " + token}
}

// JSONBody creates a request with JSON body
func JSONBody(method, path string, jsonBody string, token string) (*http.Request, *httptest.ResponseRecorder) {
	req, _ := http.NewRequest(method, path, nil)
	if jsonBody != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	return req, w
}