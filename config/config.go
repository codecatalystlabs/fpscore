package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL string
	Port        string
}

func Load() (*Config, error) {
	config := &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:pwaiswa@localhost/fpscore?sslmode=disable"),
		Port:        getEnv("PORT", "5000"),
	}

	if config.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
