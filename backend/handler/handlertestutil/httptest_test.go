package handlertestutil

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewTestRouter(t *testing.T) {
	t.Run("creates router without setup", func(t *testing.T) {
		r := NewTestRouter(nil)
		assert.NotNil(t, r)
	})

	t.Run("creates router with setup function", func(t *testing.T) {
		setupCalled := false
		setupFunc := func(engine *gin.Engine) {
			setupCalled = true
			engine.GET("/ping", func(c *gin.Context) {
				c.String(http.StatusOK, "pong")
			})
		}

		r := NewTestRouter(setupFunc)
		assert.NotNil(t, r)
		assert.True(t, setupCalled)

		// Verify the route was actually attached
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "pong", w.Body.String())
	})
}

func TestPerformRequest(t *testing.T) {
	t.Run("executes request successfully", func(t *testing.T) {
		r := gin.New()
		r.GET("/test", func(c *gin.Context) {
			assert.Equal(t, "custom-value", c.GetHeader("X-Custom-Header"))
			c.String(http.StatusOK, "success")
		})

		headers := map[string]string{
			"X-Custom-Header": "custom-value",
		}

		w := PerformRequest(r, http.MethodGet, "/test", "", headers)

		assert.NotNil(t, w)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "success", w.Body.String())
	})
}

func TestAuthHeader(t *testing.T) {
	t.Run("generates correct header map", func(t *testing.T) {
		token := "my-secret-token"
		result := AuthHeader(token)

		assert.NotNil(t, result)
		assert.Equal(t, "Bearer my-secret-token", result["Authorization"])
	})
}

func TestJSONBody(t *testing.T) {
	t.Run("sets json header and auth token", func(t *testing.T) {
		jsonBody := `{"key": "value"}`
		token := "test-token"

		req, w := JSONBody(http.MethodPost, "/api/resource", jsonBody, token)

		assert.NotNil(t, req)
		assert.NotNil(t, w)

		// Note: The original JSONBody function doesn't actually bind the jsonBody string
		// into the req.Body. It only sets the headers. We test its actual behavior here.
		assert.Equal(t, http.MethodPost, req.Method)
		assert.Equal(t, "/api/resource", req.URL.Path)
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer test-token", req.Header.Get("Authorization"))
	})

	t.Run("omits headers if empty", func(t *testing.T) {
		req, w := JSONBody(http.MethodGet, "/api/resource", "", "")

		assert.NotNil(t, req)
		assert.NotNil(t, w)

		assert.Empty(t, req.Header.Get("Content-Type"))
		assert.Empty(t, req.Header.Get("Authorization"))
	})
}
