package api_test

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/gorilla/mux"
    "github.com/rahulinux/key-server/internal/api"
    "github.com/rahulinux/key-server/internal/metrics"
    "log/slog"
    "os"
)

func TestKeyHandler(t *testing.T) {
    // ✅ Initialize metrics to avoid nil dereference
    if err := metrics.InitMetrics(16); err != nil {
        t.Fatalf("Failed to initialize metrics: %v", err)
    }

    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
    handler := api.NewKeyHandler(16, logger)

    router := mux.NewRouter()
    router.HandleFunc("/key/{length:[0-9]+}", handler.HandleKey).Methods("GET")

    tests := []struct {
        name         string
        url          string
        wantCode     int
        wantContains string
    }{
        {"valid", "/key/8", http.StatusOK, `"key":"`},
        {"too large", "/key/32", http.StatusBadRequest, "exceeds maximum allowed size"},
        {"zero", "/key/0", http.StatusBadRequest, "Length must be positive"},
        {"not a number", "/key/abc", http.StatusNotFound, "404 page not found"}, // gorilla mux won't match non-numeric
        {"negative", "/key/-5", http.StatusNotFound, "404 page not found"},
        {"missing param", "/key/", http.StatusNotFound, "404 page not found"},
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            req := httptest.NewRequest(http.MethodGet, tc.url, nil)
            rr := httptest.NewRecorder()
            router.ServeHTTP(rr, req)

            if rr.Code != tc.wantCode {
                t.Errorf("got status %d, want %d", rr.Code, tc.wantCode)
            }

            if !strings.Contains(rr.Body.String(), tc.wantContains) {
                t.Errorf("body %q does not contain %q", rr.Body.String(), tc.wantContains)
            }
        })
    }
}

