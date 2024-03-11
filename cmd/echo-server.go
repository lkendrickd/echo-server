package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/lkendrickd/echo-server/internal/server"
	"log/slog"
)

func main() {
	// Define command-line flags
	port := flag.String("port", "8080", "The port to listen on")
	logLevel := flag.String("logLevel", "info", "The log level")

	// Parse command-line flags
	flag.Parse()

	// Check for environment variables and override flag values if necessary for port
	if envPort, exists := os.LookupEnv("PORT"); exists {
		*port = envPort
	}

	// Check for environment variables and override flag values if necessary for log level
	if envLogLevel, exists := os.LookupEnv("LOG_LEVEL"); exists {
		*logLevel = envLogLevel
	}

	// Set the log level based on the provided logLevel string
	slogLevel := setLogLevel(*logLevel)

	// Initialize the logger with the determined log level
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slogLevel}))

	// Initialize the HTTP server mux
	mux := http.NewServeMux()

	// Create and start the server
	s := server.NewServer(logger, mux, fmt.Sprintf(":%s", *port))
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
