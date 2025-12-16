package config

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestLoad_Defaults(t *testing.T) {
	// Clear all relevant env vars to test defaults
	envVars := []string{
		"SERVER_HOST", "SERVER_PORT", "SERVER_READ_TIMEOUT",
		"REDIS_ADDR", "REDIS_PASSWORD",
		"POSTGRES_HOST", "POSTGRES_PORT",
		"GATEWAY_FAILURE_POLICY",
	}
	for _, key := range envVars {
		os.Unsetenv(key)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Test server defaults
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("expected Server.Host = '0.0.0.0', got %q", cfg.Server.Host)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("expected Server.Port = 8080, got %d", cfg.Server.Port)
	}

	// Test Redis defaults
	if cfg.Redis.Addr != "localhost:6379" {
		t.Errorf("expected Redis.Addr = 'localhost:6379', got %q", cfg.Redis.Addr)
	}

	// Test failure policy default
	if cfg.Gateway.FailurePolicy != "fail-closed" {
		t.Errorf("expected Gateway.FailurePolicy = 'fail-closed', got %q", cfg.Gateway.FailurePolicy)
	}
}

func TestLoad_EnvironmentVariables(t *testing.T) {
	// Set test env vars
	os.Setenv("SERVER_PORT", "9000")
	os.Setenv("SERVER_READ_TIMEOUT", "15s")
	os.Setenv("GATEWAY_FAILURE_POLICY", "fail-open")
	defer func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("SERVER_READ_TIMEOUT")
		os.Unsetenv("GATEWAY_FAILURE_POLICY")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.Server.Port != 9000 {
		t.Errorf("expected Server.Port = 9000, got %d", cfg.Server.Port)
	}
	if cfg.Server.ReadTimeout != 15*time.Second {
		t.Errorf("expected Server.ReadTimeout = 15s, got %v", cfg.Server.ReadTimeout)
	}
	if cfg.Gateway.FailurePolicy != "fail-open" {
		t.Errorf("expected Gateway.FailurePolicy = 'fail-open', got %q", cfg.Gateway.FailurePolicy)
	}
}

func TestLoad_InvalidFailurePolicy(t *testing.T) {
	os.Setenv("GATEWAY_FAILURE_POLICY", "invalid-policy")
	defer os.Unsetenv("GATEWAY_FAILURE_POLICY")

	_, err := Load()
	if err == nil {
		t.Error("expected error for invalid failure policy, got nil")
	}
}

func TestServerAddr(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Host: "127.0.0.1",
			Port: 8080,
		},
	}

	expected := "127.0.0.1:8080"
	if got := cfg.ServerAddr(); got != expected {
		t.Errorf("ServerAddr() = %q, want %q", got, expected)
	}
}

func TestPostgresDSN(t *testing.T) {
	cfg := &Config{
		Postgres: PostgresConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "testuser",
			Password: "testpass",
			Database: "testdb",
			SSLMode:  "require",
		},
	}

	dsn := cfg.PostgresDSN()
	if dsn == "" {
		t.Error("PostgresDSN() returned empty string")
	}
	// Just check it contains key components
	if !strings.Contains(dsn, "host=localhost") {
		t.Errorf("PostgresDSN() missing host: %s", dsn)
	}
	if !strings.Contains(dsn, "user=testuser") {
		t.Errorf("PostgresDSN() missing user: %s", dsn)
	}
}

