package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("ContaMed WhatsApp Chatbot - Starting server...")

	// TODO: Load configuration from environment variables
	port := "8080"
	if envPort := os.Getenv("SERVER_PORT"); envPort != "" {
		port = envPort
	}

	// TODO: Initialize MongoDB connection
	// TODO: Initialize WhatsApp client
	// TODO: Setup router and routes

	// Simple health check endpoint for now
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start HTTP server
	srv := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		fmt.Printf("Server listening on port %s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until signal is received
	<-quit
	fmt.Println("Shutting down server...")

	// TODO: Close MongoDB connection
	// TODO: Perform any necessary cleanup

	fmt.Println("Server stopped")
}
