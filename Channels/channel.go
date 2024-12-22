package channels

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// fetchDataFromAPI fetches data from a single API endpoint
func fetchDataFromAPI(apiKey, baseURL, endpoint string, params map[string]string) ([]map[string]interface{}, error) {
	client := &http.Client{}
	fullURL, _ := url.Parse(baseURL + endpoint)

	query := fullURL.Query()
	query.Add("api_key", apiKey) // Add the API key to the query parameters
	for key, value := range params {
		query.Add(key, value)
	}
	fullURL.RawQuery = query.Encode()

	var response []map[string]interface{}

	for retries := 0; retries < 5; retries++ {
		req, _ := http.NewRequest("GET", fullURL.String(), nil)
		req.Header.Set("x-api-key", apiKey)

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("request error: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests {
			// Handle rate limiting with exponential backoff
			retryAfter := math.Pow(2, float64(retries))
			time.Sleep(time.Duration(retryAfter) * time.Second)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("non-OK response: %d, endpoint: %s", resp.StatusCode, endpoint)
		}

		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, fmt.Errorf("failed to decode response from endpoint %s: %v", endpoint, err)
		}
		return response, nil
	}

	return nil, fmt.Errorf("exceeded retry limit for endpoint: %s", endpoint)
}

// FetchDataConcurrently fetches multiple API endpoints concurrently
func FetchDataConcurrently(apiKey, baseURL string, endpoints map[string]map[string]string) (map[string]interface{}, error) {
	results := make(chan map[string]interface{})
	errChan := make(chan error)

	for key, params := range endpoints {
		go func(key string, params map[string]string) {
			data, err := fetchDataFromAPI(apiKey, baseURL, key, params)
			if err != nil {
				errChan <- fmt.Errorf("failed to fetch %s: %v", key, err)
				return
			}
			results <- map[string]interface{}{key: data}
		}(key, params)
	}

	response := make(map[string]interface{})
	for i := 0; i < len(endpoints); i++ {
		select {
		case res := <-results:
			for k, v := range res {
				response[k] = v
			}
		case err := <-errChan:
			return nil, err
		}
	}

	return response, nil
}

// WorkerPool fetches data concurrently with limited workers
func WorkerPool(apiKey, baseURL string, breedIDs []string, limit int) ([]map[string]interface{}, error) {
	var wg sync.WaitGroup
	results := make(chan map[string]interface{}, len(breedIDs))
	maxRetries := 3

	for _, breedID := range breedIDs {
		wg.Add(1)
		go func(breedID string) {
			defer wg.Done()

			for retries := 0; retries < maxRetries; retries++ {
				images, err := fetchBreedImages(apiKey, baseURL, breedID, limit)
				if err != nil {
					fmt.Printf("Retry %d for breed %s failed: %v\n", retries+1, breedID, err)
					time.Sleep(2 * time.Second)
					continue
				}

				results <- map[string]interface{}{"breed_id": breedID, "images": images}
				return
			}
			fmt.Printf("Failed to fetch images for breed: %s after %d retries\n", breedID, maxRetries)
		}(breedID)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var allResults []map[string]interface{}
	for res := range results {
		allResults = append(allResults, res)
	}

	return allResults, nil
}

// fetchBreedImages fetches images for a specific breed
func fetchBreedImages(apiKey, baseURL, breedID string, limit int) ([]map[string]interface{}, error) {
	client := &http.Client{}
	endpoint := fmt.Sprintf("%s/images/search?limit=%d&breed_ids=%s", baseURL, limit, breedID)

	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Set("x-api-key", apiKey)

	fmt.Printf("Fetching images for breed: %s, URL: %s\n", breedID, endpoint)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch images for breed %s: %v", breedID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-OK response: %d for breed %s", resp.StatusCode, breedID)
	}

	var images []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&images); err != nil {
		return nil, fmt.Errorf("failed to decode images for breed %s: %v", breedID, err)
	}

	return images, nil
}
