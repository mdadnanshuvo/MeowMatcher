package channels

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// FetchDataConcurrently fetches multiple API endpoints concurrently using channels
func FetchDataConcurrently(apiKey, baseURL string, endpoints map[string]map[string]string) (map[string]interface{}, error) {
	// Channels for results and errors
	results := make(chan map[string]interface{})
	errChan := make(chan error)

	// Launch a goroutine for each endpoint
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

	// Aggregate results
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

// fetchDataFromAPI fetches data from a single API endpoint
func fetchDataFromAPI(apiKey, baseURL, endpoint string, params map[string]string) ([]map[string]interface{}, error) {
	client := &http.Client{}
	fullURL, _ := url.Parse(baseURL + endpoint)

	// Add query string parameters
	query := fullURL.Query()
	query.Add("api_key", apiKey)
	for key, value := range params {
		query.Add(key, value)
	}
	fullURL.RawQuery = query.Encode()

	req, _ := http.NewRequest("GET", fullURL.String(), nil)
	req.Header.Set("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-OK response: %d", resp.StatusCode)
	}

	var data []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}
	return data, nil
}
