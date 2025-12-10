package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Type     string // "postgres" or "sqlite"
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	SecretKey     string
	TokenDuration int // in seconds
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8081"),
		},
		Database: DatabaseConfig{
			Type:     getEnv("DB_TYPE", "sqlite"), // Default to sqlite
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "expense_user"),
			Password: getEnv("DB_PASSWORD", "expense_password"),
			DBName:   getEnv("DB_NAME", "expense_tracker.db"), // .db for sqlite
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			SecretKey:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			TokenDuration: getEnvAsInt("JWT_TOKEN_DURATION", 7*24*3600), // 7 days
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}