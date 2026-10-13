package handlertestutil

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNewTestRouter(t *testing.T) {
	t.Parallel()
	type args struct {
		setup func(*gin.Engine)
	}
	tests := []struct {
		name string
		args args
		want *gin.Engine
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := NewTestRouter(tt.args.setup); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTestRouter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPerformRequest(t *testing.T) {
	t.Parallel()
	type args struct {
		r       *gin.Engine
		method  string
		path    string
		body    string
		headers map[string]string
	}
	tests := []struct {
		name string
		args args
		want *httptest.ResponseRecorder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := PerformRequest(tt.args.r, tt.args.method, tt.args.path, tt.args.body, tt.args.headers); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PerformRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthHeader(t *testing.T) {
	t.Parallel()
	type args struct {
		token string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := AuthHeader(tt.args.token); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONBody(t *testing.T) {
	t.Parallel()
	type args struct {
		method   string
		path     string
		jsonBody string
		token    string
	}
	tests := []struct {
		name  string
		args  args
		want  *http.Request
		want1 *httptest.ResponseRecorder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, got1 := JSONBody(tt.args.method, tt.args.path, tt.args.jsonBody, tt.args.token)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSONBody() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("JSONBody() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
