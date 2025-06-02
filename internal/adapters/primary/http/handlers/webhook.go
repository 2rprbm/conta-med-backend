package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/2rprbm/conta-med-backend/config"
	"github.com/2rprbm/conta-med-backend/pkg/logger"
)

// WebhookHandler handles WhatsApp webhook requests
type WebhookHandler struct {
	config *config.Config
	logger logger.Logger
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(cfg *config.Config, log logger.Logger) *WebhookHandler {
	return &WebhookHandler{
		config: cfg,
		logger: log,
	}
}

// VerifyToken handles GET requests to verify the webhook
func (h *WebhookHandler) VerifyToken(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	mode := r.URL.Query().Get("hub.mode")
	token := r.URL.Query().Get("hub.verify_token")
	challenge := r.URL.Query().Get("hub.challenge")

	h.logger.Debug("Received webhook verification request: mode=%s, token=%s", mode, token)

	// Check mode and token
	if mode == "subscribe" && token == h.config.WhatsApp.WebhookVerifyToken {
		h.logger.Info("Webhook verified successfully")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(challenge))
		return
	}

	// Verification failed
	h.logger.Warn("Webhook verification failed: invalid token or mode")
	http.Error(w, "Verification failed", http.StatusForbidden)
}

// WebhookPayload represents the webhook payload
type WebhookPayload struct {
	Object string `json:"object"`
	Entry  []struct {
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
	} `json:"entry"`
}

// ReceiveWebhook handles POST requests from WhatsApp
func (h *WebhookHandler) ReceiveWebhook(w http.ResponseWriter, r *http.Request) {
	// Verify signature for security
	if !h.verifySignature(r) {
		h.logger.Warn("Invalid signature")
		http.Error(w, "Invalid signature", http.StatusForbidden)
		return
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("Error reading request body: %v", err)
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}

	// Parse payload
	var payload WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		h.logger.Error("Error parsing webhook payload: %v", err)
		http.Error(w, "Error parsing payload", http.StatusBadRequest)
		return
	}

	// Validate payload
	if payload.Object != "whatsapp_business_account" {
		h.logger.Warn("Received non-WhatsApp webhook: %s", payload.Object)
		http.Error(w, "Unexpected webhook object", http.StatusBadRequest)
		return
	}

	// Process messages - temporary just logging them
	for _, entry := range payload.Entry {
		for _, change := range entry.Changes {
			if change.Field == "messages" {
				for _, message := range change.Value.Messages {
					if message.Type == "text" {
						h.logger.Info("Received message from %s: %s", message.From, message.Text.Body)
						// TODO: Process message and send response via WhatsApp API
					}
				}
			}
		}
	}

	// Return 200 OK to acknowledge receipt
	w.WriteHeader(http.StatusOK)
}

// verifySignature verifies the request signature
func (h *WebhookHandler) verifySignature(r *http.Request) bool {
	// In development mode, skip signature verification if app secret is not set
	if h.config.Server.Environment == "development" && h.config.WhatsApp.AppSecret == "" {
		return true
	}

	// Get signature from header
	signature := r.Header.Get("X-Hub-Signature-256")
	if signature == "" {
		return false
	}

	// Remove 'sha256=' prefix
	signature = strings.TrimPrefix(signature, "sha256=")

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return false
	}

	// Reset body for later reading
	r.Body = io.NopCloser(strings.NewReader(string(body)))

	// Calculate expected signature
	mac := hmac.New(sha256.New, []byte(h.config.WhatsApp.AppSecret))
	mac.Write(body)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	// Compare signatures
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
