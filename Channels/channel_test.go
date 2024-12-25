package channels

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var fetchBreedImagesWrapper = fetchBreedImages
var originalFetchBreedImages = fetchBreedImages

func mockableFetchBreedImages(apiKey, baseURL, breedID string, limit int) ([]map[string]interface{}, error) {
	// Use the mock implementation for testing
	return originalFetchBreedImages(apiKey, baseURL, breedID, limit)
}

func TestWorkerPoolWithRetries(t *testing.T) {
	apiKey := "test-api-key"
	baseURL := "https://mockapi.com"
	breedIDs := []string{"error-breed"}
	limit := 2

	// Mock retry behavior
	retryCount := 0
	mockFetch := func(apiKey, baseURL, breedID string, limit int) ([]map[string]interface{}, error) {
		retryCount++
		if retryCount < 3 {
			return nil, errors.New("mock error")
		}
		return []map[string]interface{}{
			{"url": "https://example.com/image1.jpg"},
		}, nil
	}

	// Redirect the wrapper to the mock implementation
	defer func() { fetchBreedImagesWrapper = fetchBreedImages }()
	fetchBreedImagesWrapper = mockFetch

	results, err := WorkerPool(apiKey, baseURL, breedIDs, limit)
	assert.NoError(t, err, "Expected no error after retries")
	assert.Len(t, results, 1, "Expected one result after retries")
	assert.Equal(t, "error-breed", results[0]["breed_id"])
}

func TestFetchDataFromAPIWithRateLimiting(t *testing.T) {
	retryCount := 0
	r := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		retryCount++
		if retryCount < 3 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"key":"value"}]`))
	}))
	defer r.Close()

	apiKey := "test-api-key"
	data, err := fetchDataFromAPI(apiKey, r.URL, "", nil)
	assert.NoError(t, err, "Expected no error after retries")
	assert.Len(t, data, 1, "Expected one item after retries")
	assert.Equal(t, "value", data[0]["key"])
}

func TestFetchDataConcurrently(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond) // Simulate network delay
		json.NewEncoder(w).Encode([]map[string]interface{}{{"test": "data"}})
	}))
	defer server.Close()

	endpoints := map[string]map[string]string{
		"/endpoint1": {"param": "value1"},
		"/endpoint2": {"param": "value2"},
	}

	result, err := FetchDataConcurrently("test-key", server.URL, endpoints)
	if err != nil {
		t.Errorf("FetchDataConcurrently() error = %v", err)
	}
	if len(result) != len(endpoints) {
		t.Errorf("FetchDataConcurrently() got %v results, want %v", len(result), len(endpoints))
	}
}

func TestFetchBreedImages(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-api-key") != "test-key" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"url": "test.jpg", "breed_id": "test-breed"},
		})
	}))
	defer server.Close()

	images, err := fetchBreedImages("test-key", server.URL, "test-breed", 1)
	if err != nil {
		t.Errorf("fetchBreedImages() error = %v", err)
	}
	if len(images) == 0 {
		t.Error("fetchBreedImages() returned no images")
	}
}
