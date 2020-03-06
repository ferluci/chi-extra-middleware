package extmiddleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Config responsible to configure middleware
type Config struct {
	Namespace string
	Buckets   []float64
	Subsystem string
}

var DefaultBuckets = []float64{
	0.00001, // 10 microseconds
	0.0001,
	0.001, // 1ms
	0.002,
	0.005,
	0.01, // 10ms
	0.02,
	0.05,
	0.1, // 100 ms
	0.2,
	0.5,
	1.0, // 1s
	2.0,
	5.0,
}

const (
	httpRequestsCount    = "requests_total"
	httpRequestsDuration = "request_duration_seconds"
)

// DefaultConfig has the default instrumentation config
var DefaultConfig = Config{
	Namespace: "chi",
	Subsystem: "http",
	Buckets:   DefaultBuckets,
}

// NewConfig returns a new config with default values
func NewConfig() Config {
	return DefaultConfig
}

// Metrics returns an echo middleware with default config for instrumentation.
func Metrics() func(next http.Handler) http.Handler {
	return MetricsWithConfig(DefaultConfig)
}

func MetricsWithConfig(config Config) func(next http.Handler) http.Handler {
	httpRequests := promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: config.Namespace,
		Subsystem: config.Subsystem,
		Name:      httpRequestsCount,
		Help:      "Number of HTTP operations",
	}, []string{"status", "method", "pattern"})

	httpDuration := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: config.Namespace,
		Subsystem: config.Subsystem,
		Name:      httpRequestsDuration,
		Help:      "Spend time by processing a route",
		Buckets:   config.Buckets,
	}, []string{"status", "method", "pattern"})

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			routePattern := chi.RouteContext(r.Context()).RoutePattern()

			observer := httpDuration.WithLabelValues(strconv.Itoa(ww.Status()), r.Method, routePattern)
			observer.Observe(time.Since(start).Seconds())

			httpRequests.WithLabelValues(strconv.Itoa(ww.Status()), r.Method, routePattern).Inc()
		}
		return http.HandlerFunc(fn)
	}
}
