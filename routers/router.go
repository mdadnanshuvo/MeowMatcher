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

	// Favorites Data
	beego.Router("/favorites", &controllers.CatController{}, "get:Favorites")

	// Add Favorite
	beego.Router("/addfavorite", &controllers.CatController{}, "post:AddFavorite")
}
