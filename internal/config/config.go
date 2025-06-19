package config

import (
    "flag"
    "fmt"
    "os"
)

type Config struct {
    SrvPort  string
    MaxSize  int
    LogLevel string
}

func ParseFlags(args []string) (Config, error) {
    fs := flag.NewFlagSet("key-server", flag.ContinueOnError)
    fs.SetOutput(os.Stderr)

    srvPort := fs.String("srv-port", getEnvOrDefault("SERVER_PORT", "8080"), "Port to run the server on")
    maxSize := fs.Int("max-size", getEnvOrDefaultInt("MAX_SIZE", 1024), "Maximum number of bytes allowed")
    logLevel := fs.String("log-level", getEnvOrDefault("LOG_LEVEL", "info"), "Log level (debug, info, warn, error)")

    fs.Usage = func() {
        fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
        fs.PrintDefaults()
    }

    if err := fs.Parse(args); err != nil {
        if err == flag.ErrHelp {
            return Config{}, err
        }
        return Config{}, fmt.Errorf("failed to parse flags: %w", err)
    }

    cfg := Config{
        SrvPort:  *srvPort,
        MaxSize:  *maxSize,
        LogLevel: *logLevel,
    }

    if err := cfg.Validate(); err != nil {
        return Config{}, fmt.Errorf("invalid configuration: %w", err)
    }

    return cfg, nil
}

func (c Config) Validate() error {
    if c.MaxSize <= 0 {
        return fmt.Errorf("max-size must be positive")
    }
    if c.MaxSize > 1024*1024 { // 1MB limit
        return fmt.Errorf("max-size cannot exceed 1MB")
    }
    return nil
}

func getEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvOrDefaultInt(key string, defaultValue int) int {
    // Implementation would parse int from env or return default
    return defaultValue
}
