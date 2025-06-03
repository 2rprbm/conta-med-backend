package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/2rprbm/conta-med-backend/config"
	"github.com/2rprbm/conta-med-backend/internal/adapters/primary/api"
	httpserver "github.com/2rprbm/conta-med-backend/internal/adapters/primary/http"
	"github.com/2rprbm/conta-med-backend/internal/adapters/secondary/repository/mongodb"
	"github.com/2rprbm/conta-med-backend/internal/adapters/secondary/whatsapp"
	"github.com/2rprbm/conta-med-backend/internal/core/services"
	mongoclient "github.com/2rprbm/conta-med-backend/pkg/mongodb"
	"github.com/2rprbm/conta-med-backend/pkg/logger"
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

	// Initialize MongoDB client
	mongoClient, err := mongoclient.NewClient(&cfg.MongoDB, log)
	if err != nil {
		log.Fatal("Failed to initialize MongoDB client", logger.Fields{
			"error": err.Error(),
		})
	}

	// Create indexes for better performance
	ctx := context.Background()
	if err := mongoClient.CreateIndexes(ctx); err != nil {
		log.Error("Failed to create MongoDB indexes", logger.Fields{
			"error": err.Error(),
		})
		// Continue execution as this is not critical for basic functionality
	}

	// Initialize repositories with real MongoDB implementation
	messageRepo := mongodb.NewMessageRepository(mongoClient.GetDatabase(), log)
	conversationRepo := mongodb.NewConversationRepository(mongoClient.GetDatabase(), log)

	// Initialize WhatsApp service
	whatsappConfig := whatsapp.WhatsAppConfig{
		APIVersion:    cfg.WhatsApp.APIVersion,
		PhoneNumberID: cfg.WhatsApp.PhoneNumberID,
		AccessToken:   cfg.WhatsApp.AccessToken,
		WebhookToken:  cfg.WhatsApp.WebhookVerifyToken,
		BaseURL:       cfg.WhatsApp.BaseURL,
	}
	whatsAppService := whatsapp.NewWhatsAppService(whatsappConfig, log)

	// Initialize chatbot service
	chatbotService := services.NewChatbotService(messageRepo, conversationRepo, whatsAppService, log)

	// Initialize webhook handler
	webhookHandler := api.NewWebhookHandler(chatbotService, whatsAppService, log)

	// Initialize HTTP server
	server := httpserver.NewServer(cfg, log, webhookHandler, mongoClient)

	// Start server
	server.Start()

	// Wait for termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Create a context with timeout for graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownPeriod)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("Server shutdown error", logger.Fields{"error": err.Error()})
	}

	// Close MongoDB connection
	if err := mongoClient.Close(shutdownCtx); err != nil {
		log.Error("MongoDB close error", logger.Fields{"error": err.Error()})
	}

	log.Info("Server stopped")
}
