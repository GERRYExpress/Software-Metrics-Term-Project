package auth

import (
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func Test_secret(t *testing.T) {
	tests := []struct {
		name    string
		setup   func()
		want    []byte
		wantErr bool
	}{
		{
			name: "secret is configured",
			setup: func() {
				os.Setenv("JWT_SECRET", "supersecretkey")
			},
			want:    []byte("supersecretkey"),
			wantErr: false,
		},
		{
			name: "secret is missing",
			setup: func() {
				os.Unsetenv("JWT_SECRET")
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got, err := secret()
			if (err != nil) != tt.wantErr {
				t.Fatalf("secret() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("secret() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIssue(t *testing.T) {
	os.Setenv("JWT_SECRET", "supersecretkey")
	defer os.Unsetenv("JWT_SECRET")

	type args struct {
		accountID int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "successfully issue token for account ID 1",
			args:    args{accountID: 1},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Issue(tt.args.accountID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Issue() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got == "" {
				t.Errorf("Issue() returned empty token string")
			}
		})
	}
}

func TestMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("JWT_SECRET", "supersecretkey")
	defer os.Unsetenv("JWT_SECRET")

	validToken, _ := Issue(42)

	tests := []struct {
		name           string
		authHeader     string
		wantStatusCode int
	}{
		{
			name:           "valid token",
			authHeader:     "Bearer " + validToken,
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "missing bearer prefix",
			authHeader:     validToken,
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name:           "empty header",
			authHeader:     "",
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name:           "invalid token signature",
			authHeader:     "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI.invalid.signature",
			wantStatusCode: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(Middleware())
			r.GET("/test", func(c *gin.Context) {
				// Verify the middleware successfully extracted and set the accountID
				if id, exists := c.Get("accountID"); exists && id == 42 {
					c.Status(http.StatusOK)
					return
				}
				c.Status(http.StatusInternalServerError)
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("Middleware() expected status %v, got %v", tt.wantStatusCode, w.Code)
			}
		})
	}
}
