package routers

import (
	"catApiProject/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// Index Route
	beego.Router("/", &controllers.CatController{}, "get:Index")

	// Voting Data
	beego.Router("/voting", &controllers.CatController{}, "get:VotingCats")

	// Breeds Data
	beego.Router("/breeds", &controllers.CatController{}, "get:Breeds")

	// Add to Favorites (POST)
	beego.Router("/add-favourites", &controllers.CatController{}, "post:AddToFavorites")

	// Get Favorites (GET)
	beego.Router("/get-favourites", &controllers.CatController{}, "get:GetFavorites")

	// Breeds with Images
	beego.Router("/breeds-with-images", &controllers.CatController{}, "get:BreedsWithImages")
}
