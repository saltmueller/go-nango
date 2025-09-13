// Package config provides internal configuration utilities for the go-nango project.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// AppConfig holds application-wide configuration.
type AppConfig struct {
	NangoAPIKey  string
	NangoBaseURL string
	Timeout      time.Duration
	LogLevel     string
	Port         int
}

// LoadFromEnv loads configuration from environment variables.
func LoadFromEnv() (*AppConfig, error) {
	config := &AppConfig{
		NangoBaseURL: getEnvWithDefault("NANGO_BASE_URL", "https://api.nango.dev"),
		LogLevel:     getEnvWithDefault("LOG_LEVEL", "info"),
	}

	// Required environment variables
	config.NangoAPIKey = os.Getenv("NANGO_API_KEY")
	if config.NangoAPIKey == "" {
		return nil, fmt.Errorf("NANGO_API_KEY environment variable is required")
	}

	// Parse timeout
	timeoutStr := getEnvWithDefault("TIMEOUT", "30s")
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return nil, fmt.Errorf("invalid TIMEOUT format: %w", err)
	}
	config.Timeout = timeout

	// Parse port
	portStr := getEnvWithDefault("PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid PORT format: %w", err)
	}
	config.Port = port

	return config, nil
}

// getEnvWithDefault returns the value of the environment variable or a default value if not set.
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Validate checks if the configuration is valid.
func (c *AppConfig) Validate() error {
	if c.NangoAPIKey == "" {
		return fmt.Errorf("NangoAPIKey is required")
	}

	if c.NangoBaseURL == "" {
		return fmt.Errorf("NangoBaseURL is required")
	}

	if c.Timeout <= 0 {
		return fmt.Errorf("Timeout must be positive")
	}

	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("Port must be between 1 and 65535")
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if !validLogLevels[c.LogLevel] {
		return fmt.Errorf("LogLevel must be one of: debug, info, warn, error")
	}

	return nil
}
