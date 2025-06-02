package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/2rprbm/conta-med-backend/config"
	"github.com/2rprbm/conta-med-backend/pkg/logger"
)

// Client represents a WhatsApp API client
type Client struct {
	Config     *config.Config
	Logger     logger.Logger
	HttpClient *http.Client
	APIURL     string // URL template for API endpoints
}

// NewClient creates a new WhatsApp API client
func NewClient(cfg *config.Config, log logger.Logger) *Client {
	return &Client{
		Config: cfg,
		Logger: log,
		HttpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		APIURL: "https://graph.facebook.com/v18.0/%s/messages",
	}
}

// TextMessage represents a text message to send
type TextMessage struct {
	MessagingProduct string `json:"messaging_product"`
	RecipientType    string `json:"recipient_type"`
	To               string `json:"to"`
	Type             string `json:"type"`
	Text             struct {
		Body string `json:"body"`
	} `json:"text"`
}

// SendTextMessage sends a text message to a WhatsApp user
func (c *Client) SendTextMessage(to, message string) error {
	if c.Config.WhatsApp.PhoneNumberID == "" {
		return fmt.Errorf("phone number ID not configured")
	}

	// Create message payload
	payload := TextMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               to,
		Type:             "text",
	}
	payload.Text.Body = message

	// Convert to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling message: %w", err)
	}

	// Create API URL
	url := fmt.Sprintf(c.APIURL, c.Config.WhatsApp.PhoneNumberID)

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Config.WhatsApp.AccessToken)

	// Send request
	c.Logger.Debug("Sending WhatsApp message", logger.Fields{
		"to":      to,
		"message": message,
	})
	
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return fmt.Errorf("API error: %d", resp.StatusCode)
		}
		return fmt.Errorf("API error: %v", errorResp)
	}

	c.Logger.Info("Message sent successfully", logger.Fields{
		"to": to,
	})
	return nil
}
