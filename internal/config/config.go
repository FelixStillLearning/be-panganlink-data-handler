package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost            string
	DBPort            string
	DBUser            string
	DBPass            string
	DBName            string
	ServerPort        string
	JWTSecret         string
	AIServiceURL      string 
	MidtransServerKey    string 
	MidtransIsProduction bool
	AzureAccountName     string
	AzureAccountKey      string
	AzureContainerName   string
	CloudinaryCloudName  string
	CloudinaryAPIKey     string
	CloudinaryAPISecret  string
}

func LoadConfig() *Config {
	err := godotenv.Load() // Try current dir first
	if err != nil {
		_ = godotenv.Load("../../.env") // Try root dir if running from cmd/server
	}

	return &Config{
		DBHost:            getEnv("DB_HOST", "mysql"), // In docker compose it's 'mysql'
		DBPort:            getEnv("DB_PORT", "3306"),
		DBUser:            getEnv("DB_USER", "mysqluser"),
		DBPass:            getEnv("DB_PASS", "password"),
		DBName:            getEnv("DB_NAME", "panganlink_db"),
		ServerPort:        getEnv("SERVER_PORT", "8080"),
		JWTSecret:         getEnv("JWT_SECRET", "supersecretkey_panganlink"),
		AIServiceURL:      getEnv("AI_SERVICE_URL", "http://ai-service:8000/api/v1"),
		MidtransServerKey: getEnv("MIDTRANS_SERVER_KEY", "SB-Mid-server-YOURKEYHERE"),
		MidtransIsProduction: getEnv("MIDTRANS_IS_PRODUCTION", "false") == "true",
		AzureAccountName:  getEnv("AZURE_STORAGE_ACCOUNT_NAME", ""),
		AzureAccountKey:   getEnv("AZURE_STORAGE_ACCOUNT_KEY", ""),
		AzureContainerName: getEnv("AZURE_STORAGE_CONTAINER_NAME", "panganlink"),
		CloudinaryCloudName: getEnv("CLOUDINARY_CLOUD_NAME", ""),
		CloudinaryAPIKey:    getEnv("CLOUDINARY_API_KEY", ""),
		CloudinaryAPISecret: getEnv("CLOUDINARY_API_SECRET", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
