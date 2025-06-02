package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Contabilidade-Medicos/conta-med-backend/internal/core/domain"
	"github.com/Contabilidade-Medicos/conta-med-backend/internal/core/ports"
	"github.com/Contabilidade-Medicos/conta-med-backend/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementations
type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) Save(ctx context.Context, message *domain.Message) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *MockMessageRepository) FindByConversationID(ctx context.Context, conversationID string) ([]*domain.Message, error) {
	args := m.Called(ctx, conversationID)
	return args.Get(0).([]*domain.Message), args.Error(1)
}

func (m *MockMessageRepository) FindLatestByConversationID(ctx context.Context, conversationID string, limit int) ([]*domain.Message, error) {
	args := m.Called(ctx, conversationID, limit)
	return args.Get(0).([]*domain.Message), args.Error(1)
}

type MockConversationRepository struct {
	mock.Mock
}

func (m *MockConversationRepository) Save(ctx context.Context, conversation *domain.Conversation) error {
	args := m.Called(ctx, conversation)
	return args.Error(0)
}

func (m *MockConversationRepository) FindByID(ctx context.Context, id string) (*domain.Conversation, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Conversation), args.Error(1)
}

func (m *MockConversationRepository) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.Conversation, error) {
	args := m.Called(ctx, phoneNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Conversation), args.Error(1)
}

func (m *MockConversationRepository) FindActiveByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.Conversation, error) {
	args := m.Called(ctx, phoneNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Conversation), args.Error(1)
}

func (m *MockConversationRepository) Update(ctx context.Context, conversation *domain.Conversation) error {
	args := m.Called(ctx, conversation)
	return args.Error(0)
}

type MockWhatsAppService struct {
	mock.Mock
}

func (m *MockWhatsAppService) SendTextMessage(ctx context.Context, phoneNumber, message string) error {
	args := m.Called(ctx, phoneNumber, message)
	return args.Error(0)
}

func (m *MockWhatsAppService) SendOptionsMessage(ctx context.Context, phoneNumber, message string, options []string) error {
	args := m.Called(ctx, phoneNumber, message, options)
	return args.Error(0)
}

func (m *MockWhatsAppService) VerifyWebhook(mode, token, challenge string) (bool, string) {
	args := m.Called(mode, token, challenge)
	return args.Bool(0), args.String(1)
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(msg string, fields ...logger.Fields) {
	if len(fields) > 0 {
		m.Called(msg, fields[0])
	} else {
		m.Called(msg)
	}
}

func (m *MockLogger) Error(msg string, fields ...logger.Fields) {
	if len(fields) > 0 {
		m.Called(msg, fields[0])
	} else {
		m.Called(msg)
	}
}

func (m *MockLogger) Debug(msg string, fields ...logger.Fields) {
	if len(fields) > 0 {
		m.Called(msg, fields[0])
	} else {
		m.Called(msg)
	}
}

func (m *MockLogger) Warn(msg string, fields ...logger.Fields) {
	if len(fields) > 0 {
		m.Called(msg, fields[0])
	} else {
		m.Called(msg)
	}
}

func (m *MockLogger) Fatal(msg string, fields ...logger.Fields) {
	if len(fields) > 0 {
		m.Called(msg, fields[0])
	} else {
		m.Called(msg)
	}
}

func (m *MockLogger) With(fields logger.Fields) logger.Logger {
	args := m.Called(fields)
	return args.Get(0).(logger.Logger)
}

// Tests
func TestGetOrCreateConversation_ExistingConversation(t *testing.T) {
	// arrange
	mockMessageRepo := new(MockMessageRepository)
	mockConversationRepo := new(MockConversationRepository)
	mockWhatsAppService := new(MockWhatsAppService)
	mockLogger := new(MockLogger)

	phoneNumber := "+5511999999999"
	existingConversation := &domain.Conversation{
		ID:          "existing-conv-id",
		PhoneNumber: phoneNumber,
		Status:      domain.Active,
		State:       domain.MainMenu,
	}

	ctx := context.Background()
	mockConversationRepo.On("FindActiveByPhoneNumber", ctx, phoneNumber).Return(existingConversation, nil)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()

	service := NewChatbotService(mockMessageRepo, mockConversationRepo, mockWhatsAppService, mockLogger)

	// act
	result, err := service.GetOrCreateConversation(ctx, phoneNumber)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, existingConversation, result)
	mockConversationRepo.AssertExpectations(t)
}

func TestGetOrCreateConversation_NewConversation(t *testing.T) {
	// arrange
	mockMessageRepo := new(MockMessageRepository)
	mockConversationRepo := new(MockConversationRepository)
	mockWhatsAppService := new(MockWhatsAppService)
	mockLogger := new(MockLogger)

	phoneNumber := "+5511999999999"
	ctx := context.Background()

	mockConversationRepo.On("FindActiveByPhoneNumber", ctx, phoneNumber).Return(nil, errors.New("not found"))
	mockConversationRepo.On("Save", ctx, mock.AnythingOfType("*domain.Conversation")).Return(nil)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()

	service := NewChatbotService(mockMessageRepo, mockConversationRepo, mockWhatsAppService, mockLogger)

	// act
	result, err := service.GetOrCreateConversation(ctx, phoneNumber)

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, phoneNumber, result.PhoneNumber)
	assert.Equal(t, domain.Active, result.Status)
	assert.Equal(t, domain.Initial, result.State)
	mockConversationRepo.AssertExpectations(t)
}

func TestHandleIncomingMessage_NewConversation(t *testing.T) {
	// arrange
	mockMessageRepo := new(MockMessageRepository)
	mockConversationRepo := new(MockConversationRepository)
	mockWhatsAppService := new(MockWhatsAppService)
	mockLogger := new(MockLogger)

	phoneNumber := "+5511999999999"
	message := "Olá"
	ctx := context.Background()

	conversation := &domain.Conversation{
		ID:          "conv-id",
		PhoneNumber: phoneNumber,
		Status:      domain.Active,
		State:       domain.Initial,
	}

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockConversationRepo.On("FindActiveByPhoneNumber", ctx, phoneNumber).Return(conversation, nil)
	mockMessageRepo.On("Save", ctx, mock.AnythingOfType("*domain.Message")).Return(nil)
	mockConversationRepo.On("Update", ctx, mock.AnythingOfType("*domain.Conversation")).Return(nil)
	mockWhatsAppService.On("SendOptionsMessage", ctx, phoneNumber, mock.AnythingOfType("string"), mock.AnythingOfType("[]string")).Return(nil)

	service := NewChatbotService(mockMessageRepo, mockConversationRepo, mockWhatsAppService, mockLogger)

	// act
	err := service.HandleIncomingMessage(ctx, phoneNumber, message)

	// assert
	assert.NoError(t, err)
	mockMessageRepo.AssertExpectations(t)
	mockConversationRepo.AssertExpectations(t)
	mockWhatsAppService.AssertExpectations(t)
}

func TestProcessUserSelection_MainMenu_Option2(t *testing.T) {
	// arrange
	mockMessageRepo := new(MockMessageRepository)
	mockConversationRepo := new(MockConversationRepository)
	mockWhatsAppService := new(MockWhatsAppService)
	mockLogger := new(MockLogger)

	phoneNumber := "+5511999999999"
	message := "2"
	ctx := context.Background()

	conversation := &domain.Conversation{
		ID:          "conv-id",
		PhoneNumber: phoneNumber,
		Status:      domain.Active,
		State:       domain.MainMenu,
	}

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockConversationRepo.On("FindActiveByPhoneNumber", ctx, phoneNumber).Return(conversation, nil)
	mockConversationRepo.On("Update", ctx, mock.AnythingOfType("*domain.Conversation")).Return(nil)
	mockWhatsAppService.On("SendOptionsMessage", ctx, phoneNumber, mock.AnythingOfType("string"), mock.AnythingOfType("[]string")).Return(nil)
	mockMessageRepo.On("Save", ctx, mock.AnythingOfType("*domain.Message")).Return(nil)

	service := NewChatbotService(mockMessageRepo, mockConversationRepo, mockWhatsAppService, mockLogger)

	// act
	err := service.ProcessUserSelection(ctx, phoneNumber, message)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, domain.CompanyTypeSelection, conversation.State)
	mockConversationRepo.AssertExpectations(t)
	mockWhatsAppService.AssertExpectations(t)
	mockMessageRepo.AssertExpectations(t)
}

func TestProcessUserSelection_MainMenu_InvalidOption(t *testing.T) {
	// arrange
	mockMessageRepo := new(MockMessageRepository)
	mockConversationRepo := new(MockConversationRepository)
	mockWhatsAppService := new(MockWhatsAppService)
	mockLogger := new(MockLogger)

	phoneNumber := "+5511999999999"
	message := "invalid"
	ctx := context.Background()

	conversation := &domain.Conversation{
		ID:          "conv-id",
		PhoneNumber: phoneNumber,
		Status:      domain.Active,
		State:       domain.MainMenu,
	}

	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockConversationRepo.On("FindActiveByPhoneNumber", ctx, phoneNumber).Return(conversation, nil)
	mockWhatsAppService.On("SendTextMessage", ctx, phoneNumber, mock.AnythingOfType("string")).Return(nil)
	mockMessageRepo.On("Save", ctx, mock.AnythingOfType("*domain.Message")).Return(nil)

	service := NewChatbotService(mockMessageRepo, mockConversationRepo, mockWhatsAppService, mockLogger)

	// act
	err := service.ProcessUserSelection(ctx, phoneNumber, message)

	// assert
	assert.NoError(t, err)
	mockWhatsAppService.AssertExpectations(t)
	mockMessageRepo.AssertExpectations(t)
}

func TestGetGreetingByTime(t *testing.T) {
	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "should return morning greeting for 8 AM",
			time:     time.Date(2023, 1, 1, 8, 0, 0, 0, time.UTC),
			expected: "Bom dia",
		},
		{
			name:     "should return afternoon greeting for 2 PM",
			time:     time.Date(2023, 1, 1, 14, 0, 0, 0, time.UTC),
			expected: "Boa tarde",
		},
		{
			name:     "should return evening greeting for 9 PM",
			time:     time.Date(2023, 1, 1, 21, 0, 0, 0, time.UTC),
			expected: "Boa noite",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// act
			result := getGreetingByTime(tc.time)

			// assert
			assert.Equal(t, tc.expected, result)
		})
	}
} 