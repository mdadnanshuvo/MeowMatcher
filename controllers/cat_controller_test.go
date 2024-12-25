package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func setupTestController() (*CatController, *context.Context) {
	controller := &CatController{}
	ctx := context.NewContext()
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	ctx.Reset(w, req)
	controller.Init(ctx, "", "", nil)
	return controller, ctx
}

// setupMockServer creates a test server that returns predefined responses
func setupMockServer(handler http.HandlerFunc) *httptest.Server {
	server := httptest.NewServer(handler)
	web.AppConfig.Set("cat_api_base_url", server.URL)
	web.AppConfig.Set("cat_api_key", "test-key")
	web.AppConfig.Set("cat_api_sub_id", "test-sub")
	return server
}

func TestIndex(t *testing.T) {
	controller, _ := setupTestController()
	controller.Index()

	if controller.TplName != "index.tpl" {
		t.Errorf("Expected template name 'index.tpl', got %s", controller.TplName)
	}
}

func TestVotingCats(t *testing.T) {
	mockData := []map[string]interface{}{
		{"id": "test1", "url": "test.jpg"},
	}

	server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/images/search" {
			t.Errorf("Expected path /images/search, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(mockData)
	})
	defer server.Close()

	controller, _ := setupTestController()
	controller.VotingCats()

	response := controller.Data["json"]
	if response == nil {
		t.Fatal("Expected response data, got nil")
	}
}

func TestBreedsWithImages(t *testing.T) {
	breedData := []map[string]interface{}{
		{"id": "breed1", "name": "Test Breed"},
	}

	server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/breeds":
			json.NewEncoder(w).Encode(breedData)
		case "/images/search":
			json.NewEncoder(w).Encode([]map[string]interface{}{
				{"url": "test.jpg"},
			})
		}
	})
	defer server.Close()

	controller, _ := setupTestController()
	controller.BreedsWithImages()

	response := controller.Data["json"]
	if response == nil {
		t.Fatal("Expected response data, got nil")
	}
}

func TestAddToFavorites(t *testing.T) {
	server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			json.NewEncoder(w).Encode([]map[string]interface{}{})
		case "POST":
			w.WriteHeader(http.StatusCreated)
		}
	})
	defer server.Close()

	controller, ctx := setupTestController()
	payload := `{"image_id":"test123"}`
	ctx.Request = httptest.NewRequest("POST", "/", strings.NewReader(payload))

	controller.AddToFavorites()

	response := controller.Data["json"]
	if response == nil {
		t.Fatal("Expected response data, got nil")
	}
}

func TestGetFavorites(t *testing.T) {
	mockData := []map[string]interface{}{
		{"id": "fav1", "image": map[string]interface{}{"id": "img1"}},
	}

	server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(mockData)
	})
	defer server.Close()

	controller, _ := setupTestController()
	controller.GetFavorites()

	response := controller.Data["json"]
	if response == nil {
		t.Fatal("Expected response data, got nil")
	}
}

func TestDeleteFavorite(t *testing.T) {
	server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	controller, ctx := setupTestController()
	ctx.Input.SetParam(":id", "fav123")

	controller.DeleteFavorite()

	response := controller.Data["json"]
	if response == nil {
		t.Fatal("Expected response data, got nil")
	}
}

func TestPostVote(t *testing.T) {
	server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
	})
	defer server.Close()

	controller, ctx := setupTestController()
	payload := `{"image_id":"test123","value":1}`
	ctx.Request = httptest.NewRequest("POST", "/", strings.NewReader(payload))

	controller.PostVote()

	response := controller.Data["json"]
	if response == nil {
		t.Fatal("Expected response data, got nil")
	}
}

func TestGetVotes(t *testing.T) {
	mockData := []map[string]interface{}{
		{"id": "vote1", "value": 1},
	}

	server := setupMockServer(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(mockData)
	})
	defer server.Close()

	controller, ctx := setupTestController()
	ctx.Request.Form = map[string][]string{
		"sub_id": {"test-sub"},
		"order":  {"DESC"},
	}

	controller.GetVotes()

	response := controller.Data["json"]
	if response == nil {
		t.Fatal("Expected response data, got nil")
	}
}

func TestIsValidCache(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
		want bool
	}{
		{
			name: "valid_data",
			data: []map[string]interface{}{
				{"id": "1", "images": []string{"img1"}},
			},
			want: true,
		},
		{
			name: "invalid_data_no_images",
			data: []map[string]interface{}{
				{"id": "1"},
			},
			want: false,
		},
		{
			name: "invalid_data_wrong_type",
			data: "not a slice",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidCache(tt.data)
			if result != tt.want {
				t.Errorf("isValidCache() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestIsValidCacheForFav(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
		want bool
	}{
		{
			name: "valid_favorites",
			data: []map[string]interface{}{
				{"id": "1", "image": map[string]string{"url": "test.jpg"}},
			},
			want: true,
		},
		{
			name: "invalid_no_id",
			data: []map[string]interface{}{
				{"image": map[string]string{"url": "test.jpg"}},
			},
			want: false,
		},
		{
			name: "invalid_wrong_type",
			data: "not a slice",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidCacheForFav(tt.data)
			if result != tt.want {
				t.Errorf("isValidCacheForFav() = %v, want %v", result, tt.want)
			}
		})
	}
}
