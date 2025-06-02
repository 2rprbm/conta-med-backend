package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Test cases using the format described in the rules
	t.Run("should load default values when environment variables are not set", func(t *testing.T) {
		// arrange
		// Clear any environment variables that might be set
		os.Clearenv()

		// act
		cfg, err := LoadConfig()

		// assert
		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, "8080", cfg.Server.Port)
		assert.Equal(t, "localhost", cfg.Server.Host)
		assert.Equal(t, "development", cfg.Server.Environment)
		assert.Equal(t, 10*time.Second, cfg.Server.ReadTimeout)
		assert.Equal(t, 10*time.Second, cfg.Server.WriteTimeout)
		assert.Equal(t, 60*time.Second, cfg.Server.IdleTimeout)
		assert.Equal(t, 10*time.Second, cfg.Server.ShutdownPeriod)
		assert.Equal(t, "", cfg.MongoDB.URI)
		assert.Equal(t, "medical_scheduler", cfg.MongoDB.Database)
		assert.Equal(t, 10*time.Second, cfg.MongoDB.Timeout)
		assert.Equal(t, "", cfg.WhatsApp.AppID)
		assert.Equal(t, "", cfg.WhatsApp.AppSecret)
		assert.Equal(t, "", cfg.WhatsApp.AccessToken)
		assert.Equal(t, "", cfg.WhatsApp.PhoneNumberID)
		assert.Equal(t, "", cfg.WhatsApp.WebhookVerifyToken)
		assert.Equal(t, "info", cfg.Logging.Level)
	})

	t.Run("should load custom values from environment variables when set", func(t *testing.T) {
		// arrange
		os.Clearenv()
		os.Setenv("SERVER_PORT", "3000")
		os.Setenv("SERVER_HOST", "api.example.com")
		os.Setenv("ENV", "production")
		os.Setenv("SERVER_READ_TIMEOUT", "30")
		os.Setenv("MONGODB_URI", "mongodb://localhost:27017")
		os.Setenv("MONGODB_DATABASE", "test_db")
		os.Setenv("WHATSAPP_APP_ID", "12345")
		os.Setenv("WHATSAPP_WEBHOOK_VERIFY_TOKEN", "secret_token")
		os.Setenv("LOG_LEVEL", "debug")

		// act
		cfg, err := LoadConfig()

		// assert
		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, "3000", cfg.Server.Port)
		assert.Equal(t, "api.example.com", cfg.Server.Host)
		assert.Equal(t, "production", cfg.Server.Environment)
		assert.Equal(t, 30*time.Second, cfg.Server.ReadTimeout)
		assert.Equal(t, "mongodb://localhost:27017", cfg.MongoDB.URI)
		assert.Equal(t, "test_db", cfg.MongoDB.Database)
		assert.Equal(t, "12345", cfg.WhatsApp.AppID)
		assert.Equal(t, "secret_token", cfg.WhatsApp.WebhookVerifyToken)
		assert.Equal(t, "debug", cfg.Logging.Level)
	})

	t.Run("should handle invalid numeric values in environment variables", func(t *testing.T) {
		// arrange
		os.Clearenv()
		os.Setenv("SERVER_READ_TIMEOUT", "invalid")

		// act
		cfg, err := LoadConfig()

		// assert
		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		// Should use default value when the environment variable is invalid
		assert.Equal(t, 10*time.Second, cfg.Server.ReadTimeout)
	})
}
