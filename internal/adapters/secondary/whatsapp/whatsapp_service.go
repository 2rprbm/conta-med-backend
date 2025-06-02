package whatsapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/2rprbm/conta-med-backend/internal/core/ports"
	"github.com/2rprbm/conta-med-backend/pkg/logger"
)

// WhatsAppConfig holds the configuration for the WhatsApp service
type WhatsAppConfig struct {
	APIVersion     string
	PhoneNumberID  string
	AccessToken    string
	WebhookToken   string
	BaseURL        string
}

type whatsAppService struct {
	config WhatsAppConfig
	client *http.Client
	log    logger.Logger
}

// NewWhatsAppService creates a new instance of the WhatsApp service
func NewWhatsAppService(config WhatsAppConfig, log logger.Logger) ports.WhatsAppService {
	return &whatsAppService{
		config: config,
		client: &http.Client{},
		log:    log,
	}
}

// SendTextMessage sends a text message to a WhatsApp number
func (s *whatsAppService) SendTextMessage(ctx context.Context, phoneNumber, message string) error {
	s.log.Info("Sending text message", logger.Fields{
		"to":      phoneNumber,
		"message": message,
	})

	// Ensure phone number is in the correct format (no "+" prefix)
	phoneNumber = strings.TrimPrefix(phoneNumber, "+")

	// Prepare the message payload
	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                phoneNumber,
		"type":              "text",
		"text": map[string]string{
			"body": message,
		},
	}

	return s.sendMessage(ctx, payload)
}

// SendOptionsMessage sends a message with options for the user to select
func (s *whatsAppService) SendOptionsMessage(ctx context.Context, phoneNumber, message string, options []string) error {
	s.log.Info("Sending options message", logger.Fields{
		"to":      phoneNumber,
		"message": message,
		"options": options,
	})

	// For now, we'll just format the options as bullet points in the text message
	// In a production environment, you might want to use WhatsApp's interactive messages
	fullMessage := fmt.Sprintf("%s\n\n%s", message, strings.Join(options, "\n"))
	
	return s.SendTextMessage(ctx, phoneNumber, fullMessage)
}

// VerifyWebhook verifies the WhatsApp webhook signature
func (s *whatsAppService) VerifyWebhook(mode, token, challenge string) (bool, string) {
	s.log.Info("Verifying webhook", logger.Fields{
		"mode":      mode,
		"token":     token,
		"challenge": challenge,
	})

	if mode != "subscribe" || token != s.config.WebhookToken {
		return false, ""
	}

	return true, challenge
}

// sendMessage handles the actual sending of the message to the WhatsApp API
func (s *whatsAppService) sendMessage(ctx context.Context, payload map[string]interface{}) error {
	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal message payload: %w", err)
	}

	// Prepare the request
	url := fmt.Sprintf("%s/v%s/%s/messages", 
		s.config.BaseURL, 
		s.config.APIVersion, 
		s.config.PhoneNumberID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.config.AccessToken))

	// Send the request
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode >= 400 {
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return fmt.Errorf("failed with status code %d", resp.StatusCode)
		}
		return fmt.Errorf("API request failed with status code %d: %v", resp.StatusCode, errorResponse)
	}

	s.log.Info("Message sent successfully", logger.Fields{
		"status_code": resp.StatusCode,
	})

	return nil
} 