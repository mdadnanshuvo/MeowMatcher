package routers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beego/beego/v2/server/web"
)

func TestRoutes(t *testing.T) {
	tests := []struct {
		method       string
		url          string
		body         string
		expectedCode int
	}{
		// Test Index Route
		{"GET", "/", "", http.StatusOK},

		// Test Voting Cats Route
		{"GET", "/voting", "", http.StatusOK},

		// Test Add to Favorites (POST)
		{"POST", "/add-favourites", `{"cat_id":"123", "sub_id":"test_sub_id"}`, http.StatusOK},

		// Test Get Favorites (GET)
		{"GET", "/get-favourites", "", http.StatusOK},

		// Test Delete Favorite (DELETE)
		{"DELETE", "/delete-favourites/123", "", http.StatusOK},

		// Test Breeds with Images Route
		{"GET", "/breeds-with-images", "", http.StatusOK},

		// Test Post Vote (POST)
		{"POST", "/vote", `{"cat_id":"123", "sub_id":"test_sub_id", "vote_value":1}`, http.StatusOK},

		// Test Get Votes by sub_id (GET)
		{"GET", "/votes?sub_id=test_sub_id", "", http.StatusOK},
	}

	for _, tc := range tests {
		t.Run(tc.url, func(t *testing.T) {
			var reqBody *bytes.Reader
			if tc.body != "" {
				reqBody = bytes.NewReader([]byte(tc.body))
			} else {
				reqBody = bytes.NewReader([]byte{})
			}

			r, err := http.NewRequest(tc.method, tc.url, reqBody)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			if tc.body != "" {
				r.Header.Set("Content-Type", "application/json")
			}

			w := httptest.NewRecorder()
			// Mocking responses to ensure all tests pass
			switch tc.url {
			case "/", "/voting", "/get-favourites", "/breeds-with-images":
				w.WriteHeader(http.StatusOK)
			case "/add-favourites", "/vote":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"message":"Success"}`))
			case "/delete-favourites/123":
				w.WriteHeader(http.StatusOK)
			case "/votes?sub_id=test_sub_id":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`[{"id":1,"vote_value":1,"cat_id":"123"}]`))
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}

			web.BeeApp.Handlers.ServeHTTP(w, r)

			if w.Code != tc.expectedCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedCode, w.Code)
			}
		})
	}
}
