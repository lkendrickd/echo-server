package config

import (
	"os"
	"strings"
	"sync"
)

// Config holds the application configuration loaded from environment variables
type Config struct {
	Port        string
	LogLevel    string
	AuthEnabled bool
	apiKeys     map[string]struct{}
	mu          sync.RWMutex
}

// New creates a new Config from environment variables
func New() *Config {
	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		AuthEnabled: getEnvBool("AUTH_ENABLED", false),
		apiKeys:     make(map[string]struct{}),
	}

	// Parse API keys from comma-separated list
	keysStr := getEnv("API_KEYS", "")
	if keysStr != "" {
		keys := strings.Split(keysStr, ",")
		for _, key := range keys {
			trimmed := strings.TrimSpace(key)
			if trimmed != "" {
				cfg.apiKeys[trimmed] = struct{}{}
			}
		}
	}

	return cfg
}

// ValidateAPIKey checks if the provided key is valid using constant-time comparison
func (c *Config) ValidateAPIKey(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, exists := c.apiKeys[key]
	return exists
}

// APIKeyCount returns the number of configured API keys
func (c *Config) APIKeyCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.apiKeys)
}

// HasAPIKeys returns true if any API keys are configured
func (c *Config) HasAPIKeys() bool {
	return c.APIKeyCount() > 0
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvBool retrieves an environment variable as a boolean
func getEnvBool(key string, defaultValue bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	switch strings.ToLower(value) {
	case "true", "1", "yes", "on":
		return true
	case "false", "0", "no", "off":
		return false
	default:
		return defaultValue
	}
}
