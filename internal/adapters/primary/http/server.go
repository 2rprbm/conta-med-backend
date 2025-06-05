package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/2rprbm/conta-med-backend/config"
	"github.com/2rprbm/conta-med-backend/internal/adapters/primary/api"
	"github.com/2rprbm/conta-med-backend/internal/adapters/primary/http/handlers"
	"github.com/2rprbm/conta-med-backend/internal/adapters/primary/http/middleware"
	"github.com/2rprbm/conta-med-backend/pkg/logger"
	mongoclient "github.com/2rprbm/conta-med-backend/pkg/mongodb"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Server represents the HTTP server
type Server struct {
	server         *http.Server
	router         *chi.Mux
	logger         logger.Logger
	config         *config.Config
	webhookHandler *api.WebhookHandler
	healthHandler  *handlers.HealthHandler
}

// NewServer creates a new HTTP server
func NewServer(cfg *config.Config, log logger.Logger, webhookHandler *api.WebhookHandler, mongoClient *mongoclient.Client) *Server {
	r := chi.NewRouter()

	srv := &Server{
		server: &http.Server{
			Addr:         ":" + cfg.Server.Port,
			Handler:      r,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
			IdleTimeout:  cfg.Server.IdleTimeout,
		},
		router:         r,
		logger:         log,
		config:         cfg,
		webhookHandler: webhookHandler,
		healthHandler:  handlers.NewHealthHandler(mongoClient, log),
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

	// CORS configuration
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
}

// setupRoutes sets up the routes for the server
func (s *Server) setupRoutes() {
	// Health check with detailed status
	s.router.Get("/health", s.healthHandler.CheckHealth)

	// Simple health check for load balancers
	s.router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	// Adicionando suporte para HEAD em /ping
	s.router.Head("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// Não precisa escrever o corpo em HEAD
	})

	// WhatsApp webhook routes (directly at root level as per WhatsApp convention)
	s.webhookHandler.RegisterRoutes(s.router)

	// Adicionando suporte para HEAD em /webhook (handler separado)
	s.router.Head("/webhook", func(w http.ResponseWriter, r *http.Request) {
		// Registrar todos os detalhes da requisição para diagnóstico
		s.logger.Info("DIAGNÓSTICO - Recebida requisição HEAD para /webhook", logger.Fields{
			"headers":      fmt.Sprintf("%v", r.Header),
			"remote_addr":  r.RemoteAddr,
			"url":          r.URL.String(),
			"query_params": fmt.Sprintf("%v", r.URL.Query()),
		})

		// Para webhooks, usamos a mesma lógica de verificação do GET
		mode := r.URL.Query().Get("hub.mode")
		token := r.URL.Query().Get("hub.verify_token")
		challenge := r.URL.Query().Get("hub.challenge")

		s.logger.Debug("Received HEAD webhook verification request", logger.Fields{
			"mode":      mode,
			"token":     token,
			"challenge": challenge,
		})

		// Adicionando headers CORS específicos para o webhook
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, POST, OPTIONS")
		w.Header().Set("Content-Type", "text/plain")

		// Verificação com os mesmos critérios do GET
		if mode == "subscribe" && token == s.config.WhatsApp.WebhookVerifyToken {
			s.logger.Info("HEAD webhook verified successfully", logger.Fields{
				"challenge": challenge,
				"token":     token,
			})
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(challenge))

			// Log após enviar a resposta
			s.logger.Info("DIAGNÓSTICO - Resposta enviada para HEAD webhook", logger.Fields{
				"status_code":  http.StatusOK,
				"challenge":    challenge,
				"content_type": w.Header().Get("Content-Type"),
			})
			return
		}

		// Verificação falhou
		s.logger.Warn("HEAD webhook verification failed: invalid token or mode", logger.Fields{
			"expected_token": s.config.WhatsApp.WebhookVerifyToken,
			"received_token": token,
			"mode":           mode,
		})
		http.Error(w, "Verification failed", http.StatusForbidden)
	})

	// All future endpoints will be registered directly on the root router
}

// Start starts the HTTP server
func (s *Server) Start() {
	s.logger.Info("Starting server", logger.Fields{
		"port": s.config.Server.Port,
	})
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("Server error", logger.Fields{
				"error": err.Error(),
			})
		}
	}()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server...")
	return s.server.Shutdown(ctx)
}
