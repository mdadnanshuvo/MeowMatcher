package controllers

import (
	"bytes"
	channels "catApiProject/Channels"
	cache "catApiProject/caches"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/beego/beego/v2/server/web"
)

type CatController struct {
	web.Controller
}

var breedsCache = cache.NewCache(10 * time.Minute)

// Index renders the homepage
func (c *CatController) Index() {
	c.TplName = "index.tpl"
	c.Data["Title"] = "Welcome to the Cat API"
	c.Data["Message"] = "Explore voting, breeds, and favorites!"
}

// VotingCats fetches voting data
func (c *CatController) VotingCats() {
	apiKey := web.AppConfig.DefaultString("cat_api_key", "")
	baseURL := web.AppConfig.DefaultString("cat_api_base_url", "")

	data, err := channels.FetchDataConcurrently(apiKey, baseURL, map[string]map[string]string{
		"/images/search": {"limit": "10", "order": "RAND"},
	})
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
		return
	}

	c.Data["json"] = data["/images/search"]
	c.ServeJSON()
}

// Breeds fetches all breeds
func (c *CatController) Breeds() {
	apiKey := web.AppConfig.DefaultString("cat_api_key", "")
	baseURL := web.AppConfig.DefaultString("cat_api_base_url", "")

	data, err := channels.FetchDataConcurrently(apiKey, baseURL, map[string]map[string]string{
		"/breeds": nil,
	})
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
		return
	}

	c.Data["json"] = data["/breeds"]
	c.ServeJSON()
}

// BreedsWithImages fetches breeds with associated images
func (c *CatController) BreedsWithImages() {
	if cachedData, found := breedsCache.Get("breeds_with_images"); found {
		if isValidCache(cachedData) {
			c.Data["json"] = cachedData
			c.ServeJSON()
			return
		}
	}

	apiKey := web.AppConfig.DefaultString("cat_api_key", "")
	baseURL := web.AppConfig.DefaultString("cat_api_base_url", "")

	// Fetch all breeds
	breedsData, err := channels.FetchDataConcurrently(apiKey, baseURL, map[string]map[string]string{
		"/breeds": nil,
	})
	if err != nil {
		fmt.Println("Error fetching breeds:", err)
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to fetch breeds"}
		c.ServeJSON()
		return
	}

	breeds, ok := breedsData["/breeds"].([]map[string]interface{})
	if !ok {
		fmt.Println("Invalid breeds data format")
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Invalid breeds data format"}
		c.ServeJSON()
		return
	}

	// Fetch images with worker pool
	breedIDs := make([]string, len(breeds))
	for i, breed := range breeds {
		breedIDs[i] = breed["id"].(string)
	}

	imagesData, err := channels.WorkerPool(apiKey, baseURL, breedIDs, 10)
	if err != nil {
		fmt.Println("Error fetching images:", err)
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to fetch images"}
		c.ServeJSON()
		return
	}

	// Map images to breeds
	for i, breed := range breeds {
		breedID := breed["id"].(string)
		breedImages := []map[string]interface{}{}

		for _, imgData := range imagesData {
			if imgData["breed_id"] == breedID {
				breedImages = append(breedImages, imgData["images"].([]map[string]interface{})...)
			}
		}

		if len(breedImages) == 0 {
			fmt.Printf("No images mapped for breed: %s\n", breedID)
		}

		breed["images"] = breedImages
		breeds[i] = breed
	}

	breedsCache.Set("breeds_with_images", breeds)
	c.Data["json"] = breeds
	c.ServeJSON()
}

func isValidCache(data interface{}) bool {
	breeds, ok := data.([]map[string]interface{})
	if !ok {
		return false
	}

	for _, breed := range breeds {
		if breed["images"] == nil {
			return false
		}
	}
	return true
}

func (c *CatController) AddToFavorites() {
	apiKey := web.AppConfig.DefaultString("cat_api_key", "")
	subID := c.GetString("sub_id") // Optional sub_id parameter

	// Parse incoming JSON payload
	body, err := io.ReadAll(c.Ctx.Request.Body) // Ensure raw body is read for validation
	if err != nil {
		fmt.Println("Failed to read request body:", err)
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Unable to read request body"}
		c.ServeJSON()
		return
	}

	var favorite map[string]string
	if err := json.Unmarshal(body, &favorite); err != nil {
		fmt.Println("Failed to parse JSON payload:", err)
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid JSON format. Ensure 'image_id' is included."}
		c.ServeJSON()
		return
	}

	imageID, exists := favorite["image_id"]
	if !exists || imageID == "" {
		fmt.Println("Missing 'image_id' in the request body")
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Missing 'image_id' in the request body"}
		c.ServeJSON()
		return
	}

	// Construct the API URL
	url := "https://api.thecatapi.com/v1/favourites"
	if subID != "" {
		url += "?sub_id=" + subID
	}

	// Create the payload
	payload := map[string]string{"image_id": imageID}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Failed to marshal payload:", err)
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to encode favorite data"}
		c.ServeJSON()
		return
	}

	// Create and send the HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Failed to create HTTP request:", err)
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to create request"}
		c.ServeJSON()
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("HTTP request failed:", err)
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to submit favorite"}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Println("TheCatAPI Error Response:", string(bodyBytes))
		c.Ctx.Output.SetStatus(resp.StatusCode)
		c.Data["json"] = map[string]string{"error": "Failed to add to favorites. External API error."}
		c.ServeJSON()
		return
	}

	// Successfully added to favorites
	c.Data["json"] = map[string]string{"message": "Added to favorites successfully"}
	c.ServeJSON()
}

// GetFavorites retrieves a user's favorite cats using their sub_id passed as a query parameter
func (c *CatController) GetFavorites() {
	apiKey := web.AppConfig.DefaultString("cat_api_key", "")

	// Retrieve sub_id from the query parameters (passed from the frontend)
	subID := c.GetString("sub_id")
	if subID == "" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Missing 'sub_id' in the request"}
		c.ServeJSON()
		return
	}

	// Construct TheCatAPI URL with query parameters
	url := fmt.Sprintf("https://api.thecatapi.com/v1/favourites?sub_id=%s&limit=20&order=DESC", subID)

	// Create and send the HTTP GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Failed to create HTTP request:", err)
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to create request"}
		c.ServeJSON()
		return
	}
	req.Header.Set("x-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("HTTP request failed:", err)
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to fetch favorites"}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Println("TheCatAPI Error Response:", string(bodyBytes))
		c.Ctx.Output.SetStatus(resp.StatusCode)
		c.Data["json"] = map[string]string{"error": "Failed to fetch favorites. External API error."}
		c.ServeJSON()
		return
	}

	// Read and return the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to read API response"}
		c.ServeJSON()
		return
	}

	// Parse the response body into JSON and return it
	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Ctx.Output.Header("Content-Type", "application/json")
	c.Data["json"] = json.RawMessage(bodyBytes)
	c.ServeJSON()
}
