package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	SrvPort  string
	MaxSize  int
	LogLevel string
}

// ParseFlags parses command line flags and environment variables to create a Config struct.
// It returns an error if parsing fails or if the configuration is invalid.
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

// Validate checks if the configuration values are valid.
func (c Config) Validate() error {
	if c.MaxSize <= 0 {
		return fmt.Errorf("max-size must be positive")
	}
	return nil
}

// getEnvOrDefault retrieves an environment variable or returns a default value if not set.
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvOrDefaultInt retrieves an integer environment variable or returns a default value if not set or invalid.
func getEnvOrDefaultInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
