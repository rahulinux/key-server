package api

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/gorilla/mux"
)

func TestKeyHandler(t *testing.T) {
    maxSize := 16
    router := mux.NewRouter()
    router.HandleFunc("/key/{length}", KeyHandler(maxSize)).Methods("GET")

    tests := []struct {
        name         string
        url          string
        wantCode     int
        wantContains string
    }{
        {"valid", "/key/8", http.StatusOK, `"key":"`},
        {"too large", "/key/32", http.StatusBadRequest, "exceeds max-size"},
        {"not a number", "/key/abc", http.StatusBadRequest, "Invalid length"},
        {"zero", "/key/0", http.StatusBadRequest, "Invalid length"},
        {"negative", "/key/-5", http.StatusBadRequest, "Invalid length"},
        {"not found", "/key/", http.StatusNotFound, "404 page not found"},
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            req := httptest.NewRequest("GET", tc.url, nil)
            rr := httptest.NewRecorder()
            router.ServeHTTP(rr, req)
            if rr.Code != tc.wantCode {
                t.Errorf("got status %d, want %d", rr.Code, tc.wantCode)
            }
            if tc.wantContains != "" && !strings.Contains(rr.Body.String(), tc.wantContains) {
                t.Errorf("body %q does not contain %q", rr.Body.String(), tc.wantContains)
            }
        })
    }
}
