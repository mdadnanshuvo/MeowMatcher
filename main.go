package main

import (
	_ "catApiProject/routers"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	// Set route case sensitivity to false (case-insensitive routes)
	beego.BConfig.RouterCaseSensitive = false

	// Start the Beego server
	beego.Run()
}
