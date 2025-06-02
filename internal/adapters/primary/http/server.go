package http

import (
	"context"
	"net/http"
	"time"

	"github.com/2rprbm/conta-med-backend/config"
	"github.com/2rprbm/conta-med-backend/internal/adapters/primary/api"
	"github.com/2rprbm/conta-med-backend/internal/adapters/primary/http/middleware"
	"github.com/2rprbm/conta-med-backend/pkg/logger"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Server represents the HTTP server
type Server struct {
	server        *http.Server
	router        *chi.Mux
	logger        logger.Logger
	config        *config.Config
	webhookHandler *api.WebhookHandler
}

// NewServer creates a new HTTP server
func NewServer(cfg *config.Config, log logger.Logger, webhookHandler *api.WebhookHandler) *Server {
	r := chi.NewRouter()

	srv := &Server{
		server: &http.Server{
			Addr:         ":" + cfg.Server.Port,
			Handler:      r,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
			IdleTimeout:  cfg.Server.IdleTimeout,
		},
		router:        r,
		logger:        log,
		config:        cfg,
		webhookHandler: webhookHandler,
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
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
}

// setupRoutes sets up the routes for the server
func (s *Server) setupRoutes() {
	// Health check
	s.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API routes
	s.router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			// Register webhook routes
			s.webhookHandler.RegisterRoutes(r)
		})
	})
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
