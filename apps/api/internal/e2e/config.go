package e2e

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

// Config holds the test configuration
type Config struct {
	APIURL string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Try to load .env from project root
	// Get current file path to find project root relative to this file
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	// assuming e2e is in apps/api/internal/e2e, root is ../.. (apps/api)
	projectRoot := filepath.Join(basepath, "../..")
	
	if err := godotenv.Load(filepath.Join(projectRoot, ".env")); err != nil {
		// Log but don't fail, as env vars might be set directly in CI/CD
		log.Printf("Info: .env file not found or failed to load: %v", err)
	}

	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		apiURL = "http://localhost:" + port
	}

	return &Config{
		APIURL: apiURL,
	}
}
