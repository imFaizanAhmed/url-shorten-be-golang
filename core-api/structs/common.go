package structs

// ErrorResponse represents error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// ShortenResponse represents the response for URL shortening
type ShortenResponse struct {
	ShortURL  string `json:"short_url"`
	LongURL   string `json:"long_url"`
	ShortCode string `json:"short_code"`
}

// URLMapping represents the mapping between long and short URLs
type URLMapping struct {
	LongURL  string `json:"long_url"`
	ShortURL string `json:"short_url"`
}

// ShortenRequest represents the request payload for URL shortening
type ShortenRequest struct {
	LongURL string `json:"long_url"`
}
