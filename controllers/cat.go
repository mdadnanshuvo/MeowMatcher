package controllers

import (
	"bytes"
	channels "catApiProject/Channels"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/beego/beego/v2/server/web"
)

type CatController struct {
	web.Controller
}

// Render the index page
func (c *CatController) Index() {
	c.TplName = "index.tpl"
	c.Data["Title"] = "Welcome to the Cat API"
	c.Data["Message"] = "Explore voting, breeds, and favorites!"
}

// VotingCats provides voting data as JSON
func (c *CatController) VotingCats() {
	apiKey := web.AppConfig.DefaultString("cat_api_key", "")
	baseURL := web.AppConfig.DefaultString("cat_api_base_url", "")

	// Define the endpoint for voting data
	endpoints := map[string]map[string]string{
		"/images/search": {"limit": "10"},
	}

	// Fetch voting data concurrently using the channels package
	data, err := channels.FetchDataConcurrently(apiKey, baseURL, endpoints)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
		return
	}

	// Return only the voting data as JSON
	c.Data["json"] = data["/images/search"]
	c.ServeJSON()
}

// Breeds provides breeds data as JSON
func (c *CatController) Breeds() {
	apiKey := web.AppConfig.DefaultString("cat_api_key", "")
	baseURL := web.AppConfig.DefaultString("cat_api_base_url", "")

	// Define the endpoint for breeds data
	endpoints := map[string]map[string]string{
		"/breeds": nil,
	}

	// Fetch breeds data concurrently using the channels package
	data, err := channels.FetchDataConcurrently(apiKey, baseURL, endpoints)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
		return
	}

	// Return only the breeds data as JSON
	c.Data["json"] = data["/breeds"]
	c.ServeJSON()
}

// Favorites provides a placeholder list of favorite cats
func (c *CatController) Favorites() {
	favorites := []map[string]string{
		{"id": "1", "url": "https://example.com/cat1.jpg"},
		{"id": "2", "url": "https://example.com/cat2.jpg"},
	}

	// Return favorites as JSON
	c.Data["json"] = favorites
	c.ServeJSON()
}

// AddFavorite adds a new favorite cat via POST request
func (c *CatController) AddFavorite() {
	apiKey := web.AppConfig.DefaultString("cat_api_key", "")
	baseURL := web.AppConfig.DefaultString("cat_api_base_url", "")

	// Parse the incoming JSON request body
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

	// Prepare the data for the POST request
	favoriteData := map[string]interface{}{
		"image_id": requestBody.ImageID,
		"sub_id":   requestBody.SubID,
	}

	// Make the POST request to the Cat API
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

	// Return the response from the Cat API as JSON
	c.Data["json"] = result
	c.ServeJSON()
}
