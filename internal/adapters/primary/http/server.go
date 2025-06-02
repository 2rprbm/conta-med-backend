package http

import (
	"context"
	"net/http"
	"time"

	"github.com/2rprbm/conta-med-backend/config"
	"github.com/2rprbm/conta-med-backend/internal/adapters/primary/http/handlers"
	"github.com/2rprbm/conta-med-backend/internal/adapters/primary/http/middleware"
	"github.com/2rprbm/conta-med-backend/pkg/logger"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// Server represents the HTTP server
type Server struct {
	server *http.Server
	router *chi.Mux
	logger logger.Logger
	config *config.Config
}

// NewServer creates a new HTTP server
func NewServer(cfg *config.Config, log logger.Logger) *Server {
	r := chi.NewRouter()

	srv := &Server{
		server: &http.Server{
			Addr:         ":" + cfg.Server.Port,
			Handler:      r,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
			IdleTimeout:  cfg.Server.IdleTimeout,
		},
		router: r,
		logger: log,
		config: cfg,
	}

	srv.setupMiddleware()
	srv.setupRoutes()

	return srv
}

// setupMiddleware sets up the middleware for the server
func (s *Server) setupMiddleware() {
	// Basic middlewares
	s.router.Use(chimiddleware.RequestID)
	s.router.Use(chimiddleware.RealIP)
	s.router.Use(middleware.Logger(s.logger))
	s.router.Use(chimiddleware.Recoverer)

	// Set a timeout value on the request context
	s.router.Use(chimiddleware.Timeout(60 * time.Second))
}

// setupRoutes sets up the routes for the server
func (s *Server) setupRoutes() {
	// Health check
	s.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Webhook handler
	webhookHandler := handlers.NewWebhookHandler(s.config, s.logger)

	// WhatsApp webhook routes
	s.router.Route("/webhook", func(r chi.Router) {
		// GET for webhook verification
		r.Get("/whatsapp", webhookHandler.VerifyToken)

		// POST for receiving messages
		r.Post("/whatsapp", webhookHandler.ReceiveWebhook)
	})
}

// Start starts the HTTP server
func (s *Server) Start() {
	s.logger.Info("Starting server on port %s", s.config.Server.Port)
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("Server error: %v", err)
		}
	}()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server...")
	return s.server.Shutdown(ctx)
}
