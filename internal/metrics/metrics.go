package metrics

import (
    "fmt"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    KeyLengthHistogram prometheus.Histogram
    HTTPStatusCounter  *prometheus.CounterVec
    RequestDuration    *prometheus.HistogramVec
    ActiveConnections  prometheus.Gauge
)

func InitMetrics(maxSize int) error {
    if maxSize <= 0 {
        return fmt.Errorf("maxSize must be positive")
    }

    buckets, err := generateHistogramBuckets(maxSize)
    if err != nil {
        return err
    }

    // Key length distribution
    KeyLengthHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
        Name:    "key_length_distribution",
        Help:    "Distribution of requested key lengths",
        Buckets: buckets,
    })

    // HTTP status codes
    HTTPStatusCounter = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_status_codes_total",
            Help: "Count of HTTP status codes returned by the server",
        },
        []string{"code"},
    )

    // Request duration
    RequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "request_duration_seconds",
            Help:    "Time spent processing requests",
            Buckets: prometheus.DefBuckets,
        },
        []string{"endpoint"},
    )

    // Active connections
    ActiveConnections = promauto.NewGauge(prometheus.GaugeOpts{
        Name: "active_connections",
        Help: "Number of active connections",
    })

    return nil
}

// generateHistogramBuckets creates up to 20 evenly spaced histogram buckets
// based on maxSize and ensures the resulting slice is strictly increasing.
// If maxSize is smaller than the last computed bucket, it is not appended to avoid Prometheus panics.
func generateHistogramBuckets(maxSize int) ([]float64, error) {
    if maxSize <= 0 {
        return nil, fmt.Errorf("maxSize must be positive")
    }

    step := (maxSize + 19) / 20
    if step <= 0 {
        step = 1
    }

    buckets := make([]float64, 0, 21)
    for i := 0; i < 20; i++ {
        val := float64(i * step)
        buckets = append(buckets, val)
    }

    last := buckets[len(buckets)-1]
    if float64(maxSize) > last {
        buckets = append(buckets, float64(maxSize))
    }

    return buckets, nil
}

