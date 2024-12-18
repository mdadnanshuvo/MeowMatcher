package routers

import (
	"catApiProject/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// Route for the home page
	beego.Router("/", &controllers.MainController{})

	// Route for fetching random cat images (Voting Tab)
	beego.Router("/voting", &controllers.CatController{}, "get:GetVotingCats")

	// Route for listing cat breeds (Breeds Tab)
	beego.Router("/breeds", &controllers.CatController{}, "get:GetBreeds")

	// Routes for managing favorites (Favorites Tab)
	beego.Router("/favorites", &controllers.CatController{}, "get:GetFavorites;post:SaveFavorite")
}
