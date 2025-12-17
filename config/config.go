// Package config provides utilities for parsing common environment variables.
package config

import (
	"os"
	"strconv"
	"time"
)

// Common holds common configuration fields shared across services.
type Common struct {
	ServiceName string
	LogLevel    string
	Environment string
}

// LoadCommon loads common configuration from environment variables.
func LoadCommon() Common {
	return Common{
		ServiceName: GetEnv("SERVICE_NAME", "unknown"),
		LogLevel:    GetEnv("LOG_LEVEL", "info"),
		Environment: GetEnv("ENVIRONMENT", "development"),
	}
}

// GetEnv retrieves an environment variable or returns a default value.
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvInt retrieves an environment variable as an integer or returns a default.
func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

// GetEnvInt32 retrieves an environment variable as an int32 or returns a default.
func GetEnvInt32(key string, defaultValue int32) int32 {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.ParseInt(value, 10, 32); err == nil {
			return int32(i)
		}
	}
	return defaultValue
}

// GetEnvDuration retrieves an environment variable as a duration or returns a default.
func GetEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return defaultValue
}

// GetEnvBool retrieves an environment variable as a boolean or returns a default.
func GetEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return defaultValue
}

