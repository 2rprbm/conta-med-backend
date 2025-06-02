package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Contabilidade-Medicos/conta-med-backend/config"
	"github.com/Contabilidade-Medicos/conta-med-backend/internal/adapters/primary/api"
	httpserver "github.com/Contabilidade-Medicos/conta-med-backend/internal/adapters/primary/http"
	"github.com/Contabilidade-Medicos/conta-med-backend/internal/adapters/secondary/whatsapp"
	"github.com/Contabilidade-Medicos/conta-med-backend/internal/core/domain"
	"github.com/Contabilidade-Medicos/conta-med-backend/internal/core/services"
	"github.com/Contabilidade-Medicos/conta-med-backend/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log := logger.New(cfg.Logging.Level)
	log.Info("ContaMed WhatsApp Chatbot - Starting server...")

	// Initialize WhatsApp service
	whatsappConfig := whatsapp.WhatsAppConfig{
		APIVersion:    cfg.WhatsApp.APIVersion,
		PhoneNumberID: cfg.WhatsApp.PhoneNumberID,
		AccessToken:   cfg.WhatsApp.AccessToken,
		WebhookToken:  cfg.WhatsApp.WebhookVerifyToken,
		BaseURL:       cfg.WhatsApp.BaseURL,
	}
	whatsAppService := whatsapp.NewWhatsAppService(whatsappConfig, log)

	// Initialize repositories (using mocks for now)
	messageRepo := newMockMessageRepository()
	conversationRepo := newMockConversationRepository()

	// Initialize chatbot service
	chatbotService := services.NewChatbotService(messageRepo, conversationRepo, whatsAppService, log)

	// Initialize webhook handler
	webhookHandler := api.NewWebhookHandler(chatbotService, whatsAppService, log)

	// Initialize HTTP server
	server := httpserver.NewServer(cfg, log, webhookHandler)

	// Start server
	server.Start()

	// Wait for termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownPeriod)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server shutdown error: %v", err)
	}

	log.Info("Server stopped")
}

// Mock repositories for development (to be replaced with real implementations in a future sprint)

// Mock message repository
type mockMessageRepository struct{}

func newMockMessageRepository() *mockMessageRepository {
	return &mockMessageRepository{}
}

func (m *mockMessageRepository) Save(ctx context.Context, message *domain.Message) error {
	return nil
}

func (m *mockMessageRepository) FindByConversationID(ctx context.Context, conversationID string) ([]*domain.Message, error) {
	return []*domain.Message{}, nil
}

func (m *mockMessageRepository) FindLatestByConversationID(ctx context.Context, conversationID string, limit int) ([]*domain.Message, error) {
	return []*domain.Message{}, nil
}

// Mock conversation repository
type mockConversationRepository struct{}

func newMockConversationRepository() *mockConversationRepository {
	return &mockConversationRepository{}
}

func (m *mockConversationRepository) Save(ctx context.Context, conversation *domain.Conversation) error {
	return nil
}

func (m *mockConversationRepository) FindByID(ctx context.Context, id string) (*domain.Conversation, error) {
	return nil, nil
}

func (m *mockConversationRepository) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.Conversation, error) {
	return nil, nil
}

func (m *mockConversationRepository) FindActiveByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.Conversation, error) {
	return nil, nil
}

func (m *mockConversationRepository) Update(ctx context.Context, conversation *domain.Conversation) error {
	return nil
}
