package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockLogger implements the logger.Logger interface for testing
type mockLogger struct {
	buffer bytes.Buffer
	logger *log.Logger
}

func newMockLogger() *mockLogger {
	ml := &mockLogger{}
	ml.logger = log.New(&ml.buffer, "", 0)
	return ml
}

func (m *mockLogger) Debug(format string, args ...interface{}) {
	m.logger.Printf(format, args...)
}

func (m *mockLogger) Info(format string, args ...interface{}) {
	m.logger.Printf(format, args...)
}

func (m *mockLogger) Warn(format string, args ...interface{}) {
	m.logger.Printf(format, args...)
}

func (m *mockLogger) Error(format string, args ...interface{}) {
	m.logger.Printf(format, args...)
}

func (m *mockLogger) Fatal(format string, args ...interface{}) {
	m.logger.Printf(format, args...)
}

func (m *mockLogger) String() string {
	return m.buffer.String()
}

func TestLoggerMiddleware(t *testing.T) {
	t.Run("should log HTTP requests", func(t *testing.T) {
		// arrange
		mockLog := newMockLogger()
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		// Create the middleware handler
		middleware := Logger(mockLog)(nextHandler)

		// Create a test request
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("User-Agent", "test-agent")

		// Create a recorder to capture the response
		recorder := httptest.NewRecorder()

		// act
		middleware.ServeHTTP(recorder, req)

		// assert
		// Verify the response was passed through
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "OK", recorder.Body.String())

		// Check that the request was logged
		logOutput := mockLog.String()
		assert.Contains(t, logOutput, "GET")
		assert.Contains(t, logOutput, "/test")
		assert.Contains(t, logOutput, "200") // Status code
		assert.Contains(t, logOutput, "test-agent")
	})

	t.Run("should log different HTTP status codes", func(t *testing.T) {
		// arrange
		mockLog := newMockLogger()
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})

		// Create the middleware handler
		middleware := Logger(mockLog)(nextHandler)

		// Create a test request
		req := httptest.NewRequest("POST", "/missing", nil)

		// Create a recorder to capture the response
		recorder := httptest.NewRecorder()

		// act
		middleware.ServeHTTP(recorder, req)

		// assert
		// Verify the response was passed through
		assert.Equal(t, http.StatusNotFound, recorder.Code)

		// Check that the request was logged with the correct status code
		logOutput := mockLog.String()
		assert.Contains(t, logOutput, "POST")
		assert.Contains(t, logOutput, "/missing")
		assert.Contains(t, logOutput, "404") // Not found status code
	})
}
