package config

import (
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	Port              string
	Environment       string
	UserServiceURL    string
	ProductServiceURL string
	APIKey            string
	CookieDomain      string
	CookieMaxAge      int
}

// IsDev retorna true cuando el entorno es development.
func (c *Config) IsDev() bool {
	return c.Environment == "development"
}

// Load loads configuration from environment variables with defaults
func Load() *Config {
	cookieMaxAge := 3600
	if v := os.Getenv("COOKIE_MAX_AGE"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cookieMaxAge = n
		}
	}

	return &Config{
		Port:              getEnv("PORT", "8080"),
		Environment:       getEnv("ENVIRONMENT", "development"),
		UserServiceURL:    getEnv("USER_SERVICE_URL", "https://jsonplaceholder.typicode.com"),
		ProductServiceURL: getEnv("PRODUCT_SERVICE_URL", "https://fakestoreapi.com"),
		APIKey:            getEnv("API_KEY", ""),
		CookieDomain:      getEnv("COOKIE_DOMAIN", ""),
		CookieMaxAge:      cookieMaxAge,
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
