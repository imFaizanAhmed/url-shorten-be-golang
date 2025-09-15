package main

import (
	shorten_controllers "core-api/controllers"

	"github.com/beego/beego/v2/server/web"
)

func main() {
	// Register controllers
	web.Router("/", &MainController{})
	web.Router("/shorten", &shorten_controllers.ShortenController{}, "post:Post")
	web.Router("/:shortCode([A-Za-z0-9_-]+)", &shorten_controllers.ShortenController{}, "get:Get")

	// Start the Beego server
	web.Run()
}

// MainController handles the root endpoint
type MainController struct {
	web.Controller
}

// Get handles the GET request for the root URL
func (c *MainController) Get() {
	response := map[string]string{
		"message":   "URL Shortener API",
		"version":   "1.0",
		"endpoints": "POST /shorten - Create short URL, GET /{shortCode} - Redirect to long URL",
	}
	c.Data["json"] = response
	c.ServeJSON()
}
