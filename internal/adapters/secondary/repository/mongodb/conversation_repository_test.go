package mongodb

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/2rprbm/conta-med-backend/internal/core/domain"
)

func TestConversationRepository_Save(t *testing.T) {
	t.Run("should save conversation successfully when valid conversation is provided", func(t *testing.T) {
		// arrange
		db, cleanup := setupTestDB(t)
		defer cleanup()

		repo := NewConversationRepository(db, &mockLogger{})
		conversation := &domain.Conversation{
			PhoneNumber:    "+5511999999999",
			Status:         domain.Active,
			State:          domain.Initial,
			StartedAt:      time.Now(),
			LastUpdatedAt:  time.Now(),
			UserSelections: make(map[string]string),
		}

		// act
		err := repo.Save(context.Background(), conversation)

		// assert
		assert.NoError(t, err)
		assert.NotEmpty(t, conversation.ID)
	})

	t.Run("should generate ID when conversation ID is empty", func(t *testing.T) {
		// arrange
		db, cleanup := setupTestDB(t)
		defer cleanup()

		repo := NewConversationRepository(db, &mockLogger{})
		conversation := domain.NewConversation("+5511999999999")

		// act
		err := repo.Save(context.Background(), conversation)

		// assert
		assert.NoError(t, err)
		assert.NotEmpty(t, conversation.ID)
	})
}

func TestConversationRepository_FindByID(t *testing.T) {
	t.Run("should return conversation when ID exists", func(t *testing.T) {
		// arrange
		db, cleanup := setupTestDB(t)
		defer cleanup()

		repo := NewConversationRepository(db, &mockLogger{})
		conversation := domain.NewConversation("+5511999999999")
		
		err := repo.Save(context.Background(), conversation)
		require.NoError(t, err)

		// act
		found, err := repo.FindByID(context.Background(), conversation.ID)

		// assert
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, conversation.ID, found.ID)
		assert.Equal(t, conversation.PhoneNumber, found.PhoneNumber)
	})

	t.Run("should return nil when conversation ID does not exist", func(t *testing.T) {
		// arrange
		db, cleanup := setupTestDB(t)
		defer cleanup()

		repo := NewConversationRepository(db, &mockLogger{})

		// act
		found, err := repo.FindByID(context.Background(), "nonexistent")

		// assert
		assert.NoError(t, err)
		assert.Nil(t, found)
	})
}

func TestConversationRepository_FindByPhoneNumber(t *testing.T) {
	t.Run("should return conversation when phone number exists", func(t *testing.T) {
		// arrange
		db, cleanup := setupTestDB(t)
		defer cleanup()

		repo := NewConversationRepository(db, &mockLogger{})
		phoneNumber := "+5511999999999"
		conversation := domain.NewConversation(phoneNumber)
		
		err := repo.Save(context.Background(), conversation)
		require.NoError(t, err)

		// act
		found, err := repo.FindByPhoneNumber(context.Background(), phoneNumber)

		// assert
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, phoneNumber, found.PhoneNumber)
	})

	t.Run("should return nil when phone number does not exist", func(t *testing.T) {
		// arrange
		db, cleanup := setupTestDB(t)
		defer cleanup()

		repo := NewConversationRepository(db, &mockLogger{})

		// act
		found, err := repo.FindByPhoneNumber(context.Background(), "+5511888888888")

		// assert
		assert.NoError(t, err)
		assert.Nil(t, found)
	})
}

func TestConversationRepository_FindActiveByPhoneNumber(t *testing.T) {
	t.Run("should return active conversation when active conversation exists for phone number", func(t *testing.T) {
		// arrange
		db, cleanup := setupTestDB(t)
		defer cleanup()

		repo := NewConversationRepository(db, &mockLogger{})
		phoneNumber := "+5511999999999"
		conversation := domain.NewConversation(phoneNumber)
		
		err := repo.Save(context.Background(), conversation)
		require.NoError(t, err)

		// act
		found, err := repo.FindActiveByPhoneNumber(context.Background(), phoneNumber)

		// assert
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, phoneNumber, found.PhoneNumber)
		assert.Equal(t, domain.Active, found.Status)
	})

	t.Run("should return nil when no active conversation exists for phone number", func(t *testing.T) {
		// arrange
		db, cleanup := setupTestDB(t)
		defer cleanup()

		repo := NewConversationRepository(db, &mockLogger{})
		phoneNumber := "+5511999999999"
		conversation := domain.NewConversation(phoneNumber)
		conversation.Complete() // Mark as completed
		
		err := repo.Save(context.Background(), conversation)
		require.NoError(t, err)

		// act
		found, err := repo.FindActiveByPhoneNumber(context.Background(), phoneNumber)

		// assert
		assert.NoError(t, err)
		assert.Nil(t, found)
	})

	t.Run("should return nil when phone number does not exist", func(t *testing.T) {
		// arrange
		db, cleanup := setupTestDB(t)
		defer cleanup()

		repo := NewConversationRepository(db, &mockLogger{})

		// act
		found, err := repo.FindActiveByPhoneNumber(context.Background(), "+5511888888888")

		// assert
		assert.NoError(t, err)
		assert.Nil(t, found)
	})
}

func TestConversationRepository_Update(t *testing.T) {
	t.Run("should update conversation successfully when conversation exists", func(t *testing.T) {
		// arrange
		db, cleanup := setupTestDB(t)
		defer cleanup()

		repo := NewConversationRepository(db, &mockLogger{})
		conversation := domain.NewConversation("+5511999999999")
		
		err := repo.Save(context.Background(), conversation)
		require.NoError(t, err)

		// Modify conversation
		conversation.UpdateState(domain.MainMenu)
		conversation.AddUserSelection("main_menu", "1")

		// act
		err = repo.Update(context.Background(), conversation)

		// assert
		assert.NoError(t, err)

		// Verify the update
		updated, err := repo.FindByID(context.Background(), conversation.ID)
		require.NoError(t, err)
		assert.Equal(t, domain.MainMenu, updated.State)
		assert.Equal(t, "1", updated.UserSelections["main_menu"])
	})

	t.Run("should return error when conversation does not exist", func(t *testing.T) {
		// arrange
		db, cleanup := setupTestDB(t)
		defer cleanup()

		repo := NewConversationRepository(db, &mockLogger{})
		conversation := &domain.Conversation{
			ID:          "nonexistent",
			PhoneNumber: "+5511999999999",
			Status:      domain.Active,
			State:       domain.MainMenu,
		}

		// act
		err := repo.Update(context.Background(), conversation)

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "conversation not found for update")
	})
} 