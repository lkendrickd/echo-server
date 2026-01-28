package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/lkendrickd/echo-server/internal/config"
	"github.com/lkendrickd/echo-server/internal/server"
)

func main() {
	// Load configuration from environment variables
	cfg := config.New()

	// Set the log level based on the config
	slogLevel := setLogLevel(cfg.LogLevel)

	// Initialize the logger with the determined log level
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slogLevel}))

	// Log configuration (without sensitive data)
	logger.Info("configuration loaded",
		"port", cfg.Port,
		"log_level", cfg.LogLevel,
		"auth_enabled", cfg.AuthEnabled,
		"api_key_count", cfg.APIKeyCount(),
	)

	// Warn if auth is enabled but no keys configured
	if cfg.AuthEnabled && !cfg.HasAPIKeys() {
		logger.Warn("authentication enabled but no API keys configured")
	}

	// Initialize the HTTP server mux
	mux := http.NewServeMux()

	// Create and start the server
	s := server.NewServer(logger, mux, fmt.Sprintf(":%s", cfg.Port), cfg)
	if err := s.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// setLogLevel sets the log level based on the provided string
func setLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		fmt.Printf("Unknown log level: %s, defaulting to info\n", level)
		return slog.LevelInfo
	}
}
