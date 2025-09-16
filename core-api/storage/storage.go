package storage

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"strings"
	"time"

	redisClient "core-api/redis"
)

const (
	// URL_PREFIX is the prefix for URL keys in Redis
	URL_PREFIX = "url:"
	// Default TTL for URLs (24 hours)
	DEFAULT_TTL = 24 * time.Hour
)

// StoreURL stores a mapping between shortCode and longURL in Redis
func StoreURL(shortCode, longURL string) {
	key := URL_PREFIX + shortCode
	ctx := redisClient.GetContext()

	err := redisClient.Client.Set(ctx, key, longURL, DEFAULT_TTL).Err()
	if err != nil {
		log.Printf("Error storing URL in Redis: %v", err)
	}
}

// GetURL retrieves the longURL for a given shortCode from Redis
func GetURL(shortCode string) (string, bool) {
	key := URL_PREFIX + shortCode
	ctx := redisClient.GetContext()

	longURL, err := redisClient.Client.Get(ctx, key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			// Key doesn't exist
			return "", false
		}
		log.Printf("Error retrieving URL from Redis: %v", err)
		return "", false
	}

	return longURL, true
}

// GenerateShortCode generates a random short code for the URL
func GenerateShortCode() string {
	for {
		// Generate 8 random bytes for better entropy
		bytes := make([]byte, 8)
		rand.Read(bytes)

		// Encode to base64 and clean up
		encoded := base64.URLEncoding.EncodeToString(bytes)
		// Remove all special characters (dashes, underscores, padding)
		encoded = strings.ReplaceAll(encoded, "-", "")
		encoded = strings.ReplaceAll(encoded, "_", "")
		shortCode := strings.TrimRight(encoded, "=")

		// Filter to keep only alphanumeric characters
		var result strings.Builder
		for _, char := range shortCode {
			if (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') {
				result.WriteRune(char)
			}
		}

		finalCode := result.String()

		// Ensure we have at least 8 characters
		if len(finalCode) >= 8 {
			return finalCode[:8]
		}

		// If we don't have enough characters, generate again
		continue
	}
}
