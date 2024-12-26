package controllers

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func init() {
	// Attempt to load the configuration for tests
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	err := web.LoadAppConfig("ini", "conf/app.conf")
	if err != nil {
		fmt.Println("Error loading configuration during tests. Falling back to manual configuration:", err)

		// Fallback to manual configuration
		web.AppConfig.Set("appname", "catApiProject")
		web.AppConfig.Set("httpport", "8080")
		web.AppConfig.Set("runmode", "test")
		web.AppConfig.Set("cat_api_key", "testApiKey")
		web.AppConfig.Set("cat_api_base_url", "https://mockapi.com/v1")
		web.AppConfig.Set("cat_api_sub_id", "testSubID")
	}
}

func setupTestController() (*CatController, *context.Context) {

	controller := &CatController{}
	ctx := context.NewContext()
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	ctx.Reset(w, req)
	controller.Init(ctx, "", "", nil)
	return controller, ctx
}

func setupMockServer(handler http.HandlerFunc) *httptest.Server {
	server := httptest.NewServer(handler)
	web.AppConfig.Set("cat_api_base_url", server.URL)
	web.AppConfig.Set("cat_api_key", "test-key")
	web.AppConfig.Set("cat_api_sub_id", "test-sub")
	return server
}

// Test Index with multiple data scenarios
func TestIndex(t *testing.T) {
	// Test basic index
	controller, _ := setupTestController()
	controller.Index()
	if controller.Data["Title"] != "Welcome to the Cat API" {
		t.Error("Expected title not set correctly")
	}

	// Test with TplName verification
	if controller.TplName != "index.tpl" {
		t.Error("Expected template name not set correctly")
	}

	// Test message content
	if controller.Data["Message"] != "Explore voting, breeds, and favorites!" {
		t.Error("Expected message not set correctly")
	}
}

// Test VotingCats with multiple scenarios
func TestVotingCats(t *testing.T) {
	// Test successful case with multiple images
	server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"id": "cat1", "url": "url1"},
			{"id": "cat2", "url": "url2"},
		})
	})
	defer server.Close()

	controller, _ := setupTestController()
	controller.VotingCats()

	// Test with empty response
	server = setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]map[string]interface{}{})
	})
	defer server.Close()

	controller, _ = setupTestController()
	controller.VotingCats()

	// Test with malformed JSON response
	server = setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{malformed json"))
	})
	defer server.Close()

	controller, _ = setupTestController()
	controller.VotingCats()

	// Test with network error (invalid URL)
	web.AppConfig.Set("cat_api_base_url", "http://invalid-url")
	controller, _ = setupTestController()
	controller.VotingCats()
}

// Test BreedsWithImages with cache scenarios
func TestBreedsWithImages(t *testing.T) {
	// Test with valid cache
	validCache := []map[string]interface{}{
		{
			"id":     "breed1",
			"images": []string{"img1"},
		},
	}
	breedsCache.Set("breeds_with_images", validCache)
	controller, _ := setupTestController()
	controller.BreedsWithImages()

	// Test with invalid cache format
	breedsCache.Set("breeds_with_images", "invalid cache")
	controller, _ = setupTestController()
	controller.BreedsWithImages()

	// Test successful API call with multiple breeds
	server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/breeds") {
			json.NewEncoder(w).Encode([]map[string]interface{}{
				{"id": "breed1", "name": "Test Breed 1"},
				{"id": "breed2", "name": "Test Breed 2"},
			})
		} else if strings.Contains(r.URL.Path, "/images/search") {
			json.NewEncoder(w).Encode([]map[string]interface{}{
				{"breed_id": "breed1", "images": []map[string]interface{}{{"url": "url1"}}},
			})
		}
	})
	defer server.Close()

	controller, _ = setupTestController()
	controller.BreedsWithImages()

	// Test with invalid breeds response
	server = setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/breeds") {
			w.Write([]byte("invalid json"))
		}
	})
	defer server.Close()

	controller, _ = setupTestController()
	controller.BreedsWithImages()
}

// Enhanced AddToFavorites test with more scenarios
func TestAddToFavorites(t *testing.T) {
	testCases := []struct {
		name    string
		payload string
		setup   func(*httptest.Server)
	}{
		{
			name:    "Valid new favorite",
			payload: `{"image_id":"test123"}`,
			setup: func(s *httptest.Server) {
				web.AppConfig.Set("cat_api_sub_id", "test-sub")
			},
		},
		{
			name:    "Missing sub_id",
			payload: `{"image_id":"test123"}`,
			setup: func(s *httptest.Server) {
				web.AppConfig.Set("cat_api_sub_id", "")
			},
		},
		{
			name:    "Invalid JSON",
			payload: `{"invalid json`,
			setup: func(s *httptest.Server) {
				web.AppConfig.Set("cat_api_sub_id", "test-sub")
			},
		},
		{
			name:    "Empty image_id",
			payload: `{"image_id":""}`,
			setup:   func(s *httptest.Server) {},
		},
		{
			name:    "Duplicate favorite",
			payload: `{"image_id":"existing"}`,
			setup: func(s *httptest.Server) {
				web.AppConfig.Set("cat_api_sub_id", "test-sub")
			},
		},
	}

	for _, tc := range testCases {
		server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" && tc.name == "Duplicate favorite" {
				json.NewEncoder(w).Encode([]map[string]interface{}{
					{"image_id": "existing"},
				})
			} else if r.Method == "POST" {
				w.WriteHeader(http.StatusCreated)
			}
		})
		tc.setup(server)
		defer server.Close()

		controller, ctx := setupTestController()
		ctx.Request = httptest.NewRequest("POST", "/", strings.NewReader(tc.payload))
		controller.AddToFavorites()
	}
}

// Enhanced GetFavorites test with cache scenarios
func TestGetFavorites(t *testing.T) {
	// Test with valid cache
	validCache := []map[string]interface{}{
		{
			"id":    "fav1",
			"image": map[string]interface{}{"url": "test.jpg"},
		},
	}
	favoritesCache.Set("user_favorites", validCache)
	controller, _ := setupTestController()
	controller.GetFavorites()

	// Test with invalid cache
	favoritesCache.Set("user_favorites", "invalid cache")
	controller, _ = setupTestController()
	controller.GetFavorites()

	// Test with missing sub_id
	web.AppConfig.Set("cat_api_sub_id", "")
	controller, _ = setupTestController()
	controller.GetFavorites()

	// Test with network error
	web.AppConfig.Set("cat_api_base_url", "http://invalid-url")
	controller, _ = setupTestController()
	controller.GetFavorites()

	// Test with invalid response
	server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("invalid json"))
	})
	defer server.Close()
	controller, _ = setupTestController()
	controller.GetFavorites()
}

// Enhanced DeleteFavorite test with more scenarios
func TestDeleteFavorite(t *testing.T) {
	testCases := []struct {
		name   string
		favID  string
		status int
	}{
		{"Success", "valid-id", http.StatusOK},
		{"Missing ID", "", http.StatusBadRequest},
		{"Server Error", "error-id", http.StatusInternalServerError},
		{"Not Found", "not-found", http.StatusNotFound},
		{"Invalid Request", "invalid-req", http.StatusBadRequest},
	}

	for _, tc := range testCases {
		server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(tc.status)
			if tc.status == http.StatusInternalServerError {
				w.Write([]byte("Internal Server Error"))
			}
		})
		defer server.Close()

		controller, ctx := setupTestController()
		if tc.favID != "" {
			ctx.Input.SetParam(":id", tc.favID)
		}
		controller.DeleteFavorite()
	}
}

// Enhanced PostVote test with more scenarios
func TestPostVote(t *testing.T) {
	testCases := []struct {
		name    string
		payload string
		setup   func()
	}{
		{
			name:    "Valid upvote",
			payload: `{"image_id":"test123","value":1}`,
			setup:   func() { web.AppConfig.Set("cat_api_sub_id", "test-sub") },
		},
		{
			name:    "Valid downvote",
			payload: `{"image_id":"test123","value":-1}`,
			setup:   func() { web.AppConfig.Set("cat_api_sub_id", "test-sub") },
		},
		{
			name:    "Missing sub_id",
			payload: `{"image_id":"test123","value":1}`,
			setup:   func() { web.AppConfig.Set("cat_api_sub_id", "") },
		},
		{
			name:    "Invalid JSON",
			payload: `{"invalid json`,
			setup:   func() {},
		},
		{
			name:    "Empty image_id",
			payload: `{"image_id":"","value":1}`,
			setup:   func() {},
		},
		{
			name:    "Invalid value",
			payload: `{"image_id":"test123","value":2}`,
			setup:   func() {},
		},
		{
			name:    "Network error",
			payload: `{"image_id":"test123","value":1}`,
			setup:   func() { web.AppConfig.Set("cat_api_base_url", "http://invalid-url") },
		},
	}

	for _, tc := range testCases {
		tc.setup()
		server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(tc.payload, "test123") && (strings.Contains(tc.payload, `"value":1`) || strings.Contains(tc.payload, `"value":-1`)) {
				w.WriteHeader(http.StatusCreated)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		})
		defer server.Close()

		controller, ctx := setupTestController()
		ctx.Request = httptest.NewRequest("POST", "/", strings.NewReader(tc.payload))
		controller.PostVote()
	}
}

// Enhanced GetVotes test with more scenarios
func TestGetVotes(t *testing.T) {
	testCases := []struct {
		name       string
		subID      string
		order      string
		setupError bool
	}{
		{"Valid request", "test-sub", "DESC", false},
		{"Empty sub_id", "", "DESC", false},
		{"Different order", "test-sub", "ASC", false},
		{"Network error", "test-sub", "DESC", true},
		{"Invalid response", "invalid", "DESC", false},
	}

	for _, tc := range testCases {
		if tc.setupError {
			web.AppConfig.Set("cat_api_base_url", "http://invalid-url")
		} else {
			server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
				if tc.subID == "invalid" {
					w.Write([]byte("invalid json"))
				} else {
					json.NewEncoder(w).Encode([]map[string]interface{}{
						{"id": "vote1", "value": 1},
					})
				}
			})
			defer server.Close()
		}

		controller, ctx := setupTestController()
		ctx.Request.Form = map[string][]string{
			"sub_id": {tc.subID},
			"order":  {tc.order},
		}
		controller.GetVotes()
	}
}

// Test cache validation with more edge cases
func TestCacheValidation(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		isValid  bool
		testType string
	}{
		{"Valid breed cache", []map[string]interface{}{{"id": "1", "images": []string{"img1"}}}, true, "breed"},
		{"Invalid breed cache", []map[string]interface{}{{"id": "1"}}, false, "breed"},
		{"Empty breed cache", []map[string]interface{}{}, true, "breed"},
		{"Nil breed cache", nil, false, "breed"},
		{"Invalid type breed", "invalid", false, "breed"},
		{"Valid favorite cache", []map[string]interface{}{{"id": "1", "image": map[string]interface{}{"url": "test.jpg"}}}, true, "favorite"},
		{"Invalid favorite cache", []map[string]interface{}{{"image": map[string]interface{}{"url": "test.jpg"}}}, false, "favorite"},
		{"Empty favorite cache", []map[string]interface{}{}, true, "favorite"},
		{"Nil favorite cache", nil, false, "favorite"},
		{"Invalid type favorite", "invalid", false, "favorite"},
		{"Map type breed", map[string]interface{}{"test": "value"}, false, "breed"},
		{"Map type favorite", map[string]interface{}{"test": "value"}, false, "favorite"},
	}

	for _, tt := range tests {
		if tt.testType == "breed" {
			if isValidCache(tt.data) != tt.isValid {
				t.Errorf("%s: expected %v, got %v", tt.name, tt.isValid, !tt.isValid)
			}
		} else {
			if isValidCacheForFav(tt.data) != tt.isValid {
				t.Errorf("%s: expected %v, got %v", tt.name, tt.isValid, !tt.isValid)
			}
		}
	}
}
