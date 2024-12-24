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
	subID := web.AppConfig.DefaultString("cat_api_sub_id", "") // Retrieve sub_id from config file
	baseURL := web.AppConfig.DefaultString("cat_api_base_url", "https://api.thecatapi.com/v1")

	if subID == "" {
		fmt.Println("sub_id is not configured in app.conf")
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "sub_id is not configured on the server"}
		c.ServeJSON()
		return
	}

	// Parse incoming JSON payload
	body, err := io.ReadAll(c.Ctx.Request.Body)
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

	// Check if the image is already a favorite
	checkURL := fmt.Sprintf("%s/favourites?sub_id=%s", baseURL, subID)
	req, err := http.NewRequest("GET", checkURL, nil)
	if err != nil {
		fmt.Println("Failed to create GET request to check favorites:", err)
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to check existing favorites"}
		c.ServeJSON()
		return
	}

	req.Header.Set("x-api-key", apiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error checking existing favorites:", err)
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to check existing favorites"}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	// Parse the response to check if the image already exists
	var favorites []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&favorites)
	if err != nil {
		fmt.Println("Failed to decode existing favorites:", err)
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Error decoding existing favorites"}
		c.ServeJSON()
		return
	}

	for _, fav := range favorites {
		if fav["image_id"] == imageID {
			c.Ctx.Output.SetStatus(http.StatusConflict)
			c.Data["json"] = map[string]string{"error": "This image is already a favorite"}
			c.ServeJSON()
			return
		}
	}

	// Construct the API URL to add the favorite
	url := fmt.Sprintf("%s/favourites?sub_id=%s", baseURL, subID)

	// Create the payload
	payload := map[string]string{
		"image_id": imageID,
		"sub_id":   subID,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Failed to marshal payload:", err)
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to encode favorite data"}
		c.ServeJSON()
		return
	}

	// Send the POST request to TheCatAPI
	req, err = http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Failed to create POST request:", err)
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to create request"}
		c.ServeJSON()
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)

	resp, err = client.Do(req)
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

// GetFavorites retrieves a user's favorite cats using the sub_id from the configuration file
func (c *CatController) GetFavorites() {
	// Retrieve sub_id and API key from the configuration file
	subID := web.AppConfig.DefaultString("cat_api_sub_id", "")
	apiKey := web.AppConfig.DefaultString("cat_api_key", "")
	baseURL := web.AppConfig.DefaultString("cat_api_base_url", "https://api.thecatapi.com/v1")

	if subID == "" {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "sub_id is not configured on the server"}
		c.ServeJSON()
		return
	}

	// Construct the API URL
	url := fmt.Sprintf("%s/favourites?sub_id=%s", baseURL, subID)

	// Create the GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Error creating request to TheCatAPI"}
		c.ServeJSON()
		return
	}

	// Set the API key in the request headers
	req.Header.Set("x-api-key", apiKey)

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Error fetching favorites from TheCatAPI"}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.Ctx.Output.SetStatus(resp.StatusCode)
		c.Data["json"] = map[string]string{"error": fmt.Sprintf("Failed to fetch favorites from TheCatAPI: %s", string(body))}
		c.ServeJSON()
		return
	}

	// Read the response body and forward it to the client
	body, _ := io.ReadAll(resp.Body)
	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Ctx.Output.Header("Content-Type", "application/json")
	c.Ctx.WriteString(string(body))
}

// DeleteFavorite removes a favorite cat from TheCatAPI
func (c *CatController) DeleteFavorite() {
	apiKey := web.AppConfig.DefaultString("cat_api_key", "")
	baseURL := web.AppConfig.DefaultString("cat_api_base_url", "https://api.thecatapi.com/v1")

	// Get the favorite ID from the request parameters
	favID := c.Ctx.Input.Param(":id")
	if favID == "" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Favorite ID is required"}
		c.ServeJSON()
		return
	}

	// Construct the API URL for deleting the favorite
	url := fmt.Sprintf("%s/favourites/%s", baseURL, favID)

	// Create the DELETE request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Error creating request to TheCatAPI"}
		c.ServeJSON()
		return
	}

	req.Header.Set("x-api-key", apiKey)

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Error sending request to TheCatAPI"}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		c.Ctx.Output.SetStatus(resp.StatusCode)
		c.Data["json"] = map[string]string{"error": fmt.Sprintf("Failed to delete favorite: %s", string(body))}
		c.ServeJSON()
		return
	}

	// Successfully deleted the favorite
	c.Data["json"] = map[string]string{"message": "Favorite deleted successfully"}
	c.ServeJSON()
}
