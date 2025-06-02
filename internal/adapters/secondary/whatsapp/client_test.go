package whatsapp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/2rprbm/conta-med-backend/config"
	"github.com/2rprbm/conta-med-backend/pkg/logger"
	"github.com/stretchr/testify/assert"
)

// mockLogger implements the logger.Logger interface for testing
type mockLogger struct {
	debugMessages []string
	infoMessages  []string
	warnMessages  []string
	errorMessages []string
	fatalMessages []string
}

func newMockLogger() *mockLogger {
	return &mockLogger{
		debugMessages: []string{},
		infoMessages:  []string{},
		warnMessages:  []string{},
		errorMessages: []string{},
		fatalMessages: []string{},
	}
}

func (m *mockLogger) Debug(msg string, fields ...logger.Fields) {
	m.debugMessages = append(m.debugMessages, msg)
}

func (m *mockLogger) Info(msg string, fields ...logger.Fields) {
	m.infoMessages = append(m.infoMessages, msg)
}

func (m *mockLogger) Warn(msg string, fields ...logger.Fields) {
	m.warnMessages = append(m.warnMessages, msg)
}

func (m *mockLogger) Error(msg string, fields ...logger.Fields) {
	m.errorMessages = append(m.errorMessages, msg)
}

func (m *mockLogger) Fatal(msg string, fields ...logger.Fields) {
	m.fatalMessages = append(m.fatalMessages, msg)
}

func (m *mockLogger) With(fields logger.Fields) logger.Logger {
	return m
}

func TestSendTextMessage(t *testing.T) {
	t.Run("should send a text message successfully", func(t *testing.T) {
		// arrange
		// Create a test server that simulates the WhatsApp API
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request method
			assert.Equal(t, "POST", r.Method)

			// Verify request headers
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			assert.Equal(t, "Bearer test_token", r.Header.Get("Authorization"))

			// Verify request body
			var payload TextMessage
			json.NewDecoder(r.Body).Decode(&payload)
			assert.Equal(t, "whatsapp", payload.MessagingProduct)
			assert.Equal(t, "individual", payload.RecipientType)
			assert.Equal(t, "554499887766", payload.To)
			assert.Equal(t, "text", payload.Type)
			assert.Equal(t, "Hello from test", payload.Text.Body)

			// Return success response
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"success": true}`))
		}))
		defer server.Close()

		// Create client with test configuration
		logger := newMockLogger()
		cfg := &config.Config{
			WhatsApp: config.WhatsAppConfig{
				PhoneNumberID: "12345",
				AccessToken:   "test_token",
			},
		}

		client := NewClient(cfg, logger)
		// Override the HTTP client to use our test server
		client.HttpClient = server.Client()
		// Override the API URL to use our test server
		client.APIURL = server.URL + "/%s/messages"

		// act
		err := client.SendTextMessage("554499887766", "Hello from test")

		// assert
		assert.NoError(t, err)
		assert.Contains(t, logger.debugMessages[0], "Sending WhatsApp message")
		assert.Contains(t, logger.infoMessages[0], "Message sent successfully")
	})

	t.Run("should handle error when phone number ID is missing", func(t *testing.T) {
		// arrange
		logger := newMockLogger()
		cfg := &config.Config{
			WhatsApp: config.WhatsAppConfig{
				PhoneNumberID: "", // Empty phone number ID
				AccessToken:   "test_token",
			},
		}

		client := NewClient(cfg, logger)

		// act
		err := client.SendTextMessage("554499887766", "Hello from test")

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "phone number ID not configured")
	})

	t.Run("should handle API error response", func(t *testing.T) {
		// arrange
		// Create a test server that simulates the WhatsApp API returning an error
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Return error response
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": {"message": "Invalid parameter"}}`))
		}))
		defer server.Close()

		// Create client with test configuration
		logger := newMockLogger()
		cfg := &config.Config{
			WhatsApp: config.WhatsAppConfig{
				PhoneNumberID: "12345",
				AccessToken:   "test_token",
			},
		}

		client := NewClient(cfg, logger)
		// Override the HTTP client to use our test server
		client.HttpClient = server.Client()
		// Override the API URL to use our test server
		client.APIURL = server.URL + "/%s/messages"

		// act
		err := client.SendTextMessage("554499887766", "Hello from test")

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}
