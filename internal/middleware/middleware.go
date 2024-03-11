package middleware

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
)

var (
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)

	EndpointCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"path", "method"},
	)
)

// MetricsMiddleware is the middleware for capturing metrics
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := r.URL.Path
		method := r.Method

		// Increment the endpoint counter
		EndpointCount.WithLabelValues(route, method).Inc()

		// Start timer for duration metric
		timer := prometheus.NewTimer(RequestDuration.WithLabelValues(route, method))
		defer timer.ObserveDuration()

		next.ServeHTTP(w, r)
	})
}
