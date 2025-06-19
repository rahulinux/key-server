package main

import (
    "context"
    "fmt"
    "log/slog"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gorilla/handlers"
    "github.com/rahulinux/key-server/internal/config"
    "github.com/rahulinux/key-server/internal/metrics"
)

func main() {
    // Parse configuration
    cfg, err := config.ParseFlags(os.Args[1:])
    if err != nil {
        fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
        os.Exit(1)
    }

    // Setup structured logging
    logger := setupLogger(cfg.LogLevel)
    
    // Initialize metrics
    if err := metrics.InitMetrics(cfg.MaxSize); err != nil {
        logger.Error("Failed to initialize metrics", "error", err)
        os.Exit(1)
    }

    // Create and start server
    server := &http.Server{
        Addr:    fmt.Sprintf(":%s", cfg.SrvPort),
        Handler: createHandler(cfg, logger),
    }

    // Graceful shutdown
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    go func() {
        logger.Info("Starting server", "addr", server.Addr)
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Error("Server failed", "error", err)
            cancel()
        }
    }()

    // Wait for interrupt signal
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    select {
    case sig := <-sigChan:
        logger.Info("Received signal, shutting down", "signal", sig)
    case <-ctx.Done():
        logger.Info("Context cancelled, shutting down")
    }

    // Graceful shutdown with timeout
    shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer shutdownCancel()

    if err := server.Shutdown(shutdownCtx); err != nil {
        logger.Error("Server shutdown failed", "error", err)
        os.Exit(1)
    }

    logger.Info("Server stopped gracefully")
}

func setupLogger(level string) *slog.Logger {
    var logLevel slog.Level
    switch level {
    case "debug":
        logLevel = slog.LevelDebug
    case "info":
        logLevel = slog.LevelInfo
    case "warn":
        logLevel = slog.LevelWarn
    case "error":
        logLevel = slog.LevelError
    default:
        logLevel = slog.LevelInfo
    }

    opts := &slog.HandlerOptions{
        Level: logLevel,
    }

    return slog.New(slog.NewJSONHandler(os.Stdout, opts))
}

func createHandler(cfg config.Config, logger *slog.Logger) http.Handler {
    router := NewRouter(cfg, logger)
    
    // Add request logging middleware
    loggedRouter := handlers.LoggingHandler(os.Stdout, router)
    
    // Add recovery middleware
    return handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))(loggedRouter)
}
