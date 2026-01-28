package server

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lkendrickd/echo-server/internal/config"
)

func TestNewServer(t *testing.T) {
	tests := []struct {
		name             string
		port             string
		authEnabled      bool
		wantReadTimeout  time.Duration
		wantWriteTimeout time.Duration
		wantIdleTimeout  time.Duration
	}{
		{
			name:             "default port 8080 without auth",
			port:             ":8080",
			authEnabled:      false,
			wantReadTimeout:  15 * time.Second,
			wantWriteTimeout: 15 * time.Second,
			wantIdleTimeout:  60 * time.Second,
		},
		{
			name:             "custom port 9090 with auth",
			port:             ":9090",
			authEnabled:      true,
			wantReadTimeout:  15 * time.Second,
			wantWriteTimeout: 15 * time.Second,
			wantIdleTimeout:  60 * time.Second,
		},
		{
			name:             "port with host",
			port:             "localhost:8081",
			authEnabled:      false,
			wantReadTimeout:  15 * time.Second,
			wantWriteTimeout: 15 * time.Second,
			wantIdleTimeout:  60 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
			mux := http.NewServeMux()
			cfg := &config.Config{
				Port:        tt.port,
				LogLevel:    "info",
				AuthEnabled: tt.authEnabled,
			}

			s := NewServer(logger, mux, tt.port, cfg)

			if s == nil {
				t.Fatal("NewServer returned nil")
			}

			if s.logger != logger {
				t.Error("logger not set correctly")
			}

			if s.muxer != mux {
				t.Error("muxer not set correctly")
			}

			if s.port != tt.port {
				t.Errorf("port = %q, want %q", s.port, tt.port)
			}

			if s.config != cfg {
				t.Error("config not set correctly")
			}

			if s.server == nil {
				t.Fatal("http.Server is nil")
			}

			if s.server.Addr != tt.port {
				t.Errorf("server.Addr = %q, want %q", s.server.Addr, tt.port)
			}

			if s.server.ReadTimeout != tt.wantReadTimeout {
				t.Errorf("ReadTimeout = %v, want %v", s.server.ReadTimeout, tt.wantReadTimeout)
			}

			if s.server.WriteTimeout != tt.wantWriteTimeout {
				t.Errorf("WriteTimeout = %v, want %v", s.server.WriteTimeout, tt.wantWriteTimeout)
			}

			if s.server.IdleTimeout != tt.wantIdleTimeout {
				t.Errorf("IdleTimeout = %v, want %v", s.server.IdleTimeout, tt.wantIdleTimeout)
			}
		})
	}
}

func TestNewServer_NilConfig(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	mux := http.NewServeMux()

	// Should not panic with nil config
	s := NewServer(logger, mux, ":8080", nil)

	if s == nil {
		t.Fatal("NewServer returned nil")
	}
}

func TestSetupRoutes(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	mux := http.NewServeMux()
	cfg := &config.Config{
		Port:        ":8080",
		LogLevel:    "info",
		AuthEnabled: false,
	}

	s := NewServer(logger, mux, ":8080", cfg)
	s.SetupRoutes()

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		{
			name:       "health endpoint",
			method:     http.MethodGet,
			path:       "/health",
			wantStatus: http.StatusOK,
		},
		{
			name:       "metrics endpoint",
			method:     http.MethodGet,
			path:       "/metrics",
			wantStatus: http.StatusOK,
		},
		{
			name:       "echo endpoint with body",
			method:     http.MethodPost,
			path:       "/api/v1/echo",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestSetupRoutes_MethodNotAllowed(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	mux := http.NewServeMux()
	cfg := &config.Config{
		Port:        ":8080",
		LogLevel:    "info",
		AuthEnabled: false,
	}

	s := NewServer(logger, mux, ":8080", cfg)
	s.SetupRoutes()

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{
			name:   "GET on echo endpoint",
			method: http.MethodGet,
			path:   "/api/v1/echo",
		},
		{
			name:   "POST on health endpoint",
			method: http.MethodPost,
			path:   "/health",
		},
		{
			name:   "DELETE on metrics endpoint",
			method: http.MethodDelete,
			path:   "/metrics",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			// These should return 405 Method Not Allowed or 404 depending on Go version
			if rec.Code != http.StatusMethodNotAllowed && rec.Code != http.StatusNotFound {
				t.Errorf("status = %d, want %d or %d", rec.Code, http.StatusMethodNotAllowed, http.StatusNotFound)
			}
		})
	}
}

func TestSetupRoutes_NotFound(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	mux := http.NewServeMux()
	cfg := &config.Config{
		Port:        ":8080",
		LogLevel:    "info",
		AuthEnabled: false,
	}

	s := NewServer(logger, mux, ":8080", cfg)
	s.SetupRoutes()

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{
			name:   "unknown path",
			method: http.MethodGet,
			path:   "/unknown",
		},
		{
			name:   "wrong API version",
			method: http.MethodPost,
			path:   "/api/v2/echo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != http.StatusNotFound {
				t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
			}
		})
	}
}
