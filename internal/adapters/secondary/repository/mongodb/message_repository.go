package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/2rprbm/conta-med-backend/internal/core/domain"
	"github.com/2rprbm/conta-med-backend/internal/core/ports"
	"github.com/2rprbm/conta-med-backend/pkg/logger"
)

type messageRepository struct {
	collection *mongo.Collection
	log        logger.Logger
}

// NewMessageRepository creates a new MongoDB message repository
func NewMessageRepository(db *mongo.Database, log logger.Logger) ports.MessageRepository {
	return &messageRepository{
		collection: db.Collection("messages"),
		log:        log,
	}
}

// Save persists a message to MongoDB
func (r *messageRepository) Save(ctx context.Context, message *domain.Message) error {
	r.log.Info("Saving message to MongoDB", logger.Fields{
		"conversation_id": message.ConversationID,
		"phone_number":   message.PhoneNumber,
		"direction":      message.Direction,
	})

	// If ID is empty, generate a new ObjectID
	if message.ID == "" {
		message.ID = primitive.NewObjectID().Hex()
	}

	_, err := r.collection.InsertOne(ctx, message)
	if err != nil {
		r.log.Error("Failed to save message", logger.Fields{
			"error":          err.Error(),
			"conversation_id": message.ConversationID,
		})
		return fmt.Errorf("failed to save message: %w", err)
	}

	r.log.Info("Message successfully saved", logger.Fields{
		"message_id":     message.ID,
		"conversation_id": message.ConversationID,
	})

	return nil
}

// FindByConversationID retrieves all messages for a given conversation
func (r *messageRepository) FindByConversationID(ctx context.Context, conversationID string) ([]*domain.Message, error) {
	r.log.Info("Finding messages by conversation ID", logger.Fields{
		"conversation_id": conversationID,
	})

	filter := bson.M{"conversation_id": conversationID}
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}}) // Sort by timestamp ascending

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		r.log.Error("Failed to find messages", logger.Fields{
			"error":          err.Error(),
			"conversation_id": conversationID,
		})
		return nil, fmt.Errorf("failed to find messages: %w", err)
	}
	defer cursor.Close(ctx)

	var messages []*domain.Message
	if err := cursor.All(ctx, &messages); err != nil {
		r.log.Error("Failed to decode messages", logger.Fields{
			"error":          err.Error(),
			"conversation_id": conversationID,
		})
		return nil, fmt.Errorf("failed to decode messages: %w", err)
	}

	r.log.Info("Messages found", logger.Fields{
		"conversation_id": conversationID,
		"count":          len(messages),
	})

	return messages, nil
}

// FindLatestByConversationID retrieves the latest N messages for a conversation
func (r *messageRepository) FindLatestByConversationID(ctx context.Context, conversationID string, limit int) ([]*domain.Message, error) {
	r.log.Info("Finding latest messages by conversation ID", logger.Fields{
		"conversation_id": conversationID,
		"limit":          limit,
	})

	filter := bson.M{"conversation_id": conversationID}
	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}). // Sort by timestamp descending
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		r.log.Error("Failed to find latest messages", logger.Fields{
			"error":          err.Error(),
			"conversation_id": conversationID,
		})
		return nil, fmt.Errorf("failed to find latest messages: %w", err)
	}
	defer cursor.Close(ctx)

	var messages []*domain.Message
	if err := cursor.All(ctx, &messages); err != nil {
		r.log.Error("Failed to decode latest messages", logger.Fields{
			"error":          err.Error(),
			"conversation_id": conversationID,
		})
		return nil, fmt.Errorf("failed to decode latest messages: %w", err)
	}

	r.log.Info("Latest messages found", logger.Fields{
		"conversation_id": conversationID,
		"count":          len(messages),
	})

	return messages, nil
} 