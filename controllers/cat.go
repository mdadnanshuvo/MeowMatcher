package controllers

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/beego/beego/v2/server/web"
)

type CatController struct {
	web.Controller
}

func (c *CatController) RenderVotingPage() {
	c.TplName = "index.tpl"
}

// Render Breeds Page
func (c *CatController) RenderBreedsPage() {
	c.TplName = "breeds.tpl"
}

// Render Favorites Page
func (c *CatController) RenderFavoritesPage() {
	c.TplName = "favorites.tpl"
}

// Fetch random cat images for the Voting Tab
func (c *CatController) GetVotingCats() {
	apiKey := web.AppConfig.DefaultString("cat_api_key", "")
	baseURL := web.AppConfig.DefaultString("cat_api_base_url", "")

	data, err := fetchDataFromAPI(apiKey, baseURL, "/images/search", map[string]string{
		"limit": "10",
	})
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to fetch voting data"}
	} else {
		c.Data["json"] = data
	}
	c.ServeJSON()
}

// Fetch cat breeds for the Breeds Tab
func (c *CatController) GetBreeds() {
	apiKey := web.AppConfig.DefaultString("cat_api_key", "")
	baseURL := web.AppConfig.DefaultString("cat_api_base_url", "")

	data, err := fetchDataFromAPI(apiKey, baseURL, "/breeds", nil)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to fetch breeds"}
	} else {
		c.Data["json"] = data
	}
	c.ServeJSON()
}

// Fetch user's favorite cat images for the Favorites Tab
func (c *CatController) GetFavorites() {
	// Placeholder logic: Return an empty JSON array for now
	c.Data["json"] = []map[string]string{}
	c.ServeJSON()
}

// Save a favorite cat image for the Favorites Tab
func (c *CatController) SaveFavorite() {
	var favorite map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &favorite)
	if err != nil || favorite["image_id"] == nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid request payload"}
		return
	}

	// Placeholder logic for saving a favorite
	c.Data["json"] = map[string]string{"status": "Favorite saved", "image_id": favorite["image_id"].(string)}
	c.ServeJSON()
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

	var data []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)
	return data, nil
}
