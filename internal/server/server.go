package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lkendrickd/echo-server/internal/handlers"
	"github.com/lkendrickd/echo-server/internal/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server is the HTTP server
type Server struct {
	logger *slog.Logger
	muxer  *http.ServeMux
	server *http.Server
	port   string
}

// registerMetric safely registers a prometheus collector, ignoring AlreadyRegisteredError
func registerMetric(c prometheus.Collector) {
	if err := prometheus.DefaultRegisterer.Register(c); err != nil {
		var alreadyRegistered prometheus.AlreadyRegisteredError
		if !errors.As(err, &alreadyRegistered) {
			panic(err)
		}
	}
}

// NewServer creates a new Server with middleware applied
func NewServer(l *slog.Logger, mux *http.ServeMux, port string) *Server {
	// Register prometheus metrics, tolerating duplicate registrations
	registerMetric(middleware.RequestDuration)
	registerMetric(middleware.EndpointCount)

	// Wrap the existing muxer with the metricsMiddleware
	wrappedMux := middleware.MetricsMiddleware(mux)

	// Create a new http.Server using the wrapped muxer
	server := &http.Server{
		Addr:         port,
		Handler:      wrappedMux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		logger: l,
		muxer:  mux,
		server: server,
		port:   port,
	}
}

// Start starts the server and gracefully handles shutdown
func (s *Server) Start() error {
	// Setting up signal capturing
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// Add routes to the muxer
	s.logger.Debug("setting up routes")
	s.SetupRoutes()

	// Starting server in a goroutine
	go func() {
		s.logger.Info("starting server")
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("server failed to start", "error", err)
		}
	}()

	// Block until a signal is received
	<-stopChan
	s.logger.Info("shutting down server")

	// Create a deadline to wait for this is the duration the server will wait for existing connections to finish
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("server shutdown failed", "error", err)
		return err
	}
	s.logger.Info("server exited properly")
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
