package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	mongoclient "github.com/2rprbm/conta-med-backend/pkg/mongodb"
	"github.com/2rprbm/conta-med-backend/pkg/logger"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	mongoClient *mongoclient.Client
	logger      logger.Logger
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(mongoClient *mongoclient.Client, log logger.Logger) *HealthHandler {
	return &HealthHandler{
		mongoClient: mongoClient,
		logger:      log,
	}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
	Version   string            `json:"version"`
}

// CheckHealth handles GET requests to check application health
func (h *HealthHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	h.logger.Debug("Health check requested")

	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Services:  make(map[string]string),
		Version:   "v0.2.0",
	}

	// Check MongoDB health
	if h.mongoClient != nil {
		if err := h.mongoClient.Health(ctx); err != nil {
			h.logger.Error("MongoDB health check failed", logger.Fields{
				"error": err.Error(),
			})
			response.Services["mongodb"] = "unhealthy"
			response.Status = "degraded"
		} else {
			response.Services["mongodb"] = "healthy"
		}
	} else {
		response.Services["mongodb"] = "not_configured"
		response.Status = "degraded"
	}

	// Set appropriate status code
	statusCode := http.StatusOK
	if response.Status == "degraded" {
		statusCode = http.StatusServiceUnavailable
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("Failed to encode health response", logger.Fields{
			"error": err.Error(),
		})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.logger.Debug("Health check completed", logger.Fields{
		"status":     response.Status,
		"mongodb":    response.Services["mongodb"],
		"status_code": statusCode,
	})
} 