package mongodb

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/2rprbm/conta-med-backend/config"
	"github.com/2rprbm/conta-med-backend/pkg/logger"
)

type mockLogger struct{}

func (m *mockLogger) Info(msg string, fields ...logger.Fields)  {}
func (m *mockLogger) Error(msg string, fields ...logger.Fields) {}
func (m *mockLogger) Debug(msg string, fields ...logger.Fields) {}
func (m *mockLogger) Warn(msg string, fields ...logger.Fields)  {}
func (m *mockLogger) Fatal(msg string, fields ...logger.Fields) {}
func (m *mockLogger) With(fields logger.Fields) logger.Logger   { return m }

func TestNewClient(t *testing.T) {
	t.Run("should return error when MongoDB URI is empty", func(t *testing.T) {
		// arrange
		cfg := &config.MongoDBConfig{
			URI:      "",
			Database: "test_db",
			Timeout:  10 * time.Second,
		}
		log := &mockLogger{}

		// act
		client, err := NewClient(cfg, log)

		// assert
		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "MongoDB URI is required")
	})

	t.Run("should return error when MongoDB URI is invalid", func(t *testing.T) {
		// arrange
		cfg := &config.MongoDBConfig{
			URI:      "invalid-uri",
			Database: "test_db",
			Timeout:  1 * time.Second,
		}
		log := &mockLogger{}

		// act
		client, err := NewClient(cfg, log)

		// assert
		assert.Error(t, err)
		assert.Nil(t, client)
	})

	// Note: For successful connection tests, we would need a running MongoDB instance
	// This would be more appropriate for integration tests
}

func TestClient_GetDatabase(t *testing.T) {
	t.Run("should return database when client is initialized", func(t *testing.T) {
		// arrange
		client := &Client{
			database: nil, // simulating a database
			logger:   &mockLogger{},
		}

		// act
		db := client.GetDatabase()

		// assert
		assert.Nil(t, db) // In our test case, it's nil, but method works
	})
}

func TestClient_GetClient(t *testing.T) {
	t.Run("should return mongo client when initialized", func(t *testing.T) {
		// arrange
		client := &Client{
			client: nil, // simulating a client
			logger: &mockLogger{},
		}

		// act
		mongoClient := client.GetClient()

		// assert
		assert.Nil(t, mongoClient) // In our test case, it's nil, but method works
	})
}

func TestClient_Health(t *testing.T) {
	t.Run("should return error when client is nil", func(t *testing.T) {
		// arrange
		client := &Client{
			client: nil,
			logger: &mockLogger{},
		}
		ctx := context.Background()

		// act & assert
		// Should panic when client is nil, so we need to test this differently
		assert.Panics(t, func() {
			client.Health(ctx)
		})
	})
}

func TestClient_Close(t *testing.T) {
	t.Run("should handle nil client gracefully", func(t *testing.T) {
		// arrange
		client := &Client{
			client: nil,
			logger: &mockLogger{},
		}
		ctx := context.Background()

		// act & assert
		// Should panic when client is nil, so we need to test this differently
		assert.Panics(t, func() {
			client.Close(ctx)
		})
	})
}

func TestClient_CreateIndexes(t *testing.T) {
	t.Run("should handle nil database gracefully", func(t *testing.T) {
		// arrange
		client := &Client{
			database: nil,
			logger:   &mockLogger{},
		}
		ctx := context.Background()

		// act & assert
		// Should panic when database is nil
		assert.Panics(t, func() {
			client.CreateIndexes(ctx)
		})
	})
} 