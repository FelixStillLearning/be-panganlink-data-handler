package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost       string
	DBPort       string
	DBUser       string
	DBPass       string
	DBName       string
	ServerPort   string
	JWTSecret    string
	AIServiceURL string // Added for Phase 3.3
}

func LoadConfig() *Config {
	_ = godotenv.Load() // Ignore error if .env not found (using env vars from OS/Docker)

	return &Config{
		DBHost:       getEnv("DB_HOST", "mysql"), // In docker compose it's 'mysql'
		DBPort:       getEnv("DB_PORT", "3306"),
		DBUser:       getEnv("DB_USER", "mysqluser"),
		DBPass:       getEnv("DB_PASS", "password"),
		DBName:       getEnv("DB_NAME", "panganlink_db"),
		ServerPort:   getEnv("SERVER_PORT", "8080"),
		JWTSecret:    getEnv("JWT_SECRET", "supersecretkey_panganlink"),
		AIServiceURL: getEnv("AI_SERVICE_URL", "http://ai-service:8000/api/v1"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
