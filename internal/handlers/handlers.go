package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

// errorResponse represents a JSON error response
type errorResponse struct {
	Error string `json:"error"`
}

// EchoHandler is the echo handler that returns the request body
func EchoHandler(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
		return
	}

	// Write the request body back to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

// HealthHandler is the health check handler
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"healthy":true}` + "\n"))
}
