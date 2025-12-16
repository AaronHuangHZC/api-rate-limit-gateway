package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration
type Config struct {
	// Server configuration
	Server ServerConfig

	// Redis configuration
	Redis RedisConfig

	// PostgreSQL configuration
	Postgres PostgresConfig

	// Gateway behavior
	Gateway GatewayConfig
}

// ServerConfig holds HTTP server settings
type ServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// RedisConfig holds Redis connection settings
type RedisConfig struct {
	Addr         string
	Password     string
	DB           int
	MaxRetries   int
	PoolSize     int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// PostgresConfig holds PostgreSQL connection settings
type PostgresConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// GatewayConfig holds gateway-specific behavior settings
type GatewayConfig struct {
	// FailurePolicy defines behavior when Redis is unavailable
	// "fail-open" allows requests through, "fail-closed" blocks all requests
	FailurePolicy string

	// UpstreamTimeout is the timeout for proxying requests to upstream services
	UpstreamTimeout time.Duration

	// PolicyCacheTTL is how long to cache policy lookups in memory
	PolicyCacheTTL time.Duration
}

// Load reads configuration from environment variables with sensible defaults
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			Port:         getEnvAsInt("SERVER_PORT", 8080),
			ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:  getEnvAsDuration("SERVER_IDLE_TIMEOUT", 120*time.Second),
		},
		Redis: RedisConfig{
			Addr:         getEnv("REDIS_ADDR", "localhost:6379"),
			Password:     getEnv("REDIS_PASSWORD", ""),
			DB:           getEnvAsInt("REDIS_DB", 0),
			MaxRetries:   getEnvAsInt("REDIS_MAX_RETRIES", 3),
			PoolSize:     getEnvAsInt("REDIS_POOL_SIZE", 10),
			DialTimeout:  getEnvAsDuration("REDIS_DIAL_TIMEOUT", 5*time.Second),
			ReadTimeout:  getEnvAsDuration("REDIS_READ_TIMEOUT", 3*time.Second),
			WriteTimeout: getEnvAsDuration("REDIS_WRITE_TIMEOUT", 3*time.Second),
		},
		Postgres: PostgresConfig{
			Host:            getEnv("POSTGRES_HOST", "localhost"),
			Port:            getEnvAsInt("POSTGRES_PORT", 5432),
			User:            getEnv("POSTGRES_USER", "gateway"),
			Password:        getEnv("POSTGRES_PASSWORD", "gateway_password"),
			Database:        getEnv("POSTGRES_DB", "gateway_db"),
			SSLMode:         getEnv("POSTGRES_SSLMODE", "disable"),
			MaxOpenConns:    getEnvAsInt("POSTGRES_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("POSTGRES_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvAsDuration("POSTGRES_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		Gateway: GatewayConfig{
			FailurePolicy:   getEnv("GATEWAY_FAILURE_POLICY", "fail-closed"),
			UpstreamTimeout: getEnvAsDuration("GATEWAY_UPSTREAM_TIMEOUT", 30*time.Second),
			PolicyCacheTTL:  getEnvAsDuration("GATEWAY_POLICY_CACHE_TTL", 5*time.Minute),
		},
	}

	// Validate failure policy
	if cfg.Gateway.FailurePolicy != "fail-open" && cfg.Gateway.FailurePolicy != "fail-closed" {
		return nil, fmt.Errorf("invalid GATEWAY_FAILURE_POLICY: must be 'fail-open' or 'fail-closed'")
	}

	return cfg, nil
}

// ServerAddr returns the full server address (host:port)
func (c *Config) ServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// PostgresDSN returns the PostgreSQL connection string
func (c *Config) PostgresDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.User,
		c.Postgres.Password,
		c.Postgres.Database,
		c.Postgres.SSLMode,
	)
}

// Helper functions for environment variable parsing

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

