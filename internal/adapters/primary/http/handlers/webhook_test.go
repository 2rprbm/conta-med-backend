package handlers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/2rprbm/conta-med-backend/config"
	"github.com/2rprbm/conta-med-backend/pkg/logger"
	"github.com/stretchr/testify/assert"
)

// mockLogger implements the logger.Logger interface for testing
type mockLogger struct {
	buffer bytes.Buffer
	logger *log.Logger
}

func newMockLogger() *mockLogger {
	ml := &mockLogger{}
	ml.logger = log.New(&ml.buffer, "", 0)
	return ml
}

func (m *mockLogger) Debug(msg string, fields ...logger.Fields) {
	m.logger.Printf("%s %v", msg, fields)
}

func (m *mockLogger) Info(msg string, fields ...logger.Fields) {
	m.logger.Printf("%s %v", msg, fields)
}

func (m *mockLogger) Warn(msg string, fields ...logger.Fields) {
	m.logger.Printf("%s %v", msg, fields)
}

func (m *mockLogger) Error(msg string, fields ...logger.Fields) {
	m.logger.Printf("%s %v", msg, fields)
}

func (m *mockLogger) Fatal(msg string, fields ...logger.Fields) {
	m.logger.Printf("%s %v", msg, fields)
}

func (m *mockLogger) With(fields logger.Fields) logger.Logger {
	return m
}

func TestVerifyToken(t *testing.T) {
	t.Run("should verify token successfully when token matches", func(t *testing.T) {
		// arrange
		logger := newMockLogger()
		cfg := &config.Config{
			WhatsApp: config.WhatsAppConfig{
				WebhookVerifyToken: "test_token",
			},
		}

		handler := NewWebhookHandler(cfg, logger)

		// Create request with query parameters
		req := httptest.NewRequest("GET", "/webhook/whatsapp?hub.mode=subscribe&hub.verify_token=test_token&hub.challenge=challenge_value", nil)
		recorder := httptest.NewRecorder()

		// act
		handler.VerifyToken(recorder, req)

		// assert
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "challenge_value", recorder.Body.String())
	})

	t.Run("should fail verification when token doesn't match", func(t *testing.T) {
		// arrange
		logger := newMockLogger()
		cfg := &config.Config{
			WhatsApp: config.WhatsAppConfig{
				WebhookVerifyToken: "test_token",
			},
		}

		handler := NewWebhookHandler(cfg, logger)

		// Create request with incorrect token
		req := httptest.NewRequest("GET", "/webhook/whatsapp?hub.mode=subscribe&hub.verify_token=wrong_token&hub.challenge=challenge_value", nil)
		recorder := httptest.NewRecorder()

		// act
		handler.VerifyToken(recorder, req)

		// assert
		assert.Equal(t, http.StatusForbidden, recorder.Code)
	})

	t.Run("should fail verification when mode is incorrect", func(t *testing.T) {
		// arrange
		logger := newMockLogger()
		cfg := &config.Config{
			WhatsApp: config.WhatsAppConfig{
				WebhookVerifyToken: "test_token",
			},
		}

		handler := NewWebhookHandler(cfg, logger)

		// Create request with incorrect mode
		req := httptest.NewRequest("GET", "/webhook/whatsapp?hub.mode=wrong_mode&hub.verify_token=test_token&hub.challenge=challenge_value", nil)
		recorder := httptest.NewRecorder()

		// act
		handler.VerifyToken(recorder, req)

		// assert
		assert.Equal(t, http.StatusForbidden, recorder.Code)
	})
}

func TestReceiveWebhook(t *testing.T) {
	t.Run("should process valid webhook message", func(t *testing.T) {
		// arrange
		logger := newMockLogger()
		cfg := &config.Config{
			WhatsApp: config.WhatsAppConfig{
				AppSecret: "test_secret",
			},
			Server: config.ServerConfig{
				Environment: "production",
			},
		}

		handler := NewWebhookHandler(cfg, logger)

		// Create webhook payload
		payload := WebhookPayload{
			Object: "whatsapp_business_account",
			Entry: []struct {
				ID      string `json:"id"`
				Changes []struct {
					Value struct {
						MessagingProduct string `json:"messaging_product"`
						Metadata         struct {
							PhoneNumberID      string `json:"phone_number_id"`
							DisplayPhoneNumber string `json:"display_phone_number"`
						} `json:"metadata"`
						Messages []struct {
							ID        string `json:"id"`
							From      string `json:"from"`
							Timestamp string `json:"timestamp"`
							Type      string `json:"type"`
							Text      struct {
								Body string `json:"body"`
							} `json:"text,omitempty"`
						} `json:"messages,omitempty"`
					} `json:"value"`
					Field string `json:"field"`
				} `json:"changes"`
			}{
				{
					ID: "123456789",
					Changes: []struct {
						Value struct {
							MessagingProduct string `json:"messaging_product"`
							Metadata         struct {
								PhoneNumberID      string `json:"phone_number_id"`
								DisplayPhoneNumber string `json:"display_phone_number"`
							} `json:"metadata"`
							Messages []struct {
								ID        string `json:"id"`
								From      string `json:"from"`
								Timestamp string `json:"timestamp"`
								Type      string `json:"type"`
								Text      struct {
									Body string `json:"body"`
								} `json:"text,omitempty"`
							} `json:"messages,omitempty"`
						} `json:"value"`
						Field string `json:"field"`
					}{
						{
							Field: "messages",
							Value: struct {
								MessagingProduct string `json:"messaging_product"`
								Metadata         struct {
									PhoneNumberID      string `json:"phone_number_id"`
									DisplayPhoneNumber string `json:"display_phone_number"`
								} `json:"metadata"`
								Messages []struct {
									ID        string `json:"id"`
									From      string `json:"from"`
									Timestamp string `json:"timestamp"`
									Type      string `json:"type"`
									Text      struct {
										Body string `json:"body"`
									} `json:"text,omitempty"`
								} `json:"messages,omitempty"`
							}{
								MessagingProduct: "whatsapp",
								Metadata: struct {
									PhoneNumberID      string `json:"phone_number_id"`
									DisplayPhoneNumber string `json:"display_phone_number"`
								}{
									PhoneNumberID:      "123456789",
									DisplayPhoneNumber: "+1234567890",
								},
								Messages: []struct {
									ID        string `json:"id"`
									From      string `json:"from"`
									Timestamp string `json:"timestamp"`
									Type      string `json:"type"`
									Text      struct {
										Body string `json:"body"`
									} `json:"text,omitempty"`
								}{
									{
										ID:        "wamid.123456789",
										From:      "554491234567",
										Timestamp: "1617356451",
										Type:      "text",
										Text: struct {
											Body string `json:"body"`
										}{
											Body: "Hello, world!",
										},
									},
								},
							},
						},
					},
				},
			},
		}

		// Convert payload to JSON
		jsonPayload, _ := json.Marshal(payload)

		// Create request with signature
		req := httptest.NewRequest("POST", "/webhook/whatsapp", bytes.NewBuffer(jsonPayload))

		// Add signature to the request
		mac := hmac.New(sha256.New, []byte("test_secret"))
		mac.Write(jsonPayload)
		signature := hex.EncodeToString(mac.Sum(nil))
		req.Header.Set("X-Hub-Signature-256", "sha256="+signature)
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		// act
		handler.ReceiveWebhook(recorder, req)

		// assert
		assert.Equal(t, http.StatusOK, recorder.Code)
		
		// Verificar apenas se contém a mensagem principal, sem se preocupar com o formato exato
		logOutput := logger.buffer.String()
		assert.Contains(t, logOutput, "Received message")
	})

	t.Run("should reject webhook with invalid signature", func(t *testing.T) {
		// arrange
		logger := newMockLogger()
		cfg := &config.Config{
			WhatsApp: config.WhatsAppConfig{
				AppSecret: "test_secret",
			},
			Server: config.ServerConfig{
				Environment: "production",
			},
		}

		handler := NewWebhookHandler(cfg, logger)

		// Simple payload
		payload := `{"object":"whatsapp_business_account"}`

		// Create request with invalid signature
		req := httptest.NewRequest("POST", "/webhook/whatsapp", bytes.NewBufferString(payload))
		req.Header.Set("X-Hub-Signature-256", "sha256=invalid_signature")
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		// act
		handler.ReceiveWebhook(recorder, req)

		// assert
		assert.Equal(t, http.StatusForbidden, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Invalid signature")
	})

	t.Run("should reject webhook with incorrect object type", func(t *testing.T) {
		// arrange
		logger := newMockLogger()
		cfg := &config.Config{
			WhatsApp: config.WhatsAppConfig{
				AppSecret: "test_secret",
			},
			// Explicitly set to development mode with a valid app secret to bypass signature verification
			Server: config.ServerConfig{
				Environment: "development",
			},
		}

		handler := NewWebhookHandler(cfg, logger)

		// Payload with wrong object type
		payload := `{"object":"instagram"}`

		// Create request
		req := httptest.NewRequest("POST", "/webhook/whatsapp", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")

		// Add valid signature to bypass signature verification
		mac := hmac.New(sha256.New, []byte("test_secret"))
		mac.Write([]byte(payload))
		signature := hex.EncodeToString(mac.Sum(nil))
		req.Header.Set("X-Hub-Signature-256", "sha256="+signature)

		recorder := httptest.NewRecorder()

		// act
		handler.ReceiveWebhook(recorder, req)

		// assert
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Unexpected webhook object")
	})
}
