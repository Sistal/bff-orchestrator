package config

import (
	"os"
)

// Config holds the application configuration
type Config struct {
	Port              string
	Environment       string
	UserServiceURL    string
	ProductServiceURL string
}

// Load loads configuration from environment variables with defaults
func Load() *Config {
	return &Config{
		Port:              getEnv("PORT", "8080"),
		Environment:       getEnv("ENVIRONMENT", "development"),
		UserServiceURL:    getEnv("USER_SERVICE_URL", "https://jsonplaceholder.typicode.com"),
		ProductServiceURL: getEnv("PRODUCT_SERVICE_URL", "https://fakestoreapi.com"),
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
