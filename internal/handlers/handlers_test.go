package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEchoHandler(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		wantStatus  int
		wantBody    string
		contentType string
	}{
		{
			name:        "valid JSON body",
			body:        `{"test":"data"}`,
			wantStatus:  http.StatusOK,
			wantBody:    `{"test":"data"}`,
			contentType: "application/json",
		},
		{
			name:        "empty body",
			body:        "",
			wantStatus:  http.StatusOK,
			wantBody:    "",
			contentType: "application/json",
		},
		{
			name:        "plain text body",
			body:        "hello world",
			wantStatus:  http.StatusOK,
			wantBody:    "hello world",
			contentType: "application/json",
		},
		{
			name:        "large body",
			body:        strings.Repeat("a", 10000),
			wantStatus:  http.StatusOK,
			wantBody:    strings.Repeat("a", 10000),
			contentType: "application/json",
		},
		{
			name:        "special characters",
			body:        `{"emoji":"🚀","unicode":"日本語"}`,
			wantStatus:  http.StatusOK,
			wantBody:    `{"emoji":"🚀","unicode":"日本語"}`,
			contentType: "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/echo", bytes.NewBufferString(tt.body))
			rec := httptest.NewRecorder()

			EchoHandler(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			if rec.Body.String() != tt.wantBody {
				t.Errorf("body = %q, want %q", rec.Body.String(), tt.wantBody)
			}

			if ct := rec.Header().Get("Content-Type"); ct != tt.contentType {
				t.Errorf("Content-Type = %q, want %q", ct, tt.contentType)
			}
		})
	}
}

// errorReader is a reader that always returns an error
type errorReader struct{}

func (errorReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

func TestEchoHandler_ReadError(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/echo", errorReader{})
	rec := httptest.NewRecorder()

	EchoHandler(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}

	// Verify response is valid JSON with error field
	var errResp errorResponse
	if err := json.NewDecoder(rec.Body).Decode(&errResp); err != nil {
		t.Errorf("failed to decode error response: %v", err)
	}

	if errResp.Error == "" {
		t.Error("expected non-empty error message")
	}

	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want %q", ct, "application/json")
	}
}

func TestHealthHandler(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		wantStatus  int
		wantBody    string
		contentType string
	}{
		{
			name:        "GET request",
			method:      http.MethodGet,
			wantStatus:  http.StatusOK,
			wantBody:    `{"healthy":true}` + "\n",
			contentType: "application/json",
		},
		{
			name:        "HEAD request",
			method:      http.MethodHead,
			wantStatus:  http.StatusOK,
			wantBody:    `{"healthy":true}` + "\n",
			contentType: "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/health", nil)
			rec := httptest.NewRecorder()

			HealthHandler(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			if rec.Body.String() != tt.wantBody {
				t.Errorf("body = %q, want %q", rec.Body.String(), tt.wantBody)
			}

			if ct := rec.Header().Get("Content-Type"); ct != tt.contentType {
				t.Errorf("Content-Type = %q, want %q", ct, tt.contentType)
			}
		})
	}
}
