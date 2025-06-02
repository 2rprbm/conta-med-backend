package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/2rprbm/conta-med-backend/config"
	httpserver "github.com/2rprbm/conta-med-backend/internal/adapters/primary/http"
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

	// Initialize HTTP server
	server := httpserver.NewServer(cfg, log)

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
