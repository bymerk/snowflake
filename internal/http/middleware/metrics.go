package middleware

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests by path",
		},
		[]string{"path"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of request duration",
			Buckets: []float64{0.0005, 0.001, 0.002, 0.005, 0.01, 0.025, 0.05}, // от 0.5 ms до 50 ms
		},
		[]string{"path"},
	)
)

func init() {
	prometheus.MustRegister(httpRequests)
	prometheus.MustRegister(requestDuration)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		httpRequests.WithLabelValues(r.URL.Path).Inc()
		requestDuration.WithLabelValues(r.URL.Path).Observe(time.Since(start).Seconds())
	})
}
