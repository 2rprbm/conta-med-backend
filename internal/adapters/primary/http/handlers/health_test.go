package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mongoclient "github.com/2rprbm/conta-med-backend/pkg/mongodb"
)

type mockMongoClient struct {
	healthError error
}

func (m *mockMongoClient) Health(ctx context.Context) error {
	return m.healthError
}

func TestNewHealthHandler(t *testing.T) {
	t.Run("should create health handler successfully", func(t *testing.T) {
		// arrange
		log := &mockLogger{}
		var mongoClient *mongoclient.Client = nil

		// act
		handler := NewHealthHandler(mongoClient, log)

		// assert
		assert.NotNil(t, handler)
		assert.Equal(t, mongoClient, handler.mongoClient)
		assert.Equal(t, log, handler.logger)
	})
}

func TestHealthHandler_CheckHealth(t *testing.T) {
	t.Run("should return healthy status when MongoDB is working", func(t *testing.T) {
		// arrange
		log := &mockLogger{}
		// Since we can't easily mock the mongoclient.Client interface, we'll test with nil client
		handler := NewHealthHandler(nil, log)
		
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()

		// act
		handler.CheckHealth(w, req)

		// assert
		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response HealthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "degraded", response.Status)
		assert.Equal(t, "not_configured", response.Services["mongodb"])
		assert.Equal(t, "v0.2.0", response.Version)
		assert.NotZero(t, response.Timestamp)
	})

	t.Run("should return proper content type and structure", func(t *testing.T) {
		// arrange
		log := &mockLogger{}
		handler := NewHealthHandler(nil, log)
		
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()

		// act
		handler.CheckHealth(w, req)

		// assert
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		
		var response HealthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		
		// Verify structure
		assert.NotEmpty(t, response.Status)
		assert.NotNil(t, response.Services)
		assert.NotEmpty(t, response.Version)
		assert.NotZero(t, response.Timestamp)
	})
} 