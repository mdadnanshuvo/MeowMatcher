package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beego/beego"
	"github.com/beego/beego/v2/server/web"
	"github.com/magiconair/properties/assert"
)

func init() {
	web.BConfig.RunMode = "test"
}

func TestRoutes(t *testing.T) {
	routes := map[string]string{
		"/":                   "GET",
		"/voting":             "GET",
		"/add-favourites":     "POST",
		"/get-favourites":     "GET",
		"/delete-favourites":  "DELETE",
		"/breeds-with-images": "GET",
		"/vote":               "POST",
		"/votes":              "GET",
	}

	for route, method := range routes {
		t.Run(route, func(t *testing.T) {
			r := httptest.NewRequest(method, route, nil)
			w := httptest.NewRecorder()

			// Add authentication headers if needed
			if route == "/delete-favourites" || route == "/vote" {
				r.Header.Set("Authorization", "Bearer test-token")
			}

			beego.BeeApp.Handlers.ServeHTTP(w, r)
			assert.Equal(t, http.StatusOK, w.Code, "Expected status 200 for route %s", route)
		})
	}
}

func TestMethodNotAllowed(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "POST to GET route",
			method:         "POST",
			path:           "/voting",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "GET to POST route",
			method:         "GET",
			path:           "/vote",
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			web.BeeApp.Handlers.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d for method %s on route %s",
					tt.expectedStatus, w.Code, tt.method, tt.path)
			}
		})
	}
}
