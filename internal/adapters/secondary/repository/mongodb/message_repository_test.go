package mongodb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/2rprbm/conta-med-backend/internal/core/domain"
)

func TestMessageRepository_Save(t *testing.T) {
	t.Run("should save message successfully when valid message is provided", func(t *testing.T) {
		// arrange
		db, cleanup := setupTestDB(t)
		defer cleanup()

		repo := NewMessageRepository(db, &mockLogger{})
		message := &domain.Message{
			ConversationID: "conv123",
			PhoneNumber:   "+5511999999999",
			Content:       "Hello World",
			Type:          domain.TextMessage,
			Direction:     domain.Inbound,
			Timestamp:     time.Now(),
		}

		// act
		err := repo.Save(context.Background(), message)

		// assert
		assert.NoError(t, err)
		assert.NotEmpty(t, message.ID)
	})

	t.Run("should generate ID when message ID is empty", func(t *testing.T) {
		// arrange
		db, cleanup := setupTestDB(t)
		defer cleanup()

		repo := NewMessageRepository(db, &mockLogger{})
		message := &domain.Message{
			ConversationID: "conv123",
			PhoneNumber:   "+5511999999999",
			Content:       "Hello World",
			Type:          domain.TextMessage,
			Direction:     domain.Inbound,
			Timestamp:     time.Now(),
		}

		// act
		err := repo.Save(context.Background(), message)

		// assert
		assert.NoError(t, err)
		assert.NotEmpty(t, message.ID)
	})
}

func TestMessageRepository_FindByConversationID(t *testing.T) {
	t.Run("should return messages when conversation has messages", func(t *testing.T) {
		// arrange
		db, cleanup := setupTestDB(t)
		defer cleanup()

		repo := NewMessageRepository(db, &mockLogger{})
		conversationID := "conv123"

		// Save test messages
		message1 := &domain.Message{
			ConversationID: conversationID,
			PhoneNumber:   "+5511999999999",
			Content:       "First message",
			Type:          domain.TextMessage,
			Direction:     domain.Inbound,
			Timestamp:     time.Now().Add(-2 * time.Minute),
		}
		message2 := &domain.Message{
			ConversationID: conversationID,
			PhoneNumber:   "+5511999999999",
			Content:       "Second message",
			Type:          domain.TextMessage,
			Direction:     domain.Outbound,
			Timestamp:     time.Now().Add(-1 * time.Minute),
		}

		err := repo.Save(context.Background(), message1)
		require.NoError(t, err)
		err = repo.Save(context.Background(), message2)
		require.NoError(t, err)

		// act
		messages, err := repo.FindByConversationID(context.Background(), conversationID)

		// assert
		assert.NoError(t, err)
		assert.Len(t, messages, 2)
		assert.Equal(t, "First message", messages[0].Content)
		assert.Equal(t, "Second message", messages[1].Content)
	})

	t.Run("should return empty slice when no messages exist for conversation", func(t *testing.T) {
		// arrange
		db, cleanup := setupTestDB(t)
		defer cleanup()

		repo := NewMessageRepository(db, &mockLogger{})

		// act
		messages, err := repo.FindByConversationID(context.Background(), "nonexistent")

		// assert
		assert.NoError(t, err)
		assert.Empty(t, messages)
	})
}

func TestMessageRepository_FindLatestByConversationID(t *testing.T) {
	t.Run("should return latest messages with limit when multiple messages exist", func(t *testing.T) {
		// arrange
		db, cleanup := setupTestDB(t)
		defer cleanup()

		repo := NewMessageRepository(db, &mockLogger{})
		conversationID := "conv123"

		// Save test messages with different timestamps
		for i := 0; i < 5; i++ {
			message := &domain.Message{
				ConversationID: conversationID,
				PhoneNumber:   "+5511999999999",
				Content:       fmt.Sprintf("Message %d", i),
				Type:          domain.TextMessage,
				Direction:     domain.Inbound,
				Timestamp:     time.Now().Add(time.Duration(i) * time.Minute),
			}
			err := repo.Save(context.Background(), message)
			require.NoError(t, err)
		}

		// act
		messages, err := repo.FindLatestByConversationID(context.Background(), conversationID, 3)

		// assert
		assert.NoError(t, err)
		assert.Len(t, messages, 3)
		// Should return messages in descending order (latest first)
		assert.Equal(t, "Message 4", messages[0].Content)
		assert.Equal(t, "Message 3", messages[1].Content)
		assert.Equal(t, "Message 2", messages[2].Content)
	})

	t.Run("should return empty slice when no messages exist for conversation", func(t *testing.T) {
		// arrange
		db, cleanup := setupTestDB(t)
		defer cleanup()

		repo := NewMessageRepository(db, &mockLogger{})

		// act
		messages, err := repo.FindLatestByConversationID(context.Background(), "nonexistent", 5)

		// assert
		assert.NoError(t, err)
		assert.Empty(t, messages)
	})
} 