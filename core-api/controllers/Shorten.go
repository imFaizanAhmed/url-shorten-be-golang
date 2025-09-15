package shorten_controllers

import (
	"encoding/json"
	"fmt"
	"net/url"

	"core-api/storage"
	"core-api/structs"

	"github.com/beego/beego/v2/server/web"
)

// ShortenController handles URL shortening requests
type ShortenController struct {
	web.Controller
}

// Post handles POST requests to create short URLs
func (c *ShortenController) Post() {
	var req structs.ShortenRequest

	// Read request body manually
	body := c.Ctx.Input.RequestBody
	fmt.Println("Request body length:", len(body))
	fmt.Println("Request body:", string(body))

	// If RequestBody is empty, try reading directly from the request
	if len(body) == 0 {
		bodyBytes := make([]byte, c.Ctx.Request.ContentLength)
		c.Ctx.Request.Body.Read(bodyBytes)
		body = bodyBytes
		fmt.Println("Read from Request.Body:", string(body))
	}

	// Parse JSON request body
	if err := json.Unmarshal(body, &req); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.Data["json"] = structs.ErrorResponse{Error: fmt.Sprintf("Invalid JSON payload: %v", err)}
		c.ServeJSON()
		return
	}

	// Validate URL
	if req.LongURL == "" {
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.Data["json"] = structs.ErrorResponse{Error: "long_url is required"}
		c.ServeJSON()
		return
	}

	// Validate URL format
	if _, err := url.ParseRequestURI(req.LongURL); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.Data["json"] = structs.ErrorResponse{Error: "Invalid URL format"}
		c.ServeJSON()
		return
	}

	// Generate short code
	shortCode := storage.GenerateShortCode()

	// Store the mapping
	storage.StoreURL(shortCode, req.LongURL)

	// Get base URL for the response
	baseURL := getBaseURL(c)
	shortURL := fmt.Sprintf("%s/%s", baseURL, shortCode)

	// Return response
	response := structs.ShortenResponse{
		ShortURL:  shortURL,
		LongURL:   req.LongURL,
		ShortCode: shortCode,
	}

	c.Data["json"] = response
	c.ServeJSON()
}

// getBaseURL extracts the base URL from the request
func getBaseURL(c *ShortenController) string {
	scheme := "http"
	if c.Ctx.Input.IsSecure() {
		scheme = "https"
	}

	host := c.Ctx.Request.Host
	return fmt.Sprintf("%s://%s", scheme, host)
}

// Get handles GET requests for short URL redirection
func (c *ShortenController) Get() {
	shortCode := c.Ctx.Input.Param(":shortCode")

	// Look up the long URL
	longURL, exists := storage.GetURL(shortCode)

	if !exists {
		c.Ctx.ResponseWriter.WriteHeader(404)
		c.Data["json"] = structs.ErrorResponse{Error: "Short URL not found"}
		c.ServeJSON()
		return
	}

	// Redirect to the long URL
	c.Redirect(longURL, 301)
}
