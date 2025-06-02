package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/2rprbm/conta-med-backend/internal/core/ports"
	"github.com/2rprbm/conta-med-backend/pkg/logger"
	"github.com/go-chi/chi/v5"
)

// WebhookHandler handles incoming webhook requests from WhatsApp
type WebhookHandler struct {
	chatbotService ports.ChatbotService
	whatsappService ports.WhatsAppService
	log            logger.Logger
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(chatbotService ports.ChatbotService, whatsappService ports.WhatsAppService, log logger.Logger) *WebhookHandler {
	return &WebhookHandler{
		chatbotService: chatbotService,
		whatsappService: whatsappService,
		log:            log,
	}
}

// RegisterRoutes registers webhook routes
func (h *WebhookHandler) RegisterRoutes(r chi.Router) {
	r.Get("/webhook", h.VerifyWebhook)
	r.Post("/webhook", h.HandleWebhook)
}

// VerifyWebhook handles the webhook verification request from WhatsApp
func (h *WebhookHandler) VerifyWebhook(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Received webhook verification request")

	mode := r.URL.Query().Get("hub.mode")
	token := r.URL.Query().Get("hub.verify_token")
	challenge := r.URL.Query().Get("hub.challenge")

	verified, challenge := h.whatsappService.VerifyWebhook(mode, token, challenge)
	if !verified {
		h.log.Error("Failed to verify webhook", logger.Fields{
			"mode":  mode,
			"token": token,
		})
		http.Error(w, "Verification failed", http.StatusUnauthorized)
		return
	}

	h.log.Info("Webhook verified successfully")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(challenge))
}

// WhatsAppPayload represents the structure of a WhatsApp webhook payload
type WhatsAppPayload struct {
	Object string `json:"object"`
	Entry  []struct {
		ID      string `json:"id"`
		Changes []struct {
			Value struct {
				MessagingProduct string `json:"messaging_product"`
				Metadata         struct {
					PhoneNumberID string `json:"phone_number_id"`
					DisplayPhoneNumber string `json:"display_phone_number"`
				} `json:"metadata"`
				Contacts []struct {
					Profile struct {
						Name string `json:"name"`
					} `json:"profile"`
					WaID string `json:"wa_id"`
				} `json:"contacts"`
				Messages []struct {
					ID        string `json:"id"`
					From      string `json:"from"`
					Timestamp string `json:"timestamp"`
					Type      string `json:"type"`
					Text      struct {
						Body string `json:"body"`
					} `json:"text"`
				} `json:"messages"`
			} `json:"value"`
			Field string `json:"field"`
		} `json:"changes"`
	} `json:"entry"`
}

// HandleWebhook processes webhook notifications from WhatsApp
func (h *WebhookHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Received webhook notification")

	var payload WhatsAppPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.log.Error("Failed to decode webhook payload", logger.Fields{
			"error": err.Error(),
		})
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Acknowledge webhook receipt immediately
	w.WriteHeader(http.StatusOK)

	// Process the webhook asynchronously
	go h.processWebhook(payload)
}

// processWebhook handles the webhook payload processing
func (h *WebhookHandler) processWebhook(payload WhatsAppPayload) {
	if payload.Object != "whatsapp_business_account" {
		h.log.Info("Ignoring non-WhatsApp webhook", logger.Fields{
			"object": payload.Object,
		})
		return
	}

	for _, entry := range payload.Entry {
		for _, change := range entry.Changes {
			if change.Field != "messages" {
				continue
			}

			for _, message := range change.Value.Messages {
				h.log.Info("Processing message", logger.Fields{
					"from":      message.From,
					"messageId": message.ID,
					"type":      message.Type,
				})

				// Only process text messages
				if message.Type != "text" {
					h.log.Info("Ignoring non-text message", logger.Fields{
						"type": message.Type,
					})
					continue
				}

				// Clean the phone number (remove any "+" prefix)
				phoneNumber := strings.TrimPrefix(message.From, "+")

				// Process the message
				if err := h.chatbotService.HandleIncomingMessage(r.Context(), phoneNumber, message.Text.Body); err != nil {
					h.log.Error("Failed to process message", logger.Fields{
						"error": err.Error(),
						"from":  message.From,
					})
				}
			}
		}
	}
} 