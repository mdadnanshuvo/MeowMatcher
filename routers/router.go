package routers

import (
	"catApiProject/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// JSON Routes
	beego.Router("/voting", &controllers.CatController{}, "get:VotingCats")
	beego.Router("/breeds", &controllers.CatController{}, "get:Breeds")
	beego.Router("/favorites", &controllers.CatController{}, "get:Favorites")

	// HTML Routes
	beego.Router("/", &controllers.CatController{}, "get:VotingCatsHTML")              // Render voting page with HTML
	beego.Router("/breeds-html", &controllers.CatController{}, "get:BreedsHTML")       // Render breeds page with HTML
	beego.Router("/favorites-html", &controllers.CatController{}, "get:FavoritesHTML") // Render favorites page with HTML
}
