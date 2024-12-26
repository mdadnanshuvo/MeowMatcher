package channels

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWorkerPool(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return success response after first attempt
		images := []map[string]interface{}{
			{"url": "test1.jpg", "id": "1"},
			{"url": "test2.jpg", "id": "2"},
		}
		json.NewEncoder(w).Encode(images)
	}))
	defer server.Close()

	breedIDs := []string{"breed1", "breed2"}
	results, err := WorkerPool("test-api-key", server.URL, breedIDs, 2)

	if err != nil {
		t.Errorf("WorkerPool() error = %v", err)
	}
	if len(results) != len(breedIDs) {
		t.Errorf("Expected %d results, got %d", len(breedIDs), len(results))
	}
}

func TestFetchDataFromAPI(t *testing.T) {
	// Test server that simulates rate limiting then success
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount < 2 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		response := []map[string]interface{}{
			{"data": "test"},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	data, err := fetchDataFromAPI("test-key", server.URL, "/test", nil)
	if err != nil {
		t.Errorf("fetchDataFromAPI() error = %v", err)
	}
	if len(data) != 1 {
		t.Errorf("Expected 1 result, got %d", len(data))
	}
}

func TestFetchDataConcurrently(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := []map[string]interface{}{
			{"data": fmt.Sprintf("response for %s", r.URL.Path)},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	endpoints := map[string]map[string]string{
		"/test1": {"param": "value1"},
		"/test2": {"param": "value2"},
	}

	results, err := FetchDataConcurrently("test-key", server.URL, endpoints)
	if err != nil {
		t.Errorf("FetchDataConcurrently() error = %v", err)
	}
	if len(results) != len(endpoints) {
		t.Errorf("Expected %d results, got %d", len(endpoints), len(results))
	}
}

func TestFetchBreedImages(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify API key is present
		if r.Header.Get("x-api-key") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Return successful response
		images := []map[string]interface{}{
			{"url": "test.jpg", "breed_id": "test-breed"},
		}
		json.NewEncoder(w).Encode(images)
	}))
	defer server.Close()

	images, err := fetchBreedImages("test-key", server.URL, "test-breed", 1)
	if err != nil {
		t.Errorf("fetchBreedImages() error = %v", err)
	}
	if len(images) != 1 {
		t.Errorf("Expected 1 image, got %d", len(images))
	}
}
