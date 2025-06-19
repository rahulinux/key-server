package metrics_test

import (
    "testing"

    "github.com/rahulinux/key-server/internal/metrics"
)

func TestInitMetrics_Valid(t *testing.T) {
    err := metrics.InitMetrics(1024)
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
}

func TestInitMetrics_Invalid(t *testing.T) {
    err := metrics.InitMetrics(0)
    if err == nil {
        t.Error("Expected error for invalid maxSize")
    }
}

