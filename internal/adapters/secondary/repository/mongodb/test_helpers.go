package mongodb

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/2rprbm/conta-med-backend/pkg/logger"
)

type mockLogger struct{}

func (m *mockLogger) Info(msg string, fields ...logger.Fields)  {}
func (m *mockLogger) Error(msg string, fields ...logger.Fields) {}
func (m *mockLogger) Debug(msg string, fields ...logger.Fields) {}
func (m *mockLogger) Warn(msg string, fields ...logger.Fields)  {}
func (m *mockLogger) Fatal(msg string, fields ...logger.Fields) {}
func (m *mockLogger) With(fields logger.Fields) logger.Logger   { return m }

func setupTestDB(t *testing.T) (*mongo.Database, func()) {
	// Using in-memory MongoDB for tests
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)

	dbName := "test_db_" + primitive.NewObjectID().Hex()
	db := client.Database(dbName)

	cleanup := func() {
		_ = db.Drop(context.Background())
		_ = client.Disconnect(context.Background())
	}

	return db, cleanup
} 