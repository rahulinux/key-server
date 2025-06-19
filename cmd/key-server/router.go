package main

import (
    "github.com/gorilla/mux"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/rahulinux/key-server/internal/api"
    "github.com/rahulinux/key-server/internal/config"
    "log/slog"
)

func NewRouter(cfg config.Config, logger *slog.Logger) *mux.Router {
    r := mux.NewRouter()
    
    // Create handler with dependencies
    keyHandler := api.NewKeyHandler(cfg.MaxSize, logger)
    healthHandler := api.NewHealthHandler(logger)
    
    // API routes
    r.HandleFunc("/key/{length:[0-9]+}", keyHandler.HandleKey).Methods("GET")
    r.HandleFunc("/healthz", healthHandler.HandleHealth).Methods("GET")
    r.Handle("/metrics", promhttp.Handler()).Methods("GET")
    
    return r
}
