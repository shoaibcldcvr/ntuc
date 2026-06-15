package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL string
	ServerPort  string
	LogLevel    string
	JWTSecret   string
}

func Load() (*Config, error) {
	cfg := &Config{
		DatabaseURL: getEnv("DATABASE_URL", "root:password@tcp(localhost:3306)/ntuc?parseTime=true"),
		ServerPort:  getEnv("SERVER_PORT", ":8080"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
	}

	if cfg.JWTSecret == "your-secret-key" {
		fmt.Println("WARNING: Using default JWT_SECRET. Set JWT_SECRET environment variable in production.")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
