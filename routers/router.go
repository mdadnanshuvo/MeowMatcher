package routers

import (
	"catApiProject/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// Route for the home page (index page should show voting)
	beego.Router("/", &controllers.CatController{}, "get:RenderVotingPage")

	// Route for fetching random cat images (Voting Tab, optional if /voting is kept)
	beego.Router("/voting", &controllers.CatController{}, "get:RenderVotingPage")

	// Route for listing cat breeds (Breeds Tab)
	beego.Router("/breeds", &controllers.CatController{}, "get:RenderBreedsPage")

	// Routes for managing favorites (Favorites Tab)
	beego.Router("/favorites", &controllers.CatController{}, "get:RenderFavoritesPage")
}
