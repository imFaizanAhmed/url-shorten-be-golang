package storage

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
	"sync"
)

// In-memory storage for URL mappings
var (
	urlStore = make(map[string]string) // shortCode -> longURL
	mutex    = sync.RWMutex{}
)

// StoreURL stores a mapping between shortCode and longURL
func StoreURL(shortCode, longURL string) {
	mutex.Lock()
	defer mutex.Unlock()
	urlStore[shortCode] = longURL
}

// GetURL retrieves the longURL for a given shortCode
func GetURL(shortCode string) (string, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	longURL, exists := urlStore[shortCode]
	return longURL, exists
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
