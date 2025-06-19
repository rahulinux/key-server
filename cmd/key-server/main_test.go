package main

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "reflect"
    "testing"
)

func TestRouteIntegration(t *testing.T) {
    router := NewRouter(16) // use a sample maxSize

    req := httptest.NewRequest("GET", "/key/8", nil)
    rr := httptest.NewRecorder()
    router.ServeHTTP(rr, req)

    if rr.Code != http.StatusOK {
        t.Errorf("expected 200 OK, got %d", rr.Code)
    }
    if !strings.Contains(rr.Body.String(), `"key":"`) {
        t.Errorf("response body %q does not contain expected key", rr.Body.String())
    }
}

func TestParseFlags(t *testing.T) {
    args := []string{"--srv-port=9999", "--max-size=42"}
    cfg, err := ParseFlags(args)
    if err != nil {
      t.Errorf("error parsing cli flags: ParseFlags(%v)", args)
    }
    want := Config{SrvPort: "9999", MaxSize: 42}
    if !reflect.DeepEqual(cfg, want) {
        t.Errorf("ParseFlags(%v) = %+v, want %+v", args, cfg, want)
    }

    // Test defaults
    cfg, _ = ParseFlags([]string{})
    want = Config{SrvPort: "8080", MaxSize: 1024}
    if !reflect.DeepEqual(cfg, want) {
        t.Errorf("ParseFlags(defaults) = %+v, want %+v", cfg, want)
    }
}
