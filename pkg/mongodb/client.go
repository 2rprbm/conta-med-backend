package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/2rprbm/conta-med-backend/config"
	"github.com/2rprbm/conta-med-backend/pkg/logger"
)

// Client represents MongoDB client wrapper
type Client struct {
	client   *mongo.Client
	database *mongo.Database
	logger   logger.Logger
	config   *config.MongoDBConfig
}

// NewClient creates a new MongoDB client
func NewClient(cfg *config.MongoDBConfig, log logger.Logger) (*Client, error) {
	if cfg.URI == "" {
		return nil, fmt.Errorf("MongoDB URI is required")
	}

	log.Info("Connecting to MongoDB", logger.Fields{
		"database": cfg.Database,
		"timeout":  cfg.Timeout.String(),
	})

	// Set client options
	clientOptions := options.Client().ApplyURI(cfg.URI)
	clientOptions.SetConnectTimeout(cfg.Timeout)
	clientOptions.SetServerSelectionTimeout(cfg.Timeout)
	
	// Create a context with timeout for connection
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Error("Failed to connect to MongoDB", logger.Fields{
			"error": err.Error(),
			"uri":   cfg.URI,
		})
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Error("Failed to ping MongoDB", logger.Fields{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	database := client.Database(cfg.Database)

	log.Info("Successfully connected to MongoDB", logger.Fields{
		"database": cfg.Database,
	})

	return &Client{
		client:   client,
		database: database,
		logger:   log,
		config:   cfg,
	}, nil
}

// GetDatabase returns the MongoDB database instance
func (c *Client) GetDatabase() *mongo.Database {
	return c.database
}

// GetClient returns the MongoDB client instance
func (c *Client) GetClient() *mongo.Client {
	return c.client
}

// Close gracefully closes the MongoDB connection
func (c *Client) Close(ctx context.Context) error {
	c.logger.Info("Closing MongoDB connection")
	
	if err := c.client.Disconnect(ctx); err != nil {
		c.logger.Error("Failed to disconnect from MongoDB", logger.Fields{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}

	c.logger.Info("MongoDB connection closed successfully")
	return nil
}

// CreateIndexes creates necessary indexes for the application
func (c *Client) CreateIndexes(ctx context.Context) error {
	c.logger.Info("Creating MongoDB indexes")

	// Conversations indexes
	conversationsCollection := c.database.Collection("conversations")
	
	// Index for phone_number (unique)
	_, err := conversationsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{"phone_number", 1}},
		Options: options.Index().SetUnique(false), // Allow multiple conversations per phone
	})
	if err != nil {
		return fmt.Errorf("failed to create phone_number index: %w", err)
	}

	// Index for phone_number + status for active conversation lookups
	_, err = conversationsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{"phone_number", 1},
			{"status", 1},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create phone_number_status index: %w", err)
	}

	// Index for last_updated_at for sorting
	_, err = conversationsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{"last_updated_at", -1}},
	})
	if err != nil {
		return fmt.Errorf("failed to create last_updated_at index: %w", err)
	}

	// Messages indexes
	messagesCollection := c.database.Collection("messages")
	
	// Index for conversation_id
	_, err = messagesCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{"conversation_id", 1}},
	})
	if err != nil {
		return fmt.Errorf("failed to create conversation_id index: %w", err)
	}

	// Index for phone_number
	_, err = messagesCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{"phone_number", 1}},
	})
	if err != nil {
		return fmt.Errorf("failed to create phone_number index on messages: %w", err)
	}

	// Index for timestamp for sorting messages
	_, err = messagesCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{"timestamp", -1}},
	})
	if err != nil {
		return fmt.Errorf("failed to create timestamp index: %w", err)
	}

	// Compound index for conversation_id + timestamp for efficient message queries
	_, err = messagesCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{"conversation_id", 1},
			{"timestamp", -1},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create conversation_id_timestamp index: %w", err)
	}

	c.logger.Info("MongoDB indexes created successfully")
	return nil
}

// Health checks the health of the MongoDB connection
func (c *Client) Health(ctx context.Context) error {
	// Create a context with timeout
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := c.client.Ping(pingCtx, nil); err != nil {
		c.logger.Error("MongoDB health check failed", logger.Fields{
			"error": err.Error(),
		})
		return fmt.Errorf("MongoDB health check failed: %w", err)
	}

	return nil
} 