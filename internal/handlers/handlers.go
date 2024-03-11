package handlers

import (
	"fmt"
	"io"
	"net/http"
)

// EchoHandler is the echo handler that returns the request body
func EchoHandler(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			return
		}
		return
	}

	// Write the request body back to the client
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(body); err != nil {
		return
	}

	// Close the request body
	if err := r.Body.Close(); err != nil {
		return
	}
	return
}

// HealthHandler is the health check handler
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(fmt.Sprintf(`{"healthy":true}` + "\n"))); err != nil {
		return
	}
}
