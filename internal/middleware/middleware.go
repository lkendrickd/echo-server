package middleware

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method", "status"},
	)

	EndpointCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"path", "method", "status"},
	)
)

// responseWriter wraps http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode  int
	wroteHeader bool
}

// newResponseWriter creates a new responseWriter with default status 200
func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		wroteHeader:    false,
	}
}

// WriteHeader captures the status code on first call and delegates to underlying WriteHeader
func (rw *responseWriter) WriteHeader(code int) {
	if !rw.wroteHeader {
		rw.statusCode = code
		rw.wroteHeader = true
	}
	rw.ResponseWriter.WriteHeader(code)
}

// MetricsMiddleware is the middleware for capturing metrics
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := r.URL.Path
		method := r.Method

		// Wrap the response writer to capture status code
		wrapped := newResponseWriter(w)

		// Start timer for duration metric
		timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
			status := strconv.Itoa(wrapped.statusCode)
			RequestDuration.WithLabelValues(route, method, status).Observe(v)
		}))
		defer timer.ObserveDuration()

		next.ServeHTTP(wrapped, r)

		// Increment the endpoint counter with status code
		status := strconv.Itoa(wrapped.statusCode)
		EndpointCount.WithLabelValues(route, method, status).Inc()
	})
}
