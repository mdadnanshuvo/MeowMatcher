package main

import (
	"testing"

	_ "catApiProject/routers" // Ensure routers are imported

	"github.com/beego/beego/v2/server/web"
)

func TestMainApp(t *testing.T) {
	// Initialize the Beego application
	web.BConfig.RouterCaseSensitive = false // Case-insensitive routes

	// Simulate main function to cover code
	main()
}
