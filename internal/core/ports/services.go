package ports

import (
	"context"

	"github.com/2rprbm/conta-med-backend/internal/core/domain"
)

// WhatsAppService defines the interface for WhatsApp interactions
type WhatsAppService interface {
	// SendTextMessage sends a text message to a WhatsApp number
	SendTextMessage(ctx context.Context, phoneNumber, message string) error
	
	// SendOptionsMessage sends a message with options for the user to select
	SendOptionsMessage(ctx context.Context, phoneNumber, message string, options []string) error
	
	// VerifyWebhook verifies the WhatsApp webhook signature
	VerifyWebhook(mode, token, challenge string) (bool, string)
}

// ChatbotService defines the interface for chatbot functionality
type ChatbotService interface {
	// HandleIncomingMessage processes an incoming message and generates a response
	HandleIncomingMessage(ctx context.Context, phoneNumber, message string) error
	
	// GetOrCreateConversation retrieves an existing conversation or creates a new one
	GetOrCreateConversation(ctx context.Context, phoneNumber string) (*domain.Conversation, error)
	
	// SendWelcomeMessage sends the initial welcome message to a user
	SendWelcomeMessage(ctx context.Context, phoneNumber string) error
	
	// ProcessUserSelection processes a user's menu selection
	ProcessUserSelection(ctx context.Context, phoneNumber, message string) error
} 