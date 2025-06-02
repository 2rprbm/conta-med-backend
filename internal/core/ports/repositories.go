package ports

import (
	"context"

	"github.com/2rprbm/conta-med-backend/internal/core/domain"
)

// MessageRepository defines the interface for message storage operations
type MessageRepository interface {
	// Save persists a message to storage
	Save(ctx context.Context, message *domain.Message) error
	
	// FindByConversationID retrieves all messages for a given conversation
	FindByConversationID(ctx context.Context, conversationID string) ([]*domain.Message, error)
	
	// FindLatestByConversationID retrieves the latest N messages for a conversation
	FindLatestByConversationID(ctx context.Context, conversationID string, limit int) ([]*domain.Message, error)
}

// ConversationRepository defines the interface for conversation storage operations
type ConversationRepository interface {
	// Save persists a conversation to storage
	Save(ctx context.Context, conversation *domain.Conversation) error
	
	// FindByID retrieves a conversation by ID
	FindByID(ctx context.Context, id string) (*domain.Conversation, error)
	
	// FindByPhoneNumber retrieves a conversation by phone number
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.Conversation, error)
	
	// FindActiveByPhoneNumber retrieves an active conversation by phone number
	FindActiveByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.Conversation, error)
	
	// Update updates an existing conversation
	Update(ctx context.Context, conversation *domain.Conversation) error
} 