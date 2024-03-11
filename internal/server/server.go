package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/lkendrickd/echo-server/internal/handlers"
	"github.com/lkendrickd/echo-server/internal/middleware"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	// Register the prometheus metrics
	prometheus.MustRegister(middleware.RequestDuration)
	prometheus.MustRegister(middleware.EndpointCount)
}

// Server is the HTTP server
type Server struct {
	logger *slog.Logger
	muxer  *http.ServeMux
	server *http.Server
	port   string
}

// NewServer creates a new Server with middleware applied
func NewServer(l *slog.Logger, mux *http.ServeMux, port string) *Server {
	// Wrap the existing muxer with the metricsMiddleware
	wrappedMux := middleware.MetricsMiddleware(mux)

	// Create a new http.Server using the wrapped muxer
	server := &http.Server{
		Addr:    port,
		Handler: wrappedMux,
	}

	return &Server{
		logger: l,
		muxer:  mux,
		server: server,
	}
}

// Start starts the server and gracefully handles shutdown
func (s *Server) Start() error {
	// Setting up signal capturing
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// Add routes to the muxer
	s.logger.Debug(`{"message": "setting up routes"}`)
	s.SetupRoutes()

	// Starting server in a goroutine
	go func() {
		s.logger.Log(context.Background(), slog.LevelInfo, `{"message": "starting server"}`)
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error(`{"message": "server failed to start", "error": %s}`, err.Error())
		}
	}()

	// Block until a signal is received
	<-stopChan
	s.logger.Log(context.Background(), slog.LevelInfo, `{"message": "shutting down server"}`)

	// Create a deadline to wait for this is the duration the server will wait for existing connections to finish
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error(`{"message": "server shutdown failed", "error": ` + err.Error() + `}`)
		return err
	}
	s.logger.Log(context.Background(), slog.LevelInfo, `{"message": "server exited properly"}`)
	return nil
}

// SetupRoutes sets up the server routes
func (s *Server) SetupRoutes() {
	path := "/api/v1"
	s.muxer.HandleFunc(
		fmt.Sprintf(
			"%s %s/echo",
			http.MethodPost,
			path,
		), handlers.EchoHandler,
	)
	s.muxer.HandleFunc("GET /health", handlers.HealthHandler)
	s.muxer.Handle("GET /metrics", promhttp.Handler())
}
