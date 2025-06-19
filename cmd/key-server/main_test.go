package main

import (
    "net/http"
    "testing"
    "net/http/httptest"

    "github.com/rahulinux/key-server/internal/config"
)

func TestSetupLogger(t *testing.T) {
    logger := setupLogger("debug")
    if logger == nil {
        t.Error("Expected logger, got nil")
    }
}

func TestCreateHandler(t *testing.T) {
    cfg := config.Config{MaxSize: 64}
    logger := setupLogger("info")
    h := createHandler(cfg, logger)
    if h == nil {
        t.Error("Expected handler, got nil")
    }

    // Optionally test if it serves metrics
    req, _ := http.NewRequest("GET", "/metrics", nil)
    rr := httptest.NewRecorder()
    h.ServeHTTP(rr, req)

    if rr.Code != http.StatusOK {
        t.Errorf("Expected 200 from /metrics, got %d", rr.Code)
    }
}

