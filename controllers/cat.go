package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/beego/beego/v2/server/web"
)

type CatController struct {
	web.Controller
}

// Render Voting Data as JSON
func (c *CatController) VotingCats() {
	apiKey := web.AppConfig.DefaultString("cat_api_key", "")
	baseURL := web.AppConfig.DefaultString("cat_api_base_url", "")

	// Fetch data from the Cat API
	data, err := fetchDataFromAPI(apiKey, baseURL, "/images/search", map[string]string{
		"limit": "10",
	})
	if err != nil {
		// Return an error response as JSON
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to fetch voting data"}
		c.ServeJSON()
		return
	}

	// Return the fetched data as JSON
	c.Data["json"] = data
	c.ServeJSON()
}

// Render Voting Page with HTML
func (c *CatController) VotingCatsHTML() {
	apiKey := web.AppConfig.DefaultString("cat_api_key", "")
	baseURL := web.AppConfig.DefaultString("cat_api_base_url", "")

	// Fetch data from the Cat API
	data, err := fetchDataFromAPI(apiKey, baseURL, "/images/search", map[string]string{
		"limit": "10",
	})
	if err != nil {
		// Pass error message to template
		fmt.Println("Error fetching voting data:", err)
		c.Data["error"] = "Failed to fetch voting data. Please try again later."
	} else {
		// Pass data to template
		c.Data["votingCats"] = data
	}

	// Render the voting template
	c.TplName = "index.tpl"
}

// Render Breeds Data as JSON
func (c *CatController) Breeds() {
	apiKey := web.AppConfig.DefaultString("cat_api_key", "")
	baseURL := web.AppConfig.DefaultString("cat_api_base_url", "")

	// Fetch data from the Cat API
	data, err := fetchDataFromAPI(apiKey, baseURL, "/breeds", nil)
	if err != nil {
		// Return an error response as JSON
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to fetch breeds"}
		c.ServeJSON()
		return
	}

	// Return the fetched data as JSON
	c.Data["json"] = data
	c.ServeJSON()
}

// Render Breeds Page with HTML
func (c *CatController) BreedsHTML() {
	apiKey := web.AppConfig.DefaultString("cat_api_key", "")
	baseURL := web.AppConfig.DefaultString("cat_api_base_url", "")

	// Fetch data from the Cat API
	data, err := fetchDataFromAPI(apiKey, baseURL, "/breeds", nil)
	if err != nil {
		// Pass error message to template
		fmt.Println("Error fetching breeds:", err)
		c.Data["error"] = "Failed to fetch breeds. Please try again later."
	} else {
		// Pass data to template
		c.Data["breeds"] = data
	}

	// Render the breeds template
	c.TplName = "breeds.tpl"
}

// Render Favorites Data as JSON
func (c *CatController) Favorites() {
	// Placeholder logic for fetching favorites
	favorites := []map[string]string{
		{"id": "1", "url": "https://example.com/cat1.jpg"},
		{"id": "2", "url": "https://example.com/cat2.jpg"},
	}

	// Return favorites as JSON
	c.Data["json"] = favorites
	c.ServeJSON()
}

// Render Favorites Page with HTML
func (c *CatController) FavoritesHTML() {
	// Placeholder logic for fetching favorites
	favorites := []map[string]string{
		{"id": "1", "url": "https://example.com/cat1.jpg"},
		{"id": "2", "url": "https://example.com/cat2.jpg"},
	}

	// Pass favorites to template
	c.Data["favorites"] = favorites

	// Render the favorites template
	c.TplName = "favorites.tpl"
}

// Utility function to fetch data from the Cat API
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

// Add a Favorite via POST method
func (c *CatController) AddFavorite() {
	apiKey := web.AppConfig.DefaultString("cat_api_key", "")
	baseURL := web.AppConfig.DefaultString("cat_api_base_url", "")

	// Get the image_id and sub_id from the POST request (JSON body)
	var requestBody struct {
		ImageID string `json:"image_id"`
		SubID   string `json:"sub_id,omitempty"`
	}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &requestBody); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid request body"}
		c.ServeJSON()
		return
	}

	// Prepare the data to be sent to the Cat API
	favoriteData := map[string]interface{}{
		"image_id": requestBody.ImageID,
		"sub_id":   requestBody.SubID,
	}

	// Create a new favorite by sending a POST request to the Cat API
	client := &http.Client{}
	body, err := json.Marshal(favoriteData)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to encode request body"}
		c.ServeJSON()
		return
	}

	req, err := http.NewRequest("POST", baseURL+"/favourites", bytes.NewBuffer(body))
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to create request"}
		c.ServeJSON()
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to connect to Cat API"}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": fmt.Sprintf("Failed to create favorite. Status: %d", resp.StatusCode)}
		c.ServeJSON()
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to decode response from Cat API"}
		c.ServeJSON()
		return
	}

	// Return the newly created favorite ID as JSON
	c.Data["json"] = result
	c.ServeJSON()
}
