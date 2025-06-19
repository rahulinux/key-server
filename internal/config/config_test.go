package config_test

import (
	"testing"

	"github.com/rahulinux/key-server/internal/config"
)

func TestParseFlagsValid(t *testing.T) {
	args := []string{
		"-srv-port=9090",
		"-max-size=512",
		"-log-level=debug",
	}
	cfg, err := config.ParseFlags(args)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.SrvPort != "9090" || cfg.MaxSize != 512 || cfg.LogLevel != "debug" {
		t.Errorf("Unexpected config: %+v", cfg)
	}
}

func TestParseFlagsInvalidMaxSize(t *testing.T) {
	args := []string{"-max-size=-1"}
	_, err := config.ParseFlags(args)
	if err == nil {
		t.Error("Expected error for invalid max-size")
	}
}
