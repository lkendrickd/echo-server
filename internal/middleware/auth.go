package middleware

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"strings"
)

// APIKeyValidator interface for validating API keys
type APIKeyValidator interface {
	ValidateAPIKey(key string) bool
}

// authErrorResponse represents an authentication error response
type authErrorResponse struct {
	Error string `json:"error"`
}

// AuthMiddleware creates a middleware that validates API keys
// Protected paths require a valid API key in the X-API-Key header
func AuthMiddleware(validator APIKeyValidator, protectedPrefixes []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if this path needs protection
			if !isProtectedPath(r.URL.Path, protectedPrefixes) {
				next.ServeHTTP(w, r)
				return
			}

			// Get API key from header
			apiKey := r.Header.Get("X-API-Key")
			if apiKey == "" {
				writeAuthError(w, http.StatusUnauthorized, "missing API key")
				return
			}

			// Validate the API key
			if !validator.ValidateAPIKey(apiKey) {
				writeAuthError(w, http.StatusUnauthorized, "invalid API key")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// isProtectedPath checks if the given path matches any protected prefix
func isProtectedPath(path string, protectedPrefixes []string) bool {
	for _, prefix := range protectedPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

// writeAuthError writes a JSON error response for authentication failures
func writeAuthError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(authErrorResponse{Error: message})
}

// SecureCompare performs a constant-time comparison of two strings
// This prevents timing attacks when comparing API keys
func SecureCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
