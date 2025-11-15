package main

import (
	shorten_controllers "core-api/controllers"
	redisClient "core-api/redis"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

// corsFilter handles CORS (Cross-Origin Resource Sharing) headers
func corsFilter(ctx *context.Context) {
	// Set CORS headers
	ctx.Output.Header("Access-Control-Allow-Origin", "*")
	ctx.Output.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	ctx.Output.Header("Access-Control-Allow-Credentials", "true")
	ctx.Output.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With")

	// Handle preflight OPTIONS request
	if ctx.Request.Method == "OPTIONS" {
		ctx.Output.SetStatus(200)
		return
	}
}

func main() {
	// Initialize Redis connection
	redisClient.Init()

	// Add CORS filter
	web.InsertFilter("*", web.BeforeRouter, corsFilter)

	// Register controllers
	web.Router("/", &MainController{})
	web.Router("/shorten", &shorten_controllers.ShortenController{}, "post:Post")
	web.Router("/:shortCode([A-Za-z0-9_-]+)", &shorten_controllers.ShortenController{}, "get:Get")

	// Start the Beego server
	web.Run("0.0.0.0:8080")
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
