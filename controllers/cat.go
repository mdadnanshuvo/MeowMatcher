package controllers

import (
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
