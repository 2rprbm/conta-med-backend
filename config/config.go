package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	MongoDB  MongoDBConfig
	WhatsApp WhatsAppConfig
	Logging  LoggingConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port           string
	Host           string
	Environment    string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	ShutdownPeriod time.Duration
}

// MongoDBConfig holds MongoDB configuration
type MongoDBConfig struct {
	URI      string
	Database string
	Timeout  time.Duration
}

// WhatsAppConfig holds WhatsApp API configuration
type WhatsAppConfig struct {
	AppID              string
	AppSecret          string
	AccessToken        string
	PhoneNumberID      string
	WebhookVerifyToken string
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Port:        getEnv("SERVER_PORT", "8080"),
			Host:        getEnv("SERVER_HOST", "localhost"),
			Environment: getEnv("ENV", "development"),
			// Default timeouts
			ReadTimeout:    time.Duration(getEnvAsInt("SERVER_READ_TIMEOUT", 10)) * time.Second,
			WriteTimeout:   time.Duration(getEnvAsInt("SERVER_WRITE_TIMEOUT", 10)) * time.Second,
			IdleTimeout:    time.Duration(getEnvAsInt("SERVER_IDLE_TIMEOUT", 60)) * time.Second,
			ShutdownPeriod: time.Duration(getEnvAsInt("SERVER_SHUTDOWN_PERIOD", 10)) * time.Second,
		},
		MongoDB: MongoDBConfig{
			URI:      getEnv("MONGODB_URI", ""),
			Database: getEnv("MONGODB_DATABASE", "medical_scheduler"),
			Timeout:  time.Duration(getEnvAsInt("MONGODB_TIMEOUT", 10)) * time.Second,
		},
		WhatsApp: WhatsAppConfig{
			AppID:              getEnv("WHATSAPP_APP_ID", ""),
			AppSecret:          getEnv("WHATSAPP_APP_SECRET", ""),
			AccessToken:        getEnv("WHATSAPP_ACCESS_TOKEN", ""),
			PhoneNumberID:      getEnv("WHATSAPP_PHONE_NUMBER_ID", ""),
			WebhookVerifyToken: getEnv("WHATSAPP_WEBHOOK_VERIFY_TOKEN", ""),
		},
		Logging: LoggingConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}

	return cfg, nil
}

// Helper function to get an environment variable or a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper function to get an environment variable as an integer
func getEnvAsInt(key string, defaultValue int) int {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}
