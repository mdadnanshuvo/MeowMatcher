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

	// Add to Favorites (POST)
	beego.Router("/add-favourites", &controllers.CatController{}, "post:AddToFavorites")

	// Get Favorites (GET)
	beego.Router("/get-favourites", &controllers.CatController{}, "get:GetFavorites")

	// Delete Favorite (DELETE)
	beego.Router("/delete-favourites/:id", &controllers.CatController{}, "delete:DeleteFavorite")

	// Breeds with Images
	beego.Router("/breeds-with-images", &controllers.CatController{}, "get:BreedsWithImages")

	// Define the route for posting a vote (upvote or downvote)
	beego.Router("/vote", &controllers.CatController{}, "post:PostVote")

	// Define the route for retrieving votes by sub_id
	beego.Router("/votes", &controllers.CatController{}, "get:GetVotes")
}
