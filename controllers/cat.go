package controllers

import (
	channels "catApiProject/Channels"
	cache "catApiProject/caches"
	"fmt"
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
