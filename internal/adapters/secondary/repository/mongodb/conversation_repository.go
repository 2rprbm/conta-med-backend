package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/2rprbm/conta-med-backend/internal/core/domain"
	"github.com/2rprbm/conta-med-backend/internal/core/ports"
	"github.com/2rprbm/conta-med-backend/pkg/logger"
)

type conversationRepository struct {
	collection *mongo.Collection
	log        logger.Logger
}

// NewConversationRepository creates a new MongoDB conversation repository
func NewConversationRepository(db *mongo.Database, log logger.Logger) ports.ConversationRepository {
	return &conversationRepository{
		collection: db.Collection("conversations"),
		log:        log,
	}
}

// Save persists a conversation to MongoDB
func (r *conversationRepository) Save(ctx context.Context, conversation *domain.Conversation) error {
	r.log.Info("Saving conversation to MongoDB", logger.Fields{
		"phone_number": conversation.PhoneNumber,
		"state":        conversation.State,
		"status":       conversation.Status,
	})

	// If ID is empty, generate a new ObjectID
	if conversation.ID == "" {
		conversation.ID = primitive.NewObjectID().Hex()
	}

	_, err := r.collection.InsertOne(ctx, conversation)
	if err != nil {
		r.log.Error("Failed to save conversation", logger.Fields{
			"error":        err.Error(),
			"phone_number": conversation.PhoneNumber,
		})
		return fmt.Errorf("failed to save conversation: %w", err)
	}

	r.log.Info("Conversation successfully saved", logger.Fields{
		"conversation_id": conversation.ID,
		"phone_number":   conversation.PhoneNumber,
	})

	return nil
}

// FindByID retrieves a conversation by ID
func (r *conversationRepository) FindByID(ctx context.Context, id string) (*domain.Conversation, error) {
	r.log.Info("Finding conversation by ID", logger.Fields{
		"conversation_id": id,
	})

	var conversation domain.Conversation
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&conversation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.log.Info("Conversation not found", logger.Fields{
				"conversation_id": id,
			})
			return nil, nil
		}

		r.log.Error("Failed to find conversation by ID", logger.Fields{
			"error":          err.Error(),
			"conversation_id": id,
		})
		return nil, fmt.Errorf("failed to find conversation by ID: %w", err)
	}

	r.log.Info("Conversation found by ID", logger.Fields{
		"conversation_id": id,
		"phone_number":   conversation.PhoneNumber,
	})

	return &conversation, nil
}

// FindByPhoneNumber retrieves a conversation by phone number
func (r *conversationRepository) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.Conversation, error) {
	r.log.Info("Finding conversation by phone number", logger.Fields{
		"phone_number": phoneNumber,
	})

	var conversation domain.Conversation
	err := r.collection.FindOne(ctx, bson.M{"phone_number": phoneNumber}).Decode(&conversation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.log.Info("Conversation not found for phone number", logger.Fields{
				"phone_number": phoneNumber,
			})
			return nil, nil
		}

		r.log.Error("Failed to find conversation by phone number", logger.Fields{
			"error":        err.Error(),
			"phone_number": phoneNumber,
		})
		return nil, fmt.Errorf("failed to find conversation by phone number: %w", err)
	}

	r.log.Info("Conversation found by phone number", logger.Fields{
		"conversation_id": conversation.ID,
		"phone_number":   phoneNumber,
	})

	return &conversation, nil
}

// FindActiveByPhoneNumber retrieves an active conversation by phone number
func (r *conversationRepository) FindActiveByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.Conversation, error) {
	r.log.Info("Finding active conversation by phone number", logger.Fields{
		"phone_number": phoneNumber,
	})

	filter := bson.M{
		"phone_number": phoneNumber,
		"status":       domain.Active,
	}

	var conversation domain.Conversation
	err := r.collection.FindOne(ctx, filter).Decode(&conversation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.log.Info("Active conversation not found for phone number", logger.Fields{
				"phone_number": phoneNumber,
			})
			return nil, nil
		}

		r.log.Error("Failed to find active conversation by phone number", logger.Fields{
			"error":        err.Error(),
			"phone_number": phoneNumber,
		})
		return nil, fmt.Errorf("failed to find active conversation by phone number: %w", err)
	}

	r.log.Info("Active conversation found by phone number", logger.Fields{
		"conversation_id": conversation.ID,
		"phone_number":   phoneNumber,
		"state":          conversation.State,
	})

	return &conversation, nil
}

// Update updates an existing conversation
func (r *conversationRepository) Update(ctx context.Context, conversation *domain.Conversation) error {
	r.log.Info("Updating conversation", logger.Fields{
		"conversation_id": conversation.ID,
		"phone_number":   conversation.PhoneNumber,
		"state":          conversation.State,
		"status":         conversation.Status,
	})

	filter := bson.M{"_id": conversation.ID}
	update := bson.M{"$set": conversation}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		r.log.Error("Failed to update conversation", logger.Fields{
			"error":          err.Error(),
			"conversation_id": conversation.ID,
		})
		return fmt.Errorf("failed to update conversation: %w", err)
	}

	if result.MatchedCount == 0 {
		r.log.Error("Conversation not found for update", logger.Fields{
			"conversation_id": conversation.ID,
		})
		return fmt.Errorf("conversation not found for update")
	}

	r.log.Info("Conversation successfully updated", logger.Fields{
		"conversation_id": conversation.ID,
		"modified_count": result.ModifiedCount,
	})

	return nil
} 